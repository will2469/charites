package performance

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// ReactIndexAsKeyRule mendeteksi penggunaan parameter indeks array sebagai 'key'
// di dalam perulangan .map() pada koleksi data dinamis.
type ReactIndexAsKeyRule struct{}

// NewReactIndexAsKeyRule membuat instance baru dari ReactIndexAsKeyRule.
func NewReactIndexAsKeyRule() *ReactIndexAsKeyRule {
	return &ReactIndexAsKeyRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *ReactIndexAsKeyRule) ID() string {
	return "performance.react-index-as-key"
}

// Description mengembalikan ringkasan aturan.
func (r *ReactIndexAsKeyRule) Description() string {
	return "Using array index as 'key' in dynamic collection mapping breaks VDOM reconciliation when items reorder or mutate"
}

// Category mengembalikan nama kategori rule.
func (r *ReactIndexAsKeyRule) Category() string {
	return "performance"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *ReactIndexAsKeyRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ReactIndexAsKeyRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"React Official Documentation (Lists and Keys Invariants)",
			"React Reconciliation Algorithm Specification (Diffing with stable keys)",
			"Robin Pokorny Guidelines ('Index as a key is an anti-pattern')",
		},
		CoreInvariant: "Dynamic collections mapped with '.map()' must use stable, unique item identifiers (e.g. 'item.id') as the 'key' attribute rather than numeric array indexes.",
		Grounding: "React relies on the `key` attribute to identify which items in a list have changed, been added, or been removed during reconciliation.\n\n" +
			"When an array index is used (`key={index}`), rearranging, filtering, prepending, or deleting items shifts the indexes of subsequent elements.\n\n" +
			"This index drift causes React to confuse element identities, erroneously preserving local uncontrolled state (such as form inputs, focus, and CSS transitions) on the wrong items and forcing redundant re-renders of the entire subtree.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Component State Desynchronization",
				Severity: "HIGH",
				Impact:   "Internal component state (e.g. input values, selection states) remains bound to the array position rather than the underlying data entity.",
			},
			{
				Vector:   "Unnecessary Subtree DOM Repaints",
				Severity: "MEDIUM",
				Impact:   "React fails to recognize moved nodes and completely remounts DOM subtrees instead of performing lightweight repositioning.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Dynamic transactions list using array index as key",
				Code: `{transactions.map((tx, index) => (
  <TransactionRow key={index} data={tx} />
))}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Persistent unique entity identifier used as reconciliation key",
				Code: `{transactions.map((tx) => (
  <TransactionRow key={tx.id} data={tx} />
))}`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen di dalam perulangan menggunakan indeks sebagai key.
func (r *ReactIndexAsKeyRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || node.Attributes == nil {
		return nil
	}

	fileSrc := getFileSourceContent(node)
	idxIdent, isViolation := isIndexKeyViolation(node, fileSrc)
	if !isViolation {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("JSX element uses array index '%s' as 'key' in dynamic collection mapping. Using indexes as keys breaks React VDOM reconciliation during item reordering, deletion, or sorting, triggering unnecessary DOM mutations and state leakage.", idxIdent),
			Hint:     "Replace array index with a stable, unique identifier from the data entity (e.g. 'key={item.id}').",
		},
	}
}
