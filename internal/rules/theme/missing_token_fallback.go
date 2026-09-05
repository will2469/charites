package theme

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// MissingTokenFallbackRule mendeteksi penggunaan CSS custom property var(--token)
// tanpa nilai fallback cadangan (argumen kedua) di dalam markup, kelas arbitrer, atau stylesheet.
type MissingTokenFallbackRule struct{}

// NewMissingTokenFallbackRule membuat instance baru MissingTokenFallbackRule.
func NewMissingTokenFallbackRule() *MissingTokenFallbackRule {
	return &MissingTokenFallbackRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *MissingTokenFallbackRule) ID() string {
	return "theme.missing-token-fallback"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *MissingTokenFallbackRule) Description() string {
	return "Detects CSS variable references without fallback values"
}

// Category mengembalikan nama kategori rule.
func (r *MissingTokenFallbackRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *MissingTokenFallbackRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki otomatis.
func (r *MissingTokenFallbackRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C CSS Custom Properties for Cascading Variables Module Level 1",
			"WCAG 2.2 Guideline 4.1 Compatible (Robust Graceful Degradation)",
		},
		CoreInvariant: "CSS variable references in production code must supply a safe fallback value to guard against unresolved design tokens.",
		Grounding: "CSS variables evaluated via var(--name) without a fallback revert to the CSS specification's " +
			"'guaranteed-invalid value' when undefined or failing to load.\n\n" +
			"When developers write color: var(--text-brand) or bg-[var(--brand)] without a fallback:\n" +
			"1. Broken Visual Contrast: Elements render completely transparent or default black, failing WCAG AA contrast.\n" +
			"2. Unhandled CDN / Token Latency: If design tokens load asynchronously or via isolated packages, missing fallbacks cause flash of broken unstyled content (FOBUC).\n" +
			"3. Graceful Degradation Failure: Micro-frontends or embedded widgets fail without host variable injection.\n\n" +
			"Charites recommends always supplying a fallback argument: var(--name, fallback-value).",
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Missing fallback in arbitrary Tailwind utility class",
				Code:     `<div class="bg-[var(--brand)] text-[var(--text-color)]">Unsafe Variable</div>`,
			},
			{
				Language: "tsx",
				Comment:  "Missing fallback in inline style attribute",
				Code: `export function Card() {
  return <div style={{ color: "var(--brand-primary)" }}>Missing Fallback</div>;
}`,
			},
			{
				Language: "astro",
				Comment:  "Missing fallback inside style block",
				Code: `<style>
  .badge {
    background-color: var(--accent-color);
  }
</style>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Safe fallback in arbitrary Tailwind utility class",
				Code:     `<div class="bg-[var(--brand,#2563eb)] text-[var(--text-color,currentColor)]">Safe Variable</div>`,
			},
			{
				Language: "tsx",
				Comment:  "Safe fallback in inline style attribute",
				Code: `export function Card() {
  return <div style={{ color: "var(--brand-primary, #1e293b)" }}>Safe Fallback</div>;
}`,
			},
			{
				Language: "astro",
				Comment:  "Safe fallback inside style block",
				Code: `<style>
  .badge {
    background-color: var(--accent-color, #f59e0b);
  }
</style>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Guaranteed-Invalid Property Rendering",
				Severity: "MEDIUM",
				Impact:   "Missing tokens evaluate to transparent/initial CSS values, causing catastrophic unreadable contrast.",
			},
			{
				Vector:   "Micro-frontend Style Decoupling",
				Severity: "LOW",
				Impact:   "Components embedded in foreign hosts break when global tokens are not shared.",
			},
		},
	}
}

// Evaluate mengevaluasi node IR untuk mencari kemunculan var(--token) tanpa fallback.
func (r *MissingTokenFallbackRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil {
		return nil
	}

	// 1. Evaluasi pada atribut class (arbitrary class utility seperti bg-[var(--brand)])
	for _, cls := range node.Classes {
		if strings.Contains(cls, "var(--") {
			if missingVar, found := findMissingFallbackVar(cls); found {
				return []ir.Diagnostic{
					{
						Line:     node.Span.Line,
						Column:   node.Span.Column,
						Rule:     r.ID(),
						Severity: r.DefaultSeverity(),
						Message:  "CSS variable reference without fallback: \"" + missingVar + "\"",
						Hint:     "Provide a fallback value as the second argument: var(--name, <fallback-value>) to prevent broken styling if the token is undefined.",
					},
				}
			}
		}
	}

	// 2. Evaluasi pada atribut style
	if styleVal, ok := node.GetAttr("style"); ok && strings.Contains(styleVal, "var(--") {
		if missingVar, found := findMissingFallbackVar(styleVal); found {
			return []ir.Diagnostic{
				{
					Line:     node.Span.Line,
					Column:   node.Span.Column,
					Rule:     r.ID(),
					Severity: r.DefaultSeverity(),
					Message:  "CSS variable reference without fallback: \"" + missingVar + "\"",
					Hint:     "Provide a fallback value as the second argument: var(--name, <fallback-value>) to prevent broken styling if the token is undefined.",
				},
			}
		}
	}

	// 3. Evaluasi pada tag <style>
	if node.Tag == "style" {
		cssText := getStyleNodeText(node)
		if strings.Contains(cssText, "var(--") {
			cleaned := stripCSSCommentsString(cssText)
			if missingVar, found := findMissingFallbackVar(cleaned); found {
				return []ir.Diagnostic{
					{
						Line:     node.Span.Line,
						Column:   node.Span.Column,
						Rule:     r.ID(),
						Severity: r.DefaultSeverity(),
						Message:  "CSS variable reference without fallback: \"" + missingVar + "\"",
						Hint:     "Provide a fallback value as the second argument: var(--name, <fallback-value>) to prevent broken styling if the token is undefined.",
					},
				}
			}
		}
	}

	return nil
}

// findMissingFallbackVar mencari ekspresi var(--foo) yang tidak memiliki koma pemisah fallback.
func findMissingFallbackVar(s string) (string, bool) {
	idx := 0
	n := len(s)

	for idx < n {
		pos := strings.Index(s[idx:], "var(--")
		if pos == -1 {
			break
		}

		start := idx + pos
		end, hasComma := scanVarCall(s, start+len("var(--"))
		if end != -1 && !hasComma {
			return s[start : end+1], true
		}

		if end != -1 {
			idx = end + 1
		} else {
			idx = start + len("var(--")
		}
	}

	return "", false
}

// scanVarCall memindai kurung tutup penyeimbang dan mendeteksi koma pemisah fallback pada tingkat terluar.
func scanVarCall(s string, contentStart int) (int, bool) {
	parenDepth := 1
	hasComma := false

	for i := contentStart; i < len(s); i++ {
		switch {
		case s[i] == '(':
			parenDepth++
		case s[i] == ')':
			parenDepth--
			if parenDepth == 0 {
				return i, hasComma
			}
		case s[i] == ',' && parenDepth == 1:
			hasComma = true
		}
	}

	return -1, false
}
