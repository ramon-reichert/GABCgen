package web_test

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/ramon-reichert/GABCgen/internal/platform/syllabification/siteSyllabifier"
	"github.com/ramon-reichert/GABCgen/internal/platform/web"
	"github.com/ramon-reichert/GABCgen/internal/service"
	dmp "github.com/sergi/go-diff/diffmatchpatch"
	"golang.org/x/text/unicode/norm"
)

func TestGeneratePreface(t *testing.T) {
	syllabifier := siteSyllabifier.NewSyllabifier("../../../assets/syllable_databases/liturgical_syllables.json", "../../../assets/syllable_databases/user_syllables.json", "../../../assets/syllable_databases/not_syllabified.txt")

	if err := syllabifier.LoadSyllables(); err != nil {
		log.Printf("loading syllables db files: %v", err)
	}

	gabc := service.NewGabcGenAPI(syllabifier /*, render*/)
	gabcHandler := web.NewGabcHandler(gabc, time.Duration(5*time.Second))

	mux := http.NewServeMux()
	mux.HandleFunc("/preface", gabcHandler.Preface)
	server := web.NewServer(web.ServerConfig{Port: 8080, DisableRateLimit: true}, mux)

	t.Run("generates preface Páscoa I without errors", func(t *testing.T) {
		is := is.New(t)

		prefaceEntry := `{
			"dialogue": "",
			"text": "Na verdade, é digno e justo,\n é nosso dever e salvação proclamar vossa glória, ó Pai, em todo tempo,\n mas, com maior júbilo, louvar-vos nesta noite, ( neste dia ou neste tempo )\n porque Cristo, nossa Páscoa, foi imolado.\n\n É ele o verdadeiro Cordeiro, que tirou o pecado do mundo;\n morrendo, destruiu a nossa morte\n e, ressurgindo, restaurou a vida.\n\n Por isso,\n transbordando de alegria pascal, exulta a criação por toda a terra;\n também as Virtudes celestes e as Potestades angélicas proclamam um hino à vossa glória,\n cantando\n a uma só voz:"
}`
		request, _ := http.NewRequest(http.MethodPost, "/preface", strings.NewReader(prefaceEntry))
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)
		body, _ := io.ReadAll(response.Result().Body)
		body = body[:len(body)-1] // remove the last newline character

		is.True(response.Result().StatusCode == 200) // 200 OK

		expectedComposedGabcFields := `"\u003cc\u003e\u003csp\u003eV/\u003c/sp\u003e\u003c/c\u003e O(f) Se(g)nhor(h) es(h)te(h)ja(f) con(g)vos(hg)co.(g) (::) \u003cc\u003e\u003csp\u003eR/\u003c/sp\u003e\u003c/c\u003e E(f)\u003ce\u003ele\u003c/e\u003e es(g)tá(h) no(h) me(h)io(f) de(g) nós.(hg) (::) (Z) \u003cc\u003e\u003csp\u003eV/\u003c/sp\u003e\u003c/c\u003e Co(g)ra(h)ções(i) ao(h) al(gh)to.(gf) (::) \u003cc\u003e\u003csp\u003eR/\u003c/sp\u003e\u003c/c\u003e O(h) nos(h)so(h) co(g)ra(h)cão(i) es(h)tá(g) em(h) Deus.(gf) (::) (Z) \u003cc\u003e\u003csp\u003eV/\u003c/sp\u003e\u003c/c\u003e De(hg)mos(f) gra(fg)ças(h) ao(g) Se(h)nhor(ih) nos(gf)so(gh) Deus.(ghg) (::) \u003cc\u003e\u003csp\u003eR/\u003c/sp\u003e\u003c/c\u003e É(g) no(g)sso(g) de(h)ver(i) e(h) nos(h)sa(g) sal(h)va(g)ção.(gf) (::) (Z)\n\n\u003cc\u003e\u003csp\u003eV/\u003c/sp\u003e\u003c/c\u003e Na(f) ver(h)da(h)de,(h) é(h) dig(h)no(g) e(gf) jus(fg)to,(g) (;)\né(f) nos(h)so(h) de(h)ver(h) e(h) sal(h)va(h)ção(h) pro(h)cla(h)mar(h) vos(h)sa(h) gló(h)ria,(h) ó(h) Pai,(h) em(h) to(h)do(gf) tem(fg)po,(g) (;)\nmas,(g) com(g) mai(g)or(g) jú(g)bi(g)lo,(g) lou(g)var(g)-vos(g) nes(f)ta(g) noi(h)te,(g) ||\u003ci\u003e\u003cc\u003e neste dia ou neste tempo \u003c/c\u003e\u003c/i\u003e||(,)\npor(g)que(g) Cris(g)to,(g) nos(g)sa(g) Pás(g)coa,(g) foi(fe) i(ef)mo(g)la(fg)do.(f) (:)(Z)\n\nÉ(f) e(h)le(h) o(h) ver(h)da(h)dei(h)ro(h) Cor(h)dei(h)ro,(h) que(h) ti(h)rou(h) o(h) pe(h)ca(h)do(g) do(gf) mun(fg)do;(g) (;)\nmor(g)ren(g)do,(g) des(g)tru(g)iu(g) a(g) nos(f)sa(g) mor(h)te(g) (,)\ne,(g) res(g)sur(g)gin(g)do,(g) res(g)tau(fe)rou(ef) a(g) vi(fg)da.(f) (:)(Z)\n\nPor(f) is(ef)so,(f) (,)\ntrans(f)bor(h)dan(h)do(h) de(h) a(h)le(h)gri(h)a(h) pas(h)cal,(h) e(h)xul(h)ta(h) a(h) cri(h)a(h)ção(h) por(h) to(h)da(g) a(gf) ter(fg)ra;(g) (;)\ntam(f)bém(h) as(h) Vir(h)tu(h)des(h) ce(h)les(h)tes(h) e(h) as(h) Po(h)tes(h)ta(h)des(h) an(h)gé(h)li(h)cas(h) pro(h)cla(h)mam(h) um(h) hi(h)no(h) à(h) vos(h)sa(gf) gló(fg)ria,(g) (;)\ncan(g)tan(fgh)do(g) (,)\na(g) u(fe)ma(ef) só(g) voz:(fgf) (::)"`
		expectedJSONresponse := `{"gabc":` + expectedComposedGabcFields + `}`

		diffTool := dmp.New()
		diffs := diffTool.DiffMainRunes([]rune(norm.NFC.String(string(body))), []rune(norm.NFC.String(expectedJSONresponse)), false)
		if !(len(diffs) == 1 && diffs[0].Type == dmp.DiffEqual) {
			log.Println("\n\ndiffs: ", diffTool.DiffPrettyText(diffs))
		}

		is.Equal((norm.NFC.String(string(body))), (norm.NFC.String(expectedJSONresponse)))
	})

	t.Run("generates preface Páscoa I with Regional dialogue, without errors", func(t *testing.T) {
		is := is.New(t)

		prefaceEntry := `{
	   			"dialogue": "regional",
	   			"text": "Na verdade, é digno e justo,\n é nosso dever e salvação proclamar vossa glória, ó Pai, em todo tempo,\n mas, com maior júbilo, louvar-vos nesta noite, ( neste dia ou neste tempo )\n porque Cristo, nossa Páscoa, foi imolado.\n\n É ele o verdadeiro Cordeiro, que tirou o pecado do mundo;\n morrendo, destruiu a nossa morte\n e, ressurgindo, restaurou a vida.\n\n Por isso,\n transbordando de alegria pascal, exulta a criação por toda a terra;\n também as Virtudes celestes e as Potestades angélicas proclamam um hino à vossa glória,\n cantando\n a uma só voz:"
	   }`
		request, _ := http.NewRequest(http.MethodPost, "/preface", strings.NewReader(prefaceEntry))
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)
		body, _ := io.ReadAll(response.Result().Body)
		body = body[:len(body)-1] // remove the last newline character

		is.True(response.Result().StatusCode == 200) // 200 OK

		expectedComposedGabcFields := `"\u003cc\u003e\u003csp\u003eV/\u003c/sp\u003e\u003c/c\u003e O(h) Se(h)nhor(h) es(h)te(h)ja(f) con(h)vos(h)co.(h) (::) \u003cc\u003e\u003csp\u003eR/\u003c/sp\u003e\u003c/c\u003e E(h)\u003ce\u003ele\u003c/e\u003e es(h)tá(h) no(h) me(h)io(f) de(h) nós.(h) (::) (Z) \u003cc\u003e\u003csp\u003eV/\u003c/sp\u003e\u003c/c\u003e Co(h)ra(h)ções(f) ao(h) al(h)to.(h) (::) \u003cc\u003e\u003csp\u003eR/\u003c/sp\u003e\u003c/c\u003e O(h) nos(h)so(h) co(h)ra(h)cão(h) es(h)tá(f) em(h) Deus.(h) (::) (Z) \u003cc\u003e\u003csp\u003eV/\u003c/sp\u003e\u003c/c\u003e De(h)mos(h) gra(h)ças(h) ao(h) Se(h)nhor(h) nos(f)so(h) Deus.(h) (::) \u003cc\u003e\u003csp\u003eR/\u003c/sp\u003e\u003c/c\u003e É(h) nos(h)so(h) de(h)ver(h) e(h) nos(h)sa(h) sal(f)va(h)ção.(h) (::) (Z)\n\n\u003cc\u003e\u003csp\u003eV/\u003c/sp\u003e\u003c/c\u003e Na(f) ver(h)da(h)de,(h) é(h) dig(h)no(g) e(gf) jus(fg)to,(g) (;)\né(f) nos(h)so(h) de(h)ver(h) e(h) sal(h)va(h)ção(h) pro(h)cla(h)mar(h) vos(h)sa(h) gló(h)ria,(h) ó(h) Pai,(h) em(h) to(h)do(gf) tem(fg)po,(g) (;)\nmas,(g) com(g) mai(g)or(g) jú(g)bi(g)lo,(g) lou(g)var(g)-vos(g) nes(f)ta(g) noi(h)te,(g) ||\u003ci\u003e\u003cc\u003e neste dia ou neste tempo \u003c/c\u003e\u003c/i\u003e||(,)\npor(g)que(g) Cris(g)to,(g) nos(g)sa(g) Pás(g)coa,(g) foi(fe) i(ef)mo(g)la(fg)do.(f) (:)(Z)\n\nÉ(f) e(h)le(h) o(h) ver(h)da(h)dei(h)ro(h) Cor(h)dei(h)ro,(h) que(h) ti(h)rou(h) o(h) pe(h)ca(h)do(g) do(gf) mun(fg)do;(g) (;)\nmor(g)ren(g)do,(g) des(g)tru(g)iu(g) a(g) nos(f)sa(g) mor(h)te(g) (,)\ne,(g) res(g)sur(g)gin(g)do,(g) res(g)tau(fe)rou(ef) a(g) vi(fg)da.(f) (:)(Z)\n\nPor(f) is(ef)so,(f) (,)\ntrans(f)bor(h)dan(h)do(h) de(h) a(h)le(h)gri(h)a(h) pas(h)cal,(h) e(h)xul(h)ta(h) a(h) cri(h)a(h)ção(h) por(h) to(h)da(g) a(gf) ter(fg)ra;(g) (;)\ntam(f)bém(h) as(h) Vir(h)tu(h)des(h) ce(h)les(h)tes(h) e(h) as(h) Po(h)tes(h)ta(h)des(h) an(h)gé(h)li(h)cas(h) pro(h)cla(h)mam(h) um(h) hi(h)no(h) à(h) vos(h)sa(gf) gló(fg)ria,(g) (;)\ncan(g)tan(fgh)do(g) (,)\na(g) u(fe)ma(ef) só(g) voz:(fgf) (::)"`
		expectedJSONresponse := `{"gabc":` + expectedComposedGabcFields + `}`

		diffTool := dmp.New()
		diffs := diffTool.DiffMainRunes([]rune(norm.NFC.String(string(body))), []rune(norm.NFC.String(expectedJSONresponse)), false)
		if !(len(diffs) == 1 && diffs[0].Type == dmp.DiffEqual) {
			log.Println("\n\ndiffs: ", diffTool.DiffPrettyText(diffs))
		}

		is.Equal((norm.NFC.String(string(body))), (norm.NFC.String(expectedJSONresponse)))
	})

	t.Run("returns invalid json request error", func(t *testing.T) {
		is := is.New(t)

		prefaceEntry := `{
			"dialogue": ""
			"text": "missing coma after field dialogue"
}`
		expectedJSONresponse := `invalid JSON body: invalid character '"' after object key:value pair`
		request, _ := http.NewRequest(http.MethodPost, "/preface", strings.NewReader(prefaceEntry))
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)
		body, _ := io.ReadAll(response.Result().Body)
		body = body[:len(body)-1] // remove the last newline character

		is.True(response.Result().StatusCode == 400) // 400 Bad Request

		diffTool := dmp.New()
		diffs := diffTool.DiffMainRunes([]rune(norm.NFC.String(string(body))), []rune(norm.NFC.String(expectedJSONresponse)), false)
		if !(len(diffs) == 1 && diffs[0].Type == dmp.DiffEqual) {
			log.Println("\n\ndiffs: ", diffTool.DiffPrettyText(diffs))
		}

		is.Equal((norm.NFC.String(string(body))), (norm.NFC.String(expectedJSONresponse)))
	})

	t.Run("returns error from GabcGen", func(t *testing.T) {
		is := is.New(t)

		prefaceEntry := `{
			"dialogue": "",
			"text": "just one line of text, should return ErrShortParagraph"
}`
		expectedJSONresponse := `generating Preface: typing phrase: just one line of text, should return ErrShortParagraph - each paragraph must have at least three phrases, not counting the conclusion phrase - which can start the last paragraph`
		request, _ := http.NewRequest(http.MethodPost, "/preface", strings.NewReader(prefaceEntry))
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)
		body, _ := io.ReadAll(response.Result().Body)
		body = body[:len(body)-1] // remove the last newline character

		is.True(response.Result().StatusCode == 400) // 400 Bad Request

		diffTool := dmp.New()
		diffs := diffTool.DiffMainRunes([]rune(norm.NFC.String(string(body))), []rune(norm.NFC.String(expectedJSONresponse)), false)
		if !(len(diffs) == 1 && diffs[0].Type == dmp.DiffEqual) {
			log.Println("\n\ndiffs: ", diffTool.DiffPrettyText(diffs))
		}

		is.Equal((norm.NFC.String(string(body))), (norm.NFC.String(expectedJSONresponse)))
	})
}
