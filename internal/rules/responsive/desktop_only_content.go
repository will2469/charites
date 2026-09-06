package responsive

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// DesktopOnlyContentRule mendeteksi aksi penting, formulir submit, atau primary CTA
// yang disembunyikan pada baseline mobile (misal: 'hidden md:flex') tanpa menyediakan alternatif mobile.
type DesktopOnlyContentRule struct{}

// NewDesktopOnlyContentRule membuat instance baru dari DesktopOnlyContentRule.
func NewDesktopOnlyContentRule() *DesktopOnlyContentRule {
	return &DesktopOnlyContentRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *DesktopOnlyContentRule) ID() string {
	return "responsive.desktop-only-content"
}

// Description mengembalikan ringkasan aturan.
func (r *DesktopOnlyContentRule) Description() string {
	return "Warns when primary action buttons or form submit controls are hidden on mobile viewports without mobile alternatives"
}

// Category mengembalikan nama kategori rule.
func (r *DesktopOnlyContentRule) Category() string {
	return "responsive"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *DesktopOnlyContentRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *DesktopOnlyContentRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Mobile-First Responsive Web Design (Luke Wroblewski)",
			"WCAG 2.2 Guideline 3.2 (Predictable - Consistent Navigation)",
			"Google Mobile-Friendly & Core Web Vitals Guidance",
		},
		CoreInvariant: "Primary call-to-action (CTA) controls, checkout buttons, and form submissions must not be hidden on the mobile baseline ('hidden md:...') without accessible mobile parity.",
		Grounding: "In mobile-first design, hiding ancillary content (such as secondary marketing badges or decorative sidebars) on narrow viewports is standard practice.\n\n" +
			"However, hiding vital action triggers (e.g. 'Checkout', 'Bayar Sekarang', 'Kirim Berkas', or form submit buttons) via 'hidden md:flex' or 'hidden lg:block' leaves smartphone users stranded with incomplete user flows and broken core functionality.\n\n" +
			"Charites enforces that essential primary actions remain discoverable across all breakpoints, whether inline, within a bottom action sheet, or through a responsive floating bar.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Mobile User Conversion Disruption",
				Severity: "HIGH",
				Impact:   "Smartphone users cannot complete essential transactions, checkouts, or form submissions when action buttons are hidden on mobile.",
			},
			{
				Vector:   "Inconsistent Navigation Experience",
				Severity: "MEDIUM",
				Impact:   "Users switching between devices experience confusing functional disparity between desktop and mobile layouts.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Primary checkout button hidden on mobile baseline",
				Code: `<button type="submit" className="hidden md:flex items-center px-4 py-2 bg-primary text-primary-foreground">
  Bayar Sekarang
</button>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Primary checkout button visible across all viewports with responsive sizing",
				Code: `<button type="submit" className="flex items-center justify-center w-full md:w-auto px-4 py-2 bg-primary text-primary-foreground">
  Bayar Sekarang
</button>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen aksi primer disembunyikan di layar mobile.
func (r *DesktopOnlyContentRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if isPrimaryAction(node) && isDesktopHidden(node.Classes) {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  fmt.Sprintf("Primary action control <%s> is hidden on mobile baseline ('hidden md:...'). Mobile users will be unable to complete this essential action.", node.Tag),
				Hint:     "Ensure primary actions are accessible on mobile viewports (e.g. make visible on mobile or provide a sticky bottom mobile action bar).",
			},
		}
	}

	if isDesktopHidden(node.Classes) {
		for _, child := range node.Children {
			if child != nil && child.Type == ir.NodeElement && isPrimaryAction(child) {
				return []ir.Diagnostic{
					{
						Line:     child.Span.Line,
						Column:   child.Span.Column,
						Rule:     r.ID(),
						Severity: r.DefaultSeverity(),
						Message:  fmt.Sprintf("Primary action control <%s> is enclosed within a desktop-only container ('hidden md:...'). Mobile users will be locked out of this action.", child.Tag),
						Hint:     "Move the primary action outside the desktop-only container or provide a dedicated mobile action trigger.",
					},
				}
			}
		}
	}

	return nil
}
