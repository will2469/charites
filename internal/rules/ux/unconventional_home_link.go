package ux

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// UnconventionalHomeLinkRule menegakkan Hukum Jakob pada header navigasi utama:
// logo/identitas brand di dalam header wajib dibungkus oleh tautan yang mengarah ke beranda root ("/").
type UnconventionalHomeLinkRule struct{}

// NewUnconventionalHomeLinkRule membuat instance baru dari UnconventionalHomeLinkRule.
func NewUnconventionalHomeLinkRule() *UnconventionalHomeLinkRule {
	return &UnconventionalHomeLinkRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *UnconventionalHomeLinkRule) ID() string {
	return "ux.unconventional-home-link"
}

// Description mengembalikan ringkasan aturan.
func (r *UnconventionalHomeLinkRule) Description() string {
	return "Enforces Jakob's Law by ensuring header logo/brand identity links to the root home page ('/')"
}

// Category mengembalikan nama kategori rule.
func (r *UnconventionalHomeLinkRule) Category() string {
	return "ux"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *UnconventionalHomeLinkRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *UnconventionalHomeLinkRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Jakob's Law of Internet User Experience (Nielsen Norman Group)",
			"W3C Web Navigation & Landmark Architecture Guidelines",
			"ISO 9241-110 Ergonomics of Human-System Interaction (Suitability for Learning & Predictability)",
		},
		CoreInvariant: "Brand identity and logo elements in the primary header must be enclosed within an anchor ('<a>' or '<Link>') whose destination normalizes to the root homepage ('/').",
		Grounding: "Jakob's Law states that users spend most of their time on sites other than yours. " +
			"Consequently, they bring deeply ingrained mental models about standard interaction patterns. " +
			"The most universal web convention is that clicking the brand logo in the top-left header returns to the homepage ('/').\n\n" +
			"When a logo is unclickable, rendered as a passive image or plain text, or links to an unexpected secondary destination (like /about, /products, or an external portal), " +
			"users become disoriented and lose their primary visual escape hatch back to the system root.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Mental Model Disorientation",
				Severity: "MEDIUM",
				Impact:   "Users habitually click the top-left brand mark when seeking the homepage; non-functional or diverted logos induce frustration and cognitive friction.",
			},
			{
				Vector:   "Accidental Site Exit or Dead End",
				Severity: "LOW",
				Impact:   "Navigating away from the root application when attempting to reset context forces users to rely on browser history or address bar edits.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Passive brand logo in header without any enclosing link",
				Code: `<header className="flex items-center justify-between px-6 py-4 border-b">
  <img src="/brand-logo.svg" alt="Acme Corporation Logo" className="h-8 w-auto" />
  <nav className="flex gap-4">
    <a href="/features">Features</a>
    <a href="/pricing">Pricing</a>
  </nav>
</header>`,
			},
			{
				Language: "astro",
				Comment:  "Brand logo linking to an internal sub-page instead of root",
				Code: `<header class="flex items-center justify-between px-6 py-4">
  <a href="/about" class="brand-logo flex items-center gap-2">
    <img src="/logo.svg" alt="Brand Logo" />
    <span class="font-bold">Portal</span>
  </a>
</header>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Brand logo wrapped in accessible anchor linking directly to root '/'",
				Code: `<header className="flex items-center justify-between px-6 py-4 border-b">
  <a href="/" aria-label="Acme Corporation - Beranda" className="flex items-center gap-2">
    <img src="/brand-logo.svg" alt="Acme Corporation Logo" className="h-8 w-auto" />
    <span className="font-bold text-lg">Acme</span>
  </a>
  <nav className="flex gap-4">
    <a href="/features">Features</a>
  </nav>
</header>`,
			},
		},
	}
}

// Evaluate memeriksa apakah logo/identitas brand di dalam header memiliki tautan ke root.
func (r *UnconventionalHomeLinkRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	// 1. Periksa apakah node adalah tautan pembungkus brand di dalam header
	if node.Tag == "a" || stringsHasSuffixLink(node.Tag) {
		if !isInsideHeaderBanner(node) {
			return nil
		}

		if anchorContainsBrandIdentity(node) {
			href, _ := getAttrCaseInsensitive(node, "href", "to")
			if !isNormalizedRootHref(href) {
				cleanTarget := cleanAttrValue(href)
				if cleanTarget == "" {
					cleanTarget = "(empty)"
				}
				return []ir.Diagnostic{
					{
						Line:     node.Span.Line,
						Column:   node.Span.Column,
						Rule:     r.ID(),
						Severity: r.DefaultSeverity(),
						Message: fmt.Sprintf(
							"Header brand logo links to %q instead of the root home page ('/'). Jakob's Law requires primary brand anchors to navigate to the homepage.",
							cleanTarget,
						),
						Hint: "Update the link href to point directly to '/' (or the localized homepage root).",
					},
				}
			}
		}
		return nil
	}

	// 2. Periksa apakah elemen adalah brand identity tanpa tautan pembungkus di dalam header
	if isBrandIdentityElement(node) {
		if !isInsideHeaderBanner(node) {
			return nil
		}

		// Jika sudah memiliki tautan pembungkus, evaluasi tautan dilakukan pada cabang di atas
		if findEnclosingHeaderLink(node) != nil {
			return nil
		}

		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Brand logo/identity in header is not enclosed in a navigation link to the root home page ('/'). Users expect clicking the logo to return to the homepage (Jakob's Law).",
				Hint:     "Wrap the brand logo element in an <a href=\"/\"> anchor with an accessible aria-label.",
			},
		}
	}

	return nil
}

func stringsHasSuffixLink(tag string) bool {
	return len(tag) >= 4 && (tag[len(tag)-4:] == "Link" || tag[len(tag)-4:] == "link")
}

func anchorContainsBrandIdentity(node *ir.Node) bool {
	if isBrandIdentityElement(node) {
		return true
	}
	for child := range node.Walk() {
		if child == node {
			continue
		}
		if isBrandIdentityElement(child) {
			return true
		}
	}
	return false
}
