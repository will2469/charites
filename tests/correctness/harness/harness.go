package harness

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/parser/astro"
	"github.com/will2469/charites/internal/parser/tsx"
	"github.com/will2469/charites/internal/rules"
)

// ParseAndEvaluate membaca file sumber fixture, melakukan parsing AST, dan mengevaluasi rule pada setiap node IR.
func ParseAndEvaluate(t *testing.T, rule rules.Rule, filePath string) []ir.Diagnostic {
	t.Helper()

	src, err := os.ReadFile(filepath.Clean(filePath)) //nolint:gosec
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

// EvaluateDir mengevaluasi seluruh fixture Astro/TSX/JSX di dalam direktori yang ditentukan.
func EvaluateDir(t *testing.T, rule rules.Rule, dirPath string) []ir.Diagnostic {
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
		diags := ParseAndEvaluate(t, rule, path)
		allDiags = append(allDiags, diags...)
	}

	return allDiags
}
