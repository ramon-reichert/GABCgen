// Package phrases handle musical phrases composed of words.Syllable structs.
// Phrases can be typed according to the Mass part.
package phrases

import (
	"fmt"
	"strings"

	gabcErrors "github.com/ramon-reichert/GABCgen/internal/platform/errors"
)

type Paragraph struct {
	Phrases []*Phrase
}

// DistributeText takes a text with lines separated by new lines and distributes it into Paragraphs.
func DistributeText(linedText string) ([]Paragraph, error) {
	var newPhrases []*Phrase
	var paragraphs []Paragraph

	if linedText == "" {
		return nil, fmt.Errorf("distributing text to Paragraphs: %w", gabcErrors.ErrNoText)
	}

	p := 0

	for v := range strings.Lines(linedText) {
		text, _ := strings.CutSuffix(v, "\n")
		text = strings.TrimSpace(text)

		if text != "" {
			newPhrases = append(newPhrases, New(text))
		} else if newPhrases != nil {
			paragraphs = append(paragraphs, Paragraph{Phrases: newPhrases})
			p++
			newPhrases = nil
		}
	}

	if newPhrases != nil {
		paragraphs = append(paragraphs, Paragraph{Phrases: newPhrases})
	}

	if len(paragraphs) == 0 {
		return nil, fmt.Errorf("distributing text to new Phrases: %w", gabcErrors.ErrNoText)
	}

	return paragraphs, nil
}
