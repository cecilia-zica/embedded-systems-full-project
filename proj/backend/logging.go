//handlers do serviço de logging

package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type LogEntry struct {
	ID        int64   `json:"id"`
	BPM       float64 `json:"bpm"`
	SpO2      float64 `json:"spo2"`
	Class     int     `json:"class"`
	UserID    string  `json:"user_id"`
	CreatedAt string  `json:"created_at"`
}

func handlePostLogging(w http.ResponseWriter, r *http.Request) {
	var l LogEntry

	//Decode: JSON -> struct (lê a leitura que o ESP32 mandou no body)
	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		http.Error(w, "body inválido", http.StatusBadRequest)
		return
	}
	if l.UserID == "" {
		l.UserID = "unknown"
	}
	result, err := db.Exec(
		`INSERT INTO logs (bpm, spo2, class, user_id) VALUES (?, ?, ?, ?)`,
		l.BPM, l.SpO2, l.Class, l.UserID,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := result.LastInsertId()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) //201: contrato da API sucesso
	//Encode: struct/map -> JSON (devolve o id gerado pro ESP32)
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

//handleDeleteLogging: apaga todo o histórico de logs (botão "limpar" do app)
func handleDeleteLogging(w http.ResponseWriter, r *http.Request) {
	if _, err := db.Exec(`DELETE FROM logs`); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "logs apagados"})
}

func handleGetLogging(w http.ResponseWriter, r *http.Request) {
	logs := []LogEntry{}
	rows, err := db.Query(`SELECT id, bpm, spo2, class, user_id, created_at FROM logs ORDER BY created_at DESC LIMIT 50 `)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) //se der erro no Scan, devolve 500 pro cliente e sai da função
		return
	}
	defer rows.Close()
	for rows.Next() {
		var l LogEntry
		if err := rows.Scan(&l.ID, &l.BPM, &l.SpO2, &l.Class, &l.UserID, &l.CreatedAt); err != nil {
			log.Println("erro ao ler log:", err)
			continue //se er (ex: campo nulo), ignora só esse log e continua lendo os outros
		}
		logs = append(logs, l)
	}
	if err := rows.Err(); err != nil { //erro de iteração (ex: conexão caiu no meio do loop)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//resposta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
