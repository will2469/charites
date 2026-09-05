package theme_test

import (
	"testing"

	"github.com/will2469/charites/internal/token/theme"
)

func BenchmarkVar_ExtractTopLevelVarCalls(b *testing.B) {
	input := "calc(var(--base-size, 16px) * var(--scale, var(--ratio, 1.25)) + var(--foo, rgb(1, 2, 3)))"
	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		calls := theme.ExtractTopLevelVarCalls(input)
		if len(calls) != 3 {
			b.Fatalf("expected 3 calls, got %d", len(calls))
		}
	}
}

func BenchmarkVar_ExtractAllVarNames(b *testing.B) {
	input := "calc(var(--base-size, 16px) * var(--scale, var(--ratio, 1.25)) + var(--foo, rgb(1, 2, 3)))"
	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		names := theme.ExtractAllVarNames(input)
		if len(names) != 4 {
			b.Fatalf("expected 4 names, got %d", len(names))
		}
	}
}

func BenchmarkLexer_IdentifierEscape(b *testing.B) {
	src := []byte(`--tw-bg-opacity\:1: 0.85; --color\31 00: #123456; \:focus { outline: none; }`)
	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		lexer := theme.NewLexer(src)
		count := 0
		for {
			tok := lexer.NextToken()
			if tok.Type == theme.TokenEOF {
				break
			}
			count++
		}
	}
}
