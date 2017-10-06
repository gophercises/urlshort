package main

import (
	"net/http"

	"encoding/json"

	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if path, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, path, http.StatusFound)
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yamldata []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pathsToUrls []struct {
		Path string `yaml:"path"`
		URL  string `yaml:"url"`
	}
	if err := yaml.Unmarshal(yamldata, &pathsToUrls); err != nil {
		return nil, err
	}
	return func(w http.ResponseWriter, r *http.Request) {
		for _, pathtourl := range pathsToUrls {
			if pathtourl.Path == r.URL.Path {
				http.Redirect(w, r, pathtourl.URL, http.StatusFound)
				return
			}
		}
		fallback.ServeHTTP(w, r)
	}, nil
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
//     {
//       "path": "/some-path",
//       "url": "https://www.some-url.com/demo"
//     }
//
// The only errors that can be returned all related to having
// invalid JSON data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func JSONHandler(jsondata []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pathsToUrls []struct {
		Path string `json:"path"`
		URL  string `json:"url"`
	}
	if err := json.Unmarshal(jsondata, &pathsToUrls); err != nil {
		return nil, err
	}
	return func(w http.ResponseWriter, r *http.Request) {
		for _, pathtourl := range pathsToUrls {
			if pathtourl.Path == r.URL.Path {
				http.Redirect(w, r, pathtourl.URL, http.StatusFound)
				return
			}
		}
		fallback.ServeHTTP(w, r)
	}, nil
}

// BOLTHandler will use the provided BoltDB and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the BoltDB, then the
// fallback http.Handler will be called instead.
//
// BoltDB is expected to be in the format:
//
//     Bucket(pathstourls)
//         /some-path -> https://www.some-url.com/demo
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func BOLTHandler(db *bolt.DB, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte("pathstourls"))
			if bucket != nil {
				cursor := bucket.Cursor()
				for path, url := cursor.First(); path != nil; path, url = cursor.Next() {
					if string(path) == r.URL.Path {
						http.Redirect(w, r, string(url), http.StatusFound)
						return nil
					}
				}
			}
			return nil
		}); err != nil {
			panic(err)
		}
		fallback.ServeHTTP(w, r)
	}
}
