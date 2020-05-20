package yamlloader

import (
	"reflect"
	"testing"
)

func TestYamlLoader_ToURLsMapNormalYaml(t *testing.T) {
	yaml := `
- path: /path1
  url: https://url1
- path: /path2
  url: https://url2`

	loader := NewYamlLoader([]byte(yaml))

	URLMap, err := loader.ToURLsMap()

	if err != nil {
		t.Errorf("Error not expected: %s", err)
		return
	}

	expectedURLMap := map[string]string{
		"/path1": "https://url1",
		"/path2": "https://url2",
	}

	eq := reflect.DeepEqual(URLMap, expectedURLMap)

	if !eq {
		t.Errorf("Error, map got %v, expected %v\n", URLMap, expectedURLMap)
	}
}

func TestYamlLoader_ToURLsMapEmptyYaml(t *testing.T) {
	yaml := ""

	loader := NewYamlLoader([]byte(yaml))

	URLMap, err := loader.ToURLsMap()

	if err != nil {
		t.Errorf("Error not expected: %s", err)
		return
	}

	if len(URLMap) != 0 {
		t.Errorf("Error, map got %v, expected empty map\n", URLMap)
	}
}

func TestYamlLoader_ToURLsMapIncorrectYaml(t *testing.T) {
	yaml := `
	 cosa:   -- `

	loader := NewYamlLoader([]byte(yaml))

	_, err := loader.ToURLsMap()

	if err == nil {
		t.Errorf("Error expected!")
		return
	}
}
