package syllabifier

import (
	"context"
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	ErrEmptyInput = errors.New("input word is empty")
	vowels        = "aeiouáéíóúâêôãõàüAEIOUÁÉÍÓÚÂÊÔÃÕÀÜ"
	diphthongs    = map[string]bool{
		"ai": true, "au": true, "ei": true, "eu": true,
		"oi": true, "ou": true, "ui": true, "iu": true,
		"ão": true, "éu": true, "êu": true, "uê": true, "ué": true,
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

// Syllabify splits a Portuguese word into its syllables.
func Syllabify(ctx context.Context, word string) ([]string, error) {
	// 1) Context check
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

	lower := strings.ToLower(trimmed)
	origRunes := []rune(trimmed)
	runes := []rune(lower)
	n := len(runes)

	var syllables []string
	start := 0

	for i := 0; i < n; i++ {
		// cancellation inside loop
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		if !isVowel(runes[i]) {
			continue
		}

		// 2a) Diphthong or hiatus
		if i+1 < n && isVowel(runes[i+1]) {
			pair := string(runes[i : i+2])
			if diphthongs[pair] {
				i++ // keep diphthong together
			} else {
				// hiatus: split right after this vowel
				syllables = append(syllables, string(origRunes[start:i+1]))
				start = i + 1
				continue
			}
		}

		// 2b) Determine split position
		bound := i + 1
		if bound < n-1 && !isVowel(runes[bound]) && isVowel(runes[bound+1]) {
			// V‑C‑V → split before C: bound stays
		} else if bound < n-2 && !isVowel(runes[bound]) && !isVowel(runes[bound+1]) {
			cluster := string(runes[bound : bound+2])
			if !validOnsets[cluster] {
				// VC‑CV → split after first C
				bound++
			}
		} else {
			// no boundary here → continue scanning
			continue
		}

		// 3) Emit syllable
		syllables = append(syllables, string(origRunes[start:bound]))
		start = bound
	}

	// 4) Append final part
	if start < n {
		syllables = append(syllables, string(origRunes[start:n]))
	}

	// 5) Restore capitalization on first syllable
	if unicode.IsUpper(origRunes[0]) && len(syllables) > 0 {
		firstRune, size := utf8.DecodeRuneInString(syllables[0])
		syllables[0] = string(unicode.ToUpper(firstRune)) + syllables[0][size:]
	}

	return syllables, nil
}
