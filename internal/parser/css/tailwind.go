package css

import (
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

// ParseTheme memindai konten CSS dan mengekstrak seluruh deklarasi variabel di dalam blok @theme { ... }
// menggunakan CSS Lexer dan AST Parser resmi Charites.
func ParseTheme(src []byte) (*ThemeTokenRegistry, error) {
	sheet, err := Parse(src)
	if err != nil {
		return nil, err
	}

	registry := NewThemeTokenRegistry()
	collectThemeRules(sheet.Rules, registry)
	return registry, nil
}

func collectThemeRules(rules []Rule, registry *ThemeTokenRegistry) {
	for _, r := range rules {
		at, ok := r.(*AtRule)
		if !ok {
			continue
		}
		if strings.HasPrefix(strings.ToLower(at.Name), "@theme") {
			for _, decl := range at.Declarations {
				if strings.HasPrefix(decl.Property, "--") && decl.Value != "" {
					registry.Variables[decl.Property] = decl.Value
				}
			}
		}
		collectThemeRules(at.Rules, registry)
	}
}
