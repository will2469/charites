package correctness_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/will2469/charites/internal/drift"
	"github.com/will2469/charites/internal/parser/tsx"
	"github.com/will2469/charites/internal/rules/design"
	"github.com/will2469/charites/tests/correctness/harness"
)

func TestRule_DesignComponentStyleDrift_TriCorpus(t *testing.T) {
	rule := design.NewComponentStyleDriftRule()
	baseDir := filepath.Join(".", "design", "component-style-drift")
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		baseDir = "."
	}

	t.Run("Positive_Violations", func(t *testing.T) {
		posDir := filepath.Join(baseDir, "positive")
		diags := harness.EvaluateDir(t, rule, posDir)
		if len(diags) == 0 {
			t.Fatalf("expected positive violations > 0, got 0")
		}

		var hasContrastHazard bool
		for _, d := range diags {
			if strings.Contains(d.Message, "WCAG 1.4.3 Contrast Hazard") &&
				strings.Contains(d.Message, "bg-yellow-400") &&
				strings.Contains(d.Message, "text-white") {
				hasContrastHazard = true
				break
			}
		}

		if !hasContrastHazard {
			t.Errorf("expected WCAG 1.4.3 contrast hazard on bg-yellow-400 with text-white")
		}
	})

	t.Run("Negative_ZeroNoise", func(t *testing.T) {
		negDir := filepath.Join(baseDir, "negative")
		diags := harness.EvaluateDir(t, rule, negDir)
		if len(diags) != 0 {
			t.Fatalf("expected 0 negative violations, got %d: %+v", len(diags), diags)
		}
	})

	t.Run("Adversarial_BaitImmunity", func(t *testing.T) {
		advDir := filepath.Join(baseDir, "adversarial")
		diags := harness.EvaluateDir(t, rule, advDir)
		if len(diags) != 0 {
			t.Fatalf("expected 0 adversarial violations, got %d: %+v", len(diags), diags)
		}
	})

	t.Run("Repository_EmpiricalDrift", func(t *testing.T) {
		posDir := filepath.Join(baseDir, "positive")
		entries, err := os.ReadDir(posDir)
		if err != nil {
			t.Fatalf("failed to read positive dir: %v", err)
		}

		var allOccurrences []drift.StyleOccurrence
		for _, e := range entries {
			if filepath.Ext(e.Name()) != ".tsx" {
				continue
			}
			filePath := filepath.Join(posDir, e.Name())
			src, err := os.ReadFile(filepath.Clean(filePath)) //nolint:gosec
			if err != nil {
				t.Fatalf("failed to read fixture %s: %v", filePath, err)
			}
			root, err := tsx.Extract(src)
			if err != nil {
				t.Fatalf("failed to parse fixture %s: %v", filePath, err)
			}
			allOccurrences = append(allOccurrences, drift.ExtractTreeOccurrences(filePath, root)...)
		}

		analyzer := drift.NewAnalyzer(drift.DefaultOptions())
		report := analyzer.Analyze(allOccurrences)

		key := drift.ClusterKey{
			Scope: drift.ScopeKey{
				Kind:       drift.KindButton,
				Confidence: drift.ConfidenceExact,
			},
			Category: drift.CatRounded,
		}

		cluster, exists := report.GeometricClusters[key]
		if !exists {
			t.Fatalf("expected cluster for Button rounded, but none found")
		}

		if cluster.InsufficientData {
			t.Fatalf("cluster flagged with InsufficientData, expected valid statistical decision (total=%d)", cluster.Total)
		}

		if cluster.CanonicalToken != "rounded-md" {
			t.Errorf("expected canonical token 'rounded-md', got %q", cluster.CanonicalToken)
		}

		if cluster.CanonicalDominance < 80.0 {
			t.Errorf("expected canonical dominance >= 80%%, got %.1f%%", cluster.CanonicalDominance)
		}

		var foundRogue bool
		for _, out := range cluster.Outliers {
			if out.Token == "rounded-2xl" {
				foundRogue = true
				if out.Percentage > 10.0 {
					t.Errorf("expected rogue percentage <= 10%%, got %.1f%%", out.Percentage)
				}
				break
			}
		}

		if !foundRogue {
			t.Errorf("expected rogue outlier 'rounded-2xl' was not flagged in cluster: %+v", cluster.Outliers)
		}
	})
}
