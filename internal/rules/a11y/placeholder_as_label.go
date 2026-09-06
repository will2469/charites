package a11y

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// PlaceholderAsLabelRule mendeteksi kontrol isian formulir yang hanya mengandalkan
// atribut 'placeholder' sebagai penanda identitas tanpa adanya label persisten atau 'aria-label' (WCAG 3.3.2).
type PlaceholderAsLabelRule struct{}

// NewPlaceholderAsLabelRule membuat instance baru PlaceholderAsLabelRule.
func NewPlaceholderAsLabelRule() *PlaceholderAsLabelRule {
	return &PlaceholderAsLabelRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat a11y.placeholder-as-label.
func (r *PlaceholderAsLabelRule) ID() string {
	return "a11y.placeholder-as-label"
}

// Description mengembalikan deskripsi ringkas aturan.
func (r *PlaceholderAsLabelRule) Description() string {
	return "Flags form inputs relying solely on placeholder attributes without a persistent label or accessible name (WCAG 3.3.2)"
}

// Category mengembalikan nama kategori rule.
func (r *PlaceholderAsLabelRule) Category() string {
	return "a11y"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *PlaceholderAsLabelRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *PlaceholderAsLabelRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"WCAG 2.2 Success Criterion 3.3.2 (Labels or Instructions - Level A)",
			"WCAG 2.2 Success Criterion 1.3.1 (Info and Relationships - Level A)",
			"Nielsen Norman Group (Placeholders in Form Fields Are Harmful)",
		},
		CoreInvariant: "Form text inputs (<input>, <select>, <textarea>) declaring a 'placeholder' must provide a persistent label (<label>, <label htmlFor>, aria-label, or aria-labelledby).",
		Grounding: "Placeholder text is intended solely as an example of expected input format (e.g. 'e.g. user@example.com'), never as a replacement for a field label.\n\n" +
			"When developers omit <label> and rely exclusively on placeholder:\n" +
			"1. Vanishing Context: The moment a user types even one character, the placeholder vanishes. Users cannot review or verify what the field originally asked for.\n" +
			"2. Cognitive and Memory Strain: Users with memory impairments, cognitive disabilities, or mobile users distracted by notifications lose track of field meaning.\n" +
			"3. Screen Reader Gaps: Assistive technologies handle unlabelled placeholders inconsistently, frequently omitting them in forms-mode navigation.\n\n" +
			"Charites verifies that any input with a placeholder has an associated <label> or direct accessible name attribute.",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Input relying solely on placeholder without label or aria-label",
				Code:     `<input type="email" placeholder="Masukkan email Anda" className="border px-3 py-2" />`,
			},
			{
				Language: "astro",
				Comment:  "Textarea with placeholder but no associated label",
				Code:     `<textarea placeholder="Ketik pesan Anda di sini" class="border p-2"></textarea>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Persistent label associated with input ID plus descriptive placeholder",
				Code:     `<label htmlFor="user-email">Email</label><input id="user-email" type="email" placeholder="nama@domain.com" className="border px-3 py-2" />`,
			},
			{
				Language: "astro",
				Comment:  "Enclosing label wrapping input and descriptive prompt",
				Code:     `<label class="flex flex-col">Pesan <textarea placeholder="Ketik pesan..." class="border p-2"></textarea></label>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Vanishing Form Context",
				Severity: "HIGH",
				Impact:   "Users lose track of input requirements after typing begins, increasing submission errors.",
			},
			{
				Vector:   "WCAG 3.3.2 Non-Compliance",
				Severity: "HIGH",
				Impact:   "Flagged as a major Level A accessibility violation in public and enterprise audits.",
			},
		},
	}
}

// Evaluate memeriksa apakah kontrol ber-placeholder memiliki label persisten.
// Mematuhi 0 B/op, 0 allocs/op pada node bersih (QUAL-03).
func (r *PlaceholderAsLabelRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || node.Attributes == nil {
		return nil
	}

	if !IsTextualInput(node) {
		return nil
	}

	placeholder, ok := node.Attributes["placeholder"]
	if !ok || CleanAttr(placeholder) == "" {
		return nil
	}

	// 1. Periksa apakah memiliki atribut nama aksesibel langsung
	if hasDirectAccessibleName(node) {
		return nil
	}

	// 2. Periksa apakah dibungkus di dalam elemen <label>
	if HasEnclosingLabel(node) {
		return nil
	}

	// 3. Periksa apakah memiliki id yang diikat oleh <label htmlFor="id"> di dokumen
	if rawID, ok := node.Attributes["id"]; ok {
		cleanID := CleanAttr(rawID)
		if cleanID != "" {
			root := FindRoot(node)
			if HasAssociatedLabel(root, node, cleanID) {
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
			Message:  fmt.Sprintf("Input <%s> uses 'placeholder' as its sole identifier without a persistent <label> or accessible name (WCAG 3.3.2)", node.Tag),
			Hint:     "Provide a persistent <label htmlFor=\"...\"> or add 'aria-label' so the field identity remains visible after typing.",
		},
	}
}

func hasDirectAccessibleName(node *ir.Node) bool {
	if node.Attributes == nil {
		return false
	}
	for _, attr := range []string{"aria-label", "aria-labelledby", "title"} {
		if val, ok := node.Attributes[attr]; ok && CleanAttr(val) != "" {
			return true
		}
	}
	return false
}
