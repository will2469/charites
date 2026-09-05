package theme

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// GradientHardcodeRule mendeteksi penggunaan warna palet primitif, hex arbitrer, atau monokrom
// pada color stop gradien Tailwind (from-*, via-*, to-*).
type GradientHardcodeRule struct{}

// NewGradientHardcodeRule membuat instance baru GradientHardcodeRule.
func NewGradientHardcodeRule() *GradientHardcodeRule {
	return &GradientHardcodeRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *GradientHardcodeRule) ID() string {
	return "theme.gradient-hardcode"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *GradientHardcodeRule) Description() string {
	return "Detects hardcoded primitive, arbitrary hex, or monochrome colors in gradient stops"
}

// Category mengembalikan nama kategori rule.
func (r *GradientHardcodeRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *GradientHardcodeRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki otomatis.
func (r *GradientHardcodeRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Design Tokens Community Group (DTCG)",
			"Tailwind CSS Gradient Token Architecture",
		},
		CoreInvariant: "Gradient color stops must use semantic tokens (from-primary, to-accent), never primitive palette or arbitrary hex stops.",
		Grounding: "Gradients often span large hero sections or callout backgrounds. When color stops use primitive or arbitrary values (e.g. from-[#3b82f6] to-blue-500):\n\n" +
			"1. Inverted Muddy Colors: Light mode gradients rendered in dark mode produce muddy, low-contrast, or unreadable backgrounds behind text.\n" +
			"2. Theme Decoupling: Rebranding or dynamic tenant themes cannot adjust the stops without manually updating every gradient class.\n" +
			"3. Accessibility Violations: Static gradient stops cannot guarantee compliance with WCAG 2.2 text contrast across all screen areas.\n\n" +
			"Charites enforces gradient stops constructed from semantic tokens (from-primary, to-secondary, via-accent, from-transparent).",
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Gradient stops using arbitrary hex and primitive colors",
				Code:     `<div class="bg-gradient-to-r from-[#3b82f6] to-blue-500">Banner</div>`,
			},
			{
				Language: "tsx",
				Comment:  "Gradient stops using monochrome white and primitive red",
				Code: `export function Hero() {
  return <div className="bg-gradient-to-b from-white via-rose-500 to-black">Hero</div>;
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Semantic tokens for gradient stops",
				Code:     `<div class="bg-gradient-to-r from-primary to-accent">Banner</div>`,
			},
			{
				Language: "tsx",
				Comment:  "Semantic tokens adapting cleanly to dark mode",
				Code: `export function Hero() {
  return <div className="bg-gradient-to-b from-card via-primary to-background">Hero</div>;
}`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Dark Mode Breakage",
				Severity: "HIGH",
				Impact:   "Hardcoded gradient stops destroy text legibility and brand alignment in dark themes.",
			},
			{
				Vector:   "Design Token Fragmentation",
				Severity: "MEDIUM",
				Impact:   "Gradients drift out of sync with established design system tokens.",
			},
		},
	}
}

var gradientPrefixes = []string{"from-", "via-", "to-"}

// isHardcodedGradientColor memeriksa apakah segmen warna pada color stop gradien tergolong hardcoded.
func isHardcodedGradientColor(colorPart string) bool {
	if colorPart == "" || strings.HasSuffix(colorPart, "%") {
		return false
	}
	if IsTailwindPrimitiveColor(colorPart) || IsMonochromeColor(colorPart) || IsHexColor(colorPart) || IsColorFunction(colorPart) {
		return true
	}
	if _, ok := ExtractRawColorFromArbitrary(colorPart); ok {
		return true
	}
	return false
}

// Evaluate mengevaluasi node IR dan mendeteksi stop gradien yang di-hardcode.
func (r *GradientHardcodeRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || len(node.Classes) == 0 {
		return nil
	}

	var diags []ir.Diagnostic
	for _, class := range node.Classes {
		_, base := StripVariants(class)
		baseNoAlpha, _, _ := SplitAlphaModifier(base)

		for _, p := range gradientPrefixes {
			if strings.HasPrefix(baseNoAlpha, p) {
				colorPart := baseNoAlpha[len(p):]
				if isHardcodedGradientColor(colorPart) {
					diags = append(diags, ir.Diagnostic{
						Line:     node.Span.Line,
						Column:   node.Span.Column,
						Rule:     r.ID(),
						Severity: r.DefaultSeverity(),
						Message:  "Hardcoded gradient color stop: \"" + class + "\"",
						Hint:     "Use semantic tokens (e.g. \"from-primary\", \"to-accent\") for gradient color stops.",
					})
				}
				break
			}
		}
	}

	return diags
}
