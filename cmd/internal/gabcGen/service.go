package gabcGen

import (
	"context"
	"log"

	"github.com/ramon-reichert/GABCgen/cmd/internal/phrases"
	"github.com/ramon-reichert/GABCgen/cmd/internal/preface"
	"github.com/ramon-reichert/GABCgen/cmd/internal/words"
)

//	type MassPart interface {
//		New(string) *MassPart
//		DistributeTextToPhrases() error
//		ApplyGabcMelodies() (string, error)
//	}

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
	preface := preface.New(markedText)

	if err := preface.DistributeTextToPhrases(); err != nil { //MAYBE BUILD THESE PHRASES BEFORE TYPING THEM????
		return scoreFile{}, err //TODO: handle error
	}

	//Syllable := phraseTyped.GetSyllables()

	for _, v := range preface.Phrases {
		rebuiltPhrase := &phrases.Phrase{
			Raw:         v.GetRawString(),
			Syllabifier: gen.Syllabifier,
		}
		syllabs, err := rebuiltPhrase.BuildPhraseSyllables(ctx)
		if err != nil {
			//TODO handle error
		}

		v.PutSyllabes(syllabs)

	}

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
