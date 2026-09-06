package inp

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// UnboundedEffectDepsRule mendeteksi deklarasi useEffect atau useLayoutEffect tanpa array dependensi
// yang memicu re-eksekusi efek dan lonjakan CPU di setiap siklus render frame.
type UnboundedEffectDepsRule struct{}

// NewUnboundedEffectDepsRule membuat instance baru dari UnboundedEffectDepsRule.
func NewUnboundedEffectDepsRule() *UnboundedEffectDepsRule {
	return &UnboundedEffectDepsRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *UnboundedEffectDepsRule) ID() string {
	return "inp.unbounded-effect-deps"
}

// Description mengembalikan ringkasan aturan.
func (r *UnboundedEffectDepsRule) Description() string {
	return "Lifecycle hook useEffect/useLayoutEffect is missing a dependency array, triggering unbounded re-executions on every render"
}

// Category mengembalikan nama kategori rule.
func (r *UnboundedEffectDepsRule) Category() string {
	return "inp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *UnboundedEffectDepsRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *UnboundedEffectDepsRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"React Hooks Specification & Dependency Determinism",
			"W3C Cooperative Scheduling & Frame Budget Invariants",
			"Google Chrome Core Web Vitals (Input Presentation Delay)",
		},
		CoreInvariant: "React lifecycle hooks (useEffect, useLayoutEffect) must explicitly declare a dependency array as their second argument to prevent uncontrolled execution on every render cycle.",
		Grounding: "When 'useEffect' or 'useLayoutEffect' is invoked with only a callback and no second argument, React executes the effect after *every single render*.\n\n" +
			"Any state update, parent re-render, or user keystroke causes the entire effect callback to run again. If the effect queries DOM elements, reads layout properties, or synchronizes subscriptions, the main thread is constantly saturated by unnecessary computations.\n\n" +
			"Providing an explicit dependency array ('[]' for mount-only, or '[deps...]') restricts execution strictly to when dependencies change, protecting interaction frame rate.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Every-Render Effect Re-execution",
				Severity: "CRITICAL",
				Impact:   "Effects fire repeatedly on every keystroke or state change, causing severe CPU spikes and input lag.",
			},
			{
				Vector:   "Infinite Render Loops",
				Severity: "HIGH",
				Impact:   "If an unbounded effect updates state, it causes an immediate infinite re-render loop that locks the browser tab.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "useEffect without a dependency array executes on every render",
				Code: `useEffect(() => {
  recomputeHeavyLayout();
});`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Explicit empty dependency array ensures execution only on mount",
				Code: `useEffect(() => {
  recomputeHeavyLayout();
}, []);`,
			},
		},
	}
}

// Evaluate memeriksa apakah terdapat pemanggilan hook effect tanpa array dependensi.
func (r *UnboundedEffectDepsRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil {
		return nil
	}

	code := getScriptOrSourceContent(node)
	if code == "" {
		return nil
	}

	hook, line, missing := findUnboundedEffectDeps(code)
	if missing {
		targetLine := node.Span.Line
		if line > 0 {
			targetLine = line
		}
		return []ir.Diagnostic{
			{
				Line:     targetLine,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Hook '%s' is declared without a dependency array argument. Effects without dependencies re-run on every single render cycle, creating severe CPU spikes that delay frame presentation and input handling.", hook),
				Hint:     fmt.Sprintf("Pass a dependency array (e.g. '[]' for mount-only execution, or '[deps...]') as the second argument to '%s'.", hook),
			},
		}
	}

	return nil
}
