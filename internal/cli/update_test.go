package cli_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/will2469/charites/internal/cli"
)

func TestUpdateAndUpgradeEquivalence(t *testing.T) {
	// Mock server that returns no new updates
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"tag_name": "v" + cli.Version,
			"assets":   []any{},
		})
	}))
	defer server.Close()

	oldURL := cli.CharitesUpdateURL
	cli.CharitesUpdateURL = server.URL
	defer func() { cli.CharitesUpdateURL = oldURL }()

	commands := []string{"update", "upgrade"}
	for _, cmd := range commands {
		t.Run("subcommand "+cmd, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := cli.ExecuteArgs([]string{cmd}, &stdout, &stderr)

			if code != cli.ExitClean {
				t.Fatalf("expected code %d, got %d, stderr: %s", cli.ExitClean, code, stderr.String())
			}
			out := stdout.String()
			if !strings.Contains(out, "No update found. Charites is up to date.") {
				t.Errorf("expected 'No update found. Charites is up to date.', got: %q", out)
			}
		})
	}
}

func TestUpdateCheckOnly(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"tag_name": "v99.99.99",
			"assets": []map[string]any{
				{
					"name":                 fmt.Sprintf("charites-%s-%s", runtime.GOOS, runtime.GOARCH),
					"browser_download_url": "http://example.com/binary",
				},
			},
		})
	}))
	defer server.Close()

	oldURL := cli.CharitesUpdateURL
	cli.CharitesUpdateURL = server.URL
	defer func() { cli.CharitesUpdateURL = oldURL }()

	var stdout, stderr bytes.Buffer
	code := cli.ExecuteArgs([]string{"update", "--check"}, &stdout, &stderr)

	if code != cli.ExitClean {
		t.Fatalf("expected exit code %d, got %d", cli.ExitClean, code)
	}

	if !strings.Contains(stdout.String(), "Update available: v99.99.99") {
		t.Errorf("expected update available notice, got: %q", stdout.String())
	}
}

func TestUpdateNetworkFailureFallback(t *testing.T) {
	// Point to unreachable port/URL
	oldURL := cli.CharitesUpdateURL
	cli.CharitesUpdateURL = "http://127.0.0.1:0/unreachable"
	defer func() { cli.CharitesUpdateURL = oldURL }()

	var stdout, stderr bytes.Buffer
	code := cli.ExecuteArgs([]string{"update"}, &stdout, &stderr)

	if code != cli.ExitClean {
		t.Fatalf("expected exit code %d, got %d", cli.ExitClean, code)
	}

	if !strings.Contains(stdout.String(), "No update found. Charites is up to date.") {
		t.Errorf("expected clean fallback, got: %q", stdout.String())
	}
}

func TestUpdateAssetNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"tag_name": "v99.99.99",
			"assets": []map[string]any{
				{
					"name":                 "charites-unknownos-unknownarch",
					"browser_download_url": "http://example.com/binary",
				},
			},
		})
	}))
	defer server.Close()

	oldURL := cli.CharitesUpdateURL
	cli.CharitesUpdateURL = server.URL
	defer func() { cli.CharitesUpdateURL = oldURL }()

	var stdout, stderr bytes.Buffer
	code := cli.ExecuteArgs([]string{"update"}, &stdout, &stderr)

	if code != cli.ExitOperational {
		t.Fatalf("expected exit code %d, got %d", cli.ExitOperational, code)
	}
	if !strings.Contains(stderr.String(), "not found in release") {
		t.Errorf("expected asset not found error in stderr, got: %q", stderr.String())
	}
}

func TestUpdateDownloadFailure(t *testing.T) {
	dlServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer dlServer.Close()

	assetName := fmt.Sprintf("charites-%s-%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		assetName += ".exe"
	}

	releaseServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"tag_name": "v99.99.99",
			"assets": []map[string]any{
				{
					"name":                 assetName,
					"browser_download_url": dlServer.URL + "/binary",
				},
			},
		})
	}))
	defer releaseServer.Close()

	oldURL := cli.CharitesUpdateURL
	cli.CharitesUpdateURL = releaseServer.URL
	defer func() { cli.CharitesUpdateURL = oldURL }()

	var stdout, stderr bytes.Buffer
	code := cli.ExecuteArgs([]string{"update"}, &stdout, &stderr)

	if code != cli.ExitOperational {
		t.Fatalf("expected exit code %d, got %d", cli.ExitOperational, code)
	}
	if !strings.Contains(stderr.String(), "failed to download update") {
		t.Errorf("expected download failure error, got: %q", stderr.String())
	}
}

