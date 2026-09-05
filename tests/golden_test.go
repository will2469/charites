package tests_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/will2469/charites/internal/cli"
)

var update = flag.Bool("update", false, "Update golden snapshot files")

var durationRegex = regexp.MustCompile(`in \d+ms\.`)

func normalizeJSONForGolden(t *testing.T, raw []byte) []byte {
	t.Helper()
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("failed to unmarshal JSON for golden normalization: %v\nRaw output:\n%s", err, string(raw))
	}

	if summary, ok := doc["summary"].(map[string]any); ok {
		summary["duration_ms"] = float64(0)
	}

	normalized, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal normalized JSON: %v", err)
	}
	return append(normalized, '\n')
}

func normalizeTextForGolden(raw string) string {
	normalized := strings.ReplaceAll(raw, "\r\n", "\n")
	normalized = durationRegex.ReplaceAllString(normalized, "in 0ms.")
	return normalized
}

func TestPipeline_GoldenSnapshots(t *testing.T) {
	if *update && (os.Getenv("CI") == "true" || os.Getenv("GITHUB_ACTIONS") == "true") {
		t.Fatal("FATAL: Golden snapshots MUST NOT be updated in CI environments!")
	}

	scenarios := []string{
		"clean",
		"opacity_violations",
		"config_override",
		"ignore_patterns",
	}

	goldenDir := filepath.Join("golden", "projects")
	if err := os.MkdirAll(goldenDir, 0o750); err != nil {
		t.Fatalf("failed to create golden directory: %v", err)
	}

	for _, sc := range scenarios {
		t.Run(sc, func(t *testing.T) {
			projectDir := filepath.Join("fixtures", "projects", sc)
			goldenJSONPath := filepath.Join(goldenDir, sc+".golden.json")
			goldenTxtPath := filepath.Join(goldenDir, sc+".golden.txt")

			// 1. Uji Format Dokumen JSON
			var stdoutJSON, stderrJSON bytes.Buffer
			_ = cli.ExecuteArgs([]string{"scan", projectDir, "-f", "json"}, &stdoutJSON, &stderrJSON)
			actualJSON := normalizeJSONForGolden(t, stdoutJSON.Bytes())

			if *update {
				if err := os.WriteFile(goldenJSONPath, actualJSON, 0o600); err != nil {
					t.Fatalf("failed to update golden JSON file: %v", err)
				}
				t.Logf("UPDATED golden JSON file: %s", goldenJSONPath)
			} else {
				expectedJSON, err := os.ReadFile(filepath.Clean(goldenJSONPath)) //nolint:gosec // controlled golden path
				if err != nil {
					t.Fatalf("failed to read golden JSON file %s: %v. Run with -update to generate.", goldenJSONPath, err)
				}
				if !bytes.Equal(actualJSON, expectedJSON) {
					t.Fatalf("Golden JSON mismatch on scenario %s!\nExpected:\n%s\nActual:\n%s", sc, string(expectedJSON), string(actualJSON))
				}
			}

			// 2. Uji Format Terminal Text (ColorNever)
			var stdoutTxt, stderrTxt bytes.Buffer
			_ = cli.ExecuteArgs([]string{"scan", projectDir, "-f", "inline", "--no-color"}, &stdoutTxt, &stderrTxt)
			actualTxt := normalizeTextForGolden(stdoutTxt.String())

			if *update {
				if err := os.WriteFile(goldenTxtPath, []byte(actualTxt), 0o600); err != nil {
					t.Fatalf("failed to update golden TXT file: %v", err)
				}
				t.Logf("UPDATED golden TXT file: %s", goldenTxtPath)
			} else {
				expectedTxtBytes, err := os.ReadFile(filepath.Clean(goldenTxtPath)) //nolint:gosec // controlled golden path
				if err != nil {
					t.Fatalf("failed to read golden TXT file %s: %v. Run with -update to generate.", goldenTxtPath, err)
				}
				expectedTxt := strings.ReplaceAll(string(expectedTxtBytes), "\r\n", "\n")
				if actualTxt != expectedTxt {
					t.Fatalf("Golden TXT mismatch on scenario %s!\nExpected:\n%s\nActual:\n%s", sc, expectedTxt, actualTxt)
				}
			}
		})
	}
}
