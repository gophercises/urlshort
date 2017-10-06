package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gophercises/urlshort/students/latentgenius"
)

var (
	yamlPath string
)

func init() {
	flag.StringVar(&yamlPath, "yaml", "", "path to yaml file")
	flag.Parse()
}

func main() {
	mux := defaultMux()
	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := latentgenius.MapHandler(pathsToUrls, mux)

	if yamlPath != "" {
		yamlData, err := ioutil.ReadFile(yamlPath)
		if err != nil {
			log.Fatalln(err)
		}
		// Build the YAMLHandler using the mapHandler as the
		// fallback
		yamlHandler, err := latentgenius.YAMLHandler(yamlData, mapHandler)
		if err != nil {
			panic(err)
		}
		fmt.Println("Starting the server on :8080")
		http.ListenAndServe(":8080", yamlHandler)
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
