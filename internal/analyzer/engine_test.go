package analyzer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/will2469/charites/internal/analyzer"
	"github.com/will2469/charites/internal/config"
	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules/theme"
)

func TestEngine_AnalyzeFile_TSX(t *testing.T) {
	tmpDir := t.TempDir()
	tsxPath := filepath.Join(tmpDir, "Component.tsx")

	content := `
export function Button() {
  return (
    <div>
      {/* Kasus 1: Pelanggaran aktif */}
      <span className="bg-primary/10" />

      {/* Kasus 2: Pelanggaran ditekan oleh komentar inline */}
      // charites:ignore theme.hardcode-opacity-color
      <button className="bg-primary/20" />
    </div>
  );
}
`
	if err := os.WriteFile(tsxPath, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	rule := theme.NewHardcodeOpacityColorRule()
	activeRules := []config.ActiveRule{
		{
			Rule:              rule,
			EffectiveSeverity: ir.SeverityError,
		},
	}

	eng := analyzer.NewEngine(activeRules)
	if len(eng.ActiveRules()) != 1 {
		t.Fatalf("expected 1 active rule, got %d", len(eng.ActiveRules()))
	}

	diags, err := eng.AnalyzeFile(tsxPath)
	if err != nil {
		t.Fatalf("AnalyzeFile failed: %v", err)
	}

	// Hanya 1 pelanggaran yang dilaporkan (bg-primary/10), karena bg-primary/20 ditekan
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d: %+v", len(diags), diags)
	}

	if diags[0].Rule != "theme.hardcode-opacity-color" {
		t.Errorf("unexpected rule ID: %s", diags[0].Rule)
	}
	if diags[0].Severity != ir.SeverityError {
		t.Errorf("expected EffectiveSeverity Error, got %v", diags[0].Severity)
	}
}

func TestEngine_AnalyzeFile_Astro(t *testing.T) {
	tmpDir := t.TempDir()
	astroPath := filepath.Join(tmpDir, "Card.astro")

	content := `---
const title = "My Card";
---
<div class="p-4">
  <!-- Kasus 1: Pelanggaran aktif di Astro HTML template -->
  <header class="bg-primary/10">{title}</header>

  <!-- Kasus 2: Pelanggaran ditekan via HTML comment -->
  <!-- charites:ignore theme.hardcode-opacity-color -->
  <footer class="bg-destructive/10">Footer</footer>
</div>
`
	if err := os.WriteFile(astroPath, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	rule := theme.NewHardcodeOpacityColorRule()
	activeRules := []config.ActiveRule{
		{
			Rule:              rule,
			EffectiveSeverity: ir.SeverityWarn,
		},
	}

	eng := analyzer.NewEngine(activeRules)
	diags, err := eng.AnalyzeFile(astroPath)
	if err != nil {
		t.Fatalf("AnalyzeFile failed: %v", err)
	}

	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic for Astro file, got %d: %+v", len(diags), diags)
	}

	if diags[0].Severity != ir.SeverityWarn {
		t.Errorf("expected EffectiveSeverity Warn, got %v", diags[0].Severity)
	}
}

func TestEngine_AnalyzeFile_EdgeCases(t *testing.T) {
	rule := theme.NewHardcodeOpacityColorRule()
	eng := analyzer.NewEngine([]config.ActiveRule{{Rule: rule, EffectiveSeverity: ir.SeverityError}})

	// 1. File tidak ada
	_, err := eng.AnalyzeFile("/non/existent/path.tsx")
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}

	// 2. Ekstensi tidak didukung (contoh: .css)
	tmpDir := t.TempDir()
	cssPath := filepath.Join(tmpDir, "styles.css")
	_ = os.WriteFile(cssPath, []byte(".a { color: red; }"), 0o600)

	diags, err := eng.AnalyzeFile(cssPath)
	if err != nil || len(diags) != 0 {
		t.Errorf("expected nil error and 0 diags for unsupported ext, got diags=%+v, err=%v", diags, err)
	}

	// 3. AnalyzeTree dengan root nil
	diagsNil := eng.AnalyzeTree("a.tsx", nil, nil)
	if len(diagsNil) != 0 {
		t.Errorf("expected 0 diags for nil root, got %d", len(diagsNil))
	}
}
