package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gophercises/urlshort/students/liikt"
)

func main() {
	ymlpath := flag.String("ymlpath", "../storage/map.yml", "Path to a YAML file containing shortened URLs")
	jsonpath := flag.String("jsonpath", "../storage/map.json", "Path to a JSON file containing shortened URLs")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	yaml, err := ioutil.ReadFile(*ymlpath)
	if err != nil {
		panic(err)
	}

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	json, err := ioutil.ReadFile(*jsonpath)
	if err != nil {
		panic(err)
	}

	// Build the JSONHandler using the mapHandler as the
	// fallback
	jsonHandler, err := urlshort.JSONHandler([]byte(json), yamlHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
