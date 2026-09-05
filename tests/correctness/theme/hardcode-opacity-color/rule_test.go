package correctness_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/parser/astro"
	"github.com/will2469/charites/internal/parser/tsx"
	"github.com/will2469/charites/internal/rules"
	"github.com/will2469/charites/internal/rules/theme"
)

func parseAndEvaluate(t *testing.T, rule rules.Rule, filePath string) []ir.Diagnostic {
	t.Helper()

	src, err := os.ReadFile(filepath.Clean(filePath)) //nolint:gosec // test fixture path is controlled by test runner
	if err != nil {
		t.Fatalf("failed to read fixture file %s: %v", filePath, err)
	}

	var root *ir.Node
	ext := filepath.Ext(filePath)
	switch ext {
	case ".astro":
		root, err = astro.Parse(src)
	case ".tsx", ".jsx", ".ts", ".js":
		root, err = tsx.Extract(src)
	default:
		t.Fatalf("unsupported fixture extension %s for %s", ext, filePath)
	}

	if err != nil {
		t.Fatalf("failed to parse fixture %s: %v", filePath, err)
	}
	if root == nil {
		t.Fatalf("parser returned nil root for %s", filePath)
	}

	var diags []ir.Diagnostic
	for node := range root.Walk() {
		findings := rule.Evaluate(node)
		if len(findings) > 0 {
			diags = append(diags, findings...)
		}
	}

	return diags
}

func evaluateDir(t *testing.T, rule rules.Rule, dirPath string) []ir.Diagnostic {
	t.Helper()

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		t.Fatalf("failed to read directory %s: %v", dirPath, err)
	}

	var allDiags []ir.Diagnostic
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := filepath.Ext(entry.Name())
		if ext != ".astro" && ext != ".tsx" && ext != ".jsx" {
			continue
		}
		path := filepath.Join(dirPath, entry.Name())
		diags := parseAndEvaluate(t, rule, path)
		allDiags = append(allDiags, diags...)
	}

	return allDiags
}

func TestRule_ThemeHardcodeOpacityColor_TriCorpus(t *testing.T) {
	rule := theme.NewHardcodeOpacityColorRule()
	baseDir := filepath.Join(".", "theme", "hardcode-opacity-color")

	// If running directly in directory
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		baseDir = "."
	}

	t.Run("Positive_Violations", func(t *testing.T) {
		posDir := filepath.Join(baseDir, "positive")
		diags := evaluateDir(t, rule, posDir)

		if len(diags) == 0 {
			t.Fatalf("expected positive violations > 0, got 0")
		}

		// Verify specific expected violations
		expectedPatterns := []string{
			"bg-primary/10",
			"border-destructive/20",
			"text-primary/20",
			"ring-warning/10",
			"bg-primary/5",
			"hover:bg-primary/10",
			"dark:bg-primary/10",
			"md:hover:bg-primary/10",
			"dark:border-destructive/20",
			"sm:dark:hover:border-destructive/20",
		}

		foundMap := make(map[string]bool)
		for _, d := range diags {
			for _, pat := range expectedPatterns {
				if strings.Contains(d.Message, pat) {
					foundMap[pat] = true
				}
			}
		}

		for _, pat := range expectedPatterns {
			if !foundMap[pat] {
				t.Errorf("expected violation for pattern %q was not detected", pat)
			}
		}
	})

	t.Run("Negative_ZeroNoise", func(t *testing.T) {
		negDir := filepath.Join(baseDir, "negative")
		diags := evaluateDir(t, rule, negDir)

		if len(diags) != 0 {
			t.Fatalf("expected 0 negative violations (Zero Noise Invariant), got %d: %+v", len(diags), diags)
		}
	})

	t.Run("Adversarial_BaitImmunity", func(t *testing.T) {
		advDir := filepath.Join(baseDir, "adversarial")
		diags := evaluateDir(t, rule, advDir)

		if len(diags) != 0 {
			t.Fatalf("expected 0 adversarial violations (Bait Immunity Invariant), got %d: %+v", len(diags), diags)
		}
	})
}
