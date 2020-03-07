package main

import (
	"flag"
	"fmt"
	"gophercises/urlshort"
	"io/ioutil"
	"net/http"
	"os"
)

func getFromFile() ([]byte, string, error) {
	yamlFileName := flag.String("yml", "paths.yml", "a yml file with in the format of 'path,url'")
	jsonFileName := flag.String("json", "paths.json", "a json file with in the format of 'path,url'")

	flag.Parse()
	var file []byte
	var err error
	var source string

	if yamlFileName != nil {
		file, err = ioutil.ReadFile(*yamlFileName)
		source = "yaml"
		if err != nil {
			exit(fmt.Sprintf("Failed to open file %s\n", *yamlFileName))
		}
	}

	if jsonFileName != nil {
		file, err = ioutil.ReadFile(*jsonFileName)
		source = "json"
		if err != nil {
			exit(fmt.Sprintf("Failed to open file %s\n", *jsonFileName))
		}
	}

	return file, source, err
}

func exit(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func main() {
	mux := defaultMux()
	var handler http.HandlerFunc
	var err error

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	file, sourceType, err := getFromFile()

	if sourceType == "yaml" {
		handler, err = urlshort.ProxyHandler(file, urlshort.ParseYAML, mapHandler)
	}

	if sourceType == "json" {
		handler, err = urlshort.ProxyHandler(file, urlshort.ParseJSON, mapHandler)
	}

	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
