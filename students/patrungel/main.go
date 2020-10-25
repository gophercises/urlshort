package main

import (
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	urlshort "github.com/gophercises/patrungel/urlshort/handler"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	var (
		jsonMappingPath, yamlMappingPath string
		genDB                            bool
	)
	flag.StringVar(&jsonMappingPath, "json", "", "Path to json file with mappings")
	flag.StringVar(&yamlMappingPath, "yaml", "", "Path to yaml file with mappings")
	flag.BoolVar(&genDB, "gen-db", false, "Populate a bolt db and exit")
	flag.Parse()

	const (
		pathDB     = "mappings.db"
		bucketName = "mappings"
	)
	if genDB {
		fmt.Println("Generating database entries")
		err := populateDB(pathDB, bucketName)
		if err != nil {
			fmt.Printf("Failed to generate database entries: %s", err)
		} else {
			fmt.Println("Done generating database entries")
		}
		return
	}

	db, err := bolt.Open(pathDB, 0600, &bolt.Options{ReadOnly: true, Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}
	defer db.Close()

	yamlDefault := []byte(`
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`)

	yaml, err := getFromFile(yamlMappingPath, yamlDefault)
	if err != nil {
		panic(err)
	}

	json, err := getFromFile(jsonMappingPath, []byte("[]"))
	if err != nil {
		panic(err)
	}

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yamlHandler, err := urlshort.YAMLHandler(yaml, mapHandler)
	if err != nil {
		panic(err)
	}

	jsonHandler, err := urlshort.JSONHandler(json, yamlHandler)
	if err != nil {
		panic(err)
	}

	boltHandler, err := urlshort.BoltHandler(db, bucketName, jsonHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", boltHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func getFromFile(path string, defaultContents []byte) ([]byte, error) {
	contents := defaultContents
	if path != "" {
		var err error
		contents, err = ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
	}

	return contents, nil
}
