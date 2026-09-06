package mobile

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// OrientationLockRiskRule mendeteksi penguncian orientasi layar (screen orientation lock)
// yang membatasi aksesibilitas bagi pengguna dengan dudukan perangkat mobile khusus (WCAG 2.2 SC 1.3.4).
type OrientationLockRiskRule struct{}

// NewOrientationLockRiskRule membuat instance baru dari OrientationLockRiskRule.
func NewOrientationLockRiskRule() *OrientationLockRiskRule {
	return &OrientationLockRiskRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *OrientationLockRiskRule) ID() string {
	return "mobile.orientation-lock-risk"
}

// Description mengembalikan ringkasan aturan.
func (r *OrientationLockRiskRule) Description() string {
	return "Advises against rigid screen orientation locking which restricts accessibility for mounted or assistive mobile setups (WCAG 2.2 SC 1.3.4)"
}

// Category mengembalikan nama kategori rule.
func (r *OrientationLockRiskRule) Category() string {
	return "mobile"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (info).
func (r *OrientationLockRiskRule) DefaultSeverity() ir.Severity {
	return ir.SeverityInfo
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *OrientationLockRiskRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Web Content Accessibility Guidelines (WCAG) 2.2 SC 1.3.4 (Orientation - Level AA)",
			"W3C Screen Orientation API (ScreenOrientation.lock)",
			"Google Web Accessibility (Orientation Invariants)",
		},
		CoreInvariant: "Applications must not rigidly lock display orientation to portrait or landscape unless essential to the core functionality (e.g. bank check capture or piano keyboard).",
		Grounding: "Locking mobile orientation via 'screen.orientation.lock(\"portrait\")' prevents users with assistive needs from accessing content.\n\n" +
			"Users who have smartphones mounted horizontally on wheelchairs, bed frames, or vehicle dashboards cannot rotate their devices.\n\n" +
			"Web interfaces should adapt fluidly using responsive CSS (e.g. 'landscape:flex-row') rather than programmatically forbidding device rotation.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Assistive Technology Exclusion",
				Severity: "LOW",
				Impact:   "Users with fixed horizontal device mounts are unable to view or operate the application naturally.",
			},
			{
				Vector:   "Unintended Script Errors on Unsupported Browsers",
				Severity: "LOW",
				Impact:   "Calling orientation lock on Safari iOS or unsupported browsers triggers unhandled promise rejections.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Programmatic orientation lock forces portrait mode",
				Code: `useEffect(() => {
  if (screen.orientation && screen.orientation.lock) {
    screen.orientation.lock("portrait").catch(() => {});
  }
}, []);`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Fluid responsive layout adapting naturally to landscape orientation",
				Code: `<div className="flex flex-col landscape:flex-row gap-4 p-4">
  <aside className="w-full landscape:w-64">Navigasi</aside>
  <main className="flex-1">Konten Utama</main>
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah node memanggil API screen orientation lock atau meta tag terkait.
func (r *OrientationLockRiskRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil {
		return nil
	}

	if isOrientationLockMeta(node) || containsOrientationLockCode(node.RawClasses) || isOrientationLockInAttrs(node.Attributes) {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Screen orientation lock detected. Rigid orientation locking restricts accessibility for mounted or assistive mobile setups (WCAG 2.2 SC 1.3.4 Orientation).",
				Hint:     "Avoid locking orientation programmatically. Design responsive layouts that adapt naturally via CSS (e.g. 'landscape:flex-row').",
			},
		}
	}

	return nil
}

func isOrientationLockMeta(node *ir.Node) bool {
	if node.Type != ir.NodeElement || !strings.EqualFold(node.Tag, "meta") || node.Attributes == nil {
		return false
	}

	name, hasName := node.Attributes["name"]
	if !hasName {
		return false
	}
	cleanName := cleanAttrValue(name)
	if cleanName != "screen-orientation" && cleanName != "orientation" {
		return false
	}

	content, hasContent := node.Attributes["content"]
	if !hasContent {
		return false
	}
	cleanContent := cleanAttrValue(content)
	return cleanContent == "portrait" || cleanContent == "landscape"
}

func containsOrientationLockCode(code string) bool {
	if !strings.Contains(code, "orientation") {
		return false
	}
	return strings.Contains(code, "screen.orientation.lock") ||
		strings.Contains(code, "orientation.lock(") ||
		strings.Contains(code, "screen.lockOrientation(")
}

func isOrientationLockInAttrs(attrs map[string]string) bool {
	if attrs == nil {
		return false
	}
	for _, val := range attrs {
		if containsOrientationLockCode(val) {
			return true
		}
	}
	return false
}
