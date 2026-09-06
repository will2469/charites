package lcp

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// ClientOnlyLCPContentRule mendeteksi komponen pulau (island) Astro pada area pelipatan atas (hero)
// yang dideklarasikan dengan direktif client:only tanpa menyediakan SSR slot="fallback".
type ClientOnlyLCPContentRule struct{}

// NewClientOnlyLCPContentRule membuat instance baru dari ClientOnlyLCPContentRule.
func NewClientOnlyLCPContentRule() *ClientOnlyLCPContentRule {
	return &ClientOnlyLCPContentRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *ClientOnlyLCPContentRule) ID() string {
	return "lcp.client-only-lcp-content"
}

// Description mengembalikan ringkasan aturan.
func (r *ClientOnlyLCPContentRule) Description() string {
	return "Above-the-fold hero island declared with 'client:only' without an SSR fallback slot, eliminating server HTML and delaying LCP until client-side bundle execution"
}

// Category mengembalikan nama kategori rule.
func (r *ClientOnlyLCPContentRule) Category() string {
	return "lcp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *ClientOnlyLCPContentRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ClientOnlyLCPContentRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Google Chrome Core Web Vitals (Largest Contentful Paint Resource Load Delay)",
			"Astro Island Architecture & Client Directives Specification",
			"W3C Web Performance Working Group SSR Invariants",
		},
		CoreInvariant: "Above-the-fold hero island components must not bypass server-side rendering with 'client:only' unless an SSR 'slot=\"fallback\"' is provided, ensuring initial HTML contains renderable LCP content.",
		Grounding: "In Astro's island architecture, declaring 'client:only' completely skips server-side rendering of the target component, emitting an empty placeholder container into the initial server HTML.\n\n" +
			"When an above-the-fold hero banner or primary heading is wrapped in 'client:only', the browser receives zero LCP content in the initial HTML stream. The browser cannot discover or render the hero element until all client-side JavaScript bundles are fetched, parsed, and executed.\n\n" +
			"To preserve fast LCP, developers should use 'client:load' (which renders initial HTML on the server and hydrates on the client) or provide a server-rendered fallback slot using '<div slot=\"fallback\">...</div>'.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Complete SSR Elimination on Critical Path",
				Severity: "CRITICAL",
				Impact:   "LCP candidate is entirely absent from initial server HTML response, turning server-rendered Astro pages into slow client-rendered SPAs.",
			},
			{
				Vector:   "Resource Load Delay Explosion",
				Severity: "HIGH",
				Impact:   "Hero rendering is blocked behind full JS download, parsing, and execution, inflating LCP by 600ms-2500ms on low-end devices.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Hero interactive island rendered with client:only without an SSR fallback slot",
				Code: `---
import HeroInteractive from '../components/HeroInteractive.tsx';
---
<main>
  <HeroInteractive client:only="react" />
</main>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Option 1: Using client:load to render initial HTML on the server",
				Code: `---
import HeroInteractive from '../components/HeroInteractive.tsx';
---
<main>
  <HeroInteractive client:load />
</main>`,
			},
			{
				Language: "astro",
				Comment:  "Option 2: Providing an SSR fallback slot when client:only is strictly necessary",
				Code: `---
import HeroInteractive from '../components/HeroInteractive.tsx';
---
<main>
  <HeroInteractive client:only="react">
    <div slot="fallback" class="hero-skeleton">
      <h1>Welcome to Our Platform</h1>
    </div>
  </HeroInteractive>
</main>`,
			},
		},
	}
}

// Evaluate memeriksa apakah pulau komponen hero menggunakan client:only tanpa fallback slot.
func (r *ClientOnlyLCPContentRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !hasClientOnlyDirective(node.Attributes) {
		return nil
	}

	if !isHeroIsland(node) {
		return nil
	}

	if hasFallbackSlot(node) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Hero island component '<%s>' uses 'client:only' without an SSR fallback slot. This strips hero content from initial server HTML, delaying LCP until client-side hydration completes.", node.Tag),
			Hint:     "Change 'client:only' to 'client:load' to preserve server-side rendering, or provide an SSR fallback slot with '<div slot=\"fallback\">...</div>'.",
		},
	}
}
