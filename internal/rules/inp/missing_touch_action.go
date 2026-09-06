package inp

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// MissingTouchActionRule mendeteksi elemen interaktif berkait handler gestur sentuhan/pointer kustom
// yang tidak mendefinisikan kebijakan CSS touch-action eksplisit, sehingga memicu keterlambatan
// pengenalan gestur (*gesture disambiguation delay*) pada thread sentuh browser.
type MissingTouchActionRule struct{}

// NewMissingTouchActionRule membuat instance baru dari MissingTouchActionRule.
func NewMissingTouchActionRule() *MissingTouchActionRule {
	return &MissingTouchActionRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *MissingTouchActionRule) ID() string {
	return "inp.missing-touch-action"
}

// Description mengembalikan ringkasan aturan.
func (r *MissingTouchActionRule) Description() string {
	return "Interactive element with custom pointer or touch gesture handlers lacks an explicit touch-action CSS policy"
}

// Category mengembalikan nama kategori rule.
func (r *MissingTouchActionRule) Category() string {
	return "inp"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *MissingTouchActionRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *MissingTouchActionRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Google Chrome Core Web Vitals (Interaction to Next Paint Input Delay)",
			"W3C Pointer Events Level 3 (touch-action property)",
			"Tailwind CSS v4 Touch Action Utilities (touch-pan-y, touch-none)",
		},
		CoreInvariant: "Interactive elements implementing custom touch or pointer gesture handlers must declare an explicit 'touch-action' CSS policy ('touch-none', 'touch-pan-y', etc.) to eliminate browser gesture disambiguation delay on the compositor thread.",
		Grounding: "When a user touches an element with custom gesture handlers (such as 'onPointerDown' or 'onTouchStart'), the browser compositor thread cannot know whether the gesture will be handled by JavaScript or defaulted to native panning/zooming.\n\n" +
			"The browser must wait for the JavaScript event handler to execute or call 'preventDefault()', introducing a 100ms-300ms gesture disambiguation delay into every touch interaction.\n\n" +
			"Declaring an explicit CSS 'touch-action' policy (e.g. 'touch-none' for free drag handles or canvas widgets, or 'touch-pan-y' for horizontal swipe carousels) immediately signals the compositor thread to route or bypass native scrolling instantly without waiting for JavaScript.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Compositor Gesture Disambiguation Delay",
				Severity: "HIGH",
				Impact:   "Touch gestures suffer 100ms-300ms latency while the browser waits to resolve potential scrolling conflicts.",
			},
			{
				Vector:   "Scroll Contention & Touch Stutter",
				Severity: "MEDIUM",
				Impact:   "Custom drag widgets conflict with native vertical viewport scrolling on mobile touchscreens.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Custom drag handle without CSS touch-action policy",
				Code: `<div onPointerDown={handleDragStart} onPointerMove={handleDragMove} className="w-full h-48 bg-muted">
  <DragHandle />
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Explicit touch-none utility routing all gestures directly to custom handler",
				Code: `<div onPointerDown={handleDragStart} onPointerMove={handleDragMove} className="w-full h-48 bg-muted touch-none">
  <DragHandle />
</div>`,
			},
			{
				Language: "tsx",
				Comment:  "Horizontal swipe carousel declaring vertical panning freedom",
				Code: `<div onTouchStart={handleSwipeStart} className="flex overflow-x-auto touch-pan-y">
  <CarouselSlide />
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen dengan handler gestur memiliki kebijakan CSS touch-action.
func (r *MissingTouchActionRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	handler, missing := getMissingTouchAction(node.Tag, node.Attributes, node.Classes, node.RawClasses)
	if missing {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Interactive element '<%s>' handles custom gesture event '%s' without specifying an explicit CSS 'touch-action' policy. This triggers browser gesture disambiguation delay on the compositor thread.", node.Tag, handler),
				Hint:     "Add a touch-action utility (e.g. 'touch-none' for free drag/canvas widgets, or 'touch-pan-y' for horizontal sliders) to eliminate touch delay.",
			},
		}
	}

	return nil
}
