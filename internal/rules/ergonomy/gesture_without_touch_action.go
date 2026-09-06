package ergonomy

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// GestureWithoutTouchActionRule mendeteksi event handler gesture kustom tanpa deklarasi CSS touch-action.
type GestureWithoutTouchActionRule struct{}

// NewGestureWithoutTouchActionRule membuat instance baru dari GestureWithoutTouchActionRule.
func NewGestureWithoutTouchActionRule() *GestureWithoutTouchActionRule {
	return &GestureWithoutTouchActionRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *GestureWithoutTouchActionRule) ID() string {
	return "ergonomy.gesture-without-touch-action"
}

// Description mengembalikan ringkasan aturan.
func (r *GestureWithoutTouchActionRule) Description() string {
	return "Enforces CSS touch-action declaration on elements with custom gesture swipe/drag event handlers"
}

// Category mengembalikan nama kategori rule.
func (r *GestureWithoutTouchActionRule) Category() string {
	return "ergonomy"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *GestureWithoutTouchActionRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *GestureWithoutTouchActionRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Pointer Events Level 3 Section 5.2.8 (The touch-action CSS Property)",
			"Chromium & WebKit Compositor Gesture Isolation Architecture",
			"Google Chrome Developers (Touch Action Best Practices)",
		},
		CoreInvariant: "Elements attaching custom swipe or drag listeners ('onTouchMove', 'onPointerMove') must declare explicit CSS 'touch-action' ('touch-pan-y', 'touch-none') to prevent gesture cancellation by native scrolling.",
		Grounding: "When users drag or swipe an element, the browser mobile compositor thread must determine whether to handle native scrolling or yield control to JavaScript.\n\n" +
			"Without explicit CSS 'touch-action' (e.g. 'touch-pan-y' for horizontal sliders or 'touch-none' for drawing canvases), " +
			"browser vertical scrolling immediately cancels the custom touch gesture mid-drag, causing abrupt freezing or unwanted page scrolling.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Abrupt Gesture Cancellation",
				Severity: "MEDIUM",
				Impact:   "Swipeable cards and carousels stutter or lock up mid-swipe when the mobile browser takes over scrolling.",
			},
			{
				Vector:   "Accidental Page Scrolling",
				Severity: "MEDIUM",
				Impact:   "Users attempting to pan a map or slider accidentally scroll the entire page off-screen.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Horizontal swipeable container without touch-action",
				Code: `<div
  onTouchStart={handleStart}
  onTouchMove={handleMove}
  className="flex overflow-x-auto gap-4 p-4"
>
  <div className="w-64 h-40 bg-card rounded-2xl shrink-0">Kartu 1</div>
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Explicit touch-pan-y coordinates compositor axis smoothly",
				Code: `<div
  onTouchStart={handleStart}
  onTouchMove={handleMove}
  className="flex overflow-x-auto gap-4 p-4 touch-pan-y"
>
  <div className="w-64 h-40 bg-card rounded-2xl shrink-0">Kartu 1</div>
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen dengan event listener gesture memiliki utilitas touch-action.
func (r *GestureWithoutTouchActionRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || node.Attributes == nil {
		return nil
	}

	if !hasGestureEventListener(node.Attributes) {
		return nil
	}

	if hasTouchActionClass(node.Classes, node.RawClasses) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Element declares custom gesture event handlers (touch/pointer move) without CSS 'touch-action'. Native browser scrolling and gestures will clash with script touch handling.",
			Hint:     "Add a CSS touch-action utility (e.g. 'touch-pan-y' for horizontal carousels, or 'touch-none' for free canvas dragging) to isolate compositor gesture handling.",
		},
	}
}

func hasGestureEventListener(attrs map[string]string) bool {
	for k := range attrs {
		lowerKey := strings.ToLower(k)
		if lowerKey == "ontouchmove" || lowerKey == "onpointermove" {
			return true
		}
	}
	return false
}

var touchActionPrefixes = [...]string{
	"touch-none", "touch-pan-x", "touch-pan-y", "touch-pan-left", "touch-pan-right",
	"touch-pan-up", "touch-pan-down", "touch-pinch-zoom", "touch-manipulation", "touch-auto",
}

func hasTouchActionClass(classes []string, rawClasses string) bool {
	for _, pfx := range touchActionPrefixes {
		if strings.Contains(rawClasses, pfx) {
			return true
		}
	}

	for _, cls := range classes {
		for _, pfx := range touchActionPrefixes {
			if cls == pfx || strings.HasPrefix(cls, pfx) {
				return true
			}
		}
	}

	return false
}
