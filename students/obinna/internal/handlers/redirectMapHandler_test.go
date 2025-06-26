package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMapHandler(t *testing.T) {
	pathsToUrls := map[string]string{
		"/test": "https://example.com",
	}
	fallback := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	handler := MapHandler(pathsToUrls, fallback)

	tests := []struct {
		name            string
		target          string
		triggerFallback bool
	}{
		{
			name:            "Valid path redirects",
			target:          "/test",
			triggerFallback: false,
		},
		{
			name:            "Invalid path redirects",
			target:          "/invalid",
			triggerFallback: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.target, nil)
			rw := httptest.NewRecorder() // capture response

			handler.ServeHTTP(rw, req)

			if tc.triggerFallback {
				if rw.Code != http.StatusNotFound {
					t.Errorf("expected status 404 got %d", rw.Code)
				}
			} else {
				if rw.Code != http.StatusFound {
					t.Errorf("expected status 302 got %d", rw.Code)
				}

				location := rw.Header().Get("Location")
				if location != "https://example.com" {
					t.Errorf("expected redirect to got https://example.com, got %s", location)
				}
			}

		})
	}

}
