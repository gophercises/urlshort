package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Config is a struct describing the config parsed from cli arguments
type Config struct {
	PathToYAML string
	PathToJSON string
}

func main() {
	config := getConfig()

	yamlBytes := getFileBytes(config.PathToYAML)
	jsonBytes := getFileBytes(config.PathToJSON)

	mux := makeDefaultMux()
	mapHandler := makeMapHandler(mux)

	handler := mapHandler
	if yamlBytes != nil {
		handler = makeYAMLHandler(yamlBytes, &mapHandler)
	} else if jsonBytes != nil {
		handler = makeJSONHandler(jsonBytes, &mapHandler)
	}
	startServer(handler)
}

func getConfig() *Config {
	config := Config{}
	flag.StringVar(&config.PathToYAML, "yaml", "", "--yaml=path/to/file.yml")
	flag.StringVar(&config.PathToJSON, "json", "", "--json=path/to/file.json")
	flag.Parse()
	return &config
}

func getFileBytes(pathToFile string) []byte {
	bytes, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		return nil
	}
	return bytes
}

func makeDefaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloWorldHandler)
	return mux
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func makeMapHandler(mux *http.ServeMux) http.HandlerFunc {
	return MapHandler(map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}, mux)
}

func makeYAMLHandler(yamlBytes []byte, fallbackHandler *http.HandlerFunc) http.HandlerFunc {
	handler, err := YAMLHandler(yamlBytes, fallbackHandler)
	if err != nil {
		panic(err)
	}
	return handler
}

func makeJSONHandler(jsonBytes []byte, fallbackHandler *http.HandlerFunc) http.HandlerFunc {
	handler, err := JSONHandler(jsonBytes, fallbackHandler)
	if err != nil {
		panic(err)
	}
	return handler
}

func startServer(handler http.HandlerFunc) {
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}
