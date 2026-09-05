package css_test

import (
	"testing"

	"github.com/will2469/charites/internal/parser/css"
)

func TestComputeSpecificity_W3CSpec(t *testing.T) {
	tests := []struct {
		name     string
		selector string
		expected css.Specificity
	}{
		// 1. Basic selectors & universals
		{
			name:     "Universal selector has zero specificity",
			selector: "*",
			expected: css.Specificity{A: 0, B: 0, C: 0},
		},
		{
			name:     "Type/element selector",
			selector: "div",
			expected: css.Specificity{A: 0, B: 0, C: 1},
		},
		{
			name:     "Class selector",
			selector: ".card",
			expected: css.Specificity{A: 0, B: 1, C: 0},
		},
		{
			name:     "ID selector",
			selector: "#header",
			expected: css.Specificity{A: 1, B: 0, C: 0},
		},
		{
			name:     "Compound selector: element + class + ID",
			selector: "div.card#header",
			expected: css.Specificity{A: 1, B: 1, C: 1},
		},
		{
			name:     "Complex combinator hierarchy",
			selector: "body > div.main p + span.icon",
			expected: css.Specificity{A: 0, B: 2, C: 4}, // body, div, p, span (C=4); .main, .icon (B=2)
		},

		// 2. Pseudo-classes vs Elements
		{
			name:     ":root is strictly a pseudo-class (not an element)",
			selector: ":root",
			expected: css.Specificity{A: 0, B: 1, C: 0},
		},
		{
			name:     "html.dark combines element and class",
			selector: "html.dark",
			expected: css.Specificity{A: 0, B: 1, C: 1},
		},
		{
			name:     "element with pseudo-class",
			selector: "a:hover",
			expected: css.Specificity{A: 0, B: 1, C: 1},
		},

		// 3. Pseudo-elements (double colon & legacy single colon)
		{
			name:     "Double-colon pseudo-element ::before",
			selector: "::before",
			expected: css.Specificity{A: 0, B: 0, C: 1},
		},
		{
			name:     "Double-colon pseudo-element ::placeholder",
			selector: "input::placeholder",
			expected: css.Specificity{A: 0, B: 0, C: 2}, // input + ::placeholder
		},
		{
			name:     "Legacy CSS2 single-colon pseudo-element :before",
			selector: "p:before",
			expected: css.Specificity{A: 0, B: 0, C: 2}, // p + :before
		},
		{
			name:     "Legacy CSS2 single-colon pseudo-element :after",
			selector: ":after",
			expected: css.Specificity{A: 0, B: 0, C: 1},
		},

		// 4. Attribute selectors with special characters & strings
		{
			name:     "Standard attribute selector",
			selector: `[data-theme="dark"]`,
			expected: css.Specificity{A: 0, B: 1, C: 0},
		},
		{
			name:     "Attribute with colons, dots, and hashes inside quote",
			selector: `[data-theme="#dark:root.active"]`,
			expected: css.Specificity{A: 0, B: 1, C: 0},
		},
		{
			name:     "Element combined with complex attribute",
			selector: `html[data-theme=":dark#123"]`,
			expected: css.Specificity{A: 0, B: 1, C: 1},
		},

		// 5. Escaped characters in identifiers
		{
			name:     "Class with escaped colon (Tailwind style)",
			selector: `.hover\:bg-primary`,
			expected: css.Specificity{A: 0, B: 1, C: 0},
		},
		{
			name:     "Class with escaped colon imitating :root",
			selector: `.\:root`,
			expected: css.Specificity{A: 0, B: 1, C: 0},
		},
		{
			name:     "Class with escaped hash imitating #id",
			selector: `.\#header`,
			expected: css.Specificity{A: 0, B: 1, C: 0},
		},

		// 6. Selectors Level 4 Functional Pseudo-Classes
		{
			name:     ":where() always has zero specificity",
			selector: ":where(#hero, .btn, div)",
			expected: css.Specificity{A: 0, B: 0, C: 0},
		},
		{
			name:     ":is() takes maximum specificity of arguments",
			selector: ":is(div, #main, .card)",
			expected: css.Specificity{A: 1, B: 0, C: 0}, // #main wins (1, 0, 0)
		},
		{
			name:     ":not() takes maximum specificity of arguments",
			selector: ":not(.active, #sidebar)",
			expected: css.Specificity{A: 1, B: 0, C: 0},
		},
		{
			name:     ":has() takes maximum specificity of arguments",
			selector: ":has(> .icon, span)",
			expected: css.Specificity{A: 0, B: 1, C: 0}, // .icon (0, 1, 0) > span (0, 0, 1)
		},
		{
			name:     "Complex selector with :where and :not",
			selector: "ul:where(#sidebar) li:not(.active)",
			expected: css.Specificity{A: 0, B: 1, C: 2}, // ul(C=1) + where(0) + li(C=1) + not(.active)(B=1)
		},
		{
			name:     ":nth-child with of clause",
			selector: ":nth-child(2n+1 of .active, #hero)",
			expected: css.Specificity{A: 1, B: 1, C: 0}, // (0, 1, 0) + max((0, 1, 0), (1, 0, 0)) = (1, 1, 0)
		},

		// 7. Nesting & Selector Lists
		{
			name:     "Nesting selector in isolation",
			selector: "&",
			expected: css.Specificity{A: 0, B: 1, C: 0},
		},
		{
			name:     "Comma-separated selector list takes maximum",
			selector: "h1, h2, #title",
			expected: css.Specificity{A: 1, B: 0, C: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := css.ComputeSpecificity(tt.selector)
			if !actual.Equals(tt.expected) {
				t.Errorf("selector %q: expected Specificity%+v, got %+v", tt.selector, tt.expected, actual)
			}
		})
	}
}
