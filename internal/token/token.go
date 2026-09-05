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
	Type     string            // "media", "supports", "container"
	Query    string            // raw query string, e.g. "(prefers-color-scheme: dark)", "(min-width: 768px)"
	Features map[string]string // pasangan fitur kunci-nilai terurai, e.g. {"prefers-color-scheme": "dark"}
	IsDark   bool              // true jika kondisi mewajibkan dark mode
	IsLight  bool              // true jika kondisi mewajibkan light mode
}

// ParseCondition mengurai string prelude at-rule (@media, @supports, @container)
// menjadi struktur Condition dengan ekstraksi fitur terstruktur.
func ParseCondition(atName, query string) Condition {
	name := strings.ToLower(atName)
	condType := strings.TrimPrefix(name, "@")
	c := Condition{
		Type:     condType,
		Query:    strings.TrimSpace(query),
		Features: make(map[string]string),
	}

	qLower := strings.ToLower(query)
	if strings.Contains(qLower, "prefers-color-scheme") {
		if strings.Contains(qLower, "dark") {
			c.IsDark = true
			c.Features["prefers-color-scheme"] = "dark"
		} else if strings.Contains(qLower, "light") {
			c.IsLight = true
			c.Features["prefers-color-scheme"] = "light"
		}
	}

	// Ekstrak pasangan fitur (key: value) di dalam kurung
	for _, part := range strings.Split(query, ")") {
		idx := strings.Index(part, "(")
		if idx == -1 {
			continue
		}
		inner := strings.TrimSpace(part[idx+1:])
		colon := strings.Index(inner, ":")
		if colon > 0 {
			k := strings.TrimSpace(inner[:colon])
			v := strings.TrimSpace(inner[colon+1:])
			if k != "" && v != "" {
				c.Features[strings.ToLower(k)] = strings.ToLower(v)
			}
		}
	}

	return c
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

// CascadeRank merepresentasikan tuple pembobotan deklarasi CSS sesuai spesifikasi
// W3C CSS Cascading and Inheritance Level 5 (Cascade Sort).
type CascadeRank struct {
	// ConditionScore:
	// 2 = kondisi cocok persis dengan scope konteks (misal prefers-color-scheme: dark di dark context)
	// 1 = tanpa kondisi (universal fallback e.g. :root biasa)
	// 0 = kondisi berkonflik atau tidak dapat diaplikasikan (inapplicable)
	ConditionScore int

	// LayerRank:
	// Unlayered styles selalu mengalahkan layered styles (LayerRank = 1_000_000).
	// Di antara layered styles, layer yang dideklarasikan belakangan memiliki LayerRank lebih tinggi.
	LayerRank int

	// Specificity: bobot spesifisitas selektor CSS (A, B, C).
	Specificity Specificity

	// SelectorAffinity:
	// 2 = selektor persis sama dengan referer
	// 1 = selektor root (:root atau html)
	// 0 = selektor lain
	SelectorAffinity int

	// SourceOrder: urutan kemunculan deklarasi dalam berkas CSS.
	// Jika seluruh kriteria di atas imbang, deklarasi terakhir (source order lebih besar) yang menang.
	SourceOrder int
}

// GreaterThan mengembalikan true jika r memiliki prioritas cascade lebih tinggi daripada other.
func (r CascadeRank) GreaterThan(other CascadeRank) bool {
	if r.ConditionScore != other.ConditionScore {
		return r.ConditionScore > other.ConditionScore
	}
	if r.LayerRank != other.LayerRank {
		return r.LayerRank > other.LayerRank
	}
	if r.Specificity != other.Specificity {
		if r.Specificity.GreaterThan(other.Specificity) {
			return true
		}
		if other.Specificity.GreaterThan(r.Specificity) {
			return false
		}
	}
	if r.SelectorAffinity != other.SelectorAffinity {
		return r.SelectorAffinity > other.SelectorAffinity
	}
	return r.SourceOrder > other.SourceOrder
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
		for _, cond := range at.Conditions {
			if cond.IsDark {
				return true
			}
		}
		if strings.Contains(strings.ToLower(at.Prelude), "prefers-color-scheme: dark") ||
			strings.Contains(strings.ToLower(at.Prelude), "prefers-color-scheme:dark") {
			return true
		}
	}
	return false
}

// MatchesConditions memeriksa apakah seluruh kondisi pada s kompatibel dengan targetScope.
func (s Scope) MatchesConditions(target Scope) bool {
	if len(s.AtRules) == 0 {
		return true
	}

	targetIsDark := target.IsDark()
	for _, at := range s.AtRules {
		for _, cond := range at.Conditions {
			if cond.IsDark && !targetIsDark {
				return false
			}
			if cond.IsLight && targetIsDark {
				return false
			}
			// Periksa kecocokan fitur spesifik
			for k, v := range cond.Features {
				if k == "prefers-color-scheme" {
					continue
				}
				for _, targetAt := range target.AtRules {
					for _, targetCond := range targetAt.Conditions {
						if targetVal, ok := targetCond.Features[k]; ok && targetVal != v {
							return false
						}
					}
				}
			}
		}
	}

	return true
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
