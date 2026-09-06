package lcp

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// CriticalHeadStyleBloatRule mendeteksi blok <style> inline di dalam <head>
// yang menyuntikkan CSS non-kritis (seperti footer, modal, dialog) sehingga memperbesar dokumen awal.
type CriticalHeadStyleBloatRule struct{}

// NewCriticalHeadStyleBloatRule membuat instance baru dari CriticalHeadStyleBloatRule.
func NewCriticalHeadStyleBloatRule() *CriticalHeadStyleBloatRule {
	return &CriticalHeadStyleBloatRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *CriticalHeadStyleBloatRule) ID() string {
	return "lcp.critical-head-style-bloat"
}

// Description mengembalikan ringkasan aturan.
func (r *CriticalHeadStyleBloatRule) Description() string {
	return "Inline '<style>' in '<head>' contains non-critical CSS selectors (footer, modal, dialog), inflating initial HTML payload and delaying LCP paint"
}

// Category mengembalikan nama kategori rule.
func (r *CriticalHeadStyleBloatRule) Category() string {
	return "lcp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *CriticalHeadStyleBloatRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *CriticalHeadStyleBloatRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Google Chrome Core Web Vitals (Largest Contentful Paint Time to First Byte & Render Delay)",
			"W3C CSS Cascading and Inheritance Level 5",
			"Web Performance Working Group Critical CSS Separation Invariants",
		},
		CoreInvariant: "Inline '<style>' tags inside the document '<head>' should contain only essential Critical CSS required to render the above-the-fold viewport; non-critical styles must be extracted to cacheable external stylesheets.",
		Grounding: "Inlining Critical CSS directly into the HTML `<head>` is an established optimization to eliminate the render-blocking stylesheet network round-trip for initial viewport elements.\n\n" +
			"However, when monolithic application styles (such as footer links, modal overlays, dialog drawers, and below-the-fold widgets) are bundled indiscriminately into `<head>` styles, the initial HTML payload balloons in size.\n\n" +
			"Because HTML is streamed over TCP in 14KB chunks, bloated inline styles consume early round-trips before the browser even discovers hero `<img>` or heading elements, inflating TTFB and Element Render Delay without the benefit of browser HTTP caching.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Initial HTML Payload Bloat",
				Severity: "HIGH",
				Impact:   "Inflates HTML document transfer size, exhausting early TCP slow-start congestion windows before hero media tags are parsed.",
			},
			{
				Vector:   "Loss of HTTP Caching Efficiency",
				Severity: "MEDIUM",
				Impact:   "Inline CSS cannot be cached by the browser cache or CDN edge nodes across subsequent page navigations.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "html",
				Comment:  "Monolithic style in head bundling footer and modal CSS rules",
				Code: `<head>
  <style>
    .footer-links { color: #6b7280; font-size: 0.875rem; }
    .admin-modal-overlay { display: none; position: fixed; inset: 0; }
  </style>
</head>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "html",
				Comment:  "Head style strictly limited to above-the-fold critical hero layout",
				Code: `<head>
  <style>
    .hero-container { min-height: 480px; display: flex; }
  </style>
  <link rel="stylesheet" href="/assets/main.css" />
</head>`,
			},
		},
	}
}

// Evaluate memeriksa apakah blok style di head menyuntikkan CSS non-kritis berlebih.
func (r *CriticalHeadStyleBloatRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	bloatPattern, hasBloat := isHeadStyleBloat(node)
	if !hasBloat {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Inline '<style>' in '<head>' contains non-critical component CSS ('%s'). Monolithic styles in '<head>' inflate initial HTML document size, delaying stream parsing of hero LCP elements.", bloatPattern),
			Hint:     "Extract below-the-fold component styles (footer, modal, dialog) into an external cacheable stylesheet (e.g. '<link rel=\"stylesheet\" href=\"...\">') and reserve '<head>' styles strictly for initial viewport layout.",
		},
	}
}
