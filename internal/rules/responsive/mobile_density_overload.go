package responsive

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// MobileDensityOverloadRule mendeteksi bilah aksi atau toolbar yang memadatkan lebih dari 4 tombol interaktif
// dalam satu baris horizontal tanpa pengguliran horizontal atau pembungkus responsif di layar mobile.
type MobileDensityOverloadRule struct{}

// NewMobileDensityOverloadRule membuat instance baru dari MobileDensityOverloadRule.
func NewMobileDensityOverloadRule() *MobileDensityOverloadRule {
	return &MobileDensityOverloadRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *MobileDensityOverloadRule) ID() string {
	return "responsive.mobile-density-overload"
}

// Description mengembalikan ringkasan aturan.
func (r *MobileDensityOverloadRule) Description() string {
	return "Warns when toolbars or action rows cram more than 4 interactive buttons in a single unscrollable row on mobile viewports"
}

// Category mengembalikan nama kategori rule.
func (r *MobileDensityOverloadRule) Category() string {
	return "responsive"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *MobileDensityOverloadRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *MobileDensityOverloadRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Steven Hoober (Designing for Touch - Touch Target Interference)",
			"WCAG 2.2 SC 2.5.8 (Target Size - Minimum & Spacing)",
			"Material Design 3 Mobile App Bar & Toolbar Guidelines",
		},
		CoreInvariant: "Horizontal action toolbars on mobile viewports must not cram more than 4 interactive buttons in a single rigid row without 'overflow-x-auto', 'flex-wrap', or an overflow menu.",
		Grounding: "On compact smartphone screens (360px viewport width), accommodating 5 or more buttons in a single unscrollable flex row forces button widths below 48px or induces layout squishing.\n\n" +
			"This severe spatial compression leads to:\n" +
			"1. High Error Rate / Mis-taps: Users inadvertently trigger adjacent destructive or unwanted actions due to finger pad overlap.\n" +
			"2. Text/Icon Clipping: Labels are aggressively truncated, and icon hitboxes overlap.\n\n" +
			"Best practice dictates limiting direct actions to 3-4 primary controls, wrapping the toolbar in a horizontal scroll container ('overflow-x-auto'), or collapsing secondary actions into a 'More (...)' dropdown menu.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Touch Target Mis-tap Interference",
				Severity: "MEDIUM",
				Impact:   "Users frequently tap the wrong button because adjacent targets are compressed below safe physical spacing limits.",
			},
			{
				Vector:   "Mobile Visual Clutter & Overflow",
				Severity: "LOW",
				Impact:   "Rigid toolbars cause horizontal viewport tearing or text clipping on narrow mobile devices.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Five buttons packed tightly into an unscrollable horizontal flex row",
				Code: `<div className="flex items-center gap-2 p-2">
  <button type="button">Edit</button>
  <button type="button">Salin</button>
  <button type="button">Cetak</button>
  <button type="button">Bagikan</button>
  <button type="button">Hapus</button>
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Scrollable horizontal action bar accommodating many actions comfortably",
				Code: `<div className="flex items-center gap-2 p-2 overflow-x-auto">
  <button type="button">Edit</button>
  <button type="button">Salin</button>
  <button type="button">Cetak</button>
  <button type="button">Bagikan</button>
  <button type="button">Hapus</button>
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah kontainer baris memuat > 4 tombol tanpa pembatas scroll atau wrap di mobile.
func (r *MobileDensityOverloadRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !isHorizontalFlexRow(node.Classes) {
		return nil
	}

	interactiveCount := countInteractiveChildren(node)
	if interactiveCount <= 4 {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Horizontal action container <%s> crams %d interactive buttons in an unscrollable row on mobile baseline. Buttons will be compressed below safe touch targets on 360px screens.", node.Tag, interactiveCount),
			Hint:     "Add 'overflow-x-auto' to allow smooth horizontal scrolling, use 'flex-wrap', or collapse secondary actions into a dropdown menu.",
		},
	}
}
