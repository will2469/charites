package performance

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// ReactContextDomainCouplingRule mendeteksi Context.Provider yang memuat objek nilai
// dengan jumlah properti berlebih (> 5 field) atau menggabungkan domain state berbeda.
type ReactContextDomainCouplingRule struct{}

// NewReactContextDomainCouplingRule membuat instance baru dari ReactContextDomainCouplingRule.
func NewReactContextDomainCouplingRule() *ReactContextDomainCouplingRule {
	return &ReactContextDomainCouplingRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *ReactContextDomainCouplingRule) ID() string {
	return "performance.react-context-domain-coupling"
}

// Description mengembalikan ringkasan aturan.
func (r *ReactContextDomainCouplingRule) Description() string {
	return "Context.Provider bundles over-coupled multi-domain state (> 5 fields), triggering cascading re-renders across all consumers on any property change"
}

// Category mengembalikan nama kategori rule.
func (r *ReactContextDomainCouplingRule) Category() string {
	return "performance"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *ReactContextDomainCouplingRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ReactContextDomainCouplingRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"React Official Architecture (Context Modularity & State Granularity Principles)",
			"React Virtual DOM Reconciliation Invariants (Consumer tree re-render pruning)",
			"Web Performance Working Group Main-Thread Optimization Guidelines",
		},
		CoreInvariant: "React Context Providers must maintain fine-grained domain boundaries; bundling disparate application states into a monolithic 'God Context' forces all consumer components to re-render whenever any unrelated field updates.",
		Grounding: "React Context propagates state updates to all consumers without granular field-level selector filtering. When an application bundles multiple disparate domains (e.g. user authentication, shopping cart, UI modal state, notification badge count, and scroll position) into a single Provider, any update to a high-frequency field (such as a badge increment) triggers a re-render across every component subscribed to the context.\n\n" +
			"This monolithic coupling bypasses component memoization and saturates the main thread with unnecessary render passes.\n\n" +
			"Splitting state into focused, domain-specific contexts (e.g. `AuthContext`, `CartContext`, `UIContext`) ensures components only re-render when their specific domain data actually mutates.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Cascading Re-render Blast Radius",
				Severity: "HIGH",
				Impact:   "A change in a single property (e.g. notification count) forces dozens of unrelated UI components to re-render simultaneously.",
			},
			{
				Vector:   "Severe Main-Thread Jitter",
				Severity: "MEDIUM",
				Impact:   "High-frequency state mutations throttle the main thread, resulting in dropped frames and delayed user interaction responses.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Monolithic God Context bundling 7 distinct domain properties in a single provider",
				Code: `<AppContext.Provider value={{
  user,
  cart,
  theme,
  notifications,
  activeModal,
  isSidebarOpen,
  locale,
}}>
  {children}
</AppContext.Provider>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "State decoupled into domain-specific, modular context providers",
				Code: `<AuthProvider value={authValue}>
  <CartProvider value={cartValue}>
    <UIProvider value={uiValue}>
      {children}
    </UIProvider>
  </CartProvider>
</AuthProvider>`,
			},
		},
	}
}

// Evaluate memeriksa apakah Context.Provider memuat objek nilai yang ter-overcouple.
func (r *ReactContextDomainCouplingRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	propCount, isOvercoupled := isOvercoupledContextProvider(node)
	if !isOvercoupled {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Context provider '<%s>' bundles %d properties into its 'value' object. Monolithic 'God Context' providers force all subscribing components to re-render on any field update, regardless of whether that component needs the changed data.", node.Tag, propCount),
			Hint:     "Decompose the monolithic provider into domain-specific, modular contexts (e.g. 'AuthContext', 'CartContext', 'UIContext') so components only re-render for state they actually consume.",
		},
	}
}
