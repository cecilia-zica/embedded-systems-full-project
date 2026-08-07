//verificação do X-API-Key

//closure: requireAPIKey recebe o handler (next) e devolve outro que "lembra" ele
//middleware roda antes do handler; só chama next(w,r) se a API key bater, senão 401 e para ali

package main

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"os"
)

// projeto pessoal: chave fixa em código, não depende de setar env var toda vez que roda.
// se um dia isso for exposto fora do localhost, trocar por uma chave forte via API_KEY.
const defaultAPIKey = "zica123"

var apiKey = getAPIKey()

func getAPIKey() string {
	if key := os.Getenv("API_KEY"); key != "" {
		return key
	}
	return defaultAPIKey
}

// se ALLOWED_ORIGIN não estiver setada, cai pra "*" (dev local, teste no Chrome).
// em produção, setar ALLOWED_ORIGIN com o domínio real do app pra travar o CORS.
var allowedOrigin = getAllowedOrigin()

func getAllowedOrigin() string {
	if origin := os.Getenv("ALLOWED_ORIGIN"); origin != "" {
		return origin
	}
	return "*"
}

// writeJSONError centraliza a resposta de erro: garante Content-Type application/json
// antes do status (senão o cliente recebe o JSON marcado como text/plain).
func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func requireAPIKey(next http.HandlerFunc) http.HandlerFunc {
	//devolve o handler de verdade que substitui "next" nas rotas
	return func(w http.ResponseWriter, r *http.Request) {
		//pega a API key do header
		apikey := r.Header.Get("X-API-KEY")

		//compara em tempo constante (evita timing attack)
		if subtle.ConstantTimeCompare([]byte(apikey), []byte(apiKey)) != 1 {
			//não bateu: 401 e para aqui, next nunca roda (Header().Set antes do WriteHeader, senão é ignorado)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
			return
		}

		//bateu: segue pro handler de verdade
		next(w, r)
	}
}

// CORS é regra do navegador, não do servidor — sem esse header o Chrome bloqueia a resposta
// ESP32 e o app nativo no celular não passam por essa regra, é 100% coisa de browser
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin) //configurável via ALLOWED_ORIGIN; default "*" pra dev local
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key")

		//OPTIONS = preflight do navegador perguntando "posso?"; responde 200 vazio e para aqui
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
