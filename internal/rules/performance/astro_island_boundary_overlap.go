package performance

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// AstroIslandBoundaryOverlapRule mengaudit penyarangan pulau interaktif multi-framework tanpa isolasi slot.
type AstroIslandBoundaryOverlapRule struct{}

// NewAstroIslandBoundaryOverlapRule membuat instance baru dari AstroIslandBoundaryOverlapRule.
func NewAstroIslandBoundaryOverlapRule() *AstroIslandBoundaryOverlapRule {
	return &AstroIslandBoundaryOverlapRule{}
}

// ID mengembalikan identifier unik kanonikal aturan.
func (r *AstroIslandBoundaryOverlapRule) ID() string {
	return "performance.astro-island-boundary-overlap"
}

// Category mengembalikan kategori aturan ('performance').
func (r *AstroIslandBoundaryOverlapRule) Category() string {
	return "performance"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *AstroIslandBoundaryOverlapRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Description mengembalikan deskripsi tujuan dan konteks aturan.
func (r *AstroIslandBoundaryOverlapRule) Description() string {
	return "Mencegah konflik batas hidrasi pulau (island boundary overlap) dengan mewajibkan isolasi slot pada penyarangan komponen pulau interaktif."
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *AstroIslandBoundaryOverlapRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Astro Component Composition & Slots Isolation Guidelines",
			"Astro Multi-Framework Islands Architecture Invariants",
			"W3C Web Components Hydration Boundary Isolation Standards",
		},
		CoreInvariant: "Interactive Astro islands must not directly nest secondary client islands without Astro slot isolation; direct nesting blurs hydration boundaries and triggers runtime desynchronization.",
		Grounding: "Astro islands are meant to be isolated units of interactivity that hydrate independently on the page.\n\n" +
			"When an interactive island (`client:*`) nests another client island directly as a child element, the parent framework's virtual DOM attempts to manage the subtree of the child framework.\n\n" +
			"This direct nesting causes hydration mismatches, duplicate runtime overhead, and event listener conflicts. Using Astro Slots (`<div slot=\"...\">`) preserves clear HTML boundaries between distinct hydration contexts.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Hydration Mismatch & Failure",
				Severity: "HIGH",
				Impact:   "Parent framework virtual DOM reconciliation overwrites or destroys DOM nodes managed by child islands.",
			},
			{
				Vector:   "Duplicated Runtime Overhead",
				Severity: "MEDIUM",
				Impact:   "Forces multiple distinct UI framework engines to run in overlapping memory spaces on the client.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Penyarangan pulau multi-framework langsung memicu konflik hidrasi",
				Code: `<!-- Pelanggaran: Penyarangan pulau langsung -->
<ReactDashboardContainer client:load>
  <VueChartWidget client:idle />
</ReactDashboardContainer>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Memanfaatkan Astro slot untuk mengisolasi batas hidrasi",
				Code: `<!-- Patuh: Memisahkan pulau via slot terisolasi -->
<ReactDashboardContainer client:load>
  <div slot="chart-slot">
    <VueChartWidget client:idle />
  </div>
</ReactDashboardContainer>`,
			},
		},
	}
}

// Evaluate memeriksa apakah terdapat pulau anak yang disarangkan langsung di dalam pulau induk.
func (r *AstroIslandBoundaryOverlapRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || len(node.Children) == 0 {
		return nil
	}

	childIsland, childDir, hasOverlap := findNestedIslandOverlap(node)
	if !hasOverlap {
		return nil
	}

	parentDir, _ := hasClientDirective(node)
	return []ir.Diagnostic{
		{
			Line:     childIsland.Span.Line,
			Column:   childIsland.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Nested hydration island '<%s>' (%s) directly inside parent island '<%s>' (%s) without slot isolation causes hydration boundary overlap and runtime breakdown.", childIsland.Tag, childDir, node.Tag, parentDir),
			Hint:     fmt.Sprintf("Isolate '<%s>' using an Astro slot (e.g. '<div slot=\"...\"><%s ... /></div>') to preserve independent hydration boundaries.", childIsland.Tag, childIsland.Tag),
		},
	}
}
