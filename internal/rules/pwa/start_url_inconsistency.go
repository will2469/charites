package pwa

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// StartURLInconsistencyRule mendeteksi nilai start_url pada Web App Manifest yang menggunakan
// protokol tidak aman (http://), skema skrip (javascript:), atau path traversal (../).
type StartURLInconsistencyRule struct{}

// NewStartURLInconsistencyRule membuat instance baru dari StartURLInconsistencyRule.
func NewStartURLInconsistencyRule() *StartURLInconsistencyRule {
	return &StartURLInconsistencyRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *StartURLInconsistencyRule) ID() string {
	return "pwa.start-url-inconsistency"
}

// Description mengembalikan ringkasan aturan.
func (r *StartURLInconsistencyRule) Description() string {
	return "Errors when a Web App Manifest start_url uses an insecure protocol (http://), script scheme (javascript:), or path traversal (../)"
}

// Category mengembalikan nama kategori rule.
func (r *StartURLInconsistencyRule) Category() string {
	return "pwa"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *StartURLInconsistencyRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *StartURLInconsistencyRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Web App Manifest Section 5.2 (The start_url member)",
			"W3C Secure Contexts (Mixed Content Mitigation)",
			"RFC 3986 Uniform Resource Identifier (URI): Generic Syntax",
		},
		CoreInvariant: "Web App Manifest 'start_url' must not use insecure HTTP protocols, script URI schemes (javascript:), or directory traversal ('../').",
		Grounding: "The 'start_url' member defines the preferred URL that should be loaded when the user launches the web application from the mobile launcher.\n\n" +
			"According to W3C PWA specifications, PWAs must operate strictly within secure contexts (HTTPS). Setting 'start_url' to an insecure HTTP URL ('http://') causes mobile browsers to block launch execution. Setting 'start_url' to a path traversal sequence ('../') escapes the intended navigation scope and causes unpredictable routing failures.\n\n" +
			"Using a clean relative path (e.g. '/' or '/app') under the secure origin ensures consistent and secure PWA startup.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Insecure Context Launch Failure",
				Severity: "HIGH",
				Impact:   "Mobile browsers block launching PWAs whose start_url does not satisfy Secure Context requirements.",
			},
			{
				Vector:   "Path Traversal Outside Navigation Scope",
				Severity: "HIGH",
				Impact:   "Using '../' breaks origin scope confinement, leading to broken navigation and failed manifest resolution.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Insecure HTTP protocol in start_url",
				Code: `<script type="application/manifest+json">
  {JSON.stringify({
    name: "Desa Digital",
    start_url: "http://desa.id/app",
    display: "standalone",
    icons: [{ src: "/icon-512.png", sizes: "512x512", type: "image/png", purpose: "maskable" }]
  })}
</script>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Valid relative path for start_url",
				Code: `<script type="application/manifest+json">
  {JSON.stringify({
    name: "Desa Digital",
    start_url: "/",
    display: "standalone",
    icons: [{ src: "/icon-512.png", sizes: "512x512", type: "image/png", purpose: "maskable" }]
  })}
</script>`,
			},
		},
	}
}

// Evaluate memeriksa apakah nilai start_url manifest melanggar standar keamanan atau cakupan.
func (r *StartURLInconsistencyRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || !isManifestScript(node) {
		return nil
	}

	manifest, ok := parseManifest(node)
	if !ok || manifest == nil || manifest.StartURL == "" {
		return nil
	}

	reason := validateStartURL(manifest.StartURL)
	if reason == "" {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Web App Manifest 'start_url' %q is invalid or out of scope (%s). Mobile browsers will fail to launch the PWA.", manifest.StartURL, reason),
			Hint:     "Use a valid relative path (e.g. '/' or '/app') under the secure origin scope.",
		},
	}
}

func validateStartURL(startURL string) string {
	lower := strings.ToLower(startURL)
	if strings.HasPrefix(lower, "http://") {
		return "insecure HTTP protocol violates PWA Secure Context requirement"
	}
	if strings.HasPrefix(lower, "javascript:") {
		return "script URI schemes are forbidden"
	}
	if strings.HasPrefix(startURL, "..") || strings.Contains(startURL, "/../") {
		return "path traversal '../' navigates outside application scope"
	}
	return ""
}
