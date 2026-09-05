package scanner_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/will2469/charites/internal/config"
	"github.com/will2469/charites/internal/scanner"
)

func TestWalker_DirectTargetSafety(t *testing.T) {
	matcher := config.NewIgnoreMatcher(nil)
	w := scanner.NewWalker(matcher, nil)

	forbiddenTargets := []string{
		"node_modules/react/index.d.ts",
		"./node_modules/pkg/Button.tsx",
		".git/hooks/pre-commit",
		"dist/bundle.js",
		"coverage/lcov.info",
	}

	for _, target := range forbiddenTargets {
		jobs := make(chan string, 10)
		err := w.Walk(context.Background(), target, jobs)
		close(jobs)

		if err == nil {
			t.Errorf("expected error for forbidden direct target %q, got nil", target)
		}
		if len(jobs) != 0 {
			t.Errorf("expected 0 jobs queued for forbidden target %q, got %d", target, len(jobs))
		}
	}
}

func TestWalker_SingleFileTarget(t *testing.T) {
	tmpDir := t.TempDir()
	validFile := filepath.Join(tmpDir, "Valid.tsx")
	invalidExt := filepath.Join(tmpDir, "style.css")

	_ = os.WriteFile(validFile, []byte("<div />"), 0o600)
	_ = os.WriteFile(invalidExt, []byte("body {}"), 0o600)

	matcher := config.NewIgnoreMatcher(nil)
	w := scanner.NewWalker(matcher, []string{"tsx", ".astro"})

	// Target berkas valid
	jobsValid := make(chan string, 10)
	if err := w.Walk(context.Background(), validFile, jobsValid); err != nil {
		t.Fatalf("unexpected error for single valid file: %v", err)
	}
	close(jobsValid)

	queued := make([]string, 0, 10)
	for j := range jobsValid {
		queued = append(queued, j)
	}
	if len(queued) != 1 || queued[0] != validFile {
		t.Errorf("expected [%s], got %+v", validFile, queued)
	}

	// Target ekstensi tidak didukung
	jobsInvalid := make(chan string, 10)
	if err := w.Walk(context.Background(), invalidExt, jobsInvalid); err != nil {
		t.Fatalf("unexpected error for unsupported ext: %v", err)
	}
	close(jobsInvalid)
	if len(jobsInvalid) != 0 {
		t.Errorf("expected 0 jobs for invalid ext, got %d", len(jobsInvalid))
	}
}

func TestWalker_MaxFileSizeGuard(t *testing.T) {
	tmpDir := t.TempDir()
	largeFile := filepath.Join(tmpDir, "Large.tsx")

	// Buat berkas melebihi 10 MB (10 MB + 1 KB)
	f, err := os.Create(filepath.Clean(largeFile)) //nolint:gosec // controlled test file
	if err != nil {
		t.Fatalf("failed to create dummy large file: %v", err)
	}
	_ = f.Truncate(scanner.MaxScanFileSize + 1024)
	_ = f.Close()

	matcher := config.NewIgnoreMatcher(nil)
	w := scanner.NewWalker(matcher, nil)

	jobs := make(chan string, 10)
	if err := w.Walk(context.Background(), tmpDir, jobs); err != nil {
		t.Fatalf("unexpected walk error: %v", err)
	}
	close(jobs)

	if len(jobs) != 0 {
		t.Errorf("expected large file (>10MB) to be skipped from queue, got %d jobs", len(jobs))
	}
}

func TestWalker_SymlinkSafety(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	_ = os.MkdirAll(subDir, 0o750)

	targetFile := filepath.Join(subDir, "Component.tsx")
	_ = os.WriteFile(targetFile, []byte("<div />"), 0o600)

	// Buat symlink siklis: cyclic_link -> tmpDir
	cyclicLink := filepath.Join(subDir, "cyclic_link")
	if err := os.Symlink(tmpDir, cyclicLink); err != nil {
		t.Skip("symlinks not supported on current environment")
	}

	matcher := config.NewIgnoreMatcher(nil)
	w := scanner.NewWalker(matcher, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	jobs := make(chan string, 100)
	err := w.Walk(ctx, tmpDir, jobs)
	close(jobs)

	if err != nil {
		t.Fatalf("expected clean walk without infinite loop on cyclic symlink, got error: %v", err)
	}

	found := make([]string, 0, 10)
	for j := range jobs {
		found = append(found, j)
	}

	if len(found) != 1 || found[0] != targetFile {
		t.Errorf("expected only original file %q to be queued, got %+v", targetFile, found)
	}
}

func TestWalker_EarlyDirectoryPruning(t *testing.T) {
	tmpDir := t.TempDir()
	ignoredDir := filepath.Join(tmpDir, "legacy_vendor")
	_ = os.MkdirAll(ignoredDir, 0o750)
	_ = os.WriteFile(filepath.Join(ignoredDir, "Old.tsx"), []byte("<div />"), 0o600)

	keptDir := filepath.Join(tmpDir, "src")
	_ = os.MkdirAll(keptDir, 0o750)
	_ = os.WriteFile(filepath.Join(keptDir, "New.tsx"), []byte("<div />"), 0o600)

	matcher := config.NewIgnoreMatcher([]string{"legacy_vendor/"})
	w := scanner.NewWalker(matcher, nil)

	jobs := make(chan string, 10)
	if err := w.Walk(context.Background(), tmpDir, jobs); err != nil {
		t.Fatalf("unexpected walk error: %v", err)
	}
	close(jobs)

	queued := make([]string, 0, 10)
	for j := range jobs {
		queued = append(queued, j)
	}

	if len(queued) != 1 || filepath.Base(queued[0]) != "New.tsx" {
		t.Errorf("expected only New.tsx to be queued, got %+v", queued)
	}
}

func TestWalker_ContextCancellation(t *testing.T) {
	tmpDir := t.TempDir()
	for i := 0; i < 50; i++ {
		_ = os.WriteFile(filepath.Join(tmpDir, filepath.Clean(filepath.Join(".", "file.tsx"))), []byte("<div />"), 0o600)
	}

	matcher := config.NewIgnoreMatcher(nil)
	w := scanner.NewWalker(matcher, nil)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Batalkan sebelum atau saat mulai

	jobs := make(chan string, 1)
	err := w.Walk(ctx, tmpDir, jobs)
	if err == nil {
		t.Error("expected context canceled error, got nil")
	}
}
