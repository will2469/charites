package cli_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/will2469/charites/internal/cli"
	"github.com/will2469/charites/internal/reporter"
)

func setupCleanTestRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	file1 := filepath.Join(dir, "Button.tsx")
	_ = os.WriteFile(file1, []byte(`export const Button = () => <button className="bg-primary text-primary-foreground">Click</button>;`), 0o600)

	file2 := filepath.Join(dir, "Header.astro")
	_ = os.WriteFile(file2, []byte(`---
---
<header class="p-4 bg-surface"></header>`), 0o600)
	return dir
}

func setupViolationTestRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	file1 := filepath.Join(dir, "Card.astro")
	_ = os.WriteFile(file1, []byte(`---
---
<div class="bg-primary/10">Content</div>`), 0o600)
	return dir
}

func TestRunScan_CleanRepo(t *testing.T) {
	dir := setupCleanTestRepo(t)

	var stdout, stderr bytes.Buffer
	code := cli.RunScan([]string{dir}, &stdout, &stderr)

	if code != cli.ExitClean {
		t.Fatalf("expected exit clean (0), got %d, stderr: %s", code, stderr.String())
	}
	out := stdout.String()
	if !strings.Contains(out, "0 problems found (0 errors, 0 warnings)") {
		t.Errorf("stdout does not contain clean problems summary:\n%s", out)
	}
	if !strings.Contains(out, "Scanned 2 files in") {
		t.Errorf("stdout does not contain scanned files count:\n%s", out)
	}
}

func TestRunScan_ViolationsRepo(t *testing.T) {
	dir := setupViolationTestRepo(t)

	var stdout, stderr bytes.Buffer
	code := cli.RunScan([]string{"--no-color", dir}, &stdout, &stderr)

	if code != cli.ExitViolations {
		t.Fatalf("expected exit violations (1), got %d, stderr: %s", code, stderr.String())
	}

	out := stdout.String()
	if !strings.Contains(out, "[ERROR]") {
		t.Errorf("expected [ERROR] badge in output:\n%s", out)
	}
	if !strings.Contains(out, "theme.hardcode-opacity-color") {
		t.Errorf("expected rule ID in output:\n%s", out)
	}
	if !strings.Contains(out, "1 problem found (1 error, 0 warnings)") {
		t.Errorf("expected 1 problem summary in output:\n%s", out)
	}
}

func TestRunScan_JSONFormat(t *testing.T) {
	dir := setupViolationTestRepo(t)

	var stdout, stderr bytes.Buffer
	code := cli.RunScan([]string{"--format=json", dir}, &stdout, &stderr)

	if code != cli.ExitViolations {
		t.Fatalf("expected exit violations (1), got %d, stderr: %s", code, stderr.String())
	}

	var doc reporter.ScanResult
	if err := json.Unmarshal(stdout.Bytes(), &doc); err != nil {
		t.Fatalf("failed to unmarshal JSON reporter output: %v\nOutput: %s", err, stdout.String())
	}

	if doc.Summary.ErrorCount != 1 {
		t.Errorf("expected 1 error count in json, got %d", doc.Summary.ErrorCount)
	}
	if doc.Summary.Passed {
		t.Errorf("expected passed: false in json")
	}
	if len(doc.Diagnostics) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(doc.Diagnostics))
	}
	if doc.Diagnostics[0].Rule != "theme.hardcode-opacity-color" {
		t.Errorf("unexpected rule ID: %s", doc.Diagnostics[0].Rule)
	}
}

func TestRunScan_ExtNormalizationAndValidation(t *testing.T) {
	dir := setupCleanTestRepo(t)

	// Valid repeated and comma separated
	var stdout, stderr bytes.Buffer
	code := cli.RunScan([]string{"--ext", "astro,tsx", "-e", "jsx", dir}, &stdout, &stderr)
	if code != cli.ExitClean {
		t.Errorf("expected valid ext to pass with exit 0, got %d, stderr: %s", code, stderr.String())
	}

	// Empty extension
	stderr.Reset()
	code = cli.RunScan([]string{"--ext=", dir}, &stdout, &stderr)
	if code != cli.ExitOperational {
		t.Errorf("expected empty ext to exit 2, got %d", code)
	}
	if !strings.Contains(stderr.String(), "empty extension flag") {
		t.Errorf("expected empty extension error message, got: %s", stderr.String())
	}

	// Unsupported extension
	stderr.Reset()
	code = cli.RunScan([]string{"--ext=vue", dir}, &stdout, &stderr)
	if code != cli.ExitOperational {
		t.Errorf("expected unsupported ext to exit 2, got %d", code)
	}
	if !strings.Contains(stderr.String(), `unsupported extension "vue"`) {
		t.Errorf("expected unsupported extension error message, got: %s", stderr.String())
	}
}

