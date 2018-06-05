package urlshort

import (
	"io"
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
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
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
func YAMLHandler(r io.Reader, fallback http.Handler) (http.HandlerFunc, error) {
	decoder := yaml.NewDecoder(r)
	pathURLs, err := decodeYaml(decoder)
	if err != nil {
		return nil, err
	}
	pathToUrls := buildMap(pathURLs)

	mapHandler := MapHandler(pathToUrls, fallback)
	return mapHandler, nil
}

type pathURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func decodeYaml(decoder *yaml.Decoder) ([]pathURL, error) {
	var pu []pathURL
	for {
		err := decoder.Decode(&pu)
		if err == io.EOF {
			return pu, nil
		} else if err != nil {
			return nil, err
		}
	}
}

func buildMap(pathURLs []pathURL) map[string]string {
	pathToUrls := make(map[string]string)
	for _, pu := range pathURLs {
		pathToUrls[pu.Path] = pu.URL
	}

	return pathToUrls
}
