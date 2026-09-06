package reporter

import (
	"encoding/json"
	"io"
	"path/filepath"
	"strings"
)

// JSONReporter mencetak seluruh dokumen laporan dalam format JSON terformat.
type JSONReporter struct{}

// NewJSONReporter membuat instans JSONReporter baru.
func NewJSONReporter() *JSONReporter {
	return &JSONReporter{}
}

type jsonDiagnostic struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Rule     string `json:"rule"`
	Category string `json:"category"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	Hint     string `json:"hint,omitempty"`
	DocURL   string `json:"doc_url"`
}

type jsonDocument struct {
	Version     string           `json:"version"`
	Summary     ScanSummary      `json:"summary"`
	Diagnostics []jsonDiagnostic `json:"diagnostics"`
}

// Render menulis laporan hasil pemindaian ke io.Writer dalam format dokumen JSON tunggal lengkap.
func (r *JSONReporter) Render(w io.Writer, result *ScanResult) error {
	if result == nil {
		result = &ScanResult{}
	}

	doc := jsonDocument{
		Version:     result.Version,
		Summary:     result.Summary,
		Diagnostics: make([]jsonDiagnostic, 0, len(result.Diagnostics)),
	}

	if doc.Version == "" {
		doc.Version = "1.0.0"
	}

	for _, d := range result.Diagnostics {
		cat := ""
		if idx := strings.IndexByte(d.Rule, '.'); idx != -1 {
			cat = d.Rule[:idx]
		}

		doc.Diagnostics = append(doc.Diagnostics, jsonDiagnostic{
			File:     filepath.ToSlash(d.File),
			Line:     d.Line,
			Column:   d.Column,
			Rule:     d.Rule,
			Category: cat,
			Severity: string(d.Severity),
			Message:  d.Message,
			Hint:     d.Hint,
			DocURL:   "https://github.com/will2469/charites/wiki/" + d.Rule,
		})
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(doc)
}
