package theme

import (
	"strings"
)

// VarCall merepresentasikan pemanggilan fungsi var(...) dalam CSS,
// diekstrak secara grammar-aware dengan kedalaman tanda kurung/tanda kutip seimbang (balanced parentheses).
type VarCall struct {
	// Raw adalah representasi teks pemanggilan var utuh, misal: `var(--foo, rgb(1, 2, 3))`
	Raw string

	// StartOffset adalah indeks byte awal pemanggilan "var(" pada string sumber.
	StartOffset int

	// EndOffset adalah indeks byte setelah tanda kurung penutup ")" pada string sumber.
	EndOffset int

	// Name adalah nama custom property yang dirujuk, misal: `--foo` (telah di-unescape jika ada escape).
	Name string

	// Fallback adalah nilai fallback jika token tidak terdefinisi, misal: `rgb(1, 2, 3)`.
	// Kosong jika pemanggilan tidak menyertakan fallback.
	Fallback string

	// HasFallback bernilai true jika terdapat argumen kedua setelah tanda koma.
	HasFallback bool
}

func isVarCallStart(src string, i, n int) bool {
	if i+4 > n || (src[i] != 'v' && src[i] != 'V') || !strings.EqualFold(src[i:i+4], "var(") {
		return false
	}
	return i == 0 || !isIdentChar(src[i-1])
}

func scanVarTokenName(src string, start, n int) (string, int, bool) {
	i := start
	for i < n && isWhitespace(src[i]) {
		i++
	}
	if i+2 > n || src[i] != '-' || src[i+1] != '-' {
		return "", i, false
	}

	nameStart := i
	for i < n {
		if isIdentChar(src[i]) {
			i++
			continue
		}
		if src[i] == '\\' && i+1 < n && isValidEscape(src[i], src[i+1]) {
			i += 2
			continue
		}
		break
	}
	name := UnescapeCSS(src[nameStart:i])

	for i < n && isWhitespace(src[i]) {
		i++
	}
	return name, i, true
}

func scanBalancedFallback(src string, start, n int) (string, int, bool) {
	i := start
	for i < n && isWhitespace(src[i]) {
		i++
	}
	fallbackStart := i

	parenDepth := 1
	var inQuote byte

	for i < n && parenDepth > 0 {
		b := src[i]
		if inQuote != 0 {
			i = advanceQuotedChar(b, &inQuote, i, n)
			continue
		}

		switch b {
		case '"', '\'':
			inQuote = b
		case '(':
			parenDepth++
		case ')':
			parenDepth--
			if parenDepth == 0 {
				fallback := strings.TrimSpace(src[fallbackStart:i])
				return fallback, i + 1, true
			}
		}
		i++
	}

	return "", i, false
}

func advanceQuotedChar(b byte, inQuote *byte, i, n int) int {
	if b == '\\' && i+1 < n {
		return i + 2
	}
	if b == *inQuote {
		*inQuote = 0
	}
	return i + 1
}

// ExtractTopLevelVarCalls mengekstrak seluruh pemanggilan var(...) tingkat teratas (top-level)
// dari sebuah string deklarasi CSS dengan menjaga keseimbangan tanda kurung dan tanda kutip.
func ExtractTopLevelVarCalls(src string) []VarCall {
	var calls []VarCall
	n := len(src)
	i := 0

	for i < n {
		if !isVarCallStart(src, i, n) {
			i++
			continue
		}

		startOffset := i
		name, nextIdx, ok := scanVarTokenName(src, i+4, n)
		if !ok {
			i++
			continue
		}
		i = nextIdx
		if i >= n {
			break
		}

		if src[i] == ',' {
			fallback, endIdx, okFallback := scanBalancedFallback(src, i+1, n)
			if okFallback {
				calls = append(calls, VarCall{
					Raw:         src[startOffset:endIdx],
					StartOffset: startOffset,
					EndOffset:   endIdx,
					Name:        name,
					Fallback:    fallback,
					HasFallback: true,
				})
				i = endIdx
				continue
			}
		} else if src[i] == ')' {
			calls = append(calls, VarCall{
				Raw:         src[startOffset : i+1],
				StartOffset: startOffset,
				EndOffset:   i + 1,
				Name:        name,
				Fallback:    "",
				HasFallback: false,
			})
			i++
			continue
		}

		i++
	}

	return calls
}

// ExtractAllVarNames mengekstrak seluruh nama custom property yang dirujuk oleh var()
// secara rekursif (termasuk yang berada di dalam rantai fallback bersarang).
func ExtractAllVarNames(src string) []string {
	if !strings.Contains(src, "var(") && !strings.Contains(src, "VAR(") {
		return nil
	}

	calls := ExtractTopLevelVarCalls(src)
	if len(calls) == 0 {
		return nil
	}

	var names []string
	seen := make(map[string]struct{}, len(calls))

	var collect func(list []VarCall)
	collect = func(list []VarCall) {
		for _, c := range list {
			if _, ok := seen[c.Name]; !ok {
				seen[c.Name] = struct{}{}
				names = append(names, c.Name)
			}
			if c.HasFallback && (strings.Contains(c.Fallback, "var(") || strings.Contains(c.Fallback, "VAR(")) {
				collect(ExtractTopLevelVarCalls(c.Fallback))
			}
		}
	}

	collect(calls)
	return names
}
