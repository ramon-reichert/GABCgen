package gabcGen

import (
	"context"
	"log"
)

type Syllabifier interface {
	Syllabify(ctx context.Context, word string) (string, int, error)
}

type Renderer interface {
	Render(ctx context.Context, composedGABC string) (string, error)
}

type GabcGenAPI struct {
	syllabifier Syllabifier
	//	renderer    Renderer
}

func NewGabcGenAPI(syllab Syllabifier) GabcGenAPI {
	return GabcGenAPI{
		syllabifier: syllab,
	}
}

type scoreFile struct {
	Url string
}

func (gen GabcGenAPI) GeneratePreface(ctx context.Context, markedText string) (scoreFile, error) {
	preface := newPreface(markedText)
	//preface := preface.New

	err := preface.DistributeTextToPhrases(ctx, gen)
	if err != nil {
		return scoreFile{}, err //TODO: handle error
	}

	//Syllable := phraseTyped.GetSyllables()

	for _, v := range preface.phrases {
		err = v.BuildPhraseSyllables(ctx, gen)
		if err != nil {
			//TODO handle error
		}
	}

	composedGABC, err := preface.ApplyGabcMelodies(ctx)
	if err != nil {
		log.Panicln("applying gabc melodies: ", err) //TODO: handle error
	}

	var score scoreFile
	//score.url, err = gen.renderer.Render(ctx, composedGABC) //go func??
	//TODO: handle error

	score.Url = composedGABC // REMOVE LATER. JUST TO ENABLE PRE TESTING!

	return score, nil
}
