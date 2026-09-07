package reporter_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/reporter"
)

var testFixedTime = time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)

func sampleCleanResult() *reporter.ScanResult {
	return &reporter.ScanResult{
		Version:   reporter.DefaultReportVersion,
		Timestamp: testFixedTime,
		Summary: reporter.ScanSummary{
			ScannedFiles: 28,
			DurationMS:   12,
			ErrorCount:   0,
			WarningCount: 0,
			InfoCount:    0,
			Passed:       true,
		},
		Diagnostics: []ir.Diagnostic{},
	}
}

func sampleViolationsResult() *reporter.ScanResult {
	return &reporter.ScanResult{
		Version:   reporter.DefaultReportVersion,
		Timestamp: testFixedTime,
		Summary: reporter.ScanSummary{
			ScannedFiles: 28,
			DurationMS:   18,
			ErrorCount:   1,
			WarningCount: 1,
			InfoCount:    0,
			Passed:       false,
		},
		Diagnostics: []ir.Diagnostic{
			{
				File:     "src/pages/index.astro",
				Line:     14,
				Column:   8,
				Rule:     "theme.hardcode-opacity-color",
				Severity: ir.SeverityError,
				Message:  `Hardcode opacity color: "bg-primary/10"`,
				Hint:     `Use semantic token "primary-light".`,
			},
			{
				File:     "src/components/Card.tsx",
				Line:     42,
				Column:   12,
				Rule:     "theme.hardcode-color",
				Severity: ir.SeverityWarn,
				Message:  `Hardcode hex color: "#2563eb"`,
				Hint:     `Use semantic token "bg-primary".`,
			},
		},
	}
}

func readGolden(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "..", "tests", "golden", "reporters", name)
	data, err := os.ReadFile(filepath.Clean(path)) //nolint:gosec // controlled test fixture path
	if err != nil {
		t.Fatalf("failed to read golden file %s: %v", path, err)
	}
	return data
}

func TestReporter_Determinism(t *testing.T) {
	result := sampleViolationsResult()

	// 1. Determinisme JSON
	var bufJSON1, bufJSON2 bytes.Buffer
	jsonRep := reporter.NewJSONReporter()
	if err := jsonRep.Render(&bufJSON1, result); err != nil {
		t.Fatalf("json render 1 failed: %v", err)
	}
	if err := jsonRep.Render(&bufJSON2, result); err != nil {
		t.Fatalf("json render 2 failed: %v", err)
	}
	if !bytes.Equal(bufJSON1.Bytes(), bufJSON2.Bytes()) {
		t.Fatalf("JSON reporter output is not byte-for-byte identical")
	}

	// 2. Determinisme Inline
	var bufInline1, bufInline2 bytes.Buffer
	inlineRep := reporter.NewInlineReporter(reporter.ColorNever)
	if err := inlineRep.Render(&bufInline1, result); err != nil {
		t.Fatalf("inline render 1 failed: %v", err)
	}
	if err := inlineRep.Render(&bufInline2, result); err != nil {
		t.Fatalf("inline render 2 failed: %v", err)
	}
	if !bytes.Equal(bufInline1.Bytes(), bufInline2.Bytes()) {
		t.Fatalf("Inline reporter output is not byte-for-byte identical")
	}
}

func TestInlineReporter_GoldenSnapshots(t *testing.T) {
	t.Run("clean scan", func(t *testing.T) {
		var buf bytes.Buffer
		rep := reporter.NewInlineReporter(reporter.ColorNever)
		if err := rep.Render(&buf, sampleCleanResult()); err != nil {
			t.Fatalf("render failed: %v", err)
		}
		expected := readGolden(t, "inline_clean.golden")
		if !bytes.Equal(buf.Bytes(), expected) {
			t.Errorf("clean output mismatch.\nGot:\n%s\nExpected:\n%s", buf.String(), string(expected))
		}
	})

	t.Run("violations no color", func(t *testing.T) {
		var buf bytes.Buffer
		rep := reporter.NewInlineReporter(reporter.ColorNever)
		if err := rep.Render(&buf, sampleViolationsResult()); err != nil {
			t.Fatalf("render failed: %v", err)
		}
		expected := readGolden(t, "inline_no_color.golden")
		if !bytes.Equal(buf.Bytes(), expected) {
			t.Errorf("violations no color mismatch.\nGot:\n%s\nExpected:\n%s", buf.String(), string(expected))
		}
	})
}

