package mobile

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// KeyboardViewportRiskRule mendeteksi kontainer layout berukuran viewport kaku (h-screen)
// yang memuat input teks dan kontrol bawah tetap tanpa unit viewport dinamis (dvh/svh).
type KeyboardViewportRiskRule struct{}

// NewKeyboardViewportRiskRule membuat instance baru dari KeyboardViewportRiskRule.
func NewKeyboardViewportRiskRule() *KeyboardViewportRiskRule {
	return &KeyboardViewportRiskRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *KeyboardViewportRiskRule) ID() string {
	return "mobile.keyboard-viewport-risk"
}

// Description mengembalikan ringkasan aturan.
func (r *KeyboardViewportRiskRule) Description() string {
	return "Advises using dynamic viewport units (dvh/svh) on containers with inputs and fixed controls to prevent layout breaking when virtual keyboard appears"
}

// Category mengembalikan nama kategori rule.
func (r *KeyboardViewportRiskRule) Category() string {
	return "mobile"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (info).
func (r *KeyboardViewportRiskRule) DefaultSeverity() ir.Severity {
	return ir.SeverityInfo
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *KeyboardViewportRiskRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C CSS Values and Units Module Level 4 Section 6.1.2 (Small, Large, and Dynamic Viewport Units)",
			"Chromium Virtual Keyboard API & Resize Invariants",
			"Apple WebKit Form Viewport Resilience Guidelines",
		},
		CoreInvariant: "Containers enclosing active text inputs alongside bottom-pinned actions must use dynamic viewport units ('min-h-dvh', 'svh') or sticky positioning instead of rigid 'h-screen' to prevent viewport clipping when virtual keyboard opens.",
		Grounding: "When a virtual keyboard appears on smartphone touchscreens, it consumes 40% to 50% of the display height, shrinking the browser visual viewport.\n\n" +
			"Containers locked to 'h-screen' or 'h-[100vh]' do not adjust dynamically to the reduced visual viewport, causing fixed bottom action buttons or active input fields to be pushed behind the keyboard or clipped.\n\n" +
			"Adopting dynamic viewport units (such as 'min-h-dvh' or 'min-h-svh') and sticky bottom positioning guarantees smooth, scrollable adaptation across Android and iOS keyboards.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Hidden Input Fields Behind Virtual Keyboard",
				Severity: "LOW",
				Impact:   "Mobile users cannot see what they are typing because inputs remain trapped behind the active keyboard.",
			},
			{
				Vector:   "Inaccessible Fixed Bottom Submit Button",
				Severity: "LOW",
				Impact:   "Fixed bottom submit buttons get pushed below the visible viewport, preventing form completion.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Rigid h-screen container with input and fixed bottom button",
				Code: `<div className="fixed inset-0 h-screen flex flex-col justify-between">
  <input type="text" placeholder="Nama Lengkap" />
  <button className="fixed bottom-0 w-full py-3 bg-primary text-white">Simpan</button>
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Dynamic viewport height unit with sticky bottom button",
				Code: `<div className="min-h-dvh flex flex-col justify-between pb-[env(safe-area-inset-bottom)]">
  <input type="text" placeholder="Nama Lengkap" />
  <button className="sticky bottom-4 w-full py-3 bg-primary text-white rounded-xl">Simpan</button>
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah kontainer memiliki h-screen kaku dengan input dan fixed bottom control.
func (r *KeyboardViewportRiskRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !hasRigidViewportHeight(node.Classes, node.RawClasses) {
		return nil
	}

	if hasDynamicViewportUnit(node.Classes, node.RawClasses) {
		return nil
	}

	if isDesktopOnly(node) {
		return nil
	}

	hasInput, hasBottomCtrl := scanInputsAndBottomControls(node)
	if !hasInput || !hasBottomCtrl {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Container uses rigid viewport height ('h-screen' / '100vh') with active inputs and bottom controls. When the virtual keyboard appears, the visual viewport will shrink and clip controls.",
			Hint:     "Replace 'h-screen' with dynamic viewport units (e.g. 'min-h-dvh' or 'min-h-svh') and use sticky bottom positioning to accommodate virtual keyboards.",
		},
	}
}

func hasRigidViewportHeight(classes []string, rawClasses string) bool {
	if strings.Contains(rawClasses, "h-screen") ||
		strings.Contains(rawClasses, "h-[100vh]") ||
		strings.Contains(rawClasses, "min-h-screen") {
		return true
	}
	for _, cls := range classes {
		if cls == "h-screen" || cls == "h-[100vh]" || cls == "min-h-screen" {
			return true
		}
	}
	return false
}

func hasDynamicViewportUnit(classes []string, rawClasses string) bool {
	if strings.Contains(rawClasses, "dvh") || strings.Contains(rawClasses, "svh") {
		return true
	}
	for _, cls := range classes {
		if strings.Contains(cls, "dvh") || strings.Contains(cls, "svh") {
			return true
		}
	}
	return false
}

func scanInputsAndBottomControls(container *ir.Node) (hasInput bool, hasBottomCtrl bool) {
	for curr := range container.Walk() {
		if curr == container || curr.Type != ir.NodeElement {
			continue
		}

		if isTextInputElement(curr) {
			hasInput = true
		}

		if isBottomFixedOrAbsolute(curr.Classes, curr.RawClasses) {
			hasBottomCtrl = true
		}

		if hasInput && hasBottomCtrl {
			return true, true
		}
	}
	return hasInput, hasBottomCtrl
}

func isTextInputElement(node *ir.Node) bool {
	tagLower := strings.ToLower(node.Tag)
	if tagLower == "textarea" {
		return true
	}
	if tagLower == "input" {
		if node.Attributes == nil {
			return true
		}
		if t, ok := node.Attributes["type"]; ok {
			cleanType := cleanAttrValue(t)
			switch cleanType {
			case "hidden", "button", "submit", "reset", "checkbox", "radio":
				return false
			default:
				return true
			}
		}
		return true
	}
	return false
}

func isBottomFixedOrAbsolute(classes []string, rawClasses string) bool {
	isFixedOrAbs := strings.Contains(rawClasses, "fixed") || strings.Contains(rawClasses, "absolute")
	hasBottom := strings.Contains(rawClasses, "bottom-0") || strings.Contains(rawClasses, "bottom-")

	if isFixedOrAbs && hasBottom {
		return true
	}

	var hasPos, hasBot bool
	for _, cls := range classes {
		if cls == "fixed" || cls == "absolute" {
			hasPos = true
		}
		if cls == "bottom-0" || strings.HasPrefix(cls, "bottom-") {
			hasBot = true
		}
	}
	return hasPos && hasBot
}
