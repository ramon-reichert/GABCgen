// Package preface handles specific phrase types that compose the melody of the Preface of Mass.
package preface

type Dialogue string

const (
	// There is a very deep-rooted way of singing the preface dialogue throughout Brazil, which we call "regional", which actually follows the melody of the last blessing of the Mass.
	// Some priests prefer sing that regional tone because it is much easier for the people to answer it.
	Regional Dialogue = "<c><sp>V/</sp></c> O(h) Se(h)nhor(h) es(h)te(h)ja(f) con(h)vos(h)co.(h) (::) <c><sp>R/</sp></c> E(h)<e>le</e> es(h)tá(h) no(h) me(h)io(f) de(h) nós.(h) (::) (Z) <c><sp>V/</sp></c> Co(h)ra(h)ções(f) ao(h) al(h)to.(h) (::) <c><sp>R/</sp></c> O(h) nos(h)so(h) co(h)ra(h)cão(h) es(h)tá(f) em(h) Deus.(h) (::) (Z) <c><sp>V/</sp></c> De(h)mos(h) gra(h)ças(h) ao(h) Se(h)nhor(h) nos(f)so(h) Deus.(h) (::) <c><sp>R/</sp></c> É(h) nos(h)so(h) de(h)ver(h) e(h) nos(h)sa(h) sal(f)va(h)ção.(h) (::) (Z)"

	// Solemn mode is the official one that comes in the last brazilian missal.
	Solemn Dialogue = "<c><sp>V/</sp></c> O(f) Se(g)nhor(h) es(h)te(h)ja(f) con(g)vos(hg)co.(g) (::) <c><sp>R/</sp></c> E(f)<e>le</e> es(g)tá(h) no(h) me(h)io(f) de(g) nós.(hg) (::) (Z) <c><sp>V/</sp></c> Co(g)ra(h)ções(i) ao(h) al(gh)to.(gf) (::) <c><sp>R/</sp></c> O(h) nos(h)so(h) co(g)ra(h)cão(i) es(h)tá(g) em(h) Deus.(gf) (::) (Z) <c><sp>V/</sp></c> De(hg)mos(f) gra(fg)ças(h) ao(g) Se(h)nhor(ih) nos(gf)so(gh) Deus.(ghg) (::) <c><sp>R/</sp></c> É(g) no(g)sso(g) de(h)ver(i) e(h) nos(h)sa(g) sal(h)va(g)ção.(gf) (::) (Z)"
)

// SetDialogue sets the dialogue tone according to the given string.
func SetDialogueTone(d string) Dialogue {
	switch d {
	case "regional":
		return Regional
	default:
		return Solemn
	}
}
