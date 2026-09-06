package cls

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// LayoutTriggerTransitionRule mendeteksi deklarasi transisi CSS atau utilitas Tailwind transition-all
// yang menargetkan properti geometri tata letak (width, height, margin, padding, top, left) alih-alih GPU composited layer.
type LayoutTriggerTransitionRule struct{}

// NewLayoutTriggerTransitionRule membuat instance baru dari LayoutTriggerTransitionRule.
func NewLayoutTriggerTransitionRule() *LayoutTriggerTransitionRule {
	return &LayoutTriggerTransitionRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *LayoutTriggerTransitionRule) ID() string {
	return "cls.layout-trigger-transition"
}

// Description mengembalikan ringkasan aturan.
func (r *LayoutTriggerTransitionRule) Description() string {
	return "CSS transition targets layout-triggering geometry properties instead of GPU-composited transforms"
}

// Category mengembalikan nama kategori rule.
func (r *LayoutTriggerTransitionRule) Category() string {
	return "cls"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *LayoutTriggerTransitionRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *LayoutTriggerTransitionRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C CSS Transitions Level 1 (transition-property)",
			"Google Core Web Vitals (CLS & Animation Frame Stability)",
			"Tailwind CSS v4 Transition Best Practices",
		},
		CoreInvariant: "CSS transitions must avoid animating layout-triggering geometry properties ('width', 'height', 'margin', 'padding', 'top', 'left') and instead utilize GPU-composited 'transform' or 'opacity'.",
		Grounding: "Transitioning geometry properties (such as width, height, padding, or positional offsets) triggers continuous CPU layout recalculations and repaints throughout the transition duration.\n\n" +
			"When geometry transitions execute on interactive elements (e.g., hover expansion or focus state enlargement), neighboring elements are continuously pushed and shifted, generating Cumulative Layout Shift (CLS).\n\n" +
			"Transitioning 'transform' (e.g. scale, translate) or 'opacity' executes entirely on the GPU compositor thread without triggering layout reflow.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Interactive Reflow Stalling",
				Severity: "MEDIUM",
				Impact:   "Hovering or focusing interactive elements triggers continuous layout passes, causing adjacent content to jitter and shift.",
			},
			{
				Vector:   "Dropped Frames During Hover Transitions",
				Severity: "MEDIUM",
				Impact:   "Main-thread CPU layout calculations during mousemove or hover states cause micro-stutters and responsiveness lag.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "css",
				Comment:  "CSS declaration transitioning width directly",
				Code: `.sidebar {
  transition: width 300ms ease-in-out;
}`,
			},
			{
				Language: "tsx",
				Comment:  "Tailwind transition-all combined with hover geometry mutation",
				Code: `<div className="w-32 transition-all hover:w-64">
  Sidebar
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "css",
				Comment:  "CSS transition utilizing GPU-composited transform",
				Code: `.sidebar {
  transition: transform 300ms ease-in-out;
}`,
			},
			{
				Language: "tsx",
				Comment:  "Tailwind transition-transform with scale",
				Code: `<div className="w-32 transition-transform hover:scale-110">
  Sidebar
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah node (baik elemen bertata letak atau blok <style>) memiliki transisi geometri.
func (r *LayoutTriggerTransitionRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil {
		return nil
	}

	// 1. Kasus blok <style>: periksa deklarasi CSS transition: ...
	if node.Type == ir.NodeElement && strings.EqualFold(node.Tag, "style") {
		css := getStyleNodeText(node)
		violations := findCSSTransitions(css)
		if len(violations) == 0 {
			return nil
		}
		diags := make([]ir.Diagnostic, 0, len(violations))
		for _, prop := range violations {
			diags = append(diags, ir.Diagnostic{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("CSS transition targets layout-triggering geometry property '%s'. Transitioning geometry forces continuous CPU layout recalculations and risks Cumulative Layout Shift (CLS).", prop),
				Hint:     "Use 'transition: transform' with 'transform: scale(...)' or 'translate(...)', or isolate layout with 'contain: layout'.",
			})
		}
		return diags
	}

	// 2. Kasus elemen biasa: periksa kelas Tailwind transition-all / transition-[prop]
	if node.Type == ir.NodeElement {
		offending, hasRisk := hasGeometryTransition(node)
		if !hasRisk {
			return nil
		}
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Element transition targets layout-triggering geometry mutation '%s'. Transitioning geometry forces continuous CPU layout recalculations and risks Cumulative Layout Shift (CLS).", offending),
				Hint:     "Replace 'transition-all' with 'transition-transform' or 'transition-opacity', or isolate layout using 'contain-layout'.",
			},
		}
	}

	return nil
}
