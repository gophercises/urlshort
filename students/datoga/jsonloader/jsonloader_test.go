package jsonloader

import (
	"reflect"
	"testing"
)

func TestJSONLoader_ToURLsMapNormalJSON(t *testing.T) {
	json := `[
{"path": "/path1", "url": "https://url1"},
{"path": "/path2", "url": "https://url2"}]`

	loader := NewJSONLoader([]byte(json))

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

func TestJSONLoader_ToURLsMapEmptyJSON(t *testing.T) {
	json := "[]"

	loader := NewJSONLoader([]byte(json))

	URLMap, err := loader.ToURLsMap()

	if err != nil {
		t.Errorf("Error not expected: %s", err)
		return
	}

	if len(URLMap) != 0 {
		t.Errorf("Error, map got %v, expected empty map\n", URLMap)
	}
}

func TestJsonLoader_ToURLsMapIncorrectJSON(t *testing.T) {
	json := "{"

	loader := NewJSONLoader([]byte(json))

	_, err := loader.ToURLsMap()

	if err == nil {
		t.Errorf("Error expected!")
		return
	}
}
