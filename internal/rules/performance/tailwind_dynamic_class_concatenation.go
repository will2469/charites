package performance

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// TailwindDynamicClassConcatenationRule mengaudit penggabungan string kelas dinamis yang merusak static extraction Tailwind v4.
type TailwindDynamicClassConcatenationRule struct{}

// NewTailwindDynamicClassConcatenationRule membuat instance baru dari TailwindDynamicClassConcatenationRule.
func NewTailwindDynamicClassConcatenationRule() *TailwindDynamicClassConcatenationRule {
	return &TailwindDynamicClassConcatenationRule{}
}

// ID mengembalikan identifier unik kanonikal aturan.
func (r *TailwindDynamicClassConcatenationRule) ID() string {
	return "performance.tailwind-dynamic-class-concatenation"
}

// Category mengembalikan kategori aturan ('performance').
func (r *TailwindDynamicClassConcatenationRule) Category() string {
	return "performance"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *TailwindDynamicClassConcatenationRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Description mengembalikan deskripsi tujuan dan konteks aturan.
func (r *TailwindDynamicClassConcatenationRule) Description() string {
	return "Mencegah penggabungan string nama kelas dinamis parsial yang merusak deteksi compiler scanner Tailwind CSS v4 (Oxide engine)."
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *TailwindDynamicClassConcatenationRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Tailwind CSS v4 Compiler Scanner Specification (Oxide Static Extraction)",
			"Tailwind CSS Official Architecture ('Dynamic Class Names Limitations')",
			"Zero-Runtime CSS Extraction Invariants",
		},
		CoreInvariant: "Tailwind CSS utility classes must be written as complete, static string literals; dynamic string interpolation on partial class prefixes prevents the static scanner from detecting classes, resulting in missing styles in production.",
		Grounding: "Tailwind CSS v4 uses a high-performance static scanner (Oxide engine) that scans source code for complete class tokens without executing JavaScript runtime.\n\n" +
			"Constructing utility names dynamically via template literals or string concatenation (e.g. `bg-${color}-100` or `'text-' + size`) breaks static extraction completely.\n\n" +
			"Because the scanner never evaluates runtime variables, it never sees the complete utility string (like `bg-red-100`), causing the required CSS rules to be omitted from the compiled stylesheet.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Missing Production Stylesheet Rules",
				Severity: "HIGH",
				Impact:   "Utility classes generated through string concatenation are completely missing from the production CSS bundle, causing broken UI visuals.",
			},
			{
				Vector:   "Silent Runtime Failures",
				Severity: "HIGH",
				Impact:   "Classes appear functional in local environments if the class was previously cached or generated elsewhere, but fail silently upon clean production builds.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Penggabungan string parsial tidak dapat diekstrak oleh compiler",
				Code: `function Badge({ color }: { color: 'red' | 'blue' }) {
  return <span className={` + "`bg-${color}-100 text-${color}-800`" + `}>Status</span>;
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Menuliskan nama kelas secara utuh dalam kamus statis",
				Code: `const COLOR_MAP = {
  red: 'bg-red-100 text-red-800',
  blue: 'bg-blue-100 text-blue-800',
} as const;

function Badge({ color }: { color: 'red' | 'blue' }) {
  return <span className={COLOR_MAP[color]}>Status</span>;
}`,
			},
		},
	}
}

// Evaluate memeriksa apakah atribut className atau class memuat interpolasi kelas parsial dinamis.
func (r *TailwindDynamicClassConcatenationRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || len(node.Attributes) == 0 {
		return nil
	}

	for attrName, attrVal := range node.Attributes {
		if attrName != "className" && attrName != "class" {
			continue
		}

		if prefix, hasConcat := findDynamicClassConcatenation(attrVal); hasConcat {
			return []ir.Diagnostic{
				{
					Line:     node.Span.Line,
					Column:   node.Span.Column,
					Rule:     r.ID(),
					Severity: r.DefaultSeverity(),
					Message:  fmt.Sprintf("Dynamic Tailwind class concatenation around prefix '%s' breaks Tailwind CSS v4 static scanner extraction; compiled CSS rules will not be generated in production.", prefix),
					Hint:     "Write complete static utility class names or use an explicit static lookup map (e.g. 'const MAP = { ... }') instead of dynamic string splicing.",
				},
			}
		}
	}

	return nil
}
