package theme

import (
	"github.com/will2469/charites/internal/ir"
)

// HardcodeMonochromeRule mendeteksi penggunaan utility monokrom statis (bg-white, bg-black, text-white, text-black,
// border-white, ring-black, termasuk modifier alpha seperti bg-black/50 atau text-white/[0.06])
// yang menyebabkan elemen tidak adaptif terhadap perubahan tema light/dark mode.
type HardcodeMonochromeRule struct{}

// NewHardcodeMonochromeRule membuat instance baru HardcodeMonochromeRule.
func NewHardcodeMonochromeRule() *HardcodeMonochromeRule {
	return &HardcodeMonochromeRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *HardcodeMonochromeRule) ID() string {
	return "theme.hardcode-monochrome"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *HardcodeMonochromeRule) Description() string {
	return "Detects hardcoded monochrome utilities (white/black) that fail to adapt across light and dark themes"
}

// Category mengembalikan nama kategori rule.
func (r *HardcodeMonochromeRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *HardcodeMonochromeRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki otomatis.
func (r *HardcodeMonochromeRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Design Tokens Community Group (DTCG)",
			"WCAG 2.2 Relative Contrast (SC 1.4.3)",
			"Tailwind CSS Dark Mode Architecture",
		},
		CoreInvariant: "Surfaces and text must use adaptive semantic tokens (background, foreground, card, popover) rather than hardcoded static white or black.",
		Grounding: "Hardcoding white or black (e.g. bg-white, text-black, bg-black/50) creates glaring dark mode regressions:\n\n" +
			"1. Inverted Blindness: A container styled with bg-white turns into a blinding light box inside dark mode.\n" +
			"2. Invisible Text: Pairing bg-background with text-black causes black-on-black illegible text when the theme switches to dark.\n" +
			"3. Alpha Washout: Static text-white/[0.06] loses contrast completely on lighter surfaces.\n\n" +
			"Charites enforces replacing static monochrome utilities with semantic surface and typography tokens (bg-background, text-foreground, bg-card, text-muted-foreground).",
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Hardcoded static white background and black text",
				Code:     `<div class="bg-white text-black p-6 shadow-md">Un-themed Box</div>`,
			},
			{
				Language: "tsx",
				Comment:  "Static monochrome utilities with alpha modifiers",
				Code: `export function Overlay() {
  return <div className="bg-black/50 text-white/[0.06] border-white">Backdrop</div>;
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Adaptive semantic tokens for cards and text",
				Code:     `<div class="bg-card text-card-foreground p-6 shadow-md border border-border">Themed Box</div>`,
			},
			{
				Language: "tsx",
				Comment:  "Semantic tokens adapting automatically to theme state",
				Code: `export function Overlay() {
  return <div className="bg-background/80 text-muted-foreground border-border">Backdrop</div>;
}`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Contrast Failure",
				Severity: "HIGH",
				Impact:   "Black text on dark background drops contrast ratio to 1:1, completely hiding content.",
			},
			{
				Vector:   "Visual Jarring",
				Severity: "MEDIUM",
				Impact:   "Pure white cards jarringly clash against dark mode UI aesthetics.",
			},
		},
	}
}

// Evaluate mengevaluasi node IR dan mendeteksi kelas utilitas monokrom statis.
func (r *HardcodeMonochromeRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || len(node.Classes) == 0 {
		return nil
	}

	var diags []ir.Diagnostic
	for _, class := range node.Classes {
		_, base := StripVariants(class)
		baseNoAlpha, _, _ := SplitAlphaModifier(base)

		_, remainder, ok := SplitColorPrefix(baseNoAlpha)
		if ok && IsMonochromeColor(remainder) {
			diags = append(diags, ir.Diagnostic{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Hardcoded monochrome utility: \"" + class + "\"",
				Hint:     "Replace static white/black with semantic background/foreground tokens (e.g. bg-background, text-foreground, bg-card, text-card-foreground) to support dark mode.",
			})
		}
	}

	return diags
}
