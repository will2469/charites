package a11y

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// InputIOSZoomHazardRule mendeteksi kontrol input dengan ukuran font mobile di bawah 16px (< 1rem)
// yang memicu auto-zoom paksa oleh Safari iOS WebKit saat fokus.
type InputIOSZoomHazardRule struct{}

// NewInputIOSZoomHazardRule membuat instance baru InputIOSZoomHazardRule.
func NewInputIOSZoomHazardRule() *InputIOSZoomHazardRule {
	return &InputIOSZoomHazardRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat a11y.input-ios-zoom-hazard.
func (r *InputIOSZoomHazardRule) ID() string {
	return "a11y.input-ios-zoom-hazard"
}

// Description mengembalikan deskripsi ringkas aturan.
func (r *InputIOSZoomHazardRule) Description() string {
	return "Prevents forced Safari iOS viewport auto-zoom by requiring at least 16px font size on inputs on mobile viewports"
}

// Category mengembalikan nama kategori rule.
func (r *InputIOSZoomHazardRule) Category() string {
	return "a11y"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *InputIOSZoomHazardRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *InputIOSZoomHazardRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"WCAG 2.2 Success Criterion 1.4.4 (Resize Text)",
			"Apple WebKit iOS Form Viewport Behavior (16px Input Zoom Threshold)",
			"Apple Human Interface Guidelines (Form Inputs & Readability)",
		},
		CoreInvariant: "Form text inputs (<input>, <select>, <textarea>) must maintain at least 16px font size on mobile viewports ('text-base' or larger) to prevent mandatory iOS Safari viewport auto-zoom.",
		Grounding: "Apple iOS WebKit automatically forces a viewport zoom whenever a user focuses on a text input whose computed font-size is strictly below 16px (1rem).\n\n" +
			"This abrupt viewport magnification causes severe usability issues:\n" +
			"1. Broken Layout: Page content shifts horizontally, hiding submit buttons and neighboring fields off-screen.\n" +
			"2. Manual Pinch-to-Zoom Fatigue: After typing, the user is forced to perform manual two-finger pinch-to-zoom gestures to restore the original page scale.\n" +
			"3. Cognitive Disorientation: Rapid layout shifting during form completion increases cognitive friction and abandonment rates.\n\n" +
			"Responsive Defense: Providing 'text-base sm:text-sm' provides a safe 16px baseline on mobile viewports while preserving compact 14px typography on desktop screens.",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Small 14px text input on mobile triggers Safari auto-zoom",
				Code:     `<input className="text-sm px-3 py-2 border rounded" placeholder="Email" />`,
			},
			{
				Language: "astro",
				Comment:  "Dropdown select with 12px text size without mobile 16px baseline",
				Code:     `<select class="text-xs border rounded p-2"><option>Pilih Opsi</option></select>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "16px baseline on mobile, safely scaled to 14px on desktop",
				Code:     `<input className="text-base sm:text-sm px-3.5 py-2.5 border rounded" placeholder="Email" />`,
			},
			{
				Language: "astro",
				Comment:  "Standard 16px base font size on select control",
				Code:     `<select class="text-base border rounded p-2"><option>Pilih Opsi</option></select>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Forced iOS Viewport Auto-Zoom",
				Severity: "HIGH",
				Impact:   "Mobile Safari users experience jarring zoom-ins that throw layout elements out of visible viewport.",
			},
			{
				Vector:   "Form Abandonment",
				Severity: "MEDIUM",
				Impact:   "Awkward pinch-to-zoom requirements cause friction during transactional checkout and authentication forms.",
			},
		},
	}
}

// Evaluate memeriksa ukuran font mobile pada elemen input formulir.
// Mematuhi 0 B/op, 0 allocs/op pada node bersih (QUAL-03).
func (r *InputIOSZoomHazardRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || len(node.Classes) == 0 {
		return nil
	}

	if !IsTextualInput(node) {
		return nil
	}

	hasMobileFont, hasSafeBase, minSize, minClass := inspectInputFontSizes(node.Classes)

	// Jika ada deklarasi ukuran font pada mobile, dan tidak ada kelas base/lebih besar,
	// serta ukuran terkecil di bawah 16px -> Pelanggaran auto-zoom iOS
	if hasMobileFont && !hasSafeBase && minSize < 16.0 {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Input <%s> uses mobile font size '%s' (%.0fpx < 16px), triggering forced iOS Safari viewport auto-zoom", node.Tag, minClass, minSize),
				Hint:     "Use 'text-base sm:text-sm' to provide a 16px baseline on mobile while scaling down on larger viewports.",
			},
		}
	}

	return nil
}

func inspectInputFontSizes(classes []string) (hasMobileFont bool, hasSafeBase bool, minSize float64, minClass string) {
	minSize = 999.0

	for _, class := range classes {
		// Abaikan utilitas yang memiliki prefix responsive min-width (sm:, md:, lg:, dsb)
		// karena utilitas tersebut hanya berlaku di viewport desktop/tablet, bukan mobile base
		if HasMinBreakpointPrefix(class) {
			continue
		}

		base := StripVariantsOnlyBase(class)
		px, ok := ParseFontSizeToPx(base)
		if !ok {
			continue
		}

		hasMobileFont = true
		if px >= 16.0 {
			hasSafeBase = true
		} else if px < minSize {
			minSize = px
			minClass = class
		}
	}

	return hasMobileFont, hasSafeBase, minSize, minClass
}
