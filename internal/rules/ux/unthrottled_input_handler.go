package ux

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// UnthrottledInputHandlerRule mencegah pemanggilan network API / fetch langsung pada setiap keystroke input teks
// tanpa perlindungan debounce atau throttle, guna menjaga stabilitas perseptual dan mencegah layout jitter.
type UnthrottledInputHandlerRule struct{}

// NewUnthrottledInputHandlerRule membuat instance baru dari UnthrottledInputHandlerRule.
func NewUnthrottledInputHandlerRule() *UnthrottledInputHandlerRule {
	return &UnthrottledInputHandlerRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *UnthrottledInputHandlerRule) ID() string {
	return "ux.unthrottled-input-handler"
}

// Description mengembalikan ringkasan aturan.
func (r *UnthrottledInputHandlerRule) Description() string {
	return "Flags text input handlers that trigger unthrottled network calls directly on keystrokes"
}

// Category mengembalikan nama kategori rule.
func (r *UnthrottledInputHandlerRule) Category() string {
	return "ux"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *UnthrottledInputHandlerRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *UnthrottledInputHandlerRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Perceptual Stability & Doherty Threshold (< 400ms)",
			"Nielsen Norman Group: Response Times (The 3 Important Limits)",
			"WCAG 2.2 Success Criterion 2.2.4 (Interruptions)",
		},
		CoreInvariant: "Text input handlers ('onChange', 'onInput') must not trigger direct network requests without debounce or throttle protection.",
		Grounding: "Firing network requests on every single keystroke floods the network with redundant in-flight calls, " +
			"causes race conditions where earlier responses overwrite newer ones (out-of-order responses), " +
			"and produces aggressive layout thrashing / UI jitter as suggestion dropdowns flicker erratically.\n\n" +
			"Wrapping handlers in a 250-400ms debounce buffer (or throttle) stabilizes perceptual performance, " +
			"dramatically reduces server load, and guarantees that search results correspond to the user's finalized query.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Out-of-Order Race Conditions & Stale UI",
				Severity: "HIGH",
				Impact:   "Slow earlier network responses resolve after fast later responses, showing stale search results for an old keystroke.",
			},
			{
				Vector:   "UI Jitter & Layout Thrashing",
				Severity: "MEDIUM",
				Impact:   "Rapid re-rendering of dropdown popovers on each keystroke causes visual stutter and jarring jumps.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Direct unthrottled fetch call inside onChange input handler",
				Code: `<div className="relative">
  <input
    type="search"
    placeholder="Cari produk..."
    onChange={e => fetchSuggestions(e.target.value)}
    className="w-full px-4 py-2 border rounded"
  />
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Debounced handler buffer (300ms) prior to triggering network search",
				Code: `const debouncedSearch = useDebouncedCallback((query: string) => {
  fetchSuggestions(query);
}, 300);

<div className="relative">
  <input
    type="search"
    placeholder="Cari produk..."
    onChange={e => debouncedSearch(e.target.value)}
    className="w-full px-4 py-2 border rounded"
  />
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah input teks memiliki event handler yang memicu network call tanpa throttle/debounce.
func (r *UnthrottledInputHandlerRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !isTextInputElement(node) {
		return nil
	}

	for attrName, attrVal := range node.Attributes {
		if !isKeystrokeHandlerAttr(attrName) {
			continue
		}

		if callName, ok := detectUnthrottledNetworkCall(attrVal); ok {
			return []ir.Diagnostic{
				{
					Line:     node.Span.Line,
					Column:   node.Span.Column,
					Rule:     r.ID(),
					Severity: r.DefaultSeverity(),
					Message: fmt.Sprintf(
						"Input handler %q triggers unthrottled network call %q directly on keystroke.",
						attrName,
						callName,
					),
					Hint: "Wrap the network request in a debounce or throttle function (e.g. useDebouncedCallback(..., 300)) to prevent layout jitter and out-of-order race conditions.",
				},
			}
		}
	}

	return nil
}

func isTextInputElement(node *ir.Node) bool {
	tagLower := strings.ToLower(node.Tag)
	switch tagLower {
	case "textarea", "searchinput", "textfield", "searchbar", "searchbox":
		return true
	case "input":
		if t, ok := getAttrCaseInsensitive(node, "type"); ok {
			typ := cleanAttrValue(t)
			switch typ {
			case "checkbox", "radio", "file", "button", "submit", "reset", "color", "range", "hidden":
				return false
			}
		}
		return true
	}

	if val, ok := getAttrCaseInsensitive(node, "contenteditable"); ok {
		if cleanAttrValue(val) == "true" {
			return true
		}
	}

	return false
}

func isKeystrokeHandlerAttr(name string) bool {
	lower := strings.ToLower(name)
	switch lower {
	case "onchange", "oninput", "onkeyup", "onkeydown", "onkeypress":
		return true
	}
	return false
}
