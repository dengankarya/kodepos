package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//go:embed data/kodepos.json
var kodeposData []byte

func main() {
	var records []PostalCode
	if err := json.Unmarshal(kodeposData, &records); err != nil {
		log.Fatalf("failed to parse kodepos.json: %v", err)
	}

	searcher := NewSearcher(records)
	detector := NewNearestDetector(records)
	handler := NewHandler(searcher, detector)

	apiKey := os.Getenv("API_KEY")

	mux := http.NewServeMux()
	mux.Handle("/", handler)

	srv := &http.Server{
		Addr: ":" + envOr("PORT", "3000"),
		Handler: chain(mux, gzipMiddleware, func(next http.Handler) http.Handler {
			return apiKeyMiddleware(apiKey, next)
		}),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(shutdownCtx)
	}()

	log.Printf("listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