func TestUpdateSuccess(t *testing.T) {
	dummyContent := []byte("#!/bin/sh\necho updated-charites")

	dlServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(dummyContent)
	}))
	defer dlServer.Close()

	assetName := fmt.Sprintf("charites-%s-%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		assetName += ".exe"
	}

	releaseServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"tag_name": "v99.99.99",
			"assets": []map[string]any{
				{
					"name":                 assetName,
					"browser_download_url": dlServer.URL + "/binary",
				},
			},
		})
	}))
	defer releaseServer.Close()

	oldURL := cli.CharitesUpdateURL
	cli.CharitesUpdateURL = releaseServer.URL
	defer func() { cli.CharitesUpdateURL = oldURL }()

	tmpDir := t.TempDir()
	currentBin := filepath.Join(tmpDir, "current-charites")
	if err := os.WriteFile(currentBin, []byte("#!/bin/sh\necho old"), 0o600); err != nil {
		t.Fatalf("failed to create current binary: %v", err)
	}

	oldOsExec := cli.OsExecutable
	cli.OsExecutable = func() (string, error) {
		return currentBin, nil
	}
	defer func() { cli.OsExecutable = oldOsExec }()

	var stdout, stderr bytes.Buffer
	code := cli.ExecuteArgs([]string{"update"}, &stdout, &stderr)

	if code != cli.ExitClean {
		t.Fatalf("expected exit code %d, got %d (stderr: %s)", cli.ExitClean, code, stderr.String())
	}

	if !strings.Contains(stdout.String(), "Charites updated to v99.99.99 successfully.") {
		t.Errorf("expected update success message, got: %q", stdout.String())
	}

	updatedBytes, err := os.ReadFile(filepath.Clean(currentBin)) //nolint:gosec // controlled test path
	if err != nil {
		t.Fatalf("failed to read updated binary: %v", err)
	}
	if !bytes.Equal(updatedBytes, dummyContent) {
		t.Errorf("binary content not updated properly")
	}
}

func TestUpdateExecutableError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assetName := fmt.Sprintf("charites-%s-%s", runtime.GOOS, runtime.GOARCH)
		if runtime.GOOS == "windows" {
			assetName += ".exe"
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"tag_name": "v99.99.99",
			"assets": []map[string]any{
				{
					"name":                 assetName,
					"browser_download_url": "http://example.com/binary",
				},
			},
		})
	}))
	defer server.Close()

	oldURL := cli.CharitesUpdateURL
	cli.CharitesUpdateURL = server.URL
	defer func() { cli.CharitesUpdateURL = oldURL }()

	oldOsExec := cli.OsExecutable
	cli.OsExecutable = func() (string, error) {
		return "", errors.New("cannot determine executable")
	}
	defer func() { cli.OsExecutable = oldOsExec }()

	var stdout, stderr bytes.Buffer
	code := cli.ExecuteArgs([]string{"update"}, &stdout, &stderr)
	if code != cli.ExitOperational {
		t.Fatalf("expected code %d, got %d", cli.ExitOperational, code)
	}
}

func TestUpdateInvalidFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := cli.ExecuteArgs([]string{"update", "--invalid-flag"}, &stdout, &stderr)
	if code != cli.ExitOperational {
		t.Fatalf("expected code %d, got %d", cli.ExitOperational, code)
	}
}

func TestUpdateInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("not-json"))
	}))
	defer server.Close()

	oldURL := cli.CharitesUpdateURL
	cli.CharitesUpdateURL = server.URL
	defer func() { cli.CharitesUpdateURL = oldURL }()

	var stdout, stderr bytes.Buffer
	code := cli.ExecuteArgs([]string{"update"}, &stdout, &stderr)
	if code != cli.ExitClean {
		t.Fatalf("expected code %d, got %d", cli.ExitClean, code)
	}
	if !strings.Contains(stdout.String(), "No update found. Charites is up to date.") {
		t.Errorf("expected clean fallback on invalid JSON, got: %q", stdout.String())
	}
}
