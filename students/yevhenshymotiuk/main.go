package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"

	"shortener/urlshort"
)

func defaultMux(filePath string) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler(filePath))
	return mux
}

func indexHandler(filePath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, fmt.Sprintf("Using paths from %s", filePath))
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}

func run() error {
	filePath := flag.String("f", "paths.yaml", "file path")
	flag.Parse()

	mux := defaultMux(*filePath)

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	f, err := os.Open(*filePath)
	r := bufio.NewReader(f)

	var pathsHandler http.HandlerFunc

	switch {
	case path.Ext(*filePath) == ".json":
		pathsHandler, err = urlshort.JSONHandler(
			r,
			mapHandler,
		)
	default:
		pathsHandler, err = urlshort.YAMLHandler(
			r,
			mapHandler,
		)
	}
	if err != nil {
		return err
	}

	port := 8080
	fmt.Println(fmt.Sprintf("Listening on :%d", port))
	return http.ListenAndServe(fmt.Sprintf(":%d", port), pathsHandler)
}
