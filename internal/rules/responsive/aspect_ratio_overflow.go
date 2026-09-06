package responsive

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// AspectRatioOverflowRule mendeteksi elemen media atau kontainer dengan rasio aspek kaku ('aspect-*')
// yang dipadukan dengan tinggi statis ('h-96', 'h-[500px]') tanpa batas lebar fluida ('w-full', 'max-w-full'),
// memicu konflik geometris dan horizontal blowout pada viewport sempit.
type AspectRatioOverflowRule struct{}

// NewAspectRatioOverflowRule membuat instance baru dari AspectRatioOverflowRule.
func NewAspectRatioOverflowRule() *AspectRatioOverflowRule {
	return &AspectRatioOverflowRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *AspectRatioOverflowRule) ID() string {
	return "responsive.aspect-ratio-overflow"
}

// Description mengembalikan ringkasan aturan.
func (r *AspectRatioOverflowRule) Description() string {
	return "Warns against fixed aspect-ratio combined with rigid static heights without fluid width boundaries on mobile viewports"
}

// Category mengembalikan nama kategori rule.
func (r *AspectRatioOverflowRule) Category() string {
	return "responsive"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *AspectRatioOverflowRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *AspectRatioOverflowRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C CSS Box Sizing Module Level 4 (The aspect-ratio Property)",
			"WCAG 2.2 SC 1.4.10 (Reflow - Level AA)",
			"Responsive Media & Video Embed Best Practices",
		},
		CoreInvariant: "Elements specifying an explicit 'aspect-*' ratio must not pair it with a rigid fixed height without fluid width constraints ('w-full' or 'max-w-full'), which forces width computation to expand beyond narrow mobile screens.",
		Grounding: "The CSS 'aspect-ratio' property computes the corresponding dimension when one dimension is defined. When an element specifies 'aspect-video' (16/9) and also sets 'h-[450px]' without a fluid width boundary ('w-full' or 'max-w-full'), the browser calculates width as 450 * (16/9) = 800px.\n\n" +
			"On a 360px mobile screen, an 800px computed width immediately blows out the layout horizontally, forcing horizontal scrolling and clipping sibling elements.\n\n" +
			"Charites detects conflicting aspect-ratio and rigid height definitions, recommending fluid widths ('w-full aspect-video') or letting height derive naturally from fluid container width.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Massive Horizontal Layout Blowout via Derived Aspect Width",
				Severity: "HIGH",
				Impact:   "Derived width expands to 800px+ on mobile screens when static height is combined with aspect-ratio.",
			},
			{
				Vector:   "Conflicting Spatial Dimension Constraints",
				Severity: "MEDIUM",
				Impact:   "Media elements distort or overflow their parent grid/flex containers.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Aspect ratio paired with rigid fixed height forcing excessive computed width",
				Code: `<div className="aspect-video h-96 bg-black rounded-lg">
  <video src="/hero.mp4" controls />
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Fluid aspect ratio deriving height naturally from available width",
				Code: `<div className="w-full aspect-video bg-black rounded-lg">
  <video src="/hero.mp4" controls className="w-full h-full object-cover" />
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen menggabungkan aspect-ratio dengan tinggi statis kaku tanpa pembatas lebar.
func (r *AspectRatioOverflowRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !hasConflictingAspectRatio(node.Classes) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Element <%s> combines an 'aspect-*' ratio with a rigid fixed height without fluid width constraints ('w-full' or 'max-w-full'). On mobile viewports, the computed width will blow out horizontally.", node.Tag),
			Hint:     "Set a fluid width ('w-full aspect-video') and allow height to be computed dynamically from the container width.",
		},
	}
}
