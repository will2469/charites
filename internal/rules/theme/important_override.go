package theme

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// ImportantOverrideRule mendeteksi penggunaan prefix ! (important) pada utility class pewarnaan
// (seperti !bg-red-500, !text-white, !bg-primary) yang merusak hierarki spesifisitas dan cascade tema.
type ImportantOverrideRule struct{}

// NewImportantOverrideRule membuat instance baru ImportantOverrideRule.
func NewImportantOverrideRule() *ImportantOverrideRule {
	return &ImportantOverrideRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *ImportantOverrideRule) ID() string {
	return "theme.important-override"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *ImportantOverrideRule) Description() string {
	return "Detects !important modifiers on color utility classes that break theme cascade and specificity hierarchy"
}

// Category mengembalikan nama kategori rule.
func (r *ImportantOverrideRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *ImportantOverrideRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki otomatis.
func (r *ImportantOverrideRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Cascading Style Sheets (CSS) Specificity Level 4",
			"Tailwind CSS Design Token Architecture",
		},
		CoreInvariant: "Color utility classes must never use the !important modifier (!bg-*, !text-*); specificity must be managed via CSS Cascade Layers.",
		Grounding: "Using the ! modifier (e.g. !bg-red-500 or !text-white) forcefully escalates CSS declaration priority above normal cascade layers.\n\n" +
			"1. Destroys Theme Inversion: Dark mode variants (.dark bg-card) cannot override !bg-white without also adding !dark:bg-card, sparking an !important arms race.\n" +
			"2. Compromises Component Reusability: Reusable components with !important color classes cannot be customized or themed by parent containers.\n" +
			"3. Unpredictable State Styling: Hover, focus, and disabled state colors fail to trigger reliably when base colors are marked !important.\n\n" +
			"Charites enforces relying on Cascade Layers (@layer components, utilities) and semantic token definitions.",
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "!important modifier on background and text color",
				Code:     `<button class="!bg-red-500 !text-white">Delete</button>`,
			},
			{
				Language: "tsx",
				Comment:  "!important on semantic and hover colors in JSX",
				Code: `export function Action() {
  return <div className="hover:!bg-primary !border-border">Action</div>;
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Proper layer-based specificity without !important",
				Code:     `<button class="bg-destructive text-destructive-foreground">Delete</button>`,
			},
			{
				Language: "tsx",
				Comment:  "Clean semantic classes with natural CSS cascade",
				Code: `export function Action() {
  return <div className="hover:bg-primary border-border">Action</div>;
}`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Cascade Arms Race",
				Severity: "HIGH",
				Impact:   "Forces downstream theme overrides to duplicate !important, breaking modular CSS encapsulation.",
			},
			{
				Vector:   "Dark Mode Override Failure",
				Severity: "HIGH",
				Impact:   "Dark mode variants fail to override base !important styles, resulting in inverted visual glitches.",
			},
		},
	}
}

func isNonColorText(s string) bool {
	switch s {
	case "left", "center", "right", "justify", "start", "end",
		"xs", "sm", "base", "lg", "xl", "2xl", "3xl", "4xl", "5xl", "6xl", "7xl", "8xl", "9xl",
		"ellipsis", "clip", "wrap", "nowrap", "balance", "pretty":
		return true
	default:
		return false
	}
}

func isNonColorBg(s string) bool {
	switch s {
	case "bottom", "center", "left", "right", "top",
		"fixed", "local", "scroll",
		"no-repeat", "repeat", "repeat-x", "repeat-y", "repeat-round", "repeat-space",
		"auto", "cover", "contain", "none":
		return true
	default:
		return strings.HasPrefix(s, "gradient-")
	}
}

func isColorUtility(prefix, remainder string) bool {
	if prefix == "text-" && isNonColorText(remainder) {
		return false
	}
	if prefix == "bg-" && isNonColorBg(remainder) {
		return false
	}
	if strings.HasPrefix(prefix, "border") && IsNonColorBorderKeyword(remainder) {
		return false
	}
	return true
}

// Evaluate mengevaluasi node IR dan mendeteksi modifier !important pada utility pewarnaan.
func (r *ImportantOverrideRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || len(node.Classes) == 0 {
		return nil
	}

	var diags []ir.Diagnostic
	for _, class := range node.Classes {
		_, base := StripVariants(class)
		if !strings.HasPrefix(base, "!") {
			continue
		}

		baseNoExcl := base[1:]
		baseNoAlpha, _, _ := SplitAlphaModifier(baseNoExcl)

		prefix, remainder, ok := SplitColorPrefix(baseNoAlpha)
		if ok && isColorUtility(prefix, remainder) {
			diags = append(diags, ir.Diagnostic{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Important modifier on color utility class: \"" + class + "\"",
				Hint:     "Avoid !important on color utilities; rely on CSS cascade layers and semantic token specificity instead.",
			})
		}
	}

	return diags
}
