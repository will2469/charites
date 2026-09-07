package ux

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules/form"
)

// NumberInputIdentityMisuseRule mendeteksi penggunaan 'type="number"' pada field identitas diskrit,
// serial, atau kode otentikasi (misal: postal_code, account_number, phone_number, security_pin).
type NumberInputIdentityMisuseRule struct{}

// NewNumberInputIdentityMisuseRule membuat instance baru dari NumberInputIdentityMisuseRule.
func NewNumberInputIdentityMisuseRule() *NumberInputIdentityMisuseRule {
	return &NumberInputIdentityMisuseRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *NumberInputIdentityMisuseRule) ID() string {
	return "ux.number-input-identity-misuse"
}

// Description mengembalikan ringkasan maksud dan tujuan aturan.
func (r *NumberInputIdentityMisuseRule) Description() string {
	return "Flags type=\"number\" on discrete identity codes and serial identifiers where mathematical operations are meaningless"
}

// Category mengembalikan nama kategori rule.
func (r *NumberInputIdentityMisuseRule) Category() string {
	return "ux"
}

// DefaultSeverity mengembalikan tingkat keparahan default (warn).
func (r *NumberInputIdentityMisuseRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *NumberInputIdentityMisuseRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"GOV.UK Design System (Numbers vs Numeric Text)",
			"Nielsen Norman Group (Number Inputs vs. Numeric Text)",
			"W3C HTML 5.2 Section 4.10.5.1.12 (Number State vs Text State)",
		},
		CoreInvariant: "Inputs capturing discrete identity, serial, or cryptographic tokens (e.g. postal codes, account numbers, phone numbers, security PINs) must use 'type=\"text\"' with 'inputMode=\"numeric\"' rather than 'type=\"number\"'.",
		Grounding: "Browsers treat '<input type=\"number\">' as IEEE-754 floating-point mathematical numbers, which causes severe UX defects and data corruption when applied to identity codes:\n\n" +
			"1. Leading zeros are silently stripped by the browser parser (e.g. postal code '01234' becomes '1234').\n" +
			"2. Numbers exceeding 16 digits (such as bank accounts or payment cards) lose precision or collapse into exponential scientific notation ('1e+16').\n" +
			"3. Unhelpful stepper spinbuttons appear in the UI, prompting users to 'increment' or 'decrement' an identity code where arithmetic is nonsensical.\n" +
			"4. Screen readers announce spinbox increments that confuse assistive technology users.\n\n" +
			"Remediation requires switching to '<input type=\"text\" inputMode=\"numeric\" pattern=\"[0-9]*\" />' (or 'type=\"tel\"' for phone numbers), " +
			"which provides mobile numeric keypads without stripping zeros or enabling stepper arrows.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Loss of Leading Zeros",
				Severity: "CRITICAL",
				Impact:   "Postal codes, account numbers, and identity codes starting with 0 are corrupted during form serialization.",
			},
			{
				Vector:   "IEEE-754 Precision Truncation",
				Severity: "CRITICAL",
				Impact:   "Long card numbers or identifiers over 15 digits lose trailing digits or format as exponential notation.",
			},
			{
				Vector:   "Spinbox Steppers on Identity Fields",
				Severity: "LOW",
				Impact:   "Confuses users and screen readers by implying identity codes can be arithmetically adjusted.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Using type='number' on postal code strips leading zeros",
				Code: `<input
  type="number"
  name="postal_code"
  placeholder="01234"
/>`,
			},
			{
				Language: "tsx",
				Comment:  "Using type='number' on account number risks precision truncation",
				Code: `<Input
  type="number"
  name="account_number"
/>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Text input with numeric inputMode preserves leading zeros and launches mobile number keypad",
				Code: `<input
  type="text"
  inputMode="numeric"
  pattern="[0-9]*"
  name="postal_code"
  autoComplete="postal-code"
/>`,
			},
			{
				Language: "tsx",
				Comment:  "Telephone field uses dedicated type='tel' with autocomplete",
				Code: `<input
  type="tel"
  name="phone_number"
  autoComplete="tel"
/>`,
			},
		},
	}
}

// Evaluate memeriksa apakah input tipe number digunakan pada field identitas semantik.
func (r *NumberInputIdentityMisuseRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	facts := form.ExtractInputFacts(node)
	if !facts.IsInputTag || !facts.IsNumberType {
		return nil
	}

	if facts.Identifier.Class != form.IdentifierIdentity {
		return nil
	}

	fieldName := facts.Identifier.Value
	if fieldName == "" {
		fieldName = facts.Identifier.Matched
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Identity field '%s' uses 'type=\"number\"'. Number inputs strip leading zeros, produce exponential notation (e.g. 1e+16), and display useless stepper arrows.", fieldName),
			Hint:     "Use '<input type=\"text\" inputMode=\"numeric\" pattern=\"[0-9]*\" />' (or 'type=\"tel\"' for phone numbers) to preserve leading zeros and formatting (GOV.UK & Nielsen Norman Group).",
		},
	}
}
