package responsive

import (
	"github.com/will2469/charites/internal/ir"
)

// MobileTextOverflowRule mendeteksi kontainer teks dengan 'whitespace-nowrap' tanpa elipsis/scroll,
// atau blok kode inline tanpa pemenggalan kata/scroll ancestor yang berisiko merobek layout mobile.
type MobileTextOverflowRule struct{}

// NewMobileTextOverflowRule membuat instance baru dari MobileTextOverflowRule.
func NewMobileTextOverflowRule() *MobileTextOverflowRule {
	return &MobileTextOverflowRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *MobileTextOverflowRule) ID() string {
	return "responsive.mobile-text-overflow"
}

// Description mengembalikan ringkasan aturan.
func (r *MobileTextOverflowRule) Description() string {
	return "Warns when whitespace-nowrap text or code blocks lack truncation, word breaking, or horizontal scroll wrappers"
}

// Category mengembalikan nama kategori rule.
func (r *MobileTextOverflowRule) Category() string {
	return "responsive"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *MobileTextOverflowRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *MobileTextOverflowRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C CSS Text Module Level 3 (Wrapping and Breaking Text)",
			"WCAG 2.2 SC 1.4.10 (Reflow - Level AA)",
		},
		CoreInvariant: "Containers declaring 'whitespace-nowrap' must provide overflow mitigation ('truncate', 'overflow-hidden', 'overflow-x-auto'), and inline '<code>' blocks must provide word breaking ('break-all', 'break-words') or horizontal scroll ancestors.",
		Grounding: "Dynamic strings such as URLs, authentication tokens, UUIDs, IBANs, and email addresses contain no whitespace. When 'whitespace-nowrap' is declared on narrow smartphone screens (360px) without truncation or scroll containment, the text forces the container beyond the viewport.\n\n" +
			"Similarly, inline code elements ('<code>') default to unbreaking monospace text. Without 'break-all' or a scrollable parent, long code snippets tear mobile page layouts.\n\n" +
			"Using 'truncate', 'break-words', or enclosing code inside a scrollable wrapper maintains layout boundaries and satisfies WCAG Reflow requirements.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Mobile Layout Breakage via Long Unbroken Strings",
				Severity: "MEDIUM",
				Impact:   "Unbroken URLs or tokens force text containers to stretch horizontally outside the 360px mobile viewport.",
			},
			{
				Vector:   "Loss of WCAG 2.2 Reflow Compliance",
				Severity: "LOW",
				Impact:   "Users must scroll both horizontally and vertically to read content at 400% zoom or compact viewports.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "whitespace-nowrap text container without truncation or scroll containment",
				Code: `<div className="whitespace-nowrap text-sm text-foreground">
  <span>Token Transaksi: {transactionHash}</span>
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Protected text container with truncate",
				Code: `<div className="whitespace-nowrap truncate text-sm text-foreground">
  <span>Token Transaksi: {transactionHash}</span>
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen mendeklarasikan whitespace-nowrap tanpa elipsis/scroll
// atau merupakan elemen code tanpa break-words/scroll wrapper.
func (r *MobileTextOverflowRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if hasNowrap(node.Classes) && !hasTextOverflowProtection(node.Classes) {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Text container declares 'whitespace-nowrap' without overflow mitigation ('truncate', 'overflow-hidden', or 'overflow-x-auto'). Long unbroken strings (URLs, tokens, emails) will force horizontal scrolling on mobile viewports.",
				Hint:     "Add 'truncate' for single-line ellipsis or wrap inside 'overflow-x-auto' if the full text must remain accessible.",
			},
		}
	}

	if node.Tag == "code" && !hasCodeWrapOrScroll(node) {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Inline <code> element lacks word breaking ('break-all', 'break-words') or horizontal scroll container wrapper. Long snippets will blow out mobile container boundaries.",
				Hint:     "Add 'break-all' or 'break-words' to the <code> element, or wrap it inside an 'overflow-x-auto' container.",
			},
		}
	}

	return nil
}
