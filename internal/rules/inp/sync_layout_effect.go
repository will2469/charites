package inp

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// SyncLayoutEffectRule mendeteksi eksekusi komputasi non-geometris berat di dalam useLayoutEffect
// yang menahan fase paint peramban secara sinkron dan memperpanjang presentation delay.
type SyncLayoutEffectRule struct{}

// NewSyncLayoutEffectRule membuat instance baru dari SyncLayoutEffectRule.
func NewSyncLayoutEffectRule() *SyncLayoutEffectRule {
	return &SyncLayoutEffectRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *SyncLayoutEffectRule) ID() string {
	return "inp.sync-layout-effect"
}

// Description mengembalikan ringkasan aturan.
func (r *SyncLayoutEffectRule) Description() string {
	return "Synchronous non-geometrical computation in useLayoutEffect blocks browser paint and inflates presentation delay"
}

// Category mengembalikan nama kategori rule.
func (r *SyncLayoutEffectRule) Category() string {
	return "inp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *SyncLayoutEffectRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *SyncLayoutEffectRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"React useLayoutEffect Pre-Paint Execution Model",
			"W3C Presentation Timing & Frame Pipeline Invariants",
			"Google Chrome Core Web Vitals (INP Presentation Delay Optimization)",
		},
		CoreInvariant: "The 'useLayoutEffect' hook must be reserved strictly for synchronous DOM measurements; data fetching and non-geometrical state updates must reside in 'useEffect'.",
		Grounding: "Unlike 'useEffect' which runs asynchronously after the browser paints the screen, 'useLayoutEffect' fires synchronously immediately after React commits DOM mutations, *before* the browser renders pixels to the screen.\n\n" +
			"Executing non-geometrical operations (such as data fetching, localStorage I/O, or secondary state cascades) within 'useLayoutEffect' delays the browser paint phase directly, locking the main thread and dramatically increasing Presentation Delay.\n\n" +
			"Developers should restrict 'useLayoutEffect' exclusively to reading layout properties (e.g. 'getBoundingClientRect') to position popovers or tooltips without flicker, moving all other logic to 'useEffect'.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Browser Paint Phase Halting",
				Severity: "HIGH",
				Impact:   "Frame rendering is synchronously blocked while non-layout logic executes in useLayoutEffect.",
			},
			{
				Vector:   "Presentation Delay Spikes",
				Severity: "HIGH",
				Impact:   "Visual acknowledgment of user interactions is delayed by hundreds of milliseconds, breaching the 200ms INP threshold.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Data fetching inside useLayoutEffect blocks the browser paint phase",
				Code: `useLayoutEffect(() => {
  fetchUserData(userId).then(setData);
}, [userId]);`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Data fetching moved to useEffect; browser paints pixels without delay",
				Code: `useEffect(() => {
  fetchUserData(userId).then(setData);
}, [userId]);`,
			},
		},
	}
}

// Evaluate memeriksa apakah terdapat komputasi non-geometris di dalam useLayoutEffect.
func (r *SyncLayoutEffectRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil {
		return nil
	}

	code := getScriptOrSourceContent(node)
	if code == "" {
		return nil
	}

	op, line, isViolation := findSyncLayoutEffectViolation(code)
	if isViolation {
		targetLine := node.Span.Line
		if line > 0 {
			targetLine = line
		}
		return []ir.Diagnostic{
			{
				Line:     targetLine,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Hook 'useLayoutEffect' executes synchronous non-geometrical work (%s) before browser paint. Holding up the paint phase blocks visual presentation and inflates Interaction to Next Paint (INP).", op),
				Hint:     "Move asynchronous tasks, network requests, or non-layout state updates to 'useEffect', reserving 'useLayoutEffect' strictly for immediate DOM geometry measurements.",
			},
		}
	}

	return nil
}
