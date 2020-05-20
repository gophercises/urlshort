package urlshort

import (
	"net/http"

	"github.com/gophercises/urlshort/jsonloader"
	"github.com/gophercises/urlshort/loader"
	"github.com/gophercises/urlshort/yamlloader"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri, ok := pathsToUrls[r.URL.RequestURI()]

		if !ok {
			fallback.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, uri, http.StatusPermanentRedirect)
	})
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
	loader := yamlloader.NewYamlLoader(yml)

	return GenericHandlerFromLoader(loader, fallback)
}

func JSONHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {
	loader := jsonloader.NewJSONLoader(json)

	return GenericHandlerFromLoader(loader, fallback)
}

func GenericHandlerFromLoader(loader loader.Loader, fallback http.Handler) (http.HandlerFunc, error) {
	urls, err := loader.ToURLsMap()

	if err != nil {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fallback.ServeHTTP(w, r)
		})

		return handler, err
	}

	return MapHandler(urls, fallback), nil
}
