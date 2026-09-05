package cli_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
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

func TestUninstallSuccess(t *testing.T) {
	tmpDir := t.TempDir()
	dummyBin := filepath.Join(tmpDir, "dummy-charites")
	if err := os.WriteFile(dummyBin, []byte("#!/bin/sh\necho test"), 0o600); err != nil {
		t.Fatalf("failed to create dummy binary: %v", err)
	}

	oldOsExec := cli.OsExecutable
	cli.OsExecutable = func() (string, error) {
		return dummyBin, nil
	}
	defer func() { cli.OsExecutable = oldOsExec }()

	var stdout, stderr bytes.Buffer
	code := cli.ExecuteArgs([]string{"uninstall", "--yes"}, &stdout, &stderr)

	if code != cli.ExitClean {
		t.Fatalf("expected clean exit 0, got %d (stderr: %s)", code, stderr.String())
	}

	if !strings.Contains(stdout.String(), "Charites uninstalled successfully.") {
		t.Errorf("expected success message, got: %q", stdout.String())
	}

	// Pastikan file telah terhapus
	if _, err := os.Stat(dummyBin); !os.IsNotExist(err) {
		t.Errorf("expected binary to be removed, but it still exists")
	}
}

func TestUninstallExecutableError(t *testing.T) {
	oldOsExec := cli.OsExecutable
	cli.OsExecutable = func() (string, error) {
		return "", errors.New("cannot determine executable")
	}
	defer func() { cli.OsExecutable = oldOsExec }()

	var stdout, stderr bytes.Buffer
	code := cli.ExecuteArgs([]string{"uninstall", "--yes"}, &stdout, &stderr)
	if code != cli.ExitOperational {
		t.Fatalf("expected code %d, got %d", cli.ExitOperational, code)
	}
}

func TestUninstallRemoveError(t *testing.T) {
	oldOsExec := cli.OsExecutable
	cli.OsExecutable = func() (string, error) {
		return "/non/existent/path/for/removal", nil
	}
	defer func() { cli.OsExecutable = oldOsExec }()

	var stdout, stderr bytes.Buffer
	code := cli.ExecuteArgs([]string{"uninstall"}, &stdout, &stderr)
	if code != cli.ExitOperational {
		t.Fatalf("expected code %d, got %d", cli.ExitOperational, code)
	}
}
