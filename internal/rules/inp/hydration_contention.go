package inp

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// HydrationContentionRule mendeteksi penumpukan direktif hidrasi eager client:load
// pada beberapa pulau Astro dalam satu template yang memicu perebutan thread utama.
type HydrationContentionRule struct{}

// NewHydrationContentionRule membuat instance baru dari HydrationContentionRule.
func NewHydrationContentionRule() *HydrationContentionRule {
	return &HydrationContentionRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *HydrationContentionRule) ID() string {
	return "inp.hydration-contention"
}

// Description mengembalikan ringkasan aturan.
func (r *HydrationContentionRule) Description() string {
	return "Concurrently hydrating multiple Astro client:load islands saturates the main thread and spikes input delay"
}

// Category mengembalikan nama kategori rule.
func (r *HydrationContentionRule) Category() string {
	return "inp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *HydrationContentionRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *HydrationContentionRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Astro Islands Architecture & Partial Hydration Specification",
			"W3C Cooperative Scheduling & Main-Thread Budget Invariants",
			"Google Core Web Vitals (INP Input Delay & Hydration Contention)",
		},
		CoreInvariant: "Astro templates must avoid declaring multiple eager 'client:load' island directives simultaneously; non-critical islands must use deferred hydration directives ('client:idle' or 'client:visible').",
		Grounding: "The 'client:load' directive instructs the browser to immediately fetch and execute island JavaScript upon page load, before user interaction or idle periods.\n\n" +
			"When multiple islands (3 or more) declare 'client:load' on the same page, their hydration phases execute in parallel or rapid succession on the main thread. This contention monopolizes CPU resources during initial user interactions, generating severe Long Tasks and inflating Input Delay.\n\n" +
			"By reserving 'client:load' strictly for critical interactive UI (such as primary navigation) and deferring secondary components to 'client:idle' or 'client:visible', the main thread remains responsive to user taps, clicks, and keystrokes.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Initial Hydration CPU Saturation",
				Severity: "HIGH",
				Impact:   "Multiple islands running concurrent React hydration lock the main thread during the window when users attempt first interaction.",
			},
			{
				Vector:   "Severe Input Delay Spikes",
				Severity: "MEDIUM",
				Impact:   "User clicks or keystrokes are queued behind synchronous island hydration tasks, resulting in INP > 200ms.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Multiple non-critical islands concurrently hydrated with client:load",
				Code: `---
import HeaderNav from '../components/HeaderNav.tsx';
import SearchBar from '../components/SearchBar.tsx';
import PromoBanner from '../components/PromoBanner.tsx';
---
<HeaderNav client:load />
<SearchBar client:load />
<PromoBanner client:load />`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Only critical navigation uses client:load; secondary islands use deferred hydration",
				Code: `---
import HeaderNav from '../components/HeaderNav.tsx';
import SearchBar from '../components/SearchBar.tsx';
import PromoBanner from '../components/PromoBanner.tsx';
---
<HeaderNav client:load />
<SearchBar client:idle />
<PromoBanner client:visible />`,
			},
		},
	}
}

// Evaluate memeriksa apakah dokumen template memuat lebih dari 2 pulau client:load secara bersamaan.
func (r *HydrationContentionRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Parent != nil {
		return nil
	}

	firstNode, count := countClientLoadIslands(node)
	if count > 2 && firstNode != nil {
		return []ir.Diagnostic{
			{
				Line:     firstNode.Span.Line,
				Column:   firstNode.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Found %d Astro islands declared with eager 'client:load' directive in the same template. Concurrently hydrating multiple islands saturates the main thread during initial page load, spiking Input Delay and Interaction to Next Paint (INP).", count),
				Hint:     "Reserve 'client:load' only for primary critical controls, and switch secondary islands to 'client:idle' or 'client:visible'.",
			},
		}
	}

	return nil
}
