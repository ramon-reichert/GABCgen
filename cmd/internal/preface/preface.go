// Handle specific phrase types that compose the melody of the Preface of Mass.
package preface

import (
	"fmt"

	"github.com/ramon-reichert/GABCgen/cmd/internal/gabcErrors"
	"github.com/ramon-reichert/GABCgen/cmd/internal/phrases"
	"github.com/ramon-reichert/GABCgen/cmd/internal/staff"
)

type preface struct {
	MarkedText string //each line must end with a mark: =, *, $ or +.
	Phrases    []phrases.PhraseMelodyer
}

type ( // Phrase types that can occur in a Preface
	dialogue   phrases.Phrase // dialogue = whole initial dialogue (always the same); Special treatment, since it is always the same
	firsts     phrases.Phrase // firsts(of the paragraph) = intonation, reciting tone, short cadence; Must end with "="
	last       phrases.Phrase // last(of the paragraph) = reciting tone, final cadence; Must end with "$"
	mediant    phrases.Phrase // mediant = intonation, reciting tone, mediant cadence; Must end with "*"
	conclusion phrases.Phrase // conclusion = Beginning of conclusion paragraph (often "Por isso") Must end with "+"
)

// New creates a new preface struct with the marked text.
func New(markedText string) *preface { //returning a pointer because this struct is going to be modified by its methods
	return &preface{
		MarkedText: markedText,
	}
}

// TypePhrases types the already built phrases based on the given mark suffix.
func (preface *preface) TypePhrases(newPhrases []*phrases.Phrase) error {

	for _, v := range newPhrases {
		typedPhrase, err := preface.newTypedPhrase(v)
		if err != nil {
			return fmt.Errorf("typing built Phrases: %w", err)
		}

		preface.Phrases = append(preface.Phrases, typedPhrase)
	}
	return nil
}

// newTypedPhrase switches between the possible phrase types for Preface.
func (preface *preface) newTypedPhrase(ph *phrases.Phrase) (phrases.PhraseMelodyer, error) {

	switch ph.Mark {
	case "=":
		return firsts{
			Text:      ph.Text,
			Syllables: ph.Syllables,
		}, nil
	case "$":
		return last{
			Text:      ph.Text,
			Syllables: ph.Syllables,
		}, nil
	case "*":
		return mediant{
			Text:      ph.Text,
			Syllables: ph.Syllables,
		}, nil
	case "+":
		return conclusion{
			Text:      ph.Text,
			Syllables: ph.Syllables,
		}, nil
	// TODO test the other preface types cases

	default:
		return nil, gabcErrors.ErrNoMarks
	}
}

// ApplyGabcMelodies applies the GABC melodies to each phrase in the preface and returns the composed GABC string.
func (preface *preface) ApplyGabcMelodies() (string, error) {
	var composedGABC string
	for _, ph := range preface.Phrases {
		gabcPhrase, err := ph.ApplyMelody()
		if err != nil {
			return "", fmt.Errorf("applying melody to %w", err)
		}
		composedGABC = composedGABC + gabcPhrase
	}
	return composedGABC, nil

}

// applyMelody analyzes the syllables of a phrase and attaches the GABC code(note) to each one of them, following the melody rules of that specific phrase type.
func (ph firsts) ApplyMelody() (string, error) {

	//reading Syllables from the end:
	i := len(ph.Syllables) - 1
	if i < 0 {
		return "", fmt.Errorf("firsts phrase: %v: %w ", ph.Text, gabcErrors.ErrToShort)
	}

	//last unstressed Syllables:
	for !ph.Syllables[i].IsTonic {
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.Si
		i--
		if i < 0 {
			return "", fmt.Errorf("firsts phrase: %v: %w ", ph.Text, gabcErrors.ErrToShort)
		}
	}

	//last tonic syllable:
	ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.LaSi
	i--
	if i < 0 {
		return "", fmt.Errorf("firsts phrase: %v: %w ", ph.Text, gabcErrors.ErrToShort)
	}

	//syllable before the last tonic:
	ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.SiLa
	i--
	if i < 0 {
		return "", fmt.Errorf("firsts phrase: %v: %w ", ph.Text, gabcErrors.ErrToShort)
	}

	//testing the exception at last unstressed reciting syllable:
	if ph.Syllables[i].IsTonic { //default case
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.Do
		i--
		if i < 0 {
			return "", fmt.Errorf("firsts phrase: %v: %w ", ph.Text, gabcErrors.ErrToShort)
		}
	} else if ph.Syllables[i-1].IsTonic && !ph.Syllables[i-1].IsLast { //exception case
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.Si
		i--
		if i < 0 {
			return "", fmt.Errorf("firsts phrase: %v: %w ", ph.Text, gabcErrors.ErrToShort)
		}
	}

	// completing reciting Syllables:
	for i > 0 {
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.Do
		i--
		if i < 0 {
			return "", fmt.Errorf("firsts phrase: %v: %w ", ph.Text, gabcErrors.ErrToShort)
		}
	}

	//first intonation syllable:
	ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.La

	end := "(;)" //gabc code for the "half bar", to be added at the end of the phrase
	return phrases.JoinSyllables(ph.Syllables, end), nil
}

