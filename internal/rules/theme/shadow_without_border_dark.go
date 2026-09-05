package theme

import (
	"github.com/will2469/charites/internal/ir"
)

// ShadowWithoutBorderDarkRule mendeteksi kontainer berelevasi (shadow-md, shadow-lg, shadow-xl, dll.)
// yang tidak memiliki class border atau ring untuk mempertahankan batas visual di dark mode.
type ShadowWithoutBorderDarkRule struct{}

// NewShadowWithoutBorderDarkRule membuat instance baru ShadowWithoutBorderDarkRule.
func NewShadowWithoutBorderDarkRule() *ShadowWithoutBorderDarkRule {
	return &ShadowWithoutBorderDarkRule{}
}

// ID mengembalikan Charites Rule ID kanonikal.
func (r *ShadowWithoutBorderDarkRule) ID() string {
	return "theme.shadow-without-border-dark"
}

// Description mengembalikan penjelasan ringkas rule.
func (r *ShadowWithoutBorderDarkRule) Description() string {
	return "Detects elevated containers with shadow lacking border or ring indicators in dark mode"
}

// Category mengembalikan nama kategori rule.
func (r *ShadowWithoutBorderDarkRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *ShadowWithoutBorderDarkRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ShadowWithoutBorderDarkRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Material Design 3 Elevation Guidelines",
			"W3C DTCG Elevation & Shadow Tokens",
			"Dark Mode Optical Physics & Surface Boundaries",
		},
		CoreInvariant: "Elevated containers in dark mode must include a border or ring to maintain boundary perception against dark canvas backgrounds.",
		Grounding: "In dark mode, standard drop shadows (e.g. shadow-md, shadow-lg, shadow-xl) vanish because black alpha shadows cannot produce luminance contrast against dark or black canvases (optical shadow collapse):\n\n" +
			"1. Boundary Disappearance: High-elevation dialogs, popovers, and cards visually merge into the background canvas.\n" +
			"2. Loss of Spatial Hierarchy: Users lose depth perception and cannot distinguish foreground cards from background sections.\n" +
			"3. Inconsistent Multi-Theme UX: Interfaces that look well-separated in light mode become an unsegmented flat surface in dark mode.\n\n" +
			"Charites enforces pairing elevated shadows with subtle boundary tokens (e.g. border border-border or ring-1 ring-border).",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Elevated card using shadow-xl without border boundary",
				Code:     `<div className="bg-card shadow-xl rounded-xl p-6">Modal Dialog</div>`,
			},
			{
				Language: "astro",
				Comment:  "High-elevation floating panel without border or ring",
				Code:     `<div class="shadow-lg rounded-2xl bg-zinc-900 p-4">Floating Panel</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Elevated card reinforced with border-border boundary",
				Code:     `<div className="bg-card border border-border shadow-xl rounded-xl p-6">Modal Dialog</div>`,
			},
			{
				Language: "astro",
				Comment:  "Elevated panel reinforced with ring token",
				Code:     `<div class="shadow-lg ring-1 ring-border rounded-2xl bg-zinc-900 p-4">Floating Panel</div>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Dark Mode Shadow Collapse",
				Severity: "MEDIUM",
				Impact:   "Elevated elements blend completely into background surfaces in dark themes, eliminating depth cues.",
			},
			{
				Vector:   "Spatial Hierarchy Degradation",
				Severity: "MEDIUM",
				Impact:   "Users experience layout confusion between distinct interactive surfaces and parent containers.",
			},
		},
	}
}

// Evaluate mengevaluasi sebuah node IR untuk mendeteksi kontainer ber-shadow tanpa border/ring.
func (r *ShadowWithoutBorderDarkRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || len(node.Classes) == 0 {
		return nil
	}

	var shadowClass string
	for _, class := range node.Classes {
		base := StripVariantsOnlyBase(class)
		if IsElevatedShadow(base) {
			shadowClass = class
			break
		}
	}

	if shadowClass == "" {
		return nil
	}

	if HasBorderOrRing(node.Classes) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Elevated container has shadow (\"" + shadowClass + "\") without border or ring for dark mode definition",
			Hint:     "Add a border (e.g. border border-border) or ring (ring-1 ring-border) to maintain boundary definition in dark mode.",
		},
	}
}
