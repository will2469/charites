package a11y

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// FormLabelMissingControlRule memastikan komponen Shadcn UI <FormItem> yang memiliki
// <FormLabel> juga memuat <FormControl> atau elemen input interaktif yang sah.
type FormLabelMissingControlRule struct{}

// NewFormLabelMissingControlRule membuat instance baru FormLabelMissingControlRule.
func NewFormLabelMissingControlRule() *FormLabelMissingControlRule {
	return &FormLabelMissingControlRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat a11y.form-label-missing-control.
func (r *FormLabelMissingControlRule) ID() string {
	return "a11y.form-label-missing-control"
}

// Description mengembalikan deskripsi ringkas aturan.
func (r *FormLabelMissingControlRule) Description() string {
	return "Enforces that Shadcn UI <FormItem> containing <FormLabel> also contains an associated <FormControl> or input element"
}

// Category mengembalikan nama kategori rule.
func (r *FormLabelMissingControlRule) Category() string {
	return "a11y"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *FormLabelMissingControlRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *FormLabelMissingControlRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Web Content Accessibility Guidelines (WCAG) 2.2 SC 1.3.1 (Info and Relationships)",
			"W3C WAI-ARIA 1.2 Form Pattern Specifications",
			"Shadcn UI Form Component Architecture Contract",
		},
		CoreInvariant: "Every <FormItem> containing a <FormLabel> must provide an associated <FormControl> or interactive input element.",
		Grounding: "In component libraries built on Radix UI and React Hook Form (such as Shadcn UI), <FormItem> automatically generates unique IDs and wires <FormLabel> to <FormControl> via ARIA attributes.\n\n" +
			"When developers render a <FormLabel> inside a <FormItem> but forget to include <FormControl> (e.g. placing plain text, descriptions, or unassociated elements instead):\n" +
			"1. Broken Accessibility Tree: Screen reader users hear the label announcement but cannot identify or interact with the corresponding input field.\n" +
			"2. Silent Failure: React and TypeScript compile cleanly because JSX treats custom components as valid React nodes.\n" +
			"3. Form Confusion: Tapping the label does not focus any active control.\n\n" +
			"Charites inspects the component tree of each <FormItem> to guarantee relational completeness.",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "FormItem has FormLabel but lacks FormControl or input",
				Code: `<FormItem>
  <FormLabel>Nama Lengkap</FormLabel>
  <p className="text-sm text-muted-foreground">Silakan hubungi administrator</p>
</FormItem>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Proper Shadcn UI FormItem composition with FormControl",
				Code: `<FormItem>
  <FormLabel>Nama Lengkap</FormLabel>
  <FormControl>
    <Input placeholder="Nama Anda" />
  </FormControl>
</FormItem>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Orphaned Accessible Label",
				Severity: "HIGH",
				Impact:   "Assistive technology cannot pair the form label with an interactive form control.",
			},
			{
				Vector:   "Broken Click-to-Focus Affordance",
				Severity: "MEDIUM",
				Impact:   "Users tapping the label on mobile cannot focus the intended input field.",
			},
		},
	}
}

// Evaluate memeriksa hierarki komponen <FormItem> untuk memastikan keberadaan <FormControl>.
// Mematuhi 0 B/op, 0 allocs/op pada node bersih (QUAL-03).
func (r *FormLabelMissingControlRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !isFormItemTag(node.Tag) {
		return nil
	}

	hasLabel := false
	hasControl := false

	for child := range node.Walk() {
		if child == node || child.Type != ir.NodeElement {
			continue
		}

		if isFormLabelTag(child.Tag) {
			hasLabel = true
		} else if isFormControlTag(child.Tag) || isDirectInputControl(child) {
			hasControl = true
		}
	}

	// Jika memiliki label namun tidak memiliki FormControl atau input apa pun di dalamnya
	if hasLabel && !hasControl {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "<FormItem> has a <FormLabel> but is missing an associated <FormControl> or input element (Shadcn UI accessibility contract)",
				Hint:     "Wrap your form control inside <FormControl> within <FormItem> to properly associate the label with its input.",
			},
		}
	}

	return nil
}

func isFormItemTag(tag string) bool {
	return strings.EqualFold(tag, "FormItem") ||
		strings.HasSuffix(tag, ".Item") ||
		strings.HasSuffix(strings.ToLower(tag), "formitem")
}

func isFormLabelTag(tag string) bool {
	return strings.EqualFold(tag, "FormLabel") ||
		strings.HasSuffix(tag, ".Label") ||
		strings.HasSuffix(strings.ToLower(tag), "formlabel")
}

func isFormControlTag(tag string) bool {
	return strings.EqualFold(tag, "FormControl") ||
		strings.HasSuffix(tag, ".Control") ||
		strings.HasSuffix(strings.ToLower(tag), "formcontrol")
}

func isDirectInputControl(node *ir.Node) bool {
	if node == nil {
		return false
	}
	tagLower := strings.ToLower(node.Tag)
	switch tagLower {
	case "input", "select", "textarea":
		return true
	}

	// Dukungan komponen UI React populer (Input, Textarea, Select, Switch, Checkbox, Slider, RadioGroup)
	switch node.Tag {
	case "Input", "Textarea", "Select", "Switch", "Checkbox", "Slider", "RadioGroup", "Combobox":
		return true
	}

	return false
}
