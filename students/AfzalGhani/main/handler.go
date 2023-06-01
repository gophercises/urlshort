package main

import (
	"fmt"
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

func MapHandler(pathUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if dest, ok := pathUrls[path]; ok {
			http.Redirect(w, r, dest, 201)
			return

		}

		fallback.ServeHTTP(w, r)
	}
}

type YMLMAP struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var ymlMap []YMLMAP

	err := yaml.Unmarshal(yml, &ymlMap)
	fmt.Println("Inside the file")
	fmt.Printf("%v", ymlMap)

	if err != nil {
		return nil, err
	}

	paths := make(map[string]string)

	for _, value := range ymlMap {
		paths[value.Path] = value.Url
	}

	return MapHandler(paths, fallback), nil

}
