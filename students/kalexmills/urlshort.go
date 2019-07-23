package urlshort

import (
	"net/http"
	"fmt"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(out http.ResponseWriter, in *http.Request) {
		if in.Method != http.MethodGet {
			fallback.ServeHTTP(out, in)
			return
		}

		url, ok := pathsToUrls[in.URL.Path]
		if !ok {
			fallback.ServeHTTP(out, in)
			return
		}

		out.Header().Add("Location", url)
		out.WriteHeader(301)
		fmt.Printf("%s %s %d: %s\n", in.Method, in.URL.Path, 301, url)
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
//     pairs:
//     - path: /some-path
//       url: https://www.some-url.com/demo
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	type pair struct {
		Path string
		URL string
	}

	type pairs struct {
		Pairs []pair
	}

	var prs pairs
	err := yaml.Unmarshal(yml, &prs)

	pathsToUrls := make(map[string]string, len(prs.Pairs))
	for _, entry := range prs.Pairs {
		pathsToUrls[entry.Path] = entry.URL
	}

	return MapHandler(pathsToUrls, fallback), err
}