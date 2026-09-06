package cls

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// UnsizedImageRule mendeteksi elemen gambar yang tidak memiliki dimensi intrinsik eksplisit,
// rasio aspek CSS, atau utilitas dimensi Tailwind untuk mereservasi ruang rendering awal.
type UnsizedImageRule struct{}

// NewUnsizedImageRule membuat instance baru UnsizedImageRule.
func NewUnsizedImageRule() *UnsizedImageRule {
	return &UnsizedImageRule{}
}

// ID mengembalikan identifier kanonikal Semgrep rule.
func (r *UnsizedImageRule) ID() string {
	return "cls.unsized-image"
}

// Description mengembalikan ringkasan aturan.
func (r *UnsizedImageRule) Description() string {
	return "Warns when image elements lack explicit dimensions, aspect-ratio, or Tailwind box sizing"
}

// Category mengembalikan nama kategori rule.
func (r *UnsizedImageRule) Category() string {
	return "cls"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *UnsizedImageRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *UnsizedImageRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Cumulative Layout Shift (CLS) Metric Specification",
			"Google Core Web Vitals Guidelines (Target CLS < 0.1)",
			"W3C CSS Box Sizing Module Level 4 (aspect-ratio)",
			"Astro Docs: Image Optimization (astro:assets)",
		},
		CoreInvariant: "Image elements must establish a statically inferable reserved rendering box via explicit width/height attributes, CSS aspect-ratio, or Tailwind sizing utilities before the binary asset is downloaded.",
		Grounding: "When browsers parse an <img> tag without explicit dimensions or an aspect-ratio reservation, the layout engine initially allocates a 0x0 pixel box.\n\n" +
			"Once the remote image file is fetched and decoded, the browser performs a sudden reflow to accommodate the intrinsic image geometry, pushing surrounding content downward. This layout instability directly penalizes Cumulative Layout Shift (CLS) scores.\n\n" +
			"Specifying width and height attributes or utilizing Tailwind 'aspect-video' / 'aspect-square' allows modern browsers to compute the aspect ratio before network I/O completes, eliminating visual jank.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Cumulative Layout Shift (CLS)",
				Severity: "HIGH",
				Impact:   "Unsized images push subsequent content down upon load, degrading Core Web Vitals and SEO rankings.",
			},
			{
				Vector:   "Accidental User Mis-clicks",
				Severity: "MEDIUM",
				Impact:   "Users attempting to tap links or buttons near loading images accidentally trigger shifted elements.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Image with fluid width but missing height or aspect-ratio reservation",
				Code:     `<img src={heroUrl} alt="Hero Banner" className="w-full h-auto" />`,
			},
			{
				Language: "astro",
				Comment:  "Standard img tag lacking width and height attributes",
				Code:     `<img src="/pemandangan-desa.jpg" alt="Pemandangan Desa" class="rounded-lg" />`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Image with explicit numeric width and height attributes",
				Code:     `<img src={heroUrl} alt="Hero Banner" width={1200} height={600} className="w-full h-auto" />`,
			},
			{
				Language: "tsx",
				Comment:  "Image with Tailwind v4 aspect-ratio utility",
				Code:     `<img src={heroUrl} alt="Hero Banner" className="w-full aspect-video object-cover" />`,
			},
			{
				Language: "tsx",
				Comment:  "Avatar image with explicit width and height sizing utilities",
				Code:     `<img src={avatarUrl} alt="Avatar" className="w-10 h-10 rounded-full" />`,
			},
		},
	}
}

// Evaluate memeriksa apakah node elemen gambar memiliki reservasi dimensi yang cukup.
func (r *UnsizedImageRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if !isImageElement(node) {
		return nil
	}

	// Jangan beri alarm palsu pada JSX spread attributes ({...props})
	if hasSpreadProps(node.Attributes) {
		return nil
	}

	// 1. Pengecualian Komponen Astro <Image /> / <Picture /> dengan asset lokal terikat
	if node.Tag == "Image" || node.Tag == "Picture" {
		if src, ok := node.GetAttr("src"); ok {
			// Jika src berupa ekspresi variabel/impor lokal Astro (bukan URL absolut remote)
			if !strings.HasPrefix(src, "http://") && !strings.HasPrefix(src, "https://") && !strings.HasPrefix(src, "//") {
				return nil
			}
		}
	}

	// 2. Cek atribut HTML width dan height eksplisit
	if hasDimensionAttributes(node) {
		return nil
	}

	// 3. Cek utilitas Tailwind pada node itu sendiri (aspect-*, size-*, w-* + h-*)
	if hasTailwindDimensions(node) {
		return nil
	}

	// 4. Cek deklarasi inline style (aspect-ratio atau width+height)
	if hasInlineDimensionStyle(node) {
		return nil
	}

	// 5. Cek pembungkus langsung (parent level 1) dengan aspect ratio atau bounded height
	if hasAncestorDimensionOrAspect(node, 1) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Image element <%s> lacks explicit width/height attributes or an aspect-ratio constraint. Browsers cannot calculate the reserved rendering box before the image loads, risking Cumulative Layout Shift (CLS).", node.Tag),
			Hint:     "Add 'width' and 'height' attributes (e.g. width={1200} height={600}), or use Tailwind 'aspect-video' / 'aspect-square' / 'w-* h-*' to reserve vertical layout space.",
		},
	}
}
