package theme

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// MetaThemeColorMismatchRule mendeteksi tag <meta name="theme-color"> statis yang tidak
// menyertakan atribut media="(prefers-color-scheme: ...)" untuk adaptasi mobile address bar.
type MetaThemeColorMismatchRule struct{}

// NewMetaThemeColorMismatchRule membuat instance baru MetaThemeColorMismatchRule.
func NewMetaThemeColorMismatchRule() *MetaThemeColorMismatchRule {
	return &MetaThemeColorMismatchRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *MetaThemeColorMismatchRule) ID() string {
	return "theme.meta-theme-color-mismatch"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *MetaThemeColorMismatchRule) Description() string {
	return "Detects static meta theme-color tags lacking media prefers-color-scheme queries"
}

// Category mengembalikan nama kategori rule.
func (r *MetaThemeColorMismatchRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *MetaThemeColorMismatchRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki otomatis.
func (r *MetaThemeColorMismatchRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"HTML Living Standard Section 4.2.5 (The meta element)",
			"Web App Manifest & Mobile OS Theme Integration",
			"WCAG 2.2 Success Criterion 1.4.11 (Non-text Contrast)",
		},
		CoreInvariant: "Meta theme-color elements must provide media query pairs (prefers-color-scheme: light/dark) to synchronize mobile browser chrome.",
		Grounding: "Modern mobile browsers (Safari on iOS, Chrome on Android) color the operating system status bar and address bar based on the <meta name=\"theme-color\"> element in the document <head>.\n\n" +
			"When developers specify a single static theme-color without media queries (e.g. <meta name=\"theme-color\" content=\"#ffffff\">):\n" +
			"1. Blinding Address Bar: When the user toggles dark mode, the mobile address bar and status bar remain stark white, causing harsh visual glare.\n" +
			"2. Inverted Chrome: When the page switches to dark background, white status bar text collapses against white browser chrome.\n\n" +
			"Charites enforces declaring media=\"(prefers-color-scheme: light)\" and media=\"(prefers-color-scheme: dark)\" pairs on all meta theme-color definitions.",
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Static theme-color meta tag in Astro layout",
				Code:     `<meta name="theme-color" content="#ffffff" />`,
			},
			{
				Language: "tsx",
				Comment:  "Static theme-color in TSX Document head",
				Code:     `<meta name="theme-color" content="#09090b" />`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Adaptive light/dark meta theme-color pair in Astro",
				Code: `<>
  <meta name="theme-color" media="(prefers-color-scheme: light)" content="#ffffff" />
  <meta name="theme-color" media="(prefers-color-scheme: dark)" content="#09090b" />
</>`,
			},
			{
				Language: "tsx",
				Comment:  "Adaptive meta theme-color pair in TSX",
				Code: `<head>
  <meta name="theme-color" media="(prefers-color-scheme: light)" content="#ffffff" />
  <meta name="theme-color" media="(prefers-color-scheme: dark)" content="#09090b" />
</head>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Mobile Chrome Glare",
				Severity: "MEDIUM",
				Impact:   "Mobile Safari and Chrome address bars blast high-brightness white chrome when the application is viewed in dark mode.",
			},
			{
				Vector:   "Status Bar Text Invisibility",
				Severity: "LOW",
				Impact:   "OS status bar text (time, battery, Wi-Fi) becomes invisible due to poor contrast against unadapted address bar backgrounds.",
			},
		},
	}
}

// Evaluate mengevaluasi node IR untuk memeriksa apakah tag <meta name="theme-color"> memiliki query media.
func (r *MetaThemeColorMismatchRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || strings.ToLower(node.Tag) != "meta" || len(node.Attributes) == 0 {
		return nil
	}

	nameVal := strings.Trim(strings.TrimSpace(node.Attributes["name"]), "\"'`")
	if strings.ToLower(nameVal) != "theme-color" {
		return nil
	}

	mediaVal, hasMedia := node.Attributes["media"]
	cleanMedia := strings.Trim(strings.TrimSpace(mediaVal), "\"'`")
	if !hasMedia || cleanMedia == "" || !strings.Contains(strings.ToLower(cleanMedia), "prefers-color-scheme") {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Static <meta name=\"theme-color\"> tag lacks \"media=(prefers-color-scheme: ...)\" attribute",
				Hint:     "Add media=\"(prefers-color-scheme: light)\" and media=\"(prefers-color-scheme: dark)\" pairs so mobile address bars adapt to dark mode.",
			},
		}
	}

	return nil
}
