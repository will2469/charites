package inp

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// HydrationHeavyIslandRule mendeteksi komponen pulau client yang membungkus sub-elemen
// statis masif alih-alih memanfaatkan arsitektur Zero-JS Astro SSR.
type HydrationHeavyIslandRule struct{}

// NewHydrationHeavyIslandRule membuat instance baru dari HydrationHeavyIslandRule.
func NewHydrationHeavyIslandRule() *HydrationHeavyIslandRule {
	return &HydrationHeavyIslandRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *HydrationHeavyIslandRule) ID() string {
	return "inp.hydration-heavy-island"
}

// Description mengembalikan ringkasan aturan.
func (r *HydrationHeavyIslandRule) Description() string {
	return "Client island wraps excessive static DOM subtree forcing heavy virtual DOM reconciliation on the client"
}

// Category mengembalikan nama kategori rule.
func (r *HydrationHeavyIslandRule) Category() string {
	return "inp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *HydrationHeavyIslandRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *HydrationHeavyIslandRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Astro Zero-JS Server-Side Rendering (SSR) Architecture",
			"React Virtual DOM Hydration Complexity & Reconciliation Budget",
			"W3C Web Performance Working Group (Main Thread Input Delay)",
		},
		CoreInvariant: "Client-hydrated islands must remain compact and isolate only truly interactive elements; static text, articles, and decorative containers must remain in zero-JS Astro SSR.",
		Grounding: "Hydrating a monolithic React island forces the client browser to parse JavaScript, construct virtual DOM representations, and reconcile every DOM node against server HTML-even for completely static elements.\n\n" +
			"When developers wrap entire articles or document structures inside a single `<ArticleViewer client:load>`, hundreds of static paragraphs and headings are needlessly reconciled, consuming excessive main-thread CPU time.\n\n" +
			"By decomposing the UI and rendering static content through native Astro components (zero client JS), only individual interactive widgets (such as like buttons or comment inputs) are hydrated, keeping the main thread free for user interaction.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Virtual DOM Reconciliation Bloat",
				Severity: "HIGH",
				Impact:   "Large static subtrees force long synchronous VDOM tree reconciliation during island hydration.",
			},
			{
				Vector:   "Excessive Client Bundle & Parse Overhead",
				Severity: "MEDIUM",
				Impact:   "Shipping static component trees to client bundles increases script evaluation time and input latency.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Static article text wrapped inside a client-hydrated island",
				Code: `<ArticleViewerIsland client:load>
  <header><h1>Article Title</h1></header>
  <article>
    <p>Paragraph 1...</p>
    <p>Paragraph 2...</p>
    <p>Paragraph 3...</p>
  </article>
  <CommentButton />
</ArticleViewerIsland>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Static content rendered via zero-JS Astro SSR; only interactive button is an island",
				Code: `<header><h1>Article Title</h1></header>
<article>
  <p>Paragraph 1...</p>
  <p>Paragraph 2...</p>
  <p>Paragraph 3...</p>
</article>
<CommentButton client:visible />`,
			},
		},
	}
}

// Evaluate memeriksa apakah pulau client membungkus sub-elemen statis berlebihan.
func (r *HydrationHeavyIslandRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	totalNodes, isHeavy := isHydrationHeavyIsland(node)
	if !isHeavy {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Astro island '%s' with client directive wraps an excessively large static DOM subtree (%d nodes). Reconciling static content inside a client island wastes main-thread CPU during hydration, degrading INP responsiveness.", node.Tag, totalNodes),
			Hint:     "Decompose the component to keep static content in zero-JS Astro SSR, and only hydrate interactive leaves (e.g. <CommentButton client:visible />).",
		},
	}
}
