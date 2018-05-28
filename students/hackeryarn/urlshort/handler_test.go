package urlshort

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const fallbackResponse = "fallback"

func fallback(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, fallbackResponse)
}

func TestMapHandler(t *testing.T) {
	mapHandler := createMapHandler()

	t.Run("it uses the fallback for unknown routes", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/unknown", nil)
		response := httptest.NewRecorder()

		mapHandler(response, request)

		got := response.Body.String()

		if got != fallbackResponse {
			t.Errorf("Expected fallback response to be %s, got %v",
				fallbackResponse, got)
		}
	})
}

func createMapHandler() http.HandlerFunc {
	pathToUrls := map[string]string{}
	fallbackHandler := http.HandlerFunc(fallback)
	return MapHandler(pathToUrls, fallbackHandler)
}
