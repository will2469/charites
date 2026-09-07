package theme_test

import (
	"strings"
	"testing"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules/theme"
)

func TestHardcodeSizeRule_TableDrivenBoundary(t *testing.T) {
	rule := theme.NewHardcodeSizeRule()

	const wantScaleHint = "Avoid arbitrary inline scale modifiers. Define a standardized scale token/variable in global.css (e.g. --scale-press) or use standard Tailwind scale steps (scale-95, scale-105)."
	const wantFractionHint = "Use a standard integer step (e.g. p-3, p-4) or official half-step (e.g. p-3.5) instead of an off-grid decimal."
	const wantScalarHint = "Use a standard modular spacing step (e.g. p-4, p-5) or a semantic design token."

	tests := []struct {
		name         string
		classes      []string
		wantDiags    int
		wantMsgPat   string
		wantHintText string
	}{
		// 1. In-Scope: Arbitrary scale variants (Issue #3)
		{
			name:         "InScope_scale_uniform",
			classes:      []string{"scale-[0.99]"},
			wantDiags:    1,
			wantMsgPat:   `Arbitrary inline scale modifier: "scale-[0.99]"`,
			wantHintText: wantScaleHint,
		},
		{
			name:         "InScope_negative_scale_uniform",
			classes:      []string{"-scale-[0.99]"},
			wantDiags:    1,
			wantMsgPat:   `Arbitrary inline scale modifier: "-scale-[0.99]"`,
			wantHintText: wantScaleHint,
		},
		{
			name:         "InScope_scale_directional_x",
			classes:      []string{"scale-x-[0.95]", "-scale-x-[0.95]"},
			wantDiags:    2,
			wantMsgPat:   `Arbitrary inline scale modifier`,
			wantHintText: wantScaleHint,
		},
		{
			name:         "InScope_scale_directional_y",
			classes:      []string{"scale-y-[1.05]", "-scale-y-[1.05]"},
			wantDiags:    2,
			wantMsgPat:   `Arbitrary inline scale modifier`,
			wantHintText: wantScaleHint,
		},
		{
			name:         "InScope_scale_directional_z",
			classes:      []string{"scale-z-[1.05]", "-scale-z-[1.05]"},
			wantDiags:    2,
			wantMsgPat:   `Arbitrary inline scale modifier`,
			wantHintText: wantScaleHint,
		},
		{
			name:         "InScope_scale_percentage",
			classes:      []string{"scale-[98%]", "scale-[105%]"},
			wantDiags:    2,
			wantMsgPat:   `Arbitrary inline scale modifier`,
			wantHintText: wantScaleHint,
		},
		{
			name:         "InScope_scale_variants",
			classes:      []string{"active:scale-[0.99]", "hover:scale-[1.02]", "dark:active:-scale-[0.95]"},
			wantDiags:    3,
			wantMsgPat:   `Arbitrary inline scale modifier`,
			wantHintText: wantScaleHint,
		},
		{
			name:         "InScope_arbitrary_property_scale",
			classes:      []string{"[scale:0.98]", "[scale-x:0.98]", "[scale-y:1.05]", "[scale-z:1.05]"},
			wantDiags:    4,
			wantMsgPat:   `Arbitrary inline scale modifier`,
			wantHintText: wantScaleHint,
		},

		// 2. Out-of-Scope: Standard Tailwind scale steps (Clean)
		{
			name:      "OutOfScope_standard_scale_steps",
			classes:   []string{"scale-0", "scale-50", "scale-75", "scale-90", "scale-95", "scale-100", "scale-105", "scale-110", "scale-125", "scale-150"},
			wantDiags: 0,
		},
		{
			name:      "OutOfScope_standard_scale_variants",
			classes:   []string{"active:scale-95", "hover:scale-105", "focus:scale-100"},
			wantDiags: 0,
		},

		// 3. Out-of-Scope: Token-backed CSS variables (Clean)
		{
			name:      "OutOfScope_css_variable_scale",
			classes:   []string{"scale-[var(--scale-press)]", "active:scale-[var(--scale-press)]", "[scale:var(--scale-press)]", "hover:scale-[var(--scale-hover)]"},
			wantDiags: 0,
		},

		// 4. Out-of-Scope: Non-numeric / calc scale (Not owned by numeric scale classifier)
		{
			name:      "OutOfScope_non_numeric_scale",
			classes:   []string{"scale-[calc(1+2)]", "scale-[abc]", "scale-[theme(spacing.4)]"},
			wantDiags: 0,
		},

		// 5. Out-of-Scope: Strict non-ownership boundary (rotate, duration, opacity, skew)
		{
			name:      "OutOfScope_non_scale_properties",
			classes:   []string{"rotate-[45deg]", "duration-[300ms]", "opacity-[0.85]", "skew-x-[10deg]"},
			wantDiags: 0,
		},

		// 6. In-Scope: Existing spatial scalars & non-standard fractions
		{
			name:         "InScope_spatial_scalar",
			classes:      []string{"p-[19px]", "w-[320px]"},
			wantDiags:    2,
			wantMsgPat:   `Hardcoded size/spacing scalar`,
			wantHintText: wantScalarHint,
		},
		{
			name:         "InScope_non_standard_fraction",
			classes:      []string{"p-3.25", "w-2.75"},
			wantDiags:    2,
			wantMsgPat:   `Non-standard fractional scale`,
			wantHintText: wantFractionHint,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &ir.Node{
				Tag:     "div",
				Classes: tt.classes,
				Span:    ir.Span{Line: 1, Column: 1},
			}

			diags := rule.Evaluate(node)
			if len(diags) != tt.wantDiags {
				t.Fatalf("expected %d diagnostics, got %d: %+v", tt.wantDiags, len(diags), diags)
			}

			if tt.wantDiags > 0 {
				for _, d := range diags {
					if !strings.Contains(d.Message, tt.wantMsgPat) {
						t.Errorf("diagnostic message %q does not contain expected pattern %q", d.Message, tt.wantMsgPat)
					}
					if d.Hint != tt.wantHintText {
						t.Errorf("diagnostic hint %q does not match expected hint %q", d.Hint, tt.wantHintText)
					}
				}
			}
		})
	}
}

// BenchmarkHardcodeSizeRule_HardenedClean tests that a clean node containing a realistic mix of
// standard scale steps, CSS-var scale, and unowned arbitrary properties produces 0 heap allocations.
func BenchmarkHardcodeSizeRule_HardenedClean(b *testing.B) {
	rule := theme.NewHardcodeSizeRule()
	cleanNode := &ir.Node{
		Tag: "div",
		Classes: []string{
			"p-4", "w-80", "h-11", "gap-3", "text-base",
			"scale-95", "scale-100", "scale-105", "active:scale-95",
			"scale-[var(--scale-press)]", "[scale:var(--scale-press)]",
			"rotate-[45deg]", "duration-[300ms]", "opacity-[0.85]",
		},
		Span: ir.Span{Line: 1, Column: 1},
	}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		diags := rule.Evaluate(cleanNode)
		if len(diags) > 0 {
			b.Fatalf("unexpected diagnostics on clean node: %+v", diags)
		}
	}
}
