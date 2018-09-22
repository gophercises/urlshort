package urlshort

import (
	"encoding/json"
	"errors"
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

// FileMapHandler will parse the provided YAML or JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML or JSON, then the
// fallback http.Handler will be called instead.
func FileMapHandler(content []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsed, err := parse(content)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsed)
	return MapHandler(pathMap, fallback), nil
}

// parse will try to parse either yaml or json file and
// return array of PathMap or error
func parse(content []byte) ([]PathMap, error) {
	var pathmaps []PathMap
	var err error

	if err = yaml.Unmarshal(content, &pathmaps); err == nil {
		return pathmaps, nil
	}
	if err = json.Unmarshal(content, &pathmaps); err == nil {
		return pathmaps, nil
	}
	return nil, errors.New("Could not unmarshal file. Available formats: json or yaml")
}

// buildMap will convert []PathMap to map[string]string
func buildMap(parsedYaml []PathMap) map[string]string {
	resultMap := make(map[string]string)
	for _, y := range parsedYaml {
		resultMap[y.Path] = y.Url
	}
	return resultMap
}
