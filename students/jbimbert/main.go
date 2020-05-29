package main

import (
	"flag"
	"fmt"
	"jimbert/urlshort"
	"net/http"
)

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

// This function is called if all short url mappings fail
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

// run as :
//   ./Ex2_URLShort -f resources/urls.yaml
//   ./Ex2_URLShort -f resources/urls.json
func main() {
	var filename = flag.String("f", "resources/urls.yaml", "The URL file to load")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the fallback
	fileHandler, err := urlshort.FileHandler(filename, mapHandler)
	if err != nil {
		panic(err)
	}

	// Start the server
	fmt.Println("Starting the server on :7894")
	http.ListenAndServe(":7894", fileHandler)
}
