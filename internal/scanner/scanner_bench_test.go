package scanner_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/will2469/charites/internal/config"
	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/scanner"
)

func BenchmarkWalker_DirectoryTraversal(b *testing.B) {
	tmpDir := b.TempDir()

	// Buat 500 file terstruktur
	for dirIdx := 0; dirIdx < 10; dirIdx++ {
		sub := filepath.Join(tmpDir, fmt.Sprintf("dir_%d", dirIdx))
		_ = os.MkdirAll(sub, 0o750)
		for fileIdx := 0; fileIdx < 50; fileIdx++ {
			f := filepath.Join(sub, fmt.Sprintf("Component_%d.tsx", fileIdx))
			_ = os.WriteFile(f, []byte("<div className=\"p-4\" />"), 0o600)
		}
	}

	matcher := config.NewIgnoreMatcher(nil)
	w := scanner.NewWalker(matcher, nil)

	b.ReportAllocs()

	for b.Loop() {
		jobs := make(chan string, 1000)
		_ = w.Walk(context.Background(), tmpDir, jobs)
		close(jobs)
		for range jobs {
		}
	}
}

type benchDummyAnalyzer struct{}

func (b *benchDummyAnalyzer) AnalyzeFile(path string) ([]ir.Diagnostic, error) {
	return []ir.Diagnostic{
		{
			File:     path,
			Line:     1,
			Column:   1,
			Rule:     "theme.hardcode-opacity-color",
			Severity: ir.SeverityWarn,
			Message:  "violation",
		},
	}, nil
}

func BenchmarkEndToEnd_ScanPipeline(b *testing.B) {
	tmpDir := b.TempDir()

	for i := 0; i < 50; i++ {
		f := filepath.Join(tmpDir, fmt.Sprintf("Button_%d.tsx", i))
		_ = os.WriteFile(f, []byte("<button className=\"bg-primary/10\" />"), 0o600)
	}

	matcher := config.NewIgnoreMatcher(nil)
	w := scanner.NewWalker(matcher, nil)
	pool := scanner.NewPool(4)
	analyzer := &benchDummyAnalyzer{}

	b.ReportAllocs()

	for b.Loop() {
		_, _ = pool.Run(context.Background(), w, tmpDir, analyzer)
	}
}
