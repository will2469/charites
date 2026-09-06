package cls

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// UnreservedFixedHeaderRule mendeteksi elemen header navigasi berposisi fixed/sticky top-0
// yang tidak memiliki kompensasi ruang tata letak (padding top pt-* atau margin mt-* pada konten utama, atau spacer block).
type UnreservedFixedHeaderRule struct{}

// NewUnreservedFixedHeaderRule membuat instance baru dari UnreservedFixedHeaderRule.
func NewUnreservedFixedHeaderRule() *UnreservedFixedHeaderRule {
	return &UnreservedFixedHeaderRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *UnreservedFixedHeaderRule) ID() string {
	return "cls.unreserved-fixed-header"
}

// Description mengembalikan ringkasan aturan.
func (r *UnreservedFixedHeaderRule) Description() string {
	return "Fixed or sticky header lacks layout space compensation (pt/mt) on subsequent in-flow content or spacer block"
}

// Category mengembalikan nama kategori rule.
func (r *UnreservedFixedHeaderRule) Category() string {
	return "cls"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *UnreservedFixedHeaderRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *UnreservedFixedHeaderRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C CSS Positioned Layout Module Level 3 (fixed & sticky positioning)",
			"Google Core Web Vitals (View-Overlap & Content Snapping Guidelines)",
			"Responsive Layout Architecture Invariants",
		},
		CoreInvariant: "Fixed or sticky header elements taking top position must provide corresponding layout space compensation on subsequent content (such as 'pt-*' or a spacer element).",
		Grounding: "When a top navigation header is declared with 'position: fixed' or dynamically mounted as sticky, it is removed from the normal document flow.\n\n" +
			"If the subsequent sibling content (such as the main container '<main>') does not reserve equivalent top padding ('pt-16') or include an explicit spacer element, the top portion of the main document gets covered underneath the header.\n\n" +
			"Furthermore, when headers mount asynchronously or change position dynamically during hydration, uncompensated content below suddenly shifts down or up, producing Cumulative Layout Shift (CLS).",
		Risks: []ir.RiskItem{
			{
				Vector:   "Obscured Top Page Content",
				Severity: "HIGH",
				Impact:   "Primary headings, hero banners, or breadcrumbs become invisible behind fixed header overlays.",
			},
			{
				Vector:   "Hydration Content Jump",
				Severity: "MEDIUM",
				Impact:   "Subsequent in-flow content snaps vertically when dynamic headers mount or change positioning.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Fixed header without padding compensation on following main element",
				Code: `<header className="fixed top-0 left-0 w-full h-16 bg-background z-50">
  <Navbar />
</header>
<main>
  <h1>Selamat Datang</h1>
</main>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Fixed header with matching top padding on main container",
				Code: `<header className="fixed top-0 left-0 w-full h-16 bg-background z-50">
  <Navbar />
</header>
<main className="pt-16">
  <h1>Selamat Datang</h1>
</main>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen navigasi fixed/sticky memiliki kompensasi ruang pada elemen berikutnya.
func (r *UnreservedFixedHeaderRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !isFixedHeader(node) {
		return nil
	}

	if hasHeaderCompensation(node) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Fixed or sticky header element <%s> lacks layout space compensation (such as 'pt-*' or 'mt-*' on the following <main> content or a dedicated spacer block). Uncompensated fixed headers obscure top content or induce vertical layout jumps.", node.Tag),
			Hint:     "Add top padding compensation to the main content container (e.g. <main class='pt-16'>) matching header height or insert an explicit spacer element.",
		},
	}
}
