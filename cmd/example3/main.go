package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/solher/arangolite"
)

func main() {
	var url string
	flag.StringVar(&url, "url", "http://localhost:8529", "the arrango db url")
	var user string
	flag.StringVar(&user, "user", "user", "the arrango db user")
	var password string
	flag.StringVar(&password, "password", "", "the arrango db password")
	flag.Parse()

	db := arangolite.New()
	db.Connect(url, "testDB", user, password)

	t := arangolite.NewTransaction([]string{"nodes"}, nil).
		AddQuery("nodes", `
    FOR n
    IN nodes
    RETURN n
  `).AddQuery("ids", `
    FOR n
    IN {{.nodes}}
    RETURN n._id
  `).Return("ids")

	r, _ := db.Run(t)

	ids := []string{}
	json.Unmarshal(r, &ids)

	fmt.Printf("%v", ids)
}
