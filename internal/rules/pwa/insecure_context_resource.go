package pwa

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// InsecureContextResourceRule memeriksa apakah aset dimuat menggunakan protokol tidak aman (http://),
// yang memicu insiden Mixed Content dan diblokir oleh peramban dalam W3C Secure Contexts.
type InsecureContextResourceRule struct{}

// NewInsecureContextResourceRule membuat instance baru dari InsecureContextResourceRule.
func NewInsecureContextResourceRule() *InsecureContextResourceRule {
	return &InsecureContextResourceRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *InsecureContextResourceRule) ID() string {
	return "pwa.insecure-context-resource"
}

// Description mengembalikan ringkasan aturan.
func (r *InsecureContextResourceRule) Description() string {
	return "Errors when a resource element loads assets over an insecure HTTP protocol (http://) in violation of W3C Secure Contexts"
}

// Category mengembalikan nama kategori rule.
func (r *InsecureContextResourceRule) Category() string {
	return "pwa"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *InsecureContextResourceRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *InsecureContextResourceRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Secure Contexts Specification",
			"W3C Mixed Content Level 2 Specification",
			"RFC 7258 Pervasive Monitoring Is an Attack",
		},
		CoreInvariant: "Resource elements (<script>, <link>, <img>, <iframe>, <video>, <audio>) must not load external assets over insecure 'http://' (except localhost loopback).",
		Grounding: "Progressive Web Apps strictly require a Secure Context (HTTPS) to enable service workers, cache storage, and device hardware APIs.\n\n" +
			"Loading external assets (scripts, stylesheets, images, media, iframes) via an unencrypted 'http://' connection triggers Mixed Content blocking. Active mixed content (scripts, stylesheets) is blocked immediately by modern mobile browsers, breaking application functionality. Passive mixed content (images, audio) generates security warnings and can be intercepted or tampered with on public Wi-Fi networks.\n\n" +
			"All asset references must use HTTPS ('https://'), protocol-relative URLs ('//'), or local origin paths. Localhost addresses ('http://localhost' and 'http://127.0.0.1') are excepted for development purposes.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Active Mixed Content Blocking",
				Severity: "HIGH",
				Impact:   "Browsers completely block insecure scripts and stylesheets, breaking UI styling and interactive application logic.",
			},
			{
				Vector:   "Man-in-the-Middle Asset Tampering",
				Severity: "HIGH",
				Impact:   "Unencrypted HTTP traffic can be intercepted, inspected, or modified by malicious actors on untrusted networks.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Insecure HTTP external script and stylesheet links",
				Code: `<div>
  <script src="http://cdn.example.org/tracker.js" />
  <link rel="stylesheet" href="http://assets.desa.id/styles.css" />
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Secure HTTPS asset loading conforming to Secure Contexts",
				Code: `<div>
  <script src="https://cdn.example.org/tracker.js" />
  <link rel="stylesheet" href="https://assets.desa.id/styles.css" />
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen memuat resource menggunakan URL http:// yang tidak aman.
func (r *InsecureContextResourceRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || !isResourceElement(node.Tag) {
		return nil
	}

	insecureURL, attrName := findInsecureURL(node)
	if insecureURL == "" {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Resource loaded over insecure HTTP protocol %q via attribute %q in a PWA. Mixed content will be blocked by browsers in Secure Contexts.", insecureURL, attrName),
			Hint:     "Use 'https://' or protocol-relative '//' URLs to ensure secure asset loading.",
		},
	}
}

func findInsecureURL(node *ir.Node) (string, string) {
	attrs := [...]string{"src", "href", "data"}
	for _, attr := range attrs {
		if val, ok := node.GetAttr(attr); ok && isInsecureResourceURL(val) {
			return cleanAttrValue(val), attr
		}
	}
	return "", ""
}
