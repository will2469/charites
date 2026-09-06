package ux_test

import (
	"testing"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules"
)

// BenchmarkUXRules_ZeroAllocClean menguji seluruh rule kategori ux terhadap node bersih
// untuk memastikan garansi nol alokasi heap (0 B/op, 0 allocs/op) sesuai persyaratan QUAL-03.
func BenchmarkUXRules_ZeroAllocClean(b *testing.B) {
	reg := rules.NewRegistry()
	if err := rules.RegisterBuiltinRules(reg); err != nil {
		b.Fatalf("failed to register builtin rules: %v", err)
	}

	cleanNode := &ir.Node{
		Type:       ir.NodeElement,
		Tag:        "div",
		Attributes: map[string]string{"role": "region"},
		Classes:    []string{"flex", "flex-col", "gap-6", "p-4", "bg-card", "text-card-foreground"},
		RawClasses: "flex flex-col gap-6 p-4 bg-card text-card-foreground",
		Span:       ir.Span{Line: 1, Column: 1},
	}

	for _, rule := range reg.All() {
		if rule.Category() != "ux" {
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
