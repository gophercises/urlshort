package handlers

import (
	"go.etcd.io/bbolt"
	"net/http"
)

const RedirectBucket = "redirects"

func YamlHandler(path string, fallback http.Handler) (http.HandlerFunc, error) {
	return NewHandlerFromFile(path, YamlDecoder, fallback)
}

func JsonHandler(path string, fallback http.Handler) (http.HandlerFunc, error) {
	return NewHandlerFromFile(path, JsonDecoder, fallback)
}

func BoltHandler(db *bbolt.DB, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		var url []byte
		err := db.View(func(tx *bbolt.Tx) error {
			//check that bucket exists
			b := tx.Bucket([]byte(RedirectBucket))
			if b == nil {
				return nil // fallback if bucket doesn't exist
			}
			url = b.Get([]byte(path))
			return nil
		})

		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if url == nil {
			fallback.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, string(url), http.StatusFound)
	}
}
