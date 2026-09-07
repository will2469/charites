package design_test

import (
	"strings"
	"testing"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules/design"
)

func TestComponentStyleDriftRule_LocalEvaluate(t *testing.T) {
	rule := design.NewComponentStyleDriftRule()

	t.Run("ContrastHazard_Triggered", func(t *testing.T) {
		node := &ir.Node{
			Type:       ir.NodeElement,
			Tag:        "Button",
			RawClasses: "bg-yellow-400 text-white rounded-md",
			Classes:    []string{"bg-yellow-400", "text-white", "rounded-md"},
			Span:       ir.Span{Line: 12, Column: 4},
		}

		diags := rule.Evaluate(node)
		if len(diags) != 1 {
			t.Fatalf("expected 1 diagnostic for contrast hazard, got %d", len(diags))
		}
		if diags[0].Rule != "design.component-style-drift" {
			t.Errorf("expected rule design.component-style-drift, got %s", diags[0].Rule)
		}
		if !strings.Contains(diags[0].Message, "Contrast Hazard") {
			t.Errorf("expected message to mention Contrast Hazard, got: %s", diags[0].Message)
		}
	})

	t.Run("CleanAccessible_NoViolation", func(t *testing.T) {
		node := &ir.Node{
			Type:       ir.NodeElement,
			Tag:        "Button",
			RawClasses: "bg-black text-white rounded-md",
			Classes:    []string{"bg-black", "text-white", "rounded-md"},
			Span:       ir.Span{Line: 15, Column: 4},
		}

		diags := rule.Evaluate(node)
		if len(diags) != 0 {
			t.Errorf("expected 0 diagnostics for accessible button, got %d", len(diags))
		}
	})

	t.Run("NonInteractiveTag_Ignored", func(t *testing.T) {
		node := &ir.Node{
			Type:       ir.NodeElement,
			Tag:        "div",
			RawClasses: "bg-yellow-400 text-white",
			Classes:    []string{"bg-yellow-400", "text-white"},
			Span:       ir.Span{Line: 20, Column: 4},
		}

		diags := rule.Evaluate(node)
		if len(diags) != 0 {
			t.Errorf("expected 0 diagnostics for generic div, got %d", len(diags))
		}
	})
}
