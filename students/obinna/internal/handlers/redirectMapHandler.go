package handlers

import "net/http"

func MapHandler(redirectMap map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// check if path is in map
		if dest, ok := redirectMap[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound) // redirect to value gotten from key
			return
		}

		// use fallback mapping if key not found
		fallback.ServeHTTP(w, r)
	}
}
