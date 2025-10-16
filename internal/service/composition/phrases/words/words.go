// Package words provides structures and methods to handle word syllabification and related metadata.
package words

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	gabcErrors "github.com/ramon-reichert/gabcgen/internal/platform/errors"
)

type Syllabifier interface {
	Syllabify(ctx context.Context, word string) (string, int, error)
	SyllabDb
}

type SyllabDb interface {
	LoadSyllables() error //Load the syllables from a file
	SaveSyllables() error //Save the syllables to a file
}

type Syllable struct {
	Char    []rune
	IsTonic bool
	IsLast  bool   // Is it the last syllable of a word?
	IsFirst bool   // Is it the first syllable of a word? If it is an oxytone, so IsLast AND IsFirst are true.
	GABC    string // syllable text with the GABC code attached to it
}

type WordMaped struct {
	word              string       // the original word
	originalRunes     []rune       // the original word as runes
	justLetters       []rune       // only the letters of the word, all in lower case
	upperLetters      map[int]rune // the original upper case letters
	notLetters        map[int]rune // attached punctuation marks and non-letters as runes
	slashedLetters    string       // word with slashes between the syllables
	tonicIndex        int          // index of the tonic syllable in the word starting from 1
	splittedSyllables []string     // syllables of the word as slices of strings
}

// New creates a new WordMaped instance with the provided word.
func New(word string) *WordMaped {
	return &WordMaped{word: word}
}

// ParseWord populates all possible WordMap fields before syllabifying the word.
func (wMap *WordMaped) ParseWord() error {
	wMap.originalRunes = []rune(wMap.word)
	wMap.justLetters = []rune{}
	wMap.upperLetters = make(map[int]rune)
	wMap.notLetters = make(map[int]rune)

	for i, v := range wMap.originalRunes {
		if !unicode.IsLetter(v) {
			wMap.notLetters[i] = v
		} else {
			if unicode.IsUpper(v) {
				wMap.upperLetters[i] = v
				v = unicode.ToLower(v)
			}
			wMap.justLetters = append(wMap.justLetters, v)
		}
	}

	if len(wMap.justLetters) == 0 {
		return gabcErrors.ErrNoLetters
	}

	return nil
}

// Syllabify takes a word and uses the Syllabifier to split it into syllables.
func (wMap *WordMaped) Syllabify(ctx context.Context, syllabifier Syllabifier) error {
	slashed, tonicIndex, err := syllabifier.Syllabify(ctx, string(wMap.justLetters))
	if err != nil {
		return fmt.Errorf("syllabifying word %v: %w ", wMap.word, err)
	}

	wMap.slashedLetters = slashed
	wMap.tonicIndex = tonicIndex

	return nil
}

// RecomposeWord takes a word with slashes and recomposes it with the original case and punctuation marks.
func (wMap *WordMaped) RecomposeWord() {
	var recomposedWord []rune
	runeSlashed := []rune(wMap.slashedLetters)
	originalWordIndex := 0
	slashedLettersIndex := 0

	for originalWordIndex < len(wMap.originalRunes) {

		// Test if there is a ponctuation mark to put back into place
		elem, ok := wMap.notLetters[originalWordIndex]
		if ok {
			recomposedWord = append(recomposedWord, elem)
		} else {

			// Retrieve the case of each letter
			wasUpper, ok := wMap.upperLetters[originalWordIndex]
			if ok {
				runeSlashed[slashedLettersIndex] = unicode.ToUpper(wasUpper)
			}

			recomposedWord = append(recomposedWord, runeSlashed[slashedLettersIndex])
			slashedLettersIndex++

			if slashedLettersIndex < len(runeSlashed) && runeSlashed[slashedLettersIndex] == '/' {
				recomposedWord = append(recomposedWord, runeSlashed[slashedLettersIndex])
				slashedLettersIndex++
			}
		}

		originalWordIndex++
	}

	wMap.splittedSyllables = strings.Split(string(recomposedWord), "/") // using "/" instead of "-" to preserve syllables that use "-" to start speech
}

// BuildWordSyllables builds a Syllable struct with metadata from each []rune representing a syllable
func (wMap *WordMaped) BuildWordSyllables() []*Syllable {
	var wordSyllables []*Syllable

	for i, v := range wMap.splittedSyllables {
		s := &Syllable{Char: []rune(v)}

		if i+1 == wMap.tonicIndex {
			s.IsTonic = true
		}

		if i == 0 { // the first syllable
			s.IsFirst = true
		}

		if i == len(wMap.splittedSyllables)-1 { // the last syllable
			s.IsLast = true
		}

		wordSyllables = append(wordSyllables, s)
	}

	return wordSyllables
}
