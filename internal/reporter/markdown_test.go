package reporter_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/reporter"
	"github.com/will2469/charites/internal/rules"
)

func TestMarkdownReporter_Clean(t *testing.T) {
	fixedTime := time.Date(2026, 9, 6, 12, 0, 0, 0, time.UTC)
	clean := &reporter.ScanResult{
		Version:   "1.0.0",
		Timestamp: fixedTime,
		Summary: reporter.ScanSummary{
			ScannedFiles: 25,
			DurationMS:   14,
			ErrorCount:   0,
			WarningCount: 0,
			InfoCount:    0,
			Passed:       true,
		},
		Diagnostics: []ir.Diagnostic{},
	}

	rep := reporter.NewMarkdownReporter(
		reporter.WithRootDir("/workspace/project"),
		reporter.WithTimestamp(fixedTime),
	)

	var buf bytes.Buffer
	if err := rep.Render(&buf, clean); err != nil {
		t.Fatalf("render failed: %v", err)
	}

	out := buf.String()

	if !strings.Contains(out, "# Charites Frontend Static Analysis & UI Ergonomics Audit Report") {
		t.Errorf("expected header title, got:\n%s", out)
	}
	if !strings.Contains(out, "**Timestamp:** 2026-09-06T12:00:00.000Z") {
		t.Errorf("expected timestamp, got:\n%s", out)
	}
	if !strings.Contains(out, "**Status:** PASSED (Clean)") {
		t.Errorf("expected PASSED status, got:\n%s", out)
	}
	if !strings.Contains(out, "## Summary") {
		t.Errorf("expected Summary section")
	}
	if !strings.Contains(out, "| Total Berkas | 25 |") {
		t.Errorf("expected Total Berkas metric in summary")
	}
	if !strings.Contains(out, "## Detailed Info") {
		t.Errorf("expected Detailed Info section")
	}
	if !strings.Contains(out, "## Result") {
		t.Errorf("expected Result section")
	}
	if !strings.Contains(out, "No known design token or ergonomics violations found") {
		t.Errorf("expected clean message in Result")
	}
}

func TestMarkdownReporter_WithViolations(t *testing.T) {
	fixedTime := time.Date(2026, 9, 6, 12, 0, 0, 0, time.UTC)
	violations := &reporter.ScanResult{
		Version:   "1.0.0",
		Timestamp: fixedTime,
		RootDir:   "/workspace/project",
		Summary: reporter.ScanSummary{
			ScannedFiles: 10,
			DurationMS:   22,
			ErrorCount:   1,
			WarningCount: 2,
			InfoCount:    0,
			Passed:       false,
		},
		Diagnostics: []ir.Diagnostic{
			{
				File:     "src/pages/index.astro",
				Line:     15,
				Column:   7,
				Rule:     "theme.hardcode-opacity-color",
				Severity: ir.SeverityError,
				Message:  `Hardcode opacity color: "bg-primary/10"`,
				Hint:     `Use semantic token "primary-light".`,
			},
			{
				File:     "src/components/Card.tsx",
				Line:     42,
				Column:   12,
				Rule:     "theme.hardcode-size",
				Severity: ir.SeverityWarn,
				Message:  `Non-standard fractional scale: "p-3.25"`,
				Hint:     `Use standard integer step p-3 or p-4.`,
			},
			{
				File:     "src/components/Button.astro",
				Line:     8,
				Column:   4,
				Rule:     "theme.hardcode-size",
				Severity: ir.SeverityWarn,
				Message:  `Hardcoded size scalar: "w-[120px]"`,
				Hint:     `Use standard token w-28 or w-32.`,
			},
		},
	}

	rep := reporter.NewMarkdownReporter(
		reporter.WithRootDir("/workspace/project"),
		reporter.WithTimestamp(fixedTime),
	)

	var buf bytes.Buffer
	if err := rep.Render(&buf, violations); err != nil {
		t.Fatalf("render failed: %v", err)
	}

	out := buf.String()

	if !strings.Contains(out, "**Status:** FAILED (Violations Found)") {
		t.Errorf("expected FAILED status")
	}
	if !strings.Contains(out, "Found 3 violations across scanned components:") {
		t.Errorf("expected 3 violations header")
	}

	// Cek grouping per rule
	if !strings.Contains(out, "### theme.hardcode-opacity-color") {
		t.Errorf("expected theme.hardcode-opacity-color header")
	}
	if !strings.Contains(out, "### theme.hardcode-size") {
		t.Errorf("expected theme.hardcode-size header")
	}

	// Cek metadata rule
	if !strings.Contains(out, "[theme.hardcode-opacity-color Documentation](https://github.com/will2469/charites/wiki/theme.hardcode-opacity-color)") {
		t.Errorf("expected wiki link for theme.hardcode-opacity-color")
	}
	if !strings.Contains(out, "| [theme.hardcode-opacity-color](https://github.com/will2469/charites/wiki/theme.hardcode-opacity-color) |") {
		t.Errorf("expected linked rule ID in Detailed Info table for theme.hardcode-opacity-color")
	}

	// Cek supresi cerdas:
	// theme.hardcode-opacity-color hanya ada di .astro -> <!-- charites:ignore -->
	if !strings.Contains(out, "`<!-- charites:ignore theme.hardcode-opacity-color <reason> -->`") {
		t.Errorf("expected astro suppression directive for theme.hardcode-opacity-color")
	}
	// theme.hardcode-size ada di .tsx dan .astro -> both
	if !strings.Contains(out, "`<!-- charites:ignore theme.hardcode-size <reason> -->` (Astro) or `// charites:ignore theme.hardcode-size <reason>` (TSX/JSX)") {
		t.Errorf("expected dual suppression directive for theme.hardcode-size")
	}

	// Cek file links dan hash baris
	if !strings.Contains(out, "**[src/pages/index.astro:15:7](file:///workspace/project/src/pages/index.astro#L15)**") {
		t.Errorf("expected file link for index.astro with line 15")
	}
	if !strings.Contains(out, "**[src/components/Card.tsx:42:12](file:///workspace/project/src/components/Card.tsx#L42)**") {
		t.Errorf("expected file link for Card.tsx with line 42")
	}

	// Cek hint
	if !strings.Contains(out, "*Hint:* Use semantic token \"primary-light\".") {
		t.Errorf("expected hint for opacity violation")
	}
}

