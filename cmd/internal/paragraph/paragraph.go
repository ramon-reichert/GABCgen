package paragraph

import (
	"fmt"
	"log"
	"strings"

	"github.com/ramon-reichert/GABCgen/cmd/internal/gabcErrors"
	"github.com/ramon-reichert/GABCgen/cmd/internal/phrases"
)

type Paragraph struct {
	Phrases []*phrases.Phrase
}

func DistrbuteText(linedText string) ([]Paragraph, error) {
	var newPhrases []*phrases.Phrase
	var paragraphs []Paragraph

	if linedText == "" {
		return nil, fmt.Errorf("distributing text to Paragraphs: %w", gabcErrors.ErrNoText)
	}

	p := 0

	for v := range strings.Lines(linedText) {

		text, _ := strings.CutSuffix(v, "\n")
		text = strings.TrimSpace(text)

		if text != "" {

			//debug code
			fmt.Println("line: ", text)

			newPhrases = append(newPhrases, phrases.New(text))
		} else if newPhrases != nil {
			paragraphs = append(paragraphs, Paragraph{Phrases: newPhrases})
			p++
			newPhrases = nil
		}
	}
	if newPhrases != nil {
		paragraphs = append(paragraphs, Paragraph{Phrases: newPhrases})
	}

	//debug code
	log.Println("len(paragraphs): ", len(paragraphs))

	if len(paragraphs) == 0 {
		return nil, fmt.Errorf("distributing text to new Phrases: %w", gabcErrors.ErrNoText)
	}

	return paragraphs, nil
}
