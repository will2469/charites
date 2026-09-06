package lcp

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// MissingLCPImagePreloadRule menyarankan deklarasi <link rel="preload" as="image"> di dalam <head>
// untuk aset gambar kandidat LCP yang mengalami keterlambatan penemuan (delayed discovery via client script atau data attributes).
type MissingLCPImagePreloadRule struct{}

// NewMissingLCPImagePreloadRule membuat instance baru dari MissingLCPImagePreloadRule.
func NewMissingLCPImagePreloadRule() *MissingLCPImagePreloadRule {
	return &MissingLCPImagePreloadRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *MissingLCPImagePreloadRule) ID() string {
	return "lcp.missing-lcp-image-preload"
}

// Description mengembalikan ringkasan aturan.
func (r *MissingLCPImagePreloadRule) Description() string {
	return "Delayed-discovery LCP image lacks <link rel=\"preload\" as=\"image\"> in document head to initiate early asset transfer"
}

// Category mengembalikan nama kategori rule.
func (r *MissingLCPImagePreloadRule) Category() string {
	return "lcp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (info).
func (r *MissingLCPImagePreloadRule) DefaultSeverity() ir.Severity {
	return ir.SeverityInfo
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *MissingLCPImagePreloadRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Google Chrome Core Web Vitals (Largest Contentful Paint Resource Load Delay)",
			"W3C Preload Specification (<link rel=\"preload\" as=\"image\">)",
			"Document Layout Graph & Early Resource Discovery Optimization",
		},
		CoreInvariant: "LCP candidate media elements that suffer delayed discovery (dynamic data attributes, client script hydration, or CSS backgrounds) should be preloaded in '<head>' with 'fetchpriority=\"high\"'.",
		Grounding: "When an LCP image cannot be immediately parsed from a direct '<img>' element in the server-rendered HTML stream (for example, when its URL is stored in a dynamic data attribute 'data-bg-src', rendered by a client island, or defined via CSS background), its network fetch is delayed.\n\n" +
			"Injecting '<link rel=\"preload\" as=\"image\" href=\"...\" fetchpriority=\"high\">' inside '<head>' compensates for this discovery delay by instructing the browser lookahead scanner to initiate connection and download immediately during initial document streaming.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Delayed Resource Discovery",
				Severity: "MEDIUM",
				Impact:   "Image download is postponed until JavaScript hydration or style resolution finishes, inflating LCP.",
			},
			{
				Vector:   "Initial Viewport Flash",
				Severity: "LOW",
				Impact:   "Late arrival of visual hero media causes prolonged empty or placeholder hero appearance.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Hero container with dynamic data-bg-src without preload in document head",
				Code: `<head>
  <title>Product Gallery</title>
</head>
<body>
  <div id="hero-root" data-perf-role="hero" data-bg-src="https://cdn.example.com/promo.webp"></div>
</body>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Document head preloading the hero image asset with high fetch priority",
				Code: `<head>
  <title>Product Gallery</title>
  <link rel="preload" as="image" href="https://cdn.example.com/promo.webp" fetchpriority="high" />
</head>
<body>
  <div id="hero-root" data-perf-role="hero" data-bg-src="https://cdn.example.com/promo.webp"></div>
</body>`,
			},
		},
	}
}

// Evaluate memeriksa apakah aset LCP yang tertunda penemuannya memiliki preload di head.
func (r *MissingLCPImagePreloadRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	url, detected := isDelayedDiscoveryLCP(node)
	if !detected {
		return nil
	}

	if hasPreloadInHead(node, url) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Delayed-discovery LCP image asset '%s' on '<%s>' is not preloaded in document '<head>'. Adding a preload link initiates network download during early HTML streaming.", url, node.Tag),
			Hint:     "Add '<link rel=\"preload\" as=\"image\" href=\"...\" fetchpriority=\"high\" />' in '<head>' to cut Resource Load Delay.",
		},
	}
}
