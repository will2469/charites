package theme

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// opacityTokenMap memetakan utility warna dengan modifier slash opacity ke token semantik resmi pengganti.
// Bersifat unexported dan immutable-by-convention untuk menjamin purity, thread-safety, dan determinisme total.
// Sesuai SPEC-03-RULES Section 3.2.A.
var opacityTokenMap = map[string]string{
	"primary/10":     "primary-light",
	"primary/20":     "primary-light",
	"primary/5":      "primary-subtle",
	"secondary/10":   "muted-light",
	"secondary/5":    "muted-subtle",
	"destructive/10": "destructive-light",
	"destructive/20": "destructive-light",
	"destructive/5":  "destructive-subtle",
	"accent/10":      "accent-light",
	"accent/20":      "accent-light",
	"accent/5":       "accent-subtle",
	"warning/10":     "warning-light",
	"warning/5":      "warning-subtle",
	"muted/10":       "muted-light",
	"muted/5":        "muted-subtle",
	"amber/10":       "amber-light",
	"amber/5":        "amber-subtle",
	"emerald/10":     "emerald-light",
	"emerald/5":      "emerald-subtle",
}

// ReplacementFor mengembalikan token semantik pengganti untuk pasangan warna/opacity tertentu secara read-only.
// Mengembalikan replacement token dan true jika token terdaftar, atau string kosong dan false jika tidak terdaftar.
func ReplacementFor(token string) (string, bool) {
	rep, ok := opacityTokenMap[token]
	return rep, ok
}

// HardcodeOpacityColorRule mengimplementasikan aturan static analysis "theme.hardcode-opacity-color".
type HardcodeOpacityColorRule struct{}

// NewHardcodeOpacityColorRule membuat instance baru HardcodeOpacityColorRule.
func NewHardcodeOpacityColorRule() *HardcodeOpacityColorRule {
	return &HardcodeOpacityColorRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *HardcodeOpacityColorRule) ID() string {
	return "theme.hardcode-opacity-color"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *HardcodeOpacityColorRule) Description() string {
	return "Detects utility classes with hardcoded slash opacity modifiers that have official semantic token replacements"
}

// Category mengembalikan nama kategori rule.
func (r *HardcodeOpacityColorRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *HardcodeOpacityColorRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki otomatis.
func (r *HardcodeOpacityColorRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Design Tokens Community Group (DTCG)",
			"Tailwind CSS Design Token Architecture",
			"WCAG 2.2 Relative Contrast",
		},
		CoreInvariant: "Every color opacity variation that represents a semantic state or visual elevation must use a centralized semantic design token rather than an arbitrary slash modifier.",
		Grounding: "In modern design token architecture (such as Tailwind CSS with CSS Variables or OKLCH color spaces), semantic colors like primary and destructive are calibrated for foreground/background contrast against explicit color stops.\n\n" +
			"When developers append arbitrary slash modifiers (e.g. bg-primary/10), the resulting alpha-blended color:\n" +
			"1. Destroys WCAG 2.2 Contrast Predictability: Transparent alpha layers depend on whatever background color sits underneath. In dark mode or high-contrast themes, 10% opacity can drop contrast ratios below the 4.5:1 WCAG AA minimum.\n" +
			"2. Breaks Theme Export & Reusability: When exporting design tokens to mobile apps, Figma, or print styles, runtime alpha calculations cannot be resolved statically.\n" +
			"3. Creates Aesthetic Inconsistency: Different developers use varying opacities (/5, /10, /15, /20) for the same intended visual state (such as subtle hover backgrounds or tinted badge pills).\n\n" +
			"Charites enforces pre-calibrated semantic tokens (e.g. primary-light, primary-subtle, muted-light, destructive-light) that are mathematically verified for contrast and consistent across themes.",
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Direct slash opacity modifiers on semantic colors",
				Code: `<div class="card p-6 rounded-xl bg-primary/10 border border-destructive/20">
  <h2 class="text-xl font-bold text-primary/20">Card Title</h2>
  <span class="badge ring-1 ring-warning/10 bg-primary/5">Warning</span>
</div>`,
			},
			{
				Language: "tsx",
				Comment:  "Chained and single variants with hardcoded opacity",
				Code: `export function ActionCard() {
  return (
    <div className="p-4 rounded-lg hover:bg-primary/10 dark:bg-primary/10 md:hover:bg-primary/10">
      <button className="px-3 py-2 text-sm dark:border-destructive/20 sm:dark:hover:border-destructive/20">
        Delete
      </button>
    </div>
  );
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Using official semantic tokens from global.css",
				Code: `<div class="card p-6 rounded-xl bg-primary-light border border-destructive-light">
  <h2 class="text-xl font-bold text-primary">{Astro.props.title}</h2>
  <span class="badge ring-1 ring-warning-light bg-primary-subtle">Warning</span>
</div>`,
			},
			{
				Language: "tsx",
				Comment:  "Using semantic tokens with variants",
				Code: `export function ActionCard() {
  return (
    <div className="p-4 rounded-lg hover:bg-primary-light dark:bg-primary-light md:hover:bg-primary-light">
      <button className="px-3 py-2 text-sm dark:border-destructive-light">
        Delete
      </button>
    </div>
  );
}`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Accessibility Degradation",
				Severity: "HIGH",
				Impact:   "Contrast ratio drops below 4.5:1 under dark mode themes due to uncalibrated alpha blending.",
			},
			{
				Vector:   "Visual Debt & Inconsistency",
				Severity: "MEDIUM",
				Impact:   "Proliferation of slightly different opacities (/5, /10, /20) degrades product polish.",
			},
			{
				Vector:   "Theme Portability Failure",
				Severity: "MEDIUM",
				Impact:   "External design token exporters cannot map hardcoded alpha values to standalone color systems.",
			},
		},
	}
}

// stripVariants membuang seluruh prefix varian Tailwind (seperti hover:, dark:, md:hover:)
// dan mengembalikan base utility class.
func stripVariants(token string) string {
	lastColon := strings.LastIndexByte(token, ':')
	if lastColon >= 0 && lastColon < len(token)-1 {
		return token[lastColon+1:]
	}
	return token
}

// Evaluate mengevaluasi sebuah node IR dan mendeteksi utility color ber-slash opacity yang memiliki
// pemetaan semantic token pengganti resmi. Mematuhi kontrak pure function dan zero alloc pada node bersih.
func (r *HardcodeOpacityColorRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || len(node.Classes) == 0 {
		return nil
	}

	//nolint:prealloc // zero-alloc on clean nodes required by QUAL-03 allocation invariant
	var diags []ir.Diagnostic
	for _, class := range node.Classes {
		if strings.IndexByte(class, '/') == -1 {
			continue
		}

		base := stripVariants(class)
		var colorSlash string
		switch {
		case strings.HasPrefix(base, "bg-"):
			colorSlash = base[3:]
		case strings.HasPrefix(base, "text-"):
			colorSlash = base[5:]
		case strings.HasPrefix(base, "border-"):
			colorSlash = base[7:]
		case strings.HasPrefix(base, "ring-"):
			colorSlash = base[5:]
		default:
			continue
		}

		replacement, ok := opacityTokenMap[colorSlash]
		if !ok {
			continue
		}

		diags = append(diags, ir.Diagnostic{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Hardcode opacity color: \"" + class + "\"",
			Hint:     "Use semantic token \"" + replacement + "\".",
		})
	}

	return diags
}
