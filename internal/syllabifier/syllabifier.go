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

func ClassifyWordSyllables(ctx context.Context, word string) ([]definitions.Syllable, error) {
	var syllables []definitions.Syllable

	//TODO: handle case with hyphen beginning (-speech)

	//store attached punctuation marks and non-letters as runes:
	notLetters := make(map[int]rune)
	justLetters := []rune{}
	upperLetters := make(map[int]rune)
	runeWord := []rune(word)
	for i, v := range runeWord {
		if !unicode.IsLetter(v) {
			notLetters[i] = v
		} else {
			if unicode.IsUpper(v) { //store the upper case letters
				upperLetters[i] = v
				v = unicode.ToLower(v)
			}
			justLetters = append(justLetters, v)
		}
	}
	log.Println("\n entry word: ", word)
	log.Println("notLetters: ", notLetters)
	log.Println("upperLetters: ", upperLetters)
	log.Println("justLetters: ", string(justLetters))

	hyphenated, tonicIndex, err := mockSyllabify(ctx, string(justLetters))
	if err != nil {
		return syllables, fmt.Errorf("classifying syllables from word %v: %w ", word, err)
	}

	runeHyphenated := []rune(hyphenated)
	//restore letters case and ponctuation marks to exisiting Syllables:
	var recomposedWord []rune
	entryWordIndex := 0
	runeHyphenatedIndex := 0
	for entryWordIndex < len(runeWord) {
		//test if there is a ponctuation mark to put back into place:
		elem, ok := notLetters[entryWordIndex]
		if ok {
			recomposedWord = append(recomposedWord, elem)

			//DEBUG:
			log.Println(elem, " > elem rune")
			log.Println(recomposedWord, " > recomposed word runes")
			log.Println(string(recomposedWord), " > recomposed word string")

		} else {
			//retrieve the case of each letter:
			wasUpper, ok := upperLetters[entryWordIndex]
			if ok {
				runeHyphenated[runeHyphenatedIndex] = unicode.ToUpper(wasUpper)
			}

			recomposedWord = append(recomposedWord, runeHyphenated[runeHyphenatedIndex])
			log.Println("length of runeHyphenated: ", len(runeHyphenated))
			log.Println("runeHyphenatedIndex: ", runeHyphenatedIndex)
			runeHyphenatedIndex++
			if runeHyphenatedIndex < len(runeHyphenated) && runeHyphenated[runeHyphenatedIndex] == '/' {
				recomposedWord = append(recomposedWord, runeHyphenated[runeHyphenatedIndex])
				runeHyphenatedIndex++
			}

		}
		entryWordIndex++

		//DEBUG:
		log.Println("entryWordIndex: ", entryWordIndex)
		log.Println("runeHyphenatedIndex: ", runeHyphenatedIndex)
		log.Println(string(recomposedWord), " > recomposed word string")
		log.Println(recomposedWord, " > recomposed word runes")
	}

	strSyllables := strings.Split(string(recomposedWord), "/") //Using "/" instead of "-" to preserve syllables that use "-" to start speech

	for i, v := range strSyllables {
		runeSyllable := []rune(v)
		log.Println(v, " > string syllable")
		log.Println(runeSyllable, " > runes syllable")

		//build a Syllable with metadata from the []rune without letters and putctuation marks:
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
