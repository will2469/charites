package theme_test

import (
	"testing"

	"github.com/will2469/charites/internal/config"
	"github.com/will2469/charites/internal/rules/theme"
	themeengine "github.com/will2469/charites/internal/token"
)

func TestConvention_ConfigurableAndDesignAgnostic(t *testing.T) {
	// Proyek mendefinisikan skema token unik (bukan sekadar primary/secondary):
	// - Warna: "pisang" dengan fallback "kuning-muda"
	// - Opacity: /15 -> varian "-tint"
	// - Prefix custom property: "--ds-"
	cssInput := []byte(`:root {
  --ds-kuning-muda-tint: #ffeb3b;
  --ds-melon-soft: #c8e6c9;
}`)

	ctx, err := themeengine.ParseCSS(cssInput)
	if err != nil {
		t.Fatalf("ParseCSS failed: %v", err)
	}

	// Buat convention custom berbasis adapter pluggable
	customConv := theme.NewConfigurableConvention(
		theme.WithOpacityMapping("15", "-tint"),
		theme.WithOpacityMapping("30", "-soft"),
		theme.WithFallback("pisang", "kuning-muda"),
		theme.WithPrefixes("--ds-"),
	)

	// 1. Uji custom opacity mapping & fallback: "pisang" / 15 -> "--ds-kuning-muda-tint"
	candidates, ok := customConv.FindOpacityReplacement("pisang", "15", ctx)
	if !ok || len(candidates) != 1 {
		t.Fatalf("expected 1 replacement candidate for pisang/15, got %v (ok=%v)", candidates, ok)
	}
	if candidates[0].Name != "kuning-muda-tint" {
		t.Errorf("expected candidate 'kuning-muda-tint', got %q", candidates[0].Name)
	}
	if candidates[0].RawValue != "#ffeb3b" {
		t.Errorf("expected raw value '#ffeb3b', got %q", candidates[0].RawValue)
	}

	// 2. Uji direct mapping tanpa fallback: "melon" / 30 -> "--ds-melon-soft"
	melonCands, ok := customConv.FindOpacityReplacement("melon", "30", ctx)
	if !ok || len(melonCands) != 1 {
		t.Fatalf("expected 1 replacement candidate for melon/30")
	}
	if melonCands[0].Name != "melon-soft" {
		t.Errorf("expected 'melon-soft', got %q", melonCands[0].Name)
	}

	// 3. Uji ketiadaan token di graph (Banana test):
	// "jeruk" / 15 tidak ada di CSS -> tidak boleh mengarang token fiktif!
	jerukCands, ok := customConv.FindOpacityReplacement("jeruk", "15", ctx)
	if ok || len(jerukCands) > 0 {
		t.Fatalf("expected no replacement for unconfigured/missing token 'jeruk', got %v", jerukCands)
	}
}

func TestConvention_FromConfig(t *testing.T) {
	yamlContent := `
convention:
  opacity_mappings:
    "10": ["-soft", "-light"]
    "15": ["-tint"]
  fallbacks:
    accent: ["highlight"]
  prefixes:
    - "--my-theme-"
`

	cfg, err := config.Parse([]byte(yamlContent))
	if err != nil {
		t.Fatalf("failed parsing config: %v", err)
	}

	conv := theme.NewConventionFromConfig(cfg.Convention)

	cssInput := []byte(`:root {
  --my-theme-highlight-tint: #fffae6;
}`)
	ctx, err := themeengine.ParseCSS(cssInput)
	if err != nil {
		t.Fatalf("failed parsing css: %v", err)
	}

	cands, ok := conv.FindOpacityReplacement("accent", "15", ctx)
	if !ok || len(cands) != 1 {
		t.Fatalf("expected 1 replacement for accent/15 via config, got %v", cands)
	}
	if cands[0].Name != "highlight-tint" {
		t.Errorf("expected 'highlight-tint', got %q", cands[0].Name)
	}
}
