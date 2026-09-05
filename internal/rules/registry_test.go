package rules_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules"
)

type mockRule struct {
	id          string
	description string
	category    string
	severity    ir.Severity
}

func (m *mockRule) ID() string                             { return m.id }
func (m *mockRule) Description() string                    { return m.description }
func (m *mockRule) Category() string                       { return m.category }
func (m *mockRule) DefaultSeverity() ir.Severity           { return m.severity }
func (m *mockRule) Evaluate(node *ir.Node) []ir.Diagnostic { return nil }

func TestRegistry_RegisterAndGet(t *testing.T) {
	reg := rules.NewRegistry()

	r1 := &mockRule{
		id:          "theme.mock-rule",
		description: "A test rule",
		category:    "theme",
		severity:    ir.SeverityError,
	}

	if err := reg.Register(r1); err != nil {
		t.Fatalf("unexpected error registering rule: %v", err)
	}

	if reg.Count() != 1 {
		t.Fatalf("expected count 1, got %d", reg.Count())
	}

	found, ok := reg.Get("theme.mock-rule")
	if !ok {
		t.Fatalf("expected rule to be found")
	}
	if found.ID() != r1.ID() {
		t.Errorf("expected ID %q, got %q", r1.ID(), found.ID())
	}

	_, ok = reg.Get("nonexistent.rule")
	if ok {
		t.Errorf("expected nonexistent rule to not be found")
	}
}

func TestRegistry_DuplicateIDRejected(t *testing.T) {
	reg := rules.NewRegistry()

	r1 := &mockRule{id: "theme.duplicate", category: "theme"}
	r2 := &mockRule{id: "theme.duplicate", category: "theme"}

	if err := reg.Register(r1); err != nil {
		t.Fatalf("expected first registration to succeed: %v", err)
	}

	if err := reg.Register(r2); err == nil {
		t.Fatalf("expected duplicate registration to fail")
	}
}

func TestRegistry_Errors(t *testing.T) {
	reg := rules.NewRegistry()

	if err := reg.Register(nil); err != rules.ErrNilRule {
		t.Errorf("expected ErrNilRule, got %v", err)
	}

	emptyIDRule := &mockRule{id: "", category: "theme"}
	if err := reg.Register(emptyIDRule); err != rules.ErrEmptyRuleID {
		t.Errorf("expected ErrEmptyRuleID, got %v", err)
	}
}

func TestRegistry_DeterministicOrder(t *testing.T) {
	reg := rules.NewRegistry()

	// Register rules out of order
	ids := []string{
		"theme.hardcode-palette-color",
		"a11y.missing-alt",
		"responsive.touch-target",
		"theme.hardcode-opacity-color",
		"perf.inline-css",
	}

	for _, id := range ids {
		cat := id[:5]
		if err := reg.Register(&mockRule{id: id, category: cat}); err != nil {
			t.Fatalf("failed to register %s: %v", id, err)
		}
	}

	all := reg.All()
	if len(all) != len(ids) {
		t.Fatalf("expected %d rules, got %d", len(ids), len(all))
	}

	expectedAll := []string{
		"a11y.missing-alt",
		"perf.inline-css",
		"responsive.touch-target",
		"theme.hardcode-opacity-color",
		"theme.hardcode-palette-color",
	}

	for i, r := range all {
		if r.ID() != expectedAll[i] {
			t.Errorf("index %d: expected %s, got %s", i, expectedAll[i], r.ID())
		}
	}

	// Test ByCategory determinism
	themeRules := reg.ByCategory("theme")
	if len(themeRules) != 2 {
		t.Fatalf("expected 2 theme rules, got %d", len(themeRules))
	}
	if themeRules[0].ID() != "theme.hardcode-opacity-color" {
		t.Errorf("expected theme.hardcode-opacity-color first, got %s", themeRules[0].ID())
	}
	if themeRules[1].ID() != "theme.hardcode-palette-color" {
		t.Errorf("expected theme.hardcode-palette-color second, got %s", themeRules[1].ID())
	}

	// Verify defensive copy
	themeRules[0] = nil
	checkAgain := reg.ByCategory("theme")
	if checkAgain[0] == nil {
		t.Errorf("registry internal slice was mutated by external slice modification")
	}
}

func TestRegistry_ConcurrentSafety(t *testing.T) {
	reg := rules.NewRegistry()
	var wg sync.WaitGroup

	// Pre-seed some rules
	for i := 0; i < 10; i++ {
		_ = reg.Register(&mockRule{
			id:       fmt.Sprintf("theme.seed-%02d", i),
			category: "theme",
		})
	}

	// Concurrent readers
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_ = reg.All()
				_ = reg.ByCategory("theme")
				_ = reg.Count()
				_, _ = reg.Get("theme.seed-05")
			}
		}()
	}

	// Concurrent writers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		id := fmt.Sprintf("custom.writer-%02d", i)
		go func(rid string) {
			defer wg.Done()
			_ = reg.Register(&mockRule{id: rid, category: "custom"})
		}(id)
	}

	wg.Wait()

	if reg.Count() < 10 {
		t.Errorf("expected count >= 10, got %d", reg.Count())
	}
}

func TestDefaultRegistry_BuiltinAndHelpers(t *testing.T) {
	// DefaultRegistry should already have built-in rules registered via init()
	def := rules.DefaultRegistry()
	if def == nil {
		t.Fatalf("DefaultRegistry returned nil")
	}

	r, ok := rules.Get("theme.hardcode-opacity-color")
	if !ok {
		t.Fatalf("expected theme.hardcode-opacity-color in default registry")
	}
	if r.Category() != "theme" {
		t.Errorf("expected category theme, got %s", r.Category())
	}

	all := rules.All()
	if len(all) == 0 {
		t.Errorf("expected non-empty rules.All()")
	}

	themes := rules.ByCategory("theme")
	if len(themes) == 0 {
		t.Errorf("expected non-empty rules.ByCategory(\"theme\")")
	}

	if rules.Count() == 0 {
		t.Errorf("expected rules.Count() > 0")
	}

	// Test RegisterBuiltinRules helper
	freshReg := rules.NewRegistry()
	if err := rules.RegisterBuiltinRules(nil); err != rules.ErrNilRule {
		t.Errorf("expected ErrNilRule, got %v", err)
	}
	if err := rules.RegisterBuiltinRules(freshReg); err != nil {
		t.Fatalf("RegisterBuiltinRules failed: %v", err)
	}
	if freshReg.Count() != 1 {
		t.Errorf("expected 1 rule in freshReg, got %d", freshReg.Count())
	}
}
