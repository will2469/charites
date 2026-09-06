package cls

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// DynamicTableReflowRule mendeteksi elemen tabel (<table>) yang merender baris data secara dinamis
// tanpa strategi penentuan ukuran kolom statis (kelas table-fixed, deklarasi <colgroup>, atau lebar <th> terdefinisi).
type DynamicTableReflowRule struct{}

// NewDynamicTableReflowRule membuat instance baru dari DynamicTableReflowRule.
func NewDynamicTableReflowRule() *DynamicTableReflowRule {
	return &DynamicTableReflowRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *DynamicTableReflowRule) ID() string {
	return "cls.dynamic-table-reflow"
}

// Description mengembalikan ringkasan aturan.
func (r *DynamicTableReflowRule) Description() string {
	return "Dynamic <table> lacks a statically inferable column sizing strategy, risking continuous column reflow"
}

// Category mengembalikan nama kategori rule.
func (r *DynamicTableReflowRule) Category() string {
	return "cls"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *DynamicTableReflowRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *DynamicTableReflowRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C HTML Living Standard (Table Rendering & Column Sizing)",
			"CSS Table Module Level 3 (table-layout: fixed)",
			"Google Core Web Vitals (CLS Prevention in Dynamic Data Tables)",
		},
		CoreInvariant: "Dynamic data <table> elements must declare a statically inferable column sizing strategy via 'table-fixed', a <colgroup> block, or explicit width on all header cells.",
		Grounding: "By default, HTML tables operate under 'table-layout: auto', where column widths are continuously recalculated based on the widest cell content across all loaded rows.\n\n" +
			"When tables render dynamic data (e.g. streaming responses, paginated arrays, or WebSocket feeds), incoming rows with varying text lengths force the browser to recalculate and shift every column boundary on each render pass.\n\n" +
			"Using 'table-fixed' (CSS 'table-layout: fixed') or declaring explicit column widths via '<colgroup><col className=\"w-1/3\" />...</colgroup>' ensures the browser determines column boundaries immediately from the first row or colgroup specification, eliminating table reflow entirely.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Continuous Column Boundary Snapping",
				Severity: "MEDIUM",
				Impact:   "As dynamic data streams or updates, table column borders snap horizontally, producing high cumulative layout shift.",
			},
			{
				Vector:   "Delayed Initial Table Paint",
				Severity: "LOW",
				Impact:   "Under table-layout: auto, browser rendering of table rows is deferred until content lengths across all cells are computed.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Dynamic table rendering data items without column sizing strategy",
				Code: `<table className="w-full">
  <tbody>
    {items.map(it => (
      <tr key={it.id}>
        <td>{it.name}</td>
        <td>{it.price}</td>
      </tr>
    ))}
  </tbody>
</table>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Dynamic table locked with Tailwind table-fixed class",
				Code: `<table className="w-full table-fixed">
  <tbody>
    {items.map(it => (
      <tr key={it.id}>
        <td>{it.name}</td>
        <td>{it.price}</td>
      </tr>
    ))}
  </tbody>
</table>`,
			},
			{
				Language: "tsx",
				Comment:  "Dynamic table with explicit colgroup definition",
				Code: `<table className="w-full">
  <colgroup>
    <col className="w-3/4" />
    <col className="w-1/4" />
  </colgroup>
  <tbody>
    {items.map(it => (
      <tr key={it.id}>
        <td>{it.name}</td>
        <td>{it.price}</td>
      </tr>
    ))}
  </tbody>
</table>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen <table> dinamis memiliki strategi penentuan ukuran kolom statis.
func (r *DynamicTableReflowRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || !strings.EqualFold(node.Tag, "table") {
		return nil
	}

	// 1. Periksa apakah tabel merupakan tabel data dinamis
	if !isDynamicTable(node) {
		return nil
	}

	// 2. Periksa apakah tabel memiliki strategi penentuan ukuran kolom statis
	if hasTableFixed(node) || hasColGroupWithCols(node) || hasSizedHeaders(node) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Dynamic <table> lacks a statically inferable column sizing strategy. Unconstrained dynamic tables trigger continuous column reflow and Cumulative Layout Shift (CLS) as rows stream or load.",
			Hint:     "Add class 'table-fixed' (CSS 'table-layout: fixed') or define explicit column widths using <colgroup><col class='w-*' /></colgroup>.",
		},
	}
}