func TestInlineReporter_SingularGrammarAndSeverities(t *testing.T) {
	result := &reporter.ScanResult{
		Version: "1.0.0",
		Summary: reporter.ScanSummary{
			ScannedFiles: 1,
			DurationMS:   5,
			ErrorCount:   0,
			WarningCount: 0,
			InfoCount:    1,
			Passed:       true,
		},
		Diagnostics: []ir.Diagnostic{
			{
				File:     "src/index.tsx",
				Line:     1,
				Column:   1,
				Rule:     "a11y.alt",
				Severity: ir.SeverityInfo,
				Message:  "Consider adding aria-label",
			},
		},
	}

	var buf bytes.Buffer
	rep := reporter.NewInlineReporter(reporter.ColorNever)
	if err := rep.Render(&buf, result); err != nil {
		t.Fatalf("render failed: %v", err)
	}

	out := buf.String()
	expectedBadge := "[INFO] src/index.tsx:1:1 [a11y.alt]\n  Consider adding aria-label\n\n 1 problem found (0 errors, 0 warnings)\n  Scanned 1 file in 5ms.\n"
	if out != expectedBadge {
		t.Errorf("singular output mismatch.\nGot:\n%s\nExpected:\n%s", out, expectedBadge)
	}

	// Test ColorAlways
	var bufColor bytes.Buffer
	repColor := reporter.NewInlineReporter(reporter.ColorAlways)
	if err := repColor.Render(&bufColor, result); err != nil {
		t.Fatalf("render failed: %v", err)
	}
	if !bytes.Contains(bufColor.Bytes(), []byte("\033[1;36m[INFO]\033[0m")) {
		t.Errorf("expected colored info badge in output: %s", bufColor.String())
	}
}

func TestJSONReporter_GoldenSnapshots(t *testing.T) {
	t.Run("clean json", func(t *testing.T) {
		var buf bytes.Buffer
		rep := reporter.NewJSONReporter()
		if err := rep.Render(&buf, sampleCleanResult()); err != nil {
			t.Fatalf("render failed: %v", err)
		}
		expected := readGolden(t, "json_clean.golden")
		if !bytes.Equal(buf.Bytes(), expected) {
			t.Errorf("clean json mismatch.\nGot:\n%s\nExpected:\n%s", buf.String(), string(expected))
		}
	})

	t.Run("violations json", func(t *testing.T) {
		var buf bytes.Buffer
		rep := reporter.NewJSONReporter()
		if err := rep.Render(&buf, sampleViolationsResult()); err != nil {
			t.Fatalf("render failed: %v", err)
		}
		expected := readGolden(t, "json_violations.golden")
		if !bytes.Equal(buf.Bytes(), expected) {
			t.Errorf("violations json mismatch.\nGot:\n%s\nExpected:\n%s", buf.String(), string(expected))
		}
	})
}

