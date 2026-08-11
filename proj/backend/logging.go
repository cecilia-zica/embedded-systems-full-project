// Handlers for the logging (readings) service.

package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// LogEntry is a single sensor reading stored and served by the API.
type LogEntry struct {
	ID        int64   `json:"id"`
	BPM       float64 `json:"bpm"`
	SpO2      float64 `json:"spo2"`
	Class     int     `json:"class"`
	UserID    string  `json:"user_id"`
	CreatedAt string  `json:"created_at"`
}

// handlePostLogging stores one reading sent by the device and returns its id.
func handlePostLogging(w http.ResponseWriter, r *http.Request) {
	var l LogEntry

	// cap the body at 1 MB so a client cannot exhaust memory with a huge POST
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if l.UserID == "" {
		l.UserID = "unknown"
	}

	// reject readings outside a plausible range instead of storing garbage
	if l.BPM <= 0 || l.BPM > 300 {
		writeJSONError(w, http.StatusBadRequest, "bpm out of range (0 < bpm <= 300)")
		return
	}
	if l.SpO2 < 0 || l.SpO2 > 100 {
		writeJSONError(w, http.StatusBadRequest, "spo2 out of range (0 <= spo2 <= 100)")
		return
	}

	result, err := db.Exec(
		`INSERT INTO logs (bpm, spo2, class, user_id) VALUES (?, ?, ?, ?)`,
		l.BPM, l.SpO2, l.Class, l.UserID,
	)

	if err != nil {
		log.Println("insert log failed:", err)
		writeJSONError(w, http.StatusInternalServerError, "internal error saving reading")
		return
	}
	id, _ := result.LastInsertId()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

// handleDeleteLogging clears the entire reading history (the app's "clear"
// button). The deletion is irreversible.
func handleDeleteLogging(w http.ResponseWriter, r *http.Request) {
	if _, err := db.Exec(`DELETE FROM logs`); err != nil {
		log.Println("delete logs failed:", err)
		writeJSONError(w, http.StatusInternalServerError, "internal error deleting logs")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "logs deleted"})
}

// handleGetLogging returns the 50 most recent readings, newest first.
func handleGetLogging(w http.ResponseWriter, r *http.Request) {
	logs := []LogEntry{}
	rows, err := db.Query(`SELECT id, bpm, spo2, COALESCE(class, 0), user_id, created_at FROM logs ORDER BY created_at DESC LIMIT 50 `)
	if err != nil {
		log.Println("query logs failed:", err)
		writeJSONError(w, http.StatusInternalServerError, "internal error reading logs")
		return
	}
	defer rows.Close()
	for rows.Next() {
		var l LogEntry
		if err := rows.Scan(&l.ID, &l.BPM, &l.SpO2, &l.Class, &l.UserID, &l.CreatedAt); err != nil {
			log.Println("scan log failed:", err)
			continue // skip only the bad row and keep reading the rest
		}
		logs = append(logs, l)
	}
	if err := rows.Err(); err != nil { // iteration error (e.g. connection dropped mid-loop)
		log.Println("iterate logs failed:", err)
		writeJSONError(w, http.StatusInternalServerError, "internal error reading logs")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
