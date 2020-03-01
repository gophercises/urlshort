package urlshort

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/etcd-io/bbolt"
	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if targetURL, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, targetURL, http.StatusPermanentRedirect)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

type yamlFormats struct {
	URL  string `yaml:"url"`
	Path string `yaml:"path"`
}

type jsonFormats struct {
	URL  string `json:"url"`
	Path string `json:"path"`
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
func YAMLHandler(yml []byte, fallback http.Handler) (http.Handler, error) {
	toParseData := make([]yamlFormats, 0)
	pathsToUrls := make(map[string]string)
	if err := yaml.Unmarshal(yml, &toParseData); err != nil {
		return nil, err
	}
	for _, v := range toParseData {
		pathsToUrls[v.Path] = v.URL
	}
	return MapHandler(pathsToUrls, fallback), nil
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
//     {
//       "path": "/some-path",
//       "url": "https://www.some-url.com/demo"
//     }
//
// The only errors that can be returned all related to having
// invalid JSON data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func JSONHandler(jsn []byte, fallback http.Handler) (http.Handler, error) {
	toParseData := make([]jsonFormats, 0)
	pathsToUrls := make(map[string]string)
	if err := json.Unmarshal(jsn, &toParseData); err != nil {
		return nil, err
	}
	for _, v := range toParseData {
		pathsToUrls[v.Path] = v.URL
	}
	return MapHandler(pathsToUrls, fallback), nil
}

func LoadYAML(path string, fallback http.Handler) (http.Handler, error) {
	f, err := os.Open(path)
	if err != nil {
		return fallback, err
	}
	registerClose(f.Close)
	yml, err := ioutil.ReadAll(f)
	if err != nil {
		return fallback, err
	}
	return YAMLHandler(yml, fallback)
}

func LoadJSON(path string, fallback http.Handler) (http.Handler, error) {
	f, err := os.Open(path)
	if err != nil {
		return fallback, err
	}
	registerClose(f.Close)
	jsond, err := ioutil.ReadAll(f)
	if err != nil {
		return fallback, err
	}
	return JSONHandler(jsond, fallback)
}

const (
	BBOLT_DB_NAME     = "BBOLT_DB"
	BBOLT_BUCKET_NAME = "BBOLT_BUCKET"
)

func LoadBBolt(path string, fallback http.Handler) (http.Handler, error) {
	db, err := bbolt.Open(path, 0666, &bbolt.Options{ReadOnly: true, NoFreelistSync: true})
	if err != nil {
		return fallback, err
	}
	registerClose(db.Close)
	tx, err := db.Begin(false)
	if err != nil {
		return fallback, err
	}
	// returns an empty pointer if it does not exist
	bucket := tx.Bucket([]byte(BBOLT_BUCKET_NAME))
	if bucket == nil {
		return fallback, nil
	}
	cursor := bucket.Cursor()
	urlMaps := make(map[string]string)
	for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
		urlMaps[string(k)] = string(v)
	}
	return MapHandler(urlMaps, fallback), nil
}

func registerClose(f func() error) {
	closeHandle = append(closeHandle, f)
}

var closeHandle = make([]func() error, 0)

// Close
func Close() error {
	var err error
	for _, f := range closeHandle {
		err = errors.Unwrap(f())
	}

	return err
}
