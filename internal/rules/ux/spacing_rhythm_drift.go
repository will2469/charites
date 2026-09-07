package ux

import (
	"github.com/will2469/charites/internal/ir"
)

// SpacingRhythmDriftRule mendeteksi sekuens spasi yang membentuk ritme lokal tidak konsisten
// di dalam satu grup tata letak (layout group) yang sama.
type SpacingRhythmDriftRule struct{}

// NewSpacingRhythmDriftRule membuat instance baru dari SpacingRhythmDriftRule.
func NewSpacingRhythmDriftRule() *SpacingRhythmDriftRule {
	return &SpacingRhythmDriftRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *SpacingRhythmDriftRule) ID() string {
	return "ux.spacing-rhythm-drift"
}

// Description mengembalikan ringkasan aturan.
func (r *SpacingRhythmDriftRule) Description() string {
	return "Detects spacing sequences that form an inconsistent local rhythm within the same layout group"
}

// Category mengembalikan nama kategori rule.
func (r *SpacingRhythmDriftRule) Category() string {
	return "ux"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *SpacingRhythmDriftRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *SpacingRhythmDriftRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Gestalt Law of Proximity & Visual Rhythm",
			"W3C Design Tokens Community Group (DTCG v2025.10 - Spatial Cadence)",
			"Evidence-Based Static Analysis for Design Systems",
		},
		CoreInvariant: "Spacing values serving the same spatial role within the same layout group should normally come from a coherent local spacing family. When a dominant local rhythm is established (>= 3 occurrences, dominance ratio >= 60%), isolated unexplained outliers are reported as warning.",
		Grounding: "In modern user interface engineering, consistent spacing intervals between peer components establish visual cadence, cognitive grouping, and predictable layout rhythm.\n\n" +
			"When a layout container contains repeated sibling relationships with homogeneous spacing, but one relationship unexpectedly deviates (e.g. mb-4, mb-4, mb-7, mb-4), users perceive visual jarring and broken alignment.\n\n" +
			"Rather than enforcing arbitrary universal moduli (e.g. strict 4px or 8px grid), ux.spacing-rhythm-drift evaluates local relational consistency within verified layout groups, reporting isolated unexplained outliers while remaining immune to intentional responsive variants, distinct semantic roles, and small sample sizes.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Visual Cadence Disruption",
				Severity: "MEDIUM",
				Impact:   "Irregular spacing among peer items erodes user trust and creates visual jarring without intentional hierarchy.",
			},
			{
				Vector:   "Cognitive Grouping Confusion",
				Severity: "LOW",
				Impact:   "Users misperceive arbitrary spacing outliers as deliberate semantic distinctions or broken layout states.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Repeated form fields use uniform mb-4 spacing, but one field introduces an unexplained mb-7 outlier",
				Code: `<div>
  <Field className="mb-4" />
  <Field className="mb-4" />
  <Field className="mb-7" />
  <Field className="mb-4" />
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Repeated form fields adhere strictly to a coherent local cadence with homogeneous mb-4 spacing",
				Code: `<div>
  <Field className="mb-4" />
  <Field className="mb-4" />
  <Field className="mb-4" />
  <Field className="mb-4" />
</div>`,
			},
		},
	}
}

// Evaluate menganalisis node untuk mendeteksi drift ritme spasi lokal pada layout owner.
func (r *SpacingRhythmDriftRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	kind := ClassifyLayoutOwner(node)
	if kind == LayoutOwnerNone {
		return nil
	}

	groups := BuildRhythmGroups(node, kind)
	if len(groups) == 0 {
		return nil
	}

	return DetectRhythmDrift(groups, DefaultRhythmConfig())
}
