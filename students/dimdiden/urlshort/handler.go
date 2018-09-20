package dimdiden

import (
	"fmt"
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

type PathMap struct {
	Path string
	Url  string
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if val, ok := pathsToUrls[r.URL.String()]; ok {
			fmt.Println("Mached: ", val)
			http.Redirect(w, r, val, 301)
			return
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
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

// parseYAML will parse the yaml file into []PathMap
func parseYAML(yml []byte) ([]PathMap, error) {
	// https://codebeautify.org/yaml-to-json-xml-csv
	var parsedYaml []PathMap
	err := yaml.Unmarshal(yml, &parsedYaml)
	if err != nil {
		return nil, err
	}
	return parsedYaml, nil
}

// buildMap will convert []PathMap to map[string]string
func buildMap(parsedYaml []PathMap) map[string]string {
	resultMap := make(map[string]string)
	for _, y := range parsedYaml {
		resultMap[y.Path] = y.Url
	}
	return resultMap
}
