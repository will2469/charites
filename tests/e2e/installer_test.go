package e2e_test

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestInstaller_Shellcheck memverifikasi skrip scripts/install.sh lulus linter tanpa peringatan bashism.
func TestInstaller_Shellcheck(t *testing.T) {
	if _, err := exec.LookPath("shellcheck"); err != nil {
		t.Skip("shellcheck not found in PATH, skipping lint test")
	}

	repoRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("failed to resolve repo root: %v", err)
	}
	scriptPath := filepath.Join(repoRoot, "scripts", "install.sh")

	// #nosec G204 -- test executes test-controlled shellcheck
	cmd := exec.Command("shellcheck", "-s", "sh", scriptPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("shellcheck failed: %v\nOutput:\n%s", err, string(out))
	}
}

// TestInstaller_SecurityInvariants memverifikasi bahwa skrip instalasi menolak checksum palsu dan tarball berisi path traversal.
func TestInstaller_SecurityInvariants(t *testing.T) {
	repoRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("failed to resolve repo root: %v", err)
	}
	scriptPath := filepath.Join(repoRoot, "scripts", "install.sh")

	// 1. Uji Coba Penolakan Checksum Palsu: Server mock menyajikan tarball dengan hash berbeda dari checksums.txt
	t.Run("RejectCorruptChecksum", func(t *testing.T) {
		dummyTarball := createDummyTarball(t, "charites", []byte("#!/bin/sh\necho test\n"))
		wrongChecksum := fmt.Sprintf("%064d  charites_v1.0.0_linux_amd64.tar.gz\n", 9999)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "checksums.txt") {
				_, _ = w.Write([]byte(wrongChecksum))
				return
			}
			if strings.HasSuffix(r.URL.Path, ".tar.gz") {
				_, _ = w.Write(dummyTarball)
				return
			}
			http.NotFound(w, r)
		}))
		defer server.Close()

		// Modifikasi sementara skrip via sed/awk atau jalankan sub-shell dengan GITHUB_URL di-override
		installDir := t.TempDir()
		// #nosec G204 -- test executes test-controlled installer script
		cmd := exec.Command("/bin/sh", scriptPath, "v1.0.0")
		cmd.Env = append(os.Environ(),
			"GITHUB_URL="+server.URL,
			"CHARITES_INSTALL_DIR="+installDir,
			"CHARITES_VERSION=v1.0.0",
		)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err == nil {
			t.Fatalf("expected installer to reject corrupted checksum, but it succeeded!")
		}
		if !strings.Contains(stderr.String(), "Checksum mismatch") {
			t.Errorf("expected stderr to contain 'Checksum mismatch', got: %s", stderr.String())
		}
	})

	// 2. Uji Coba Pencegahan Path Traversal: Server mock menyajikan tarball dengan path traversal (../usr/bin)
	t.Run("RejectPathTraversalArchive", func(t *testing.T) {
		maliciousTarball := createDummyTarball(t, "../usr/bin/evil", []byte("evil payload"))
		hash := fmt.Sprintf("%x", sha256.Sum256(maliciousTarball))
		checksumContent := fmt.Sprintf("%s  charites_v1.0.0_linux_amd64.tar.gz\n", hash)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "checksums.txt") {
				_, _ = w.Write([]byte(checksumContent))
				return
			}
			if strings.HasSuffix(r.URL.Path, ".tar.gz") {
				_, _ = w.Write(maliciousTarball)
				return
			}
			http.NotFound(w, r)
		}))
		defer server.Close()

		installDir := t.TempDir()
		// #nosec G204 -- test executes test-controlled installer script
		cmd := exec.Command("/bin/sh", scriptPath, "v1.0.0")
		cmd.Env = append(os.Environ(),
			"GITHUB_URL="+server.URL,
			"CHARITES_INSTALL_DIR="+installDir,
			"CHARITES_VERSION=v1.0.0",
		)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err == nil {
			t.Fatalf("expected installer to reject path traversal archive, but it succeeded!")
		}
		if !strings.Contains(stderr.String(), "path traversal") {
			t.Errorf("expected stderr to mention path traversal, got: %s", stderr.String())
		}
	})
}

func createDummyTarball(t *testing.T, filename string, content []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	hdr := &tar.Header{
		Name: filename,
		Mode: 0755,
		Size: int64(len(content)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		t.Fatalf("failed to write tar header: %v", err)
	}
	if _, err := tw.Write(content); err != nil {
		t.Fatalf("failed to write tar body: %v", err)
	}
	_ = tw.Close()
	_ = gw.Close()

	return buf.Bytes()
}
