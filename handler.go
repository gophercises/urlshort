package urlshort

import (
	"gopkg.in/yaml.v2"
	"net/http"
)

func MapHandler(paths map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if destination, ok := paths[path]; ok {
			http.Redirect(w, r, destination, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func YAMLHandler(yamlinbytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pathandurl []urlpath
	err := yaml.Unmarshal(yamlinbytes, &pathandurl)
	if err != nil {
		return nil, err
	}
	paths := make(map[string]string)
	for _, pu := range pathandurl {
		paths[pu.Path] = pu.URL
	}
	return MapHandler(paths, fallback), nil
}

type urlpath struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}
