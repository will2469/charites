package theme

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// HydrationThemeMismatchRule mendeteksi root layout SSR (<head>) yang tidak menyertakan
// blocking inline script untuk menginisialisasi tema sebelum first paint, yang memicu Theme FOUC.
type HydrationThemeMismatchRule struct{}

// NewHydrationThemeMismatchRule membuat instance baru HydrationThemeMismatchRule.
func NewHydrationThemeMismatchRule() *HydrationThemeMismatchRule {
	return &HydrationThemeMismatchRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *HydrationThemeMismatchRule) ID() string {
	return "theme.hydration-theme-mismatch"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *HydrationThemeMismatchRule) Description() string {
	return "Detects SSR root layouts lacking blocking inline script for theme initialization"
}

// Category mengembalikan nama kategori rule.
func (r *HydrationThemeMismatchRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *HydrationThemeMismatchRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki otomatis.
func (r *HydrationThemeMismatchRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Web Performance Working Group (Core Web Vitals FOUC Prevention)",
			"React 18/19 Hydration Boundary Specification",
			"Astro SSR Zero-JS Script Tag Standards",
		},
		CoreInvariant: "Root SSR document layouts (<head>) must include a render-blocking inline theme script to resolve theme state before first paint and prevent theme FOUC.",
		Grounding: "In Server-Side Rendered (SSR) architectures (such as Astro, Next.js, or Remix):\n\n" +
			"1. Flash of Unstyled Theme (FOUC): If theme detection runs only after deferred client hydration (e.g. inside useEffect), the browser paints a blinding white default page before snapping jarringly to dark mode.\n" +
			"2. React Hydration Mismatch: Inconsistent theme attributes between server-rendered HTML and client hydration trigger React warning cascades and forced DOM re-mounts.\n" +
			"3. Cumulative Layout Shift (CLS): Font, border, or icon shifts caused by late theme flipping harm Core Web Vitals.\n\n" +
			"Charites enforces placing an inline render-blocking theme initialization script directly in the SSR root <head>.",
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "SSR root head in Astro without inline theme script",
				Code: `<html>
  <head>
    <meta charset="utf-8" />
    <title>Application</title>
  </head>
  <body>
    <slot />
  </body>
</html>`,
			},
			{
				Language: "tsx",
				Comment:  "Root head in TSX missing blocking theme initializer",
				Code: `<head>
  <meta charSet="utf-8" />
  <title>Dashboard</title>
</head>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Blocking inline theme script in Astro head",
				Code: `<head>
  <meta charset="utf-8" />
  <script is:inline>
    const theme = localStorage.getItem('theme') || (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');
    document.documentElement.classList.toggle('dark', theme === 'dark');
  </script>
</head>`,
			},
			{
				Language: "tsx",
				Comment:  "Blocking dangerouslySetInnerHTML theme script in TSX head",
				Code: `<head>
  <script
    dangerouslySetInnerHTML={{
      __html: "document.documentElement.classList.add(localStorage.getItem('theme') || 'light');",
    }}
  />
</head>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Theme FOUC Glare",
				Severity: "HIGH",
				Impact:   "Users in dark environments experience a painful full-screen white flash on every page navigation.",
			},
			{
				Vector:   "Hydration Error Cascade",
				Severity: "MEDIUM",
				Impact:   "React discards server-rendered DOM nodes upon encountering mismatched class attributes, increasing TTI.",
			},
		},
	}
}

// Evaluate mengevaluasi node IR untuk memeriksa apakah tag <head> memiliki script inisialisasi tema.
func (r *HydrationThemeMismatchRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || strings.ToLower(node.Tag) != "head" {
		return nil
	}

	// Hanya evaluasi jika dokumen mendukung tema / dark mode
	if !isThemeAwareDocument(node) {
		return nil
	}

	if hasBlockingThemeScript(node) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "SSR root <head> lacks blocking inline theme initialization script",
			Hint:     "Add a blocking inline script in <head> to synchronize document class with localStorage or prefers-color-scheme before first paint to prevent FOUC.",
		},
	}
}

func isThemeAwareDocument(head *ir.Node) bool {
	if head.Parent == nil {
		return false
	}

	// 1. Periksa atribut html induk
	if strings.Contains(head.Parent.RawClasses, "dark") || len(head.Parent.Attributes["data-theme"]) > 0 {
		return true
	}

	// 2. Periksa sibling body dan subtree di bawahnya untuk keberadaan utilitas dark: atau .dark
	for _, sibling := range head.Parent.Children {
		if sibling == head {
			continue
		}
		for n := range sibling.Walk() {
			if strings.Contains(n.RawClasses, "dark:") || strings.Contains(n.RawClasses, "dark") ||
				len(n.Attributes["data-theme"]) > 0 {
				return true
			}
		}
	}

	return false
}

func hasBlockingThemeScript(head *ir.Node) bool {
	for _, child := range head.Children {
		if strings.ToLower(child.Tag) != "script" {
			continue
		}

		// Periksa atribut dangerouslySetInnerHTML (React)
		if val, ok := child.Attributes["dangerouslySetInnerHTML"]; ok {
			lower := strings.ToLower(val)
			if strings.Contains(lower, "theme") || strings.Contains(lower, "dark") ||
				strings.Contains(lower, "color-scheme") || strings.Contains(lower, "localstorage") {
				return true
			}
		}

		// Periksa teks script di dalam child NodeText
		scriptText := getStyleNodeText(child)
		lower := strings.ToLower(scriptText)
		if strings.Contains(lower, "theme") || strings.Contains(lower, "dark") ||
			strings.Contains(lower, "color-scheme") || strings.Contains(lower, "localstorage") ||
			strings.Contains(lower, "prefers-color-scheme") {
			return true
		}
	}

	return false
}
