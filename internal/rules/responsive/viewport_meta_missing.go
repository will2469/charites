package responsive

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// ViewportMetaMissingRule memeriksa apakah tag <meta name="viewport"> menyertakan
// parameter konfigurasi mobile penting: 'width=device-width' dan 'viewport-fit=cover'.
type ViewportMetaMissingRule struct{}

// NewViewportMetaMissingRule membuat instance baru dari ViewportMetaMissingRule.
func NewViewportMetaMissingRule() *ViewportMetaMissingRule {
	return &ViewportMetaMissingRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *ViewportMetaMissingRule) ID() string {
	return "responsive.viewport-meta-missing"
}

// Description mengembalikan ringkasan aturan.
func (r *ViewportMetaMissingRule) Description() string {
	return "Warns when <meta name=\"viewport\"> is missing width=device-width or viewport-fit=cover"
}

// Category mengembalikan nama kategori rule.
func (r *ViewportMetaMissingRule) Category() string {
	return "responsive"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *ViewportMetaMissingRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ViewportMetaMissingRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"HTML Living Standard (Viewport Meta Element)",
			"Apple WebKit Safe Area Viewport Expansion Guidelines",
			"W3C CSS Device Adaptation Module Level 1",
		},
		CoreInvariant: "<meta name=\"viewport\"> elements must declare both 'width=device-width' (preventing 980px virtual desktop zoom fallback) and 'viewport-fit=cover' (enabling safe area inset expansion on notched displays).",
		Grounding: "Omitting 'width=device-width' causes mobile browsers (WebKit and Chromium) to fall back to a 980px virtual desktop viewport, forcing users to pinch-zoom and rendering responsive media queries ineffective.\n\n" +
			"Omitting 'viewport-fit=cover' causes CSS safe area variables (env(safe-area-inset-*)) to evaluate to 0px on iOS devices, resulting in white letterboxing around sensor cutouts and disabling hardware-safe bottom docks.\n\n" +
			"Declaring both parameters ensures proportionate rendering across all smartphone screen densities and full hardware edge-to-edge layout immersion.",
		Risks: []ir.RiskItem{
			{
				Vector:   "980px Virtual Desktop Zoom Fallback",
				Severity: "HIGH",
				Impact:   "Mobile browsers scale down pages to fit a 980px virtual width, making text unreadable and disabling responsive layouts.",
			},
			{
				Vector:   "Safe Area Inset Failure and Letterboxing",
				Severity: "MEDIUM",
				Impact:   "CSS env(safe-area-inset-bottom) evaluates to 0px, causing bottom bars to be obscured by hardware home indicators and displaying letterbox bars.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Viewport meta tag missing viewport-fit=cover",
				Code:     `<meta name="viewport" content="width=device-width, initial-scale=1.0" />`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Complete mobile viewport configuration with device width and safe-area expansion",
				Code:     `<meta name="viewport" content="width=device-width, initial-scale=1.0, viewport-fit=cover" />`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen <meta name="viewport"> memiliki atribut content yang memuat
// 'width=device-width' dan 'viewport-fit=cover'.
func (r *ViewportMetaMissingRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || node.Tag != "meta" || len(node.Attributes) == 0 {
		return nil
	}

	nameVal := cleanAttrValue(node.Attributes["name"])
	if !strings.EqualFold(nameVal, "viewport") {
		return nil
	}

	contentVal := cleanAttrValue(node.Attributes["content"])
	hasWidth := hasDeviceWidth(contentVal)
	hasCover := hasViewportFitCover(contentVal)

	if hasWidth && hasCover {
		return nil
	}

	msg, hint := formatViewportMetaDiagnostic(contentVal, hasWidth, hasCover)

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  msg,
			Hint:     hint,
		},
	}
}

func formatViewportMetaDiagnostic(content string, hasWidth, hasCover bool) (string, string) {
	if !hasWidth && !hasCover {
		return fmt.Sprintf("<meta name=\"viewport\"> tag is missing both 'width=device-width' and 'viewport-fit=cover' in content %q. Mobile browsers will scale to 980px desktop width and safe area insets will be disabled.", content),
			"Update content to include both parameters, e.g. 'width=device-width, initial-scale=1.0, viewport-fit=cover'."
	}
	if !hasWidth {
		return fmt.Sprintf("<meta name=\"viewport\"> tag is missing 'width=device-width' in content %q. Mobile browsers will scale page to 980px desktop width instead of using the physical device width.", content),
			"Add 'width=device-width' to the viewport content attribute."
	}
	return fmt.Sprintf("<meta name=\"viewport\"> tag is missing 'viewport-fit=cover' in content %q. Safe area insets (env(safe-area-inset-*)) will evaluate to 0px on iOS devices, causing letterboxing and overlapping bottom docks.", content),
		"Add 'viewport-fit=cover' to the viewport content attribute to enable full-bleed display and safe area insets."
}
