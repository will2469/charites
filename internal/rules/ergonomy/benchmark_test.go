package ergonomy_test

import (
	"testing"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules"
)

// BenchmarkErgonomyRules_ZeroAllocClean menguji seluruh rule kategori ergonomy terhadap node bersih
// untuk memastikan garansi nol alokasi heap (0 B/op, 0 allocs/op) sesuai persyaratan QUAL-03.
func BenchmarkErgonomyRules_ZeroAllocClean(b *testing.B) {
	reg := rules.NewRegistry()
	if err := rules.RegisterBuiltinRules(reg); err != nil {
		b.Fatalf("failed to register builtin rules: %v", err)
	}

	cleanNode := &ir.Node{
		Type:       ir.NodeElement,
		Tag:        "button",
		Attributes: map[string]string{"type": "button"},
		Classes:    []string{"h-11", "px-4", "bg-primary", "text-primary-foreground", "rounded-xl", "active:scale-95"},
		RawClasses: "h-11 px-4 bg-primary text-primary-foreground rounded-xl active:scale-95",
		Span:       ir.Span{Line: 1, Column: 1},
	}

	for _, rule := range reg.All() {
		if rule.Category() != "ergonomy" {
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
