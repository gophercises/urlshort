package handler

import (
	"fmt"
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

func lookupRedirectInMapAndActOnIt(mapOfRedirects map[string]string, writer http.ResponseWriter, req *http.Request, fallback http.Handler) {
	redirectUrl, found := mapOfRedirects[req.URL.Path]
	if found {
		fmt.Println(fmt.Sprintf("Redirecting to %s", redirectUrl))
		http.Redirect(writer, req, redirectUrl, 317)
	} else {
		fallback.ServeHTTP(writer, req)
	}
}

func newHandler(urlsMap map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		lookupRedirectInMapAndActOnIt(urlsMap, writer, req, fallback)
	}
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return newHandler(pathsToUrls, fallback)
}

type redirectConfig struct {
	Url  string `yaml:"url"`
	Path string `yaml:"path"`
}

type redirectConfigs []redirectConfig

func (c *redirectConfigs) parseFromYaml(data []byte) error {
	return yaml.Unmarshal(data, c)
}

func (c redirectConfigs) convertToMap() map[string]string {
	redirectsMap := make(map[string]string)
	for _, r := range c {
		redirectsMap[r.Path] = r.Url
	}
	return redirectsMap
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var redirects redirectConfigs
	err := redirects.parseFromYaml(yml)
	redirectsMap := redirects.convertToMap()
	return newHandler(redirectsMap, fallback), err
}
