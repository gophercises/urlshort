package urlshort

import (
	"fmt"
	"net/http"
	"encoding/json"

	"gopkg.in/yaml.v2"
	"github.com/boltdb/bolt"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//	TODO: Implement this...
	return func(w http.ResponseWriter, r *http.Request) {
			for path, url := range pathsToUrls {
				if(r.URL.Path == path){
					http.Redirect(w, r, url, 302)
				}
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

type urlMap struct  {
	Path string `yaml:"path"`
	Url string `yaml:"url"`
}

func parseYAML(yml []byte) (urls []urlMap, err error) {
	err = yaml.Unmarshal(yml, &urls)
	return urls, err
}

func buildMap(urls []urlMap) (map[string]string) {
	pathToUrls := make(map[string]string)
	for _, url := range urls {
		pathToUrls[url.Path] = url.Url
	}
	return pathToUrls
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yml)
	if(err != nil) {
		return nil, err
	}
	pathsToUrls := buildMap(parsedYaml)
	fmt.Printf("%v\n", pathsToUrls)
	return MapHandler(pathsToUrls, fallback), nil
}

func parseJSON(jsonBlob []byte) (urls []urlMap, err error) {
	err = json.Unmarshal(jsonBlob, &urls)
	return urls, err
}

func JSONHandler(jsonBlob []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJson, err := parseJSON(jsonBlob)
	if err != nil {
		return nil, err
	}

	pathToUrls := buildMap(parsedJson)
	fmt.Printf("%v\n", pathToUrls)
	return MapHandler(pathToUrls, fallback), nil
}

func DBHandler(db *bolt.DB, fallback http.Handler)(http.HandlerFunc, error) {
	
	return func(w http.ResponseWriter, r *http.Request) {
		var url []byte

		db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte("urlshort"))
			url = bucket.Get([]byte(r.URL.Path))
			return nil
		})

		if url != nil {
			http.Redirect(w, r, string(url), 302) 
			return 
		}
	fallback.ServeHTTP(w, r)
	}, nil
}