package css_test

import (
	"testing"

	"github.com/will2469/charites/internal/parser/css"
)

func TestCSS_ParserRobustness(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		checkFunc func(t *testing.T, sheet *css.StyleSheet)
	}{
		{
			name: "Semicolon inside quotes",
			input: `.foo {
  --value: "hello;world";
}`,
			checkFunc: func(t *testing.T, sheet *css.StyleSheet) {
				if len(sheet.Rules) != 1 {
					t.Fatalf("expected 1 rule, got %d", len(sheet.Rules))
				}
				style, ok := sheet.Rules[0].(*css.StyleRule)
				if !ok {
					t.Fatalf("expected StyleRule")
				}
				if len(style.Declarations) != 1 {
					t.Fatalf("expected 1 declaration, got %d", len(style.Declarations))
				}
				decl := style.Declarations[0]
				if decl.Property != "--value" {
					t.Errorf("expected property '--value', got %q", decl.Property)
				}
				if decl.Value != `"hello;world"` {
					t.Errorf("expected value '\"hello;world\"', got %q", decl.Value)
				}
			},
		},
		{
			name: "Data URI with colons and semicolons",
			input: `.foo {
  --image: url("data:image/svg+xml;utf8,<svg>hello</svg>");
}`,
			checkFunc: func(t *testing.T, sheet *css.StyleSheet) {
				style := sheet.Rules[0].(*css.StyleRule)
				if len(style.Declarations) != 1 {
					t.Fatalf("expected 1 declaration, got %d", len(style.Declarations))
				}
				decl := style.Declarations[0]
				if decl.Property != "--image" {
					t.Errorf("expected property '--image', got %q", decl.Property)
				}
				if decl.Value != `url("data:image/svg+xml;utf8,<svg>hello</svg>")` {
					t.Errorf("expected verbatim url value, got %q", decl.Value)
				}
			},
		},
		{
			name: "CSS Nesting & Pseudo-selectors",
			input: `.card {
  --surface: red;
  &:hover {
    --surface: blue;
  }
}`,
			checkFunc: func(t *testing.T, sheet *css.StyleSheet) {
				style := sheet.Rules[0].(*css.StyleRule)
				if style.Selector != ".card" {
					t.Errorf("expected selector '.card', got %q", style.Selector)
				}
				if len(style.Declarations) != 1 || style.Declarations[0].Value != "red" {
					t.Errorf("expected surface: red")
				}
				if len(style.Rules) != 1 {
					t.Fatalf("expected 1 nested rule, got %d", len(style.Rules))
				}
				nested := style.Rules[0].(*css.StyleRule)
				if nested.Selector != "&:hover" {
					t.Errorf("expected nested selector '&:hover', got %q", nested.Selector)
				}
				if len(nested.Declarations) != 1 || nested.Declarations[0].Value != "blue" {
					t.Errorf("expected nested surface: blue")
				}
			},
		},
		{
			name: "Nested AtRules: @layer theme inside @media",
			input: `@layer theme {
  @media (prefers-color-scheme: dark) {
    :root {
      --brand: oklch(0.65 0.20 260);
    }
  }
}`,
			checkFunc: func(t *testing.T, sheet *css.StyleSheet) {
				if len(sheet.Rules) != 1 {
					t.Fatalf("expected 1 top-level rule, got %d", len(sheet.Rules))
				}
				layer := sheet.Rules[0].(*css.AtRule)
				if layer.Name != "@layer" || layer.Prelude != "theme" {
					t.Errorf("expected @layer theme, got %s %s", layer.Name, layer.Prelude)
				}
				if len(layer.Rules) != 1 {
					t.Fatalf("expected 1 nested rule in layer, got %d", len(layer.Rules))
				}
				media := layer.Rules[0].(*css.AtRule)
				if media.Name != "@media" || media.Prelude != "(prefers-color-scheme: dark)" {
					t.Errorf("expected @media query, got %s %s", media.Name, media.Prelude)
				}
				if len(media.Rules) != 1 {
					t.Fatalf("expected 1 nested rule in media, got %d", len(media.Rules))
				}
				root := media.Rules[0].(*css.StyleRule)
				if root.Selector != ":root" {
					t.Errorf("expected :root, got %q", root.Selector)
				}
				if len(root.Declarations) != 1 || root.Declarations[0].Property != "--brand" {
					t.Errorf("expected --brand declaration")
				}
			},
		},
		{
			name: "Banana test: arbitrary tokens without opinion",
			input: `:root {
  --banana: #123456;
  --thing-that-is-definitely-not-primary: red;
  --super-special-design-token: var(--banana);
}`,
			checkFunc: func(t *testing.T, sheet *css.StyleSheet) {
				root := sheet.Rules[0].(*css.StyleRule)
				if len(root.Declarations) != 3 {
					t.Fatalf("expected 3 declarations, got %d", len(root.Declarations))
				}
				if root.Declarations[0].Property != "--banana" || root.Declarations[0].Value != "#123456" {
					t.Errorf("failed banana declaration: %+v", root.Declarations[0])
				}
				if root.Declarations[1].Property != "--thing-that-is-definitely-not-primary" {
					t.Errorf("failed second declaration: %+v", root.Declarations[1])
				}
				if root.Declarations[2].Property != "--super-special-design-token" || root.Declarations[2].Value != "var(--banana)" {
					t.Errorf("failed var reference declaration: %+v", root.Declarations[2])
				}
			},
		},
		{
			name: "Semicolons inside parentheses and brackets do not terminate declaration",
			input: `:root {
  --func-val: fn(first; second);
  --bracket-val: [col-1; col-2];
  --clean: #ffffff;
}`,
			checkFunc: func(t *testing.T, sheet *css.StyleSheet) {
				root := sheet.Rules[0].(*css.StyleRule)
				if len(root.Declarations) != 3 {
					t.Fatalf("expected 3 declarations, got %d", len(root.Declarations))
				}
				if root.Declarations[0].Property != "--func-val" || root.Declarations[0].Value != "fn(first; second)" {
					t.Errorf("expected fn(first; second), got %+v", root.Declarations[0])
				}
				if root.Declarations[1].Property != "--bracket-val" || root.Declarations[1].Value != "[col-1; col-2]" {
					t.Errorf("expected [col-1; col-2], got %+v", root.Declarations[1])
				}
				if root.Declarations[2].Property != "--clean" || root.Declarations[2].Value != "#ffffff" {
					t.Errorf("expected #ffffff, got %+v", root.Declarations[2])
				}
			},
		},
		{
			name: "Pseudo-class selectors and colon distinction",
			input: `a:hover {
  color: red;
}
:is(.btn, .link):focus {
  outline: 2px;
}`,
			checkFunc: func(t *testing.T, sheet *css.StyleSheet) {
				if len(sheet.Rules) != 2 {
					t.Fatalf("expected 2 rules, got %d", len(sheet.Rules))
				}
				rule1, ok1 := sheet.Rules[0].(*css.StyleRule)
				if !ok1 || rule1.Selector != "a:hover" {
					t.Errorf("expected style rule with selector 'a:hover', got %+v", sheet.Rules[0])
				}
				rule2, ok2 := sheet.Rules[1].(*css.StyleRule)
				if !ok2 || rule2.Selector != ":is(.btn, .link):focus" {
					t.Errorf("expected style rule with selector ':is(.btn, .link):focus', got %+v", sheet.Rules[1])
				}
			},
		},
		{
			name: "Final declaration without semicolon before closing brace",
			input: `.box {
  --pad: 10px;
  --margin: 20px
}`,
			checkFunc: func(t *testing.T, sheet *css.StyleSheet) {
				rule := sheet.Rules[0].(*css.StyleRule)
				if len(rule.Declarations) != 2 {
					t.Fatalf("expected 2 declarations, got %d", len(rule.Declarations))
				}
				if rule.Declarations[1].Property != "--margin" || rule.Declarations[1].Value != "20px" {
					t.Errorf("expected --margin: 20px, got %+v", rule.Declarations[1])
				}
			},
		},
		{
			name: "Multiple empty semicolons",
			input: `.box {
  ;; --a: 1; ; ; --b: 2;;
}`,
			checkFunc: func(t *testing.T, sheet *css.StyleSheet) {
				rule := sheet.Rules[0].(*css.StyleRule)
				if len(rule.Declarations) != 2 {
					t.Fatalf("expected 2 declarations, got %d", len(rule.Declarations))
				}
				if rule.Declarations[0].Property != "--a" || rule.Declarations[1].Property != "--b" {
					t.Errorf("unexpected declarations: %+v", rule.Declarations)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sheet, err := css.Parse([]byte(tt.input))
			if err != nil {
				t.Fatalf("unexpected parse error: %v", err)
			}
			tt.checkFunc(t, sheet)
		})
	}
}
