package responsive

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// UnwrappedTableOverflowRule mendeteksi elemen tabel HTML (<table>) yang tidak dibungkus
// oleh kontainer pengguliran horizontal (overflow-x-auto) atau tanpa transformasi responsif.
type UnwrappedTableOverflowRule struct{}

// NewUnwrappedTableOverflowRule membuat instance baru dari UnwrappedTableOverflowRule.
func NewUnwrappedTableOverflowRule() *UnwrappedTableOverflowRule {
	return &UnwrappedTableOverflowRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *UnwrappedTableOverflowRule) ID() string {
	return "responsive.unwrapped-table-overflow"
}

// Description mengembalikan ringkasan aturan.
func (r *UnwrappedTableOverflowRule) Description() string {
	return "Warns when an HTML table element lacks a responsive horizontal scroll wrapper (overflow-x-auto) or responsive display transformation"
}

// Category mengembalikan nama kategori rule.
func (r *UnwrappedTableOverflowRule) Category() string {
	return "responsive"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *UnwrappedTableOverflowRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *UnwrappedTableOverflowRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C HTML Living Standard (Table Rendering & Intrinsic Sizing)",
			"Responsive Web Design Data Table Patterns",
			"Mobile Touch Usability Guidelines (Scroll Container Isolation)",
		},
		CoreInvariant: "HTML <table> elements must be enclosed within an ancestor container providing horizontal scrolling (overflow-x-auto) or declared with responsive display styling (hidden md:table).",
		Grounding: "On compact smartphone viewports (360px-390px), data tables possess an intrinsic min-content sizing model (table-layout: auto) that prevents columns from shrinking beyond their widest words.\n\n" +
			"Placing a naked <table> element directly into normal document flow forces the entire webpage to blow out horizontally, inducing unwanted page-level horizontal sway and breaking swipe navigation.\n\n" +
			"Wrapping data tables in a dedicated scroll container (<div class=\"overflow-x-auto\">) isolates horizontal scrolling to the table boundaries without disrupting page flow.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Document-Wide Horizontal Blowout",
				Severity: "MEDIUM",
				Impact:   "Entire mobile page wobbles horizontally during scrolling because intrinsic table width exceeds screen boundary.",
			},
			{
				Vector:   "Hidden Data Columns Without Scroll Affordance",
				Severity: "MEDIUM",
				Impact:   "Users on compact devices cannot view right-hand table columns without an explicit horizontal scroll container.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Unwrapped data table directly inside layout causing mobile horizontal overflow",
				Code: `<table className="w-full border">
  <thead>
    <tr><th>Nama</th><th>NIK</th><th>Alamat</th><th>Status</th></tr>
  </thead>
  <tbody>...</tbody>
</table>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Table enclosed within an overflow-x-auto scroll container",
				Code: `<div className="w-full overflow-x-auto">
  <table className="w-full border">
    <thead>
      <tr><th>Nama</th><th>NIK</th><th>Alamat</th><th>Status</th></tr>
    </thead>
    <tbody>...</tbody>
  </table>
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen <table> memiliki scroll wrapper atau penataan responsif.
func (r *UnwrappedTableOverflowRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || !strings.EqualFold(node.Tag, "table") {
		return nil
	}

	if hasResponsiveTableDisplay(node) || hasScrollContainerAncestor(node) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "HTML <table> element lacks a responsive scroll wrapper (e.g. <div class='overflow-x-auto'>) or responsive display classes (e.g. 'hidden md:table'). Unwrapped data tables exceed mobile viewport boundaries and induce horizontal page scrolling.",
			Hint:     "Wrap the table in a container with 'w-full overflow-x-auto' or provide a mobile card alternative using 'hidden md:table'.",
		},
	}
}
