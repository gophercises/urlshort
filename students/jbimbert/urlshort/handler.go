package urlshort

import (
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"encoding/json"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths (keys in the map) to their corresponding URL
// (values that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if newPath, found := pathsToUrls[path]; found {
			http.Redirect(w, r, newPath, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// FileHandler will parse the provided YAML/JSON and then return an http.HandlerFunc
// (which also implements http.Handler)
// that will attempt to map any paths to their corresponding URL.
// If the path is not provided in the YAML/JSON, then the fallback http.Handler will be called instead.
//
// YAML/JSON is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via a mapping of paths to urls.
func FileHandler(filename *string, fallback http.Handler) (http.HandlerFunc, error) {
	// Load the file mappings into a []Url
	urls, err := parseFile(*filename)
	if err != nil {
		return nil, err
	}
	// Transform the []Urls into a map[string]string
	pathMap := buildMap(urls)
	// Pass the map to the MapHandler
	return MapHandler(pathMap, fallback), nil
}

// Structure of the JSON or YAML file
type Url struct {
	Path string
	Url  string
}

// Transform a []Url into a map[string]string
func buildMap(urls []Url) map[string]string {
	m := make(map[string]string)
	for _, v := range urls {
		m[v.Path] = v.Url
	}
	return m
}

// Load a file (YAML or JSON) into a []Url
func parseFile(filename string) ([]Url, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	switch filepath.Ext(filename) {
	case ".json":
		return parseJson(bytes)
	case ".yaml":
		return parseYaml(bytes)
	default:
		return nil, errors.New("Bad file extension")
	}
}

func parseYaml(a []byte) ([]Url, error) {
	var us []struct {
		Path string `yaml:"path"`
		Url  string `yaml:"url"`
	}
	err := yaml.Unmarshal(a, &us)
	if err != nil {
		return nil, err
	}
	// Transform the unmarshalled []us into a []Url
	urls := make([]Url, len(us))
	for i, u := range us {
		urls[i].Path = u.Path
		urls[i].Url = u.Url
	}
	return urls, nil
}

func parseJson(a []byte) ([]Url, error) {
	var us []struct {
		Path string `json:"path"`
		Url  string `json:"url"`
	}
	err := json.Unmarshal(a, &us)
	if err != nil {
		return nil, err
	}
	// Transform the unmarshalled []us into a []Url
	urls := make([]Url, len(us))
	for i, u := range us {
		urls[i].Path = u.Path
		urls[i].Url = u.Url
	}
	return urls, nil
}
