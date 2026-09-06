package a11y

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// ButtonTypeMissingRule memastikan elemen <button> di dalam formulir memiliki
// atribut type eksplisit (type="button", type="submit", atau type="reset").
type ButtonTypeMissingRule struct{}

// NewButtonTypeMissingRule membuat instance baru ButtonTypeMissingRule.
func NewButtonTypeMissingRule() *ButtonTypeMissingRule {
	return &ButtonTypeMissingRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat a11y.button-type-missing.
func (r *ButtonTypeMissingRule) ID() string {
	return "a11y.button-type-missing"
}

// Description mengembalikan deskripsi ringkas aturan.
func (r *ButtonTypeMissingRule) Description() string {
	return "Enforces explicit type attribute on <button> elements inside forms to prevent unintended form submission"
}

// Category mengembalikan nama kategori rule.
func (r *ButtonTypeMissingRule) Category() string {
	return "a11y"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *ButtonTypeMissingRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ButtonTypeMissingRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Web Content Accessibility Guidelines (WCAG) 2.2 SC 3.2.2 (On Input)",
			"HTML Living Standard 4.10.6 (The button element)",
			"React SPA Form Submission Invariants",
		},
		CoreInvariant: "<button> elements nested within a <form> must explicitly declare their type attribute (type=\"button\", type=\"submit\", or type=\"reset\").",
		Grounding: "In HTML and React applications, a <button> without an explicit type attribute defaults to type=\"submit\".\n\n" +
			"When secondary action buttons (such as \"Batal\", \"Close\", \"Preview\", or modal dismiss buttons) are placed inside a <form> without type=\"button\":\n" +
			"1. Unintended Submission: Clicking the button triggers the enclosing form's onSubmit handler and validation logic.\n" +
			"2. Unwanted Full-Page Reload: If event.preventDefault() is missing in the click handler, the browser reloads the entire page.\n" +
			"3. State Mutation Hazards: Form fields may be serialized or posted inadvertently to the server.\n\n" +
			"Charites traverses element ancestors to detect <button> tags rendered inside <form> lacking explicit type definitions.",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Button inside form lacks type, defaulting to type=\"submit\"",
				Code: `<form onSubmit={handleSubmit}>
  <input name="query" className="h-11 px-3.5 text-base border rounded-lg" />
  <button onClick={handlePreview} className="h-11 px-4 text-sm font-medium">Pratinjau</button>
  <button type="submit" className="h-11 px-4 text-sm font-medium bg-primary text-primary-foreground">Simpan</button>
</form>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Explicit type=\"button\" prevents accidental submission",
				Code: `<form onSubmit={handleSubmit}>
  <input name="query" className="h-11 px-3.5 text-base border rounded-lg" />
  <button type="button" onClick={handlePreview} className="h-11 px-4 text-sm font-medium">Pratinjau</button>
  <button type="submit" className="h-11 px-4 text-sm font-medium bg-primary text-primary-foreground">Simpan</button>
</form>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Accidental Form Submission",
				Severity: "HIGH",
				Impact:   "Clicking secondary buttons submits the form, triggers validation errors, or mutates data.",
			},
			{
				Vector:   "Full-Page Browser Reload",
				Severity: "MEDIUM",
				Impact:   "React SPAs reload the entire page unexpectedly if event.preventDefault is absent.",
			},
		},
	}
}

// Evaluate memeriksa apakah elemen <button> di dalam <form> memiliki atribut type eksplisit.
func (r *ButtonTypeMissingRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !strings.EqualFold(node.Tag, "button") {
		return nil
	}

	if t, ok := node.GetAttr("type"); ok && strings.TrimSpace(t) != "" {
		return nil
	}

	if hasSpreadProps(node.Attributes) || isDecorativeOrHidden(node.Attributes) {
		return nil
	}

	if !isInsideForm(node) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Missing explicit type attribute on <button> inside <form>. In HTML/React, buttons default to type=\"submit\" which can trigger unintended form submissions.",
			Hint:     "Add explicit type=\"button\" (or type=\"submit\" / type=\"reset\") to make the button behavior intentional.",
		},
	}
}

// isInsideForm memeriksa apakah node berada di dalam elemen <form> atau <Form> pada rantai ancestor.
func isInsideForm(node *ir.Node) bool {
	for p := node.Parent; p != nil; p = p.Parent {
		if p.Type == ir.NodeElement && (strings.EqualFold(p.Tag, "form") || p.Tag == "Form") {
			return true
		}
	}
	return false
}
