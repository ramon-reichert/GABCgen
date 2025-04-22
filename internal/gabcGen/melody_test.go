package gabcGen_test

import (
	"context"
	"log"
	"testing"

	"github.com/matryer/is"
	"github.com/ramon-reichert/GABCgen/internal/gabcGen"
	"github.com/ramon-reichert/GABCgen/internal/syllabification"
)

var ctx context.Context = context.Background()

func TestApplyMelody(t *testing.T) {

	t.Run("apply gabc melody to firsts phrase", func(t *testing.T) {
		is := is.New(t)

		phrase, err := gabcGen.NewGabcGen(syllabification.NewSyllabifier()).BuildPhrase(ctx, "-Na: verd'ade, é .digno e justo,=")
		is.NoErr(err)

		for _, v := range phrase.Syllables { //DEBUG code
			log.Println("\n Char: ", v.Char)
			log.Println("Char string: ", string(v.Char))
			log.Println("GABC: ", v.GABC)
			log.Println("IsTonic: ", v.IsTonic)
			log.Println("IsFirst: ", v.IsFirst)
			log.Println("IsLast: ", v.IsLast)
		}

		expectedGABC := "-Na:(f) ver(h)d'a(h)de,(h) é(h) .dig(h)no(g) e(gf) jus(fg)to,(g) (;)"

		composedGABC, err := phrase.ApplyMelodyGABC(ctx)
		is.NoErr(err)
		is.Equal(composedGABC, expectedGABC)
	})
}
