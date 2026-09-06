package browser

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// SafariOnlyAPIRule mendeteksi pemanggilan API eksklusif Apple WebKit/Safari tanpa pengecekan ketersediaan atau fallback W3C standar.
type SafariOnlyAPIRule struct{}

// NewSafariOnlyAPIRule membuat instance baru dari SafariOnlyAPIRule.
func NewSafariOnlyAPIRule() *SafariOnlyAPIRule {
	return &SafariOnlyAPIRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *SafariOnlyAPIRule) ID() string {
	return "browser.safari-only-api"
}

// Description mengembalikan ringkasan aturan.
func (r *SafariOnlyAPIRule) Description() string {
	return "Flags unguarded Apple WebKit/Safari-proprietary APIs without universal web platform fallbacks"
}

// Category mengembalikan nama kategori rule.
func (r *SafariOnlyAPIRule) Category() string {
	return "browser"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *SafariOnlyAPIRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *SafariOnlyAPIRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Web App Manifest (display-mode: standalone)",
			"W3C Pointer Events Level 3",
			"Apple Pay on the Web Guidelines (Feature Detection Requirements)",
		},
		CoreInvariant: "Direct invocation of Apple Safari-exclusive APIs ('navigator.standalone', 'ApplePaySession', iOS gesture events) must provide W3C standard fallbacks for Android and desktop platforms.",
		Grounding: "Apple WebKit includes proprietary features designed exclusively for iOS/macOS Safari.\n\n" +
			"Calling 'navigator.standalone' directly will always return undefined on Android Chrome, while calling 'ApplePaySession.canMakePayments()' without checking 'window.ApplePaySession' throws ReferenceError crashes.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Crash on Non-Apple Platforms",
				Severity: "MEDIUM",
				Impact:   "Calling 'ApplePaySession.canMakePayments()' on Android Chrome, Windows Edge, or Linux Firefox throws 'ReferenceError: ApplePaySession is not defined'.",
			},
			{
				Vector:   "Broken PWA Detection on Android",
				Severity: "MEDIUM",
				Impact:   "Using 'navigator.standalone' fails to detect installed PWAs on Android, where 'display-mode: standalone' is the universal standard.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "javascript",
				Comment:  "Direct invocation of ApplePaySession without availability check",
				Code: `if (ApplePaySession.canMakePayments()) {
  showApplePayButton();
}`,
			},
			{
				Language: "javascript",
				Comment:  "Relying solely on iOS-proprietary navigator.standalone",
				Code:     `const isAppMode = navigator.standalone;`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "javascript",
				Comment:  "Defensive feature guard before ApplePaySession invocation",
				Code: `if (typeof window !== "undefined" && window.ApplePaySession && window.ApplePaySession.canMakePayments()) {
  showApplePayButton();
}`,
			},
			{
				Language: "javascript",
				Comment:  "Standard W3C display-mode with legacy iOS fallback",
				Code: `const isAppMode = (typeof window !== "undefined" && window.matchMedia("(display-mode: standalone)").matches) ||
  (typeof navigator !== "undefined" && Boolean(navigator.standalone));`,
			},
		},
	}
}

type safariAPIItem struct {
	trigger     string
	name        string
	standardAPI string
}

var safariAPIs = [...]safariAPIItem{
	{trigger: "navigator.standalone", name: "navigator.standalone", standardAPI: "window.matchMedia('(display-mode: standalone)')"},
	{trigger: "applepaysession", name: "ApplePaySession", standardAPI: "window.ApplePaySession guard"},
	{trigger: "gesturestart", name: "gesturestart event", standardAPI: "Pointer Events"},
	{trigger: "gesturechange", name: "gesturechange event", standardAPI: "Pointer Events"},
	{trigger: "gestureend", name: "gestureend event", standardAPI: "Pointer Events"},
	{trigger: "window.safari.pushnotification", name: "window.safari.pushNotification", standardAPI: "Standard Web Push (PushManager)"},
}

// Evaluate memeriksa apakah kode script atau event handler memanggil Safari API tanpa fallback.
func (r *SafariOnlyAPIRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil {
		return nil
	}

	// 1. Cek atribut event handler pada elemen JSX / HTML
	if node.Attributes != nil {
		for attrName, attrVal := range node.Attributes {
			if isScriptAttribute(attrName) {
				if diags := r.evaluateScriptContent(node, attrVal, false); len(diags) > 0 {
					return diags
				}
			}
		}
	}

	// 2. Cek blok <script>
	if strings.ToLower(node.Tag) == "script" {
		scriptText := getStyleNodeText(node)
		return r.evaluateScriptContent(node, scriptText, true)
	}

	return nil
}

func (r *SafariOnlyAPIRule) evaluateScriptContent(node *ir.Node, script string, isScriptBlock bool) []ir.Diagnostic {
	lower := strings.ToLower(script)
	if lower == "" || !hasAnySafariAPI(lower) {
		return nil
	}

	if strings.Contains(lower, "try {") && strings.Contains(lower, "} catch") {
		return nil
	}

	//nolint:prealloc // zero-alloc on clean nodes required by QUAL-03
	var diags []ir.Diagnostic

	for _, item := range safariAPIs {
		if !strings.Contains(lower, item.trigger) {
			continue
		}

		if isSafariGuarded(lower, item) {
			continue
		}

		line := resolveTriggerLine(node, script, item.trigger, isScriptBlock)

		diags = append(diags, ir.Diagnostic{
			Line:     line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Unguarded usage of Safari/Apple WebKit proprietary API '" + item.name + "'. Fails or returns undefined on Android and Windows.",
			Hint:     "Guard with runtime check (e.g. 'window." + item.name + "') or use standard W3C alternative ('" + item.standardAPI + "').",
		})
	}

	return diags
}

func hasAnySafariAPI(lower string) bool {
	for i := range safariAPIs {
		if strings.Contains(lower, safariAPIs[i].trigger) {
			return true
		}
	}
	return false
}

func isSafariGuarded(lower string, item safariAPIItem) bool {
	// 1. Optional chaining: e.g. navigator.standalone?.
	if strings.Contains(lower, item.trigger+"?.") {
		return true
	}

	// 2. Pengecekan guard ApplePaySession: window.ApplePaySession atau 'ApplePaySession' in window
	if item.trigger == "applepaysession" {
		if strings.Contains(lower, "window.applepaysession") ||
			strings.Contains(lower, "'applepaysession' in") ||
			strings.Contains(lower, "\"applepaysession\" in") ||
			strings.Contains(lower, "typeof applepaysession") {
			return true
		}
	}

	// 3. Pengecekan navigator.standalone jika didampingi matchMedia display-mode
	if item.trigger == "navigator.standalone" {
		if strings.Contains(lower, "display-mode: standalone") ||
			strings.Contains(lower, "'standalone' in navigator") ||
			strings.Contains(lower, "\"standalone\" in navigator") {
			return true
		}
	}

	// 4. Gesture events guarded jika ada pointerdown atau touchstart
	if strings.HasPrefix(item.trigger, "gesture") {
		if strings.Contains(lower, "pointer") || strings.Contains(lower, "touch") {
			return true
		}
	}

	// 5. Guard window.safari
	if item.trigger == "window.safari.pushnotification" {
		if strings.Contains(lower, "typeof window.safari") ||
			strings.Contains(lower, "'safari' in window") ||
			strings.Contains(lower, "window.safari?.") {
			return true
		}
	}

	return false
}
