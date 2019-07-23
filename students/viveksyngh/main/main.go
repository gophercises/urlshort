package main

import (
	"io/ioutil"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/gophercises/urlshort"
	"github.com/boltdb/bolt"
)


func main() {
	var ymlFile string
	var jsonFile string
	flag.StringVar(&ymlFile, "yml", "urls.yml", "YAML file having path to url mapping")
	flag.StringVar(&jsonFile, "json", "urls.json", "JSON file for path to url mapping")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml, err := ioutil.ReadFile(ymlFile) 
	if err != nil {
			panic(err)
	}
	
	jsonBlob, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		panic(err)
	}

	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}
		
	jsonHandler, err := urlshort.JSONHandler([]byte(jsonBlob), yamlHandler)
	if err != nil {
		panic(err)
	}
	
	makeDatabaseEntry()

	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}
	defer db.Close()
	dbHandler, err := urlshort.DBHandler(db, jsonHandler)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", dbHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func makeDatabaseEntry() {
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("urlshort"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		err = bucket.Put([]byte("/my-urlshort"), []byte("https://www.github.com/viveksyngh/urlshort"))
		err = bucket.Put([]byte("/my-github"), []byte("https://www.github.com/viveksyngh"))
		return nil
	})
}