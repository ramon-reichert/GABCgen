package gabcGen

import (
	"context"
	"strings"

	"github.com/ramon-reichert/GABCgen/cmd/internal/gabcErrors"
	"github.com/ramon-reichert/GABCgen/cmd/internal/phrases"
	"github.com/ramon-reichert/GABCgen/cmd/internal/preface"
	"github.com/ramon-reichert/GABCgen/cmd/internal/words"
)

type Renderer interface {
	Render(ctx context.Context, composedGABC string) (string, error)
}

type GabcGenAPI struct {
	Syllabifier words.Syllabifier
	//	renderer    Renderer
}

func NewGabcGenAPI(syllab words.Syllabifier) GabcGenAPI {
	return GabcGenAPI{
		Syllabifier: syllab,
	}
}

type scoreFile struct {
	Url string
}

func (gen GabcGenAPI) GeneratePreface(ctx context.Context, markedText string) (scoreFile, error) {
	marks := "=+*$" //Possible preface marks

	newPhrases, err := gen.distributeTextToPhrases(markedText, marks)
	if err != nil {
		return scoreFile{}, err //TODO: handle error
	}

	for _, v := range newPhrases {
		v.Syllabifier = gen.Syllabifier

		if err := v.BuildPhraseSyllables(ctx); err != nil {
			//TODO handle error
		}
	}

	preface := preface.New(markedText)

	if err := preface.TypePhrases(newPhrases); err != nil {
		return scoreFile{}, err //TODO: handle error
	}

	composedGABC, err := preface.ApplyGabcMelodies()
	if err != nil {
		return scoreFile{}, err //TODO: handle error
	}

	var score scoreFile
	//score.url, err = gen.renderer.Render(ctx, composedGABC) //go func??
	//TODO: handle error

	score.Url = composedGABC // REMOVE LATER. JUST TO ENABLE PRE TESTING!

	return score, nil
}

func (gen GabcGenAPI) distributeTextToPhrases(MarkedText, marks string) ([]*phrases.Phrase, error) {
	var newPhrases []*phrases.Phrase

	for v := range strings.Lines(MarkedText) {
		//TODO handle errors with empty lines between pharagraphs

		//Parse suffix mark:
		index := strings.LastIndexAny(v, marks)
		if index == -1 {
			return newPhrases, gabcErrors.ErrNoMarks
		}
		mark := string(v[index])
		text, _ := strings.CutSuffix(v, mark)

		newPhrases = append(newPhrases, phrases.New(text, mark))
	}

	if len(newPhrases) == 0 {
		return newPhrases, gabcErrors.ErrNoText
	}

	return newPhrases, nil
}
