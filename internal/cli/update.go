package cli

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// CharitesUpdateURL mendefinisikan endpoint GitHub API untuk rilis terbaru.
// Dapat di-override dalam pengujian unit via CharitesUpdateURL.
var CharitesUpdateURL = "https://api.github.com/repos/will2469/charites/releases/latest"

type githubRelease struct {
	TagName string        `json:"tag_name"`
	Assets  []githubAsset `json:"assets"`
}

type githubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func fetchLatestRelease(client *http.Client) (*githubRelease, error) {
	req, err := http.NewRequest(http.MethodGet, CharitesUpdateURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Charites-CLI")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		if resp != nil {
			_ = resp.Body.Close()
		}
		return nil, fmt.Errorf("failed to fetch release")
	}
	defer func() { _ = resp.Body.Close() }()

	var rel githubRelease
	if decodeErr := json.NewDecoder(resp.Body).Decode(&rel); decodeErr != nil || rel.TagName == "" {
		return nil, fmt.Errorf("invalid release format")
	}
	return &rel, nil
}

func findAssetDownloadURL(rel *githubRelease) string {
	expectedAsset := fmt.Sprintf("charites-%s-%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		expectedAsset += ".exe"
	}

	for _, asset := range rel.Assets {
		if strings.EqualFold(asset.Name, expectedAsset) {
			return asset.BrowserDownloadURL
		}
	}
	return ""
}

func downloadAndReplaceBinary(client *http.Client, downloadURL string) error {
	dlReq, reqErr := http.NewRequest(http.MethodGet, downloadURL, nil)
	if reqErr != nil {
		return fmt.Errorf("failed to create download request: %w", reqErr)
	}
	dlReq.Header.Set("User-Agent", "Charites-CLI")

	dlResp, dlErr := client.Do(dlReq)
	if dlErr != nil || dlResp.StatusCode != http.StatusOK {
		if dlResp != nil {
			_ = dlResp.Body.Close()
		}
		return fmt.Errorf("failed to download update: %w", dlErr)
	}
	defer func() { _ = dlResp.Body.Close() }()

	execPath, execErr := OsExecutable()
	if execErr != nil {
		return fmt.Errorf("failed to determine executable path: %w", execErr)
	}
	execPath, symErr := filepath.EvalSymlinks(execPath)
	if symErr != nil {
		return fmt.Errorf("failed to resolve executable symlinks: %w", symErr)
	}

	execDir := filepath.Dir(execPath)
	tmpFile, tmpErr := os.CreateTemp(execDir, "charites-update-*")
	if tmpErr != nil {
		return fmt.Errorf("failed to create temporary update file: %w", tmpErr)
	}
	tmpName := tmpFile.Name()

	_, copyErr := io.Copy(tmpFile, dlResp.Body)
	closeErr := tmpFile.Close()
	if copyErr != nil || closeErr != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("failed to write update file")
	}

	// #nosec G302 -- binary executable permissions
	if chmodErr := os.Chmod(tmpName, 0o755); chmodErr != nil { //nolint:gosec
		_ = os.Remove(tmpName)
		return fmt.Errorf("failed to set executable permissions: %w", chmodErr)
	}

	if renErr := os.Rename(tmpName, execPath); renErr != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("failed to replace binary (permission denied): %w", renErr)
	}

	return nil
}

// RunUpdate mengorkestrasi pengecekan dan pembaruan biner ke versi rilis terbaru.
func RunUpdate(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("update", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.Usage = func() {}

	var checkOnly bool
	fs.BoolVar(&checkOnly, "check", false, "Hanya periksa pembaruan tanpa mengunduh")
	fs.BoolVar(&checkOnly, "c", false, "Hanya periksa pembaruan (shorthand)")

	if parseErr := fs.Parse(args); parseErr != nil {
		if errors.Is(parseErr, flag.ErrHelp) {
			_, _ = fmt.Fprint(stdout, UsageString())
			return ExitClean
		}
		_, _ = fmt.Fprintf(stderr, "charites: error: %v. Run 'charites --help' for usage.\n", parseErr)
		return ExitOperational
	}

	client := &http.Client{Timeout: 5 * time.Second}
	rel, err := fetchLatestRelease(client)
	if err != nil {
		_, _ = fmt.Fprintln(stdout, "No update found. Charites is up to date.")
		return ExitClean
	}

	latestVersion := strings.TrimPrefix(rel.TagName, "v")
	currentVersion := strings.TrimPrefix(Version, "v")

	if latestVersion == currentVersion || latestVersion == "" {
		_, _ = fmt.Fprintln(stdout, "No update found. Charites is up to date.")
		return ExitClean
	}

	if checkOnly {
		_, _ = fmt.Fprintf(stdout, "Update available: %s (current: %s). Run 'charites update' to upgrade.\n", rel.TagName, Version)
		return ExitClean
	}

	downloadURL := findAssetDownloadURL(rel)
	if downloadURL == "" {
		_, _ = fmt.Fprintf(stderr, "charites: error: release asset for %s-%s not found in release %s.\n", runtime.GOOS, runtime.GOARCH, rel.TagName)
		return ExitOperational
	}

	if err := downloadAndReplaceBinary(client, downloadURL); err != nil {
		_, _ = fmt.Fprintf(stderr, "charites: error: %v\n", err)
		return ExitOperational
	}

	_, _ = fmt.Fprintf(stdout, "Charites updated to %s successfully.\n", rel.TagName)
	return ExitClean
}
