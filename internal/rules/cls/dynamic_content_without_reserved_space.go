package cls

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// DynamicContentWithoutReservedSpaceRule mendeteksi penyuntikan konten dinamis atau komponen banner/widget
// dalam alur dokumen tanpa adanya reservasi dimensi geometris (min-h-* / h-*) atau perlindungan Suspense.
type DynamicContentWithoutReservedSpaceRule struct{}

// NewDynamicContentWithoutReservedSpaceRule membuat instance baru dari DynamicContentWithoutReservedSpaceRule.
func NewDynamicContentWithoutReservedSpaceRule() *DynamicContentWithoutReservedSpaceRule {
	return &DynamicContentWithoutReservedSpaceRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *DynamicContentWithoutReservedSpaceRule) ID() string {
	return "cls.dynamic-content-without-reserved-space"
}

// Description mengembalikan ringkasan aturan.
func (r *DynamicContentWithoutReservedSpaceRule) Description() string {
	return "Dynamic widget or banner injected in document flow lacks reserved vertical dimensions (min-h/h), risking content reflow"
}

// Category mengembalikan nama kategori rule.
func (r *DynamicContentWithoutReservedSpaceRule) Category() string {
	return "cls"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *DynamicContentWithoutReservedSpaceRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *DynamicContentWithoutReservedSpaceRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Google Core Web Vitals (Dynamic Content Injection Guidelines)",
			"React Suspense & Skeleton Architecture Invariants",
			"W3C Cumulative Layout Shift Mitigation",
		},
		CoreInvariant: "Dynamic in-flow widgets, promotional banners, or asynchronously injected components must be enclosed in containers with reserved vertical dimensions ('min-h-*') or guarded by Suspense fallback skeletons.",
		Grounding: "When asynchronous content (such as personalization widgets, promotional announcements, or dynamic notification bars) loads after initial page paint, injecting it directly into normal document flow without pre-allocated vertical space forces all content below it to shift abruptly downward.\n\n" +
			"This post-load displacement is one of the single largest real-world contributors to poor Cumulative Layout Shift (CLS) scores, frequently leading to accidental user miss-clicks and navigation disorientation.\n\n" +
			"Enclosing dynamic elements in a container with an explicit minimum height ('min-h-[120px]') or using a matching skeleton placeholder ensures document flow remains completely stable before and after data resolution.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Post-Load In-Flow Content Snapping",
				Severity: "HIGH",
				Impact:   "Asynchronous widgets popping into document flow push down articles, forms, or buttons while users are reading or interacting.",
			},
			{
				Vector:   "Accidental Miss-Clicks",
				Severity: "HIGH",
				Impact:   "Sudden vertical shifts cause users to accidentally click unintended links or submit buttons.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Unreserved dynamic promotional banner in document flow",
				Code: `<main>
  <h1>Artikel</h1>
  <PromoBanner />
  <Content />
</main>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Dynamic banner enclosed in container with reserved min-height",
				Code: `<main>
  <h1>Artikel</h1>
  <div className="min-h-[120px]">
    <PromoBanner />
  </div>
  <Content />
</main>`,
			},
		},
	}
}

// Evaluate memeriksa apakah komponen dinamis in-flow kekurangan reservasi dimensi vertikal.
func (r *DynamicContentWithoutReservedSpaceRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	offending, isUnreserved := isUnreservedDynamicContent(node)
	if !isUnreserved {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Dynamic content or widget '%s' injected into document flow lacks reserved vertical dimensions ('min-h-*' or 'h-*'). Asynchronous in-flow injections push down adjacent content and trigger severe Cumulative Layout Shift (CLS).", offending),
			Hint:     "Enclose the dynamic component in a container with a reserved min-height (e.g. <div class='min-h-[120px]'>) or wrap it with a <Suspense fallback={<Skeleton />}> boundary.",
		},
	}
}
