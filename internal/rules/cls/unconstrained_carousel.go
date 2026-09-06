package cls

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// UnconstrainedCarouselRule mendeteksi track carousel atau slider yang tidak memiliki
// pembatas tinggi vertikal atau rasio aspek slide, berisiko memicu lonjakan tinggi saat slide berganti.
type UnconstrainedCarouselRule struct{}

// NewUnconstrainedCarouselRule membuat instance baru UnconstrainedCarouselRule.
func NewUnconstrainedCarouselRule() *UnconstrainedCarouselRule {
	return &UnconstrainedCarouselRule{}
}

// ID mengembalikan identifier kanonikal Semgrep rule.
func (r *UnconstrainedCarouselRule) ID() string {
	return "cls.unconstrained-carousel"
}

// Description mengembalikan ringkasan aturan.
func (r *UnconstrainedCarouselRule) Description() string {
	return "Warns when carousel or slider containers lack bounded height or slide aspect-ratio constraints"
}

// Category mengembalikan nama kategori rule.
func (r *UnconstrainedCarouselRule) Category() string {
	return "cls"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *UnconstrainedCarouselRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *UnconstrainedCarouselRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Cumulative Layout Shift (CLS) Metric Specification",
			"W3C CSS Scroll Snap Module Level 1",
			"W3C CSS Box Sizing Module Level 4 (aspect-ratio)",
		},
		CoreInvariant: "Carousel and slider viewport tracks must constrain container height or bind slide items to fixed aspect ratios to prevent vertical reflow during slide transitions.",
		Grounding: "Horizontal scrolling tracks and carousels render dynamic collections of cards, banners, or images.\n\n" +
			"When the carousel track lacks an explicit height (e.g. 'h-64' or 'min-h-[300px]') and slides do not have locked aspect ratios, incoming slides with varying image proportions or dynamic text will force the entire container to expand or collapse vertically.\n\n" +
			"Fixing the container height or assigning 'aspect-video' / 'aspect-square' to slide items ensures layout stability throughout horizontal panning.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Vertical Container Height Jitter",
				Severity: "MEDIUM",
				Impact:   "Slide transitions with varying content heights push subsequent page content up and down.",
			},
			{
				Vector:   "Cumulative Layout Shift (CLS)",
				Severity: "HIGH",
				Impact:   "Carousel height adjustments contribute cumulative shift points during user scrolling.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Horizontal snap container without container height or slide aspect-ratio",
				Code: `<div className="flex overflow-x-auto snap-x">
  {slides.map(s => <img key={s.id} src={s.url} alt={s.title} />)}
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Carousel container with explicit height constraint",
				Code: `<div className="flex overflow-x-auto snap-x h-64 md:h-96 w-full">
  {slides.map(s => (
    <div key={s.id} className="snap-center shrink-0 w-full h-full">
      <img src={s.url} alt={s.title} className="w-full h-full object-cover" />
    </div>
  ))}
</div>`,
			},
			{
				Language: "tsx",
				Comment:  "Carousel slide items locked with aspect-video utility",
				Code: `<div className="flex overflow-x-auto snap-x w-full">
  {slides.map(s => (
    <div key={s.id} className="snap-center shrink-0 w-80 aspect-video">
      <img src={s.url} alt={s.title} className="w-full h-full object-cover" />
    </div>
  ))}
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah track carousel memiliki pembatas ketinggian.
func (r *UnconstrainedCarouselRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if !isCarouselTrack(node) {
		return nil
	}

	// Jangan beri alarm palsu pada JSX spread attributes ({...props})
	if hasSpreadProps(node.Attributes) {
		return nil
	}

	// 1. Cek pembatas ketinggian pada kontainer itu sendiri
	if hasBoundedHeight(node) {
		return nil
	}

	// 2. Cek apakah slide anak memiliki ketinggian terikat atau rasio aspek
	for _, child := range node.Children {
		if child.Type == ir.NodeElement {
			if hasBoundedHeight(child) || hasTailwindDimensions(child) || hasDimensionAttributes(child) {
				return nil
			}
		}
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Carousel/slider container <%s> lacks bounded vertical dimensions ('h-*', 'min-h-*', or 'aspect-*'). Dynamically loaded slides or slide transitions may cause vertical layout jumps.", node.Tag),
			Hint:     "Set an explicit container height (e.g. 'h-64 md:h-96') or lock slide aspect ratios with 'aspect-video' or 'aspect-square'.",
		},
	}
}
