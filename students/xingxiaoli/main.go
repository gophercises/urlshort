package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"xingxiaoli/urlshort"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler, err := urlshort.MapHandler(pathsToUrls, mux)
	if err != nil {
		panic(err)
	}

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := parseFile("./urls.yaml")
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func parseFile(filename string) []byte {
	yaml, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}
	return yaml
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
