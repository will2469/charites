package tests_test

import (
	"bytes"
	"io"
	"path/filepath"
	"testing"

	"github.com/will2469/charites/internal/cli"
)

// BenchmarkFullPipeline_Monorepo mengukur latensi pemindaian menyeluruh dari ujung ke ujung
// sesuai dengan spesifikasi BENCH-06-E2E-001 pada TEST-06-GOLDEN.
func BenchmarkFullPipeline_Monorepo(b *testing.B) {
	fixtureDir := filepath.Join("fixtures", "projects")

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		var buf bytes.Buffer
		code := cli.ExecuteArgs([]string{"scan", fixtureDir, "-f", "json"}, &buf, io.Discard)
		if code != 0 && code != 1 {
			b.Fatalf("unexpected exit code: %d", code)
		}
	}
}
