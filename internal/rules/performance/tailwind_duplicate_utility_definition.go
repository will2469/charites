package performance

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// TailwindDuplicateUtilityDefinitionRule mengaudit deklarasi @utility kustom yang menduplikasi utilitas core bawaan.
type TailwindDuplicateUtilityDefinitionRule struct{}

// NewTailwindDuplicateUtilityDefinitionRule membuat instance baru dari TailwindDuplicateUtilityDefinitionRule.
func NewTailwindDuplicateUtilityDefinitionRule() *TailwindDuplicateUtilityDefinitionRule {
	return &TailwindDuplicateUtilityDefinitionRule{}
}

// ID mengembalikan identifier unik kanonikal aturan.
func (r *TailwindDuplicateUtilityDefinitionRule) ID() string {
	return "performance.tailwind-duplicate-utility-definition"
}

// Category mengembalikan kategori aturan ('performance').
func (r *TailwindDuplicateUtilityDefinitionRule) Category() string {
	return "performance"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *TailwindDuplicateUtilityDefinitionRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Description mengembalikan deskripsi tujuan dan konteks aturan.
func (r *TailwindDuplicateUtilityDefinitionRule) Description() string {
	return "Mencegah duplikasi deklarasi utilitas CSS kustom (@utility) yang properti dan nilainya sudah disediakan oleh utilitas core bawaan Tailwind CSS v4."
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *TailwindDuplicateUtilityDefinitionRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Tailwind CSS v4 '@utility' Directive Specification",
			"Compiled CSS Output Deduplication & Bundle Hygiene",
			"Atomic CSS Design Invariants",
		},
		CoreInvariant: "Custom '@utility' declarations must not duplicate built-in Tailwind CSS core utilities; redundant definitions generate unnecessary stylesheet bytes and break atomic CSS composability.",
		Grounding: "The `@utility` directive in Tailwind CSS v4 is designed to register brand-new utilities for modern or proprietary CSS features not yet included in core Tailwind.\n\n" +
			"Defining custom `@utility` blocks for combinations already covered by core utilities (such as `@utility center-flex { display: flex; align-items: center; }`) produces duplicate CSS rules in the compiled stylesheet.\n\n" +
			"Composing canonical core utilities (`flex items-center`) directly in markup preserves atomic stylesheet economy and avoids redundant selector bloat.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Compiled Stylesheet Bloat",
				Severity: "MEDIUM",
				Impact:   "Adds unnecessary custom CSS rules to the production build that duplicate pre-existing atomic utilities.",
			},
			{
				Vector:   "Bypassed Utility Composability",
				Severity: "LOW",
				Impact:   "Custom wrapper utilities fracture atomic consistency and make component classes harder to refactor.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "css",
				Comment:  "Mendefinisikan @utility yang menduplikasi utilitas core flexbox",
				Code: `@utility center-flex {
  display: flex;
  align-items: center;
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Menggunakan kombinasi utilitas core native langsung di markup",
				Code:     `<div className="flex items-center">Konten</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah terdapat deklarasi @utility yang menduplikasi utilitas core.
func (r *TailwindDuplicateUtilityDefinitionRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil {
		return nil
	}

	var fileSrc string
	isStyle := false
	switch {
	case strings.EqualFold(node.Tag, "style"):
		fileSrc = getStyleNodeText(node)
		isStyle = true
	case isSourceRootOrScript(node):
		fileSrc = getFileSourceContent(node)
	default:
		return nil
	}

	if len(fileSrc) == 0 {
		return nil
	}

	violations := findDuplicateUtilityDefinitions(fileSrc)
	if len(violations) == 0 {
		return nil
	}

	diags := make([]ir.Diagnostic, 0, len(violations))
	for _, v := range violations {
		line := v.Line
		if isStyle && node.Span.Line > 0 {
			line = node.Span.Line + v.Line - 1
		}
		diags = append(diags, ir.Diagnostic{
			Line:     line,
			Column:   1,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Custom '@utility %s' duplicates native built-in Tailwind CSS utilities ('%s'). Redundant utility definitions bloat compiled CSS bundles.", v.UtilityName, v.CoreEquiv),
			Hint:     fmt.Sprintf("Remove custom '@utility %s' and compose canonical utilities ('%s') directly in markup.", v.UtilityName, v.CoreEquiv),
		})
	}

	return diags
}
