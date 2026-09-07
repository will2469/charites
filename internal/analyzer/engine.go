package analyzer

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/will2469/charites/internal/config"
	"github.com/will2469/charites/internal/drift"
	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/parser/astro"
	"github.com/will2469/charites/internal/parser/tsx"
)

// Engine bertanggung jawab melakukan traversal pohon AST menggunakan iterator Go 1.26
// dan mengoordinasikan evaluasi active rules serta penekanan direktif inline ignore.
type Engine struct {
	activeRules []config.ActiveRule
}

// NewEngine membuat instans Engine baru dengan daftar ActiveRule yang sudah teresolusi.
func NewEngine(activeRules []config.ActiveRule) *Engine {
	return &Engine{
		activeRules: activeRules,
	}
}

// ActiveRules mengembalikan salinan daftar active rules yang dikonfigurasi pada engine.
func (e *Engine) ActiveRules() []config.ActiveRule {
	out := make([]config.ActiveRule, len(e.activeRules))
	copy(out, e.activeRules)
	return out
}

// AnalyzeFile membaca berkas dari disk, mengekstrak direktif inline ignore,
// mem-parse ke pohon IR terpadu (*ir.Node), dan mengevaluasi active rules.
func (e *Engine) AnalyzeFile(path string) ([]ir.Diagnostic, error) {
	diags, _, err := e.AnalyzeFileWithOccurrences(path)
	return diags, err
}

// AnalyzeFileWithOccurrences membaca berkas, mem-parse IR satu kali, mengevaluasi active rules,
// dan mengekstrak StyleOccurrences untuk repository-level drift analysis.
func (e *Engine) AnalyzeFileWithOccurrences(path string) ([]ir.Diagnostic, []drift.StyleOccurrence, error) {
	src, err := os.ReadFile(filepath.Clean(path)) //nolint:gosec // controlled scan target path
	if err != nil {
		return nil, nil, err
	}

	inlineIgnores := ParseDirectives(src)
	ext := strings.ToLower(filepath.Ext(path))

	var root *ir.Node
	var parseErr error

	switch ext {
	case ".astro":
		root, parseErr = astro.Parse(src)
	case ".tsx", ".jsx":
		root, parseErr = tsx.Extract(src)
	default:
		// Format tidak didukung: kembalikan tanpa error
		return nil, nil, nil
	}

	if parseErr != nil {
		return nil, nil, parseErr
	}

	diags := e.AnalyzeTree(path, root, inlineIgnores)
	occs := drift.ExtractTreeOccurrences(path, root)
	return diags, occs, nil
}

// AnalyzeTree mengevaluasi active rules pada pohon IR in-memory.
// Digunakan langsung oleh AnalyzeFile dan suite pengujian/benchmark berkecepatan tinggi.
func (e *Engine) AnalyzeTree(path string, root *ir.Node, inlineIgnores map[int][]string) []ir.Diagnostic {
	if root == nil || len(e.activeRules) == 0 {
		return nil
	}

	ctx := NewContext(path, inlineIgnores)

	for node := range root.Walk() {
		for _, active := range e.activeRules {
			diags := active.Rule.Evaluate(node)
			for _, d := range diags {
				// Normalisasi path berkas dan timpa severity dengan EffectiveSeverity
				d.File = path
				d.Severity = active.EffectiveSeverity

				// Evaluasi penekanan inline ignore (signature baku: IsIgnored(d, node))
				if !ctx.IsIgnored(d, node) {
					ctx.AddDiagnostic(d)
				}
			}
		}
	}

	return ctx.DiagnosticsList()
}
