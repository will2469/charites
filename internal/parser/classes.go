package parser

import "strings"

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

	// Tangani format JSX kurung kurawal: class={...}
	if len(trimmed) >= 2 && strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}") {
		inner := strings.TrimSpace(trimmed[1 : len(trimmed)-1])
		if inner == "" {
			return nil, false
		}

		// String literal di dalam kurung kurawal: {"foo bar"} atau {'foo bar'}
		if len(inner) >= 2 && ((strings.HasPrefix(inner, "\"") && strings.HasSuffix(inner, "\"")) ||
			(strings.HasPrefix(inner, "'") && strings.HasSuffix(inner, "'"))) {
			content := inner[1 : len(inner)-1]
			return strings.Fields(content), false
		}

		// Template literal: {`foo ${bar} baz`}
		if len(inner) >= 2 && strings.HasPrefix(inner, "`") && strings.HasSuffix(inner, "`") {
			return extractTemplateLiteralClasses(inner[1 : len(inner)-1])
		}

		// Ekspresi variabel atau pemanggilan fungsi JS dinamis lainnya (misal class={cls})
		return nil, true
	}

	// Template literal langsung tanpa kurung kurawal terluar: `foo ${bar} baz`
	if len(trimmed) >= 2 && strings.HasPrefix(trimmed, "`") && strings.HasSuffix(trimmed, "`") {
		return extractTemplateLiteralClasses(trimmed[1 : len(trimmed)-1])
	}

	// String literal berkutip ganda atau tunggal: "foo bar" atau 'foo bar'
	if len(trimmed) >= 2 && ((strings.HasPrefix(trimmed, "\"") && strings.HasSuffix(trimmed, "\"")) ||
		(strings.HasPrefix(trimmed, "'") && strings.HasSuffix(trimmed, "'"))) {
		content := trimmed[1 : len(trimmed)-1]
		return strings.Fields(content), false
	}

	// Nilai atribut tanpa kutip (unquoted string)
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
