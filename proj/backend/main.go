// Command backend is the monitoring API: a small HTTP service backed by SQLite
// that stores device readings and serves the alert configuration.
package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	initDB()

	mux := http.NewServeMux()

	// Health check: unauthenticated, used by Docker/orchestrators for readiness.
	mux.HandleFunc("GET /healthz", handleHealthz)

	// Logging service (writes go through the per-IP rate limiter).
	mux.HandleFunc("POST /api/v1/logging", rateLimit(requireAPIKey(handlePostLogging)))
	mux.HandleFunc("GET /api/v1/logging", requireAPIKey(handleGetLogging))
	mux.HandleFunc("DELETE /api/v1/logging", rateLimit(requireAPIKey(handleDeleteLogging)))

	// Config service.
	mux.HandleFunc("GET /api/v1/controle", requireAPIKey(handleGetControle))
	mux.HandleFunc("POST /api/v1/controle", rateLimit(requireAPIKey(handlePostControle)))

	// Explicit timeouts instead of a bare ListenAndServe.
	server := &http.Server{
		Addr:              ":8080",
		Handler:           withCORS(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Graceful shutdown: on SIGINT/SIGTERM, stop accepting new connections and
	// let in-flight requests drain before exiting.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Println("server listening on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done() // block until a termination signal arrives
	log.Println("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("forced shutdown: %v", err)
	}
	if db != nil {
		db.Close()
	}
	log.Println("server shut down cleanly")
}

// handleHealthz returns 200 when the process is up and the database responds.
func handleHealthz(w http.ResponseWriter, r *http.Request) {
	if db == nil || db.Ping() != nil {
		http.Error(w, "database unavailable", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}
