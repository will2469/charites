package lcp

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// UnhintedLCPImagePriorityRule mendeteksi elemen gambar kandidat LCP di dalam aliran HTML awal
// yang tidak menyertakan atribut fetchpriority="high" untuk menginstruksikan browser preload scanner.
type UnhintedLCPImagePriorityRule struct{}

// NewUnhintedLCPImagePriorityRule membuat instance baru dari UnhintedLCPImagePriorityRule.
func NewUnhintedLCPImagePriorityRule() *UnhintedLCPImagePriorityRule {
	return &UnhintedLCPImagePriorityRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *UnhintedLCPImagePriorityRule) ID() string {
	return "lcp.unhinted-lcp-image-priority"
}

// Description mengembalikan ringkasan aturan.
func (r *UnhintedLCPImagePriorityRule) Description() string {
	return "Above-the-fold LCP candidate image lacks fetchpriority=\"high\", delaying bandwidth allocation in early network stream"
}

// Category mengembalikan nama kategori rule.
func (r *UnhintedLCPImagePriorityRule) Category() string {
	return "lcp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *UnhintedLCPImagePriorityRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *UnhintedLCPImagePriorityRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Google Chrome Core Web Vitals (Largest Contentful Paint Resource Load Delay)",
			"W3C Priority Hints Specification (fetchpriority attribute)",
			"Chrome Preload Scanner Network Bandwidth Scheduling",
		},
		CoreInvariant: "Above-the-fold LCP candidate images must declare 'fetchpriority=\"high\"' to prioritize early network bandwidth ahead of non-critical stylesheets and scripts.",
		Grounding: "By default, browsers assign an initial fetch priority of 'Low' to image resources discovered in the HTML stream.\n\n" +
			"For the primary hero image (the LCP element), this default low priority forces the image download to compete with or yield to lower-priority scripts, stylesheets, and fonts.\n\n" +
			"Declaring 'fetchpriority=\"high\"' instructs the speculative preload scanner to immediately elevate the resource to the highest network tier, initiating the TCP/TLS transfer with maximum allocated bandwidth.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Bandwidth Starvation by Non-Critical Assets",
				Severity: "HIGH",
				Impact:   "Hero image bytes are delayed behind non-critical deferred scripts and fonts, inflating LCP by 150ms-400ms.",
			},
			{
				Vector:   "Sub-optimal Browser Network Scheduling",
				Severity: "MEDIUM",
				Impact:   "Browsers under HTTP/2 or HTTP/3 multiplexing prioritize other resources unless explicitly hinted.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Above-the-fold hero banner image lacking priority hint",
				Code: `<header className="hero-banner" data-perf-role="hero">
  <img src="/hero.webp" alt="Primary Banner" className="w-full aspect-video" />
</header>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Hero image explicitly prioritized with fetchpriority='high'",
				Code: `<header className="hero-banner" data-perf-role="hero">
  <img src="/hero.webp" alt="Primary Banner" fetchpriority="high" className="w-full aspect-video" />
</header>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen kandidat LCP kekurangan atribut fetchpriority="high".
func (r *UnhintedLCPImagePriorityRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	isHero, _ := isLCPCandidate(node)
	if !isHero {
		return nil
	}

	// Precedence Guard: Jika gambar masih memiliki loading="lazy", biarkan rule lazy-loaded-lcp-image
	// menangani terlebih dahulu untuk mencegah kombinasi kontradiktif (lazy + high priority).
	if isLazyLoaded(node.Attributes) {
		return nil
	}

	if !hasHighFetchPriority(node.Attributes) {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Above-the-fold LCP candidate '<%s>' lacks 'fetchpriority=\"high\"'. Without an explicit priority hint, the browser preload scanner treats hero media with low priority, prolonging Resource Load Delay.", node.Tag),
				Hint:     "Add 'fetchpriority=\"high\"' (or JSX 'fetchPriority=\"high\"') to prioritize hero image network transfer.",
			},
		}
	}

	return nil
}
