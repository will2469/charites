package theme

import (
	"slices"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// ChartColorHardcodeRule mendeteksi penggunaan warna mentah (heksadesimal atau rgb/hsl)
// pada atribut komponen chart visualisasi data (seperti <Bar fill="#3b82f6" />).
type ChartColorHardcodeRule struct{}

// NewChartColorHardcodeRule membuat instance baru ChartColorHardcodeRule.
func NewChartColorHardcodeRule() *ChartColorHardcodeRule {
	return &ChartColorHardcodeRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *ChartColorHardcodeRule) ID() string {
	return "theme.chart-color-hardcode"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *ChartColorHardcodeRule) Description() string {
	return "Detects hardcoded color values on chart visualization components"
}

// Category mengembalikan nama kategori rule.
func (r *ChartColorHardcodeRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *ChartColorHardcodeRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki otomatis.
func (r *ChartColorHardcodeRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"WCAG 2.2 Success Criterion 1.4.3 (Contrast Minimum)",
			"WCAG 2.2 Success Criterion 1.4.11 (Non-text Contrast)",
			"Accessible Data Visualization Design Tokens",
		},
		CoreInvariant: "Chart components must reference semantic theme tokens (e.g. var(--chart-1)) rather than hardcoded hex or color literals.",
		Grounding: "Data visualization libraries (such as Recharts, Chart.js, or Nivo) rely on SVG fill and stroke attributes to render bars, lines, and areas.\n\n" +
			"When developers hardcode hex colors onto chart elements (e.g. <Bar dataKey=\"sales\" fill=\"#3b82f6\" />):\n" +
			"1. Dark Mode Contrast Inversion: The hardcoded colors clash with dark card backgrounds, failing accessibility contrast minimums.\n" +
			"2. Theme Blindness: Visualizations fail to adapt when switching between light, dark, or high-contrast themes.\n" +
			"3. Fragmented Visual Identity: Brand colors drift between charts and surrounding interface tokens.\n\n" +
			"Charites enforces using CSS custom properties (fill=\"var(--chart-1)\") or dynamic theme mappings.",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Hardcoded hex fill and stroke on Recharts Bar and Line",
				Code: `<>
  <Bar dataKey="revenue" fill="#3b82f6" />
  <Line dataKey="profit" stroke="#10b981" />
</>`,
			},
			{
				Language: "astro",
				Comment:  "Hardcoded color on Area and Cell components",
				Code: `<Area dataKey="uv" fill="#8884d8" stroke="#82ca9d" />
<Cell fill="#f43f5e" />`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Semantic chart tokens from design system",
				Code: `<>
  <Bar dataKey="revenue" fill="var(--chart-1)" />
  <Line dataKey="profit" stroke="var(--chart-2)" />
</>`,
			},
			{
				Language: "astro",
				Comment:  "CSS variable references on Area and Cell",
				Code: `<Area dataKey="uv" fill="var(--chart-1)" stroke="var(--chart-2)" />
<Cell fill="var(--chart-destructive)" />`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Chart Contrast Invalidation",
				Severity: "HIGH",
				Impact:   "Chart bars and lines become illegible against inverted dark backgrounds, obscuring critical analytics.",
			},
			{
				Vector:   "Theme Desynchronization",
				Severity: "MEDIUM",
				Impact:   "Data visualizations remain locked to legacy colors while the rest of the application adapts dynamically.",
			},
		},
	}
}

// Evaluate mengevaluasi node IR untuk mendeteksi warna hardcoded pada elemen chart.
func (r *ChartColorHardcodeRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || len(node.Attributes) == 0 || !isChartElement(node) {
		return nil
	}

	var diags []ir.Diagnostic
	sortedKeys := make([]string, 0, len(node.Attributes))
	for k := range node.Attributes {
		sortedKeys = append(sortedKeys, k)
	}
	slices.Sort(sortedKeys)

	for _, k := range sortedKeys {
		if !isChartColorAttr(k) {
			continue
		}
		v := node.Attributes[k]
		cleanVal := strings.Trim(strings.TrimSpace(v), "\"'`")
		if isHardcodedColorValue(cleanVal) {
			diags = append(diags, ir.Diagnostic{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Hardcoded chart color attribute: " + k + "=\"" + cleanVal + "\"",
				Hint:     "Use semantic CSS variables (e.g. \"var(--chart-1)\") so charts dynamically adapt to light and dark themes.",
			})
		}
	}

	return diags
}

func isChartElement(node *ir.Node) bool {
	switch node.Tag {
	case "Bar", "Line", "Area", "Pie", "Cell", "Scatter", "Radar", "RadialBar", "Treemap", "Funnel":
		return true
	}
	_, hasDataKey := node.Attributes["dataKey"]
	return hasDataKey
}

func isChartColorAttr(attr string) bool {
	switch strings.ToLower(attr) {
	case "fill", "stroke", "stopcolor", "stop-color":
		return true
	}
	return false
}

func isHardcodedColorValue(val string) bool {
	if val == "" || val == "none" || val == "currentColor" || val == "transparent" {
		return false
	}
	if strings.HasPrefix(val, "var(") || strings.HasPrefix(val, "url(") {
		return false
	}
	// Deteksi heksadesimal (#fff, #123456, dll.)
	if strings.HasPrefix(val, "#") {
		return true
	}
	// Deteksi rgb, rgba, hsl, hsla, oklch literal
	lower := strings.ToLower(val)
	if strings.HasPrefix(lower, "rgb(") || strings.HasPrefix(lower, "rgba(") ||
		strings.HasPrefix(lower, "hsl(") || strings.HasPrefix(lower, "hsla(") ||
		strings.HasPrefix(lower, "oklch(") {
		return true
	}
	return false
}
