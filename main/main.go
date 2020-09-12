package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"../../urlshort"
	"github.com/tidwall/buntdb"
)

func openBuntDB() *buntdb.DB {
	// open in diskless mode
	db, err := buntdb.Open(":memory:")
	if err != nil {
		log.Fatal(err)
	}

	seedBuntDB(db)
	return db
}

func seedBuntDB(db *buntdb.DB) {
	err := db.Update(func(tx *buntdb.Tx) error {
		tx.Set("/db", "https://github.com/tidwall/buntdb", nil)
		tx.Set("/db-docs", "https://pkg.go.dev/github.com/tidwall/buntdb", nil)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	configFile := flag.String("config", "", "Name of the yml or json file containing redirects")
	configType := flag.String("type", "", "Type of the config file. Valid types:\n\t- yaml\n\t- json\n\t- db")
	flag.Parse()

	if (*configType == "") {
		fmt.Println("Parse type must be specified with the '--type' flag")
		os.Exit(1)
	}

	if (*configType == "db" && *configFile != "") {
		fmt.Println("When using the db parse type you cannot specify a file name with '--config'")
		os.Exit(1)
	} else if (*configFile == "" && *configType == "json" || *configType == "yaml") {
		fmt.Println("A config file must be specified with '--config' when parsing types json or yaml")
		os.Exit(1)
	}

	var config []byte
	var err error
	if (*configType == "yaml" || *configType == "json") {
		if (*configType == "yaml" && !strings.Contains(*configFile, "yml")) {
			fmt.Println("When type is yaml a *.yml file must be specified")
			os.Exit(1)
		}

		if (*configType == "json" && !strings.Contains(*configFile, "json")) {
			fmt.Println("When type is json a *.json file must be specified")
			os.Exit(1)
		}

		workDir, err := os.Getwd()
		if err != nil {
			 log.Fatal(err)
		}

		config, err = ioutil.ReadFile(workDir + "/" + *configFile) // For read access.
		if err != nil {
			 log.Fatal(err)
		}
	}

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	var redirectHandler http.HandlerFunc
	var db *buntdb.DB
	switch *configType {
	case "yaml":
		redirectHandler, err = urlshort.YAMLHandler([]byte(config), mapHandler)
	case "json":
		redirectHandler, err = urlshort.JSONHandler([]byte(config), mapHandler)
	case "db":
		db = openBuntDB()
		redirectHandler, err = urlshort.BuntDBHandler(db, mapHandler)
	}
	if err != nil {
		panic(err)
	}

	defer db.Close()
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", redirectHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
