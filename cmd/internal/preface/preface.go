// Project: gabcgen - GABC generator for Gregorian chant
package preface

import (
	"fmt"
	"log"
	"strings"

	"github.com/ramon-reichert/GABCgen/cmd/internal/gabcErrors"
	"github.com/ramon-reichert/GABCgen/cmd/internal/phrases" //realy needed
)

type preface struct {
	MarkedText string //each line must end with a mark: =, *, // or +.
	Phrases    []phrases.PhraseMelodyer
}

type ( // Phrase types that can occur in a Preface
	dialogue   phrases.Phrase // dialogue = whole initial dialogue (always the same); Special treatment, since it is always the same
	firsts     phrases.Phrase // firsts(of the paragraph) = intonation, reciting tone, short cadence; Must end with "="
	last       phrases.Phrase // last(of the paragraph) = reciting tone, final cadence; Must end with "//"
	mediant    phrases.Phrase // mediant = intonation, reciting tone, mediant cadence; Must end with "*"
	conclusion phrases.Phrase // conclusion = Beginning of conclusion paragraph (often "Por isso") Must end with "+"
)

// New creates a new preface struct with the marked text.
func New(markedText string) *preface { //returning a pointer because this struct is going to be modified by its methods
	return &preface{
		MarkedText: markedText,
	}
}

func (preface *preface) TypePhrases(newPhrases []*phrases.Phrase) error {

	for _, v := range newPhrases {
		typedPhrase, err := preface.newTypedPhrase(v)
		if err != nil {
			return err //TODO handle error
		}

		aFirsts, ok := typedPhrase.(firsts) //DEBUG code
		if ok {
			log.Printf("On preface.TypePhrases(): firsts: %v\n", aFirsts)
		} else {
			log.Printf("On preface.TypePhrases(): not firsts: %v\n", typedPhrase)
		} //DEBUG code

		preface.Phrases = append(preface.Phrases, typedPhrase)
	}
	return nil
}

// newTypedPhrase creates a new Phrase struct based on the given string.
func (preface *preface) newTypedPhrase(ph *phrases.Phrase) (phrases.PhraseMelodyer, error) {

	switch {
	case strings.HasSuffix(ph.Raw, "="):
		ph.Raw, _ = strings.CutSuffix(ph.Raw, "=")
		return firsts{
			Raw:       ph.Raw,
			Syllables: ph.Syllables,
		}, nil
		/*	case strings.HasSuffix(s, "*"):
				s, _ = strings.CutSuffix(s, "*")
				return mediant{Raw: s}, nil
			case strings.HasSuffix(s, "//"):
				s, _ = strings.CutSuffix(s, "//")
				return last{Raw: s}, nil
			case strings.HasSuffix(s, "+"):
				s, _ = strings.CutSuffix(s, "+")
				return conclusion{Raw: s}, nil */
	default:
		return nil, fmt.Errorf("defining Phrase type from line: %w ", gabcErrors.ErrNoMarks)
	}
}

// ApplyGabcMelodies applies the GABC melodies to each phrase in the preface and returns the composed GABC string.
func (preface *preface) ApplyGabcMelodies() (string, error) {
	var composedGABC string
	for _, ph := range preface.Phrases {
		gabcPhrase, err := ph.ApplyMelody()
		if err != nil {
			return "", fmt.Errorf("applying melody: %w ", err)
		}
		composedGABC = composedGABC + gabcPhrase
	}
	return composedGABC, nil

}

// applyMelodySwitch route to the correct applyMelody function based on the phrase type.
/*func (preface *preface) applyMelodySwitch(ph *phrases.Phrase) (string, error) {

	switch ph.PhraseTyped {
	case firsts:
		gabcPhrase, err := ph.applyMelodyFirsts()
		if err != nil {
			return "", err
		}
		return gabcPhrase, nil

	case "last":
		/*	for i := len(ph.Syllables) - 1; ph.Syllables[i].IsTonic; i-- { //reading Syllables from the end
					if ph.Syllables[i].IsLast && ph.Syllables[i].IsFirst { //it means it's an oxytone
						ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(fgf)"
					} else {
						ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(fgf)"
					}
				}

		}
	}
	return "", fmt.Errorf("Phrase type is none of the accepted ones: %v ", ph.PhraseTyped)
}*/

func (ph firsts) GetRawString() string {
	return ph.Raw
}

// applyMelody analyzes the syllables of a phrase and attaches the GABC code(note) to each one of them, following the melody rules of that specific phrase type.
func (ph firsts) ApplyMelody() (string, error) {
	log.Printf("On firsts.ApplyMelody(): firsts.Raw: %v\n len(Syllables): %v \nSyllables: %v\n", ph.Raw, len(ph.Syllables), ph.Syllables) //DEBUG code

	i := len(ph.Syllables) - 1 //reading Syllables from the end:
	if i < 0 {
		return "", fmt.Errorf("error at firsts phrase: %v: %w ", ph.Raw, gabcErrors.ErrToShort)
	}

	//last unstressed Syllables:
	for !ph.Syllables[i].IsTonic {
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(g)"
		i--
		if i < 0 {
			return "", fmt.Errorf("error at firsts phrase: %v: %w ", ph.Raw, gabcErrors.ErrToShort)
		}
	}

	//last tonic syllable:
	ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(fg)"
	i--
	if i < 0 {
		return "", fmt.Errorf("error at firsts phrase: %v: %w ", ph.Raw, gabcErrors.ErrToShort)
	}

	//syllable before the last tonic:
	ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(gf)"
	i--
	if i < 0 {
		return "", fmt.Errorf("error at firsts phrase: %v: %w ", ph.Raw, gabcErrors.ErrToShort)
	}

	//testing the exception at last unstressed reciting syllable:
	if ph.Syllables[i].IsTonic { //default case
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(h)"
		i--
		if i < 0 {
			return "", fmt.Errorf("error at firsts phrase: %v: %w ", ph.Raw, gabcErrors.ErrToShort)
		}
	} else if ph.Syllables[i-1].IsTonic && !ph.Syllables[i-1].IsLast { //exception case
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(g)"
		i--
		if i < 0 {
			return "", fmt.Errorf("error at firsts phrase: %v: %w ", ph.Raw, gabcErrors.ErrToShort)
		}
	}

	// completing reciting Syllables:
	for i > 0 {
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(h)"
		i--
		if i < 0 {
			return "", fmt.Errorf("error at firsts phrase: %v: %w ", ph.Raw, gabcErrors.ErrToShort)
		}
	}

	//first intonation syllable:
	ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(f)"

	end := "(;)"
	return phrases.JoinSyllables(ph.Syllables, end), nil
}
