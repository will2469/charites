package theme

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// PrimitiveInComponentRule mendeteksi konsumsi langsung palet primitif Tailwind (misal: bg-blue-600, text-slate-800)
// di dalam komponen UI, yang melanggar hierarki token 3-Tier W3C DTCG.
type PrimitiveInComponentRule struct{}

// NewPrimitiveInComponentRule membuat instance baru PrimitiveInComponentRule.
func NewPrimitiveInComponentRule() *PrimitiveInComponentRule {
	return &PrimitiveInComponentRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *PrimitiveInComponentRule) ID() string {
	return "theme.primitive-in-component"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *PrimitiveInComponentRule) Description() string {
	return "Detects direct usage of Tailwind primitive palette colors in component classes instead of semantic tokens"
}

// Category mengembalikan nama kategori rule.
func (r *PrimitiveInComponentRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *PrimitiveInComponentRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki otomatis.
func (r *PrimitiveInComponentRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Design Tokens Community Group (DTCG) 3-Tier Architecture",
			"Tailwind CSS Design Token Architecture",
		},
		CoreInvariant: "UI components must consume Tier 2 Semantic Tokens (e.g. bg-primary, text-muted-foreground), never Tier 1 Primitive Palette tokens directly.",
		Grounding: "The W3C Design Tokens Community Group establishes a 3-tier hierarchy:\n\n" +
			"1. Tier 1 (Primitive/Base): Raw palette colors (blue-600, slate-800) defining available color DNA.\n" +
			"2. Tier 2 (Semantic/Alias): Role-based intents (primary, destructive, card, muted) that map differently across themes.\n" +
			"3. Tier 3 (Component-Specific): Optional scoped overrides.\n\n" +
			"When components consume Tier 1 colors directly:\n" +
			"- Dark mode parity breaks because blue-600 has no semantic relationship to surface contrast.\n" +
			"- Multi-tenant white-labeling is impossible without modifying every template.\n" +
			"- Intent is lost: a developer cannot tell if blue-600 represents an interactive action, info state, or brand accent.",
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Direct primitive colors in button",
				Code:     `<button class="bg-blue-600 hover:bg-blue-700 text-white">Submit</button>`,
			},
			{
				Language: "tsx",
				Comment:  "Primitive text and border colors in card",
				Code: `export function Card() {
  return <div className="border-gray-200 text-slate-800 bg-zinc-50">Content</div>;
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Semantic tokens mapped from global.css",
				Code:     `<button class="bg-primary hover:bg-primary/90 text-primary-foreground">Submit</button>`,
			},
			{
				Language: "tsx",
				Comment:  "Semantic tokens for theme consistency",
				Code: `export function Card() {
  return <div className="border-border text-card-foreground bg-card">Content</div>;
}`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Broken Dark Mode",
				Severity: "HIGH",
				Impact:   "Components with hardcoded primitive colors fail to invert or adapt when switching between light and dark modes.",
			},
			{
				Vector:   "Architectural Decay",
				Severity: "HIGH",
				Impact:   "Violating DTCG token layering forces ad-hoc overrides, leading to widespread design debt.",
			},
		},
	}
}

// extractPrimitiveColorFromToken mengekstrak nama warna primitif jika token menggunakan warna primitif Tailwind.
func extractPrimitiveColorFromToken(class string) (string, bool) {
	_, base := StripVariants(class)
	baseNoAlpha, _, _ := SplitAlphaModifier(base)

	_, remainder, ok := SplitColorPrefix(baseNoAlpha)
	if ok && IsTailwindPrimitiveColor(remainder) {
		return remainder, true
	}

	// Periksa arbitrary CSS variable yang membungkus primitif: var(--blue-500)
	if strings.Contains(baseNoAlpha, "var(--") {
		start := strings.Index(baseNoAlpha, "var(--") + 6
		end := strings.IndexAny(baseNoAlpha[start:], "),")
		if end != -1 {
			varName := baseNoAlpha[start : start+end]
			if IsTailwindPrimitiveColor(varName) {
				return varName, true
			}
		}
	}

	return "", false
}

// Evaluate mengevaluasi node IR dan mendeteksi kelas utilitas yang mengonsumsi palet primitif.
func (r *PrimitiveInComponentRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || len(node.Classes) == 0 {
		return nil
	}

	var diags []ir.Diagnostic
	for _, class := range node.Classes {
		if primColor, found := extractPrimitiveColorFromToken(class); found {
			diags = append(diags, ir.Diagnostic{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Direct primitive color usage in component: \"" + class + "\"",
				Hint:     "Replace primitive color \"" + primColor + "\" with a semantic token (e.g. bg-primary, text-foreground, border-border).",
			})
		}
	}

	return diags
}
