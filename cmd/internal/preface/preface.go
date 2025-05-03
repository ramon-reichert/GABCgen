// Project: gabcgen - GABC generator for Gregorian chant
package preface

import (
	"fmt"
	"log"
	"strings"

	"github.com/ramon-reichert/GABCgen/cmd/internal/gabcErrors"
	"github.com/ramon-reichert/GABCgen/cmd/internal/phrases" //realy needed
	"github.com/ramon-reichert/GABCgen/cmd/internal/words"
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

// DistributeTextToPhrases takes the marked text and distributes it into phrases based on the marks at the end of each line.
// It creates a new Phrase struct for each line and appends it to the preface's phrases slice.
func (preface *preface) DistributeTextToPhrases() /*(PhraseMelodyer,*/ error {

	for v := range strings.Lines(preface.MarkedText) {
		typedPhrase, err := preface.newTypedPhrase(v)
		if err != nil {
			return err //TODO handle error
		}

		preface.Phrases = append(preface.Phrases, typedPhrase)
	}

	return nil
}

// newTypedPhrase creates a new Phrase struct based on the given string.
func (preface *preface) newTypedPhrase(s string) (phrases.PhraseMelodyer, error) {

	switch {
	case strings.HasSuffix(s, "="):
		s, _ = strings.CutSuffix(s, "=")
		return firsts{Raw: s}, nil
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
		return nil, fmt.Errorf("defining Phrase type from line: %w ", gabcErrors.ErrResponseNoMarks)
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

func (ph firsts) PutSyllables(sylls []*words.Syllable) {
	ph.Syllables = sylls
}

// applyMelody analyzes the syllables of a phrase and attaches the GABC code(note) to each one of them, following the melody rules of that specific phrase type.
func (ph firsts) ApplyMelody() (string, error) {
	log.Printf("On firsts.ApplyMelody(): firsts.Raw: %v\n len(Syllables): %v \nSyllables: %v\n", ph.Raw, len(ph.Syllables), ph.Syllables) //DEBUG code

	i := len(ph.Syllables) - 1 //reading Syllables from the end:
	if i < 0 {
		return "", fmt.Errorf("error at firsts phrase: %v: %w ", ph.Raw, gabcErrors.ErrResponseToShort)
	}

	//last unstressed Syllables:
	for !ph.Syllables[i].IsTonic {
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(g)"
		i--
		if i < 0 {
			return "", fmt.Errorf("error at firsts phrase: %v: %w ", ph.Raw, gabcErrors.ErrResponseToShort)
		}
	}

	//last tonic syllable:
	ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(fg)"
	i--
	if i < 0 {
		return "", fmt.Errorf("error at firsts phrase: %v: %w ", ph.Raw, gabcErrors.ErrResponseToShort)
	}

	//syllable before the last tonic:
	ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(gf)"
	i--
	if i < 0 {
		return "", fmt.Errorf("error at firsts phrase: %v: %w ", ph.Raw, gabcErrors.ErrResponseToShort)
	}

	//testing the exception at last unstressed reciting syllable:
	if ph.Syllables[i].IsTonic { //default case
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(h)"
		i--
		if i < 0 {
			return "", fmt.Errorf("error at firsts phrase: %v: %w ", ph.Raw, gabcErrors.ErrResponseToShort)
		}
	} else if ph.Syllables[i-1].IsTonic && !ph.Syllables[i-1].IsLast { //exception case
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(g)"
		i--
		if i < 0 {
			return "", fmt.Errorf("error at firsts phrase: %v: %w ", ph.Raw, gabcErrors.ErrResponseToShort)
		}
	}

	// completing reciting Syllables:
	for i > 0 {
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(h)"
		i--
		if i < 0 {
			return "", fmt.Errorf("error at firsts phrase: %v: %w ", ph.Raw, gabcErrors.ErrResponseToShort)
		}
	}

	//first intonation syllable:
	ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(f)"

	end := "(;)"
	return phrases.JoinSyllables(ph.Syllables, end), nil
}
