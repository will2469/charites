package browser

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// NonPassiveScrollListenerRule menegakkan penggunaan opsi { passive: true } pada event listener sentuh dan scroll.
type NonPassiveScrollListenerRule struct{}

// NewNonPassiveScrollListenerRule membuat instance baru dari NonPassiveScrollListenerRule.
func NewNonPassiveScrollListenerRule() *NonPassiveScrollListenerRule {
	return &NonPassiveScrollListenerRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *NonPassiveScrollListenerRule) ID() string {
	return "browser.non-passive-scroll-listener"
}

// Description mengembalikan ringkasan aturan.
func (r *NonPassiveScrollListenerRule) Description() string {
	return "Enforces { passive: true } option on touch and wheel event listeners to prevent main thread scroll blocking"
}

// Category mengembalikan nama kategori rule.
func (r *NonPassiveScrollListenerRule) Category() string {
	return "browser"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *NonPassiveScrollListenerRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *NonPassiveScrollListenerRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C DOM Level 4 Events Specification (Passive Event Listeners)",
			"Chromium & WebKit Compositor Scrolling Pipeline Guidelines",
			"Google Lighthouse Best Practices (Does not use passive listeners)",
		},
		CoreInvariant: "Event listeners for 'touchstart', 'touchmove', 'wheel', or 'mousewheel' must declare '{ passive: true }' to ensure non-blocking compositor scrolling.",
		Grounding: "Browsers execute smooth scrolling on a dedicated compositor thread. Without '{ passive: true }', " +
			"the compositor must block and wait for JavaScript execution on the main thread to see if 'preventDefault()' is called.\n\n" +
			"This introduces severe touch response latency and frame rate drops (scroll jank) on mobile devices.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Mobile Scroll Jank & Latency",
				Severity: "MEDIUM",
				Impact:   "Users experience jerky, lagging scrolling and delayed touch gestures on Safari iOS and Android Chrome.",
			},
			{
				Vector:   "Lighthouse Performance Penalty",
				Severity: "LOW",
				Impact:   "Fails Lighthouse 'Does not use passive listeners to improve scrolling performance' audit.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "javascript",
				Comment:  "Adding touchmove listener without passive: true option",
				Code: `window.addEventListener("touchmove", (e) => {
  trackTouchPosition(e);
});`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "javascript",
				Comment:  "Specifying { passive: true } to unblock the compositor thread",
				Code: `window.addEventListener("touchmove", (e) => {
  trackTouchPosition(e);
}, { passive: true });`,
			},
		},
	}
}

var targetScrollEvents = [...]string{
	"\"touchstart\"", "'touchstart'", "`touchstart`",
	"\"touchmove\"", "'touchmove'", "`touchmove`",
	"\"wheel\"", "'wheel'", "`wheel`",
	"\"mousewheel\"", "'mousewheel'", "`mousewheel`",
}

// Evaluate memeriksa apakah addEventListener dipanggil pada event sentuh/wheel tanpa opsi passive.
func (r *NonPassiveScrollListenerRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil {
		return nil
	}

	// 1. Cek pada elemen <script>
	if strings.ToLower(node.Tag) == "script" {
		scriptText := getStyleNodeText(node)
		return r.evalScript(node, scriptText)
	}

	// 2. Cek pada atribut inline yang memuat addEventListener
	if node.Attributes != nil {
		for _, attrVal := range node.Attributes {
			if strings.Contains(strings.ToLower(attrVal), "addeventlistener") {
				if diags := r.evalScript(node, attrVal); len(diags) > 0 {
					return diags
				}
			}
		}
	}

	return nil
}

func (r *NonPassiveScrollListenerRule) evalScript(node *ir.Node, script string) []ir.Diagnostic {
	lower := strings.ToLower(script)
	if !strings.Contains(lower, "addeventlistener") {
		return nil
	}

	hasTargetEvent := false
	for _, evt := range targetScrollEvents {
		if strings.Contains(lower, evt) {
			hasTargetEvent = true
			break
		}
	}
	if !hasTargetEvent {
		return nil
	}

	var diags []ir.Diagnostic
	idx := 0

	for {
		match := strings.Index(lower[idx:], "addeventlistener")
		if match == -1 {
			break
		}
		callStart := idx + match

		// Cari batas akhir pemanggilan addEventListener dengan balanced parenthesis tracking
		callEnd := findCallEnd(lower, callStart)
		callText := lower[callStart:callEnd]

		for _, evt := range targetScrollEvents {
			if strings.Contains(callText, evt) {
				// Cek apakah memiliki opsi passive: true/false atau memanggil preventDefault
				if strings.Contains(callText, "passive:") || strings.Contains(callText, "passive :") || strings.Contains(callText, "preventdefault") {
					continue
				}

				lineOffset := strings.Count(script[:callStart], "\n")
				diags = append(diags, ir.Diagnostic{
					Line:     node.Span.Line + lineOffset,
					Column:   node.Span.Column,
					Rule:     r.ID(),
					Severity: r.DefaultSeverity(),
					Message:  "Scroll-blocking event listener on " + evt + " lacks '{ passive: true }'. This delays compositor thread scrolling and causes stutter on touch devices.",
					Hint:     "Add '{ passive: true }' as the third argument to addEventListener to guarantee smooth 60fps scrolling.",
				})
			}
		}

		idx = callStart + 16
	}

	return diags
}

// findCallEnd mencari batas penutup ')' dari pemanggilan fungsi dengan menghitung kedalaman tanda kurung.
func findCallEnd(s string, startPos int) int {
	parenIdx := strings.IndexByte(s[startPos:], '(')
	if parenIdx == -1 {
		return len(s)
	}

	depth := 1
	i := startPos + parenIdx + 1
	inString := byte(0)

	for i < len(s) {
		ch := s[i]
		if inString != 0 {
			if ch == inString && (i == 0 || s[i-1] != '\\') {
				inString = 0
			}
			i++
			continue
		}

		if ch == '"' || ch == '\'' || ch == '`' {
			inString = ch
			i++
			continue
		}

		switch ch {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i + 1
			}
		}
		i++
	}

	return len(s)
}
