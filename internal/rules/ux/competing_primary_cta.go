package ux

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// CompetingPrimaryCTARule mendeteksi bilah aksi atau kontainer interaktif yang memuat lebih dari satu
// tombol primary CTA, memecah fokus visual dan memicu kelumpuhan keputusan (Von Restorff Effect & Hick-Hyman Law).
type CompetingPrimaryCTARule struct{}

// NewCompetingPrimaryCTARule membuat instance baru dari CompetingPrimaryCTARule.
func NewCompetingPrimaryCTARule() *CompetingPrimaryCTARule {
	return &CompetingPrimaryCTARule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *CompetingPrimaryCTARule) ID() string {
	return "ux.competing-primary-cta"
}

// Description mengembalikan ringkasan aturan.
func (r *CompetingPrimaryCTARule) Description() string {
	return "Warns when an action group or interactive container contains more than one primary call-to-action button"
}

// Category mengembalikan nama kategori rule.
func (r *CompetingPrimaryCTARule) Category() string {
	return "ux"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *CompetingPrimaryCTARule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *CompetingPrimaryCTARule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Von Restorff Effect (The Isolation Effect / Visual Dominance)",
			"Hick-Hyman Law (Logarithmic Decision Latency)",
			"Nielsen Norman Group (Visual Hierarchy for Action Buttons)",
		},
		CoreInvariant: "An action container or button group must contain at most one primary call-to-action button, ensuring a clear visual focal point and zero decision ambiguity.",
		Grounding: "The Von Restorff Effect (Isolation Effect) predicts that when multiple similar items are presented, the one that differs from the rest is most likely to be remembered and acted upon.\n\n" +
			"When an interface presents two or more buttons styled identically as primary actions (e.g., two 'bg-primary' or 'variant=\"primary\"' buttons side by side in a modal footer or form actions), visual hierarchy collapses.\n\n" +
			"This competing prominence causes choice paralysis (Hick-Hyman Law), forces users to pause and re-read labels carefully, and drastically increases the probability of accidental mis-clicks. Every decision context must have exactly one visually distinct primary action, while supporting actions should be styled with 'outline', 'secondary', or 'ghost' variants.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Choice Paralysis & Decision Latency",
				Severity: "HIGH",
				Impact:   "Users hesitate when confronted with equal-weight visual cues, increasing conversion drop-off rates.",
			},
			{
				Vector:   "Accidental Action Slips",
				Severity: "MEDIUM",
				Impact:   "Users mistake secondary or cancel triggers for primary confirmation due to identical color and elevation styling.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Dialog footer with two competing primary buttons creating visual ambiguity",
				Code: `<div className="flex justify-end gap-3 p-4">
  <button type="button" className="bg-primary text-white px-4 py-2 rounded-md">Simpan Draf</button>
  <button type="submit" className="bg-primary text-white px-4 py-2 rounded-md">Publikasikan</button>
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Clear hierarchy: one primary button paired with a secondary outline button",
				Code: `<div className="flex justify-end gap-3 p-4">
  <button type="button" className="border border-input bg-transparent px-4 py-2 rounded-md">Simpan Draf</button>
  <button type="submit" className="bg-primary text-white px-4 py-2 rounded-md">Publikasikan</button>
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah grup aksi memuat lebih dari satu primary CTA button.
func (r *CompetingPrimaryCTARule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if !isActionContainer(node) {
		return nil
	}

	primaryCount := 0
	for _, child := range node.Children {
		if child.Type != ir.NodeElement {
			continue
		}
		if isPrimaryCTAButton(child) {
			primaryCount++
		}
	}

	if primaryCount > 1 {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message: fmt.Sprintf(
					"Action container contains %d competing primary CTA buttons. In accordance with the Von Restorff Effect, an action group must present at most 1 primary action; demote supporting actions to 'outline', 'secondary', or 'ghost'.",
					primaryCount,
				),
				Hint: "Keep only 1 primary action button and style the others as secondary or outline variants.",
			},
		}
	}

	return nil
}
