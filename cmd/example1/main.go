package main

import (
	"flag"
	"log"

	ara "github.com/diegogub/aranGO"
)

type DocTest struct {
	ara.Document // Must include arango Document in every struct you want to save id, key, rev after saving it
	Name         string
	Age          int
	Likes        []string
}

func run(url, user, password string) error {
	s, err := setupDB(url, user, password)
	if err != nil {
		return err
	}
	err = createAndRelateDocs(s)
	if err != nil {
		return err
	}

	return nil
}

func setupDB(url, user, password string) (*ara.Session, error) {
	log.Print("setupDB start")
	defer log.Print("setupDB end")

	s, err := ara.Connect(url, user, password, false)
	if err != nil {
		return nil, err
	}

	exists, err := dbExist(s, "test")
	if err != nil {
		return nil, err
	}
	if !exists {
		err = s.CreateDB("test", nil)
		if err != nil {
			return nil, err
		}
	}

	// create Collections test if exist
	if !s.DB("test").ColExist("docs1") {
		// CollectionOptions has much more options, here we just define name , sync
		docs1 := ara.NewCollectionOptions("docs1", true)
		err = s.DB("test").CreateCollection(docs1)
		if err != nil {
			return nil, err
		}
	}

	if !s.DB("test").ColExist("docs2") {
		docs2 := ara.NewCollectionOptions("docs2", true)
		err = s.DB("test").CreateCollection(docs2)
		if err != nil {
			return nil, err
		}
	}

	if !s.DB("test").ColExist("ed") {
		edges := ara.NewCollectionOptions("ed", true)
		edges.IsEdge() // set to Edge
		err = s.DB("test").CreateCollection(edges)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func dbExist(s *ara.Session, name string) (bool, error) {
	dbs, err := s.AvailableDBs()
	if err != nil {
		return false, err
	}
	for _, db := range dbs {
		if db == name {
			return true, nil
		}
	}
	return false, nil
}

func createAndRelateDocs(s *ara.Session) error {
	log.Print("createAndRelateDocs start")
	defer log.Print("createAndRelateDocs end")

	var d1, d2 DocTest
	d1.Name = "Diego"
	d1.Age = 22
	d1.Likes = []string{"arangodb", "golang", "linux"}

	d2.Name = "Facundo"
	d2.Age = 25
	d2.Likes = []string{"php", "linux", "python"}

	err := s.DB("test").Col("docs1").Save(&d1)
	// NOTE: It blocks here.
	if err != nil {
		return err
	}
	log.Print("createAndRelateDocs saved d1")

	err = s.DB("test").Col("docs1").Save(&d2)
	if err != nil {
		return err
	}
	log.Print("createAndRelateDocs saved d2")

	// could also check error in document
	/*
	  if d1.Error {
	    panic(d1.Message)
	  }
	*/

	// update document
	d1.Age = 23
	err = s.DB("test").Col("docs1").Replace(d1.Key, d1)
	if err != nil {
		return err
	}
	log.Print("createAndRelateDocs updated d1")

	// Relate documents
	err = s.DB("test").Col("ed").Relate(d1.Id, d2.Id, map[string]interface{}{"is": "friend"})
	if err != nil {
		return err
	}
	log.Print("createAndRelateDocs related d1 to d2")

	return nil
}

func main() {
	var url string
	flag.StringVar(&url, "url", "http://localhost:8529", "the arrango db url")
	var user string
	flag.StringVar(&user, "user", "", "the arrango db user")
	var password string
	flag.StringVar(&password, "password", "", "the arrango db password")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	err := run(url, user, password)
	if err != nil {
		panic(err)
	}
}
