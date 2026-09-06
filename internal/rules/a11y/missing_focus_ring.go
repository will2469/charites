package a11y

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// MissingFocusRingRule mendeteksi elemen interaktif yang menekan cincin fokus bawaan browser
// (misal: 'outline-none') tanpa menyediakan indikator fokus visual pengganti ('focus-visible:ring-*').
type MissingFocusRingRule struct{}

// NewMissingFocusRingRule membuat instance baru MissingFocusRingRule.
func NewMissingFocusRingRule() *MissingFocusRingRule {
	return &MissingFocusRingRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat a11y.missing-focus-ring.
func (r *MissingFocusRingRule) ID() string {
	return "a11y.missing-focus-ring"
}

// Description mengembalikan deskripsi ringkas aturan.
func (r *MissingFocusRingRule) Description() string {
	return "Enforces visible focus indicator when suppressing default outline with outline-none (WCAG 2.4.7)"
}

// Category mengembalikan nama kategori rule.
func (r *MissingFocusRingRule) Category() string {
	return "a11y"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *MissingFocusRingRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *MissingFocusRingRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"WCAG 2.2 Success Criterion 2.4.7 (Focus Visible - Level AA)",
			"WCAG 2.2 Success Criterion 2.4.11 (Focus Appearance - Level AAA)",
			"W3C WAI-ARIA Authoring Practices (Keyboard Navigation & Focus Ring Management)",
		},
		CoreInvariant: "Interactive controls that suppress default browser outlines ('outline-none') must provide an alternative focus indicator (e.g. 'focus-visible:ring-2').",
		Grounding: "Web developers frequently add `outline-none` or `focus:outline-none` to eliminate the default browser blue halo when clicking buttons with a mouse.\n\n" +
			"However, stripping this outline without providing an alternative completely blinds keyboard and screen reader users:\n" +
			"1. Complete Loss of Focus Context: Users pressing Tab cannot see which button, link, or input is currently active.\n" +
			"2. Form & Transaction Abandonment: Users accidentally activate unexpected actions (like Delete instead of Next) because focus state is invisible.\n" +
			"3. Legal Accessibility Violations: Directly violates WCAG 2.4.7 Level AA criteria across ADA and European Accessibility Act (EAA) audits.\n\n" +
			"Best Practice: Use `outline-none focus-visible:ring-2 focus-visible:ring-primary` so mouse users see no distraction while keyboard users receive a crisp, visible focus indicator.",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Interactive button stripping outline with no replacement indicator",
				Code:     `<button className="outline-none bg-primary text-white px-4 py-2 rounded">Simpan</button>`,
			},
			{
				Language: "astro",
				Comment:  "Navigation link silencing focus outline",
				Code:     `<a href="/settings" class="focus:outline-none text-blue-600">Pengaturan</a>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Compliant replacement focus ring on keyboard navigation",
				Code:     `<button className="outline-none focus-visible:ring-2 focus-visible:ring-primary bg-primary text-white px-4 py-2 rounded">Simpan</button>`,
			},
			{
				Language: "astro",
				Comment:  "Link with replacement focus-visible outline indicator",
				Code:     `<a href="/settings" class="focus:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 text-blue-600">Pengaturan</a>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Keyboard Navigation Blindness",
				Severity: "HIGH",
				Impact:   "Users relying on keyboard or switch devices lose track of the active element.",
			},
			{
				Vector:   "WCAG 2.4.7 Compliance Audit Failure",
				Severity: "HIGH",
				Impact:   "Flagged as a major non-compliance violation in corporate accessibility audits.",
			},
		},
	}
}

// Evaluate memeriksa apakah penekanan outline disertai indikator fokus pengganti.
// Mematuhi 0 B/op, 0 allocs/op pada node bersih (QUAL-03).
func (r *MissingFocusRingRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || len(node.Classes) == 0 {
		return nil
	}

	if !IsInteractiveElement(node) {
		return nil
	}

	suppressed, suppressClass := hasOutlineSuppression(node.Classes)
	if !suppressed {
		return nil
	}

	if hasFocusReplacement(node.Classes) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Interactive <%s> suppresses default focus outline with '%s' without providing a visible replacement ring (WCAG 2.4.7)", node.Tag, suppressClass),
			Hint:     "Add 'focus-visible:ring-2 focus-visible:ring-offset-2' to guarantee keyboard focus visibility.",
		},
	}
}

func hasOutlineSuppression(classes []string) (bool, string) {
	for _, class := range classes {
		base := StripVariantsOnlyBase(class)
		if base == "outline-none" || base == "outline-0" || base == "[outline:none]" {
			return true, class
		}
	}
	return false, ""
}

func hasFocusReplacement(classes []string) bool {
	for _, class := range classes {
		// Periksa keberadaan modifier focus-visible atau focus yang memberi ring/outline/border/shadow
		if !strings.Contains(class, "focus-visible:") && !strings.Contains(class, "focus:") {
			continue
		}

		base := StripVariantsOnlyBase(class)

		// Ring pengganti: focus-visible:ring, focus:ring-2, dsb.
		if strings.HasPrefix(base, "ring") {
			return true
		}

		// Outline pengganti: focus-visible:outline-2, dsb (selain outline-none/0)
		if strings.HasPrefix(base, "outline-") && base != "outline-none" && base != "outline-0" {
			return true
		}

		// Border pengganti: focus-visible:border-2, focus:border-primary
		if strings.HasPrefix(base, "border-") || strings.HasPrefix(base, "border") {
			return true
		}

		// Shadow pengganti: focus-visible:shadow, focus:shadow-outline
		if strings.HasPrefix(base, "shadow") {
			return true
		}
	}

	return false
}
