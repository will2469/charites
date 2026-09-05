package wiki_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/will2469/charites/internal/rules"
	"github.com/will2469/charites/internal/wiki"
)

func TestGenerator_Generate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "wiki-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	gen := wiki.NewGenerator(nil) // uses DefaultRegistry

	err = gen.Generate(tmpDir)
	if err != nil {
		t.Fatalf("gen.Generate failed: %v", err)
	}

	// 1. Verify Home.md
	homePath := filepath.Join(tmpDir, "Home.md")
	homeBytes, readErr := os.ReadFile(filepath.Clean(homePath)) //nolint:gosec // controlled test path
	if readErr != nil {
		t.Fatalf("Home.md was not generated: %v", readErr)
	}
	homeContent := string(homeBytes)
	if !strings.Contains(homeContent, "theme") {
		t.Errorf("Home.md missing 'theme' category")
	}
	if !strings.Contains(homeContent, "theme/hardcode-opacity-color.md") {
		t.Errorf("Home.md missing 'theme/hardcode-opacity-color.md'")
	}

	// 2. Verify theme.md
	themePath := filepath.Join(tmpDir, "theme.md")
	themeBytes, readErr := os.ReadFile(filepath.Clean(themePath)) //nolint:gosec // controlled test path
	if readErr != nil {
		t.Fatalf("theme.md was not generated: %v", readErr)
	}
	themeContent := string(themeBytes)
	if !strings.Contains(themeContent, "# Theme Rules (`theme`)") {
		t.Errorf("theme.md missing header")
	}
	if !strings.Contains(themeContent, "theme/hardcode-opacity-color.md") {
		t.Errorf("theme.md missing link to theme/hardcode-opacity-color.md")
	}

	// 3. Verify theme/hardcode-opacity-color.md
	rulePath := filepath.Join(tmpDir, "theme", "hardcode-opacity-color.md")
	ruleBytes, readErr := os.ReadFile(filepath.Clean(rulePath)) //nolint:gosec // controlled test path
	if readErr != nil {
		t.Fatalf("theme/hardcode-opacity-color.md was not generated: %v", readErr)
	}
	ruleContent := string(ruleBytes)
	if !strings.Contains(ruleContent, "# theme.hardcode-opacity-color") {
		t.Errorf("rule doc missing header")
	}
}

func TestGenerator_WithCustomRegistry(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "wiki-custom-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	reg := rules.NewRegistry()
	_ = rules.RegisterBuiltinRules(reg)

	gen := wiki.NewGenerator(reg)
	err = gen.Generate(tmpDir)
	if err != nil {
		t.Fatalf("Generate with custom registry failed: %v", err)
	}

	themePath := filepath.Join(tmpDir, "theme.md")
	if _, statErr := os.Stat(themePath); os.IsNotExist(statErr) {
		t.Errorf("expected theme.md to be created")
	}

	rulePath := filepath.Join(tmpDir, "theme", "hardcode-opacity-color.md")
	if _, statErr := os.Stat(rulePath); os.IsNotExist(statErr) {
		t.Errorf("expected theme/hardcode-opacity-color.md to be created")
	}
}

func TestGenerator_RegenerateWiki(t *testing.T) {
	gen := wiki.NewGenerator(nil)
	repoRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("failed to get repo root: %v", err)
	}
	wikiDir := filepath.Join(repoRoot, "wiki")
	if genErr := gen.Generate(wikiDir); genErr != nil {
		t.Fatalf("failed to generate wiki into %s: %v", wikiDir, genErr)
	}
}
