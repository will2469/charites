package responsive

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// FlexChildOverflowRule mendeteksi elemen anak langsung dari kontainer flex yang memuat teks
// atau konten dinamis tetapi tidak mendeklarasikan batas peredam 'min-w-0', memicu gotcha 'min-width: auto'.
type FlexChildOverflowRule struct{}

// NewFlexChildOverflowRule membuat instance baru dari FlexChildOverflowRule.
func NewFlexChildOverflowRule() *FlexChildOverflowRule {
	return &FlexChildOverflowRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *FlexChildOverflowRule) ID() string {
	return "responsive.flex-child-overflow"
}

// Description mengembalikan ringkasan aturan.
func (r *FlexChildOverflowRule) Description() string {
	return "Warns when a flex child containing text or dynamic content lacks min-w-0, causing min-width: auto container blowout"
}

// Category mengembalikan nama kategori rule.
func (r *FlexChildOverflowRule) Category() string {
	return "responsive"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *FlexChildOverflowRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *FlexChildOverflowRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C CSS Flexible Box Layout Module Level 1 (Section 4.5: Implied Minimum Size of Flex Items)",
			"MDN Flexbox Gotchas: Min-Width Auto on Flex Items",
		},
		CoreInvariant: "Direct flex item children containing text or dynamic content must declare 'min-w-0' (or 'overflow-hidden') to override the CSS default 'min-width: auto' behavior.",
		Grounding: "The CSS Flexbox specification defines that flex items default to 'min-width: auto' rather than 'min-width: 0'. Consequently, a flex child will refuse to shrink below the intrinsic width of its text or content.\n\n" +
			"When a flex child encloses long paragraphs, code blocks, or dynamic strings, the flex child forces the parent container and mobile viewport to expand beyond 100vw, completely breaking text truncation ('truncate') and causing horizontal overflow.\n\n" +
			"Adding 'min-w-0' to the flex item overrides the implied minimum size, allowing text truncation and responsive shrinkage to function correctly.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Flex Container Viewport Blowout",
				Severity: "MEDIUM",
				Impact:   "Flex items refuse to shrink below long content strings, blowing out the parent container beyond 100vw.",
			},
			{
				Vector:   "Broken Text Truncation",
				Severity: "LOW",
				Impact:   "CSS 'truncate' fails completely on nested text elements because the parent flex item has no minimum width boundary.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Flex child with text content lacking min-w-0",
				Code: `<div className="flex items-center gap-4">
  <div className="w-full">
    <p className="truncate">{userDescription}</p>
  </div>
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Flex child protected with min-w-0",
				Code: `<div className="flex items-center gap-4">
  <div className="min-w-0 w-full">
    <p className="truncate">{userDescription}</p>
  </div>
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen merupakan flex child yang memuat teks/dinamis tanpa min-w-0.
func (r *FlexChildOverflowRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || node.Parent == nil {
		return nil
	}

	if !isFlexContainer(node.Parent) {
		return nil
	}

	if hasFlexChildMinBoundary(node) {
		return nil
	}

	if !hasPotentiallyOverflowingContent(node) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Flex item child <%s> containing text or dynamic content lacks 'min-w-0'. In CSS Flexbox, flex items default to 'min-width: auto', causing container blowout and breaking nested text truncation.", node.Tag),
			Hint:     "Add 'min-w-0' to this flex item to override the default minimum size and enable responsive shrinkage.",
		},
	}
}
