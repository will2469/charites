package theme

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// DynamicClassRule mendeteksi konstruksi kelas dinamis berbasis template literal
// yang merekatkan ekspresi variabel langsung pada prefix atau suffix utility Tailwind (seperti `text-${color}-500`),
// yang merusak deteksi statis Tailwind JIT compiler pada production build.
type DynamicClassRule struct{}

// NewDynamicClassRule membuat instance baru DynamicClassRule.
func NewDynamicClassRule() *DynamicClassRule {
	return &DynamicClassRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *DynamicClassRule) ID() string {
	return "theme.dynamic-class"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *DynamicClassRule) Description() string {
	return "Detects unpadded dynamic template strings breaking Tailwind JIT class generation"
}

// Category mengembalikan nama kategori rule.
func (r *DynamicClassRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *DynamicClassRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki otomatis.
func (r *DynamicClassRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Tailwind CSS JIT Static Analysis & Extraction Guidelines",
			"Build-Time CSS Zero-Runtime Architecture",
			"W3C Web Performance & Production Reliability",
		},
		CoreInvariant: "Utility classes must be written as complete static string literals so the Tailwind build compiler can reliably extract and generate them.",
		Grounding: "Tailwind CSS searches source files using regular expressions looking for complete class strings at build time. It does not evaluate JavaScript at runtime.\n\n" +
			"When developers dynamically construct utility classes using template literal slicing (e.g. className={`text-${color}-500`} or `bg-${variant}`):\n" +
			"1. Missing Production CSS: The Tailwind compiler never matches the interpolated string, leaving the utility completely absent from the production stylesheet.\n" +
			"2. Silent Visual Degradation: The component appears broken or unstyled in production while appearing to work intermittently in dev if another component imported that class.\n" +
			"3. Inscrutable Debugging: Developers struggle to trace why specific color variants intermittently fail to render.\n\n" +
			"Charites enforces using static class maps or complete utility strings within conditional expressions.",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Dynamic class string splicing in JSX className",
				Code:     `<div className={` + "`" + `text-${color}-500 font-bold` + "`" + `}>Status</div>`,
			},
			{
				Language: "astro",
				Comment:  "Dynamic background variant splicing in Astro",
				Code:     `<button class={` + "`" + `px-4 py-2 bg-${variant}` + "`" + `}>Action</button>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Static class lookup map for dynamic variants",
				Code: `const colorMap: Record<string, string> = {
  red: "text-red-500",
  blue: "text-blue-500",
  green: "text-green-500",
};
<div className={` + "`" + `${colorMap[color]} font-bold` + "`" + `}>Status</div>`,
			},
			{
				Language: "tsx",
				Comment:  "Complete utility class strings in ternary expression",
				Code:     `<button className={` + "`" + `px-4 py-2 ${isActive ? "bg-primary text-primary-foreground" : "bg-muted"}` + "`" + `}>Action</button>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Missing Production Styles",
				Severity: "CRITICAL",
				Impact:   "Tailwind JIT engine strips un-scanned utility classes from production bundles, breaking layout and colors.",
			},
			{
				Vector:   "Heisenbug UI Regressions",
				Severity: "HIGH",
				Impact:   "Styles intermittently vanish or break depending on which other files are compiled in the same build chunk.",
			},
		},
	}
}

// Evaluate mengevaluasi node IR untuk memeriksa apakah ada template literal class yang merusak ekstraksi Tailwind.
func (r *DynamicClassRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil {
		return nil
	}

	raw := node.RawClasses
	if raw == "" {
		raw = node.Attributes["className"]
		if raw == "" {
			raw = node.Attributes["class"]
		}
	}

	if raw == "" || !strings.Contains(raw, "${") {
		return nil
	}

	if hasBrokenDynamicClass(raw) {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Dynamic template literal class splicing breaks Tailwind JIT compiler static extraction",
				Hint:     "Write complete utility names (e.g. use a static lookup map or ternary with full classes like 'text-red-500').",
			},
		}
	}

	return nil
}

func hasBrokenDynamicClass(raw string) bool {
	idx := 0
	for {
		dollarIdx := strings.Index(raw[idx:], "${")
		if dollarIdx == -1 {
			break
		}
		pos := idx + dollarIdx

		// 1. Periksa karakter sebelum ${: jika merupakan tanda hubung '-' (seperti text-${, bg-${)
		if pos > 0 && raw[pos-1] == '-' {
			return true
		}

		// Cari penutup kurung '}' yang sesuai
		closeIdx := strings.IndexByte(raw[pos+2:], '}')
		if closeIdx != -1 {
			endPos := pos + 2 + closeIdx
			// 2. Periksa karakter setelah }: jika ada tanda hubung '-' langsung menempel (seperti }-${ atau }-500)
			if endPos+1 < len(raw) && raw[endPos+1] == '-' {
				return true
			}
			idx = endPos + 1
		} else {
			idx = pos + 2
		}
	}

	return false
}
