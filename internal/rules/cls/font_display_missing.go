package cls

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// FontDisplayMissingRule mendeteksi deklarasi @font-face yang tidak memiliki
// deskriptor 'font-display' sah (swap, optional, atau fallback).
type FontDisplayMissingRule struct{}

// NewFontDisplayMissingRule membuat instance baru FontDisplayMissingRule.
func NewFontDisplayMissingRule() *FontDisplayMissingRule {
	return &FontDisplayMissingRule{}
}

// ID mengembalikan identifier kanonikal Semgrep rule.
func (r *FontDisplayMissingRule) ID() string {
	return "cls.font-display-missing"
}

// Description mengembalikan ringkasan aturan.
func (r *FontDisplayMissingRule) Description() string {
	return "Requires font-display descriptor on custom @font-face declarations to prevent FOIT reflow"
}

// Category mengembalikan nama kategori rule.
func (r *FontDisplayMissingRule) Category() string {
	return "cls"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *FontDisplayMissingRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *FontDisplayMissingRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C CSS Fonts Module Level 4 (@font-face font-display descriptor)",
			"Google Core Web Vitals Guidelines (Cumulative Layout Shift & FOUT/FOIT)",
			"Web.dev Font Best Practices",
		},
		CoreInvariant: "All custom @font-face declarations must declare an explicit, valid 'font-display' descriptor ('swap', 'optional', or 'fallback') to ensure continuous text visibility and prevent layout reflow.",
		Grounding: "When a browser encounters a custom @font-face without a 'font-display' descriptor, it defaults to 'font-display: auto' (often identical to 'block').\n\n" +
			"Under the 'block' period, the browser hides text completely (Flash of Invisible Text / FOIT) for up to 3 seconds while waiting for the web font file. Once the font arrives, the browser suddenly swaps the font and recalculates line wrapping, triggering Cumulative Layout Shift (CLS).\n\n" +
			"Using 'font-display: swap' renders system fallback fonts immediately and swaps gracefully when the font finishes loading, ensuring accessibility and predictable rendering.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Flash of Invisible Text (FOIT)",
				Severity: "HIGH",
				Impact:   "Users stare at blank spaces on slow networks while waiting for fonts to load.",
			},
			{
				Vector:   "Cumulative Layout Shift (CLS)",
				Severity: "HIGH",
				Impact:   "Late font swaps cause text wrapping reflow that pushes subsequent content down.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Custom @font-face rule missing font-display descriptor",
				Code: `<style>
  @font-face {
    font-family: 'GeistSans';
    src: url('/fonts/geist.woff2') format('woff2');
  }
</style>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Custom @font-face declaring font-display: swap",
				Code: `<style>
  @font-face {
    font-family: 'GeistSans';
    src: url('/fonts/geist.woff2') format('woff2');
    font-display: swap;
  }
</style>`,
			},
		},
	}
}

// Evaluate memeriksa apakah blok @font-face menyertakan deskriptor font-display yang sah.
func (r *FontDisplayMissingRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Tag != "style" {
		return nil
	}

	cssText := getStyleNodeText(node)
	if !strings.Contains(cssText, "@font-face") {
		return nil
	}

	blocks := extractFontFaceBlocks(cssText)
	if len(blocks) == 0 {
		return nil
	}

	diags := make([]ir.Diagnostic, 0, len(blocks))
	for _, block := range blocks {
		if !hasValidFontDisplay(block) {
			diags = append(diags, ir.Diagnostic{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Declaration '@font-face' lacks a valid 'font-display' descriptor ('swap', 'optional', or 'fallback'). Missing font-display causes invisible text during loading (FOIT) and subsequent Cumulative Layout Shift (CLS).",
				Hint:     "Add 'font-display: swap;' or 'font-display: optional;' to the '@font-face' declaration to guarantee text visibility.",
			})
		}
	}

	return diags
}
