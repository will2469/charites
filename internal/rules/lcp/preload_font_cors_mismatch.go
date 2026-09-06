package lcp

import (
	"github.com/will2469/charites/internal/ir"
)

// PreloadFontCORSMismatchRule mendeteksi elemen <link rel="preload" as="font"> yang tidak menyertakan
// atribut crossorigin, pola yang melanggar spesifikasi W3C dan menyebabkan browser membuang hasil preload.
type PreloadFontCORSMismatchRule struct{}

// NewPreloadFontCORSMismatchRule membuat instance baru dari PreloadFontCORSMismatchRule.
func NewPreloadFontCORSMismatchRule() *PreloadFontCORSMismatchRule {
	return &PreloadFontCORSMismatchRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *PreloadFontCORSMismatchRule) ID() string {
	return "lcp.preload-font-cors-mismatch"
}

// Description mengembalikan ringkasan aturan.
func (r *PreloadFontCORSMismatchRule) Description() string {
	return "Font preload '<link rel=\"preload\" as=\"font\">' lacks 'crossorigin' attribute, triggering browser cache discard and double network downloads"
}

// Category mengembalikan nama kategori rule.
func (r *PreloadFontCORSMismatchRule) Category() string {
	return "lcp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *PreloadFontCORSMismatchRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *PreloadFontCORSMismatchRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Preload Specification (Font Preload CORS Requirements)",
			"HTML Living Standard Crossorigin Attribute Specification",
			"Google Chrome Core Web Vitals (Largest Contentful Paint Resource Optimization)",
		},
		CoreInvariant: "All '<link rel=\"preload\" as=\"font\">' tags must specify the 'crossorigin' attribute to ensure the preloaded font binary is accepted by the browser font cache.",
		Grounding: "The W3C CSS Fonts specification mandates that web fonts must be fetched using anonymous Cross-Origin Resource Sharing (CORS) mode, even when the font is hosted on the same origin as the page.\n\n" +
			"When a '<link rel=\"preload\" as=\"font\">' tag omits the 'crossorigin' attribute, the preload scanner fetches the font using standard non-CORS mode.\n\n" +
			"When the CSS parser subsequently requests the font in anonymous CORS mode, the browser detects a CORS mode mismatch, discards the preloaded resource from cache, and executes a second, redundant network download.\n\n" +
			"Adding 'crossorigin' (or 'crossorigin=\"anonymous\"') ensures the preloaded font matches the CSS font engine request key, avoiding double downloads.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Redundant Network Double Download",
				Severity: "CRITICAL",
				Impact:   "The font asset is downloaded twice over the network, completely defeating the purpose of preloading and inflating bandwidth usage.",
			},
			{
				Vector:   "LCP Text Rendering Delay",
				Severity: "HIGH",
				Impact:   "The second fetch is queued after CSS parsing, adding hundreds of milliseconds to text block paint times.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "html",
				Comment:  "Font preloaded without crossorigin attribute triggering cache discard",
				Code: `<head>
  <link rel="preload" href="/fonts/inter.woff2" as="font" type="font/woff2" />
</head>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "html",
				Comment:  "Font preload declared with crossorigin attribute matching W3C anonymous CORS requirement",
				Code: `<head>
  <link rel="preload" href="/fonts/inter.woff2" as="font" type="font/woff2" crossorigin />
</head>`,
			},
		},
	}
}

// Evaluate memeriksa apakah tag preload font tidak menyertakan atribut crossorigin.
func (r *PreloadFontCORSMismatchRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if isFontPreloadMissingCORS(node) {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Font preload '<link rel=\"preload\" as=\"font\">' is missing 'crossorigin' attribute. Per W3C spec, fonts are fetched via anonymous CORS; without 'crossorigin', the preloaded cache is discarded, causing double downloads and inflating LCP.",
				Hint:     "Add 'crossorigin' or 'crossorigin=\"anonymous\"' to the '<link rel=\"preload\" as=\"font\">' tag.",
			},
		}
	}

	return nil
}
