package pwa

import (
	"github.com/will2469/charites/internal/ir"
)

// ServiceWorkerMissingRule memeriksa apakah dokumen HTML yang mendeklarasikan link manifest
// juga mendaftarkan Service Worker untuk mengaktifkan offline caching dan instalasi PWA.
type ServiceWorkerMissingRule struct{}

// NewServiceWorkerMissingRule membuat instance baru dari ServiceWorkerMissingRule.
func NewServiceWorkerMissingRule() *ServiceWorkerMissingRule {
	return &ServiceWorkerMissingRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *ServiceWorkerMissingRule) ID() string {
	return "pwa.service-worker-missing"
}

// Description mengembalikan ringkasan aturan.
func (r *ServiceWorkerMissingRule) Description() string {
	return "Warns when an HTML document head links to a Web App Manifest but lacks a Service Worker registration in the document"
}

// Category mengembalikan nama kategori rule.
func (r *ServiceWorkerMissingRule) Category() string {
	return "pwa"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *ServiceWorkerMissingRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ServiceWorkerMissingRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Service Workers (Document Integration)",
			"W3C Web App Manifest (Installation Requirements)",
			"Google Chrome PWA Criteria (Offline Capability)",
		},
		CoreInvariant: "When an HTML document head links to a Web App Manifest, the document must register a Service Worker via navigator.serviceWorker.register or an external worker script.",
		Grounding: "A Progressive Web App requires a registered Service Worker to cache shell assets, intercept network requests during outages, and satisfy full mobile installability audits.\n\n" +
			"Linking a manifest file without registering a Service Worker leaves the application behaving like a conventional static website, incapable of offline execution or background updates.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Lack of Offline Capability",
				Severity: "MEDIUM",
				Impact:   "Users cannot open or use the application when device connectivity is intermittent or offline.",
			},
			{
				Vector:   "Failed Installability Criteria",
				Severity: "LOW",
				Impact:   "Modern mobile browsers will not trigger automated installation banners without an active Service Worker.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Document head links to manifest but no Service Worker is registered",
				Code: `<head>
  <title>Layanan Desa</title>
  <link rel="manifest" href="/manifest.webmanifest" />
</head>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Document head links to manifest and registers a Service Worker",
				Code: `<head>
  <title>Layanan Desa</title>
  <link rel="manifest" href="/manifest.webmanifest" />
  <script>
    if ('serviceWorker' in navigator) {
      navigator.serviceWorker.register('/sw.js').catch(console.error);
    }
  </script>
</head>`,
			},
		},
	}
}

// Evaluate memeriksa apakah dokumen dengan manifest mendaftarkan Service Worker.
func (r *ServiceWorkerMissingRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || !isHeadElement(node.Tag) {
		return nil
	}

	if !hasManifestLink(node) {
		return nil
	}

	if hasServiceWorkerRegistration(node) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "PWA document head declares a manifest link but lacks a Service Worker registration in the document. Without a registered Service Worker, the application cannot cache shell assets or operate offline.",
			Hint:     "Register a Service Worker in a <script> tag (e.g. navigator.serviceWorker.register('/sw.js')) or include a service worker registration script.",
		},
	}
}
