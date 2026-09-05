package e2e_test

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const binaryPath = "../../bin/charites"

func getBinaryPath(t *testing.T) string {
	t.Helper()
	abs, err := filepath.Abs(binaryPath)
	if err != nil {
		t.Fatalf("failed to get absolute binary path: %v", err)
	}
	return abs
}

func runBinary(t *testing.T, args []string, env ...string) (stdout string, stderr string, exitCode int) {
	t.Helper()
	// #nosec G204 -- test executes compiled binary with controlled arguments
	cmd := exec.Command(getBinaryPath(t), args...)
	if len(env) > 0 {
		cmd.Env = append(os.Environ(), env...)
	}

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	exitCode = 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("failed to run command: %v", err)
		}
	}

	return outBuf.String(), errBuf.String(), exitCode
}

func setupFixtureRepo(t *testing.T) (cleanDir, violationDir string) {
	t.Helper()
	cleanDir = t.TempDir()
	violationDir = t.TempDir()

	// Clean file
	_ = os.WriteFile(filepath.Join(cleanDir, "Card.tsx"), []byte(`export const Card = () => <div className="bg-primary text-primary-foreground" />;`), 0o600)

	// Violation file (theme.hardcode-opacity-color)
	_ = os.WriteFile(filepath.Join(violationDir, "Card.astro"), []byte(`---
---
<div class="bg-primary/10">Card</div>`), 0o600)

	return cleanDir, violationDir
}

// Skenario 1: charites -> scan .
func TestScenario01_EmptyArgsDefaultsToScan(t *testing.T) {
	cleanDir, _ := setupFixtureRepo(t)

	// #nosec G204
	cmd := exec.Command(getBinaryPath(t))
	cmd.Dir = cleanDir
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	if err != nil {
		t.Fatalf("expected exit 0, got %v, stderr: %s", err, errBuf.String())
	}

	if !strings.Contains(outBuf.String(), "0 problems found") {
		t.Errorf("expected clean scan output, got: %s", outBuf.String())
	}
}

// Skenario 2: Path/Flag tanpa subcommand
func TestScenario02_DirectPathAndFlag(t *testing.T) {
	cleanDir, violationDir := setupFixtureRepo(t)

	// Direct path
	stdout, stderr, code := runBinary(t, []string{cleanDir})
	if code != 0 {
		t.Fatalf("expected direct path to exit 0, got %d, stderr: %s", code, stderr)
	}
	if !strings.Contains(stdout, "0 problems found") {
		t.Errorf("expected clean summary, got: %s", stdout)
	}

	// Direct flag with path
	stdout, stderr, code = runBinary(t, []string{"--format=json", violationDir})
	if code != 1 {
		t.Fatalf("expected direct flag to exit 1 on violation, got %d, stderr: %s", code, stderr)
	}
	if !strings.Contains(stdout, `"theme.hardcode-opacity-color"`) {
		t.Errorf("expected json diagnostic in output: %s", stdout)
	}
}

// Skenario 3: Ekuivalensi scan, check, run
func TestScenario03_SubcommandEquivalence(t *testing.T) {
	_, violationDir := setupFixtureRepo(t)

	outScan, errScan, codeScan := runBinary(t, []string{"scan", "--format=json", violationDir})
	outCheck, errCheck, codeCheck := runBinary(t, []string{"check", "--format=json", violationDir})
	outRun, errRun, codeRun := runBinary(t, []string{"run", "--format=json", violationDir})

	if codeScan != 1 || codeCheck != 1 || codeRun != 1 {
		t.Fatalf("expected all commands to exit 1, got scan=%d, check=%d, run=%d", codeScan, codeCheck, codeRun)
	}

	if errScan != errCheck || errCheck != errRun {
		t.Errorf("stderrs differ between aliases")
	}

	// Unmarshal json to compare equivalence ignoring dynamic duration_ms
	var docScan, docCheck, docRun map[string]interface{}
	_ = json.Unmarshal([]byte(outScan), &docScan)
	_ = json.Unmarshal([]byte(outCheck), &docCheck)
	_ = json.Unmarshal([]byte(outRun), &docRun)

	delete(docScan["summary"].(map[string]interface{}), "duration_ms")
	delete(docCheck["summary"].(map[string]interface{}), "duration_ms")
	delete(docRun["summary"].(map[string]interface{}), "duration_ms")

	bScan, _ := json.Marshal(docScan)
	bCheck, _ := json.Marshal(docCheck)
	bRun, _ := json.Marshal(docRun)

	if !bytes.Equal(bScan, bCheck) || !bytes.Equal(bCheck, bRun) {
		t.Fatalf("scan, check, run produce divergent JSON payloads")
	}
}

