package responsive

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// DynamicViewportInconsistencyRule mendeteksi pencampuran unit viewport statis (vh, h-screen)
// dan unit viewport dinamis (dvh, svh) pada hierarki layout yang sama, yang memicu lonjakan layout
// dan scrollbar ganda pada browser mobile.
type DynamicViewportInconsistencyRule struct{}

// NewDynamicViewportInconsistencyRule membuat instance baru dari DynamicViewportInconsistencyRule.
func NewDynamicViewportInconsistencyRule() *DynamicViewportInconsistencyRule {
	return &DynamicViewportInconsistencyRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *DynamicViewportInconsistencyRule) ID() string {
	return "responsive.dynamic-viewport-inconsistency"
}

// Description mengembalikan ringkasan aturan.
func (r *DynamicViewportInconsistencyRule) Description() string {
	return "Warns when static viewport units (100vh, h-screen) and modern dynamic units (dvh, svh) are mixed inconsistently across layout hierarchies"
}

// Category mengembalikan nama kategori rule.
func (r *DynamicViewportInconsistencyRule) Category() string {
	return "responsive"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *DynamicViewportInconsistencyRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *DynamicViewportInconsistencyRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C CSS Values and Units Module Level 4 (Small, Large, and Dynamic Viewport Units)",
			"WebKit Dynamic Viewport Sizing Specification",
			"Chrome for Android URL Bar Scroll Resize Guidelines",
		},
		CoreInvariant: "Components nested within a dynamic viewport container ('dvh', 'svh') must not use static viewport units ('100vh', 'h-screen'), and conflicting viewport dimensions must not be declared on the same element.",
		Grounding: "Modern mobile browsers (Safari iOS and Chrome Android) feature dynamic interface chrome (URL address bar and bottom navigation toolbar) that expand and collapse during user scrolling.\n\n" +
			"The dynamic viewport unit 'dvh' continuously tracks the active visible viewport height. In contrast, classical '100vh' and 'h-screen' are fixed to the Large Viewport (the maximum screen height assuming all browser chrome is collapsed).\n\n" +
			"When an outer wrapper uses 'min-h-dvh' while an inner component specifies 'h-screen' or 'h-[100vh]', the child height exceeds the visible parent area whenever the address bar is visible, causing unexpected double scrollbars, layout clipping, and jarring viewport jitter.\n\n" +
			"Charites enforces consistent viewport unit adoption across component hierarchies.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Mobile Double Scrollbar & Viewport Jitter",
				Severity: "MEDIUM",
				Impact:   "Inner components sized with 100vh exceed the dvh container, causing double scrollbars and layout jerking during scroll.",
			},
			{
				Vector:   "Content Clipping Behind Browser Chrome",
				Severity: "LOW",
				Impact:   "Bottom actions and footers are pushed offscreen beneath mobile browser toolbars.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Inner child with h-screen nested inside an outer min-h-dvh container",
				Code: `<main className="min-h-dvh flex flex-col">
  <div className="h-screen bg-surface">
    <h2>Konten Terpotong di Mobile</h2>
  </div>
</main>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Consistent dynamic viewport units across parent and child",
				Code: `<main className="min-h-dvh flex flex-col">
  <div className="h-full bg-surface">
    <h2>Konten Selaras Mengikuti Viewport</h2>
  </div>
</main>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen mencampurkan unit statis dan dinamis secara bertentangan.
func (r *DynamicViewportInconsistencyRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	hasStatic := hasClassicalViewportHeight(node.Classes)
	hasDynamic := hasDynamicViewportHeight(node.Classes)

	if hasStatic && hasDynamic {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Element <%s> mixes conflicting static ('h-screen', '100vh') and dynamic ('dvh', 'svh') viewport units on the same node.", node.Tag),
				Hint:     "Standardize on dynamic viewport units ('min-h-dvh' or 'h-dvh') to ensure smooth mobile browser address bar responsiveness.",
			},
		}
	}

	if hasStatic {
		curr := node.Parent
		for curr != nil {
			if hasDynamicViewportHeight(curr.Classes) {
				return []ir.Diagnostic{
					{
						Line:     node.Span.Line,
						Column:   node.Span.Column,
						Rule:     r.ID(),
						Severity: r.DefaultSeverity(),
						Message:  fmt.Sprintf("Element <%s> uses static viewport unit ('h-screen' or '100vh') while ancestor <%s> uses dynamic unit ('dvh'). This causes the child to exceed visible boundaries and force double scrollbars on mobile browsers.", node.Tag, curr.Tag),
						Hint:     "Use 'h-full' or dynamic viewport units ('h-dvh') on the child element to stay consistent with the parent container.",
					},
				}
			}
			curr = curr.Parent
		}
	}

	return nil
}
