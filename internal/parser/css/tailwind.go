package css

import (
	"bytes"
	"strings"
)

// ThemeTokenRegistry menyimpan metadata deklarasi variabel CSS mentah dari blok @theme Tailwind CSS v4.
// Bersifat netral terhadap aturan (rule-agnostic), tanpa melakukan pemetaan semantik ekuivalensi opacity.
type ThemeTokenRegistry struct {
	Variables map[string]string // Key: "--color-primary", Value: "#2563eb"
}

// NewThemeTokenRegistry menginisialisasi registri token tema kosong.
func NewThemeTokenRegistry() *ThemeTokenRegistry {
	return &ThemeTokenRegistry{
		Variables: make(map[string]string),
	}
}

// ParseTheme memindai konten CSS dan mengekstrak seluruh deklarasi variabel di dalam blok @theme { ... }.
// Mengabaikan komentar CSS /* ... */ dan whitespace.
func ParseTheme(src []byte) (*ThemeTokenRegistry, error) {
	registry := NewThemeTokenRegistry()
	cleaned := stripCSSComments(src)

	// Cari setiap kemunculan blok @theme { ... }
	idx := 0
	for {
		themeIdx := bytes.Index(cleaned[idx:], []byte("@theme"))
		if themeIdx == -1 {
			break
		}
		themeStart := idx + themeIdx + len("@theme")

		// Cari kurung kurawal pembuka '{'
		openBrace := bytes.IndexByte(cleaned[themeStart:], '{')
		if openBrace == -1 {
			break
		}
		blockStart := themeStart + openBrace + 1

		// Cari kurung kurawal penutup '}' penyeimbang
		braceDepth := 1
		blockEnd := -1
		for i := blockStart; i < len(cleaned); i++ {
			if cleaned[i] == '{' {
				braceDepth++
			} else if cleaned[i] == '}' {
				braceDepth--
				if braceDepth == 0 {
					blockEnd = i
					break
				}
			}
		}

		if blockEnd == -1 {
			// Blok @theme tidak ditutup sempurna, ekstraksi hingga akhir berkas
			extractDeclarations(cleaned[blockStart:], registry)
			break
		}

		extractDeclarations(cleaned[blockStart:blockEnd], registry)
		idx = blockEnd + 1
	}

	return registry, nil
}

// stripCSSComments menghapus seluruh komentar /* ... */ dari buffer CSS.
func stripCSSComments(src []byte) []byte {
	var buf bytes.Buffer
	buf.Grow(len(src))

	i := 0
	for i < len(src) {
		if i+1 < len(src) && src[i] == '/' && src[i+1] == '*' {
			// Mulai komentar
			end := bytes.Index(src[i+2:], []byte("*/"))
			if end == -1 {
				// Komentar tidak ditutup, abaikan sisa berkas
				break
			}
			i += end + 4
			continue
		}
		buf.WriteByte(src[i])
		i++
	}

	return buf.Bytes()
}

// extractDeclarations memindai baris deklarasi properti kustom CSS (--var: val;).
func extractDeclarations(block []byte, registry *ThemeTokenRegistry) {
	lines := bytes.Split(block, []byte(";"))
	for _, rawLine := range lines {
		trimmed := strings.TrimSpace(string(rawLine))
		if trimmed == "" {
			continue
		}

		colonIdx := strings.IndexByte(trimmed, ':')
		if colonIdx == -1 {
			continue
		}

		key := strings.TrimSpace(trimmed[:colonIdx])
		val := strings.TrimSpace(trimmed[colonIdx+1:])

		// Hanya simpan deklarasi variabel CSS custom properties (--*)
		if strings.HasPrefix(key, "--") && val != "" {
			registry.Variables[key] = val
		}
	}
}
