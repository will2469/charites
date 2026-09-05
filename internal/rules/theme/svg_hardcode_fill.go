package theme

import (
	"slices"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// SVGHardcodeFillRule mendeteksi penggunaan warna mentah (heksadesimal atau palet primitif)
// pada atribut elemen SVG (seperti fill="#000", stroke="#3b82f6", stop-color="#fff").
type SVGHardcodeFillRule struct{}

// NewSVGHardcodeFillRule membuat instance baru SVGHardcodeFillRule.
func NewSVGHardcodeFillRule() *SVGHardcodeFillRule {
	return &SVGHardcodeFillRule{}
}

// ID mengembalikan Charites Rule ID kanonikal.
func (r *SVGHardcodeFillRule) ID() string {
	return "theme.svg-hardcode-fill"
}

// Description mengembalikan penjelasan ringkas rule.
func (r *SVGHardcodeFillRule) Description() string {
	return "Detects hardcoded color attributes on SVG markup preventing theme adaptation"
}

// Category mengembalikan nama kategori rule.
func (r *SVGHardcodeFillRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *SVGHardcodeFillRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *SVGHardcodeFillRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C SVG 2 Specification (Styling & currentColor)",
			"WCAG 2.2 Success Criterion 1.4.11 (Non-text Contrast)",
			"Design System Scalable Iconography Architecture",
		},
		CoreInvariant: "SVG vector elements must inherit colors dynamically via currentColor or semantic CSS variables, never hardcoded hex or primitive colors.",
		Grounding: "Directly hardcoding raw colors onto SVG elements (such as <path fill=\"#000000\"> or <stop stop-color=\"#3b82f6\">) locks graphics to a static palette:\n\n" +
			"1. Theme Blindness: Dark icons with fill=\"#000\" vanish when the user toggles dark mode.\n" +
			"2. Inverted Hover/Active States: Hardcoded stroke attributes prevent buttons and navigation links from changing icon color on hover or focus.\n" +
			"3. Reusability Breakdown: Components cannot share identical SVG glyphs across varying semantic surfaces without duplicating markup.\n\n" +
			"Charites enforces dynamic inheritance using fill=\"currentColor\", stroke=\"currentColor\", or semantic design tokens (var(--primary)).",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Hardcoded hex fill on SVG path in TSX",
				Code:     `<path fill="#000000" d="M10 10 H 90 V 90 H 10 Z" />`,
			},
			{
				Language: "astro",
				Comment:  "Primitive hex stop-color and stroke in Astro SVG",
				Code: `<svg viewBox="0 0 100 100">
  <stop stop-color="#3b82f6" offset="100%" />
  <circle cx="50" cy="50" r="40" stroke="#ef4444" fill="none" />
</svg>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Adaptive currentColor fill in TSX",
				Code:     `<path fill="currentColor" d="M10 10 H 90 V 90 H 10 Z" />`,
			},
			{
				Language: "astro",
				Comment:  "Dynamic CSS variable in gradient stop and currentColor stroke",
				Code: `<svg viewBox="0 0 100 100">
  <stop stop-color="var(--primary)" offset="100%" />
  <circle cx="50" cy="50" r="40" stroke="currentColor" fill="none" />
</svg>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Dark Mode Icon Invisibility",
				Severity: "HIGH",
				Impact:   "Vector icons hardcoded to black or dark shades become completely invisible against dark backgrounds.",
			},
			{
				Vector:   "Broken State Affordance",
				Severity: "MEDIUM",
				Impact:   "Icons fail to inherit hover, focus, and disabled states from parent interactive components.",
			},
		},
	}
}

// Evaluate mengevaluasi sebuah node IR untuk mendeteksi atribut warna SVG yang di-hardcode.
func (r *SVGHardcodeFillRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || !IsSVGElement(node.Tag) || len(node.Attributes) == 0 {
		return nil
	}

	var diags []ir.Diagnostic
	sortedKeys := make([]string, 0, len(node.Attributes))
	for k := range node.Attributes {
		sortedKeys = append(sortedKeys, k)
	}
	slices.Sort(sortedKeys)

	for _, k := range sortedKeys {
		v := node.Attributes[k]
		cleanVal := strings.Trim(strings.TrimSpace(v), "\"'`")
		if IsHardcodedSVGAttribute(k, cleanVal) {
			diags = append(diags, ir.Diagnostic{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Hardcoded SVG color attribute: " + k + "=\"" + cleanVal + "\"",
				Hint:     "Use \"currentColor\" or \"var(--...)\" so vector graphics dynamically inherit theme colors.",
			})
		}
	}

	return diags
}
