package browser

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// FirefoxOnlyAPIRule mendeteksi penggunaan properti dan method legacy berprefix Gecko tanpa alternatif standar W3C.
type FirefoxOnlyAPIRule struct{}

// NewFirefoxOnlyAPIRule membuat instance baru dari FirefoxOnlyAPIRule.
func NewFirefoxOnlyAPIRule() *FirefoxOnlyAPIRule {
	return &FirefoxOnlyAPIRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *FirefoxOnlyAPIRule) ID() string {
	return "browser.firefox-only-api"
}

// Description mengembalikan ringkasan aturan.
func (r *FirefoxOnlyAPIRule) Description() string {
	return "Flags usage of legacy Gecko/Firefox-exclusive DOM extensions and APIs without standard W3C equivalents"
}

// Category mengembalikan nama kategori rule.
func (r *FirefoxOnlyAPIRule) Category() string {
	return "browser"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *FirefoxOnlyAPIRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *FirefoxOnlyAPIRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Fullscreen API Specification",
			"W3C DOM Standards (Gecko Extension Deprecations)",
			"MDN Web Docs (Gecko-Specific DOM Interfaces)",
		},
		CoreInvariant: "Gecko-prefixed DOM methods and proprietary APIs ('mozRequestFullScreen', 'InstallTrigger', etc.) must provide standard W3C equivalents for Blink and WebKit.",
		Grounding: "Mozilla Firefox historically exposed vendor-prefixed APIs such as 'mozRequestFullScreen' and browser-specific globals like 'InstallTrigger'.\n\n" +
			"Calling these directly without standard W3C methods causes instant crashes or undefined behavior in Blink (Chrome/Edge) and WebKit (Safari).",
		Risks: []ir.RiskItem{
			{
				Vector:   "Crash in Chrome and Safari",
				Severity: "MEDIUM",
				Impact:   "Invoking 'element.mozRequestFullScreen()' throws 'TypeError: element.mozRequestFullScreen is not a function' in all non-Gecko browsers.",
			},
			{
				Vector:   "Obsolete Browser Sniffing",
				Severity: "LOW",
				Impact:   "Relying on 'InstallTrigger' to detect Firefox is deprecated and breaks as Firefox modernizes its engine.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "javascript",
				Comment:  "Direct invocation of mozRequestFullScreen without standard check",
				Code: `function enterFullscreen(element) {
  element.mozRequestFullScreen();
}`,
			},
			{
				Language: "javascript",
				Comment:  "Direct access to Gecko-specific inner screen property",
				Code:     `const screenX = window.mozInnerScreenX;`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "javascript",
				Comment:  "Prioritizing standard W3C fullscreen method",
				Code: `function enterFullscreen(element) {
  if (element.requestFullscreen) {
    element.requestFullscreen();
  } else if (element.mozRequestFullScreen) {
    element.mozRequestFullScreen();
  }
}`,
			},
			{
				Language: "javascript",
				Comment:  "Standard fullscreenElement fallback chain",
				Code:     `const fsElement = document.fullscreenElement || document.mozFullScreenElement;`,
			},
		},
	}
}

type firefoxAPIItem struct {
	trigger     string
	name        string
	standardAPI string
}

var firefoxAPIs = [...]firefoxAPIItem{
	{trigger: "mozrequestfullscreen", name: "mozRequestFullScreen", standardAPI: "requestFullscreen"},
	{trigger: "mozcancelfullscreen", name: "mozCancelFullScreen", standardAPI: "exitFullscreen"},
	{trigger: "mozfullscreenelement", name: "mozFullScreenElement", standardAPI: "fullscreenElement"},
	{trigger: "mozfullscreenenabled", name: "mozFullScreenEnabled", standardAPI: "fullscreenEnabled"},
	{trigger: "installtrigger", name: "InstallTrigger", standardAPI: "Feature Detection"},
	{trigger: "window.sidebar", name: "window.sidebar", standardAPI: "Standard Bookmark API"},
	{trigger: "navigator.buildid", name: "navigator.buildID", standardAPI: "Feature Detection"},
	{trigger: "mozinnerscreenx", name: "mozInnerScreenX", standardAPI: "screenX"},
	{trigger: "mozinnerscreeny", name: "mozInnerScreenY", standardAPI: "screenY"},
}

// Evaluate memeriksa apakah kode script atau event handler memanggil Firefox/Gecko API tanpa fallback.
func (r *FirefoxOnlyAPIRule) Evaluate(node *ir.Node) []ir.Diagnostic {
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

func (r *FirefoxOnlyAPIRule) evaluateScriptContent(node *ir.Node, script string, isScriptBlock bool) []ir.Diagnostic {
	lower := strings.ToLower(script)
	if lower == "" || !hasAnyFirefoxAPI(lower) {
		return nil
	}

	if strings.Contains(lower, "try {") && strings.Contains(lower, "} catch") {
		return nil
	}

	//nolint:prealloc // zero-alloc on clean nodes required by QUAL-03
	var diags []ir.Diagnostic

	for _, item := range firefoxAPIs {
		if !strings.Contains(lower, item.trigger) {
			continue
		}

		if isFirefoxGuarded(lower, item) {
			continue
		}

		line := resolveTriggerLine(node, script, item.trigger, isScriptBlock)

		diags = append(diags, ir.Diagnostic{
			Line:     line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Usage of legacy Gecko/Firefox-exclusive API '" + item.name + "'. Unsupported on Chrome, Safari, and standard environments.",
			Hint:     "Provide standard W3C equivalent (e.g. '" + item.standardAPI + " || " + item.name + "').",
		})
	}

	return diags
}

func hasAnyFirefoxAPI(lower string) bool {
	for i := range firefoxAPIs {
		if strings.Contains(lower, firefoxAPIs[i].trigger) {
			return true
		}
	}
	return false
}

func isFirefoxGuarded(lower string, item firefoxAPIItem) bool {
	// 1. Optional chaining
	if strings.Contains(lower, item.trigger+"?.") {
		return true
	}

	// 2. Guard 'in' operator: e.g. 'mozRequestFullScreen' in element
	if strings.Contains(lower, "'"+item.trigger+"' in") ||
		strings.Contains(lower, "\""+item.trigger+"\" in") {
		return true
	}

	// 3. Fallback chain dengan operator || atau ?? bersamaan dengan standard API
	stdLower := strings.ToLower(item.standardAPI)
	if (strings.Contains(lower, "||") || strings.Contains(lower, "??")) && strings.Contains(lower, stdLower) {
		return true
	}

	// 4. Pengujian if eksplisit: if (element.requestFullscreen) ... else if (element.mozRequestFullScreen)
	if strings.Contains(lower, "if") && strings.Contains(lower, stdLower) && strings.Contains(lower, item.trigger) {
		return true
	}

	return false
}
