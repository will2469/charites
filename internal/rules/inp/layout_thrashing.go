package inp

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// LayoutThrashingRule mendeteksi penulisan gaya/kelas DOM yang langsung diikuti oleh pembacaan
// properti geometri tata letak (offsetWidth, getBoundingClientRect) di dalam alur eksekusi sinkron yang sama.
type LayoutThrashingRule struct{}

// NewLayoutThrashingRule membuat instance baru dari LayoutThrashingRule.
func NewLayoutThrashingRule() *LayoutThrashingRule {
	return &LayoutThrashingRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *LayoutThrashingRule) ID() string {
	return "inp.layout-thrashing"
}

// Description mengembalikan ringkasan aturan.
func (r *LayoutThrashingRule) Description() string {
	return "Sequential DOM style mutation followed by layout geometry reading triggers forced synchronous reflow"
}

// Category mengembalikan nama kategori rule.
func (r *LayoutThrashingRule) Category() string {
	return "inp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *LayoutThrashingRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *LayoutThrashingRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Web Performance Working Group (Interaction to Next Paint - INP)",
			"Google Chrome Rendering Engine Pipeline (Forced Synchronous Layout)",
			"Browser Main-Thread Event Loop Scheduling",
		},
		CoreInvariant: "Imperative JavaScript execution must separate layout queries from style mutations, avoiding forced synchronous reflow passes (read-then-write batching).",
		Grounding: "When JavaScript mutates DOM styles or class names (e.g. 'el.style.width = ...') and subsequently reads a layout geometry property (e.g. 'el.offsetHeight' or 'getBoundingClientRect()') within the same synchronous execution block, the browser is forced to flush pending style changes and perform an immediate, blocking layout recalculation.\n\n" +
			"This phenomenon, known as 'Layout Thrashing' or 'Forced Synchronous Reflow', locks the browser main thread, preventing user interaction processing and drastically inflating Interaction to Next Paint (INP) latency.\n\n" +
			"Batching all layout reads before performing style writes, or deferring updates via 'requestAnimationFrame', prevents synchronous recalculations and keeps the main thread responsive.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Forced Synchronous Reflow Stalling",
				Severity: "HIGH",
				Impact:   "Synchronous layout computation blocks the main thread during interaction handling, causing dropped frames and severe INP degradation.",
			},
			{
				Vector:   "Cascading Rendering Bottleneck",
				Severity: "HIGH",
				Impact:   "Interleaved write-read loops exponentially degrade interaction responsiveness on complex DOM trees.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Style mutation immediately followed by geometry reading (forced reflow)",
				Code: `function resizeBox(el: HTMLElement) {
  el.style.width = '200px';
  const height = el.offsetHeight; // Memaksa kalkulasi layout sinkron!
  el.style.height = (height * 2) + 'px';
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Read-then-write batching to prevent forced layout calculation",
				Code: `function resizeBox(el: HTMLElement) {
  const currentHeight = el.offsetHeight; // Baca di awal
  el.style.width = '200px';              // Tulis serentak
  el.style.height = (currentHeight * 2) + 'px';
}`,
			},
		},
	}
}

// Evaluate memeriksa apakah terdapat sekuens layout thrashing pada simpul skrip atau atribut fungsi.
func (r *LayoutThrashingRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil {
		return nil
	}

	// 1. Periksa elemen <script>
	if node.Type == ir.NodeElement && strings.EqualFold(node.Tag, "script") {
		scriptText := extractScriptNodeText(node)
		writePattern, readPattern, hasThrashing := hasLayoutThrashingSequence(scriptText)
		if hasThrashing {
			return []ir.Diagnostic{
				{
					Line:     node.Span.Line,
					Column:   node.Span.Column,
					Rule:     r.ID(),
					Severity: r.DefaultSeverity(),
					Message:  fmt.Sprintf("Detected layout thrashing sequence: DOM style mutation '%s' followed by geometry read '%s' in synchronous flow. This forces synchronous layout recalculation on the main thread, severely degrading Interaction to Next Paint (INP).", writePattern, readPattern),
					Hint:     "Batch all layout readings (e.g. offsetHeight) before performing style mutations, or isolate mutations using requestAnimationFrame.",
				},
			}
		}
	}

	// 2. Periksa atribut handler event (e.g. onClick, onKeyDown)
	if node.Type == ir.NodeElement {
		for _, attrVal := range node.Attributes {
			writePattern, readPattern, hasThrashing := hasLayoutThrashingSequence(attrVal)
			if hasThrashing {
				return []ir.Diagnostic{
					{
						Line:     node.Span.Line,
						Column:   node.Span.Column,
						Rule:     r.ID(),
						Severity: r.DefaultSeverity(),
						Message:  fmt.Sprintf("Detected layout thrashing sequence in event handler: DOM style mutation '%s' followed by geometry read '%s'. This triggers forced synchronous reflow, spiking INP latency.", writePattern, readPattern),
						Hint:     "Separate layout readings from style mutations or batch DOM writes via requestAnimationFrame.",
					},
				}
			}
		}
	}

	return nil
}
