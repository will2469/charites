package reporter

import (
	"encoding/json"
	"io"
	"strings"
)

// JSONReporter mencetak seluruh dokumen laporan dalam format JSON terformat.
type JSONReporter struct{}

// NewJSONReporter membuat instans JSONReporter baru.
func NewJSONReporter() *JSONReporter {
	return &JSONReporter{}
}

// Render menulis laporan hasil pemindaian ke io.Writer dalam format dokumen JSON tunggal lengkap.
func (r *JSONReporter) Render(w io.Writer, result *ScanResult) error {
	if result == nil {
		result = &ScanResult{}
	}

	version := result.Version
	if version == "" {
		version = DefaultReportVersion
	}

	canonicalDiags := cloneAndSortDiagnostics(result.Diagnostics)
	fileGroups := GroupByFile(canonicalDiags)

	doc := jsonDocument{
		Version:     version,
		Summary:     result.Summary,
		Files:       make([]jsonFileGroup, 0, len(fileGroups)),
		Diagnostics: make([]jsonDiagnostic, 0, len(canonicalDiags)),
	}

	for _, fg := range fileGroups {
		violations := make([]jsonViolationItem, 0, len(fg.Diagnostics))
		for _, d := range fg.Diagnostics {
			cat := ""
			if idx := strings.IndexByte(d.Rule, '.'); idx != -1 {
				cat = d.Rule[:idx]
			}

			violations = append(violations, jsonViolationItem{
				Line:        d.Line,
				Column:      d.Column,
				Rule:        d.Rule,
				Category:    cat,
				Severity:    string(d.Severity),
				Message:     d.Message,
				Hint:        d.Hint,
				DocURL:      "https://github.com/will2469/charites/wiki/" + d.Rule,
				Suppression: SuppressionDirective(fg.File, d.Rule),
			})
		}

		doc.Files = append(doc.Files, jsonFileGroup{
			File:         fg.File,
			ErrorCount:   fg.ErrorCount,
			WarningCount: fg.WarningCount,
			InfoCount:    fg.InfoCount,
			TotalIssues:  fg.TotalIssues,
			Violations:   violations,
		})
	}

	for _, d := range canonicalDiags {
		cat := ""
		if idx := strings.IndexByte(d.Rule, '.'); idx != -1 {
			cat = d.Rule[:idx]
		}

		doc.Diagnostics = append(doc.Diagnostics, jsonDiagnostic{
			File:     d.File,
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
