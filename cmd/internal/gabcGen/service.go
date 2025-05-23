// The core of the GABC Generator. Coordinates internal and external packages interactions.
package gabcGen

import (
	"context"
	"fmt"
	"strings"

	"github.com/ramon-reichert/GABCgen/cmd/internal/gabcErrors"
	"github.com/ramon-reichert/GABCgen/cmd/internal/phrases"
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

// GeneratePreface attaches GABC code to each syllable of the incomming marked text following the preface melody rules.
func (gen GabcGen) GeneratePreface(ctx context.Context, markedText string) (string, error) {
	var composedGABC string
	//TODO: PARSE THE INCOMING PARAGRAPH TO DEFINE PHRASE  TYPES ACORDING JUST TO THE LINES, INSTEAD OF THE MARKS
	marks := "=+*$" //Possible preface marks

	newPhrases, err := gen.distributeTextToPhrases(markedText, marks)
	if err != nil {
		return composedGABC, fmt.Errorf("generating Preface: %w", err)
	}

	for _, v := range newPhrases {
		v.Syllabifier = gen.Syllabifier

		if err := v.BuildPhraseSyllables(ctx); err != nil {
			return composedGABC, fmt.Errorf("generating Preface: %w", err)
		}
	}

	preface := preface.New(markedText)

	if err := preface.TypePhrases(newPhrases); err != nil {
		return composedGABC, fmt.Errorf("generating Preface: %w", err)
	}

	composedGABC, err = preface.ApplyGabcMelodies()
	if err != nil {
		return composedGABC, fmt.Errorf("generating Preface: %w", err)
	}

	return composedGABC, nil
}

// distributeTextToPhrases takes a marked text and stores each line in a new Phrase struct.
// The last character of each line is considered a mark and is removed from the text.
func (gen GabcGen) distributeTextToPhrases(MarkedText, marks string) ([]*phrases.Phrase, error) {
	var newPhrases []*phrases.Phrase

	if MarkedText == "" {
		return newPhrases, fmt.Errorf("distributing text to new Phrases: %w", gabcErrors.ErrNoText)
	}

	for v := range strings.Lines(MarkedText) {
		//TODO handle empty lines between pharagraphs

		//Remove suffix mark from the text and store it in another field of the Phrase struct.
		index := strings.LastIndexAny(v, marks)
		if index == -1 {
			return newPhrases, fmt.Errorf("distributing text to new Phrases: %w", gabcErrors.ErrNoMarks)
		}
		mark := string(v[index])
		text, _ := strings.CutSuffix(v, mark)

		newPhrases = append(newPhrases, phrases.New(text, mark))
	}

	if len(newPhrases) == 0 {
		return newPhrases, fmt.Errorf("distributing text to new Phrases: %w", gabcErrors.ErrNoText)
	}

	return newPhrases, nil
}

func (gen GabcGen) RenderPDF(ctx context.Context, markedText string) (scoreFile, error) {
	var score scoreFile
	//score.url, err = gen.renderer.Render(ctx, composedGABC) //go func??
	//TODO: handle error
	return score, nil
}
