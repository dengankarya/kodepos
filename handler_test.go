package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var handlerRecords = []PostalCode{
	{Province: "Jawa Barat", Regency: "Ciamis", District: "Cisaga", Village: "Danasari", Code: 46386, Latitude: -7.327, Longitude: 108.457, Elevation: 110, Timezone: "WIB"},
	{Province: "Jawa Tengah", Regency: "Purbalingga", District: "Karangjambu", Village: "Danasari", Code: 53357, Latitude: -7.185, Longitude: 109.436, Elevation: 705, Timezone: "WIB"},
	{Province: "Jawa Barat", Regency: "Purwakarta", District: "Jatiluhur", Village: "Kembangkuning", Code: 41152, Latitude: -6.549, Longitude: 107.412, Elevation: 112, Timezone: "WIB"},
}

func newTestHandler() *Handler {
	s := NewSearcher(handlerRecords)
	d := NewNearestDetector(handlerRecords)
	return NewHandler(s, d)
}

func TestSearchEndpoint(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest("GET", "/search?q=danasari", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp APIResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 || resp.Code != "OK" {
		t.Errorf("unexpected response: %+v", resp)
	}

	data, ok := resp.Data.([]interface{})
	if !ok || len(data) == 0 {
		t.Error("expected non-empty data array")
	}
}

func TestSearchEndpoint_MissingQ(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest("GET", "/search", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestSearchEndpoint_EmptyQ(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest("GET", "/search?q=", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestSearchEndpoint_CacheHeaders(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest("GET", "/search?q=danasari", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	cc := rr.Header().Get("Cache-Control")
	if cc != "s-maxage=86400, stale-while-revalidate=604800" {
		t.Errorf("unexpected Cache-Control: %q", cc)
	}
}

func TestDetectEndpoint(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest("GET", "/detect?latitude=-6.547&longitude=107.398", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp APIResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 || resp.Code != "OK" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestDetectEndpoint_MissingParams(t *testing.T) {
	h := newTestHandler()

	tests := []struct {
		name string
		url  string
	}{
		{"missing both", "/detect"},
		{"missing longitude", "/detect?latitude=-6.547"},
		{"missing latitude", "/detect?longitude=107.398"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			rr := httptest.NewRecorder()
			h.ServeHTTP(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Errorf("expected 400, got %d", rr.Code)
			}
		})
	}
}

func TestDetectEndpoint_CacheHeaders(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest("GET", "/detect?latitude=-6.547&longitude=107.398", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	cc := rr.Header().Get("Cache-Control")
	if cc != "s-maxage=86400, stale-while-revalidate=604800" {
		t.Errorf("unexpected Cache-Control: %q", cc)
	}
}

func TestHomeRedirect_WithQ(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest("GET", "/?q=danasari", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusMovedPermanently {
		t.Errorf("expected 301, got %d", rr.Code)
	}
	loc := rr.Header().Get("Location")
	if loc != "/search/?q=danasari" {
		t.Errorf("expected redirect to /search/?q=danasari, got %q", loc)
	}
}

func TestHomeRedirect_WithoutQ(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusMovedPermanently {
		t.Errorf("expected 301, got %d", rr.Code)
	}
	loc := rr.Header().Get("Location")
	if loc != "https://github.com/dengankarya/kodepos" {
		t.Errorf("expected redirect to GitHub, got %q", loc)
	}
}

func TestNotFound(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest("GET", "/unknown", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}
