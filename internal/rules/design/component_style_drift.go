package design

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/drift"
	"github.com/will2469/charites/internal/ir"
)

// ComponentStyleDriftRule mendeteksi bahaya kontras lokal WCAG 1.4.3 pada node komponen
// dan menyediakan antarmuka evaluasi untuk rule design.component-style-drift.
type ComponentStyleDriftRule struct{}

// NewComponentStyleDriftRule membuat instance baru ComponentStyleDriftRule.
func NewComponentStyleDriftRule() *ComponentStyleDriftRule {
	return &ComponentStyleDriftRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *ComponentStyleDriftRule) ID() string {
	return "design.component-style-drift"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *ComponentStyleDriftRule) Description() string {
	return "Clusters and flags design system style sprawl, rogue geometric outliers, and chromatic color drift per UI component"
}

// Category mengembalikan nama kategori rule.
func (r *ComponentStyleDriftRule) Category() string {
	return "design"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *ComponentStyleDriftRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki otomatis.
func (r *ComponentStyleDriftRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Design Tokens Community Group (DTCG) 3-Tier Architecture",
			"W3C Web Content Accessibility Guidelines (WCAG) 2.2 Criterion 1.4.3 (Contrast Minimum)",
			"Concentric Border Radius Geometry (R_inner = R_outer - Padding)",
		},
		CoreInvariant: "UI components must maintain uniform geometric tokens (radii, elevation, press feedback) and accessible chromatic contrast across the repository without rogue style drift.",
		Grounding: "In large monorepos and evolving design systems, developers frequently introduce Style Drift (Design Entropy).\n\n" +
			"While static rules catch raw hex codes on a single-file basis, they are blind to cross-file design entropy where the same UI primitive splinters into subtle variations:\n" +
			"- <Button> using rounded-md (98% dominant), but isolated features use rounded-2xl or rounded-none.\n" +
			"- <Button> using active:scale-95, but rogue pages use active:scale-[0.98].\n" +
			"- Lack of first-class variants (e.g. missing 'warning' variant on <Button>) forcing developers to improvise with competing raw Tailwind hues (yellow-400, amber-500, orange-500), creating severe WCAG 1.4.3 contrast failures (bg-yellow-400 with white text = 1.53:1).\n\n" +
			"Concentric border radius geometry dictates that outer surfaces (<Card> using rounded-xl/2xl) must legitimately differ from inner elements (<Button> using rounded-md). Therefore, Charites enforces Component-Scoped statistical clustering.\n\n" +
			"Repository analysis groups styles by (ComponentKind, PropertyCategory), applies a minimum sample gate (10 occurrences), determines canonical standards (>= 80%), and flags rogue outliers (<= 10%).",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Severe WCAG contrast hazard with primitive background and white text",
				Code:     `<Button className="bg-yellow-400 text-white">Warning Action</Button>`,
			},
			{
				Language: "tsx",
				Comment:  "Rogue outlier corner radius on button diverging from design system canonical rounded-md",
				Code:     `<Button className="rounded-2xl">Checkout</Button>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Canonical button styling with accessible contrast and semantic tokens",
				Code:     `<Button className="rounded-md active:scale-95 bg-primary text-primary-foreground">Submit</Button>`,
			},
			{
				Language: "tsx",
				Comment:  "Card surface container legitimately using larger outer radius",
				Code:     `<Card className="rounded-xl border border-border/40 p-6">Card Content</Card>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Visual & Brand Incoherence",
				Severity: "MEDIUM",
				Impact:   "Fragmented corner radii and micro-interactions degrade perceived app polish and user confidence.",
			},
			{
				Vector:   "WCAG Accessibility Violations",
				Severity: "HIGH",
				Impact:   "Ad-hoc yellow and amber button backgrounds paired with white text fail WCAG AA contrast (minimum 4.5:1), rendering labels unreadable for low-vision users.",
			},
		},
	}
}

// Evaluate mengevaluasi sebuah node IR untuk memeriksa invariant lokal mutlak (misal: WCAG contrast hazard).
// Evaluasi empiris statistik repositori (canonical vs rogue) dilakukan di tingkat repositori oleh package internal/drift.
func (r *ComponentStyleDriftRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || len(node.Classes) == 0 {
		return nil
	}

	scope, ok := drift.ResolveScope(node)
	if !ok {
		return nil
	}

	// Evaluasi hanya pada komponen interaktif atau indikator status
	switch scope.Kind {
	case drift.KindButton, drift.KindInput, drift.KindBadge, drift.KindAlert, drift.KindCard:
	default:
		return nil
	}

	var bgClass, textClass string
	for _, c := range node.Classes {
		if strings.HasPrefix(c, "bg-") && !strings.Contains(c, ":") {
			bgClass = c
		}
		if strings.HasPrefix(c, "text-") && !strings.Contains(c, ":") {
			textClass = c
		}
	}

	if bgClass == "" || textClass == "" {
		return nil
	}

	ratio, failsAA, ok := drift.CheckContrastHazard(bgClass, textClass)
	if ok && failsAA {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message: fmt.Sprintf(
					"WCAG 1.4.3 Contrast Hazard: '%s' with '%s' yields contrast ratio %.2f:1 (fails minimum 4.5:1)",
					bgClass,
					textClass,
					ratio,
				),
				Hint: "Use a dark text foreground token (e.g. text-black) or define a high-contrast semantic token in global.css.",
			},
		}
	}

	return nil
}
