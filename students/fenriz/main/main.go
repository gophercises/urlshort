package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/fenriz07/urlshort"
)

func main() {
	mux := defaultMux()

	yamlfile := flag.String("yamlfile", "redirect.yaml", "name to yaml file examplñe redirect.yaml")
	flag.Parse()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	mapHandler = urlshort.YAMLHandler(*yamlfile, mapHandler)

	mapHandler = urlshort.JSONHandler("redirect.json", mapHandler)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", mapHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
