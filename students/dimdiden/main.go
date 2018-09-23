package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dimdiden/gophercises/urlshort/students/dimdiden/urlshort"
)

// DEFAULTFILE is the yaml file expected to be loaded by default
const DEFAULTFILE = "map.yaml"

func main() {
	// Flag block
	file := flag.String("f", DEFAULTFILE, "specify the path to json or yaml file")
	useDB := flag.Bool("d", false, "Enable DB usage. Any provided files will be ignored")
	flag.Parse()

	mux := defaultMux()
	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// pairProducer will be used in the MainHandler
	var pairProducer urlshort.PairProducer
	// get content from DB or from files
	if !*useDB {
		// get content from file
		var err error
		pairProducer, err = getContent(*file)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// Open db
		db, err := urlshort.OpenBDB("my.db", 0600)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		// load test data
		if err := db.LoadInitData(); err != nil {
			log.Fatal(err)
		}
		pairProducer = db
	}

	// mainHandler will be used as in ListenAndServe
	mainHandler, err := urlshort.MainHandler(pairProducer, mapHandler)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", mainHandler)
}

// getContent opens file and returns urlshort.Content
func getContent(file string) (urlshort.Content, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return urlshort.Content(content), nil
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
