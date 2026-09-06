package lcp

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// ImageSourceDensityMismatchRule mendeteksi elemen gambar kandidat LCP berdimensi tetap
// yang tidak menyertakan deskriptor densitas piksel (1x, 2x) pada atribut srcset.
type ImageSourceDensityMismatchRule struct{}

// NewImageSourceDensityMismatchRule membuat instance baru dari ImageSourceDensityMismatchRule.
func NewImageSourceDensityMismatchRule() *ImageSourceDensityMismatchRule {
	return &ImageSourceDensityMismatchRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *ImageSourceDensityMismatchRule) ID() string {
	return "lcp.image-source-density-mismatch"
}

// Description mengembalikan ringkasan aturan.
func (r *ImageSourceDensityMismatchRule) Description() string {
	return "Fixed-dimension LCP candidate image lacks aligned '1x, 2x' pixel density descriptors in 'srcset', risking blurry rendering or unoptimized asset delivery on high-DPI screens"
}

// Category mengembalikan nama kategori rule.
func (r *ImageSourceDensityMismatchRule) Category() string {
	return "lcp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (info/advisory).
func (r *ImageSourceDensityMismatchRule) DefaultSeverity() ir.Severity {
	return ir.SeverityInfo
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ImageSourceDensityMismatchRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Google Chrome Core Web Vitals (Largest Contentful Paint Resource Optimization)",
			"HTML Living Standard Pixel Density Descriptors (1x, 2x)",
			"W3C Web Performance Working Group High-DPI Media Guidelines",
		},
		CoreInvariant: "Fixed-dimension LCP candidate images must specify aligned '1x' and '2x' pixel density descriptors in 'srcset' to prevent blurry rendering on high-DPI displays while avoiding oversized single asset downloads on standard displays.",
		Grounding: "Fixed-dimension images (such as brand masthead logos, author avatar badges, or feature icons with fixed width and height) do not scale fluidly with viewport width.\n\n" +
			"Serving a single resolution asset forces high-DPI (Retina) screens to upscale lower-resolution images, causing visual blurriness, or forces standard 1x screens to download an unnecessarily large 2x/3x asset.\n\n" +
			"Providing a 'srcset' attribute with '1x' and '2x' density descriptors enables the browser to automatically select the optimal resolution based on the device pixel ratio (DPR).",
		Risks: []ir.RiskItem{
			{
				Vector:   "Visual Degradation on High-DPI Screens",
				Severity: "MEDIUM",
				Impact:   "Single 1x assets appear blurry or pixelated on modern smartphone and laptop displays with DPR >= 2.",
			},
			{
				Vector:   "Wasted Bandwidth on Standard Displays",
				Severity: "LOW",
				Impact:   "Single 2x assets download double the required byte payload on standard 1x displays.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Fixed-dimension logo in hero masthead loading a single oversized 2000px asset without 1x/2x descriptors",
				Code: `<header data-perf-role="hero">
  <img
    src="/assets/logo-2000.webp"
    width="120"
    height="40"
    alt="Corporate Logo"
    fetchpriority="high"
  />
</header>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Fixed-dimension logo configured with 1x and 2x pixel density descriptors",
				Code: `<header data-perf-role="hero">
  <img
    src="/assets/logo-120.webp"
    srcset="/assets/logo-120.webp 1x, /assets/logo-240.webp 2x"
    width="120"
    height="40"
    alt="Corporate Logo"
    fetchpriority="high"
  />
</header>`,
			},
		},
	}
}

// Evaluate memeriksa apakah gambar kandidat LCP berdimensi tetap tidak memiliki deskriptor densitas 1x, 2x.
func (r *ImageSourceDensityMismatchRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !isMediaTag(node.Tag) {
		return nil
	}

	isHero, _ := isLCPCandidate(node)
	if !isHero {
		return nil
	}

	if node.Attributes == nil {
		return nil
	}

	rawSrc, ok := node.Attributes["src"]
	if !ok || len(cleanAttrVal(rawSrc)) == 0 {
		return nil
	}

	// Format SVG bersifat vektor dan tidak memerlukan varian densitas piksel
	if isSVGFile(rawSrc) {
		return nil
	}

	// Aturan ini HANYA berlaku untuk elemen dengan dimensi piksel tetap (fixed dimensions)
	w, h, fixed := isFixedDimensionImage(node)
	if !fixed {
		return nil
	}

	rawSrcset, hasSrcset := getSrcsetAttribute(node.Attributes)
	if hasSrcset {
		if hasDensityDescriptors(rawSrcset) {
			return nil
		}
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Fixed-dimension LCP candidate image (%dx%d) uses width descriptors in 'srcset' instead of '1x, 2x' pixel density descriptors.", w, h),
				Hint:     "Replace width descriptors with density descriptors (e.g. 'srcset=\"image.webp 1x, image@2x.webp 2x\"') for fixed-dimension media.",
			},
		}
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Fixed-dimension LCP candidate image (%dx%d) lacks '1x, 2x' pixel density descriptors in 'srcset', risking blurry rendering on high-DPI screens or unoptimized asset delivery.", w, h),
			Hint:     "Provide a 'srcset' attribute with '1x' and '2x' descriptors (e.g. 'srcset=\"/logo.webp 1x, /logo@2x.webp 2x\"') to support high-DPI displays optimally.",
		},
	}
}
