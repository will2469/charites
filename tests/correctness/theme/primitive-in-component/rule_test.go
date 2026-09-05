package correctness_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/will2469/charites/internal/rules/theme"
	"github.com/will2469/charites/tests/correctness/harness"
)

func TestRule_ThemePrimitiveInComponent_TriCorpus(t *testing.T) {
	rule := theme.NewPrimitiveInComponentRule()
	baseDir := filepath.Join(".", "theme", "primitive-in-component")
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		baseDir = "."
	}

	t.Run("Positive_Violations", func(t *testing.T) {
		posDir := filepath.Join(baseDir, "positive")
		diags := harness.EvaluateDir(t, rule, posDir)
		if len(diags) == 0 {
			t.Fatalf("expected positive violations > 0, got 0")
		}

		expectedPatterns := []string{
			"bg-blue-600",
			"hover:bg-blue-700",
			"text-slate-800",
			"border-gray-200",
			"ring-emerald-500",
			"fill-sky-400",
			"bg-[var(--blue-500)]",
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
