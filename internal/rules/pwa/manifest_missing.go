package pwa

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// ManifestMissingRule memeriksa keberadaan tag <link rel="manifest" href="...">
// di dalam elemen <head> atau root layout dokumen web.
type ManifestMissingRule struct{}

// NewManifestMissingRule membuat instance baru dari ManifestMissingRule.
func NewManifestMissingRule() *ManifestMissingRule {
	return &ManifestMissingRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *ManifestMissingRule) ID() string {
	return "pwa.manifest-missing"
}

// Description mengembalikan ringkasan aturan.
func (r *ManifestMissingRule) Description() string {
	return "Warns when the HTML document <head> is missing a <link rel=\"manifest\" href=\"...\"> declaration"
}

// Category mengembalikan nama kategori rule.
func (r *ManifestMissingRule) Category() string {
	return "pwa"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *ManifestMissingRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ManifestMissingRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Web App Manifest Section 4 (Linking to a Manifest)",
			"HTML Living Standard Section 4.2.4 (The link element)",
			"Chromium Progressive Web App Discovery Engine",
		},
		CoreInvariant: "The HTML document <head> or root layout must include a '<link rel=\"manifest\" href=\"...\">' tag with a non-empty href attribute.",
		Grounding: "For mobile and desktop browsers to locate and parse a web application's manifest file, the root HTML document must link to it via a <link rel=\"manifest\" href=\"...\"> tag within the <head> section.\n\n" +
			"Without this explicit link element, browsers cannot discover the manifest, and consequently will never offer the install banner ('Add to Home Screen') or configure standalone display mode.\n\n" +
			"Including a valid manifest link in the document head ensures seamless PWA discovery across Chromium, Safari, and Gecko engines.",
		Risks: []ir.RiskItem{
			{
				Vector:   "PWA Feature Invisibility",
				Severity: "HIGH",
				Impact:   "Browsers treat the site as a traditional desktop webpage and never offer PWA installation or offline capabilities.",
			},
			{
				Vector:   "Missing Homescreen Install Capability",
				Severity: "MEDIUM",
				Impact:   "Users on mobile devices cannot install the web app to their home screen or application launcher.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "HTML head without a manifest link element",
				Code: `<head>
  <title>Layanan Surat Desa</title>
  <meta name="viewport" content="width=device-width, initial-scale=1" />
</head>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "HTML head declaring a manifest link with valid href",
				Code: `<head>
  <title>Layanan Surat Desa</title>
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <link rel="manifest" href="/manifest.webmanifest" />
</head>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen <head> memiliki link rel="manifest" yang sah.
func (r *ManifestMissingRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !isHeadElement(strings.ToLower(node.Tag)) {
		return nil
	}

	if hasManifestLink(node) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "HTML <head> is missing a '<link rel=\"manifest\" href=\"...\">' declaration. Mobile browsers cannot detect this site as an installable PWA.",
			Hint:     "Add '<link rel=\"manifest\" href=\"/manifest.webmanifest\">' (or .json) inside the <head> element.",
		},
	}
}
