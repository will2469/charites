package mobile

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// FixedActionObstructionRule mendeteksi elemen fixed bottom (bottom nav, floating CTA)
// yang tidak memiliki padding bawah kompensasi pada parent atau sibling konten utama.
type FixedActionObstructionRule struct{}

// NewFixedActionObstructionRule membuat instance baru dari FixedActionObstructionRule.
func NewFixedActionObstructionRule() *FixedActionObstructionRule {
	return &FixedActionObstructionRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *FixedActionObstructionRule) ID() string {
	return "mobile.fixed-action-obstruction"
}

// Description mengembalikan ringkasan aturan.
func (r *FixedActionObstructionRule) Description() string {
	return "Warns when fixed bottom elements lack compensating bottom padding on parent or content siblings, risking content obstruction"
}

// Category mengembalikan nama kategori rule.
func (r *FixedActionObstructionRule) Category() string {
	return "mobile"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *FixedActionObstructionRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *FixedActionObstructionRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Apple Human Interface Guidelines (Bottom Toolbars & Screen Clearance)",
			"Google Material Design 3 (Bottom App Bars & Safe Content Boundaries)",
			"W3C CSS Positioned Layout Module Level 3 (Fixed Positioning)",
		},
		CoreInvariant: "Fixed bottom bars and floating action buttons must be accompanied by compensating bottom padding ('pb-16', 'pb-20', 'pb-24', 'pb-safe') on parent layouts or content siblings to prevent content obstruction.",
		Grounding: "Elements anchored with 'fixed bottom-0' float out of normal document flow, permanently covering the lower portion of the viewport.\n\n" +
			"Without compensating bottom padding (such as 'pb-24' or 'pb-[env(safe-area-inset-bottom)]') on the layout container or content siblings (<main>, <article>, <form>), the final rows of text, interactive inputs, or submit buttons will be permanently hidden behind the fixed bar.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Obstructed Content & Trapped Form Inputs",
				Severity: "MEDIUM",
				Impact:   "The bottom-most form fields or submit controls are occluded by the fixed bar, blocking user progress.",
			},
			{
				Vector:   "Accidental Clicks on Fixed Bar",
				Severity: "LOW",
				Impact:   "Users attempting to tap the bottom of the page accidentally trigger bottom navigation items instead.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Fixed bottom nav without compensating bottom padding on main content",
				Code: `<div className="min-h-screen bg-background">
  <main className="p-4 space-y-4">
    <p>Konten formulir paling bawah...</p>
  </main>
  <nav className="fixed bottom-0 inset-x-0 h-16 bg-card border-t">
    <button type="button">Beranda</button>
  </nav>
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Compensating bottom padding (pb-24) ensures full content clearance",
				Code: `<div className="min-h-screen bg-background">
  <main className="p-4 space-y-4 pb-24">
    <p>Konten formulir paling bawah...</p>
  </main>
  <nav className="fixed bottom-0 inset-x-0 h-16 bg-card border-t pb-[env(safe-area-inset-bottom)]">
    <button type="button">Beranda</button>
  </nav>
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen fixed bottom memiliki kompensasi padding di layout atau sibling.
func (r *FixedActionObstructionRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !isFixedBottomElement(node.Classes, node.RawClasses) {
		return nil
	}

	if isDesktopOnly(node) {
		return nil
	}

	if node.Parent == nil {
		return nil
	}

	if hasCompensatingBottomPadding(node.Parent.Classes, node.Parent.RawClasses) {
		return nil
	}

	if hasSiblingWithBottomPadding(node.Parent.Children, node) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Fixed bottom element (<" + node.Tag + ">) floats over content without compensating bottom padding on parent or content siblings, risking content obstruction on mobile screens.",
			Hint:     "Add compensating bottom padding (e.g. 'pb-20', 'pb-24', or 'pb-[env(safe-area-inset-bottom)]') to the parent container or content sibling (<main>) to provide clearance.",
		},
	}
}

func isFixedBottomElement(classes []string, rawClasses string) bool {
	if strings.Contains(rawClasses, "fixed") &&
		(strings.Contains(rawClasses, "bottom-0") || strings.Contains(rawClasses, "bottom-")) {
		return true
	}

	var hasFixed, hasBottom bool
	for _, cls := range classes {
		if cls == "fixed" {
			hasFixed = true
		}
		if cls == "bottom-0" || strings.HasPrefix(cls, "bottom-") {
			hasBottom = true
		}
	}
	return hasFixed && hasBottom
}

func hasCompensatingBottomPadding(classes []string, rawClasses string) bool {
	if strings.Contains(rawClasses, "pb-") ||
		strings.Contains(rawClasses, "pb-[") ||
		strings.Contains(rawClasses, "p-safe") ||
		strings.Contains(rawClasses, "pb-safe") {
		// Validasi bukan hanya pb-0
		if !isOnlyZeroPadding(rawClasses) {
			return true
		}
	}

	for _, cls := range classes {
		if strings.HasPrefix(cls, "pb-") && cls != "pb-0" {
			return true
		}
		if cls == "pb-safe" || cls == "p-safe" {
			return true
		}
	}

	return false
}

func isOnlyZeroPadding(raw string) bool {
	return strings.Contains(raw, "pb-0") && !strings.Contains(raw, "pb-1") &&
		!strings.Contains(raw, "pb-2") && !strings.Contains(raw, "pb-3") &&
		!strings.Contains(raw, "pb-4") && !strings.Contains(raw, "pb-6") &&
		!strings.Contains(raw, "pb-8") && !strings.Contains(raw, "pb-safe")
}

func hasSiblingWithBottomPadding(siblings []*ir.Node, self *ir.Node) bool {
	for _, sib := range siblings {
		if sib == self || sib.Type != ir.NodeElement {
			continue
		}
		if hasCompensatingBottomPadding(sib.Classes, sib.RawClasses) {
			return true
		}
	}
	return false
}
