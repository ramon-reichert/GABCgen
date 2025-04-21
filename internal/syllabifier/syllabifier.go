package syllabifier

import (
	"context"
)

func Syllabify(ctx context.Context, word string) (string, int, error) {
	var hyphen string
	var tonic int

	//Mocking syllabification to allow testing the core application:
	switch word { //"Na verdade, é digno e justo,="
	case "na":
		hyphen = "na"
		tonic = 1
	case "verdade":
		hyphen = "ver/da/de"
		tonic = 2
	case "é":
		hyphen = "é"
		tonic = 1
	case "digno":
		hyphen = "dig/no"
		tonic = 1
	case "e":
		hyphen = "e"
		tonic = 1
	case "justo":
		hyphen = "jus/to"
		tonic = 1
	}

	//TODO: fetch the word in a list of already used words. Could be a map["palavra"]Syllab{hyphen: "pa-la-vra", tonic: "la"}
	//TODO: if it is not there, ask it to an AI API or another solution

	return hyphen, tonic, nil
}
