// Handle specific phrase types that compose the melody of the Preface of Mass.
package preface

import (
	"fmt"
	"strings"

	"github.com/ramon-reichert/GABCgen/cmd/internal/gabcErrors"
	"github.com/ramon-reichert/GABCgen/cmd/internal/paragraph"
	"github.com/ramon-reichert/GABCgen/cmd/internal/phrases"
	"github.com/ramon-reichert/GABCgen/cmd/internal/staff"
)

type preface struct {
	LinedText string
	Phrases   []phrases.PhraseMelodyer
}

type ( // Phrase types that can occur in a Preface
	dialogue   phrases.Phrase // dialogue = whole initial dialogue (always the same text); Special treatment: just the melody can differ between simple or solemn tones
	firsts     phrases.Phrase // firsts(of the paragraph) = intonation, reciting tone, short cadence;
	last       phrases.Phrase // last(of the paragraph) = reciting tone, final cadence;
	mediant    phrases.Phrase // mediant = intonation, reciting tone, mediant cadence;
	conclusion phrases.Phrase // conclusion = Beginning of conclusion paragraph (often "Por isso")
)

// New creates a new preface struct with the lined text.
func New(linedText string) *preface { //returning a pointer because this struct is going to be modified by its methods
	return &preface{
		LinedText: linedText,
	}
}

// TypePhrases types the already built phrases based on the position of the phrases in the paragraph.
func (preface *preface) TypePhrases(newParagraphs []paragraph.Paragraph) error {

	for n, p := range newParagraphs {

		if n == len(newParagraphs)-1 && strings.HasPrefix(p.Phrases[0].Text, "Por isso") { //"Por isso" is a special conclusion expression that can start the last paragraph, and has its own melody
			preface.Phrases = append(preface.Phrases, conclusion{
				Text:      p.Phrases[0].Text,
				Syllables: p.Phrases[0].Syllables,
			})

			p.Phrases = p.Phrases[1:] //removing the conclusion phrase from the paragraph, so it won't be processed again
		}

		if len(p.Phrases) < 3 { //each paragraph must have at least three phrases to enable applying the melody, not counting the conclusion phrase - which can start the last paragraph
			return fmt.Errorf("typing phrase: %v - %w", p.Phrases[0].Text, gabcErrors.ErrShortParagraph)
		}

		for i := 0; i < len(p.Phrases); i++ {

			if i < len(p.Phrases)-2 {
				preface.Phrases = append(preface.Phrases, firsts{
					Text:      p.Phrases[i].Text,
					Syllables: p.Phrases[i].Syllables,
				})
				continue
			}
			if i == len(p.Phrases)-2 {
				preface.Phrases = append(preface.Phrases, mediant{
					Text:      p.Phrases[i].Text,
					Syllables: p.Phrases[i].Syllables,
				})
				continue
			}
			preface.Phrases = append(preface.Phrases, last{
				Text:      p.Phrases[i].Text,
				Syllables: p.Phrases[i].Syllables,
			})
		}
	}
	return nil
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
		return "", fmt.Errorf("firsts phrase: %v: %w ", ph.Text, gabcErrors.ErrShortPhrase)
	}

	//last unstressed Syllables:
	for !ph.Syllables[i].IsTonic {
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.Si
		i--
		if i < 0 {
			return "", fmt.Errorf("firsts phrase: %v: %w ", ph.Text, gabcErrors.ErrShortPhrase)
		}
	}

	//last tonic syllable:
	ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.LaSi
	i--
	if i < 0 {
		return "", fmt.Errorf("firsts phrase: %v: %w ", ph.Text, gabcErrors.ErrShortPhrase)
	}

	//syllable before the last tonic:
	ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.SiLa
	i--
	if i < 0 {
		return "", fmt.Errorf("firsts phrase: %v: %w ", ph.Text, gabcErrors.ErrShortPhrase)
	}

	//testing the exception at last unstressed reciting syllable:
	if ph.Syllables[i].IsTonic { //default case
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.Do
		i--
		if i < 0 {
			return "", fmt.Errorf("firsts phrase: %v: %w ", ph.Text, gabcErrors.ErrShortPhrase)
		}
	} else if ph.Syllables[i-1].IsTonic && !ph.Syllables[i-1].IsLast { //exception case
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.Si
		i--
		if i < 0 {
			return "", fmt.Errorf("firsts phrase: %v: %w ", ph.Text, gabcErrors.ErrShortPhrase)
		}
	}

	// completing reciting Syllables:
	for i > 0 {
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.Do
		i--
		if i < 0 {
			return "", fmt.Errorf("firsts phrase: %v: %w ", ph.Text, gabcErrors.ErrShortPhrase)
		}
	}

	//first intonation syllable:
	ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.La

	end := "(;)\n" //gabc code for the "half bar", to be added at the end of the phrase
	return phrases.JoinSyllables(ph.Syllables, end), nil
}

