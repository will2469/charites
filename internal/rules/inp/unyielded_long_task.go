package inp

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// UnyieldedLongTaskRule mendeteksi fungsi komputasi panjang atau iterasi data besar
// yang tidak menyertakan batas penjadwalan kooperatif (await scheduler.yield?.()).
type UnyieldedLongTaskRule struct{}

// NewUnyieldedLongTaskRule membuat instance baru dari UnyieldedLongTaskRule.
func NewUnyieldedLongTaskRule() *UnyieldedLongTaskRule {
	return &UnyieldedLongTaskRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *UnyieldedLongTaskRule) ID() string {
	return "inp.unyielded-long-task"
}

// Description mengembalikan ringkasan aturan.
func (r *UnyieldedLongTaskRule) Description() string {
	return "Long task processing large arrays without cooperative scheduling yields stalls main-thread responsiveness"
}

// Category mengembalikan nama kategori rule.
func (r *UnyieldedLongTaskRule) Category() string {
	return "inp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *UnyieldedLongTaskRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *UnyieldedLongTaskRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Cooperative Scheduling Controller (scheduler.yield)",
			"Google Chrome Core Web Vitals (Long Tasks & Input Responsiveness)",
			"Main-Thread Cooperative Concurrency Invariants",
		},
		CoreInvariant: "Long execution tasks triggered by or affecting user interactions must periodically yield control to the browser event loop via cooperative scheduling boundaries.",
		Grounding: "Long tasks running uninterrupted on the main thread (> 50ms) prevent the browser from acknowledging new user inputs (clicks, keypresses, taps) or rendering visual updates.\n\n" +
			"When user actions initiate extensive batch computations, running the entire process synchronously locks the page until completion, producing high Interaction to Next Paint (INP) latency.\n\n" +
			"By periodically pausing execution using modern cooperative scheduling: 'await (window.scheduler?.yield?.() ?? new Promise(r => setTimeout(r, 0)))', the browser is given immediate opportunities to handle pending user inputs and paint frames before continuing task work.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Main-Thread Input Starvation",
				Severity: "HIGH",
				Impact:   "Long uninterrupted execution loops starve the browser input queue, leaving pages unresponsive during batch processing.",
			},
			{
				Vector:   "High INP & Dropped Frames",
				Severity: "MEDIUM",
				Impact:   "Presentation of user feedback is blocked for hundreds of milliseconds, breaching the 200ms INP threshold.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "ts",
				Comment:  "Long calculation loop over large dataset without cooperative yield",
				Code: `function processLargeArray(items: string[]) {
  for (let i = 0; i < items.length; i++) {
    heavyCalculation(items[i]);
  }
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "ts",
				Comment:  "Periodic cooperative yielding to maintain input responsiveness",
				Code: `async function processLargeArray(items: string[]) {
  for (let i = 0; i < items.length; i++) {
    heavyCalculation(items[i]);
    if (i % 50 === 0) {
      await (window.scheduler?.yield?.() ?? new Promise(r => setTimeout(r, 0)));
    }
  }
}`,
			},
		},
	}
}

// Evaluate memeriksa apakah terdapat loop komputasi panjang yang kekurangan batas jeda kooperatif.
func (r *UnyieldedLongTaskRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil {
		return nil
	}

	// 1. Periksa elemen <script>
	if node.Type == ir.NodeElement && strings.EqualFold(node.Tag, "script") {
		scriptText := extractScriptNodeText(node)
		task, hasLongTask := hasUnyieldedLongLoop(scriptText)
		if hasLongTask {
			return []ir.Diagnostic{
				{
					Line:     node.Span.Line,
					Column:   node.Span.Column,
					Rule:     r.ID(),
					Severity: r.DefaultSeverity(),
					Message:  fmt.Sprintf("Detected %s processing large items without cooperative scheduling yields. Long uninterrupted loops starve the main thread and spike Interaction to Next Paint (INP).", task),
					Hint:     "Periodically yield execution back to the browser using 'await scheduler.yield?.()' or offload to a Web Worker.",
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
			task, hasLongTask := hasUnyieldedLongLoop(attrVal)
			if hasLongTask {
				return []ir.Diagnostic{
					{
						Line:     node.Span.Line,
						Column:   node.Span.Column,
						Rule:     r.ID(),
						Severity: r.DefaultSeverity(),
						Message:  fmt.Sprintf("Interactive handler '%s' contains %s without cooperative yields. This locks the main thread and causes severe interaction latency.", attrName, task),
						Hint:     "Break up long loops using 'await scheduler.yield?.()' or move heavy calculations to a Web Worker.",
					},
				}
			}
		}
	}

	return nil
}
