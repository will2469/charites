package ux_test

import (
	"strings"
	"testing"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules"
)

// TestUXRules_CanonicalContract menguji kontrak kanonikal untuk seluruh rule di kategori ux:
// 1. Format ID kanonikal (ux.<slug>)
// 2. Kategori wajib "ux"
// 3. Deskripsi non-kosong
// 4. Tingkat keparahan sah (Error, Warn, Info)
// 5. Kelengkapan dokumentasi 8-Pillars (TargetStandards, CoreInvariant, Grounding, BadExamples, GoodExamples, Risks)
// 6. Fail-safe invariant: evaluasi node nil atau kosong tidak boleh panic dan mengembalikan 0 diagnostik.
func TestUXRules_CanonicalContract(t *testing.T) {
	reg := rules.NewRegistry()
	if err := rules.RegisterBuiltinRules(reg); err != nil {
		t.Fatalf("failed to register builtin rules: %v", err)
	}

	uxRules := make([]rules.Rule, 0)
	for _, r := range reg.All() {
		if r.Category() == "ux" {
			uxRules = append(uxRules, r)
		}
	}

	if len(uxRules) == 0 {
		t.Fatalf("no ux rules registered in builtin registry")
	}

	for _, rule := range uxRules {
		t.Run(rule.ID(), func(t *testing.T) {
			// 1. ID check
			if !strings.HasPrefix(rule.ID(), "ux.") || len(rule.ID()) <= 3 {
				t.Errorf("rule ID %q does not conform to canonical format 'ux.<slug>'", rule.ID())
			}

			// 2. Category check
			if rule.Category() != "ux" {
				t.Errorf("rule %s has category %q, want 'ux'", rule.ID(), rule.Category())
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
				t.Fatalf("rule %s does not implement Doc() ir.RuleDocumentation (required for 8-Pillars wiki)", rule.ID())
			}
			doc := docRule.Doc()

			if len(doc.TargetStandards) == 0 {
				t.Errorf("rule %s has empty TargetStandards", rule.ID())
			}
			if len(strings.TrimSpace(doc.CoreInvariant)) == 0 {
				t.Errorf("rule %s has empty CoreInvariant", rule.ID())
			}
			if len(strings.TrimSpace(doc.Grounding)) == 0 {
				t.Errorf("rule %s has empty Grounding", rule.ID())
			}
			if len(doc.Risks) == 0 {
				t.Errorf("rule %s has empty Risks", rule.ID())
			}
			if len(doc.BadExamples) == 0 {
				t.Errorf("rule %s has empty BadExamples", rule.ID())
			}
			if len(doc.GoodExamples) == 0 {
				t.Errorf("rule %s has empty GoodExamples", rule.ID())
			}

			// 6. Fail-safe invariant
			func() {
				defer func() {
					if rec := recover(); rec != nil {
						t.Errorf("rule %s panicked on nil node: %v", rule.ID(), rec)
					}
				}()
				diags := rule.Evaluate(nil)
				if len(diags) != 0 {
					t.Errorf("rule %s returned %d diagnostics on nil node, want 0", rule.ID(), len(diags))
				}
			}()

			func() {
				defer func() {
					if rec := recover(); rec != nil {
						t.Errorf("rule %s panicked on empty node: %v", rule.ID(), rec)
					}
				}()
				emptyNode := &ir.Node{}
				diags := rule.Evaluate(emptyNode)
				if len(diags) != 0 {
					t.Errorf("rule %s returned %d diagnostics on empty node, want 0", rule.ID(), len(diags))
				}
			}()
		})
	}
}
