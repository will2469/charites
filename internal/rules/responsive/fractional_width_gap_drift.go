package responsive

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// FractionalWidthGapDriftRule mendeteksi flex container yang mendeklarasikan celah horizontal (gap-*)
// dengan anak-anak bernilai lebar fraksional/persentase total >= 100% yang tidak dapat diserap
// secara native oleh Flexbox (karena flex-wrap aktif atau flex-shrink dimatikan).
type FractionalWidthGapDriftRule struct{}

// NewFractionalWidthGapDriftRule membuat instance baru dari FractionalWidthGapDriftRule.
func NewFractionalWidthGapDriftRule() *FractionalWidthGapDriftRule {
	return &FractionalWidthGapDriftRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *FractionalWidthGapDriftRule) ID() string {
	return "responsive.fractional-width-gap-drift"
}

// Description mengembalikan ringkasan aturan.
func (r *FractionalWidthGapDriftRule) Description() string {
	return "Warns when flex container declares a horizontal gap alongside fractional-width children totaling >= 100% with flex-wrap or disabled shrinking, causing unexpected wrapping or container blowout"
}

// Category mengembalikan nama kategori rule.
func (r *FractionalWidthGapDriftRule) Category() string {
	return "responsive"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *FractionalWidthGapDriftRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *FractionalWidthGapDriftRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C CSS Flexible Box Layout Module Level 1 (Section 9.3: Line Breaking & Section 9.7: Resolving Flexible Lengths)",
			"W3C CSS Grid Layout Module Level 2 (Section 7.2: Flexible Track Sizing)",
			"Tailwind CSS Responsive Design & Spacing Geometry",
		},
		CoreInvariant: "Flex containers declaring horizontal gaps alongside fractional-width children totaling 100% or more must not enable 'flex-wrap' or disable flex shrinking ('shrink-0' / 'flex-none'), which causes unexpected line wrapping or container blowout.",
		Grounding: "In CSS Flexbox, percentage/fractional widths (e.g. 'w-1/2', 'w-2/3', 'w-1/3') are computed against the container's inner content box, but they do NOT deduct the horizontal gap space declared on the container.\n\n" +
			"When fractional children sum to 100% (e.g. 50% + 50% = 100%) and a positive horizontal gap is declared (e.g. 'gap-4' = 16px), the total required row width becomes 100% + 16px.\n\n" +
			"Under default 'flex-wrap: nowrap' and default 'flex-shrink: 1', browser Flexbox automatically shrinks items to absorb the 16px negative free space, making it visually safe. However, when 'flex-wrap' is active, the second item is forced onto a new line (unexpected line wrapping). Similarly, when flex shrinking is disabled on fractional items ('shrink-0' or 'flex-none'), the row refuses to shrink and overflows the container.\n\n" +
			"The modern, defect-free alternative is CSS Grid (e.g. 'grid grid-cols-2 gap-4' or 'grid grid-cols-1 md:grid-cols-3 gap-4') where fraction units ('fr') automatically distribute remaining track space after deducting gap sizes.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Unexpected Line Wrapping",
				Severity: "HIGH",
				Impact:   "When 'flex-wrap' is active, columns intended to appear side-by-side break awkwardly onto new lines because fractional widths + gap exceed 100%.",
			},
			{
				Vector:   "Container Blowout & Horizontal Leak",
				Severity: "MEDIUM",
				Impact:   "When 'shrink-0' is applied to fractional flex items, the row cannot absorb gap spacing and blows out the container horizontally.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Flex container with flex-wrap and fractional children summing to 100%",
				Code: `<div className="flex flex-wrap gap-4">
  <div className="w-1/2">Kolom 1</div>
  <div className="w-1/2">Kolom 2 (terdorong ke baris baru karena gap-4)</div>
</div>`,
			},
			{
				Language: "tsx",
				Comment:  "Flex container with shrink-disabled fractional items causing horizontal blowout",
				Code: `<div className="flex gap-4">
  <div className="w-2/3 shrink-0">Kolom Utama</div>
  <div className="w-1/3 shrink-0">Sidebar</div>
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "CSS Grid allocation where fraction units (fr) automatically deduct gap width",
				Code: `<div className="grid grid-cols-1 md:grid-cols-3 gap-4">
  <div className="md:col-span-2">Kolom Utama</div>
  <div>Sidebar</div>
</div>`,
			},
			{
				Language: "tsx",
				Comment:  "Fluid flex distribution using flex-1 / basis-0",
				Code: `<div className="flex gap-4">
  <div className="flex-1 basis-0">Kolom 1</div>
  <div className="flex-1 basis-0">Kolom 2</div>
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah kontainer flex mengalami flexbox spacing math drift.
func (r *FractionalWidthGapDriftRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || len(node.Children) == 0 {
		return nil
	}

	geom := ResolveGeometry(node)
	if geom == nil {
		return nil
	}

	finding := ClassifyFlexGapDrift(geom)
	if finding == nil {
		return nil
	}

	tierLabel := "on mobile baseline"
	if finding.Tier != TierBaseline {
		tierLabel = fmt.Sprintf("at breakpoint modifier %q", finding.Tier.String()+":")
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message: fmt.Sprintf(
				"Flex container %s declares horizontal gap (%s) with child items whose fractional widths total 100%% or more, resulting in %s.",
				tierLabel,
				finding.GapVal,
				finding.Reason,
			),
			Hint: "Use CSS Grid (e.g. 'grid grid-cols-2 gap-4' or 'grid grid-cols-1 md:grid-cols-3 gap-4') where fraction units ('fr') automatically allocate space after deducting gaps, or use 'flex-1 basis-0' for fluid flex distribution.",
		},
	}
}
