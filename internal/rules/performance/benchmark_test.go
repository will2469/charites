package performance_test

import (
	"strings"
	"testing"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules"
)

type docProvider interface {
	Doc() ir.RuleDocumentation
}

// TestPerformanceRules_CanonicalContract memverifikasi bahwa seluruh rule kategori performance
// mematuhi antarmuka kanonikal ir.Rule dan menyediakan dokumentasi 8-Pillars lengkap.
func TestPerformanceRules_CanonicalContract(t *testing.T) {
	reg := rules.NewRegistry()
	if err := rules.RegisterBuiltinRules(reg); err != nil {
		t.Fatalf("failed to register builtin rules: %v", err)
	}

	for _, rule := range reg.All() {
		if rule.Category() != "performance" {
			continue
		}
		t.Run(rule.ID(), func(t *testing.T) {
			assertRuleBasics(t, rule)
			assertRuleDoc(t, rule)
		})
	}
}

func assertRuleBasics(t *testing.T, rule rules.Rule) {
	if !strings.HasPrefix(rule.ID(), "performance.") {
		t.Fatalf("rule ID %q must start with 'performance.'", rule.ID())
	}
	if rule.Description() == "" {
		t.Fatalf("rule Description cannot be empty")
	}
	if rule.Category() != "performance" {
		t.Fatalf("expected category 'performance', got %q", rule.Category())
	}
}

func assertRuleDoc(t *testing.T, rule rules.Rule) {
	docRule, ok := rule.(docProvider)
	if !ok {
		t.Fatalf("rule %s must implement Doc() method", rule.ID())
	}
	doc := docRule.Doc()
	if len(doc.TargetStandards) == 0 {
		t.Errorf("rule %s missing TargetStandards", rule.ID())
	}
	if doc.CoreInvariant == "" {
		t.Errorf("rule %s missing CoreInvariant", rule.ID())
	}
	if doc.Grounding == "" {
		t.Errorf("rule %s missing Grounding", rule.ID())
	}
	if len(doc.Risks) == 0 {
		t.Errorf("rule %s missing Risks", rule.ID())
	}
	if len(doc.BadExamples) == 0 {
		t.Errorf("rule %s missing BadExamples", rule.ID())
	}
	if len(doc.GoodExamples) == 0 {
		t.Errorf("rule %s missing GoodExamples", rule.ID())
	}
}

// BenchmarkPerformanceRules_ZeroAllocClean menguji seluruh rule kategori performance terhadap node bersih
// untuk memastikan garansi nol alokasi heap (0 B/op, 0 allocs/op).
func BenchmarkPerformanceRules_ZeroAllocClean(b *testing.B) {
	reg := rules.NewRegistry()
	if err := rules.RegisterBuiltinRules(reg); err != nil {
		b.Fatalf("failed to register builtin rules: %v", err)
	}

	cleanNode := &ir.Node{
		Type:       ir.NodeElement,
		Tag:        "div",
		Attributes: map[string]string{"className": "container mx-auto p-4", "id": "main"},
		Classes:    []string{"container", "mx-auto", "p-4"},
		RawClasses: "container mx-auto p-4",
		Span:       ir.Span{Line: 1, Column: 1},
	}

	for _, rule := range reg.All() {
		if rule.Category() != "performance" {
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
