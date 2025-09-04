package web

import (
	"fmt"
	"net/http"
	"time"
)

type ServerConfig struct {
	Port             int
	DisableRateLimit bool
	Timeout          time.Duration
}

// NewServer applies generic middleware (CORS, timeout, rate limit) to the given handler.
// It does NOT register any application routes.
func NewServer(cfg ServerConfig, handler http.Handler) *http.Server {
	h := timeoutMiddleware(cfg.Timeout)(handler)
	h = rateLimitMiddleware(cfg.DisableRateLimit)(h)
	h = corsMiddleware(h)

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: h,
	}
}
