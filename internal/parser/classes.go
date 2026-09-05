package parser

import "strings"

// ExtractClasses mengekstrak token kelas CSS dari nilai atribut class/className mentah.
// Mematuhi kontrak Option B:
// - Literal string biasa diekstrak dan ditokenisasi.
// - Template literal (backticks) mengekstrak segmen statis di luar ${...}.
// - Segmen di dalam ${...} diisolasi secara buram (opaque) tanpa parsing JS AST.
// - Flag hasDynamic bernilai true jika terdapat ekspresi ${...} atau ekspresi variabel dinamis.
func isQuotedString(s string) bool {
	return len(s) >= 2 && ((s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\''))
}

func isTemplateLiteral(s string) bool {
	return len(s) >= 2 && s[0] == '`' && s[len(s)-1] == '`'
}

func extractJSXExpressionClasses(inner string) ([]string, bool) {
	inner = strings.TrimSpace(inner)
	if inner == "" {
		return nil, false
	}
	if isQuotedString(inner) {
		return strings.Fields(inner[1 : len(inner)-1]), false
	}
	if isTemplateLiteral(inner) {
		return extractTemplateLiteralClasses(inner[1 : len(inner)-1])
	}
	return nil, true
}

// ExtractClasses mengekstrak token kelas CSS dari nilai atribut class/className mentah.
// Mematuhi kontrak Option B:
// - Literal string biasa diekstrak dan ditokenisasi.
// - Template literal (backticks) mengekstrak segmen statis di luar ${...}.
// - Segmen di dalam ${...} diisolasi secara buram (opaque) tanpa parsing JS AST.
// - Flag hasDynamic bernilai true jika terdapat ekspresi ${...} atau ekspresi variabel dinamis.
func ExtractClasses(raw string) (classes []string, hasDynamic bool) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, false
	}

	if len(trimmed) >= 2 && strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}") {
		return extractJSXExpressionClasses(trimmed[1 : len(trimmed)-1])
	}

	if isTemplateLiteral(trimmed) {
		return extractTemplateLiteralClasses(trimmed[1 : len(trimmed)-1])
	}

	if isQuotedString(trimmed) {
		return strings.Fields(trimmed[1 : len(trimmed)-1]), false
	}

	return strings.Fields(trimmed), false
}

// extractTemplateLiteralClasses memisahkan teks statis di luar ${...} dan menandai keberadaan segmen dinamis.
func extractTemplateLiteralClasses(templateStr string) ([]string, bool) {
	var staticParts []string
	hasDynamic := false
	idx := 0

	for idx < len(templateStr) {
		dollarIdx := strings.Index(templateStr[idx:], "${")
		if dollarIdx == -1 {
			// Sisa string adalah bagian statis murni
			staticParts = append(staticParts, templateStr[idx:])
			break
		}

		// Bagian statis sebelum ${
		staticParts = append(staticParts, templateStr[idx:idx+dollarIdx])
		hasDynamic = true

		// Cari kurung kurawal penutup '}' penyeimbang untuk ekspresi ${...}
		exprStart := idx + dollarIdx + 2
		braceDepth := 1
		exprEnd := -1

		for j := exprStart; j < len(templateStr); j++ {
			if templateStr[j] == '{' {
				braceDepth++
			} else if templateStr[j] == '}' {
				braceDepth--
				if braceDepth == 0 {
					exprEnd = j
					break
				}
			}
		}

		if exprEnd == -1 {
			// Ekspresi tidak ditutup sempurna, lewati sisa string
			break
		}

		idx = exprEnd + 1
	}

	var classes []string
	for _, part := range staticParts {
		classes = append(classes, strings.Fields(part)...)
	}

	return classes, hasDynamic
}
