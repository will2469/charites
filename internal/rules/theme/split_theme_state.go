package theme

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// SplitThemeStateRule mendeteksi pembacaan atau penulisan status tema secara langsung
// melalui localStorage di dalam komponen atau event handler di luar ThemeProvider / <head>.
type SplitThemeStateRule struct{}

// NewSplitThemeStateRule membuat instance baru SplitThemeStateRule.
func NewSplitThemeStateRule() *SplitThemeStateRule {
	return &SplitThemeStateRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *SplitThemeStateRule) ID() string {
	return "theme.split-theme-state"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *SplitThemeStateRule) Description() string {
	return "Detects ad-hoc direct access to theme state via localStorage outside ThemeProvider"
}

// Category mengembalikan nama kategori rule.
func (r *SplitThemeStateRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *SplitThemeStateRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki otomatis.
func (r *SplitThemeStateRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Design System Single Source of Truth (SSOT) Architecture",
			"React State Management Best Practices (Context & Hooks)",
			"WCAG 2.2 Predictable Navigation & Consistency",
		},
		CoreInvariant: "Component UI state must consume theme through a unified ThemeProvider context or custom hook, never querying localStorage directly in component bodies or handlers.",
		Grounding: "When developers directly access or mutate localStorage.getItem('theme') or localStorage.theme in scattered components:\n\n" +
			"1. Fragmented State: Component A reads localStorage while Component B listens to React Context, causing disparate parts of the UI to display inconsistent themes.\n" +
			"2. Missing Reactivity: Updates directly to localStorage do not trigger React or framework re-renders across sibling components.\n" +
			"3. Testability Breakdown: Components cannot be unit tested or rendered in isolation without mocking global browser APIs.\n\n" +
			"Charites enforces routing all theme state access through a unified Theme Provider / useTheme hook, permitting direct localStorage access only in root <head> bootstrap scripts.",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Direct localStorage mutation in button onClick handler",
				Code:     `<button onClick={() => localStorage.setItem('theme', 'dark')}>Toggle</button>`,
			},
			{
				Language: "astro",
				Comment:  "Direct localStorage inspection in Astro component body",
				Code:     `<div data-theme={localStorage.getItem('theme')}>Container</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Using unified useTheme hook from ThemeProvider",
				Code: `const { theme, setTheme } = useTheme();
<button onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}>Toggle</button>`,
			},
			{
				Language: "astro",
				Comment:  "Permitted inline bootstrap script inside root <head>",
				Code: `<head>
  <script is:inline>
    const theme = localStorage.getItem('theme') || 'light';
    document.documentElement.classList.toggle('dark', theme === 'dark');
  </script>
</head>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Theme Desynchronization Across UI",
				Severity: "MEDIUM",
				Impact:   "Different page regions display discordant color schemes due to uncoordinated local state reads.",
			},
			{
				Vector:   "Broken Component Reactivity",
				Severity: "MEDIUM",
				Impact:   "Theme switches fail to re-render affected components without a full browser page refresh.",
			},
		},
	}
}

// Evaluate mengevaluasi node IR untuk memeriksa apakah ada akses langsung ke localStorage tema di luar <head>.
func (r *SplitThemeStateRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil {
		return nil
	}

	// Lewati jika node adalah <head> atau berada di dalam <head> (bootstrap script sah)
	for p := node; p != nil; p = p.Parent {
		if strings.ToLower(p.Tag) == "head" {
			return nil
		}
	}

	// 1. Periksa atribut (seperti onClick, data-theme, dll.)
	for _, v := range node.Attributes {
		if isAdHocThemeLocalStorage(v) {
			return []ir.Diagnostic{
				{
					Line:     node.Span.Line,
					Column:   node.Span.Column,
					Rule:     r.ID(),
					Severity: r.DefaultSeverity(),
					Message:  "Ad-hoc direct access to theme state via localStorage outside ThemeProvider",
					Hint:     "Use a unified ThemeProvider context or useTheme hook to maintain a single source of truth for theme state.",
				},
			}
		}
	}

	// 2. Jika node adalah <script> di luar <head>, periksa isi teksnya
	if strings.ToLower(node.Tag) == "script" {
		text := getStyleNodeText(node)
		if isAdHocThemeLocalStorage(text) {
			return []ir.Diagnostic{
				{
					Line:     node.Span.Line,
					Column:   node.Span.Column,
					Rule:     r.ID(),
					Severity: r.DefaultSeverity(),
					Message:  "Ad-hoc direct access to theme state via localStorage outside ThemeProvider",
					Hint:     "Use a unified ThemeProvider context or useTheme hook to maintain a single source of truth for theme state.",
				},
			}
		}
	}

	return nil
}

func isAdHocThemeLocalStorage(text string) bool {
	lower := strings.ToLower(text)
	if !strings.Contains(lower, "localstorage") {
		return false
	}

	return strings.Contains(lower, "getitem('theme'") ||
		strings.Contains(lower, "getitem(\"theme\"") ||
		strings.Contains(lower, "setitem('theme'") ||
		strings.Contains(lower, "setitem(\"theme\"") ||
		strings.Contains(lower, "localstorage.theme")
}
