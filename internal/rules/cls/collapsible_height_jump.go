package cls

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// CollapsibleHeightJumpRule mendeteksi komponen akordeon atau menu lipat yang menganimasikan
// dimensi max-height secara arbitrer (max-h-[1000px] / max-h-0) dengan transition-all,
// alih-alih menggunakan teknik CSS Grid zero-shift yang mulus (grid-template-rows: 0fr -> 1fr).
type CollapsibleHeightJumpRule struct{}

// NewCollapsibleHeightJumpRule membuat instance baru dari CollapsibleHeightJumpRule.
func NewCollapsibleHeightJumpRule() *CollapsibleHeightJumpRule {
	return &CollapsibleHeightJumpRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *CollapsibleHeightJumpRule) ID() string {
	return "cls.collapsible-height-jump"
}

// Description mengembalikan ringkasan aturan.
func (r *CollapsibleHeightJumpRule) Description() string {
	return "Collapsible accordion or drawer animates arbitrary max-height bounds instead of zero-shift CSS Grid"
}

// Category mengembalikan nama kategori rule.
func (r *CollapsibleHeightJumpRule) Category() string {
	return "cls"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *CollapsibleHeightJumpRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *CollapsibleHeightJumpRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"CSS Grid Module Level 3 (grid-template-rows interpolation)",
			"Google Core Web Vitals (Interactive Animation CLS Invariants)",
			"Modern Zero-Shift Accordion Architectural Standards",
		},
		CoreInvariant: "Collapsible content drawers and accordions must avoid animating arbitrary max-height bounds and instead adopt zero-shift CSS Grid (grid-template-rows: 0fr -> 1fr).",
		Grounding: "A common legacy technique for animating collapsible elements involves transitioning 'max-height' from 0 to an arbitrarily large value (such as 'max-h-[1000px]').\n\n" +
			"Because CSS transitions interpolate over the declared boundary rather than actual content height, the animation duration becomes severely distorted: closing appears delayed and snapping occurs at the end of the transition, forcing layout reflow on surrounding elements.\n\n" +
			"The modern zero-shift solution utilizes CSS Grid: '<div class=\"grid transition-[grid-template-rows] duration-300 grid-rows-[0fr]\"><div class=\"overflow-hidden\">...</div></div>'. This allows CSS to smoothly interpolate intrinsic content height from 0fr to 1fr without any duration distortion or layout jumps.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Duration Distortion & Layout Snapping",
				Severity: "MEDIUM",
				Impact:   "Collapsing animations finish before the transition duration elapses, causing abrupt layout snaps and visual hitching.",
			},
			{
				Vector:   "Continuous Main-Thread Reflow During Accordion Expansion",
				Severity: "MEDIUM",
				Impact:   "Transitioning max-height triggers continuous layout passes across all frames during accordion interactions.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Accordion drawer animating arbitrary max-height bounds",
				Code: `<div className="transition-all duration-300 overflow-hidden max-h-[1000px]">
  <AccordionBody />
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Modern zero-shift CSS Grid accordion interpolation",
				Code: `<div className="grid transition-[grid-template-rows] duration-300 grid-rows-[1fr]">
  <div className="overflow-hidden">
    <AccordionBody />
  </div>
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen menganimasikan max-height arbitrer alih-alih teknik CSS Grid zero-shift.
func (r *CollapsibleHeightJumpRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	offending, hasRisk := isCollapsibleHeightAnimation(node)
	if !hasRisk {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Collapsible element animates layout dimensions using '%s' with arbitrary max-height bounds. Transitioning arbitrary max-height causes duration distortion, layout snapping, and Cumulative Layout Shift (CLS). Use CSS Grid zero-shift interpolation ('grid-template-rows: 0fr -> 1fr') instead.", offending),
			Hint:     "Wrap collapsible content in '<div class=\"grid transition-[grid-template-rows] duration-300 grid-rows-[0fr]\"><div class=\"overflow-hidden\">...</div></div>' or use modern 'interpolate-size: allow-keywords'.",
		},
	}
}
