package httpGabcGen

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/ramon-reichert/GABCgen/cmd/internal/gabcErrors"

	"golang.org/x/time/rate"
)

type ServerConfig struct {
	Port             int
	DisableRateLimit bool
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
				responseJSON(w, http.StatusTooManyRequests, errR)
				return

			}
			next.ServeHTTP(w, r)
		})
	}
}

func NewServer(config ServerConfig, h GabcHandler) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ping)
	mux.HandleFunc("/preface", h.preface)

	// Apply the middleware to all routes
	limitedMux := rateLimitMiddleware(config.DisableRateLimit)(mux)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: limitedMux,
	}
	return &server
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
