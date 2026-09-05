package regression_test

import (
	"os"
	"testing"

	"github.com/will2469/charites/internal/parser/astro"
	"github.com/will2469/charites/internal/parser/css"
	"github.com/will2469/charites/internal/parser/tsx"
)

func TestRegression_Corpus(t *testing.T) {
	// 1. Astro single quote class regression
	astroSrc, err := os.ReadFile("single_quote_class.astro")
	if err != nil {
		t.Fatalf("failed to read astro regression fixture: %v", err)
	}
	if _, parseErr := astro.Parse(astroSrc); parseErr != nil {
		t.Fatalf("unexpected error parsing single_quote_class.astro: %v", parseErr)
	}

	// 2. TSX generic arrow regression
	tsxSrc, err := os.ReadFile("generic_arrow.tsx")
	if err != nil {
		t.Fatalf("failed to read tsx regression fixture: %v", err)
	}
	root, err := tsx.Extract(tsxSrc)
	if err != nil {
		t.Fatalf("unexpected error extracting generic_arrow.tsx: %v", err)
	}
	if root.Tag != "div" {
		t.Errorf("expected root tag 'div', got %q", root.Tag)
	}

	// 3. Tailwind unclosed @theme regression
	cssSrc, err := os.ReadFile("unclosed_theme.css")
	if err != nil {
		t.Fatalf("failed to read css regression fixture: %v", err)
	}
	reg, err := css.ParseTheme(cssSrc)
	if err != nil {
		t.Fatalf("unexpected error parsing unclosed_theme.css: %v", err)
	}
	if reg.Variables["--color-unclosed"] != "#2563eb" {
		t.Errorf("expected --color-unclosed to be extracted, got %q", reg.Variables["--color-unclosed"])
	}
}
