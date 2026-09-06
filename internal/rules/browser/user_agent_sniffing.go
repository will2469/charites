package browser

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// UserAgentSniffingRule mendeteksi percabangan logika aplikasi berbasis string navigator.userAgent.
type UserAgentSniffingRule struct{}

// NewUserAgentSniffingRule membuat instance baru dari UserAgentSniffingRule.
func NewUserAgentSniffingRule() *UserAgentSniffingRule {
	return &UserAgentSniffingRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *UserAgentSniffingRule) ID() string {
	return "browser.user-agent-sniffing"
}

// Description mengembalikan ringkasan aturan.
func (r *UserAgentSniffingRule) Description() string {
	return "Flags conditional branching based on navigator.userAgent string sniffing and enforces W3C capability/feature detection"
}

// Category mengembalikan nama kategori rule.
func (r *UserAgentSniffingRule) Category() string {
	return "browser"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *UserAgentSniffingRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *UserAgentSniffingRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C HTML Design Principles (Avoid Browser Sniffing)",
			"Chromium Client Hints & User-Agent Reduction Guidelines",
			"MDN Web Docs (Browser Detection Using the User Agent - Best Practices)",
		},
		CoreInvariant: "Application logic and responsive branching must not rely on substring or regex matching of 'navigator.userAgent'. Use W3C capability detection instead.",
		Grounding: "User-Agent strings are historically fragile, frequently spoofed, and currently frozen across major browsers (Chrome, Safari, Edge).\n\n" +
			"For example, Chrome contains 'Safari' and 'WebKit', Edge contains 'Chrome', and iPadOS reports as macOS 'Macintosh'. " +
			"Branching on User-Agent strings leads to silent feature failures and broken responsive layouts on newer devices.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Frozen & Spoofed UA Strings",
				Severity: "MEDIUM",
				Impact:   "Browsers freeze version numbers or disguise platform tokens, causing browser sniffing logic to misclassify modern mobile devices as desktop or vice versa.",
			},
			{
				Vector:   "Cross-Browser Engine Breakage",
				Severity: "MEDIUM",
				Impact:   "Alternative browsers (Brave, Vivaldi, Arc, Firefox Focus) or tablets (iPadOS) receive crippled mobile views or desktop-only controls that cannot be touched.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "javascript",
				Comment:  "Branching layout or feature based on navigator.userAgent regex",
				Code: `if (/android|iphone|ipad/i.test(navigator.userAgent)) {
  initMobileLayout();
}`,
			},
			{
				Language: "typescript",
				Comment:  "Checking browser brand via userAgent.includes",
				Code: `if (navigator.userAgent.includes("Chrome")) {
  enableChromeFeature();
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "javascript",
				Comment:  "Using W3C CSS Media Queries for pointer capability detection",
				Code: `if (window.matchMedia("(pointer: coarse)").matches) {
  initMobileLayout();
}`,
			},
			{
				Language: "typescript",
				Comment:  "Using feature detection instead of browser sniffing",
				Code: `if ("visualViewport" in window) {
  enableViewportFeature();
}`,
			},
			{
				Language: "typescript",
				Comment:  "Telemetry logging is allowed and not flagged",
				Code: `logger.sendMetrics({
  userAgent: navigator.userAgent,
  timestamp: Date.now(),
});`,
			},
		},
	}
}

// Evaluate memeriksa apakah kode script atau event handler melakukan user-agent sniffing.
func (r *UserAgentSniffingRule) Evaluate(node *ir.Node) []ir.Diagnostic {
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

func (r *UserAgentSniffingRule) evaluateScriptContent(node *ir.Node, script string, isScriptBlock bool) []ir.Diagnostic {
	lower := strings.ToLower(script)
	if lower == "" || (!strings.Contains(lower, "useragent") && !strings.Contains(lower, "useragentdata")) {
		return nil
	}

	if !hasUASniffing(lower) {
		return nil
	}

	line := node.Span.Line
	if isScriptBlock {
		line = findUASniffingLine(node, script)
	}

	return []ir.Diagnostic{
		{
			Line:     line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Fragile browser sniffing via navigator.userAgent. User-Agent strings are spoofed, frozen, and inaccurate across devices. Use W3C capability detection ('feature' in window, CSS.supports, or matchMedia) instead.",
			Hint:     "Replace navigator.userAgent string inspection with capability detection (e.g. 'window.matchMedia(\"(pointer: coarse)\").matches' or '\"visualViewport\" in window').",
		},
	}
}

var uaBrandKeywords = [...]string{
	"chrome", "safari", "firefox", "opera", "edge", "android", "iphone", "ipad", "mobile", "macintosh", "windows", "msie", "trident",
}

func hasUASniffing(lower string) bool {
	// 1. Method pemanggilan langsung pada navigator.userAgent
	if strings.Contains(lower, "useragent.includes(") ||
		strings.Contains(lower, "useragent.indexof(") ||
		strings.Contains(lower, "useragent.match(") ||
		strings.Contains(lower, "useragent.search(") {
		return true
	}

	// 2. Pengujian regex: .test(navigator.userAgent) atau .test(window.navigator.userAgent)
	if strings.Contains(lower, ".test(navigator.useragent") ||
		strings.Contains(lower, ".test(window.navigator.useragent") {
		return true
	}

	// 3. Substring pemeriksaan brand jika terdapat variabel ua dan pengujian brand
	if strings.Contains(lower, "useragent") &&
		(strings.Contains(lower, ".includes(") || strings.Contains(lower, ".indexof(") || strings.Contains(lower, ".test(")) {
		for _, brand := range uaBrandKeywords {
			if strings.Contains(lower, brand) {
				return true
			}
		}
	}

	return false
}

func findUASniffingLine(node *ir.Node, script string) int {
	lines := strings.Split(script, "\n")
	for idx, l := range lines {
		lowerLine := strings.ToLower(l)
		if hasUASniffing(lowerLine) || (strings.Contains(lowerLine, "useragent") && (strings.Contains(lowerLine, "if ") || strings.Contains(lowerLine, "match") || strings.Contains(lowerLine, "test"))) {
			return node.Span.Line + idx
		}
	}
	return node.Span.Line
}
