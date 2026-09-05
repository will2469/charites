package token_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/will2469/charites/internal/token"
)

func TestTheme_BananaTest(t *testing.T) {
	input := []byte(`:root {
  --banana: #123456;
  --thing-that-is-definitely-not-primary: red;
  --super-special-design-token: var(--banana);
}`)

	ctx, err := token.ParseCSS(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	bananaTokens := ctx.ByName("--banana")
	if len(bananaTokens) != 1 {
		t.Fatalf("expected 1 --banana token, got %d", len(bananaTokens))
	}
	if bananaTokens[0].RawValue != "#123456" {
		t.Errorf("expected raw value #123456, got %q", bananaTokens[0].RawValue)
	}

	notPrimary := ctx.ByName("--thing-that-is-definitely-not-primary")
	if len(notPrimary) != 1 || notPrimary[0].RawValue != "red" {
		t.Errorf("failed extracting --thing-that-is-definitely-not-primary")
	}

	specialToken := ctx.ByName("--super-special-design-token")
	if len(specialToken) != 1 {
		t.Fatalf("expected 1 --super-special-design-token")
	}
	if specialToken[0].RawValue != "var(--banana)" {
		t.Errorf("expected raw value var(--banana), got %q", specialToken[0].RawValue)
	}
	if len(specialToken[0].References) != 1 || specialToken[0].References[0] != "--banana" {
		t.Errorf("expected reference to --banana, got %v", specialToken[0].References)
	}

	// Test resolution
	resolved, ok, err := ctx.Resolve(specialToken[0].ID, token.ResolveOptions{})
	if err != nil || !ok {
		t.Fatalf("failed to resolve --super-special-design-token: %v", err)
	}
	if resolved != "#123456" {
		t.Errorf("expected resolved value #123456, got %q", resolved)
	}
}

func TestTheme_CycleDetection(t *testing.T) {
	input := []byte(`:root {
  --a: var(--b);
  --b: var(--a);
}`)

	ctx, err := token.ParseCSS(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	aToks := ctx.ByName("--a")
	if len(aToks) != 1 {
		t.Fatalf("expected 1 --a token")
	}

	_, _, err = ctx.Resolve(aToks[0].ID, token.ResolveOptions{})
	if err == nil {
		t.Fatalf("expected cycle error, got nil")
	}
	if !errors.Is(err, token.ErrCycleDetected) {
		t.Errorf("expected ErrCycleDetected, got %v", err)
	}

	cycles := ctx.FindCycles()
	if len(cycles) == 0 {
		t.Errorf("expected FindCycles to return cycle nodes")
	}
}

func TestTheme_MultiHopResolution(t *testing.T) {
	input := []byte(`:root {
  --base: #abcdef;
  --alias-1: var(--base);
  --alias-2: var(--alias-1);
  --alias-3: var(--alias-2);
}`)

	ctx, err := token.ParseCSS(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	a3 := ctx.ByName("--alias-3")
	if len(a3) != 1 {
		t.Fatalf("expected 1 --alias-3 token")
	}

	resolved, ok, err := ctx.Resolve(a3[0].ID, token.ResolveOptions{})
	if err != nil || !ok {
		t.Fatalf("failed multi-hop resolution: %v", err)
	}
	if resolved != "#abcdef" {
		t.Errorf("expected #abcdef, got %q", resolved)
	}
}

func TestTheme_EvaluationBudget(t *testing.T) {
	input := []byte(`:root {
  --step-1: #111;
  --step-2: var(--step-1);
  --step-3: var(--step-2);
  --step-4: var(--step-3);
}`)

	ctx, err := token.ParseCSS(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s4 := ctx.ByName("--step-4")
	if len(s4) != 1 {
		t.Fatalf("expected --step-4")
	}

	// Set budget to 2 nodes (should exceed because chain has 4 nodes)
	_, _, err = ctx.Resolve(s4[0].ID, token.ResolveOptions{MaxNodes: 2})
	if err == nil {
		t.Fatalf("expected budget error, got nil")
	}
	if !errors.Is(err, token.ErrEvaluationBudgetExceeded) {
		t.Errorf("expected ErrEvaluationBudgetExceeded, got %v", err)
	}
}

func TestTheme_ScopeSeparationAndSpecificity(t *testing.T) {
	input := []byte(`:root {
  --brand: #111111;
}

.card {
  --brand: #222222;
}

#header.hero .card {
  --brand: #333333;
}`)

	ctx, err := token.ParseCSS(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	brandTokens := ctx.ByName("--brand")
	if len(brandTokens) != 3 {
		t.Fatalf("expected 3 distinct declarations of --brand, got %d", len(brandTokens))
	}

	// Verify each token preserves its scope
	if brandTokens[0].Scope.Selector != ":root" || brandTokens[0].RawValue != "#111111" {
		t.Errorf("unexpected root token: %+v", brandTokens[0])
	}
	if brandTokens[1].Scope.Selector != ".card" || brandTokens[1].RawValue != "#222222" {
		t.Errorf("unexpected card token: %+v", brandTokens[1])
	}
	if brandTokens[2].Scope.Selector != "#header.hero .card" || brandTokens[2].RawValue != "#333333" {
		t.Errorf("unexpected hero token: %+v", brandTokens[2])
	}

	// Check specificity
	rootSpec := brandTokens[0].Scope.Specificity
	heroSpec := brandTokens[2].Scope.Specificity
	if !heroSpec.GreaterThan(rootSpec) {
		t.Errorf("expected hero specificity to be greater than root specificity")
	}
	if heroSpec.A != 1 || heroSpec.B != 2 {
		t.Errorf("expected hero specificity A=1, B=2, got %+v", heroSpec)
	}
}

func TestTheme_ExtractGlobalCSSFixture(t *testing.T) {
	fixturePath := filepath.Join("..", "..", "tests", "fixtures", "global.css")
	data, err := os.ReadFile(filepath.Clean(fixturePath))
	if err != nil {
		t.Fatalf("failed to read global.css fixture: %v", err)
	}

	ctx, err := token.ParseCSS(data)
	if err != nil {
		t.Fatalf("unexpected error parsing global.css: %v", err)
	}

	// 1. Structural Checks
	if !ctx.HasScopeProperty("color-scheme", "light") {
		t.Errorf("expected color-scheme light to be declared")
	}

	// 2. Token extraction check
	colorTokens := ctx.ByPrefix("--color-")
	if len(colorTokens) == 0 {
		t.Errorf("expected --color- tokens to be extracted from global.css")
	}

	radiusTokens := ctx.ByPrefix("--radius-")
	if len(radiusTokens) == 0 {
		t.Errorf("expected --radius- tokens to be extracted from global.css")
	}

	// 3. Check specific tokens
	primaryTokens := ctx.ByName("--color-primary")
	if len(primaryTokens) == 0 {
		t.Errorf("expected --color-primary token to exist")
	}

	primaryLightTokens := ctx.ByName("--color-primary-light")
	if len(primaryLightTokens) == 0 {
		t.Errorf("expected --color-primary-light token to exist")
	}
}

func TestTheme_DiscoverAndLoad(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("ZeroConfigFallback", func(t *testing.T) {
		ctx, err := token.DiscoverAndLoad(tempDir, "")
		if err != nil {
			t.Fatalf("unexpected error on empty dir: %v", err)
		}
		if len(ctx.Tokens()) != 0 {
			t.Errorf("expected empty context tokens, got %d", len(ctx.Tokens()))
		}
	})

	t.Run("StandardDiscovery", func(t *testing.T) {
		styleDir := filepath.Join(tempDir, "src", "style")
		if err := os.MkdirAll(styleDir, 0750); err != nil {
			t.Fatalf("failed to create style dir: %v", err)
		}
		cssContent := []byte(`:root { --color-primary: oklch(0.5 0.2 250); --color-primary-light: oklch(0.5 0.2 250 / 10%); }`)
		if err := os.WriteFile(filepath.Join(styleDir, "global.css"), cssContent, 0600); err != nil {
			t.Fatalf("failed to write global.css: %v", err)
		}

		discovered, err := token.DiscoverAndLoad(tempDir, "")
		if err != nil {
			t.Fatalf("failed to discover standard global.css: %v", err)
		}
		toks := discovered.ByName("--color-primary")
		if len(toks) == 0 {
			t.Errorf("expected --color-primary to be discovered")
		}
	})

	t.Run("CustomPath", func(t *testing.T) {
		customDir := filepath.Join(tempDir, "custom")
		if err := os.MkdirAll(customDir, 0750); err != nil {
			t.Fatalf("failed to create custom dir: %v", err)
		}
		customFile := filepath.Join(customDir, "tokens.css")
		if err := os.WriteFile(customFile, []byte(`:root { --color-brand: #ff0000; }`), 0600); err != nil {
			t.Fatalf("failed to write tokens.css: %v", err)
		}

		customCtx, err := token.DiscoverAndLoad(tempDir, filepath.Join("custom", "tokens.css"))
		if err != nil {
			t.Fatalf("failed to load custom path: %v", err)
		}
		toks := customCtx.ByName("--color-brand")
		if len(toks) == 0 {
			t.Errorf("expected --color-brand to be discovered from custom path")
		}
	})
}

func TestTheme_ComplexVarFallbacks(t *testing.T) {
	input := []byte(`:root {
  --color-with-rgb-fallback: var(--non-existent, rgb(1, 2, 3));
  --color-with-nested-fallback: var(--missing-a, var(--missing-b, #123456));
  --tw-bg-opacity\:1: 0.85;
  --consumer: var(--tw-bg-opacity\:1);
}`)

	ctx, err := token.ParseCSS(input)
	if err != nil {
		t.Fatalf("ParseCSS failed: %v", err)
	}

	// 1. Uji fallback rgb(1, 2, 3) tidak rusak oleh tanda kurung
	rgbToks := ctx.ByName("--color-with-rgb-fallback")
	if len(rgbToks) != 1 {
		t.Fatalf("expected 1 token for --color-with-rgb-fallback, got %d", len(rgbToks))
	}
	rgbVal, ok, err := ctx.Resolve(rgbToks[0].ID, token.ResolveOptions{})
	if err != nil || !ok {
		t.Fatalf("failed to resolve rgb fallback: %v", err)
	}
	if rgbVal != "rgb(1, 2, 3)" {
		t.Errorf("expected resolved 'rgb(1, 2, 3)', got %q", rgbVal)
	}

	// 2. Uji nested fallback resolution
	nestedToks := ctx.ByName("--color-with-nested-fallback")
	if len(nestedToks) != 1 {
		t.Fatalf("expected 1 token for --color-with-nested-fallback, got %d", len(nestedToks))
	}
	nestedVal, ok, err := ctx.Resolve(nestedToks[0].ID, token.ResolveOptions{})
	if err != nil || !ok {
		t.Fatalf("failed to resolve nested fallback: %v", err)
	}
	if nestedVal != "#123456" {
		t.Errorf("expected resolved '#123456', got %q", nestedVal)
	}

	// 3. Uji CSS escaped custom property name dan referensinya
	escapedToks := ctx.ByName("--tw-bg-opacity:1")
	if len(escapedToks) != 1 {
		t.Fatalf("expected 1 token for --tw-bg-opacity:1, got %d", len(escapedToks))
	}
	consumerToks := ctx.ByName("--consumer")
	if len(consumerToks) != 1 {
		t.Fatalf("expected 1 token for --consumer, got %d", len(consumerToks))
	}
	consumerVal, ok, err := ctx.Resolve(consumerToks[0].ID, token.ResolveOptions{})
	if err != nil || !ok {
		t.Fatalf("failed to resolve consumer of escaped var: %v", err)
	}
	if consumerVal != "0.85" {
		t.Errorf("expected resolved '0.85', got %q", consumerVal)
	}
}
