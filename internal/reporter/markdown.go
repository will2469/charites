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
// mengadopsi standar tata letak dan hierarki dari Argus audit engine.
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

	attachedRules := result.AttachedRules
	if len(attachedRules) == 0 {
		attachedRules = r.buildDynamicRuleAuditInfo(result)
	}

	var sb strings.Builder
	r.renderHeader(&sb, timeStr, status)
	r.renderSummary(&sb, result.Summary, len(attachedRules), len(result.Diagnostics))
	r.renderDetailedInfo(&sb, attachedRules)
	r.renderResults(&sb, result.Diagnostics, rootDir)

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

func (r *MarkdownReporter) renderResults(sb *strings.Builder, diags []ir.Diagnostic, rootDir string) {
	sb.WriteString("## Result\n\n")
	if len(diags) == 0 {
		sb.WriteString("No known design token or ergonomics violations found\n")
		return
	}

	sb.WriteString(fmt.Sprintf("Found %d violations across scanned components:\n\n", len(diags)))

	grouped := make(map[string][]ir.Diagnostic)
	for _, d := range diags {
		grouped[d.Rule] = append(grouped[d.Rule], d)
	}

	sortedRules := make([]string, 0, len(grouped))
	for ruleID := range grouped {
		sortedRules = append(sortedRules, ruleID)
	}
	sort.Strings(sortedRules)

	for _, ruleID := range sortedRules {
		r.renderRuleGroup(sb, ruleID, grouped[ruleID], rootDir)
	}
}

func (r *MarkdownReporter) renderRuleGroup(sb *strings.Builder, ruleID string, diags []ir.Diagnostic, rootDir string) {
	firstDiag := diags[0]
	category, desc, severity := r.resolveRuleMeta(ruleID, firstDiag)

	sb.WriteString(fmt.Sprintf("### %s\n\n", ruleID))
	sb.WriteString(fmt.Sprintf("- **Severity:** %s\n", severity))
	sb.WriteString(fmt.Sprintf("- **Category:** %s\n", category))
	if desc != "" {
		sb.WriteString(fmt.Sprintf("- **Description:** %s\n", desc))
	}
	sb.WriteString(fmt.Sprintf("- **Wiki:** [%s Documentation](https://github.com/will2469/charites/wiki/%s)\n", ruleID, ruleID))

	hasAstro, hasJSX, allSameMessage := inspectDiagnostics(diags)
	sb.WriteString(determineSuppressionHint(ruleID, hasAstro, hasJSX))

	if allSameMessage && firstDiag.Message != "" {
		sb.WriteString(fmt.Sprintf("- **Message:** %s\n", firstDiag.Message))
	}
	sb.WriteString("\n")

	for _, d := range diags {
		renderOccurrence(sb, d, rootDir, !allSameMessage)
	}
	sb.WriteString("\n")
}

func (r *MarkdownReporter) resolveRuleMeta(ruleID string, fallback ir.Diagnostic) (category, desc, severity string) {
	severity = string(fallback.Severity)
	if r.reg != nil {
		if ruleObj, ok := r.reg.Get(ruleID); ok {
			category = ruleObj.Category()
			desc = ruleObj.Description()
			if severity == "" {
				severity = string(ruleObj.DefaultSeverity())
			}
		}
	}
	if category == "" {
		if idx := strings.IndexByte(ruleID, '.'); idx != -1 {
			category = ruleID[:idx]
		} else {
			category = "general"
		}
	}
	return category, desc, severity
}

func inspectDiagnostics(diags []ir.Diagnostic) (hasAstro, hasJSX, allSameMessage bool) {
	allSameMessage = true
	firstMsg := diags[0].Message

	for _, d := range diags {
		ext := strings.ToLower(filepath.Ext(d.File))
		if ext == ".astro" {
			hasAstro = true
		} else {
			hasJSX = true
		}
		if d.Message != firstMsg {
			allSameMessage = false
		}
	}
	return hasAstro, hasJSX, allSameMessage
}

func determineSuppressionHint(ruleID string, hasAstro, hasJSX bool) string {
	switch {
	case hasAstro && hasJSX:
		return fmt.Sprintf("- **Suppression:** `<!-- charites:ignore %s <reason> -->` (Astro) or `// charites:ignore %s <reason>` (TSX/JSX)\n", ruleID, ruleID)
	case hasAstro:
		return fmt.Sprintf("- **Suppression:** `<!-- charites:ignore %s <reason> -->`\n", ruleID)
	default:
		return fmt.Sprintf("- **Suppression:** `// charites:ignore %s <reason>`\n", ruleID)
	}
}

func renderOccurrence(sb *strings.Builder, d ir.Diagnostic, rootDir string, printMessage bool) {
	absPath := d.File
	if !filepath.IsAbs(absPath) {
		if _, err := os.Stat(d.File); err == nil {
			if p, err := filepath.Abs(d.File); err == nil {
				absPath = p
			} else {
				absPath = filepath.Join(rootDir, d.File)
			}
		} else {
			absPath = filepath.Join(rootDir, d.File)
		}
	}

	relPath := d.File
	if filepath.IsAbs(relPath) {
		if rel, err := filepath.Rel(rootDir, relPath); err == nil {
			relPath = rel
		}
	}

	posixRel := filepath.ToSlash(relPath)
	posixAbs := "/" + strings.TrimPrefix(filepath.ToSlash(absPath), "/")
	link := fmt.Sprintf("file://%s#L%d", posixAbs, d.Line)

	if d.Column > 0 {
		sb.WriteString(fmt.Sprintf("- **[%s:%d:%d](%s)**\n", posixRel, d.Line, d.Column, link))
	} else {
		sb.WriteString(fmt.Sprintf("- **[%s:%d](%s)**\n", posixRel, d.Line, link))
	}

	if printMessage && d.Message != "" {
		sb.WriteString(fmt.Sprintf("  - *Message:* %s\n", d.Message))
	}
	if d.Hint != "" {
		sb.WriteString(fmt.Sprintf("  - *Hint:* %s\n", d.Hint))
	}
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
