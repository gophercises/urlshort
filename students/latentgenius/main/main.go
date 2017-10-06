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
	jsonPath string
)

func init() {
	flag.StringVar(&yamlPath, "yaml", "", "path to yaml file")
	flag.StringVar(&jsonPath, "json", "", "path to json file")
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
	} else if jsonPath != "" {
		jsonData, err := ioutil.ReadFile(jsonPath)
		if err != nil {
			log.Fatalln(err)
		}
		// Build the JSONHandler using the mapHandler as the
		// fallback
		jsonHandler, err := latentgenius.JSONHandler(jsonData, mapHandler)
		if err != nil {
			panic(err)
		}
		fmt.Println("Starting the server on :8080")
		http.ListenAndServe(":8080", jsonHandler)
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
