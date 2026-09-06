package inp

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// UnboundedCollectionRenderRule mendeteksi elemen kontainer scroll yang merender pemetaan
// koleksi dinamis berkardinalitas tak terhingga (.map()) langsung ke DOM tanpa paginasi atau virtualisasi.
type UnboundedCollectionRenderRule struct{}

// NewUnboundedCollectionRenderRule membuat instance baru dari UnboundedCollectionRenderRule.
func NewUnboundedCollectionRenderRule() *UnboundedCollectionRenderRule {
	return &UnboundedCollectionRenderRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *UnboundedCollectionRenderRule) ID() string {
	return "inp.unbounded-collection-render"
}

// Description mengembalikan ringkasan aturan.
func (r *UnboundedCollectionRenderRule) Description() string {
	return "Scrollable collection container renders unbounded dynamic data via .map() without window virtualization or pagination limits"
}

// Category mengembalikan nama kategori rule.
func (r *UnboundedCollectionRenderRule) Category() string {
	return "inp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *UnboundedCollectionRenderRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *UnboundedCollectionRenderRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Google Chrome Core Web Vitals (Interaction to Next Paint Presentation Delay)",
			"W3C DOM Performance & Rendering Subsystem Scaling",
			"React Virtual List Windowing Patterns (@tanstack/react-virtual)",
		},
		CoreInvariant: "Scrollable collection containers must not render arbitrarily large dynamic collections directly into the DOM; virtualization windowing or explicit pagination limits must be applied to cap active DOM node count.",
		Grounding: "When dynamic lists or tables render an unbounded number of items directly via '.map()', every item creates multiple nested DOM elements.\n\n" +
			"In scrollable containers, users scroll while interacting. If hundreds or thousands of DOM nodes reside in the tree, every user interaction triggers recalculations across the massive DOM tree, inflating browser Presentation Delay well beyond the 200ms INP threshold.\n\n" +
			"Window virtualization (e.g. '@tanstack/react-virtual') or explicit pagination (e.g. '.slice(0, 20)') limits rendered elements strictly to the visible viewport, keeping the DOM lightweight and presentation latency minimal.",
		Risks: []ir.RiskItem{
			{
				Vector:   "DOM Node Count Explosion",
				Severity: "HIGH",
				Impact:   "Massive collections mapped directly to DOM cause layout tree bloat, degrading styling calculations and memory usage.",
			},
			{
				Vector:   "Excessive Presentation Delay",
				Severity: "HIGH",
				Impact:   "Browser rendering engine spends excessive frame time recalibrating off-screen nodes during user interactions.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Scrollable container rendering full dynamic collection without virtualization or limits",
				Code: `<div className="h-96 overflow-y-auto">
  {dynamicDataFromApi.map(item => (
    <InteractiveItemRow key={item.id} data={item} onSelect={handleSelect} />
  ))}
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Virtual list windowing rendering only visible items in viewport",
				Code: `<div ref={parentRef} className="h-96 overflow-y-auto">
  <div style={{ height: ` + "`${rowVirtualizer.getTotalSize()}px`" + ` }}>
    {rowVirtualizer.getVirtualItems().map(virtualRow => (
      <InteractiveItemRow key={virtualRow.index} data={dynamicDataFromApi[virtualRow.index]} />
    ))}
  </div>
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen kontainer scroll merender koleksi dinamis tanpa virtualisasi/paginasi.
func (r *UnboundedCollectionRenderRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !isScrollableContainer(node.Tag, node.Classes, node.RawClasses) {
		return nil
	}

	src := getScriptOrSourceContent(node)
	if src == "" {
		return nil
	}

	mapExpr, detected := hasUnboundedCollectionMapping(src, node.Span.Line)
	if detected {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Scrollable container '<%s>' renders dynamic collection '%s' without window virtualization or pagination limits. Unbounded DOM node generation severely inflates browser Presentation Delay during user interaction.", node.Tag, mapExpr),
				Hint:     "Integrate virtual windowing (e.g. '@tanstack/react-virtual') or restrict rendered items using pagination (e.g. '.slice(0, 20)').",
			},
		}
	}

	return nil
}
