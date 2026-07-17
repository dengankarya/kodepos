package main

import (
	"net/http"
	"strconv"
	"strings"
)

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	searcher *Searcher
	detector *NearestDetector
}

// NewHandler creates a handler with the given services.
func NewHandler(searcher *Searcher, detector *NearestDetector) *Handler {
	return &Handler{searcher: searcher, detector: detector}
}

// ServeHTTP routes requests to the appropriate handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/" && r.Method == http.MethodGet:
		h.handleHome(w, r)
	case r.URL.Path == "/search" && r.Method == http.MethodGet:
		h.handleSearch(w, r)
	case r.URL.Path == "/detect" && r.Method == http.MethodGet:
		h.handleDetect(w, r)
	default:
		respondJSON(w, http.StatusNotFound, "NOT_FOUND", "This endpoint cannot be found.")
	}
}

func (h *Handler) handleHome(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q != "" {
		target := "/search/?" + r.URL.RawQuery
		http.Redirect(w, r, target, http.StatusMovedPermanently)
		return
	}
	http.Redirect(w, r, "https://github.com/sooluh/kodepos", http.StatusMovedPermanently)
}

func (h *Handler) handleSearch(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q == "" {
		respondJSON(w, http.StatusBadRequest, "BAD_REQUEST", "The 'q' parameter is required.")
		return
	}

	query := ParseSearchQuery(q)
	results := h.searcher.Search(query, 20)

	w.Header().Set("Cache-Control", "s-maxage=86400, stale-while-revalidate=604800")
	writeJSON(w, http.StatusOK, newOK(results))
}

func (h *Handler) handleDetect(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("latitude")
	lngStr := r.URL.Query().Get("longitude")

	if latStr == "" || lngStr == "" {
		respondJSON(w, http.StatusBadRequest, "BAD_REQUEST", "The 'latitude' and 'longitude' parameters is required.")
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, "BAD_REQUEST", "The 'latitude' and 'longitude' parameters is required.")
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, "BAD_REQUEST", "The 'latitude' and 'longitude' parameters is required.")
		return
	}

	result := h.detector.Detect(lat, lng)
	if result == nil {
		respondJSON(w, http.StatusNotFound, "NOT_FOUND", "This endpoint cannot be found.")
		return
	}

	w.Header().Set("Cache-Control", "s-maxage=86400, stale-while-revalidate=604800")
	writeJSON(w, http.StatusOK, newOK(result))
}
