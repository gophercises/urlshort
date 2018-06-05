package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gophercises/urlshort/students/hackeryarn/urlshort"
)

const (
	// YamlFlag is used to set a file for the questions
	YamlFlag = "yaml"
	// YamlFlagValue is the value used when no YamlFlag is provided
	YamlFlagValue = "urls.yaml"
	// YamlFlagUsage is the help string for the YamlFlag
	YamlFlagUsage = "URLs file in yaml format"
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

// ConfigFlags will configure the flags used by the application
func ConfigFlags(flagger Flagger) {
	flagger.StringVar(&yaml, YamlFlag, YamlFlagValue, YamlFlagUsage)
}

func main() {
	flagger := &urlshortFlagger{}
	ConfigFlags(flagger)

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)
	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yamlFile, err := os.Open(yaml)
	if err != nil {
		panic(err)
	}

	yamlHandler, err := urlshort.YAMLHandler(yamlFile, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
