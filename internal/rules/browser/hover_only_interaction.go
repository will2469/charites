package browser

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// HoverOnlyInteractionRule memastikan elemen interaktif atau aksi penting yang disembunyikan
// dan dimunculkan lewat hover memiliki alternatif fokus keyboard dan sentuhan (focus-visible, group-focus-within).
type HoverOnlyInteractionRule struct{}

// NewHoverOnlyInteractionRule membuat instance baru HoverOnlyInteractionRule.
func NewHoverOnlyInteractionRule() *HoverOnlyInteractionRule {
	return &HoverOnlyInteractionRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat browser.hover-only-interaction.
func (r *HoverOnlyInteractionRule) ID() string {
	return "browser.hover-only-interaction"
}

// Description mengembalikan deskripsi ringkas aturan.
func (r *HoverOnlyInteractionRule) Description() string {
	return "Ensures interactive actions and state reveals have keyboard and touch counterparts instead of relying solely on hover"
}

// Category mengembalikan nama kategori rule.
func (r *HoverOnlyInteractionRule) Category() string {
	return "browser"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *HoverOnlyInteractionRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *HoverOnlyInteractionRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Web Content Accessibility Guidelines (WCAG) 2.2 SC 2.1.1 (Keyboard)",
			"WICG / WHATWG Touch Events & Pointer Events Level 3 (Touch vs Hover Ergonomics)",
			"Apple Human Interface Guidelines for iOS Touch Interactions",
		},
		CoreInvariant: "Interactive controls and revealed elements must not rely exclusively on ':hover' or 'group-hover:' without keyboard/touch counterparts ('focus-visible:', 'group-focus-within:').",
		Grounding: "Touchscreen devices (the majority of web traffic on Safari iOS and Chrome Android) have no physical cursor and cannot perform genuine mouse hover.\n\n" +
			"When critical action buttons (e.g. delete, edit, copy) are hidden by default with 'opacity-0' and only revealed via 'group-hover:opacity-100':\n" +
			"1. Total Mobile Inaccessibility: Touchscreen users cannot see or activate the controls because hovering does not exist on mobile.\n" +
			"2. iOS Sticky Hover Bug: Tapping an element on Safari iOS triggers an inconsistent 'sticky hover' state, requiring multiple confusing taps.\n" +
			"3. Keyboard Navigation Failure: Users navigating with the 'Tab' key cannot discover or focus hidden controls unless focus-within or focus-visible is bound.\n\n" +
			"Charites enforces that any hover-revealed element provides accessible keyboard and touch parity.",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Delete button hidden by default and only revealed on group-hover (invisible on touchscreens)",
				Code: `<div className="group flex items-center justify-between p-3 border rounded-xl">
  <span>Berkas_KTP.pdf</span>
  <button
    type="button"
    onClick={handleDelete}
    className="opacity-0 group-hover:opacity-100 text-destructive text-sm"
  >
    Hapus
  </button>
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Button accessible via hover, keyboard Tab navigation, and touch focus",
				Code: `<div className="group flex items-center justify-between p-3 border rounded-xl">
  <span>Berkas_KTP.pdf</span>
  <button
    type="button"
    onClick={handleDelete}
    className="opacity-0 group-hover:opacity-100 group-focus-within:opacity-100 focus-visible:opacity-100 text-destructive text-sm"
  >
    Hapus
  </button>
</div>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Mobile Touch Exclusion",
				Severity: "HIGH",
				Impact:   "Critical action controls are completely invisible and unreachable on smartphones and tablets.",
			},
			{
				Vector:   "Keyboard Accessibility Barrier",
				Severity: "MEDIUM",
				Impact:   "Fails WCAG 2.2 Level A keyboard navigation audits when hidden controls cannot be focused with the Tab key.",
			},
		},
	}
}

// Evaluate memeriksa apakah ada pengungkapan konten berbasis hover murni tanpa pasangan fokus.
func (r *HoverOnlyInteractionRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil {
		return nil
	}

	classes := node.Classes
	if len(classes) == 0 && node.Attributes != nil {
		if raw, ok := node.Attributes["class"]; ok {
			classes = strings.Fields(raw)
		} else if raw, ok := node.Attributes["className"]; ok {
			classes = strings.Fields(raw)
		}
	}

	if len(classes) == 0 {
		return nil
	}

	hasDefaultHidden := false
	hasHoverReveal := false
	hasFocusReveal := false

	for _, cls := range classes {
		// 1. Deteksi apakah elemen disembunyikan secara default
		if cls == "opacity-0" || cls == "invisible" || cls == "hidden" || cls == "scale-0" {
			hasDefaultHidden = true
		}

		// 2. Deteksi apakah dimunculkan saat hover
		if isHoverRevealClass(cls) {
			hasHoverReveal = true
		}

		// 3. Deteksi apakah ada penanganan alternatif fokus keyboard / sentuh
		if isFocusRevealClass(cls) {
			hasFocusReveal = true
		}
	}

	// Pelanggaran terjadi jika elemen disembunyikan secara default, dimunculkan lewat hover,
	// tetapi TIDAK memiliki padanan fokus (focus-visible / group-focus-within).
	if hasDefaultHidden && hasHoverReveal && !hasFocusReveal {
		tag := strings.ToLower(node.Tag)
		severity := r.DefaultSeverity()

		// Jika bukan tombol/input/link langsung, turunkan ke warning
		if tag != "button" && tag != "a" && tag != "input" && node.Attributes["role"] != "button" {
			severity = ir.SeverityWarn
		}

		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: severity,
				Message:  "Interactive element is hidden by default and only revealed on hover. Touch devices (iOS Safari, Android Chrome) lack hover capabilities, rendering this control inaccessible.",
				Hint:     "Add 'group-focus-within:opacity-100' and 'focus-visible:opacity-100' (or keep controls visible permanently) to ensure mobile and keyboard accessibility.",
			},
		}
	}

	return nil
}

func isHoverRevealClass(cls string) bool {
	return strings.HasPrefix(cls, "group-hover:opacity-") ||
		strings.HasPrefix(cls, "hover:opacity-") ||
		cls == "group-hover:block" ||
		cls == "hover:block" ||
		cls == "group-hover:visible" ||
		cls == "hover:visible" ||
		cls == "group-hover:scale-100" ||
		cls == "hover:scale-100"
}

func isFocusRevealClass(cls string) bool {
	return strings.HasPrefix(cls, "group-focus-within:") ||
		strings.HasPrefix(cls, "focus-within:") ||
		strings.HasPrefix(cls, "focus-visible:") ||
		strings.HasPrefix(cls, "focus:") ||
		strings.HasPrefix(cls, "active:") ||
		strings.HasPrefix(cls, "group-focus:")
}
