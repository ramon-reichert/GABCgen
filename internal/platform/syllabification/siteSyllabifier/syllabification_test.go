package sitesyllabifier_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/matryer/is"
	"github.com/ramon-reichert/gabcgen/internal/platform/syllabification/sitesyllabifier"
)

var ctx context.Context = context.Background()

func TestSyllabify(t *testing.T) {
	is := is.New(t)
	syllabifier := sitesyllabifier.NewSyllabifier("test_liturgical_syllables.json", "test_user_syllables.json", "test_not_syllabified.txt")
	is.NoErr(os.WriteFile("test_user_syllables.json", []byte("{}"), 0644)) //write an empty json file to the user syllables path

	t.Run("fetch syllables from words that are already at liturgical syllabs db file", func(t *testing.T) {
		is := is.New(t)

		jsonWord := map[string]sitesyllabifier.SyllableInfo{
			"litúrgicas": {
				Slashed:    "fetched/in/liturgical/db",
				TonicIndex: 2,
			}}
		data, err := json.MarshalIndent(jsonWord, "", "  ")
		is.NoErr(err)
		is.NoErr(os.WriteFile("test_liturgical_syllables.json", data, 0644))
		is.NoErr(syllabifier.LoadSyllables())

		slashed, tonicIndex, err := syllabifier.Syllabify(ctx, "litúrgicas")
		is.NoErr(err)
		is.Equal(slashed, "fetched/in/liturgical/db") // proposital wrong answer, to ensure that the syllables were fetched from the liturgical db
		is.Equal(tonicIndex, 2)
	})

	t.Run("fetch syllables from new words at external site", func(t *testing.T) {
		is := is.New(t)

		newWord := "externo"
		slashed, tonicIndex, err := syllabifier.Syllabify(ctx, newWord)
		is.NoErr(err)
		is.Equal(slashed, "ex/ter/no") // word not present in liturgical db or user db, so it was fetched from external site
		is.Equal(tonicIndex, 2)

		syllabifier.SaveSyllables()

		jsonWord := map[string]sitesyllabifier.SyllableInfo{
			newWord: {
				Slashed:    slashed,
				TonicIndex: tonicIndex,
			},
		}
		data, err := json.MarshalIndent(jsonWord, "", "  ")
		is.NoErr(err)
		fileContent, err := os.ReadFile("test_user_syllables.json")
		is.NoErr(err)
		is.Equal(fileContent, data) // check if the user db file was created with the new word
	})
}
