package a11y

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// KeyboardTrapMissingEscapeRule memastikan custom modal dialog menyediakan
// mekanisme penutupan berbasis keyboard (Escape key listener atau tombol close yang aksesibel).
type KeyboardTrapMissingEscapeRule struct{}

// NewKeyboardTrapMissingEscapeRule membuat instance baru KeyboardTrapMissingEscapeRule.
func NewKeyboardTrapMissingEscapeRule() *KeyboardTrapMissingEscapeRule {
	return &KeyboardTrapMissingEscapeRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat a11y.keyboard-trap-missing-escape.
func (r *KeyboardTrapMissingEscapeRule) ID() string {
	return "a11y.keyboard-trap-missing-escape"
}

// Description mengembalikan deskripsi ringkas aturan.
func (r *KeyboardTrapMissingEscapeRule) Description() string {
	return "Enforces that custom modal dialogs provide an Escape key listener or an accessible dismiss mechanism"
}

// Category mengembalikan nama kategori rule.
func (r *KeyboardTrapMissingEscapeRule) Category() string {
	return "a11y"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *KeyboardTrapMissingEscapeRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *KeyboardTrapMissingEscapeRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Web Content Accessibility Guidelines (WCAG) 2.2 SC 2.1.2 (No Keyboard Trap)",
			"W3C WAI-ARIA Authoring Practices Guide (APG) Modal Dialog Pattern",
		},
		CoreInvariant: "Custom modal dialogs (role=\"dialog\" / role=\"alertdialog\") must provide keyboard dismissibility (onKeyDown Escape listener) or an accessible close button.",
		Grounding: "When web applications render custom modal overlays using <div> containers instead of accessible headless dialog primitives:\n" +
			"1. Keyboard Trap: Users navigating solely via keyboard or switch devices can become trapped inside the modal overlay with no mechanism to exit.\n" +
			"2. Violation of WCAG 2.1.2: Users must always be able to untrap themselves by pressing the Escape key without navigating through complex hierarchies.\n" +
			"3. Screen Reader Disorientation: If an overlay cannot be dismissed via standard keyboard conventions, assistive technology users are forced to refresh the application.\n\n" +
			"Charites inspects custom dialog containers to ensure they wire onKeyDown handlers or expose accessible close triggers.",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Custom dialog overlay without onKeyDown Escape listener or close button",
				Code: `<div role="dialog" aria-modal="true" aria-labelledby="modal-title" className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
  <div className="bg-background p-6 rounded-xl">
    <h2 id="modal-title" className="text-lg font-semibold">Konfirmasi Hapus</h2>
    <p className="text-sm text-muted-foreground mt-2">Apakah Anda yakin ingin menghapus data ini?</p>
  </div>
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Compliant custom dialog with onKeyDown Escape listener and accessible close button",
				Code: `<div role="dialog" aria-modal="true" aria-labelledby="modal-title" onKeyDown={handleKeyDown} className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
  <div className="bg-background p-6 rounded-xl relative">
    <button type="button" onClick={onClose} aria-label="Tutup dialog" className="size-11 absolute top-2 right-2 flex items-center justify-center">
      <XIcon className="size-5" />
    </button>
    <h2 id="modal-title" className="text-lg font-semibold">Konfirmasi Hapus</h2>
    <p className="text-sm text-muted-foreground mt-2">Apakah Anda yakin ingin menghapus data ini?</p>
  </div>
</div>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Permanent Keyboard Trap",
				Severity: "CRITICAL",
				Impact:   "Keyboard and switch users become unable to escape the dialog overlay (WCAG SC 2.1.2 Level A failure).",
			},
			{
				Vector:   "Loss of Focus State",
				Severity: "HIGH",
				Impact:   "Users are forced to refresh the browser tab, abandoning in-flight data entry.",
			},
		},
	}
}

// Evaluate memeriksa apakah custom modal dialog menyediakan penanganan keyboard atau tombol close.
func (r *KeyboardTrapMissingEscapeRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !isCustomDialogElement(node) {
		return nil
	}

	if hasSpreadProps(node.Attributes) || isDecorativeOrHidden(node.Attributes) {
		return nil
	}

	if hasKeyboardListener(node) || hasAccessibleCloseTrigger(node) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Custom modal dialog lacks an Escape key listener (onKeyDown) or accessible close trigger, causing a keyboard trap hazard.",
			Hint:     "Add an onKeyDown listener to handle the Escape key and ensure an accessible close button (aria-label=\"Close\") is available.",
		},
	}
}

func isCustomDialogElement(node *ir.Node) bool {
	if role, ok := node.GetAttr("role"); ok {
		cleanRole := strings.ToLower(CleanAttr(role))
		if cleanRole == "dialog" || cleanRole == "alertdialog" {
			return true
		}
	}

	tagLower := strings.ToLower(node.Tag)
	return tagLower == "modal" || tagLower == "customdialog"
}

func hasKeyboardListener(node *ir.Node) bool {
	if node.Attributes != nil {
		for k := range node.Attributes {
			kLower := strings.ToLower(k)
			if kLower == "onkeydown" || kLower == "onkeyup" || kLower == "onkeypress" {
				return true
			}
		}
	}

	for _, child := range node.Children {
		if hasKeyboardListener(child) {
			return true
		}
	}

	return false
}

func hasAccessibleCloseTrigger(node *ir.Node) bool {
	for n := range node.Walk() {
		if n == node || n.Type != ir.NodeElement {
			continue
		}

		if isCloseButtonNode(n) {
			return true
		}
	}
	return false
}

func isCloseButtonNode(n *ir.Node) bool {
	tagLower := strings.ToLower(n.Tag)
	if tagLower != "button" && tagLower != "a" && !strings.Contains(tagLower, "close") {
		if role, ok := n.GetAttr("role"); !ok || !strings.EqualFold(CleanAttr(role), "button") {
			return false
		}
	}

	if label, ok := n.GetAttr("aria-label"); ok && containsCloseIntent(label) {
		return true
	}
	if title, ok := n.GetAttr("title"); ok && containsCloseIntent(title) {
		return true
	}

	return hasCloseTextChild(n)
}

func containsCloseIntent(val string) bool {
	lower := strings.ToLower(CleanAttr(val))
	return strings.Contains(lower, "close") ||
		strings.Contains(lower, "tutup") ||
		strings.Contains(lower, "batal") ||
		strings.Contains(lower, "cancel") ||
		strings.Contains(lower, "dismiss") ||
		strings.Contains(lower, "kembali")
}

func hasCloseTextChild(n *ir.Node) bool {
	for _, child := range n.Children {
		if child.Type == ir.NodeText {
			text := strings.TrimSpace(child.RawClasses)
			if text == "" || text == "×" || text == "X" || text == "x" || containsCloseIntent(text) {
				return true
			}
		}
	}
	return false
}
