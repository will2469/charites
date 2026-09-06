package lcp

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// UndiscoverableLCPImageRule mendeteksi kontainer hero pelipatan atas yang memuat citra visual utama
// melalui CSS background-image atau utilitas Tailwind bg-[url(...)] tanpa disertai tag <link rel="preload">
// di dalam elemen <head> dokumen.
type UndiscoverableLCPImageRule struct{}

// NewUndiscoverableLCPImageRule membuat instance baru dari UndiscoverableLCPImageRule.
func NewUndiscoverableLCPImageRule() *UndiscoverableLCPImageRule {
	return &UndiscoverableLCPImageRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *UndiscoverableLCPImageRule) ID() string {
	return "lcp.undiscoverable-lcp-image"
}

// Description mengembalikan ringkasan aturan.
func (r *UndiscoverableLCPImageRule) Description() string {
	return "Above-the-fold hero container loads primary image via CSS background without <link rel=\"preload\"> in document head"
}

// Category mengembalikan nama kategori rule.
func (r *UndiscoverableLCPImageRule) Category() string {
	return "lcp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *UndiscoverableLCPImageRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *UndiscoverableLCPImageRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Google Chrome Core Web Vitals (Largest Contentful Paint Resource Load Delay)",
			"Chromium Speculative Preload Scanner Discovery Architecture",
			"W3C Preload Specification (<link rel=\"preload\" as=\"image\">)",
		},
		CoreInvariant: "Hero visual assets must not be embedded exclusively via CSS background-image without a corresponding '<link rel=\"preload\">' in '<head>'; CSS-based images are undiscoverable by the HTML preload scanner.",
		Grounding: "When an image is referenced inside CSS ('background-image: url(...)' or Tailwind 'bg-[url(...)]'), the browser's speculative preload scanner cannot discover the asset URL while streaming the HTML.\n\n" +
			"The browser must first download all render-blocking CSS, parse the cascade, and run the style computation step before it even learns that the image URL exists. This creates massive Resource Load Delay for LCP.\n\n" +
			"Migrating the visual background to a native '<img>' element (e.g. with 'absolute inset-0 w-full h-full object-cover -z-10') makes it immediately discoverable in HTML. If CSS background is necessary, injecting '<link rel=\"preload\" as=\"image\" href=\"...\" fetchpriority=\"high\">' into '<head>' bridges the discovery gap.",
		Risks: []ir.RiskItem{
			{
				Vector:   "CSS Cascade Dependency Block",
				Severity: "HIGH",
				Impact:   "Hero image discovery is delayed until external CSS stylesheets are downloaded and parsed, adding 300ms-1000ms to LCP.",
			},
			{
				Vector:   "Speculative Scanner Blindness",
				Severity: "HIGH",
				Impact:   "Lookahead scanner cannot prefetch the LCP asset during initial document TCP streaming.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Hero container using CSS background-image without head preload",
				Code: `<header className="w-full h-[480px] bg-[url('/hero.webp')] bg-cover" data-perf-role="hero">
  <h1 className="text-white">Galactic Exploration</h1>
</header>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Native <img> with object-cover immediately discoverable by preload scanner",
				Code: `<header className="relative w-full h-[480px] overflow-hidden" data-perf-role="hero">
  <img src="/hero.webp" alt="Hero Background" fetchpriority="high" className="absolute inset-0 w-full h-full object-cover -z-10" />
  <h1 className="relative z-10 text-white p-8">Galactic Exploration</h1>
</header>`,
			},
		},
	}
}

// Evaluate memeriksa apakah kontainer hero memuat gambar via CSS background tanpa preload di head.
func (r *UndiscoverableLCPImageRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !isCandidateHeroContainer(node) {
		return nil
	}

	url, found := hasCSSBackgroundImage(node.RawClasses, node.Attributes)
	if !found {
		return nil
	}

	// Periksa apakah sudah didampingi oleh <link rel="preload" as="image"> di head
	if hasPreloadInHead(node, url) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Hero container '<%s>' embeds primary image '%s' via CSS background without a companion '<link rel=\"preload\">' in '<head>'. The browser speculative preload scanner cannot discover CSS background assets until stylesheet parsing finishes.", node.Tag, url),
			Hint:     "Migrate to a native '<img>' with 'object-cover -z-10' or inject '<link rel=\"preload\" as=\"image\" href=\"...\" fetchpriority=\"high\">' into '<head>'.",
		},
	}
}
