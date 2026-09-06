package cls

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// ClientOnlyHydrationPopRule mendeteksi penggunaan direktif client:only pada pulau Astro
// yang tidak menyertakan slot="fallback" atau wadah pembungkus dengan reservasi tinggi minimum (min-h-*).
type ClientOnlyHydrationPopRule struct{}

// NewClientOnlyHydrationPopRule membuat instance baru dari ClientOnlyHydrationPopRule.
func NewClientOnlyHydrationPopRule() *ClientOnlyHydrationPopRule {
	return &ClientOnlyHydrationPopRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *ClientOnlyHydrationPopRule) ID() string {
	return "cls.client-only-hydration-pop"
}

// Description mengembalikan ringkasan aturan.
func (r *ClientOnlyHydrationPopRule) Description() string {
	return "Astro client:only island lacks a slot='fallback' shell or reserved min-height container, causing hydration layout shift"
}

// Category mengembalikan nama kategori rule.
func (r *ClientOnlyHydrationPopRule) Category() string {
	return "cls"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *ClientOnlyHydrationPopRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ClientOnlyHydrationPopRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Astro Islands Architecture (client:only directives & fallback slots)",
			"W3C Core Web Vitals (Cumulative Layout Shift Prevention)",
			"Progressive Enhancement & Skeleton Shell Invariants",
		},
		CoreInvariant: "Astro components utilizing 'client:only' must define an official fallback shell (<div slot='fallback'>) or be enclosed within a container with reserved min-height.",
		Grounding: "In Astro's island architecture, the 'client:only' directive explicitly opts out of server-side rendering (SSR), omitting initial HTML markup for the component during build time.\n\n" +
			"Without a server-rendered placeholder or designated fallback shell, the browser initially renders an empty 0-height space. When the client-side JavaScript bundle finishes downloading, parsing, and executing, the rendered component abruptly expands and pushes all subsequent document content downward.\n\n" +
			"Providing a dedicated fallback shell via '<div slot=\"fallback\" class=\"min-h-[...]\">...</div>' ensures that the space is permanently reserved in initial server HTML, completely neutralizing Cumulative Layout Shift upon client hydration.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Post-Hydration Content Displacement",
				Severity: "HIGH",
				Impact:   "Delayed hydration of client-only islands causes sudden vertical document jumping when interactive components finish booting.",
			},
			{
				Vector:   "Blank Hole Flash",
				Severity: "MEDIUM",
				Impact:   "Users experience an empty white space where interactive widgets or charts belong prior to JavaScript execution.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "client:only island without fallback slot or reserved height",
				Code: `<main class="space-y-4">
  <h1>Dashboard</h1>
  <AnalyticsChart client:only="react" />
  <p>Live stats</p>
</main>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "client:only island with dedicated fallback slot shell",
				Code: `<main class="space-y-4">
  <h1>Dashboard</h1>
  <AnalyticsChart client:only="react">
    <div slot="fallback" class="w-full min-h-[350px] bg-muted/20 animate-pulse rounded-lg flex items-center justify-center">
      <span>Memuat grafik...</span>
    </div>
  </AnalyticsChart>
  <p>Live stats</p>
</main>`,
			},
		},
	}
}

// Evaluate memeriksa apakah komponen pulau client:only memiliki fallback shell atau penahan tinggi.
func (r *ClientOnlyHydrationPopRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !isClientOnlyIsland(node) {
		return nil
	}

	if hasAstroFallbackSlot(node) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Astro island component '%s' uses directive 'client:only' without an official fallback slot (<div slot='fallback'>) or reserved min-height container. Bypassing SSR leaves an empty hole on initial load, popping and shifting layout when client hydration finishes.", node.Tag),
			Hint:     "Provide a fallback shell using '<div slot=\"fallback\" class=\"min-h-[...]\">...</div>' or wrap the component in a container with a fixed or min-height class.",
		},
	}
}
