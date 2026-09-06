package inp

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// LargeInteractionLayoutScopeRule mendeteksi panel drawer, modal, atau overlay besar yang beroperasi
// di dalam lingkup tata letak dokumen terbuka tanpa isolasi layout (contain: layout atau native <dialog>),
// sehingga memicu rekalkulasi reflow pada seluruh pohon DOM saat interaksi buka/tutup.
type LargeInteractionLayoutScopeRule struct{}

// NewLargeInteractionLayoutScopeRule membuat instance baru dari LargeInteractionLayoutScopeRule.
func NewLargeInteractionLayoutScopeRule() *LargeInteractionLayoutScopeRule {
	return &LargeInteractionLayoutScopeRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *LargeInteractionLayoutScopeRule) ID() string {
	return "inp.large-interaction-layout-scope"
}

// Description mengembalikan ringkasan aturan.
func (r *LargeInteractionLayoutScopeRule) Description() string {
	return "Interactive overlay or drawer element lacks layout containment or native dialog isolation, triggering document-wide reflow on toggle"
}

// Category mengembalikan nama kategori rule.
func (r *LargeInteractionLayoutScopeRule) Category() string {
	return "inp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *LargeInteractionLayoutScopeRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *LargeInteractionLayoutScopeRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Google Chrome Core Web Vitals (Interaction to Next Paint Presentation Delay)",
			"W3C CSS Containment Module Level 3 (contain: layout / contain: strict)",
			"HTML Living Standard HTMLDialogElement Top-Layer Architecture",
		},
		CoreInvariant: "Large interactive overlays and drawers must establish layout containment ('contain: layout') or utilize the browser top-layer ('<dialog>') to prevent whole-page layout recalculations during interactions.",
		Grounding: "When a large overlay, slide-over drawer, or modal toggles its visibility in the standard document flow, the browser's layout engine must invalidate ancestor and sibling boxes, triggering a full document reflow.\n\n" +
			"For complex interfaces with thousands of elements, this layout recalculation stalls the main thread and inflates Presentation Delay well past 100ms.\n\n" +
			"Using the HTML5 '<dialog>' element places the modal in the browser's isolated 'top-layer', preventing any layout impact on the document tree. Alternatively, applying CSS layout containment (e.g. 'contain-layout' or '[contain:layout]') constrains layout recalculations strictly inside the overlay container.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Document-Wide Layout Invalidation",
				Severity: "HIGH",
				Impact:   "Toggling modal/drawer state forces the layout engine to recalculate geometry for every element on the page.",
			},
			{
				Vector:   "Presentation Delay Frame Drops",
				Severity: "HIGH",
				Impact:   "Users experience visible stutters and sluggishness when expanding or collapsing sidebars, sheets, or dialogs.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Unconstrained fixed drawer in normal document flow triggering document reflow",
				Code: `<div className={` + "`fixed inset-y-0 right-0 w-96 ${isOpen ? \"block\" : \"hidden\"}`" + `}>
  <HeavySidebar />
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Native HTML5 dialog rendered in the browser's isolated top-layer",
				Code: `<dialog ref={dialogRef} className="fixed inset-y-0 right-0 w-96">
  <HeavySidebar />
</dialog>`,
			},
			{
				Language: "tsx",
				Comment:  "Explicit CSS layout containment isolating reflows to the panel",
				Code: `<div className={` + "`fixed inset-y-0 right-0 w-96 contain-layout ${isOpen ? \"block\" : \"hidden\"}`" + `}>
  <HeavySidebar />
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen overlay/drawer besar beroperasi tanpa isolasi layout.
func (r *LargeInteractionLayoutScopeRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if isUnconstrainedLayoutScope(node.Tag, node.Classes, node.RawClasses) {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Interactive overlay/drawer '<%s>' operates in the unconstrained document layout flow without layout isolation. Toggling its visibility triggers full-document reflow recalculation, degrading interaction responsiveness.", node.Tag),
				Hint:     "Isolate layout scope by migrating to native '<dialog>' (which lives in the browser top-layer) or applying layout containment ('contain-layout' / '[contain:layout]').",
			},
		}
	}

	return nil
}
