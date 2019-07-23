package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/DennisVis/urlshort/students/dennisvis/urlshort"
	"github.com/boltdb/bolt"
)

var (
	pathsFile = flag.String("pathsFile", "paths.yml", "The file containing shortened paths to URL's")
	initDB    = flag.Bool("initDB", true, "Whether or not to initialize the paths database")
)

func getFileBytes(fileName string) []byte {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Could not open file %s", fileName)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(f)
	if err != nil {
		log.Fatalf("Could not read file %s", fileName)
	}

	return buf.Bytes()
}

func initBoltDB(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("paths"))
		if err != nil {
			return err
		}

		err = b.Put([]byte("/urlshort"), []byte("https://github.com/gophercises/urlshort"))
		err = b.Put([]byte("/urlshort-final"), []byte("https://github.com/gophercises/urlshort/tree/solution"))

		return err
	})
}

func main() {
	mux := defaultMux()

	flag.Parse()

	ext := filepath.Ext(*pathsFile)

	var handler http.Handler
	var err error
	if ext == ".yml" {
		handler, err = urlshort.YAMLHandler(getFileBytes(*pathsFile), mux)
		if err != nil {
			panic(err)
		}
	} else if ext == ".json" {
		handler, err = urlshort.JSONHandler(getFileBytes(*pathsFile), mux)
		if err != nil {
			panic(err)
		}
	} else if ext == ".db" {
		db, err := bolt.Open(*pathsFile, 0600, nil)
		if err != nil {
			panic(err)
		}

		defer db.Close()

		if *initDB {
			if err = initBoltDB(db); err != nil {
				panic(err)
			}
		}

		handler = urlshort.DBHandler(db, mux)
	} else {
		log.Fatal("Paths file needs to be either a YAML, a JSON or a bolt DB file")
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
