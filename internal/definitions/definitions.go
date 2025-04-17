package definitions

type Syllable struct { //TODO: MOVE THIS TO PROPER FILE LATER!!
	Char    []rune
	IsTonic bool
	IsPunct bool //TODO: Cut it off since the punctuation handling was built in syllabifier file!!
	IsLast  bool //If it is the last syllable of a word.
	IsFirst bool //If it is the first syllable of a word. If it is an oxytone, so IsLast an Is First are true.
	GABC    string
}

type Phrase struct {
	PhraseType string //Types can be:
	//  dialogue = whole initial dialogue (always the same); Special treatment, since it is always the same
	//  firsts(of the paragraph) = intonation, reciting tone, short cadence; Must end with "="
	//  mediant = intonation, reciting tone, mediant cadence; Must end with "*"
	//  last(of the paragraph) = reciting tone, final cadence; Must end with "//"
	//	conclusion = Beginning of conclusion paragraph (often "Por isso") Must end with "+"
	Syllables []Syllable
}
