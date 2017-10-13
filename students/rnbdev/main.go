package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/boltdb/bolt"
)

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func main() {
	var jsonFile, yamlFile, boltFile string
	flag.StringVar(&jsonFile, "json", "", "path to json file.")
	flag.StringVar(&yamlFile, "yaml", "", "path to yaml file.")
	flag.StringVar(&boltFile, "bolt", "", "path to boltdb file.")

	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	mapHandler := MapHandler(pathsToUrls, mux)

	if jsonFile != "" {
		jsonData, err := ioutil.ReadFile(jsonFile)
		if err != nil {
			panic(err)
		}
		// Build the JSONHandler using the mapHandler as the
		// fallback
		jsonHandler, err := JSONHandler([]byte(jsonData), mapHandler)
		if err != nil {
			panic(err)
		}
		fmt.Println("Starting the server on :8080")
		http.ListenAndServe(":8080", jsonHandler)
	} else if yamlFile != "" {
		yamlData, err := ioutil.ReadFile(yamlFile)
		if err != nil {
			panic(err)
		}
		// Build the YAMLHandler using the mapHandler as the
		// fallback
		yamlHandler, err := YAMLHandler([]byte(yamlData), mapHandler)
		if err != nil {
			panic(err)
		}
		fmt.Println("Starting the server on :8080")
		http.ListenAndServe(":8080", yamlHandler)
	} else if boltFile != "" {
		db, err := bolt.Open(boltFile, 0600, nil)
		if err != nil {
			panic(err)
		}
		defer db.Close()
		// Build the BOLTHandler using the mapHandler as the
		// fallback
		boltHandler := BOLTHandler(db, mapHandler)
		fmt.Println("Starting the server on :8080")
		http.ListenAndServe(":8080", boltHandler)
	} else {
		fmt.Println("Starting the server on :8080")
		http.ListenAndServe(":8080", mapHandler)
	}
}
