package syllabifier_test

import (
	"context"
	"testing"

	"github.com/ramon-reichert/GABCgen/internal/syllabifier"
)

func TestSyllabify_BasicWords(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"alegria", []string{"a", "le", "gri", "a"}},
		{"oração", []string{"o", "ra", "ção"}},
		{"sagrado", []string{"sa", "gra", "do"}},
		{"misericórdia", []string{"mi", "se", "ri", "cór", "di", "a"}},
		{"glória", []string{"gló", "ri", "a"}},
		{"português", []string{"por", "tu", "guês"}},
		{"Alegria", []string{"A", "le", "gri", "a"}},
		//	{"Espírito-Santo", []string{"Es", "pí", "ri", "to", "-", "San", "to"}},
		//	{"mãe-do-salvador", []string{"mãe", "-", "do", "-", "sal", "va", "dor"}},
	}

	ctx := context.Background()
	for _, tt := range tests {
		got, err := syllabifier.Syllabify(ctx, tt.input)
		if err != nil {
			t.Errorf("unexpected error for %q: %v", tt.input, err)
		}
		if len(got) != len(tt.want) {
			t.Errorf("Syllabify(%q) = %v; want %v", tt.input, got, tt.want)
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("Syllabify(%q)[%d] = %q; want %q", tt.input, i, got[i], tt.want[i])
			}
		}
	}
}
