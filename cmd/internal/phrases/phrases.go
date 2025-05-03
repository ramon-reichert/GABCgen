// Handle musical phrases composed of words.Syllable structs.
// Phrases can be typed according to the Mass part.
package phrases

import (
	"context"
	"fmt"
	"strings"

	"github.com/ramon-reichert/GABCgen/cmd/internal/words"
)

type Phrase struct {
	Raw         string            // the original phrase
	Syllables   []*words.Syllable // the syllables of the phrase
	Syllabifier words.Syllabifier // the Syllabifier to be used to syllabify the words of the phrase
}

type PhraseMelodyer interface {
	ApplyMelody() (string, error)   // Applying the Open/Closed principle from SOLID so we can always have new types of Phrases
	GetRawString() string           //TODO: Split this interface to apply Interface Segregation SOLID principle
	PutSyllables([]*words.Syllable) // Put the built Syllables back to the original typed Phrase
}

func New(raw string, typedPhrase PhraseMelodyer) *Phrase {
	return &Phrase{
		Raw: raw,

		//	PhraseType: phraseType,
	}
}

// BuildSyllabes populates a Phrase.Syllables creating Syllable structs from each word of the Phrase.
func (ph *Phrase) BuildPhraseSyllables(ctx context.Context) ([]*words.Syllable, error) {

	words := strings.Fields(ph.Raw)
	for _, v := range words {
		//TODO: verify if word is composed and divide it at the hyphen
		syllables, err := ph.classifyWordSyllables(ctx, v)
		if err != nil {
			return syllables, fmt.Errorf("building Phrase Syllables: %w ", err)
		}
		ph.Syllables = append(ph.Syllables, syllables...)
	}

	return ph.Syllables, nil
}

// classifyWordSyllables divides the syllables of a word and builds a Syllable struct from each one of them.
func (ph *Phrase) classifyWordSyllables(ctx context.Context, word string) ([]*words.Syllable, error) {

	wordMaped := words.New(word)

	wordMaped.ParseWord()

	if err := wordMaped.Syllabify(ctx, ph.Syllabifier); err != nil {
		return []*words.Syllable{}, err //TODO: handle error
	}

	wordMaped.RecomposeWord()

	return wordMaped.BuildWordSyllables(), nil
}

// joinSyllables is a helper function that joins the GABC of all Syllables in a Phrase and adds the end string to it.
func JoinSyllables(syl []*words.Syllable, end string) string {
	var result string
	for _, v := range syl {
		result = result + v.GABC
		if v.IsLast {
			result = result + " "
		}
	}

	return result + end
}
