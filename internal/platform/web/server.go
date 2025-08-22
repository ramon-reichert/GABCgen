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
	h := corsMiddleware(handler)
	h = timeoutMiddleware(cfg.Timeout)(h)
	h = rateLimitMiddleware(cfg.DisableRateLimit)(h)

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: h,
	}
}
