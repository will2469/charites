package inp

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// ExpensiveStyleMutationRule mendeteksi manipulasi gaya imperatif pada properti peka-cat berbiaya tinggi
// (seperti boxShadow, filter, backdropFilter) di dalam penangan interaksi kontinu (onPointerMove, onTouchMove, onScroll).
type ExpensiveStyleMutationRule struct{}

// NewExpensiveStyleMutationRule membuat instance baru dari ExpensiveStyleMutationRule.
func NewExpensiveStyleMutationRule() *ExpensiveStyleMutationRule {
	return &ExpensiveStyleMutationRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *ExpensiveStyleMutationRule) ID() string {
	return "inp.expensive-style-mutation"
}

// Description mengembalikan ringkasan aturan.
func (r *ExpensiveStyleMutationRule) Description() string {
	return "Continuous interaction handler imperatively mutates high-cost paint-sensitive style properties (boxShadow, filter, etc.)"
}

// Category mengembalikan nama kategori rule.
func (r *ExpensiveStyleMutationRule) Category() string {
	return "inp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *ExpensiveStyleMutationRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ExpensiveStyleMutationRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Google Chrome Core Web Vitals (Interaction to Next Paint Presentation Delay)",
			"W3C CSS Compositing and Blending Level 2",
			"Hardware-Accelerated CSS Transforms & Opacity Subsystem",
		},
		CoreInvariant: "Continuous interaction handlers (onPointerMove, onTouchMove, onScroll) must not imperatively mutate high-cost paint-sensitive CSS properties ('boxShadow', 'filter', 'backdropFilter', etc.); GPU-accelerated composited properties ('transform', 'opacity') should be used instead.",
		Grounding: "Properties such as 'box-shadow', 'filter', 'backdrop-filter', and 'background-image' require software or GPU rasterization passes every time their values change.\n\n" +
			"When mutated inside high-frequency continuous interaction handlers (e.g. 'onPointerMove', 'onTouchMove', or 'onScroll' which fire at 60Hz-120Hz), the browser is forced to discard rasterized layer caches and repaint damaged regions continuously.\n\n" +
			"This raster contention causes heavy frame drops and delays Presentation Delay. Replacing dynamic shadow or blur mutations with GPU-composited 'transform' or discrete CSS class toggles avoids CPU/GPU raster churn completely.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Continuous Paint Cache Invalidation",
				Severity: "HIGH",
				Impact:   "High-frequency pointer movements force continual repainting of heavy blur or shadow layers.",
			},
			{
				Vector:   "Frame Drops & Touch Presentation Delay",
				Severity: "HIGH",
				Impact:   "Rasterizer stalls degrade input responsiveness and drop interaction frames down to 15-30 FPS on mobile.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Imperative box-shadow and blur mutation on every pointer move event",
				Code: `<div onPointerMove={(e) => {
  e.currentTarget.style.boxShadow = ` + "`0 ${e.clientY / 10}px 30px rgba(0,0,0,0.5)`" + `;
  e.currentTarget.style.filter = ` + "`blur(${e.clientX / 50}px)`" + `;
}}>
  Interactive Card
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "GPU-accelerated transform without triggering rasterization cache invalidation",
				Code: `<div onPointerMove={(e) => {
  e.currentTarget.style.transform = ` + "`translateY(${e.clientY / 10}px)`" + `;
}}>
  Interactive Card
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah penangan interaksi kontinu memutasi properti peka-cat berbiaya tinggi.
func (r *ExpensiveStyleMutationRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	handler, prop, detected := findExpensiveStyleMutation(node.Tag, node.Attributes)
	if detected {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Continuous interaction handler '%s' imperatively mutates paint-sensitive property '%s'. Mutating raster-heavy properties at 60Hz-120Hz forces continual raster cache invalidation, severely degrading Presentation Delay.", handler, prop),
				Hint:     "Replace dynamic shadow/filter mutations with GPU-accelerated composited properties like 'transform' or discrete CSS class toggles.",
			},
		}
	}

	return nil
}
