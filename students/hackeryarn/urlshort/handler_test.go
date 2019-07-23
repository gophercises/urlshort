package urlshort

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	fallbackResponse = "fallback"
	path             = "/test"
	dest             = "https://test.com"
)

func fallback(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, fallbackResponse)
}

func TestMapHandler(t *testing.T) {
	pathToUrls := map[string]string{path: dest}

	t.Run("it uses the fallback for unknown routes", func(t *testing.T) {
		result := runMapHandler(pathToUrls, "/unknown")

		assertBody(t, result, fallbackResponse)
	})

	t.Run("it redirects for found url", func(t *testing.T) {
		result := runMapHandler(pathToUrls, path)

		assertStatus(t, result, http.StatusFound)
		assertURL(t, result, dest)
	})
}

func TestYAMLHandler(t *testing.T) {
	yaml := fmt.Sprintf(`
  - path: %s
    url: %s
  `, path, dest)

	t.Run("it uses the fallback for unknown routes", func(t *testing.T) {
		result := runReaderHandler(YAMLHandler, yaml, "/unknown")

		assertBody(t, result, fallbackResponse)
	})

	t.Run("it redirects for found url", func(t *testing.T) {
		result := runReaderHandler(YAMLHandler, yaml, path)

		assertStatus(t, result, http.StatusFound)
		assertURL(t, result, dest)
	})
}

func TestJSONHandler(t *testing.T) {
	json := fmt.Sprintf(`[ 
	{
		"path": "%s",
		"url": "%s"
	}
]`, path, dest)

	t.Run("it uses the fallback for unknown routes", func(t *testing.T) {
		result := runReaderHandler(JSONHandler, json, "/unknown")

		assertBody(t, result, fallbackResponse)
	})

	t.Run("it redirects for found url", func(t *testing.T) {
		result := runReaderHandler(JSONHandler, json, path)

		assertStatus(t, result, http.StatusFound)
		assertURL(t, result, dest)
	})
}

func runMapHandler(pathToUrls map[string]string, path string) *http.Response {
	request, _ := http.NewRequest(http.MethodGet, path, nil)
	response := httptest.NewRecorder()

	mapHandler := createMapHandler(pathToUrls)
	mapHandler(response, request)

	return response.Result()
}

func createMapHandler(pathToUrls map[string]string) http.HandlerFunc {
	fallbackHandler := http.HandlerFunc(fallback)
	return MapHandler(pathToUrls, fallbackHandler)
}

type handlerFunc func(io.Reader, http.Handler) (http.HandlerFunc, error)

func runReaderHandler(handler handlerFunc, body, path string) *http.Response {
	request, _ := http.NewRequest(http.MethodGet, path, nil)
	response := httptest.NewRecorder()

	handlerF := createReaderHandler(handler, body)
	handlerF(response, request)

	return response.Result()
}

func createReaderHandler(handler handlerFunc, yaml string) http.HandlerFunc {
	yamlReader := bytes.NewBufferString(yaml)
	fallbackHandler := http.HandlerFunc(fallback)

	handlerF, err := handler(yamlReader, fallbackHandler)
	if err != nil {
		fmt.Println(err)
		panic("could not create a Reader handler")
	}

	return handlerF
}

func assertStatus(t *testing.T, resp *http.Response, want int) {
	t.Helper()
	if resp.StatusCode != want {
		t.Errorf("Expected status to be %d, got %d",
			want, resp.StatusCode)
	}
}

func assertBody(t *testing.T, resp *http.Response, want string) {
	t.Helper()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal("Could not ready response body", err)
	}

	got := string(body)
	if want != got {
		t.Errorf("Expected response body to be %s, got %s",
			want, got)
	}
}

func assertURL(t *testing.T, resp *http.Response, want string) {
	t.Helper()
	url, err := resp.Location()

	if err != nil {
		t.Fatal("Could not read location", err)
	}

	if url.String() != want {
		t.Errorf("Expected url to be %s, got %s", url, want)
	}
}
