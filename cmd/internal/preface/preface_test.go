package preface_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/matryer/is"
	"github.com/ramon-reichert/GABCgen/cmd/internal/gabcGen"
	"github.com/ramon-reichert/GABCgen/cmd/internal/preface"
	"github.com/ramon-reichert/GABCgen/cmd/internal/syllabification/mockSyllabifier"
	"golang.org/x/text/unicode/norm"
)

var ctx context.Context = context.Background()

func TestGeneratePreface(t *testing.T) {
	syllabifier := mockSyllabifier.NewSyllabifier()

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

		composedPreface, err := gabcGen.NewGabcGenAPI(syllabifier).GeneratePreface(ctx, preface.Preface{Text: preface.PrefaceText{LinedText: inputText}})
		is.NoErr(err)

		composedGABC := composedPreface.Text.ComposedGABC

		expectedGABC := "<c><sp>V/</sp></c> -Na:(f) ver(h)d'a(h)de,(h) é(h) .dig(h)no(g) e(gf) jus(fg)to,(g) (;)\nNa(f) ver(h)da(h)de,(h) dig(h)no,(gf) jus(fg)to,(g) (;)\nNa(f) ver(h)da(h)de,(h) dig(h)no(h) e(h) jus(h)to(gf) é,(fg) (;)\nNa(g) ver(g)da(g)de,(g) é(g) dig(g)no(f) e(g) jus(h)to(g) (,)\n-Na:(g) ver(g)d'a(g)de,(g) é(g) .dig(fe)no(ef) e(g) jus(fg)to,(f) (:)(Z)\n\nPor(f) is(h)so,(h) na(h) ver(gf)da(fg)de,(g) (;)\n-Na:(f) ver(h)d'a(h)de,(h) é(h) .dig(h)no(g) e(gf) jus(fg)to,(g) (;)\nNa(g) ver(g)da(fgh)de(g) (,)\nNa(g) ver(g)da(g)de,(g) dig(g)no(g) e(fe) jus(ef)to(g) é,(fgf) (:)(Z)\n\nPor(f) is(f)so,(f) na(f) ver(f)da(ef)de,(f) (,)\nNa(f) ver(h)da(h)de,(h) dig(h)no,(gf) jus(fg)to,(g) (;)\ndig(fgh)no(g) (,)\n-Na:(g) ver(g)d'a(g)de,(g) é(g) .dig(fe)no(ef) e(g) jus(fg)to,(f) (::)"

		is.Equal(composedGABC, expectedGABC)

	})

	t.Run("exception attempt to apply 'last' melody to short phrase like 'Senhor nosso'", func(t *testing.T) {
		is := is.New(t)

		inputText := "Na verdade, é digno e justo,\n por Cristo,\n Senhor nosso."

		composedPreface, err := gabcGen.NewGabcGenAPI(syllabifier).GeneratePreface(ctx, preface.Preface{Text: preface.PrefaceText{LinedText: inputText}})
		is.NoErr(err)

		composedGABC := composedPreface.Text.ComposedGABC

		expectedGABC := "<c><sp>V/</sp></c> Na(f) ver(h)da(h)de,(h) é(h) dig(h)no(g) e(gf) jus(fg)to,(g) (;)\npor(g) Cris(fgh)to,(g) (,)\nSe(fe)nhor(efg) nos(fg)so.(f) (::)"

		//dmp := diffmatchpatch.New()
		//diffs := dmp.DiffMainRunes([]rune(norm.NFC.String(composedGABC)), []rune(norm.NFC.String(expectedGABC)), false)
		//log.Println("\n\ndiffs: ", diffs)

		is.Equal(norm.NFC.String(composedGABC), norm.NFC.String(expectedGABC))

	})
}
