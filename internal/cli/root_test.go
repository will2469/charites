package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/will2469/charites/internal/cli"
)

func TestExecuteArgs(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedCode   int
		containsStdout string
		containsStderr string
	}{
		{
			name:           "flag --version",
			args:           []string{"--version"},
			expectedCode:   0,
			containsStdout: "charites version",
		},
		{
			name:           "flag -v",
			args:           []string{"-v"},
			expectedCode:   0,
			containsStdout: "charites version",
		},
		{
			name:           "subcommand version",
			args:           []string{"version"},
			expectedCode:   0,
			containsStdout: "charites version",
		},
		{
			name:           "flag --help",
			args:           []string{"--help"},
			expectedCode:   0,
			containsStdout: "Usage: charites",
		},
		{
			name:           "flag -h",
			args:           []string{"-h"},
			expectedCode:   0,
			containsStdout: "Usage: charites",
		},
		{
			name:           "subcommand help",
			args:           []string{"help"},
			expectedCode:   0,
			containsStdout: "Usage: charites",
		},
		{
			name:           "unknown command exits 2 to stderr",
			args:           []string{"unknown-command"},
			expectedCode:   2,
			containsStderr: `unknown command "unknown-command"`,
		},
		{
			name:           "unknown flag exits 2 to stderr",
			args:           []string{"--bogus-flag"},
			expectedCode:   2,
			containsStderr: "flag provided but not defined",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := cli.ExecuteArgs(tc.args, &stdout, &stderr)

			if code != tc.expectedCode {
				t.Fatalf("expected exit code %d, got %d (stderr: %s)", tc.expectedCode, code, stderr.String())
			}

			stdoutStr := stdout.String()
			stderrStr := stderr.String()

			if tc.expectedCode == 0 {
				if stderrStr != "" {
					t.Errorf("expected clean stderr on exit 0, got: %q", stderrStr)
				}
			} else {
				if stdoutStr != "" {
					t.Errorf("expected clean stdout on exit %d, got: %q", tc.expectedCode, stdoutStr)
				}
			}

			if tc.containsStdout != "" && !strings.Contains(stdoutStr, tc.containsStdout) {
				t.Errorf("stdout %q does not contain %q", stdoutStr, tc.containsStdout)
			}
			if tc.containsStderr != "" && !strings.Contains(stderrStr, tc.containsStderr) {
				t.Errorf("stderr %q does not contain %q", stderrStr, tc.containsStderr)
			}
		})
	}
}

func TestExecute_Entrypoint(t *testing.T) {
	// Memverifikasi bahwa trampoline Execute() dapat dipanggil dengan flag help tanpa panik
	code := cli.Execute([]string{"--help"})
	if code != cli.ExitClean {
		t.Errorf("expected Execute(--help) to return 0, got %d", code)
	}

	// 0 args defaults to scan .
	var stdout, stderr bytes.Buffer
	code0 := cli.ExecuteArgs([]string{}, &stdout, &stderr)
	if code0 != cli.ExitClean {
		t.Errorf("expected 0 args to return 0, got %d", code0)
	}
}

func TestExecuteArgs_PathResolution(t *testing.T) {
	paths := []string{
		".",
		"..",
		"./src",
		"../charites",
		"src/components",
		"file.astro",
		"file.tsx",
		"file.jsx",
		"file.ts",
		"file.js",
		"file.json",
		"file.yaml",
		"file.yml",
	}

	for _, p := range paths {
		t.Run("path_"+p, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			// We only care that isPath returns true and routes to RunScan (even if target does not exist, it routes to scan and returns exit 2 or 0)
			_ = cli.ExecuteArgs([]string{p}, &stdout, &stderr)
			// Ensure it did NOT output "unknown command"
			if strings.Contains(stderr.String(), "unknown command") {
				t.Errorf("expected path %q to be recognized as path, got unknown command", p)
			}
		})
	}
}
