package gabcGen

import (
	"context"
	"fmt"
	"strings"
)

type preface struct {
	markedText string    //each line must end with a mark: =, *, // or +.
	phrases    []*Phrase //TODO: change to PhraseMelodyer
	//phrases    []PhraseMelodyer
}

type PrefacePhraseType string

const ( // Phrase types that can occur in a Preface
	dialogue   PrefacePhraseType = "dialogue"   // dialogue = whole initial dialogue (always the same); Special treatment, since it is always the same
	firsts     PrefacePhraseType = "firsts"     // firsts(of the paragraph) = intonation, reciting tone, short cadence; Must end with "="
	last       PrefacePhraseType = "last"       // last(of the paragraph) = reciting tone, final cadence; Must end with "//"
	mediant    PrefacePhraseType = "mediant"    // mediant = intonation, reciting tone, mediant cadence; Must end with "*"
	conclusion PrefacePhraseType = "conclusion" // conclusion = Beginning of conclusion paragraph (often "Por isso") Must end with "+"
)

// NewPreface creates a new preface struct with the marked text.
func newPreface(markedText string) *preface { //returning a pointer because this struct is going to be modified by its methods
	return &preface{
		markedText: markedText,
	}
}

// DistributeTextToPhrases takes the marked text and distributes it into phrases based on the marks at the end of each line.
// It creates a new Phrase struct for each line and appends it to the preface's phrases slice.
func (preface *preface) DistributeTextToPhrases(ctx context.Context, gen GabcGenAPI) /*(PhraseMelodyer,*/ error {

	for v := range strings.Lines(preface.markedText) {
		typedPhrase, err := preface.newTypedPhrase(v, gen.syllabifier)
		if err != nil {
			return err //TODO handle error
		}

		preface.phrases = append(preface.phrases, typedPhrase)
	}

	return nil
}

// newTypedPhrase creates a new Phrase struct based on the given string.
func (preface *preface) newTypedPhrase(s string, syllab Syllabifier) (*Phrase, error) {

	switch {
	case strings.HasSuffix(s, "="):
		s, _ = strings.CutSuffix(s, "=")
		return &Phrase{
			Raw:         s,
			PhraseTyped: firsts,
			syllabifier: syllab,
		}, nil
	case strings.HasSuffix(s, "*"):
		s, _ = strings.CutSuffix(s, "*")
		return &Phrase{
			Raw:         s,
			PhraseTyped: mediant,
			syllabifier: syllab,
		}, nil
	case strings.HasSuffix(s, "//"):
		s, _ = strings.CutSuffix(s, "//")
		return &Phrase{
			Raw:         s,
			PhraseTyped: last,
			syllabifier: syllab,
		}, nil
	case strings.HasSuffix(s, "+"):
		s, _ = strings.CutSuffix(s, "+")
		return &Phrase{
			Raw:         s,
			PhraseTyped: conclusion,
			syllabifier: syllab,
		}, nil
	default:
		return nil, fmt.Errorf("defining Phrase type from line: %w ", ErrResponseNoMarks)
	}
}

// ApplyGabcMelodies applies the GABC melodies to each phrase in the preface and returns the composed GABC string.
func (preface *preface) ApplyGabcMelodies(ctx context.Context) (string, error) {
	var composedGABC string
	for _, ph := range preface.phrases {
		gabcPhrase, err := preface.applyMelodySwitch(ph) //ph.ApplyMelody()
		if err != nil {
			return "", fmt.Errorf("applying melody: %w ", err)
		}
		composedGABC = composedGABC + gabcPhrase
	}
	return composedGABC, nil

}

// applyMelodySwitch route to the correct applyMelody function based on the phrase type.
func (preface *preface) applyMelodySwitch(ph *Phrase) (string, error) {

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

		}*/
	}
	return "", fmt.Errorf("Phrase type is none of the accepted ones: %v ", ph.PhraseTyped)
}

//type Firsts Phrase

//func (ph *Firsts) Type() PrefacePhraseType {
//	return firsts
//}

// applyMelodyFirsts applies the GABC melody to the firsts phrases of a preface.
func (ph *Phrase) applyMelodyFirsts() (string, error) {
	i := len(ph.Syllables) - 1 //reading Syllables from the end:
	if i < 0 {
		return "", fmt.Errorf("error at firsts phrase: %v: %w ", ph.Raw, ErrResponseToShort)
	}

	//last unstressed Syllables:
	for !ph.Syllables[i].IsTonic {
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(g)"
		i--
		if i < 0 {
			return "", fmt.Errorf("error at firsts phrase: %v: %w ", ph.Raw, ErrResponseToShort)
		}
	}

	//last tonic syllable:
	ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(fg)"
	i--
	if i < 0 {
		return "", fmt.Errorf("error at firsts phrase: %v: %w ", ph.Raw, ErrResponseToShort)
	}

	//syllable before the last tonic:
	ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(gf)"
	i--
	if i < 0 {
		return "", fmt.Errorf("error at firsts phrase: %v: %w ", ph.Raw, ErrResponseToShort)
	}

	//testing the exception at last unstressed reciting syllable:
	if ph.Syllables[i].IsTonic { //default case
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(h)"
		i--
		if i < 0 {
			return "", fmt.Errorf("error at firsts phrase: %v: %w ", ph.Raw, ErrResponseToShort)
		}
	} else if ph.Syllables[i-1].IsTonic && !ph.Syllables[i-1].IsLast { //exception case
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(g)"
		i--
		if i < 0 {
			return "", fmt.Errorf("error at firsts phrase: %v: %w ", ph.Raw, ErrResponseToShort)
		}
	}

	// completing reciting Syllables:
	for i > 0 {
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(h)"
		i--
		if i < 0 {
			return "", fmt.Errorf("error at firsts phrase: %v: %w ", ph.Raw, ErrResponseToShort)
		}
	}

	//first intonation syllable:
	ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(f)"

	end := "(;)"
	return joinSyllables(ph.Syllables, end), nil
}

// joinSyllables is a helper function that joins the GABC of all Syllables in a Phrase and adds the end string to it.
func joinSyllables(syl []*Syllable, end string) string {
	var result string
	for _, v := range syl {
		result = result + v.GABC
		if v.IsLast {
			result = result + " "
		}
	}

	return result + end
}
