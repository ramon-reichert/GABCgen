// The core of the GABC Generator. Coordinates internal and external packages interactions.
package gabcGen

import (
	"context"
	"fmt"
	"log"
	"strings"

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
// Each line is a phrase with its corresponding melody. Pharagraphs are separated by a double newline.
func (gen GabcGen) GeneratePreface(ctx context.Context, p preface.Preface) (preface.Preface, error) {
	linedText := p.Text.LinedText
	var composedGABC string

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

	composedGABC, err = prefaceText.ApplyGabcMelodies()
	if err != nil {
		return preface.Preface{}, fmt.Errorf("generating Preface: %w", err)
	}

	//Adjust the ending of the composed GABC string:
	composedGABC = strings.TrimSuffix(composedGABC, "(:)(Z)\n\n") + "(::)"

	return preface.Preface{Text: preface.PrefaceText{ComposedGABC: composedGABC}}, nil
}

func (gen GabcGen) RenderPDF(ctx context.Context, markedText string) (scoreFile, error) {
	var score scoreFile
	//score.url, err = gen.renderer.Render(ctx, composedGABC) //go func??
	//TODO: handle error
	return score, nil
}
