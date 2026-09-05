package analyzer_test

import (
	"fmt"
	"testing"

	"github.com/will2469/charites/internal/analyzer"
	"github.com/will2469/charites/internal/config"
	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules/theme"
)

func BenchmarkAnalyzer_EngineTraversal(b *testing.B) {
	// Bangun pohon IR in-memory dengan 100 node
	root := &ir.Node{
		Type: ir.NodeElement,
		Tag:  "div",
		Span: ir.Span{Line: 1, Column: 1, EndLine: 100, EndColumn: 10},
	}

	curr := root
	for i := 2; i <= 100; i++ {
		classes := []string{"flex", "items-center", "p-4"}
		if i%5 == 0 {
			classes = append(classes, "bg-primary/10")
		}
		child := &ir.Node{
			Type:       ir.NodeElement,
			Tag:        "span",
			Classes:    classes,
			RawClasses: fmt.Sprintf("span-%d", i),
			Span:       ir.Span{Line: i, Column: 1, EndLine: i, EndColumn: 20},
		}
		curr.Children = append(curr.Children, child)
	}

	rule := theme.NewHardcodeOpacityColorRule()
	activeRules := []config.ActiveRule{
		{Rule: rule, EffectiveSeverity: ir.SeverityError},
	}
	eng := analyzer.NewEngine(activeRules)

	b.ReportAllocs()

	for b.Loop() {
		_ = eng.AnalyzeTree("test.tsx", root, nil)
	}
}
