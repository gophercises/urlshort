package urlshort

import (
	"fmt"
	"net/http"
)

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

// MainHandler works with PairProducer interface to get pairs, convert
// them to map and invokes MapHandler with this map
func MainHandler(pp PairProducer, fallback http.Handler) (http.HandlerFunc, error) {
	pairs, err := pp.Pair()
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(pairs)
	return MapHandler(pathMap, fallback), nil
}

// buildMap will convert []PathMap to map[string]string
func buildMap(pairs []Pair) map[string]string {
	resultMap := make(map[string]string)
	for _, p := range pairs {
		resultMap[p.Path] = p.Url
	}
	return resultMap
}
