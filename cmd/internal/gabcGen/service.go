package gabcGen

import (
	"context"
	"log"

	"github.com/ramon-reichert/GABCgen/cmd/internal/preface"
)

//	type MassPart interface {
//		New(string) *MassPart
//		DistributeTextToPhrases() error
//		ApplyGabcMelodies() (string, error)
//	}
type Syllabifier interface {
	Syllabify(ctx context.Context, word string) (string, int, error)
}

type Renderer interface {
	Render(ctx context.Context, composedGABC string) (string, error)
}

type GabcGenAPI struct {
	Syllabifier Syllabifier
	//	renderer    Renderer
}

func NewGabcGenAPI(syllab Syllabifier) GabcGenAPI {
	return GabcGenAPI{
		Syllabifier: syllab,
	}
}

type scoreFile struct {
	Url string
}

func (gen GabcGenAPI) GeneratePreface(ctx context.Context, markedText string) (scoreFile, error) {
	preface := preface.New(markedText)

	if err := preface.DistributeTextToPhrases(); err != nil {
		return scoreFile{}, err //TODO: handle error
	}

	//Syllable := phraseTyped.GetSyllables()

	//for _, v := range preface.Phrases {
	//	phraseValues := v.GetValues()

	//	if err := v.BuildPhraseSyllables(ctx, gen); err != nil {
	//TODO handle error
	//	}
	//}

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
