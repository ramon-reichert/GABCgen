package syllabifier

import (
	"context"
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	// ErrEmptyInput is returned when input is empty or whitespace
	ErrEmptyInput = errors.New("input word is empty")

	vowels = "aeiouáéíóúâêôãõàüAEIOUÁÉÍÓÚÂÊÔÃÕÀÜ"

	diphthongs = map[string]bool{
		"ai": true, "au": true, "ei": true, "eu": true,
		"oi": true, "ou": true, "ui": true, "iu": true,
		"ão": true, "éu": true, "êu": true, "uê": true, "ué": true,
	}

	triphthongs = map[string]bool{
		"uai": true, "uão": true, "uei": true, "uõe": true,
	}

	validOnsets = map[string]bool{
		"pr": true, "pl": true, "br": true, "bl": true,
		"cr": true, "cl": true, "dr": true, "fr": true,
		"gr": true, "gl": true, "tr": true, "vr": true,
		"fl": true, "ch": true, "lh": true, "nh": true,
		"qu": true, "gu": true,
	}
)

func isVowel(r rune) bool {
	return strings.ContainsRune(vowels, r)
}

// Syllabify splits a Portuguese word into syllables with basic diphthong/triphthong support and exception handling.
func Syllabify(ctx context.Context, word string) ([]string, error) {
	if ctx == nil {
		return nil, errors.New("context is nil")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	trimmed := strings.TrimSpace(word)
	if trimmed == "" {
		return nil, ErrEmptyInput
	}

	// Check exceptions
	if syl, ok := exceptions[strings.ToLower(trimmed)]; ok {
		if unicode.IsUpper([]rune(trimmed)[0]) {
			r, size := utf8.DecodeRuneInString(syl[0])
			syl[0] = string(unicode.ToUpper(r)) + syl[0][size:]
		}
		return syl, nil
	}

	// Handle hyphenated words
	if strings.Contains(trimmed, "-") {
		parts := strings.Split(trimmed, "-")
		var result []string
		for i, part := range parts {
			if part == "" {
				continue
			}
			sylls, err := Syllabify(ctx, part)
			if err != nil {
				return nil, err
			}
			result = append(result, sylls...)
			if i < len(parts)-1 {
				result = append(result, "-")
			}
		}
		return result, nil
	}

	original := trimmed
	lower := strings.ToLower(trimmed)
	runes := []rune(lower)
	origRunes := []rune(original)
	n := len(runes)
	var syllables []string
	start := 0

	for i := 0; i < n; i++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if !isVowel(runes[i]) {
			continue
		}

		// Triphthong
		if i+2 < n {
			tri := string(runes[i : i+3])
			if triphthongs[tri] {
				syllables = append(syllables, string(origRunes[start:i+3]))
				start = i + 3
				i += 2
				continue
			}
		}

		// Diphthong
		if i+1 < n {
			di := string(runes[i : i+2])
			if diphthongs[di] {
				syllables = append(syllables, string(origRunes[start:i+2]))
				start = i + 2
				i += 1
				continue
			}
		}

		// Single vowel
		next := i + 1
		if next < n-1 && !isVowel(runes[next]) && isVowel(runes[next+1]) {
			// V-C-V: split before C
			syllables = append(syllables, string(origRunes[start:next]))
			start = next
			i = next - 1
		} else if next < n-2 && !isVowel(runes[next]) && !isVowel(runes[next+1]) {
			cluster := string(runes[next : next+2])
			if validOnsets[cluster] {
				syllables = append(syllables, string(origRunes[start:next]))
				start = next
				i = next - 1
			} else {
				syllables = append(syllables, string(origRunes[start:next+1]))
				start = next + 1
				i = next
			}
		}
	}

	if start < n {
		syllables = append(syllables, string(origRunes[start:n]))
	}

	// Restore first letter casing
	if unicode.IsUpper(origRunes[0]) && len(syllables) > 0 {
		firstRune, size := utf8.DecodeRuneInString(syllables[0])
		syllables[0] = string(unicode.ToUpper(firstRune)) + syllables[0][size:]
	}

	return syllables, nil
}
