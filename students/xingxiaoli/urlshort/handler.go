package urlshort

import (
	"net/http"

	json "github.com/json-iterator/go"
	yaml "gopkg.in/yaml.v2"
)

type UrlPath struct {
	Path string `yaml:"path" json:"path"`
	Url  string `yaml:"url" json:"url"`
}

func MapHandler(maps map[string]string, fallback http.Handler) (http.HandlerFunc, error) {
	return func(w http.ResponseWriter, r *http.Request) {
		if url, ok := maps[r.URL.Path]; ok {
			http.Redirect(w, r, url, 301)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}, nil
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var ymls []UrlPath
	err := yaml.Unmarshal(yml, &ymls)
	if err != nil {
		return nil, err
	}

	return Handler(ymls, fallback)
}

func JsonHandler(data []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var res []UrlPath
	err := json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	return Handler(res, fallback)
}

func Handler(urlPath []UrlPath, fallback http.Handler) (http.HandlerFunc, error) {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, up := range urlPath {
			if up.Path == r.URL.Path {
				http.Redirect(w, r, up.Url, 301)
				return
			}
		}
		fallback.ServeHTTP(w, r)
	}, nil
}
