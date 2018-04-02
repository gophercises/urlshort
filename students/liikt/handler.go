package urlshort

import (
	"encoding/json"
	"net/http"

	"github.com/boltdb/bolt"

	yaml "gopkg.in/yaml.v2"
)

type urlMap struct {
	paths map[string]string
}

type mapItem struct {
	Path string
	Url  string
}

var globalMap *urlMap = &urlMap{make(map[string]string)}

func (*urlMap) redirect(w http.ResponseWriter, r *http.Request) {
	if url, ok := globalMap.paths[r.URL.String()]; ok {
		http.Redirect(w, r, url, 307)
		return
	}
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback *http.ServeMux) {
	for k, v := range pathsToUrls {
		globalMap.paths[k] = v
		fallback.HandleFunc(k, globalMap.redirect)
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
func YAMLHandler(yml []byte, fallback *http.ServeMux) error {
	var list []mapItem
	err := yaml.Unmarshal(yml, &list)
	if err != nil {
		return err
	}

	for _, item := range list {
		globalMap.paths[item.Path] = item.Url
		fallback.HandleFunc(item.Path, globalMap.redirect)
	}

	return nil
}

func JSONHandler(jsn []byte, fallback *http.ServeMux) error {
	var list []mapItem
	err := json.Unmarshal(jsn, &list)
	if err != nil {
		return err
	}

	for _, item := range list {
		globalMap.paths[item.Path] = item.Url
		fallback.HandleFunc(item.Path, globalMap.redirect)
	}

	return nil
}

func BoltHandler(path string, fallback *http.ServeMux) error {
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil
	}
	defer db.Close()

	// Insert testdata into a bucket.
	if err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("paths"))
		if err != nil {
			return err
		}

		if err := b.Put([]byte("/git"), []byte("https://github.com/")); err != nil {
			return err
		}
		if err := b.Put([]byte("/radare"), []byte("http://radare.today/posts/using-radare2/")); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("paths"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			globalMap.paths[string(k)] = string(v)
			fallback.HandleFunc(string(k), globalMap.redirect)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
