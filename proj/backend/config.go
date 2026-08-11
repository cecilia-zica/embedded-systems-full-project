// Handlers for the configuration service.

package main

import (
	"encoding/json"
	"net/http"
)

// Config is the alert configuration shared with the device.
type Config struct {
	BPMThreshold int  `json:"bpm_threshold"`
	AlertEnabled bool `json:"alert_enabled"`
}

// handleGetConfig returns the current configuration as JSON.
func handleGetConfig(w http.ResponseWriter, r *http.Request) {
	var cfg Config
	var alertEnabledInt int

	err := db.QueryRow(`SELECT bpm_threshold, alert_enabled FROM config WHERE id = 1`).
		Scan(&cfg.BPMThreshold, &alertEnabledInt)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to read config")
		return
	}
	cfg.AlertEnabled = alertEnabledInt == 1 // SQLite has no bool; stored as 0/1

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cfg)
}

// handlePostConfig validates a configuration from the request body and
// persists it.
func handlePostConfig(w http.ResponseWriter, r *http.Request) {
	var req Config

	// cap the body at 1 MB (same rationale as handlePostLogging)
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid body")
		return
	}

	if req.BPMThreshold <= 0 {
		writeJSONError(w, http.StatusBadRequest, "bpm_threshold must be greater than 0")
		return
	}

	alertEnabledInt := 0
	if req.AlertEnabled {
		alertEnabledInt = 1
	}

	_, err := db.Exec(
		`UPDATE config SET bpm_threshold = ?, alert_enabled = ?, updated_at = datetime('now') WHERE id = 1`,
		req.BPMThreshold, alertEnabledInt,
	)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to save config")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "updated"})
}
