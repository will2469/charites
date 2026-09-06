package lcp

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// LegacyCriticalFontResourceRule mendeteksi deklarasi @font-face yang hanya merujuk pada format
// font lawas tak terkompresi (.ttf, .otf, .eot) atau tidak memprioritaskan format modern WOFF2.
type LegacyCriticalFontResourceRule struct{}

// NewLegacyCriticalFontResourceRule membuat instance baru dari LegacyCriticalFontResourceRule.
func NewLegacyCriticalFontResourceRule() *LegacyCriticalFontResourceRule {
	return &LegacyCriticalFontResourceRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *LegacyCriticalFontResourceRule) ID() string {
	return "lcp.legacy-critical-font-resource"
}

// Description mengembalikan ringkasan aturan.
func (r *LegacyCriticalFontResourceRule) Description() string {
	return "Custom '@font-face' declaration provides only legacy uncompressed font formats (.ttf, .otf, .eot) or deprioritizes WOFF2 in 'src:', inflating byte transfer payload"
}

// Category mengembalikan nama kategori rule.
func (r *LegacyCriticalFontResourceRule) Category() string {
	return "lcp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *LegacyCriticalFontResourceRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *LegacyCriticalFontResourceRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Google Chrome Core Web Vitals (Largest Contentful Paint Resource Load Duration)",
			"W3C WOFF File Format 2.0 (WOFF2) Recommendation",
			"IETF Brotli Compressed Data Format Specification",
		},
		CoreInvariant: "Custom @font-face declarations for web fonts must specify the modern WOFF2 format as the first item in the 'src' descriptor to maximize compression efficiency.",
		Grounding: "Legacy font formats such as raw TrueType (.ttf), OpenType (.otf), and Embedded OpenType (.eot) lack modern compression algorithms, resulting in file sizes ranging from 200KB to 800KB per font weight.\n\n" +
			"WOFF2 utilizes the Brotli compression algorithm, reducing font binary size by 50% to 80% compared to TTF/OTF and approximately 30% compared to WOFF 1.0 without loss of font hinting or OpenType layout features.\n\n" +
			"Browsers evaluate 'src' declarations in sequential order. Declaring WOFF2 first guarantees that modern browsers download the most compressed variant, accelerating the Resource Load Duration of LCP text elements.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Massive Font Transfer Payload",
				Severity: "HIGH",
				Impact:   "Downloading uncompressed 500KB+ TTF/OTF font files inflates Resource Load Duration on mobile networks.",
			},
			{
				Vector:   "Bandwidth Competition with Hero Media",
				Severity: "MEDIUM",
				Impact:   "Bulky font files compete for socket bandwidth against hero images and critical CSS stylesheets during early page loading.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "css",
				Comment:  "Font declaration only provides raw uncompressed TTF format",
				Code: `@font-face {
  font-family: 'HeadingDisplay';
  src: url('/fonts/heading.ttf') format('truetype');
  font-display: swap;
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "css",
				Comment:  "WOFF2 declared as primary format with progressive TTF fallback",
				Code: `@font-face {
  font-family: 'HeadingDisplay';
  src: url('/fonts/heading.woff2') format('woff2'),
       url('/fonts/heading.ttf') format('truetype');
  font-display: swap;
}`,
			},
		},
	}
}

// Evaluate memeriksa apakah blok @font-face menggunakan format lawas tanpa WOFF2 sebagai prioritas pertama.
func (r *LegacyCriticalFontResourceRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || !strings.EqualFold(node.Tag, "style") {
		return nil
	}

	cssText := getStyleNodeText(node)
	if !strings.Contains(cssText, "@font-face") {
		return nil
	}

	blocks := extractFontFaceBlocks(cssText)
	if len(blocks) == 0 {
		return nil
	}

	var diags []ir.Diagnostic
	for _, block := range blocks {
		family := extractFontFamilyName(block)

		// Pengecualian: Font ikon
		if isIconFontFamily(family) {
			continue
		}

		// Pengecualian: Font sistem lokal
		if isLocalOnlyFontFace(block) {
			continue
		}

		// Pengecualian: Komentar ignore Charites
		if strings.Contains(block, "charites:ignore") {
			continue
		}

		if legacyFmt, violated := hasLegacyFontFormatViolation(block); violated {
			diags = append(diags, ir.Diagnostic{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Custom '@font-face' declaration for '%s' uses legacy uncompressed font formats ('%s') without WOFF2 as first priority. WOFF2 with Brotli compression reduces font payload by 50-80%%, accelerating LCP text block paint.", family, legacyFmt),
				Hint:     "Provide WOFF2 as the first format in the 'src' declaration (e.g. 'url(\"font.woff2\") format(\"woff2\")') with legacy formats as progressive fallbacks.",
			})
		}
	}

	return diags
}
