package lcp

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// OversizedLCPResourceSelectionRule mendeteksi elemen gambar kandidat LCP berdimensi fluida
// yang hanya menyediakan atribut src tunggal tanpa atribut srcset dan sizes responsif.
type OversizedLCPResourceSelectionRule struct{}

// NewOversizedLCPResourceSelectionRule membuat instance baru dari OversizedLCPResourceSelectionRule.
func NewOversizedLCPResourceSelectionRule() *OversizedLCPResourceSelectionRule {
	return &OversizedLCPResourceSelectionRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *OversizedLCPResourceSelectionRule) ID() string {
	return "lcp.oversized-lcp-resource-selection"
}

// Description mengembalikan ringkasan aturan.
func (r *OversizedLCPResourceSelectionRule) Description() string {
	return "Fluid responsive LCP candidate image lacks responsive 'srcset' and 'sizes' attributes, forcing mobile viewports to download oversized desktop assets"
}

// Category mengembalikan nama kategori rule.
func (r *OversizedLCPResourceSelectionRule) Category() string {
	return "lcp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *OversizedLCPResourceSelectionRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *OversizedLCPResourceSelectionRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Google Chrome Core Web Vitals (Largest Contentful Paint Resource Load Duration)",
			"Responsive Images Community Group (RICG) Responsive Images Specification",
			"HTML Living Standard srcset and sizes Attributes Specification",
		},
		CoreInvariant: "Fluid responsive LCP candidate images must provide responsive 'srcset' with width descriptors and a 'sizes' attribute to prevent mobile devices from downloading oversized desktop assets.",
		Grounding: "When a fluid image (such as a full-width hero banner) only specifies a single large 'src' attribute, mobile devices with small viewports are forced to download the same high-resolution asset designed for 4K desktop screens.\n\n" +
			"This unnecessary byte payload directly prolongs the Resource Load Duration component of LCP over cellular networks.\n\n" +
			"By providing a 'srcset' attribute with width descriptors (e.g. '400w, 800w, 1200w') alongside a matching 'sizes' attribute (or using the '<Image />' component from 'astro:assets'), the browser can accurately select the optimal image variant for the user's viewport and device pixel ratio.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Resource Load Duration Bloat",
				Severity: "HIGH",
				Impact:   "Mobile devices download 2MB-5MB desktop-resolution assets over cellular connections, adding 500ms-2500ms to LCP.",
			},
			{
				Vector:   "Excess Mobile Data Consumption",
				Severity: "MEDIUM",
				Impact:   "Users on metered cellular data plans consume excessive bandwidth downloading unneeded pixels.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Fluid hero image only provides a single massive desktop asset without srcset and sizes",
				Code: `<section className="hero-section" data-perf-role="hero">
  <img
    src="/images/hero-3840x2160.jpg"
    alt="Hero Banner"
    className="w-full h-auto"
    fetchpriority="high"
  />
</section>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Fluid hero image configured with responsive srcset width descriptors and sizes attribute",
				Code: `<section className="hero-section" data-perf-role="hero">
  <img
    src="/images/hero-1200.webp"
    srcset="/images/hero-400.webp 400w, /images/hero-800.webp 800w, /images/hero-1200.webp 1200w"
    sizes="(max-width: 768px) 100vw, 1200px"
    alt="Hero Banner"
    className="w-full h-auto"
    fetchpriority="high"
  />
</section>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen gambar kandidat LCP berdimensi fluida tidak memiliki atribut srcset dan sizes.
func (r *OversizedLCPResourceSelectionRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !isMediaTag(node.Tag) {
		return nil
	}

	// Komponen <Image /> bawaan framework (Astro/Next) meng-generate srcset otomatis
	if strings.EqualFold(node.Tag, "Image") {
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

	// Format SVG bersifat vektor dan resolusi-independen
	if isSVGFile(rawSrc) {
		return nil
	}

	// Jika dibungkus dalam <picture> dengan <source srcset="...">, browser memilih sumber responsif
	if isInsidePictureWithResponsiveSource(node) {
		return nil
	}

	// Aturan ini HANYA mengaudit elemen berdimensi fluida; elemen dimensi tetap adalah domain image-source-density-mismatch
	if !isFluidImage(node) {
		return nil
	}

	if !hasResponsiveSrcset(node.Attributes) {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Fluid LCP candidate image '<%s>' lacks responsive 'srcset' and 'sizes' attributes. Mobile viewports will download oversized desktop-resolution assets, heavily inflating Resource Load Duration.", node.Tag),
				Hint:     "Provide a responsive 'srcset' with width descriptors (e.g. '400w, 800w, 1200w') and a 'sizes' attribute, or use the '<Image />' component from 'astro:assets'.",
			},
		}
	}

	return nil
}
