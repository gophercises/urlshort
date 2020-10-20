package main

import (
	"flag"
	"fmt"
	urlshort "github.com/gophercises/patrungel/urlshort/handler"
	"io/ioutil"
	"net/http"
)

func main() {
	var jsonMappingPath, yamlMappingPath string
	flag.StringVar(&jsonMappingPath, "json", "", "Path to json file with mappings")
	flag.StringVar(&yamlMappingPath, "yaml", "", "Path to yaml file with mappings")
	flag.Parse()

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
