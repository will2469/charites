package browser

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// WebKitOnlyAPIRule mendeteksi pemanggilan API legacy berprefix WebKit tanpa alternatif standar W3C.
type WebKitOnlyAPIRule struct{}

// NewWebKitOnlyAPIRule membuat instance baru dari WebKitOnlyAPIRule.
func NewWebKitOnlyAPIRule() *WebKitOnlyAPIRule {
	return &WebKitOnlyAPIRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *WebKitOnlyAPIRule) ID() string {
	return "browser.webkit-only-api"
}

// Description mengembalikan ringkasan aturan.
func (r *WebKitOnlyAPIRule) Description() string {
	return "Flags direct invocation of WebKit-prefixed legacy APIs without standard W3C equivalents or graceful fallbacks"
}

// Category mengembalikan nama kategori rule.
func (r *WebKitOnlyAPIRule) Category() string {
	return "browser"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *WebKitOnlyAPIRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *WebKitOnlyAPIRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Fullscreen API Specification",
			"W3C Web Audio API Specification",
			"W3C Web Speech API Specification",
		},
		CoreInvariant: "Direct invocation of WebKit-prefixed methods ('webkitRequestFullscreen', 'webkitAudioContext', etc.) must provide standard W3C fallbacks for non-WebKit browsers.",
		Grounding: "WebKit vendor-prefixed APIs were transitional features during early HTML5 standardization.\n\n" +
			"Calling them directly without checking W3C standard methods causes instant runtime crashes in Mozilla Firefox desktop and modern Chromium environments.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Crash in Firefox and Non-WebKit Browsers",
				Severity: "MEDIUM",
				Impact:   "Methods like 'element.webkitRequestFullscreen()' throw 'TypeError: not a function' in Firefox because Gecko does not implement WebKit prefixes.",
			},
			{
				Vector:   "Missed W3C Standard Improvements",
				Severity: "LOW",
				Impact:   "Legacy WebKit methods lack updated return types (such as Promises) supported by modern standard W3C APIs.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "javascript",
				Comment:  "Direct invocation of WebKit fullscreen without standard check",
				Code: `function enterFullscreen(element) {
  element.webkitRequestFullscreen();
}`,
			},
			{
				Language: "javascript",
				Comment:  "Direct instantiation of webkitAudioContext",
				Code:     `const audioCtx = new webkitAudioContext();`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "javascript",
				Comment:  "Prioritizing standard W3C method with WebKit fallback",
				Code: `function enterFullscreen(element) {
  if (element.requestFullscreen) {
    element.requestFullscreen();
  } else if (element.webkitRequestFullscreen) {
    element.webkitRequestFullscreen();
  }
}`,
			},
			{
				Language: "javascript",
				Comment:  "Standard AudioContext fallback chain",
				Code: `const AudioContextClass = window.AudioContext || window.webkitAudioContext;
const audioCtx = new AudioContextClass();`,
			},
		},
	}
}

type webkitAPIItem struct {
	trigger     string
	name        string
	standardAPI string
}

var webkitAPIs = [...]webkitAPIItem{
	{trigger: "webkitrequestfullscreen", name: "webkitRequestFullscreen", standardAPI: "requestFullscreen"},
	{trigger: "webkitexitfullscreen", name: "webkitExitFullscreen", standardAPI: "exitFullscreen"},
	{trigger: "webkitcancelfullscreen", name: "webkitCancelFullScreen", standardAPI: "exitFullscreen"},
	{trigger: "webkitaudiocontext", name: "webkitAudioContext", standardAPI: "AudioContext"},
	{trigger: "webkitspeechrecognition", name: "webkitSpeechRecognition", standardAPI: "SpeechRecognition"},
	{trigger: "webkitrequestanimationframe", name: "webkitRequestAnimationFrame", standardAPI: "requestAnimationFrame"},
	{trigger: "webkitcancelanimationframe", name: "webkitCancelAnimationFrame", standardAPI: "cancelAnimationFrame"},
	{trigger: "webkiturl", name: "webkitURL", standardAPI: "URL"},
	{trigger: "webkitindexeddb", name: "webkitIndexedDB", standardAPI: "indexedDB"},
}

// Evaluate memeriksa apakah kode script atau event handler memanggil WebKit API tanpa fallback.
func (r *WebKitOnlyAPIRule) Evaluate(node *ir.Node) []ir.Diagnostic {
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

func (r *WebKitOnlyAPIRule) evaluateScriptContent(node *ir.Node, script string, isScriptBlock bool) []ir.Diagnostic {
	lower := strings.ToLower(script)
	if lower == "" || !strings.Contains(lower, "webkit") {
		return nil
	}

	if strings.Contains(lower, "try {") && strings.Contains(lower, "} catch") {
		return nil
	}

	//nolint:prealloc // zero-alloc on clean nodes required by QUAL-03
	var diags []ir.Diagnostic

	for _, item := range webkitAPIs {
		if !strings.Contains(lower, item.trigger) {
			continue
		}

		if isWebKitGuarded(lower, item) {
			continue
		}

		line := resolveTriggerLine(node, script, item.trigger, isScriptBlock)

		diags = append(diags, ir.Diagnostic{
			Line:     line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Direct invocation of legacy WebKit-prefixed API '" + item.name + "'. Unsupported on Firefox and non-WebKit standard browsers.",
			Hint:     "Provide standard W3C fallback (e.g. '" + item.standardAPI + " || " + item.name + "').",
		})
	}

	return diags
}

func isWebKitGuarded(lower string, item webkitAPIItem) bool {
	// 1. Optional chaining
	if strings.Contains(lower, item.trigger+"?.") {
		return true
	}

	// 2. Guard 'in' operator: e.g. 'webkitRequestFullscreen' in element
	if strings.Contains(lower, "'"+item.trigger+"' in") ||
		strings.Contains(lower, "\""+item.trigger+"\" in") {
		return true
	}

	// 3. Fallback chain dengan operator || atau ?? bersamaan dengan standard API
	stdLower := strings.ToLower(item.standardAPI)
	if (strings.Contains(lower, "||") || strings.Contains(lower, "??")) && strings.Contains(lower, stdLower) {
		return true
	}

	// 4. Pengujian if eksplisit: if (element.requestFullscreen) ... else if (element.webkitRequestFullscreen)
	if strings.Contains(lower, "if") && strings.Contains(lower, stdLower) && strings.Contains(lower, item.trigger) {
		return true
	}

	return false
}
