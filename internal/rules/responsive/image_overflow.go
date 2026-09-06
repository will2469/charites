package responsive

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// ImageOverflowRule mendeteksi elemen media (img, video, svg, canvas) yang memiliki dimensi statis
// lebih besar dari 320px tanpa kelas penskalaan fluida responsif 'max-w-full' atau 'w-full'.
type ImageOverflowRule struct{}

// NewImageOverflowRule membuat instance baru dari ImageOverflowRule.
func NewImageOverflowRule() *ImageOverflowRule {
	return &ImageOverflowRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *ImageOverflowRule) ID() string {
	return "responsive.image-overflow"
}

// Description mengembalikan ringkasan aturan.
func (r *ImageOverflowRule) Description() string {
	return "Warns when media elements with large fixed dimensions lack responsive max-w-full scaling"
}

// Category mengembalikan nama kategori rule.
func (r *ImageOverflowRule) Category() string {
	return "responsive"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *ImageOverflowRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ImageOverflowRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"HTML Living Standard (Embedded Content: img, video, svg)",
			"Web.dev Responsive Media & Core Web Vitals (CLS Prevention)",
		},
		CoreInvariant: "Media elements with explicit width dimensions exceeding 320px must declare 'max-w-full' or 'w-full' to prevent horizontal viewport tearing on mobile screens.",
		Grounding: "Specifying explicit 'width' and 'height' attributes on media elements is recommended for Core Web Vitals to reserve aspect ratio boxes and prevent Cumulative Layout Shift (CLS).\n\n" +
			"However, when large static dimensions (e.g. width={1200}) lack responsive CSS scaling ('max-w-full h-auto'), mobile browsers render the media at full absolute physical pixel width, breaking outside narrow 360px viewports and causing severe horizontal scrolling.\n\n" +
			"Applying 'max-w-full h-auto' preserves CLS aspect ratio benefits while ensuring the media smoothly downsizes to fit compact screens.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Mobile Viewport Tearing",
				Severity: "MEDIUM",
				Impact:   "Large unconstrained images expand outside mobile viewport boundaries, forcing horizontal scrollbars.",
			},
			{
				Vector:   "Aspect Ratio Distortion",
				Severity: "LOW",
				Impact:   "Images constrained by height but not width stretch disproportionately on narrow viewports.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Media element with large width attribute lacking max-w-full",
				Code:     `<img src="/hero-desa.jpg" width={1200} height={800} alt="Pemandangan Desa" />`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Responsive media element with max-w-full and h-auto",
				Code:     `<img className="max-w-full h-auto" src="/hero-desa.jpg" width={1200} height={800} alt="Pemandangan Desa" />`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen media memiliki lebar > 320px tanpa max-w-full atau w-full.
func (r *ImageOverflowRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || !isMediaTag(node.Tag) {
		return nil
	}

	w, ok := extractMediaWidth(node)
	if !ok || w <= 320 {
		return nil
	}

	if hasResponsiveMediaScaling(node) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Media element <%s> has fixed width (%dpx) exceeding mobile viewport (320px) without responsive scaling ('max-w-full' or 'w-full'). The image will render oversized and break mobile layouts.", node.Tag, w),
			Hint:     "Add 'max-w-full h-auto' to ensure the media element scales down responsively on narrow viewports while preserving its aspect ratio.",
		},
	}
}
