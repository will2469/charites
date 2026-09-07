package reporter

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules"
)

// MarkdownReporter memformat temuan analisis dalam format laporan audit Markdown terstruktur,
// mengadopsi standar hierarki file-first (File -> Violations).
type MarkdownReporter struct {
	rootDir   string
	timestamp time.Time
	reg       *rules.Registry
}

// MarkdownOption adalah functional option untuk mengonfigurasi MarkdownReporter.
type MarkdownOption func(*MarkdownReporter)

// WithRootDir menyetel direktori akar (root workspace) untuk resolusi path berkas relatif dan absolut.
func WithRootDir(rootDir string) MarkdownOption {
	return func(r *MarkdownReporter) {
		r.rootDir = rootDir
	}
}

// WithTimestamp menyetel timestamp deterministik laporan (berguna untuk testing).
func WithTimestamp(t time.Time) MarkdownOption {
	return func(r *MarkdownReporter) {
		r.timestamp = t
	}
}

// WithRegistry menyetel registry rule kustom untuk resolusi metadata (deskripsi, kategori, severity).
func WithRegistry(reg *rules.Registry) MarkdownOption {
	return func(r *MarkdownReporter) {
		r.reg = reg
	}
}

// NewMarkdownReporter membuat instans MarkdownReporter baru dengan opsi opsional.
func NewMarkdownReporter(opts ...MarkdownOption) *MarkdownReporter {
	r := &MarkdownReporter{
		rootDir: ".",
		reg:     rules.DefaultRegistry(),
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Render menulis laporan hasil pemindaian dalam format dokumen Markdown lengkap ke io.Writer.
func (r *MarkdownReporter) Render(w io.Writer, result *ScanResult) error {
	if result == nil {
		result = &ScanResult{}
	}

	timeStr, status := r.resolveTimeAndStatus(result)
	rootDir := r.resolveRootDir(result)

	sortedDiags := SortDiagnosticsCanonical(result.Diagnostics)
	fileGroups := GroupByFile(sortedDiags)
	NormalizeSummary(&result.Summary, fileGroups, result.Summary.ScannedFiles)

	attachedRules := result.AttachedRules
	if len(attachedRules) == 0 {
		attachedRules = r.buildDynamicRuleAuditInfo(result)
	}

	var sb strings.Builder
	r.renderHeader(&sb, timeStr, status)
	r.renderSummary(&sb, result.Summary, len(attachedRules), len(sortedDiags))
	r.renderDetailedInfo(&sb, attachedRules)
	r.renderResults(&sb, fileGroups, len(sortedDiags), rootDir)

	_, err := io.WriteString(w, sb.String())
	return err
}

func (r *MarkdownReporter) resolveTimeAndStatus(result *ScanResult) (string, string) {
	reportTime := result.Timestamp
	if reportTime.IsZero() {
		reportTime = r.timestamp
	}
	if reportTime.IsZero() {
		reportTime = time.Now().UTC()
	}
	timeStr := reportTime.Format("2006-01-02T15:04:05.000Z")

	status := "PASSED (Clean)"
	if len(result.Diagnostics) > 0 {
		status = "FAILED (Violations Found)"
	}
	return timeStr, status
}

func (r *MarkdownReporter) resolveRootDir(result *ScanResult) string {
	if result.RootDir != "" {
		return result.RootDir
	}
	if r.rootDir != "" {
		return r.rootDir
	}
	return "."
}

func (r *MarkdownReporter) renderHeader(sb *strings.Builder, timeStr, status string) {
	sb.WriteString("# Charites Frontend Static Analysis & UI Ergonomics Audit Report\n\n")
	sb.WriteString(fmt.Sprintf("**Timestamp:** %s  \n", timeStr))
	sb.WriteString(fmt.Sprintf("**Status:** %s  \n\n", status))
}

func (r *MarkdownReporter) renderSummary(sb *strings.Builder, s ScanSummary, attachedCount, issuesCount int) {
	sb.WriteString("## Summary\n\n")
	sb.WriteString("| Metric | Jumlah |\n")
	sb.WriteString("| :--- | :--- |\n")
	sb.WriteString(fmt.Sprintf("| Total Berkas | %d |\n", s.ScannedFiles))
	if s.FilesWithIssues > 0 || s.CleanFiles > 0 {
		sb.WriteString(fmt.Sprintf("| Berkas Bermasalah | %d |\n", s.FilesWithIssues))
		sb.WriteString(fmt.Sprintf("| Berkas Bersih | %d |\n", s.CleanFiles))
	}
	sb.WriteString(fmt.Sprintf("| Durasi Pemindaian | %dms |\n", s.DurationMS))
	sb.WriteString(fmt.Sprintf("| Rules Attached | %d |\n", attachedCount))
	sb.WriteString(fmt.Sprintf("| Total Issues | %d |\n", issuesCount))
	sb.WriteString(fmt.Sprintf("| Errors | %d |\n", s.ErrorCount))
	sb.WriteString(fmt.Sprintf("| Warnings | %d |\n", s.WarningCount))
	sb.WriteString(fmt.Sprintf("| Info | %d |\n\n", s.InfoCount))
}

func (r *MarkdownReporter) renderDetailedInfo(sb *strings.Builder, rules []RuleAuditInfo) {
	sb.WriteString("## Detailed Info\n\n")
	sb.WriteString("| ID | Category | Description | Issues Found | Status |\n")
	sb.WriteString("| :--- | :--- | :--- | :---: | :---: |\n")
	for _, ar := range rules {
		sb.WriteString(fmt.Sprintf("| [%s](https://github.com/will2469/charites/wiki/%s) | %s | %s | %d | %s |\n",
			ar.ID, ar.ID, ar.Category, ar.Description, ar.IssuesFound, ar.Status))
	}
	sb.WriteString("\n")
}

func (r *MarkdownReporter) renderResults(sb *strings.Builder, fileGroups []FileGroup, totalViolations int, rootDir string) {
	sb.WriteString("## Results by File\n\n")
	if len(fileGroups) == 0 {
		sb.WriteString("No known design token or ergonomics violations found\n")
		return
	}

	totalFiles := len(fileGroups)
	violationWord := "violations"
	if totalViolations == 1 {
		violationWord = "violation"
	}
	fileWord := "files"
	if totalFiles == 1 {
		fileWord = "file"
	}

	sb.WriteString(fmt.Sprintf("Found %d %s across %d %s:\n\n", totalViolations, violationWord, totalFiles, fileWord))

	for i, fg := range fileGroups {
		if i > 0 {
			sb.WriteString("---\n\n")
		}
		r.renderFileGroup(sb, fg, rootDir)
	}
}

func (r *MarkdownReporter) renderFileGroup(sb *strings.Builder, fg FileGroup, rootDir string) {
	posixRel, posixAbs := resolveAbsAndRelPath(fg.File, rootDir)
	firstLine := 1
	if len(fg.Diagnostics) > 0 && fg.Diagnostics[0].Line > 0 {
		firstLine = fg.Diagnostics[0].Line
	}
	firstLink := fmt.Sprintf("file://%s#L%d", posixAbs, firstLine)

	issueWord := "issues"
	if fg.TotalIssues == 1 {
		issueWord = "issue"
	}
	errorWord := "errors"
	if fg.ErrorCount == 1 {
		errorWord = "error"
	}
	warnWord := "warnings"
	if fg.WarningCount == 1 {
		warnWord = "warning"
	}

	sb.WriteString(fmt.Sprintf("### [%s](%s) (%d %s: %d %s, %d %s)\n\n",
		posixRel, firstLink, fg.TotalIssues, issueWord, fg.ErrorCount, errorWord, fg.WarningCount, warnWord))

	for _, d := range fg.Diagnostics {
		r.renderFileViolation(sb, d, posixAbs)
	}
}

func (r *MarkdownReporter) renderFileViolation(sb *strings.Builder, d ir.Diagnostic, posixAbs string) {
	link := fmt.Sprintf("file://%s#L%d", posixAbs, d.Line)
	posStr := fmt.Sprintf("L%d", d.Line)
	if d.Column > 0 {
		posStr = fmt.Sprintf("L%d:C%d", d.Line, d.Column)
	}

	sevTag := strings.ToUpper(string(d.Severity))
	if sevTag == "" {
		sevTag = "WARN"
	}

	sb.WriteString(fmt.Sprintf("- **[%s](%s)** • `[%s]` • [`%s`](https://github.com/will2469/charites/wiki/%s)\n",
		posStr, link, sevTag, d.Rule, d.Rule))

	if d.Message != "" {
		sb.WriteString(fmt.Sprintf("  - **Message:** %s\n", d.Message))
	}
	if d.Hint != "" {
		sb.WriteString(fmt.Sprintf("  - **Hint:** %s\n", d.Hint))
	}
	supp := SuppressionDirective(d.File, d.Rule)
	if supp != "" {
		sb.WriteString(fmt.Sprintf("  - **Suppression:** `%s`\n", supp))
	}
	sb.WriteString("\n")
}

func resolveAbsAndRelPath(filePath, rootDir string) (posixRel, posixAbs string) {
	absPath := filePath
	if !filepath.IsAbs(absPath) {
		if _, err := os.Stat(filePath); err == nil {
			if p, err := filepath.Abs(filePath); err == nil {
				absPath = p
			} else {
				absPath = filepath.Join(rootDir, filePath)
			}
		} else {
			absPath = filepath.Join(rootDir, filePath)
		}
	}

	relPath := filePath
	if filepath.IsAbs(relPath) {
		if rel, err := filepath.Rel(rootDir, relPath); err == nil {
			relPath = rel
		}
	}

	posixRel = normalizePOSIXPath(relPath)
	posixAbs = "/" + strings.TrimPrefix(normalizePOSIXPath(absPath), "/")
	return posixRel, posixAbs
}

func (r *MarkdownReporter) buildDynamicRuleAuditInfo(result *ScanResult) []RuleAuditInfo {
	counts := make(map[string]int)
	for _, d := range result.Diagnostics {
		counts[d.Rule]++
	}

	var infos []RuleAuditInfo

	if r.reg != nil && r.reg.Count() > 0 {
		for _, ruleObj := range r.reg.All() {
			c := counts[ruleObj.ID()]
			st := "PASS"
			if c > 0 {
				st = "FAILED"
			}
			infos = append(infos, RuleAuditInfo{
				ID:          ruleObj.ID(),
				Category:    ruleObj.Category(),
				Description: ruleObj.Description(),
				Severity:    string(ruleObj.DefaultSeverity()),
				IssuesFound: c,
				Status:      st,
			})
		}
	} else {
		for ruleID, c := range counts {
			cat := ""
			if idx := strings.IndexByte(ruleID, '.'); idx != -1 {
				cat = ruleID[:idx]
			}
			infos = append(infos, RuleAuditInfo{
				ID:          ruleID,
				Category:    cat,
				Description: "",
				Severity:    "warn",
				IssuesFound: c,
				Status:      "FAILED",
			})
		}
	}

	sort.Slice(infos, func(i, j int) bool {
		return infos[i].ID < infos[j].ID
	})

	return infos
}
