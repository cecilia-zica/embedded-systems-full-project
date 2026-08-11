// API-key authentication, CORS and shared HTTP helpers.

package main

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"os"
)

// defaultAPIKey is the development fallback used when API_KEY is unset. It must
// never be relied on in production; set a strong API_KEY instead.
const defaultAPIKey = "zica123"

var apiKey = getAPIKey()

// getAPIKey returns the key from the API_KEY environment variable, or the
// development default when it is unset.
func getAPIKey() string {
	if key := os.Getenv("API_KEY"); key != "" {
		return key
	}
	return defaultAPIKey
}

var allowedOrigin = getAllowedOrigin()

// getAllowedOrigin returns the allowed CORS origin from ALLOWED_ORIGIN, or "*"
// for local development. Set it to the app's domain in production.
func getAllowedOrigin() string {
	if origin := os.Getenv("ALLOWED_ORIGIN"); origin != "" {
		return origin
	}
	return "*"
}

// writeJSONError writes msg as a JSON error body, setting Content-Type before
// the status so the response is not mislabeled as text/plain.
func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// requireAPIKey wraps next and rejects requests whose X-API-Key header does not
// match, comparing in constant time to avoid a timing side channel.
func requireAPIKey(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apikey := r.Header.Get("X-API-KEY")

		if subtle.ConstantTimeCompare([]byte(apikey), []byte(apiKey)) != 1 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
			return
		}

		next(w, r)
	}
}

// withCORS adds the CORS headers browsers require and short-circuits preflight
// OPTIONS requests. Non-browser clients (the device, native apps) are unaffected.
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
