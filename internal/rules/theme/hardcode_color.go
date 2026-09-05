package theme

import (
	"github.com/will2469/charites/internal/ir"
)

// HardcodeColorRule mendeteksi penggunaan literal warna heksadesimal atau fungsi warna arbitrer
// di dalam kelas utilitas Tailwind atau arbitrary property (misal: bg-[#2563eb], [color:#fff]).
type HardcodeColorRule struct{}

// NewHardcodeColorRule membuat instance baru HardcodeColorRule.
func NewHardcodeColorRule() *HardcodeColorRule {
	return &HardcodeColorRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *HardcodeColorRule) ID() string {
	return "theme.hardcode-color"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *HardcodeColorRule) Description() string {
	return "Detects hardcoded arbitrary hex or rgb color literals in Tailwind utility classes and arbitrary properties"
}

// Category mengembalikan nama kategori rule.
func (r *HardcodeColorRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *HardcodeColorRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki otomatis.
func (r *HardcodeColorRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Design Tokens Community Group (DTCG)",
			"Tailwind CSS Design Token Architecture",
			"WCAG 2.2 Contrast Predictability",
		},
		CoreInvariant: "Color declarations in markup must use centralized semantic design tokens or CSS variables, never arbitrary raw hex or color function literals.",
		Grounding: "Directly embedding raw hex or rgb colors (e.g. bg-[#2563eb] or [color:#fff]) inside UI components creates serious maintenance barriers:\n\n" +
			"1. Theme Blindness: Arbitrary color values cannot respond to dark mode, high-contrast, or tenant theme switching.\n" +
			"2. Design Drift: Slight variations in hex codes (e.g. #2563eb vs #2564ea) fracture visual consistency.\n" +
			"3. Inflexible Rebranding: Global style updates require searching and replacing thousands of isolated class strings.\n\n" +
			"Charites enforces migrating arbitrary color literals to semantic tokens defined in global.css (e.g. bg-primary, text-card-foreground).",
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Arbitrary hex color in class attribute",
				Code:     `<div class="bg-[#1e293b] text-[#f8fafc] [color:#fff]">Un-tokenized Card</div>`,
			},
			{
				Language: "tsx",
				Comment:  "Arbitrary rgb and hex literals in JSX",
				Code: `export function Badge() {
  return <span className="hover:bg-[#2563eb] text-[rgb(255,0,0)]">Status</span>;
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Using semantic tokens and CSS variables",
				Code:     `<div class="bg-card text-card-foreground">Tokenized Card</div>`,
			},
			{
				Language: "tsx",
				Comment:  "Semantic token utility with dark mode support",
				Code: `export function Badge() {
  return <span className="hover:bg-primary text-destructive">Status</span>;
}`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Theme Inflexibility",
				Severity: "HIGH",
				Impact:   "Hardcoded hex values remain static during dark mode toggle, causing illegible text and broken contrast.",
			},
			{
				Vector:   "Maintenance Bloat",
				Severity: "MEDIUM",
				Impact:   "Scattered arbitrary colors prevent centralized palette changes and design system updates.",
			},
		},
	}
}

// Evaluate mengevaluasi sebuah node IR dan mendeteksi utility class dengan nilai warna arbitrer.
// Mematuhi kontrak pure function dan zero-alloc pada node bersih (QUAL-03).
func (r *HardcodeColorRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || len(node.Classes) == 0 {
		return nil
	}

	var diags []ir.Diagnostic
	for _, class := range node.Classes {
		_, base := StripVariants(class)
		baseNoAlpha, _, _ := SplitAlphaModifier(base)

		rawColor, isArbitrary := ExtractRawColorFromArbitrary(baseNoAlpha)
		if !isArbitrary {
			prop, val, isProp := ParseArbitraryProperty(baseNoAlpha)
			if isProp && IsColorProperty(prop) && (IsHexColor(val) || IsColorFunction(val)) {
				rawColor = val
				isArbitrary = true
			}
		}

		if isArbitrary {
			diags = append(diags, ir.Diagnostic{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Hardcoded color literal: \"" + class + "\"",
				Hint:     "Use a semantic design token or CSS variable instead of arbitrary color \"" + rawColor + "\".",
			})
		}
	}

	return diags
}
