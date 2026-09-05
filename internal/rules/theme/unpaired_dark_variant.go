package theme

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// UnpairedDarkVariantRule mendeteksi deklarasi varian dark mode sepihak yang menyebabkan
// tabrakan kontras ekstrem (black-on-black atau white-on-white).
type UnpairedDarkVariantRule struct{}

// NewUnpairedDarkVariantRule membuat instance baru UnpairedDarkVariantRule.
func NewUnpairedDarkVariantRule() *UnpairedDarkVariantRule {
	return &UnpairedDarkVariantRule{}
}

// ID mengembalikan Charites Rule ID kanonikal.
func (r *UnpairedDarkVariantRule) ID() string {
	return "theme.unpaired-dark-variant"
}

// Description mengembalikan penjelasan ringkas rule.
func (r *UnpairedDarkVariantRule) Description() string {
	return "Detects one-sided dark theme variant declarations causing severe contrast collisions"
}

// Category mengembalikan nama kategori rule.
func (r *UnpairedDarkVariantRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *UnpairedDarkVariantRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *UnpairedDarkVariantRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"WCAG 2.2 Success Criterion 1.4.3 (Contrast Minimum)",
			"W3C Design Tokens Community Group (DTCG)",
			"Tailwind CSS Dark Mode Variant Architecture",
		},
		CoreInvariant: "Background and text theme variants must be paired symmetrically, or use adaptive semantic tokens (bg-card, text-card-foreground) to guarantee contrast.",
		Grounding: "Declaring one-sided dark mode classes (such as dark:bg-zinc-900 without a light base background, or inverting container backgrounds without adapting child text colors) causes catastrophic contrast collapses:\n\n" +
			"1. Black-on-Black Illegibility: An element that inverts to dark:bg-zinc-900 while child text remains text-zinc-900 renders completely unreadable text in dark mode.\n" +
			"2. Incomplete State Inversion: Specifying dark:bg-* without a default bg-* causes unpredictable transparency blending over parent containers.\n" +
			"3. Accessibility Failures: Contrast ratios plummet below 1.5:1, triggering immediate WCAG Level AA and AAA violations.\n\n" +
			"Charites enforces symmetric pairing (e.g. bg-white dark:bg-zinc-900 with text-zinc-900 dark:text-zinc-100) or using theme-adaptive semantic tokens (bg-card text-card-foreground).",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Container inverts background but child text remains dark-mode blind",
				Code: `<div className="bg-white dark:bg-zinc-900">
  <span className="text-zinc-900">Title</span>
</div>`,
			},
			{
				Language: "astro",
				Comment:  "Unpaired dark background variant without base background",
				Code:     `<div class="dark:bg-zinc-900"><span>Content</span></div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Using semantic tokens that adapt automatically across themes",
				Code: `<div className="bg-card text-card-foreground">
  <span>Title</span>
</div>`,
			},
			{
				Language: "astro",
				Comment:  "Symmetrically paired background and text variants",
				Code: `<div class="bg-white dark:bg-zinc-900">
  <span class="text-zinc-900 dark:text-zinc-100">Title</span>
</div>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Contrast Collapse (Black-on-Black / White-on-White)",
				Severity: "HIGH",
				Impact:   "Users are unable to read text or interact with controls when switching theme modes.",
			},
			{
				Vector:   "Theme State Fragmentation",
				Severity: "MEDIUM",
				Impact:   "Unpaired utility modifiers lead to unpredictable cascading color bugs across nested layouts.",
			},
		},
	}
}

// Evaluate mengevaluasi sebuah node IR untuk mendeteksi varian tema dark yang tidak berpasangan.
func (r *UnpairedDarkVariantRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || len(node.Classes) == 0 {
		return nil
	}

	var diags []ir.Diagnostic
	var darkBgClass string
	hasBaseBg := false
	hasDarkText := false
	var hardcodedTextClass string

	for _, class := range node.Classes {
		if isDarkBgClass(class) {
			darkBgClass = class
		}
		if isBaseBgClass(class) {
			hasBaseBg = true
		}
		if isDarkTextClass(class) {
			hasDarkText = true
		}
		if isHardcodedTextClass(class) {
			hardcodedTextClass = class
		}
	}

	// 1. Cek dark:bg-* tanpa base bg-*
	if darkBgClass != "" && !hasBaseBg {
		diags = append(diags, ir.Diagnostic{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Unpaired dark background variant: \"" + darkBgClass + "\" lacks base light background class",
			Hint:     "Add a base background class (e.g. bg-white) or use adaptive semantic token bg-card.",
		})
	}

	// 2. Cek container yang menginversi background (bg-* & dark:bg-*) dengan teks sendiri yang hardcoded tanpa dark:text-*
	if darkBgClass != "" && hasBaseBg && hardcodedTextClass != "" && !hasDarkText {
		diags = append(diags, ir.Diagnostic{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Inverted container background has hardcoded text: \"" + hardcodedTextClass + "\" without dark:text-* variant",
			Hint:     "Add a dark:text-* variant or use semantic token text-card-foreground.",
		})
	}

	// 3. Cek container yang menginversi background dengan anak teks yang hardcoded tanpa dark:text-*
	if darkBgClass != "" && hasBaseBg && len(node.Children) > 0 {
		diags = append(diags, r.checkChildTextParity(node)...)
	}

	return diags
}

func (r *UnpairedDarkVariantRule) checkChildTextParity(node *ir.Node) []ir.Diagnostic {
	var diags []ir.Diagnostic
	for _, child := range node.Children {
		if child == nil || len(child.Classes) == 0 {
			continue
		}
		childHasDarkText := false
		var childHardcodedText string
		for _, c := range child.Classes {
			if isDarkTextClass(c) {
				childHasDarkText = true
			}
			if isHardcodedTextClass(c) {
				childHardcodedText = c
			}
		}
		if childHardcodedText != "" && !childHasDarkText {
			diags = append(diags, ir.Diagnostic{
				Line:     child.Span.Line,
				Column:   child.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Child text \"" + childHardcodedText + "\" in inverted container lacks dark:text-* variant",
				Hint:     "Add a matching dark:text-* variant or use semantic token text-card-foreground.",
			})
		}
	}
	return diags
}

func isDarkBgClass(class string) bool {
	if strings.HasPrefix(class, "dark:bg-") {
		rem := class[8:]
		return rem != "transparent" && rem != "none"
	}
	return false
}

func isBaseBgClass(class string) bool {
	if strings.HasPrefix(class, "dark:") || strings.Contains(class, ":dark:") {
		return false
	}
	base := StripVariantsOnlyBase(class)
	return strings.HasPrefix(base, "bg-") && base != "bg-transparent" && base != "bg-none"
}

func isDarkTextClass(class string) bool {
	return strings.HasPrefix(class, "dark:text-") || strings.Contains(class, ":dark:text-")
}

func isHardcodedTextClass(class string) bool {
	if strings.HasPrefix(class, "dark:") || strings.Contains(class, ":dark:") {
		return false
	}
	base := StripVariantsOnlyBase(class)
	if !strings.HasPrefix(base, "text-") {
		return false
	}
	rem := base[5:]
	baseColor, _, _ := SplitAlphaModifier(rem)
	if IsMonochromeColor(baseColor) || IsTailwindPrimitiveColor(baseColor) || IsHexColor(baseColor) {
		return true
	}
	if strings.HasPrefix(rem, "[#") || strings.HasPrefix(rem, "[rgb") {
		return true
	}
	return false
}
