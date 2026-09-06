package performance

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// TailwindUntrackedPackageSourceRule mengaudit berkas CSS root Tailwind v4 yang tidak memuat direktif @source.
type TailwindUntrackedPackageSourceRule struct{}

// NewTailwindUntrackedPackageSourceRule membuat instance baru dari TailwindUntrackedPackageSourceRule.
func NewTailwindUntrackedPackageSourceRule() *TailwindUntrackedPackageSourceRule {
	return &TailwindUntrackedPackageSourceRule{}
}

// ID mengembalikan identifier unik kanonikal aturan.
func (r *TailwindUntrackedPackageSourceRule) ID() string {
	return "performance.tailwind-untracked-package-source"
}

// Category mengembalikan kategori aturan ('performance').
func (r *TailwindUntrackedPackageSourceRule) Category() string {
	return "performance"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *TailwindUntrackedPackageSourceRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Description mengembalikan deskripsi tujuan dan konteks aturan.
func (r *TailwindUntrackedPackageSourceRule) Description() string {
	return "Mewajibkan pendaftaran direktif @source pada berkas CSS root Tailwind v4 ketika mengimpor paket workspace monorepo eksternal."
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *TailwindUntrackedPackageSourceRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Tailwind CSS v4 Configuration Architecture (@source Directive Specification)",
			"Monorepo Multi-Package Style Discovery Standards",
			"Oxide Engine Workspace Scanning Invariants",
		},
		CoreInvariant: "Tailwind CSS v4 root stylesheets importing external monorepo packages must declare '@source' path directives; without '@source', the Oxide scanner skips external package directories, silently dropping all utility styles from compiled builds.",
		Grounding: "In Tailwind CSS v4, the legacy `tailwind.config.js` `content` array is replaced by CSS-first `@source` directives in the main stylesheet.\n\n" +
			"By default, Tailwind v4 only scans files in the immediate project directory. If the project imports components from external workspace packages (e.g. `@repo/ui` or `../../packages/...`), those package directories are ignored by default.\n\n" +
			"Failing to add `@source \"../../packages/ui\";` causes all Tailwind utility classes used inside those shared packages to be completely absent from the final CSS bundle.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Missing Monorepo Component Styles",
				Severity: "HIGH",
				Impact:   "Shared monorepo UI components render completely unstyled in production because utility classes inside them were never scanned.",
			},
			{
				Vector:   "Silent Build Failures",
				Severity: "HIGH",
				Impact:   "No build errors are thrown; stylesheets simply compile without the required utility declarations.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "css",
				Comment:  "Berkas CSS root tidak menyertakan direktif @source untuk paket monorepo",
				Code: `/* Pelanggaran: Mengimpor tailwindcss tanpa @source untuk monorepo packages */
@import "tailwindcss";`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "css",
				Comment:  "Mendaftarkan path paket eksternal via @source",
				Code: `/* Patuh: Menyertakan direktif @source untuk paket monorepo */
@import "tailwindcss";
@source "../../packages/ui";`,
			},
		},
	}
}

// Evaluate memeriksa apakah berkas stylesheet root Tailwind v4 kehilangan direktif @source.
func (r *TailwindUntrackedPackageSourceRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil {
		return nil
	}

	var fileSrc string
	line := 1
	switch {
	case strings.EqualFold(node.Tag, "style"):
		fileSrc = getStyleNodeText(node)
		if node.Span.Line > 0 {
			line = node.Span.Line
		}
	case isSourceRootOrScript(node):
		fileSrc = getFileSourceContent(node)
	default:
		return nil
	}

	if len(fileSrc) == 0 {
		return nil
	}

	missingSource, isMissing := findUntrackedPackageSource(fileSrc)
	if !isMissing {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     line,
			Column:   1,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Tailwind CSS v4 root stylesheet imports 'tailwindcss' but is missing '@source' directives for external %s, preventing Oxide compiler from scanning package utility classes.", missingSource),
			Hint:     "Add '@source \"<relative-path-to-package>\";' to ensure the Tailwind v4 compiler scans external packages for utility extraction.",
		},
	}
}
