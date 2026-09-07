package theme_test

import (
	"testing"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules/theme"
)

func TestImportantOverride_ShadowRegression(t *testing.T) {
	rule := theme.NewImportantOverrideRule()

	tests := []struct {
		name      string
		classes   []string
		wantDiags int
	}{
		{
			name:      "Shadow_Elevations_Ignored",
			classes:   []string{"!shadow-sm", "!shadow-md", "!shadow-lg", "!shadow-inner", "!shadow-none", "!shadow-sm/20"},
			wantDiags: 0,
		},
		{
			name:      "Shadow_SemanticColor_Flagged",
			classes:   []string{"!shadow-primary"},
			wantDiags: 1,
		},
		{
			name:      "Shadow_SemanticColor_WithOpacity_Flagged",
			classes:   []string{"!shadow-primary/20"},
			wantDiags: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &ir.Node{
				Span:    ir.Span{Line: 10, Column: 1},
				Classes: tt.classes,
			}
			diags := rule.Evaluate(node)
			if len(diags) != tt.wantDiags {
				t.Fatalf("expected %d diagnostics, got %d: %+v", tt.wantDiags, len(diags), diags)
			}
		})
	}
}
