//SQLite: conexão + CREATE TABLE

package main

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite" //registra o driver
)

var db *sql.DB

func initDB() {
	var err error

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./app.db" //fallback pra rodar local, fora do Docker
	}

	db, err = sql.Open("sqlite", dbPath) //abre a conexão

	if err != nil {
		log.Fatal(err) //lança erro e encerra com 1
	}

	//teste com db.Ping()
	err = db.Ping() //verifica se a conexão está ativa, boa prática
	if err != nil {
		log.Fatal(err)
	}

	//create table confg
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

	//create table logs
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			bpm REAL NOT NULL,
			spo2 REAL NOT NULL,
			class INT,
			user_id TEXT NOT NULL default 'unknown',
			created_at DATETIME default CURRENT_TIMESTAMP
		)
	`) //cada leitura -> update novo
	if err != nil {
		log.Fatal(err)
	}

	// insert or ignore leituras
	_, err = db.Exec(`INSERT OR IGNORE INTO config (id) VALUES (1)`)
	if err != nil {
		log.Fatal(err)
	}

}
