package pwa

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// IconMaskableMissingRule mendeteksi Web App Manifest yang memiliki ikon tetapi tidak menyertakan
// ikon adaptif dengan purpose: "maskable" untuk Android homescreen launcher.
type IconMaskableMissingRule struct{}

// NewIconMaskableMissingRule membuat instance baru dari IconMaskableMissingRule.
func NewIconMaskableMissingRule() *IconMaskableMissingRule {
	return &IconMaskableMissingRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *IconMaskableMissingRule) ID() string {
	return "pwa.icon-maskable-missing"
}

// Description mengembalikan ringkasan aturan.
func (r *IconMaskableMissingRule) Description() string {
	return "Warns when a Web App Manifest defines icons but none has purpose: 'maskable' for Android adaptive launcher icons"
}

// Category mengembalikan nama kategori rule.
func (r *IconMaskableMissingRule) Category() string {
	return "pwa"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *IconMaskableMissingRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *IconMaskableMissingRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Web App Manifest Specification (Adaptive Icon Masking)",
			"Google Android Maskable Icons Specification",
			"Android Oreo+ Adaptive Launcher Icon Architecture",
		},
		CoreInvariant: "When a Web App Manifest defines icons, at least one icon must declare 'purpose: \"maskable\"' to prevent Android launcher letterboxing.",
		Grounding: "Starting in Android 8.0 Oreo, native device launchers crop application icons according to user-selected device masks (circles, squircles, rounded rectangles).\n\n" +
			"When a PWA provides only standard icons (purpose: 'any' or omitted purpose), modern Android launchers place the icon inside a small white square box (letterboxing) to fit the mask. This disrupts the visual consistency of native mobile app trays.\n\n" +
			"Providing at least one icon with 'purpose: \"maskable\"' (with an appropriate safe zone margin) ensures the launcher can scale and mask the icon seamlessly to fill the full shape.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Letterboxed Android Launcher Icons",
				Severity: "MEDIUM",
				Impact:   "PWA icon appears inside an awkward white square box on Android device home screens.",
			},
			{
				Vector:   "Degraded Native Visual Immersion",
				Severity: "LOW",
				Impact:   "Breaks aesthetic parity with native Android apps installed from Google Play.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Manifest defines icons without any maskable purpose",
				Code: `<script type="application/manifest+json">
  {JSON.stringify({
    name: "Desa Digital",
    start_url: "/",
    display: "standalone",
    icons: [
      { src: "/icon-512.png", sizes: "512x512", type: "image/png" }
    ]
  })}
</script>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Manifest includes an adaptive icon with purpose: maskable",
				Code: `<script type="application/manifest+json">
  {JSON.stringify({
    name: "Desa Digital",
    start_url: "/",
    display: "standalone",
    icons: [
      { src: "/icon-512.png", sizes: "512x512", type: "image/png" },
      { src: "/icon-512-maskable.png", sizes: "512x512", type: "image/png", purpose: "maskable" }
    ]
  })}
</script>`,
			},
		},
	}
}

// Evaluate memeriksa apakah manifest yang mendefinisikan icons memiliki ikon maskable.
func (r *IconMaskableMissingRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || !isManifestScript(node) {
		return nil
	}

	manifest, ok := parseManifest(node)
	if !ok || manifest == nil {
		return nil
	}

	if len(manifest.Icons) == 0 {
		return nil
	}

	if hasMaskableIcon(manifest.Icons) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Web App Manifest icons do not define an adaptive icon with purpose: 'maskable'. Android launcher will display a letterboxed icon.",
			Hint:     "Add an icon entry with purpose: 'maskable' (recommended 512x512 PNG) to adapt seamlessly to Android icon masks.",
		},
	}
}

func hasMaskableIcon(icons []ManifestIcon) bool {
	for _, icon := range icons {
		purpose := strings.ToLower(icon.Purpose)
		if strings.Contains(purpose, "maskable") {
			return true
		}
	}
	return false
}
