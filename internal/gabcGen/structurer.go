package gabcGen

import (
	"context"
	"fmt"
	"strings"
)

type Syllable struct {
	Char    []rune
	IsTonic bool
	IsLast  bool //If it is the last syllable of a word.
	IsFirst bool //If it is the first syllable of a word. If it is an oxytone, so IsLast an Is First are true.
	GABC    string
}

type Phrase struct {
	phraseType string //Types can be:
	//  dialogue = whole initial dialogue (always the same); Special treatment, since it is always the same
	//  firsts(of the paragraph) = intonation, reciting tone, short cadence; Must end with "="
	//  mediant = intonation, reciting tone, mediant cadence; Must end with "*"
	//  last(of the paragraph) = reciting tone, final cadence; Must end with "//"
	//	conclusion = Beginning of conclusion paragraph (often "Por isso") Must end with "+"
	Syllables []Syllable
}

func (gabc GabcGen) BuildPhrase(ctx context.Context, s string) (Phrase, error) {
	var phrase Phrase

	switch {
	case strings.HasSuffix(s, "="):
		s, _ = strings.CutSuffix(s, "=")
		phrase.phraseType = "firsts"
	case strings.HasSuffix(s, "*"):
		s, _ = strings.CutSuffix(s, "*")
		phrase.phraseType = "mediant"
	case strings.HasSuffix(s, "//"):
		s, _ = strings.CutSuffix(s, "//")
		phrase.phraseType = "last"
	case strings.HasSuffix(s, "+"):
		s, _ = strings.CutSuffix(s, "+")
		phrase.phraseType = "conclusion"
	default:
		return Phrase{}, fmt.Errorf("building Phrase from sentence: %w ", ErrResponseNoMarks)
	}

	words := strings.Fields(s)
	for _, v := range words {
		//TODO: verify if word is composed and divide it at the hyphen
		syllables, err := gabc.classifyWordSyllables(ctx, v)
		if err != nil {
			return Phrase{}, fmt.Errorf("building Phrase from sentence: %w ", err)
		}
		phrase.Syllables = append(phrase.Syllables, syllables...)
	}

	return phrase, nil
}
