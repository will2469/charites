package inp

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// MissingStartTransitionRule mendeteksi penangan interaksi yang menggabungkan pembaruan input
// mendesak dengan pembaruan state sekunder berat tanpa pembungkus startTransition.
type MissingStartTransitionRule struct{}

// NewMissingStartTransitionRule membuat instance baru dari MissingStartTransitionRule.
func NewMissingStartTransitionRule() *MissingStartTransitionRule {
	return &MissingStartTransitionRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *MissingStartTransitionRule) ID() string {
	return "inp.missing-start-transition"
}

// Description mengembalikan ringkasan aturan.
func (r *MissingStartTransitionRule) Description() string {
	return "Secondary non-urgent state update inside interactive handler should be wrapped in startTransition to prevent input lag"
}

// Category mengembalikan nama kategori rule.
func (r *MissingStartTransitionRule) Category() string {
	return "inp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (info).
func (r *MissingStartTransitionRule) DefaultSeverity() ir.Severity {
	return ir.SeverityInfo
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *MissingStartTransitionRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"React 18/19 Concurrent Mode Architecture (startTransition & Transitions API)",
			"W3C User Timing & Cooperative Scheduling Invariants",
			"Google Chrome Core Web Vitals (Input to Paint Responsiveness)",
		},
		CoreInvariant: "Secondary non-urgent state updates triggered alongside urgent user input must be scheduled as transitions via 'React.startTransition' to preserve typing responsiveness.",
		Grounding: "In modern user interfaces, an interactive event (such as typing in a search bar or clicking a filter tab) often triggers two types of updates: an urgent update (updating the input text cursor) and a non-urgent secondary update (filtering a large list or fetching preview cards).\n\n" +
			"When both updates are processed synchronously without transitions, React treats the expensive secondary re-render with the same high priority as the keystroke, blocking the main thread and causing noticeable keystroke stutter.\n\n" +
			"Wrapping secondary updates in 'React.startTransition' informs the scheduler that the secondary render is interruptible. React will immediately paint the user's keystroke, keeping INP low while deferring list rendering.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Keystroke Input Lag",
				Severity: "MEDIUM",
				Impact:   "Synchronous secondary re-renders block subsequent keystroke frames, creating sluggish typing feedback.",
			},
			{
				Vector:   "Main Thread Presentation Delays",
				Severity: "LOW",
				Impact:   "The browser cannot acknowledge user interactions within the 200ms INP threshold.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Synchronously combining urgent input setter with heavy list filtering",
				Code: `function handleSearch(e: React.ChangeEvent<HTMLInputElement>) {
  setSearchQuery(e.target.value);
  setFilteredLargeList(expensiveFilter(e.target.value));
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Urgent input text is updated immediately; secondary list is wrapped in startTransition",
				Code: `function handleSearch(e: React.ChangeEvent<HTMLInputElement>) {
  setSearchQuery(e.target.value);
  React.startTransition(() => {
    setFilteredLargeList(expensiveFilter(e.target.value));
  });
}`,
			},
		},
	}
}

// Evaluate memeriksa apakah handler interaksi memicu pembaruan sekunder berat tanpa startTransition.
func (r *MissingStartTransitionRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	for attrName, attrVal := range node.Attributes {
		if !isInteractiveHandlerAttr(attrName) {
			continue
		}

		heavySetter, missing := hasMissingStartTransition(attrVal)
		if missing {
			return []ir.Diagnostic{
				{
					Line:     node.Span.Line,
					Column:   node.Span.Column,
					Rule:     r.ID(),
					Severity: r.DefaultSeverity(),
					Message:  fmt.Sprintf("Interactive handler '%s' combines an urgent input state update with a secondary non-urgent update '%s' without wrapping it in 'startTransition'. This stalls typing responsiveness and inflates Interaction to Next Paint (INP).", attrName, heavySetter),
					Hint:     "Wrap the secondary state update in 'React.startTransition(() => { ... })' so React prioritizes keystrokes over expensive re-renders.",
				},
			}
		}
	}

	return nil
}
