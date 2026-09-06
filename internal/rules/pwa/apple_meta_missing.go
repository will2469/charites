package pwa

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// AppleMetaMissingRule memeriksa apakah dokumen HTML yang mendeklarasikan link manifest
// juga menyertakan meta tag WebKit Apple untuk mendukung mode standalone dan icon di iOS Safari.
type AppleMetaMissingRule struct{}

// NewAppleMetaMissingRule membuat instance baru dari AppleMetaMissingRule.
func NewAppleMetaMissingRule() *AppleMetaMissingRule {
	return &AppleMetaMissingRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *AppleMetaMissingRule) ID() string {
	return "pwa.apple-meta-missing"
}

// Description mengembalikan ringkasan aturan.
func (r *AppleMetaMissingRule) Description() string {
	return "Warns when an HTML document head with a Web App Manifest is missing Apple WebKit standalone meta tags (apple-mobile-web-app-capable and apple-touch-icon)"
}

// Category mengembalikan nama kategori rule.
func (r *AppleMetaMissingRule) Category() string {
	return "pwa"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *AppleMetaMissingRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *AppleMetaMissingRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Apple Safari Web Content Guide (Configuring Web Applications)",
			"WebKit Standalone PWA Architecture",
			"W3C Web App Manifest (Apple Ecosystem Compatibility)",
		},
		CoreInvariant: "When an HTML document <head> links to a Web App Manifest, it must declare '<meta name=\"apple-mobile-web-app-capable\" content=\"yes\">' and '<link rel=\"apple-touch-icon\" href=\"...\">'.",
		Grounding: "On Apple iOS (iPhone and iPad), Mobile Safari historically ignores the W3C Web App Manifest 'display: standalone' and 'icons' array when a user taps 'Add to Home Screen'.\n\n" +
			"To ensure the web app launches in an immersive fullscreen standalone mode without browser chrome (URL bar and bottom toolbar) and displays a sharp, high-resolution app icon on the iOS springboard, developers must declare Apple WebKit meta tags.\n\n" +
			"Providing both 'apple-mobile-web-app-capable' and 'apple-touch-icon' guarantees native-feeling PWA experiences on Apple devices.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Browser Chrome Intrusion on iOS",
				Severity: "MEDIUM",
				Impact:   "PWA launched from iOS Home Screen opens inside a regular Safari browser tab with URL navigation bars.",
			},
			{
				Vector:   "Degraded Springboard Branding",
				Severity: "LOW",
				Impact:   "iOS displays a shrunken screenshot placeholder instead of the official high-resolution application icon.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Head with manifest link but missing Apple WebKit meta tags",
				Code: `<head>
  <title>Layanan Desa</title>
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <link rel="manifest" href="/manifest.webmanifest" />
</head>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Head declaring both WebKit standalone meta and apple-touch-icon",
				Code: `<head>
  <title>Layanan Desa</title>
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <link rel="manifest" href="/manifest.webmanifest" />
  <meta name="apple-mobile-web-app-capable" content="yes" />
  <link rel="apple-touch-icon" href="/apple-touch-icon.png" />
</head>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen <head> yang memiliki manifest menyertakan meta tag WebKit Apple.
func (r *AppleMetaMissingRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || !isHeadElement(node.Tag) {
		return nil
	}

	if !hasManifestLink(node) {
		return nil
	}

	hasCapable := hasAppleCapableMeta(node)
	hasIcon := hasAppleTouchIcon(node)
	if hasCapable && hasIcon {
		return nil
	}

	missing := collectMissingAppleTags(hasCapable, hasIcon)
	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("PWA document <head> declares a manifest link but is missing Apple WebKit standalone tags: %s. iOS Safari will launch the web app inside a browser tab with a generic screenshot icon.", strings.Join(missing, " and ")),
			Hint:     "Add '<meta name=\"apple-mobile-web-app-capable\" content=\"yes\">' and '<link rel=\"apple-touch-icon\" href=\"/apple-touch-icon.png\">' inside the <head> element.",
		},
	}
}

func collectMissingAppleTags(hasCapable, hasIcon bool) []string {
	var missing []string
	if !hasCapable {
		missing = append(missing, "<meta name=\"apple-mobile-web-app-capable\" content=\"yes\">")
	}
	if !hasIcon {
		missing = append(missing, "<link rel=\"apple-touch-icon\" href=\"...\">")
	}
	return missing
}
