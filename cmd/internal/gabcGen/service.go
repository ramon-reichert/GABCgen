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

func (gen GabcGen) GeneratePreface(ctx context.Context, markedText string) (scoreFile, error) {
	var score scoreFile

	marks := "=+*$" //Possible preface marks

	newPhrases, err := gen.distributeTextToPhrases(markedText, marks)
	if err != nil {
		return score, fmt.Errorf("generating Preface: %w", err)
	}

	for _, v := range newPhrases {
		v.Syllabifier = gen.Syllabifier

		if err := v.BuildPhraseSyllables(ctx); err != nil {
			return score, fmt.Errorf("generating Preface: %w", err)
		}
	}

	preface := preface.New(markedText)

	if err := preface.TypePhrases(newPhrases); err != nil {
		return score, fmt.Errorf("generating Preface: %w", err)
	}

	composedGABC, err := preface.ApplyGabcMelodies()
	if err != nil {
		return score, fmt.Errorf("generating Preface: %w", err)
	}

	//score.url, err = gen.renderer.Render(ctx, composedGABC) //go func??
	//TODO: handle error

	score.Url = composedGABC // REMOVE LATER. JUST TO ENABLE PRE TESTING!

	return score, nil
}

func (gen GabcGen) distributeTextToPhrases(MarkedText, marks string) ([]*phrases.Phrase, error) {
	var newPhrases []*phrases.Phrase

	if MarkedText == "" {
		return newPhrases, fmt.Errorf("distributing text to new Phrases: %w", gabcErrors.ErrNoText)
	}

	for v := range strings.Lines(MarkedText) {
		//TODO handle empty lines between pharagraphs

		//Parse suffix mark:
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
