package theme_test

import (
	"testing"

	"github.com/will2469/charites/internal/token/theme"
)

func TestCSS_ParserRobustness(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		checkFunc func(t *testing.T, sheet *theme.StyleSheet)
	}{
		{
			name: "Semicolon inside quotes",
			input: `.foo {
  --value: "hello;world";
}`,
			checkFunc: func(t *testing.T, sheet *theme.StyleSheet) {
				if len(sheet.Rules) != 1 {
					t.Fatalf("expected 1 rule, got %d", len(sheet.Rules))
				}
				style, ok := sheet.Rules[0].(*theme.StyleRule)
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
			checkFunc: func(t *testing.T, sheet *theme.StyleSheet) {
				style := sheet.Rules[0].(*theme.StyleRule)
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
			checkFunc: func(t *testing.T, sheet *theme.StyleSheet) {
				style := sheet.Rules[0].(*theme.StyleRule)
				if style.Selector != ".card" {
					t.Errorf("expected selector '.card', got %q", style.Selector)
				}
				if len(style.Declarations) != 1 || style.Declarations[0].Value != "red" {
					t.Errorf("expected surface: red")
				}
				if len(style.Rules) != 1 {
					t.Fatalf("expected 1 nested rule, got %d", len(style.Rules))
				}
				nested := style.Rules[0].(*theme.StyleRule)
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
			checkFunc: func(t *testing.T, sheet *theme.StyleSheet) {
				if len(sheet.Rules) != 1 {
					t.Fatalf("expected 1 top-level rule, got %d", len(sheet.Rules))
				}
				layer := sheet.Rules[0].(*theme.AtRule)
				if layer.Name != "@layer" || layer.Prelude != "theme" {
					t.Errorf("expected @layer theme, got %s %s", layer.Name, layer.Prelude)
				}
				if len(layer.Rules) != 1 {
					t.Fatalf("expected 1 nested rule in layer, got %d", len(layer.Rules))
				}
				media := layer.Rules[0].(*theme.AtRule)
				if media.Name != "@media" || media.Prelude != "(prefers-color-scheme: dark)" {
					t.Errorf("expected @media query, got %s %s", media.Name, media.Prelude)
				}
				if len(media.Rules) != 1 {
					t.Fatalf("expected 1 nested rule in media, got %d", len(media.Rules))
				}
				root := media.Rules[0].(*theme.StyleRule)
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
			checkFunc: func(t *testing.T, sheet *theme.StyleSheet) {
				root := sheet.Rules[0].(*theme.StyleRule)
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sheet, err := theme.Parse([]byte(tt.input))
			if err != nil {
				t.Fatalf("unexpected parse error: %v", err)
			}
			tt.checkFunc(t, sheet)
		})
	}
}
