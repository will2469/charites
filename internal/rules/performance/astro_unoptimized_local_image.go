package performance

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// AstroUnoptimizedLocalImageRule mengaudit penggunaan tag <img> mentah pada gambar lokal di proyek Astro.
type AstroUnoptimizedLocalImageRule struct{}

// NewAstroUnoptimizedLocalImageRule membuat instance baru dari AstroUnoptimizedLocalImageRule.
func NewAstroUnoptimizedLocalImageRule() *AstroUnoptimizedLocalImageRule {
	return &AstroUnoptimizedLocalImageRule{}
}

// ID mengembalikan identifier unik kanonikal aturan.
func (r *AstroUnoptimizedLocalImageRule) ID() string {
	return "performance.astro-unoptimized-local-image"
}

// Category mengembalikan kategori aturan ('performance').
func (r *AstroUnoptimizedLocalImageRule) Category() string {
	return "performance"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (info/advisory).
func (r *AstroUnoptimizedLocalImageRule) DefaultSeverity() ir.Severity {
	return ir.SeverityInfo
}

// Description mengembalikan deskripsi tujuan dan konteks aturan.
func (r *AstroUnoptimizedLocalImageRule) Description() string {
	return "Menganjurkan pemakaian komponen <Image /> dari astro:assets pada gambar lokal guna mengaktifkan konversi format modern dan kompresi build."
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *AstroUnoptimizedLocalImageRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Astro Asset Pipeline Best Practices ('astro:assets' Image & Picture)",
			"Core Web Vitals Largest Contentful Paint (LCP) Image Payload Optimization",
			"W3C Next-Gen Responsive Image Delivery Standards (WebP/AVIF)",
		},
		CoreInvariant: "Local raster image assets in Astro templates should be rendered via '<Image />' from 'astro:assets' rather than raw '<img>' tags to leverage automated build-time format conversion and dimension inference.",
		Grounding: "Astro provides a native image optimization pipeline through the `astro:assets` module.\n\n" +
			"Using a raw HTML `<img>` tag pointing to a local file path (`src=\"../assets/banner.png\"`) completely bypasses this pipeline, serving uncompressed, legacy formats (PNG/JPEG) with no automatic width/height dimension injection.\n\n" +
			"Migrating to `<Image />` allows Astro to automatically convert images to AVIF/WebP, generate responsive srcset attributes, and prevent Cumulative Layout Shift (CLS) by inferring exact dimensions at build time.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Inflated Asset Payload",
				Severity: "LOW",
				Impact:   "Serves unoptimized PNG/JPEG images that are 40-70% larger than modern WebP/AVIF equivalents.",
			},
			{
				Vector:   "Missing Intrinsic Aspect Ratio",
				Severity: "LOW",
				Impact:   "Raw img tags without width and height attributes cause layout shifts during image load.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Tag img mentah melewatkan kompresi build-time Astro",
				Code: `<!-- Advisory: Tag img mentah pada path lokal -->
<img src="../assets/product-hero.png" alt="Produk Baru" />`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Memanfaatkan komponen Image bawaan astro:assets",
				Code: `---
import { Image } from 'astro:assets';
import productImg from '../assets/product-hero.png';
---
<Image src={productImg} alt="Produk Baru" />`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen <img> menggunakan path gambar lokal mentah.
func (r *AstroUnoptimizedLocalImageRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || node.Tag != "img" {
		return nil
	}

	src, isLocal := isLocalRawImage(node)
	if !isLocal {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Raw '<img>' tag references local image asset '%s' directly, bypassing Astro build-time image compression, modern format conversion (WebP/AVIF), and automatic dimension extraction.", src),
			Hint:     "Import the image and render it via '<Image />' from 'astro:assets' to enable automatic build-time optimization.",
		},
	}
}
