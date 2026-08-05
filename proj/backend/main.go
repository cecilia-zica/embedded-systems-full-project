//entry point: rotas + servidor :8080

package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	initDB() 

	mux := http.NewServeMux()

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
		ReadHeaderTimeout: 5 * time.Second, //cliente terminar headers de requisicao
		ReadTimeout:       10 * time.Second, //ler reqs
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second, //conection pode ficar aberta ate 1 min 
	}

	log.Println("Servidor rodando em :8080")
	log.Fatal(server.ListenAndServe())
}
