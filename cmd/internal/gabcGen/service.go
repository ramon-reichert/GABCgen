package gabcGen

import "context"

type Syllabifier interface {
	Syllabify(ctx context.Context, word string) (string, int, error)
}

type Renderer interface {
	Render(ctx context.Context, composedGABC string) (string, error)
}

type GabcGenAPI struct {
	syllabifier Syllabifier
	renderer    Renderer
}

func NewGabcGenAPI(syllab Syllabifier) GabcGenAPI {
	return GabcGenAPI{
		syllabifier: syllab,
	}
}

type scoreFile struct {
	url string
}

func (gen GabcGenAPI) GeneratePreface(ctx context.Context, markedText string) (scoreFile, error) {
	preface := newPreface(markedText)
	err := preface.StructurePhrases(ctx)
	//TODO: handle error
	composedGABC, err := preface.ApplyMelodyGABC(ctx)
	//TODO: handle error
	var score scoreFile
	score.url, err = gen.renderer.Render(ctx, composedGABC) //go func??
	//TODO: handle error

	return score, nil
}
