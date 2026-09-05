package theme

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// MissingColorSchemeRule mendeteksi deklarasi selektor mode gelap (.dark atau [data-theme="dark"])
// di dalam style tag atau stylesheet yang tidak menyertakan properti color-scheme: dark.
type MissingColorSchemeRule struct{}

// NewMissingColorSchemeRule membuat instance baru MissingColorSchemeRule.
func NewMissingColorSchemeRule() *MissingColorSchemeRule {
	return &MissingColorSchemeRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *MissingColorSchemeRule) ID() string {
	return "theme.missing-color-scheme"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *MissingColorSchemeRule) Description() string {
	return "Detects dark theme definitions (.dark, [data-theme=\"dark\"]) missing color-scheme property"
}

// Category mengembalikan nama kategori rule.
func (r *MissingColorSchemeRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *MissingColorSchemeRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki otomatis.
func (r *MissingColorSchemeRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C CSS Color Adjustment Module Level 1",
			"HTML Living Standard Section 4.2.5.5 (color-scheme)",
		},
		CoreInvariant: "Dark mode theme selectors must declare color-scheme: dark to synchronize native browser UI elements with the theme.",
		Grounding: "When developers configure dark mode using .dark or [data-theme='dark'] CSS rules without declaring color-scheme: dark:\n\n" +
			"1. White Form Controls: Native form elements (<select> dropdown popovers, <input type='date'> calendars, checkboxes) remain bright white.\n" +
			"2. Inverted Scrollbars: Operating system scrollbars fail to enter dark mode, glaring against dark page content.\n" +
			"3. Blinding Autofill Scrims: Browser credential autofill backgrounds turn blinding yellow-white.\n\n" +
			"Charites enforces declaring color-scheme: dark within all dark mode selector scopes.",
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Dark mode class without color-scheme declaration",
				Code: `<style>
  .dark {
    --background: #09090b;
    --foreground: #fafafa;
  }
</style>`,
			},
			{
				Language: "tsx",
				Comment:  "Data-theme attribute without color-scheme",
				Code: `export function DarkTheme() {
  return (
    <style>{` + "`" + `
      [data-theme="dark"] {
        --bg-main: #121212;
      }
    ` + "`" + `}</style>
  );
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Dark selector with explicit color-scheme",
				Code: `<style>
  .dark {
    color-scheme: dark;
    --background: #09090b;
    --foreground: #fafafa;
  }
</style>`,
			},
			{
				Language: "astro",
				Comment:  "Global color-scheme on root",
				Code: `<style>
  :root {
    color-scheme: light dark;
  }
</style>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Native UI Flash",
				Severity: "MEDIUM",
				Impact:   "Native controls, scrollbars, and dropdowns render in high-contrast light theme inside dark mode.",
			},
			{
				Vector:   "Form Usability Degradation",
				Severity: "LOW",
				Impact:   "Datepickers and select menus become illegible due to mismatched browser chrome defaults.",
			},
		},
	}
}

// Evaluate mengevaluasi node IR untuk memeriksa apakah selektor dark mode memiliki color-scheme.
func (r *MissingColorSchemeRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Tag != "style" {
		return nil
	}

	cssText := getStyleNodeText(node)
	if !strings.Contains(cssText, ".dark") && !strings.Contains(cssText, "data-theme") {
		return nil
	}

	cleaned := stripCSSCommentsString(cssText)
	if hasDarkSelectorWithoutColorScheme(cleaned) {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Dark theme definition lacks \"color-scheme: dark\" declaration",
				Hint:     "Add \"color-scheme: dark;\" inside your .dark or [data-theme=\"dark\"] selector so native browser controls (scrollbars, popovers) adapt properly.",
			},
		}
	}

	return nil
}

func hasDarkSelectorWithoutColorScheme(css string) bool {
	hasDarkSelector := strings.Contains(css, ".dark") ||
		strings.Contains(css, "data-theme=\"dark\"") ||
		strings.Contains(css, "data-theme='dark'") ||
		strings.Contains(css, "data-theme=dark")
	if !hasDarkSelector {
		return false
	}

	remaining := css
	for {
		idx := strings.Index(remaining, "color-scheme:")
		if idx == -1 {
			break
		}
		val := remaining[idx+len("color-scheme:"):]
		semi := strings.IndexByte(val, ';')
		decl := val
		if semi != -1 {
			decl = val[:semi]
			remaining = val[semi+1:]
		} else {
			remaining = ""
		}
		if strings.Contains(decl, "dark") {
			return false
		}
	}

	return true
}
