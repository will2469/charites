package fuzz

import (
	"os"
	"testing"

	"github.com/will2469/charites/internal/parser/css"
	"github.com/will2469/charites/internal/token"
)

// FuzzCSSParser menguji kekebalan parser leksikal CSS terhadap input malformed, byte acak,
// dan unclosed token untuk menjamin Zero-Panic Invariant.
func FuzzCSSParser(f *testing.F) {
	// 1. Seed dari fixture nyata
	if sample, err := os.ReadFile("../fixtures/global.css"); err == nil {
		f.Add(sample)
	}

	// 2. Seeds skenario ekstrem dan malformed
	seeds := [][]byte{
		[]byte(""),
		[]byte(":root { --banana: #123456; --super: var(--banana); }"),
		[]byte(".card { --brand: red; &:hover { --brand: blue; } }"),
		[]byte("@layer theme { @media (prefers-color-scheme: dark) { :root { --brand: oklch(0.6 0.2 250); } } }"),
		[]byte("/* unclosed comment"),
		[]byte(`--unclosed-quote: "hello world;`),
		[]byte(`--unclosed-single: 'single quote;`),
		[]byte(`--data-uri: url("data:image/svg+xml;utf8,<svg><path d='M0 0h10v10H0z'/></svg>");`),
		[]byte(`--broken-url: url("unclosed url`),
		[]byte(`--escaped: \:hover { --color\:blue: #00f; }`),
		[]byte(`{{{{}}}};;;::::(((())))[[[[[]]]]]`),
		[]byte(`:root { --a: var(--b); --b: var(--a); }`),
		[]byte(`calc(var(--a) + var(--b, 10px) * 2)`),
		[]byte("@theme { --font-*: initial; }"),
		[]byte("@container (min-width: 400px) { :root { --layout: flex; } }"),
		[]byte("@supports (backdrop-filter: blur(10px)) { .glass { --blur: 10px; } }"),
		[]byte("\x00\x01\x02\xff\xfe\xfd--null-byte: 123;"),
		[]byte(`,,,,,;;;;:::::......`),
		[]byte(`@import "tailwindcss";`),
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		// Zero-panic invariant: Parser dan Extractor dilarang memicu runtime panic
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("panic triggered in css parser/token extractor with input %q: %v", string(data), r)
			}
		}()

		// 1. Uji Generic CSS Parser (Layer 1)
		sheet, _ := css.Parse(data)
		if sheet != nil {
			// Pastikan aturan dapat diiterasi tanpa panic
			for _, r := range sheet.Rules {
				_ = r.GetSpan()
			}
		}

		// 2. Uji Token Extractor & Dependency Graph (Layer 2 & 3)
		ctx, err := token.ParseCSS(data)
		if err == nil && ctx != nil {
			tokens := ctx.Tokens()
			_ = ctx.Scopes()

			// Deteksi siklus harus aman terhadap input acak
			_ = ctx.FindCycles()

			// Resolusi token harus aman terhadap kedalaman apa pun
			maxResolve := 10
			if len(tokens) < maxResolve {
				maxResolve = len(tokens)
			}
			for i := 0; i < maxResolve; i++ {
				_, _, _ = ctx.Resolve(tokens[i].ID, token.ResolveOptions{MaxNodes: 50})
			}
		}
	})
}
