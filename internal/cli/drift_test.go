package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/will2469/charites/internal/cli"
)

func TestRunDrift(t *testing.T) {
	tempDir := t.TempDir()

	// Buat beberapa file komponen uji
	file1 := `export function Buttons() {
  return (
    <div>
      <Button className="rounded-md bg-amber-500 text-white">Save 1</Button>
      <Button className="rounded-md bg-amber-500 text-white">Save 2</Button>
      <Button className="rounded-md bg-amber-500 text-white">Save 3</Button>
      <Button className="rounded-md bg-amber-500 text-white">Save 4</Button>
      <Button className="rounded-md bg-amber-500 text-white">Save 5</Button>
      <Button className="rounded-md bg-amber-500 text-white">Save 6</Button>
      <Button className="rounded-md bg-amber-500 text-white">Save 7</Button>
      <Button className="rounded-md bg-amber-500 text-white">Save 8</Button>
      <Button className="rounded-md bg-amber-500 text-white">Save 9</Button>
      <Button className="rounded-md bg-amber-500 text-white">Save 10</Button>
      <Button className="rounded-2xl bg-yellow-400 text-white">Rogue</Button>
    </div>
  );
}`
	if err := os.WriteFile(filepath.Join(tempDir, "Buttons.tsx"), []byte(file1), 0o600); err != nil {
		t.Fatalf("failed to write fixture file: %v", err)
	}

	t.Run("InlineReport", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		code := cli.RunDrift([]string{"--no-color", tempDir}, &stdout, &stderr)
		if code != cli.ExitClean {
			t.Fatalf("expected exit clean (0), got %d, stderr: %s", code, stderr.String())
		}
		out := stdout.String()
		if !strings.Contains(out, "Style Drift Report") {
			t.Errorf("expected report header in output: %s", out)
		}
		if !strings.Contains(out, "CANONICAL") || !strings.Contains(out, "rounded-md") {
			t.Errorf("expected canonical rounded-md in output: %s", out)
		}
		if !strings.Contains(out, "ROGUE DRIFT") || !strings.Contains(out, "rounded-2xl") {
			t.Errorf("expected rogue rounded-2xl in output: %s", out)
		}
	})

	t.Run("MarkdownReportToFile", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		reportPath := filepath.Join(tempDir, "drift-report.md")
		code := cli.RunDrift([]string{"-o", reportPath, "-f", "markdown", tempDir}, &stdout, &stderr)
		if code != cli.ExitClean {
			t.Fatalf("expected exit clean (0), got %d, stderr: %s", code, stderr.String())
		}
		data, err := os.ReadFile(filepath.Clean(reportPath))
		if err != nil {
			t.Fatalf("failed to read report file: %v", err)
		}
		content := string(data)
		if !strings.Contains(content, "# Charites Component-Scoped Style Drift Report") {
			t.Errorf("missing markdown header: %s", content)
		}
	})
}
