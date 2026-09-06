package lcp

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// BlockedCriticalFontRule mendeteksi deklarasi @font-face untuk jenis huruf kustom
// yang tidak menyertakan deskriptor font-display: swap atau font-display: optional.
type BlockedCriticalFontRule struct{}

// NewBlockedCriticalFontRule membuat instance baru dari BlockedCriticalFontRule.
func NewBlockedCriticalFontRule() *BlockedCriticalFontRule {
	return &BlockedCriticalFontRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *BlockedCriticalFontRule) ID() string {
	return "lcp.blocked-critical-font"
}

// Description mengembalikan ringkasan aturan.
func (r *BlockedCriticalFontRule) Description() string {
	return "Custom '@font-face' declaration lacks 'font-display: swap' or 'font-display: optional', risking FOIT (Flash of Invisible Text) and delaying LCP text paint"
}

// Category mengembalikan nama kategori rule.
func (r *BlockedCriticalFontRule) Category() string {
	return "lcp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *BlockedCriticalFontRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *BlockedCriticalFontRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Google Chrome Core Web Vitals (Largest Contentful Paint Text Block Paint)",
			"W3C CSS Fonts Module Level 4 (font-display descriptor specification)",
			"Web Performance Working Group FOIT Minimization Guidelines",
		},
		CoreInvariant: "Custom @font-face declarations for web fonts must specify 'font-display: swap' or 'font-display: optional' to prevent Flash of Invisible Text (FOIT) on LCP text blocks.",
		Grounding: "When a browser discovers text styled with a custom web font, it evaluates the @font-face 'font-display' descriptor.\n\n" +
			"By default ('font-display: auto' or 'font-display: block'), modern browsers enter a 3-second 'block period' during which text is rendered with invisible transparent glyphs while the font binary is fetched from the network.\n\n" +
			"If the primary heading (<h1> or hero banner text) is the LCP candidate, this block period directly delays the Largest Contentful Paint until font download completes.\n\n" +
			"Specifying 'font-display: swap' enables immediate text rendering with a system fallback font followed by an in-place swap once the custom font arrives.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Flash of Invisible Text (FOIT)",
				Severity: "HIGH",
				Impact:   "LCP candidate text remains completely invisible for up to 3000ms on cellular or high-latency networks.",
			},
			{
				Vector:   "Element Render Delay Inflation",
				Severity: "MEDIUM",
				Impact:   "Directly inflates LCP duration by coupling text paint to third-party or remote font network latency.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "css",
				Comment:  "Custom @font-face without font-display causing FOIT on hero headings",
				Code: `@font-face {
  font-family: 'CabinetGrotesk';
  src: url('/fonts/cabinet.woff2') format('woff2');
}
h1 {
  font-family: 'CabinetGrotesk', sans-serif;
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "css",
				Comment:  "Custom @font-face configured with font-display: swap to ensure immediate text rendering",
				Code: `@font-face {
  font-family: 'CabinetGrotesk';
  src: url('/fonts/cabinet.woff2') format('woff2');
  font-display: swap;
}
h1 {
  font-family: 'CabinetGrotesk', sans-serif;
}`,
			},
		},
	}
}

// Evaluate memeriksa apakah blok @font-face tidak memiliki font-display swap atau optional.
func (r *BlockedCriticalFontRule) Evaluate(node *ir.Node) []ir.Diagnostic {
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

		// Pengecualian: Font ikon sengaja menggunakan font-display block untuk mencegah character jitter
		if isIconFontFamily(family) {
			continue
		}

		// Pengecualian: Font sistem lokal murni tanpa unduhan font web
		if isLocalOnlyFontFace(block) {
			continue
		}

		// Pengecualian: Komentar ignore Charites
		if strings.Contains(block, "charites:ignore") {
			continue
		}

		if !hasValidSwapFontDisplay(block) {
			diags = append(diags, ir.Diagnostic{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Custom '@font-face' declaration for '%s' lacks 'font-display: swap' or 'font-display: optional'. Missing or block font-display causes FOIT (Flash of Invisible Text), blocking LCP text rendering during font network download.", family),
				Hint:     "Add 'font-display: swap;' to the '@font-face' declaration to guarantee immediate text visibility using a fallback font.",
			})
		}
	}

	return diags
}
