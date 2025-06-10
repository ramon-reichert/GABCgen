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

	"github.com/ramon-reichert/GABCgen/cmd/httpGabcGen"
	"github.com/ramon-reichert/GABCgen/cmd/internal/gabcGen"
	"github.com/ramon-reichert/GABCgen/cmd/internal/syllabification/siteSyllabifier"
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
	syllabifier := siteSyllabifier.NewSyllabifier("cmd/syllable_databases/liturgical_syllables.json", "cmd/syllable_databases/user_syllables.json", "cmd/syllable_databases/not_syllabified.txt")
	err := syllabifier.LoadSyllables()
	if err != nil {
		return fmt.Errorf("loading syllables db files: %w", err)
	}

	//Init service with its dependencies:
	gabc := gabcGen.NewGabcGenAPI(syllabifier /*, render*/)
	gabcHandler := httpGabcGen.NewGabcHandler(gabc)

	//create and init http server:
	server := httpGabcGen.NewServer(httpGabcGen.ServerConfig{Port: 8080}, gabcHandler)

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
