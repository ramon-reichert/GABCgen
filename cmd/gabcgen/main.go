package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ramon-reichert/GABCgen/internal/platform/syllabification/siteSyllabifier"
	"github.com/ramon-reichert/GABCgen/internal/platform/web"
	"github.com/ramon-reichert/GABCgen/internal/service"
)

func main() {
	if err := run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func run() error {
	// Initialize dependencies
	syllabifier := siteSyllabifier.NewSyllabifier("assets/syllable_databases/liturgical_syllables.json", "assets/syllable_databases/user_syllables.json", "assets/syllable_databases/not_syllabified.txt")

	if err := syllabifier.LoadSyllables(); err != nil {
		return fmt.Errorf("loading syllables db files: %w", err)
	}

	// Initialize service with dependencies
	generatorAPI := service.NewGabcGenAPI(syllabifier /*, render*/)

	// Initialize http handler with service dependency
	gabcHandler := web.NewGabcHandler(generatorAPI, time.Duration(10*time.Second))

	// Setup http routes
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", web.Ping)
	mux.HandleFunc("/preface", gabcHandler.Preface)

	// Initialize http server
	disableRate := os.Getenv("DISABLE_RATE_LIMIT") == "true"
	server := web.NewServer(web.ServerConfig{Port: 8080, DisableRateLimit: disableRate, Timeout: 10 * time.Second}, mux)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("unexpected http server error: %v", err)
		}
		log.Println("stopped serving new requests.")
	}()

	// Wait for termination signal and shutdown gracefully
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGINT, syscall.SIGTERM)
	<-stopSignal

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("HTTP shutdown error: %w", err)
	}
	log.Println("Graceful shutdown complete.")

	return nil
}
