package urlshort

import (
	"log/slog"
	"net/http"
)

type ShortenedUrl struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

type ShortenedUrls []ShortenedUrl

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request url: " + r.URL.String())
		path := r.URL.Path
		if url, exists := pathsToUrls[path]; exists {
			slog.Info("Redirecting to " + url)
			http.Redirect(w, r, url, http.StatusMovedPermanently)
			return
		}
		slog.Warn("No url in map")
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
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yamlInput []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yamlInput)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func parseYAML(yamlInput []byte) (ShortenedUrls, error) {
	return nil, nil
}
func buildMap(urls ShortenedUrls) map[string]string {
	return nil
}
