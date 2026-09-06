package lcp

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// ExternalFontDiscoveryDelayRule mendeteksi tag <link rel="stylesheet"> yang memuat font pihak ketiga
// (seperti Google Fonts) tanpa menyertakan tag <link rel="preconnect"> ke origin biner font.
type ExternalFontDiscoveryDelayRule struct{}

// NewExternalFontDiscoveryDelayRule membuat instance baru dari ExternalFontDiscoveryDelayRule.
func NewExternalFontDiscoveryDelayRule() *ExternalFontDiscoveryDelayRule {
	return &ExternalFontDiscoveryDelayRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *ExternalFontDiscoveryDelayRule) ID() string {
	return "lcp.external-font-discovery-delay"
}

// Description mengembalikan ringkasan aturan.
func (r *ExternalFontDiscoveryDelayRule) Description() string {
	return "External font stylesheet loaded without '<link rel=\"preconnect\">' hints, adding 200ms-400ms connection handshake latency to LCP font discovery"
}

// Category mengembalikan nama kategori rule.
func (r *ExternalFontDiscoveryDelayRule) Category() string {
	return "lcp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *ExternalFontDiscoveryDelayRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ExternalFontDiscoveryDelayRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Google Chrome Core Web Vitals (Largest Contentful Paint Resource Load Delay)",
			"W3C Resource Hints (preconnect and dns-prefetch specification)",
			"Web Performance Working Group Connection Handshake Optimization Invariants",
		},
		CoreInvariant: "External cross-origin font stylesheets must be preceded by '<link rel=\"preconnect\">' hints to eliminate DNS, TCP, and TLS handshake round-trips before font binaries are requested.",
		Grounding: "Loading web fonts from third-party CDNs (such as Google Fonts or Adobe Typekit) involves a multi-origin dependency chain: the CSS stylesheet is fetched from one domain (e.g. 'fonts.googleapis.com'), while the font binaries (.woff2) are hosted on a separate storage origin (e.g. 'fonts.gstatic.com').\n\n" +
			"Without preconnect hints, the browser cannot initiate the DNS resolution, TCP three-way handshake, and TLS negotiation with the font storage origin until the stylesheet is completely downloaded and parsed.\n\n" +
			"Adding '<link rel=\"preconnect\" href=\"https://fonts.gstatic.com\" crossorigin>' allows the browser to perform socket setup in parallel during initial document streaming, shaving 200ms-400ms off LCP font discovery.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Connection Handshake Serialization",
				Severity: "HIGH",
				Impact:   "Sequential DNS/TCP/TLS handshakes to external font origins add 200ms-400ms round-trip latency to critical text LCP.",
			},
			{
				Vector:   "Cellular Network Latency Amplification",
				Severity: "MEDIUM",
				Impact:   "On mobile 3G/4G networks with high RTT (Round Trip Time), serialized connection setup severely degrades user perceptual paint times.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "html",
				Comment:  "External Google Fonts imported without preconnect hints to stylesheet and font binary origins",
				Code: `<head>
  <link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@700&display=swap" />
</head>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "html",
				Comment:  "Preconnect hints declared early to pre-warm DNS and TLS sockets before stylesheet parsing",
				Code: `<head>
  <link rel="preconnect" href="https://fonts.googleapis.com" />
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
  <link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@700&display=swap" />
</head>`,
			},
		},
	}
}

// Evaluate memeriksa apakah stylesheet font pihak ketiga dimuat tanpa preconnect ke origin biner font.
func (r *ExternalFontDiscoveryDelayRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	href, isThirdParty := isThirdPartyFontStylesheet(node)
	if !isThirdParty {
		return nil
	}

	requiredOrigin := getRequiredFontPreconnectOrigin(href)
	if hasPreconnectHint(node, requiredOrigin) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("External font stylesheet '%s' loaded without preconnect hints. Lack of '<link rel=\"preconnect\">' to font origin '%s' adds 200ms-400ms of connection handshake latency to LCP font discovery.", href, requiredOrigin),
			Hint:     "Add '<link rel=\"preconnect\" href=\"https://fonts.googleapis.com\" />' and '<link rel=\"preconnect\" href=\"https://fonts.gstatic.com\" crossorigin />' before the stylesheet link.",
		},
	}
}
