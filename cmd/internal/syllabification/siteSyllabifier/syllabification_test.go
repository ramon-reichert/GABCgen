package siteSyllabifier_test

import (
	"context"
	"testing"
)

var ctx context.Context = context.Background()

func TestSyllabify(t *testing.T) {

	t.Run("fetch syllables from words that are already at liturgical syllabs db file", func(t *testing.T) {
		//	is := is.New(t)
		//		syllabifier := syllabification.NewSyllabifier("B:/dev/GABCgen/cmd/user_syllables.json")

		//	is.NoErr(syllabifier.LoadSyllables())

	})

	t.Run("fetch syllables from new words at external site", func(t *testing.T) {
		//	is := is.New(t)

	})

	t.Run("fetch syllables inserted at user syllabs db file", func(t *testing.T) {
		//	is := is.New(t)

	})
}
