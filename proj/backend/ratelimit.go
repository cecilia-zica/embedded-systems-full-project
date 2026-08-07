//rate limiting simples por IP nas rotas de escrita: token bucket em memória.
//protege contra flood de POST/DELETE, já que a API key é pública nos clientes.

package main

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type tokenBucket struct {
	tokens float64
	last   time.Time
}

type ipRateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*tokenBucket
	rps     float64 // tokens repostos por segundo
	burst   float64 // capacidade máxima do balde
}

func newIPRateLimiter(rps, burst float64) *ipRateLimiter {
	return &ipRateLimiter{
		buckets: make(map[string]*tokenBucket),
		rps:     rps,
		burst:   burst,
	}
}

// allow debita 1 token do balde do IP; false = estourou o limite.
func (l *ipRateLimiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	b := l.buckets[ip]
	if b == nil {
		b = &tokenBucket{tokens: l.burst, last: now}
		l.buckets[ip] = b
	}
	//repõe tokens proporcional ao tempo passado, até o teto (burst)
	b.tokens += now.Sub(b.last).Seconds() * l.rps
	if b.tokens > l.burst {
		b.tokens = l.burst
	}
	b.last = now

	if b.tokens >= 1 {
		b.tokens--
		return true
	}
	return false
}

// clientIP prefere o 1o IP do X-Forwarded-For (atrás do proxy do Fly),
// senão o host do RemoteAddr.
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if i := strings.IndexByte(xff, ','); i >= 0 {
			return strings.TrimSpace(xff[:i])
		}
		return strings.TrimSpace(xff)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// limiter compartilhado pelas rotas de escrita: 2 req/s por IP, rajada de 10.
// Folgado pro uso real (ESP32 posta a cada ~5s), mas corta flood.
var writeLimiter = newIPRateLimiter(2, 10)

// rateLimit responde 429 quando o IP estoura o limite, senão segue.
func rateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !writeLimiter.allow(clientIP(r)) {
			writeJSONError(w, http.StatusTooManyRequests, "muitas requisições, tente novamente em instantes")
			return
		}
		next(w, r)
	}
}