func TestColorResolution(t *testing.T) {
	var buf bytes.Buffer

	// Flag no-color override
	if mode := reporter.ResolveColorMode(true, &buf); mode != reporter.ColorNever {
		t.Errorf("expected ColorNever for noColor=true, got %v", mode)
	}

	// Environment variable NO_COLOR
	t.Setenv("NO_COLOR", "1")
	if mode := reporter.ResolveColorMode(false, &buf); mode != reporter.ColorNever {
		t.Errorf("expected ColorNever for NO_COLOR=1, got %v", mode)
	}

	// Non-terminal buffer with NO_COLOR unset
	t.Setenv("NO_COLOR", "")
	if mode := reporter.ResolveColorMode(false, &buf); mode != reporter.ColorNever {
		t.Errorf("expected ColorNever for non-terminal writer, got %v", mode)
	}

	if reporter.IsTerminal(&buf) {
		t.Errorf("bytes.Buffer should not be a terminal")
	}

	// Regular file check (not a terminal character device)
	tmpFile, err := os.CreateTemp(t.TempDir(), "test-file-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer func() { _ = tmpFile.Close() }()

	if reporter.IsTerminal(tmpFile) {
		t.Errorf("regular file should not be reported as terminal")
	}
}

type errWriter struct {
	failAfter int
	written   int
}

func (e *errWriter) Write(p []byte) (n int, err error) {
	if e.written >= e.failAfter {
		return 0, os.ErrClosed
	}
	e.written++
	return len(p), nil
}

func TestInlineReporter_EdgeCasesAndErrors(t *testing.T) {
	rep := reporter.NewInlineReporter(reporter.ColorNever)

	// Nil result handling
	var buf bytes.Buffer
	if err := rep.Render(&buf, nil); err != nil {
		t.Errorf("expected nil result to render without error: %v", err)
	}

	// Unknown severity badge with and without color
	resUnknown := &reporter.ScanResult{
		Diagnostics: []ir.Diagnostic{
			{
				File:     "unknown.tsx",
				Line:     1,
				Column:   1,
				Rule:     "custom.unknown",
				Severity: "fatal",
				Message:  "Unknown severity test",
			},
		},
	}
	var bufUnknown bytes.Buffer
	repAlways := reporter.NewInlineReporter(reporter.ColorAlways)
	if err := repAlways.Render(&bufUnknown, resUnknown); err != nil {
		t.Fatalf("render unknown failed: %v", err)
	}
	if !bytes.Contains(bufUnknown.Bytes(), []byte("[ERROR]")) {
		t.Errorf("expected default [ERROR] badge for unknown severity")
	}

	var bufUnknownNoColor bytes.Buffer
	if err := rep.Render(&bufUnknownNoColor, resUnknown); err != nil {
		t.Fatalf("render unknown no color failed: %v", err)
	}

	// Colored render for error and warn
	var bufViolationsColored bytes.Buffer
	if err := repAlways.Render(&bufViolationsColored, sampleViolationsResult()); err != nil {
		t.Fatalf("render violations colored failed: %v", err)
	}
	if !bytes.Contains(bufViolationsColored.Bytes(), []byte("\033[1;31m[ERROR]\033[0m")) ||
		!bytes.Contains(bufViolationsColored.Bytes(), []byte("\033[1;33m[WARN]\033[0m")) {
		t.Errorf("expected colored badges in output: %s", bufViolationsColored.String())
	}

	// Writer failure error paths
	for i := 0; i < 6; i++ {
		ew := &errWriter{failAfter: i}
		_ = rep.Render(ew, sampleViolationsResult())
	}
}

func TestJSONReporter_EdgeCasesAndErrors(t *testing.T) {
	rep := reporter.NewJSONReporter()

	// Nil result handling
	var buf bytes.Buffer
	if err := rep.Render(&buf, nil); err != nil {
		t.Errorf("expected nil result to render without error: %v", err)
	}

	// Rule without category dot and custom version
	resCustom := &reporter.ScanResult{
		Version: "2.0.0",
		Diagnostics: []ir.Diagnostic{
			{
				File:     "src/App.tsx",
				Line:     1,
				Column:   1,
				Rule:     "nodotrule",
				Severity: ir.SeverityError,
				Message:  "Rule without dot",
			},
		},
	}
	var bufCustom bytes.Buffer
	if err := rep.Render(&bufCustom, resCustom); err != nil {
		t.Fatalf("custom render failed: %v", err)
	}
	if !bytes.Contains(bufCustom.Bytes(), []byte(`"category": ""`)) {
		t.Errorf("expected empty category for rule without dot")
	}

	// Writer failure
	ew := &errWriter{failAfter: 0}
	if err := rep.Render(ew, resCustom); err == nil {
		t.Errorf("expected error from failing writer, got nil")
	}
}

type parsedTestJSONDoc struct {
	Version string `json:"version"`
	Summary struct {
		ScannedFiles    int   `json:"scanned_files"`
		FilesWithIssues int   `json:"files_with_issues"`
		CleanFiles      int   `json:"clean_files"`
		DurationMS      int64 `json:"duration_ms"`
		ErrorCount      int   `json:"error_count"`
		WarningCount    int   `json:"warning_count"`
		InfoCount       int   `json:"info_count"`
		TotalIssues     int   `json:"total_issues"`
		Passed          bool  `json:"passed"`
	} `json:"summary"`
	Files []struct {
		File         string `json:"file"`
		ErrorCount   int    `json:"error_count"`
		WarningCount int    `json:"warning_count"`
		InfoCount    int    `json:"info_count"`
		TotalIssues  int    `json:"total_issues"`
		Violations   []struct {
			Line        int    `json:"line"`
			Column      int    `json:"column"`
			Rule        string `json:"rule"`
			Category    string `json:"category"`
			Severity    string `json:"severity"`
			Message     string `json:"message"`
			Hint        string `json:"hint"`
			DocURL      string `json:"doc_url"`
			Suppression string `json:"suppression"`
		} `json:"violations"`
	} `json:"files"`
	Diagnostics []struct {
		File     string `json:"file"`
		Line     int    `json:"line"`
		Column   int    `json:"column"`
		Rule     string `json:"rule"`
		Category string `json:"category"`
		Severity string `json:"severity"`
		Message  string `json:"message"`
		Hint     string `json:"hint"`
		DocURL   string `json:"doc_url"`
	} `json:"diagnostics"`
}

func renderAndParseJSON(t *testing.T, res *reporter.ScanResult) parsedTestJSONDoc {
	t.Helper()
	var buf bytes.Buffer
	rep := reporter.NewJSONReporter()
	if err := rep.Render(&buf, res); err != nil {
		t.Fatalf("failed to render JSON: %v", err)
	}
	var doc parsedTestJSONDoc
	if err := json.Unmarshal(buf.Bytes(), &doc); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v\nJSON:\n%s", err, buf.String())
	}
	return doc
}

