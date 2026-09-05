package theme

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// HardcodeBorderColorRule mendeteksi penggunaan warna garis batas (border / divide) yang di-hardcode
// menggunakan palet primitif Tailwind, hex/rgb arbitrer, atau monokrom statis.
type HardcodeBorderColorRule struct{}

// NewHardcodeBorderColorRule membuat instance baru HardcodeBorderColorRule.
func NewHardcodeBorderColorRule() *HardcodeBorderColorRule {
	return &HardcodeBorderColorRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *HardcodeBorderColorRule) ID() string {
	return "theme.hardcode-border-color"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *HardcodeBorderColorRule) Description() string {
	return "Detects hardcoded border and divider colors using primitive palettes, raw hex literals, or static monochrome"
}

// Category mengembalikan nama kategori rule.
func (r *HardcodeBorderColorRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *HardcodeBorderColorRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki otomatis.
func (r *HardcodeBorderColorRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Design Tokens Community Group (DTCG)",
			"Tailwind CSS Border Token Architecture",
		},
		CoreInvariant: "Component borders and dividers must use semantic tokens (border-border, border-input), never primitive palette or arbitrary hex colors.",
		Grounding: "Border lines define container elevation, separation, and affordance. When border colors are hardcoded (e.g. border-gray-200, border-[#e5e5e5]):\n\n" +
			"1. Invisibility in Dark Mode: A light gray border (#e5e5e5) provides zero contrast or turns into an inverted stark line in dark themes.\n" +
			"2. Theme Disconnect: When the primary or brand palette changes, borders remain pinned to legacy gray scales.\n" +
			"3. Inconsistent Boundaries: Disparate components end up using gray-200, slate-200, zinc-300 arbitrarily for identical UI dividers.\n\n" +
			"Charites enforces using centralized border tokens (border-border, border-input, divide-border).",
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Hardcoded border primitive and arbitrary hex",
				Code:     `<div class="border border-gray-200 divide-y divide-[#e5e5e5]">List</div>`,
			},
			{
				Language: "tsx",
				Comment:  "Primitive directional border in JSX",
				Code: `export function Card() {
  return <div className="border-t-slate-300 border-x-[#cccccc]">Content</div>;
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Using semantic border and divider tokens",
				Code:     `<div class="border border-border divide-y divide-border">List</div>`,
			},
			{
				Language: "tsx",
				Comment:  "Semantic border tokens with dark mode adaptability",
				Code: `export function Card() {
  return <div className="border-t border-border">Content</div>;
}`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Dark Mode Invisibility",
				Severity: "HIGH",
				Impact:   "Hardcoded light borders vanish or glow unnaturally on dark theme backgrounds.",
			},
			{
				Vector:   "Visual Fragmentation",
				Severity: "MEDIUM",
				Impact:   "Different shades of gray borders destroy cohesive surface elevation hierarchy.",
			},
		},
	}
}

var borderPrefixes = []string{
	"border-t-", "border-r-", "border-b-", "border-l-", "border-x-", "border-y-",
	"divide-x-", "divide-y-", "divide-",
	"border-",
}

// isHardcodedBorderColor memeriksa apakah bagian warna pada utility border tergolong hardcoded.
func isHardcodedBorderColor(colorPart string) bool {
	if colorPart == "" || IsNonColorBorderKeyword(colorPart) {
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

// Evaluate mengevaluasi node IR dan mendeteksi border / divide color yang di-hardcode.
func (r *HardcodeBorderColorRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || len(node.Classes) == 0 {
		return nil
	}

	var diags []ir.Diagnostic
	for _, class := range node.Classes {
		_, base := StripVariants(class)
		baseNoAlpha, _, _ := SplitAlphaModifier(base)

		for _, p := range borderPrefixes {
			if strings.HasPrefix(baseNoAlpha, p) {
				colorPart := baseNoAlpha[len(p):]
				if isHardcodedBorderColor(colorPart) {
					diags = append(diags, ir.Diagnostic{
						Line:     node.Span.Line,
						Column:   node.Span.Column,
						Rule:     r.ID(),
						Severity: r.DefaultSeverity(),
						Message:  "Hardcoded border color: \"" + class + "\"",
						Hint:     "Use semantic border token (e.g. \"border-border\", \"border-input\") instead of hardcoded border color.",
					})
				}
				break
			}
		}
	}

	return diags
}
