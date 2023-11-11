package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"urlshort/urlshort"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	urlMapFile := "assets/urlmap.yaml"
	yamlFile, err := os.ReadFile(urlMapFile)
	if err != nil {
		msg := fmt.Sprintf("Error while reading file '%s': ", urlMapFile)
		logger.Error(fmt.Sprintf("%s%v", msg, err))
	}

	yamlInput := string(yamlFile)

	// Build the YAMLHandler using the mapHandler as the fallback
	yamlHandler, err := urlshort.YAMLHandler([]byte(yamlInput), mapHandler)
	if err != nil {
		panic(err)
	}
	logger.Info("Starting the server on :8080")
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
