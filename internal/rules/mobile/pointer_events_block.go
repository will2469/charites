package mobile

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// PointerEventsBlockRule mendeteksi elemen interaktif di bawah kontainer berkelas pointer-events-none
// tanpa pemulihan eksplisit pointer-events-auto pada perangkat sentuh mobile.
type PointerEventsBlockRule struct{}

// NewPointerEventsBlockRule membuat instance baru dari PointerEventsBlockRule.
func NewPointerEventsBlockRule() *PointerEventsBlockRule {
	return &PointerEventsBlockRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *PointerEventsBlockRule) ID() string {
	return "mobile.pointer-events-block"
}

// Description mengembalikan ringkasan aturan.
func (r *PointerEventsBlockRule) Description() string {
	return "Warns when an ancestor declares pointer-events-none over interactive children without restoring pointer-events-auto on mobile"
}

// Category mengembalikan nama kategori rule.
func (r *PointerEventsBlockRule) Category() string {
	return "mobile"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *PointerEventsBlockRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *PointerEventsBlockRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Pointer Events Level 3 (Pointer Event Processing Model)",
			"CSS Basic User Interface Module Level 4 (The pointer-events Property)",
			"Chromium Touch Action & Pointer Hierarchy Engine",
		},
		CoreInvariant: "Interactive descendants (<button>, <a>, <input>) nested under a 'pointer-events-none' ancestor must explicitly declare 'pointer-events-auto' so mobile touch taps are dispatched.",
		Grounding: "Applying CSS 'pointer-events-none' to an ancestor wrapper disables hit-testing for the element and all its children.\n\n" +
			"When developers nest interactive controls (<button>, <a>, <input>) inside such wrappers (often used for visual backdrop filters or transition overlays) without restoring 'pointer-events-auto', touchscreen taps and mouse clicks are completely ignored by the browser.\n\n" +
			"Restoring 'pointer-events-auto' directly on the interactive control re-enables event capture while preserving the pass-through behavior of the parent.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Completely Unresponsive Touch Buttons",
				Severity: "MEDIUM",
				Impact:   "Users tap buttons or links repeatedly with zero visual feedback or event dispatch on mobile browsers.",
			},
			{
				Vector:   "Silently Broken Form Submissions",
				Severity: "MEDIUM",
				Impact:   "Submit controls become inactive, giving the illusion that the application is broken or frozen.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Interactive button blocked under pointer-events-none parent",
				Code: `<div className="pointer-events-none opacity-90 p-4">
  <button onClick={handleSave} className="bg-primary text-white px-4 py-2">
    Simpan Data
  </button>
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Explicit pointer-events-auto restores touch interactivity",
				Code: `<div className="pointer-events-none opacity-90 p-4">
  <button onClick={handleSave} className="pointer-events-auto bg-primary text-white px-4 py-2 rounded-xl">
    Simpan Data
  </button>
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen interaktif terblokir oleh ancestor ber-pointer-events-none.
func (r *PointerEventsBlockRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !isInteractiveElement(node) {
		return nil
	}

	if hasPointerEventsAuto(node.Classes, node.RawClasses) {
		return nil
	}

	if isExplicitlyDisabled(node) {
		return nil
	}

	if !isBlockedByAncestorPointerEvents(node) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Interactive element (<" + node.Tag + ">) is nested under a 'pointer-events-none' ancestor without restoring 'pointer-events-auto'. Touch taps and clicks will be ignored.",
			Hint:     "Add 'pointer-events-auto' to the interactive element to restore tap and click event dispatch on mobile screens.",
		},
	}
}

func hasPointerEventsAuto(classes []string, rawClasses string) bool {
	if strings.Contains(rawClasses, "pointer-events-auto") {
		return true
	}
	for _, cls := range classes {
		if cls == "pointer-events-auto" {
			return true
		}
	}
	return false
}

func hasPointerEventsNone(classes []string, rawClasses string) bool {
	if strings.Contains(rawClasses, "pointer-events-none") {
		return true
	}
	for _, cls := range classes {
		if cls == "pointer-events-none" {
			return true
		}
	}
	return false
}

func isExplicitlyDisabled(node *ir.Node) bool {
	if node.Attributes == nil {
		return false
	}
	if _, ok := node.Attributes["disabled"]; ok {
		return true
	}
	if hidden, ok := node.Attributes["aria-hidden"]; ok && cleanAttrValue(hidden) == "true" {
		return true
	}
	return false
}

func isBlockedByAncestorPointerEvents(node *ir.Node) bool {
	for p := node.Parent; p != nil; p = p.Parent {
		if p.Type != ir.NodeElement {
			continue
		}
		if hasPointerEventsAuto(p.Classes, p.RawClasses) {
			return false
		}
		if hasPointerEventsNone(p.Classes, p.RawClasses) {
			return true
		}
	}
	return false
}
