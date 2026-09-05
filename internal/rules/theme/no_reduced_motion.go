package theme

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// NoReducedMotionRule mendeteksi transisi tema (warna, background, atau all) di dalam <style>
// yang tidak dibungkus atau dimitigasi oleh media query prefers-reduced-motion.
type NoReducedMotionRule struct{}

// NewNoReducedMotionRule membuat instance baru NoReducedMotionRule.
func NewNoReducedMotionRule() *NoReducedMotionRule {
	return &NoReducedMotionRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *NoReducedMotionRule) ID() string {
	return "theme.no-reduced-motion"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *NoReducedMotionRule) Description() string {
	return "Detects global theme transitions without prefers-reduced-motion media query wrapping"
}

// Category mengembalikan nama kategori rule.
func (r *NoReducedMotionRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *NoReducedMotionRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki otomatis.
func (r *NoReducedMotionRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"WCAG 2.2 Success Criterion 2.3.3 (Animation from Interactions)",
			"W3C Media Queries Level 5 (prefers-reduced-motion)",
			"Accessible Web Animation & Vestibular Safety Guidelines",
		},
		CoreInvariant: "Global theme and color transitions must be scoped within prefers-reduced-motion: no-preference or mitigated with reduced-motion overrides.",
		Grounding: "Smooth CSS transitions applied to root or theme switching (such as * { transition: background-color 0.3s, color 0.3s; } or transition: all 0.2s) can cause dizziness, headaches, and nausea for users with vestibular disorders.\n\n" +
			"WCAG 2.2 Success Criterion 2.3.3 requires that non-essential animations triggered by user interaction can be turned off or respect system accessibility preferences.\n\n" +
			"Charites enforces wrapping theme transitions in @media (prefers-reduced-motion: no-preference) or providing an explicit @media (prefers-reduced-motion: reduce) { transition: none; } fallback.",
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Unmitigated global theme transition in Astro style",
				Code: `<style>
  * {
    transition: background-color 0.3s ease, color 0.3s ease;
  }
</style>`,
			},
			{
				Language: "tsx",
				Comment:  "Broad transition all without motion preference in TSX",
				Code: `<style>{` + "`" + `
  body {
    transition: all 0.25s ease-in-out;
  }
` + "`" + `}</style>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Theme transition scoped to no-preference media query",
				Code: `<style>
  @media (prefers-reduced-motion: no-preference) {
    * {
      transition: background-color 0.3s ease, color 0.3s ease;
    }
  }
</style>`,
			},
			{
				Language: "astro",
				Comment:  "Explicit reduced-motion override",
				Code: `<style>
  body {
    transition: background-color 0.3s ease;
  }
  @media (prefers-reduced-motion: reduce) {
    body {
      transition: none;
    }
  }
</style>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Vestibular Distress",
				Severity: "MEDIUM",
				Impact:   "Rapid or uncontrolled surface transitions induce disorientation or motion sickness for sensitive users.",
			},
			{
				Vector:   "WCAG 2.2 SC 2.3.3 Non-Compliance",
				Severity: "MEDIUM",
				Impact:   "Failure to honor OS-level accessibility preferences prevents compliance with regulatory accessibility standards.",
			},
		},
	}
}

// Evaluate mengevaluasi node IR untuk memeriksa apakah ada transisi tema tanpa perlindungan reduced-motion.
func (r *NoReducedMotionRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Tag != "style" {
		return nil
	}

	cssText := getStyleNodeText(node)
	if !strings.Contains(cssText, "transition") {
		return nil
	}

	cleaned := stripCSSCommentsString(cssText)
	if hasThemeTransitionWithoutReducedMotion(cleaned) {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Theme transition defined without prefers-reduced-motion media query",
				Hint:     "Wrap transition in \"@media (prefers-reduced-motion: no-preference)\" or add a reduced-motion fallback.",
			},
		}
	}

	return nil
}

func hasThemeTransitionWithoutReducedMotion(css string) bool {
	if strings.Contains(css, "prefers-reduced-motion") {
		return false
	}
	return containsThemeTransition(css)
}

func containsThemeTransition(css string) bool {
	remaining := css
	for {
		idx := strings.Index(remaining, "transition")
		if idx == -1 {
			break
		}
		after := remaining[idx+len("transition"):]
		colonIdx := strings.IndexByte(after, ':')
		if colonIdx == -1 {
			remaining = after
			continue
		}

		// Ensure it's transition: or transition-property:
		propPrefix := strings.TrimSpace(after[:colonIdx])
		if propPrefix != "" && propPrefix != "-property" {
			remaining = after[colonIdx+1:]
			continue
		}

		valRest := after[colonIdx+1:]
		semiIdx := strings.IndexByte(valRest, ';')
		closeBrace := strings.IndexByte(valRest, '}')
		endIdx := semiIdx
		if endIdx == -1 || (closeBrace != -1 && closeBrace < endIdx) {
			endIdx = closeBrace
		}

		var declVal string
		if endIdx != -1 {
			declVal = valRest[:endIdx]
			remaining = valRest[endIdx:]
		} else {
			declVal = valRest
			remaining = ""
		}

		declVal = strings.ToLower(declVal)
		if strings.Contains(declVal, "background") ||
			strings.Contains(declVal, "color") ||
			strings.Contains(declVal, "all") ||
			strings.Contains(declVal, "theme") {
			return true
		}
	}
	return false
}
