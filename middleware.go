package main

import (
	"compress/gzip"
	"crypto/subtle"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
)

// writeJSON writes a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// respondJSON writes a standardized error response.
func respondJSON(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, APIResponse{
		StatusCode: status,
		Code:       code,
		Message:    message,
	})
}

// apiKeyMiddleware validates the X-API-KEY header using constant-time comparison.
// If key is empty, auth is skipped.
func apiKeyMiddleware(key string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if key != "" {
			got := r.Header.Get("X-API-KEY")
			if subtle.ConstantTimeCompare([]byte(got), []byte(key)) != 1 {
				respondJSON(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or missing API key.")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// gzip pool to reuse writers and their ~32KB internal buffers.
var gzipPool = sync.Pool{
	New: func() interface{} {
		w, _ := gzip.NewWriterLevel(nil, gzip.BestSpeed)
		return w
	},
}

// gzipMiddleware compresses responses when the client accepts gzip.
func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz := gzipPool.Get().(*gzip.Writer)
		gz.Reset(w)
		defer func() {
			gz.Close()
			gzipPool.Put(gz)
		}()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(&gzipResponseWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

// gzipResponseWriter wraps http.ResponseWriter to write through a gzip writer.
type gzipResponseWriter struct {
	http.ResponseWriter
	io.Writer
}

func (grw *gzipResponseWriter) Write(b []byte) (int, error) {
	return grw.Writer.Write(b)
}

// chain applies middlewares in order: first listed = outermost.
func chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
