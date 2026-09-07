package scanner

import (
	"context"
	"runtime"
	"sync"

	"github.com/will2469/charites/internal/drift"
	"github.com/will2469/charites/internal/ir"
)

// FileAnalyzer mendefinisikan interface pemrosesan analisis berkas mandiri.
// Mengurangi kopling langsung antara paket scanner dan analyzer.
type FileAnalyzer interface {
	AnalyzeFile(path string) ([]ir.Diagnostic, error)
}

// FileAnalysisResult merepresentasikan temuan diagnostik dan kemunculan gaya dari satu berkas.
type FileAnalysisResult struct {
	Diagnostics []ir.Diagnostic
	Occurrences []drift.StyleOccurrence
}

// OccurrencesAnalyzer mendefinisikan interface pemrosesan analisis berkas yang mengekstrak StyleOccurrence.
type OccurrencesAnalyzer interface {
	FileAnalyzer
	AnalyzeFileWithOccurrences(path string) ([]ir.Diagnostic, []drift.StyleOccurrence, error)
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

// Run mengeksekusi pipeline pemindaian paralel standar.
func (p *Pool) Run(ctx context.Context, walker *Walker, target string, analyzer FileAnalyzer) ([]ir.Diagnostic, error) {
	diags, _, err := p.RunWithOccurrences(ctx, walker, target, analyzer)
	return diags, err
}

// RunWithOccurrences mengeksekusi pipeline pemindaian paralel dengan transfer kepemilikan data (race-free by ownership):
// 1. Walker goroutine memproduksi path berkas ke channel jobs.
// 2. N Worker goroutines mengonsumsi jobs dan mengirim FileAnalysisResult ke channel results.
// 3. Coordinator goroutine tersinkronisasi sync.WaitGroup menutup channel results.
// 4. Aggregator tunggal mengumpulkan diagnostik dan style occurrences tanpa shared mutex.
func (p *Pool) RunWithOccurrences(ctx context.Context, walker *Walker, target string, analyzer FileAnalyzer) ([]ir.Diagnostic, []drift.StyleOccurrence, error) {
	bufferSize := p.workers * 2
	jobs := make(chan string, bufferSize)
	results := make(chan FileAnalysisResult, bufferSize)

	var wg sync.WaitGroup

	// 1. Luncurkan N Worker Goroutines
	for i := 0; i < p.workers; i++ {
		wg.Add(1)
		go p.startWorker(ctx, &wg, jobs, results, analyzer)
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

	// 4. Aggregator: Kumpulkan seluruh diagnostic dan occurrences
	var allDiags []ir.Diagnostic
	var allOccurrences []drift.StyleOccurrence
	for res := range results {
		if len(res.Diagnostics) > 0 {
			allDiags = append(allDiags, res.Diagnostics...)
		}
		if len(res.Occurrences) > 0 {
			allOccurrences = append(allOccurrences, res.Occurrences...)
		}
	}

	walkErr := <-walkErrChan

	// Jika terjadi pembatalan context, buang hasil parsial demi integritas laporan
	if ctx.Err() != nil {
		return nil, nil, ctx.Err()
	}

	if walkErr != nil {
		return nil, nil, walkErr
	}

	return ir.SortDiagnostics(allDiags), allOccurrences, nil
}

func (p *Pool) startWorker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan string, results chan<- FileAnalysisResult, analyzer FileAnalyzer) {
	defer wg.Done()
	occAnalyzer, hasOcc := analyzer.(OccurrencesAnalyzer)

	for {
		select {
		case <-ctx.Done():
			return
		case path, ok := <-jobs:
			if !ok {
				return
			}

			var diags []ir.Diagnostic
			var occs []drift.StyleOccurrence
			var err error

			if hasOcc {
				diags, occs, err = occAnalyzer.AnalyzeFileWithOccurrences(path)
			} else {
				diags, err = analyzer.AnalyzeFile(path)
			}

			if err == nil && (len(diags) > 0 || len(occs) > 0) {
				select {
				case results <- FileAnalysisResult{Diagnostics: diags, Occurrences: occs}:
				case <-ctx.Done():
					return
				}
			}
		}
	}
}
