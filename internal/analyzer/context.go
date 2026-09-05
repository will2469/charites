package analyzer

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// Context menyimpan status dan hasil analisis per berkas dalam prinsip Share-Nothing.
type Context struct {
	FilePath      string
	InlineIgnores map[int][]string // Nomor baris (1-indexed) -> daftar Rule ID yang ditekan
	diagnostics   []ir.Diagnostic
}

// NewContext membuat instans Context baru yang terisolasi untuk sebuah berkas.
func NewContext(filePath string, inlineIgnores map[int][]string) *Context {
	if inlineIgnores == nil {
		inlineIgnores = make(map[int][]string)
	}
	return &Context{
		FilePath:      filePath,
		InlineIgnores: inlineIgnores,
		diagnostics:   make([]ir.Diagnostic, 0),
	}
}

// AddDiagnostic menambahkan diagnostic temuan ke dalam konteks berkas.
func (c *Context) AddDiagnostic(d ir.Diagnostic) {
	c.diagnostics = append(c.diagnostics, d)
}

// DiagnosticsList mengembalikan seluruh diagnostic yang terkumpul pada konteks ini.
func (c *Context) DiagnosticsList() []ir.Diagnostic {
	return c.diagnostics
}

// IsIgnored mengevaluasi apakah suatu temuan diagnostik ditekan oleh direktif inline ignore.
// Mendukung same-line trailing comment, next-line preceding comment, serta cakupan AST node span.
func (c *Context) IsIgnored(diag ir.Diagnostic, node *ir.Node) bool {
	// 1. Same-Line: Komentar berada di baris temuan yang sama
	if rules, ok := c.InlineIgnores[diag.Line]; ok {
		if matchesRule(rules, diag.Rule) {
			return true
		}
	}

	// 2. Next-Line: Komentar berada di baris tepat sebelum temuan (N-1)
	if diag.Line > 1 {
		if rules, ok := c.InlineIgnores[diag.Line-1]; ok {
			if matchesRule(rules, diag.Rule) {
				return true
			}
		}
	}

	// 3. AST Node Span Scope: Jika node membentang multi-baris
	if node != nil {
		// Evaluasi komentar tepat sebelum node dimulai (Line - 1)
		if node.Span.Line > 1 {
			if rules, ok := c.InlineIgnores[node.Span.Line-1]; ok {
				if diag.Line >= node.Span.Line && diag.Line <= node.Span.EndLine {
					if matchesRule(rules, diag.Rule) {
						return true
					}
				}
			}
		}
		// Evaluasi komentar pada opening tag baris pertama node (Line)
		if rules, ok := c.InlineIgnores[node.Span.Line]; ok {
			if diag.Line >= node.Span.Line && diag.Line <= node.Span.EndLine {
				if matchesRule(rules, diag.Rule) {
					return true
				}
			}
		}
	}

	return false
}

func matchesRule(rules []string, ruleID string) bool {
	for _, r := range rules {
		if r == "*" || r == ruleID {
			return true
		}
	}
	return false
}

// ParseDirectives memindai berkas sumber baris-per-baris (1-indexed) dan mengekstrak
// direktif penekanan `charites:ignore`.
// Mendukung format JavaScript/TS (`//`, `/* ... */`) dan HTML/Astro (`<!-- ... -->`).
func ParseDirectives(src []byte) map[int][]string {
	result := make(map[int][]string)
	scanner := bufio.NewScanner(bytes.NewReader(src))
	lineNum := 0

	const directiveMarker = "charites:ignore"

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		idx := strings.Index(line, directiveMarker)
		if idx == -1 {
			continue
		}

		// Ambil substring setelah "charites:ignore"
		after := line[idx+len(directiveMarker):]

		// Bersihkan penutup komentar jika ada
		after = strings.TrimSuffix(after, "-->")
		after = strings.TrimSuffix(after, "*/")
		after = strings.TrimSpace(after)

		if after == "" {
			// Direktif kosong: tidak menekan apa pun
			continue
		}

		// Parse daftar rule ID yang dipisahkan koma
		tokens := strings.Split(after, ",")
		seen := make(map[string]bool)
		var ruleList []string

		for _, tok := range tokens {
			cleanID := strings.TrimSpace(tok)
			if cleanID != "" && !seen[cleanID] {
				seen[cleanID] = true
				ruleList = append(ruleList, cleanID)
			}
		}

		if len(ruleList) > 0 {
			result[lineNum] = ruleList
		}
	}

	return result
}
