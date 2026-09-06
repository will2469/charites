package a11y

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// LabelMissingControlRule memverifikasi bahwa atribut 'htmlFor' atau 'for' pada elemen <label>
// benar-benar merujuk ke ID kontrol input yang eksis di dalam berkas yang sama (WCAG 1.3.1).
type LabelMissingControlRule struct{}

// NewLabelMissingControlRule membuat instance baru LabelMissingControlRule.
func NewLabelMissingControlRule() *LabelMissingControlRule {
	return &LabelMissingControlRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat a11y.label-missing-control.
func (r *LabelMissingControlRule) ID() string {
	return "a11y.label-missing-control"
}

// Description mengembalikan deskripsi ringkas aturan.
func (r *LabelMissingControlRule) Description() string {
	return "Ensures label htmlFor attributes match an existing input control ID in the same document (WCAG 1.3.1)"
}

// Category mengembalikan nama kategori rule.
func (r *LabelMissingControlRule) Category() string {
	return "a11y"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *LabelMissingControlRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *LabelMissingControlRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"WCAG 2.2 Success Criterion 1.3.1 (Info and Relationships - Level A)",
			"W3C WAI-ARIA Authoring Practices (Form Labeling)",
			"HTML5 Specification (Section 4.10.4 The label element)",
		},
		CoreInvariant: "Every <label> declaring an 'htmlFor' (or 'for') attribute must match the 'id' of an existing element in the same file.",
		Grounding: "HTML form labels use the `htmlFor` (or `for`) attribute to create a programmatic bond between text instructions and form controls.\n\n" +
			"When developers introduce typos (e.g. `htmlFor=\"user_id\"` vs `<input id=\"userId\">`) or delete inputs without updating labels:\n" +
			"1. Broken Click Targets: Clicking the label fails to focus or toggle the control, frustrating desktop and mobile touch users alike.\n" +
			"2. Screen Reader Disconnect: Screen readers read the label as isolated static text, leaving the actual input unannounced and unlabelled.\n" +
			"3. Linter False Negatives: Conventional linters only check that `htmlFor` exists as a string, completely ignoring that the target ID is missing.\n\n" +
			"Charites traverses the document AST symbol table to confirm that every declared `htmlFor` target actually exists.",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Typo in htmlFor target (user_id vs userId)",
				Code:     `<label htmlFor="user_id">ID Pengguna</label><input id="userId" className="border px-3 py-2" />`,
			},
			{
				Language: "astro",
				Comment:  "Label referencing non-existent control ID",
				Code:     `<label for="missing-input">Nama Lengkap</label><input id="fullname" class="border p-2" />`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Matching htmlFor and input id",
				Code:     `<label htmlFor="userId">ID Pengguna</label><input id="userId" className="border px-3 py-2" />`,
			},
			{
				Language: "astro",
				Comment:  "Accurate for reference matching input element",
				Code:     `<label for="fullname">Nama Lengkap</label><input id="fullname" class="border p-2" />`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Broken Label Association",
				Severity: "HIGH",
				Impact:   "Screen reader users receive unlabelled inputs; clicking labels fails to focus controls.",
			},
			{
				Vector:   "WCAG 1.3.1 Non-Compliance",
				Severity: "HIGH",
				Impact:   "Critical Level A accessibility non-compliance.",
			},
		},
	}
}

// Evaluate memeriksa kecocokan target htmlFor/for terhadap elemen ber-id di dokumen.
// Mematuhi 0 B/op, 0 allocs/op pada node bersih (QUAL-03).
func (r *LabelMissingControlRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || !strings.EqualFold(node.Tag, "label") || node.Attributes == nil {
		return nil
	}

	targetID := ""
	if v, ok := node.Attributes["htmlFor"]; ok {
		targetID = v
	} else if v, ok := node.Attributes["for"]; ok {
		targetID = v
	}

	if targetID == "" {
		return nil
	}

	cleanTarget := CleanAttr(targetID)
	if cleanTarget == "" {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "<label> declares an empty 'htmlFor' attribute without a target control ID",
				Hint:     "Specify a valid control ID in 'htmlFor' or remove the empty attribute.",
			},
		}
	}

	// Jangan berikan false-positive jika nilai htmlFor merupakan ekspresi dinamis interpolasi
	if strings.ContainsAny(cleanTarget, "${}+()?:") {
		return nil
	}

	root := FindRoot(node)
	if root == nil {
		return nil
	}

	if !HasDocumentElementWithID(root, cleanTarget) {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("<label> attribute htmlFor=\"%s\" does not match any element ID in this document (WCAG 1.3.1)", cleanTarget),
				Hint:     fmt.Sprintf("Ensure an element exists with matching id=\"%s\" or fix the ID typo in htmlFor.", cleanTarget),
			},
		}
	}

	return nil
}
