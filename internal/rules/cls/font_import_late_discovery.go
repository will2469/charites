package cls

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// FontImportLateDiscoveryRule mendeteksi penggunaan direktif CSS @import yang memuat font eksternal
// alih-alih menggunakan tag <link rel="preconnect"> dan <link rel="stylesheet"> pada <head>.
type FontImportLateDiscoveryRule struct{}

// NewFontImportLateDiscoveryRule membuat instance baru FontImportLateDiscoveryRule.
func NewFontImportLateDiscoveryRule() *FontImportLateDiscoveryRule {
	return &FontImportLateDiscoveryRule{}
}

// ID mengembalikan identifier kanonikal Semgrep rule.
func (r *FontImportLateDiscoveryRule) ID() string {
	return "cls.font-import-late-discovery"
}

// Description mengembalikan ringkasan aturan.
func (r *FontImportLateDiscoveryRule) Description() string {
	return "Warns when CSS @import is used for external font loading, delaying discovery and risking layout shift"
}

// Category mengembalikan nama kategori rule.
func (r *FontImportLateDiscoveryRule) Category() string {
	return "cls"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *FontImportLateDiscoveryRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *FontImportLateDiscoveryRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Cumulative Layout Shift (CLS) Metric Specification",
			"Google Core Web Vitals Guidelines (Render-Blocking Resources & Critical Path)",
			"Tailwind CSS v4 Import Specifications",
		},
		CoreInvariant: "External web fonts must be discovered and preconnected in HTML/Astro <head> rather than imported via cascading CSS @import directives, while whitelisting Tailwind CSS and local stylesheets.",
		Grounding: "Placing @import rules referencing external fonts (such as Google Fonts or Typekit) inside CSS creates a cascading waterfall of render-blocking requests.\n\n" +
			"The browser must download the HTML, parse the stylesheet, discover the nested @import, download the font CSS, and only then start downloading the binary font file. This profound delay forces long periods of fallback rendering and dramatic late layout shifts.\n\n" +
			"Declaring '<link rel=\"preconnect\">' alongside '<link rel=\"stylesheet\">' in the HTML layout starts DNS preconnection and font loading at the earliest possible phase of page load.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Cascading Network Waterfall",
				Severity: "HIGH",
				Impact:   "Nested font CSS imports delay font delivery by several hundred milliseconds on mobile networks.",
			},
			{
				Vector:   "Severe Late Layout Shift",
				Severity: "MEDIUM",
				Impact:   "Delayed font swapping abruptly reorganizes the text geometry long after initial paint.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Late discovery font import inside style block",
				Code: `<style>
  @import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;700&display=swap');
</style>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Fonts loaded via preconnect and stylesheet link in head",
				Code: `<head>
  <link rel="preconnect" href="https://fonts.googleapis.com" />
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
  <link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Inter:wght@400;700&display=swap" />
</head>`,
			},
			{
				Language: "astro",
				Comment:  "Whitelisted Tailwind CSS and local file imports",
				Code: `<style>
  @import "tailwindcss";
  @import "./local-tokens.css";
</style>`,
			},
		},
	}
}

// Evaluate memeriksa apakah terdapat deklarasi @import yang memuat font eksternal.
func (r *FontImportLateDiscoveryRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Tag != "style" {
		return nil
	}

	cssText := getStyleNodeText(node)
	if !strings.Contains(cssText, "@import") {
		return nil
	}

	violations := findExternalFontImports(cssText)
	if len(violations) == 0 {
		return nil
	}

	diags := make([]ir.Diagnostic, 0, len(violations))
	for _, v := range violations {
		diags = append(diags, ir.Diagnostic{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "CSS '@import' references external font stylesheet (" + v + "). Importing fonts inside CSS creates a cascading render-blocking dependency chain and delays font discovery, risking late layout shift.",
			Hint:     "Replace the CSS '@import' with '<link rel=\"preconnect\">' and '<link rel=\"stylesheet\">' tags directly in your HTML/Astro '<head>' layout.",
		})
	}

	return diags
}
