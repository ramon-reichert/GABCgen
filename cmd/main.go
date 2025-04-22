package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ramon-reichert/GABCgen/internal/gabcGen"
	"github.com/ramon-reichert/GABCgen/internal/syllabification"
)

var ctx context.Context = context.Background()

func main() {
	syllabifier := syllabification.NewSyllabifier()
	gabc := gabcGen.NewGabcGen(syllabifier /*, render*/)

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
}
