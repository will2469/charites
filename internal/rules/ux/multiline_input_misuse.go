package ux

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules/form"
)

// MultilineInputMisuseRule mendeteksi penggunaan input teks satu baris (<Input> atau <input>)
// untuk field catatan, keterangan, deskripsi, atau alasan yang semestinya memakai <Textarea>.
type MultilineInputMisuseRule struct{}

// NewMultilineInputMisuseRule membuat instance baru dari MultilineInputMisuseRule.
func NewMultilineInputMisuseRule() *MultilineInputMisuseRule {
	return &MultilineInputMisuseRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *MultilineInputMisuseRule) ID() string {
	return "ux.multiline-input-misuse"
}

// Description mengembalikan ringkasan maksud dan tujuan aturan.
func (r *MultilineInputMisuseRule) Description() string {
	return "Flags single-line <Input> elements used for open-ended commentary, notes, or descriptions instead of <Textarea>"
}

// Category mengembalikan nama kategori rule.
func (r *MultilineInputMisuseRule) Category() string {
	return "ux"
}

// DefaultSeverity mengembalikan tingkat keparahan default (warn).
func (r *MultilineInputMisuseRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *MultilineInputMisuseRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"GOV.UK Design System (Text input vs Textarea for multi-line user answers)",
			"Nielsen Norman Group (Textareas for Multi-Line Text Input)",
			"W3C WCAG 2.2 Supporting Context (Success Criterion 3.3.2 Labels or Instructions)",
		},
		CoreInvariant: "Form fields intended for open-ended commentary, descriptions, or notes must use multiline <Textarea> controls, never single-line <input>.",
		Grounding: "Form fields intended for open-ended commentary, descriptions, reasons, or notes (identified by keywords such as 'keterangan', 'catatan', 'notes', 'deskripsi', 'description', 'alasan', 'komentar') frequently misuse single-line '<input>' controls.\n\n" +
			"When users type explanatory sentences into a single-line text box:\n" +
			"1. Text Overflow: The text scrolls horizontally out of view, preventing users from reviewing or verifying their responses without manual cursor scrolling.\n" +
			"2. Line Break Inability: Users cannot press 'Enter' to structure their thoughts into paragraphs or distinct points.\n" +
			"3. Usability Failure: According to GOV.UK Design System and Nielsen Norman Group form design guidelines, text inputs must only be used for short single-line entries, whereas open-ended commentary and notes require multiline '<Textarea rows={3}>' controls.\n\n" +
			"Charites inspects single-line input controls across attribute channels ('name', 'id', 'placeholder', 'aria-label') while respecting single-line qualifiers (such as 'note_title', 'reason_code', 'short_description') to ensure zero false positives on legitimate single-line fields.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Input Review Obstruction",
				Severity: "HIGH",
				Impact:   "Users cannot read their full input at a glance, increasing submission errors on critical notes and reasons.",
			},
			{
				Vector:   "Cognitive Load & Text Formatting Frustration",
				Severity: "MEDIUM",
				Impact:   "Users are unable to insert line breaks or structure thoughts into paragraphs in single-line controls.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Single-line input for open-ended notes/keterangan (Astrades specimen)",
				Code: `<Label htmlFor="keterangan-input">Keterangan (Opsional)</Label>
<Input id="keterangan-input" placeholder="Catatan tambahan (bila ada)" />`,
			},
			{
				Language: "astro",
				Comment:  "Single-line text input for rejection reason in Astro template",
				Code: `<label for="rejection-reason">Alasan Penolakan</label>
<input id="rejection-reason" name="alasan_penolakan" type="text" />`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Multiline Textarea for readable notes",
				Code: `<Label htmlFor="keterangan-input">Keterangan (Opsional)</Label>
<Textarea id="keterangan-input" placeholder="Catatan tambahan (bila ada)" rows={3} />`,
			},
			{
				Language: "astro",
				Comment:  "Multiline textarea for cancellation reason",
				Code: `<label for="rejection-reason">Alasan Penolakan</label>
<textarea id="rejection-reason" name="alasan_penolakan" rows={3}></textarea>`,
			},
		},
	}
}

// Evaluate memeriksa apakah kontrol input teks satu baris digunakan untuk konten multiline.
func (r *MultilineInputMisuseRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	facts := form.ExtractInputFacts(node)
	if !facts.IsInputTag {
		return nil
	}
	if facts.TypeClass != form.InputTypeText {
		return nil
	}
	if facts.Content.Intent != form.ContentIntentMultiline {
		return nil
	}

	sourceName := facts.Content.Value
	if sourceName == "" {
		sourceName = facts.Content.Matched
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Single-line input <%s> is used for multiline commentary ('%s'); use <Textarea> instead.", node.Tag, sourceName),
			Hint:     fmt.Sprintf("Replace <%s> with <Textarea rows={3}> to allow users to review and format multiline text without horizontal scrolling (GOV.UK & Nielsen Norman Group).", node.Tag),
		},
	}
}
