package gabcGen_test

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"github.com/ramon-reichert/GABCgen/cmd/internal/gabcGen"
	"github.com/ramon-reichert/GABCgen/cmd/internal/syllabification"
)

var ctx context.Context = context.Background()

func TestGeneratePreface(t *testing.T) {

	t.Run("apply gabc melody to a preface firsts phrase", func(t *testing.T) {
		is := is.New(t)

		firstsPhrase, err := gabcGen.NewGabcGenAPI(syllabification.NewSyllabifier()).GeneratePreface(ctx, "-Na: verd'ade, é .digno e justo,=")
		is.NoErr(err)

		expectedGABC := "-Na:(f) ver(h)d'a(h)de,(h) é(h) .dig(h)no(g) e(gf) jus(fg)to,(g) (;)"

		is.Equal(firstsPhrase.Url, expectedGABC)
	})
}
