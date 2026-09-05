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

// ExtractTopLevelVarCalls mengekstrak seluruh pemanggilan var(...) tingkat teratas (top-level)
// dari sebuah string deklarasi CSS dengan menjaga keseimbangan tanda kurung dan tanda kutip.
func ExtractTopLevelVarCalls(src string) []VarCall {
	var calls []VarCall
	n := len(src)
	i := 0

	for i < n {
		// Cari "var(" secara case-insensitive
		if i+4 <= n && (src[i] == 'v' || src[i] == 'V') && strings.EqualFold(src[i:i+4], "var(") {
			// Pastikan tidak didahului oleh karakter ident (misal 'myvar(')
			if i > 0 && isIdentChar(src[i-1]) {
				i++
				continue
			}

			startOffset := i
			i += 4 // lewati "var("

			// Lewati whitespace di depan nama token
			for i < n && isWhitespace(src[i]) {
				i++
			}

			// Nama token harus diawali tanda hubung ganda "--"
			if i+2 > n || src[i] != '-' || src[i+1] != '-' {
				continue
			}

			// Kumpulkan nama token (mendukung CSS escape seperti \:)
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

			// Lewati whitespace setelah nama token
			for i < n && isWhitespace(src[i]) {
				i++
			}

			if i >= n {
				break
			}

			if src[i] == ',' {
				i++ // lewati ','
				// Lewati whitespace sebelum nilai fallback
				for i < n && isWhitespace(src[i]) {
					i++
				}
				fallbackStart := i

				// Pindai fallback dengan pelacakan kedalaman kurung seimbang (balanced parenthesis)
				parenDepth := 1
				bracketDepth := 0
				braceDepth := 0
				var inQuote byte

				for i < n && parenDepth > 0 {
					b := src[i]
					if inQuote != 0 {
						if b == '\\' && i+1 < n {
							i += 2
							continue
						}
						if b == inQuote {
							inQuote = 0
						}
						i++
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
							// Tanda kurung penutup var() utama tercapai
							fallback := strings.TrimSpace(src[fallbackStart:i])
							i++ // lewati ')'
							calls = append(calls, VarCall{
								Raw:         src[startOffset:i],
								StartOffset: startOffset,
								EndOffset:   i,
								Name:        name,
								Fallback:    fallback,
								HasFallback: true,
							})
							goto nextOuter
						}
					case '[':
						bracketDepth++
					case ']':
						if bracketDepth > 0 {
							bracketDepth--
						}
					case '{':
						braceDepth++
					case '}':
						if braceDepth > 0 {
							braceDepth--
						}
					}
					i++
				}
				continue
			} else if src[i] == ')' {
				i++ // lewati ')'
				calls = append(calls, VarCall{
					Raw:         src[startOffset:i],
					StartOffset: startOffset,
					EndOffset:   i,
					Name:        name,
					Fallback:    "",
					HasFallback: false,
				})
				continue
			}
		}
		i++
	nextOuter:
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
