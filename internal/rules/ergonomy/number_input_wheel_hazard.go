package ergonomy

import (
	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules/form"
)

// NumberInputWheelHazardRule mendeteksi <input type="number"> yang tidak memiliki penanganan event wheel,
// mencegah mutasi nilai secara tidak sengaja ketika pengguna melakukan scrolling halaman (Scroll Hijacking).
type NumberInputWheelHazardRule struct{}

// NewNumberInputWheelHazardRule membuat instance baru dari NumberInputWheelHazardRule.
func NewNumberInputWheelHazardRule() *NumberInputWheelHazardRule {
	return &NumberInputWheelHazardRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *NumberInputWheelHazardRule) ID() string {
	return "ergonomy.number-input-wheel-hazard"
}

// Description mengembalikan ringkasan maksud dan tujuan aturan.
func (r *NumberInputWheelHazardRule) Description() string {
	return "Flags number inputs without wheel prevention to prevent unintended value mutations during scrolling (Scroll Hijacking)"
}

// Category mengembalikan nama kategori rule.
func (r *NumberInputWheelHazardRule) Category() string {
	return "ergonomy"
}

// DefaultSeverity mengembalikan tingkat keparahan default (warn).
func (r *NumberInputWheelHazardRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *NumberInputWheelHazardRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Nielsen Norman Group (Input Steppers and Number Inputs)",
			"GOV.UK Design System (Numbers and Form Inputs)",
			"W3C HTML 5.2 Section 4.10.5.1.12 (Number State)",
		},
		CoreInvariant: "Number inputs (<input type=\"number\">) must declare an explicit wheel event handler (e.g. 'onWheel={(e) => e.currentTarget.blur()}') to claim ownership and prevent accidental value mutations during page scrolling.",
		Grounding: "When a user focuses on an input with type=\"number\" and scrolls the document using a mouse wheel or trackpad, " +
			"desktop browsers intercept the scroll delta to increment or decrement the numeric value rather than moving the viewport.\n\n" +
			"This phenomenon (Scroll Hijacking) silently mutates quantitative user data (such as item quantities or monetary amounts) without the user's conscious intent. " +
			"Disabled and readOnly inputs are excluded from this rule because they are not intended for direct user mutation.\n\n" +
			"Note: Static analysis checks whether the developer explicitly declared a wheel handler (such as 'onWheel' or 'onWheelCapture') as a static protection claim. " +
			"It does not perform runtime interprocedural flow analysis on the callback body.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Accidental Value Mutation (Scroll Hijacking)",
				Severity: "HIGH",
				Impact:   "Users scrolling past focused inputs unintentionally increment or decrement orders, quantities, or amounts.",
			},
			{
				Vector:   "Silent Submission Corruption",
				Severity: "HIGH",
				Impact:   "Mutated values submit without visual warning because the user assumes their gesture only shifted the viewport.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Unprotected number input vulnerable to mouse wheel mutation during page scroll",
				Code: `<input
  type="number"
  name="quantity"
  min="1"
/>`,
			},
			{
				Language: "tsx",
				Comment:  "Unprotected component-based number input",
				Code: `<Input
  type="number"
  name="item_count"
  min="0"
/>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Blur on wheel prevents mouse wheel value increment while preserving page scrolling",
				Code: `<input
  type="number"
  name="quantity"
  min="1"
  onWheel={(e) => e.currentTarget.blur()}
/>`,
			},
			{
				Language: "tsx",
				Comment:  "Component-based input with designated wheel handler callback",
				Code: `<Input
  type="number"
  name="item_count"
  min="0"
  onWheel={handleWheel}
/>`,
			},
		},
	}
}

// Evaluate memeriksa apakah input tipe number dilindungi oleh handler wheel.
func (r *NumberInputWheelHazardRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	facts := form.ExtractInputFacts(node)
	if !facts.IsInputTag || !facts.IsNumberType {
		return nil
	}

	// Pengecualian: disabled atau readOnly tidak ditujukan untuk mutasi langsung oleh user.
	// Jika developer telah mendeklarasikan wheel handler (onWheel/onWheelCapture),
	// elemen dianggap dilindungi berdasarkan static ownership claim.
	if facts.IsDisabled || facts.IsReadOnly || facts.HasWheelHandler {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Number input lacks mouse wheel mutation prevention. Scrolling over focused number inputs causes accidental value changes (Scroll Hijacking).",
			Hint:     "Add 'onWheel={(e) => e.currentTarget.blur()}' or a wheel event handler to prevent unintended value mutations during page scrolling (GOV.UK & Nielsen Norman Group).",
		},
	}
}
