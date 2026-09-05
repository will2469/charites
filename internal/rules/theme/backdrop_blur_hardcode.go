package theme

import (
	"github.com/will2469/charites/internal/ir"
)

// BackdropBlurHardcodeRule mendeteksi penggunaan nilai blur atau backdrop-blur arbitrer
// di dalam kelas utilitas Tailwind atau arbitrary property (misal: backdrop-blur-[5px], blur-[12px]).
type BackdropBlurHardcodeRule struct{}

// NewBackdropBlurHardcodeRule membuat instance baru BackdropBlurHardcodeRule.
func NewBackdropBlurHardcodeRule() *BackdropBlurHardcodeRule {
	return &BackdropBlurHardcodeRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *BackdropBlurHardcodeRule) ID() string {
	return "theme.backdrop-blur-hardcode"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *BackdropBlurHardcodeRule) Description() string {
	return "Detects hardcoded arbitrary blur and backdrop-blur scalars in Tailwind utility classes"
}

// Category mengembalikan nama kategori rule.
func (r *BackdropBlurHardcodeRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *BackdropBlurHardcodeRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki otomatis.
func (r *BackdropBlurHardcodeRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Filter & Backdrop Filter Specification",
			"Design System Glassmorphism Standards",
			"Hardware-Accelerated Compositing Guidelines",
		},
		CoreInvariant: "Glassmorphism and surface blur effects must adhere to standardized blur tokens, never arbitrary scalar lengths.",
		Grounding: "Using arbitrary blur values (e.g. backdrop-blur-[5px] or blur-[12px]) produces inconsistent glassmorphism and performance bottlenecks:\n\n" +
			"1. GPU Overdraw Fragility: Arbitrary blur radii bypass optimized compositor layer pooling, causing unnecessary GPU rasterization penalties on mobile devices.\n" +
			"2. Glassmorphism Fragmentation: Slightly differing blur radii (e.g. 5px vs 8px vs 10px) across headers, dialogs, and drawer sheets ruin visual polish.\n" +
			"3. Inflexible Accessibility Adjustments: Standard tokens allow globally disabling or tuning blurs for users requesting reduced motion or low-power modes.\n\n" +
			"Charites enforces utilizing standard blur scale tokens (e.g. backdrop-blur-sm, backdrop-blur-md, backdrop-blur-lg) or CSS variables.",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Arbitrary backdrop-blur on navigation header",
				Code:     `<header className="backdrop-blur-[5px] bg-background/80">Sticky Nav</header>`,
			},
			{
				Language: "astro",
				Comment:  "Arbitrary filter blur in Astro component",
				Code:     `<div class="blur-[12px] [backdrop-filter:blur(7px)]">Frosted Panel</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Using standard backdrop blur token",
				Code:     `<header className="backdrop-blur-md bg-background/80">Sticky Nav</header>`,
			},
			{
				Language: "astro",
				Comment:  "Standard filter blur token",
				Code:     `<div class="blur-md backdrop-blur-sm">Frosted Panel</div>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Glassmorphism Visual Discordance",
				Severity: "MEDIUM",
				Impact:   "Irregular blur intensity breaks cohesive layering and depth hierarchy across interface overlays.",
			},
			{
				Vector:   "Mobile GPU Performance Stutter",
				Severity: "HIGH",
				Impact:   "Unstandardized backdrop-filter passes induce dropped frames during touch scrolling and bottom sheet gestures.",
			},
		},
	}
}

// Evaluate mengevaluasi sebuah node IR dan mendeteksi kelas blur arbitrer.
// Mematuhi kontrak pure function dan zero-alloc pada node bersih (QUAL-03).
func (r *BackdropBlurHardcodeRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || len(node.Classes) == 0 {
		return nil
	}

	var diags []ir.Diagnostic
	for _, class := range node.Classes {
		base := StripVariantsOnlyBase(class)
		if IsHardcodedBackdropBlur(base) {
			diags = append(diags, ir.Diagnostic{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Hardcoded blur scalar: \"" + class + "\"",
				Hint:     "Use a standard blur token (e.g. backdrop-blur-md, blur-sm) or CSS variable.",
			})
		}
	}

	return diags
}
