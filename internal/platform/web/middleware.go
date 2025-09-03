package web

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// ---- CORS ----

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w, r)

		// Handle preflight here so handlers don't need to.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func setCORSHeaders(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	// Keep the two origins your tests & current code allow.
	if origin == "http://localhost:5173" || origin == "https://ramon-reichert.github.io" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// ---- Timeout ----

func timeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	if timeout <= 0 {
		// No-op wrapper if timeout is disabled/zero
		return func(next http.Handler) http.Handler { return next }
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ---- Rate limiting (per-remote IP) ----

var (
	visitors = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

func rateLimitMiddleware(disable bool) func(http.Handler) http.Handler {
	if disable {
		return func(next http.Handler) http.Handler { return next }
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Skip rate limiting for preflight requests
			if r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			limiter := getVisitorLimiter(ip)
			if !limiter.Allow() {
				setCORSHeaders(w, r)
				http.Error(w, "rate limit exceeded. try again after 30 seconds.", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func getVisitorLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	limiter, ok := visitors[ip]
	if !ok {
		// Matches your current behavior roughly: 1 req/sec, burst 3.
		limiter = rate.NewLimiter(rate.Every(60*time.Second), 2)
		visitors[ip] = limiter
	}
	return limiter
}
