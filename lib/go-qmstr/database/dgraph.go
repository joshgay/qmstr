package database

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"text/template"
	"time"

	"encoding/json"

	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
	client "github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"

	"google.golang.org/grpc"
)

const schema = `
path:string @index(hash,trigram) .
hash:string @index(exact) .
type:string @index(hash) .
name:string @index(hash) .
fileData:uid .
dataNodes:[uid] @reverse .
data:string @index(hash) .
projectNodeType:string @index(hash) .
trustLevel:string @index(hash) .
packageNodeType:string @index(hash) .
fileNodeType:string @index(hash) .
fileDataNodeType:string @index(hash) .
infoNodeType:string @index(hash) .
diagnosticNodeType:string @index(hash) .
severity:string @index(hash) .
dataNodeType:string @index(hash) .
analyzerNodeType:string @index(hash) .
`

type insertQueueState int8

const (
	openQueue insertQueueState = iota
	closeQueue
)

type DataBase struct {
	client       *client.Dgraph
	insertQueue  chan interface{}
	insertMutex  *sync.Mutex
	pending      uint64
	queueWorkers uint16
}

func CheckSchema(checkSchema string) bool {
	scanner := bufio.NewScanner(strings.NewReader(schema))
	for scanner.Scan() {
		if !strings.Contains(checkSchema, scanner.Text()) {
			log.Printf("Required schema not found: %s", scanner.Text())
			return false
		}
	}
	return true
}

// Setup connects to dgraph and returns the instance
func Setup(dbAddr string, queueWorkers int) (*DataBase, error) {
	log.Println("Setting up database connection")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	conn, err := grpc.DialContext(ctx, dbAddr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(100<<20),
			grpc.MaxCallRecvMsgSize(100<<20)))
	if err != nil {
		cancel()
		return nil, fmt.Errorf("Failed to connect to the dgraph server: %v", err)
	}

	db := &DataBase{
		client:       client.NewDgraphClient(api.NewDgraphClient(conn)),
		insertMutex:  &sync.Mutex{},
		queueWorkers: uint16(queueWorkers),
	}

	for {
		err = db.client.Alter(ctx, &api.Operation{
			Schema: schema,
		})
		if err == nil {
			cancel()
			break
		}
	}
	return db, nil
}

func (db *DataBase) GetPendingInserts() uint64 {
	return atomic.LoadUint64(&db.pending)
}

func (db *DataBase) Sync() {
	// TODO replace busy waiting with proper signaling
	log.Println("Waiting for inserts")
	for {
		pendingInserts := atomic.LoadUint64(&db.pending)
		if pendingInserts == 0 {
			break
		}
		log.Printf("Pending inserts %d", pendingInserts)
		time.Sleep(2 * time.Second)
	}
}

func (db *DataBase) CloseInsertQueue() {
	db.Sync()
	db.insertMutex.Lock()
	defer db.insertMutex.Unlock()
	if db.insertQueue == nil {
		return
	}
	close(db.insertQueue)
	db.insertQueue = nil
	log.Println("Closed insert queue")
}

func (db *DataBase) OpenInsertQueue() {
	db.insertMutex.Lock()
	defer db.insertMutex.Unlock()
	if db.insertQueue == nil {
		db.insertQueue = make(chan interface{}, 1000)
		log.Println("Opened insert queue")

		for i := uint16(0); i < db.queueWorkers; i++ {
			go db.queueWorker()
		}
	}
}

// the queueWorker runs in a go routine and inserts the nodes from the insert queue into the database
func (db *DataBase) queueWorker() {
	for {
		node := <-db.insertQueue
		if node == nil {
			return
		}

		nodeType := reflect.TypeOf(node)
		switch nodeType {
		case reflect.TypeOf((*service.ProjectNode)(nil)):
			db.insertProjectNode(node.(*service.ProjectNode))
		case reflect.TypeOf((*service.PackageNode)(nil)):
			db.insertPkgNode(node.(*service.PackageNode))
		case reflect.TypeOf((*service.FileNode)(nil)):
			db.insertFileNode(node.(*service.FileNode))
		default:
			log.Printf("Wrong node type %s trying to be inserted in the database", nodeType.String())
		}
	}
}

func (db *DataBase) checkFileData(node *service.FileNode) error {
	hash := node.FileData.GetHash()
	if hash == "" {
		return fmt.Errorf("hash cannot be empty for node: %v", node)
	}
	fDataUID, err := db.GetFileDataUID(hash)
	if err != nil {
		return err
	}
	if fDataUID != "" {
		node.FileData.Uid = fDataUID
	}
	return nil
}

