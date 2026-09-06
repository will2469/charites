package responsive

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// GridMinColumnRule mendeteksi penggunaan 'minmax()' pada CSS grid dengan batas minimum yang kaku (> 320px)
// tanpa penyesuaian breakpoint atau 'min(100%, ...)', yang memicu horizontal overflow pada smartphone.
type GridMinColumnRule struct{}

// NewGridMinColumnRule membuat instance baru dari GridMinColumnRule.
func NewGridMinColumnRule() *GridMinColumnRule {
	return &GridMinColumnRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *GridMinColumnRule) ID() string {
	return "responsive.grid-min-column"
}

// Description mengembalikan ringkasan aturan.
func (r *GridMinColumnRule) Description() string {
	return "Warns against CSS grid minmax column definitions with rigid minimum sizes (> 320px) that cause horizontal overflow on mobile viewports"
}

// Category mengembalikan nama kategori rule.
func (r *GridMinColumnRule) Category() string {
	return "responsive"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *GridMinColumnRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *GridMinColumnRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C CSS Grid Layout Module Level 1 (The minmax() Function)",
			"WCAG 2.2 SC 1.4.10 (Reflow - Level AA)",
			"Mobile Web Best Practices: Preventing Horizontal Viewport Blowout",
		},
		CoreInvariant: "CSS grid column minmax tracks on mobile baseline must not enforce rigid minimum widths greater than 320px without dynamic clamping ('min(100%, <size>)') or desktop breakpoint scoping ('md:grid-cols-...').",
		Grounding: "A common CSS grid pattern for auto-fit cards is 'repeat(auto-fit, minmax(350px, 1fr))' or 'repeat(auto-fill, minmax(400px, 1fr))'. While this looks great on desktop and tablet monitors, the minimum column track width of 350px or 400px exceeds the 360px physical width of most smartphones and the 320px minimum WCAG reflow baseline.\n\n" +
			"Because CSS grid does not shrink tracks below their minimum minmax threshold, the grid blows out horizontally, introducing an unintended horizontal scrollbar across the entire mobile page.\n\n" +
			"Charites detects rigid minmax tracks and suggests using the standard modern CSS clamp idiom: 'minmax(min(100%, 20rem), 1fr)' or scoping multi-column grids behind 'md:' breakpoints.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Mobile Viewport Horizontal Scrollbar Blowout",
				Severity: "HIGH",
				Impact:   "The entire website scrolls sideways on mobile phones because a single card grid enforces 350px+ minimum track size.",
			},
			{
				Vector:   "Broken Touch Gestures & Visual Glitches",
				Severity: "MEDIUM",
				Impact:   "Accidental horizontal swiping triggers page drift instead of vertical scroll.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Grid specifying 400px minimum column width on mobile baseline",
				Code: `<div className="grid grid-cols-[repeat(auto-fit,minmax(400px,1fr))] gap-4">
  <div className="card">Kartu Layanan 1</div>
  <div className="card">Kartu Layanan 2</div>
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Clamped minmax ensuring column never exceeds 100% on narrow screens",
				Code: `<div className="grid grid-cols-[repeat(auto-fit,minmax(min(100%,20rem),1fr))] gap-4">
  <div className="card">Kartu Layanan 1</div>
  <div className="card">Kartu Layanan 2</div>
</div>`,
			},
			{
				Language: "tsx",
				Comment:  "Mobile single-column with desktop-scoped multi-column minmax",
				Code: `<div className="grid grid-cols-1 md:grid-cols-[repeat(auto-fit,minmax(350px,1fr))] gap-4">
  <div className="card">Kartu Layanan 1</div>
  <div className="card">Kartu Layanan 2</div>
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah kelas grid-cols menggunakan minmax dengan lebar minimum yang melebihi batas mobile.
func (r *GridMinColumnRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	cls, found := hasExcessiveGridMinColumn(node.Classes)
	if !found {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Grid column class '%s' specifies a rigid minimum column size exceeding 320px. On smartphones (360px), this triggers severe horizontal viewport blowout.", cls),
			Hint:     "Wrap the minimum threshold with 'min(100%, <size>)' (e.g. 'minmax(min(100%, 20rem), 1fr)') or scope the minmax grid behind responsive desktop prefixes ('md:grid-cols-...').",
		},
	}
}
