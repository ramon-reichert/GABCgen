// Package siteSyllabifier is an adapter that fetches syllables from an external website.
package siteSyllabifier

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type SiteSyllabifier struct {
	userSyllabs            map[string]SyllableInfo // map of user syllables, loaded from a file and new words added during runtime
	userFilePath           string                  // path to the user syllables file
	liturgicalSyllabs      map[string]SyllableInfo // map of liturgical syllables, loaded from a file
	liturgicalFilePath     string                  // path to the liturgical syllables file
	NotSyllabifiedWords    string                  // list of words that were not syllabified, to be saved to a file later
	notSyllabifiedFilePath string                  // path to the file where the not syllabified words will be saved
}

// NewSyllabifier creates a new SiteSyllabifier instance.
func NewSyllabifier(liturgicalSyllabsPath, userSyllabsPath, notSyllabifiedPath string) *SiteSyllabifier {
	return &SiteSyllabifier{
		userFilePath:           userSyllabsPath,
		liturgicalFilePath:     liturgicalSyllabsPath,
		notSyllabifiedFilePath: notSyllabifiedPath,
	}
}

type SyllableInfo struct {
	Slashed    string `json:"slashed"`
	TonicIndex int    `json:"tonic_index"`
}

// Syllabify syllabifies a word, first checking the user and liturgical databases, then fetching from an external website if not found.
func (s *SiteSyllabifier) Syllabify(ctx context.Context, word string) (string, int, error) {
	// Check if the word is already syllabified in the user database of new words
	if info, ok := s.userSyllabs[word]; ok {
		return info.Slashed, info.TonicIndex, nil
	}

	// Check if the word is already syllabified in the embedded json database of liturgical words
	if info, ok := s.liturgicalSyllabs[word]; ok {
		return info.Slashed, info.TonicIndex, nil
	}

	// Fetch the word from a external website, if not found in the databases
	info, err := fetchSyllabs(ctx, word)
	if err != nil {
		// Put the word into a list of non-syllabified words
		s.NotSyllabifiedWords = s.NotSyllabifiedWords + "\n" + word
		return "", 0, fmt.Errorf("syllabifying new word: %w", err)
	}

	// Add the word to the user database of new words
	s.userSyllabs[word] = info

	return info.Slashed, info.TonicIndex, nil
}

// LoadSyllables loads the syllables from the liturgical and user files.
func (s *SiteSyllabifier) LoadSyllables() error {
	dataL, err := os.ReadFile(s.liturgicalFilePath)
	if err != nil {
		return err
	}

	if json.Unmarshal(dataL, &s.liturgicalSyllabs) != nil {
		return fmt.Errorf("unmarshaling file %v: %w", s.liturgicalFilePath, err)
	}

	dataU, err := os.ReadFile(s.userFilePath)
	if err != nil {
		return err
	}

	if json.Unmarshal(dataU, &s.userSyllabs) != nil {
		return fmt.Errorf("unmarshaling file %v: %w", s.userFilePath, err)
	}

	dataNS, err := os.ReadFile(s.notSyllabifiedFilePath)
	if err != nil {
		return err
	}

	s.NotSyllabifiedWords = string(dataNS)

	return nil
}

// SaveSyllables saves the user syllables and the not syllabified words to their respective files.
func (s *SiteSyllabifier) SaveSyllables() error {
	data, err := json.MarshalIndent(s.userSyllabs, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling syllables to JSON: %w", err)
	}

	err = os.WriteFile(s.userFilePath, data, 0644)
	if err != nil {
		return fmt.Errorf("writing syllables to file %s: %w", s.userFilePath, err)
	}

	err = os.WriteFile(s.notSyllabifiedFilePath, []byte(s.NotSyllabifiedWords), 0644)
	if err != nil {
		return fmt.Errorf("writing syllables to file %s: %w", s.notSyllabifiedFilePath, err)
	}
	return nil
}

// fetchSyllabs fetches the syllables of a word from the external website.
func fetchSyllabs(ctx context.Context, word string) (SyllableInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.separaremsilabas.com/index.php?lang=index.php&p="+word+"&button=Separa%C3%A7%C3%A3o+das+s%C3%ADlabas", nil)
	if err != nil {
		return SyllableInfo{}, fmt.Errorf("fetching syllables of %v from separaremsilabas.com: %w", word, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return SyllableInfo{}, fmt.Errorf("fetching syllables of %v from separaremsilabas.com: %w", word, err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return SyllableInfo{}, fmt.Errorf("fetching syllables of %v from separaremsilabas.com: %w", word, err)
	}

	// Regex to extract content between matching tags
	re := regexp.MustCompile(`1\.9em">(.*?)</font>`)
	matches := re.FindStringSubmatch(string(body))
	if matches == nil {
		return SyllableInfo{}, fmt.Errorf("fetching syllables of %v from separaremsilabas.com: no syllables found", word)
	}

	// Define the tonic syllable:
	tonicIndex := 0
	syllabs := strings.Split(matches[1], "-")
	for i, s := range syllabs {
		if strings.HasPrefix(s, "<strong>") {
			s = strings.TrimPrefix(s, "<strong>")
			s = strings.TrimSuffix(s, "</strong>")
			tonicIndex = i + 1 // tonic syllable is 1-based index
			syllabs[i] = s     // replace the strong tag with the syllable
		}
	}

	if tonicIndex == 0 {
		return SyllableInfo{}, fmt.Errorf("unable to define tonic syllable in: %v", word)
	}

	// Build slashed syllable string
	return SyllableInfo{Slashed: strings.Join(syllabs, "/"), TonicIndex: tonicIndex}, nil
}
