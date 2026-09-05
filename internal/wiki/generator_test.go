package wiki_test

import (
	"bytes"
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
	if !strings.Contains(homeContent, "[`theme.hardcode-opacity-color`](theme.hardcode-opacity-color)") {
		t.Errorf("Home.md missing link to theme.hardcode-opacity-color")
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
	if !strings.Contains(themeContent, "[`theme.hardcode-opacity-color`](theme.hardcode-opacity-color)") {
		t.Errorf("theme.md missing link to theme.hardcode-opacity-color")
	}

	// 3. Verify theme.hardcode-opacity-color.md
	rulePath := filepath.Join(tmpDir, "theme.hardcode-opacity-color.md")
	ruleBytes, readErr := os.ReadFile(filepath.Clean(rulePath)) //nolint:gosec // controlled test path
	if readErr != nil {
		t.Fatalf("theme.hardcode-opacity-color.md was not generated: %v", readErr)
	}
	ruleContent := string(ruleBytes)
	if !strings.Contains(ruleContent, "# theme.hardcode-opacity-color") {
		t.Errorf("rule doc missing header")
	}

	// 4. Verify _Sidebar.md for GitHub Wiki hierarchical navigation
	sidebarPath := filepath.Join(tmpDir, "_Sidebar.md")
	sidebarBytes, readErr := os.ReadFile(filepath.Clean(sidebarPath)) //nolint:gosec // controlled test path
	if readErr != nil {
		t.Fatalf("_Sidebar.md was not generated: %v", readErr)
	}
	sidebarContent := string(sidebarBytes)
	if !strings.Contains(sidebarContent, "* [**Home**](Home)") {
		t.Errorf("_Sidebar.md missing Home link")
	}
	if !strings.Contains(sidebarContent, "<details open>") {
		t.Errorf("_Sidebar.md missing collapsible details")
	}
	if !strings.Contains(sidebarContent, "<summary><b>Theme</b> (15 rules)") {
		t.Errorf("_Sidebar.md missing Theme category summary")
	}
	if !strings.Contains(sidebarContent, "* [**Theme Overview**](theme)") {
		t.Errorf("_Sidebar.md missing Theme Overview link")
	}
	if !strings.Contains(sidebarContent, "* [`theme.hardcode-opacity-color`](theme.hardcode-opacity-color)") {
		t.Errorf("_Sidebar.md missing rule entry in sidebar")
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

	rulePath := filepath.Join(tmpDir, "theme.hardcode-opacity-color.md")
	if _, statErr := os.Stat(rulePath); os.IsNotExist(statErr) {
		t.Errorf("expected theme.hardcode-opacity-color.md to be created")
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

func TestWikiGenerator_DynamicCategoriesAndAtomic(t *testing.T) {
	tmpTarget := t.TempDir()
	reg := rules.DefaultRegistry()

	gen := wiki.NewGenerator(reg)
	err := gen.Generate(tmpTarget)
	if err != nil {
		t.Fatalf("failed to generate wiki: %v", err)
	}

	homeContent, err := os.ReadFile(filepath.Clean(filepath.Join(tmpTarget, "Home.md"))) //nolint:gosec // controlled test path
	if err != nil || !strings.Contains(string(homeContent), "theme.hardcode-opacity-color") {
		t.Errorf("Home.md missing rule entry: %v", err)
	}

	ruleDoc, err := os.ReadFile(filepath.Clean(filepath.Join(tmpTarget, "theme.hardcode-opacity-color.md"))) //nolint:gosec // controlled test path
	if err != nil || !strings.Contains(string(ruleDoc), "## 1. Overview & Core Invariant") {
		t.Errorf("rule 8-pillars document missing or corrupted: %v", err)
	}

	// Determinism test: generating twice produces byte-for-byte identical output
	secondTarget := t.TempDir()
	if err := gen.Generate(secondTarget); err != nil {
		t.Fatalf("second generation failed: %v", err)
	}

	firstBytes, _ := os.ReadFile(filepath.Clean(filepath.Join(tmpTarget, "theme.hardcode-opacity-color.md")))     //nolint:gosec // controlled test path
	secondBytes, _ := os.ReadFile(filepath.Clean(filepath.Join(secondTarget, "theme.hardcode-opacity-color.md"))) //nolint:gosec // controlled test path
	if !bytes.Equal(firstBytes, secondBytes) {
		t.Errorf("Wiki output is not byte-for-byte identical across runs")
	}
}

func TestRenderRuleDoc(t *testing.T) {
	reg := rules.DefaultRegistry()
	rule, ok := reg.Get("theme.hardcode-opacity-color")
	if !ok {
		t.Fatalf("expected rule theme.hardcode-opacity-color to exist")
	}

	doc, err := wiki.RenderRuleDoc(rule)
	if err != nil {
		t.Fatalf("RenderRuleDoc failed: %v", err)
	}
	if !strings.Contains(doc, "# theme.hardcode-opacity-color") {
		t.Errorf("expected doc to contain header, got: %s", doc)
	}
	if !strings.Contains(doc, "## 1. Overview & Core Invariant") {
		t.Errorf("expected doc to contain overview section, got: %s", doc)
	}

	// Nil rule handling
	_, err = wiki.RenderRuleDoc(nil)
	if err == nil {
		t.Errorf("expected error when rendering nil rule, got nil")
	}
}
