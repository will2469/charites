package ux

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules/form"
)

// NumberInputMissingBoundsRule mendeteksi <input type="number"> yang tidak memiliki batas bawah domain (min)
// eksplisit untuk mencegah kesalahan pengiriman nilai di luar domain yang valid (misal: kuantitas negatif).
type NumberInputMissingBoundsRule struct{}

// NewNumberInputMissingBoundsRule membuat instance baru dari NumberInputMissingBoundsRule.
func NewNumberInputMissingBoundsRule() *NumberInputMissingBoundsRule {
	return &NumberInputMissingBoundsRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *NumberInputMissingBoundsRule) ID() string {
	return "ux.number-input-missing-bounds"
}

// Description mengembalikan ringkasan maksud dan tujuan aturan.
func (r *NumberInputMissingBoundsRule) Description() string {
	return "Flags mathematical number inputs lacking explicit domain lower bounds (min) to prevent accidental out-of-range submissions"
}

// Category mengembalikan nama kategori rule.
func (r *NumberInputMissingBoundsRule) Category() string {
	return "ux"
}

// DefaultSeverity mengembalikan tingkat keparahan default (warn).
func (r *NumberInputMissingBoundsRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *NumberInputMissingBoundsRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Charites Design System & Form Ergonomy Policy",
			"WCAG 2.2 Success Criterion 3.3.4 (Error Prevention - Contextual Usability Guidance)",
			"ISO 9241-110:2020 (Interaction Ergonomics - Error Tolerance)",
		},
		CoreInvariant: "Number inputs representing domain quantities should declare an explicit domain lower bound ('min') to prevent erroneous or accidental out-of-range submissions.",
		Grounding: "While the HTML5 specification permits '<input type=\"number\">' without a 'min' attribute, " +
			"Charites UX design system policy requires developers to declare an explicit domain lower bound for quantitative fields.\n\n" +
			"Without a declared bound, users can inadvertently scroll or step into nonsensical negative numbers (e.g. quantity: -5, tickets: -1), " +
			"or submit values below business thresholds.\n\n" +
			"Cross-Rule Yield Hierarchy: When an input field is classified as a discrete identity code (such as postal code, account number, or NIK), " +
			"this rule intentionally yields and suppresses itself. Recommending 'min=\"0\"' on an identity field is erroneous pseudo-remediation; " +
			"the field must instead be converted to text via 'ux.number-input-identity-misuse'.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Accidental Negative Value Submission",
				Severity: "MEDIUM",
				Impact:   "Users submit nonsensical negative quantities or counts when stepper controls decrement below zero.",
			},
			{
				Vector:   "Unbounded Stepper Traversal",
				Severity: "LOW",
				Impact:   "Spinbox controls allow decrementing past business logic boundaries without visual constraint.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Quantity input lacking explicit lower bound can be decremented into negative numbers",
				Code: `<input
  type="number"
  name="quantity"
  onWheel={(e) => e.currentTarget.blur()}
/>`,
			},
			{
				Language: "tsx",
				Comment:  "Discount percentage input without lower bound",
				Code: `<Input
  type="number"
  name="discount_percent"
  onWheel={handleWheel}
/>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Explicit lower bound min='0' prevents stepping into negative numbers",
				Code: `<input
  type="number"
  name="quantity"
  min="0"
  onWheel={(e) => e.currentTarget.blur()}
/>`,
			},
			{
				Language: "tsx",
				Comment:  "Domain temperature input with valid negative bound",
				Code: `<input
  type="number"
  name="temperature"
  min="-50"
  onWheel={(e) => e.currentTarget.blur()}
/>`,
			},
		},
	}
}

// Evaluate memeriksa apakah input numerik memiliki deklarasi lower bound yang eksplisit.
func (r *NumberInputMissingBoundsRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	facts := form.ExtractInputFacts(node)
	if !facts.IsInputTag || !facts.IsNumberType {
		return nil
	}

	// Pengecualian: disabled atau readOnly tidak ditujukan untuk mutasi langsung oleh user.
	if facts.IsDisabled || facts.IsReadOnly {
		return nil
	}

	// YIELD HIERARCHY: Jika field terbukti sebagai kode identitas (identity misuse),
	// aturan missing-bounds menekan diri sendiri karena remediasi yang benar adalah mengubah type ke text,
	// bukan menambahkan min="0".
	if facts.Identifier.Class == form.IdentifierIdentity {
		return nil
	}

	// Jika developer sudah mendeklarasikan bound yang valid (static atau dynamic), lolos.
	if facts.HasDeclaredMin {
		return nil
	}

	fieldName := facts.Identifier.Value
	if fieldName == "" {
		fieldName = "input"
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Number input '%s' lacks an explicit domain lower bound.", fieldName),
			Hint:     "Declare the appropriate domain lower bound, e.g. min=\"0\", when negative values are invalid for this field (Charites UX policy; WCAG 3.3.4 & ISO 9241-110).",
		},
	}
}
