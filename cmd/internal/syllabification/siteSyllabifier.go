package syllabification

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type SiteSyllabifier struct {
	userSyllabs  map[string]SyllableInfo
	userFilePath string //path to the user syllables file
}

func NewSyllabifier(userSyllabsPath string) *SiteSyllabifier {
	return &SiteSyllabifier{
		userFilePath: userSyllabsPath,
	}
}

type SyllableInfo struct {
	Slashed    string `json:"slashed"`
	TonicIndex int    `json:"tonic_index"`
}

func (s *SiteSyllabifier) Syllabify(ctx context.Context, word string) (string, int, error) {

	//TODO: check if the word is already syllabified in the embedded json database of liturgical words

	//TODO: check if the word is already syllabified in the user database of new words
	if info, ok := s.userSyllabs[word]; ok {
		return info.Slashed, info.TonicIndex, nil
	}

	//if the word is not found in the databases, fetch it from a external website:
	info, err := fetchSyllabs(word)
	if err != nil {
		return "", 0, fmt.Errorf("syllabifying new word: %w", err)
	}

	//debug code:
	log.Printf("Syllabified new word %v: %v", word, info)
	log.Println("userSyllabs before adding new word: ", s.userSyllabs)
	log.Printf("adress of userSyllabs: %p", &s.userSyllabs)

	//add the word to the user database of new words
	s.userSyllabs[word] = info

	return info.Slashed, info.TonicIndex, nil
}

func (s *SiteSyllabifier) LoadSyllables() error {
	data, err := os.ReadFile(s.userFilePath)
	if err != nil {
		return err
	}

	if json.Unmarshal(data, &s.userSyllabs) != nil {
		return fmt.Errorf("unmarshaling file %v: %w", s.userFilePath, err)
	}

	//debug code:
	log.Println("Loaded syllables from userSyllabs: ", s.userSyllabs)
	log.Printf("adress of userSyllabs: %p", &s.userSyllabs)

	return nil
}

func (s *SiteSyllabifier) SaveSyllables() error {
	data, err := json.MarshalIndent(s.userSyllabs, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling syllables to JSON: %w", err)
	}

	err = os.WriteFile(s.userFilePath, data, 0644)
	if err != nil {
		return fmt.Errorf("writing syllables to file %s: %w", s.userFilePath, err)
	}
	return nil
}

// fetchSyllabs fetches the syllables of a word from separaremsilabas.com
func fetchSyllabs(word string) (SyllableInfo, error) {
	// Fetch the HTML
	resp, err := http.Get("https://www.separaremsilabas.com/index.php?lang=index.php&p=" + word + "&button=Separa%C3%A7%C3%A3o+das+s%C3%ADlabas")
	if err != nil {
		return SyllableInfo{}, fmt.Errorf("fetching syllables of %v from separaremsilabas.com: %w", word, err)
	}
	defer resp.Body.Close()

	// Read the body
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

	//define the tonic syllable:
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

	//build slashed syllable string
	return SyllableInfo{Slashed: strings.Join(syllabs, "/"), TonicIndex: tonicIndex}, nil

}
