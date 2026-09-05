package theme

import (
	"bytes"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// UnlayeredTokenDefinitionRule mendeteksi deklarasi CSS custom properties (:root, [data-theme])
// di luar blok @layer theme atau @layer base sesuai CSS Cascade Layers Level 5.
type UnlayeredTokenDefinitionRule struct{}

// NewUnlayeredTokenDefinitionRule membuat instance baru UnlayeredTokenDefinitionRule.
func NewUnlayeredTokenDefinitionRule() *UnlayeredTokenDefinitionRule {
	return &UnlayeredTokenDefinitionRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *UnlayeredTokenDefinitionRule) ID() string {
	return "theme.unlayered-token-definition"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *UnlayeredTokenDefinitionRule) Description() string {
	return "Detects CSS custom property definitions declared outside @layer theme or @layer base"
}

// Category mengembalikan nama kategori rule.
func (r *UnlayeredTokenDefinitionRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *UnlayeredTokenDefinitionRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki otomatis.
func (r *UnlayeredTokenDefinitionRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Cascading Style Sheets (CSS) Level 5 (Cascade Layers)",
			"W3C CSS Custom Properties for Cascading Variables Module Level 1",
		},
		CoreInvariant: "CSS custom properties representing theme tokens must be declared within @layer theme or @layer base to ensure deterministic cascade resolution.",
		Grounding: "In modern frontend architectures and Tailwind CSS v4, unlayered CSS custom properties automatically take precedence " +
			"over all layered styles regardless of specificity.\n\n" +
			"When developers declare :root { --primary: #... } without @layer theme or @layer base:\n" +
			"1. Cascade Inversion: Unlayered rules override framework layers and variant cascades unexpectedly.\n" +
			"2. Dark Mode Clashes: Nested dark mode themes defined within layers cannot reliably override unlayered root variables.\n" +
			"3. Specificity Pollution: Subsequent theme overrides require !important or higher specificity hacks to function.\n\n" +
			"Charites enforces encapsulating theme custom property definitions inside @layer theme or @layer base.",
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Unlayered :root custom property definition in style tag",
				Code: `<style>
  :root {
    --primary: #2563eb;
    --background: #ffffff;
  }
</style>`,
			},
			{
				Language: "tsx",
				Comment:  "Unlayered [data-theme] custom property definition",
				Code: `export function GlobalStyles() {
  return (
    <style>{` + "`" + `
      :root {
        --brand-color: #3b82f6;
      }
    ` + "`" + `}</style>
  );
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Enclosed within @layer theme",
				Code: `<style>
  @layer theme {
    :root {
      --primary: #2563eb;
      --background: #ffffff;
    }
  }
</style>`,
			},
			{
				Language: "astro",
				Comment:  "Enclosed within @layer base",
				Code: `<style>
  @layer base {
    :root {
      --primary: #2563eb;
    }
  }
</style>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Cascade Priority Inversion",
				Severity: "HIGH",
				Impact:   "Unlayered properties override all cascade layers, preventing dark mode and variant styles from taking effect.",
			},
			{
				Vector:   "Theme Specificity Escalation",
				Severity: "MEDIUM",
				Impact:   "Teams resort to !important declarations to override unlayered variables, causing style degradation.",
			},
		},
	}
}

// Evaluate mengevaluasi node IR untuk memeriksa apakah deklarasi custom property berada di luar @layer.
func (r *UnlayeredTokenDefinitionRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Tag != "style" {
		return nil
	}

	cssText := getStyleNodeText(node)
	if cssText == "" || !strings.Contains(cssText, "--") {
		return nil
	}

	cleaned := stripCSSCommentsString(cssText)
	if !hasUnlayeredCustomProperty(cleaned) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "CSS custom properties defined outside @layer theme or @layer base",
			Hint:     "Wrap root token declarations inside @layer theme { ... } or @layer base { ... } to maintain predictable cascade layers.",
		},
	}
}

