package cls

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// LayoutTriggerAnimationRule mendeteksi animasi CSS @keyframes yang memutasi properti geometri tata letak
// pemicu reflow CPU (top, left, width, height, margin, padding) alih-alih transformasi GPU composited (transform, opacity).
type LayoutTriggerAnimationRule struct{}

// NewLayoutTriggerAnimationRule membuat instance baru dari LayoutTriggerAnimationRule.
func NewLayoutTriggerAnimationRule() *LayoutTriggerAnimationRule {
	return &LayoutTriggerAnimationRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *LayoutTriggerAnimationRule) ID() string {
	return "cls.layout-trigger-animation"
}

// Description mengembalikan ringkasan aturan.
func (r *LayoutTriggerAnimationRule) Description() string {
	return "CSS @keyframes animation mutates layout-triggering geometry properties instead of GPU-composited transforms"
}

// Category mengembalikan nama kategori rule.
func (r *LayoutTriggerAnimationRule) Category() string {
	return "cls"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *LayoutTriggerAnimationRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *LayoutTriggerAnimationRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C CSS Animations Level 1 (@keyframes declaration blocks)",
			"Google Core Web Vitals (CLS Compositor Thread Guidelines)",
			"High Performance Mobile Web (GPU Compositing vs CPU Reflow)",
		},
		CoreInvariant: "CSS @keyframes animations must mutate GPU-composited layer properties ('transform', 'opacity') rather than layout-triggering geometry properties ('top', 'left', 'width', 'height', 'margin', 'padding').",
		Grounding: "When CSS keyframes animate geometry properties (such as top, left, width, height, margin, or padding), the browser is forced to execute full layout recalculations (reflow) and repaint stages on the main CPU thread for every animation frame (typically 60-120 times per second).\n\n" +
			"This continuous geometry invalidation directly triggers Cumulative Layout Shift (CLS) for neighboring elements and causes noticeable frame jank (dropped frames) on mobile and low-power hardware.\n\n" +
			"Modern browser rendering pipelines offload 'transform' and 'opacity' mutations directly to the GPU compositor thread, executing smooth, 60fps animations that never invalidate document geometry or cause layout shifts.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Continuous CPU Layout Reflow",
				Severity: "HIGH",
				Impact:   "Animating geometry properties causes browser recalculation of surrounding elements on every frame, generating Cumulative Layout Shift.",
			},
			{
				Vector:   "Rendering Pipeline Jank & Dropped Frames",
				Severity: "MEDIUM",
				Impact:   "High CPU load from continuous layout reflow stalls main thread execution, resulting in choppy animations and poor touch responsiveness.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "css",
				Comment:  "Keyframe animation mutating positional and margin geometry properties",
				Code: `@keyframes slideIn {
  from { top: -20px; margin-top: 10px; }
  to { top: 0; margin-top: 0; }
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "css",
				Comment:  "GPU-composited keyframe animation using transform and opacity",
				Code: `@keyframes slideIn {
  from { transform: translateY(-20px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}`,
			},
		},
	}
}

// Evaluate memeriksa apakah blok <style> memiliki animasi @keyframes yang memutasi properti geometri.
func (r *LayoutTriggerAnimationRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || !strings.EqualFold(node.Tag, "style") {
		return nil
	}

	css := getStyleNodeText(node)
	violations := findLayoutTriggerKeyframes(css)
	if len(violations) == 0 {
		return nil
	}

	diags := make([]ir.Diagnostic, 0, len(violations))
	for _, v := range violations {
		diags = append(diags, ir.Diagnostic{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("CSS @keyframes animation '%s' mutates layout-triggering geometry property '%s'. Animating geometry properties forces CPU reflow on every frame, risking Cumulative Layout Shift (CLS). Mutate GPU-composited layer properties ('transform', 'opacity') instead.", v.AnimationName, v.Property),
			Hint:     "Replace positional mutations (top/left) with 'transform: translate(...)', or use 'transform: scale(...)' / 'opacity' for smooth GPU animation without reflow.",
		})
	}

	return diags
}
