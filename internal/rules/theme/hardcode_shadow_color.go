package theme

import (
	"github.com/will2469/charites/internal/ir"
)

// HardcodeShadowColorRule mendeteksi penggunaan warna hardcoded di dalam deklarasi box-shadow
// pada kelas utilitas Tailwind atau arbitrary property (misal: shadow-[0_4px_10px_#00000040]).
type HardcodeShadowColorRule struct{}

// NewHardcodeShadowColorRule membuat instance baru HardcodeShadowColorRule.
func NewHardcodeShadowColorRule() *HardcodeShadowColorRule {
	return &HardcodeShadowColorRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *HardcodeShadowColorRule) ID() string {
	return "theme.hardcode-shadow-color"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *HardcodeShadowColorRule) Description() string {
	return "Detects hardcoded color literals embedded in box-shadow declarations"
}

// Category mengembalikan nama kategori rule.
func (r *HardcodeShadowColorRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *HardcodeShadowColorRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki otomatis.
func (r *HardcodeShadowColorRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C DTCG Elevation Tokens",
			"Tailwind CSS Box Shadow Specification",
			"Dark Mode Optical Physics & Contrast",
		},
		CoreInvariant: "Elevation shadows must not embed raw hex or arbitrary color literals; shadow tints must adapt dynamically across light and dark modes via semantic tokens.",
		Grounding: "Embedding raw color literals inside arbitrary shadow brackets (e.g. shadow-[0_4px_10px_#00000040]) introduces major theme defects:\n\n" +
			"1. Dark Mode Disappearance: Dark shadows (black/gray with alpha) disappear completely when rendered over dark backgrounds (e.g. #09090b), leaving elevated cards looking flat.\n" +
			"2. Unadaptive Tints: Brand theme colors cannot tint shadows realistically when hardcoded hex codes are baked into individual classes.\n" +
			"3. Specificity Collisions: Overriding arbitrary shadow strings requires higher specificity or duplicate classes.\n\n" +
			"Charites enforces using standard shadow scale tokens (e.g. shadow-sm, shadow-md, shadow-lg) or semantic elevation tokens defined in global.css.",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Arbitrary shadow with embedded hex color",
				Code:     `<div className="shadow-[0_4px_10px_#00000040] p-6">Floating Card</div>`,
			},
			{
				Language: "astro",
				Comment:  "Arbitrary property box-shadow with rgb",
				Code:     `<section class="[box-shadow:0_10px_15px_rgba(0,0,0,0.1)]">Elevated Panel</section>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Using standard elevation shadow tokens",
				Code:     `<div className="shadow-md p-6">Adaptive Floating Card</div>`,
			},
			{
				Language: "astro",
				Comment:  "CSS variable shadow color",
				Code:     `<section class="shadow-[0_4px_6px_var(--shadow-color)]">Elevated Panel</section>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Dark Mode Elevation Invisibility",
				Severity: "HIGH",
				Impact:   "Hardcoded dark shadows become completely invisible against dark canvases, collapsing visual depth.",
			},
			{
				Vector:   "Inconsistent Ambient Occlusion",
				Severity: "MEDIUM",
				Impact:   "Disparate shadow colors across components destroy uniform light-source perception in the design system.",
			},
		},
	}
}

// Evaluate mengevaluasi sebuah node IR dan mendeteksi kelas bayangan dengan warna hardcoded.
// Mematuhi kontrak pure function dan zero-alloc pada node bersih (QUAL-03).
func (r *HardcodeShadowColorRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || len(node.Classes) == 0 {
		return nil
	}

	var diags []ir.Diagnostic
	for _, class := range node.Classes {
		base := StripVariantsOnlyBase(class)
		if IsHardcodedShadowColor(base) {
			diags = append(diags, ir.Diagnostic{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Hardcoded shadow color: \"" + class + "\"",
				Hint:     "Use a standard elevation shadow token (e.g. shadow-md) or a semantic CSS variable.",
			})
		}
	}

	return diags
}
