package analyzer

import (
	"github.com/will2469/charites/internal/ir"
)

// SortDiagnostics mengurutkan temuan diagnostik menggunakan 7-level Canonical Total Ordering:
// File -> Line -> Column -> Rule -> Severity -> Message -> Hint
// dan melakukan deduplikasi identik.
// Memanfaatkan secara langsung implementasi SSOT dari paket internal/ir.
func SortDiagnostics(diags []ir.Diagnostic) []ir.Diagnostic {
	return ir.SortDiagnostics(diags)
}
