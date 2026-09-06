package pwa_test

import (
	"strings"
	"testing"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules"
)

// TestPwaRules_CanonicalContract menguji kontrak kanonikal untuk seluruh rule di kategori pwa:
// 1. Format ID kanonikal (pwa.<slug>)
// 2. Kategori wajib "pwa"
// 3. Deskripsi non-kosong
// 4. Tingkat keparahan sah (Error, Warn, Info)
// 5. Kelengkapan dokumentasi 8-Pillars
// 6. Fail-safe invariant: evaluasi node nil atau kosong tidak boleh panic dan mengembalikan 0 diagnostik.
func TestPwaRules_CanonicalContract(t *testing.T) {
	reg := rules.NewRegistry()
	if err := rules.RegisterBuiltinRules(reg); err != nil {
		t.Fatalf("failed to register builtin rules: %v", err)
	}

	pwaRules := getPwaRules(reg)
	if len(pwaRules) == 0 {
		t.Fatalf("no pwa rules registered in builtin registry")
	}

	for _, rule := range pwaRules {
		t.Run(rule.ID(), func(t *testing.T) {
			assertRuleBasics(t, rule)
			assertDocContract(t, rule)
			assertFailSafeContract(t, rule)
		})
	}
}

func getPwaRules(reg *rules.Registry) []rules.Rule {
	var res []rules.Rule
	for _, r := range reg.All() {
		if r.Category() == "pwa" {
			res = append(res, r)
		}
	}
	return res
}

func assertRuleBasics(t *testing.T, rule rules.Rule) {
	if !strings.HasPrefix(rule.ID(), "pwa.") || len(rule.ID()) <= 4 {
		t.Errorf("rule ID %q does not conform to canonical format 'pwa.<slug>'", rule.ID())
	}
	if rule.Category() != "pwa" {
		t.Errorf("rule %s has category %q, want 'pwa'", rule.ID(), rule.Category())
	}
	if len(strings.TrimSpace(rule.Description())) == 0 {
		t.Errorf("rule %s has empty description", rule.ID())
	}
	sev := rule.DefaultSeverity()
	if sev != ir.SeverityError && sev != ir.SeverityWarn && sev != ir.SeverityInfo {
		t.Errorf("rule %s has invalid severity: %v", rule.ID(), sev)
	}
}

func assertDocContract(t *testing.T, rule rules.Rule) {
	docRule, ok := rule.(interface{ Doc() ir.RuleDocumentation })
	if !ok {
		t.Errorf("rule %s does not implement Doc() ir.RuleDocumentation", rule.ID())
		return
	}
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

func assertFailSafeContract(t *testing.T, rule rules.Rule) {
	if diags := rule.Evaluate(nil); len(diags) != 0 {
		t.Errorf("rule %s returned %d diagnostics on nil node, want 0", rule.ID(), len(diags))
	}
	if diags := rule.Evaluate(&ir.Node{}); len(diags) != 0 {
		t.Errorf("rule %s returned %d diagnostics on empty node, want 0", rule.ID(), len(diags))
	}
}
