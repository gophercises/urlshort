package jsonloader

import (
	"encoding/json"

	"github.com/gophercises/urlshort/model"
)

type JSONLoader struct {
	json []byte
}

func NewJSONLoader(json []byte) *JSONLoader {
	return &JSONLoader{json: json}
}

func (loader JSONLoader) ToURLsMap() (map[string]string, error) {
	itemList := model.ItemList{}

	err := json.Unmarshal(loader.json, &itemList)

	if err != nil {
		return nil, err
	}

	urls := map[string]string{}

	for _, item := range itemList {
		urls[item.Path] = item.URL
	}

	return urls, nil
}
