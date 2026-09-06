package cls

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// UnstableScrollbarGutterRule mendeteksi deklarasi overflow-y: auto pada root dokumen (html, body, :root)
// yang tidak menyertakan scrollbar-gutter: stable, yang menyebabkan pergeseran horizontal (15-17px) saat scrollbar muncul.
type UnstableScrollbarGutterRule struct{}

// NewUnstableScrollbarGutterRule membuat instance baru dari UnstableScrollbarGutterRule.
func NewUnstableScrollbarGutterRule() *UnstableScrollbarGutterRule {
	return &UnstableScrollbarGutterRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *UnstableScrollbarGutterRule) ID() string {
	return "cls.unstable-scrollbar-gutter"
}

// Description mengembalikan ringkasan aturan.
func (r *UnstableScrollbarGutterRule) Description() string {
	return "Root document scroller declares overflow-y: auto without scrollbar-gutter: stable, risking horizontal layout shifts"
}

// Category mengembalikan nama kategori rule.
func (r *UnstableScrollbarGutterRule) Category() string {
	return "cls"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (info).
func (r *UnstableScrollbarGutterRule) DefaultSeverity() ir.Severity {
	return ir.SeverityInfo
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *UnstableScrollbarGutterRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C CSS Box Model Module Level 4 (scrollbar-gutter)",
			"Google Core Web Vitals (Horizontal Layout Shift Prevention)",
			"Desktop & Mobile Multi-Platform Viewport Consistency",
		},
		CoreInvariant: "Root document scrollers (html, body, :root) with dynamic overflow should specify 'scrollbar-gutter: stable' to permanently reserve viewport space for scrollbars.",
		Grounding: "When a webpage renders with 'overflow-y: auto' at the document root level, the operating system initially renders content spanning the full viewport width.\n\n" +
			"As dynamic content streams in, hydration completes, or the user navigates to a longer page, the vertical scrollbar suddenly appears. On non-overlay desktop platforms (Windows, Linux, non-overlay macOS), this vertical scrollbar consumes 15-17px of width.\n\n" +
			"This sudden shrinkage of available client width causes all centered layouts, responsive grids, and full-bleed headers to snap and shift horizontally, registering an instant Cumulative Layout Shift.\n\n" +
			"Adding 'scrollbar-gutter: stable;' reserves the scrollbar space permanently, ensuring viewport dimensions remain completely invariant regardless of page height.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Document-Wide Horizontal Snapping",
				Severity: "LOW",
				Impact:   "Sudden appearance of vertical scrollbars causes all centered containers and flex items to jump 15-17px horizontally.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "css",
				Comment:  "Root html scroller with auto overflow but no reserved scrollbar gutter",
				Code: `html {
  overflow-y: auto;
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "css",
				Comment:  "Root html scroller with stable scrollbar gutter reservation",
				Code: `html {
  overflow-y: auto;
  scrollbar-gutter: stable;
}`,
			},
		},
	}
}

// Evaluate memeriksa apakah blok <style> mendefinisikan root scroller tanpa scrollbar-gutter: stable.
func (r *UnstableScrollbarGutterRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || !strings.EqualFold(node.Tag, "style") {
		return nil
	}

	css := getStyleNodeText(node)
	violations := checkUnstableScrollbarGutter(css)
	if len(violations) == 0 {
		return nil
	}

	diags := make([]ir.Diagnostic, 0, len(violations))
	for _, sel := range violations {
		diags = append(diags, ir.Diagnostic{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Root document scroller '%s' declares 'overflow-y: auto' without 'scrollbar-gutter: stable'. When dynamic content expands past the viewport height, sudden appearance of the vertical scrollbar causes an unexpected 15-17px horizontal layout shift.", sel),
			Hint:     "Add 'scrollbar-gutter: stable;' or 'overflow-y: scroll;' to the root stylesheet to reserve space for the scrollbar permanently.",
		})
	}

	return diags
}
