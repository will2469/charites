package ir

import (
	"cmp"
	"slices"
	"strings"
)

// Severity merepresentasikan tingkat keparahan suatu temuan diagnostik.
type Severity string

const (
	// SeverityError merepresentasikan pelanggaran fatal (exit code 1).
	SeverityError Severity = "error"
	// SeverityWarn merepresentasikan peringatan perbaikan (exit code 1).
	SeverityWarn Severity = "warn"
	// SeverityInfo merepresentasikan informasi atau saran perbaikan (exit code 0 jika hanya info).
	SeverityInfo Severity = "info"
)

// Diagnostic merepresentasikan temuan diagnostik tunggal dari hasil evaluasi rule.
// Direpresentasikan sebagai struktur flat JSON deterministik.
type Diagnostic struct {
	File     string   `json:"file"`           // Path relatif berkas sumber (normalisasi POSIX '/')
	Line     int      `json:"line"`           // Baris lokasi temuan (1-indexed)
	Column   int      `json:"column"`         // Kolom lokasi temuan (1-indexed)
	Rule     string   `json:"rule"`           // Charites Rule ID (contoh: "theme.hardcode-opacity-color")
	Severity Severity `json:"severity"`       // "error" | "warn" | "info"
	Message  string   `json:"message"`        // Deskripsi pelanggaran ringkas
	Hint     string   `json:"hint,omitempty"` // Rekomendasi tindakan remediasi (opsional)
}

// severityRank memetakan tingkat keparahan ke bobot integer untuk perbandingan deterministik.
func severityRank(s Severity) int {
	switch s {
	case SeverityError:
		return 0
	case SeverityWarn:
		return 1
	case SeverityInfo:
		return 2
	default:
		return 3
	}
}

// CompareDiagnostics membandingkan dua Diagnostic menggunakan 7-level Canonical Diagnostic Total Ordering (DiagnosticOrderKey):
// File -> Line -> Column -> Rule -> Severity -> Message -> Hint.
// Mengembalikan -1 jika a < b, 1 jika a > b, dan 0 jika a == b (identik mutlak).
func CompareDiagnostics(a, b Diagnostic) int {
	if c := strings.Compare(a.File, b.File); c != 0 {
		return c
	}
	if c := cmp.Compare(a.Line, b.Line); c != 0 {
		return c
	}
	if c := cmp.Compare(a.Column, b.Column); c != 0 {
		return c
	}
	if c := strings.Compare(a.Rule, b.Rule); c != 0 {
		return c
	}
	if c := cmp.Compare(severityRank(a.Severity), severityRank(b.Severity)); c != 0 {
		return c
	}
	if c := strings.Compare(a.Message, b.Message); c != 0 {
		return c
	}
	return strings.Compare(a.Hint, b.Hint)
}

// SortDiagnostics mengurutkan slice diagnosis secara in-place menggunakan 7-level total order
// dan memangkas duplikat identik (idempotent deduplication).
func SortDiagnostics(diags []Diagnostic) []Diagnostic {
	slices.SortFunc(diags, CompareDiagnostics)
	return slices.CompactFunc(diags, func(a, b Diagnostic) bool {
		return CompareDiagnostics(a, b) == 0
	})
}
