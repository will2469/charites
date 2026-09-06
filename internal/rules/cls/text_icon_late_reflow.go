package cls

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// TextIconLateReflowRule mendeteksi elemen font ikon berbasis teks ligatur
// yang tidak memiliki kotak dimensi terkunci (inline-block, size-*, overflow-hidden).
type TextIconLateReflowRule struct{}

// NewTextIconLateReflowRule membuat instance baru TextIconLateReflowRule.
func NewTextIconLateReflowRule() *TextIconLateReflowRule {
	return &TextIconLateReflowRule{}
}

// ID mengembalikan identifier kanonikal Semgrep rule.
func (r *TextIconLateReflowRule) ID() string {
	return "cls.text-icon-late-reflow"
}

// Description mengembalikan ringkasan aturan.
func (r *TextIconLateReflowRule) Description() string {
	return "Requires locked bounding dimensions on text-ligature icon elements to prevent text reflow"
}

// Category mengembalikan nama kategori rule.
func (r *TextIconLateReflowRule) Category() string {
	return "cls"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (info).
func (r *TextIconLateReflowRule) DefaultSeverity() ir.Severity {
	return ir.SeverityInfo
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *TextIconLateReflowRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Cumulative Layout Shift (CLS) Metric Specification",
			"Google Material Icons & Material Symbols Integration Guide",
		},
		CoreInvariant: "Text-ligature icon elements must lock their bounding box via 'inline-block' (or block/flex), explicit width/height (or 'size-*'), and 'overflow-hidden' to prevent raw word ligature text from expanding the layout before the icon font is loaded.",
		Grounding: "Icon fonts like Material Icons or Material Symbols render icons by substituting raw text strings (such as 'shopping_cart', 'account_circle', or 'arrow_back') with icon glyphs via OpenType ligatures.\n\n" +
			"Before the web font finishes downloading, the browser displays the fallback word text ('shopping_cart') at full length (spanning 80-120px).\n\n" +
			"When the web font suddenly loads, the word shrinks into a 24x24px glyph, causing surrounding navigation bars, buttons, and text to collapse backward and triggering Cumulative Layout Shift (CLS).\n\n" +
			"Locking the container dimensions to 'inline-block size-6 overflow-hidden' ensures the element occupies exactly 24x24px regardless of font loading state.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Button and Header Layout Contraction",
				Severity: "MEDIUM",
				Impact:   "Long ligature strings expand buttons initially, then contract suddenly when font arrives.",
			},
			{
				Vector:   "Cumulative Layout Shift (CLS)",
				Severity: "LOW",
				Impact:   "Shifts around interactive icons trigger layout recalculations in headers and toolbars.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Material icon with raw text ligature without locked box dimensions",
				Code: `<button className="flex items-center gap-2">
  <span className="material-icons">shopping_cart</span> Beli
</button>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Locked icon bounding box with inline-block, size-6, and overflow-hidden",
				Code: `<button className="flex items-center gap-2">
  <span className="material-icons inline-block size-6 overflow-hidden">shopping_cart</span> Beli
</button>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen ligatur ikon mengunci kotak dimensinya.
func (r *TextIconLateReflowRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if hasSpreadProps(node.Attributes) {
		return nil
	}

	if !isIconLigatureElement(node) {
		return nil
	}

	if hasLockedLigatureBox(node) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Text-ligature icon element <%s> lacks locked box dimensions ('inline-block', 'size-*', 'overflow-hidden'). Before the icon web font loads, raw word ligature text expands surrounding layout causing Cumulative Layout Shift (CLS).", node.Tag),
			Hint:     "Add 'inline-block size-6 overflow-hidden' (or 'w-6 h-6 inline-block overflow-hidden') to clamp the ligature bounding box to the expected icon size.",
		},
	}
}
