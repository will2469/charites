package responsive

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// FixedWidthOverflowRule mendeteksi lebar kontainer statis (w-[...px] atau min-w-[...px])
// yang melebihi batas fisik layar smartphone (320px) tanpa pembatas fluida (max-w-full).
type FixedWidthOverflowRule struct{}

// NewFixedWidthOverflowRule membuat instance baru dari FixedWidthOverflowRule.
func NewFixedWidthOverflowRule() *FixedWidthOverflowRule {
	return &FixedWidthOverflowRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *FixedWidthOverflowRule) ID() string {
	return "responsive.fixed-width-overflow"
}

// Description mengembalikan ringkasan aturan.
func (r *FixedWidthOverflowRule) Description() string {
	return "Detects static fixed container widths exceeding 320px that cause horizontal overflow on mobile viewports"
}

// Category mengembalikan nama kategori rule.
func (r *FixedWidthOverflowRule) Category() string {
	return "responsive"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *FixedWidthOverflowRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *FixedWidthOverflowRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C CSS Box Sizing & Fluid Layout Standards",
			"Mobile-First Responsive Layout Dimensions (320px Minimum Screen Width)",
			"Tailwind CSS Arbitrary Values & Constrained Width Geometry",
		},
		CoreInvariant: "Static widths and min-widths exceeding 320px on mobile baseline must be bounded by fluid constraints (max-w-full) or scoped to desktop breakpoints.",
		Grounding: "Compact and foldable smartphones feature viewport widths starting at 320px (e.g. early iPhone SE or folded Galaxy Z Fold).\n\n" +
			"Declaring rigid static widths greater than 320px (such as w-[500px] or min-w-[400px]) directly on the mobile baseline mechanically exceeds the physical screen boundaries, causing the page to tear and creating an unwanted horizontal scrollbar.\n\n" +
			"Using fluid widths with maximum caps (w-full max-w-lg) ensures full responsiveness across all screen dimensions.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Mobile Horizontal Layout Tear",
				Severity: "HIGH",
				Impact:   "Container forces document width beyond screen borders, creating horizontal scrolling and broken edge-swipe gestures.",
			},
			{
				Vector:   "Cutoff Touch Targets",
				Severity: "MEDIUM",
				Impact:   "Buttons on the right edge of the fixed container become inaccessible without panning horizontally.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Static fixed width exceeding 320px without fluid boundary",
				Code: `<div className="w-[500px] bg-card p-4">
  <p>Kartu Informasi Desa</p>
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Fluid mobile width with max-width ceiling on larger screens",
				Code: `<div className="w-full max-w-lg bg-card p-4">
  <p>Kartu Informasi Desa</p>
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen mendeklarasikan lebar statis > 320px tanpa pembatas fluida.
func (r *FixedWidthOverflowRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || len(node.Classes) == 0 {
		return nil
	}

	if hasFluidWidthBoundary(node.Classes) {
		return nil
	}

	var diags []ir.Diagnostic
	for _, cls := range node.Classes {
		px, ok := extractStaticPixelWidth(cls)
		if ok && px > 320 {
			diags = append(diags, ir.Diagnostic{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Static container width %q (%dpx) exceeds the minimum smartphone viewport width (320px) without a fluid boundary (max-w-full). This triggers unconstrained horizontal overflow and breaks mobile swipe gestures.", cls, px),
				Hint:     "Replace static width with fluid width constraints, e.g. 'w-full max-w-lg' or scope with breakpoint modifiers 'w-full md:w-[500px]'.",
			})
		}
	}

	return diags
}
