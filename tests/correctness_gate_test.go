package tests_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/parser/astro"
	"github.com/will2469/charites/internal/parser/tsx"
	"github.com/will2469/charites/internal/rules"
)

// getCorpusDir resolves the golden corpus directory for a rule, supporting both nested
// category/slug paths (e.g. correctness/theme/hardcode-opacity-color) and flat paths.
func getCorpusDir(rule rules.Rule) string {
	slug := strings.TrimPrefix(rule.ID(), rule.Category()+".")
	nested := filepath.Join("correctness", rule.Category(), slug)
	if _, err := os.Stat(nested); err == nil {
		return nested
	}
	return filepath.Join("correctness", rule.ID())
}

// evaluateNodeTree walks the IR tree and collects all diagnostics for a rule.
func evaluateNodeTree(rule rules.Rule, root *ir.Node) []ir.Diagnostic {
	var diags []ir.Diagnostic
	for node := range root.Walk() {
		findings := rule.Evaluate(node)
		if len(findings) > 0 {
			diags = append(diags, findings...)
		}
	}
	return diags
}

// evaluateFile parses and evaluates a single fixture file.
func evaluateFile(t *testing.T, rule rules.Rule, filePath string) []ir.Diagnostic {
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

	return evaluateNodeTree(rule, root)
}

// evaluateCorpusDirectory evaluates all fixture files within a sub-corpus directory.
func evaluateCorpusDirectory(t *testing.T, rule rules.Rule, dirPath string) []ir.Diagnostic {
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
		diags := evaluateFile(t, rule, path)
		allDiags = append(allDiags, diags...)
	}

	return allDiags
}

// TestGoldenCorpus_AdoptionMatrix asserts Layer 1: Structural Presence Gate.
// Every registered rule must have positive/, negative/, and adversarial/ fixture directories with at least 1 fixture.
func TestGoldenCorpus_AdoptionMatrix(t *testing.T) {
	allRules := rules.All()
	if len(allRules) == 0 {
		t.Fatalf("no rules registered in rules.All()")
	}

	for _, rule := range allRules {
		ruleID := rule.ID()
		t.Run(ruleID, func(t *testing.T) {
			corpusDir := getCorpusDir(rule)
			if _, err := os.Stat(corpusDir); os.IsNotExist(err) {
				t.Fatalf("rule %s is missing golden corpus directory: %s", ruleID, corpusDir)
			}

			subdirs := []string{"positive", "negative", "adversarial"}
			for _, sub := range subdirs {
				p := filepath.Join(corpusDir, sub)
				info, err := os.Stat(p)
				if os.IsNotExist(err) || !info.IsDir() {
					t.Fatalf("rule %s is missing required sub-corpus directory: %s", ruleID, p)
				}

				entries, err := os.ReadDir(p)
				if err != nil || len(entries) == 0 {
					t.Fatalf("rule %s has empty sub-corpus directory: %s", ruleID, p)
				}
			}
		})
	}
}

// TestCorrectnessGate asserts Layer 2: Rule Correctness Metric Gate.
// RuleCorrectnessMetric = (PositiveViolations > 0) && (NegativeViolations == 0) && (AdversarialViolations == 0).
func TestCorrectnessGate(t *testing.T) {
	allRules := rules.All()
	if len(allRules) == 0 {
		t.Fatalf("no rules registered in rules.All()")
	}

	for _, rule := range allRules {
		ruleID := rule.ID()
		t.Run(ruleID, func(t *testing.T) {
			corpusDir := getCorpusDir(rule)

			posDiags := evaluateCorpusDirectory(t, rule, filepath.Join(corpusDir, "positive"))
			negDiags := evaluateCorpusDirectory(t, rule, filepath.Join(corpusDir, "negative"))
			advDiags := evaluateCorpusDirectory(t, rule, filepath.Join(corpusDir, "adversarial"))

			posCount := len(posDiags)
			negCount := len(negDiags)
			advCount := len(advDiags)

			t.Logf("Rule %s evaluation results: Positive=%d, Negative=%d, Adversarial=%d",
				ruleID, posCount, negCount, advCount)

			if posCount == 0 {
				t.Errorf("FAIL [ROAD-03-GATE-003]: %s produced 0 positive violations (want > 0)", ruleID)
			}
			if negCount > 0 {
				t.Errorf("FAIL [ROAD-03-GATE-003]: %s produced %d negative violations (want 0 - Zero Noise Invariant): %+v",
					ruleID, negCount, negDiags)
			}
			if advCount > 0 {
				t.Errorf("FAIL [ROAD-03-GATE-003]: %s produced %d adversarial violations (want 0 - Bait Immunity Invariant): %+v",
					ruleID, advCount, advDiags)
			}

			if posCount > 0 && negCount == 0 && advCount == 0 {
				t.Logf("PASS [ROAD-03-GATE-003]: %s satisfies RuleCorrectnessMetric 100%%", ruleID)
			}
		})
	}
}
