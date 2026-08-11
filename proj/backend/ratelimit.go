// Simple per-IP rate limiting for the write endpoints: an in-memory token
// bucket that guards against POST/DELETE floods, since the API key is public
// in the clients.

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

// ipRateLimiter keeps one token bucket per client IP.
type ipRateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*tokenBucket
	rps     float64 // tokens refilled per second
	burst   float64 // maximum bucket capacity
}

func newIPRateLimiter(rps, burst float64) *ipRateLimiter {
	return &ipRateLimiter{
		buckets: make(map[string]*tokenBucket),
		rps:     rps,
		burst:   burst,
	}
}

// allow spends one token from the IP's bucket; it returns false when the limit
// is exceeded.
func (l *ipRateLimiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	b := l.buckets[ip]
	if b == nil {
		b = &tokenBucket{tokens: l.burst, last: now}
		l.buckets[ip] = b
	}
	// refill proportionally to the elapsed time, capped at burst
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

// clientIP prefers the first X-Forwarded-For hop (behind Fly's proxy) and falls
// back to the RemoteAddr host.
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

// writeLimiter is shared by the write routes: 2 req/s per IP, burst of 10.
// Generous for real use (the device posts every ~5s) but enough to cut floods.
var writeLimiter = newIPRateLimiter(2, 10)

// rateLimit responds 429 when the IP exceeds its limit, otherwise calls next.
func rateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !writeLimiter.allow(clientIP(r)) {
			writeJSONError(w, http.StatusTooManyRequests, "too many requests, try again shortly")
			return
		}
		next(w, r)
	}
}