// Skenario 4: Valid vs Invalid --ext
func TestScenario04_ValidVsInvalidExt(t *testing.T) {
	cleanDir, _ := setupFixtureRepo(t)

	// Valid
	_, _, code := runBinary(t, []string{"scan", "--ext=astro,tsx", cleanDir})
	if code != 0 {
		t.Errorf("expected valid ext to exit 0, got %d", code)
	}

	// Invalid
	_, stderr, code := runBinary(t, []string{"scan", "--ext=vue", cleanDir})
	if code != 2 {
		t.Errorf("expected invalid ext to exit 2, got %d", code)
	}
	if !strings.Contains(stderr, `unsupported extension "vue"`) {
		t.Errorf("unexpected stderr: %s", stderr)
	}
}

// Skenario 5: Repeated + Comma-separated --ext
func TestScenario05_RepeatedAndCommaExt(t *testing.T) {
	cleanDir, _ := setupFixtureRepo(t)

	_, stderr, code := runBinary(t, []string{"scan", "--ext", "astro,tsx", "--ext", "jsx", cleanDir})
	if code != 0 {
		t.Errorf("expected repeated/comma ext to exit 0, got %d, stderr: %s", code, stderr)
	}
}

// Skenario 6: --category x --rule Intersection
func TestScenario06_CategoryRuleIntersection(t *testing.T) {
	cleanDir, _ := setupFixtureRepo(t)

	// Valid intersection
	_, stderr, code := runBinary(t, []string{"scan", "--category=theme", "--rule=theme.hardcode-opacity-color", cleanDir})
	if code != 0 {
		t.Errorf("expected valid category/rule intersection to exit 0, got %d, stderr: %s", code, stderr)
	}

	// Unknown category
	_, stderr, code = runBinary(t, []string{"scan", "--category=bogus_cat", cleanDir})
	if code != 2 {
		t.Errorf("expected unknown category to exit 2, got %d", code)
	}
	if !strings.Contains(stderr, `unknown category "bogus_cat"`) {
		t.Errorf("unexpected stderr: %s", stderr)
	}
}

// Skenario 7: Config off vs Explicit --rule
func TestScenario07_ConfigPolicyBeatsCLI(t *testing.T) {
	_, violationDir := setupFixtureRepo(t)

	cfgPath := filepath.Join(violationDir, "charites.yaml")
	_ = os.WriteFile(cfgPath, []byte("rules:\n  theme.hardcode-opacity-color: off\n"), 0o600)

	_, stderr, code := runBinary(t, []string{"scan", "--config", cfgPath, "--rule=theme.hardcode-opacity-color", violationDir})
	if code != 0 {
		t.Fatalf("expected config off to beat CLI --rule and exit 0, got %d, stderr: %s", code, stderr)
	}
}

// Skenario 8: Builtin Hard Exclusion + --ignore
func TestScenario08_BuiltinHardExclusionImmunity(t *testing.T) {
	tmpDir := t.TempDir()
	nodeModulesDir := filepath.Join(tmpDir, "node_modules", "pkg")
	_ = os.MkdirAll(nodeModulesDir, 0o750)
	_ = os.WriteFile(filepath.Join(nodeModulesDir, "Violator.astro"), []byte(`<div class="bg-primary/10" />`), 0o600)

	// Try to un-ignore node_modules using '!node_modules'
	stdout, stderr, code := runBinary(t, []string{"scan", "--ignore=!node_modules", tmpDir})
	if code != 0 {
		t.Fatalf("expected node_modules to remain excluded (exit 0), got %d, stderr: %s", code, stderr)
	}
	if !strings.Contains(stdout, "0 problems found") {
		t.Errorf("node_modules was scanned despite builtin hard exclusion: %s", stdout)
	}
}

