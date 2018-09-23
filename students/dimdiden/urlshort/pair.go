package urlshort

import (
	"encoding/json"
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

// Pair is the main structure used in the package logic
type Pair struct {
	// Path is a short url representation
	Path string
	// Url is full url address
	Url string
}

// PairProducer is the general interface to produce array of Pair structs
type PairProducer interface {
	Pair() ([]Pair, error)
}

// Content is used to give []byte the possibility of implementing PairProducer
type Content []byte

// Pair tries to parse either yaml or json content
// and converts it to array of Pair structs. If non of the available
// umarshaling methods succeed, a simple error is returned
func (c Content) Pair() ([]Pair, error) {
	var pairs []Pair
	var err error

	if err = yaml.Unmarshal(c, &pairs); err == nil {
		return pairs, nil
	}
	if err = json.Unmarshal(c, &pairs); err == nil {
		return pairs, nil
	}
	return nil, fmt.Errorf("Could not unmarshal file. Available formats: json or yaml: %v", err)
}
