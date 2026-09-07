package reporter

import (
	"github.com/will2469/charites/internal/ir"
)

// FileGroup mengelompokkan temuan diagnostik berdasarkan berkas.
type FileGroup struct {
	File         string
	ErrorCount   int
	WarningCount int
	InfoCount    int
	TotalIssues  int
	Diagnostics  []ir.Diagnostic
}

// cloneAndSortDiagnostics membuat salinan slice diagnostik dan mengurutkannya secara kanonikal
// menggunakan implementasi SSOT dari paket internal/ir.
func cloneAndSortDiagnostics(diags []ir.Diagnostic) []ir.Diagnostic {
	if len(diags) == 0 {
		return nil
	}
	cloned := make([]ir.Diagnostic, len(diags))
	copy(cloned, diags)
	return ir.SortDiagnostics(cloned)
}

// GroupByFile mempartisi slice diagnostik yang SUDAH terurut secara kanonikal ke dalam slice FileGroup.
// Fungsi ini HANYA melakukan partisi linier satu lintasan O(N) tanpa melakukan pengurutan ulang.
func GroupByFile(canonicalDiags []ir.Diagnostic) []FileGroup {
	if len(canonicalDiags) == 0 {
		return nil
	}

	var groups []FileGroup
	var current *FileGroup

	for _, d := range canonicalDiags {
		if current == nil || current.File != d.File {
			groups = append(groups, FileGroup{
				File: d.File,
			})
			current = &groups[len(groups)-1]
		}

		current.Diagnostics = append(current.Diagnostics, d)
		current.TotalIssues++

		switch d.Severity {
		case ir.SeverityError:
			current.ErrorCount++
		case ir.SeverityWarn:
			current.WarningCount++
		case ir.SeverityInfo:
			current.InfoCount++
		}
	}

	return groups
}
