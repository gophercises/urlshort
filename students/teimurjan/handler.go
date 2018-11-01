package main

import (
	"encoding/json"
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		originalURL, ok := pathsToUrls[r.URL.Path]
		if ok {
			http.Redirect(w, r, originalURL, 301)
		} else {
			fallback.ServeHTTP(w, r)
		}
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

type URLMapper struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

func YAMLHandler(YAML []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var mappers []URLMapper
	err := yaml.Unmarshal(YAML, &mappers)
	if err != nil {
		return nil, err
	}
	return func(w http.ResponseWriter, r *http.Request) {
		for _, mapper := range mappers {
			if mapper.Path == r.URL.Path {
				http.Redirect(w, r, mapper.URL, 301)
				return
			}
		}
		fallback.ServeHTTP(w, r)
	}, nil
}

func JSONHandler(JSON []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var mappers []URLMapper
	err := json.Unmarshal(JSON, &mappers)
	if err != nil {
		return nil, err
	}
	return func(w http.ResponseWriter, r *http.Request) {
		for _, mapper := range mappers {
			if mapper.Path == r.URL.Path {
				http.Redirect(w, r, mapper.URL, 301)
				return
			}
		}
		fallback.ServeHTTP(w, r)
	}, nil
}
