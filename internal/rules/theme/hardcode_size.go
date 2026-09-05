package theme

import (
	"github.com/will2469/charites/internal/ir"
)

// HardcodeSizeRule mendeteksi penggunaan ukuran spasial, spacing, posisi, atau tipografi arbitrer
// di dalam kelas utilitas Tailwind atau arbitrary property (misal: p-[19px], w-[320px], text-[15px]).
type HardcodeSizeRule struct{}

// NewHardcodeSizeRule membuat instance baru HardcodeSizeRule.
func NewHardcodeSizeRule() *HardcodeSizeRule {
	return &HardcodeSizeRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *HardcodeSizeRule) ID() string {
	return "theme.hardcode-size"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *HardcodeSizeRule) Description() string {
	return "Detects hardcoded arbitrary size, spacing, or typography scalars in Tailwind utility classes"
}

// Category mengembalikan nama kategori rule.
func (r *HardcodeSizeRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *HardcodeSizeRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki otomatis.
func (r *HardcodeSizeRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C DTCG Spatial Scale Standard",
			"8pt Modular Grid Rhythm",
			"Tailwind CSS Spacing Architecture",
		},
		CoreInvariant: "Spatial dimensions, spacing intervals, and typography sizes must use standardized modular scale tokens or CSS variables, never arbitrary raw scalar values.",
		Grounding: "Embedding arbitrary scalar dimensions (e.g. p-[19px], w-[320px], or text-[15px]) introduces severe UI design regressions:\n\n" +
			"1. Spatial Rhythm Degradation: Arbitrary pixel/rem values shatter the visual harmony of 4px/8px modular grid systems.\n" +
			"2. Typography Drift: Off-scale text sizes break proportional line-height and vertical rhythm standards across viewports.\n" +
			"3. Maintenance Overhead: Dispersed magic numbers make global layout scaling and device adaptations cumbersome.\n\n" +
			"Charites enforces migrating arbitrary sizing and spacing utilities to standard token steps (e.g. p-5, w-80, text-base) or token variables.",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Arbitrary padding, width, and text size in JSX",
				Code:     `<div className="p-[19px] w-[320px] text-[15px]">Hardcoded container</div>`,
			},
			{
				Language: "astro",
				Comment:  "Arbitrary spacing and margin in Astro component",
				Code:     `<section class="gap-[13px] mt-[27px] [padding:19px]">Arbitrary layout</section>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Standard modular scale tokens",
				Code:     `<div className="p-5 w-80 text-base">Standard modular container</div>`,
			},
			{
				Language: "astro",
				Comment:  "System tokens and CSS variables",
				Code:     `<section class="gap-3 mt-6 p-5">Standard layout</section>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Visual Rhythm Breakdown",
				Severity: "MEDIUM",
				Impact:   "Inconsistent micro-spacing across components causes fragmented alignment and sloppy UI rendering.",
			},
			{
				Vector:   "Typography Scale Drift",
				Severity: "HIGH",
				Impact:   "Unchecked font sizes degrade readability, leading calculation, and accessibility scaling.",
			},
		},
	}
}

// Evaluate mengevaluasi sebuah node IR dan mendeteksi kelas spasial dengan skalar arbitrer.
// Mematuhi kontrak pure function dan zero-alloc pada node bersih (QUAL-03).
func (r *HardcodeSizeRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || len(node.Classes) == 0 {
		return nil
	}

	var diags []ir.Diagnostic
	for _, class := range node.Classes {
		base := StripVariantsOnlyBase(class)
		if IsHardcodedSizeUtility(base) {
			diags = append(diags, ir.Diagnostic{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Hardcoded size/spacing scalar: \"" + class + "\"",
				Hint:     "Use a standard modular spacing step (e.g. p-4, p-5) or a semantic design token.",
			})
		}
	}

	return diags
}
