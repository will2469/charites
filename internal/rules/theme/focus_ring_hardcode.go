package theme

import (
	"github.com/will2469/charites/internal/ir"
)

// FocusRingHardcodeRule mendeteksi penggunaan warna mentah atau palet primitif
// pada indikator fokus ring dan outline (misal: focus:ring-[#3b82f6], ring-blue-500).
type FocusRingHardcodeRule struct{}

// NewFocusRingHardcodeRule membuat instance baru FocusRingHardcodeRule.
func NewFocusRingHardcodeRule() *FocusRingHardcodeRule {
	return &FocusRingHardcodeRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *FocusRingHardcodeRule) ID() string {
	return "theme.focus-ring-hardcode"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *FocusRingHardcodeRule) Description() string {
	return "Detects hardcoded primitive palette or arbitrary hex colors on focus rings and outlines"
}

// Category mengembalikan nama kategori rule.
func (r *FocusRingHardcodeRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *FocusRingHardcodeRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki otomatis.
func (r *FocusRingHardcodeRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"WCAG 2.2 Success Criterion 2.4.11 (Focus Not Obscured)",
			"WCAG 2.2 Success Criterion 2.4.13 (Focus Appearance)",
			"W3C DTCG State & Focus Tokens",
		},
		CoreInvariant: "Keyboard focus indicator colors must be driven by semantic ring design tokens (e.g. ring-ring), never primitive palette or hardcoded hex colors.",
		Grounding: "Specifying raw hex literals or primitive colors on focus rings (e.g. focus:ring-[#3b82f6] or ring-blue-500) creates severe accessibility and theme regressions:\n\n" +
			"1. WCAG Contrast Failures: Static blue or hex rings fail the minimum 3:1 contrast ratio against dark or tinted component backgrounds.\n" +
			"2. Theme Blindness: A ring-offset-white class flashes a glaring white halo when tabbed in dark mode.\n" +
			"3. Fragmented Keyboard Affordance: Keyboard navigation users experience jarringly different focus indicators across distinct views.\n\n" +
			"Charites enforces using semantic focus tokens (e.g. focus-visible:ring-ring or ring-ring) and token-driven offsets (ring-offset-background).",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Arbitrary hex and primitive focus ring in JSX",
				Code:     `<button className="focus:ring-[#3b82f6] focus:ring-2">Sign in</button>`,
			},
			{
				Language: "astro",
				Comment:  "Primitive ring and static offset in Astro",
				Code:     `<input class="ring-blue-500 ring-offset-white focus:outline-blue-500" />`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Using semantic focus ring token",
				Code:     `<button className="focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2">Sign in</button>`,
			},
			{
				Language: "astro",
				Comment:  "Semantic ring and background-adaptive offset",
				Code:     `<input class="focus:ring-2 focus:ring-ring ring-offset-background" />`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "WCAG 2.4.13 Non-Compliance",
				Severity: "HIGH",
				Impact:   "Low-vision and keyboard users cannot perceive the active focus indicator due to inadequate contrast ratios.",
			},
			{
				Vector:   "Dark Mode Halo Inversion",
				Severity: "MEDIUM",
				Impact:   "Hardcoded light offsets create blinding light borders around active inputs on dark surfaces.",
			},
		},
	}
}

// Evaluate mengevaluasi sebuah node IR dan mendeteksi kelas ring fokus hardcoded.
// Mematuhi kontrak pure function dan zero-alloc pada node bersih (QUAL-03).
func (r *FocusRingHardcodeRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || len(node.Classes) == 0 {
		return nil
	}

	var diags []ir.Diagnostic
	for _, class := range node.Classes {
		base := StripVariantsOnlyBase(class)
		if IsHardcodedFocusRing(base) {
			diags = append(diags, ir.Diagnostic{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Hardcoded focus ring/outline color: \"" + class + "\"",
				Hint:     "Use a semantic focus token (e.g. ring-ring, ring-offset-background) instead of hardcoded colors.",
			})
		}
	}

	return diags
}
