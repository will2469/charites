package rules

import (
	"github.com/will2469/charites/internal/ir"
)

// Rule mendefinisikan interface baku evaluasi aturan statis Charites.
// Seluruh implementasi aturan wajib mematuhi kontrak evaluasi murni (pure function tanpa side-effect I/O).
type Rule interface {
	// ID mengembalikan Charites Rule ID tunggal berformat <category>.<rule-slug> (misal: "theme.hardcode-opacity-color").
	ID() string
	// Description mengembalikan penjelasan ringkas maksud dan tujuan aturan.
	Description() string
	// Category mengembalikan nama kategori aturan (theme, a11y, perf, layout, seo).
	Category() string
	// DefaultSeverity mengembalikan tingkat keparahan default jika tidak di-override oleh konfigurasi pengguna.
	DefaultSeverity() ir.Severity
	// Evaluate mengevaluasi sebuah node IR dan mengembalikan daftar pelanggaran diagnosis yang ditemukan.
	// Fungsi ini wajib bebas dari operasi I/O dan tidak boleh memodifikasi node IR.
	Evaluate(node *ir.Node) []ir.Diagnostic
}
