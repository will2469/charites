package semantic

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// ConsecutiveBRSpacingRule mendeteksi penggunaan tag <br> berturut-turut (consecutive <br>)
// sebagai pola simulasi spasi antar-paragraf yang melanggar WCAG 1.3.1 (Technique F34).
type ConsecutiveBRSpacingRule struct{}

// NewConsecutiveBRSpacingRule membuat instance baru dari ConsecutiveBRSpacingRule.
func NewConsecutiveBRSpacingRule() *ConsecutiveBRSpacingRule {
	return &ConsecutiveBRSpacingRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *ConsecutiveBRSpacingRule) ID() string {
	return "semantic.consecutive-br-spacing"
}

// Description mengembalikan ringkasan maksud dan tujuan aturan.
func (r *ConsecutiveBRSpacingRule) Description() string {
	return "Flags consecutive <br> tags used for paragraph spacing violating WCAG 1.3.1 (Technique F34)"
}

// Category mengembalikan nama kategori rule.
func (r *ConsecutiveBRSpacingRule) Category() string {
	return "semantic"
}

// DefaultSeverity mengembalikan tingkat keparahan default (warn).
func (r *ConsecutiveBRSpacingRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ConsecutiveBRSpacingRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C WCAG 2.2 Technique F34 (Failure of Success Criterion 1.3.1 due to using characters to make space or using <br> tags instead of structural markup like <p>)",
			"HTML Living Standard Section 4.5.27 (The br Element)",
			"WebAIM Screen Reader Usability Survey (Navigation and Empty Breaks)",
		},
		CoreInvariant: "Paragraph spacing and vertical separation must use structural container elements (<p>, <div>) with CSS spacing, never consecutive <br> tags.",
		Grounding: "The '<br>' element represents a thematic line break within a single phrase or address block (such as in postal addresses or poetry stanzas).\n\n" +
			"Charites treats consecutive '<br>' siblings as a static indicator of spacing-by-markup. When consecutive '<br>' tags are stacked to push content downward, " +
			"screen readers announce each break as a blank line or acoustic break sound, forcing assistive technology users to navigate through phantom empty elements.\n\n" +
			"Furthermore, hardcoded '<br>' line heights bypass design system spacing tokens ('space-y-*', 'gap-*') and fail to scale with responsive typography or container queries. " +
			"Paragraphs should instead be wrapped in semantic '<p>' tags and spaced using CSS or Tailwind vertical rhythm classes.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Assistive Technology Reading Degradation",
				Severity: "HIGH",
				Impact:   "Screen reader users encounter repeated 'blank' announcements that disrupt document reading comprehension.",
			},
			{
				Vector:   "Responsive Spacing Rigidity",
				Severity: "MEDIUM",
				Impact:   "Hardcoded line breaks bypass design system spacing tokens, preventing dynamic vertical rhythm scaling.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Consecutive <br> tags used as paragraph spacing (Astrades dialog specimen)",
				Code: `<>
  Apakah Anda yakin memilih kandidat ini sebagai Pimpinan Utama?
  <br />
  <br />
  <strong className="text-destructive">PERINGATAN:</strong> Tindakan ini tidak dapat dibatalkan.
</>`,
			},
			{
				Language: "astro",
				Comment:  "Multiple <br> tags simulating paragraph gaps in Astro template",
				Code: `<div>
  <p>Pengumuman penting untuk seluruh staf.</p>
  <br />
  <br />
  <br />
  <p>Harap berkumpul di aula utama pukul 09:00 WIB.</p>
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Semantic paragraphs wrapped in a container with Tailwind vertical spacing",
				Code: `<div className="space-y-3">
  <p>Apakah Anda yakin memilih kandidat ini sebagai Pimpinan Utama?</p>
  <p>
    <strong className="text-destructive">PERINGATAN:</strong> Tindakan ini tidak dapat dibatalkan.
  </p>
</div>`,
			},
			{
				Language: "tsx",
				Comment:  "Permitted single line break within address semantics",
				Code: `<address>
  Jl. Jenderal Sudirman No. 42<br />
  Jakarta Pusat, 10220<br />
  Indonesia
</address>`,
			},
		},
	}
}

// Evaluate memeriksa apakah direct children dari elemen/fragment memuat rentetan consecutive <br>.
// Mematuhi finite-state run detector:
// - Hanya memeriksa direct children (tanpa recursive flattening).
// - Ignorable siblings (NodeComment dan NodeText whitespace-only) tidak memutus run.
// - Menghasilkan tepat 1 diagnostik per rentetan consecutive <br> (saat runCount mencapai 2).
func (r *ConsecutiveBRSpacingRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || (node.Type != ir.NodeElement && node.Type != ir.NodeFragment) {
		return nil
	}
	if len(node.Children) < 2 {
		return nil
	}

	var diags []ir.Diagnostic
	runCount := 0

	for _, child := range node.Children {
		if child == nil {
			continue
		}

		// Ignorable siblings: komentar dan teks whitespace-only tidak memutus atau menambah run
		if child.Type == ir.NodeComment {
			continue
		}
		if child.Type == ir.NodeText && strings.TrimSpace(child.RawClasses) == "" {
			continue
		}

		// Direct <br> element menambah run count
		if child.Type == ir.NodeElement && strings.EqualFold(child.Tag, "br") {
			runCount++
			// Tepat 1 diagnostik per rentetan run, berjangkar pada <br> kedua
			if runCount == 2 {
				diags = append(diags, ir.Diagnostic{
					Line:     child.Span.Line,
					Column:   child.Span.Column,
					Rule:     r.ID(),
					Severity: r.DefaultSeverity(),
					Message:  "Consecutive '<br>' elements detected as a spacing pattern; use structural elements and CSS for paragraph separation.",
					Hint:     "Replace consecutive '<br>' tags with structural markup (e.g. '<p>' elements) styled with spacing classes ('space-y-*', 'gap-*') (WCAG 1.3.1 Technique F34).",
				})
			}
			continue
		}

		// Teks bermakna atau elemen non-<br> me-reset run
		runCount = 0
	}

	return diags
}
