package urlshort

import (
	"encoding/json"
	"net/http"

	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v2"
)

type pathToURL struct {
	Path string
	URL  string
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		url := pathsToUrls[req.URL.Path]
		if url != "" {
			http.Redirect(res, req, url, http.StatusPermanentRedirect)
		} else {
			fallback.ServeHTTP(res, req)
		}
	})
}

func buildMap(pathsToURLs []pathToURL) (builtMap map[string]string) {
	builtMap = make(map[string]string)
	for _, ptu := range pathsToURLs {
		builtMap[ptu.Path] = ptu.URL
	}
	return
}

func parseYAML(yamlData []byte) (pathsToURLs []pathToURL, err error) {
	err = yaml.Unmarshal(yamlData, &pathsToURLs)
	return
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
func YAMLHandler(yamlData []byte, fallback http.Handler) (yamlHandler http.HandlerFunc, err error) {
	parsedYaml, err := parseYAML(yamlData)
	if err != nil {
		return
	}
	pathMap := buildMap(parsedYaml)
	yamlHandler = MapHandler(pathMap, fallback)
	return
}

func parseJSON(jsonData []byte) (pathsToURLs []pathToURL, err error) {
	err = json.Unmarshal(jsonData, &pathsToURLs)
	return
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
func JSONHandler(jsonData []byte, fallback http.Handler) (jsonHandler http.HandlerFunc, err error) {
	parsedJSON, err := parseJSON(jsonData)
	if err != nil {
		return
	}
	pathMap := buildMap(parsedJSON)
	jsonHandler = MapHandler(pathMap, fallback)
	return
}

// DBHandler will use the provided Bolt database and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the DB, then the
// fallback http.Handler will be called instead.
func DBHandler(db *bolt.DB, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var url string
		err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("paths"))
			bts := b.Get([]byte(req.URL.Path))
			if bts != nil {
				url = string(bts)
			}
			return nil
		})

		if err == nil && url != "" {
			http.Redirect(res, req, url, http.StatusTemporaryRedirect)
		} else {
			fallback.ServeHTTP(res, req)
		}
	})
}
