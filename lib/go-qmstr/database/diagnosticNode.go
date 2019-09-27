package database

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"text/template"

	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
)

// AddDiagnosticNodes stores the given DiagnosticNodes in a PackageNode or FileNode identified by the nodeID
func (db *DataBase) AddDiagnosticNode(nodeID string, diagnosticnode *service.DiagnosticNode) error {
	db.insertMutex.Lock()
	defer db.insertMutex.Unlock()

	const q = `
	query Node($id: string){
		node(func: uid($id)) @filter(has(packageNodeType) or has(fileDataNodeType)) @recurse(loop: false) {
			uid
			packageNodeType
			fileDataNodeType
			diagnosticInfo
			analyzer
			name
		}
	}
	`
	vars := map[string]string{"$id": nodeID}
	var result map[string][]interface{}

	resp, err := db.client.NewTxn().QueryWithVars(context.Background(), q, vars)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resp.Json, &result)
	if err != nil {
		log.Fatal(err)
	}

	if len(result["node"]) < 1 {
		return fmt.Errorf("No package or file node with uid %s found", nodeID)
	}

	receiverNode := result["node"][0].(map[string]interface{})
	if diagnosticInfoInter, ok := receiverNode["diagnosticInfo"]; ok {
		// each analyzer should create one diagnostic node for each file node
		// so check if diagnostic node already exists
		for _, diagnosticInfo := range diagnosticInfoInter.([]interface{}) {
			for attrName, attrValue := range diagnosticInfo.(map[string]interface{}) {
				if attrName == "analyzer" {
					for _, analyzer := range attrValue.([]interface{}) {
						for name, value := range analyzer.(map[string]interface{}) {
							if name == "name" {
								if value.(string) == diagnosticnode.Analyzer[0].Name {
									log.Printf("Already created diagnostic node for file %s, skipping insert..", nodeID)
									return nil
								}
							}
						}
					}
				}
			}
		}
	}
	var diagnosticInfo []*service.DiagnosticNode
	diagnosticInfo = append(diagnosticInfo, diagnosticnode)

	if _, ok := receiverNode["packageNodeType"]; ok {
		packageNode := service.PackageNode{}
		packageNode.Uid = nodeID
		packageNode.DiagnosticInfo = []*service.DiagnosticNode{diagnosticnode}
		_, err = dbInsert(db.client, &packageNode)
		if err != nil {
			return err
		}
	} else if _, ok := receiverNode["fileDataNodeType"]; ok {
		fileDataNode := service.FileNode_FileDataNode{}
		fileDataNode.Uid = nodeID
		fileDataNode.DiagnosticInfo = []*service.DiagnosticNode{diagnosticnode}
		_, err = dbInsert(db.client, &fileDataNode)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("wrong type of node. Can't connect diagnostic nodes to it")
	}
	return nil
}

//GetDiagnosticNodeBySeverity queries diagnostic nodes on a specific severity
func (db *DataBase) GetDiagnosticNodeBySeverity(diagnosticNode *service.DiagnosticNode) ([]*service.DiagnosticNode, error) {
	var ret map[string][]*service.DiagnosticNode

	const q = `query DiagnosticData($Severity: int){
		getDiagnosticData(func: has(diagnosticNodeType)) @filter(eq(severity, $Severity)) {
			diagnosticInfo
			message
		}}`

	queryTmpl, err := template.New("diagnosticnodebyseverity").Parse(q)

	type QueryParams struct {
		Severity int
	}

	qp := QueryParams{}
	//get the int value from the enumeration
	t := service.DiagnosticNode_Severity_value[diagnosticNode.Severity.String()]
	nt := int(t)
	qp.Severity = nt

	//convert it to string to query it
	vars := map[string]string{"$Severity": strconv.Itoa(nt)}

	var b bytes.Buffer
	err = queryTmpl.Execute(&b, qp)
	if err != nil {
		return nil, err
	}
	err = db.queryNodes(b.String(), vars, &ret)
	if err != nil {
		return nil, err
	}

	messages := ret["getDiagnosticData"]
	if len(messages) < 1 {
		return nil, fmt.Errorf("No diagnostic node %v found in the database", strconv.Itoa(nt))
	}
	return messages, nil
}
