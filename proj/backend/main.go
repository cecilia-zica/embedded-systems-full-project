//entry point: rotas + servidor :8080

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

	// Health check: sem auth, usado pelo Docker/orquestrador pra saber se subiu
	mux.HandleFunc("GET /healthz", handleHealthz)

	// Rotas do serviço de Logging
	mux.HandleFunc("POST /api/v1/logging", requireAPIKey(handlePostLogging))
	mux.HandleFunc("GET /api/v1/logging", requireAPIKey(handleGetLogging))
	mux.HandleFunc("DELETE /api/v1/logging", requireAPIKey(handleDeleteLogging))

	// Rotas do serviço de Controle
	mux.HandleFunc("GET /api/v1/controle", requireAPIKey(handleGetControle))
	mux.HandleFunc("POST /api/v1/controle", requireAPIKey(handlePostControle))

	//servidor com timeouts em vez de ListenAndServe puro
	server := &http.Server{
		Addr:              ":8080",
		Handler:           withCORS(mux),
		ReadHeaderTimeout: 5 * time.Second,  //cliente terminar headers de requisicao
		ReadTimeout:       10 * time.Second, //ler reqs
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second, //conection pode ficar aberta ate 1 min
	}

	//graceful shutdown: no SIGINT/SIGTERM, para de aceitar conexões novas e
	//deixa as requisições em voo terminarem antes de fechar o processo.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Println("Servidor rodando em :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done() //bloqueia até chegar um sinal de encerramento
	log.Println("Encerrando servidor...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown forçado: %v", err)
	}
	if db != nil {
		db.Close()
	}
	log.Println("Servidor encerrado com sucesso.")
}

//handleHealthz: 200 se o processo está de pé e o banco responde.
func handleHealthz(w http.ResponseWriter, r *http.Request) {
	if db == nil || db.Ping() != nil {
		http.Error(w, "db indisponível", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}
