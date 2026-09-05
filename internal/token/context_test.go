package token_test

import (
	"testing"

	"github.com/will2469/charites/internal/token"
)

func TestContext_Immutability(t *testing.T) {
	input := []byte(`:root {
  --color-primary: #123456;
  --color-secondary: #654321;
}`)

	ctx, err := token.ParseCSS(input)
	if err != nil {
		t.Fatalf("ParseCSS failed: %v", err)
	}

	// 1. Uji defensive copy Tokens()
	tokens1 := ctx.Tokens()
	if len(tokens1) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(tokens1))
	}

	// Mutasi slice hasil copy
	tokens1[0].Name = "--mutated-token"
	tokens1[0].RawValue = "#000000"

	// Ambil snapshot baru dan pastikan node asli di graph tidak berubah
	tokens2 := ctx.Tokens()
	if tokens2[0].Name == "--mutated-token" || tokens2[0].RawValue == "#000000" {
		t.Fatalf("mutation leak detected! ctx.Tokens() must return a defensive snapshot")
	}

	tok, ok := ctx.LookupToken("--color-primary")
	if !ok || tok.Name != "--color-primary" || tok.RawValue != "#123456" {
		t.Fatalf("internal graph was corrupted by external slice mutation!")
	}
}

func TestContext_Iterators(t *testing.T) {
	input := []byte(`:root {
  --color-primary: #123456;
  --color-secondary: #654321;
  --space-sm: 4px;
  --space-md: 8px;
}
@media (prefers-color-scheme: dark) {
  :root {
    --color-primary: #ffffff;
  }
}`)

	ctx, err := token.ParseCSS(input)
	if err != nil {
		t.Fatalf("ParseCSS failed: %v", err)
	}

	// 1. AllTokens iterator
	allCount := 0
	for tok := range ctx.AllTokens() {
		if tok.Name == "" {
			t.Errorf("expected non-empty token name")
		}
		allCount++
	}
	if allCount != ctx.TokenCount() {
		t.Errorf("expected %d tokens from AllTokens(), got %d", ctx.TokenCount(), allCount)
	}

	// 2. TokensByName iterator
	primaryCount := 0
	for tok := range ctx.TokensByName("--color-primary") {
		if tok.Name != "--color-primary" {
			t.Errorf("expected --color-primary, got %s", tok.Name)
		}
		primaryCount++
	}
	if primaryCount != 2 {
		t.Errorf("expected 2 declarations for --color-primary, got %d", primaryCount)
	}

	// 3. TokensByPrefix iterator
	colorCount := 0
	for tok := range ctx.TokensByPrefix("--color-") {
		colorCount++
		if tok.RawValue == "" {
			t.Errorf("empty raw value in prefix iterator")
		}
	}
	if colorCount != 3 {
		t.Errorf("expected 3 --color- tokens, got %d", colorCount)
	}
}

func BenchmarkContext_LookupToken(b *testing.B) {
	input := []byte(`:root {
  --color-primary: #123456;
  --color-secondary: #654321;
  --color-accent: #abcdef;
}`)
	ctx, err := token.ParseCSS(input)
	if err != nil {
		b.Fatalf("failed to parse: %v", err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		tok, ok := ctx.LookupToken("--color-accent")
		if !ok || tok.Name == "" {
			b.Fatalf("lookup failed")
		}
	}
}

func BenchmarkContext_HasToken(b *testing.B) {
	input := []byte(`:root {
  --color-primary: #123456;
}`)
	ctx, err := token.ParseCSS(input)
	if err != nil {
		b.Fatalf("failed to parse: %v", err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		if !ctx.HasToken("--color-primary") {
			b.Fatalf("hasToken failed")
		}
	}
}

func BenchmarkContext_AllTokens_Iter(b *testing.B) {
	input := []byte(`:root {
  --color-primary: #123456;
  --color-secondary: #654321;
  --color-accent: #abcdef;
}`)
	ctx, err := token.ParseCSS(input)
	if err != nil {
		b.Fatalf("failed to parse: %v", err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		count := 0
		for tok := range ctx.AllTokens() {
			if tok.Name != "" {
				count++
			}
		}
		if count != 3 {
			b.Fatalf("unexpected count %d", count)
		}
	}
}
