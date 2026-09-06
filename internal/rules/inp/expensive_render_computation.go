package inp

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// ExpensiveRenderComputationRule mendeteksi operasi komputasi derivasi data koleksi berat
// (seperti rantai .filter().sort()) yang dieksekusi langsung di badan render tanpa useMemo.
type ExpensiveRenderComputationRule struct{}

// NewExpensiveRenderComputationRule membuat instance baru dari ExpensiveRenderComputationRule.
func NewExpensiveRenderComputationRule() *ExpensiveRenderComputationRule {
	return &ExpensiveRenderComputationRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *ExpensiveRenderComputationRule) ID() string {
	return "inp.expensive-render-computation"
}

// Description mengembalikan ringkasan aturan.
func (r *ExpensiveRenderComputationRule) Description() string {
	return "Expensive data transformations (chained .filter() and .sort()) execute synchronously in the render path without useMemo"
}

// Category mengembalikan nama kategori rule.
func (r *ExpensiveRenderComputationRule) Category() string {
	return "inp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *ExpensiveRenderComputationRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ExpensiveRenderComputationRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"React Render Phase Performance Optimization Principles",
			"W3C Cooperative Scheduling & Frame Execution Invariants",
			"Google Chrome Core Web Vitals (Input Processing Delay)",
		},
		CoreInvariant: "Heavy collection derivations involving sequential filtering and sorting in component render paths must be memoized using 'useMemo' to prevent recomputation on every user keystroke.",
		Grounding: "The body of a functional React component executes synchronously on every render cycle-including every keystroke inside controlled form fields or hover interactions.\n\n" +
			"When developers write heavy array transformations (such as 'users.filter(...).sort(...)') directly within the render body or inside JSX props without 'useMemo', the browser re-filters and re-sorts the entire collection on every single frame.\n\n" +
			"Wrapping the computation in 'useMemo(() => ..., [deps])' ensures the expensive algorithm only recalculates when source items or filter criteria change, eliminating hundreds of milliseconds of processing delay.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Keystroke Render Stutter",
				Severity: "HIGH",
				Impact:   "Synchronous collection sorting on every keystroke freezes input acknowledgment and breaches 200ms INP.",
			},
			{
				Vector:   "Unnecessary Garbage Collection",
				Severity: "MEDIUM",
				Impact:   "Creating intermediate filtered and sorted array instances on every frame causes heavy GC pauses.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Unmemoized chained filter and sort running on every render",
				Code: `function UserList({ users, filterText }: Props) {
  const visibleUsers = users.filter(u => u.name.includes(filterText)).sort((a, b) => b.score - a.score);
  return <ul>{visibleUsers.map(u => <li key={u.id}>{u.name}</li>)}</ul>;
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Computation wrapped in useMemo to execute only when inputs change",
				Code: `function UserList({ users, filterText }: Props) {
  const visibleUsers = useMemo(() => {
    return users.filter(u => u.name.includes(filterText)).sort((a, b) => b.score - a.score);
  }, [users, filterText]);
  return <ul>{visibleUsers.map(u => <li key={u.id}>{u.name}</li>)}</ul>;
}`,
			},
		},
	}
}

// Evaluate memeriksa apakah terdapat transformasi data berat tanpa useMemo.
func (r *ExpensiveRenderComputationRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil {
		return nil
	}

	// 1. Periksa atribut elemen JSX (misal items={users.filter(...).sort(...)})
	if node.Type == ir.NodeElement {
		for attrName, attrVal := range node.Attributes {
			if attrName == "__script__" {
				continue
			}
			if pattern, hasExpensive := hasExpensiveRenderComputation(attrVal); hasExpensive {
				return []ir.Diagnostic{
					{
						Line:     node.Span.Line,
						Column:   node.Span.Column,
						Rule:     r.ID(),
						Severity: r.DefaultSeverity(),
						Message:  fmt.Sprintf("Expensive data transformation (%s) executes synchronously in prop '%s' on every render without 'useMemo'. This computation runs during keystrokes and interactions, causing main-thread stutter.", pattern, attrName),
						Hint:     "Wrap the expensive data transformation in 'useMemo(() => ..., [deps])' or precompute before render.",
					},
				}
			}
		}
	}

	// 2. Periksa skrip atau badan komponen di root / <script>
	code := getScriptOrSourceContent(node)
	if code != "" {
		if pattern, hasExpensive := hasExpensiveRenderComputation(code); hasExpensive {
			return []ir.Diagnostic{
				{
					Line:     node.Span.Line,
					Column:   node.Span.Column,
					Rule:     r.ID(),
					Severity: r.DefaultSeverity(),
					Message:  fmt.Sprintf("Expensive data transformation (%s) executes synchronously during render without 'useMemo'. This computation runs on every render cycle, spiking processing delay on user input.", pattern),
					Hint:     "Wrap the computation in 'useMemo(() => ..., [deps])' to memoize results across render cycles.",
				},
			}
		}
	}

	return nil
}
