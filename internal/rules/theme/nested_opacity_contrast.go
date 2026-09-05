package theme

import (
	"github.com/will2469/charites/internal/ir"
)

// NestedOpacityContrastRule mendeteksi penumpukan hierarkis modifier transparansi/opacity
// yang melipatgandakan penurunan kontras hingga di bawah ambang batas WCAG AA (4.5:1).
type NestedOpacityContrastRule struct{}

// NewNestedOpacityContrastRule membuat instance baru NestedOpacityContrastRule.
func NewNestedOpacityContrastRule() *NestedOpacityContrastRule {
	return &NestedOpacityContrastRule{}
}

// ID mengembalikan Charites Rule ID kanonikal.
func (r *NestedOpacityContrastRule) ID() string {
	return "theme.nested-opacity-contrast"
}

// Description mengembalikan penjelasan ringkas rule.
func (r *NestedOpacityContrastRule) Description() string {
	return "Detects nested opacity modifiers that compound to cause catastrophic text contrast degradation"
}

// Category mengembalikan nama kategori rule.
func (r *NestedOpacityContrastRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *NestedOpacityContrastRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *NestedOpacityContrastRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"WCAG 2.2 Success Criterion 1.4.3 (Contrast Minimum - 4.5:1)",
			"W3C DTCG State & Opacity Token Architecture",
			"Hardware-Accelerated Compositing & Alpha Blending",
		},
		CoreInvariant: "Containers with opacity or semi-transparent backgrounds must not enclose child elements with compounded opacity modifiers.",
		Grounding: "When a parent container declares opacity (e.g. opacity-80 or bg-muted/40) and encloses child text or elements with another opacity modifier (e.g. text-foreground/50 or opacity-60), the browser multiplies effective alpha channels (0.8 × 0.5 = 0.40):\n\n" +
			"1. WCAG Contrast Catastrophe: Text that was theoretically compliant plummets below 2.5:1 contrast against the surface.\n" +
			"2. Inverted Washed-Out Appearance: Nested semi-transparency produces muddy, unreadable grey layers in dark mode.\n" +
			"3. Unpredictable Compositing: Nested opacity triggers extra GPU compositing passes and subpixel rendering degradation.\n\n" +
			"Charites enforces using pre-calibrated solid semantic tokens (e.g. text-muted-foreground instead of compounding text-foreground/50 over an opacity container).",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Container opacity compounded with child slash opacity in TSX",
				Code: `<div className="bg-muted/40 opacity-80">
  <p className="text-foreground/50">Notice</p>
</div>`,
			},
			{
				Language: "astro",
				Comment:  "Nested opacity on parent and child text in Astro",
				Code: `<section class="opacity-75">
  <span class="text-white/60">Subtle</span>
</section>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Using solid semantic container and pre-calibrated muted text token",
				Code: `<div className="bg-muted">
  <p className="text-muted-foreground">Notice</p>
</div>`,
			},
			{
				Language: "astro",
				Comment:  "Solid background token and semantic foreground",
				Code: `<section class="bg-card">
  <span class="text-foreground">Subtle</span>
</section>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Multiplicative Alpha Collapse",
				Severity: "HIGH",
				Impact:   "Compounded opacity causes text contrast to fail WCAG AA 4.5:1 accessibility requirements.",
			},
			{
				Vector:   "Compositing Performance Overhead",
				Severity: "LOW",
				Impact:   "Nested alpha layers force browser rasterization pipelines into multiple offscreen passes.",
			},
		},
	}
}

// Evaluate mengevaluasi sebuah node IR untuk mendeteksi penumpukan modifier opacity.
func (r *NestedOpacityContrastRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || len(node.Children) == 0 || len(node.Classes) == 0 {
		return nil
	}

	containerOpacity, hasContainer := HasContainerOpacity(node.Classes)
	if !hasContainer {
		return nil
	}

	var diags []ir.Diagnostic
	for _, child := range node.Children {
		if child == nil || len(child.Classes) == 0 {
			continue
		}
		if childOpacity, hasChild := HasTextOrInnerOpacity(child.Classes); hasChild {
			diags = append(diags, ir.Diagnostic{
				Line:     child.Span.Line,
				Column:   child.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Nested opacity contrast degradation: container has \"" + containerOpacity + "\" and child has \"" + childOpacity + "\"",
				Hint:     "Avoid compounding opacity. Use a solid semantic background and text-muted-foreground instead.",
			})
		}
	}

	return diags
}
