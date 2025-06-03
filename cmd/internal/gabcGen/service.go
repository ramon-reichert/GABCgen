// The core of the GABC Generator. Coordinates internal and external packages interactions.
package gabcGen

import (
	"context"
	"fmt"

	"github.com/ramon-reichert/GABCgen/cmd/internal/paragraph"
	"github.com/ramon-reichert/GABCgen/cmd/internal/preface"
	"github.com/ramon-reichert/GABCgen/cmd/internal/words"
)

type Renderer interface {
	Render(ctx context.Context, composedGABC string) (string, error)
}

type GabcGen struct {
	Syllabifier words.Syllabifier
	//	renderer    Renderer
}

func NewGabcGenAPI(syllab words.Syllabifier) GabcGen {
	return GabcGen{
		Syllabifier: syllab,
	}
}

type scoreFile struct {
	Url string
}

// GeneratePreface attaches GABC code to each syllable of the incomming lined text following the preface melody rules.
func (gen GabcGen) GeneratePreface(ctx context.Context, linedText string) (string, error) {
	var composedGABC string

	newParagraphs, err := paragraph.DistributeText(linedText)
	if err != nil {
		return composedGABC, fmt.Errorf("generating Preface: %w", err)
	}

	for _, p := range newParagraphs {
		for _, v := range p.Phrases {
			v.Syllabifier = gen.Syllabifier

			if err := v.BuildPhraseSyllables(ctx); err != nil {
				return composedGABC, fmt.Errorf("generating Preface: %w", err)
			}
		}
	}

	//save the user syllables to the file at once with all new words
	err = gen.Syllabifier.SaveSyllables("user_syllables.json")
	if err != nil {
		return composedGABC, fmt.Errorf("saving user syllables: %w", err)
	}

	preface := preface.New(linedText)

	if err := preface.TypePhrases(newParagraphs); err != nil {
		return composedGABC, fmt.Errorf("generating Preface: %w", err)
	}

	composedGABC, err = preface.ApplyGabcMelodies()
	if err != nil {
		return composedGABC, fmt.Errorf("generating Preface: %w", err)
	}

	return composedGABC, nil
}

func (gen GabcGen) RenderPDF(ctx context.Context, markedText string) (scoreFile, error) {
	var score scoreFile
	//score.url, err = gen.renderer.Render(ctx, composedGABC) //go func??
	//TODO: handle error
	return score, nil
}
