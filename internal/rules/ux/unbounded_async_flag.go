package ux

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// UnboundedAsyncFlagRule mendeteksi handler asinkron yang mengaktifkan status loading (setLoading(true))
// sebelum pemanggilan await tanpa menjamin reset status (setLoading(false)) pada blok finally atau exit path error,
// guna mencegah antarmuka terkunci dalam status spinner tak berujung saat jaringan gagal.
type UnboundedAsyncFlagRule struct{}

// NewUnboundedAsyncFlagRule membuat instance baru dari UnboundedAsyncFlagRule.
func NewUnboundedAsyncFlagRule() *UnboundedAsyncFlagRule {
	return &UnboundedAsyncFlagRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *UnboundedAsyncFlagRule) ID() string {
	return "ux.unbounded-async-flag"
}

// Description mengembalikan ringkasan aturan.
func (r *UnboundedAsyncFlagRule) Description() string {
	return "Detects async handlers setting loading flags without guaranteed reset in finally/catch exit paths"
}

// Category mengembalikan nama kategori rule.
func (r *UnboundedAsyncFlagRule) Category() string {
	return "ux"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *UnboundedAsyncFlagRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *UnboundedAsyncFlagRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Mental Model Continuity & Deadlock Prevention",
			"Nielsen Heuristic #1: Visibility of System Status",
			"ISO 9241-110 Ergonomics of Human-System Interaction (Error Tolerance)",
		},
		CoreInvariant: "Async operations setting loading state before 'await' must guarantee state reset in all exit paths or a 'finally' block.",
		Grounding: "When asynchronous functions activate loading state (e.g. 'setLoading(true)') before awaiting a network operation " +
			"and fail to ensure that the flag is reset in a 'finally' block or error handler, any unexpected rejection (500 internal server error, " +
			"timeout, network dropout) leaves the UI permanently frozen in a loading spinner.\n\n" +
			"Users can neither re-try the action nor interact with adjacent controls, creating an unrecoverable dead end.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Permanent UI Lockup & Infinite Spinner Deadlock",
				Severity: "HIGH",
				Impact:   "Failed API requests leave spinners active indefinitely, blocking subsequent interactions and forcing users to hard-refresh.",
			},
			{
				Vector:   "Silent Failure Masking",
				Severity: "HIGH",
				Impact:   "Users believe an operation is still in flight even after the underlying request failed and aborted.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Loading state activated before await without reset in catch/finally block",
				Code: `<button
  onClick={async () => {
    setLoading(true);
    try {
      await api.fetchUsers();
    } catch (err) {
      console.error(err);
      // setLoading(false) terlupakan! Spinner berputar selamanya saat API gagal.
    }
  }}
>
  Muat Data
</button>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Guaranteed loading reset in finally block ensuring UI unlock under all outcomes",
				Code: `<button
  onClick={async () => {
    setLoading(true);
    try {
      await api.fetchUsers();
    } finally {
      setLoading(false);
    }
  }}
>
  Muat Data
</button>`,
			},
		},
	}
}

// Evaluate memeriksa apakah ada handler yang menyalakan loading flag tanpa reset terjamin.
func (r *UnboundedAsyncFlagRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	for attrName, attrVal := range node.Attributes {
		if !isEventHandlerOrActionAttr(attrName) {
			continue
		}

		if detectUnboundedAsyncLoading(attrVal) {
			return []ir.Diagnostic{
				{
					Line:     node.Span.Line,
					Column:   node.Span.Column,
					Rule:     r.ID(),
					Severity: r.DefaultSeverity(),
					Message: fmt.Sprintf(
						"Event handler %q sets async loading flag before 'await' without guaranteed reset in a 'finally' block or error exit path.",
						attrName,
					),
					Hint: "Wrap the async operation in a try...finally block and reset the loading state (e.g. 'finally { setLoading(false); }') to prevent infinite loading deadlocks.",
				},
			}
		}
	}

	return nil
}
