package lcp

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// LazyLoadedLCPImageRule mendeteksi elemen gambar kandidat LCP pelipatan atas yang diberi atribut loading="lazy",
// pola yang menunda inisiasi pengunduhan aset hingga fase tata letak selesai dan merusak metrik LCP.
type LazyLoadedLCPImageRule struct{}

// NewLazyLoadedLCPImageRule membuat instance baru dari LazyLoadedLCPImageRule.
func NewLazyLoadedLCPImageRule() *LazyLoadedLCPImageRule {
	return &LazyLoadedLCPImageRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *LazyLoadedLCPImageRule) ID() string {
	return "lcp.lazy-loaded-lcp-image"
}

// Description mengembalikan ringkasan aturan.
func (r *LazyLoadedLCPImageRule) Description() string {
	return "Critical above-the-fold LCP candidate image has loading=\"lazy\", delaying resource discovery and load initiation"
}

// Category mengembalikan nama kategori rule.
func (r *LazyLoadedLCPImageRule) Category() string {
	return "lcp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *LazyLoadedLCPImageRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *LazyLoadedLCPImageRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Google Chrome Core Web Vitals (Largest Contentful Paint Resource Load Delay)",
			"HTML Living Standard Lazy Loading Specification",
			"W3C Web Performance Working Group Invariants",
		},
		CoreInvariant: "Above-the-fold LCP candidate images must not be configured with loading='lazy'; lazy loading defers image download until layout completion, directly adding hundreds of milliseconds to LCP.",
		Grounding: "When a browser encounters an '<img>' with 'loading=\"lazy\"', it deliberately pauses fetching the image resource until the page layout is calculated and the element is verified to be within or near the viewport.\n\n" +
			"For hero images and above-the-fold content that constitute the Largest Contentful Paint (LCP), this artificial pause wastes the initial network idle period. The browser speculative preload scanner is effectively blocked from fetching the hero asset early.\n\n" +
			"Removing 'loading=\"lazy\"' or declaring 'loading=\"eager\"' combined with 'fetchpriority=\"high\"' allows the browser to initiate the network download immediately upon parsing the HTML token.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Resource Load Delay Inflation",
				Severity: "CRITICAL",
				Impact:   "Hero image download is postponed until stylesheet download, CSS parsing, and layout pass complete, adding 200ms-800ms to LCP.",
			},
			{
				Vector:   "Speculative Preload Scanner Suppression",
				Severity: "HIGH",
				Impact:   "The browser's high-speed HTML lookahead parser skips downloading the hero asset during early stream processing.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Above-the-fold hero banner image configured with loading='lazy'",
				Code: `<section className="hero-section" data-perf-role="hero">
  <h1>Welcome to Our Platform</h1>
  <img src="/assets/hero.webp" alt="Hero Banner" loading="lazy" className="w-full h-auto" />
</section>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Hero image configured with loading='eager' and high fetch priority",
				Code: `<section className="hero-section" data-perf-role="hero">
  <h1>Welcome to Our Platform</h1>
  <img src="/assets/hero.webp" alt="Hero Banner" loading="eager" fetchpriority="high" className="w-full h-auto" />
</section>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen kandidat LCP memiliki atribut loading="lazy".
func (r *LazyLoadedLCPImageRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	isHero, _ := isLCPCandidate(node)
	if !isHero {
		return nil
	}

	if isLazyLoaded(node.Attributes) {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Above-the-fold LCP candidate '<%s>' has 'loading=\"lazy\"'. Lazy loading hero media delays resource fetch initiation until layout calculation completes, severely inflating LCP.", node.Tag),
				Hint:     "Remove 'loading=\"lazy\"' or change to 'loading=\"eager\"' and add 'fetchpriority=\"high\"'.",
			},
		}
	}

	return nil
}
