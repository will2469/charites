package performance

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// ReactDerivedStateInEffectRule mendeteksi sinkronisasi derived state via useEffect yang memicu render sekunder ganda.
type ReactDerivedStateInEffectRule struct{}

// NewReactDerivedStateInEffectRule membuat instance baru dari ReactDerivedStateInEffectRule.
func NewReactDerivedStateInEffectRule() *ReactDerivedStateInEffectRule {
	return &ReactDerivedStateInEffectRule{}
}

// ID mengembalikan identifier unik kanonikal aturan.
func (r *ReactDerivedStateInEffectRule) ID() string {
	return "performance.react-derived-state-in-effect"
}

// Category mengembalikan kategori aturan ('performance').
func (r *ReactDerivedStateInEffectRule) Category() string {
	return "performance"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *ReactDerivedStateInEffectRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Description mengembalikan deskripsi tujuan dan konteks aturan.
func (r *ReactDerivedStateInEffectRule) Description() string {
	return "Mencegah sinkronisasi derived state dari props atau state yang sudah ada melalui useEffect, yang memicu siklus perenderan sekunder ganda."
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ReactDerivedStateInEffectRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"React Official Documentation ('You Might Not Need an Effect')",
			"React Reconciliation Lifecycle (Avoiding Cascading Secondary Renders)",
			"React Best Practices on Pure Render-Phase Computations",
		},
		CoreInvariant: "Derived values computed synchronously from props or existing state must be calculated directly in the component body during render; updating state inside 'useEffect' triggers redundant secondary render passes.",
		Grounding: "Updating state within a `useEffect` callback causes React to first render the component with the stale value, commit it to the DOM, and immediately schedule a second render pass to apply the updated state.\n\n" +
			"When derived state is calculated synchronously (e.g. concatenating names, filtering a list, or calculating a total), calculating it in an effect needlessly burns main thread time on layout calculations and DOM diffing twice per interaction.\n\n" +
			"Computing the value directly during the render pass completely eliminates the secondary render pass and keeps state management minimal.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Cascading Secondary Render Cycles",
				Severity: "HIGH",
				Impact:   "Forces React to run duplicate render and diff cycles on every prop change, directly degrading interaction responsiveness (INP).",
			},
			{
				Vector:   "Visual Stutter / Layout Shift",
				Severity: "MEDIUM",
				Impact:   "May momentarily display stale computed values before the effect updates, causing brief content flicker or layout shifts.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Sinkronisasi derived state via useEffect memicu render ganda",
				Code: `const [fullName, setFullName] = useState('');
useEffect(() => {
  setFullName(firstName + ' ' + lastName);
}, [firstName, lastName]);`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Dihitung secara sinkron dalam satu kali fase render",
				Code:     `const fullName = firstName + ' ' + lastName;`,
			},
		},
	}
}

// Evaluate memeriksa apakah useEffect digunakan untuk sinkronisasi nilai turunan secara murni.
func (r *ReactDerivedStateInEffectRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || !isSourceRootOrScript(node) {
		return nil
	}

	fileSrc := getFileSourceContent(node)
	if len(fileSrc) == 0 {
		return nil
	}

	violations := findDerivedStateInEffects(fileSrc)
	if len(violations) == 0 {
		return nil
	}

	diags := make([]ir.Diagnostic, 0, len(violations))
	for _, v := range violations {
		diags = append(diags, ir.Diagnostic{
			Line:     v.Line,
			Column:   1,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Hook '%s' synchronizes derived state via '%s' from existing props or state. Updating state inside an effect forces a redundant secondary render cycle; compute derived values synchronously during the render phase.", v.EffectName, v.StateName),
			Hint:     fmt.Sprintf("Eliminate '%s' and derive the state synchronously directly in the component body (e.g. 'const derivedValue = ...').", v.EffectName),
		})
	}

	return diags
}
