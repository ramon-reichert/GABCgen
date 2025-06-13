package preface_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/matryer/is"
	"github.com/ramon-reichert/GABCgen/cmd/internal/gabcGen"
	"github.com/ramon-reichert/GABCgen/cmd/internal/phrases"
	"github.com/ramon-reichert/GABCgen/cmd/internal/syllabification/mockSyllabifier"
)

var ctx context.Context = context.Background()

func TestGeneratePreface(t *testing.T) {

	t.Run("apply gabc melodies to a group of phrases using mockSyllabifier", func(t *testing.T) {
		is := is.New(t)

		a := "\n"                                 //spare new line
		b := "-Na: verd'ade, é .digno e justo,\n" // paroxytone with exception - firsts
		c := "Na verdade, digno, justo,\n"        //paroxytone without exception" - firsts
		d := "Na verdade, digno e justo é,\n"     //oxytone with exception - firsts
		e := "Na verdade, é digno e justo\n"      //paraxytone - mediant
		f := "-Na: verd'ade, é .digno e justo,\n" //paroxytone - last
		g := "\n"                                 //new line separating paragraphs
		h := " \n"                                //spare new line and space
		l := "Por isso, na verdade,\n"            //conclusion phrase short
		//b again
		i := "Na verdade\n"                   //mediant phrase - 3 syllables
		j := "Na verdade, digno e justo é,\n" //last phrase - oxytone
		//g again
		//l again
		//c again
		k := "digno\n" //mediant phrase - 1 syllable
		//f again
		//a again

		inputText := fmt.Sprint(a + b + c + d + e + f + g + h + l + b + i + j + g + l + c + k + f + a)
		//log.Println("inputText: ", inputText)

		syllabifier := mockSyllabifier.NewSyllabifier()

		composedGABC, err := gabcGen.NewGabcGenAPI(syllabifier).GeneratePreface(ctx, inputText)
		is.NoErr(err)

		expectedGABC := "-Na:(f) ver(h)d'a(h)de,(h) é(h) .dig(h)no(g) e(gf) jus(fg)to,(g) (;)\nNa(f) ver(h)da(h)de,(h) dig(h)no,(gf) jus(fg)to,(g) (;)\nNa(f) ver(h)da(h)de,(h) dig(h)no(h) e(h) jus(h)to(gf) é,(fg) (;)\nNa(g) ver(g)da(g)de,(g) é(g) dig(g)no(f) e(g) jus(h)to(g) (,)\n-Na:(g) ver(g)d'a(g)de,(g) é(g) .dig(fe)no(ef) e(g) jus(fg)to,(f) (:)(Z)\n\nPor(f) is(h)so,(h) na(h) ver(gf)da(fg)de,(g) (;)\n-Na:(f) ver(h)d'a(h)de,(h) é(h) .dig(h)no(g) e(gf) jus(fg)to,(g) (;)\nNa(g) ver(g)da(fgh)de(g) (,)\nNa(g) ver(g)da(g)de,(g) dig(g)no(g) e(fe) jus(ef)to(g) é,(fgf) (:)(Z)\n\nPor(f) is(f)so,(f) na(f) ver(f)da(ef)de,(f) (,)\nNa(f) ver(h)da(h)de,(h) dig(h)no,(gf) jus(fg)to,(g) (;)\ndig(fgh)no(g) (,)\n-Na:(g) ver(g)d'a(g)de,(g) é(g) .dig(fe)no(ef) e(g) jus(fg)to,(f) (::)"

		is.Equal(composedGABC, expectedGABC)

	})
}

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
