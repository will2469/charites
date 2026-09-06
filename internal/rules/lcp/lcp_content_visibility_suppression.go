package lcp

import (
	"github.com/will2469/charites/internal/ir"
)

// ContentVisibilitySuppressionRule mendeteksi elemen kontainer seksi hero atau pelipatan awal
// yang menerapkan properti content-visibility: auto atau class Tailwind content-auto.
type ContentVisibilitySuppressionRule struct{}

// NewContentVisibilitySuppressionRule membuat instance baru dari ContentVisibilitySuppressionRule.
func NewContentVisibilitySuppressionRule() *ContentVisibilitySuppressionRule {
	return &ContentVisibilitySuppressionRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *ContentVisibilitySuppressionRule) ID() string {
	return "lcp.lcp-content-visibility-suppression"
}

// Description mengembalikan ringkasan aturan.
func (r *ContentVisibilitySuppressionRule) Description() string {
	return "Above-the-fold hero container specifies 'content-visibility: auto' or 'content-auto', suppressing initial paint and severely inflating LCP"
}

// Category mengembalikan nama kategori rule.
func (r *ContentVisibilitySuppressionRule) Category() string {
	return "lcp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *ContentVisibilitySuppressionRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ContentVisibilitySuppressionRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Google Chrome Core Web Vitals (Largest Contentful Paint Element Render Delay)",
			"W3C CSS Containment Module Level 2 (content-visibility property specification)",
			"Chromium Blink Rendering Pipeline Initial Viewport Invariants",
		},
		CoreInvariant: "Initial viewport hero containers are strictly forbidden from specifying 'content-visibility: auto' or 'content-auto' as it instructs the browser engine to suppress early layout and paint passes.",
		Grounding: "The CSS property 'content-visibility: auto' is a powerful rendering performance feature for below-the-fold content, allowing browsers to skip layout and painting for elements until they approach the viewport.\n\n" +
			"However, when applied to above-the-fold hero elements (such as `<header>`, hero sections, or containers with `data-perf-role=\"hero\"`), the browser initially skips rendering the element entirely during the first layout pass.\n\n" +
			"Only after subsequent scrolling or intersection calculations does the engine lay out and paint the content, resulting in massive Element Render Delay and catastrophically failing LCP benchmarks.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Initial Layout Pass Suppression",
				Severity: "CRITICAL",
				Impact:   "Blink skips the initial layout and paint of the primary hero container, delaying LCP registration until post-hydration intersection observer checks.",
			},
			{
				Vector:   "Blank Initial Viewport Flash",
				Severity: "HIGH",
				Impact:   "Users experience an empty screen or blank space in the initial viewport on fast network connections.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Hero section using content-visibility: auto in the initial viewport",
				Code: `<section className="hero-section content-auto" data-perf-role="hero">
  <h1>Solusi Cloud Enterprise</h1>
  <img src="/hero.webp" fetchpriority="high" />
</section>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Hero section rendered immediately without initial paint suppression",
				Code: `<section className="hero-section" data-perf-role="hero">
  <h1>Solusi Cloud Enterprise</h1>
  <img src="/hero.webp" fetchpriority="high" />
</section>`,
			},
		},
	}
}

// Evaluate memeriksa apakah kontainer hero di pelipatan awal menetapkan content-visibility auto.
func (r *ContentVisibilitySuppressionRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !hasContentVisibilityAuto(node) {
		return nil
	}

	if !isAboveFoldHeroContainer(node) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Initial viewport hero container specifies 'content-visibility: auto' (or 'content-auto'). Applying 'content-visibility: auto' to above-the-fold content instructs the browser to skip early layout and painting, severely delaying LCP.",
			Hint:     "Remove 'content-auto' or 'content-visibility: auto' from initial viewport hero containers and restrict its usage exclusively to below-the-fold content.",
		},
	}
}
