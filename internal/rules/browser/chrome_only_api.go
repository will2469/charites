package browser

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// ChromeOnlyAPIRule mendeteksi pemanggilan API eksklusif Chromium tanpa pengujian ketersediaan atau fallback untuk Firefox dan Safari.
type ChromeOnlyAPIRule struct{}

// NewChromeOnlyAPIRule membuat instance baru dari ChromeOnlyAPIRule.
func NewChromeOnlyAPIRule() *ChromeOnlyAPIRule {
	return &ChromeOnlyAPIRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *ChromeOnlyAPIRule) ID() string {
	return "browser.chrome-only-api"
}

// Description mengembalikan ringkasan aturan.
func (r *ChromeOnlyAPIRule) Description() string {
	return "Flags reliance on Chromium-exclusive APIs without cross-browser fallbacks for Firefox and Safari"
}

// Category mengembalikan nama kategori rule.
func (r *ChromeOnlyAPIRule) Category() string {
	return "browser"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *ChromeOnlyAPIRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ChromeOnlyAPIRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Device Memory Specification (Chromium-only draft)",
			"W3C Network Information API (Excluded from Safari and Firefox)",
			"W3C Web Platform Status (Mozilla & WebKit Positions on Chromium APIs)",
		},
		CoreInvariant: "Reliance on Chromium-exclusive APIs ('navigator.deviceMemory', 'navigator.connection', Web Serial/USB, etc.) must include runtime guards and cross-browser fallbacks.",
		Grounding: "Chromium exposes several proprietary or non-consensus APIs that Mozilla Firefox and Apple WebKit have formally opposed due to privacy or architecture concerns.\n\n" +
			"Direct, unguarded access to 'navigator.deviceMemory' or 'navigator.connection' will fail or return 'undefined', throwing TypeErrors when accessing nested properties.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Runtime Crash on Safari and Firefox",
				Severity: "MEDIUM",
				Impact:   "Accessing 'navigator.connection.effectiveType' throws 'TypeError: Cannot read properties of undefined' on Safari and desktop Firefox.",
			},
			{
				Vector:   "Feature Lockout for Non-Chromium Users",
				Severity: "MEDIUM",
				Impact:   "Users on Safari (iOS/macOS) and Firefox cannot complete critical workflows if File System Access or Web Serial is required without a standard fallback.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "javascript",
				Comment:  "Direct unguarded access to navigator.connection",
				Code:     `const effectiveSpeed = navigator.connection.effectiveType;`,
			},
			{
				Language: "javascript",
				Comment:  "Direct access to navigator.deviceMemory",
				Code:     `const memoryGiB = navigator.deviceMemory;`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "javascript",
				Comment:  "Guarded access with default fallback for non-Chromium browsers",
				Code:     `const effectiveSpeed = (typeof navigator !== "undefined" && navigator.connection?.effectiveType) || "4g";`,
			},
			{
				Language: "javascript",
				Comment:  "Capability guard using in operator",
				Code:     `const memoryGiB = (typeof navigator !== "undefined" && "deviceMemory" in navigator) ? navigator.deviceMemory : 4;`,
			},
		},
	}
}

type chromeAPIItem struct {
	trigger   string
	name      string
	shortName string
}

var chromeAPIs = [...]chromeAPIItem{
	{trigger: "navigator.devicememory", name: "navigator.deviceMemory", shortName: "deviceMemory"},
	{trigger: "navigator.connection", name: "navigator.connection", shortName: "connection"},
	{trigger: "window.chrome.webstore", name: "window.chrome.webstore", shortName: "chrome"},
	{trigger: "navigator.serial", name: "navigator.serial", shortName: "serial"},
	{trigger: "navigator.usb", name: "navigator.usb", shortName: "usb"},
	{trigger: "navigator.bluetooth", name: "navigator.bluetooth", shortName: "bluetooth"},
	{trigger: "navigator.hid", name: "navigator.hid", shortName: "hid"},
}

// Evaluate memeriksa apakah kode script atau event handler memanggil Chromium API tanpa fallback.
func (r *ChromeOnlyAPIRule) Evaluate(node *ir.Node) []ir.Diagnostic {
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

func (r *ChromeOnlyAPIRule) evaluateScriptContent(node *ir.Node, script string, isScriptBlock bool) []ir.Diagnostic {
	lower := strings.ToLower(script)
	if lower == "" || !hasAnyChromeAPI(lower) {
		return nil
	}

	if strings.Contains(lower, "try {") && strings.Contains(lower, "} catch") {
		return nil
	}

	//nolint:prealloc // zero-alloc on clean nodes required by QUAL-03
	var diags []ir.Diagnostic

	for _, item := range chromeAPIs {
		if !strings.Contains(lower, item.trigger) {
			continue
		}

		if isChromeGuarded(lower, item) {
			continue
		}

		line := resolveTriggerLine(node, script, item.trigger, isScriptBlock)

		diags = append(diags, ir.Diagnostic{
			Line:     line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Unguarded usage of Chromium-exclusive API '" + item.name + "'. Unsupported on Firefox and Safari.",
			Hint:     "Guard with runtime capability check (e.g. 'if (\"" + item.shortName + "\" in navigator)') and provide a standard fallback.",
		})
	}

	return diags
}

func hasAnyChromeAPI(lower string) bool {
	for i := range chromeAPIs {
		if strings.Contains(lower, chromeAPIs[i].trigger) {
			return true
		}
	}
	return false
}

func isChromeGuarded(lower string, item chromeAPIItem) bool {
	// 1. Optional chaining: e.g. navigator.connection?.
	if strings.Contains(lower, item.trigger+"?.") || strings.Contains(lower, item.trigger+" ?.") {
		return true
	}

	// 2. 'in' operator: e.g. 'connection' in navigator
	if strings.Contains(lower, "'"+strings.ToLower(item.shortName)+"' in") ||
		strings.Contains(lower, "\""+strings.ToLower(item.shortName)+"\" in") {
		return true
	}

	// 3. typeof check: typeof navigator.connection !== "undefined"
	if strings.Contains(lower, "typeof "+item.trigger) {
		return true
	}

	// 4. If condition or ternary guard
	if strings.Contains(lower, "if ("+item.trigger+")") ||
		strings.Contains(lower, "if ("+item.trigger+" ") ||
		strings.Contains(lower, item.trigger+" ? ") ||
		strings.Contains(lower, item.trigger+" ?") {
		return true
	}

	// 5. Logical AND guard: navigator.connection && navigator.connection.effectiveType
	if strings.Contains(lower, item.trigger+" &&") || strings.Contains(lower, item.trigger+"&&") {
		return true
	}

	return false
}
