package urlshort

import (
	"encoding/json"
	"net/http"

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
		if targetURL, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, targetURL, http.StatusPermanentRedirect)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

type yamlFormats struct {
	URL  string `yaml:"url"`
	Path string `yaml:"path"`
}

type jsonFormats struct {
	URL  string `json:"url"`
	Path string `json:"path"`
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
	toParseData := make([]yamlFormats, 0)
	pathsToUrls := make(map[string]string)
	if err := yaml.Unmarshal(yml, &toParseData); err != nil {
		return nil, err
	}
	for _, v := range toParseData {
		pathsToUrls[v.Path] = v.URL
	}
	return MapHandler(pathsToUrls, fallback), nil
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
func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	toParseData := make([]jsonFormats, 0)
	pathsToUrls := make(map[string]string)
	if err := json.Unmarshal(jsn, &toParseData); err != nil {
		return nil, err
	}
	for _, v := range toParseData {
		pathsToUrls[v.Path] = v.URL
	}
	return MapHandler(pathsToUrls, fallback), nil
}
