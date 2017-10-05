package main

import (
	"fmt"
	"net/http"
	"urlshort"
	"log"
	"flag"
)

func main() {
	yamlFilename := flag.String("yaml-file", "redirect.yaml", "Yaml file name with redirection URLs")
	flag.Parse()

	mux := defaultMux()

	mapHandler := urlshort.NewHttpRedirectHandler(
		urlshort.NewBaseUrlMapper(map[string]string{
			"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
			"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
		}), mux)

	yamlUrlMapper, err := urlshort.NewYamlUrlMapper(*yamlFilename)
	if err != nil {
		log.Fatalf("Can't create YAML redirect URL provider. %v", err)
	}

	yamlHandler := urlshort.NewHttpRedirectHandler(yamlUrlMapper, mapHandler)

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
