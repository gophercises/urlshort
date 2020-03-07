package urlshort

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"
)

type tt []map[string]string

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//	TODO: Implement this...

	//return nil
	return func(w http.ResponseWriter, r *http.Request) {
		var uri string = r.URL.RequestURI()
		if completeURL, ok := pathsToUrls[uri]; ok {
			http.Redirect(w, r, completeURL, http.StatusPermanentRedirect)
		}

		fallback.ServeHTTP(w, r)
	}
}

//ProxyHandler works like a proxy, parse stream data using a parse funcion and call MapHandler
func ProxyHandler(strm []byte, parser func(strm []byte) (tt, error), fallback http.Handler) (http.HandlerFunc, error) {

	parseStreamData, err := parser(strm)

	if err != nil {
		return nil, err
	}

	pathsToUrls := buildMap(parseStreamData)
	return MapHandler(pathsToUrls, fallback), nil
}

func ParseYAML(yml []byte) (tt, error) {

	var yTm tt

	err := yaml.Unmarshal(yml, &yTm)

	if err != nil {
		panic(err)
	}

	return yTm, err
}

func ParseJSON(jsn []byte) (tt, error) {

	var yTm tt

	err := json.Unmarshal(jsn, &yTm)

	if err != nil {
		panic(err)
	}

	return yTm, err
}

func buildMap(yTm tt) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, ptu := range yTm {
		pathsToUrls[ptu["path"]] = ptu["url"]
	}

	return pathsToUrls
}
