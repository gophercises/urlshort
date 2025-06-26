package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggingMiddleware_LogsPathAndMethod(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)      //capture logs in buffer
	defer log.SetOutput(nil) // reset log output

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := LoggingMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/log-me", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	output := buf.String()
	if !strings.Contains(output, "GET") || !strings.Contains(output, "/log-me") {
		t.Errorf("expected log to contain method and path, got: %s", output)
	}
}
