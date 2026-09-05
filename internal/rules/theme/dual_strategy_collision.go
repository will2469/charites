package theme

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// DualStrategyCollisionRule mendeteksi percampuran strategi tema (@media prefers-color-scheme
// bercampur dengan class .dark atau atribut [data-theme]) di dalam blok style yang sama.
type DualStrategyCollisionRule struct{}

// NewDualStrategyCollisionRule membuat instance baru DualStrategyCollisionRule.
func NewDualStrategyCollisionRule() *DualStrategyCollisionRule {
	return &DualStrategyCollisionRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *DualStrategyCollisionRule) ID() string {
	return "theme.dual-strategy-collision"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *DualStrategyCollisionRule) Description() string {
	return "Detects conflicting dark mode strategies (@media vs .dark/[data-theme]) in the same style scope"
}

// Category mengembalikan nama kategori rule.
func (r *DualStrategyCollisionRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *DualStrategyCollisionRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki otomatis.
func (r *DualStrategyCollisionRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C CSS Cascading and Inheritance Level 5",
			"Design System Theming Strategy Alignment",
			"Tailwind CSS Dark Mode Strategy Selector vs Media",
		},
		CoreInvariant: "Stylesheets must adhere to a single unified dark mode strategy (either media query or selector-based), avoiding contradictory cascade conflicts.",
		Grounding: "When developers mix @media (prefers-color-scheme: dark) with class (.dark) or attribute ([data-theme=\"dark\"]) selectors within the same scope:\n\n" +
			"1. Frankenstein Interface: System dark mode triggers media queries while manual theme toggles toggle classes, producing a fractured, half-dark layout.\n" +
			"2. Specificity Inversion: Class selectors have higher specificity than unnested media query elements, creating unpredictable styling overrides.\n" +
			"3. State Desynchronization: Manual UI theme switches fail to override hardcoded @media rules.\n\n" +
			"Charites enforces choosing a single, coherent dark mode switching strategy across each style scope.",
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Mixing prefers-color-scheme media query with .dark class selector",
				Code: `<style>
  @media (prefers-color-scheme: dark) {
    body {
      background: #121212;
    }
  }
  .dark {
    --bg-main: #000000;
  }
</style>`,
			},
			{
				Language: "tsx",
				Comment:  "Mixing media query with data-theme attribute in TSX style",
				Code: `<style>{` + "`" + `
  @media (prefers-color-scheme: dark) {
    :root { --card: #18181b; }
  }
  [data-theme="dark"] {
    --card: #09090b;
  }
` + "`" + `}</style>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Single coherent class-based strategy",
				Code: `<style>
  :root {
    --bg-main: #ffffff;
  }
  .dark {
    color-scheme: dark;
    --bg-main: #09090b;
  }
</style>`,
			},
			{
				Language: "tsx",
				Comment:  "Single coherent media-query-based strategy",
				Code: `<style>{` + "`" + `
  :root {
    --bg-main: #ffffff;
  }
  @media (prefers-color-scheme: dark) {
    :root {
      color-scheme: dark;
      --bg-main: #09090b;
    }
  }
` + "`" + `}</style>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Frankenstein UI Collision",
				Severity: "HIGH",
				Impact:   "System dark mode and application theme toggles conflict, resulting in partially inverted and illegible components.",
			},
			{
				Vector:   "Cascade Specificity Wars",
				Severity: "MEDIUM",
				Impact:   "Rules under @media cannot be overridden by user-selected theme classes without high-specificity hacks.",
			},
		},
	}
}

// Evaluate mengevaluasi node IR untuk memeriksa tabrakan strategi dark mode.
func (r *DualStrategyCollisionRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Tag != "style" {
		return nil
	}

	cssText := getStyleNodeText(node)
	if !strings.Contains(cssText, "prefers-color-scheme") {
		return nil
	}

	cleaned := stripCSSCommentsString(cssText)
	if hasDualStrategyCollision(cleaned) {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Dual theme strategy collision: stylesheet mixes '@media (prefers-color-scheme)' with '.dark' / '[data-theme]'",
				Hint:     "Adopt a single coherent dark mode strategy (e.g. class-based or data-attribute on root) to prevent conflicting theme cascades.",
			},
		}
	}

	return nil
}

func hasDualStrategyCollision(css string) bool {
	hasMediaDark := strings.Contains(css, "prefers-color-scheme")
	if !hasMediaDark {
		return false
	}

	hasSelectorDark := strings.Contains(css, ".dark") ||
		strings.Contains(css, "data-theme=\"dark\"") ||
		strings.Contains(css, "data-theme='dark'") ||
		strings.Contains(css, "data-theme=dark")

	return hasSelectorDark
}
