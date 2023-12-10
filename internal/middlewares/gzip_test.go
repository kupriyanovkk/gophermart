package middlewares

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGzipMiddleware(t *testing.T) {
	t.Run("NoCompression", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello, World!"))
		})

		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Accept-Encoding", "identity")

		rr := httptest.NewRecorder()

		Gzip(handler).ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		if rr.Body.String() != "Hello, World!" {
			t.Errorf("Expected body %s, got %s", "Hello, World!", rr.Body.String())
		}

		if rr.Header().Get("Content-Encoding") != "" {
			t.Errorf("Expected no Content-Encoding header, got %s", rr.Header().Get("Content-Encoding"))
		}
	})

	t.Run("WithCompression", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cw := NewCompressWriter(w)
			defer cw.Close()

			w.Header().Set("Content-Type", "text/plain")

			cw.WriteHeader(http.StatusOK)
			cw.Write([]byte("Hello, World!"))
		})

		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("Accept-Encoding", "gzip")

		rr := httptest.NewRecorder()

		Gzip(handler).ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		body, err := gzip.NewReader(rr.Body)
		if err != nil {
			t.Errorf("Error creating gzip reader: %v", err)
		}
		defer body.Close()

		var buf bytes.Buffer
		buf.ReadFrom(body)

		if buf.String() != "Hello, World!" {
			t.Errorf("Expected body %s, got %s", "Hello, World!", buf.String())
		}

		if rr.Header().Get("Content-Encoding") != "gzip" {
			t.Errorf("Expected Content-Encoding header to be gzip, got %s", rr.Header().Get("Content-Encoding"))
		}
	})
}
