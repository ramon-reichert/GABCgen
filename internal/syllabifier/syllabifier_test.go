package syllabifier_test

import (
	"context"
	"testing"

	syllabifier "github.com/ramon-reichert/GABCgen/internal/draft"
)

func TestSyllabify_BasicWords(t *testing.T) {
	tests := []struct {
		word     string
		expected []string
	}{
		{"alegria", []string{"a", "le", "gri", "a"}},
		{"oração", []string{"o", "ra", "ção"}},
		{"sagrado", []string{"sa", "gra", "do"}},
		{"misericórdia", []string{"mi", "se", "ri", "cór", "di", "a"}},
		{"glória", []string{"gló", "ri", "a"}},
		{"português", []string{"por", "tu", "guês"}},
		{"Alegria", []string{"A", "le", "gri", "a"}},
	}

	for _, tc := range tests {
		t.Run(tc.word, func(t *testing.T) {
			got, err := syllabifier.Syllabify(context.Background(), tc.word)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tc.expected) {
				t.Fatalf("Syllabify(%q) = %v; want %v", tc.word, got, tc.expected)
			}
			for i := range got {
				if got[i] != tc.expected[i] {
					t.Errorf("Syllabify(%q)[%d] = %q; want %q", tc.word, i, got[i], tc.expected[i])
				}
			}
		})
	}
}

func TestSyllabify_Triphthongs(t *testing.T) {
	tests := []struct {
		word     string
		expected []string
	}{
		{"Paraguaio", []string{"Pa", "ra", "gua", "io"}},
		{"enxaguei", []string{"en", "xa", "guei"}},
		{"Uruguaio", []string{"U", "ru", "gua", "io"}},
		{"sosseguei", []string{"sos", "se", "guei"}},

		// Genuine triphthongs / edge cases
		{"quais", []string{"quais"}},                 // uai
		{"iguais", []string{"i", "guais"}},           // guais
		{"enxaguou", []string{"en", "xa", "guou"}},   // uou
		{"ameaçou", []string{"a", "me", "a", "çou"}}, // a‑me‑a‑çou
		{"saguão", []string{"sa", "guão"}},           // uão
	}

	for _, tc := range tests {
		t.Run(tc.word, func(t *testing.T) {
			got, err := syllabifier.Syllabify(context.Background(), tc.word)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != len(tc.expected) {
				t.Fatalf("Syllabify(%q) = %v; want %v", tc.word, got, tc.expected)
			}
			for i := range got {
				if got[i] != tc.expected[i] {
					t.Errorf("Syllabify(%q)[%d] = %q; want %q", tc.word, i, got[i], tc.expected[i])
				}
			}
		})
	}
}
