package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/will2469/charites/internal/cli"
)

func TestDoctorCommand(t *testing.T) {
	commands := []string{"doctor", "-doctor", "--doctor"}
	for _, cmd := range commands {
		t.Run("command "+cmd, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := cli.ExecuteArgs([]string{cmd}, &stdout, &stderr)

			// Doctor exits 0 when clean or with warnings only
			if code != cli.ExitClean && code != cli.ExitOperational {
				t.Fatalf("unexpected exit code %d, stderr: %s", code, stderr.String())
			}

			out := stdout.String()
			if !strings.Contains(out, "Charites Doctor") {
				t.Errorf("expected header 'Charites Doctor', got: %q", out)
			}
			if !strings.Contains(out, "Host Platform") {
				t.Errorf("expected 'Host Platform' check, got: %q", out)
			}
			if !strings.Contains(out, "Doctor Verdict") {
				t.Errorf("expected 'Doctor Verdict', got: %q", out)
			}
		})
	}
}

func TestDoctorInvalidFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := cli.ExecuteArgs([]string{"doctor", "--invalid-flag"}, &stdout, &stderr)
	if code != cli.ExitOperational {
		t.Fatalf("expected code %d, got %d", cli.ExitOperational, code)
	}
	if !strings.Contains(stderr.String(), "flag provided but not defined") {
		t.Errorf("expected invalid flag error in stderr, got: %q", stderr.String())
	}
}
