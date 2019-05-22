package main

import (
	"context"
	"crypto/sha1"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/QMSTR/qmstr/lib/go-qmstr/common"
	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
	"github.com/spdx/tools-golang/v0/spdx"
	"github.com/spdx/tools-golang/v0/tvsaver"
)

const (
	ModuleName  = "reporter-spdx"
	outFileName = "%s.spdx"
)

type SPDXReporter struct {
	enableWarnings bool
	enableErrors   bool
	outputdir      string
	nsURI          string
}

func (r *SPDXReporter) Configure(config map[string]string) error {
	if outDir, ok := config["outputdir"]; ok {
		r.outputdir = outDir
	} else {
		return fmt.Errorf("no output directory configured")
	}

	if nsURI, ok := config["namespaceURI"]; ok {
		r.nsURI = nsURI
	} else {
		// see: https://spdx.org/spdx-specification-21-web-version#h.1gdfkutofa90
		r.nsURI = "http://spdx.org/spdxdocs/"
	}

	return nil
}

func (r *SPDXReporter) Report(cserv service.ControlServiceClient, rserv service.ReportServiceClient) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pkgStream, err := cserv.GetPackageNode(ctx, &service.PackageNode{})
	if err != nil {
		return err
	}

	for {
		pkg, err := pkgStream.Recv()
		if err != nil {
			return err
		}
		if err = r.generateSPDX(pkg); err != nil {
			return err
		}

	}
}

func (r *SPDXReporter) generateSPDX(pkgNode *service.PackageNode) error {
	files := []*spdx.File2_1{}
	hashes := []string{}

	for _, trgt := range pkgNode.Targets {
		fl := &spdx.File2_1{
			FileName: trgt.Name,
			// this should be unique
			FileSPDXIdentifier: "SPDXRef-file-" + trgt.Name,
			FileChecksumSHA1:   trgt.Hash,
			LicenseConcluded:   "NOASSERTION",
			LicenseInfoInFile:  []string{"NOASSERTION"},
			FileCopyrightText:  "NOASSERTION",
		}
		files = append(files, fl)
		hashes = append(hashes, trgt.Hash)
	}

	dnldLocation := pkgNode.GetMetaData("SourceURL", "NOASSERTION")
	pkg := &spdx.Package2_1{
		PackageName: pkgNode.Name,
		// this should be unique
		PackageSPDXIdentifier:   "SPDXRef-pkg-" + pkgNode.Name,
		PackageDownloadLocation: dnldLocation,
		FilesAnalyzed:           true,
		PackageVerificationCode: calcSHA1Hash(hashes),
		// this is the license we detect
		PackageLicenseConcluded:     "NOASSERTION",
		PackageLicenseInfoFromFiles: []string{"NOASSERTION"},
		PackageLicenseDeclared:      pkgNode.GetMetaData("LicenseDeclared", "NOASSERTION"),
		PackageCopyrightText:        "NOASSERTION",
		Files:                       files,
	}

	doc := &spdx.Document2_1{
		CreationInfo: &spdx.CreationInfo2_1{
			SPDXVersion:       "SPDX-2.1",
			DataLicense:       "CC0-1.0",
			SPDXIdentifier:    "SPDXRef-DOCUMENT",
			DocumentName:      pkgNode.Name,
			DocumentNamespace: r.nsURI + url.PathEscape(pkgNode.Name),
			Created:           time.Now().Format(time.RFC3339),
			CreatorTools: []string{
				"QMSTR",
			},
		},
		Packages: []*spdx.Package2_1{pkg},
	}

	fName := filepath.Join(r.outputdir, fmt.Sprintf(outFileName, common.GetPosixFullyPortableFilename(pkgNode.Name)))
	out, err := os.Create(fName)
	if err != nil {
		return fmt.Errorf("failed to create out file %q: %v", fName, err)
	}
	defer out.Close()

	// export to tag-value format
	err = tvsaver.Save2_1(doc, out)
	if err != nil {
		return fmt.Errorf("failed to export SPDX doc: %v", err)
	}
	return nil
}

func (r *SPDXReporter) PostReport() error {
	return nil
}

func calcSHA1Hash(l []string) string {
	// sort the list
	sort.Strings(l)

	h := sha1.New()
	for _, s := range l {
		h.Write([]byte(s))
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}