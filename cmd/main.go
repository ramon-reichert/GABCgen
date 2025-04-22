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

	"github.com/ramon-reichert/GABCgen/cmd/httpGabcGen"
	"github.com/ramon-reichert/GABCgen/cmd/internal/gabcGen"
	"github.com/ramon-reichert/GABCgen/cmd/internal/syllabification"
)

var ctx context.Context = context.Background()

func main() {
	err := run()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func run() error {
	//Init dependencies:
	syllabifier := syllabification.NewSyllabifier()

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

	incomingPhrase := "-Na: verd'ade, é .digno e justo,=" //TODO: Pass the entire Preface text to general method called BuildPreface, and it will return the entire GABC text.
	phrase, err := gabc.BuildPhrase(ctx, incomingPhrase)  //TODO: BuildPhrase should be an internal method of GABCgen.
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("incomingPhrase: ", incomingPhrase)

	composedGABC, err := phrase.ApplyMelodyGABC(ctx)
	fmt.Println("composedGABC", composedGABC)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
