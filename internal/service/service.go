// The core of the GABC Generator. Coordinates internal and external packages interactions.
package service

import (
	"context"
	"fmt"
	"log"

	"github.com/ramon-reichert/GABCgen/internal/service/composition/phrases"
	"github.com/ramon-reichert/GABCgen/internal/service/composition/phrases/words"
	"github.com/ramon-reichert/GABCgen/internal/service/preface"
)

type GabcGen struct {
	Syllabifier words.Syllabifier
	//	renderer    Renderer
}

func NewGabcGenAPI(syllab words.Syllabifier) GabcGen {
	return GabcGen{
		Syllabifier: syllab,
	}
}

type Service interface {
	GeneratePreface(ctx context.Context, dialogue, text string) (gabc string, err error)
}

// GeneratePreface attaches GABC code to each syllable of the incomming lined text following the preface melody rules.
// Each line is a phrase with its corresponding melody. Pharagraphs are separated by a double newline.
func (gen GabcGen) GeneratePreface(ctx context.Context, dialogue, linedText string) (string, error) {

	newParagraphs, err := phrases.DistributeText(linedText)
	if err != nil {
		return "", fmt.Errorf("generating Preface: %w", err)
	}

	for _, p := range newParagraphs {
		for _, ph := range p.Phrases {

			if ph.ExtractDirectives() != nil {
				log.Println(err)
			}

			ph.Syllabifier = gen.Syllabifier

			if err := ph.BuildPhraseSyllables(ctx); err != nil {
				return "", fmt.Errorf("generating Preface: %w", err)
			}
		}
	}

	//save the user syllables to the file at once with all new words
	err = gen.Syllabifier.SaveSyllables()
	if err != nil {
		return "", fmt.Errorf("saving user syllables: %w", err)
	}

	prefaceText := preface.New(linedText)

	if err := prefaceText.TypePhrases(newParagraphs); err != nil {
		return "", fmt.Errorf("generating Preface: %w", err)
	}

	if prefaceText.ApplyGabcMelodies() != nil {
		return "", fmt.Errorf("generating Preface: %w", err)
	}

	// join preface dialogue and generated GABC text:
	s := string(preface.SetDialogueTone(dialogue)) + "\n\n" + prefaceText.ComposedGABC
	return fmt.Sprintf(`%v`, s), nil
}
