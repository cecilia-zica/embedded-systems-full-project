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

// setupTestDB points the backend at a temporary SQLite file and runs the real
// initDB, so tests exercise the same schema as production.
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

// TestLoggingRoundTrip verifies a reading POSTed to the API reappears on GET.
func TestLoggingRoundTrip(t *testing.T) {
	setupTestDB(t)

	body := `{"bpm":72.5,"spo2":98.0,"class":0,"user_id":"tester"}`
	postRec := httptest.NewRecorder()
	postReq := httptest.NewRequest(http.MethodPost, "/api/v1/logging", strings.NewReader(body))
	handlePostLogging(postRec, postReq)

	if postRec.Code != http.StatusCreated {
		t.Fatalf("POST logging: expected 201, got %d (%s)", postRec.Code, postRec.Body)
	}
	var created map[string]int64
	if err := json.Unmarshal(postRec.Body.Bytes(), &created); err != nil {
		t.Fatalf("POST response is not valid JSON: %v", err)
	}
	if created["id"] <= 0 {
		t.Fatalf("expected id > 0, got %d", created["id"])
	}

	getRec := httptest.NewRecorder()
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/logging", nil)
	handleGetLogging(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("GET logging: expected 200, got %d", getRec.Code)
	}
	var logs []LogEntry
	if err := json.Unmarshal(getRec.Body.Bytes(), &logs); err != nil {
		t.Fatalf("GET response is not valid JSON: %v", err)
	}
	if len(logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(logs))
	}
	if logs[0].BPM != 72.5 || logs[0].SpO2 != 98.0 || logs[0].UserID != "tester" {
		t.Fatalf("stored log differs from the one sent: %+v", logs[0])
	}
}

// TestPostLoggingRejectsInvalidBody expects a malformed body to yield 400, not 500.
func TestPostLoggingRejectsInvalidBody(t *testing.T) {
	setupTestDB(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/logging", strings.NewReader("{not json"))
	handlePostLogging(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid body, got %d", rec.Code)
	}
}

// TestPostLoggingRejectsOutOfRange expects out-of-range readings to yield 400.
func TestPostLoggingRejectsOutOfRange(t *testing.T) {
	setupTestDB(t)

	cases := []string{
		`{"bpm":-5,"spo2":98}`,  // negative bpm
		`{"bpm":0,"spo2":98}`,   // zero bpm
		`{"bpm":500,"spo2":98}`, // absurd bpm
		`{"bpm":72,"spo2":150}`, // spo2 > 100
		`{"bpm":72,"spo2":-1}`,  // negative spo2
	}
	for _, body := range cases {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/logging", strings.NewReader(body))
		handlePostLogging(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("body %s: expected 400, got %d", body, rec.Code)
		}
	}
}

// TestConfigRoundTrip verifies what POST /controle saves is read back on GET.
func TestConfigRoundTrip(t *testing.T) {
	setupTestDB(t)

	body := `{"bpm_threshold":120,"alert_enabled":true}`
	postRec := httptest.NewRecorder()
	postReq := httptest.NewRequest(http.MethodPost, "/api/v1/controle", strings.NewReader(body))
	handlePostControle(postRec, postReq)
	if postRec.Code != http.StatusOK {
		t.Fatalf("POST controle: expected 200, got %d (%s)", postRec.Code, postRec.Body)
	}

	getRec := httptest.NewRecorder()
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/controle", nil)
	handleGetControle(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("GET controle: expected 200, got %d", getRec.Code)
	}
	var cfg Config
	if err := json.Unmarshal(getRec.Body.Bytes(), &cfg); err != nil {
		t.Fatalf("GET response is not valid JSON: %v", err)
	}
	if cfg.BPMThreshold != 120 || !cfg.AlertEnabled {
		t.Fatalf("saved config differs: %+v", cfg)
	}
}

// TestPostControleRejectsNonPositiveThreshold expects threshold <= 0 to yield 400.
func TestPostControleRejectsNonPositiveThreshold(t *testing.T) {
	setupTestDB(t)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/controle", strings.NewReader(`{"bpm_threshold":0}`))
	handlePostControle(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for threshold 0, got %d", rec.Code)
	}
}

// TestIPRateLimiter: within burst it passes, once exhausted it blocks, and each
// IP has an independent bucket.
func TestIPRateLimiter(t *testing.T) {
	l := newIPRateLimiter(0.0001, 3) // negligible refill, burst 3
	ip := "10.0.0.1"

	for i := 1; i <= 3; i++ {
		if !l.allow(ip) {
			t.Fatalf("request %d should pass (within burst)", i)
		}
	}
	if l.allow(ip) {
		t.Fatal("4th request from the same IP should be blocked (burst exhausted)")
	}
	if !l.allow("10.0.0.2") {
		t.Fatal("a different IP should have an independent bucket")
	}
}

// TestClientIP: X-Forwarded-For takes precedence; otherwise the RemoteAddr host.
func TestClientIP(t *testing.T) {
	r1 := httptest.NewRequest(http.MethodGet, "/", nil)
	r1.RemoteAddr = "203.0.113.9:5555"
	if got := clientIP(r1); got != "203.0.113.9" {
		t.Fatalf("RemoteAddr: expected 203.0.113.9, got %q", got)
	}

	r2 := httptest.NewRequest(http.MethodGet, "/", nil)
	r2.Header.Set("X-Forwarded-For", "198.51.100.7, 10.0.0.1")
	if got := clientIP(r2); got != "198.51.100.7" {
		t.Fatalf("XFF: expected 198.51.100.7, got %q", got)
	}
}

// TestMain ensures tests don't leak env state to other packages.
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
