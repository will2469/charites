package inp

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// RepeatedStateUpdateRule mendeteksi pemanggilan state updater React (setState / setFoo)
// di dalam perulangan yang memecah batas batching otomatis React 18+ (loop asinkron dengan await atau flushSync).
type RepeatedStateUpdateRule struct{}

// NewRepeatedStateUpdateRule membuat instance baru dari RepeatedStateUpdateRule.
func NewRepeatedStateUpdateRule() *RepeatedStateUpdateRule {
	return &RepeatedStateUpdateRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *RepeatedStateUpdateRule) ID() string {
	return "inp.repeated-state-update"
}

// Description mengembalikan ringkasan aturan.
func (r *RepeatedStateUpdateRule) Description() string {
	return "Repeated state updater calls inside loops breaking automatic batching trigger cascading re-renders"
}

// Category mengembalikan nama kategori rule.
func (r *RepeatedStateUpdateRule) Category() string {
	return "inp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *RepeatedStateUpdateRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *RepeatedStateUpdateRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"React 18+ Automatic Batching Specification",
			"W3C Web Performance Working Group (Interaction to Next Paint - INP)",
			"Concurrent React Scheduling & Reconciliation Cost",
		},
		CoreInvariant: "React state setters must not be repeatedly invoked within loop iterations that break automatic batching (such as asynchronous loops containing 'await' or 'flushSync').",
		Grounding: "While React 18 automatically batches multiple state updates within standard synchronous handlers, asynchronous loops (e.g. 'for ... of' with 'await' inside) or explicit 'flushSync' blocks break automatic batching.\n\n" +
			"Calling a state updater on every iteration of an asynchronous loop causes React to trigger a full re-render, VDOM diffing, and reconciliation cycle on every microtask tick.\n\n" +
			"This creates an enormous render queue backlog on the main thread, stalling user interactions and causing high Interaction to Next Paint (INP) latency.\n\n" +
			"Accumulating results locally into an array and issuing a single state update after the loop completes ensures a single, batched render pass.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Per-Iteration Re-render Cascades",
				Severity: "HIGH",
				Impact:   "Each iteration of an async loop schedules a separate render pass, saturating the React scheduler and freezing UI input.",
			},
			{
				Vector:   "Presentation Delay Ballooning",
				Severity: "MEDIUM",
				Impact:   "Successive re-renders continuously postpone the browser paint phase, severely degrading INP.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "State updater called on each iteration of an async loop",
				Code: `for (const item of items) {
  const detail = await fetchDetail(item.id);
  setItems(prev => [...prev, detail]); // Memicu re-render pada setiap iterasi!
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Accumulating all results and updating state once after loop completion",
				Code: `const results = [];
for (const item of items) {
  results.push(await fetchDetail(item.id));
}
setItems(prev => [...prev, ...results]); // Hanya satu siklus render`,
			},
		},
	}
}

// Evaluate memeriksa apakah terdapat pemanggilan state updater di dalam loop pemecah batching.
func (r *RepeatedStateUpdateRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil {
		return nil
	}

	// 1. Periksa elemen <script>
	if node.Type == ir.NodeElement && strings.EqualFold(node.Tag, "script") {
		scriptText := extractScriptNodeText(node)
		setter, hasRepeated := hasRepeatedStateUpdateInLoop(scriptText)
		if hasRepeated {
			return []ir.Diagnostic{
				{
					Line:     node.Span.Line,
					Column:   node.Span.Column,
					Rule:     r.ID(),
					Severity: r.DefaultSeverity(),
					Message:  fmt.Sprintf("State updater '%s' is invoked inside a loop breaking React batching (await or flushSync). This schedules cascading re-renders on every iteration, severely degrading Interaction to Next Paint (INP).", setter),
					Hint:     "Accumulate values into a local array or object within the loop, then call the state updater once after iteration completes.",
				},
			}
		}
	}

	// 2. Periksa atribut handler event
	if node.Type == ir.NodeElement {
		for attrName, attrVal := range node.Attributes {
			if !isInteractiveHandlerAttr(attrName) {
				continue
			}
			setter, hasRepeated := hasRepeatedStateUpdateInLoop(attrVal)
			if hasRepeated {
				return []ir.Diagnostic{
					{
						Line:     node.Span.Line,
						Column:   node.Span.Column,
						Rule:     r.ID(),
						Severity: r.DefaultSeverity(),
						Message:  fmt.Sprintf("State updater '%s' is called inside a loop breaking React batching in handler '%s'. This triggers cascading re-renders on every iteration, causing interaction lag.", setter, attrName),
						Hint:     "Collect all values first and issue a single state update after the loop.",
					},
				}
			}
		}
	}

	return nil
}
