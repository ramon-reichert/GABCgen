package gabcGen

import (
	"context"
	"fmt"
	"strings"
)

// TODO: study how to implement this:
//type PhraseMelodyer interface {
//	applyMelody() (string, error) // Applying the Open/Closed principle from SOLID so we can always have new types of Phrases
//}

type Phrase struct {
	Raw         string
	PhraseTyped PrefacePhraseType
	Syllables   []*Syllable
	syllabifier Syllabifier
}

// BuildSyllabes populates a Phrase.Syllables creating Syllable structs from each word of the Phrase.
func (ph *Phrase) BuildPhraseSyllables(ctx context.Context, gen GabcGenAPI) error {

	words := strings.Fields(ph.Raw)
	for _, v := range words {
		//TODO: verify if word is composed and divide it at the hyphen
		syllables, err := ph.classifyWordSyllables(ctx, v)
		if err != nil {
			return fmt.Errorf("building Phrase Syllables: %w ", err)
		}
		ph.Syllables = append(ph.Syllables, syllables...)
	}

	return nil
}

// classifyWordSyllables divides the syllables of a word and builds a Syllable struct from each one of them.
func (ph *Phrase) classifyWordSyllables(ctx context.Context, word string) ([]*Syllable, error) {

	wordMaped := WordMaped{word: word}
	//wordMaped := word.New(word)

	wordMaped.parseWord()

	if err := wordMaped.syllabify(ctx, ph.syllabifier); err != nil {
		return []*Syllable{}, err //TODO: handle error
	}

	wordMaped.recomposeWord()

	return wordMaped.buildWordSyllables(), nil
}
