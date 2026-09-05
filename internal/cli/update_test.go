package cli_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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
