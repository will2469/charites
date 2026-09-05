package tests_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/will2469/charites/internal/cli"
	"github.com/will2469/charites/internal/rules/theme"
	"github.com/will2469/charites/internal/token"
)

func TestTokenIntegration_AutoDiscoveryAndResolution(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Buat struktur proyek realistis: src/style/global.css
	styleDir := filepath.Join(tempDir, "src", "style")
	if err := os.MkdirAll(styleDir, 0o750); err != nil {
		t.Fatalf("failed to create style dir: %v", err)
	}

	globalCSS := []byte(`@layer theme {
  :root {
    --banana: #ffcc00;
    --color-brand: oklch(0.60 0.22 260);
    --color-brand-light: oklch(0.60 0.22 260 / 15%);
    --color-brand-subtle: oklch(0.60 0.22 260 / 5%);
    --spacing-card: 1.5rem;
    --card-padding: var(--spacing-card);
  }

  .dark {
    --color-brand: oklch(0.75 0.18 260);
    --color-brand-light: oklch(0.75 0.18 260 / 25%);
  }
}`)

	if err := os.WriteFile(filepath.Join(styleDir, "global.css"), globalCSS, 0o600); err != nil {
		t.Fatalf("failed to write global.css: %v", err)
	}

	// 2. Uji DiscoverAndLoad dari subdirektori dalam (simulasi upward walk)
	subDir := filepath.Join(tempDir, "src", "components", "dashboard", "cards")
	if err := os.MkdirAll(subDir, 0o750); err != nil {
		t.Fatalf("failed to create sub dir: %v", err)
	}

	ctx, err := token.DiscoverAndLoad(subDir, "")
	if err != nil {
		t.Fatalf("DiscoverAndLoad failed from nested subDir: %v", err)
	}

	if ctx == nil || len(ctx.Tokens()) == 0 {
		t.Fatalf("expected tokens to be discovered from nested path, got 0 tokens")
	}

	// 3. Verifikasi Token Graph & Query Facade
	brandTokens := ctx.ByName("--color-brand")
	if len(brandTokens) != 2 {
		t.Fatalf("expected 2 declarations of --color-brand (:root and .dark), got %d", len(brandTokens))
	}

	var rootBrand, darkBrand token.Token
	for _, tok := range brandTokens {
		if tok.Scope.IsDark() {
			darkBrand = tok
		} else {
			rootBrand = tok
		}
	}

	if rootBrand.RawValue != "oklch(0.60 0.22 260)" {
		t.Errorf("unexpected root brand raw value: %q", rootBrand.RawValue)
	}
	if darkBrand.RawValue != "oklch(0.75 0.18 260)" {
		t.Errorf("unexpected dark brand raw value: %q", darkBrand.RawValue)
	}

	// 4. Verifikasi Resolusi Dependensi Token (var(--spacing-card) -> 1.5rem)
	cardPaddingTokens := ctx.ByName("--card-padding")
	if len(cardPaddingTokens) != 1 {
		t.Fatalf("expected 1 --card-padding token")
	}
	resolvedVal, ok, err := ctx.Resolve(cardPaddingTokens[0].ID, token.ResolveOptions{})
	if err != nil || !ok {
		t.Fatalf("failed to resolve --card-padding: %v", err)
	}
	if resolvedVal != "1.5rem" {
		t.Errorf("expected 1.5rem, got %q", resolvedVal)
	}

	// 5. Verifikasi Inferensi Konvensi Layer 4 (TokenConvention) terhadap Context Proyek Nyata
	conv := theme.NewDefaultCharitesConvention()
	cands, found := conv.FindOpacityReplacement("brand", "10", ctx)
	if !found || len(cands) == 0 {
		t.Fatalf("expected brand/10 replacement to be found in project tokens")
	}
	if cands[0].Name != "brand-light" {
		t.Errorf("expected brand-light, got %q", cands[0].Name)
	}

	// 6. Verifikasi Zero False Positive: Warna tanpa token -light di proyek TIDAK boleh mengembalikan pengganti
	_, foundUnmapped := conv.FindOpacityReplacement("banana", "10", ctx)
	if foundUnmapped {
		t.Errorf("banana does not have --banana-light, expected found=false")
	}
}

