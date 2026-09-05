package theme

import (
	"github.com/will2469/charites/internal/ir"
)

// HardcodeZIndexRule mendeteksi penggunaan z-index arbitrer skalar
// di dalam kelas utilitas Tailwind atau arbitrary property (misal: z-[9999], z-[100], [z-index:999]).
type HardcodeZIndexRule struct{}

// NewHardcodeZIndexRule membuat instance baru HardcodeZIndexRule.
func NewHardcodeZIndexRule() *HardcodeZIndexRule {
	return &HardcodeZIndexRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *HardcodeZIndexRule) ID() string {
	return "theme.hardcode-z-index"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *HardcodeZIndexRule) Description() string {
	return "Detects hardcoded arbitrary z-index scalars that trigger stacking context wars"
}

// Category mengembalikan nama kategori rule.
func (r *HardcodeZIndexRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *HardcodeZIndexRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki otomatis.
func (r *HardcodeZIndexRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"CSS Stacking Context Specification",
			"Design System Elevation Hierarchy",
			"Modal & Overlay Governance Standards",
		},
		CoreInvariant: "Element stacking context elevation must be declared using semantic elevation tokens or CSS variables, never arbitrary numerical z-index scalars.",
		Grounding: "Using arbitrary z-index values (e.g. z-[9999] or [z-index:1000]) triggers destructive 'z-index wars':\n\n" +
			"1. Stacking Context Escalation: When engineers pick arbitrary large numbers (999, 9999, 99999) to force elements to the top, other elements inevitably get occluded.\n" +
			"2. Overlay Clashes: Modals, tooltips, dropdown menus, toast notifications, and sticky navigation headers collide unpredictably.\n" +
			"3. Unmaintainable Layering: Without a centralized hierarchy, debugging stacking context bugs requires inspecting the entire DOM tree.\n\n" +
			"Charites enforces utilizing structured elevation tokens (e.g. z-dropdown, z-modal, z-toast) or CSS custom properties (e.g. z-[var(--z-modal)]).",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Arbitrary runaway z-index in fixed modal",
				Code:     `<div className="fixed top-0 z-[9999]">Escalated Modal</div>`,
			},
			{
				Language: "astro",
				Comment:  "Arbitrary property z-index",
				Code:     `<nav class="sticky top-0 [z-index:1000]">Sticky Header</nav>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Semantic elevation token or standard scale",
				Code:     `<div className="fixed top-0 z-50">Controlled Modal</div>`,
			},
			{
				Language: "astro",
				Comment:  "Token variable elevation",
				Code:     `<nav class="sticky top-0 z-[var(--z-sticky)]">Sticky Header</nav>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Z-Index Escalation Wars",
				Severity: "HIGH",
				Impact:   "Engineers continually increase z-index numbers, eventually breaking native select popovers and dialogs.",
			},
			{
				Vector:   "Overlay Occlusion",
				Severity: "HIGH",
				Impact:   "Tooltips and toasts become permanently trapped behind sticky navigations or dropdown menus.",
			},
		},
	}
}

// Evaluate mengevaluasi sebuah node IR dan mendeteksi kelas z-index arbitrer.
// Mematuhi kontrak pure function dan zero-alloc pada node bersih (QUAL-03).
func (r *HardcodeZIndexRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || len(node.Classes) == 0 {
		return nil
	}

	var diags []ir.Diagnostic
	for _, class := range node.Classes {
		base := StripVariantsOnlyBase(class)
		if IsHardcodedZIndex(base) {
			diags = append(diags, ir.Diagnostic{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Hardcoded z-index scalar: \"" + class + "\"",
				Hint:     "Use a semantic elevation token or design token variable (e.g. z-[var(--z-modal)]).",
			})
		}
	}

	return diags
}
