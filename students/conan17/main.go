package main

import (
	"flag"
	"fmt"
	"net/http"

	urlshort "./urlshort"
)

var (
	boltPath = flag.String("bolt", "./urls.db", "short linked data using the bolt database. ")
	yamlPath = flag.String("yaml", "./urls.yml", "short linked data using the yaml file. ")
	jsonPath = flag.String("json", "./urls.json", "short linked data using the json file. ")
)

func main() {
	flag.Parse()
	defer urlshort.Close()
	mux := defaultMux()
	handle := handleErr(urlshort.LoadBBolt(*boltPath, mux))
	handle = handleErr(urlshort.LoadYAML(*yamlPath, handle))
	handle = handleErr(urlshort.LoadJSON(*jsonPath, handle))
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handle)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func handleErr(handler http.Handler, err error) http.Handler {
	if err != nil {
		panic(err)
	}
	return handler
}
