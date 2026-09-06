package mobile

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// ModalViewportLockRule mendeteksi kontainer modal atau dialog yang mengunci viewport dengan overflow-hidden
// tanpa menyediakan region scroll vertikal internal (overflow-y-auto) pada perangkat mobile.
type ModalViewportLockRule struct{}

// NewModalViewportLockRule membuat instance baru dari ModalViewportLockRule.
func NewModalViewportLockRule() *ModalViewportLockRule {
	return &ModalViewportLockRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *ModalViewportLockRule) ID() string {
	return "mobile.modal-viewport-lock"
}

// Description mengembalikan ringkasan aturan.
func (r *ModalViewportLockRule) Description() string {
	return "Detects modal dialog containers locked with overflow-hidden without an internal scrollable region on mobile viewports"
}

// Category mengembalikan nama kategori rule.
func (r *ModalViewportLockRule) Category() string {
	return "mobile"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *ModalViewportLockRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ModalViewportLockRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C ARIA Authoring Practices Guide (Modal Dialog Design Pattern)",
			"WCAG 2.2 Success Criterion 2.1.2 (No Keyboard Trap)",
			"Apple Human Interface Guidelines (Modals and Sheets on Mobile)",
		},
		CoreInvariant: "Modal dialog containers declaring 'overflow-hidden' must provide an internal scrollable region ('overflow-y-auto') so content remains accessible on short mobile screens.",
		Grounding: "Full-screen modal dialogs or bottom sheets often lock body scrolling with 'overflow-hidden'.\n\n" +
			"If the modal container itself lacks an internal vertical scrollable container ('overflow-y-auto' or 'overflow-y-scroll'), content that exceeds the screen height (such as on short smartphones, landscape orientation, or when the virtual keyboard opens) is permanently cropped.\n\n" +
			"Users cannot scroll to reach bottom form inputs or confirm/cancel action buttons, resulting in a critical UX failure.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Unreachable Submit & Dismiss Actions",
				Severity: "HIGH",
				Impact:   "Users are locked in the modal with no ability to reach submission or close buttons on smaller mobile screens.",
			},
			{
				Vector:   "Form Inaccessibility on Keyboard Activation",
				Severity: "HIGH",
				Impact:   "When virtual keyboard expands, form fields below the keyboard cannot be scrolled into view.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Modal dialog container locked with overflow-hidden without scrollable region",
				Code: `<div role="dialog" aria-modal="true" className="fixed inset-0 overflow-hidden flex items-center justify-center p-4">
  <div className="bg-card p-6 rounded-2xl w-full max-w-md h-screen">
    <h2>Form Permohonan Bantuan</h2>
    <div className="space-y-4">...isi form panjang...</div>
    <button type="submit">Kirim</button>
  </div>
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Internal scroll region (overflow-y-auto) allows smooth scrolling on mobile screens",
				Code: `<div role="dialog" aria-modal="true" className="fixed inset-0 overflow-y-auto flex items-center justify-center p-4">
  <div className="bg-card p-6 rounded-2xl w-full max-w-md my-auto">
    <h2>Form Permohonan Bantuan</h2>
    <div className="space-y-4">...isi form panjang...</div>
    <button type="submit">Kirim</button>
  </div>
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah kontainer modal mengunci overflow tanpa region scroll internal.
func (r *ModalViewportLockRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !isModalContainer(node) {
		return nil
	}

	if !hasOverflowHidden(node.Classes, node.RawClasses) {
		return nil
	}

	if isDesktopOnly(node) {
		return nil
	}

	if hasInternalScrollRegion(node) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Modal dialog container locks viewport with 'overflow-hidden' without an internal scrollable region ('overflow-y-auto'). On short mobile screens, content and actions will be cropped.",
			Hint:     "Add 'overflow-y-auto' to the modal wrapper or an inner dialog body container to allow mobile scrolling.",
		},
	}
}

func isModalContainer(node *ir.Node) bool {
	if node.Attributes != nil {
		if role, ok := node.Attributes["role"]; ok {
			cleanRole := cleanAttrValue(role)
			if cleanRole == "dialog" || cleanRole == "alertdialog" {
				return true
			}
		}
	}

	tagLower := strings.ToLower(node.Tag)
	if tagLower == "dialog" || tagLower == "modal" || strings.Contains(tagLower, "dialog") {
		return true
	}

	// Deteksi kelas fixed inset-0 full screen overlay
	hasFixed := strings.Contains(node.RawClasses, "fixed")
	hasInset0 := strings.Contains(node.RawClasses, "inset-0")
	if hasFixed && hasInset0 {
		return true
	}

	for _, cls := range node.Classes {
		if cls == "inset-0" && hasFixed {
			return true
		}
	}

	return false
}

func hasOverflowHidden(classes []string, rawClasses string) bool {
	if strings.Contains(rawClasses, "overflow-hidden") {
		return true
	}
	for _, cls := range classes {
		if cls == "overflow-hidden" {
			return true
		}
	}
	return false
}

func hasInternalScrollRegion(container *ir.Node) bool {
	for curr := range container.Walk() {
		if curr == container || curr.Type != ir.NodeElement {
			continue
		}
		if isScrollableClass(curr.Classes, curr.RawClasses) {
			return true
		}
	}
	return false
}

func isScrollableClass(classes []string, rawClasses string) bool {
	if strings.Contains(rawClasses, "overflow-y-auto") ||
		strings.Contains(rawClasses, "overflow-y-scroll") ||
		strings.Contains(rawClasses, "overflow-auto") {
		return true
	}
	for _, cls := range classes {
		if cls == "overflow-y-auto" || cls == "overflow-y-scroll" || cls == "overflow-auto" {
			return true
		}
	}
	return false
}
