package urlshort

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMapHandler(t *testing.T) {
	ohandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<html><body>Hello World!</body></html>")
	})

	//Default Case Test
	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	handler := MapHandler(pathsToUrls, ohandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `<html><body>Hello World!</body></html>`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	//Check if Map actually mapped
	req, err = http.NewRequest("GET", "/yaml-godoc", nil)
	handler = MapHandler(pathsToUrls, ohandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	notExpected := `<html><body>Hello World!</body></html>`
	if rr.Body.String() == notExpected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}

func TestYAMLHandler(t *testing.T) {
	ohandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<html><body>Hello World!</body></html>")
	})

	//Default Case Test
	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	handler, err := YAMLHandler([]byte(yaml), ohandler)
	if err != nil {
		t.Errorf("YAML parsing Failed")
	}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `<html><body>Hello World!</body></html>`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	//Check if YAML actually mapped
	req, err = http.NewRequest("GET", "/yaml-godoc", nil)
	handler, err = YAMLHandler([]byte(yaml), ohandler)
	if err != nil {
		t.Errorf("YAML parsing Failed")
	}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	notExpected := `<html><body>Hello World!</body></html>`
	if rr.Body.String() == notExpected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	//Check if YAML parsing Failed
	yaml = `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
	 path: /urlshort-final
`
	req, err = http.NewRequest("GET", "/yaml-godoc", nil)
	handler, err = YAMLHandler([]byte(yaml), ohandler)
	if err == nil {
		t.Errorf("YAML parsing Failed")
	}

}
