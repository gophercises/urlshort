package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/jtschelling/urlshort/students/jtschelling/pkg"
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
	var config []byte
	var err error
	var redirectHandler http.HandlerFunc
	var db *buntdb.DB

	configFile := flag.String("config", "", "Name of the yml or json file containing redirects")
	configType := flag.String("type", "", "Type of the config file. Valid types:\n\t- yaml\n\t- json\n\t- db")
	flag.Parse()

	if (*configType == "") {
		commandLineError("Parse type must be specified with the '--type' flag")
	}

	if (*configType == "db" && *configFile != "") {
		commandLineError("When using the db parse type you cannot specify a file name with '--config'")
	} else if (*configFile == "" && *configType == "json" || *configType == "yaml") {
		commandLineError("A config file must be specified with '--config' when parsing types json or yaml")
	}

	if (*configType == "yaml" || *configType == "json") {
		if (*configType == "yaml" && !strings.Contains(*configFile, "yml")) {
			commandLineError("When type is yaml a *.yml file must be specified")
		}

		if (*configType == "json" && !strings.Contains(*configFile, "json")) {
			commandLineError("When type is json a *.json file must be specified")
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
		"/github": "https://github.com",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

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

func commandLineError(s string) {
	fmt.Println(s)
	os.Exit(1)
}
