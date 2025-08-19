package httpGabcGen

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	gabcErrors "github.com/ramon-reichert/GABCgen/internal/platform/errors"

	"golang.org/x/time/rate"
)

type ServerConfig struct {
	Port             int
	DisableRateLimit bool
}

func NewServer(config ServerConfig, h GabcHandler) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ping)
	mux.HandleFunc("/preface", h.preface)

	handlerWithCors := corsMiddleware(mux)
	handlerWithTimeout := timeoutMiddleware(time.Duration(h.requestTimeout))(handlerWithCors)
	handlerWithRateLimit := rateLimitMiddleware(config.DisableRateLimit)(handlerWithTimeout)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: handlerWithRateLimit,
	}
	return &server
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setCorsHeaders(w, r)

		// Handle preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func setCorsHeaders(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "http://localhost:5173" || origin == "https://ramon-reichert.github.io" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func timeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			// Replace request context with the timed context
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

var (
	visitors = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

// Get or create a rate limiter for a specific IP
func getVisitorLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		// 1 request per 60 seconds, burst of 1
		limiter = rate.NewLimiter(rate.Every(60*time.Second), 2)
		visitors[ip] = limiter
	}
	return limiter
}

func rateLimitMiddleware(disable bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Bypass rate limit for preflight CORS requests
			if disable || r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}
			limiter := getVisitorLimiter(ip)

			if !limiter.Allow() {

				errR := gabcErrors.ErrResponse{
					Code:    gabcErrors.ErrToManyRequests.Code,
					Message: gabcErrors.ErrToManyRequests.Message,
				}
				setCorsHeaders(w, r)
				responseJSON(w, http.StatusTooManyRequests, errR)
				return

			}
			next.ServeHTTP(w, r)
		})
	}
}

/* Tests the http server connection.  */
func ping(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == http.MethodGet {
		w.Write([]byte("pong"))
		return
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
