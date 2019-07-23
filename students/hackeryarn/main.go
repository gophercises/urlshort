package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gophercises/urlshort/students/hackeryarn/urlshort"
)

const (
	// YAMLFlag is used to set urls for yaml
	YAMLFlag = "yaml"
	// YAMLFlagValue is the value used when no YAMLFlag is provided
	YAMLFlagValue = "urls.yaml"
	// YAMLFlagUsage is the help string for the YAMLFlag
	YAMLFlagUsage = "URLs file in yaml format"

	// JSONFlag is used to set a file for the questions
	JSONFlag = "json"
	// JSONFlagValue is the value used when no JSONFlag is provided
	JSONFlagValue = "urls.json"
	// JSONFlagUsage is the help string for the JSONFlag
	JSONFlagUsage = "URLs file in json format"
)

// Flagger is an interface for configuring various flags
type Flagger interface {
	StringVar(p *string, name, value, usage string)
}

type urlshortFlagger struct{}

func (uf *urlshortFlagger) StringVar(p *string, name, value, usage string) {
	flag.StringVar(p, name, value, usage)
}

var yaml string
var json string

// ConfigFlags will configure the flags used by the application
func ConfigFlags(flagger Flagger) {
	flagger.StringVar(&yaml, YAMLFlag, YAMLFlagValue, YAMLFlagUsage)
	flagger.StringVar(&json, JSONFlag, JSONFlagValue, JSONFlagUsage)
}

func main() {
	flagger := &urlshortFlagger{}
	ConfigFlags(flagger)

	mapHandler := createMapHandler()
	yamlHandler := createYAMLHandler(mapHandler)
	jsonHandler := createJSONHandler(yamlHandler)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

// Build the MapHandler using the mux as the fallback
var pathsToUrls = map[string]string{
	"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
	"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
}

func createMapHandler() http.HandlerFunc {
	mux := defaultMux()
	return urlshort.MapHandler(pathsToUrls, mux)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func createYAMLHandler(fallback http.HandlerFunc) http.HandlerFunc {
	yamlFile, err := os.Open(yaml)
	if err != nil {
		panic(err)
	}

	yamlHandler, err := urlshort.YAMLHandler(yamlFile, fallback)
	if err != nil {
		panic(err)
	}

	return yamlHandler
}

func createJSONHandler(fallback http.HandlerFunc) http.HandlerFunc {
	jsonFile, err := os.Open(json)
	if err != nil {
		panic(err)
	}

	jsonHandler, err := urlshort.JSONHandler(jsonFile, fallback)
	if err != nil {
		panic(err)
	}

	return jsonHandler
}
