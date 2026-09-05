package tailwind

import (
	"reflect"
	"testing"
)

func TestTailwind_ThemeExtractor(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name:     "empty input",
			input:    "",
			expected: map[string]string{},
		},
		{
			name:     "no theme block",
			input:    "body { color: red; }",
			expected: map[string]string{},
		},
		{
			name: "basic theme block with variables",
			input: `
@theme {
  --color-primary: #2563eb;
  --color-primary-light: rgba(37, 99, 235, 0.1);
  --font-sans: Inter, sans-serif;
}
`,
			expected: map[string]string{
				"--color-primary":       "#2563eb",
				"--color-primary-light": "rgba(37, 99, 235, 0.1)",
				"--font-sans":           "Inter, sans-serif",
			},
		},
		{
			name: "oklch color space tokens",
			input: `
@theme {
  --color-accent: oklch(0.623 0.214 259.815);
  --color-accent-light: oklch(0.623 0.214 259.815 / 10%);
}
`,
			expected: map[string]string{
				"--color-accent":       "oklch(0.623 0.214 259.815)",
				"--color-accent-light": "oklch(0.623 0.214 259.815 / 10%)",
			},
		},
		{
			name: "theme with comments and non-variable declarations",
			input: `
/* Global styles */
@theme {
  /* Primary brand colors */
  --color-brand: #3b82f6; /* inline comment */
  ignored-prop: 10px;
  --spacing-sm: 0.5rem;
}
/* Trailing comment
`,
			expected: map[string]string{
				"--color-brand": "#3b82f6",
				"--spacing-sm":  "0.5rem",
			},
		},
		{
			name: "multiple @theme blocks",
			input: `
@theme {
  --color-a: #111;
}
@theme {
  --color-b: #222;
}
`,
			expected: map[string]string{
				"--color-a": "#111",
				"--color-b": "#222",
			},
		},
		{
			name: "unclosed @theme block recovery",
			input: `
@theme {
  --color-unclosed: #999;
`,
			expected: map[string]string{
				"--color-unclosed": "#999",
			},
		},
		{
			name: "nested braces inside comments or malformed lines",
			input: `
@theme {
  /* { fake block } */
  --color-test: #fff;
  empty-colon: ;
  invalid line without colon;
}
`,
			expected: map[string]string{
				"--color-test": "#fff",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg, err := ParseTheme([]byte(tt.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(reg.Variables, tt.expected) {
				t.Errorf("got %v, want %v", reg.Variables, tt.expected)
			}
		})
	}
}
