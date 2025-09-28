// Package phrases handle musical phrases composed of words.Syllable structs.
// Phrases can be typed according to the Mass part.
package phrases

import (
	"context"
	"fmt"
	"strings"

	gabcErrors "github.com/ramon-reichert/GABCgen/internal/platform/errors"
	"github.com/ramon-reichert/GABCgen/internal/service/composition/phrases/words"
)

type Phrase struct {
	Text        string            // text of the phrase
	Syllables   []*words.Syllable // syllables of the phrase
	Syllabifier words.Syllabifier // Syllabifier to be used to syllabify the words of the phrase
	Directives  []Directive       // possible singing directives may come between parentheses and are not to be sung. They are removed from the text before the syllabification and should be put back again after the melody is applied.
}

type Directive struct {
	Text   string
	Before string
}

// New creates a new Phrase instance with the provided text.
func New(text string) *Phrase {
	return &Phrase{
		Text: text,
	}
}

// BuildPhraseSyllables populates a Phrase.Syllables iterating over each word of the Phrase.
func (ph *Phrase) BuildPhraseSyllables(ctx context.Context) error {
	words := strings.Fields(ph.Text)

	for _, v := range words {
		syllables, err := ph.classifyWordSyllables(ctx, v)
		if err != nil {
			return fmt.Errorf("building Phrase Syllables: %w ", err)
		}
		ph.Syllables = append(ph.Syllables, syllables...)
	}

	return nil
}

// classifyWordSyllables divides the syllables of a word and builds a Syllable struct from each one of them.
func (ph *Phrase) classifyWordSyllables(ctx context.Context, word string) ([]*words.Syllable, error) {
	var wordSyllables []*words.Syllable
	wordMaped := words.New(word)

	if err := wordMaped.ParseWord(); err != gabcErrors.ErrNoLetters { // Early scape to avoid trying to syllabify a "non-letter word"
		if err := wordMaped.Syllabify(ctx, ph.Syllabifier); err != nil {
			return wordSyllables, fmt.Errorf("classifying word syllables: %w", err)
		}
	}

	wordMaped.RecomposeWord()
	wordSyllables = wordMaped.BuildWordSyllables()

	return wordSyllables, nil
}

// JoinSyllables is a helper function that joins the GABC of all Syllables in a Phrase and adds the end string to it.
// It also attempts to put the directives back into the right place.
func JoinSyllables(syl []*words.Syllable, end string, d []Directive) string {
	var result string
	var pool string
	dirIndex := 0

	for i, v := range syl {
		// Join the GABC of each syllable
		result += v.GABC

		if v.IsLast {
			result += " "
		}

		// Put the directive - if it exists - back into the right place
		if i < len(syl)-1 { // skip last syllable to avoid conflicts with the end marker
			pool += string(v.Char)

			if dirIndex < len(d) && strings.HasSuffix(pool, d[dirIndex].Before) { // compares with the letters that were before the directive at the moment it was removed from the original phrase.
				result += "||<i><c>" + d[dirIndex].Text + "</c></i>||" + "(,) "
				dirIndex++
			}
		}
	}

	for dirIndex < len(d) { // if there is no match, the directive is put at the end of the phrase
		result += "||<i><c>" + d[dirIndex].Text + "</c></i>||"
		dirIndex++
	}

	return result + end
}

// ExtractDirectives looks for directives between parentheses in the Phrase text, extracts them and stores them in the Phrase.Directives slice.
func (ph *Phrase) ExtractDirectives() error {
	var pool string

	for leftOriginal, after, open := strings.Cut(ph.Text, "("); open; {
		extracted, rightOriginal, close := strings.Cut(after, ")")

		if !close {
			return fmt.Errorf("missing closer parentheses in: %v", ph.Text)
		}

		for v := range strings.FieldsSeq(leftOriginal) {
			pool += v
		}

		ph.Directives = append(ph.Directives, Directive{Text: extracted, Before: pool})

		ph.Text = strings.TrimSpace(leftOriginal) + " " + strings.TrimSpace(rightOriginal)
		leftOriginal, after, open = strings.Cut(ph.Text, "(")
	}

	return nil
}
