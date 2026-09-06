package a11y

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// FormLabelCompositeControlRule mendeteksi <FormLabel> yang terikat langsung ke komponen
// komposit multi-field (misal <DateRangePicker>, <AddressFields>) tanpa <fieldset>/<legend>,
// yang menyebabkan ambiguitas nama aksesibel pada pembaca layar (screen reader).
type FormLabelCompositeControlRule struct{}

// NewFormLabelCompositeControlRule membuat instance baru FormLabelCompositeControlRule.
func NewFormLabelCompositeControlRule() *FormLabelCompositeControlRule {
	return &FormLabelCompositeControlRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat a11y.form-label-composite-control.
func (r *FormLabelCompositeControlRule) ID() string {
	return "a11y.form-label-composite-control"
}

// Description mengembalikan deskripsi ringkas aturan.
func (r *FormLabelCompositeControlRule) Description() string {
	return "Warns when <FormLabel> is directly bound to a composite multi-field control causing screen reader ambiguity"
}

// Category mengembalikan nama kategori rule.
func (r *FormLabelCompositeControlRule) Category() string {
	return "a11y"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *FormLabelCompositeControlRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *FormLabelCompositeControlRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Web Content Accessibility Guidelines (WCAG) 2.2 SC 1.3.1 (Info and Relationships)",
			"W3C WAI-ARIA 1.2 Grouping Controls (<fieldset> / <legend> pattern)",
			"HTML Living Standard Section 4.10.16 (The fieldset element)",
		},
		CoreInvariant: "A single <FormLabel> must bind to exactly one discrete interactive input; composite multi-field components must use <fieldset> with <legend>.",
		Grounding: "In Shadcn UI and Radix-based form architectures, <FormLabel> generates an HTML <label> linked via htmlFor to the single ID assigned to <FormControl>.\n\n" +
			"When developers bind a <FormLabel> to a composite multi-field component (such as <DateRangePicker>, <AddressFields>, or multiple <input> fields within one <FormControl>):\n" +
			"1. Ambiguous Screen Reader Announcement: The screen reader announces the label only for the first input or fails to associate it properly with subsequent controls.\n" +
			"2. Inability to Navigate: Keyboard users cannot tell which part of the composite input the label describes.\n" +
			"3. Accessibility Violation: WCAG 1.3.1 requires clear programmatic grouping for related multi-field data.\n\n" +
			"Charites detects composite component naming patterns and multiple inputs under a single <FormItem>.",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Single FormLabel bound directly to composite DateRangePicker",
				Code: `<FormItem>
  <FormLabel>Rentang Tanggal</FormLabel>
  <FormControl>
    <DateRangePicker fromDate={from} toDate={to} />
  </FormControl>
</FormItem>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Separated into individual FormItem controls or wrapped with fieldset",
				Code: `<div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
  <FormItem>
    <FormLabel>Tanggal Mulai</FormLabel>
    <FormControl><DatePicker date={from} /></FormControl>
  </FormItem>
  <FormItem>
    <FormLabel>Tanggal Selesai</FormLabel>
    <FormControl><DatePicker date={to} /></FormControl>
  </FormItem>
</div>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Screen Reader Control Ambiguity",
				Severity: "MEDIUM",
				Impact:   "Users of screen readers cannot determine which discrete sub-field is labeled.",
			},
			{
				Vector:   "Form Autofill Fragmentation",
				Severity: "LOW",
				Impact:   "Browser password managers and autofill engines fail to parse multi-part composite inputs.",
			},
		},
	}
}

// Evaluate memeriksa apakah <FormItem> mengikat kontrol komposit multi-field.
// Mematuhi 0 B/op, 0 allocs/op pada node bersih (QUAL-03).
func (r *FormLabelCompositeControlRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !isFormItemTag(node.Tag) {
		return nil
	}

	hasLabel := false
	var formControlNode *ir.Node

	for child := range node.Walk() {
		if child == node || child.Type != ir.NodeElement {
			continue
		}
		if isFormLabelTag(child.Tag) {
			hasLabel = true
		} else if isFormControlTag(child.Tag) {
			formControlNode = child
		}
	}

	if !hasLabel || formControlNode == nil {
		return nil
	}

	// Periksa apakah kontrol di dalam FormControl sudah dibungkus secara semantik dengan <fieldset>
	if hasFieldsetWrapper(formControlNode) {
		return nil
	}

	// 1. Periksa apakah terdapat komponen komposit yang dikenal
	compositeTag := findCompositeChild(formControlNode)
	if compositeTag != "" {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("<FormLabel> is bound to composite control <%s>, causing ambiguous screen reader labeling", compositeTag),
				Hint:     "Group multi-field inputs inside a <fieldset> with <legend>, or split them into separate <FormItem> controls.",
			},
		}
	}

	// 2. Periksa apakah terdapat lebih dari satu elemen input interaktif di dalam FormControl yang sama
	inputCount := countDirectInputElements(formControlNode)
	if inputCount > 1 {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("<FormLabel> is bound to <FormControl> containing %d distinct input elements without grouping", inputCount),
				Hint:     "Wrap multiple inputs in a <fieldset> with <legend> or give each input its own <FormItem> with dedicated <FormLabel>.",
			},
		}
	}

	return nil
}

func hasFieldsetWrapper(node *ir.Node) bool {
	if node == nil {
		return false
	}
	for child := range node.Walk() {
		if child.Type == ir.NodeElement && strings.EqualFold(child.Tag, "fieldset") {
			return true
		}
	}
	return false
}

func findCompositeChild(node *ir.Node) string {
	if node == nil {
		return ""
	}
	for child := range node.Walk() {
		if child == node || child.Type != ir.NodeElement {
			continue
		}
		tag := child.Tag
		tagLower := strings.ToLower(tag)
		// Pola komponen komposit multi-field
		if strings.HasSuffix(tagLower, "rangepicker") ||
			strings.HasSuffix(tagLower, "daterangepicker") ||
			strings.HasSuffix(tagLower, "timerangepicker") ||
			strings.HasSuffix(tagLower, "addressfields") ||
			strings.HasSuffix(tagLower, "compositefield") ||
			strings.HasSuffix(tagLower, "fieldsgroup") ||
			strings.HasSuffix(tagLower, "multiselect") ||
			strings.HasSuffix(tagLower, "phoneinputwithcountry") {
			return tag
		}
	}
	return ""
}

func countDirectInputElements(node *ir.Node) int {
	if node == nil {
		return 0
	}
	count := 0
	for child := range node.Walk() {
		if child == node || child.Type != ir.NodeElement {
			continue
		}
		tagLower := strings.ToLower(child.Tag)
		switch tagLower {
		case "input", "select", "textarea":
			count++
		}
	}
	return count
}
