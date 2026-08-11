// SQLite connection setup and schema creation.

package main

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite" // registers the driver
)

var db *sql.DB

// initDB opens the SQLite database, verifies the connection and creates the
// schema when missing. It fails fast (log.Fatal) on any error.
func initDB() {
	var err error

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./app.db" // local fallback when running outside Docker
	}

	// WAL plus a busy_timeout avoid "database is locked" when the device writes
	// (POST /logging every ~5s) while the app reads or deletes concurrently.
	dsn := dbPath + "?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)"
	db, err = sql.Open("sqlite", dsn)

	if err != nil {
		log.Fatal(err)
	}

	// SQLite allows a single writer; serializing connections trades lock errors
	// for a short, predictable wait, which is safe at this scale.
	db.SetMaxOpenConns(1)

	// verify the connection is alive
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// config: single row (id = 1) holding the alert settings
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS config (
			id INTEGER PRIMARY KEY default 1,
			bpm_threshold INT NOT NULL default 0,
			alert_enabled INT NOT NULL default 1,
			updated_at DATETIME default CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// logs: one row per reading reported by the device
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			bpm REAL NOT NULL,
			spo2 REAL NOT NULL,
			class INT NOT NULL DEFAULT 0,
			user_id TEXT NOT NULL default 'unknown',
			created_at DATETIME default CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// seed the single config row if it isn't there yet
	_, err = db.Exec(`INSERT OR IGNORE INTO config (id) VALUES (1)`)
	if err != nil {
		log.Fatal(err)
	}

}
