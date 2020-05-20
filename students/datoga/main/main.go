package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gophercises/urlshort"
	"github.com/gophercises/urlshort/dbloader"
	"github.com/gophercises/urlshort/jsonloader"
	"github.com/gophercises/urlshort/loader"
	"github.com/gophercises/urlshort/yamlloader"
)

func main() {

	if *yamlFileFlag != "" && *jsonFileFlag != "" {
		fmt.Println("Both json and yaml must not be provided")
		os.Exit(1)
	}

	var loader loader.Loader
	var err error

	if *yamlFileFlag != "" {
		loader, err = loaderFromYamlFile(*yamlFileFlag)
	} else if *jsonFileFlag != "" {
		loader, err = loaderFromJSONFile(*jsonFileFlag)
	} else {
		dbloader, err := dbloader.NewBoltDBLoader()

		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		defer dbloader.Close()

		dbloader.AddURL("/cosa", "https://test1.largo.es")
		dbloader.AddURL("/cosa2", "https://test2.largo.es")

		loader = dbloader
	}

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback

	genericHandler, err := urlshort.GenericHandlerFromLoader(loader, mapHandler)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("Starting the server on :8080")

	fmt.Println(http.ListenAndServe(":8080", genericHandler))
}

func loaderFromYamlFile(file string) (loader.Loader, error) {

	yaml, err := ioutil.ReadFile(file)

	if err != nil {
		return nil, err
	}

	return yamlloader.NewYamlLoader(yaml), nil

}

func loaderFromJSONFile(file string) (loader.Loader, error) {
	json, err := ioutil.ReadFile(file)

	if err != nil {
		return nil, err
	}

	return jsonloader.NewJSONLoader(json), nil

}

func loaderDefaultYaml() loader.Loader {
	yaml := []byte(`
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
	`)

	return yamlloader.NewYamlLoader(yaml)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
