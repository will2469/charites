package token

import (
	"strings"

	"github.com/will2469/charites/internal/token/theme"
)

// ID adalah identifier numerik unik untuk setiap node deklarasi token.
type ID uint32

// TokenID adalah alias untuk ID agar konsisten lintas package.
//
//nolint:revive // alias retained for backward compatibility
type TokenID = ID

// Condition merepresentasikan kueri kondisi pada at-rule (@media, @supports, @container).
type Condition struct {
	Type  string // "media", "supports", "container"
	Query string // "(prefers-color-scheme: dark)", "(min-width: 768px)"
}

// AtRule merepresentasikan konteks at-rule hierarkis tempat token dideklarasikan.
type AtRule struct {
	Name       string      // "@media", "@layer", "@supports", "@container", "@theme"
	Prelude    string      // "theme", "(prefers-color-scheme: dark)"
	Conditions []Condition // Rincian kondisi terurai
}

// Specificity merepresentasikan bobot spesifisitas selektor CSS (A, B, C).
// A: ID selectors (#header)
// B: Class, attribute, & pseudo-classes (.card, [data-theme], :root, :hover)
// C: Elements & pseudo-elements (div, html, ::before)
type Specificity struct {
	A int
	B int
	C int
}

// GreaterThan mengembalikan true jika s memiliki spesifisitas lebih tinggi dari other.
func (s Specificity) GreaterThan(other Specificity) bool {
	if s.A != other.A {
		return s.A > other.A
	}
	if s.B != other.B {
		return s.B > other.B
	}
	return s.C > other.C
}

// Scope merepresentasikan konteks deklarasi token lengkap, termasuk selektor,
// rantai at-rule pembungkus, CSS cascade layer, urutan kemunculan, dan spesifisitas.
type Scope struct {
	Selector    string
	AtRules     []AtRule
	Layers      []string
	SourceOrder int
	Specificity Specificity
}

// IsRoot mengembalikan true jika scope merupakan root level (:root atau html).
func (s Scope) IsRoot() bool {
	sel := strings.TrimSpace(s.Selector)
	return sel == ":root" || sel == "html" || strings.HasPrefix(sel, ":root") || strings.HasPrefix(sel, "html")
}

// IsDark mengembalikan true jika scope berada di dalam konteks dark mode
// (.dark, [data-theme="dark"], atau @media (prefers-color-scheme: dark)).
func (s Scope) IsDark() bool {
	sel := strings.ToLower(s.Selector)
	if strings.Contains(sel, ".dark") || strings.Contains(sel, "dark") {
		return true
	}
	for _, at := range s.AtRules {
		if strings.Contains(strings.ToLower(at.Prelude), "prefers-color-scheme: dark") ||
			strings.Contains(strings.ToLower(at.Prelude), "prefers-color-scheme:dark") {
			return true
		}
	}
	return false
}

// Token merepresentasikan deklarasi token desain tunggal sebagai fakta murni.
// Mempertahankan identitas deklarasi individual meskipun beberapa deklarasi memiliki nama yang sama
// (misal: :root { --brand: red; } vs .card { --brand: blue; }).
type Token struct {
	ID         TokenID
	Name       string           // Nama custom property, e.g. "--banana", "--color-primary"
	RawValue   string           // Nilai verbatim CSS, e.g. "#123456", "var(--banana)"
	Scope      Scope            // Konteks scope deklarasi
	Span       theme.SourceSpan // Posisi baris & kolom pada berkas sumber
	References []string         // Nama-nama token lain yang direferensikan melalui var(--...)
}

// ComputeSpecificity menghitung bobot spesifisitas CSS dari string selektor.
func ComputeSpecificity(selector string) Specificity {
	s := strings.TrimSpace(selector)
	if s == "" {
		return Specificity{}
	}

	var spec Specificity
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == ' ' || r == '>' || r == '+' || r == '~' || r == ','
	})

	for _, part := range parts {
		if part == "" {
			continue
		}
		// Hitung ID (#)
		spec.A += strings.Count(part, "#")

		// Hitung Classes (.)
		spec.B += strings.Count(part, ".")

		// Hitung Attributes ([...])
		spec.B += strings.Count(part, "[")

		// Hitung Pseudo-classes & Pseudo-elements (:)
		colons := strings.Count(part, ":")
		doubleColons := strings.Count(part, "::")
		spec.C += doubleColons                  // Pseudo-elements (::before)
		spec.B += (colons - (doubleColons * 2)) // Pseudo-classes (:hover, :root)

		// Elemen dasar (html, div, body, svg, etc.)
		cleaned := strings.Map(func(r rune) rune {
			if r == '.' || r == '#' || r == ':' || r == '[' || r == ']' || r == '*' {
				return ' '
			}
			return r
		}, part)
		for _, el := range strings.Fields(cleaned) {
			if el != "" && !strings.HasPrefix(el, "&") {
				spec.C++
			}
		}
	}

	return spec
}
