package handlers

import (
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestBuildMap(t *testing.T) {
	expectedResult := map[string]string{
		"/test":  "https://github.com/gophercises/urlshort",
		"/godoc": "https://pkg.go.dev/std",
	}

	urlPaths := getPaths()
	result := buildMap(urlPaths)

	if !cmp.Equal(expectedResult, result) {
		t.Errorf("Result mismatch:\n%s", cmp.Diff(expectedResult, result))
	}
}

func TestNewHandlerFromFile(t *testing.T) {
	type testCase struct {
		name           string
		fileName       string
		path           string
		decoder        DecoderFunc
		expectErr      bool
		expectedErr    string
		expectLocation string
		expectStatus   int
	}

	tests := []testCase{
		{
			name:           "valid YAML file",
			fileName:       getTestFilePath("valid.yaml"),
			decoder:        YamlDecoder,
			expectErr:      false,
			path:           "/test",
			expectLocation: "https://github.com/gophercises/urlshort",
			expectStatus:   http.StatusFound,
		},
		{
			name:           "valid JSON file",
			fileName:       getTestFilePath("valid.json"),
			decoder:        JsonDecoder,
			expectErr:      false,
			path:           "/test",
			expectLocation: "https://github.com/gophercises/urlshort",
			expectStatus:   http.StatusFound,
		},
		{
			name:        "malformed YAML file",
			fileName:    getTestFilePath("malformed.yaml"),
			decoder:     YamlDecoder,
			expectErr:   true,
			expectedErr: "not allowed in this context",
		},
		{
			name:        "malformed JSON file",
			fileName:    getTestFilePath("malformed.json"),
			decoder:     JsonDecoder,
			expectErr:   true,
			expectedErr: "invalid character",
		},
		{
			name:        "nonexistent file",
			fileName:    "doesnotexist.yaml",
			decoder:     YamlDecoder,
			expectErr:   true,
			expectedErr: "no such file or directory",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fallback := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			})

			handler, err := NewHandlerFromFile(tc.fileName, tc.decoder, fallback)

			if tc.expectErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}

				if !strings.Contains(err.Error(), tc.expectedErr) {
					t.Fatalf("expected %s\n got %s\n", tc.expectedErr, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			req := httptest.NewRequest("GET", tc.path, nil)
			rw := httptest.NewRecorder()

			handler.ServeHTTP(rw, req)

			if rw.Code != tc.expectStatus {
				t.Errorf("expected status %d, got %d", tc.expectStatus, rw.Code)
			}

			if tc.expectStatus == http.StatusFound {
				loc := rw.Header().Get("Location")
				if !strings.EqualFold(loc, tc.expectLocation) {
					t.Errorf("expected redirect to %q, got %q", tc.expectLocation, loc)
				}
			}
		})
	}
}

func TestUnmarshalFile(t *testing.T) {
	tests := []struct {
		name        string
		fileName    string
		decoder     DecoderFunc
		expectErr   bool
		expectedLen int
		expectedErr string
	}{
		{
			name:        "valid YAML",
			fileName:    getTestFilePath("valid.yaml"),
			decoder:     YamlDecoder,
			expectErr:   false,
			expectedLen: 2,
		},
		{
			name:        "valid JSON",
			fileName:    getTestFilePath("valid.json"),
			decoder:     JsonDecoder,
			expectErr:   false,
			expectedLen: 1,
		},
		{
			name:        "malformed YAML",
			fileName:    getTestFilePath("malformed.yaml"),
			decoder:     YamlDecoder,
			expectErr:   true,
			expectedErr: "not allowed in this context",
		},
		{
			name:        "malformed JSON",
			fileName:    getTestFilePath("malformed.json"),
			decoder:     JsonDecoder,
			expectErr:   true,
			expectedErr: "invalid character",
		},
		{
			name:      "file not found",
			fileName:  "notfound.yaml",
			decoder:   YamlDecoder,
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var result []PathUrl
			err := UnmarshalFile(tc.fileName, &result, tc.decoder)

			if tc.expectErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}

				if !strings.Contains(err.Error(), tc.expectedErr) {
					t.Fatalf("expected %s\n got %s\n", tc.expectedErr, err.Error())
				}

			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if len(result) != tc.expectedLen {
					t.Fatalf("expected %d elements, got %d", tc.expectedLen, len(result))
				}
			}
		})
	}
}

func getTestFilePath(name string) string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..", "testdata", name)
}

func getPaths() []PathUrl {
	return []PathUrl{
		{
			Path: "/test",
			Url:  "https://github.com/gophercises/urlshort",
		},
		{
			Path: "/godoc",
			Url:  "https://pkg.go.dev/std",
		},
	}
}
