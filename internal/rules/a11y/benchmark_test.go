package a11y_test

import (
	"testing"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules"
)

// BenchmarkA11yRules_ZeroAllocClean menguji seluruh rule kategori a11y terhadap node bersih
// untuk memastikan garansi nol alokasi heap (0 B/op, 0 allocs/op) sesuai persyaratan QUAL-03.
func BenchmarkA11yRules_ZeroAllocClean(b *testing.B) {
	reg := rules.NewRegistry()
	if err := rules.RegisterBuiltinRules(reg); err != nil {
		b.Fatalf("failed to register builtin rules: %v", err)
	}

	cleanNode := &ir.Node{
		Type:    ir.NodeElement,
		Tag:     "button",
		Classes: []string{"h-11", "w-11", "min-h-[44px]", "min-w-[44px]", "p-2.5", "focus-visible:ring-2", "text-base"},
		Span:    ir.Span{Line: 1, Column: 1},
	}

	for _, rule := range reg.All() {
		if rule.Category() != "a11y" {
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
