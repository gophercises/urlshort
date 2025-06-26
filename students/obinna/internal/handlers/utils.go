package handlers

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
)

type PathUrl struct {
	Path string `json:"path" yaml:"path"`
	Url  string `json:"url" yaml:"url"`
}

func buildMap(pathURLs []PathUrl) map[string]string {
	m := make(map[string]string)
	for _, pu := range pathURLs {
		m[pu.Path] = pu.Url
	}
	return m
}

type DecoderFunc func(r io.Reader, out interface{}) error

func YamlDecoder(r io.Reader, out interface{}) error {
	return yaml.NewDecoder(r).Decode(out)
}

func JsonDecoder(r io.Reader, out interface{}) error {
	return json.NewDecoder(r).Decode(out)
}

func UnmarshalFile(path string, out interface{}, decoder DecoderFunc) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return decoder(f, out)
}

func NewHandlerFromFile(path string, decoder DecoderFunc, fallback http.Handler) (http.HandlerFunc, error) {
	var pathURLs []PathUrl
	err := UnmarshalFile(path, &pathURLs, decoder)
	if err != nil {
		return nil, err
	}
	return MapHandler(buildMap(pathURLs), fallback), nil
}
