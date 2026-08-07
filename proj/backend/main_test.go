package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// setupTestDB aponta o backend pra um SQLite temporário e roda o initDB real,
// então os testes exercitam o mesmo schema/migrations da produção.
func setupTestDB(t *testing.T) {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test.db")
	t.Setenv("DB_PATH", dbPath)
	initDB()
	t.Cleanup(func() {
		if db != nil {
			db.Close()
		}
	})
}

// TestLoggingRoundTrip: uma leitura enviada por POST tem que reaparecer no GET.
func TestLoggingRoundTrip(t *testing.T) {
	setupTestDB(t)

	body := `{"bpm":72.5,"spo2":98.0,"class":0,"user_id":"tester"}`
	postRec := httptest.NewRecorder()
	postReq := httptest.NewRequest(http.MethodPost, "/api/v1/logging", strings.NewReader(body))
	handlePostLogging(postRec, postReq)

	if postRec.Code != http.StatusCreated {
		t.Fatalf("POST logging: esperava 201, veio %d (%s)", postRec.Code, postRec.Body)
	}
	var created map[string]int64
	if err := json.Unmarshal(postRec.Body.Bytes(), &created); err != nil {
		t.Fatalf("resposta do POST não é JSON válido: %v", err)
	}
	if created["id"] <= 0 {
		t.Fatalf("esperava id > 0, veio %d", created["id"])
	}

	getRec := httptest.NewRecorder()
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/logging", nil)
	handleGetLogging(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("GET logging: esperava 200, veio %d", getRec.Code)
	}
	var logs []LogEntry
	if err := json.Unmarshal(getRec.Body.Bytes(), &logs); err != nil {
		t.Fatalf("resposta do GET não é JSON válido: %v", err)
	}
	if len(logs) != 1 {
		t.Fatalf("esperava 1 log, veio %d", len(logs))
	}
	if logs[0].BPM != 72.5 || logs[0].SpO2 != 98.0 || logs[0].UserID != "tester" {
		t.Fatalf("log gravado diferente do enviado: %+v", logs[0])
	}
}

// TestPostLoggingRejectsInvalidBody: body malformado tem que dar 400, não 500.
func TestPostLoggingRejectsInvalidBody(t *testing.T) {
	setupTestDB(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/logging", strings.NewReader("{nao é json"))
	handlePostLogging(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperava 400 pra body inválido, veio %d", rec.Code)
	}
}

// TestPostLoggingRejectsOutOfRange: leituras fora de faixa plausível dão 400.
func TestPostLoggingRejectsOutOfRange(t *testing.T) {
	setupTestDB(t)

	cases := []string{
		`{"bpm":-5,"spo2":98}`,  // bpm negativo
		`{"bpm":0,"spo2":98}`,   // bpm zero
		`{"bpm":500,"spo2":98}`, // bpm absurdo
		`{"bpm":72,"spo2":150}`, // spo2 > 100
		`{"bpm":72,"spo2":-1}`,  // spo2 negativo
	}
	for _, body := range cases {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/logging", strings.NewReader(body))
		handlePostLogging(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("body %s: esperava 400, veio %d", body, rec.Code)
		}
	}
}

// TestConfigRoundTrip: o que o POST /controle salva tem que voltar no GET.
func TestConfigRoundTrip(t *testing.T) {
	setupTestDB(t)

	body := `{"bpm_threshold":120,"alert_enabled":true}`
	postRec := httptest.NewRecorder()
	postReq := httptest.NewRequest(http.MethodPost, "/api/v1/controle", strings.NewReader(body))
	handlePostControle(postRec, postReq)
	if postRec.Code != http.StatusOK {
		t.Fatalf("POST controle: esperava 200, veio %d (%s)", postRec.Code, postRec.Body)
	}

	getRec := httptest.NewRecorder()
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/controle", nil)
	handleGetControle(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("GET controle: esperava 200, veio %d", getRec.Code)
	}
	var cfg Config
	if err := json.Unmarshal(getRec.Body.Bytes(), &cfg); err != nil {
		t.Fatalf("resposta do GET não é JSON válido: %v", err)
	}
	if cfg.BPMThreshold != 120 || !cfg.AlertEnabled {
		t.Fatalf("config salva diferente: %+v", cfg)
	}
}

// TestPostControleRejectsNonPositiveThreshold: threshold <= 0 tem que dar 400.
func TestPostControleRejectsNonPositiveThreshold(t *testing.T) {
	setupTestDB(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/controle", strings.NewReader(`{"bpm_threshold":0}`))
	handlePostControle(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperava 400 pra threshold 0, veio %d", rec.Code)
	}
}

// TestRequireAPIKey: sem chave 401, com a chave certa passa pro handler.
func TestRequireAPIKey(t *testing.T) {
	called := false
	next := func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}
	guarded := requireAPIKey(next)

	// sem chave
	rec := httptest.NewRecorder()
	guarded(rec, httptest.NewRequest(http.MethodGet, "/api/v1/logging", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("sem chave: esperava 401, veio %d", rec.Code)
	}
	if called {
		t.Fatal("handler não deveria ter sido chamado sem chave")
	}

	// com a chave correta
	rec = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/logging", nil)
	req.Header.Set("X-API-Key", apiKey)
	guarded(rec, req)
	if rec.Code != http.StatusOK || !called {
		t.Fatalf("com chave: esperava 200 e handler chamado, veio %d called=%v", rec.Code, called)
	}
}

// TestIPRateLimiter: dentro do burst passa, estourou bloqueia, e cada IP tem
// balde independente.
func TestIPRateLimiter(t *testing.T) {
	l := newIPRateLimiter(0.0001, 3) // reposição desprezível, burst 3
	ip := "10.0.0.1"

	for i := 1; i <= 3; i++ {
		if !l.allow(ip) {
			t.Fatalf("req %d deveria passar (dentro do burst)", i)
		}
	}
	if l.allow(ip) {
		t.Fatal("4a req do mesmo IP deveria ser bloqueada (burst esgotado)")
	}
	if !l.allow("10.0.0.2") {
		t.Fatal("IP diferente deveria ter balde independente")
	}
}

// TestClientIP: X-Forwarded-For tem prioridade; senão usa o host do RemoteAddr.
func TestClientIP(t *testing.T) {
	r1 := httptest.NewRequest(http.MethodGet, "/", nil)
	r1.RemoteAddr = "203.0.113.9:5555"
	if got := clientIP(r1); got != "203.0.113.9" {
		t.Fatalf("RemoteAddr: esperava 203.0.113.9, veio %q", got)
	}

	r2 := httptest.NewRequest(http.MethodGet, "/", nil)
	r2.Header.Set("X-Forwarded-For", "198.51.100.7, 10.0.0.1")
	if got := clientIP(r2); got != "198.51.100.7" {
		t.Fatalf("XFF: esperava 198.51.100.7, veio %q", got)
	}
}

// garante que os testes não deixem lixo de env pra outros pacotes
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
