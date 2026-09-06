package responsive

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// ContainerOverconstraintRule mendeteksi kontainer dengan batasan lebar dan padding horizontal berlebih
// pada mobile baseline yang menyisakan area konten terlalu sempit (< 280px) pada layar smartphone.
type ContainerOverconstraintRule struct{}

// NewContainerOverconstraintRule membuat instance baru dari ContainerOverconstraintRule.
func NewContainerOverconstraintRule() *ContainerOverconstraintRule {
	return &ContainerOverconstraintRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *ContainerOverconstraintRule) ID() string {
	return "responsive.container-overconstraint"
}

// Description mengembalikan ringkasan aturan.
func (r *ContainerOverconstraintRule) Description() string {
	return "Warns against excessive mobile horizontal padding or overconstrained widths that pinch usable content width below 280px on smartphone viewports"
}

// Category mengembalikan nama kategori rule.
func (r *ContainerOverconstraintRule) Category() string {
	return "responsive"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *ContainerOverconstraintRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ContainerOverconstraintRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"WCAG 2.2 SC 1.4.10 (Reflow - Level AA)",
			"Responsive Web Design Usable Width Baseline (320px - 360px)",
			"Tailwind CSS Layout Container Best Practices",
		},
		CoreInvariant: "Mobile baseline containers must not combine narrow width constraints with excessive horizontal padding (e.g. 'px-16', 'px-20', 'max-w-xs px-12') without responsive breakpoint prefixes, ensuring usable width stays above 280px.",
		Grounding: "On standard smartphones with a 360px wide screen (such as Galaxy A series and baseline Android devices), excessive horizontal padding like 'px-16' (64px each side = 128px total) reduces the usable reading width to just 232px.\n\n" +
			"When combined with narrow constraints like 'max-w-xs' (320px) and large padding, content gets severely cramped, forcing awkward line breaks, clipped tables, and unreadable text.\n\n" +
			"Charites flags unprefixed heavy horizontal padding on container elements, urging developers to start with compact padding on mobile (e.g. 'px-4') and scale up via responsive prefixes ('md:px-16').",
		Risks: []ir.RiskItem{
			{
				Vector:   "Severe Content Cramping & Layout Distortion",
				Severity: "MEDIUM",
				Impact:   "Text blocks and interactive widgets become vertically stretched with single-word line breaks.",
			},
			{
				Vector:   "Unnecessary Mobile Space Wastage",
				Severity: "LOW",
				Impact:   "More than 35% of the mobile screen width is wasted on dead whitespace padding margins.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Container applying desktop-sized horizontal padding on mobile baseline",
				Code: `<div className="container mx-auto px-16 py-8">
  <h1 className="text-2xl font-bold">Judul Halaman Warga</h1>
  <p>Deskripsi layanan kependudukan desa.</p>
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Fluid padding scaling smoothly from mobile to desktop",
				Code: `<div className="container mx-auto px-4 md:px-16 py-8">
  <h1 className="text-2xl font-bold">Judul Halaman Warga</h1>
  <p>Deskripsi layanan kependudukan desa.</p>
</div>`,
			},
		},
	}
}

// Evaluate memeriksa keberadaan padding horizontal berlebih tanpa prefix breakpoint pada kontainer.
func (r *ContainerOverconstraintRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !hasExcessiveMobilePadding(node.Classes) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Container <%s> applies excessive horizontal padding on mobile baseline (%s). On a 360px viewport, this pinches usable content width below 280px.", node.Tag, node.RawClasses),
			Hint:     "Adopt a mobile-first approach: use compact padding on mobile (e.g. 'px-4') and scale up with responsive breakpoint prefixes (e.g. 'md:px-12', 'lg:px-16').",
		},
	}
}
