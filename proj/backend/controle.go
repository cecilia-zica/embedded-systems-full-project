//handlers do serviço de configuração

package main

import (
	"encoding/json"
	"net/http"
)

type Config struct {
	BPMThreshold int  `json:"bpm_threshold"`
	AlertEnabled bool `json:"alert_enabled"`
}

// handleGetControle: busca a config atual (bpm_threshold, alert_enabled) salva no banco e devolve como JSON
func handleGetControle(w http.ResponseWriter, r *http.Request) {
	var cfg Config
	var alertEnabledInt int

	//QueryRow busca 1 linha só; Scan joga cada coluna do SELECT numa variável, na ordem em que aparecem
	err := db.QueryRow(`SELECT bpm_threshold, alert_enabled FROM config WHERE id = 1`).
		Scan(&cfg.BPMThreshold, &alertEnabledInt)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "erro ao ler config")
		return
	}
	cfg.AlertEnabled = alertEnabledInt == 1 //SQLite não tem bool, guarda como 0/1

	w.Header().Set("Content-Type", "application/json")
	//Encode: struct/map Go -> JSON (escreve a resposta pro cliente)
	json.NewEncoder(w).Encode(cfg)
}

// handlePostControle: recebe uma nova config no body (JSON), valida e atualiza o banco
func handlePostControle(w http.ResponseWriter, r *http.Request) {
	var req Config

	//limita o body a 1 MB (mesmo motivo do handlePostLogging)
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	//Decode: JSON -> struct/map Go (lê o que o cliente mandou no body)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "body inválido")
		return
	}

	if req.BPMThreshold <= 0 {
		writeJSONError(w, http.StatusBadRequest, "bpm_threshold deve ser maior que 0")
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
		writeJSONError(w, http.StatusInternalServerError, "erro ao salvar config")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	//Encode: struct/map Go -> JSON (confirma que salvou)
	json.NewEncoder(w).Encode(map[string]string{"message": "updated"})
}
