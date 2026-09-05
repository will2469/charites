package reporter

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// InlineReporter memformat temuan analisis dalam bentuk teks ANSI untuk terminal konsol.
type InlineReporter struct {
	colorMode ColorMode
}

// NewInlineReporter membuat instans InlineReporter dengan ColorMode tertentu.
func NewInlineReporter(mode ColorMode) *InlineReporter {
	return &InlineReporter{
		colorMode: mode,
	}
}

// Render menulis laporan hasil pemindaian dalam format inline ANSI ke io.Writer.
func (r *InlineReporter) Render(w io.Writer, result *ScanResult) error {
	if result == nil {
		result = &ScanResult{}
	}

	useColor := r.colorMode == ColorAlways || (r.colorMode == ColorAuto && IsTerminal(w))

	// 1. Render setiap temuan diagnostik
	for i, d := range result.Diagnostics {
		if i > 0 {
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}

		badge := r.formatBadge(d.Severity, useColor)
		posixPath := filepath.ToSlash(d.File)
		header := fmt.Sprintf("%s %s:%d:%d [%s]", badge, posixPath, d.Line, d.Column, d.Rule)
		if _, err := fmt.Fprintln(w, header); err != nil {
			return err
		}

		msgLine := fmt.Sprintf("  %s", d.Message)
		if _, err := fmt.Fprintln(w, msgLine); err != nil {
			return err
		}

		if d.Hint != "" {
			hintPrefix := "Hint:"
			if useColor {
				hintPrefix = "\033[2mHint:\033[0m"
			}
			hintLine := fmt.Sprintf("  %s %s", hintPrefix, d.Hint)
			if _, err := fmt.Fprintln(w, hintLine); err != nil {
				return err
			}
		}
	}

	// 2. Jika ada diagnostik, beri baris pemisah sebelum ringkasan footer
	if len(result.Diagnostics) > 0 {
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}

	// 3. Render Ringkasan Footer
	totalProblems := result.Summary.ErrorCount + result.Summary.WarningCount + result.Summary.InfoCount

	problemsWord := "problems"
	if totalProblems == 1 {
		problemsWord = "problem"
	}

	errorsWord := "errors"
	if result.Summary.ErrorCount == 1 {
		errorsWord = "error"
	}

	warningsWord := "warnings"
	if result.Summary.WarningCount == 1 {
		warningsWord = "warning"
	}

	summaryLine := fmt.Sprintf(" %d %s found (%d %s, %d %s)",
		totalProblems,
		problemsWord,
		result.Summary.ErrorCount,
		errorsWord,
		result.Summary.WarningCount,
		warningsWord,
	)
	if _, err := fmt.Fprintln(w, summaryLine); err != nil {
		return err
	}

	filesWord := "files"
	if result.Summary.ScannedFiles == 1 {
		filesWord = "file"
	}

	scannedLine := fmt.Sprintf("  Scanned %d %s in %dms.",
		result.Summary.ScannedFiles,
		filesWord,
		result.Summary.DurationMS,
	)
	if _, err := fmt.Fprintln(w, scannedLine); err != nil {
		return err
	}

	return nil
}

func (r *InlineReporter) formatBadge(sev ir.Severity, useColor bool) string {
	s := strings.ToLower(string(sev))
	switch s {
	case "error":
		if useColor {
			return "\033[1;31m[ERROR]\033[0m"
		}
		return "[ERROR]"
	case "warn", "warning":
		if useColor {
			return "\033[1;33m[WARN]\033[0m"
		}
		return "[WARN]"
	case "info":
		if useColor {
			return "\033[1;36m[INFO]\033[0m"
		}
		return "[INFO]"
	default:
		if useColor {
			return "\033[1;31m[ERROR]\033[0m"
		}
		return "[ERROR]"
	}
}
