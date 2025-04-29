package gabcGen

import (
	"context"
	"fmt"
	"strings"
	"unicode"
)

type wordMap struct {
	word         string
	justLetters  []rune
	upperLetters map[int]rune
	notLetters   map[int]rune
}

// createWordMap takes a word and returns a wordMap struct with the word, its letters, upper case letters, and non-letter characters.
// It separates letters from non-letters and stores the original case of the letters.
func createWordMap(word string) wordMap {

	wMap := wordMap{
		word:         word,
		justLetters:  []rune{},           //store only the letters of the word, all in lower case
		upperLetters: make(map[int]rune), //store the original upper case letters
		notLetters:   make(map[int]rune), //store attached punctuation marks and non-letters as runes
	}

	runeWord := []rune(word)
	for i, v := range runeWord {
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
	return wMap
}

// recomposeWord takes a word with slashes and recomposes it with the original case and punctuation marks.
func recomposeWord(runeSlashed []rune, wordMap wordMap) string {
	var recomposedWord []rune

	entryWordIndex := 0
	runeHyphenatedIndex := 0
	for entryWordIndex < len([]rune(wordMap.word)) {
		//test if there is a ponctuation mark to put back into place:
		elem, ok := wordMap.notLetters[entryWordIndex]
		if ok {
			recomposedWord = append(recomposedWord, elem)

		} else {
			//retrieve the case of each letter:
			wasUpper, ok := wordMap.upperLetters[entryWordIndex]
			if ok {
				runeSlashed[runeHyphenatedIndex] = unicode.ToUpper(wasUpper)
			}

			recomposedWord = append(recomposedWord, runeSlashed[runeHyphenatedIndex])

			runeHyphenatedIndex++
			if runeHyphenatedIndex < len(runeSlashed) && runeSlashed[runeHyphenatedIndex] == '/' {
				recomposedWord = append(recomposedWord, runeSlashed[runeHyphenatedIndex])
				runeHyphenatedIndex++
			}

		}
		entryWordIndex++
	}
	return string(recomposedWord)
}

// classifyWordSyllables takes a word and returns its syllables with metadata.
func (gabc GabcGenAPI) classifyWordSyllables(ctx context.Context, word string) ([]Syllable, error) {
	var syllables []Syllable

	wordMap := createWordMap(word)

	slashed, tonicIndex, err := gabc.syllabifier.Syllabify(ctx, string(wordMap.justLetters))
	if err != nil {
		return syllables, fmt.Errorf("classifying syllables from word %v: %w ", word, err)
	}

	recomposedWord := recomposeWord([]rune(slashed), wordMap)

	strSyllables := strings.Split(recomposedWord, "/") //Using "/" instead of "-" to preserve syllables that use "-" to start speech

	//build a gabcGen.Syllable with metadata from each []rune syllable:
	for i, v := range strSyllables {
		runeSyllable := []rune(v)

		s := Syllable{Char: runeSyllable}
		if i+1 == tonicIndex {
			s.IsTonic = true
		}

		if i == 0 { //the first syllable
			s.IsFirst = true
		}
		if i == len(strSyllables)-1 { //the last syllable
			s.IsLast = true
		}
		syllables = append(syllables, s)
	}

	return syllables, nil
}
