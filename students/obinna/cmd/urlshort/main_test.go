package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestValidatePath(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expectError bool
	}{
		{"empty path", "", false},
		{"invalid extension", "config.txt", true},
		{"non-existent file", "nonexistent.yaml", true},
		{"valid path", getTestFilePath("valid.yaml"), false}, // assuming this exists
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validatePath(tc.path)
			if tc.expectError && err == nil {
				t.Errorf("expected error for path %q but got nil", tc.path)
			}
			if !tc.expectError && err != nil {
				t.Errorf("did not expect error for path %q but got: %v", tc.path, err)
			}
		})
	}
}

func TestSetupHandler(t *testing.T) {

	tests := []struct {
		name             string
		yamlPath         string
		requestPath      string
		expectRedirect   bool
		expectStatusCode int
		expectedLocation string
		format           string
	}{
		{
			name:             "valid yaml with matching path",
			yamlPath:         getTestFilePath("valid.yaml"),
			format:           "yaml",
			requestPath:      "/test",
			expectRedirect:   true,
			expectStatusCode: http.StatusFound,
			expectedLocation: "https://github.com/gophercises/urlshort",
		},
		{
			name:             "valid yaml but fallback route",
			yamlPath:         getTestFilePath("valid.yaml"),
			format:           "yaml",
			requestPath:      "/notfound",
			expectRedirect:   false,
			expectStatusCode: http.StatusOK,
		},
		{
			name:             "no yaml path given, fallback to map",
			yamlPath:         "",
			requestPath:      "/godoc",
			expectRedirect:   true,
			expectStatusCode: http.StatusFound,
			expectedLocation: "https://pkg.go.dev/std",
		},
		{
			name:             "no yaml path, fallback default",
			yamlPath:         "",
			requestPath:      "/notfound",
			expectRedirect:   false,
			expectStatusCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler, err := setupHandler(tc.yamlPath, tc.format)
			if err != nil {
				t.Fatalf("unexpected error from setupHandler: %v", err)
			}

			req := httptest.NewRequest(http.MethodGet, tc.requestPath, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tc.expectStatusCode {
				t.Errorf("expected status code %d, got %d", tc.expectStatusCode, w.Code)
			}

			if tc.expectRedirect {
				location := w.Header().Get("Location")
				if location != tc.expectedLocation {
					t.Errorf("expected redirect to %s, got %s", tc.expectedLocation, location)
				}
			} else {
				body := w.Body.String()
				fmt.Printf("body: %s\n", body)
				if !strings.Contains(body, "fallback") {
					t.Errorf("expected fallback handler to be hit, got body: %s", body)
				}
			}
		})
	}
}

func getTestFilePath(name string) string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..", "testdata", name)
}
