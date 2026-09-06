package performance

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// ReactInlinePropMemoRule mendeteksi pengiriman inline object/array/function
// pada pemanggilan komponen yang dibungkus React.memo / memo(...).
type ReactInlinePropMemoRule struct{}

// NewReactInlinePropMemoRule membuat instance baru dari ReactInlinePropMemoRule.
func NewReactInlinePropMemoRule() *ReactInlinePropMemoRule {
	return &ReactInlinePropMemoRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *ReactInlinePropMemoRule) ID() string {
	return "performance.react-inline-prop-memo"
}

// Description mengembalikan ringkasan aturan.
func (r *ReactInlinePropMemoRule) Description() string {
	return "Passing inline object, array, or function literal to memoized component bypasses shallow memoization on every parent render"
}

// Category mengembalikan nama kategori rule.
func (r *ReactInlinePropMemoRule) Category() string {
	return "performance"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *ReactInlinePropMemoRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ReactInlinePropMemoRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"React Official Architecture (Reconciliation Identity & Referential Equality Invariants)",
			"React Memoization Contract ('React.memo' shallow prop comparison specification)",
			"Web Performance Working Group VDOM Re-render Minimization Guidelines",
		},
		CoreInvariant: "Memoized React components must receive referentially stable props; inline object literals, array literals, and arrow functions allocate new heap memory on each parent render pass, completely nullifying 'React.memo'.",
		Grounding: "When a component is wrapped in `React.memo()`, React evaluates whether to skip rendering by performing shallow equality checks (`prevProps[key] === nextProps[key]`) across all incoming props.\n\n" +
			"Passing an inline object literal (`prop={{ ... }}`), inline array (`prop={[ ... ]}`), or arrow function (`prop={() => ...}`) directly at the JSX call-site instantiates a brand-new memory reference on every parent render.\n\n" +
			"Because `===` on different object references always evaluates to `false`, React is forced to re-render the memoized component every time, incurring the performance penalty of shallow comparison without receiving any of its caching benefits.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Bypassed Component Memoization",
				Severity: "HIGH",
				Impact:   "Forces expensive memoized component trees to re-render unnecessarily on every parent state mutation.",
			},
			{
				Vector:   "Garbage Collection Churn",
				Severity: "MEDIUM",
				Impact:   "Allocates short-lived transient objects and closures in heap memory during rapid interaction loops.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Memoized UserCard receives inline object and arrow function props",
				Code: `const UserCard = React.memo(({ user, config, onSelect }: UserCardProps) => {
  return <div>{user.name}</div>;
});

function Parent({ currentUser }: { currentUser: User }) {
  return (
    <UserCard
      user={currentUser}
      config={{ theme: 'dark', compact: true }}
      onSelect={() => console.log('selected')}
    />
  );
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Props stabilized via external constants or useCallback/useMemo hooks",
				Code: `const USER_CONFIG = { theme: 'dark', compact: true } as const;

function Parent({ currentUser }: { currentUser: User }) {
  const handleSelect = useCallback(() => {
    console.log('selected');
  }, []);

  return (
    <UserCard
      user={currentUser}
      config={USER_CONFIG}
      onSelect={handleSelect}
    />
  );
}`,
			},
		},
	}
}

// Evaluate memeriksa apakah pemanggilan komponen memoized menerima prop inline.
func (r *ReactInlinePropMemoRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || len(node.Attributes) == 0 {
		return nil
	}

	// Hanya audit komponen kustom berhuruf kapital
	if len(node.Tag) == 0 || node.Tag[0] < 'A' || node.Tag[0] > 'Z' {
		return nil
	}

	fileSrc := getFileSourceContent(node)
	if !isMemoizedComponent(node.Tag, fileSrc) {
		return nil
	}

	var diags []ir.Diagnostic
	for attrName, attrVal := range node.Attributes {
		if kind, isInline := isInlineLiteralProp(attrVal); isInline {
			diags = append(diags, ir.Diagnostic{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Memoized component '<%s>' receives inline %s on prop '%s'. Fresh references allocated at the call-site defeat 'React.memo' shallow equality, forcing re-renders on every parent pass.", node.Tag, kind, attrName),
				Hint:     "Stabilize the prop reference using 'useMemo' (for objects/arrays), 'useCallback' (for functions), or lift it out as a module-level constant.",
			})
		}
	}

	return diags
}
