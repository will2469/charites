package responsive

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// HorizontalOverflowRule mendeteksi penggunaan 'overflow-x-scroll' tanpa pembatas fluida 'w-full'
// yang memaksa rel scrollbar kaku permanen dan merusak rantai gestur usap pada layar sentuh.
type HorizontalOverflowRule struct{}

// NewHorizontalOverflowRule membuat instance baru dari HorizontalOverflowRule.
func NewHorizontalOverflowRule() *HorizontalOverflowRule {
	return &HorizontalOverflowRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *HorizontalOverflowRule) ID() string {
	return "responsive.horizontal-overflow"
}

// Description mengembalikan ringkasan aturan.
func (r *HorizontalOverflowRule) Description() string {
	return "Warns when unconstrained overflow-x-scroll is declared without fluid width boundary or dynamic auto-scrolling"
}

// Category mengembalikan nama kategori rule.
func (r *HorizontalOverflowRule) Category() string {
	return "responsive"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *HorizontalOverflowRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *HorizontalOverflowRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C CSS Overflow Module Level 3",
			"Google Web Vitals (Preventing Unintended Layout Shifts)",
			"iOS WebKit Touch Gesture Chaining Guidelines",
		},
		CoreInvariant: "Horizontal scroll containers on mobile baseline must declare fluid boundaries ('w-full', 'max-w-full') and use dynamic scrolling ('overflow-x-auto') rather than forced permanent scrollbar rails ('overflow-x-scroll').",
		Grounding: "Declaring 'overflow-x-scroll' directly on mobile baseline forces WebKit and Chromium browsers to render a permanent, unyielding scrollbar rail even when content fits within the viewport.\n\n" +
			"Furthermore, when horizontal scrolling lacks explicit boundary containment ('w-full' or 'max-w-full'), touch drag events can bleed into root document scrolling, causing disorienting horizontal page wobble.\n\n" +
			"Using 'overflow-x-auto w-full' ensures content only scrolls when overflowing, preserves natural gesture chaining, and prevents horizontal page wobble.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Forced Permanent Scrollbar Rails",
				Severity: "MEDIUM",
				Impact:   "Rigid horizontal scrollbar tracks clutter compact mobile viewports even when content fits completely.",
			},
			{
				Vector:   "Broken Vertical Touch Gesture Chaining",
				Severity: "LOW",
				Impact:   "Users attempting vertical swipe navigation get trapped inside unconstrained horizontal scroll containers.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Unbounded overflow-x-scroll forcing persistent scrollbar rail",
				Code: `<div className="overflow-x-scroll">
  <div className="flex gap-4">
    <div className="p-4 bg-card">Item 1</div>
  </div>
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Fluid container using dynamic overflow-x-auto",
				Code: `<div className="w-full overflow-x-auto">
  <div className="flex gap-4 min-w-max">
    <div className="p-4 bg-card">Item 1</div>
  </div>
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen mendeklarasikan overflow-x-scroll tanpa pembatas fluida w-full.
func (r *HorizontalOverflowRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || len(node.Classes) == 0 {
		return nil
	}

	hasUnprefixedScroll := false
	for _, cls := range node.Classes {
		if cls == "overflow-x-scroll" {
			hasUnprefixedScroll = true
			break
		}
	}

	if !hasUnprefixedScroll {
		return nil
	}

	if hasFluidWidthBoundary(node.Classes) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Unconstrained %q declared on mobile baseline without fluid boundary ('w-full' or 'max-w-full'). Forced horizontal scrollbars cause permanent rails on touch devices and disrupt swipe gesture chaining.", "overflow-x-scroll"),
			Hint:     "Use 'overflow-x-auto w-full' to enable dynamic touch scrolling without permanent scrollbar rails or page wobble.",
		},
	}
}
