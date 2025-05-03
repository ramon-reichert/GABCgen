package gabcGen

import (
	"context"
	"log"
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
	newPhrases, err := gen.distributeTextToPhrases(markedText)
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
	}DO THIS!

	composedGABC, err := preface.ApplyGabcMelodies()
	if err != nil {
		log.Panicln("applying gabc melodies: ", err) //TODO: handle error
	}

	var score scoreFile
	//score.url, err = gen.renderer.Render(ctx, composedGABC) //go func??
	//TODO: handle error

	score.Url = composedGABC // REMOVE LATER. JUST TO ENABLE PRE TESTING!

	return score, nil
}

func (gen GabcGenAPI) distributeTextToPhrases(MarkedText string) ([]*phrases.Phrase, error) {
	var newPhrases []*phrases.Phrase

	for v := range strings.Lines(MarkedText) {
		//TODO handle errors with empty lines between pharagraphs
		newPhrases = append(newPhrases, phrases.New(v))
	}

	if len(newPhrases) == 0 {
		return newPhrases, gabcErrors.ErrNoText
	}

	return newPhrases, nil
}
