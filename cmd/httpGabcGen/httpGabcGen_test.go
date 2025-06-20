package httpGabcGen_test

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/ramon-reichert/GABCgen/cmd/httpGabcGen"
	"github.com/ramon-reichert/GABCgen/cmd/internal/gabcGen"
	"github.com/ramon-reichert/GABCgen/cmd/internal/syllabification/siteSyllabifier"
	"golang.org/x/text/unicode/norm"
)

func TestGeneratePreface(t *testing.T) {

	syllabifier := siteSyllabifier.NewSyllabifier("../syllable_databases/liturgical_syllables.json", "../syllable_databases/user_syllables.json", "../syllable_databases/not_syllabified.txt")
	err := syllabifier.LoadSyllables()
	if err != nil {
		log.Printf("loading syllables db files: %v", err)
	}
	gabc := gabcGen.NewGabcGenAPI(syllabifier /*, render*/)
	gabcHandler := httpGabcGen.NewGabcHandler(gabc, time.Duration(5*time.Second))
	server := httpGabcGen.NewServer(httpGabcGen.ServerConfig{Port: 8080}, gabcHandler)
	t.Run("generates preface Páscoa I without errors", func(t *testing.T) {
		is := is.New(t)

		prefaceEntry := `{
			"header": {},
			"dialogue": "",
			"text": "Na verdade, é digno e justo,\n é nosso dever e salvação proclamar vossa glória, ó Pai, em todo tempo,\n mas, com maior júbilo, louvar-vos nesta noite, ( neste dia ou neste tempo )\n porque Cristo, nossa Páscoa, foi imolado.\n\n É ele o verdadeiro Cordeiro, que tirou o pecado do mundo;\n morrendo, destruiu a nossa morte\n e, ressurgindo, restaurou a vida.\n\n Por isso,\n transbordando de alegria pascal, exulta a criação por toda a terra;\n também as Virtudes celestes e as Potestades angélicas proclamam um hino à vossa glória,\n cantando\n a uma só voz:"
}`

		expectedJSONresponse := `{"header":{},"dialogue":"","text":"Na(f) ver(h)da(h)de,(h) é(h) dig(h)no(g) e(gf) jus(fg)to,(g) (;)\né(f) nos(h)so(h) de(h)ver(h) e(h) sal(h)va(h)ção(h) pro(h)cla(h)mar(h) vos(h)sa(h) gló(h)ria,(h) ó(h) Pai,(h) em(h) to(h)do(gf) tem(fg)po,(g) (;)\nmas,(g) com(g) mai(g)or(g) jú(g)bi(g)lo,(g) lou(g)var(g)-vos(g) nes(f)ta(g) noi(h)te,(g) ||\u003ci\u003e\u003cc\u003e neste dia ou neste tempo \u003c/c\u003e\u003c/i\u003e||(,)\npor(g)que(g) Cris(g)to,(g) nos(g)sa(g) Pás(g)coa,(g) foi(fe) i(ef)mo(g)la(fg)do.(f) (:)(Z)\n\nÉ(f) e(h)le(h) o(h) ver(h)da(h)dei(h)ro(h) Cor(h)dei(h)ro,(h) que(h) ti(h)rou(h) o(h) pe(h)ca(h)do(g) do(gf) mun(fg)do;(g) (;)\nmor(g)ren(g)do,(g) des(g)tru(g)iu(g) a(g) nos(f)sa(g) mor(h)te(g) (,)\ne,(g) res(g)sur(g)gin(g)do,(g) res(g)tau(fe)rou(ef) a(g) vi(fg)da.(f) (:)(Z)\n\nPor(f) is(ef)so,(f) (,)\ntrans(f)bor(h)dan(h)do(h) de(h) a(h)le(h)gri(h)a(h) pas(h)cal,(h) e(h)xul(h)ta(h) a(h) cri(h)a(h)ção(h) por(h) to(h)da(g) a(gf) ter(fg)ra;(g) (;)\ntam(f)bém(h) as(h) Vir(h)tu(h)des(h) ce(h)les(h)tes(h) e(h) as(h) Po(h)tes(h)ta(h)des(h) an(h)gé(h)li(h)cas(h) pro(h)cla(h)mam(h) um(h) hi(h)no(h) à(h) vos(h)sa(gf) gló(fg)ria,(g) (;)\ncan(g)tan(fgh)do(g) (,)\na(g) u(fe)ma(ef) só(g) voz:(fgf) (::)"}`
		request, _ := http.NewRequest(http.MethodPost, "/preface", strings.NewReader(prefaceEntry))
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)
		body, _ := io.ReadAll(response.Result().Body)

		body = body[:len(body)-1] // remove the last newline character

		is.True(response.Result().StatusCode == 200)

		//	dmp := diffmatchpatch.New()
		//	diffs := dmp.DiffMainRunes([]rune(norm.NFC.String(string(body))), []rune(norm.NFC.String(expectedJSONresponse)), false)
		//	log.Println("\n\ndiffs: ", diffs)

		is.Equal((norm.NFC.String(string(body))), (norm.NFC.String(expectedJSONresponse)))

	})
}
