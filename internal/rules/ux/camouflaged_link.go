package ux

import (
	"github.com/will2469/charites/internal/ir"
)

// CamouflagedLinkRule mendeteksi tautan inline dalam paragraf atau teks prosa yang hanya mengandalkan warna
// tanpa affordance non-warna persisten (underline atau border-b), melanggar WCAG 2.2 SC 1.4.1 dan prinsip affordance.
type CamouflagedLinkRule struct{}

// NewCamouflagedLinkRule membuat instance baru dari CamouflagedLinkRule.
func NewCamouflagedLinkRule() *CamouflagedLinkRule {
	return &CamouflagedLinkRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *CamouflagedLinkRule) ID() string {
	return "ux.camouflaged-link"
}

// Description mengembalikan ringkasan aturan.
func (r *CamouflagedLinkRule) Description() string {
	return "Warns when inline prose links rely solely on color without persistent underline or non-color affordance"
}

// Category mengembalikan nama kategori rule.
func (r *CamouflagedLinkRule) Category() string {
	return "ux"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *CamouflagedLinkRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *CamouflagedLinkRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Web Content Accessibility Guidelines (WCAG 2.2 SC 1.4.1 Use of Color - Level A)",
			"Gestalt Law of Similarity & Visual Affordance Principles (Norman, 2013)",
			"Nielsen Norman Group (Hyperlink Affordance Guidelines: Persistent Underlines)",
		},
		CoreInvariant: "Inline prose hyperlinks embedded within body text must provide a persistent non-color visual cue ('underline' or 'border-b') in idle state rather than relying solely on text color or hover-only transitions.",
		Grounding: "WCAG 2.2 Success Criterion 1.4.1 mandates that color must not be used as the only visual means of conveying information, indicating an action, prompting a response, or distinguishing a visual element.\n\n" +
			"When an inline link inside body copy removes underlines ('no-underline') or only shows underlines on hover ('hover:underline') while displaying text in primary brand color, users with color vision deficiency (protanopia, deuteranopia, tritanopia) or those using monitors with non-calibrated contrast cannot perceive the text as an interactive link.\n\n" +
			"Furthermore, according to NN/g research, static text that lacks standard underline affordance forces users to hunt and peck with the cursor, drastically diminishing reading fluency and scan efficiency. Persistent underlines or bottom border decorations provide immediate visual affordance.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Inaccessible Color-Dependent Affordance",
				Severity: "HIGH",
				Impact:   "Color-blind users cannot distinguish clickable inline links from regular static prose text.",
			},
			{
				Vector:   "Reduced Reading & Scanning Fluency",
				Severity: "MEDIUM",
				Impact:   "Users fail to discover important inline reference links, increasing task failure rates.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Inline prose link with no-underline relying only on text color and hover",
				Code: `<p className="text-base text-neutral-700">
  Untuk informasi lebih lengkap, silakan kunjungi
  <a href="/panduan" className="text-primary hover:underline"> buku panduan warga</a>.
</p>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Persistent underline in idle state ensures immediate non-color affordance",
				Code: `<p className="text-base text-neutral-700">
  Untuk informasi lebih lengkap, silakan kunjungi
  <a href="/panduan" className="text-primary underline decoration-primary/50 hover:decoration-primary"> buku panduan warga</a>.
</p>`,
			},
		},
	}
}

// Evaluate menganalisis tautan inline apakah memiliki affordance garis bawah persisten.
func (r *CamouflagedLinkRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !isNavLinkNode(node) {
		return nil
	}

	// Lewatkan jika ditata menyerupai tombol fisik
	if isStyledAsButton(node) {
		return nil
	}

	// Hanya evaluasi jika berada dalam konteks bacaan paragraf / teks prosa
	if !isProseContext(node) {
		return nil
	}

	// Cek apakah memiliki affordance underline persisten di state idle
	if hasPersistentUnderlineAffordance(node.Classes) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Inline link inside prose lacks persistent non-color affordance ('underline' or 'border-b') in idle state, relying solely on color (violates WCAG 2.2 SC 1.4.1).",
			Hint:     "Add 'underline' or 'border-b' in idle state to provide a persistent non-color affordance for readability and color-blind users.",
		},
	}
}
