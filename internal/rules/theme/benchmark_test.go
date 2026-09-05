package theme_test

import (
	"testing"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules"
)

// BenchmarkThemeRules_ZeroAllocClean menguji seluruh rule kategori theme terhadap node bersih
// untuk memastikan garansi nol alokasi heap (0 B/op, 0 allocs/op) sesuai persyaratan QUAL-03.
func BenchmarkThemeRules_ZeroAllocClean(b *testing.B) {
	reg := rules.NewRegistry()
	if err := rules.RegisterBuiltinRules(reg); err != nil {
		b.Fatalf("failed to register builtin rules: %v", err)
	}

	cleanNode := &ir.Node{
		Tag:     "div",
		Classes: []string{"bg-card", "border", "border-border", "text-card-foreground", "p-6", "rounded-xl", "shadow-sm"},
		Span:    ir.Span{Line: 1, Column: 1},
	}

	for _, rule := range reg.All() {
		if rule.Category() != "theme" {
			continue
		}
		b.Run(rule.ID(), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = rule.Evaluate(cleanNode)
			}
		})
	}
}
