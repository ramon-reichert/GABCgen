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
	Raw       string            // the original phrase
	Syllables []*words.Syllable // the syllables of the phrase
	//	Syllabifier gabcGen.Syllabifier // the Syllabifier to be used to syllabify the words of the phrase
}

// BuildSyllabes populates a Phrase.Syllables creating Syllable structs from each word of the Phrase.
func (ph *Phrase) BuildPhraseSyllables(ctx context.Context) error {

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
func (ph *Phrase) classifyWordSyllables(ctx context.Context, word string) ([]*words.Syllable, error) {

	wordMaped := words.New(word)

	wordMaped.ParseWord()

	//if err := wordMaped.Syllabify(ctx, ph.Syllabifier); err != nil {
	//	return []*words.Syllable{}, err //TODO: handle error
	//}

	wordMaped.RecomposeWord()

	return wordMaped.BuildWordSyllables(), nil
}
