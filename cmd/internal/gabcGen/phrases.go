package gabcGen

import (
	"context"
	"fmt"
	"strings"
)

//type PhraseMelodyer interface {
//	applyMelody() (string, error) // Applying the Open/Closed principle from SOLID so we can always have new types of Phrases
//}

type Phrase struct {
	Raw         string
	PhraseTyped PrefacePhraseType
	Syllables   []Syllable
}

// BuildSyllabes populates a Phrase.Syllables creating Syllable structs from each word of the Phrase.
func (ph *Phrase) BuildSyllables(ctx context.Context, gen GabcGenAPI) error {

	words := strings.Fields(ph.Raw)
	for _, v := range words {
		//TODO: verify if word is composed and divide it at the hyphen
		syllables, err := gen.classifyWordSyllables(ctx, v)
		if err != nil {
			return fmt.Errorf("building Phrase Syllables: %w ", err)
		}
		ph.Syllables = append(ph.Syllables, syllables...)
	}

	return nil
}
