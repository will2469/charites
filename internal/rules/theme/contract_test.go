package theme_test

import (
	"strings"
	"testing"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules"
)

// TestThemeRules_CanonicalContract menguji kontrak kanonikal untuk seluruh rule di kategori theme:
// 1. Format ID kanonikal (<category>.<slug>)
// 2. Kategori wajib "theme"
// 3. Deskripsi non-kosong
// 4. Kelengkapan dokumentasi 8-Pillars
// 5. Fail-safe invariant: evaluasi node nil atau kosong tidak boleh panic dan mengembalikan 0 diagnostik.
func TestThemeRules_CanonicalContract(t *testing.T) {
	reg := rules.NewRegistry()
	if err := rules.RegisterBuiltinRules(reg); err != nil {
		t.Fatalf("failed to register builtin rules: %v", err)
	}

	themeRules := make([]rules.Rule, 0)
	for _, r := range reg.All() {
		if r.Category() == "theme" {
			themeRules = append(themeRules, r)
		}
	}

	if len(themeRules) == 0 {
		t.Fatalf("no theme rules registered in builtin registry")
	}

	for _, rule := range themeRules {
		t.Run(rule.ID(), func(t *testing.T) {
			// 1. ID check
			if !strings.HasPrefix(rule.ID(), "theme.") || len(rule.ID()) <= 6 {
				t.Errorf("rule ID %q does not conform to canonical format 'theme.<slug>'", rule.ID())
			}

			// 2. Category check
			if rule.Category() != "theme" {
				t.Errorf("rule %s has category %q, want 'theme'", rule.ID(), rule.Category())
			}

			// 3. Description check
			if len(strings.TrimSpace(rule.Description())) == 0 {
				t.Errorf("rule %s has empty description", rule.ID())
			}

			// 4. Severity check
			sev := rule.DefaultSeverity()
			if sev != ir.SeverityError && sev != ir.SeverityWarn && sev != ir.SeverityInfo {
				t.Errorf("rule %s has invalid severity: %v", rule.ID(), sev)
			}

			// 5. 8-Pillars Documentation check
			docRule, ok := rule.(interface{ Doc() ir.RuleDocumentation })
			if !ok {
				t.Errorf("rule %s does not implement Doc() ir.RuleDocumentation", rule.ID())
			} else {
				doc := docRule.Doc()
				if len(doc.TargetStandards) == 0 {
					t.Errorf("rule %s missing TargetStandards in Doc", rule.ID())
				}
				if len(strings.TrimSpace(doc.CoreInvariant)) == 0 {
					t.Errorf("rule %s missing CoreInvariant in Doc", rule.ID())
				}
				if len(strings.TrimSpace(doc.Grounding)) == 0 {
					t.Errorf("rule %s missing Grounding in Doc", rule.ID())
				}
				if len(doc.BadExamples) == 0 {
					t.Errorf("rule %s missing BadExamples in Doc", rule.ID())
				}
				if len(doc.GoodExamples) == 0 {
					t.Errorf("rule %s missing GoodExamples in Doc", rule.ID())
				}
				if len(doc.Risks) == 0 {
					t.Errorf("rule %s missing Risks in Doc", rule.ID())
				}
			}

			// 6. Fail-safe edge cases
			if diags := rule.Evaluate(nil); len(diags) != 0 {
				t.Errorf("rule %s returned %d diagnostics on nil node, want 0", rule.ID(), len(diags))
			}
			if diags := rule.Evaluate(&ir.Node{}); len(diags) != 0 {
				t.Errorf("rule %s returned %d diagnostics on empty node, want 0", rule.ID(), len(diags))
			}
		})
	}
}
