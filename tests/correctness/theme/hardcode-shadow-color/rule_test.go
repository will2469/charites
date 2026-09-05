package correctness_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/will2469/charites/internal/rules/theme"
	"github.com/will2469/charites/tests/correctness/harness"
)

func TestRule_ThemeHardcodeShadowColor_TriCorpus(t *testing.T) {
	rule := theme.NewHardcodeShadowColorRule()
	baseDir := filepath.Join(".", "theme", "hardcode-shadow-color")
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
			"shadow-[0_4px_10px_#00000040]",
			"shadow-[0_10px_15px_rgba(0,0,0,0.1)]",
			"shadow-[0_2px_4px_#333]",
			"[box-shadow:0_4px_6px_#000]",
			"shadow-[0_8px_16px_#1e293b]",
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
