package gabcGen

type GabcGenAPI struct {
	syllabifier Syllabifier
	//renderer             Renderer
}

func NewGabcGenAPI(syllab Syllabifier) GabcGenAPI {
	return GabcGenAPI{
		syllabifier: syllab,
	}
}
