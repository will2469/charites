package cli

import (
	"github.com/will2469/charites/internal/reporter"
)

const (
	// ExitClean menunjukkan pemindaian selesai tanpa pelanggaran bertingkat error (exit 0).
	ExitClean = 0
	// ExitViolations menunjukkan ditemukan pelanggaran bertingkat error atau warning dengan --fail-on-warn (exit 1).
	ExitViolations = 1
	// ExitOperational menunjukkan kesalahan operasional CLI atau masukan tidak sah (exit 2).
	ExitOperational = 2
	// ExitInterrupted menunjukkan proses dihentikan oleh sinyal pengguna SIGINT/SIGTERM (exit 130).
	ExitInterrupted = 130
)

// ResolveExitCode menentukan exit code POSIX deterministik berdasarkan hasil pemindaian.
// Sesuai invarian SPEC-05 & QUAL-05: temuan pelanggaran diagnostik dilarang keras menghasilkan exit 2.
func ResolveExitCode(summary *reporter.ScanSummary, failOnWarn bool) int {
	if summary == nil {
		return ExitClean
	}
	if summary.ErrorCount > 0 {
		return ExitViolations
	}
	if failOnWarn && summary.WarningCount > 0 {
		return ExitViolations
	}
	return ExitClean
}
