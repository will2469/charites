package a11y

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// ErrorNotAnnouncedRule mendeteksi kontrol input dengan atribut 'aria-invalid'
// yang tidak terhubung secara programatis ke wadah pesan kesalahan via 'aria-describedby' (WCAG 3.3.1).
type ErrorNotAnnouncedRule struct{}

// NewErrorNotAnnouncedRule membuat instance baru ErrorNotAnnouncedRule.
func NewErrorNotAnnouncedRule() *ErrorNotAnnouncedRule {
	return &ErrorNotAnnouncedRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat a11y.error-not-announced.
func (r *ErrorNotAnnouncedRule) ID() string {
	return "a11y.error-not-announced"
}

// Description mengembalikan deskripsi ringkas aturan.
func (r *ErrorNotAnnouncedRule) Description() string {
	return "Ensures form controls with aria-invalid are programmatically linked to error messages via aria-describedby (WCAG 3.3.1)"
}

// Category mengembalikan nama kategori rule.
func (r *ErrorNotAnnouncedRule) Category() string {
	return "a11y"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *ErrorNotAnnouncedRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ErrorNotAnnouncedRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"WCAG 2.2 Success Criterion 3.3.1 (Error Identification - Level A)",
			"WCAG 2.2 Success Criterion 1.3.1 (Info and Relationships - Level A)",
			"WAI-ARIA 1.2 Form Validation Authoring Practices",
		},
		CoreInvariant: "Form input controls (<input>, <select>, <textarea>) declaring 'aria-invalid' must provide an 'aria-describedby' attribute linking to the error message element ID.",
		Grounding: "When form validation fails, declaring `aria-invalid` prompts screen readers to announce that the field contains an error.\n\n" +
			"However, if the control lacks `aria-describedby` referencing the error description container:\n" +
			"1. Blind Error State: Screen reader users are informed that their input is invalid, but are left completely unaware of *why* or what format is expected.\n" +
			"2. Inaccessible Visual Feedback: Visual error banners with `text-destructive` are visible to sighted users but completely disconnected in the accessibility tree.\n" +
			"3. Severe WCAG 3.3.1 Level A Non-Compliance: Automatic audit failure in enterprise accessibility and European Accessibility Act (EAA) reviews.\n\n" +
			"Charites verifies that any input with active or dynamic `aria-invalid` includes an accompanying `aria-describedby`.",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Dynamic aria-invalid without aria-describedby to the error text",
				Code:     `<input id="email" aria-invalid={!!errors.email} className="border-destructive" />`,
			},
			{
				Language: "astro",
				Comment:  "Static aria-invalid without programmatic error connection",
				Code:     `<input id="username" aria-invalid="true" class="border-red-500" />`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Input wired to error message ID via aria-describedby",
				Code:     `<input id="email" aria-invalid={!!errors.email} aria-describedby={errors.email ? "email-error" : undefined} className="border-destructive" />`,
			},
			{
				Language: "astro",
				Comment:  "Input linked to error container with role='alert'",
				Code:     `<input id="username" aria-invalid="true" aria-describedby="username-error" class="border-red-500" />`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Unannounced Form Validation Errors",
				Severity: "HIGH",
				Impact:   "Screen reader users cannot determine why form submission failed or how to correct invalid inputs.",
			},
			{
				Vector:   "WCAG 3.3.1 Level A Non-Compliance",
				Severity: "HIGH",
				Impact:   "Immediate compliance failure during accessibility audits.",
			},
		},
	}
}

// Evaluate memeriksa apakah kontrol ber-aria-invalid memiliki aria-describedby.
// Mematuhi 0 B/op, 0 allocs/op pada node bersih (QUAL-03).
func (r *ErrorNotAnnouncedRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || node.Attributes == nil {
		return nil
	}

	if !IsTextualInput(node) {
		return nil
	}

	ariaInvalid, ok := node.Attributes["aria-invalid"]
	if !ok {
		return nil
	}

	cleanInvalid := CleanAttr(ariaInvalid)
	if cleanInvalid == "" || strings.EqualFold(cleanInvalid, "false") {
		return nil
	}

	ariaDescribedBy, hasDescribedBy := node.Attributes["aria-describedby"]
	if !hasDescribedBy || CleanAttr(ariaDescribedBy) == "" {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Input <%s> declares 'aria-invalid' but lacks 'aria-describedby' linking to its error message container (WCAG 3.3.1)", node.Tag),
				Hint:     "Add 'aria-describedby=\"<error-message-id>\"' to programmatically link the input to its validation error.",
			},
		}
	}

	return nil
}
