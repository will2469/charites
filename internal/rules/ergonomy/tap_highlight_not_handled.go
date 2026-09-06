package ergonomy

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// TapHighlightNotHandledRule mendeteksi elemen kustom interaktif non-native tanpa feedback sentuh active atau tap-highlight CSS.
type TapHighlightNotHandledRule struct{}

// NewTapHighlightNotHandledRule membuat instance baru dari TapHighlightNotHandledRule.
func NewTapHighlightNotHandledRule() *TapHighlightNotHandledRule {
	return &TapHighlightNotHandledRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *TapHighlightNotHandledRule) ID() string {
	return "ergonomy.tap-highlight-not-handled"
}

// Description mengembalikan ringkasan aturan.
func (r *TapHighlightNotHandledRule) Description() string {
	return "Flags clickable non-native custom elements lacking tactile tap feedback or tap-highlight management"
}

// Category mengembalikan nama kategori rule.
func (r *TapHighlightNotHandledRule) Category() string {
	return "ergonomy"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (info).
func (r *TapHighlightNotHandledRule) DefaultSeverity() ir.Severity {
	return ir.SeverityInfo
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *TapHighlightNotHandledRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Touch Events Community Group Guidelines",
			"Chromium Android Tap Feedback UX Standards",
			"Google Material Design (Tactile States & Surface Elevation)",
		},
		CoreInvariant: "Non-native clickable elements (<div onClick>, <span role=\"button\">) must declare deliberate active feedback or suppress the default Android Chrome grey tap highlight box.",
		Grounding: "On Chromium Android, tapping an element without a native button role causes the browser to flash a rigid semi-transparent grey overlay box.\n\n" +
			"Without deliberate 'active:' micro-interactions (such as 'active:scale-[0.99]' or 'active:bg-muted') or setting '[-webkit-tap-highlight-color:transparent]', " +
			"the application exhibits noticeable visual glitches and lacks native tactile responsiveness.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Visual Glitches on Android Chrome",
				Severity: "LOW",
				Impact:   "Rigid grey highlight rectangles flash abruptly over custom cards, badges, and list rows during touch.",
			},
			{
				Vector:   "Poor Tactile Feedback",
				Severity: "LOW",
				Impact:   "Users cannot perceive if a touch tap registered, leading to repeated tapping and accidental duplicate submissions.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Non-native clickable div without active feedback",
				Code: `<div
  role="button"
  tabIndex={0}
  onClick={handleSelectCard}
  className="p-4 bg-card border rounded-2xl"
>
  <span>Pilihan Layanan</span>
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Deliberate active feedback with suppressed grey highlight",
				Code: `<div
  role="button"
  tabIndex={0}
  onClick={handleSelectCard}
  className="p-4 bg-card border rounded-2xl active:bg-muted/60 active:scale-[0.99] transition-transform [-webkit-tap-highlight-color:transparent]"
>
  <span>Pilihan Layanan</span>
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen non-native interaktif menangani feedback tap highlight.
func (r *TapHighlightNotHandledRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || node.Attributes == nil {
		return nil
	}

	tag := strings.ToLower(node.Tag)
	if isNativeInteractiveTag(tag) {
		return nil
	}

	if !hasInteractiveEventOrRole(node.Attributes) {
		return nil
	}

	if hasActiveFeedbackOrTapHighlight(node.Classes, node.RawClasses) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Non-native interactive element has click/touch handler without active feedback or tap-highlight styling. Android Chrome will show a rigid grey tap highlight box.",
			Hint:     "Add tactile feedback (e.g. 'active:bg-muted/60 active:scale-[0.99]') and consider '[-webkit-tap-highlight-color:transparent]'.",
		},
	}
}

func isNativeInteractiveTag(tag string) bool {
	switch tag {
	case "button", "a", "input", "select", "textarea", "summary", "option", "label":
		return true
	default:
		return false
	}
}

func hasInteractiveEventOrRole(attrs map[string]string) bool {
	for k, v := range attrs {
		lowerKey := strings.ToLower(k)
		if lowerKey == "onclick" || lowerKey == "ontouchstart" || lowerKey == "onpointerdown" {
			return true
		}
		if lowerKey == "role" && strings.ToLower(v) == "button" {
			return true
		}
	}
	return false
}

func hasActiveFeedbackOrTapHighlight(classes []string, rawClasses string) bool {
	if strings.Contains(rawClasses, "active:") ||
		strings.Contains(rawClasses, "tap-highlight") ||
		strings.Contains(rawClasses, "-webkit-tap-highlight-color") {
		return true
	}

	for _, cls := range classes {
		if strings.HasPrefix(cls, "active:") ||
			strings.Contains(cls, "tap-highlight") ||
			strings.Contains(cls, "-webkit-tap-highlight-color") {
			return true
		}
	}

	return false
}
