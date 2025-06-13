package gabcGen_test

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"github.com/ramon-reichert/GABCgen/cmd/internal/gabcGen"
	"github.com/ramon-reichert/GABCgen/cmd/internal/syllabification/siteSyllabifier"
	"golang.org/x/text/unicode/norm"
)

var ctx context.Context = context.Background()

func TestIntegrationGeneratePreface(t *testing.T) {
	is := is.New(t)
	// Initialize the syllabifier with the necessary files:
	syllabifier := siteSyllabifier.NewSyllabifier("../../syllable_databases/liturgical_syllables.json", "../../syllable_databases/user_syllables.json", "../../syllable_databases/not_syllabified.txt")
	is.NoErr(syllabifier.LoadSyllables())

	t.Run("generate preface Páscoa I", func(t *testing.T) {
		is := is.New(t)

		inputText := "Na verdade, é digno e justo,\n é nosso dever e salvação proclamar vossa glória, ó Pai, em todo tempo,\n mas, com maior júbilo, louvar-vos nesta noite,\n ( neste dia ou neste tempo ) porque Cristo, nossa Páscoa, foi imolado.\n\n É ele o verdadeiro Cordeiro, que tirou o pecado do mundo;\n morrendo, destruiu a nossa morte\n e, ressurgindo, restaurou a vida.\n\n Por isso,\n transbordando de alegria pascal, exulta a criação por toda a terra;\n também as Virtudes celestes e as Potestades angélicas proclamam um hino à vossa glória,\n cantando\n a uma só voz:"
		//log.Println("inputText: ", inputText)

		composedGABC, err := gabcGen.NewGabcGenAPI(syllabifier).GeneratePreface(ctx, inputText)
		is.NoErr(err)
		// ||<i><c>neste dia ou neste tempo</c></i>||(,)
		expectedGABC := `Na(f) ver(h)da(h)de,(h) é(h) dig(h)no(g) e(gf) jus(fg)to,(g) (;)
é(f) nos(h)so(h) de(h)ver(h) e(h) sal(h)va(h)ção(h) pro(h)cla(h)mar(h) vos(h)sa(h) gló(h)ria,(h) ó(h) Pai,(h) em(h) to(h)do(gf) tem(fg)po,(g) (;)
mas,(g) com(g) mai(g)or(g) jú(g)bi(g)lo,(g) lou(g)var(g)-vos(g) nes(f)ta(g) noi(h)te,(g) (,)
por(g)que(g) Cris(g)to,(g) nos(g)sa(g) Pás(g)coa,(g) foi(fe) i(ef)mo(g)la(fg)do.(f) (:)(Z)

É(f) e(h)le(h) o(h) ver(h)da(h)dei(h)ro(h) Cor(h)dei(h)ro,(h) que(h) ti(h)rou(h) o(h) pe(h)ca(h)do(g) do(gf) mun(fg)do;(g) (;)
mor(g)ren(g)do,(g) des(g)tru(g)iu(g) a(g) nos(f)sa(g) mor(h)te(g) (,)
e,(g) res(g)sur(g)gin(g)do,(g) res(g)tau(fe)rou(ef) a(g) vi(fg)da.(f) (:)(Z)

Por(f) is(ef)so,(f) (,)
trans(f)bor(h)dan(h)do(h) de(h) a(h)le(h)gri(h)a(h) pas(h)cal,(h) e(h)xul(h)ta(h) a(h) cri(h)a(h)ção(h) por(h) to(h)da(g) a(gf) ter(fg)ra;(g) (;)
tam(f)bém(h) as(h) Vir(h)tu(h)des(h) ce(h)les(h)tes(h) e(h) as(h) Po(h)tes(h)ta(h)des(h) an(h)gé(h)li(h)cas(h) pro(h)cla(h)mam(h) um(h) hi(h)no(h) à(h) vos(h)sa(gf) gló(fg)ria,(g) (;)
can(g)tan(fgh)do(g) (,)
a(g) u(fe)ma(ef) só(g) voz:(fgf) (::)`

		is.Equal(norm.NFC.String(composedGABC), norm.NFC.String(expectedGABC))
	})
}
