package performance

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// AstroOverPrefetchingRule mengaudit prefetch agresif pada tautan sekunder atau tautan footer.
type AstroOverPrefetchingRule struct{}

// NewAstroOverPrefetchingRule membuat instance baru dari AstroOverPrefetchingRule.
func NewAstroOverPrefetchingRule() *AstroOverPrefetchingRule {
	return &AstroOverPrefetchingRule{}
}

// ID mengembalikan identifier unik kanonikal aturan.
func (r *AstroOverPrefetchingRule) ID() string {
	return "performance.astro-over-prefetching"
}

// Category mengembalikan kategori aturan ('performance').
func (r *AstroOverPrefetchingRule) Category() string {
	return "performance"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *AstroOverPrefetchingRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Description mengembalikan deskripsi tujuan dan konteks aturan.
func (r *AstroOverPrefetchingRule) Description() string {
	return "Mencegah pemborosan kuota data seluler dengan melarang penempatan strategi prefetch agresif (viewport/load) pada tautan navigasi sekunder atau footer."
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *AstroOverPrefetchingRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Astro Prefetch Configuration Best Practices ('data-astro-prefetch')",
			"W3C Resource Hints & Speculative Parsing Bandwidth Economy",
			"Mobile Web Data Saver & Cellular Network Latency Guidelines",
		},
		CoreInvariant: "Aggressive 'viewport' or 'load' prefetch strategies must not be assigned to secondary or low-conversion navigation links; secondary links should use passive 'hover' or 'tap' prefetching.",
		Grounding: "Astro provides link prefetching via `data-astro-prefetch`.\n\n" +
			"Using aggressive strategies like `data-astro-prefetch=\"viewport\"` causes the browser to immediately fetch all linked documents as soon as their anchors enter the viewport.\n\n" +
			"When applied to secondary links (such as legal terms, privacy policies, or footer menus), this aggressively consumes user bandwidth and saturates the network connection, starving critical assets such as images and analytical payloads on slow mobile networks.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Cellular Data Waste",
				Severity: "MEDIUM",
				Impact:   "Preemptively downloads full pages that users rarely click, depleting metered mobile data connections.",
			},
			{
				Vector:   "Network Queue Contention",
				Severity: "MEDIUM",
				Impact:   "Prefetch network requests crowd the HTTP queue and delay high-priority above-the-fold assets.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Prefetch agresif pada tautan footer sekunder",
				Code: `<footer>
  <a href="/terms" data-astro-prefetch="viewport">Syarat & Ketentuan</a>
  <a href="/privacy" data-astro-prefetch="viewport">Kebijakan Privasi</a>
</footer>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Prefetch pasif saat hover untuk tautan footer",
				Code: `<footer>
  <a href="/terms" data-astro-prefetch="hover">Syarat & Ketentuan</a>
  <a href="/privacy" data-astro-prefetch="hover">Kebijakan Privasi</a>
</footer>`,
			},
		},
	}
}

// Evaluate memeriksa apakah tautan sekunder atau footer menggunakan prefetch agresif.
func (r *AstroOverPrefetchingRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || node.Tag != "a" {
		return nil
	}

	href, prefetch, isAggressive := isAggressiveSecondaryPrefetch(node)
	if !isAggressive {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Secondary or low-conversion link '%s' uses aggressive prefetch strategy '%s', wasting mobile user data and competing for network bandwidth with critical page resources.", href, prefetch),
			Hint:     "Change 'data-astro-prefetch' to 'hover' or 'tap' so secondary pages are only prefetched when the user intends to navigate.",
		},
	}
}