func (db *DataBase) insertFileNode(node *service.FileNode) {
	ready := true

	for idx, dep := range node.DerivedFrom {
		// check if fileData already exist in db
		err := db.checkFileData(dep)
		if err != nil {
			log.Fatalf("checkFileData failed: %v", err)
		}
		if dep.FileData.Uid == "" {
			ready = false
		}
		if dep.Uid == "" {
			// missing dep
			ready = false
			// look up dep in db
			uid, err := db.GetFileNodeUid(dep.Path, dep.FileData.GetHash())
			if err != nil {
				log.Fatalf("getFileNodeUid failed for node: %v: %v", node, err)
			}
			// found uid
			if uid != "" {
				node.DerivedFrom[idx].Uid = uid
			}
		}
	}

	for idx, dep := range node.Dependencies {
		// check if fileData already exist in db
		err := db.checkFileData(dep)
		if err != nil {
			log.Fatalf("checkFileData failed: %v", err)
		}
		if dep.FileData.Uid == "" {
			ready = false
		}

		if dep.Uid == "" {
			// missing dep
			ready = false
			// look up dep in db
			uid, err := db.GetFileNodeUid(dep.Path, dep.FileData.GetHash())
			if err != nil {
				log.Fatalf("getFileNodeUid failed for node: %v: %v", node, err)
			}
			// found uid
			if uid != "" {
				node.Dependencies[idx].Uid = uid
			}
		}
	}

	if !ready {
		// put node back to queue
		go func() { db.insertQueue <- node }()
		return
	}

	// we are ready to insert the node
	db.insertMutex.Lock()

	// check if fileData already exist in db
	err := db.checkFileData(node)
	if err != nil {
		log.Fatalf("checkFileData failed: %v", err)
	}

	uid, err := db.GetFileNodeUid(node.Path, node.FileData.GetHash())
	if err != nil {
		log.Fatalf("getFileNodeUid failed for node: %v: %v", node, err)
	}
	if uid != "" {
		node.Uid = uid
	}
	_, err = dbInsert(db.client, node)
	if err != nil {
		log.Fatalf("failed inserting file node %v: %v", node, err)
	}
	atomic.AddUint64(&db.pending, ^uint64(0))
	db.insertMutex.Unlock()
}

func (db *DataBase) insertPkgNode(node *service.PackageNode) {
	ready := true
	for idx, dep := range node.Targets {
		// check if fileData already exist in db
		err := db.checkFileData(dep)
		if err != nil {
			log.Fatalf("checkFileData failed: %v", err)
		}
		if dep.FileData.Uid == "" {
			ready = false
		}
		if dep.Uid == "" {
			// missing dep
			ready = false
			// look up dep in db
			uid, err := db.GetFileNodeUid(dep.Path, dep.FileData.GetHash())
			if err != nil {
				log.Fatalf("getFileNodeUid failed for node: %v: %v", node, err)
			}
			// found uid
			if uid != "" {
				node.Targets[idx].Uid = uid
			}
		}
	}

	if !ready {
		// put node back to queue
		go func() { db.insertQueue <- node }()
		return
	}

	// we are ready to insert the node
	db.insertMutex.Lock()
	packageNode, err := db.GetPackageNodeByName(node.Name)
	if err == nil {
		node.Uid = packageNode.Uid
		node.Targets = append(packageNode.Targets, node.Targets...)
	}

	_, err = dbInsert(db.client, node)
	if err != nil {
		log.Fatalf("failed inserting package node failed: %v", err)
	}
	atomic.AddUint64(&db.pending, ^uint64(0))
	db.insertMutex.Unlock()
}

func (db *DataBase) insertProjectNode(node *service.ProjectNode) {
	db.insertMutex.Lock()

	_, err := dbInsert(db.client, node)
	if err != nil {
		log.Fatalf("failed inserting package node failed: %v", err)
	}
	atomic.AddUint64(&db.pending, ^uint64(0))
	db.insertMutex.Unlock()
}

func fixNodeFields(field *reflect.Value) error {
	switch field.Kind() {
	case reflect.Struct:
		for i := 0; i < field.NumField(); i++ {
			fieldVal := field.Field(i)
			if fieldVal.Kind() == reflect.String && strings.HasSuffix(field.Type().Field(i).Name, "NodeType") {
				if fieldVal.CanSet() {
					fieldVal.SetString("_")
					continue
				}
				return fmt.Errorf("%s not settable", field.Type().Field(i).Name)
			} else if fieldVal.Kind() == reflect.String && field.Type().Field(i).Name == "Timestamp" {
				if fieldVal.CanSet() {
					timeInBytes, err := time.Now().MarshalText()
					if err != nil {
						return err
					}
					time := string(timeInBytes)
					fieldVal.SetString(time)
					continue
				}
				return fmt.Errorf("%s not settable", field.Type().Field(i).Name)
			}
			fixNodeFields(&fieldVal)
		}
	case reflect.Ptr:
		fieldVal := field.Elem()
		fixNodeFields(&fieldVal)
	case reflect.Slice, reflect.Array:
		for i := 0; i < field.Len(); i++ {
			fieldVal := field.Index(i).Elem()
			fixNodeFields(&fieldVal)
		}
	}
	return nil
}

