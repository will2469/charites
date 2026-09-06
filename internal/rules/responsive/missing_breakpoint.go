package responsive

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// MissingBreakpointRule mendeteksi deklarasi layout grid multi-kolom (grid-cols-3 s/d 12)
// atau tipografi raksasa (text-5xl s/d 9xl) pada baseline mobile tanpa modifier breakpoint responsif.
type MissingBreakpointRule struct{}

// NewMissingBreakpointRule membuat instance baru dari MissingBreakpointRule.
func NewMissingBreakpointRule() *MissingBreakpointRule {
	return &MissingBreakpointRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *MissingBreakpointRule) ID() string {
	return "responsive.missing-breakpoint"
}

// Description mengembalikan ringkasan aturan.
func (r *MissingBreakpointRule) Description() string {
	return "Warns when multi-column grids or giant font sizes are declared on mobile baseline without responsive breakpoint modifiers"
}

// Category mengembalikan nama kategori rule.
func (r *MissingBreakpointRule) Category() string {
	return "responsive"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *MissingBreakpointRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *MissingBreakpointRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Mobile-First Responsive Web Design Principles",
			"W3C CSS Grid Layout Module Level 2",
			"Tailwind CSS Responsive Design Specification",
		},
		CoreInvariant: "Multi-column grids (grid-cols-[3-12]) and giant font sizes (text-[5-9]xl) must not be defined on mobile baseline without responsive breakpoint prefixes (sm:, md:, lg:).",
		Grounding: "On compact smartphone screens (360px-390px), defining 3 or more columns directly on mobile baseline squeezes individual columns below 100px width, causing severe card distortion and text wrapping.\n\n" +
			"Similarly, declaring giant typography (e.g. text-6xl) on mobile baseline causes single words to span multiple vertical lines, breaking header visual balance.\n\n" +
			"Adhering to mobile-first progression requires starting from single-column baselines (grid-cols-1) and scaling up via responsive modifiers (sm:grid-cols-2 md:grid-cols-4).",
		Risks: []ir.RiskItem{
			{
				Vector:   "Severe Column Squeeze on Mobile",
				Severity: "MEDIUM",
				Impact:   "Multi-column cards become unreadable and distorted when squeezed into 360px phone screens.",
			},
			{
				Vector:   "Typography Layout Blowout",
				Severity: "LOW",
				Impact:   "Giant font headings wrap awkwardly into 4-5 lines on narrow mobile viewports.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Multi-column grid on mobile baseline without responsive modifier",
				Code: `<div className="grid grid-cols-4 gap-4">
  <div className="bg-card p-4">Item 1</div>
  <div className="bg-card p-4">Item 2</div>
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Mobile-first progression starting from 1 column to multi-column on desktop",
				Code: `<div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4">
  <div className="bg-card p-4">Item 1</div>
  <div className="bg-card p-4">Item 2</div>
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen mendeklarasikan grid multi-kolom atau font raksasa tanpa breakpoint.
func (r *MissingBreakpointRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || len(node.Classes) == 0 {
		return nil
	}

	var diags []ir.Diagnostic
	for _, cls := range node.Classes {
		if hasBreakpointPrefix(cls) {
			continue
		}
		if isMultiColGrid(cls) {
			diags = append(diags, ir.Diagnostic{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Multi-column grid %q declared on mobile baseline without responsive breakpoint prefix (e.g. sm:, md:). On compact mobile viewports (360px), multi-column grids compress content below 100px per column.", cls),
				Hint:     "Start with a single column baseline 'grid-cols-1' and scale with breakpoint modifiers, e.g. 'grid-cols-1 sm:grid-cols-2 md:grid-cols-4'.",
			})
		} else if isGiantFont(cls) {
			diags = append(diags, ir.Diagnostic{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Giant font size %q declared on mobile baseline without responsive breakpoint prefix. Oversized text causes extreme line-wrapping and layout distortion on compact mobile screens.", cls),
				Hint:     "Use a smaller mobile baseline and scale on larger viewports, e.g. 'text-2xl md:text-5xl'.",
			})
		}
	}

	return diags
}
