package theme_test

import (
	"testing"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules/theme"
)

func TestHardcodeOpacityColorRule_Metadata(t *testing.T) {
	rule := theme.NewHardcodeOpacityColorRule()

	if rule.ID() != "theme.hardcode-opacity-color" {
		t.Errorf("expected ID 'theme.hardcode-opacity-color', got %q", rule.ID())
	}
	if rule.Category() != "theme" {
		t.Errorf("expected Category 'theme', got %q", rule.Category())
	}
	if rule.DefaultSeverity() != ir.SeverityError {
		t.Errorf("expected DefaultSeverity 'error', got %v", rule.DefaultSeverity())
	}
	if len(rule.Description()) == 0 {
		t.Errorf("expected non-empty description")
	}
}

func TestHardcodeOpacityColorRule_ReplacementFor(t *testing.T) {
	tests := []struct {
		token     string
		wantRep   string
		wantFound bool
	}{
		{"primary/10", "primary-light", true},
		{"primary/20", "primary-light", true},
		{"primary/5", "primary-subtle", true},
		{"destructive/10", "destructive-light", true},
		{"nonexistent/10", "", false},
		{"primary/30", "", false},
	}

	for _, tt := range tests {
		got, ok := theme.ReplacementFor(tt.token)
		if ok != tt.wantFound {
			t.Errorf("ReplacementFor(%q) ok = %v, want %v", tt.token, ok, tt.wantFound)
		}
		if got != tt.wantRep {
			t.Errorf("ReplacementFor(%q) = %q, want %q", tt.token, got, tt.wantRep)
		}
	}
}

