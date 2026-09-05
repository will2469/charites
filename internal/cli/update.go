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

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, CharitesUpdateURL, nil)
	if err != nil {
		_, _ = fmt.Fprintln(stdout, "No update found. Charites is up to date.")
		return ExitClean
	}
	req.Header.Set("User-Agent", "Charites-CLI")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		if resp != nil {
			_ = resp.Body.Close()
		}
		_, _ = fmt.Fprintln(stdout, "No update found. Charites is up to date.")
		return ExitClean
	}
	defer func() { _ = resp.Body.Close() }()

	var rel githubRelease
	if decodeErr := json.NewDecoder(resp.Body).Decode(&rel); decodeErr != nil || rel.TagName == "" {
		_, _ = fmt.Fprintln(stdout, "No update found. Charites is up to date.")
		return ExitClean
	}

	latestVersion := strings.TrimPrefix(rel.TagName, "v")
	currentVersion := strings.TrimPrefix(Version, "v")

	// Jika versi rilis sama atau saat ini berjalan pada versi dev/terbaru
	if latestVersion == currentVersion || currentVersion == latestVersion || latestVersion == "" {
		_, _ = fmt.Fprintln(stdout, "No update found. Charites is up to date.")
		return ExitClean
	}

	if checkOnly {
		_, _ = fmt.Fprintf(stdout, "Update available: %s (current: %s). Run 'charites update' to upgrade.\n", rel.TagName, Version)
		return ExitClean
	}

	// Cari asset biner yang cocok dengan GOOS dan GOARCH
	expectedAsset := fmt.Sprintf("charites-%s-%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		expectedAsset += ".exe"
	}

	var downloadURL string
	for _, asset := range rel.Assets {
		if strings.EqualFold(asset.Name, expectedAsset) {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		_, _ = fmt.Fprintf(stderr, "charites: error: release asset %q not found in release %s.\n", expectedAsset, rel.TagName)
		return ExitOperational
	}

	// Unduh biner baru
	dlReq, reqErr := http.NewRequest(http.MethodGet, downloadURL, nil)
	if reqErr != nil {
		_, _ = fmt.Fprintf(stderr, "charites: error: failed to create download request: %v\n", reqErr)
		return ExitOperational
	}
	dlReq.Header.Set("User-Agent", "Charites-CLI")

	dlResp, dlErr := client.Do(dlReq)
	if dlErr != nil || dlResp.StatusCode != http.StatusOK {
		if dlResp != nil {
			_ = dlResp.Body.Close()
		}
		_, _ = fmt.Fprintf(stderr, "charites: error: failed to download update: %v\n", dlErr)
		return ExitOperational
	}
	defer func() { _ = dlResp.Body.Close() }()

	// Dapatkan path eksekutabel saat ini
	execPath, execErr := os.Executable()
	if execErr != nil {
		_, _ = fmt.Fprintf(stderr, "charites: error: failed to determine executable path: %v\n", execErr)
		return ExitOperational
	}
	execPath, symErr := filepath.EvalSymlinks(execPath)
	if symErr != nil {
		_, _ = fmt.Fprintf(stderr, "charites: error: failed to resolve executable symlinks: %v\n", symErr)
		return ExitOperational
	}

	// Tulis ke berkas sementara di direktori yang sama agar rename bersifat atomik
	execDir := filepath.Dir(execPath)
	tmpFile, tmpErr := os.CreateTemp(execDir, "charites-update-*")
	if tmpErr != nil {
		_, _ = fmt.Fprintf(stderr, "charites: error: failed to create temporary update file: %v\n", tmpErr)
		return ExitOperational
	}
	tmpName := tmpFile.Name()

	_, copyErr := io.Copy(tmpFile, dlResp.Body)
	closeErr := tmpFile.Close()
	if copyErr != nil || closeErr != nil {
		_ = os.Remove(tmpName)
		_, _ = fmt.Fprintf(stderr, "charites: error: failed to write update file.\n")
		return ExitOperational
	}

	if chmodErr := os.Chmod(tmpName, 0o755); chmodErr != nil { //nolint:gosec // binary executable permissions
		_ = os.Remove(tmpName)
		_, _ = fmt.Fprintf(stderr, "charites: error: failed to set executable permissions: %v\n", chmodErr)
		return ExitOperational
	}

	// Timpa biner secara atomik
	if renErr := os.Rename(tmpName, execPath); renErr != nil {
		_ = os.Remove(tmpName)
		_, _ = fmt.Fprintf(stderr, "charites: error: failed to replace binary (permission denied): %v\n", renErr)
		return ExitOperational
	}

	_, _ = fmt.Fprintf(stdout, "Charites updated to %s successfully.\n", rel.TagName)
	return ExitClean
}
