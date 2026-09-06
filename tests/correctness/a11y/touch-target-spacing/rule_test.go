package a11y_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/will2469/charites/internal/rules/a11y"
	"github.com/will2469/charites/tests/correctness/harness"
)

func TestRule_A11yTouchTargetSpacing_TriCorpus(t *testing.T) {
	rule := a11y.NewTouchTargetSpacingRule()
	baseDir := "."
	if _, err := os.Stat("positive"); os.IsNotExist(err) {
		baseDir = filepath.Join("tests", "correctness", "a11y", "touch-target-spacing")
	}

	t.Run("Positive_Violations", func(t *testing.T) {
		posDir := filepath.Join(baseDir, "positive")
		diags := harness.EvaluateDir(t, rule, posDir)
		if len(diags) == 0 {
			t.Fatalf("expected positive violations > 0, got 0")
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
}
