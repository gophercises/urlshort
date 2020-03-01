package urlshort

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestYAMLHandler(t *testing.T) {
	type args struct {
		yml      []byte
		fallback http.Handler
		req      *http.Request
	}
	tests := []struct {
		name     string
		args     args
		wantResp func(*http.Request) http.ResponseWriter
		wantErr  bool
	}{
		{"success", args{yml: []byte(`
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution`), fallback: nil, req: httptest.NewRequest("GET", "http://localhost:8080/urlshort", nil)},
			func(req *http.Request) http.ResponseWriter {
				resp := httptest.NewRecorder()
				http.Redirect(resp, req, "https://github.com/gophercises/urlshort", http.StatusPermanentRedirect)
				return resp
			}, false},
		{"can't find", args{yml: []byte(`
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution`), fallback: func() http.HandlerFunc {
			return func(resp http.ResponseWriter, req *http.Request) {
				fmt.Fprintln(resp, "Hello, world!")
			}
		}(), req: httptest.NewRequest("GET", "http://localhost:8080/test", nil)},
			func(req *http.Request) http.ResponseWriter {
				resp := httptest.NewRecorder()
				fmt.Fprintln(resp, "Hello, world!")
				return resp
			}, false},
		{"error yaml", args{yml: []byte(`
path: /urlshort
  urls: https://github.com/gophercises/urlshort`), fallback: func() http.HandlerFunc {
			return func(resp http.ResponseWriter, req *http.Request) {
				fmt.Fprintln(resp, "Hello, world!")
			}
		}(), req: httptest.NewRequest("GET", "http://localhost:8080/test", nil)},
			func(req *http.Request) http.ResponseWriter {
				resp := httptest.NewRecorder()
				fmt.Fprintln(resp, "Hello, world!")
				return resp
			}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := YAMLHandler(tt.args.yml, tt.args.fallback)
			if (err != nil) != tt.wantErr {
				t.Errorf("YAMLHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}
			resp := httptest.NewRecorder()
			got.ServeHTTP(resp, tt.args.req)
			wantResp := tt.wantResp(tt.args.req)
			if !reflect.DeepEqual(resp, wantResp) {
				t.Errorf("YAMLHandler() resp = %v, wantResp %v", resp, wantResp)
			}
		})
	}
}

func TestJSONHandler(t *testing.T) {
	type args struct {
		jsons    []byte
		fallback http.Handler
		req      *http.Request
	}
	tests := []struct {
		name     string
		args     args
		wantResp func(*http.Request) http.ResponseWriter
		wantErr  bool
	}{
		{"success", args{jsons: []byte(`
[{"path": "/urlshort", "url":"https://github.com/gophercises/urlshort"}, {"path":"/urlshort-final", "url":"https://github.com/gophercises/urlshort/tree/solution"}]
`), fallback: nil, req: httptest.NewRequest("GET", "http://localhost:8080/urlshort", nil)},
			func(req *http.Request) http.ResponseWriter {
				resp := httptest.NewRecorder()
				http.Redirect(resp, req, "https://github.com/gophercises/urlshort", http.StatusPermanentRedirect)
				return resp
			}, false},
		{"can't find", args{jsons: []byte(`
		[{"path": "/urlshort", "url":"https://github.com/gophercises/urlshort"}, {"path":"/urlshort-final", "url":"https://github.com/gophercises/urlshort/tree/solution"}]
		`), fallback: func() http.HandlerFunc {
			return func(resp http.ResponseWriter, req *http.Request) {
				fmt.Fprintln(resp, "Hello, world!")
			}
		}(), req: httptest.NewRequest("GET", "http://localhost:8080/test", nil)},
			func(req *http.Request) http.ResponseWriter {
				resp := httptest.NewRecorder()
				fmt.Fprintln(resp, "Hello, world!")
				return resp
			}, false},
		{"error json", args{jsons: []byte(`
		zq124
		`), fallback: func() http.HandlerFunc {
			return func(resp http.ResponseWriter, req *http.Request) {
				fmt.Fprintln(resp, "Hello, world!")
			}
		}(), req: httptest.NewRequest("GET", "http://localhost:8080/test", nil)},
			func(req *http.Request) http.ResponseWriter {
				resp := httptest.NewRecorder()
				fmt.Fprintln(resp, "Hello, world!")
				return resp
			}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JSONHandler(tt.args.jsons, tt.args.fallback)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}
			resp := httptest.NewRecorder()
			got.ServeHTTP(resp, tt.args.req)
			wantResp := tt.wantResp(tt.args.req)
			if !reflect.DeepEqual(resp, wantResp) {
				t.Errorf("JSONHandler() resp = %v, wantResp %v", resp, wantResp)
			}
		})
	}
}