func TestMarkdownReporter_Determinism(t *testing.T) {
	fixedTime := time.Date(2026, 9, 6, 12, 0, 0, 0, time.UTC)
	result := sampleViolationsResult()
	result.Timestamp = fixedTime

	rep := reporter.NewMarkdownReporter(
		reporter.WithRootDir("/workspace/project"),
		reporter.WithTimestamp(fixedTime),
	)

	var buf1, buf2 bytes.Buffer
	if err := rep.Render(&buf1, result); err != nil {
		t.Fatalf("render 1 failed: %v", err)
	}
	if err := rep.Render(&buf2, result); err != nil {
		t.Fatalf("render 2 failed: %v", err)
	}

	if !bytes.Equal(buf1.Bytes(), buf2.Bytes()) {
		t.Fatalf("Markdown reporter output is not byte-for-byte identical across runs")
	}
}

func TestMarkdownReporter_CustomAttachedRules(t *testing.T) {
	result := &reporter.ScanResult{
		Version: "1.0.0",
		Summary: reporter.ScanSummary{ScannedFiles: 5, DurationMS: 8, Passed: true},
		AttachedRules: []reporter.RuleAuditInfo{
			{
				ID:          "theme.hardcode-color",
				Category:    "theme",
				Description: "Detects hardcoded color",
				Severity:    "warn",
				IssuesFound: 0,
				Status:      "PASS",
			},
		},
	}

	var buf bytes.Buffer
	rep := reporter.NewMarkdownReporter()
	if err := rep.Render(&buf, result); err != nil {
		t.Fatalf("render failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "| Rules Attached | 1 |") {
		t.Errorf("expected 1 attached rule in summary, got:\n%s", out)
	}
	if !strings.Contains(out, "| [theme.hardcode-color](https://github.com/will2469/charites/wiki/theme.hardcode-color) | theme | Detects hardcoded color | 0 | PASS |") {
		t.Errorf("expected detailed row for theme.hardcode-color, got:\n%s", out)
	}
}

func TestMarkdownReporter_EmptyRegistryFallback(t *testing.T) {
	emptyReg := rules.NewRegistry()
	fixedTime := time.Date(2026, 9, 6, 12, 0, 0, 0, time.UTC)
	result := &reporter.ScanResult{
		Timestamp: fixedTime,
		Diagnostics: []ir.Diagnostic{
			{
				File:     "test.tsx",
				Line:     1,
				Rule:     "custom.dummy-rule",
				Severity: ir.SeverityWarn,
				Message:  "Dummy violation",
			},
		},
	}

	rep := reporter.NewMarkdownReporter(
		reporter.WithRegistry(emptyReg),
		reporter.WithTimestamp(fixedTime),
	)

	var buf bytes.Buffer
	if err := rep.Render(&buf, result); err != nil {
		t.Fatalf("render failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "custom.dummy-rule") {
		t.Errorf("expected custom.dummy-rule in output")
	}
	if !strings.Contains(out, "| [custom.dummy-rule](https://github.com/will2469/charites/wiki/custom.dummy-rule) | custom |  | 1 | FAILED |") {
		t.Errorf("expected custom.dummy-rule row in detailed info")
	}
}