// applyMelody analyzes the syllables of a phrase and attaches the GABC code(note) to each one of them, following the melody rules of that specific phrase type.
func (ph last) ApplyMelody() (string, error) {

	//reading Syllables from the end:
	i := len(ph.Syllables) - 1
	if i < 0 {
		return "", fmt.Errorf("last phrase: %v: %w ", ph.Text, gabcErrors.ErrShortPhrase)
	}

	//last unstressed Syllables:
	for !ph.Syllables[i].IsTonic {
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.La
		i--
		if i < 0 {
			return "", fmt.Errorf("last phrase: %v: %w ", ph.Text, gabcErrors.ErrShortPhrase)
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
		return "", fmt.Errorf("last phrase: %v: %w ", ph.Text, gabcErrors.ErrShortPhrase)
	}

	//first syllable before the last tonic:
	ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.Si
	i--
	if i < 0 {
		return "", fmt.Errorf("last phrase: %v: %w ", ph.Text, gabcErrors.ErrShortPhrase)
	}

	//second syllable before the last tonic:
	ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.SolLa
	i--
	if i < 0 {
		return "", fmt.Errorf("last phrase: %v: %w ", ph.Text, gabcErrors.ErrShortPhrase)
	}

	//third syllable before the last tonic:
	ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.LaSol
	i--
	if i < 0 {
		return "", fmt.Errorf("last phrase: %v: %w ", ph.Text, gabcErrors.ErrShortPhrase)
	}

	// completing reciting Syllables:
	for i > 0 {
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.Si
		i--
		if i < 0 {
			return "", fmt.Errorf("last phrase: %v: %w ", ph.Text, gabcErrors.ErrShortPhrase)
		}
	}

	//first syllable of the phrase:
	ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.Si

	end := "(:)\n\n" //gabc code for the "whole bar", to be added at the end of the phrase
	return phrases.JoinSyllables(ph.Syllables, end), nil
}

// applyMelody analyzes the syllables of a phrase and attaches the GABC code(note) to each one of them, following the melody rules of that specific phrase type.
// There is no clear and strict rule observable in the Latin Prefaces on exactly "where" to apply the notes of the the mediant melody,
// so one possible solution was chosen here, which sounds natural according to current pratices in Brazilian liturgy.
func (ph mediant) ApplyMelody() (string, error) {

	//reading Syllables from the end:
	i := len(ph.Syllables) - 1
	if i < 0 {
		return "", fmt.Errorf("mediant phrase: %v: %w ", ph.Text, gabcErrors.ErrShortPhrase)
	}

	//last unstressed Syllables:
	for !ph.Syllables[i].IsTonic {
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.Si
		i--
		if i < 0 {
			return "", fmt.Errorf("mediant phrase: %v: %w ", ph.Text, gabcErrors.ErrShortPhrase)
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

	end := "(,)\n" //gabc code for the "quarter bar", to be added at the end of the phrase
	return phrases.JoinSyllables(ph.Syllables, end), nil
}

// applyMelody analyzes the syllables of a phrase and attaches the GABC code(note) to each one of them, following the melody rules of that specific phrase type.
func (ph conclusion) ApplyMelody() (string, error) {

	//reading Syllables from the end:
	i := len(ph.Syllables) - 1
	if i < 0 {
		return "", fmt.Errorf("conclusion phrase: %v: %w ", ph.Text, gabcErrors.ErrShortPhrase)
	}

	//last unstressed Syllables:
	for !ph.Syllables[i].IsTonic {
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.La
		i--
		if i < 0 {
			return "", fmt.Errorf("conclusion phrase: %v: %w ", ph.Text, gabcErrors.ErrShortPhrase)
		}
	}

	//last tonic syllable:
	ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.SolLa
	i--
	if i < 0 {
		return "", fmt.Errorf("conclusion phrase: %v: %w ", ph.Text, gabcErrors.ErrShortPhrase)
	}

	// completing reciting Syllables:
	for i > 0 {
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.La
		i--
	}
	if i == 0 {
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + staff.La
	}

	end := "(,)\n" //gabc code for the "quarter bar", to be added at the end of the phrase
	return phrases.JoinSyllables(ph.Syllables, end), nil
}
