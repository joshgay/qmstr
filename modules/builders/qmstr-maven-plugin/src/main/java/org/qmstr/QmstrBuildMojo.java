package org.qmstr;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Collections;
import java.util.List;
import java.util.Set;
import java.util.stream.Collectors;

import org.apache.maven.artifact.handler.ArtifactHandler;
import org.apache.maven.execution.MavenSession;
import org.apache.maven.plugin.AbstractMojo;
import org.apache.maven.plugin.MojoExecutionException;
import org.apache.maven.plugins.annotations.LifecyclePhase;
import org.apache.maven.plugins.annotations.Mojo;
import org.apache.maven.plugins.annotations.Parameter;
import org.apache.maven.project.MavenProject;
import org.qmstr.client.BuildServiceClient;
import org.qmstr.grpc.service.Datamodel;
import org.qmstr.util.FilenodeUtils;
import org.qmstr.util.transformations.*;

/**
 * add built files to qmstr master server
 */
@Mojo(name = "qmstrbuild", defaultPhase = LifecyclePhase.PROCESS_CLASSES)
public class QmstrBuildMojo extends AbstractMojo {
    private final String qmstrMasterAddress = System.getenv("QMSTR_MASTER");

    /**
     * address to connect to
     */
    @Parameter(defaultValue = "localhost", property = "qmstrAddress", required = true)
    private String qmstrAddress;

    @Parameter(defaultValue = "${session}", readonly = true)
    private MavenSession session;

    @Parameter(defaultValue = "${project}", readonly = true)
    private MavenProject project;

    /**
     * The directory for compiled classes.
     *
     */
    @Parameter(readonly = true, required = true, defaultValue = "${project.build.outputDirectory}")
    private File outputDirectory;

    public void execute() throws MojoExecutionException {
        BuildServiceClient bsc = new BuildServiceClient(qmstrAddress);

        ArtifactHandler artifactHandler = project.getArtifact().getArtifactHandler();
        if (!artifactHandler.getLanguage().equalsIgnoreCase("java")) {
            getLog().warn("Not a java project");
            return;
        }

        Set<File> sourceDirs = project.getCompileSourceRoots().stream().map(sds -> Paths.get(sds).toFile())
                .collect(Collectors.toSet());

        Set<File> sources = getSourceFiles(project.getCompileSourceRoots());

        try {
            Set<Datamodel.FileNode> fileNodes = FilenodeUtils.processSourceFiles(Transform.COMPILEJAVA, sources,
                    sourceDirs, Collections.singleton(outputDirectory));
            bsc.SendBuildFileNodes(fileNodes);
        } catch (TransformationException e) {
            // ("qmstr plugin could not transform source to target " + e.getMessage());
        } catch (FileNotFoundException fnfe) {
            //throw new MojoExecutionException("qmstr plugin could not find the source file " + fnfe.getMessage());
        }

        try {
            bsc.close();
        } catch (InterruptedException e) {
            throw new MojoExecutionException("qmstr: failed to close grpc channel " + e.getMessage());
        }

    }

    private Set<File> getSourceFiles(List<String> sourceDirs) {
        return sourceDirs.stream().map(sds -> Paths.get(sds)).flatMap(sd -> getSourcesFromDir(sd).stream())
                .collect(Collectors.toSet());
    }

    private Set<File> getSourcesFromDir(Path sourceDir) {
        try {
            return Files.walk(sourceDir).filter(Files::isRegularFile)
                    .filter(p -> p.getFileName().toString().endsWith(".java")).map(p -> p.toFile())
                    .collect(Collectors.toSet());
        } catch (IOException ioe) {
            return Collections.emptySet();
        }
    }
}
