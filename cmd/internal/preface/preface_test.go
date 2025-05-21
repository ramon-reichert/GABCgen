package preface_test

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"github.com/ramon-reichert/GABCgen/cmd/internal/gabcGen"
	"github.com/ramon-reichert/GABCgen/cmd/internal/syllabification"
)

var ctx context.Context = context.Background()

func TestGeneratePreface(t *testing.T) {

	t.Run("apply gabc melody to a preface firsts phrase - paroxytone with exception", func(t *testing.T) {
		is := is.New(t)

		firstsPhrase, err := gabcGen.NewGabcGenAPI(syllabification.NewSyllabifier()).GeneratePreface(ctx, "-Na: verd'ade, é .digno e justo,=")
		is.NoErr(err)

		expectedGABC := "-Na:(f) ver(h)d'a(h)de,(h) é(h) .dig(h)no(g) e(gf) jus(fg)to,(g) (;)"

		is.Equal(firstsPhrase.Url, expectedGABC)
	})

	t.Run("apply gabc melody to a preface firsts phrase - paroxytone without exception", func(t *testing.T) {
		is := is.New(t)

		firstsPhrase, err := gabcGen.NewGabcGenAPI(syllabification.NewSyllabifier()).GeneratePreface(ctx, "Na verdade, digno, justo,=")
		is.NoErr(err)

		expectedGABC := "Na(f) ver(h)da(h)de,(h) dig(h)no,(gf) jus(fg)to,(g) (;)"

		is.Equal(firstsPhrase.Url, expectedGABC)
	})

	t.Run("apply gabc melody to a preface firsts phrase - oxytone", func(t *testing.T) {
		is := is.New(t)

		firstsPhrase, err := gabcGen.NewGabcGenAPI(syllabification.NewSyllabifier()).GeneratePreface(ctx, "Na verdade, digno e justo é,=")
		is.NoErr(err)

		expectedGABC := "Na(f) ver(h)da(h)de,(h) dig(h)no(h) e(h) jus(h)to(gf) é,(fg) (;)"

		is.Equal(firstsPhrase.Url, expectedGABC)
	})

	t.Run("apply gabc melody to a preface last phrase - paroxytone", func(t *testing.T) {
		is := is.New(t)

		lastPhrase, err := gabcGen.NewGabcGenAPI(syllabification.NewSyllabifier()).GeneratePreface(ctx, "-Na: verd'ade, é .digno e justo,$")
		is.NoErr(err)

		expectedGABC := "-Na:(g) ver(g)d'a(g)de,(g) é(g) .dig(fe)no(ef) e(g) jus(fg)to,(f) (:)"

		is.Equal(lastPhrase.Url, expectedGABC)
	})

	t.Run("apply gabc melody to a preface last phrase - oxytone", func(t *testing.T) {
		is := is.New(t)

		lastPhrase, err := gabcGen.NewGabcGenAPI(syllabification.NewSyllabifier()).GeneratePreface(ctx, "Na verdade, digno e justo é,$")
		is.NoErr(err)

		expectedGABC := "Na(g) ver(g)da(g)de,(g) dig(g)no(g) e(fe) jus(ef)to(g) é,(fgf) (:)"

		is.Equal(lastPhrase.Url, expectedGABC)
	})

}
