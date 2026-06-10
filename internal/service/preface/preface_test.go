package preface_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/matryer/is"
	"github.com/ramon-reichert/gabcgen/internal/platform/syllabification/mocksyllabifier"
	"github.com/ramon-reichert/gabcgen/internal/service"
	"golang.org/x/text/unicode/norm"
)

var ctx context.Context = context.Background()

func TestGeneratePreface(t *testing.T) {
	syllabifier := mocksyllabifier.NewSyllabifier()

	t.Run("apply gabc melodies to a group of phrases using mockSyllabifier", func(t *testing.T) {
		is := is.New(t)

		a := "\n"                                 // spare new line
		b := "-Na: verd'ade, é .digno e justo,\n" // paroxytone with exception - firsts
		c := "Na verdade, digno, justo,\n"        // paroxytone without exception" - firsts
		d := "Na verdade, digno e justo é,\n"     // oxytone with exception - firsts
		e := "Na verdade, é digno e justo\n"      // paraxytone - mediant
		f := "-Na: verd'ade, é .digno e justo,\n" // paroxytone - last
		g := "\n"                                 // new line separating paragraphs
		h := " \n"                                // spare new line and space
		l := "Por isso, na verdade,\n"            // conclusion phrase short
		// b again
		i := "Na verdade\n"                   // mediant phrase - 3 syllables
		j := "Na verdade, digno e justo é,\n" // last phrase - oxytone
		// g again
		// l again
		// c again
		k := "digno\n" // mediant phrase - 1 syllable
		// f again
		// a again

		inputText := fmt.Sprint(a + b + c + d + e + f + g + h + l + b + i + j + g + l + c + k + f + a)

		composedGABC, err := service.NewGabcGenAPI(syllabifier).GeneratePreface(ctx, "", inputText)
		is.NoErr(err)

		expectedGABC := "<c><sp>V/</sp></c> O(e) Se(f)nhor(g) es(g)te(g)ja(e) con(f)vos(gf)co.(f) (::) <c><sp>R/</sp></c> E(e)<e>le</e> es(f)tá(g) no(g) me(g)io(e) de(f) nós.(gf) (::) (Z) <c><sp>V/</sp></c> Co(f)ra(g)ções(h) ao(g) al(fg)to.(fe) (::) <c><sp>R/</sp></c> O(g) nos(g)so(g) co(f)ra(g)cão(h) es(g)tá(f) em(g) Deus.(fe) (::) (Z) <c><sp>V/</sp></c> De(gf)mos(e) gra(ef)ças(g) ao(f) Se(g)nhor(hg) nos(fe)so(fg) Deus.(fgf) (::) <c><sp>R/</sp></c> É(f) no(f)sso(f) de(g)ver(h) e(g) nos(g)sa(f) sal(g)va(f)ção.(fe) (::) (Z)\n\n<c><sp>V/</sp></c> -Na:(f) ver(h)d'a(h)de,(h) é(h) .dig(h)no(g) e(gf) jus(fg)to,(g) (;)\nNa(f) ver(h)da(h)de,(h) dig(h)no,(gf) jus(fg)to,(g) (;)\nNa(f) ver(h)da(h)de,(h) dig(h)no(h) e(h) jus(h)to(gf) é,(fg) (;)\nNa(g) ver(g)da(g)de,(g) é(g) dig(g)no(f) e(g) jus(h)to(g) (,)\n-Na:(g) ver(g)d'a(g)de,(g) é(g) .dig(fe)no(ef) e(g) jus(fg)to,(f) (:)(Z)\n\nPor(f) is(h)so,(h) na(h) ver(gf)da(fg)de,(g) (;)\n-Na:(f) ver(h)d'a(h)de,(h) é(h) .dig(h)no(g) e(gf) jus(fg)to,(g) (;)\nNa(g) ver(g)da(fgh)de(g) (,)\nNa(g) ver(g)da(g)de,(g) dig(g)no(g) e(fe) jus(ef)to(g) é,(fgf) (:)(Z)\n\nPor(f) is(f)so,(f) na(f) ver(f)da(ef)de,(f) (,)\nNa(f) ver(h)da(h)de,(h) dig(h)no,(gf) jus(fg)to,(g) (;)\ndig(fgh)no(g) (,)\n-Na:(g) ver(g)d'a(g)de,(g) é(g) .dig(fe)no(ef) e(g) jus(fg)to,(f) (::)"

		is.Equal(composedGABC, expectedGABC)
	})

	t.Run("exception attempt to apply 'last' melody to short phrase like 'Senhor nosso'", func(t *testing.T) {
		is := is.New(t)

		inputText := "Na verdade, é digno e justo,\n por Cristo,\n Senhor nosso."

		composedGABC, err := service.NewGabcGenAPI(syllabifier).GeneratePreface(ctx, "", inputText)
		is.NoErr(err)

		expectedGABC := "<c><sp>V/</sp></c> O(e) Se(f)nhor(g) es(g)te(g)ja(e) con(f)vos(gf)co.(f) (::) <c><sp>R/</sp></c> E(e)<e>le</e> es(f)tá(g) no(g) me(g)io(e) de(f) nós.(gf) (::) (Z) <c><sp>V/</sp></c> Co(f)ra(g)ções(h) ao(g) al(fg)to.(fe) (::) <c><sp>R/</sp></c> O(g) nos(g)so(g) co(f)ra(g)cão(h) es(g)tá(f) em(g) Deus.(fe) (::) (Z) <c><sp>V/</sp></c> De(gf)mos(e) gra(ef)ças(g) ao(f) Se(g)nhor(hg) nos(fe)so(fg) Deus.(fgf) (::) <c><sp>R/</sp></c> É(f) no(f)sso(f) de(g)ver(h) e(g) nos(g)sa(f) sal(g)va(f)ção.(fe) (::) (Z)\n\n<c><sp>V/</sp></c> Na(f) ver(h)da(h)de,(h) é(h) dig(h)no(g) e(gf) jus(fg)to,(g) (;)\npor(g) Cris(fgh)to,(g) (,)\nSe(fe)nhor(efg) nos(fg)so.(f) (::)"

		is.Equal(norm.NFC.String(composedGABC), norm.NFC.String(expectedGABC))
	})
}