func TestHardcodeOpacityColorRule_TableDrivenBoundary(t *testing.T) {
	rule := theme.NewHardcodeOpacityColorRule()

	tests := []struct {
		name           string
		node           *ir.Node
		wantCount      int
		wantClasses    []string
		wantHints      []string
		isViolation    bool
		classification string
	}{
		// 1. In-Scope Base Violations
		{
			name: "InScope_bg_primary_10",
			node: &ir.Node{
				Span:    ir.Span{Line: 10, Column: 5},
				Classes: []string{"bg-primary/10"},
			},
			wantCount:      1,
			wantClasses:    []string{"bg-primary/10"},
			wantHints:      []string{`Use semantic token "primary-light".`},
			isViolation:    true,
			classification: "In-Scope Base",
		},
		{
			name: "InScope_text_primary_20",
			node: &ir.Node{
				Span:    ir.Span{Line: 12, Column: 2},
				Classes: []string{"text-primary/20"},
			},
			wantCount:      1,
			wantClasses:    []string{"text-primary/20"},
			wantHints:      []string{`Use semantic token "primary-light".`},
			isViolation:    true,
			classification: "In-Scope Base",
		},
		{
			name: "InScope_border_destructive_20",
			node: &ir.Node{
				Span:    ir.Span{Line: 15, Column: 4},
				Classes: []string{"border-destructive/20"},
			},
			wantCount:      1,
			wantClasses:    []string{"border-destructive/20"},
			wantHints:      []string{`Use semantic token "destructive-light".`},
			isViolation:    true,
			classification: "In-Scope Base",
		},
		{
			name: "InScope_ring_warning_10",
			node: &ir.Node{
				Span:    ir.Span{Line: 18, Column: 8},
				Classes: []string{"ring-warning/10"},
			},
			wantCount:      1,
			wantClasses:    []string{"ring-warning/10"},
			wantHints:      []string{`Use semantic token "warning-light".`},
			isViolation:    true,
			classification: "In-Scope Base",
		},
		{
			name: "InScope_bg_primary_5",
			node: &ir.Node{
				Span:    ir.Span{Line: 20, Column: 1},
				Classes: []string{"bg-primary/5"},
			},
			wantCount:      1,
			wantClasses:    []string{"bg-primary/5"},
			wantHints:      []string{`Use semantic token "primary-subtle".`},
			isViolation:    true,
			classification: "In-Scope Base",
		},
		{
			name: "InScope_bg_secondary_10",
			node: &ir.Node{
				Span:    ir.Span{Line: 22, Column: 3},
				Classes: []string{"bg-secondary/10"},
			},
			wantCount:      1,
			wantClasses:    []string{"bg-secondary/10"},
			wantHints:      []string{`Use semantic token "muted-light".`},
			isViolation:    true,
			classification: "In-Scope Base",
		},
		{
			name: "InScope_bg_muted_5",
			node: &ir.Node{
				Span:    ir.Span{Line: 24, Column: 6},
				Classes: []string{"bg-muted/5"},
			},
			wantCount:      1,
			wantClasses:    []string{"bg-muted/5"},
			wantHints:      []string{`Use semantic token "muted-subtle".`},
			isViolation:    true,
			classification: "In-Scope Base",
		},
		{
			name: "InScope_text_accent_10",
			node: &ir.Node{
				Span:    ir.Span{Line: 26, Column: 1},
				Classes: []string{"text-accent/10"},
			},
			wantCount:      1,
			wantClasses:    []string{"text-accent/10"},
			wantHints:      []string{`Use semantic token "accent-light".`},
			isViolation:    true,
			classification: "In-Scope Base",
		},
		{
			name: "InScope_border_amber_10",
			node: &ir.Node{
				Span:    ir.Span{Line: 28, Column: 2},
				Classes: []string{"border-amber/10"},
			},
			wantCount:      1,
			wantClasses:    []string{"border-amber/10"},
			wantHints:      []string{`Use semantic token "amber-light".`},
			isViolation:    true,
			classification: "In-Scope Base",
		},
		{
			name: "InScope_ring_emerald_5",
			node: &ir.Node{
				Span:    ir.Span{Line: 30, Column: 4},
				Classes: []string{"ring-emerald/5"},
			},
			wantCount:      1,
			wantClasses:    []string{"ring-emerald/5"},
			wantHints:      []string{`Use semantic token "emerald-subtle".`},
			isViolation:    true,
			classification: "In-Scope Base",
		},

		// 2. In-Scope Variants (Single & Chained)
		{
			name: "InScope_hover_bg_primary_10",
			node: &ir.Node{
				Span:    ir.Span{Line: 35, Column: 2},
				Classes: []string{"hover:bg-primary/10"},
			},
			wantCount:      1,
			wantClasses:    []string{"hover:bg-primary/10"},
			wantHints:      []string{`Use semantic token "primary-light".`},
			isViolation:    true,
			classification: "In-Scope Variant",
		},
		{
			name: "InScope_dark_bg_primary_10",
			node: &ir.Node{
				Span:    ir.Span{Line: 37, Column: 5},
				Classes: []string{"dark:bg-primary/10"},
			},
			wantCount:      1,
			wantClasses:    []string{"dark:bg-primary/10"},
			wantHints:      []string{`Use semantic token "primary-light".`},
			isViolation:    true,
			classification: "In-Scope Variant",
		},
		{
			name: "InScope_md_hover_bg_primary_10",
			node: &ir.Node{
				Span:    ir.Span{Line: 40, Column: 1},
				Classes: []string{"md:hover:bg-primary/10"},
			},
			wantCount:      1,
			wantClasses:    []string{"md:hover:bg-primary/10"},
			wantHints:      []string{`Use semantic token "primary-light".`},
			isViolation:    true,
			classification: "In-Scope Chained Variant",
		},
		{
			name: "InScope_sm_dark_hover_border_destructive_20",
			node: &ir.Node{
				Span:    ir.Span{Line: 42, Column: 7},
				Classes: []string{"sm:dark:hover:border-destructive/20"},
			},
			wantCount:      1,
			wantClasses:    []string{"sm:dark:hover:border-destructive/20"},
			wantHints:      []string{`Use semantic token "destructive-light".`},
			isViolation:    true,
			classification: "In-Scope Multi-Chained Variant",
		},
		{
			name: "InScope_multiple_violations_single_node",
			node: &ir.Node{
				Span:    ir.Span{Line: 45, Column: 10},
				Classes: []string{"p-4", "bg-primary/10", "text-destructive/20", "flex"},
			},
			wantCount:      2,
			wantClasses:    []string{"bg-primary/10", "text-destructive/20"},
			wantHints:      []string{`Use semantic token "primary-light".`, `Use semantic token "destructive-light".`},
			isViolation:    true,
			classification: "In-Scope Multiple",
		},

		// 3. Clean Negatives (Valid Semantic Tokens & Standard Classes)
		{
			name: "Clean_bg_primary_light",
			node: &ir.Node{
				Span:    ir.Span{Line: 50, Column: 1},
				Classes: []string{"bg-primary-light"},
			},
			wantCount:      0,
			isViolation:    false,
			classification: "Clean Semantic Token",
		},
		{
			name: "Clean_bg_primary_subtle",
			node: &ir.Node{
				Span:    ir.Span{Line: 52, Column: 1},
				Classes: []string{"bg-primary-subtle"},
			},
			wantCount:      0,
			isViolation:    false,
			classification: "Clean Semantic Token",
		},
		{
			name: "Clean_text_muted",
			node: &ir.Node{
				Span:    ir.Span{Line: 54, Column: 1},
				Classes: []string{"text-muted", "text-muted-light"},
			},
			wantCount:      0,
			isViolation:    false,
			classification: "Clean Semantic Token",
		},
		{
			name: "Clean_standard_layout_classes",
			node: &ir.Node{
				Span:    ir.Span{Line: 56, Column: 1},
				Classes: []string{"flex", "items-center", "justify-between", "p-4", "rounded-lg"},
			},
			wantCount:      0,
			isViolation:    false,
			classification: "Clean Standard Utilities",
		},

		// 4. Out-of-Scope Baits (Must produce 0 diagnostics)
		{
			name: "OutOfScope_layout_width_fraction",
			node: &ir.Node{
				Span:    ir.Span{Line: 60, Column: 1},
				Classes: []string{"w-1/2", "w-2/3", "max-w-1/2"},
			},
			wantCount:      0,
			isViolation:    false,
			classification: "Out-of-Scope Layout Fraction",
		},
		{
			name: "OutOfScope_layout_height_fraction",
			node: &ir.Node{
				Span:    ir.Span{Line: 62, Column: 1},
				Classes: []string{"h-1/3", "h-2/4"},
			},
			wantCount:      0,
			isViolation:    false,
			classification: "Out-of-Scope Layout Fraction",
		},
		{
			name: "OutOfScope_aspect_ratios",
			node: &ir.Node{
				Span:    ir.Span{Line: 64, Column: 1},
				Classes: []string{"aspect-16/9", "aspect-4/3", "aspect-1/1"},
			},
			wantCount:      0,
			isViolation:    false,
			classification: "Out-of-Scope Aspect Ratio",
		},
		{
			name: "OutOfScope_grid_cols_fraction",
			node: &ir.Node{
				Span:    ir.Span{Line: 66, Column: 1},
				Classes: []string{"grid-cols-2/3"},
			},
			wantCount:      0,
			isViolation:    false,
			classification: "Out-of-Scope Grid Fraction",
		},
		{
			name: "OutOfScope_line_height_modifiers",
			node: &ir.Node{
				Span:    ir.Span{Line: 68, Column: 1},
				Classes: []string{"text-sm/6", "text-xs/relaxed", "text-base/7"},
			},
			wantCount:      0,
			isViolation:    false,
			classification: "Out-of-Scope Line-Height",
		},
		{
			name: "OutOfScope_unmapped_opacities",
			node: &ir.Node{
				Span:    ir.Span{Line: 70, Column: 1},
				Classes: []string{"bg-primary/30", "bg-primary/50", "bg-primary/100", "bg-primary/[0.1]"},
			},
			wantCount:      0,
			isViolation:    false,
			classification: "Out-of-Scope Unmapped Opacity",
		},
		{
			name: "OutOfScope_arbitrary_hex_color",
			node: &ir.Node{
				Span:    ir.Span{Line: 72, Column: 1},
				Classes: []string{"bg-[#123456]/10", "text-[#ff0000]/20", "border-[#abcdef]/5"},
			},
			wantCount:      0,
			isViolation:    false,
			classification: "Out-of-Scope Arbitrary Hex Color",
		},
		{
			name: "OutOfScope_raw_palette_color",
			node: &ir.Node{
				Span:    ir.Span{Line: 74, Column: 1},
				Classes: []string{"bg-red-500/10", "text-blue-600/20", "border-gray-300/5"},
			},
			wantCount:      0,
			isViolation:    false,
			classification: "Out-of-Scope Raw Palette Color",
		},
		{
			name: "OutOfScope_raw_black_white_color",
			node: &ir.Node{
				Span:    ir.Span{Line: 76, Column: 1},
				Classes: []string{"bg-black/10", "text-white/20"},
			},
			wantCount:      0,
			isViolation:    false,
			classification: "Out-of-Scope Raw Black/White",
		},

		// 5. Structural Edge Cases
		{
			name:           "EdgeCase_nil_node",
			node:           nil,
			wantCount:      0,
			isViolation:    false,
			classification: "Edge Case Nil",
		},
		{
			name: "EdgeCase_empty_classes",
			node: &ir.Node{
				Span:    ir.Span{Line: 80, Column: 1},
				Classes: []string{},
			},
			wantCount:      0,
			isViolation:    false,
			classification: "Edge Case Empty",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			diags := rule.Evaluate(tc.node)

			if len(diags) != tc.wantCount {
				t.Fatalf("[%s] expected %d diagnostics, got %d: %+v", tc.classification, tc.wantCount, len(diags), diags)
			}

			if !tc.isViolation {
				if diags != nil {
					t.Errorf("[%s] expected nil slice on non-violation, got allocated slice: %+v", tc.classification, diags)
				}
				return
			}

			for i, d := range diags {
				if d.Rule != rule.ID() {
					t.Errorf("expected Rule %q, got %q", rule.ID(), d.Rule)
				}
				if d.Severity != rule.DefaultSeverity() {
					t.Errorf("expected Severity %v, got %v", rule.DefaultSeverity(), d.Severity)
				}
				if d.Line != tc.node.Span.Line {
					t.Errorf("expected Line %d, got %d", tc.node.Span.Line, d.Line)
				}
				if d.Column != tc.node.Span.Column {
					t.Errorf("expected Column %d, got %d", tc.node.Span.Column, d.Column)
				}

				expectedMsg := `Hardcode opacity color: "` + tc.wantClasses[i] + `"`
				if d.Message != expectedMsg {
					t.Errorf("diag %d: expected Message %q, got %q", i, expectedMsg, d.Message)
				}

				if d.Hint != tc.wantHints[i] {
					t.Errorf("diag %d: expected Hint %q, got %q", i, tc.wantHints[i], d.Hint)
				}
			}
		})
	}
}

func BenchmarkEvaluateHardcodeOpacityColor_Clean(b *testing.B) {
	rule := theme.NewHardcodeOpacityColorRule()
	node := &ir.Node{
		Tag:     "div",
		Classes: []string{"p-4", "flex", "items-center", "justify-between", "rounded-lg"},
	}

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = rule.Evaluate(node)
	}
}

func BenchmarkEvaluateHardcodeOpacityColor_Violation(b *testing.B) {
	rule := theme.NewHardcodeOpacityColorRule()
	node := &ir.Node{
		Span:    ir.Span{Line: 1, Column: 1},
		Tag:     "div",
		Classes: []string{"p-4", "flex", "bg-primary/10", "rounded-lg"},
	}

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = rule.Evaluate(node)
	}
}
