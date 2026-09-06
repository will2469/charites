package inp

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// HeavyEventHandlerRule mendeteksi penangan interaksi pengguna (onClick, onKeyDown, dsb.)
// yang memuat operasi komputasi sinkron berat (JSON.parse, .sort, loop bertingkat) tanpa pembagian tugas kooperatif.
type HeavyEventHandlerRule struct{}

// NewHeavyEventHandlerRule membuat instance baru dari HeavyEventHandlerRule.
func NewHeavyEventHandlerRule() *HeavyEventHandlerRule {
	return &HeavyEventHandlerRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *HeavyEventHandlerRule) ID() string {
	return "inp.heavy-event-handler"
}

// Description mengembalikan ringkasan aturan.
func (r *HeavyEventHandlerRule) Description() string {
	return "Interactive event handler executes heavy synchronous operations (JSON.parse, Array.sort) without cooperative yields"
}

// Category mengembalikan nama kategori rule.
func (r *HeavyEventHandlerRule) Category() string {
	return "inp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *HeavyEventHandlerRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *HeavyEventHandlerRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Web Performance Working Group (Interaction to Next Paint - INP)",
			"Google Chrome Core Web Vitals (Input Delay & Processing Duration)",
			"Browser Cooperative Scheduling Guidelines (scheduler.yield)",
		},
		CoreInvariant: "Interactive event handlers must avoid heavy synchronous computations on the main thread, adopting cooperative task yielding or Web Worker offloading.",
		Grounding: "When users tap, click, or type, the browser expects the main thread to quickly acknowledge the interaction and schedule the next paint frame (ideally within 50ms, with INP target <= 200ms).\n\n" +
			"Executing heavy synchronous operations (such as large JSON parsing, array sorting, or complex synchronous data manipulation) directly inside event handler callbacks blocks the main thread during the crucial input processing phase.\n\n" +
			"This delays the presentation of visual feedback (e.g. active button states, loading spinners) and directly inflates the INP metric.\n\n" +
			"Breaking long tasks with 'await scheduler.yield?.()' or offloading computation to a dedicated Web Worker allows the browser to present visual feedback immediately before executing intensive processing.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Processing Phase Thread Saturation",
				Severity: "HIGH",
				Impact:   "Synchronous algorithms in click/key handlers block the main thread, exceeding the 200ms INP threshold.",
			},
			{
				Vector:   "Frozen Visual Feedback",
				Severity: "MEDIUM",
				Impact:   "Buttons appear unresponsive or stuck because UI rendering is starved by long synchronous handler execution.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Synchronous heavy data parsing and sorting directly inside onClick handler",
				Code: `<button onClick={() => {
  const data = JSON.parse(hugePayload);
  const sorted = data.sort((a, b) => b.score - a.score);
  setResults(sorted);
}}>
  Urutkan Data
</button>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Cooperative yielding to acknowledge user input before intensive processing",
				Code: `<button onClick={async () => {
  setLoading(true);
  await (window.scheduler?.yield?.() ?? new Promise(r => setTimeout(r, 0)));
  const data = JSON.parse(hugePayload);
  setResults(data.sort((a, b) => b.score - a.score));
  setLoading(false);
}}>
  Urutkan Data
</button>`,
			},
		},
	}
}

// Evaluate memeriksa apakah atribut penangan interaksi memuat komputasi sinkron berat.
func (r *HeavyEventHandlerRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	for attrName, attrVal := range node.Attributes {
		if !isInteractiveHandlerAttr(attrName) {
			continue
		}

		offending, isHeavy := hasHeavySynchronousOps(attrVal)
		if isHeavy {
			return []ir.Diagnostic{
				{
					Line:     node.Span.Line,
					Column:   node.Span.Column,
					Rule:     r.ID(),
					Severity: r.DefaultSeverity(),
					Message:  fmt.Sprintf("Interactive event handler '%s' executes heavy synchronous operation '%s'. Running intensive synchronous work in interaction callbacks monopolizes the main thread, spiking Interaction to Next Paint (INP).", attrName, offending),
					Hint:     "Yield to the browser event loop using 'await scheduler.yield?.()' before heavy processing, or offload the algorithm to a Web Worker.",
				},
			}
		}
	}

	return nil
}
