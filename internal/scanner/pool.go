package scanner

import (
	"context"
	"runtime"
	"sync"

	"github.com/will2469/charites/internal/ir"
)

// FileAnalyzer mendefinisikan interface pemrosesan analisis berkas mandiri.
// Mengurangi kopling langsung antara paket scanner dan analyzer.
type FileAnalyzer interface {
	AnalyzeFile(path string) ([]ir.Diagnostic, error)
}

// Pool mengelola konkurensi pemrosesan berkas menggunakan goroutine worker pool.
type Pool struct {
	workers int
}

// NewPool menginisialisasi worker pool baru dengan batas kapasitas [1, 256].
// Jika workers <= 0, default mengalokasikan runtime.GOMAXPROCS(0).
func NewPool(workers int) *Pool {
	if workers <= 0 {
		workers = runtime.GOMAXPROCS(0)
	}
	if workers < 1 {
		workers = 1
	}
	if workers > 256 {
		workers = 256
	}

	return &Pool{
		workers: workers,
	}
}

// Workers mengembalikan jumlah goroutine aktif yang dikonfigurasi pada pool.
func (p *Pool) Workers() int {
	return p.workers
}

// Run mengeksekusi pipeline pemindaian paralel dengan invarian Single Producer = Single Closer:
// 1. Walker goroutine memproduksi dan menutup channel jobs.
// 2. N Worker goroutines mengonsumsi jobs dan memproduksi hasil ke channel results.
// 3. Coordinator goroutine tersinkronisasi sync.WaitGroup menutup channel results.
// 4. Aggregator mengumpulkan temuan dan mengurutkan secara total ordering (ir.SortDiagnostics).
// Pada saat pembatalan context (SIGINT/SIGTERM), temuan parsial dibuang dan mengembalikan ctx.Err().
func (p *Pool) Run(ctx context.Context, walker *Walker, target string, analyzer FileAnalyzer) ([]ir.Diagnostic, error) {
	bufferSize := p.workers * 2
	jobs := make(chan string, bufferSize)
	results := make(chan []ir.Diagnostic, bufferSize)

	var wg sync.WaitGroup

	// 1. Luncurkan N Worker Goroutines
	for i := 0; i < p.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case path, ok := <-jobs:
					if !ok {
						return
					}
					diags, err := analyzer.AnalyzeFile(path)
					if err == nil && len(diags) > 0 {
						select {
						case results <- diags:
						case <-ctx.Done():
							return
						}
					}
				}
			}
		}()
	}

	// 2. Luncurkan Walker Goroutine (Single Producer & Closer untuk channel jobs)
	walkErrChan := make(chan error, 1)
	go func() {
		defer close(jobs)
		walkErrChan <- walker.Walk(ctx, target, jobs)
	}()

	// 3. Luncurkan Coordinator Goroutine (Single Closer untuk channel results)
	go func() {
		wg.Wait()
		close(results)
	}()

	// 4. Aggregator: Kumpulkan seluruh diagnostic
	var allDiags []ir.Diagnostic
	for diags := range results {
		allDiags = append(allDiags, diags...)
	}

	walkErr := <-walkErrChan

	// Jika terjadi pembatalan context, buang hasil parsial demi integritas laporan
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if walkErr != nil {
		return nil, walkErr
	}

	return ir.SortDiagnostics(allDiags), nil
}
