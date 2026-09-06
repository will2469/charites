package responsive

import (
	"github.com/will2469/charites/internal/ir"
)

// SafeAreaMissingRule mendeteksi elemen fixed atau sticky yang menempel di dasar viewport (bottom-0)
// tanpa bantalan safe area (pb-[env(safe-area-inset-bottom)] atau pb-safe).
type SafeAreaMissingRule struct{}

// NewSafeAreaMissingRule membuat instance baru dari SafeAreaMissingRule.
func NewSafeAreaMissingRule() *SafeAreaMissingRule {
	return &SafeAreaMissingRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *SafeAreaMissingRule) ID() string {
	return "responsive.safe-area-missing"
}

// Description mengembalikan ringkasan aturan.
func (r *SafeAreaMissingRule) Description() string {
	return "Warns when bottom-docked fixed or sticky elements lack safe-area-inset-bottom padding for modern mobile home indicators"
}

// Category mengembalikan nama kategori rule.
func (r *SafeAreaMissingRule) Category() string {
	return "responsive"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *SafeAreaMissingRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *SafeAreaMissingRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C CSS Mobile Safe Area Insets (env(safe-area-inset-bottom))",
			"Apple Human Interface Guidelines (Display Cutouts & Home Indicator)",
			"Android Full-Screen Gesture Navigation Guidelines",
		},
		CoreInvariant: "Elements docked to the bottom of the viewport (fixed bottom-0 or sticky bottom-0) must include safe-area bottom padding (pb-[env(safe-area-inset-bottom)] or pb-safe).",
		Grounding: "Modern smartphones without physical home buttons utilize system-level gesture bars (the iPhone Home Indicator and Android gesture pill) at the bottom edge of the screen.\n\n" +
			"Positioning bottom navigation bars or action buttons flush against the bottom edge (bottom-0) without safe-area padding causes controls to collide directly with the operating system navigation bar.\n\n" +
			"Providing safe-area bottom padding ensures interactive controls are elevated cleanly above system navigation indicators.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Home Indicator Collision & Mis-Taps",
				Severity: "MEDIUM",
				Impact:   "Users attempting to tap bottom navigation items accidentally trigger the OS home swipe gesture instead.",
			},
			{
				Vector:   "Visual Element Occlusion",
				Severity: "LOW",
				Impact:   "Bottom buttons and labels appear partially obscured behind the white/black system gesture bar.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Bottom fixed navigation bar lacking safe-area padding",
				Code: `<nav className="fixed bottom-0 left-0 right-0 h-16 bg-surface flex items-center justify-around">
  <a href="/home">Beranda</a>
  <a href="/layanan">Layanan</a>
</nav>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Bottom fixed navigation with safe-area padding lifting content above home indicator",
				Code: `<nav className="fixed bottom-0 left-0 right-0 pb-[env(safe-area-inset-bottom)] bg-surface flex items-center justify-around">
  <a href="/home">Beranda</a>
  <a href="/layanan">Layanan</a>
</nav>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen fixed/sticky di dasar layar memiliki safe-area padding.
func (r *SafeAreaMissingRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || len(node.Classes) == 0 {
		return nil
	}

	if !isBottomDocked(node.Classes) {
		return nil
	}

	if hasSafeAreaBottomPadding(node.Classes) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Bottom-docked container (fixed/sticky bottom-0) lacks safe-area inset bottom padding (e.g. pb-[env(safe-area-inset-bottom)] or pb-safe). On modern mobile devices (iPhone Home Bar, Android gesture navigation), interactive elements will be obscured by the system home indicator.",
			Hint:     "Add safe-area bottom padding, e.g. 'pb-[env(safe-area-inset-bottom)]' or 'pb-safe' to lift navigation above the home indicator.",
		},
	}
}
