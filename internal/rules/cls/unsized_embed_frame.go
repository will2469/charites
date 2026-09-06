package cls

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// UnsizedEmbedFrameRule mendeteksi media embed (iframe, video, embed, object, dsb.)
// yang tidak memiliki dimensi intrinsik eksplisit atau kontainer pembungkus aspect-ratio.
type UnsizedEmbedFrameRule struct{}

// NewUnsizedEmbedFrameRule membuat instance baru UnsizedEmbedFrameRule.
func NewUnsizedEmbedFrameRule() *UnsizedEmbedFrameRule {
	return &UnsizedEmbedFrameRule{}
}

// ID mengembalikan identifier kanonikal Semgrep rule.
func (r *UnsizedEmbedFrameRule) ID() string {
	return "cls.unsized-embed-frame"
}

// Description mengembalikan ringkasan aturan.
func (r *UnsizedEmbedFrameRule) Description() string {
	return "Warns when embedded media frames lack explicit dimensions or an aspect-ratio container wrapper"
}

// Category mengembalikan nama kategori rule.
func (r *UnsizedEmbedFrameRule) Category() string {
	return "cls"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *UnsizedEmbedFrameRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *UnsizedEmbedFrameRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Cumulative Layout Shift (CLS) Metric Specification",
			"HTML Living Standard (iframe and media embedding)",
			"W3C CSS Box Sizing Module Level 4 (aspect-ratio)",
		},
		CoreInvariant: "Embedded media frames must define explicit width/height dimensions or be enclosed in an ancestor container with an aspect-ratio or bounded height reservation.",
		Grounding: "Third-party embedded frames (such as YouTube videos, Vimeo players, interactive maps, and external iframes) take significant time to establish network handshakes and negotiate player dimensions.\n\n" +
			"When an iframe is placed in the DOM without reserved box sizing, it renders at initial zero or default browser dimensions (typically 300x150px) before snapping to full player proportions, causing substantial layout shift.\n\n" +
			"Enclosing embedded frames inside a container with 'aspect-video' or providing explicit 'width' and 'height' attributes reserves the exact layout footprint in the rendering tree immediately.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Severe Layout Instability (CLS)",
				Severity: "HIGH",
				Impact:   "Late-loading iframes pop into the document flow, shifting subsequent content by hundreds of pixels.",
			},
			{
				Vector:   "Broken Responsive Player Scaling",
				Severity: "MEDIUM",
				Impact:   "Embeds lacking proper aspect-ratio wrappers can overflow narrow mobile screens.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Iframe with fluid width but missing height or aspect-ratio wrapper",
				Code:     `<iframe src="https://www.youtube.com/embed/xyz" title="Video Profil Desa" className="w-full" />`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Iframe wrapped in a container with aspect-video utility",
				Code: `<div className="w-full aspect-video">
  <iframe src="https://www.youtube.com/embed/xyz" title="Video Profil Desa" className="w-full h-full" />
</div>`,
			},
			{
				Language: "tsx",
				Comment:  "Video element with explicit width and height attributes",
				Code:     `<video src="/promo.mp4" width={640} height={360} controls className="w-full h-auto" />`,
			},
		},
	}
}

// Evaluate memeriksa apakah media embed memiliki dimensi atau kontainer pembungkus rasio aspek.
func (r *UnsizedEmbedFrameRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if !isEmbedElement(node) {
		return nil
	}

	// Jangan beri alarm palsu pada JSX spread attributes ({...props})
	if hasSpreadProps(node.Attributes) {
		return nil
	}

	// 1. Cek dimensi langsung pada elemen
	if hasDimensionAttributes(node) {
		return nil
	}

	// 2. Cek utilitas Tailwind pada elemen itu sendiri
	if hasTailwindDimensions(node) {
		return nil
	}

	// 3. Cek deklarasi inline style pada elemen
	if hasInlineDimensionStyle(node) {
		return nil
	}

	// 4. L3 Structural Graph: Cek simpul leluhur hingga 3 tingkat (wrapper aspect-video / min-h-*)
	if hasAncestorDimensionOrAspect(node, 3) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Embedded frame <%s> lacks explicit dimensions or an aspect-ratio container wrapper. Unreserved embed frames cause significant layout shift when the player or document loads.", node.Tag),
			Hint:     "Set explicit 'width' and 'height' attributes, or wrap the embed in a container with 'aspect-video' or an explicit 'min-h-*' constraint.",
		},
	}
}
