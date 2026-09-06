package lcp

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// MissingCriticalOriginHintRule mendeteksi aset visual kandidat LCP kritis
// yang dimuat dari origin eksternal tanpa adanya tag <link rel="preconnect"> pada dokumen.
type MissingCriticalOriginHintRule struct{}

// NewMissingCriticalOriginHintRule membuat instance baru dari MissingCriticalOriginHintRule.
func NewMissingCriticalOriginHintRule() *MissingCriticalOriginHintRule {
	return &MissingCriticalOriginHintRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *MissingCriticalOriginHintRule) ID() string {
	return "lcp.missing-critical-origin-hint"
}

// Description mengembalikan ringkasan aturan.
func (r *MissingCriticalOriginHintRule) Description() string {
	return "Critical LCP visual asset loaded from third-party CDN origin without '<link rel=\"preconnect\">' connection hint in '<head>'"
}

// Category mengembalikan nama kategori rule.
func (r *MissingCriticalOriginHintRule) Category() string {
	return "lcp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (info advisory).
func (r *MissingCriticalOriginHintRule) DefaultSeverity() ir.Severity {
	return ir.SeverityInfo
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *MissingCriticalOriginHintRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Google Chrome Core Web Vitals (Largest Contentful Paint Resource Load Delay)",
			"W3C Resource Hints (preconnect and dns-prefetch specification)",
			"Web Performance Working Group Network Socket Pre-warming Guidelines",
		},
		CoreInvariant: "Critical visual LCP assets hosted on external cross-origin domains should have an early '<link rel=\"preconnect\">' hint in '<head>' to eliminate DNS, TCP, and TLS socket handshake round-trips.",
		Grounding: "When an above-the-fold hero image is loaded from a third-party domain (e.g. 'images.unsplash.com', 'cdn.shopify.com', or 'res.cloudinary.com'), the browser cannot initiate the network connection until the `<img>` tag is discovered and parsed.\n\n" +
			"Establishing a new HTTPS connection to an external origin requires three sequential round-trips: DNS resolution, TCP three-way handshake, and TLS cryptographic negotiation, adding 150ms to 400ms of idle latency on cellular connections.\n\n" +
			"Declaring `<link rel=\"preconnect\" href=\"https://images.unsplash.com\" />` in `<head>` instructs the browser to open the socket in the background during early HTML document streaming.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Connection Handshake Latency Spike",
				Severity: "MEDIUM",
				Impact:   "Adds 150ms-400ms of socket negotiation delay to LCP Resource Load Delay before image bytes start streaming.",
			},
			{
				Vector:   "Cellular Round-Trip Time Penalty",
				Severity: "LOW",
				Impact:   "Multi-RTT connection setup noticeably degrades mobile performance metrics in emerging market networks.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "html",
				Comment:  "Critical hero image loaded from external CDN without preconnect hint in head",
				Code: `<head>
  <title>Product Showcase</title>
</head>
<body>
  <img src="https://images.unsplash.com/photo-hero" fetchpriority="high" data-perf-role="hero" />
</body>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "html",
				Comment:  "Preconnect hint declared in head to pre-warm CDN connection sockets early",
				Code: `<head>
  <title>Product Showcase</title>
  <link rel="preconnect" href="https://images.unsplash.com" />
</head>
<body>
  <img src="https://images.unsplash.com/photo-hero" fetchpriority="high" data-perf-role="hero" />
</body>`,
			},
		},
	}
}

// Evaluate memeriksa apakah aset LCP eksternal dimuat tanpa deklarasi preconnect di head.
func (r *MissingCriticalOriginHintRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || !strings.EqualFold(node.Tag, "img") {
		return nil
	}

	// Hanya audit jika elemen merupakan kandidat LCP
	isCandidate, _ := isLCPCandidate(node)
	if !isCandidate {
		return nil
	}

	srcVal, hasSrc := node.Attributes["src"]
	if !hasSrc {
		return nil
	}

	origin, host, isExt := extractExternalOrigin(srcVal)
	if !isExt {
		return nil
	}

	if hasOriginPreconnect(node, host) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Critical LCP image '%s' loaded from external origin '%s' without preconnect hints. Adding '<link rel=\"preconnect\" href=\"%s\">' in '<head>' eliminates 150ms-400ms of DNS/TLS handshake delay.", srcVal, host, origin),
			Hint:     fmt.Sprintf("Add '<link rel=\"preconnect\" href=\"%s\" />' inside the document '<head>'.", origin),
		},
	}
}
