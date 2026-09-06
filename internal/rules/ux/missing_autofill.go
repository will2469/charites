package ux

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// MissingAutofillRule menegakkan penggunaan atribut autocomplete HTML Living Standard
// pada field data pribadi (PII), kredensial login (password), dan pembayaran (credit card) sesuai Hukum Tesler.
type MissingAutofillRule struct{}

// NewMissingAutofillRule membuat instance baru dari MissingAutofillRule.
func NewMissingAutofillRule() *MissingAutofillRule {
	return &MissingAutofillRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *MissingAutofillRule) ID() string {
	return "ux.missing-autofill"
}

// Description mengembalikan ringkasan aturan.
func (r *MissingAutofillRule) Description() string {
	return "Enforces W3C Living Standard autocomplete attributes on personal identity, credential, and payment form inputs (Tesler's Law)"
}

// Category mengembalikan nama kategori rule.
func (r *MissingAutofillRule) Category() string {
	return "ux"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *MissingAutofillRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *MissingAutofillRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"HTML Living Standard Section 4.10.18.7 (Autofill / The autocomplete attribute)",
			"WCAG 2.2 Success Criterion 3.3.7 (Redundant Entry - Level A)",
			"Tesler's Law of Conservation of Complexity",
		},
		CoreInvariant: "Form controls collecting Personally Identifiable Information (PII), authentication credentials, or financial payment details must declare valid HTML autocomplete attributes and never disable autofill on sensitive fields.",
		Grounding: "According to Tesler's Law, every application has an inherent amount of irreducible complexity. " +
			"The design decision is whether the software absorbs this complexity or forces it upon the human user.\n\n" +
			"Entering email addresses, physical street addresses, telephone numbers, and generated complex passwords manually on every single website " +
			"forces extreme cognitive friction and typing typos upon users. " +
			"Modern browsers, OS keychains, and third-party password managers rely on standardized W3C 'autocomplete' tokens " +
			"(e.g. 'current-password', 'new-password', 'email', 'tel', 'cc-number') to securely fill verified data.\n\n" +
			"Explicitly setting 'autocomplete=\"off\"' on password or credit card fields is a severe antipattern that breaks password generation, encourages weak password reuse, and degrades account security.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Weak Password Reuse & Credential Hijacking",
				Severity: "HIGH",
				Impact:   "Blocking or omitting password autofill forces users to type passwords manually, leading to short, memorable, and easily cracked credentials.",
			},
			{
				Vector:   "Form Abandonment & Redundant Data Entry Friction",
				Severity: "MEDIUM",
				Impact:   "Users abandon multi-field checkout and registration flows when forced to re-type address and phone details manually.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Password input lacking autocomplete attribute",
				Code: `<input
  type="password"
  name="password"
  placeholder="Masukkan kata sandi..."
  className="border rounded px-3 py-2"
/>`,
			},
			{
				Language: "astro",
				Comment:  "Payment input explicitly disabling autocomplete",
				Code: `<input
  type="text"
  name="cc-number"
  autocomplete="off"
  placeholder="Nomor Kartu Kredit"
  class="border rounded px-3 py-2"
/>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Explicit current-password token assisting password managers",
				Code: `<input
  type="password"
  name="password"
  autoComplete="current-password"
  placeholder="Masukkan kata sandi..."
  className="border rounded px-3 py-2"
/>`,
			},
			{
				Language: "astro",
				Comment:  "Compliant contact field with valid email autocomplete token",
				Code: `<input
  type="email"
  name="user_email"
  autocomplete="email"
  placeholder="name@example.com"
  class="border rounded px-3 py-2"
/>`,
			},
		},
	}
}

// Evaluate memeriksa apakah input PII/kredensial/pembayaran memiliki atribut autocomplete yang sesuai.
func (r *MissingAutofillRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	category, expectedToken, isSevere, identifier := detectAutofillCategory(node)
	if category == "" {
		return nil
	}

	autoVal, hasAuto := getAttrCaseInsensitive(node, "autocomplete", "autoComplete")
	autoClean := cleanAttrValue(autoVal)

	// Kasus 1: Atribut autocomplete="off" sengaja diset pada field sensitif
	if autoClean == "off" {
		sev := ir.SeverityInfo
		if isSevere {
			sev = ir.SeverityWarn
		}

		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: sev,
				Message: fmt.Sprintf(
					"Form field %q explicitly disables autofill ('autocomplete=\"off\"'). Modern password managers and browser autofill rely on autocomplete tokens to protect credentials.",
					identifier,
				),
				Hint: fmt.Sprintf("Remove 'autocomplete=\"off\"' and declare 'autocomplete=%q' instead (Tesler's Law).", expectedToken),
			},
		}
	}

	// Kasus 2: Atribut autocomplete tidak ada atau kosong
	if !hasAuto || autoClean == "" {
		sev := ir.SeverityInfo
		if isSevere {
			sev = ir.SeverityWarn
		}

		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: sev,
				Message: fmt.Sprintf(
					"Form input %q collecting %s data lacks an 'autocomplete' attribute. Password managers and browser autofill cannot assist users without standard W3C tokens.",
					identifier, category,
				),
				Hint: fmt.Sprintf("Add 'autocomplete=%q' to enable browser credential and autofill management (Tesler's Law).", expectedToken),
			},
		}
	}

	// Kasus 3: Memiliki token autocomplete non-empty selain "off" -> Compliant
	return nil
}
