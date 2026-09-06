package pwa

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// ManifestRequiredFieldsMissingRule memeriksa apakah Web App Manifest mendefinisikan seluruh
// field wajib W3C untuk memenuhi kriteria instalasi PWA pada perangkat seluler.
type ManifestRequiredFieldsMissingRule struct{}

// NewManifestRequiredFieldsMissingRule membuat instance baru dari ManifestRequiredFieldsMissingRule.
func NewManifestRequiredFieldsMissingRule() *ManifestRequiredFieldsMissingRule {
	return &ManifestRequiredFieldsMissingRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *ManifestRequiredFieldsMissingRule) ID() string {
	return "pwa.manifest-required-fields-missing"
}

// Description mengembalikan ringkasan aturan.
func (r *ManifestRequiredFieldsMissingRule) Description() string {
	return "Errors when a Web App Manifest definition is missing required fields (name/short_name, start_url, display, icons)"
}

// Category mengembalikan nama kategori rule.
func (r *ManifestRequiredFieldsMissingRule) Category() string {
	return "pwa"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *ManifestRequiredFieldsMissingRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ManifestRequiredFieldsMissingRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Web App Manifest Specification Section 5 (Manifest Members)",
			"Google Chrome Web App Installability Criteria",
			"W3C Application Lifecycle & Installation Architecture",
		},
		CoreInvariant: "Web App Manifest declarations (<script type=\"application/manifest+json\">) must declare required installability fields: 'name' (or 'short_name'), 'start_url', 'display', and at least one icon in 'icons'.",
		Grounding: "For mobile operating systems (Android, iOS) and modern web engines to recognize a website as an installable Progressive Web App, the Web App Manifest must declare minimum installability metadata.\n\n" +
			"Omitting 'name' or 'short_name' leaves the OS homescreen launcher with an empty or broken app label. Missing 'start_url' prevents the launcher from knowing which route to boot into. Omitting 'display' forces the app to open inside standard browser tabs rather than immersive standalone mode. Omitting 'icons' results in missing or broken asset icons on the user's home screen.\n\n" +
			"Declaring all four fundamental fields ensures reliable PWA install prompts and clean native OS integration.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Browser Installation Prompt Suppression",
				Severity: "HIGH",
				Impact:   "Mobile browsers silently suppress the 'Add to Home Screen' / install banner when manifest required fields are absent.",
			},
			{
				Vector:   "Broken Application Branding",
				Severity: "MEDIUM",
				Impact:   "If installed manually, the web application displays placeholder text and fallback generic browser icons.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Manifest missing start_url, display, and icons",
				Code: `<script type="application/manifest+json">
  {JSON.stringify({
    name: "Desa Digital"
  })}
</script>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Manifest declares all required installability members",
				Code: `<script type="application/manifest+json">
  {JSON.stringify({
    name: "Desa Digital",
    short_name: "Desa",
    start_url: "/",
    display: "standalone",
    icons: [
      { src: "/icon-192.png", sizes: "192x192", type: "image/png" },
      { src: "/icon-512.png", sizes: "512x512", type: "image/png", purpose: "maskable" }
    ]
  })}
</script>`,
			},
		},
	}
}

// Evaluate memeriksa apakah deklarasi manifest script memiliki field wajib lengkap.
func (r *ManifestRequiredFieldsMissingRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || !isManifestScript(node) {
		return nil
	}

	manifest, ok := parseManifest(node)
	if !ok || manifest == nil {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message:  "Web App Manifest script contains empty or unparseable manifest content. PWA installation will fail.",
				Hint:     "Provide a valid JSON manifest declaring 'name', 'start_url', 'display', and 'icons'.",
			},
		}
	}

	missing := checkMissingFields(manifest)
	if len(missing) == 0 {
		return nil
	}

	missingList := strings.Join(missing, ", ")
	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Web App Manifest is missing required installability fields: %s. PWA installation will be rejected by mobile browsers.", missingList),
			Hint:     "Provide 'name' (or 'short_name'), 'start_url', 'display' (e.g. 'standalone'), and at least one icon in 'icons'.",
		},
	}
}

func checkMissingFields(m *WebAppManifest) []string {
	var missing []string
	if m.Name == "" && m.ShortName == "" {
		missing = append(missing, "name/short_name")
	}
	if m.StartURL == "" {
		missing = append(missing, "start_url")
	}
	if m.Display == "" {
		missing = append(missing, "display")
	}
	if len(m.Icons) == 0 {
		missing = append(missing, "icons")
	}
	return missing
}
