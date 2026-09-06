package ux

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// SpacingInversionRule mendeteksi inversi hierarki spasi Gestalt di mana spasi internal anak
// lebih renggang atau sama dengan gap pemisah parent, serta tabrakan spesifisitas Tailwind v3 space-y.
type SpacingInversionRule struct{}

// NewSpacingInversionRule membuat instance baru dari SpacingInversionRule.
func NewSpacingInversionRule() *SpacingInversionRule {
	return &SpacingInversionRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *SpacingInversionRule) ID() string {
	return "ux.spacing-inversion"
}

// Description mengembalikan ringkasan aturan.
func (r *SpacingInversionRule) Description() string {
	return "Warns when child element intra-spacing exceeds parent gap or when space-y conflicts with child mt margin in Tailwind v3"
}

// Category mengembalikan nama kategori rule.
func (r *SpacingInversionRule) Category() string {
	return "ux"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *SpacingInversionRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *SpacingInversionRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Gestalt Law of Proximity (Visual Perceptual Hierarchy)",
			"Tailwind CSS v3 Space-Between Sibling Selector Specificity Quirks",
			"W3C Design Tokens Community Group (DTCG v2025.10 - Spatial Scale)",
		},
		CoreInvariant: "Child element intra-spacing must be strictly tighter than the inter-element gap separating parent sibling groups, and parent 'space-y-*' must not conflict with child 'mt-*' overrides.",
		Grounding: "According to the Gestalt Law of Proximity, elements that belong to the same logical group must have smaller internal spacing than the boundary gap between distinct sibling groups.\n\n" +
			"When a child card or section applies an internal margin or gap that is larger than or equal to the parent container's gap (e.g., parent has 'space-y-4' while child has 'mb-8'), the visual cohesion dissolves, confusing users about which headings, texts, or actions belong together.\n\n" +
			"Furthermore, in Tailwind CSS v3, 'space-y-*' generates a complex sibling selector '> :not([hidden]) ~ :not([hidden])' with specificity (0, 3, 0), which silently overrides any child 'mt-*' utility (0, 1, 0) without a compiler error. Switching the parent to 'flex flex-col gap-*' restores deterministic CSS cascade and spatial clarity.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Cognitive Grouping Disruption",
				Severity: "MEDIUM",
				Impact:   "Users misattribute subheadings and actions to unrelated neighbouring cards due to broken proximity cues.",
			},
			{
				Vector:   "Silent CSS Specificity Override",
				Severity: "MEDIUM",
				Impact:   "Tailwind v3 sibling selectors override child margins without error, leading to unintended layout shifts.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Parent uses space-y-4 while child card specifies mb-8, causing Gestalt proximity inversion and v3 specificity clash",
				Code: `<section className="space-y-4">
  <div className="mb-8">
    <h3 className="text-sm font-semibold">Grup A</h3>
  </div>
  <div className="mb-8">
    <h3 className="text-sm font-semibold">Grup B</h3>
  </div>
</section>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Parent sets wider gap-8 separating groups, while child maintains tighter mb-3 intra-spacing",
				Code: `<section className="flex flex-col gap-8">
  <div className="mb-3">
    <h3 className="text-sm font-semibold">Grup A</h3>
  </div>
  <div className="mb-3">
    <h3 className="text-sm font-semibold">Grup B</h3>
  </div>
</section>`,
			},
		},
	}
}

// Evaluate menganalisis node untuk mendeteksi tabrakan spesifisitas v3 dan inversi proximity Gestalt.
func (r *SpacingInversionRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	var diags []ir.Diagnostic
	diags = r.checkTailwindV3Collision(node, diags)
	diags = r.checkGestaltProximityInversion(node, diags)
	return diags
}

func (r *SpacingInversionRule) checkTailwindV3Collision(node *ir.Node, diags []ir.Diagnostic) []ir.Diagnostic {
	spaceCls, hasSpace := hasTailwindV3SpaceY(node.Classes)
	if !hasSpace {
		return diags
	}

	for _, child := range node.Children {
		if child.Type != ir.NodeElement {
			continue
		}
		if childCls, _, found := extractChildMarginTop(child.Classes); found {
			diags = append(diags, ir.Diagnostic{
				Line:     child.Span.Line,
				Column:   child.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message: fmt.Sprintf(
					"Tailwind v3 specificity collision: parent uses '%s' which overrides child's '%s' via higher selector specificity (0, 3, 0 vs 0, 1, 0). Use 'flex flex-col gap-*' on parent instead.",
					spaceCls, childCls,
				),
				Hint: "Replace parent 'space-y-*' with 'flex flex-col gap-*' to avoid sibling selector collision.",
			})
		}
	}
	return diags
}

func (r *SpacingInversionRule) checkGestaltProximityInversion(node *ir.Node, diags []ir.Diagnostic) []ir.Diagnostic {
	parentGap, hasGap := extractVerticalGapParent(node.Classes)
	if !hasGap || parentGap <= 0 {
		return diags
	}

	for _, child := range node.Children {
		if child.Type != ir.NodeElement {
			continue
		}
		childCls, childSpacing, found := extractChildIntraSpacing(child.Classes)
		if !found || childSpacing < parentGap {
			continue
		}
		if isAlreadyFlagged(diags, child.Span.Line, child.Span.Column) {
			continue
		}
		diags = append(diags, ir.Diagnostic{
			Line:     child.Span.Line,
			Column:   child.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message: fmt.Sprintf(
				"Gestalt proximity inversion: child intra-spacing '%s' (%.2f rem) is greater than or equal to parent gap (%.2f rem), breaking perceptual hierarchy.",
				childCls, childSpacing, parentGap,
			),
			Hint: "Ensure child intra-spacing is strictly smaller than the parent gap separating sibling groups.",
		})
	}
	return diags
}

func isAlreadyFlagged(diags []ir.Diagnostic, line, col int) bool {
	for _, d := range diags {
		if d.Line == line && d.Column == col {
			return true
		}
	}
	return false
}
