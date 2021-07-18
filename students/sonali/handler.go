package urlshort

import (
	"fmt"
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
	return func(w http.ResponseWriter, req *http.Request) {
		path := req.URL.Path
		if URL, ok := pathsToUrls[path]; ok {
			http.Redirect(w, req, URL, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, req)
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
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pathUrlMap1 []map[string]string
	errUnmarshal := yaml.Unmarshal(yml, &pathUrlMap1)
	if errUnmarshal != nil {
		fmt.Println("Unmarshal Error: " + errUnmarshal.Error())
		return nil, errUnmarshal
	}
	finalMap := map[string]string{}
	for _, mapping := range pathUrlMap1 {
		finalMap[mapping["path"]] = mapping["url"]
	}
	return MapHandler(finalMap, fallback), nil

}
