package browser

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// AppearanceNativeOverrideRule mendeteksi elemen form kontrol native (select, checkbox, radio, range, date)
// yang diberi styling kustom (border, background, rounded) tanpa menyertakan utility 'appearance-none'
// untuk mereset tampilan native WebKit/Safari.
type AppearanceNativeOverrideRule struct{}

// NewAppearanceNativeOverrideRule membuat instance baru AppearanceNativeOverrideRule.
func NewAppearanceNativeOverrideRule() *AppearanceNativeOverrideRule {
	return &AppearanceNativeOverrideRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat browser.appearance-native-override.
func (r *AppearanceNativeOverrideRule) ID() string {
	return "browser.appearance-native-override"
}

// Description mengembalikan deskripsi ringkas aturan.
func (r *AppearanceNativeOverrideRule) Description() string {
	return "Enforces explicit appearance-none on form controls with custom styling to prevent WebKit/Safari native UI clashes"
}

// Category mengembalikan nama kategori rule.
func (r *AppearanceNativeOverrideRule) Category() string {
	return "browser"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *AppearanceNativeOverrideRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *AppearanceNativeOverrideRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C CSS Basic User Interface Module Level 4 (appearance: none)",
			"HTML Living Standard Section 4.10.5 (Form Controls & Native Rendering)",
			"WebKit Form Control Styling Compatibility Guidelines",
		},
		CoreInvariant: "Native form controls (<select>, <input type=\"checkbox|radio|range|date|time|datetime-local\">) with custom styling classes must explicitly declare 'appearance-none' to prevent WebKit/Safari OS widget collisions.",
		Grounding: "Blink (Chrome/Edge) and Gecko (Firefox) automatically strip most native platform widget decorations when custom border, background, or border-radius properties are defined on form controls.\n\n" +
			"In contrast, WebKit (Safari macOS and iOS) retains native glossy gradients, 3D rounded bezels, and OS-level indicator graphics unless 'appearance: none' (-webkit-appearance: none) is explicitly specified.\n\n" +
			"When developers test only in desktop Chrome, a custom-styled <select> appears sleek and modern. However, on iOS Safari, the custom border and background clash with native pickers and glossy overlays.\n\n" +
			"Charites enforces explicit 'appearance-none' on all custom-styled native form controls, ensuring visual cross-engine parity.",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Custom styled select without appearance-none (causes glossy bezel clash on iOS Safari)",
				Code: `<select className="h-11 px-3.5 py-2.5 bg-background border border-input rounded-xl text-sm font-medium">
  <option value="1">Layanan Surat</option>
  <option value="2">Layanan Kependudukan</option>
</select>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Select with appearance-none resetting native WebKit styling cleanly",
				Code: `<select className="appearance-none h-11 px-3.5 py-2.5 bg-background border border-input rounded-xl text-sm font-medium">
  <option value="1">Layanan Surat</option>
  <option value="2">Layanan Kependudukan</option>
</select>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "WebKit Bezel Collision",
				Severity: "MEDIUM",
				Impact:   "Severe visual inconsistency on Safari iOS where native OS glossy gradients render on top of Tailwind theme styling.",
			},
			{
				Vector:   "Dropdown Arrow Misalignment",
				Severity: "LOW",
				Impact:   "Unaligned custom dropdown arrows and clipped options inside custom-sized containers.",
			},
		},
	}
}

// Evaluate memeriksa apakah kontrol form native memiliki styling kustom tanpa appearance-none.
func (r *AppearanceNativeOverrideRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	tag := node.Tag
	if len(tag) == 0 || (tag[0] >= 'A' && tag[0] <= 'Z') {
		return nil
	}
	tag = strings.ToLower(tag)
	if !isTargetFormControl(tag, node.Attributes) {
		return nil
	}

	// Cek apakah sudah memiliki appearance-none atau reset appearance
	if hasAppearanceReset(node) {
		return nil
	}

	// Cek apakah memiliki kelas styling kustom yang memicu kebutuhan reset WebKit
	if !hasCustomStyling(node) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Form control '<" + tag + ">' has custom styling but lacks 'appearance-none'. WebKit/Safari will retain native OS glossy bezels and widgets unless explicitly reset.",
			Hint:     "Add 'appearance-none' to the className list to ensure cross-engine rendering parity on Safari iOS and macOS.",
		},
	}
}

func isTargetFormControl(tag string, attrs map[string]string) bool {
	if tag == "select" || tag == "progress" || tag == "meter" {
		return true
	}

	if tag == "input" {
		inputType := strings.ToLower(attrs["type"])
		switch inputType {
		case "checkbox", "radio", "range", "date", "time", "datetime-local":
			return true
		default:
			return false
		}
	}

	return false
}

func hasAppearanceReset(node *ir.Node) bool {
	for _, cls := range node.Classes {
		if cls == "appearance-none" || cls == "appearance-auto" {
			return true
		}
	}

	if rawClass, ok := node.Attributes["class"]; ok {
		if strings.Contains(rawClass, "appearance-none") {
			return true
		}
	}

	if rawClass, ok := node.Attributes["className"]; ok {
		if strings.Contains(rawClass, "appearance-none") {
			return true
		}
	}

	if style, ok := node.Attributes["style"]; ok {
		if strings.Contains(style, "appearance") {
			return true
		}
	}

	return false
}

func hasCustomStyling(node *ir.Node) bool {
	for _, cls := range node.Classes {
		if isCustomStyleClass(cls) {
			return true
		}
	}

	if rawClass, ok := node.Attributes["class"]; ok {
		fields := strings.Fields(rawClass)
		for _, f := range fields {
			if isCustomStyleClass(f) {
				return true
			}
		}
	}

	if rawClass, ok := node.Attributes["className"]; ok {
		fields := strings.Fields(rawClass)
		for _, f := range fields {
			if isCustomStyleClass(f) {
				return true
			}
		}
	}

	if _, ok := node.Attributes["style"]; ok {
		return true
	}

	return false
}

func isCustomStyleClass(cls string) bool {
	if cls == "" || cls == "block" || cls == "inline-block" || cls == "inline" || cls == "flex" {
		return false
	}

	base := cls
	if idx := strings.LastIndex(cls, ":"); idx != -1 {
		base = cls[idx+1:]
	}

	return strings.HasPrefix(base, "border") ||
		strings.HasPrefix(base, "bg-") ||
		strings.HasPrefix(base, "rounded") ||
		strings.HasPrefix(base, "shadow") ||
		strings.HasPrefix(base, "p-") ||
		strings.HasPrefix(base, "px-") ||
		strings.HasPrefix(base, "py-") ||
		strings.HasPrefix(base, "h-") ||
		strings.HasPrefix(base, "w-") ||
		strings.HasPrefix(base, "ring-")
}
