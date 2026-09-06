package pwa

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// ServiceWorkerRegistrationRule memeriksa apakah pendaftaran Service Worker
// dilindungi dengan pemeriksaan feature detection dan penanganan error yang andal.
type ServiceWorkerRegistrationRule struct{}

// NewServiceWorkerRegistrationRule membuat instance baru dari ServiceWorkerRegistrationRule.
func NewServiceWorkerRegistrationRule() *ServiceWorkerRegistrationRule {
	return &ServiceWorkerRegistrationRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *ServiceWorkerRegistrationRule) ID() string {
	return "pwa.service-worker-registration"
}

// Description mengembalikan ringkasan aturan.
func (r *ServiceWorkerRegistrationRule) Description() string {
	return "Warns when Service Worker registration lacks feature detection ('serviceWorker' in navigator) or error handling (.catch)"
}

// Category mengembalikan nama kategori rule.
func (r *ServiceWorkerRegistrationRule) Category() string {
	return "pwa"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *ServiceWorkerRegistrationRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ServiceWorkerRegistrationRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Service Workers (Registration Lifecycle & Error Handling)",
			"MDN Progressive Web App Guides (Registering a Service Worker Safely)",
			"Google Web Fundamentals (Service Worker Reliability)",
		},
		CoreInvariant: "Calls to navigator.serviceWorker.register must be guarded by feature detection ('serviceWorker' in navigator) and handled with error callbacks (.catch or try/catch).",
		Grounding: "Calling navigator.serviceWorker.register() without feature detection triggers fatal runtime TypeErrors in legacy browsers, restricted WebViews, or non-secure HTTP contexts.\n\n" +
			"Furthermore, failing to handle registration failure (.catch or try/catch) causes unhandled promise rejections that can disrupt analytics scripts and break client-side bootstrapping.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Runtime Script Crash on Older Browsers",
				Severity: "MEDIUM",
				Impact:   "Browsers or WebViews lacking Service Worker support crash with Uncaught TypeError: Cannot read properties of undefined.",
			},
			{
				Vector:   "Silent Unhandled Promise Rejections",
				Severity: "LOW",
				Impact:   "Registration rejections (e.g. 404 or SSL errors) pollute telemetry logs and fail to notify diagnostic listeners.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Unsafe registration without feature detection and error handling",
				Code: `<script>
  navigator.serviceWorker.register('/sw.js');
</script>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Safe registration with feature detection and error handler",
				Code: `<script>
  if ('serviceWorker' in navigator) {
    navigator.serviceWorker.register('/sw.js')
      .then((reg) => console.log('SW registered:', reg.scope))
      .catch((err) => console.error('SW registration failed:', err));
  }
</script>`,
			},
		},
	}
}

// Evaluate memeriksa apakah pendaftaran Service Worker memiliki guard dan error handling.
func (r *ServiceWorkerRegistrationRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || !strings.EqualFold(node.Tag, "script") {
		return nil
	}

	txt := extractScriptText(node)
	if !strings.Contains(txt, "serviceWorker.register") {
		return nil
	}

	hasFeature := hasSWFeatureDetection(txt)
	hasError := hasSWErrorHandling(txt)
	if hasFeature && hasError {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Service Worker registration is missing feature detection ('serviceWorker' in navigator) or error handling (.catch / try-catch). Unchecked registration causes runtime TypeErrors in unsupported environments and unhandled promise rejections.",
			Hint:     "Wrap registration in 'if (\x27serviceWorker\x27 in navigator)' and append '.catch(err => console.error(err))' or wrap in a try/catch block.",
		},
	}
}
