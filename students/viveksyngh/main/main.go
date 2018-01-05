package main

import (
	"io/ioutil"
	"flag"
	"fmt"
	"net/http"

	"github.com/gophercises/urlshort"
)


func main() {
	var ymlFile string
	flag.StringVar(&ymlFile, "yml", "urls.yml", "YAML file having path to url mapping")
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
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	var jsonBlob = []byte(`[
		{"path": "/my-urlshort", "url": "https://github.com/viveksyngh/urlshort"},
		{"path": "/my-github",    "url": "https://github.com/viveksyngh"}
	]`)
	
	jsonHandler, err := urlshort.JSONHandler([]byte(jsonBlob), yamlHandler)
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
