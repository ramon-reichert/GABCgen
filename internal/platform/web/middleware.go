// Package web is the http layer adapter for the GABCgen service.
package web

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// corsMiddleware adds CORS headers to responses and handles preflight requests.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w, r)

		// Handle preflight here so handlers don't need to
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// setCORSHeaders sets the CORS headers for the response.
func setCORSHeaders(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	if origin == "http://localhost:5173" || origin == "https://ramon-reichert.github.io" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// timeoutMiddleware adds a timeout to the request context.
func timeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	if timeout <= 0 {
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

var (
	visitors = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

// rateLimitMiddleware adds rate limiting based on the client's IP address.
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
				http.Error(w, "rate limit exceeded. try again after 30 seconds.", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// getVisitorLimiter returns the rate limiter for the given IP address, creating one if it doesn't exist.
func getVisitorLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	limiter, ok := visitors[ip]
	if !ok {
		limiter = rate.NewLimiter(rate.Every(60*time.Second), 2) // allows 2 requests per minute
		visitors[ip] = limiter
	}

	return limiter
}
