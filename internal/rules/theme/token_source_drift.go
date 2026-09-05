package theme

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// TokenSourceDriftRule mendeteksi penugasan nilai warna mentah (hex/rgb) langsung ke CSS custom properties
// di dalam markup komponen atau stylesheet lokal, yang melanggar prinsip Single Source of Truth (SSOT).
type TokenSourceDriftRule struct{}

// NewTokenSourceDriftRule membuat instance baru TokenSourceDriftRule.
func NewTokenSourceDriftRule() *TokenSourceDriftRule {
	return &TokenSourceDriftRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *TokenSourceDriftRule) ID() string {
	return "theme.token-source-drift"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *TokenSourceDriftRule) Description() string {
	return "Detects hardcoded color values bypassing the single source of truth design token pipeline"
}

// Category mengembalikan nama kategori rule.
func (r *TokenSourceDriftRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *TokenSourceDriftRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki otomatis.
func (r *TokenSourceDriftRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Design Tokens Community Group (DTCG)",
			"Single Source of Truth (SSOT) Architecture",
		},
		CoreInvariant: "Custom properties representing theme tokens must not be assigned raw color literals in component scopes; they must resolve to SSOT token references.",
		Grounding: "Assigning raw hex/rgb color values directly to theme custom properties inside components or local stylesheets " +
			"fractures the design token pipeline.\n\n" +
			"When developers write style=\"--primary: #2563eb\" or declare local --color-brand: #3b82f6:\n" +
			"1. Drift from Global SSOT: The component diverges from centralized theme tokens (global.css), creating fragmented brand colors.\n" +
			"2. Theme Switching Failure: Dynamic theme changes (e.g. high-contrast, dark mode, multi-tenant branding) cannot override local hardcoded values.\n" +
			"3. Design System Audit Blind Spot: Design linters fail to track where rogue colors enter the application.\n\n" +
			"Charites enforces binding theme variables to global design tokens via var(--...) instead of raw literals.",
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Hardcoded hex assigned to theme token in inline style",
				Code:     `<div style="--primary: #2563eb; --background: #ffffff;">Drifting Tokens</div>`,
			},
			{
				Language: "tsx",
				Comment:  "Hardcoded rgb assigned to custom property in JSX style",
				Code: `export function Header() {
  return <header style={{ '--color-brand': 'rgb(37, 99, 235)' }}>Drifted Header</header>;
}`,
			},
			{
				Language: "astro",
				Comment:  "Raw color assigned to theme custom property in style tag",
				Code: `<style>
  .card {
    --card-bg: #1e293b;
  }
</style>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Theme token mapped via SSOT variable reference",
				Code:     `<div style="--primary: var(--color-blue-600);">SSOT Aligned</div>`,
			},
			{
				Language: "tsx",
				Comment:  "Non-color numeric custom property",
				Code: `export function Tabs() {
  return <div style={{ '--tab-index': '2' }}>Safe Property</div>;
}`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Token SSOT Incoherence",
				Severity: "HIGH",
				Impact:   "Hardcoded local variable assignments decouple components from global design system updates.",
			},
			{
				Vector:   "Theme Switch Blind Spot",
				Severity: "HIGH",
				Impact:   "Local variable assignments prevent dynamic color schemes and tenant styling from cascading.",
			},
		},
	}
}

// Evaluate mengevaluasi node IR untuk mencari deklarasi custom property dengan nilai warna mentah.
func (r *TokenSourceDriftRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil {
		return nil
	}

	// 1. Periksa style attribute
	if styleVal, ok := node.GetAttr("style"); ok && strings.Contains(styleVal, "--") {
		if decl, found := findDriftingTokenDeclaration(styleVal); found {
			return []ir.Diagnostic{
				{
					Line:     node.Span.Line,
					Column:   node.Span.Column,
					Rule:     r.ID(),
					Severity: r.DefaultSeverity(),
					Message:  "Theme token source drift: raw color literal assigned to custom property \"" + decl + "\"",
					Hint:     "Reference standard design tokens via var(...) or configure tokens in the design system SSOT instead of assigning raw color literals directly.",
				},
			}
		}
	}

	// 2. Periksa <style> tag
	if node.Tag == "style" {
		cssText := getStyleNodeText(node)
		if strings.Contains(cssText, "--") {
			cleaned := stripCSSCommentsString(cssText)
			if decl, found := findDriftingTokenDeclaration(cleaned); found {
				return []ir.Diagnostic{
					{
						Line:     node.Span.Line,
						Column:   node.Span.Column,
						Rule:     r.ID(),
						Severity: r.DefaultSeverity(),
						Message:  "Theme token source drift: raw color literal assigned to custom property \"" + decl + "\"",
						Hint:     "Reference standard design tokens via var(...) or configure tokens in the design system SSOT instead of assigning raw color literals directly.",
					},
				}
			}
		}
	}

	return nil
}

// findDriftingTokenDeclaration mencari deklarasi --prop: <raw-color>;
func findDriftingTokenDeclaration(s string) (string, bool) {
	idx := 0
	n := len(s)

	for idx < n {
		pos := strings.Index(s[idx:], "--")
		if pos == -1 {
			break
		}

		start := idx + pos
		colon := strings.IndexByte(s[start:], ':')
		if colon == -1 {
			break
		}

		propName := strings.TrimSpace(s[start : start+colon])
		valStart := start + colon + 1

		// Cari batas akhir deklarasi (titik koma ';' atau kurung kurawal '}' atau newline)
		end := n
		for i := valStart; i < n; i++ {
			if s[i] == ';' || s[i] == '}' || s[i] == '\n' {
				end = i
				break
			}
		}

		val := strings.TrimSpace(s[valStart:end])
		if isRawColorLiteral(val) {
			return propName + ": " + val, true
		}

		idx = end + 1
	}

	return "", false
}

// isRawColorLiteral memeriksa apakah nilai deklarasi merupakan warna mentah (hex, rgb, hsl, oklch).
func isRawColorLiteral(val string) bool {
	// Jika menggunakan var(...) atau keyword transparan/inherit, bukan drifting raw color
	if strings.Contains(val, "var(") || val == "transparent" || val == "currentColor" || val == "inherit" {
		return false
	}

	// 1. Hex: #fff, #2563eb, dsb.
	if strings.HasPrefix(val, "#") && len(val) >= 4 {
		return true
	}

	// 2. Fungsi warna fungsional: rgb(...), rgba(...), hsl(...), hsla(...), oklch(...), lab(...)
	if strings.HasPrefix(val, "rgb(") || strings.HasPrefix(val, "rgba(") ||
		strings.HasPrefix(val, "hsl(") || strings.HasPrefix(val, "hsla(") ||
		strings.HasPrefix(val, "oklch(") || strings.HasPrefix(val, "lab(") {
		return true
	}

	return false
}
