package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/will2469/charites/internal/cli"
)

func TestUninstallHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := cli.ExecuteArgs([]string{"uninstall", "--help"}, &stdout, &stderr)

	if code != cli.ExitClean {
		t.Fatalf("expected code %d, got %d, stderr: %s", cli.ExitClean, code, stderr.String())
	}

	if !strings.Contains(stdout.String(), "Usage: charites") {
		t.Errorf("expected usage string in stdout, got: %q", stdout.String())
	}
}

func TestUninstallUnknownFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := cli.ExecuteArgs([]string{"uninstall", "--unknown-flag"}, &stdout, &stderr)

	if code != cli.ExitOperational {
		t.Fatalf("expected code %d, got %d", cli.ExitOperational, code)
	}

	if !strings.Contains(stderr.String(), "flag provided but not defined") {
		t.Errorf("expected flag error in stderr, got: %q", stderr.String())
	}
}
