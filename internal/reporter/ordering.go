package reporter

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// normalizePOSIXPath memastikan pemisah direktori senantiasa '/' di semua sistem operasi.
func normalizePOSIXPath(p string) string {
	return strings.ReplaceAll(filepath.ToSlash(p), "\\", "/")
}

// FileGroup mengelompokkan temuan diagnostik berdasarkan berkas normalisasi.
type FileGroup struct {
	File         string
	ErrorCount   int
	WarningCount int
	InfoCount    int
	TotalIssues  int
	Diagnostics  []ir.Diagnostic
}

// SortDiagnosticsCanonical melakukan total ordering kanonikal 7-dimensi deterministik:
// File ASC -> Line ASC -> Column ASC -> Rule ASC -> Severity ASC -> Message ASC -> Hint ASC.
func SortDiagnosticsCanonical(diags []ir.Diagnostic) []ir.Diagnostic {
	if len(diags) == 0 {
		return nil
	}

	sorted := make([]ir.Diagnostic, len(diags))
	copy(sorted, diags)

	sort.SliceStable(sorted, func(i, j int) bool {
		a, b := sorted[i], sorted[j]
		fileA := normalizePOSIXPath(a.File)
		fileB := normalizePOSIXPath(b.File)
		if fileA != fileB {
			return fileA < fileB
		}
		if a.Line != b.Line {
			return a.Line < b.Line
		}
		if a.Column != b.Column {
			return a.Column < b.Column
		}
		if a.Rule != b.Rule {
			return a.Rule < b.Rule
		}
		if a.Severity != b.Severity {
			return a.Severity < b.Severity
		}
		if a.Message != b.Message {
			return a.Message < b.Message
		}
		return a.Hint < b.Hint
	})

	return sorted
}

// GroupByFile mengelompokkan diagnostik ke dalam slice FileGroup yang terurut secara kanonikal.
// Diagnostik diurutkan sekali sebelum dipartisi sehingga tiap FileGroup dan daftar pelanggaran
// di dalamnya otomatis mewarisi pengurutan kanonikal (Invarian I1, I2, dan I3).
func GroupByFile(diags []ir.Diagnostic) []FileGroup {
	if len(diags) == 0 {
		return nil
	}

	sorted := SortDiagnosticsCanonical(diags)
	var groups []FileGroup
	var current *FileGroup

	for _, d := range sorted {
		normFile := normalizePOSIXPath(d.File)
		if current == nil || current.File != normFile {
			groups = append(groups, FileGroup{
				File: normFile,
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

// NormalizeSummary menyelaraskan field agregasi ringkasan ScanSummary terhadap temuan aktual
// sehingga invarian I4-I8 senantiasa terpenuhi secara matematis.
func NormalizeSummary(s *ScanSummary, files []FileGroup, totalScanned int) {
	if s == nil {
		return
	}

	var errCount, warnCount, infoCount, totalIssues int
	for _, fg := range files {
		errCount += fg.ErrorCount
		warnCount += fg.WarningCount
		infoCount += fg.InfoCount
		totalIssues += fg.TotalIssues
	}

	s.ErrorCount = errCount
	s.WarningCount = warnCount
	s.InfoCount = infoCount
	s.TotalIssues = totalIssues
	s.FilesWithIssues = len(files)

	if totalScanned > 0 {
		s.ScannedFiles = totalScanned
	}
	if s.ScannedFiles < s.FilesWithIssues {
		s.ScannedFiles = s.FilesWithIssues
	}
	s.CleanFiles = s.ScannedFiles - s.FilesWithIssues
	if s.CleanFiles < 0 {
		s.CleanFiles = 0
	}
	s.Passed = s.ErrorCount == 0
}
