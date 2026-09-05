package theme

import (
	"github.com/will2469/charites/internal/ir"
)

// PseudoHardcodeColorRule mendeteksi penggunaan warna primitif, hex arbitrer, atau monokrom
// di dalam varian pseudo-element atau pseudo-class (misal: placeholder:text-gray-400, selection:bg-blue-200).
type PseudoHardcodeColorRule struct{}

// NewPseudoHardcodeColorRule membuat instance baru PseudoHardcodeColorRule.
func NewPseudoHardcodeColorRule() *PseudoHardcodeColorRule {
	return &PseudoHardcodeColorRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *PseudoHardcodeColorRule) ID() string {
	return "theme.pseudo-hardcode-color"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *PseudoHardcodeColorRule) Description() string {
	return "Detects hardcoded primitive, arbitrary hex, or monochrome colors inside pseudo-element and pseudo-class variants"
}

// Category mengembalikan nama kategori rule.
func (r *PseudoHardcodeColorRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *PseudoHardcodeColorRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki otomatis.
func (r *PseudoHardcodeColorRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Design Tokens Community Group (DTCG)",
			"Tailwind CSS Pseudo-Element Architecture",
		},
		CoreInvariant: "Pseudo-elements (placeholder, selection, file, marker) must consume semantic tokens, never raw primitive or arbitrary colors.",
		Grounding: "Pseudo-element styling often slips past generic linters that only inspect top-level classes.\n\n" +
			"When developers specify placeholder:text-gray-400 or selection:bg-blue-200:\n" +
			"1. Input Readability Degradation: A placeholder styled with light gray-400 becomes completely invisible on light input surfaces or garish on dark inputs.\n" +
			"2. Selection Contrast Clashes: Static blue-200 selection background can fail WCAG contrast against the text color in dark mode.\n" +
			"3. Inconsistent State Branding: File inputs and list markers fail to reflect global theme tokens.\n\n" +
			"Charites enforces using semantic tokens (placeholder:text-muted-foreground, selection:bg-primary-light).",
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Hardcoded primitive colors in placeholder and selection",
				Code:     `<input class="placeholder:text-gray-400 selection:bg-blue-200" />`,
			},
			{
				Language: "tsx",
				Comment:  "Arbitrary hex in pseudo variants",
				Code: `export function Input() {
  return <input className="placeholder:text-[#94a3b8] file:bg-slate-100" />;
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Semantic tokens for pseudo styling",
				Code:     `<input class="placeholder:text-muted-foreground selection:bg-primary-light" />`,
			},
			{
				Language: "tsx",
				Comment:  "Semantic tokens adapting to active theme",
				Code: `export function Input() {
  return <input className="placeholder:text-muted-foreground file:bg-secondary" />;
}`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Form Accessibility Failure",
				Severity: "HIGH",
				Impact:   "Low-contrast placeholder text fails WCAG minimum ratio, making form inputs confusing for users.",
			},
			{
				Vector:   "Selection Highlight Glitch",
				Severity: "MEDIUM",
				Impact:   "Hardcoded selection backgrounds obliterate text visibility under dark themes.",
			},
		},
	}
}

// isHardcodedColorPart memeriksa apakah bagian pewarnaan utilitas adalah hardcoded.
func isHardcodedColorPart(colorPart string) bool {
	if colorPart == "" {
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

// Evaluate mengevaluasi node IR dan mendeteksi warna hardcoded di dalam varian pseudo.
func (r *PseudoHardcodeColorRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || len(node.Classes) == 0 {
		return nil
	}

	var diags []ir.Diagnostic
	for _, class := range node.Classes {
		variants, base := StripVariants(class)
		if !HasPseudoVariant(variants) {
			continue
		}

		baseNoAlpha, _, _ := SplitAlphaModifier(base)
		_, remainder, ok := SplitColorPrefix(baseNoAlpha)
		if !ok {
			remainder = baseNoAlpha
		}

		if isHardcodedColorPart(remainder) {
			diags = append(diags, ir.Diagnostic{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Hardcoded color in pseudo variant: \"" + class + "\"",
				Hint:     "Use semantic tokens for pseudo-element/class styling (e.g. \"placeholder:text-muted-foreground\", \"selection:bg-primary-light\").",
			})
		}
	}

	return diags
}
