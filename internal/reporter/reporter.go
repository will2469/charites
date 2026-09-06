package reporter

import (
	"io"
	"time"

	"github.com/will2469/charites/internal/ir"
)

// RuleAuditInfo merepresentasikan status audit per-rule untuk pelaporan detil.
type RuleAuditInfo struct {
	ID          string `json:"id"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	IssuesFound int    `json:"issues_found"`
	Status      string `json:"status"` // "PASS" | "FAILED"
}

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
	Version       string          `json:"version"`
	Timestamp     time.Time       `json:"timestamp,omitempty"`
	RootDir       string          `json:"root_dir,omitempty"`
	Summary       ScanSummary     `json:"summary"`
	Diagnostics   []ir.Diagnostic `json:"diagnostics"`
	AttachedRules []RuleAuditInfo `json:"attached_rules,omitempty"`
}

// Reporter mendefinisikan interface abstraksi presenter dokumen laporan pemindaian.
type Reporter interface {
	Render(w io.Writer, result *ScanResult) error
}
