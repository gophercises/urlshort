package yamlloader

import (
	"github.com/gophercises/urlshort/model"
	"gopkg.in/yaml.v2"
)

type YamlLoader struct {
	yaml []byte
}

func NewYamlLoader(yaml []byte) *YamlLoader {
	return &YamlLoader{yaml: yaml}
}

func (loader YamlLoader) ToURLsMap() (map[string]string, error) {
	itemList := model.ItemList{}

	err := yaml.Unmarshal(loader.yaml, &itemList)

	if err != nil {
		return nil, err
	}

	urls := map[string]string{}

	for _, item := range itemList {
		urls[item.Path] = item.URL
	}

	return urls, nil
}