func TestTokenIntegration_CircularDependencyProtection(t *testing.T) {
	tempDir := t.TempDir()

	// Proyek dengan siklus sirkular di global.css
	cyclicCSS := []byte(`:root {
  --token-a: var(--token-b);
  --token-b: var(--token-c);
  --token-c: var(--token-a);
  --safe-token: #333333;
}`)

	if err := os.WriteFile(filepath.Join(tempDir, "global.css"), cyclicCSS, 0o600); err != nil {
		t.Fatalf("failed to write global.css: %v", err)
	}

	ctx, err := token.DiscoverAndLoad(tempDir, "")
	if err != nil {
		t.Fatalf("DiscoverAndLoad failed: %v", err)
	}

	// Graph harus menemukan siklus tanpa deadlock atau stack overflow
	cycles := ctx.Graph().FindCycles()
	if len(cycles) == 0 {
		t.Fatalf("expected circular tokens to be detected by FindCycles()")
	}

	// Resolusi token siklus harus mengembalikan ErrCycleDetected secara deterministik
	tokA := ctx.ByName("--token-a")
	if len(tokA) == 0 {
		t.Fatalf("expected --token-a to exist")
	}
	_, _, resolveErr := ctx.Resolve(tokA[0].ID, token.ResolveOptions{})
	if resolveErr == nil {
		t.Fatalf("expected cycle error, got nil")
	}

	// Resolusi safe token tetap harus berhasil
	tokSafe := ctx.ByName("--safe-token")
	if len(tokSafe) == 0 {
		t.Fatalf("expected --safe-token to exist")
	}
	safeVal, ok, err := ctx.Resolve(tokSafe[0].ID, token.ResolveOptions{})
	if err != nil || !ok || safeVal != "#333333" {
		t.Errorf("safe token resolution failed: %v, val=%q", err, safeVal)
	}
}

func TestTokenIntegration_EndToEndCLIScan(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Berkas SSOT global.css
	styleDir := filepath.Join(tempDir, "src", "styles")
	if err := os.MkdirAll(styleDir, 0o750); err != nil {
		t.Fatalf("failed to create styles dir: %v", err)
	}
	cssData := []byte(`:root {
  --color-primary: #3b82f6;
  --color-primary-light: #60a5fa;
}`)
	if err := os.WriteFile(filepath.Join(styleDir, "global.css"), cssData, 0o600); err != nil {
		t.Fatalf("failed to write global.css: %v", err)
	}

	// 2. Berkas target Astro dengan pelanggaran bg-primary/10
	componentsDir := filepath.Join(tempDir, "src", "components")
	if err := os.MkdirAll(componentsDir, 0o750); err != nil {
		t.Fatalf("failed to create components dir: %v", err)
	}
	astroFile := filepath.Join(componentsDir, "Badge.astro")
	astroContent := []byte(`---
const { label } = Astro.props;
---
<span class="inline-flex px-2 py-1 bg-primary/10 rounded">
  {label}
</span>
`)
	if err := os.WriteFile(astroFile, astroContent, 0o600); err != nil {
		t.Fatalf("failed to write Badge.astro: %v", err)
	}

	// 3. Eksekusi CLI scan format JSON
	var stdout, stderr bytes.Buffer
	exitCode := cli.ExecuteArgs([]string{"scan", tempDir, "-f", "json"}, &stdout, &stderr)
	if exitCode != 1 {
		t.Fatalf("expected exit code 1 (rule violations found), got %d. stderr: %s", exitCode, stderr.String())
	}

	var result struct {
		Diagnostics []struct {
			Rule    string `json:"rule"`
			Message string `json:"message"`
			Hint    string `json:"hint"`
		} `json:"diagnostics"`
	}

	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse CLI JSON output: %v\nOutput was: %s", err, stdout.String())
	}

	if len(result.Diagnostics) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(result.Diagnostics))
	}

	diag := result.Diagnostics[0]
	if diag.Rule != "theme.hardcode-opacity-color" {
		t.Errorf("expected rule 'theme.hardcode-opacity-color', got %q", diag.Rule)
	}
	if diag.Hint != "Use semantic token \"primary-light\"." {
		t.Errorf("expected hint with primary-light, got %q", diag.Hint)
	}
}
