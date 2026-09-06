package lcp

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// HeavyRasterLCPAssetRule mendeteksi elemen gambar kandidat LCP yang merujuk pada berkas
// format raster kuno (.png, .bmp, .tiff, .gif) tanpa penyediaan varian modern seperti WebP atau AVIF.
type HeavyRasterLCPAssetRule struct{}

// NewHeavyRasterLCPAssetRule membuat instance baru dari HeavyRasterLCPAssetRule.
func NewHeavyRasterLCPAssetRule() *HeavyRasterLCPAssetRule {
	return &HeavyRasterLCPAssetRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *HeavyRasterLCPAssetRule) ID() string {
	return "lcp.heavy-raster-lcp-asset"
}

// Description mengembalikan ringkasan aturan.
func (r *HeavyRasterLCPAssetRule) Description() string {
	return "LCP candidate image uses legacy uncompressed raster format (.png, .bmp, .tiff, .gif); modern formats like WebP or AVIF should be served to reduce transfer size"
}

// Category mengembalikan nama kategori rule.
func (r *HeavyRasterLCPAssetRule) Category() string {
	return "lcp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *HeavyRasterLCPAssetRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *HeavyRasterLCPAssetRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Google Chrome Core Web Vitals (Largest Contentful Paint Resource Load Duration)",
			"W3C Web Performance Working Group Media Optimization Guidelines",
			"IETF AVIF / WebP Image Compression Standards",
		},
		CoreInvariant: "Above-the-fold LCP candidate images must utilize next-generation compressed formats (WebP, AVIF) rather than legacy uncompressed raster formats (.png, .bmp, .tiff, .gif) to minimize byte transfer payload.",
		Grounding: "Serving high-resolution photographs or hero imagery in legacy raster formats such as PNG or uncompressed BMP results in massive byte payloads (often 2MB-5MB per image).\n\n" +
			"Next-generation formats such as WebP and AVIF provide superior lossy and lossless compression algorithms, reducing image file sizes by 30% to 70% compared to PNG and JPEG without perceptual visual degradation.\n\n" +
			"For above-the-fold hero images that dictate the LCP metric, reducing file transfer size directly accelerates the Resource Load Duration phase over bandwidth-constrained networks.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Resource Load Duration Delay",
				Severity: "HIGH",
				Impact:   "Downloading uncompressed 2MB-5MB PNG/BMP images over 4G/cellular connections introduces 800ms-3000ms delay to LCP.",
			},
			{
				Vector:   "Memory Footprint & GPU Texture Pressure",
				Severity: "MEDIUM",
				Impact:   "Large uncompressed raster graphics consume excessive client RAM and GPU texture memory during decode and compositing.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Critical hero image served as an uncompressed 3MB PNG file",
				Code: `<section className="hero-section" data-perf-role="hero">
  <img
    src="/assets/hero-banner.png"
    alt="Hero Showcase"
    fetchpriority="high"
    className="w-full h-auto"
  />
</section>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Hero image converted to compressed modern WebP format",
				Code: `<section className="hero-section" data-perf-role="hero">
  <img
    src="/assets/hero-banner.webp"
    alt="Hero Showcase"
    fetchpriority="high"
    className="w-full h-auto"
  />
</section>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen kandidat LCP menggunakan format raster kuno (.png, .bmp, .tiff, .gif).
func (r *HeavyRasterLCPAssetRule) Evaluate(node *ir.Node) []ir.Diagnostic {
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
	if !ok {
		return nil
	}

	ext, heavy := isHeavyLegacyRasterFormat(rawSrc)
	if !heavy {
		return nil
	}

	// Pengecualian: Ikon atau lencana kecil (<= 64px) yang biasanya menggunakan PNG transparan
	if w, h, fixed := isFixedDimensionImage(node); fixed && w > 0 && w <= 64 && (h == 0 || h <= 64) {
		return nil
	}

	// Pengecualian: Jika dibungkus dalam <picture> dengan <source type="image/avif"> atau type="image/webp"
	if isPictureWithModernSource(node) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("LCP candidate image '<%s>' uses legacy uncompressed raster format ('%s'). Legacy formats substantially increase network transfer payload and inflate Resource Load Duration.", node.Tag, ext),
			Hint:     "Convert the asset to a modern format such as WebP or AVIF, or wrap it in a '<picture>' element providing modern '<source>' elements with 'type=\"image/webp\"' or 'type=\"image/avif\"'.",
		},
	}
}
