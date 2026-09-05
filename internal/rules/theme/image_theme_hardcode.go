package theme

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// ImageThemeHardcodeRule mendeteksi penggunaan tag <img> yang merujuk pada aset grafis/logo
// tanpa mekanisme adaptasi tema (seperti varian dark:hidden / dark:block, filter dark:invert, atau elemen <picture>).
type ImageThemeHardcodeRule struct{}

// NewImageThemeHardcodeRule membuat instance baru ImageThemeHardcodeRule.
func NewImageThemeHardcodeRule() *ImageThemeHardcodeRule {
	return &ImageThemeHardcodeRule{}
}

// ID mengembalikan Charites Rule ID kanonikal.
func (r *ImageThemeHardcodeRule) ID() string {
	return "theme.image-theme-hardcode"
}

// Description mengembalikan penjelasan ringkas rule.
func (r *ImageThemeHardcodeRule) Description() string {
	return "Detects graphic assets and logos in img tags lacking dark mode theme adaptation"
}

// Category mengembalikan nama kategori rule.
func (r *ImageThemeHardcodeRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *ImageThemeHardcodeRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ImageThemeHardcodeRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"WCAG 2.2 Success Criterion 1.4.11 (Non-text Contrast)",
			"W3C Responsive Images & Art Direction Specification",
			"Tailwind CSS Dark Mode Graphic Switching Guidelines",
		},
		CoreInvariant: "Graphic assets, logos, and diagrams in img tags must provide theme-adaptive variants via picture, dark: utility classes, or invert filters.",
		Grounding: "Embedding graphical assets (such as brand logos, SVG diagrams, and charts) via static <img> tags without dark mode adaptation triggers severe visual breakage:\n\n" +
			"1. Asset Invisibility: A dark or black logo rendered against a dark mode background becomes completely invisible.\n" +
			"2. Inverted Eye-Strain: High-glare white background diagrams blast excessive light on dark UI themes.\n" +
			"3. Inflexible Art Direction: Projects without responsive image pairing cannot tailor vector artwork to dark contrast requirements.\n\n" +
			"Charites enforces theme-aware graphic handling using dark:hidden / dark:block class pairs, dark:invert filters, or responsive <picture> elements.",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Static logo in img tag without dark mode variant in TSX",
				Code:     `<img src="/images/logo-black.svg" alt="Company Logo" />`,
			},
			{
				Language: "astro",
				Comment:  "Vector architecture diagram without theme switching in Astro",
				Code:     `<img src="/assets/diagram.svg" alt="Architecture Flow" />`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Theme-paired image switching using Tailwind dark utilities",
				Code: `<>
  <img src="/images/logo-light.svg" className="dark:hidden" alt="Logo" />
  <img src="/images/logo-dark.svg" className="hidden dark:block" alt="Logo" />
</>`,
			},
			{
				Language: "astro",
				Comment:  "Using dark:invert filter for vector diagrams",
				Code:     `<img src="/assets/diagram.svg" class="dark:invert" alt="Architecture Flow" />`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Asset Disappearance in Dark Themes",
				Severity: "MEDIUM",
				Impact:   "Brand logos, technical diagrams, and icon artwork become illegible on dark surfaces.",
			},
			{
				Vector:   "Non-text Contrast Failure (WCAG 1.4.11)",
				Severity: "MEDIUM",
				Impact:   "Visual cues necessary for interface understanding fail accessibility contrast requirements.",
			},
		},
	}
}

// Evaluate mengevaluasi sebuah node IR untuk mendeteksi tag <img> grafis tanpa adaptasi tema.
func (r *ImageThemeHardcodeRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || strings.ToLower(node.Tag) != "img" {
		return nil
	}

	src := node.Attributes["src"]
	if src == "" || !IsThemeGraphicAsset(src) {
		return nil
	}

	// 1. Periksa apakah berada di dalam elemen <picture> (art direction adaptif)
	for p := node.Parent; p != nil; p = p.Parent {
		if strings.ToLower(p.Tag) == "picture" {
			return nil
		}
	}

	// 2. Periksa apakah memiliki varian dark: (dark:hidden, dark:block, dark:invert, dll.)
	if HasDarkVariant(node.Classes) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Graphic asset in <img> tag lacks dark mode theme adaptation: \"" + src + "\"",
			Hint:     "Pair with dark mode variant (e.g. dark:hidden / dark:block), apply dark:invert, or wrap in <picture>.",
		},
	}
}
