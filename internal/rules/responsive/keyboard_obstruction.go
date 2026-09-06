package responsive

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// KeyboardObstructionRule mendeteksi penempatan aksi 'fixed bottom-0' di dalam form
// tanpa pembungkus scroll vertikal ('overflow-y-auto'), yang berisiko tertutup oleh virtual keyboard pada smartphone.
type KeyboardObstructionRule struct{}

// NewKeyboardObstructionRule membuat instance baru dari KeyboardObstructionRule.
func NewKeyboardObstructionRule() *KeyboardObstructionRule {
	return &KeyboardObstructionRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *KeyboardObstructionRule) ID() string {
	return "responsive.keyboard-obstruction"
}

// Description mengembalikan ringkasan aturan.
func (r *KeyboardObstructionRule) Description() string {
	return "Warns against fixed bottom action bars in forms lacking vertical scroll containers, which can be obstructed by the mobile virtual keyboard"
}

// Category mengembalikan nama kategori rule.
func (r *KeyboardObstructionRule) Category() string {
	return "responsive"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *KeyboardObstructionRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *KeyboardObstructionRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"WCAG 2.2 Guideline 2.1 (Keyboard Accessible)",
			"Material Design 3 Mobile Form Guidelines",
			"iOS Human Interface Guidelines (Managing the Virtual Keyboard)",
		},
		CoreInvariant: "Forms containing text inputs and fixed/sticky bottom action bars must provide a scrollable container ('overflow-y-auto') so inputs and actions are never obscured when the virtual keyboard expands.",
		Grounding: "When a user taps an input on a smartphone, the virtual software keyboard slides up from the bottom of the screen, consuming 40% to 50% of the visible viewport.\n\n" +
			"Elements styled with 'fixed bottom-0' or 'sticky bottom-0' remain pinned above the viewport bottom or above the keyboard. If the parent form is not wrapped in a vertical scroll container ('overflow-y-auto'), the active input field gets pushed behind the keyboard or under the fixed button, leaving the user unable to view their input or complete submission.\n\n" +
			"Charites enforces scrollable viewport resilience for mobile form layouts.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Mobile Virtual Keyboard Input Obstruction",
				Severity: "HIGH",
				Impact:   "Users cannot see text being typed into lower form inputs because fixed bottom bars pin directly over them.",
			},
			{
				Vector:   "Form Abandonment & Submission Blockers",
				Severity: "MEDIUM",
				Impact:   "When keyboard expansion pushes inputs offscreen without scroll capabilities, conversion rates drop.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Fixed bottom submit bar in a rigid form lacking scrollable region",
				Code: `<form className="h-screen flex flex-col justify-between">
  <div className="p-4 space-y-4">
    <input type="text" placeholder="Nama Lengkap" />
    <input type="email" placeholder="Alamat Surel" />
    <textarea placeholder="Pesan Anda" />
  </div>
  <div className="fixed bottom-0 inset-x-0 p-4 bg-surface border-t">
    <button type="submit" className="w-full bg-primary text-white py-3 rounded">Kirim</button>
  </div>
</form>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Scrollable body with fixed bottom bar allowing smooth mobile keyboard reflow",
				Code: `<form className="h-screen flex flex-col">
  <div className="flex-1 overflow-y-auto p-4 space-y-4">
    <input type="text" placeholder="Nama Lengkap" />
    <input type="email" placeholder="Alamat Surel" />
    <textarea placeholder="Pesan Anda" />
  </div>
  <div className="p-4 bg-surface border-t">
    <button type="submit" className="w-full bg-primary text-white py-3 rounded">Kirim</button>
  </div>
</form>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen fixed-bottom berada di dalam formulir tanpa kontainer scroll.
func (r *KeyboardObstructionRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !isFixedBottom(node.Classes) {
		return nil
	}

	formAncestor := findFormAncestor(node)
	if formAncestor == nil {
		return nil
	}

	if !hasFormInputDescendant(formAncestor) {
		return nil
	}

	if hasScrollableRegionAncestor(node) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Fixed/sticky bottom bar <%s> inside <%s> lacks a scrollable container ('overflow-y-auto'). On mobile devices, the virtual keyboard will obscure inputs or push submit actions offscreen.", node.Tag, formAncestor.Tag),
			Hint:     "Wrap form fields in a scrollable container ('flex-1 overflow-y-auto') or use standard static layout flow instead of 'fixed bottom-0'.",
		},
	}
}

func findFormAncestor(node *ir.Node) *ir.Node {
	curr := node.Parent
	for curr != nil {
		if curr.Tag == "form" {
			return curr
		}
		curr = curr.Parent
	}
	return nil
}
