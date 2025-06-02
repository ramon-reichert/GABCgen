package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type SyllableInfo struct {
	Slashed    string `json:"slashed"`
	TonicIndex int    `json:"tonic_index"`
}

var userSyllabs map[string]SyllableInfo

func main() {
	var err error
	userSyllabs, err = LoadSyllables("user_syllables.json")
	if err != nil {
		log.Println("loading user syllables file:", err)
		return
	}

	word := "abraçado"
	sInfo, err := Syllabify(word)
	if err != nil {
		log.Printf("syllabifying word %v: %v", word, err)
		return
	}

	fmt.Println(sInfo)
}

func Syllabify(word string) (SyllableInfo, error) {

	//TODO: check if the word is already syllabified in the embedded json database

	//TODO: check if the word is already syllabified in the user database of new words
	if info, ok := userSyllabs[word]; ok {
		return info, nil
	}

	//if the word is not found in the databases, fetch it from a external website:
	info, err := fetchSyllabs(word)
	if err != nil {
		return SyllableInfo{}, fmt.Errorf("syllabifying new word: %w", err)
	}

	//add the word to the user database of new words
	userSyllabs[word] = info

	//save the user syllables to the file TODO: make it at once with all new words
	if SaveSyllables("user_syllables.json", userSyllabs) != nil {
		return SyllableInfo{}, fmt.Errorf("saving user syllables: %w", err)
	}

	return info, nil
}

func LoadSyllables(jsonFileName string) (map[string]SyllableInfo, error) {
	data, err := os.ReadFile(jsonFileName)
	if err != nil {
		return nil, err
	}

	var infos map[string]SyllableInfo
	err = json.Unmarshal(data, &infos)
	if err != nil {
		return nil, err
	}
	return infos, nil
}

func SaveSyllables(jsonFileName string, syllabs map[string]SyllableInfo) error {
	data, err := json.MarshalIndent(syllabs, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling syllables to JSON: %w", err)
	}

	err = os.WriteFile(jsonFileName, data, 0644)
	if err != nil {
		return fmt.Errorf("writing syllables to file %s: %w", jsonFileName, err)
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
