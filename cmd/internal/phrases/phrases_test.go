package phrases_test

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"github.com/ramon-reichert/GABCgen/cmd/internal/phrases"
)

var ctx context.Context = context.Background()

func TestExtractDirectives(t *testing.T) {
	t.Run("extract a single directive from a phrase", func(t *testing.T) {
		is := is.New(t)

		ph := phrases.New("before parentheses (inside parentheses) after parentheses")
		err := ph.ExtractDirectives()
		is.NoErr(err)
		is.Equal(ph.Directives[0], "inside parentheses")
		is.Equal(ph.Text, "before parentheses after parentheses")

	})

	t.Run("extract a many directives from a single phrase", func(t *testing.T) {
		is := is.New(t)

		ph := phrases.New("before parentheses 1 (inside parentheses 1) after parentheses 1 and before 2 (inside 2) after 2")
		err := ph.ExtractDirectives()
		is.NoErr(err)
		is.Equal(ph.Directives[0], "inside parentheses 1")
		is.Equal(ph.Directives[1], "inside 2")
		is.Equal(ph.Text, "before parentheses 1 after parentheses 1 and before 2 after 2")

	})
}
