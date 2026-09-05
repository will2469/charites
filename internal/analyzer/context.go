package analyzer

import (
	"bytes"
	"slices"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// Context menyimpan status dan hasil analisis per berkas dalam prinsip Share-Nothing.
// Mengisolasi data konfigurasi dan temuan diagnostik secara penuh.
type Context struct {
	FilePath      string
	InlineIgnores map[int][]string // Nomor baris (1-indexed) -> daftar Rule ID yang ditekan
	diagnostics   []ir.Diagnostic
}

// NewContext membuat instans Context baru yang terisolasi untuk sebuah berkas.
// Melakukan deep-copy pada map dan slice inlineIgnores demi menjamin integritas memori caller.
func NewContext(filePath string, inlineIgnores map[int][]string) *Context {
	copiedIgnores := make(map[int][]string, len(inlineIgnores))
	for lineNum, rules := range inlineIgnores {
		copiedRules := make([]string, len(rules))
		copy(copiedRules, rules)
		copiedIgnores[lineNum] = copiedRules
	}

	return &Context{
		FilePath:      filePath,
		InlineIgnores: copiedIgnores,
		diagnostics:   make([]ir.Diagnostic, 0),
	}
}

// AddDiagnostic menambahkan diagnostic temuan ke dalam konteks berkas.
func (c *Context) AddDiagnostic(d ir.Diagnostic) {
	c.diagnostics = append(c.diagnostics, d)
}

// DiagnosticsList mengembalikan klon independen dari seluruh diagnostic yang terkumpul pada konteks ini.
// Memutus pointer aliasing sehingga mutasi eksternal tidak mengubah state internal Context.
func (c *Context) DiagnosticsList() []ir.Diagnostic {
	return slices.Clone(c.diagnostics)
}

// IsIgnored mengevaluasi apakah suatu temuan diagnostik ditekan oleh direktif inline ignore.
// Mendukung same-line trailing comment, next-line preceding comment, serta cakupan AST node span.
func (c *Context) IsIgnored(diag ir.Diagnostic, node *ir.Node) bool {
	// 1. Same-Line: Komentar berada di baris temuan yang sama
	if rules, ok := c.InlineIgnores[diag.Line]; ok && matchesRule(rules, diag.Rule) {
		return true
	}

	// 2. Next-Line: Komentar berada di baris tepat sebelum temuan (N-1)
	if diag.Line > 1 {
		if rules, ok := c.InlineIgnores[diag.Line-1]; ok && matchesRule(rules, diag.Rule) {
			return true
		}
	}

	// 3. AST Node Span Scope: Jika node membentang multi-baris
	if node != nil {
		return c.isNodeSpanIgnored(diag, node)
	}

	return false
}

func (c *Context) isNodeSpanIgnored(diag ir.Diagnostic, node *ir.Node) bool {
	if diag.Line < node.Span.Line || diag.Line > node.Span.EndLine {
		return false
	}

	// Evaluasi komentar tepat sebelum node dimulai (Line - 1)
	if node.Span.Line > 1 {
		if rules, ok := c.InlineIgnores[node.Span.Line-1]; ok && matchesRule(rules, diag.Rule) {
			return true
		}
	}
	// Evaluasi komentar pada opening tag baris pertama node (Line)
	if rules, ok := c.InlineIgnores[node.Span.Line]; ok && matchesRule(rules, diag.Rule) {
		return true
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

// ParseDirectives memindai berkas sumber secara leksikal dan mengekstrak
// direktif penekanan `charites:ignore`.
// Hanya mengenali penanda yang berada di dalam komentar asli (Line Comment `//`,
// Block Comment `/* ... */`, atau HTML Comment `<!-- ... -->`).
// Menolak tegas substring yang berada di dalam string literal, template literal,
// atribut JSX, maupun plain HTML/text.
func ParseDirectives(src []byte) map[int][]string {
	result := make(map[int][]string)
	n := len(src)
	pos := 0
	line := 1

	for pos < n {
		ch := src[pos]

		if ch == '\n' {
			line++
			pos++
			continue
		}

		if parseCommentOrLiteral(src, &pos, &line, result, true) {
			continue
		}

		pos++
	}

	return result
}

func parseCommentOrLiteral(src []byte, pos *int, line *int, result map[int][]string, allowHTML bool) bool {
	rem := src[*pos:]

	if bytes.HasPrefix(rem, []byte("//")) {
		*pos += 2
		parseLineComment(src, pos, *line, result)
		return true
	}

	if bytes.HasPrefix(rem, []byte("/*")) {
		*pos += 2
		parseBlockComment(src, pos, line, result)
		return true
	}

	if allowHTML && bytes.HasPrefix(rem, []byte("<!--")) {
		*pos += 4
		parseHTMLComment(src, pos, line, result)
		return true
	}

	switch rem[0] {
	case '"', '\'':
		quote := rem[0]
		*pos++
		skipString(src, pos, line, quote)
		return true
	case '`':
		*pos++
		skipTemplateLiteral(src, pos, line, result)
		return true
	}

	return false
}

func parseLineComment(src []byte, pos *int, line int, result map[int][]string) {
	n := len(src)
	startComment := *pos
	for *pos < n && src[*pos] != '\n' {
		(*pos)++
	}
	text := string(src[startComment:*pos])
	if rules := extractDirectiveRules(text); len(rules) > 0 {
		addDirectiveRules(result, line, rules)
	}
}

func parseBlockComment(src []byte, pos *int, line *int, result map[int][]string) {
	n := len(src)
	lineStart := *pos
	currLine := *line

	for *pos < n {
		if *pos+1 < n && src[*pos] == '*' && src[*pos+1] == '/' {
			text := string(src[lineStart:*pos])
			if rules := extractDirectiveRules(text); len(rules) > 0 {
				addDirectiveRules(result, currLine, rules)
			}
			*pos += 2
			return
		}

		if src[*pos] == '\n' {
			text := string(src[lineStart:*pos])
			if rules := extractDirectiveRules(text); len(rules) > 0 {
				addDirectiveRules(result, currLine, rules)
			}
			(*line)++
			(*pos)++
			lineStart = *pos
			currLine = *line
			continue
		}

		(*pos)++
	}

	// Tangani sisa blok jika berkas berakhir tanpa penutup */
	if lineStart < *pos {
		text := string(src[lineStart:*pos])
		if rules := extractDirectiveRules(text); len(rules) > 0 {
			addDirectiveRules(result, currLine, rules)
		}
	}
}

func parseHTMLComment(src []byte, pos *int, line *int, result map[int][]string) {
	n := len(src)
	lineStart := *pos
	currLine := *line

	for *pos < n {
		if *pos+2 < n && src[*pos] == '-' && src[*pos+1] == '-' && src[*pos+2] == '>' {
			text := string(src[lineStart:*pos])
			if rules := extractDirectiveRules(text); len(rules) > 0 {
				addDirectiveRules(result, currLine, rules)
			}
			*pos += 3
			return
		}

		if src[*pos] == '\n' {
			text := string(src[lineStart:*pos])
			if rules := extractDirectiveRules(text); len(rules) > 0 {
				addDirectiveRules(result, currLine, rules)
			}
			(*line)++
			(*pos)++
			lineStart = *pos
			currLine = *line
			continue
		}

		(*pos)++
	}

	if lineStart < *pos {
		text := string(src[lineStart:*pos])
		if rules := extractDirectiveRules(text); len(rules) > 0 {
			addDirectiveRules(result, currLine, rules)
		}
	}
}

func skipString(src []byte, pos *int, line *int, quote byte) {
	n := len(src)
	for *pos < n {
		if src[*pos] == quote {
			(*pos)++
			return
		}
		if src[*pos] == '\\' && *pos+1 < n {
			(*pos) += 2
			continue
		}
		if src[*pos] == '\n' {
			(*line)++
		}
		(*pos)++
	}
}

func skipTemplateLiteral(src []byte, pos *int, line *int, result map[int][]string) {
	n := len(src)
	for *pos < n {
		if src[*pos] == '`' {
			(*pos)++
			return
		}
		if src[*pos] == '\\' && *pos+1 < n {
			(*pos) += 2
			continue
		}
		if src[*pos] == '\n' {
			(*line)++
			(*pos)++
			continue
		}
		if src[*pos] == '$' && *pos+1 < n && src[*pos+1] == '{' {
			*pos += 2
			parseTemplateExpr(src, pos, line, result)
			continue
		}
		(*pos)++
	}
}

func parseTemplateExpr(src []byte, pos *int, line *int, result map[int][]string) {
	n := len(src)
	braceDepth := 1

	for *pos < n && braceDepth > 0 {
		ch := src[*pos]

		if ch == '\n' {
			(*line)++
			(*pos)++
			continue
		}

		if ch == '{' {
			braceDepth++
			(*pos)++
			continue
		}

		if ch == '}' {
			braceDepth--
			(*pos)++
			if braceDepth == 0 {
				return
			}
			continue
		}

		if parseCommentOrLiteral(src, pos, line, result, false) {
			continue
		}

		(*pos)++
	}
}

func extractDirectiveRules(commentLine string) []string {
	const marker = "charites:ignore"
	idx := strings.Index(commentLine, marker)
	if idx == -1 {
		return nil
	}

	after := commentLine[idx+len(marker):]
	after = strings.TrimSpace(after)
	after = strings.TrimSuffix(after, "-->")
	after = strings.TrimSuffix(after, "*/")
	after = strings.TrimSuffix(after, "}")
	after = strings.TrimSpace(after)

	if after == "" {
		return nil
	}

	tokens := strings.Split(after, ",")
	seen := make(map[string]bool)
	var ruleList []string

	for _, tok := range tokens {
		cleanID := strings.TrimSpace(tok)
		cleanID = strings.TrimSuffix(cleanID, "-->")
		cleanID = strings.TrimSuffix(cleanID, "*/")
		cleanID = strings.TrimSuffix(cleanID, "}")
		cleanID = strings.TrimSpace(cleanID)
		if cleanID != "*" {
			cleanID = strings.TrimSuffix(cleanID, "*")
			cleanID = strings.TrimSpace(cleanID)
		}

		if cleanID != "" && !seen[cleanID] {
			seen[cleanID] = true
			ruleList = append(ruleList, cleanID)
		}
	}

	return ruleList
}

func addDirectiveRules(result map[int][]string, line int, rules []string) {
	if len(rules) == 0 {
		return
	}
	existing := result[line]
	seen := make(map[string]bool, len(existing)+len(rules))
	for _, r := range existing {
		seen[r] = true
	}
	for _, r := range rules {
		if !seen[r] {
			seen[r] = true
			existing = append(existing, r)
		}
	}
	result[line] = existing
}
