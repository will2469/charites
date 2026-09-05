package reporter

import (
	"io"

	"github.com/will2469/charites/internal/ir"
)

// ScanSummary merepresentasikan ringkasan agregasi metrik dari eksekusi pemindaian.
type ScanSummary struct {
	ScannedFiles int   `json:"scanned_files"`
	DurationMS   int64 `json:"duration_ms"`
	ErrorCount   int   `json:"error_count"`
	WarningCount int   `json:"warning_count"`
	InfoCount    int   `json:"info_count"`
	Passed       bool  `json:"passed"`
}

// ScanResult merepresentasikan struktur dokumen lengkap hasil analisis kode.
type ScanResult struct {
	Version     string          `json:"version"`
	Summary     ScanSummary     `json:"summary"`
	Diagnostics []ir.Diagnostic `json:"diagnostics"`
}

// Reporter mendefinisikan interface abstraksi presenter dokumen laporan pemindaian.
type Reporter interface {
	Render(w io.Writer, result *ScanResult) error
}
