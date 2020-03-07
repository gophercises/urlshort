package urlshort

import (
	"encoding/json"
	"net/http"

	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v2"
)

//TT ...
type TT []map[string]string

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls *bolt.DB, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var uri string = r.URL.RequestURI()
		var completeURL []byte

		err := pathsToUrls.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte("paths"))
			completeURL = bucket.Get([]byte(uri))
			return nil
		})

		if err != nil {
			fallback.ServeHTTP(w, r)
		}
		http.Redirect(w, r, string(completeURL), http.StatusPermanentRedirect)
	}
}

//ProxyHandler works like a proxy, parse stream data using a parse funcion and call MapHandler
func ProxyHandler(db *bolt.DB, parseStreamData TT, fallback http.Handler) (http.HandlerFunc, error) {
	var err error
	err = buildMap(parseStreamData, db)

	if err != nil {
		return nil, err
	}

	return MapHandler(db, fallback), nil
}

//ParseYAML ...
func ParseYAML(yml []byte) (TT, error) {

	var yTm TT
	var err error
	err = yaml.Unmarshal(yml, &yTm)

	if err != nil {
		panic(err)
	}

	return yTm, err
}

//ParseJSON ...
func ParseJSON(jsn []byte) (TT, error) {

	var yTm TT
	var err error

	err = json.Unmarshal(jsn, &yTm)

	if err != nil {
		panic(err)
	}

	return yTm, err
}

func buildMap(yTm TT, db *bolt.DB) error {
	var err error
	for _, ptu := range yTm {
		//store in a db
		err = db.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists([]byte("paths"))

			if err != nil {
				return err
			}

			err = bucket.Put([]byte(ptu["path"]), []byte(ptu["url"]))

			if err != nil {
				return err
			}

			return nil

		})
	}

	return err
}
