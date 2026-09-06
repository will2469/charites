package performance

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// TailwindDuplicateArbitraryRulesRule mengaudit pemakaian nilai arbitrary yang menduplikasi skala bawaan Tailwind v4.
type TailwindDuplicateArbitraryRulesRule struct{}

// NewTailwindDuplicateArbitraryRulesRule membuat instance baru dari TailwindDuplicateArbitraryRulesRule.
func NewTailwindDuplicateArbitraryRulesRule() *TailwindDuplicateArbitraryRulesRule {
	return &TailwindDuplicateArbitraryRulesRule{}
}

// ID mengembalikan identifier unik kanonikal aturan.
func (r *TailwindDuplicateArbitraryRulesRule) ID() string {
	return "performance.tailwind-duplicate-arbitrary-rules"
}

// Category mengembalikan kategori aturan ('performance').
func (r *TailwindDuplicateArbitraryRulesRule) Category() string {
	return "performance"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *TailwindDuplicateArbitraryRulesRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Description mengembalikan deskripsi tujuan dan konteks aturan.
func (r *TailwindDuplicateArbitraryRulesRule) Description() string {
	return "Menganjurkan penggunaan utilitas skala inti bawaan Tailwind v4 alih-alih nilai arbitrary sembarang yang menghasilkan deklarasi CSS duplikat."
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *TailwindDuplicateArbitraryRulesRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Tailwind CSS v4 Design Tokens & Spacing Scale Standards",
			"Compiled CSS Output Deduplication & Payload Economy",
			"W3C Stylesheet Declarative Optimization Guidelines",
		},
		CoreInvariant: "Arbitrary value utilities (e.g. 'p-[16px]', 'mt-[1rem]') that match standard Tailwind core scale tokens should use the canonical core utility (e.g. 'p-4', 'mt-4') to avoid duplicate CSS rule generation in compiled bundles.",
		Grounding: "Tailwind CSS includes a refined, consistent default spacing and sizing scale.\n\n" +
			"When developers write ad-hoc arbitrary values like `p-[16px]` alongside `p-4` (which also resolves to `padding: 1rem / 16px`), Tailwind generates separate unique CSS selector rules for both.\n\n" +
			"Consolidating arbitrary values to their core scale equivalents eliminates redundant rule definitions, shrinks the production CSS footprint, and ensures consistent visual rhythm across the application.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Compiled CSS Bloat",
				Severity: "MEDIUM",
				Impact:   "Inflates stylesheet size with duplicate CSS selector blocks that declare identical CSS properties and values.",
			},
			{
				Vector:   "Visual Rhythm Inconsistency",
				Severity: "LOW",
				Impact:   "Ad-hoc arbitrary values drift away from the design system's cohesive 4px/8px modular grid.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Menggunakan nilai arbitrary yang menduplikasi utilitas core p-4 dan mt-4",
				Code:     `<div className="p-[16px] mt-[1rem]">Konten</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Menggunakan utilitas skala core standar",
				Code:     `<div className="p-4 mt-4">Konten</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah terdapat kelas arbitrary yang menduplikasi utilitas core.
func (r *TailwindDuplicateArbitraryRulesRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || len(node.Classes) == 0 {
		return nil
	}

	arb, coreEquiv, hasDuplicate := findDuplicateArbitraryUtility(node.Classes, node.RawClasses)
	if !hasDuplicate {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Arbitrary value utility '%s' duplicates built-in Tailwind scale utility '%s', bloating compiled CSS output.", arb, coreEquiv),
			Hint:     fmt.Sprintf("Replace arbitrary utility '%s' with canonical core utility '%s'.", arb, coreEquiv),
		},
	}
}
