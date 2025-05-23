// Handle musical phrases composed of words.Syllable structs.
// Phrases can be typed according to the Mass part.
package phrases

import (
	"context"
	"fmt"
	"strings"

	"github.com/ramon-reichert/GABCgen/cmd/internal/gabcErrors"
	"github.com/ramon-reichert/GABCgen/cmd/internal/words"
)

type Phrase struct {
	Text        string            // the text of the phrase
	Syllables   []*words.Syllable // the syllables of the phrase
	Syllabifier words.Syllabifier // the Syllabifier to be used to syllabify the words of the phrase
}

type PhraseMelodyer interface {
	ApplyMelody() (string, error) //Applying the Open/Closed principle from SOLID so we can always have new types of Phrases
}

func New(text string) *Phrase {
	return &Phrase{
		Text: text,
	}
}

// BuildSyllabes populates a Phrase.Syllables creating Syllable structs from each word of the Phrase.
func (ph *Phrase) BuildPhraseSyllables(ctx context.Context) error {

	words := strings.Fields(ph.Text)
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
func (ph *Phrase) classifyWordSyllables(ctx context.Context, word string) ([]*words.Syllable, error) {
	var wordSyllables []*words.Syllable

	wordMaped := words.New(word)

	if err := wordMaped.ParseWord(); err != gabcErrors.ErrNoLetters { // Early scape to avoid trying to syllabify a "non-letter word"

		if err := wordMaped.Syllabify(ctx, ph.Syllabifier); err != nil {
			return wordSyllables, fmt.Errorf("classifying word syllables: %w", err)
		}

	}

	wordMaped.RecomposeWord()

	wordSyllables = wordMaped.BuildWordSyllables()

	return wordSyllables, nil
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
