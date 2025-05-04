// Maps Words and handle Syllable structs that compose musical Phrases.
package words

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	"github.com/ramon-reichert/GABCgen/cmd/internal/gabcErrors"
)

type Syllabifier interface {
	Syllabify(ctx context.Context, word string) (string, int, error)
}

type Syllable struct {
	Char    []rune
	IsTonic bool
	IsLast  bool //If it is the last syllable of a word.
	IsFirst bool //If it is the first syllable of a word. If it is an oxytone, so IsLast an Is First are true.
	GABC    string
}
type WordMaped struct {
	word              string       //the original word
	originalRunes     []rune       //the original word as runes
	justLetters       []rune       //store only the letters of the word, all in lower case
	upperLetters      map[int]rune //store the original upper case letters
	notLetters        map[int]rune //store attached punctuation marks and non-letters as runes
	slashedLetters    string       //the word with slashes between the syllables
	tonicIndex        int          //the index of the tonic syllable in the word starting from 1
	splittedSyllables []string     //the syllables of the word as slices of strings
}

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

// recomposeWord takes a word with slashes and recomposes it with the original case and punctuation marks.
func (wMap *WordMaped) RecomposeWord() {
	var recomposedWord []rune
	runeSlashed := []rune(wMap.slashedLetters)

	originalWordIndex := 0
	slashedLettersIndex := 0
	for originalWordIndex < len(wMap.originalRunes) {
		//test if there is a ponctuation mark to put back into place:
		elem, ok := wMap.notLetters[originalWordIndex]
		if ok {
			recomposedWord = append(recomposedWord, elem)

		} else {
			//retrieve the case of each letter:
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

	wMap.splittedSyllables = strings.Split(string(recomposedWord), "/") //Using "/" instead of "-" to preserve syllables that use "-" to start speech
}

// buildWordSyllables builds a Syllable struct with metadata from each []rune representing a syllable:
func (wMap *WordMaped) BuildWordSyllables() []*Syllable {

	var wordSyllables []*Syllable

	for i, v := range wMap.splittedSyllables {

		s := &Syllable{Char: []rune(v)}
		if i+1 == wMap.tonicIndex {
			s.IsTonic = true
		}

		if i == 0 { //the first syllable
			s.IsFirst = true
		}
		if i == len(wMap.splittedSyllables)-1 { //the last syllable
			s.IsLast = true
		}
		wordSyllables = append(wordSyllables, s)
	}

	return wordSyllables
}
