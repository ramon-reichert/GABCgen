// Handle musical phrases composed of words.Syllable structs.
// Phrases can be typed according to the Mass part.
package phrases

import (
	"context"
	"fmt"
	"strings"

	"github.com/ramon-reichert/GABCgen/internal/domain/composition/paragraph/phrases/words"
	gabcErrors "github.com/ramon-reichert/GABCgen/internal/platform/errors"
)

type Phrase struct {
	Text        string            // the text of the phrase
	Syllables   []*words.Syllable // the syllables of the phrase
	Syllabifier words.Syllabifier // the Syllabifier to be used to syllabify the words of the phrase
	Directives  []Directive       // possible singing directives may come between parentheses and are not to be sung. They are removed from the text before the syllabification and should be put back again after the melody is applied.
}

type PhraseMelodyer interface {
	ApplyMelody() (string, error) //Applying the Open/Closed principle from SOLID so we can always have new types of Phrases
}

type Directive struct {
	Text   string
	Before string
}

func New(text string) *Phrase {
	return &Phrase{
		Text: text,
	}
}

// BuildSyllabes populates a Phrase.Syllables creating Syllable structs from each word of the Phrase.
func (ph *Phrase) BuildPhraseSyllables(ctx context.Context) error {

	words := strings.Fields(ph.Text)
	for _, v := range words {

		//TODO: MAYBE! verify if word is composed and divide it at the hyphen, or not (ave-maria, compará-lo)

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
	//debug code:
	//for i, v := range d {
	//	log.Printf("directive just received in JoinSyllables:\n d.[%v]: %+v", i, v)
	//}

	var result string
	var pool string
	dirIndex := 0
	for i, v := range syl {
		// join the GABC of each syllable:
		result += v.GABC
		if v.IsLast {
			result += " "

		}

		// put the directive - if it exists - back into the right place:
		if i < len(syl)-1 { //skip last syllable to avoid conflicts with the end marker
			pool += string(v.Char)
			if dirIndex < len(d) && strings.HasSuffix(pool, d[dirIndex].Before) { //compares with the letters that were before the directive at the moment it was removed from the original phrase.
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
		//debug code:
		//log.Printf("leftOriginal: %v\n pool: %v\n", leftOriginal, pool)

		ph.Directives = append(ph.Directives, Directive{Text: extracted, Before: pool})

		ph.Text = strings.TrimSpace(leftOriginal) + " " + strings.TrimSpace(rightOriginal)
		leftOriginal, after, open = strings.Cut(ph.Text, "(")
	}

	//debug:
	//log.Printf("directives extracted: %+v", ph.Directives)

	return nil
}