// fill in the type and timestamp value to work around omitting empty values on serialization
func fillNodeFields(data interface{}) error {
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Ptr {
		return errors.New("you shall not call fillNodeFields by value")
	}
	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		val = val.Elem()
	}
	return fixNodeFields(&val)
}

// the data should be JSON marshalable
func dbInsert(c *client.Dgraph, data interface{}) (map[string]string, error) {
	txn := c.NewTxn()
	defer txn.Discard(context.Background())

	if err := fillNodeFields(&data); err != nil {
		return nil, err
	}

	jNode, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	ret, err := txn.Mutate(context.Background(), &api.Mutation{CommitNow: true, SetJson: jNode})
	if err != nil {
		return nil, err
	}

	return ret.Uids, nil
}

// Delete deletes nodes from the db.
// The data should be JSON marshalable
func Delete(db *DataBase, nodeUID string) (string, error) {
	txn := db.client.NewTxn()
	defer txn.Discard(context.Background())

	// we just need the uid of the node
	d := map[string]string{"uid": nodeUID}
	pb, err := json.Marshal(d)
	if err != nil {
		log.Fatal(err)
	}
	ret, err := txn.Mutate(context.Background(), &api.Mutation{CommitNow: true, DeleteJson: pb})
	if err != nil {
		return "", err
	}
	uid := ret.Uids["blank-0"]
	return uid, nil
}

// DeleteEdge deletes an edge from a node
func DeleteEdge(db *DataBase, nodeUID string, predicate string) (string, error) {
	txn := db.client.NewTxn()
	defer txn.Discard(context.Background())

	mu := &api.Mutation{CommitNow: true}
	client.DeleteEdges(mu, nodeUID, predicate)
	ret, err := txn.Mutate(context.Background(), mu)
	if err != nil {
		return "", err
	}
	uid := ret.Uids["blank-0"]
	return uid, nil
}

func getVarName(index int) string {
	var result string
	for index > 25 {
		result = result + string(rune(65))
		index = index - 26
	}
	result = result + string(rune(65+index))
	return result
}

func (db *DataBase) GetNodesByType(valuetype string, recursive bool, namefilter string) ([]*service.FileNode, error) {

	var ret map[string]interface{}

	q := `query NodeByType($Type: string, $Name: string){
		  getNodeByType(func: eq(type, $Type)) {{.Filter}} {{.Recurse}}{
			uid
			hash
			path
			derivedFrom
		  }}`

	queryTmpl, err := template.New("nodesbytype").Parse(q)

	type QueryParams struct {
		Recurse string
		Filter  string
	}

	qp := QueryParams{}
	if recursive {
		qp.Recurse = "@recurse(loop: false)"
	}
	if namefilter != "" {
		qp.Filter = "@filter(eq(name, $Name))"
	}

	var b bytes.Buffer
	err = queryTmpl.Execute(&b, qp)
	if err != nil {
		return nil, fmt.Errorf("GetNodesByType failed: %v", err)
	}

	vars := map[string]string{"$Type": valuetype, "$Name": namefilter}

	err = db.queryNodes(b.String(), vars, &ret)
	if err != nil {
		return nil, err
	}

	return ret["getNodeByType"].([]*service.FileNode), nil
}

func (db *DataBase) queryNodes(query string, queryVars map[string]string, resultMap interface{}) error {
	resp, err := db.client.NewTxn().QueryWithVars(context.Background(), query, queryVars)
	if err != nil {
		return fmt.Errorf("Could not query for node with: \n\n%s\n\nVars:\n\n%v\n\nError: %v", query, queryVars, err)
	}

	if err = json.Unmarshal(resp.Json, resultMap); err != nil {
		return fmt.Errorf("Could not unmashal query response: %v", err)
	}
	return nil
}

func (db *DataBase) queryNodesSimple(query string, resultMap interface{}) error {
	resp, err := db.client.NewTxn().Query(context.Background(), query)
	if err != nil {
		return fmt.Errorf("Could not query for node with: \n\n%s\n\nError: %v", query, err)
	}

	if err = json.Unmarshal(resp.Json, resultMap); err != nil {
		return fmt.Errorf("Could not unmashal query response: %v", err)
	}
	return nil
}
