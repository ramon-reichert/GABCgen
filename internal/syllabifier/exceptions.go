package syllabifier

var exceptions = map[string][]string{
	"alegria":      {"a", "le", "gri", "a"},
	"português":    {"por", "tu", "guês"},
	"paraguaio":    {"pa", "ra", "gua", "io"},
	"enxaguei":     {"en", "xa", "guei"},
	"sosseguei":    {"sos", "se", "guei"},
	"misericórdia": {"mi", "se", "ri", "cór", "di", "a"},
	"glória":       {"gló", "ri", "a"},
	"uruguaio":     {"u", "ru", "gua", "io"},
	"quais":        {"quais"},
	"iguais":       {"i", "guais"},
	"ameaçou":      {"a", "me", "a", "çou"},
}
