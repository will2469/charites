package theme

import (
	"github.com/will2469/charites/internal/ir"
)

// InlineStyleHardcodeRule mendeteksi penggunaan warna hardcoded di dalam atribut style HTML/JSX
// yang mengabaikan pewarisan cascade dan memotong adaptabilitas tema.
type InlineStyleHardcodeRule struct{}

// NewInlineStyleHardcodeRule membuat instance baru InlineStyleHardcodeRule.
func NewInlineStyleHardcodeRule() *InlineStyleHardcodeRule {
	return &InlineStyleHardcodeRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *InlineStyleHardcodeRule) ID() string {
	return "theme.inline-style-hardcode"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *InlineStyleHardcodeRule) Description() string {
	return "Detects hardcoded color literals inside HTML/JSX style attributes that prevent theme cascade"
}

// Category mengembalikan nama kategori rule.
func (r *InlineStyleHardcodeRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *InlineStyleHardcodeRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki otomatis.
func (r *InlineStyleHardcodeRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Cascading Style Sheets (CSS) Level 3",
			"W3C Design Tokens Community Group (DTCG)",
		},
		CoreInvariant: "Color properties must not be declared as raw literals inside inline style attributes; they must use semantic classes or CSS variables.",
		Grounding: "Inline style attributes have the highest specificity in CSS, superseding all class selectors and theme cascades.\n\n" +
			"When developers write style=\"color: #2563eb\" or style={{ background: '#fff' }}:\n" +
			"1. Impossible Dark Mode: The inline declaration cannot be targeted or overridden by .dark or [data-theme='dark'] class rules.\n" +
			"2. Broken Theming Pipeline: Token transformations (such as high-contrast mode or tenant styling) fail completely.\n" +
			"3. Maintenance Pitfall: Colors hidden in inline style strings avoid static analysis tools unless specifically parsed.\n\n" +
			"Charites enforces moving inline colors into utility classes or CSS custom properties.",
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Hardcoded hex in HTML inline style",
				Code:     `<div style="color: #2563eb; background: #ffffff;">Inline Color</div>`,
			},
			{
				Language: "tsx",
				Comment:  "Hardcoded rgb in JSX style object",
				Code: `export function Card() {
  return <div style={{ color: '#2563eb', backgroundColor: 'rgb(255, 0, 0)' }}>Bad Style</div>;
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Semantic utility classes instead of inline style",
				Code:     `<div class="text-primary bg-background">Themed Color</div>`,
			},
			{
				Language: "tsx",
				Comment:  "CSS variable in inline style for dynamic calculations",
				Code: `export function Card() {
  return <div style={{ color: 'var(--primary)' }}>Safe Style</div>;
}`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Theme Specificity Lockout",
				Severity: "HIGH",
				Impact:   "Inline style specificity completely disables dark mode and stylesheet theming.",
			},
			{
				Vector:   "Accessibility Barrier",
				Severity: "HIGH",
				Impact:   "High-contrast mode and accessibility themes cannot override inline hardcoded styles.",
			},
		},
	}
}

// Evaluate mengevaluasi node IR dan mendeteksi warna hardcoded di dalam atribut style.
func (r *InlineStyleHardcodeRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || len(node.Attributes) == 0 {
		return nil
	}

	styleVal, ok := node.GetAttr("style")
	if !ok || styleVal == "" {
		return nil
	}

	hardcodedDecl, hasHardcoded := ExtractInlineStyleHardcodedColor(styleVal)
	if !hasHardcoded {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Hardcoded color in inline style: \"" + hardcodedDecl + "\"",
			Hint:     "Move inline color styles to semantic utility classes or CSS variables to allow theme cascading.",
		},
	}
}
