package urlshort

import (
	"net/http"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"github.com/tidwall/buntdb"
)

type RedirectRecord struct {
	Path string
	URL string
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		redirect, exists := pathsToUrls[r.URL.Path]

		if exists {
			http.RedirectHandler(redirect, 301).ServeHTTP(w, r)
		} else {
			fallback.ServeHTTP(w, r)
		}
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
	parsedYaml, err := parseYAML(yml)
  if err != nil {
    return nil, err
  }
  pathMap := buildMap(parsedYaml)
  return MapHandler(pathMap, fallback), nil
}

func parseYAML(yml []byte) ([]RedirectRecord, error) {
	var ret_parsed []RedirectRecord

	err := yaml.Unmarshal(yml, &ret_parsed)

	return ret_parsed, err
}

func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJson, err := parseJSON(jsn)
	if err != nil {
    return nil, err
  }
  pathMap := buildMap(parsedJson)
	return MapHandler(pathMap, fallback), nil
}

func parseJSON(jsn []byte) ([]RedirectRecord, error) {
	var ret_parsed []RedirectRecord

	err := json.Unmarshal(jsn, &ret_parsed)

	return ret_parsed, err
}

func BuntDBHandler(db *buntdb.DB, fallback http.Handler) (http.HandlerFunc, error) {
	parsedBuntDB, err := readBuntDB(db)
	if err != nil {
    return nil, err
  }
  pathMap := buildMap(parsedBuntDB)
	return MapHandler(pathMap, fallback), nil
}

func readBuntDB(db *buntdb.DB) ([]RedirectRecord, error) {
	var ret_parsed []RedirectRecord

	err := db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get("/db")
		ret_parsed = append(ret_parsed, RedirectRecord {
						Path: "/db",
						URL: val,
					})

		val, err = tx.Get("/db-docs")
		ret_parsed = append(ret_parsed, RedirectRecord {
						Path: "/db-docs",
						URL: val,
					})
		return err
	})
	if err != nil {
    return nil, err
  }

	return ret_parsed, nil
}

func buildMap(redirects []RedirectRecord) map[string]string {
	ret_map := make(map[string]string)

	for _, redirect := range redirects {
		ret_map[redirect.Path] = redirect.URL
	}

	return ret_map
}
