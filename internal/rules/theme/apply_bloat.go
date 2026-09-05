package theme

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// ApplyBloatRule mendeteksi penggunaan berlebih direktif @apply dengan lebih dari 8 utility classes
// di dalam blok style, yang menyebabkan pembengkakan bundle CSS dan merusak keterbacaan kode.
type ApplyBloatRule struct{}

// NewApplyBloatRule membuat instance baru ApplyBloatRule.
func NewApplyBloatRule() *ApplyBloatRule {
	return &ApplyBloatRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *ApplyBloatRule) ID() string {
	return "theme.apply-bloat"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *ApplyBloatRule) Description() string {
	return "Detects excessive use of @apply with more than 8 utility classes in CSS or style blocks"
}

// Category mengembalikan nama kategori rule.
func (r *ApplyBloatRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *ApplyBloatRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki otomatis.
func (r *ApplyBloatRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Tailwind CSS v3/v4 Architectural Best Practices",
			"W3C Web Performance & CSS Bundle Size Guidelines",
		},
		CoreInvariant: "The @apply directive must not aggregate more than 8 utility classes in a single declaration to prevent CSS bloat and abstraction decay.",
		Grounding: "The @apply directive in Tailwind CSS was designed for small semantic abstractions (such as buttons or form inputs). " +
			"Overusing @apply by stacking dozens of utility classes recreates the worst aspects of monolithic CSS.\n\n" +
			"When developers write @apply flex items-center justify-between p-4 bg-white rounded-lg shadow-md border border-gray-200 text-sm font-medium:\n" +
			"1. Bundle Size Inflation: Utility classes are duplicated into individual CSS selectors, negating Tailwind's atomic deduplication benefits.\n" +
			"2. Loss of Utility Ergonomics: Developers lose the ability to override individual styles via props or conditional classes.\n" +
			"3. Maintenance Decay: Giant @apply strings become unreadable 'css-in-css' dumping grounds.\n\n" +
			"Charites enforces a maximum threshold of 8 utility classes per @apply directive.",
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Overloaded @apply declaration with 11 utility classes",
				Code: `<style>
  .card {
    @apply flex items-center justify-between p-4 bg-white rounded-lg shadow-md border border-gray-200 text-sm font-medium;
  }
</style>`,
			},
			{
				Language: "tsx",
				Comment:  "Bloated @apply inside TSX style tag",
				Code: `export function Widget() {
  return (
    <style>{` + "`" + `
      .btn-primary {
        @apply inline-flex items-center justify-center px-4 py-2 text-sm font-semibold rounded-md shadow-sm text-white bg-primary hover:bg-primary/90;
      }
    ` + "`" + `}</style>
  );
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Concise @apply declaration with 4 utility classes",
				Code: `<style>
  .card {
    @apply flex items-center justify-between p-4;
  }
</style>`,
			},
			{
				Language: "tsx",
				Comment:  "Utilities applied directly to JSX markup",
				Code: `export function Card() {
  return <div className="flex items-center justify-between p-4 bg-white rounded-lg shadow-md">Markup Utilities</div>;
}`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "CSS Bundle Bloat",
				Severity: "MEDIUM",
				Impact:   "Overloaded @apply directives balloon production stylesheet size and defeat atomic CSS compression.",
			},
			{
				Vector:   "Component Maintainability Decay",
				Severity: "LOW",
				Impact:   "Massive CSS helper blocks reduce readability and make conditional variant overrides difficult.",
			},
		},
	}
}

// Evaluate mengevaluasi node IR untuk mencari kemunculan @apply dengan lebih dari 8 classes.
func (r *ApplyBloatRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Tag != "style" {
		return nil
	}

	cssText := getStyleNodeText(node)
	if !strings.Contains(cssText, "@apply") {
		return nil
	}

	cleaned := stripCSSCommentsString(cssText)
	count, found := findExcessiveApply(cleaned, 8)
	if !found {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Excessive @apply bloat: %d utility classes detected (threshold is 8)", count),
			Hint:     "Refactor overloaded @apply directives into standard component JSX markup or extract semantic CSS rules.",
		},
	}
}

// findExcessiveApply memeriksa apakah ada deklarasi @apply dengan jumlah token melebihi threshold.
func findExcessiveApply(s string, threshold int) (int, bool) {
	idx := 0
	n := len(s)

	for idx < n {
		pos := strings.Index(s[idx:], "@apply")
		if pos == -1 {
			break
		}

		start := idx + pos + len("@apply")
		semi := strings.IndexByte(s[start:], ';')
		if semi == -1 {
			// Jika tidak ada titik koma, periksa hingga kurung kurawal penutup '}'
			semi = strings.IndexByte(s[start:], '}')
			if semi == -1 {
				semi = n - start
			}
		}

		classBlock := s[start : start+semi]
		count := countWhitespaceTokens(classBlock)
		if count > threshold {
			return count, true
		}

		idx = start + semi + 1
	}

	return 0, false
}

// countWhitespaceTokens menghitung jumlah token string yang dipisahkan oleh whitespace tanpa alokasi slice.
func countWhitespaceTokens(s string) int {
	count := 0
	inToken := false

	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			if inToken {
				count++
				inToken = false
			}
		} else {
			inToken = true
		}
	}
	if inToken {
		count++
	}

	return count
}
