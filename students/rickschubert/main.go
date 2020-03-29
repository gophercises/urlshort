package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gophersizes/urlshort/students/rickschubert/handler"
)

func main() {
	mux := createDefaultRedirectHandler()

	pathsToUrls := map[string]string{
		"/urlshort": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml":     "https://godoc.org/gopkg.in/yaml.v2",
		"/rick":     "https://rick-schubert.com",
	}
	mapHandler := handler.MapHandler(pathsToUrls, mux)

	fileLocation := "./exampleyaml.yaml"
	yamlFile, err := ioutil.ReadFile(fileLocation)
	exitOnError(err, fmt.Sprintf("Unable to read yaml file at location %s", fileLocation))

	yamlHandler, err := handler.YAMLHandler(yamlFile, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func createDefaultRedirectHandler() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "You are doing really well.")
}

func exitOnError(err error, message string) {
	if err != nil {
		log.Fatal(message)
	}
}
