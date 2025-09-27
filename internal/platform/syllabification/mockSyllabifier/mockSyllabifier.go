// Package mockSyllabifier provides a mock implementation of a syllabifier for testing purposes.
package mockSyllabifier

import (
	"context"
)

type MockSyllabifier struct {
}

func NewSyllabifier() *MockSyllabifier {
	return &MockSyllabifier{}
}

// Syllabify provides a mock syllabification of given words for testing.
func (syllab MockSyllabifier) Syllabify(ctx context.Context, word string) (string, int, error) {
	var slashed string
	var tonic int // index of the tonic syllable, beginning with 1

	switch word { // "Na verdade, é digno e justo"
	case "na":
		slashed = "na"
		tonic = 1
	case "verdade":
		slashed = "ver/da/de"
		tonic = 2
	case "é":
		slashed = "é"
		tonic = 1
	case "digno":
		slashed = "dig/no"
		tonic = 1
	case "e":
		slashed = "e"
		tonic = 1
	case "justo":
		slashed = "jus/to"
		tonic = 1
	case "por":
		slashed = "por"
		tonic = 1
	case "isso":
		slashed = "is/so"
		tonic = 1
	case "cristo":
		slashed = "cris/to"
		tonic = 1
	case "senhor":
		slashed = "se/nhor"
		tonic = 2
	case "nosso":
		slashed = "nos/so"
		tonic = 1
	}

	return slashed, tonic, nil
}

func (syllab MockSyllabifier) LoadSyllables() error {
	// Mocking the loading of syllables, no actual file operations
	return nil
}

func (syllab MockSyllabifier) SaveSyllables() error {
	// Mocking the saving of syllables, no actual file operations
	return nil
}
