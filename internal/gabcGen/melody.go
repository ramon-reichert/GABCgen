package gabcGen

import (
	"context"
	"fmt"
)

func (ph Phrase) ApplyMelodyGABC(ctx context.Context) (string, error) {
	switch ph.phraseType {
	case "firsts":
		composedPhrase, err := applyFirsts(ph)
		if err != nil {
			return "", fmt.Errorf("error applying firsts phrase: %w ", err)
		}
		end := "(;)"
		return joinSyllables(composedPhrase, end), nil

	case "last":
		/*	for i := len(ph.Syllables) - 1; ph.Syllables[i].IsTonic; i-- { //reading Syllables from the end
				if ph.Syllables[i].IsLast && ph.Syllables[i].IsFirst { //it means it's an oxytone
					ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(fgf)"
				} else {
					ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(fgf)"
				}
			}
		*/
	}
	return "", fmt.Errorf("someerror") //TODO: HANDLE ERROR CASES
}

func joinSyllables(ph Phrase, end string) string {
	var result string
	for _, v := range ph.Syllables {
		result = result + v.GABC
		if v.IsLast {
			result = result + " "
		}
	}

	return result + end
}

func applyFirsts(ph Phrase) (Phrase, error) {
	i := len(ph.Syllables) - 1 //reading Syllables from the end:
	if i < 0 {
		return Phrase{}, ErrResponseToShort
	}

	//last unstressed Syllables:
	for !ph.Syllables[i].IsTonic {
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(g)"
		i--
		if i < 0 {
			return Phrase{}, ErrResponseToShort
		}
	}

	//last tonic syllable:
	ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(fg)"
	i--
	if i < 0 {
		return Phrase{}, ErrResponseToShort
	}

	//syllable before the last tonic:
	ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(gf)"
	i--
	if i < 0 {
		return Phrase{}, ErrResponseToShort
	}

	//testing the exception at last unstressed reciting syllable:
	if ph.Syllables[i].IsTonic { //default case
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(h)"
		i--
		if i < 0 {
			return Phrase{}, ErrResponseToShort
		}
	} else if ph.Syllables[i-1].IsTonic && !ph.Syllables[i-1].IsLast { //exception case
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(g)"
		i--
		if i < 0 {
			return Phrase{}, ErrResponseToShort
		}
	}

	// completing reciting Syllables:
	for i > 0 {
		ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(h)"
		i--
		if i < 0 {
			return Phrase{}, ErrResponseToShort
		}
	}

	//first intonation syllable:
	ph.Syllables[i].GABC = string(ph.Syllables[i].Char) + "(f)"

	return ph, nil
}
