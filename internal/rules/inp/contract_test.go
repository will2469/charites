package inp_test

import (
	"strings"
	"testing"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules"
)

// TestINPRules_CanonicalContract menguji kepatuhan kontrak kanonikal seluruh rule di kategori inp:
// 1. Format ID kanonikal (inp.<slug>)
// 2. Kategori wajib "inp"
// 3. Deskripsi non-kosong
// 4. Tingkat keparahan sah (Error, Warn, Info)
// 5. Kelengkapan dokumentasi 8-Pillars
// 6. Fail-safe invariant: evaluasi node nil atau kosong tidak boleh panic dan mengembalikan 0 diagnostik.
func TestINPRules_CanonicalContract(t *testing.T) {
	reg := rules.NewRegistry()
	if err := rules.RegisterBuiltinRules(reg); err != nil {
		t.Fatalf("failed to register builtin rules: %v", err)
	}

	inpRules := getINPRules(reg)
	if len(inpRules) == 0 {
		t.Fatalf("no inp rules registered in builtin registry")
	}

	for _, rule := range inpRules {
		t.Run(rule.ID(), func(t *testing.T) {
			assertRuleBasics(t, rule)
			assertDocContract(t, rule)
			assertFailSafeContract(t, rule)
		})
	}
}

func getINPRules(reg *rules.Registry) []rules.Rule {
	var res []rules.Rule
	for _, r := range reg.All() {
		if r.Category() == "inp" {
			res = append(res, r)
		}
	}
	return res
}

func assertRuleBasics(t *testing.T, rule rules.Rule) {
	if !strings.HasPrefix(rule.ID(), "inp.") || len(rule.ID()) <= 4 {
		t.Errorf("rule ID %q does not conform to canonical format 'inp.<slug>'", rule.ID())
	}
	if rule.Category() != "inp" {
		t.Errorf("rule %s has category %q, want 'inp'", rule.ID(), rule.Category())
	}
	if len(strings.TrimSpace(rule.Description())) == 0 {
		t.Errorf("rule %s has empty description", rule.ID())
	}
	sev := rule.DefaultSeverity()
	if sev != ir.SeverityError && sev != ir.SeverityWarn && sev != ir.SeverityInfo {
		t.Errorf("rule %s has invalid severity: %v", rule.ID(), sev)
	}
}

type docProvider interface {
	Doc() ir.RuleDocumentation
}

func assertDocContract(t *testing.T, rule rules.Rule) {
	docRule, ok := rule.(docProvider)
	if !ok {
		t.Errorf("rule %s does not implement Doc() method for 8-Pillars wiki generation", rule.ID())
		return
	}
	doc := docRule.Doc()
	if len(doc.TargetStandards) == 0 {
		t.Errorf("rule %s has empty TargetStandards in Doc()", rule.ID())
	}
	if len(strings.TrimSpace(doc.CoreInvariant)) == 0 {
		t.Errorf("rule %s has empty CoreInvariant in Doc()", rule.ID())
	}
	if len(strings.TrimSpace(doc.Grounding)) == 0 {
		t.Errorf("rule %s has empty Grounding in Doc()", rule.ID())
	}
	if len(doc.Risks) == 0 {
		t.Errorf("rule %s has empty Risks in Doc()", rule.ID())
	}
	if len(doc.BadExamples) == 0 {
		t.Errorf("rule %s has empty BadExamples in Doc()", rule.ID())
	}
	if len(doc.GoodExamples) == 0 {
		t.Errorf("rule %s has empty GoodExamples in Doc()", rule.ID())
	}
}

func assertFailSafeContract(t *testing.T, rule rules.Rule) {
	// 1. Nil node
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("rule %s panicked on nil node: %v", rule.ID(), r)
		}
	}()
	if diags := rule.Evaluate(nil); len(diags) != 0 {
		t.Errorf("rule %s returned %d diagnostics on nil node, want 0", rule.ID(), len(diags))
	}

	// 2. Empty node
	emptyNode := &ir.Node{}
	if diags := rule.Evaluate(emptyNode); len(diags) != 0 {
		t.Errorf("rule %s returned %d diagnostics on empty node, want 0", rule.ID(), len(diags))
	}
}
