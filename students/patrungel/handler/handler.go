package handler

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v2"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		from := r.URL.Path
		to, ok := pathsToUrls[from]
		if !ok {
			fallback.ServeHTTP(w, r)
			return
		}
		http.Redirect(w, r, to, http.StatusFound)
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
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathsToUrls, err := yamlToMap(yml)
	if err != nil {
		return nil, err
	}
	return MapHandler(pathsToUrls, fallback), nil
}

func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathsToUrls, err := jsonToMap(jsn)
	if err != nil {
		return nil, err
	}
	return MapHandler(pathsToUrls, fallback), nil
}

func BoltHandler(db *bolt.DB, bucketName string, fallback http.Handler) (http.HandlerFunc, error) {
	return func(w http.ResponseWriter, r *http.Request) {
		var to []byte
		from := r.URL.Path

		err := db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(bucketName))
			if bucket == nil {
				return fmt.Errorf("bucket %s not found", bucketName)
			}
			to = bucket.Get([]byte(from))
			return nil
		})
		if err != nil {
			panic(err)
		}

		if to == nil {
			fallback.ServeHTTP(w, r)
			return
		}
		http.Redirect(w, r, string(to), http.StatusFound)
	}, nil
}

func yamlToMap(yml []byte) (map[string]string, error) {
	var records []ShortenerRecord
	err := yaml.Unmarshal(yml, &records)
	if err != nil {
		return nil, err
	}

	return buildMap(records), nil
}

func jsonToMap(jsn []byte) (map[string]string, error) {
	var records []ShortenerRecord
	err := json.Unmarshal(jsn, &records)
	if err != nil {
		return nil, err
	}

	return buildMap(records), nil
}

func buildMap(records []ShortenerRecord) map[string]string {
	mappings := make(map[string]string, len(records))
	for _, record := range records {
		mappings[record.Path] = record.URL
	}
	return mappings
}

type ShortenerRecord struct {
	Path string `json:"path" yaml:"path"`
	URL  string `json:"url" yaml:"url"`
}
