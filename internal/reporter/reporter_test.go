package reporter_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/reporter"
)

func sampleCleanResult() *reporter.ScanResult {
	return &reporter.ScanResult{
		Version: "1.0.0",
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
		Version: "1.0.0",
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