// applyMelody analyzes the syllables of a phrase and attaches the GABC code(note) to each one of them, following the melody rules of that specific phrase type.
func (ph last) ApplyMelody() (string, error) {

	//reading Syllables from the end:
	i := len(ph.Syllables) - 1
	if i < 0 {
		return "", fmt.Errorf("last phrase: %v: %w ", ph.Text, gabcErrors.ErrToShort)
	}

	//last unstressed Syllables:
	for !ph.Syllables[i].IsTonic {
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.La
		i--
		if i < 0 {
			return "", fmt.Errorf("last phrase: %v: %w ", ph.Text, gabcErrors.ErrToShort)
		}
	}

	//last tonic syllable from oxytone:
	if ph.Syllables[i].IsTonic && ph.Syllables[i].IsLast {
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.LaSiLa
	} else { //last tonic syllable from non oxytone:
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.LaSi
	}
	i--
	if i < 0 {
		return "", fmt.Errorf("last phrase: %v: %w ", ph.Text, gabcErrors.ErrToShort)
	}

	//first syllable before the last tonic:
	ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.Si
	i--
	if i < 0 {
		return "", fmt.Errorf("last phrase: %v: %w ", ph.Text, gabcErrors.ErrToShort)
	}

	//second syllable before the last tonic:
	ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.SolLa
	i--
	if i < 0 {
		return "", fmt.Errorf("last phrase: %v: %w ", ph.Text, gabcErrors.ErrToShort)
	}

	//third syllable before the last tonic:
	ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.LaSol
	i--
	if i < 0 {
		return "", fmt.Errorf("last phrase: %v: %w ", ph.Text, gabcErrors.ErrToShort)
	}

	// completing reciting Syllables:
	for i > 0 {
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.Si
		i--
		if i < 0 {
			return "", fmt.Errorf("last phrase: %v: %w ", ph.Text, gabcErrors.ErrToShort)
		}
	}

	//first syllable of the phrase:
	ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.Si

	end := "(:)" //gabc code for the "whole bar", to be added at the end of the phrase
	return phrases.JoinSyllables(ph.Syllables, end), nil
}

// applyMelody analyzes the syllables of a phrase and attaches the GABC code(note) to each one of them, following the melody rules of that specific phrase type.
// There is no clear and strict rule observable in the Latin Prefaces on exactly "where" to apply the notes of the the mediant melody,
// so one possible solution was chosen here, which sounds natural according to current pratices in Brazilian liturgy.
func (ph mediant) ApplyMelody() (string, error) {

	//reading Syllables from the end:
	i := len(ph.Syllables) - 1
	if i < 0 {
		return "", fmt.Errorf("mediant phrase: %v: %w ", ph.Text, gabcErrors.ErrToShort)
	}

	//last unstressed Syllables:
	for !ph.Syllables[i].IsTonic {
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.Si
		i--
		if i < 0 {
			return "", fmt.Errorf("mediant phrase: %v: %w ", ph.Text, gabcErrors.ErrToShort)
		}
	}

	//last tonic syllable:
	if i >= 3 { //means that the melody can be spared throughout 3 syllables
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.Do
		i--
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.Si
		i--
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.La
		i--
	} else { //means that the whole melody must be upon this last tonic syllable
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.LaSiDo
		i--
	}

	// completing reciting Syllables:
	for i > 0 {
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.Si
		i--
	}
	if i == 0 {
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.Si
	}

	end := "(,)" //gabc code for the "quarter bar", to be added at the end of the phrase
	return phrases.JoinSyllables(ph.Syllables, end), nil
}

// applyMelody analyzes the syllables of a phrase and attaches the GABC code(note) to each one of them, following the melody rules of that specific phrase type.
func (ph conclusion) ApplyMelody() (string, error) {

	//reading Syllables from the end:
	i := len(ph.Syllables) - 1
	if i < 0 {
		return "", fmt.Errorf("conclusion phrase: %v: %w ", ph.Text, gabcErrors.ErrToShort)
	}

	//last unstressed Syllables:
	for !ph.Syllables[i].IsTonic {
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.La
		i--
		if i < 0 {
			return "", fmt.Errorf("conclusion phrase: %v: %w ", ph.Text, gabcErrors.ErrToShort)
		}
	}

	//last tonic syllable:
	ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.SolLa
	i--
	if i < 0 {
		return "", fmt.Errorf("conclusion phrase: %v: %w ", ph.Text, gabcErrors.ErrToShort)
	}

	// completing reciting Syllables:
	for i > 0 {
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.La
		i--
	}
	if i == 0 {
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.La
	}

	end := "(,)" //gabc code for the "quarter bar", to be added at the end of the phrase
	return phrases.JoinSyllables(ph.Syllables, end), nil
}
