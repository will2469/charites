package cls

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// UnadjustedFontMetricRule mendeteksi deklarasi @font-face fallback lokal (src: local(...))
// yang tidak menyertakan deskriptor penyesuaian metrik font (size-adjust, ascent-override, descent-override).
type UnadjustedFontMetricRule struct{}

// NewUnadjustedFontMetricRule membuat instance baru UnadjustedFontMetricRule.
func NewUnadjustedFontMetricRule() *UnadjustedFontMetricRule {
	return &UnadjustedFontMetricRule{}
}

// ID mengembalikan identifier kanonikal Semgrep rule.
func (r *UnadjustedFontMetricRule) ID() string {
	return "cls.unadjusted-font-metric"
}

// Description mengembalikan ringkasan aturan.
func (r *UnadjustedFontMetricRule) Description() string {
	return "Recommends font metric overrides on fallback font declarations to reduce swap CLS"
}

// Category mengembalikan nama kategori rule.
func (r *UnadjustedFontMetricRule) Category() string {
	return "cls"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (info).
func (r *UnadjustedFontMetricRule) DefaultSeverity() ir.Severity {
	return ir.SeverityInfo
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *UnadjustedFontMetricRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C CSS Fonts Module Level 4 (size-adjust, ascent-override, descent-override)",
			"Google Chrome Font Metric Override Guidelines",
		},
		CoreInvariant: "Local fallback @font-face declarations (using 'src: local(...)') should specify metric adjustment descriptors ('size-adjust', 'ascent-override', or 'descent-override') to align bounding boxes with the principal web font.",
		Grounding: "When a web font downloads and replaces a system fallback font, variations in glyph x-height, ascent, and descent alter the computed bounding boxes of every text line.\n\n" +
			"This disparity causes sudden vertical expansion or contraction of paragraphs and headers, contributing to Cumulative Layout Shift.\n\n" +
			"By declaring 'size-adjust', 'ascent-override', and 'descent-override' on the fallback @font-face, developers can calibrate the system font's metrics to match the web font, creating a near zero-shift swap.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Font Swap Layout Jitter",
				Severity: "LOW",
				Impact:   "Paragraphs and navigation bars visibly shift lines when the web font swaps in.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Local fallback font-face without metric override descriptors",
				Code: `<style>
  @font-face {
    font-family: 'InterFallback';
    src: local('Arial');
  }
</style>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Local fallback font-face with size-adjust and ascent-override",
				Code: `<style>
  @font-face {
    font-family: 'InterFallback';
    src: local('Arial');
    ascent-override: 90%;
    descent-override: 22%;
    size-adjust: 107%;
  }
</style>`,
			},
		},
	}
}

// Evaluate memeriksa apakah deklarasi fallback font menyertakan penyesuaian metrik.
func (r *UnadjustedFontMetricRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Tag != "style" {
		return nil
	}

	cssText := getStyleNodeText(node)
	if !strings.Contains(cssText, "@font-face") || !strings.Contains(cssText, "local(") {
		return nil
	}

	blocks := extractFontFaceBlocks(cssText)
	if len(blocks) == 0 {
		return nil
	}

	diags := make([]ir.Diagnostic, 0, len(blocks))
	for _, block := range blocks {
		if hasLocalFontSource(block) && !hasFontMetricOverrides(block) {
			diags = append(diags, ir.Diagnostic{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Local fallback '@font-face' ('src: local(...)') lacks font metric overrides ('size-adjust', 'ascent-override', or 'descent-override'). Disparities in glyph bounding boxes between system and web fonts will cause text reflow on swap.",
				Hint:     "Declare 'size-adjust', 'ascent-override', or 'descent-override' on the fallback font-face to calibrate bounding boxes with the principal web font.",
			})
		}
	}

	return diags
}
