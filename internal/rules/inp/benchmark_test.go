package inp_test

import (
	"testing"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules"
)

// BenchmarkINPRules_ZeroAllocClean menguji seluruh rule kategori inp terhadap node bersih
// untuk memastikan garansi nol alokasi heap (0 B/op, 0 allocs/op).
func BenchmarkINPRules_ZeroAllocClean(b *testing.B) {
	reg := rules.NewRegistry()
	if err := rules.RegisterBuiltinRules(reg); err != nil {
		b.Fatalf("failed to register builtin rules: %v", err)
	}

	cleanNode := &ir.Node{
		Type:       ir.NodeElement,
		Tag:        "button",
		Attributes: map[string]string{"type": "button", "onClick": "() => handleAction()"},
		Classes:    []string{"px-4", "py-2", "bg-primary", "text-white", "rounded-lg"},
		RawClasses: "px-4 py-2 bg-primary text-white rounded-lg",
		Span:       ir.Span{Line: 1, Column: 1},
	}

	for _, rule := range reg.All() {
		if rule.Category() != "inp" {
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