func TestRunScan_CategoryAndRuleValidation(t *testing.T) {
	dir := setupCleanTestRepo(t)

	// Unknown category
	var stdout, stderr bytes.Buffer
	code := cli.RunScan([]string{"--category=unknown_category", dir}, &stdout, &stderr)
	if code != cli.ExitOperational {
		t.Errorf("expected unknown category to exit 2, got %d", code)
	}
	if !strings.Contains(stderr.String(), `unknown category "unknown_category"`) {
		t.Errorf("unexpected stderr: %s", stderr.String())
	}

	// Unknown rule
	stderr.Reset()
	code = cli.RunScan([]string{"--rule=theme.nonexistent_rule", dir}, &stdout, &stderr)
	if code != cli.ExitOperational {
		t.Errorf("expected unknown rule to exit 2, got %d", code)
	}
	if !strings.Contains(stderr.String(), `unknown rule "theme.nonexistent_rule"`) {
		t.Errorf("unexpected stderr: %s", stderr.String())
	}

	// Intersection conflict
	stderr.Reset()
	// Catatan: saat ini category terdaftar adalah theme
	code = cli.RunScan([]string{"--category=theme", "--rule=theme.hardcode-opacity-color", dir}, &stdout, &stderr)
	if code != cli.ExitClean {
		t.Errorf("expected valid intersection to pass, got %d, stderr: %s", code, stderr.String())
	}
}

func TestRunScan_ConfigHandlingAndPrecedence(t *testing.T) {
	dir := setupViolationTestRepo(t)

	// Missing custom config
	var stdout, stderr bytes.Buffer
	code := cli.RunScan([]string{"--config=nonexistent_config.yaml", dir}, &stdout, &stderr)
	if code != cli.ExitOperational {
		t.Errorf("expected missing config to exit 2, got %d", code)
	}
	if !strings.Contains(stderr.String(), "config file not found") {
		t.Errorf("unexpected stderr: %s", stderr.String())
	}

	// Config turning rule "off" (3-tier precedence: policy 'off' beats explicit CLI --rule)
	cfgFile := filepath.Join(dir, "charites.yaml")
	_ = os.WriteFile(cfgFile, []byte("rules:\n  theme.hardcode-opacity-color: off\n"), 0o600)

	stdout.Reset()
	stderr.Reset()
	code = cli.RunScan([]string{"--config", cfgFile, "--rule=theme.hardcode-opacity-color", dir}, &stdout, &stderr)
	if code != cli.ExitClean {
		t.Errorf("expected policy off to override CLI rule and exit 0, got %d, stderr: %s", code, stderr.String())
	}
}

func TestRunScan_TargetValidation(t *testing.T) {
	var stdout, stderr bytes.Buffer

	// Multiple targets
	code := cli.RunScan([]string{"target1", "target2"}, &stdout, &stderr)
	if code != cli.ExitOperational {
		t.Errorf("expected multiple targets to exit 2, got %d", code)
	}
	if !strings.Contains(stderr.String(), "multiple scan targets not supported") {
		t.Errorf("unexpected stderr: %s", stderr.String())
	}

	// Non-existent target
	stderr.Reset()
	code = cli.RunScan([]string{"nonexistent_directory_target_xyz"}, &stdout, &stderr)
	if code != cli.ExitOperational {
		t.Errorf("expected non-existent target to exit 2, got %d", code)
	}
	if !strings.Contains(stderr.String(), "does not exist") {
		t.Errorf("unexpected stderr: %s", stderr.String())
	}

	// Direct-target safety on builtin hard exclusion
	tmpDir := t.TempDir()
	nodeModulesTarget := filepath.Join(tmpDir, "node_modules", "package", "Component.tsx")
	_ = os.MkdirAll(filepath.Dir(nodeModulesTarget), 0o750)
	_ = os.WriteFile(nodeModulesTarget, []byte("<div />"), 0o600)

	stderr.Reset()
	code = cli.RunScan([]string{nodeModulesTarget}, &stdout, &stderr)
	if code != cli.ExitOperational {
		t.Errorf("expected builtin hard exclusion target to exit 2, got %d", code)
	}
	if !strings.Contains(stderr.String(), "within excluded directory (builtin hard exclusion)") {
		t.Errorf("unexpected stderr: %s", stderr.String())
	}
}

func TestRunScan_UnsupportedFormat(t *testing.T) {
	dir := setupCleanTestRepo(t)
	var stdout, stderr bytes.Buffer
	code := cli.RunScan([]string{"--format=yaml", dir}, &stdout, &stderr)
	if code != cli.ExitOperational {
		t.Errorf("expected unsupported format to exit 2, got %d", code)
	}
	if !strings.Contains(stderr.String(), `unsupported format "yaml"`) {
		t.Errorf("unexpected stderr: %s", stderr.String())
	}
}

func TestResolveExitCode(t *testing.T) {
	if code := cli.ResolveExitCode(nil, false); code != cli.ExitClean {
		t.Errorf("expected nil summary to return ExitClean, got %d", code)
	}

	clean := &reporter.ScanSummary{ErrorCount: 0, WarningCount: 0}
	if code := cli.ResolveExitCode(clean, false); code != cli.ExitClean {
		t.Errorf("expected clean to return ExitClean, got %d", code)
	}

	errSummary := &reporter.ScanSummary{ErrorCount: 1, WarningCount: 0}
	if code := cli.ResolveExitCode(errSummary, false); code != cli.ExitViolations {
		t.Errorf("expected error to return ExitViolations, got %d", code)
	}

	warnSummary := &reporter.ScanSummary{ErrorCount: 0, WarningCount: 2}
	if code := cli.ResolveExitCode(warnSummary, false); code != cli.ExitClean {
		t.Errorf("expected warning without failOnWarn to return ExitClean, got %d", code)
	}
	if code := cli.ResolveExitCode(warnSummary, true); code != cli.ExitViolations {
		t.Errorf("expected warning with failOnWarn to return ExitViolations, got %d", code)
	}
}
