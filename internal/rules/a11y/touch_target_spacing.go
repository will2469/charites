package a11y

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// TouchTargetSpacingRule mendeteksi kontainer dengan kontrol interaktif bersebelahan
// yang memiliki jarak spasial terlalu sempit (< 8px), memicu risiko salah sentuh (miss-tap).
type TouchTargetSpacingRule struct{}

// NewTouchTargetSpacingRule membuat instance baru TouchTargetSpacingRule.
func NewTouchTargetSpacingRule() *TouchTargetSpacingRule {
	return &TouchTargetSpacingRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat a11y.touch-target-spacing.
func (r *TouchTargetSpacingRule) ID() string {
	return "a11y.touch-target-spacing"
}

// Description mengembalikan deskripsi ringkas aturan.
func (r *TouchTargetSpacingRule) Description() string {
	return "Enforces at least 8px spacing between adjacent interactive elements to prevent miss-taps (WCAG 2.5.8)"
}

// Category mengembalikan nama kategori rule.
func (r *TouchTargetSpacingRule) Category() string {
	return "a11y"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *TouchTargetSpacingRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *TouchTargetSpacingRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"WCAG 2.2 Success Criterion 2.5.8 (Target Size & Spacing)",
			"Apple Human Interface Guidelines (Hit Target Separation)",
		},
		CoreInvariant: "Adjacent interactive elements within flex/grid containers must provide at least 8px physical spacing (e.g. gap-2) to avoid accidental miss-taps.",
		Grounding: "When interactive controls (such as Delete and Cancel buttons) are placed adjacent to each other with `gap-0`, `gap-1` (4px), or no spacing at all, users frequently suffer from miss-taps:\n\n" +
			"1. Destructive Activation: A user aiming for 'Cancel' inadvertently triggers 'Delete'.\n" +
			"2. Touch Jitter: Natural hand tremor or unstable mobile contexts (walking, transit) compound target activation errors.\n\n" +
			"Charites enforces a minimum 8px boundary (Tailwind `gap-2` or equivalent) across sibling interactive controls.",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Zero gap between action buttons in a flex row",
				Code:     `<div className="flex gap-0"><button className="bg-destructive text-white">Hapus</button><button>Batal</button></div>`,
			},
			{
				Language: "astro",
				Comment:  "Cramped 4px gap between icon buttons",
				Code:     `<div class="flex gap-1"><button aria-label="Edit"><EditIcon /></button><button aria-label="Delete"><TrashIcon /></button></div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Proper 12px separation between destructive and neutral actions",
				Code:     `<div className="flex gap-3"><button>Batal</button><button className="bg-destructive text-white">Hapus</button></div>`,
			},
			{
				Language: "astro",
				Comment:  "Standard 8px spacing between toolbar buttons",
				Code:     `<div class="flex gap-2"><button aria-label="Edit"><EditIcon /></button><button aria-label="Delete"><TrashIcon /></button></div>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Accidental Destructive Actions",
				Severity: "HIGH",
				Impact:   "Users trigger irreversible data loss due to cramped button spacing.",
			},
		},
	}
}

// Evaluate mengevaluasi jarak antar elemen interaktif di dalam kontainer flex/grid.
// Mematuhi 0 B/op, 0 allocs/op pada node bersih (QUAL-03).
func (r *TouchTargetSpacingRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || len(node.Children) < 2 {
		return nil
	}

	// Periksa apakah kontainer adalah flex atau grid
	isFlexOrGrid := false
	for _, c := range node.Classes {
		base := StripVariantsOnlyBase(c)
		if base == "flex" || base == "inline-flex" || base == "grid" || base == "inline-grid" {
			isFlexOrGrid = true
			break
		}
	}

	if !isFlexOrGrid {
		return nil
	}

	// Hitung jumlah anak interaktif langsung
	interactiveCount := countInteractiveChildren(node)
	if interactiveCount < 2 {
		return nil
	}

	// Evaluasi kelas gap
	hasGap := false
	gapPx := 0.0

	for _, c := range node.Classes {
		base := StripVariantsOnlyBase(c)
		if strings.HasPrefix(base, "gap-") || strings.HasPrefix(base, "gap-x-") {
			hasGap = true
			sub := strings.TrimPrefix(base, "gap-x-")
			sub = strings.TrimPrefix(sub, "gap-")
			if px, ok := ParseTailwindSizeToPx(sub); ok {
				gapPx = px
			}
		}
	}

	// Pelanggaran jika ada gap eksplisit < 8px
	if hasGap && gapPx < 8.0 {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Interactive controls inside container have insufficient touch spacing (gap %.0fpx < 8px)", gapPx),
				Hint:     "Provide at least 'gap-2' (8px) between adjacent interactive controls to prevent accidental miss-taps (WCAG 2.5.8).",
			},
		}
	}

	// Jika tidak ada gap sama sekali dan tidak ada margin kompensasi pada anak
	if !hasGap && !hasChildrenMarginCompensation(node) {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Multiple adjacent interactive controls inside flex/grid container without separation gap",
				Hint:     "Add at least 'gap-2' (8px) to the container to ensure touch safety between sibling buttons (WCAG 2.5.8).",
			},
		}
	}

	return nil
}

func countInteractiveChildren(node *ir.Node) int {
	count := 0
	for _, child := range node.Children {
		if child.Type == ir.NodeElement && IsInteractiveElement(child) {
			count++
		}
	}
	return count
}

func hasChildrenMarginCompensation(node *ir.Node) bool {
	for _, child := range node.Children {
		if child.Type != ir.NodeElement || !IsInteractiveElement(child) {
			continue
		}
		for _, c := range child.Classes {
			base := StripVariantsOnlyBase(c)
			if strings.HasPrefix(base, "m-") || strings.HasPrefix(base, "mx-") || strings.HasPrefix(base, "my-") ||
				strings.HasPrefix(base, "mr-") || strings.HasPrefix(base, "ml-") || strings.HasPrefix(base, "space-x-") {
				return true
			}
		}
	}
	return false
}
