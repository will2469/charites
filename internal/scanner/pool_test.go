package scanner_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/will2469/charites/internal/config"
	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/scanner"
)

type mockFileAnalyzer struct {
	processedCount int64
	delay          time.Duration
}

func (m *mockFileAnalyzer) AnalyzeFile(path string) ([]ir.Diagnostic, error) {
	atomic.AddInt64(&m.processedCount, 1)
	if m.delay > 0 {
		time.Sleep(m.delay)
	}
	return []ir.Diagnostic{
		{
			File:     path,
			Line:     1,
			Column:   1,
			Rule:     "theme.hardcode-opacity-color",
			Severity: ir.SeverityWarn,
			Message:  "mock violation in " + filepath.Base(path),
		},
	}, nil
}

func TestPool_WorkerLimits(t *testing.T) {
	pLow := scanner.NewPool(-5)
	if pLow.Workers() < 1 {
		t.Errorf("expected pool workers >= 1, got %d", pLow.Workers())
	}

	pHigh := scanner.NewPool(500)
	if pHigh.Workers() != 256 {
		t.Errorf("expected pool workers clamped to 256, got %d", pHigh.Workers())
	}
}

func TestPool_ConcurrentExecutionAndRace(t *testing.T) {
	tmpDir := t.TempDir()
	numFiles := 200

	for i := 0; i < numFiles; i++ {
		fname := filepath.Join(tmpDir, fmt.Sprintf("Component%03d.tsx", i))
		_ = os.WriteFile(fname, []byte("<div />"), 0o600)
	}

	matcher := config.NewIgnoreMatcher(nil)
	w := scanner.NewWalker(matcher, nil)
	pool := scanner.NewPool(8)

	analyzer := &mockFileAnalyzer{}
	diags, err := pool.Run(context.Background(), w, tmpDir, analyzer)
	if err != nil {
		t.Fatalf("pool.Run failed: %v", err)
	}

	if len(diags) != numFiles {
		t.Fatalf("expected %d diagnostics, got %d", numFiles, len(diags))
	}

	if atomic.LoadInt64(&analyzer.processedCount) != int64(numFiles) {
		t.Errorf("expected %d processed files, got %d", numFiles, analyzer.processedCount)
	}

	// Verifikasi deterministic ordering (File menaik)
	for i := 1; i < len(diags); i++ {
		if diags[i].File < diags[i-1].File {
			t.Errorf("diagnostics not deterministically sorted: [%s] < [%s]", diags[i].File, diags[i-1].File)
		}
	}
}

func TestPool_ContextCancellationCleanExit(t *testing.T) {
	tmpDir := t.TempDir()
	numFiles := 100

	for i := 0; i < numFiles; i++ {
		fname := filepath.Join(tmpDir, fmt.Sprintf("Slow%03d.tsx", i))
		_ = os.WriteFile(fname, []byte("<div />"), 0o600)
	}

	matcher := config.NewIgnoreMatcher(nil)
	w := scanner.NewWalker(matcher, nil)
	pool := scanner.NewPool(4)

	analyzer := &mockFileAnalyzer{delay: 20 * time.Millisecond}

	ctx, cancel := context.WithCancel(context.Background())
	// Batalkan segera setelah peluncuran
	go func() {
		time.Sleep(30 * time.Millisecond)
		cancel()
	}()

	diags, err := pool.Run(ctx, w, tmpDir, analyzer)
	if err == nil {
		t.Fatal("expected context cancellation error, got nil")
	}

	// Hasil parsial wajib dibuang
	if len(diags) != 0 {
		t.Errorf("expected partial diagnostics to be discarded on cancellation, got %d", len(diags))
	}
}

func TestPool_WalkErrorPropagation(t *testing.T) {
	matcher := config.NewIgnoreMatcher(nil)
	w := scanner.NewWalker(matcher, nil)
	pool := scanner.NewPool(2)

	// Target terlarang (direct-target safety) memicu error di walker
	_, err := pool.Run(context.Background(), w, "node_modules/foo/Button.tsx", &mockFileAnalyzer{})
	if err == nil {
		t.Fatal("expected error from walker, got nil")
	}
}
