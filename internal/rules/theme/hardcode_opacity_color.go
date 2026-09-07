package theme

import (
	"strings"
	"sync"

	"github.com/will2469/charites/internal/ir"
	themeengine "github.com/will2469/charites/internal/token"
)

var (
	discoveredTheme     themeengine.Context
	discoveredThemeOnce sync.Once
)

func getDiscoveredTheme() themeengine.Context {
	discoveredThemeOnce.Do(func() {
		ctx, err := themeengine.DiscoverAndLoad("", "")
		if err == nil && ctx != nil {
			discoveredTheme = ctx
		} else {
			discoveredTheme = themeengine.NewEmptyContext()
		}
	})
	return discoveredTheme
}

// ReplacementFor mengembalikan token semantik pengganti untuk pasangan warna/opacity tertentu secara read-only
// dari Context yang terdeteksi dari CSS SSOT menggunakan konvensi default.
func ReplacementFor(token string) (string, bool) {
	slashIdx := strings.IndexByte(token, '/')
	if slashIdx == -1 {
		return "", false
	}
	colorBase := token[:slashIdx]
	opacity := token[slashIdx+1:]
	cands, ok := NewDefaultCharitesConvention().FindOpacityReplacement(colorBase, opacity, getDiscoveredTheme())
	if !ok || len(cands) == 0 {
		return "", false
	}
	return cands[0].Name, true
}

// HardcodeOpacityColorRule mengimplementasikan aturan static analysis "theme.hardcode-opacity-color".
type HardcodeOpacityColorRule struct {
	themeCtx   themeengine.Context
	convention TokenConvention
}

// NewHardcodeOpacityColorRule membuat instance baru HardcodeOpacityColorRule dengan Context yang ditemukan otomatis.
func NewHardcodeOpacityColorRule() *HardcodeOpacityColorRule {
	return &HardcodeOpacityColorRule{
		themeCtx:   getDiscoveredTheme(),
		convention: NewDefaultCharitesConvention(),
	}
}

// NewHardcodeOpacityColorRuleWithTheme membuat instance baru HardcodeOpacityColorRule dengan dynamic Context.
func NewHardcodeOpacityColorRuleWithTheme(themeCtx themeengine.Context) *HardcodeOpacityColorRule {
	if themeCtx == nil {
		themeCtx = getDiscoveredTheme()
	}
	return &HardcodeOpacityColorRule{
		themeCtx:   themeCtx,
		convention: NewDefaultCharitesConvention(),
	}
}

// WithTheme menetapkan context tema dinamis ke rule.
func (r *HardcodeOpacityColorRule) WithTheme(themeCtx themeengine.Context) *HardcodeOpacityColorRule {
	if themeCtx == nil {
		themeCtx = getDiscoveredTheme()
	}
	r.themeCtx = themeCtx
	return r
}

// WithConvention menetapkan adapter konvensi semantik ke rule.
func (r *HardcodeOpacityColorRule) WithConvention(conv TokenConvention) *HardcodeOpacityColorRule {
	if conv != nil {
		r.convention = conv
	}
	return r
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *HardcodeOpacityColorRule) ID() string {
	return "theme.hardcode-opacity-color"
}

const uncalibratedOpacityHint = "Use an existing semantic token or declare a calibrated semantic token in global.css (e.g. --<base>-<state>) instead of using arbitrary slash opacity modifiers."

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *HardcodeOpacityColorRule) Description() string {
	return "Detects utility classes with hardcoded or uncalibrated slash opacity modifiers bypassing global.css SSOT"
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
				Comment:  "Arbitrary uncalibrated slash opacities and shadow utilities bypassing SSOT",
				Code: `export function ActionCard() {
  return (
    <div className="shadow-primary/20 hover:border-primary/50 bg-muted/20 border-warning/40 text-warning/90">
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

// Evaluate mengevaluasi sebuah node IR dan mendeteksi utility color ber-slash opacity yang melanggar SSOT.
// isNonColorOrDelegated determines whether a class prefix and color base should be ignored
// as a non-color utility or delegated to an orthogonal Charites theme rule.
func isNonColorOrDelegated(prefix, colorBase string) bool {
	// 1. Tolak keyword non-color utility (zero noise invariant)
	if prefix == "text-" && IsTailwindFontSize(colorBase) {
		return true
	}
	if strings.HasPrefix(prefix, "border") && IsNonColorBorderKeyword(colorBase) {
		return true
	}
	if prefix == "shadow-" && IsShadowSizeKeyword(colorBase) {
		return true
	}

	// 2. Delegasikan ke rule ortogonal sesuai Tri-Corpus SSOT
	if IsTailwindPrimitiveColor(colorBase) {
		return true
	}
	if IsMonochromeColor(colorBase) {
		return true
	}
	if strings.HasPrefix(colorBase, "[") || IsHexColor(colorBase) {
		return true
	}

	return false
}

func (r *HardcodeOpacityColorRule) evaluateClass(class string, span ir.Span) (ir.Diagnostic, bool) {
	// 1. Strip Tailwind variants (mendukung arbitrary variants ber-bracket seperti [&>svg]:...)
	base := StripVariantsOnlyBase(class)

	// 2. Split alpha modifier menggunakan parser resmi bracket-safe
	baseNoAlpha, alpha, hasAlpha := SplitAlphaModifier(base)
	if !hasAlpha {
		return ir.Diagnostic{}, false
	}

	// 3. Pisahkan prefix pewarnaan resmi menggunakan registry terpusat
	prefix, colorBase, ok := SplitColorPrefix(baseNoAlpha)
	if !ok || isNonColorOrDelegated(prefix, colorBase) {
		return ir.Diagnostic{}, false
	}

	// 4. Seluruh basis warna yang tersisa adalah kandidat semantik (semantic candidate).
	// Cari apakah terdapat token pengganti resmi terkalibrasi di SSOT global.css.
	tCtx := r.themeCtx
	if tCtx == nil {
		tCtx = getDiscoveredTheme()
	}

	conv := r.convention
	if conv == nil {
		conv = NewDefaultCharitesConvention()
	}

	hint := uncalibratedOpacityHint
	cands, found := conv.FindOpacityReplacement(colorBase, alpha, tCtx)
	if found && len(cands) > 0 {
		hint = "Use semantic token \"" + cands[0].Name + "\"."
	}

	return ir.Diagnostic{
		Line:     span.Line,
		Column:   span.Column,
		Rule:     r.ID(),
		Severity: r.DefaultSeverity(),
		Message:  "Hardcode opacity color: \"" + class + "\"",
		Hint:     hint,
	}, true
}

// Evaluate memeriksa apakah node mengandung utility class dengan arbitrary slash opacity.
// Mematuhi kontrak pure function, classification boundary bertingkat, dan zero alloc pada node bersih.
func (r *HardcodeOpacityColorRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || len(node.Classes) == 0 {
		return nil
	}

	//nolint:prealloc // zero-alloc on clean nodes required by QUAL-03
	var diags []ir.Diagnostic
	for _, class := range node.Classes {
		diag, ok := r.evaluateClass(class, node.Span)
		if ok {
			diags = append(diags, diag)
		}
	}

	return diags
}