// Skenario 9: Inline ANSI Reporter
func TestScenario09_InlineReporter(t *testing.T) {
	_, violationDir := setupFixtureRepo(t)

	stdout, _, code := runBinary(t, []string{"scan", "--format=inline", violationDir})
	if code != 1 {
		t.Fatalf("expected exit 1 on violation, got %d", code)
	}
	if !strings.Contains(stdout, "[ERROR]") {
		t.Errorf("expected [ERROR] badge in inline reporter: %s", stdout)
	}
	if !strings.Contains(stdout, "1 problem found (1 error, 0 warnings)") {
		t.Errorf("expected problem summary in inline reporter: %s", stdout)
	}
}

// Skenario 10: JSON Reporter & Schema
func TestScenario10_JSONReporter(t *testing.T) {
	_, violationDir := setupFixtureRepo(t)

	stdout, _, code := runBinary(t, []string{"scan", "--format=json", violationDir})
	if code != 1 {
		t.Fatalf("expected exit 1 on violation, got %d", code)
	}

	var res struct {
		Version string `json:"version"`
		Summary struct {
			ScannedFiles int   `json:"scanned_files"`
			DurationMS   int64 `json:"duration_ms"`
			ErrorCount   int   `json:"error_count"`
			WarningCount int   `json:"warning_count"`
			InfoCount    int   `json:"info_count"`
			Passed       bool  `json:"passed"`
		} `json:"summary"`
		Diagnostics []struct {
			File     string `json:"file"`
			Line     int    `json:"line"`
			Column   int    `json:"column"`
			Rule     string `json:"rule"`
			Category string `json:"category"`
			Severity string `json:"severity"`
			Message  string `json:"message"`
		} `json:"diagnostics"`
	}

	if err := json.Unmarshal([]byte(stdout), &res); err != nil {
		t.Fatalf("failed to parse JSON reporter output: %v\nOutput: %s", err, stdout)
	}

	if res.Summary.ErrorCount != 1 || res.Summary.Passed {
		t.Errorf("unexpected summary values: %+v", res.Summary)
	}
	if len(res.Diagnostics) != 1 || res.Diagnostics[0].Category != "theme" {
		t.Errorf("unexpected diagnostics values: %+v", res.Diagnostics)
	}
}

// Skenario 11: TTY / Pipe / NO_COLOR
func TestScenario11_NoColorResolution(t *testing.T) {
	_, violationDir := setupFixtureRepo(t)

	// In exec.Command without pseudo-terminal, stdout is piped -> Mode Never ANSI
	stdout, _, code := runBinary(t, []string{"scan", violationDir})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if strings.Contains(stdout, "\033[") {
		t.Errorf("piped output contained ANSI escape sequences: %q", stdout)
	}

	// Explicit NO_COLOR=1
	stdout, _, code = runBinary(t, []string{"scan", violationDir}, "NO_COLOR=1")
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if strings.Contains(stdout, "\033[") {
		t.Errorf("NO_COLOR=1 output contained ANSI escape sequences: %q", stdout)
	}
}

