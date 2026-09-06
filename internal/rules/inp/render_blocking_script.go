package inp

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// RenderBlockingScriptRule mendeteksi tag script eksternal sinkron yang memblokir
// parser HTML dan menunda kesiapan thread utama dalam menerima interaksi pengguna.
type RenderBlockingScriptRule struct{}

// NewRenderBlockingScriptRule membuat instance baru dari RenderBlockingScriptRule.
func NewRenderBlockingScriptRule() *RenderBlockingScriptRule {
	return &RenderBlockingScriptRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *RenderBlockingScriptRule) ID() string {
	return "inp.render-blocking-script"
}

// Description mengembalikan ringkasan aturan.
func (r *RenderBlockingScriptRule) Description() string {
	return "External script element without defer, async, or type=\"module\" synchronously blocks rendering and input responsiveness"
}

// Category mengembalikan nama kategori rule.
func (r *RenderBlockingScriptRule) Category() string {
	return "inp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *RenderBlockingScriptRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *RenderBlockingScriptRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"HTML Living Standard (The script element & execution pipeline)",
			"W3C Web Performance & Navigation Timing Specification",
			"Google Chrome Core Web Vitals (Eliminating Render-Blocking Resources)",
		},
		CoreInvariant: "External script elements must declare 'defer', 'async', or 'type=\"module\"' to avoid synchronously blocking HTML parsing and main-thread readiness.",
		Grounding: "When the browser encounters a synchronous `<script src=\"...\">` tag, it must pause HTML parsing, establish a network connection, download the script, and execute it before resuming document rendering.\n\n" +
			"In Astro, standard `<script>` tags are automatically bundled into deferred ES modules. However, scripts marked with `is:inline` or raw external scripts in HTML document heads bypass bundling and execute synchronously.\n\n" +
			"Adding 'defer' or 'type=\"module\"' ensures the script is downloaded in the background and executed without halting the parser, keeping the browser immediately receptive to early user taps and clicks.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Synchronous Parser Halting",
				Severity: "HIGH",
				Impact:   "HTML parsing and initial rendering are paused until external scripts download and execute.",
			},
			{
				Vector:   "Delayed Main-Thread Input Availability",
				Severity: "MEDIUM",
				Impact:   "The browser input event loop is delayed, resulting in unacknowledged early user taps.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Synchronous external inline script blocking HTML parser",
				Code:     `<script is:inline src="https://analytics.example.com/heavy-bundle.js"></script>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "External inline script deferred to prevent parser blocking",
				Code:     `<script is:inline src="https://analytics.example.com/heavy-bundle.js" defer></script>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen script eksternal memblokir rendering sinkron.
func (r *RenderBlockingScriptRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if !isRenderBlockingScript(node) {
		return nil
	}

	srcVal := node.Attributes["src"]
	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("External script '%s' lacks 'defer', 'async', or 'type=\"module\"'. Synchronously evaluating external scripts pauses the HTML parser and delays main-thread availability for user interactions.", srcVal),
			Hint:     "Add 'defer' or 'type=\"module\"' to the script tag to ensure asynchronous non-blocking loading.",
		},
	}
}
