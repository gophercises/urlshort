package urlshort

import (
	"encoding/json"
	"io"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(
	pathsToUrls map[string]string,
	fallback http.Handler,
) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		url := pathsToUrls[req.URL.Path]
		if url == "" {
			fallback.ServeHTTP(res, req)
		} else {
			http.Redirect(res, req, url, http.StatusPermanentRedirect)
		}
	}
}

// PathToURL provides data for a redirect
type PathToURL struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url"  json:"url"`
}

func pathsToURLsToMap(pathsToURLs []PathToURL) map[string]string {
	result := make(map[string]string)

	for _, p := range pathsToURLs {
		result[p.Path] = p.URL
	}

	return result
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
func YAMLHandler(
	reader io.Reader,
	fallback http.Handler,
) (http.HandlerFunc, error) {
	pathsToURLs := []PathToURL{}

	decoder := yaml.NewDecoder(reader)
	err := decoder.Decode(&pathsToURLs)
	if err != nil {
		return nil, err
	}

	pathToURLsMap := pathsToURLsToMap(pathsToURLs)

	return MapHandler(pathToURLsMap, fallback), nil
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
//     [
//	     {
//         "path": "/some-path",
//         "url": "https://www.some-url.com/demo"
//       }
//	   ]
//
// The only errors that can be returned all related to having
// invalid JSON data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func JSONHandler(
	reader io.Reader,
	fallback http.Handler,
) (http.HandlerFunc, error) {
	pathsToURLs := []PathToURL{}

	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&pathsToURLs)
	if err != nil {
		return nil, err
	}

	pathsToURLsMap := pathsToURLsToMap(pathsToURLs)

	return MapHandler(pathsToURLsMap, fallback), nil
}