// Skenario 12: Deterministic Total Ordering
func TestScenario12_DeterministicOrdering(t *testing.T) {
	_, violationDir := setupFixtureRepo(t)

	// Create multiple violating files
	for i := 0; i < 5; i++ {
		_ = os.WriteFile(filepath.Join(violationDir, string(rune('A'+i))+".astro"), []byte(`<div class="bg-primary/10" />`), 0o600)
	}

	out1, _, _ := runBinary(t, []string{"scan", "--format=json", violationDir})
	out2, _, _ := runBinary(t, []string{"scan", "--format=json", violationDir})

	var doc1, doc2 map[string]interface{}
	_ = json.Unmarshal([]byte(out1), &doc1)
	_ = json.Unmarshal([]byte(out2), &doc2)

	delete(doc1["summary"].(map[string]interface{}), "duration_ms")
	delete(doc2["summary"].(map[string]interface{}), "duration_ms")

	b1, _ := json.Marshal(doc1)
	b2, _ := json.Marshal(doc2)

	if !bytes.Equal(b1, b2) {
		t.Fatalf("two runs produced non-deterministic diagnostic ordering")
	}
}

// Skenario 13: Exit Codes Taxonomy
func TestScenario13_ExitCodesTaxonomy(t *testing.T) {
	cleanDir, violationDir := setupFixtureRepo(t)

	// Clean -> 0
	_, _, code0 := runBinary(t, []string{"scan", cleanDir})
	if code0 != 0 {
		t.Errorf("expected clean repo to exit 0, got %d", code0)
	}

	// Violations -> 1
	_, _, code1 := runBinary(t, []string{"scan", violationDir})
	if code1 != 1 {
		t.Errorf("expected violation repo to exit 1, got %d", code1)
	}

	// Operational error -> 2
	_, _, code2 := runBinary(t, []string{"scan", "--invalid-flag", cleanDir})
	if code2 != 2 {
		t.Errorf("expected invalid flag to exit 2, got %d", code2)
	}
}

// Skenario 14: Zero Host Footprint & Residual Audit
func TestScenario14_ZeroHostFootprintAudit(t *testing.T) {
	cleanDir, _ := setupFixtureRepo(t)

	// Audit snapshot pre-execution
	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".config", "charites")
	cacheDir := filepath.Join(home, ".cache", "charites")

	_, _, code := runBinary(t, []string{"scan", cleanDir})
	if code != 0 {
		t.Fatalf("scan failed: %d", code)
	}

	// Assert no global residual directories created
	if _, err := os.Stat(configDir); err == nil {
		t.Errorf("illegal host pollution: directory %s was created", configDir)
	}
	if _, err := os.Stat(cacheDir); err == nil {
		t.Errorf("illegal host pollution: directory %s was created", cacheDir)
	}

	// Assert no temporary swap files created in system temp
	matches, _ := filepath.Glob(filepath.Join(os.TempDir(), "charites*"))
	if len(matches) > 0 {
		t.Errorf("illegal residual files in temp: %v", matches)
	}
}

// Skenario 15: Update & Upgrade Equivalence dan Uninstall Help
func TestScenario15_UpdateUpgradeEquivalenceAndUninstallHelp(t *testing.T) {
	// 1. Update and Upgrade equivalence
	stdoutUpdate, stderrUpdate, codeUpdate := runBinary(t, []string{"update"})
	stdoutUpgrade, stderrUpgrade, codeUpgrade := runBinary(t, []string{"upgrade"})

	if codeUpdate != 0 {
		t.Fatalf("expected update to exit 0, got %d, stderr: %s", codeUpdate, stderrUpdate)
	}
	if codeUpgrade != 0 {
		t.Fatalf("expected upgrade to exit 0, got %d, stderr: %s", codeUpgrade, stderrUpgrade)
	}
	if stdoutUpdate != stdoutUpgrade {
		t.Errorf("update output (%q) differs from upgrade output (%q)", stdoutUpdate, stdoutUpgrade)
	}
	if !strings.Contains(stdoutUpdate, "No update found. Charites is up to date.") {
		t.Errorf("expected clean up-to-date message, got: %q", stdoutUpdate)
	}

	// 2. Uninstall help
	stdoutUninst, _, codeUninst := runBinary(t, []string{"uninstall", "-h"})
	if codeUninst != 0 {
		t.Fatalf("expected uninstall -h to exit 0, got %d", codeUninst)
	}
	if !strings.Contains(stdoutUninst, "Usage: charites") {
		t.Errorf("expected usage string for uninstall -h, got: %q", stdoutUninst)
	}
}
