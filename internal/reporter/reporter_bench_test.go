package reporter_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/reporter"
)

func generate1000FindingsResult() *reporter.ScanResult {
	diags := make([]ir.Diagnostic, 1000)
	for i := 0; i < 1000; i++ {
		sev := ir.SeverityError
		if i%2 == 0 {
			sev = ir.SeverityWarn
		}
		diags[i] = ir.Diagnostic{
			File:     fmt.Sprintf("src/components/Card%04d.tsx", i),
			Line:     i + 1,
			Column:   10,
			Rule:     "theme.hardcode-opacity-color",
			Severity: sev,
			Message:  "Hardcoded opacity color violation",
			Hint:     "Use design system semantic token",
		}
	}

	return &reporter.ScanResult{
		Version: "1.0.0",
		Summary: reporter.ScanSummary{
			ScannedFiles: 500,
			DurationMS:   150,
			ErrorCount:   500,
			WarningCount: 500,
			InfoCount:    0,
			Passed:       false,
		},
		Diagnostics: diags,
	}
}

// BenchmarkJSONReporter_Render_1000Findings menguji anggaran performa QUAL-05-PERF-001 (<= 5ms).
func BenchmarkJSONReporter_Render_1000Findings(b *testing.B) {
	result := generate1000FindingsResult()
	rep := reporter.NewJSONReporter()

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = rep.Render(io.Discard, result)
	}
}

// BenchmarkInlineReporter_Render_1000Findings menguji anggaran performa QUAL-05-PERF-001 (<= 10ms).
func BenchmarkInlineReporter_Render_1000Findings(b *testing.B) {
	result := generate1000FindingsResult()
	rep := reporter.NewInlineReporter(reporter.ColorNever)

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_ = rep.Render(io.Discard, result)
	}
}
