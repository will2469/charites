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
			"Tailwind CSS Transform Scale Architecture",
		},
		CoreInvariant: "Spatial dimensions, spacing intervals, typography sizes, and transform scale modifiers must use standardized modular scale tokens or CSS variables, never arbitrary raw scalar values, non-standard fractional steps, or inline arbitrary scale modifiers.",
		Grounding: "Embedding arbitrary scalar dimensions (e.g. p-[19px], w-[320px], or text-[15px]), non-standard fractional scales (e.g. p-3.25, w-2.75), or arbitrary inline scale modifiers (e.g. active:scale-[0.99], hover:scale-[1.02]) introduces severe UI design regressions:\n\n" +
			"1. Spatial Rhythm Degradation: Arbitrary pixel/rem values and fractional decimals shatter the visual harmony of 4px/8px modular grid systems.\n" +
			"2. Sub-pixel Anti-Aliasing Blur: Fractional step dimensions like p-3.25 (13px) or w-2.75 (11px) fail to align cleanly on various mobile device pixel ratios (DPR), causing fuzzy borders and sub-pixel rendering artifacts.\n" +
			"3. False Conformance: Tailwind IntelliSense suggests shorthand decimals (e.g. 'p-[13px] can be written as p-3.25') which pass Tailwind validation without warning, but violate design system consistency.\n" +
			"4. Maintenance Overhead: Dispersed magic numbers make global layout scaling and responsive adaptation cumbersome.\n" +
			"5. Transform Scale Drift & Subpixel Blurring: Arbitrary inline scale modifiers (e.g. active:scale-[0.99], hover:scale-[1.02], -scale-x-[0.95], [scale:0.98]) fragment tactile physics and micro-interactions across pages, trigger blurry text and 1px borders due to fractional rasterization, and violate centralized token governance.\n\n" +
			"Charites enforces migrating arbitrary sizing, non-standard fractional utilities, and arbitrary inline scales to standard token steps (e.g. p-3, p-3.5, p-4, w-80, text-base, scale-95, scale-105) or token-backed CSS variables (e.g. scale-[var(--scale-press)]).",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Arbitrary padding, width, and non-standard fractional scale in JSX",
				Code:     `<div className="p-[19px] w-[320px] p-3.25 w-2.75 text-[15px]">Hardcoded container</div>`,
			},
			{
				Language: "tsx",
				Comment:  "Arbitrary inline transform scale modifiers in JSX",
				Code:     `<div className="active:scale-[0.99] hover:scale-[1.02] -scale-x-[0.95] [scale:0.98]">Arbitrary scale</div>`,
			},
			{
				Language: "astro",
				Comment:  "Arbitrary spacing and non-standard decimal in Astro component",
				Code:     `<section class="gap-1.25 mt-[27px] [padding:19px]">Arbitrary layout</section>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Standard modular scale tokens (integers and canonical .5 half-steps)",
				Code:     `<div className="p-3.5 p-4 w-80 text-base">Standard modular container</div>`,
			},
			{
				Language: "tsx",
				Comment:  "Standard Tailwind scale steps and token-backed CSS variable",
				Code:     `<div className="active:scale-95 hover:scale-105 active:scale-[var(--scale-press)]">Standard scale</div>`,
			},
			{
				Language: "astro",
				Comment:  "System tokens and CSS variables",
				Code:     `<section class="gap-3 mt-6 p-4">Standard layout</section>`,
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
			{
				Vector:   "Sub-pixel Rendering Artifacts",
				Severity: "MEDIUM",
				Impact:   "Non-standard off-grid decimal scales (e.g. 11px, 13px) produce rounding errors and blurry borders across fractional display scales.",
			},
			{
				Vector:   "Physics & Motion Fragmentation",
				Severity: "MEDIUM",
				Impact:   "Arbitrary inline scale modifiers produce erratic tactile micro-interactions and blurry subpixel text rendering on standard DPI displays.",
			},
		},
	}
}

// Evaluate mengevaluasi sebuah node IR dan mendeteksi kelas spasial dengan skalar arbitrer,
// pecahan desimal non-standar (p-3.25, w-2.75), atau modifier skala arbitrer (scale-[0.99], [scale:0.98]).
// Mematuhi kontrak pure function dan zero-alloc pada node bersih (QUAL-03).
func (r *HardcodeSizeRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || len(node.Classes) == 0 {
		return nil
	}

	//nolint:prealloc // zero-alloc on clean nodes required by QUAL-03
	var diags []ir.Diagnostic
	for _, class := range node.Classes {
		base := StripVariantsOnlyBase(class)
		kind := ClassifySizeUtility(base)
		if kind == SizeClassNone {
			continue
		}

		var msg, hint string
		switch kind {
		case SizeClassArbitraryScale:
			msg = "Arbitrary inline scale modifier: \"" + class + "\""
			hint = "Avoid arbitrary inline scale modifiers. Define a standardized scale token/variable in global.css (e.g. --scale-press) or use standard Tailwind scale steps (scale-95, scale-105)."
		case SizeClassNonStandardFraction:
			msg = "Non-standard fractional scale: \"" + class + "\" (violates 4px/8px modular rhythm)"
			hint = "Use a standard integer step (e.g. p-3, p-4) or official half-step (e.g. p-3.5) instead of an off-grid decimal."
		case SizeClassArbitraryScalar:
			msg = "Hardcoded size/spacing scalar: \"" + class + "\""
			hint = "Use a standard modular spacing step (e.g. p-4, p-5) or a semantic design token."
		}

		diags = append(diags, ir.Diagnostic{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  msg,
			Hint:     hint,
		})
	}

	return diags
}
