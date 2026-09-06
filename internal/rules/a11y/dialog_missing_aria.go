package a11y

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// DialogMissingAriaRule memastikan custom modal dialog dengan role="dialog"
// atau role="alertdialog" memiliki aria-modal="true" dan nama aksesibel (aria-labelledby atau aria-label).
type DialogMissingAriaRule struct{}

// NewDialogMissingAriaRule membuat instance baru DialogMissingAriaRule.
func NewDialogMissingAriaRule() *DialogMissingAriaRule {
	return &DialogMissingAriaRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat a11y.dialog-missing-aria.
func (r *DialogMissingAriaRule) ID() string {
	return "a11y.dialog-missing-aria"
}

// Description mengembalikan deskripsi ringkas aturan.
func (r *DialogMissingAriaRule) Description() string {
	return "Enforces that custom modal dialogs declare aria-modal=\"true\" and have an accessible name"
}

// Category mengembalikan nama kategori rule.
func (r *DialogMissingAriaRule) Category() string {
	return "a11y"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *DialogMissingAriaRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *DialogMissingAriaRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Web Content Accessibility Guidelines (WCAG) 2.2 SC 4.1.2 (Name, Role, Value)",
			"W3C WAI-ARIA 1.2 Dialog (Modal) Pattern Specification",
			"W3C WAI-ARIA Authoring Practices Guide (APG)",
		},
		CoreInvariant: "Any element with role=\"dialog\" or role=\"alertdialog\" must specify aria-modal=\"true\" and declare an accessible name via aria-labelledby or aria-label.",
		Grounding: "In web applications, declaring role=\"dialog\" on a <div> establishes the container as an accessible dialog in the accessibility tree.\n\n" +
			"However, role=\"dialog\" alone is insufficient without two mandatory attributes:\n" +
			"1. Modal Boundary (aria-modal=\"true\"): Without aria-modal=\"true\", screen readers treat the dialog as a non-modal popup, allowing reading cursor navigation to bleed into background page elements behind the backdrop.\n" +
			"2. Accessible Name (aria-labelledby or aria-label): Screen readers announce \"dialog\" but cannot state what the dialog is about (e.g. \"Konfirmasi Hapus Data\") if aria-labelledby or aria-label is omitted.\n\n" +
			"Charites inspects custom dialog nodes to verify both boundary demarcation and accessible naming.",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Custom dialog missing aria-modal and accessible name",
				Code: `<div role="dialog" className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
  <div className="bg-background p-6 rounded-xl">
    <h2>Konfirmasi Tindakan</h2>
    <p>Apakah Anda ingin melanjutkan?</p>
  </div>
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Compliant custom modal dialog with aria-modal and aria-labelledby",
				Code: `<div role="dialog" aria-modal="true" aria-labelledby="dialog-title" className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
  <div className="bg-background p-6 rounded-xl">
    <h2 id="dialog-title" className="text-lg font-semibold">Konfirmasi Tindakan</h2>
    <p className="text-sm text-muted-foreground mt-2">Apakah Anda ingin melanjutkan?</p>
  </div>
</div>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Screen Reader Background Bleed",
				Severity: "HIGH",
				Impact:   "Screen readers navigate into background DOM elements behind the dialog overlay.",
			},
			{
				Vector:   "Unnamed Modal Context",
				Severity: "HIGH",
				Impact:   "Assistive technologies announce an unnamed generic dialog without purpose or context.",
			},
		},
	}
}

// Evaluate memeriksa apakah elemen dialog memiliki aria-modal="true" dan nama aksesibel.
func (r *DialogMissingAriaRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	role, ok := node.GetAttr("role")
	if !ok {
		return nil
	}

	cleanRole := strings.ToLower(CleanAttr(role))
	if cleanRole != "dialog" && cleanRole != "alertdialog" {
		return nil
	}

	if hasSpreadProps(node.Attributes) || isDecorativeOrHidden(node.Attributes) {
		return nil
	}

	hasModal := isAriaModalDeclared(node)
	hasName := hasAccessibleDialogName(node)

	if hasModal && hasName {
		return nil
	}

	var msg, hint string
	switch {
	case !hasModal && !hasName:
		msg = "Custom dialog lacks aria-modal=\"true\" and an accessible name (aria-labelledby or aria-label)."
		hint = "Add aria-modal=\"true\" and wire aria-labelledby to the dialog title element."
	case !hasModal:
		msg = "Custom dialog lacks aria-modal=\"true\" attribute to prevent assistive technology focus bleed."
		hint = "Add aria-modal=\"true\" to properly isolate dialog contents from background elements."
	default:
		msg = "Custom dialog lacks an accessible name (aria-labelledby or aria-label)."
		hint = "Add aria-labelledby pointing to the heading element ID, or add an explicit aria-label."
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  msg,
			Hint:     hint,
		},
	}
}

func isAriaModalDeclared(node *ir.Node) bool {
	modalVal, ok := node.GetAttr("aria-modal")
	if !ok {
		return false
	}
	clean := strings.ToLower(CleanAttr(modalVal))
	return clean == "true" || strings.Contains(modalVal, "{")
}

func hasAccessibleDialogName(node *ir.Node) bool {
	if label, ok := node.GetAttr("aria-label"); ok && CleanAttr(label) != "" {
		return true
	}
	if labelledby, ok := node.GetAttr("aria-labelledby"); ok && CleanAttr(labelledby) != "" {
		return true
	}
	if title, ok := node.GetAttr("title"); ok && CleanAttr(title) != "" {
		return true
	}

	return hasDialogTitleChild(node)
}

func hasDialogTitleChild(node *ir.Node) bool {
	for _, child := range node.Children {
		if child.Type == ir.NodeElement {
			tagLower := strings.ToLower(child.Tag)
			if strings.Contains(tagLower, "dialogtitle") || strings.Contains(tagLower, "modaltitle") {
				return true
			}
		}
	}
	return false
}
