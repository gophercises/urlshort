package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"urlshort/urlshort"
)

const urlMapYaml = "assets/urlmap.yaml"
const urlMapJson = "assets/urlmap.json"
var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	yamlFile, err := os.ReadFile(urlMapYaml)
	if err != nil {
		msg := fmt.Sprintf("Error while reading file '%s': ", urlMapYaml)
		logger.Error(fmt.Sprintf("%s%v", msg, err))
	}

	jsonFile, err := os.ReadFile(urlMapJson)
	if err != nil {
		msg := fmt.Sprintf("Error while reading file '%s': ", urlMapJson)
		logger.Error(fmt.Sprintf("%s%v", msg, err))
	}

	yamlInput := string(yamlFile)
	jsonInput := string(jsonFile)

	// Build the YAMLHandler using the mapHandler as the fallback
	yamlHandler, err := urlshort.YAMLHandler([]byte(yamlInput), mapHandler)
	if err != nil {
		panic(err)
	}

	// Build the JSONHandler using the YAMLHandler as the fallback
	jsonHandler, err := urlshort.JSONHandler([]byte(jsonInput), yamlHandler)
	if err != nil {
		panic(err)
	}

	logger.Info("Starting the server on :8080")
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
