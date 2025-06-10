package gabcGen_test

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"github.com/ramon-reichert/GABCgen/cmd/internal/gabcGen"
	"github.com/ramon-reichert/GABCgen/cmd/internal/syllabification/siteSyllabifier"
)

var ctx context.Context = context.Background()

func TestIntegrationGeneratePreface(t *testing.T) {
	is := is.New(t)
	// Initialize the syllabifier with the necessary files:
	syllabifier := siteSyllabifier.NewSyllabifier("../../syllable_databases/liturgical_syllables.json", "../../syllable_databases/user_syllables.json", "../../syllable_databases/not_syllabified.txt")
	is.NoErr(syllabifier.LoadSyllables())

	t.Run("generate preface Páscoa I", func(t *testing.T) {
		is := is.New(t)

		inputText := "Na verdade, é digno e justo,\n é nosso dever e salvação proclamar vossa glória, ó Pai, em todo tempo,\n mas, com maior júbilo, louvar-vos nesta noite\n ( neste dia ou neste tempo ) , porque Cristo, nossa Páscoa, foi imolado.\n\n É ele o verdadeiro Cordeiro, que tirou o pecado do mundo;\n morrendo, destruiu a nossa morte\n e, ressurgindo, restaurou a vida.\n\n Por isso,\n transbordando de alegria pascal, exulta a criação por toda a terra;\n também as Virtudes celestes e as Potestades angélicas proclamam um hino à vossa glória,\n cantando\n a uma só voz:"
		//log.Println("inputText: ", inputText)

		composedGABC, err := gabcGen.NewGabcGenAPI(syllabifier).GeneratePreface(ctx, inputText)
		is.NoErr(err)

		expectedGABC := `Na(f) verdade, é(hr0)  dig(h)no(g) e(gf) jus(fg)to,(g) (;)
é(f) nosso dever e salvação proclamar vossa glória, ó Pai, em(hr0) to(h)do(gf) tem(fg)po,(g) (;)
mas, com maior júbilo, louvar(gr0)-vos(g) nes(f)ta(g) noi(h)te,(g) ||<i><c>neste dia ou neste tempo</c></i>|| (,) por(g)que(g) Cristo nossa(gr0) Pás(g)coa(g) foi(fe) i(ef)mo(g)la(fg)do.(f) (:) (Z)

É(f) e(h)le(h) o(h) ver(h)da(h)dei(h)ro(h) Cor(h)dei(h)ro(h) que(h) ti(h)rou(h) o(h) pe(h)ca(h)do(g) do(gf) mun(fg)do;(g)(;) 
mor(g)ren(g)do,(g) des(g)tru(g)iu(g) a(g) nos(f)sa(g) mor(h)te,(g) (,) e(g) res(g)sur(g)gin(g)do(g) res(g)tau(fe)rou(ef) a(g) vi(fg)da.(f) (:) (Z)

Por(f) is(ef)so(f) (,) trans(f)bordando de alegria pascal, exulta a criação por(hr0) to(h)da(h) a(gf) ter(fg)ra;(g) (:) (z) tam(f)bém as Virtudes celestes e as Potestades angélicas proclamam um hino(hr0) à(h) vo(h)ssa(gf) gló(fg)ria(g) (;) (z) can(g)tan(f!gwh)do(g) a(g) u(fe)ma(ef) só(g) voz:(fgf) (::)`

		is.Equal(composedGABC, expectedGABC)
	})
}
