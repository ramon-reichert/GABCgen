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

	"github.com/ramon-reichert/GABCgen/internal/generator"
	"github.com/ramon-reichert/GABCgen/internal/handlers"
	"github.com/ramon-reichert/GABCgen/internal/platform/web"
	"github.com/ramon-reichert/GABCgen/internal/syllabification/siteSyllabifier"
)

func main() {
	err := run()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func run() error {
	//Init dependencies:
	syllabifier := siteSyllabifier.NewSyllabifier("assets/syllable_databases/liturgical_syllables.json", "assets/syllable_databases/user_syllables.json", "assets/syllable_databases/not_syllabified.txt")
	err := syllabifier.LoadSyllables()
	if err != nil {
		return fmt.Errorf("loading syllables db files: %w", err)
	}

	//Init service with its dependencies:
	generatorAPI := generator.NewGabcGenAPI(syllabifier /*, render*/)
	gabcHandler := handlers.NewGabcHandler(generatorAPI, time.Duration(10*time.Second))

	// router:
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", handlers.Ping)
	mux.HandleFunc("/preface", gabcHandler.Preface)

	//create and init http server:
	disableRate := os.Getenv("DISABLE_RATE_LIMIT") == "true"
	server := web.NewServer(web.ServerConfig{Port: 8080, DisableRateLimit: disableRate, Timeout: 10 * time.Second}, mux)

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("unexpected http server error: %v", err)
		}
		log.Println("stopped serving new requests.")
	}()

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
