package pwa

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// ServiceWorkerNoOfflineFallbackRule mendeteksi skrip Service Worker yang mencegat
// event fetch tetapi tidak menyediakan strategi cadangan cache offline atau penanganan error.
type ServiceWorkerNoOfflineFallbackRule struct{}

// NewServiceWorkerNoOfflineFallbackRule membuat instance baru dari ServiceWorkerNoOfflineFallbackRule.
func NewServiceWorkerNoOfflineFallbackRule() *ServiceWorkerNoOfflineFallbackRule {
	return &ServiceWorkerNoOfflineFallbackRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *ServiceWorkerNoOfflineFallbackRule) ID() string {
	return "pwa.service-worker-no-offline-fallback"
}

// Description mengembalikan ringkasan aturan.
func (r *ServiceWorkerNoOfflineFallbackRule) Description() string {
	return "Warns when a Service Worker intercepts fetch events without providing an offline cache fallback or failure handler"
}

// Category mengembalikan nama kategori rule.
func (r *ServiceWorkerNoOfflineFallbackRule) Category() string {
	return "pwa"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *ServiceWorkerNoOfflineFallbackRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ServiceWorkerNoOfflineFallbackRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Service Workers 1 (Offline Resilience Architecture)",
			"W3C Cache Storage Specification (Offline Asset Fallback)",
			"Google Chrome PWA Reliability Criteria (Offline Support)",
		},
		CoreInvariant: "Service Worker fetch event handlers must implement an offline cache fallback (e.g. caches.match) or failure catch handler instead of bare pass-through fetch interception.",
		Grounding: "In spotty or rural mobile network conditions (3G/4G signal drops), a Service Worker that intercepts fetch events without an offline cache strategy causes the browser to immediately display a connection-lost screen (the offline dinosaur page).\n\n" +
			"Providing a resilient cache-first or network-first fallback mechanism guarantees that the application shell and cached pages remain accessible even when completely disconnected from the Internet.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Immediate Offline Blackout",
				Severity: "MEDIUM",
				Impact:   "Users opening the PWA without network connectivity encounter a browser network failure screen instead of cached application content.",
			},
			{
				Vector:   "PWA Installability Rejection",
				Severity: "LOW",
				Impact:   "Mobile browsers may downgrade or reject full PWA installation status due to failing offline resilience audits.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Pass-through fetch interception without offline cache fallback",
				Code: `<script>
  self.addEventListener("fetch", (event) => {
    event.respondWith(fetch(event.request));
  });
</script>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Cache fallback provided via caches.match with network fallback",
				Code: `<script>
  self.addEventListener("fetch", (event) => {
    event.respondWith(
      caches.match(event.request).then((cached) => {
        return cached || fetch(event.request).catch(() => caches.match("/offline.html"));
      })
    );
  });
</script>`,
			},
		},
	}
}

// Evaluate memeriksa apakah skrip Service Worker yang mencegat fetch memiliki strategi fallback cache.
func (r *ServiceWorkerNoOfflineFallbackRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || !strings.EqualFold(node.Tag, "script") {
		return nil
	}

	txt := extractScriptText(node)
	if !hasSWFetchListener(txt) || !interceptsFetch(txt) {
		return nil
	}

	if hasSWOfflineFallback(txt) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Service Worker fetch listener intercepts requests with respondWith() but provides no offline cache fallback (caches.match) or failure handler (.catch). Disconnected devices will display a network failure screen.",
			Hint:     "Implement an offline cache-first strategy (caches.match(event.request)) or append a network failure fallback (.catch(() => caches.match('/offline.html'))).",
		},
	}
}
