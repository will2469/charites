package a11y

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// FormInputMissingNameRule mendeteksi kontrol formulir yang tidak memiliki atribut pengenal
// 'name' maupun 'id' yang diperlukan untuk autofill browser dan serialisasi form (WCAG 4.1.2).
type FormInputMissingNameRule struct{}

// NewFormInputMissingNameRule membuat instance baru FormInputMissingNameRule.
func NewFormInputMissingNameRule() *FormInputMissingNameRule {
	return &FormInputMissingNameRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat a11y.form-input-missing-name.
func (r *FormInputMissingNameRule) ID() string {
	return "a11y.form-input-missing-name"
}

// Description mengembalikan deskripsi ringkas aturan.
func (r *FormInputMissingNameRule) Description() string {
	return "Ensures form input controls declare an identifying name or id attribute for form submission and autofill (WCAG 4.1.2)"
}

// Category mengembalikan nama kategori rule.
func (r *FormInputMissingNameRule) Category() string {
	return "a11y"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *FormInputMissingNameRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *FormInputMissingNameRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"WCAG 2.2 Success Criterion 4.1.2 (Name, Role, Value - Level A)",
			"HTML5 Specification (Section 4.10.5 The input element & form submission)",
			"W3C Form Usability & Autofill Standards",
		},
		CoreInvariant: "Form input controls (<input>, <select>, <textarea>) must declare an identifying 'name' or 'id' attribute.",
		Grounding: "Browsers, password managers, and native form processors rely on 'name' and 'id' attributes to identify field context.\n\n" +
			"When form controls lack both 'name' and 'id':\n" +
			"1. Autofill Failure: Chrome, Safari, and password managers (1Password, Bitwarden) cannot identify whether an input is for email, username, or credit card, forcing users to type manually.\n" +
			"2. Silent Form Data Loss: Standard HTML <form> submission constructs FormData from controls with a 'name' attribute. Inputs without 'name' are completely omitted from the payload.\n" +
			"3. Unidentifiable Tree Nodes: Screen readers and automated testing frameworks cannot reference or locate the input deterministically.\n\n" +
			"Charites flags form controls that fail to specify at least one identifying key ('name' or 'id').",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Input missing both name and id attributes",
				Code:     `<input type="email" className="w-full border px-3 py-2" placeholder="nama@domain.com" />`,
			},
			{
				Language: "astro",
				Comment:  "Select dropdown without name or id identifier",
				Code:     `<select class="border p-2"><option value="id">Indonesia</option></select>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Properly identified input with matching name and id",
				Code:     `<input id="user-email" name="email" type="email" className="w-full border px-3 py-2" placeholder="nama@domain.com" />`,
			},
			{
				Language: "astro",
				Comment:  "Select dropdown with descriptive name and id",
				Code:     `<select name="country" id="country" class="border p-2"><option value="id">Indonesia</option></select>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Silent Form Submission Data Loss",
				Severity: "HIGH",
				Impact:   "Form inputs without 'name' attributes are excluded from native FormData submissions.",
			},
			{
				Vector:   "Password Manager & Autofill Breakdown",
				Severity: "MEDIUM",
				Impact:   "Browser autofill cannot populate user credentials, causing login drop-off.",
			},
		},
	}
}

// Evaluate memeriksa keberadaan atribut name atau id pada kontrol isian formulir.
// Mematuhi 0 B/op, 0 allocs/op pada node bersih (QUAL-03).
func (r *FormInputMissingNameRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !IsTextualInput(node) {
		return nil
	}

	if node.Attributes != nil {
		if nameVal, ok := node.Attributes["name"]; ok && CleanAttr(nameVal) != "" {
			return nil
		}
		if idVal, ok := node.Attributes["id"]; ok && CleanAttr(idVal) != "" {
			return nil
		}

		// Pengecualian library form dan spread props (react-hook-form {...register('field')}, {...props})
		for k, v := range node.Attributes {
			if strings.HasPrefix(k, "{...") || strings.Contains(k, "register") || strings.Contains(v, "register") {
				return nil
			}
		}
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Form control <%s> is missing both 'name' and 'id' attributes, breaking autofill and form serialization", node.Tag),
			Hint:     "Add 'name=\"...\"' and 'id=\"...\"' so browsers and password managers can identify and autofill this field.",
		},
	}
}
