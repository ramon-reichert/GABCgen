// The core of the GABC Generator. Coordinates internal and external packages interactions.
package gabcGen

import (
	"context"
	"fmt"
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
func (gen GabcGen) GeneratePreface(ctx context.Context, linedText string) (string, error) {
	var composedGABC string

	newParagraphs, err := paragraph.DistributeText(linedText)
	if err != nil {
		return composedGABC, fmt.Errorf("generating Preface: %w", err)
	}

	for _, p := range newParagraphs {
		for _, v := range p.Phrases {

			//TODO: Parse each line looking for singing directives (between parenthesss), which must not be syllabified nor "noteted"
			//Hold its value and position and put it back at final composedGABC as ||<i><c>directive not to be sung</c></i>||

			v.Syllabifier = gen.Syllabifier

			if err := v.BuildPhraseSyllables(ctx); err != nil {
				return composedGABC, fmt.Errorf("generating Preface: %w", err)
			}
		}
	}

	//save the user syllables to the file at once with all new words
	err = gen.Syllabifier.SaveSyllables()
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

	//Adjust the ending of the composed GABC string:
	composedGABC = strings.TrimSuffix(composedGABC, "(:)(Z)\n\n") + "(::)"

	return composedGABC, nil
}

func (gen GabcGen) RenderPDF(ctx context.Context, markedText string) (scoreFile, error) {
	var score scoreFile
	//score.url, err = gen.renderer.Render(ctx, composedGABC) //go func??
	//TODO: handle error
	return score, nil
}
