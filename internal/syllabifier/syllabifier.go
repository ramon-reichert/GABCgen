package syllabifier

import (
	"context"
	"fmt"
	"log"
	"strings"
	"unicode"

	"github.com/ramon-reichert/GABCgen/internal/definitions"
)

func mockSyllabify(ctx context.Context, word string) (string, int, error) {
	var hyphen string
	var tonic int

	//Mocking syllabification to allow testing the core application:
	switch word { //"Na verdade, é digno e justo,="
	case "na":
		hyphen = "na"
		tonic = 1
	case "verdade":
		hyphen = "ver/da/de"
		tonic = 2
	case "é":
		hyphen = "é"
		tonic = 1
	case "digno":
		hyphen = "dig/no"
		tonic = 1
	case "e":
		hyphen = "e"
		tonic = 1
	case "justo":
		hyphen = "jus/to"
		tonic = 1
	}

	//TODO: fetch the word in a list of already used words. Could be a map["palavra"]Syllab{hyphen: "pa-la-vra", tonic: "la"}
	//TODO: if it is not there, ask it to an AI API or another solution

	return hyphen, tonic, nil
}

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

			//DEBUG:
			log.Println(elem, " > elem rune")
			log.Println(recomposedWord, " > recomposed word runes")
			log.Println(string(recomposedWord), " > recomposed word string")

		} else {
			//retrieve the case of each letter:
			wasUpper, ok := wordMap.upperLetters[entryWordIndex]
			if ok {
				runeSlashed[runeHyphenatedIndex] = unicode.ToUpper(wasUpper)
			}

			recomposedWord = append(recomposedWord, runeSlashed[runeHyphenatedIndex])
			log.Println("length of runeHyphenated: ", len(runeSlashed))
			log.Println("runeHyphenatedIndex: ", runeHyphenatedIndex)
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

// ClassifyWordSyllables takes a word and returns its syllables with metadata.
func ClassifyWordSyllables(ctx context.Context, word string) ([]definitions.Syllable, error) {
	var syllables []definitions.Syllable

	wordMap := createWordMap(word)

	slashed, tonicIndex, err := mockSyllabify(ctx, string(wordMap.justLetters))
	if err != nil {
		return syllables, fmt.Errorf("classifying syllables from word %v: %w ", word, err)
	}

	recomposedWord := recomposeWord([]rune(slashed), wordMap)

	strSyllables := strings.Split(recomposedWord, "/") //Using "/" instead of "-" to preserve syllables that use "-" to start speech

	//build a definitions.Syllable with metadata from each []rune syllable:
	for i, v := range strSyllables {
		runeSyllable := []rune(v)

		s := definitions.Syllable{Char: runeSyllable}
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

func BuildPhrase(ctx context.Context, s string) (definitions.Phrase, error) {
	var phrase definitions.Phrase

	switch {
	case strings.HasSuffix(s, "="):
		s, _ = strings.CutSuffix(s, "=")
		phrase.PhraseType = "firsts"
	case strings.HasSuffix(s, "*"):
		s, _ = strings.CutSuffix(s, "*")
		phrase.PhraseType = "mediant"
	case strings.HasSuffix(s, "//"):
		s, _ = strings.CutSuffix(s, "//")
		phrase.PhraseType = "last"
	case strings.HasSuffix(s, "+"):
		s, _ = strings.CutSuffix(s, "+")
		phrase.PhraseType = "conclusion"
	default:
		return definitions.Phrase{}, fmt.Errorf("building Phrase from sentence: %w ", definitions.ErrResponseNoMarks)
	}

	words := strings.Fields(s)
	for _, v := range words {
		//TODO: verify if word is composed and divide it at the hyphen
		syllables, err := ClassifyWordSyllables(ctx, v)
		if err != nil {
			return definitions.Phrase{}, fmt.Errorf("building Phrase from sentence: %w ", err)
		}
		phrase.Syllables = append(phrase.Syllables, syllables...)
	}

	return phrase, nil
}
