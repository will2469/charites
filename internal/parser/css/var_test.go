package css_test

import (
	"reflect"
	"testing"

	"github.com/will2469/charites/internal/parser/css"
)

func TestVar_BalancedParsing(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedCalls []css.VarCall
		expectedNames []string
	}{
		{
			name:  "User case: var with rgb fallback containing parentheses",
			input: "var(--foo, rgb(1, 2, 3))",
			expectedCalls: []css.VarCall{
				{
					Raw:         "var(--foo, rgb(1, 2, 3))",
					StartOffset: 0,
					EndOffset:   24,
					Name:        "--foo",
					Fallback:    "rgb(1, 2, 3)",
					HasFallback: true,
				},
			},
			expectedNames: []string{"--foo"},
		},
		{
			name:  "Nested var calls in fallback",
			input: "var(--primary, var(--fallback-color, #123456))",
			expectedCalls: []css.VarCall{
				{
					Raw:         "var(--primary, var(--fallback-color, #123456))",
					StartOffset: 0,
					EndOffset:   46,
					Name:        "--primary",
					Fallback:    "var(--fallback-color, #123456)",
					HasFallback: true,
				},
			},
			expectedNames: []string{"--primary", "--fallback-color"},
		},
		{
			name:  "Fallback with quotes and commas",
			input: `var(--font, "Helvetica, Arial", sans-serif)`,
			expectedCalls: []css.VarCall{
				{
					Raw:         `var(--font, "Helvetica, Arial", sans-serif)`,
					StartOffset: 0,
					EndOffset:   43,
					Name:        "--font",
					Fallback:    `"Helvetica, Arial", sans-serif`,
					HasFallback: true,
				},
			},
			expectedNames: []string{"--font"},
		},
		{
			name:  "Multiple vars in calc expression",
			input: "calc(var(--base-size, 16px) * var(--scale, 1.25))",
			expectedCalls: []css.VarCall{
				{
					Raw:         "var(--base-size, 16px)",
					StartOffset: 5,
					EndOffset:   27,
					Name:        "--base-size",
					Fallback:    "16px",
					HasFallback: true,
				},
				{
					Raw:         "var(--scale, 1.25)",
					StartOffset: 30,
					EndOffset:   48,
					Name:        "--scale",
					Fallback:    "1.25",
					HasFallback: true,
				},
			},
			expectedNames: []string{"--base-size", "--scale"},
		},
		{
			name:  "Escaped colon in var name",
			input: `var(--tw-bg-opacity\:1, 0.8)`,
			expectedCalls: []css.VarCall{
				{
					Raw:         `var(--tw-bg-opacity\:1, 0.8)`,
					StartOffset: 0,
					EndOffset:   28,
					Name:        "--tw-bg-opacity:1",
					Fallback:    "0.8",
					HasFallback: true,
				},
			},
			expectedNames: []string{"--tw-bg-opacity:1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := css.ExtractTopLevelVarCalls(tt.input)
			if len(calls) != len(tt.expectedCalls) {
				t.Fatalf("expected %d calls, got %d", len(tt.expectedCalls), len(calls))
			}
			for i, c := range calls {
				exp := tt.expectedCalls[i]
				if c.Raw != exp.Raw {
					t.Errorf("call[%d] Raw: expected %q, got %q", i, exp.Raw, c.Raw)
				}
				if c.StartOffset != exp.StartOffset || c.EndOffset != exp.EndOffset {
					t.Errorf("call[%d] offsets: expected [%d:%d], got [%d:%d]", i, exp.StartOffset, exp.EndOffset, c.StartOffset, c.EndOffset)
				}
				if c.Name != exp.Name {
					t.Errorf("call[%d] Name: expected %q, got %q", i, exp.Name, c.Name)
				}
				if c.Fallback != exp.Fallback {
					t.Errorf("call[%d] Fallback: expected %q, got %q", i, exp.Fallback, c.Fallback)
				}
				if c.HasFallback != exp.HasFallback {
					t.Errorf("call[%d] HasFallback: expected %v, got %v", i, exp.HasFallback, c.HasFallback)
				}
			}

			names := css.ExtractAllVarNames(tt.input)
			if !reflect.DeepEqual(names, tt.expectedNames) {
				t.Errorf("expected var names %v, got %v", tt.expectedNames, names)
			}
		})
	}
}
