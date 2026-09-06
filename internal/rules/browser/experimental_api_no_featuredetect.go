package browser

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// ExperimentalAPINoFeaturedetectRule mendeteksi pemanggilan Web API eksperimental tanpa feature detection guard.
type ExperimentalAPINoFeaturedetectRule struct{}

// NewExperimentalAPINoFeaturedetectRule membuat instance baru dari ExperimentalAPINoFeaturedetectRule.
func NewExperimentalAPINoFeaturedetectRule() *ExperimentalAPINoFeaturedetectRule {
	return &ExperimentalAPINoFeaturedetectRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *ExperimentalAPINoFeaturedetectRule) ID() string {
	return "browser.experimental-api-no-featuredetect"
}

// Description mengembalikan ringkasan aturan.
func (r *ExperimentalAPINoFeaturedetectRule) Description() string {
	return "Detects invocation of experimental Web APIs without runtime feature detection guards"
}

// Category mengembalikan nama kategori rule.
func (r *ExperimentalAPINoFeaturedetectRule) Category() string {
	return "browser"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *ExperimentalAPINoFeaturedetectRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ExperimentalAPINoFeaturedetectRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"WICG Web Share API Specification (navigator.share)",
			"WICG File System Access API (showOpenFilePicker)",
			"W3C CSS View Transitions Module Level 1 (startViewTransition)",
			"ECMA-262 Feature Detection & Defensive JavaScript Guidelines",
		},
		CoreInvariant: "Experimental or non-universal Web APIs must be guarded with runtime capability checks ('prop' in obj, if (obj.prop), optional chaining, or try/catch) before invocation.",
		Grounding: "Modern Web APIs are adopted unevenly across browser vendors. For instance, 'showOpenFilePicker' is exclusive to Chromium and crashes instantly on Firefox and Safari.\n\n" +
			"'navigator.share' throws an uncaught TypeError on desktop Firefox or non-secure contexts.\n\n" +
			"Directly calling these APIs without feature guards results in severe runtime exceptions that crash SPAs.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Runtime Exception Crash",
				Severity: "HIGH",
				Impact:   "Uncaught TypeError: undefined is not a function immediately terminates JavaScript execution in non-supporting browsers.",
			},
			{
				Vector:   "Broken Core User Action",
				Severity: "HIGH",
				Impact:   "Primary user actions like document sharing or file importing fail completely without informative fallback UX.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Direct invocation of navigator.share without runtime feature guard (crashes on desktop Firefox)",
				Code: `<button onClick={() => {
  navigator.share({ title: "Surat", url: window.location.href });
}}>
  Bagikan
</button>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Defensive feature detection with fallback to clipboard copy",
				Code: `<button onClick={async () => {
  if (typeof navigator !== "undefined" && navigator.share) {
    await navigator.share({ title: "Surat", url: window.location.href });
  } else {
    await navigator.clipboard?.writeText(window.location.href);
  }
}}>
  Bagikan
</button>`,
			},
		},
	}
}

// Evaluate memeriksa apakah kode script atau event handler memanggil API eksperimental tanpa guard.
func (r *ExperimentalAPINoFeaturedetectRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil {
		return nil
	}

	// 1. Cek atribut event handler pada elemen JSX / HTML (onClick, onChange, dangerouslySetInnerHTML)
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

func isScriptAttribute(name string) bool {
	lower := strings.ToLower(name)
	return strings.HasPrefix(lower, "on") || lower == "dangerouslysetinnerhtml"
}

type experimentalAPI struct {
	trigger  string
	name     string
	hintText string
}

var experimentalAPIs = [...]experimentalAPI{
	{trigger: "navigator.share", name: "navigator.share()", hintText: "Check 'if (navigator.share)' or 'if (\"share\" in navigator)' before calling."},
	{trigger: "showopenfilepicker", name: "window.showOpenFilePicker()", hintText: "Check 'if (window.showOpenFilePicker)' and provide an <input type=\"file\"> fallback for Firefox & Safari."},
	{trigger: "showsavefilepicker", name: "window.showSaveFilePicker()", hintText: "Check 'if (window.showSaveFilePicker)' before calling."},
	{trigger: "showdirectorypicker", name: "window.showDirectoryPicker()", hintText: "Check 'if (window.showDirectoryPicker)' before calling."},
	{trigger: "startviewtransition", name: "document.startViewTransition()", hintText: "Check 'if (document.startViewTransition)' or use optional chaining 'document.startViewTransition?.()' before invoking."},
	{trigger: "eyedropper", name: "new EyeDropper()", hintText: "Check 'if (\"EyeDropper\" in window)' before instantiating."},
	{trigger: "virtualkeyboard", name: "navigator.virtualKeyboard", hintText: "Check 'if (\"virtualKeyboard\" in navigator)' before using."},
	{trigger: "navigator.bluetooth", name: "navigator.bluetooth", hintText: "Check 'if (\"bluetooth\" in navigator)' before accessing."},
	{trigger: "navigator.serial", name: "navigator.serial", hintText: "Check 'if (\"serial\" in navigator)' before accessing."},
	{trigger: "navigator.usb", name: "navigator.usb", hintText: "Check 'if (\"usb\" in navigator)' before accessing."},
}

func hasAnyExperimentalAPI(lower string) bool {
	for i := range experimentalAPIs {
		if strings.Contains(lower, experimentalAPIs[i].trigger) {
			return true
		}
	}
	return false
}

func resolveTriggerLine(node *ir.Node, script, trigger string, isScriptBlock bool) int {
	if !isScriptBlock {
		return node.Span.Line
	}
	lines := strings.Split(script, "\n")
	for lineIdx, l := range lines {
		if strings.Contains(strings.ToLower(l), trigger) {
			return node.Span.Line + lineIdx
		}
	}
	return node.Span.Line
}

func (r *ExperimentalAPINoFeaturedetectRule) evaluateScriptContent(node *ir.Node, script string, isScriptBlock bool) []ir.Diagnostic {
	lower := strings.ToLower(script)
	if lower == "" || !hasAnyExperimentalAPI(lower) {
		return nil
	}

	// Jika seluruh script dibungkus try/catch global, dianggap aman
	if strings.Contains(lower, "try {") && strings.Contains(lower, "} catch") {
		return nil
	}

	//nolint:prealloc // zero-alloc on clean nodes required by QUAL-03
	var diags []ir.Diagnostic

	for _, api := range experimentalAPIs {
		if !strings.Contains(lower, api.trigger) || isGuardedInvocation(lower, api.trigger) {
			continue
		}

		line := resolveTriggerLine(node, script, api.trigger, isScriptBlock)

		diags = append(diags, ir.Diagnostic{
			Line:     line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Direct invocation of experimental Web API '" + api.name + "' without runtime feature detection. Unsupported browsers (Firefox, Safari, or non-secure contexts) will crash with an uncaught TypeError.",
			Hint:     api.hintText,
		})
	}

	return diags
}

func isGuardedInvocation(scriptLower, trigger string) bool {
	// 1. Optional chaining: e.g. .share?. or startViewTransition?. or showOpenFilePicker?.
	if strings.Contains(scriptLower, trigger+"?.") ||
		strings.Contains(scriptLower, trigger+" ?.") {
		return true
	}

	shortName := trigger
	if idx := strings.LastIndex(trigger, "."); idx != -1 {
		shortName = trigger[idx+1:]
	}

	if strings.Contains(scriptLower, shortName+"?.") ||
		strings.Contains(scriptLower, shortName+" ?.") {
		return true
	}

	// 2. Feature detection via "in" operator: e.g. 'share' in navigator, "showOpenFilePicker" in window
	if strings.Contains(scriptLower, "'"+shortName+"' in") ||
		strings.Contains(scriptLower, "\""+shortName+"\" in") ||
		strings.Contains(scriptLower, "`"+shortName+"` in") {
		return true
	}

	// 3. Conditional guard with if, &&, or ternary:
	if strings.Contains(scriptLower, "if (") || strings.Contains(scriptLower, "if(") {
		if strings.Contains(scriptLower, "&& "+trigger) ||
			strings.Contains(scriptLower, "&& "+shortName) ||
			strings.Contains(scriptLower, trigger+" &&") ||
			strings.Contains(scriptLower, shortName+" &&") ||
			strings.Contains(scriptLower, "if ("+trigger) ||
			strings.Contains(scriptLower, "if ("+shortName) ||
			strings.Contains(scriptLower, "if (document."+shortName) ||
			strings.Contains(scriptLower, "if (window."+shortName) ||
			strings.Contains(scriptLower, "if (navigator."+shortName) ||
			strings.Contains(scriptLower, "if("+trigger) ||
			strings.Contains(scriptLower, "if("+shortName) ||
			strings.Contains(scriptLower, "if(document."+shortName) ||
			strings.Contains(scriptLower, "if(window."+shortName) ||
			strings.Contains(scriptLower, "if(navigator."+shortName) {
			return true
		}
	}

	// 4. Logical AND guard: navigator.share && navigator.share(...)
	if strings.Contains(scriptLower, trigger+" &&") ||
		strings.Contains(scriptLower, shortName+" &&") ||
		strings.Contains(scriptLower, "&& "+trigger) ||
		strings.Contains(scriptLower, "&& "+shortName) {
		return true
	}

	return false
}
