package httpGabcGen

import (
	"fmt"
	"net/http"
)

type ServerConfig struct {
	Port int
}

func NewServer(config ServerConfig, h GabcHandler) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ping)
	mux.HandleFunc("/preface", h.preface)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: mux,
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
