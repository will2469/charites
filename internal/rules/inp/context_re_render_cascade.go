package inp

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// ContextReRenderCascadeRule mendeteksi instansiasi objek literal inline pada prop value
// di Context.Provider yang memicu cascade re-render seluruh pohon komponen konsumen.
type ContextReRenderCascadeRule struct{}

// NewContextReRenderCascadeRule membuat instance baru dari ContextReRenderCascadeRule.
func NewContextReRenderCascadeRule() *ContextReRenderCascadeRule {
	return &ContextReRenderCascadeRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *ContextReRenderCascadeRule) ID() string {
	return "inp.context-re-render-cascade"
}

// Description mengembalikan ringkasan aturan.
func (r *ContextReRenderCascadeRule) Description() string {
	return "Passing an unmemoized inline object literal to Context.Provider value triggers cascading re-renders across all consumers"
}

// Category mengembalikan nama kategori rule.
func (r *ContextReRenderCascadeRule) Category() string {
	return "inp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *ContextReRenderCascadeRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ContextReRenderCascadeRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"React Context API Referential Equality Invariants",
			"React Virtual DOM Reconciliation & Tree Pruning Standards",
			"Google Chrome Core Web Vitals (Input to Next Paint Interaction Optimization)",
		},
		CoreInvariant: "React Context Provider 'value' props must receive referentially stable references (via 'useMemo' or external constants) rather than freshly instantiated inline object literals.",
		Grounding: "React determines whether consumer components of a Context must re-render by performing a strict reference equality check ('Object.is(prevValue, nextValue)').\n\n" +
			"When an inline object literal ('value={{ user, token }}') is passed directly into a Provider, a completely new object in heap memory is allocated on *every single render* of the parent component.\n\n" +
			"Because the memory reference changes each time, React bypasses 'React.memo' optimizations in all descendant consumers, forcing the entire subtree to re-render simultaneously and causing severe interaction lag on user input.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Cascading Consumer Re-renders",
				Severity: "HIGH",
				Impact:   "Every component consuming the context is forced to re-render on any parent state change.",
			},
			{
				Vector:   "Heap Allocation & GC Churn",
				Severity: "MEDIUM",
				Impact:   "Repeated object allocations in the render path trigger garbage collection pauses during rapid user interactions.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Fresh inline object allocated on every render triggers full consumer re-renders",
				Code: `<AuthContext.Provider value={{ user, isAuthenticated, login }}>
  {children}
</AuthContext.Provider>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Context value wrapped in useMemo to preserve reference equality",
				Code: `const authValue = useMemo(() => ({ user, isAuthenticated, login }), [user, isAuthenticated]);
return (
  <AuthContext.Provider value={authValue}>
    {children}
  </AuthContext.Provider>
);`,
			},
		},
	}
}

// Evaluate memeriksa apakah Context.Provider menerima objek literal inline pada prop value.
func (r *ContextReRenderCascadeRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	providerTag, isInline := isInlineContextObjectValue(node.Tag, node.Attributes)
	if isInline {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Context provider '%s' passes a newly instantiated object literal directly to the 'value' prop. This creates a new object reference on every render, forcing all consumer components to re-render and spiking interaction latency.", providerTag),
				Hint:     "Wrap the context value in 'useMemo(() => ({ ... }), [dependencies])' to preserve referential stability across renders.",
			},
		}
	}

	return nil
}
