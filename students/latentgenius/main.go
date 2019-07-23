package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gophercises/urlshort/students/latentgenius/handlers"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	yamlPath string
	jsonPath string
	flagDB   string
)

func init() {
	flag.StringVar(&yamlPath, "yaml", "", "path to yaml file")
	flag.StringVar(&jsonPath, "json", "", "path to json file")
	flag.StringVar(&flagDB, "db", "urls.db", "path to sqlite3 database file")
	flag.Parse()
}

func main() {
	mux := defaultMux()
	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := handlers.MapHandler(pathsToUrls, mux)

	if yamlPath != "" {
		yamlData, err := ioutil.ReadFile(yamlPath)
		if err != nil {
			log.Fatalln(err)
		}
		// Build the YAMLHandler using the mapHandler as the
		// fallback
		yamlHandler, err := handlers.YAMLHandler(yamlData, mapHandler)
		if err != nil {
			log.Fatalln("Something went wrong: ", err)
		}
		log.Println("Starting the server on :8080")
		http.ListenAndServe(":8080", yamlHandler)
	} else if jsonPath != "" {
		jsonData, err := ioutil.ReadFile(jsonPath)
		if err != nil {
			log.Fatalf("Could not read file %s: %v\n", jsonPath, err)
		}
		// Build the JSONHandler using the mapHandler as the
		// fallback
		jsonHandler, err := handlers.JSONHandler(jsonData, mapHandler)
		if err != nil {
			log.Fatalln("Something went wrong: ", err)
		}
		log.Println("Starting the server on :8080")
		http.ListenAndServe(":8080", jsonHandler)
	} else {
		db, err := gorm.Open("sqlite3", flagDB)
		if err != nil {
			log.Fatalf("Could not open database: %v", err)
		}
		defer db.Close()
		dbHandler, err := handlers.DBHandler(db, mapHandler)
		if err != nil {
			log.Fatalln(err)
		}
		http.ListenAndServe(":8080", dbHandler)
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
