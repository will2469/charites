package responsive

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// ViewportUnitLeakRule mendeteksi penggunaan unit tinggi viewport statis (h-screen, min-h-screen, 100vh)
// yang memicu pergeseran tata letak (layout shift) saat bilah navigasi browser mobile muncul atau menghilang.
type ViewportUnitLeakRule struct{}

// NewViewportUnitLeakRule membuat instance baru dari ViewportUnitLeakRule.
func NewViewportUnitLeakRule() *ViewportUnitLeakRule {
	return &ViewportUnitLeakRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *ViewportUnitLeakRule) ID() string {
	return "responsive.viewport-unit-leak"
}

// Description mengembalikan ringkasan aturan.
func (r *ViewportUnitLeakRule) Description() string {
	return "Warns when viewport height relies on static 100vh instead of modern dynamic dvh or svh units"
}

// Category mengembalikan nama kategori rule.
func (r *ViewportUnitLeakRule) Category() string {
	return "responsive"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *ViewportUnitLeakRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ViewportUnitLeakRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C CSS Values and Units Module Level 4 (Small, Large, and Dynamic Viewport Units)",
			"WebKit Safari iOS Dynamic Viewport Sizing Specification",
			"Core Web Vitals Cumulative Layout Shift (CLS) Mitigation",
		},
		CoreInvariant: "Viewport height declarations should use CSS Level 4 dynamic units (dvh, svh) rather than static 100vh (h-screen, min-h-screen) to eliminate mobile layout shifts.",
		Grounding: "On mobile browsers (Safari iOS and Chrome Android), the browser address bar and bottom toolbar dynamically expand and collapse during vertical scrolling.\n\n" +
			"The classic 100vh unit uses the Large Viewport Height, which does not account for the visible URL bar. This causes bottom-anchored content to be covered by browser chrome and leads to disruptive layout jumps when the address bar toggles.\n\n" +
			"Utilizing dynamic viewport units (min-h-dvh or h-dvh) ensures the layout adapts smoothly to the actual visible viewport height.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Mobile Browser Layout Jumps (CLS)",
				Severity: "MEDIUM",
				Impact:   "Content suddenly shifts when mobile address bar hides or appears during scroll.",
			},
			{
				Vector:   "Occluded Bottom CTA Buttons",
				Severity: "LOW",
				Impact:   "Bottom buttons in a 100vh container are partially covered beneath Safari's bottom navigation bar.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Static 100vh height causing layout jumps on mobile browsers",
				Code: `<main className="min-h-screen bg-background flex flex-col justify-between">
  <h1>Beranda Desa</h1>
  <button className="h-11 px-4 bg-primary text-primary-foreground">Lanjutkan</button>
</main>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Dynamic viewport height unit adapting smoothly to mobile address bar",
				Code: `<main className="min-h-dvh bg-background flex flex-col justify-between">
  <h1>Beranda Desa</h1>
  <button className="h-11 px-4 bg-primary text-primary-foreground">Lanjutkan</button>
</main>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen menggunakan unit viewport statis h-screen atau 100vh.
func (r *ViewportUnitLeakRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || len(node.Classes) == 0 {
		return nil
	}

	var diags []ir.Diagnostic
	for _, cls := range node.Classes {
		if isViewportHeightClass(cls) {
			diags = append(diags, ir.Diagnostic{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Viewport height utility %q relies on static 100vh calculation, causing layout jumps and cut-off content on mobile browsers as dynamic toolbars expand/collapse.", cls),
				Hint:     "Modernize to CSS Values Level 4 dynamic viewport units, e.g. 'min-h-dvh' or 'h-dvh' (or 'min-h-svh' for small viewport fallback).",
			})
		}
	}

	return diags
}
