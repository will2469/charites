package e2e_test

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestBinarySmoke(t *testing.T) {
	binPath := "../../bin/charites"

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
			name:           "empty args runs scan",
			args:           []string{},
			expectedCode:   0,
			containsStdout: "problems found",
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
			// #nosec G204 -- smoke test executes test-controlled CLI binary with constant test arguments
			cmd := exec.Command(binPath, tc.args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()
			stdoutStr := stdout.String()
			stderrStr := stderr.String()

			if tc.expectedCode == 0 {
				if err != nil {
					t.Fatalf("expected exit 0, got err: %v, stderr: %s", err, stderrStr)
				}
				if stderrStr != "" {
					t.Errorf("expected clean stderr on exit 0, got: %s", stderrStr)
				}
			} else {
				if exitErr, ok := err.(*exec.ExitError); ok {
					if exitErr.ExitCode() != tc.expectedCode {
						t.Fatalf("expected exit %d, got %d (stderr: %s)", tc.expectedCode, exitErr.ExitCode(), stderrStr)
					}
				} else {
					t.Fatalf("expected exit code %d, got: %v", tc.expectedCode, err)
				}
			}

			if tc.containsStdout != "" && !strings.Contains(stdoutStr, tc.containsStdout) {
				t.Errorf("stdout does not contain %q: %s", tc.containsStdout, stdoutStr)
			}
			if tc.containsStderr != "" && !strings.Contains(stderrStr, tc.containsStderr) {
				t.Errorf("stderr does not contain %q: %s", tc.containsStderr, stderrStr)
			}
		})
	}
}
