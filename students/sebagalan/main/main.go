package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	urlshort "gophercises/urlshort/students/sebagalan"

	"github.com/boltdb/bolt"
)

func getFromFile() ([]byte, string, error) {
	nameFileName := flag.String("name", "paths.yml", "file with in the format of 'path,url'")
	typeFileName := flag.String("type", "yaml", "type format")

	flag.Parse()
	var file []byte
	var err error
	var source string

	if nameFileName != nil {
		file, err = ioutil.ReadFile(*nameFileName)
		if err != nil {
			exit(fmt.Sprintf("Failed to open file %s\n", *nameFileName))
		}
	}

	source = *typeFileName
	return file, source, err
}

func startUpDatabase(databaseName string) *bolt.DB {

	db, err := bolt.Open(databaseName, 600, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func closeDB(db *bolt.DB) {
	db.Close()
}

func exit(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func main() {
	mux := defaultMux()
	var handler http.HandlerFunc
	var err error
	var parser func(file []byte) (urlshort.TT, error)

	db := startUpDatabase("./urlshort.db")
	defer closeDB(db)

	mapHandler := urlshort.MapHandler(db, mux)
	file, sourceType, err := getFromFile()

	if sourceType == "yaml" {
		parser = urlshort.ParseYAML
	}

	if sourceType == "json" {
		parser = urlshort.ParseJSON
	}

	parseStreamData, err := func(
		parser func(file []byte) (urlshort.TT, error),
		strm []byte) (urlshort.TT, error) {

		parseStreamData, err := parser(strm)

		if err != nil {
			return nil, err
		}
		return parseStreamData, nil
	}(parser, file)

	handler, err = urlshort.ProxyHandler(db, parseStreamData, mapHandler)

	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
