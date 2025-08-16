// The core of the GABC Generator. Coordinates internal and external packages interactions.
package gabcGen

import (
	"context"
	"fmt"
	"log"

	"github.com/ramon-reichert/GABCgen/internal/paragraph"
	"github.com/ramon-reichert/GABCgen/internal/preface"
	"github.com/ramon-reichert/GABCgen/internal/words"
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

// GeneratePreface attaches GABC code to each syllable of the incomming lined text following the preface melody rules.
// Each line is a phrase with its corresponding melody. Pharagraphs are separated by a double newline.
func (gen GabcGen) GeneratePreface(ctx context.Context, p preface.Preface) (preface.Preface, error) {
	linedText := p.Text.LinedText

	newParagraphs, err := paragraph.DistributeText(linedText)
	if err != nil {
		return preface.Preface{}, fmt.Errorf("generating Preface: %w", err)
	}

	for _, p := range newParagraphs {
		for _, ph := range p.Phrases {

			if ph.ExtractDirectives() != nil {
				log.Println(err)
			}

			ph.Syllabifier = gen.Syllabifier

			if err := ph.BuildPhraseSyllables(ctx); err != nil {
				return preface.Preface{}, fmt.Errorf("generating Preface: %w", err)
			}
		}
	}

	//save the user syllables to the file at once with all new words
	err = gen.Syllabifier.SaveSyllables()
	if err != nil {
		return preface.Preface{}, fmt.Errorf("saving user syllables: %w", err)
	}

	prefaceText := preface.New(linedText)

	if err := prefaceText.TypePhrases(newParagraphs); err != nil {
		return preface.Preface{}, fmt.Errorf("generating Preface: %w", err)
	}

	if prefaceText.ApplyGabcMelodies() != nil {
		return preface.Preface{}, fmt.Errorf("generating Preface: %w", err)
	}

	p.Text = *prefaceText
	p.Gabc = p.JoinPrefaceFields()

	return p, nil
}