// getStyleNodeText mengekstrak seluruh teks CSS dari child node NodeText dalam <style>.
func getStyleNodeText(node *ir.Node) string {
	if len(node.Children) == 0 {
		return node.RawClasses
	}

	var sb strings.Builder
	for _, child := range node.Children {
		if child.Type == ir.NodeText {
			sb.WriteString(child.RawClasses)
		}
	}
	res := sb.String()
	if res == "" {
		return node.RawClasses
	}
	return res
}

// stripCSSCommentsString menghapus komentar /* ... */ dari CSS string.
func stripCSSCommentsString(src string) string {
	if !strings.Contains(src, "/*") {
		return src
	}

	var buf bytes.Buffer
	buf.Grow(len(src))

	i := 0
	n := len(src)
	for i < n {
		if i+1 < n && src[i] == '/' && src[i+1] == '*' {
			end := strings.Index(src[i+2:], "*/")
			if end == -1 {
				break
			}
			i += end + 4
			continue
		}
		buf.WriteByte(src[i])
		i++
	}
	return buf.String()
}

// hasUnlayeredCustomProperty memeriksa apakah ada deklarasi CSS custom property (--*)
// yang berada di luar blok @layer.
func hasUnlayeredCustomProperty(css string) bool {
	// Jika tidak ada variabel CSS sama sekali, aman
	if !strings.Contains(css, "--") {
		return false
	}

	// Buang seluruh blok @layer { ... } terlebih dahulu
	filtered := removeLayerBlocks(css)

	// Sekarang periksa apakah sisa CSS masih memuat deklarasi custom property (--foo: ...)
	return containsCustomPropertyDeclaration(filtered)
}

// removeLayerBlocks menghapus blok @layer ... { ... } bersarang dari CSS.
func removeLayerBlocks(css string) string {
	var buf strings.Builder
	idx := 0
	n := len(css)

	for idx < n {
		layerIdx := strings.Index(css[idx:], "@layer")
		if layerIdx == -1 {
			buf.WriteString(css[idx:])
			break
		}

		buf.WriteString(css[idx : idx+layerIdx])
		curr := idx + layerIdx + len("@layer")

		openBrace := strings.IndexByte(css[curr:], '{')
		if openBrace == -1 {
			semi := strings.IndexByte(css[curr:], ';')
			if semi != -1 {
				idx = curr + semi + 1
			} else {
				idx = n
			}
			continue
		}

		blockEnd := findMatchingBraceEnd(css, curr+openBrace+1)
		if blockEnd == -1 {
			break
		}

		idx = blockEnd + 1
	}

	return buf.String()
}

// findMatchingBraceEnd mencari indeks penutup kurung kurawal '}' yang seimbang.
func findMatchingBraceEnd(s string, start int) int {
	depth := 1
	for i := start; i < len(s); i++ {
		if s[i] == '{' {
			depth++
		} else if s[i] == '}' {
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// containsCustomPropertyDeclaration memeriksa apakah teks memuat deklarasi properti --name: value;
func containsCustomPropertyDeclaration(css string) bool {
	idx := 0
	n := len(css)

	for idx < n {
		varIdx := strings.Index(css[idx:], "--")
		if varIdx == -1 {
			break
		}

		pos := idx + varIdx + 2
		// Validasi nama variabel CSS
		identStart := pos
		for pos < n && isCSSIdentChar(css[pos]) {
			pos++
		}

		if pos > identStart {
			// Lewati spasi sebelum titik dua ':'
			scan := pos
			for scan < n && (css[scan] == ' ' || css[scan] == '\t' || css[scan] == '\r' || css[scan] == '\n') {
				scan++
			}
			if scan < n && css[scan] == ':' {
				// Ditemukan deklarasi properti custom di luar layer!
				return true
			}
		}

		idx = pos
	}

	return false
}

func isCSSIdentChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_'
}
