package cls_test

import (
	"testing"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules"
)

// BenchmarkCLSRules_ZeroAllocClean menguji seluruh rule kategori cls terhadap node bersih
// untuk memastikan garansi nol alokasi heap (0 B/op, 0 allocs/op).
func BenchmarkCLSRules_ZeroAllocClean(b *testing.B) {
	reg := rules.NewRegistry()
	if err := rules.RegisterBuiltinRules(reg); err != nil {
		b.Fatalf("failed to register builtin rules: %v", err)
	}

	cleanNode := &ir.Node{
		Type:       ir.NodeElement,
		Tag:        "img",
		Attributes: map[string]string{"src": "/hero.jpg", "alt": "Hero", "width": "1200", "height": "600"},
		Classes:    []string{"w-full", "aspect-video", "object-cover", "rounded-xl"},
		RawClasses: "w-full aspect-video object-cover rounded-xl",
		Span:       ir.Span{Line: 1, Column: 1},
	}

	for _, rule := range reg.All() {
		if rule.Category() != "cls" {
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