func TestReporter_Invariants(t *testing.T) {
	t.Run("1_EmptyResult", func(t *testing.T) {
		res := &reporter.ScanResult{
			Summary: reporter.ScanSummary{ScannedFiles: 10},
		}
		doc := renderAndParseJSON(t, res)
		if len(doc.Files) != 0 || len(doc.Diagnostics) != 0 {
			t.Fatalf("expected 0 files and 0 diagnostics, got files=%d diags=%d", len(doc.Files), len(doc.Diagnostics))
		}
		if doc.Summary.FilesWithIssues != 0 || doc.Summary.CleanFiles != 10 || doc.Summary.TotalIssues != 0 || !doc.Summary.Passed {
			t.Fatalf("summary mismatch on empty result: %+v", doc.Summary)
		}
	})

	t.Run("2_SingleDiagnostic", func(t *testing.T) {
		res := &reporter.ScanResult{
			Summary: reporter.ScanSummary{ScannedFiles: 5},
			Diagnostics: []ir.Diagnostic{
				{File: "src/A.tsx", Line: 10, Column: 2, Rule: "theme.hardcode-color", Severity: ir.SeverityWarn, Message: "color"},
			},
		}
		doc := renderAndParseJSON(t, res)
		if len(doc.Files) != 1 || len(doc.Diagnostics) != 1 {
			t.Fatalf("expected 1 file and 1 diag, got files=%d diags=%d", len(doc.Files), len(doc.Diagnostics))
		}
		if doc.Files[0].File != "src/A.tsx" || doc.Files[0].TotalIssues != 1 {
			t.Fatalf("unexpected file group: %+v", doc.Files[0])
		}
	})

	t.Run("3_SingleFileMultipleRules", func(t *testing.T) {
		res := &reporter.ScanResult{
			Diagnostics: []ir.Diagnostic{
				{File: "src/Card.tsx", Line: 5, Column: 1, Rule: "theme.hardcode-size", Severity: ir.SeverityWarn},
				{File: "src/Card.tsx", Line: 10, Column: 2, Rule: "theme.hardcode-color", Severity: ir.SeverityError},
			},
		}
		doc := renderAndParseJSON(t, res)
		if len(doc.Files) != 1 || len(doc.Files[0].Violations) != 2 {
			t.Fatalf("expected 1 file group with 2 violations, got %d files", len(doc.Files))
		}
	})

	t.Run("4_MultipleFilesMultipleRules", func(t *testing.T) {
		res := &reporter.ScanResult{
			Diagnostics: []ir.Diagnostic{
				{File: "src/B.tsx", Line: 1, Rule: "theme.b"},
				{File: "src/A.tsx", Line: 1, Rule: "theme.a"},
			},
		}
		doc := renderAndParseJSON(t, res)
		if len(doc.Files) != 2 || doc.Files[0].File != "src/A.tsx" || doc.Files[1].File != "src/B.tsx" {
			t.Fatalf("files not sorted alphabetically: %+v", doc.Files)
		}
	})

	t.Run("5_SameLineDifferentColumns", func(t *testing.T) {
		res := &reporter.ScanResult{
			Diagnostics: []ir.Diagnostic{
				{File: "src/A.tsx", Line: 10, Column: 20, Rule: "theme.c"},
				{File: "src/A.tsx", Line: 10, Column: 5, Rule: "theme.a"},
				{File: "src/A.tsx", Line: 10, Column: 12, Rule: "theme.b"},
			},
		}
		doc := renderAndParseJSON(t, res)
		v := doc.Files[0].Violations
		if v[0].Column != 5 || v[1].Column != 12 || v[2].Column != 20 {
			t.Fatalf("columns not in ascending order: %d, %d, %d", v[0].Column, v[1].Column, v[2].Column)
		}
	})

	t.Run("6_SamePositionDifferentRules", func(t *testing.T) {
		res := &reporter.ScanResult{
			Diagnostics: []ir.Diagnostic{
				{File: "src/A.tsx", Line: 10, Column: 5, Rule: "theme.z"},
				{File: "src/A.tsx", Line: 10, Column: 5, Rule: "theme.a"},
			},
		}
		doc := renderAndParseJSON(t, res)
		v := doc.Files[0].Violations
		if v[0].Rule != "theme.a" || v[1].Rule != "theme.z" {
			t.Fatalf("rules not in ascending order: %s, %s", v[0].Rule, v[1].Rule)
		}
	})

	t.Run("7_IdenticalPositionSameRuleDifferentMessages", func(t *testing.T) {
		res := &reporter.ScanResult{
			Diagnostics: []ir.Diagnostic{
				{File: "src/A.tsx", Line: 10, Column: 5, Rule: "theme.a", Message: "zebra"},
				{File: "src/A.tsx", Line: 10, Column: 5, Rule: "theme.a", Message: "apple"},
			},
		}
		doc := renderAndParseJSON(t, res)
		v := doc.Files[0].Violations
		if v[0].Message != "apple" || v[1].Message != "zebra" {
			t.Fatalf("messages not in ascending order: %s, %s", v[0].Message, v[1].Message)
		}
	})

	t.Run("8_ShuffledFilesCanonicalized", func(t *testing.T) {
		res := &reporter.ScanResult{
			Diagnostics: []ir.Diagnostic{
				{File: "src/z.tsx", Line: 1},
				{File: "src/a.tsx", Line: 1},
				{File: "src/m.tsx", Line: 1},
			},
		}
		doc := renderAndParseJSON(t, res)
		if doc.Files[0].File != "src/a.tsx" || doc.Files[1].File != "src/m.tsx" || doc.Files[2].File != "src/z.tsx" {
			t.Fatalf("files not canonicalized: %s, %s, %s", doc.Files[0].File, doc.Files[1].File, doc.Files[2].File)
		}
	})

	t.Run("9_ShuffledDiagnosticsCanonicalized", func(t *testing.T) {
		res := &reporter.ScanResult{
			Diagnostics: []ir.Diagnostic{
				{File: "src/a.tsx", Line: 50},
				{File: "src/a.tsx", Line: 5},
				{File: "src/a.tsx", Line: 20},
			},
		}
		doc := renderAndParseJSON(t, res)
		v := doc.Files[0].Violations
		if v[0].Line != 5 || v[1].Line != 20 || v[2].Line != 50 {
			t.Fatalf("diagnostics not sorted by line: %d, %d, %d", v[0].Line, v[1].Line, v[2].Line)
		}
	})

	t.Run("10_I3_FlattenConsistency", func(t *testing.T) {
		res := sampleViolationsResult()
		doc := renderAndParseJSON(t, res)

		var flattened []struct {
			File   string
			Line   int
			Column int
			Rule   string
		}
		for _, f := range doc.Files {
			for _, v := range f.Violations {
				flattened = append(flattened, struct {
					File   string
					Line   int
					Column int
					Rule   string
				}{File: f.File, Line: v.Line, Column: v.Column, Rule: v.Rule})
			}
		}

		if len(flattened) != len(doc.Diagnostics) {
			t.Fatalf("flatten length %d != diagnostics length %d", len(flattened), len(doc.Diagnostics))
		}
		for i := range flattened {
			if flattened[i].File != doc.Diagnostics[i].File ||
				flattened[i].Line != doc.Diagnostics[i].Line ||
				flattened[i].Column != doc.Diagnostics[i].Column ||
				flattened[i].Rule != doc.Diagnostics[i].Rule {
				t.Fatalf("diagnostic item %d mismatch: %+v vs %+v", i, flattened[i], doc.Diagnostics[i])
			}
		}
	})

	t.Run("11_I4_I5_PerFileCounters", func(t *testing.T) {
		res := &reporter.ScanResult{
			Diagnostics: []ir.Diagnostic{
				{File: "src/A.tsx", Line: 1, Severity: ir.SeverityError},
				{File: "src/A.tsx", Line: 2, Severity: ir.SeverityWarn},
				{File: "src/A.tsx", Line: 3, Severity: ir.SeverityInfo},
			},
		}
		doc := renderAndParseJSON(t, res)
		f := doc.Files[0]
		if f.TotalIssues != 3 || f.ErrorCount != 1 || f.WarningCount != 1 || f.InfoCount != 1 {
			t.Fatalf("per-file counters mismatch: %+v", f)
		}
	})

	t.Run("12_I6_I7_I8_GlobalCounters", func(t *testing.T) {
		res := &reporter.ScanResult{
			Summary: reporter.ScanSummary{ScannedFiles: 10},
			Diagnostics: []ir.Diagnostic{
				{File: "src/A.tsx", Line: 1, Severity: ir.SeverityError},
				{File: "src/B.tsx", Line: 1, Severity: ir.SeverityWarn},
			},
		}
		doc := renderAndParseJSON(t, res)
		if doc.Summary.TotalIssues != 2 || doc.Summary.ErrorCount != 1 || doc.Summary.WarningCount != 1 {
			t.Fatalf("global severity counters mismatch: %+v", doc.Summary)
		}
		if doc.Summary.FilesWithIssues != 2 || doc.Summary.CleanFiles != 8 {
			t.Fatalf("global file counters mismatch: %+v", doc.Summary)
		}
	})

	t.Run("13_PathNormalization", func(t *testing.T) {
		res := &reporter.ScanResult{
			Diagnostics: []ir.Diagnostic{
				{File: "src/components/Card.tsx", Line: 1},
			},
		}
		doc := renderAndParseJSON(t, res)
		if doc.Files[0].File != "src/components/Card.tsx" {
			t.Fatalf("expected POSIX path, got %s", doc.Files[0].File)
		}
	})

	t.Run("14_WindowsSeparatorNormalization", func(t *testing.T) {
		res := &reporter.ScanResult{
			Diagnostics: []ir.Diagnostic{
				{File: `src\components\Card.tsx`, Line: 1},
			},
		}
		doc := renderAndParseJSON(t, res)
		if doc.Files[0].File != "src/components/Card.tsx" {
			t.Fatalf("expected normalized POSIX path, got %s", doc.Files[0].File)
		}
	})

	t.Run("15_SuppressionAstro", func(t *testing.T) {
		res := &reporter.ScanResult{
			Diagnostics: []ir.Diagnostic{
				{File: "src/pages/index.astro", Line: 1, Rule: "theme.color"},
			},
		}
		doc := renderAndParseJSON(t, res)
		expected := "<!-- charites:ignore theme.color <reason> -->"
		if doc.Files[0].Violations[0].Suppression != expected {
			t.Fatalf("expected %s, got %s", expected, doc.Files[0].Violations[0].Suppression)
		}
	})

	t.Run("16_SuppressionTSX", func(t *testing.T) {
		res := &reporter.ScanResult{
			Diagnostics: []ir.Diagnostic{
				{File: "src/Card.tsx", Line: 1, Rule: "theme.color"},
			},
		}
		doc := renderAndParseJSON(t, res)
		expected := "// charites:ignore theme.color <reason>"
		if doc.Files[0].Violations[0].Suppression != expected {
			t.Fatalf("expected %s, got %s", expected, doc.Files[0].Violations[0].Suppression)
		}
	})

	t.Run("17_SuppressionJSX", func(t *testing.T) {
		res := &reporter.ScanResult{
			Diagnostics: []ir.Diagnostic{
				{File: "src/Card.jsx", Line: 1, Rule: "theme.color"},
			},
		}
		doc := renderAndParseJSON(t, res)
		expected := "// charites:ignore theme.color <reason>"
		if doc.Files[0].Violations[0].Suppression != expected {
			t.Fatalf("expected %s, got %s", expected, doc.Files[0].Violations[0].Suppression)
		}
	})

	t.Run("18_MarkdownFileOrdering", func(t *testing.T) {
		res := &reporter.ScanResult{
			Diagnostics: []ir.Diagnostic{
				{File: "src/Z.tsx", Line: 1},
				{File: "src/A.tsx", Line: 1},
			},
		}
		var buf bytes.Buffer
		rep := reporter.NewMarkdownReporter()
		_ = rep.Render(&buf, res)
		out := buf.String()
		posA := strings.Index(out, "src/A.tsx")
		posZ := strings.Index(out, "src/Z.tsx")
		if posA == -1 || posZ == -1 || posA >= posZ {
			t.Fatalf("Markdown files not sorted: posA=%d, posZ=%d", posA, posZ)
		}
	})

	t.Run("19_MarkdownViolationOrdering", func(t *testing.T) {
		res := &reporter.ScanResult{
			Diagnostics: []ir.Diagnostic{
				{File: "src/A.tsx", Line: 40},
				{File: "src/A.tsx", Line: 10},
			},
		}
		var buf bytes.Buffer
		rep := reporter.NewMarkdownReporter()
		_ = rep.Render(&buf, res)
		out := buf.String()
		pos10 := strings.Index(out, "L10")
		pos40 := strings.Index(out, "L40")
		if pos10 == -1 || pos40 == -1 || pos10 >= pos40 {
			t.Fatalf("Markdown violations not sorted by line: pos10=%d, pos40=%d", pos10, pos40)
		}
	})

	t.Run("20_MarkdownCleanScan", func(t *testing.T) {
		res := &reporter.ScanResult{
			Summary: reporter.ScanSummary{ScannedFiles: 5, Passed: true},
		}
		var buf bytes.Buffer
		rep := reporter.NewMarkdownReporter()
		_ = rep.Render(&buf, res)
		out := buf.String()
		if !strings.Contains(out, "## Results by File\n\nNo known design token or ergonomics violations found") {
			t.Fatalf("expected clean scan banner, got:\n%s", out)
		}
	})

	t.Run("21_I9_Determinism", func(t *testing.T) {
		res := sampleViolationsResult()
		repJSON := reporter.NewJSONReporter()
		repMD := reporter.NewMarkdownReporter(reporter.WithTimestamp(res.Timestamp))

		var firstJSON, firstMD []byte
		for i := 0; i < 5; i++ {
			var bufJ, bufM bytes.Buffer
			_ = repJSON.Render(&bufJ, res)
			_ = repMD.Render(&bufM, res)
			if i == 0 {
				firstJSON = bufJ.Bytes()
				firstMD = bufM.Bytes()
			} else {
				if !bytes.Equal(firstJSON, bufJ.Bytes()) {
					t.Fatalf("JSON render nondeterministic on iteration %d", i)
				}
				if !bytes.Equal(firstMD, bufM.Bytes()) {
					t.Fatalf("Markdown render nondeterministic on iteration %d", i)
				}
			}
		}
	})
}
