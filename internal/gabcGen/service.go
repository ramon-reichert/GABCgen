package gabcGen

type GabcGen struct {
	syllabifier Syllabifier
	//renderer             Renderer
}

func NewGabcGen(syllab Syllabifier) GabcGen {
	return GabcGen{
		syllabifier: syllab,
	}
}
