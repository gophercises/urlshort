package model

type Item struct {
	URL  string `yaml:"url" json:"url"`
	Path string `yaml:"path" json:"path"`
}

type ItemList []Item
