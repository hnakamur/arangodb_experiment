package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/solher/arangolite"
)

type Node struct {
	arangolite.Document
}

func main() {
	var rootPassword string
	flag.StringVar(&rootPassword, "root-password", "", `the arrango password for "root" user`)
	var userPassword string
	flag.StringVar(&userPassword, "user-password", "", `the arrango password for "user" user`)
	flag.Parse()

	db := arangolite.New().
		LoggerOptions(false, false, false).
		Connect("http://localhost:8529", "_system", "root", rootPassword)

	_, _ = db.Run(&arangolite.CreateDatabase{
		Name: "testDB",
		Users: []map[string]interface{}{
			{"username": "root", "passwd": rootPassword},
			{"username": "user", "passwd": userPassword},
		},
	})

	db.SwitchDatabase("testDB").SwitchUser("user", userPassword)

	_, _ = db.Run(&arangolite.CreateCollection{Name: "nodes"})

	key := "48765564346"

	q := arangolite.NewQuery(`
    FOR n
    IN nodes
    FILTER n._key == %s
    RETURN n
  `, key).Cache(true).BatchSize(500) // The caching feature is unavailable prior to ArangoDB 2.7

	// The Run method returns all the query results of every batches
	// available in the cursor as a slice of byte.
	r, _ := db.Run(q)

	nodes := []Node{}
	json.Unmarshal(r, &nodes)

	// The RunAsync method returns a Result struct allowing to handle batches as they
	// are retrieved from the database.
	async, _ := db.RunAsync(q)

	nodes = []Node{}
	decoder := json.NewDecoder(async.Buffer())

	for async.HasMore() {
		batch := []Node{}
		decoder.Decode(&batch)
		nodes = append(nodes, batch...)
	}

	fmt.Printf("%v", nodes)
}
