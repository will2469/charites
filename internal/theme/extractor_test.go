package theme_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/will2469/charites/internal/theme"
)

func TestTheme_ExtractGlobalCSSFixture(t *testing.T) {
	fixturePath := filepath.Join("..", "..", "tests", "fixtures", "global.css")
	data, err := os.ReadFile(filepath.Clean(fixturePath)) //nolint:gosec // test fixture path is controlled
	if err != nil {
		t.Fatalf("failed to read global.css fixture: %v", err)
	}

	ctx, err := theme.ParseCSS(data)
	if err != nil {
		t.Fatalf("unexpected error parsing global.css: %v", err)
	}

	// 1. Structural Checks
	if !ctx.IsLayered() {
		t.Errorf("expected IsLayered() to be true (@layer theme)")
	}
	if !ctx.HasColorScheme() {
		t.Errorf("expected HasColorScheme() to be true")
	}
	if ctx.RootColorScheme != "light" {
		t.Errorf("expected RootColorScheme to be 'light', got %q", ctx.RootColorScheme)
	}
	if !ctx.HasDarkColorScheme {
		t.Errorf("expected HasDarkColorScheme to be true")
	}
	if ctx.DarkColorScheme != "dark" {
		t.Errorf("expected DarkColorScheme to be 'dark', got %q", ctx.DarkColorScheme)
	}

	// 2. Base Color Tokens (OKLCH First-Class Citizen & Multi-Format)
	expectedColors := []string{
		"primary",
		"secondary",
		"muted",
		"accent",
		"destructive",
		"warning",
		"amber",
		"emerald",
		"background",
		"foreground",
		"card",
		"popover",
		"border",
		"input",
		"ring",
		"chart-1",
	}

	for _, col := range expectedColors {
		if !ctx.IsKnownColor(col) {
			t.Errorf("expected IsKnownColor(%q) to be true", col)
		}
	}

	// 3. Opacity Variants & Semantic Replacements
	opacityTests := []struct {
		slash   string
		wantRep string
		wantOk  bool
	}{
		{"primary/10", "primary-light", true},
		{"primary/20", "primary-light", true},
		{"primary/5", "primary-subtle", true},
		{"destructive/10", "destructive-light", true},
		{"destructive/20", "destructive-light", true},
		{"destructive/5", "destructive-subtle", true},
		{"warning/10", "warning-light", true},
		{"warning/5", "warning-subtle", true},
		{"amber/10", "amber-light", true},
		{"amber/5", "amber-subtle", true},
		{"emerald/10", "emerald-light", true},
		{"emerald/5", "emerald-subtle", true},
		{"accent/10", "accent-light", true},
		{"accent/5", "accent-subtle", true},
		{"primary/30", "", false},     // Uncalibrated opacity
		{"nonexistent/10", "", false}, // Unknown color
	}

	for _, tt := range opacityTests {
		got, ok := ctx.ReplacementForSlash(tt.slash)
		if ok != tt.wantOk {
			t.Errorf("ReplacementForSlash(%q) ok = %v, want %v", tt.slash, ok, tt.wantOk)
		}
		if got != tt.wantRep {
			t.Errorf("ReplacementForSlash(%q) = %q, want %q", tt.slash, got, tt.wantRep)
		}
	}

	// 4. Shape & Border Radius Tokens
	expectedRadii := []string{"sm", "md", "lg", "xl", "full"}
	for _, rad := range expectedRadii {
		if !ctx.IsKnownRadius(rad) {
			t.Errorf("expected IsKnownRadius(%q) to be true", rad)
		}
	}
	if ctx.IsKnownRadius("2xl") {
		t.Errorf("expected IsKnownRadius('2xl') to be false (not in fixture)")
	}

	// 5. Semantic Z-Index Scale
	expectedZ := []string{"dropdown", "sticky", "modal", "popover", "tooltip", "toast"}
	for _, z := range expectedZ {
		if !ctx.IsKnownZIndex(z) {
			t.Errorf("expected IsKnownZIndex(%q) to be true", z)
		}
	}
	if ctx.IsKnownZIndex("9999") {
		t.Errorf("expected IsKnownZIndex('9999') to be false")
	}

	// 6. Elevation Shadows & Blurs
	if !ctx.IsKnownShadow("sm") || !ctx.IsKnownShadow("md") || !ctx.IsKnownShadow("lg") {
		t.Errorf("expected shadow scales sm, md, lg to be known")
	}
	if !ctx.IsKnownBlur("sm") || !ctx.IsKnownBlur("md") {
		t.Errorf("expected blur scales sm, md to be known")
	}

	// 7. Dark Mode Parity Check
	if !ctx.CheckDarkParity("--color-primary") {
		t.Errorf("expected CheckDarkParity('--color-primary') to be true")
	}
	if !ctx.CheckDarkParity("--color-card") {
		t.Errorf("expected CheckDarkParity('--color-card') to be true")
	}

	missing := ctx.MissingDarkParityTokens()
	// In global.css, chart-1 through chart-5 only exist in :root, so they should be detected as missing in dark
	var foundChart1 bool
	for _, m := range missing {
		if m == "--color-chart-1" {
			foundChart1 = true
			break
		}
	}
	if !foundChart1 {
		t.Errorf("expected missing dark parity to contain --color-chart-1, got %v", missing)
	}
}

func TestTheme_DiscoverAndLoad(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("ZeroConfigFallback", func(t *testing.T) {
		ctx, err := theme.DiscoverAndLoad(tempDir, "")
		if err != nil {
			t.Fatalf("unexpected error on empty dir: %v", err)
		}
		if ctx.IsKnownColor("primary") {
			t.Errorf("expected empty context")
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

		discovered, err := theme.DiscoverAndLoad(tempDir, "")
		if err != nil {
			t.Fatalf("failed to discover standard global.css: %v", err)
		}
		if !discovered.IsKnownColor("primary") {
			t.Errorf("expected primary to be discovered")
		}
		rep, ok := discovered.ReplacementForSlash("primary/10")
		if !ok || rep != "primary-light" {
			t.Errorf("expected primary-light replacement, got %q (ok=%v)", rep, ok)
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

		customCtx, err := theme.DiscoverAndLoad(tempDir, filepath.Join("custom", "tokens.css"))
		if err != nil {
			t.Fatalf("failed to load custom path: %v", err)
		}
		if !customCtx.IsKnownColor("brand") {
			t.Errorf("expected brand to be discovered from custom path")
		}
	})
}

func BenchmarkTheme_Queries(b *testing.B) {
	fixturePath := filepath.Join("..", "..", "tests", "fixtures", "global.css")
	data, err := os.ReadFile(filepath.Clean(fixturePath)) //nolint:gosec // test fixture path is controlled
	if err != nil {
		b.Fatalf("failed to read global.css: %v", err)
	}
	ctx, err := theme.ParseCSS(data)
	if err != nil {
		b.Fatalf("failed to parse css: %v", err)
	}

	b.Run("IsKnownColor", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = ctx.IsKnownColor("primary")
		}
	})

	b.Run("ReplacementForSlash", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = ctx.ReplacementForSlash("primary/10")
		}
	})

	b.Run("IsKnownRadius", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = ctx.IsKnownRadius("md")
		}
	})

	b.Run("IsKnownZIndex", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = ctx.IsKnownZIndex("modal")
		}
	})
}
