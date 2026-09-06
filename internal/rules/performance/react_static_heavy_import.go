package performance

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// ReactStaticHeavyImportRule mengaudit pernyataan impor statis pustaka berbobot besar di tingkat atas.
type ReactStaticHeavyImportRule struct{}

// NewReactStaticHeavyImportRule membuat instance baru dari ReactStaticHeavyImportRule.
func NewReactStaticHeavyImportRule() *ReactStaticHeavyImportRule {
	return &ReactStaticHeavyImportRule{}
}

// ID mengembalikan identifier unik kanonikal aturan.
func (r *ReactStaticHeavyImportRule) ID() string {
	return "performance.react-static-heavy-import"
}

// Category mengembalikan kategori aturan ('performance').
func (r *ReactStaticHeavyImportRule) Category() string {
	return "performance"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *ReactStaticHeavyImportRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Description mengembalikan deskripsi tujuan dan konteks aturan.
func (r *ReactStaticHeavyImportRule) Description() string {
	return "Mengaudit pernyataan impor statis modul berukuran besar di tingkat atas yang membengkakkan bundel JavaScript awal dan mewajibkan pemisahan kode via React.lazy() dan <Suspense>."
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ReactStaticHeavyImportRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Web Performance & Initial Route Code Splitting Best Practices",
			"React Official Documentation (Code-Splitting with React.lazy and Suspense)",
			"Chrome DevTools Lighthouse JavaScript Bundle Size Optimization Guidelines",
		},
		CoreInvariant: "Heavy visualization, editing, and utility libraries must be asynchronously loaded via 'React.lazy()' and wrapped in '<Suspense>'; top-level static imports bloat critical initial bundles and degrade FCP/TBT.",
		Grounding: "Top-level static import statements (`import { Chart } from 'chart.js'`) force modern JavaScript bundlers (Webpack, Vite, Rollup, esbuild) to include the imported module directly in the entry-point chunk.\n\n" +
			"Heavy third-party libraries such as `monaco-editor`, `chart.js`, `echarts`, `quill`, `pdfjs-dist`, `three`, or `xlsx` often weigh several hundreds of kilobytes compressed.\n\n" +
			"Because these libraries are typically secondary to initial above-the-fold content, loading them synchronously forces mobile devices to spend hundreds of milliseconds downloading, parsing, and compiling JavaScript before the user can interact with the page.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Initial Bundle Size Bloat",
				Severity: "HIGH",
				Impact:   "Drastically increases First Contentful Paint (FCP) and Total Blocking Time (TBT) due to massive synchronous script downloads.",
			},
			{
				Vector:   "Mobile CPU Decompression Bottlenecks",
				Severity: "MEDIUM",
				Impact:   "Consumes limited mobile device CPU and memory parsing large scripts that are not immediately rendered on initial screen view.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Mengimpor modul grafik berat secara statis di tingkat atas",
				Code: `import { Chart } from 'chart.js';

export function Dashboard() {
  return <Chart data={stats} />;
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Memisahkan modul grafik menggunakan React.lazy() dan Suspense",
				Code: `import { Suspense, lazy } from 'react';
const Chart = lazy(() => import('chart.js'));

export function Dashboard() {
  return (
    <Suspense fallback={<ChartSkeleton />}>
      <Chart data={stats} />
    </Suspense>
  );
}`,
			},
		},
	}
}

// Evaluate memeriksa apakah berkas memuat pustaka berbobot besar secara statis di tingkat atas.
func (r *ReactStaticHeavyImportRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || !isSourceRootOrScript(node) {
		return nil
	}

	fileSrc := getFileSourceContent(node)
	if len(fileSrc) == 0 {
		return nil
	}

	violations := findStaticHeavyImports(fileSrc)
	if len(violations) == 0 {
		return nil
	}

	diags := make([]ir.Diagnostic, 0, len(violations))
	for _, v := range violations {
		diags = append(diags, ir.Diagnostic{
			Line:     v.Line,
			Column:   1,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Heavy module '%s' is imported statically at the top level. Large visualization or processing libraries bloat the initial JavaScript bundle; dynamically load with 'React.lazy()' and wrap in '<Suspense>'.", v.Module),
			Hint:     fmt.Sprintf("Refactor to 'const %s = React.lazy(() => import(\"%s\"));' and render inside a '<Suspense fallback={...}>' boundary.", v.Module, v.Module),
		})
	}

	return diags
}
