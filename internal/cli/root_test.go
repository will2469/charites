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

func TestExecuteArgs_WikiCommand(t *testing.T) {
	tmpDir := t.TempDir()
	var stdout, stderr bytes.Buffer

	code := cli.ExecuteArgs([]string{"wiki", tmpDir}, &stdout, &stderr)
	if code != cli.ExitClean {
		t.Fatalf("expected wiki command to return 0, got %d, stderr: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "wiki documentation successfully generated") {
		t.Errorf("expected success message on stdout, got: %s", stdout.String())
	}
}

func TestExecuteWithStreams_MCPCommand(t *testing.T) {
	tmpDir := t.TempDir()
	var stdout, stderr bytes.Buffer
	// empty stdin triggers EOF -> graceful exit 0
	inR := strings.NewReader("")

	code := cli.ExecuteWithStreams([]string{"mcp", "--workspace", tmpDir}, inR, &stdout, &stderr)
	if code != cli.ExitClean {
		t.Fatalf("expected mcp command to return 0 on EOF, got %d, stderr: %s", code, stderr.String())
	}

	// MCP help
	var hStdout, hStderr bytes.Buffer
	codeHelp := cli.ExecuteWithStreams([]string{"mcp", "--help"}, inR, &hStdout, &hStderr)
	if codeHelp != cli.ExitClean {
		t.Fatalf("expected mcp --help to return 0, got %d", codeHelp)
	}
	if !strings.Contains(hStdout.String(), "Usage: charites mcp") {
		t.Errorf("expected mcp help string, got: %s", hStdout.String())
	}
}

func TestUsageString_HasMCPAndWiki(t *testing.T) {
	usage := cli.UsageString()
	if !strings.Contains(usage, "mcp") {
		t.Errorf("expected usage string to contain 'mcp'")
	}
	if !strings.Contains(usage, "wiki") {
		t.Errorf("expected usage string to contain 'wiki'")
	}
}
