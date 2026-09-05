package parser_test

import (
	"os"
	"runtime"
	"testing"

	"github.com/will2469/charites/internal/parser/astro"
	"github.com/will2469/charites/internal/parser/tsx"
)

// TestParser_MemoryBudgetProtocol memvalidasi protokol anggaran memori kanonikal:
// baseline = MemStats.HeapAlloc sebelum parse
// parse canonical fixture (~100 KB)
// runtime.GC()
// measure HeapAlloc setelah parse
// delta = post - baseline
// live heap delta SHOULD maintain <= 4x source size.
func TestParser_MemoryBudgetProtocol(t *testing.T) {
	src, err := os.ReadFile("../../tests/fixtures/canonical_sample.tsx")
	if err != nil {
		t.Fatalf("failed to read canonical fixture: %v", err)
	}

	// 1. Ambil baseline alokasi heap sebelum parse setelah membersihkan garbage
	runtime.GC()
	var mPre runtime.MemStats
	runtime.ReadMemStats(&mPre)

	// 2. Parse canonical fixture
	root, err := tsx.Extract(src)
	if err != nil {
		t.Fatalf("unexpected extract error: %v", err)
	}
	if root == nil {
		t.Fatalf("expected non-nil root")
	}

	// Jaga agar objek root tetap hidup saat pengukuran
	runtime.KeepAlive(root)

	// 3. Garbage collection untuk membersihkan alokasi sementara
	runtime.GC()
	var mPost runtime.MemStats
	runtime.ReadMemStats(&mPost)

	// 4. Hitung delta heap hidup (live heap delta)
	var delta uint64
	if mPost.HeapAlloc > mPre.HeapAlloc {
		delta = mPost.HeapAlloc - mPre.HeapAlloc
	}

	sourceSize := uint64(len(src))
	maxAllowed := 4 * sourceSize
	ratio := float64(delta) / float64(sourceSize)

	t.Logf("Canonical fixture size: %d bytes (%.2f KB)", sourceSize, float64(sourceSize)/1024)
	t.Logf("Post-GC live heap delta: %d bytes (%.2f KB)", delta, float64(delta)/1024)
	t.Logf("Heap-to-source ratio: %.2fx (budget ceiling: 4.00x)", ratio)

	if delta > maxAllowed {
		t.Errorf("canonical parser workload violated memory budget: delta %d bytes exceeds 4x source size %d bytes (ratio: %.2fx)", delta, maxAllowed, ratio)
	}
}

func BenchmarkParser_MemoryBudget(b *testing.B) {
	src, err := os.ReadFile("../../tests/fixtures/canonical_sample.tsx")
	if err != nil {
		b.Fatalf("failed to read canonical fixture: %v", err)
	}

	b.ReportAllocs()
	b.SetBytes(int64(len(src)))
	b.ResetTimer()

	for b.Loop() {
		node, err := tsx.Extract(src)
		if err != nil {
			b.Fatalf("parse error: %v", err)
		}
		if node == nil {
			b.Fatal("nil node returned")
		}
	}
}

func BenchmarkParser_Astro(b *testing.B) {
	src, err := os.ReadFile("../../tests/fixtures/sample.astro")
	if err != nil {
		b.Fatalf("failed to read astro fixture: %v", err)
	}

	b.ReportAllocs()
	b.SetBytes(int64(len(src)))
	b.ResetTimer()

	for b.Loop() {
		node, err := astro.Parse(src)
		if err != nil {
			b.Fatalf("parse error: %v", err)
		}
		if node == nil {
			b.Fatal("nil node returned")
		}
	}
}
