package theme

import (
	"github.com/will2469/charites/internal/ir"
)

// HardcodeBorderRadiusRule mendeteksi penggunaan border-radius arbitrer
// di dalam kelas utilitas Tailwind atau arbitrary property (misal: rounded-[7px], rounded-t-[14px]).
type HardcodeBorderRadiusRule struct{}

// NewHardcodeBorderRadiusRule membuat instance baru HardcodeBorderRadiusRule.
func NewHardcodeBorderRadiusRule() *HardcodeBorderRadiusRule {
	return &HardcodeBorderRadiusRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *HardcodeBorderRadiusRule) ID() string {
	return "theme.hardcode-border-radius"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *HardcodeBorderRadiusRule) Description() string {
	return "Detects hardcoded arbitrary border-radius scalars in Tailwind utility classes"
}

// Category mengembalikan nama kategori rule.
func (r *HardcodeBorderRadiusRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *HardcodeBorderRadiusRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki otomatis.
func (r *HardcodeBorderRadiusRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C DTCG Shape & Radius Tokens",
			"Design System Shape Hierarchy",
			"Nested Curvature Optics Standard",
		},
		CoreInvariant: "Corner rounding and shape curvature must use standardized shape tokens or CSS variables, never arbitrary raw radius scalars.",
		Grounding: "Specifying arbitrary border-radius values (e.g. rounded-[7px] or rounded-t-[14px]) harms UI consistency:\n\n" +
			"1. Geometric Discordance: Components with off-scale radii look disjointed when nested or placed side-by-side.\n" +
			"2. Outer/Inner Radius Mismatch: Nested cards require deliberate radius proportion calculations (R_inner = R_outer - padding) defined by the shape system.\n" +
			"3. Rebranding Vulnerability: Global shape system updates (e.g. switching from square to rounded theme) cannot adapt arbitrary bracket classes.\n\n" +
			"Charites enforces using standard shape tokens (e.g. rounded-sm, rounded-md, rounded-xl) or token variables (e.g. rounded-[var(--radius)]).",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Arbitrary border radius on button and card",
				Code:     `<button className="rounded-[7px] p-3">Submit</button>`,
			},
			{
				Language: "astro",
				Comment:  "Directional arbitrary radius in Astro component",
				Code:     `<div class="rounded-t-[14px] [border-radius:9px]">Modal Header</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Using design system shape tokens",
				Code:     `<button className="rounded-md p-3">Submit</button>`,
			},
			{
				Language: "astro",
				Comment:  "Standard directional rounded tokens",
				Code:     `<div class="rounded-t-xl rounded-b-none">Modal Header</div>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Geometric Incoherence",
				Severity: "MEDIUM",
				Impact:   "Arbitrary corner radii make cards, buttons, and inputs clash visually across user interfaces.",
			},
			{
				Vector:   "Theme Rigidity",
				Severity: "HIGH",
				Impact:   "Hardcoded radius prevents sweeping design system theme modernizations or brand shape adjustments.",
			},
		},
	}
}

// Evaluate mengevaluasi sebuah node IR dan mendeteksi kelas border-radius arbitrer.
// Mematuhi kontrak pure function dan zero-alloc pada node bersih (QUAL-03).
func (r *HardcodeBorderRadiusRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || len(node.Classes) == 0 {
		return nil
	}

	var diags []ir.Diagnostic
	for _, class := range node.Classes {
		base := StripVariantsOnlyBase(class)
		if IsHardcodedBorderRadius(base) {
			diags = append(diags, ir.Diagnostic{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Hardcoded border-radius scalar: \"" + class + "\"",
				Hint:     "Use a standard shape token (e.g. rounded-md, rounded-lg) or CSS variable.",
			})
		}
	}

	return diags
}
