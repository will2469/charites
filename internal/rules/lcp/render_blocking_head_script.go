package lcp

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// RenderBlockingHeadScriptRule mendeteksi elemen <script src="..."> di dalam <head>
// yang dieksekusi secara sinkron (misal dengan is:inline) tanpa atribut defer, async, atau type="module".
type RenderBlockingHeadScriptRule struct{}

// NewRenderBlockingHeadScriptRule membuat instance baru dari RenderBlockingHeadScriptRule.
func NewRenderBlockingHeadScriptRule() *RenderBlockingHeadScriptRule {
	return &RenderBlockingHeadScriptRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *RenderBlockingHeadScriptRule) ID() string {
	return "lcp.render-blocking-head-script"
}

// Description mengembalikan ringkasan aturan.
func (r *RenderBlockingHeadScriptRule) Description() string {
	return "External script in '<head>' without 'defer', 'async', or 'type=\"module\"' synchronously blocks HTML parser and delays LCP candidate paint"
}

// Category mengembalikan nama kategori rule.
func (r *RenderBlockingHeadScriptRule) Category() string {
	return "lcp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *RenderBlockingHeadScriptRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *RenderBlockingHeadScriptRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Google Chrome Core Web Vitals (Largest Contentful Paint Element Render Delay)",
			"HTML Living Standard (The script element & parser-blocking execution pipeline)",
			"W3C Web Performance Working Group Critical Path Minimization Guidelines",
		},
		CoreInvariant: "External scripts declared in the document '<head>' must specify 'defer', 'async', or 'type=\"module\"' to prevent halting HTML tokenization and per-frame rendering before the LCP candidate is displayed.",
		Grounding: "When the browser HTML parser encounters a synchronous `<script src=\"...\">` tag in the `<head>`, it must halt DOM construction, initiate a TCP/TLS connection to the script origin, download the JavaScript payload, and execute it before resuming document rendering.\n\n" +
			"In Astro, standard `<script>` tags are automatically processed by the bundler into deferred ES modules. However, external scripts tagged with `is:inline` or raw `<script src>` tags in document layouts bypass bundling and execute synchronously.\n\n" +
			"Adding 'defer' or 'type=\"module\"' allows the browser to download the script in parallel in the background while continuing HTML parsing and per-frame rendering of the hero LCP element.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Synchronous HTML Parser Halting",
				Severity: "HIGH",
				Impact:   "Halts DOM construction and suppresses layout passes, directly inflating LCP Element Render Delay by the full script network latency.",
			},
			{
				Vector:   "Head-of-Line Network Contention",
				Severity: "MEDIUM",
				Impact:   "Competes with critical hero media and external stylesheets for initial HTTP connection bandwidth.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Synchronous external inline script blocking HTML parsing in document head",
				Code: `<head>
  <script is:inline src="https://analytics.example.com/tracker.js"></script>
</head>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "External inline script deferred to allow non-blocking HTML parsing",
				Code: `<head>
  <script is:inline src="https://analytics.example.com/tracker.js" defer></script>
</head>`,
			},
		},
	}
}

// Evaluate memeriksa apakah tag script di head memblokir rendering sinkron.
func (r *RenderBlockingHeadScriptRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if !isRenderBlockingHeadScript(node) {
		return nil
	}

	srcVal := cleanAttrVal(node.Attributes["src"])
	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("External script '%s' in '<head>' lacks 'defer', 'async', or 'type=\"module\"'. Synchronously executing external scripts blocks the HTML parser and delays LCP candidate element rendering.", srcVal),
			Hint:     "Add 'defer' or 'type=\"module\"' to the '<script>' tag to allow parallel background script fetching without blocking initial page paint.",
		},
	}
}
