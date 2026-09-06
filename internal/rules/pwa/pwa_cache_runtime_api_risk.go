package pwa

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// CacheRuntimeAPIRiskRule mendeteksi akses ke API thread utama (window, document, localStorage)
// di dalam lingkungan eksekusi Service Worker yang akan memicu ReferenceError pada runtime.
type CacheRuntimeAPIRiskRule struct{}

// NewCacheRuntimeAPIRiskRule membuat instance baru dari CacheRuntimeAPIRiskRule.
func NewCacheRuntimeAPIRiskRule() *CacheRuntimeAPIRiskRule {
	return &CacheRuntimeAPIRiskRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *CacheRuntimeAPIRiskRule) ID() string {
	return "pwa.pwa-cache-runtime-api-risk"
}

// Description mengembalikan ringkasan aturan.
func (r *CacheRuntimeAPIRiskRule) Description() string {
	return "Prevents access to main-thread DOM and synchronous Web Storage APIs (window, document, localStorage) inside Service Worker scripts"
}

// Category mengembalikan nama kategori rule.
func (r *CacheRuntimeAPIRiskRule) Category() string {
	return "pwa"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *CacheRuntimeAPIRiskRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *CacheRuntimeAPIRiskRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Service Workers (ServiceWorkerGlobalScope Execution Context)",
			"HTML Living Standard (Dedicated Worker & Web Worker Isolation)",
			"W3C Web Storage Specification (Thread Affinity Limitations)",
		},
		CoreInvariant: "Service Worker scripts must not access main-thread DOM or synchronous storage APIs (window, document, localStorage, sessionStorage, alert, confirm, prompt).",
		Grounding: "Service Workers run in a distinct background worker thread (ServiceWorkerGlobalScope) that is entirely decoupled from the browser UI thread.\n\n" +
			"Attempting to access DOM APIs (window, document) or synchronous storage (localStorage, sessionStorage) in a Service Worker throws an immediate fatal ReferenceError at runtime, aborting worker installation and breaking all offline caching capabilities.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Immediate Service Worker Installation Crash",
				Severity: "HIGH",
				Impact:   "Worker script fails during parsing/evaluation with Uncaught ReferenceError: window is not defined, completely disabling offline caching.",
			},
			{
				Vector:   "Broken Background Push/Sync Functionality",
				Severity: "HIGH",
				Impact:   "Background sync and push notifications fail to initialize because the worker thread crashed upon bootstrap.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Service Worker script attempting to access window and localStorage",
				Code: `<script>
  self.addEventListener("install", (event) => {
    const token = localStorage.getItem("token");
    window.location.reload();
  });
</script>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Compliant Service Worker using Cache Storage and Worker primitives",
				Code: `<script>
  self.addEventListener("install", (event) => {
    event.waitUntil(
      caches.open("v1").then((cache) => cache.addAll(["/", "/offline.html"]))
    );
  });
</script>`,
			},
		},
	}
}

// Evaluate memeriksa apakah skrip Service Worker mengakses API terlarang thread utama.
func (r *CacheRuntimeAPIRiskRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || !strings.EqualFold(node.Tag, "script") {
		return nil
	}

	txt := extractScriptText(node)
	if !isWorkerScope(txt) {
		return nil
	}

	forbidden := collectForbiddenWorkerAPIs(txt)
	if len(forbidden) == 0 {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Forbidden main-thread API access detected inside Service Worker script: %s. ServiceWorkerGlobalScope has no access to DOM or synchronous Web Storage and will throw runtime ReferenceErrors.", strings.Join(forbidden, ", ")),
			Hint:     "Use Service Worker Cache Storage API (caches.open), IndexedDB, or clients.matchAll() with postMessage() to communicate with window clients.",
		},
	}
}
