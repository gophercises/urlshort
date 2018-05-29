package urlshort

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const fallbackResponse = "fallback"

func fallback(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, fallbackResponse)
}

func TestMapHandler(t *testing.T) {
	t.Run("it uses the fallback for unknown routes", func(t *testing.T) {
		result := runMapHandler("", "")

		assertBody(t, result, fallbackResponse)
	})

	t.Run("it redirects for found url", func(t *testing.T) {
		path := "/test"
		dest := "https://test.com"

		result := runMapHandler(path, dest)

		assertStatus(t, result, http.StatusFound)
		assertURL(t, result, dest)
	})
}

func createMapHandler(pathToUrls map[string]string) http.HandlerFunc {
	fallbackHandler := http.HandlerFunc(fallback)
	return MapHandler(pathToUrls, fallbackHandler)
}

func runMapHandler(path, dest string) *http.Response {
	pathToUrls := map[string]string{path: dest}
	mapHandler := createMapHandler(pathToUrls)

	if path == "" {
		path = "/unknown"
	}
	request, _ := http.NewRequest(http.MethodGet, path, nil)
	response := httptest.NewRecorder()

	mapHandler(response, request)
	return response.Result()
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
