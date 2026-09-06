package browser

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// ScrollbarVendorIncompleteRule memastikan kustomisasi scrollbar dideklarasikan berpasangan
// antara pseudo-elemen WebKit (::-webkit-scrollbar) dan standar W3C (scrollbar-width, scrollbar-color).
type ScrollbarVendorIncompleteRule struct{}

// NewScrollbarVendorIncompleteRule membuat instance baru ScrollbarVendorIncompleteRule.
func NewScrollbarVendorIncompleteRule() *ScrollbarVendorIncompleteRule {
	return &ScrollbarVendorIncompleteRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat browser.scrollbar-vendor-incomplete.
func (r *ScrollbarVendorIncompleteRule) ID() string {
	return "browser.scrollbar-vendor-incomplete"
}

// Description mengembalikan deskripsi ringkas aturan.
func (r *ScrollbarVendorIncompleteRule) Description() string {
	return "Enforces bidirectional cross-engine scrollbar styling pairing between WebKit pseudo-elements and W3C standard properties"
}

// Category mengembalikan nama kategori rule.
func (r *ScrollbarVendorIncompleteRule) Category() string {
	return "browser"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *ScrollbarVendorIncompleteRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ScrollbarVendorIncompleteRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C CSS Scrollbars Styling Module Level 1 (scrollbar-width, scrollbar-color)",
			"WebKit Proprietary Scrollbar Styling Documentation",
			"MDN Browser Compatibility Matrix for Scrollbar Customization",
		},
		CoreInvariant: "Scrollbar styling declarations must be bidirectional: declaring '::-webkit-scrollbar*' requires declaring W3C standard 'scrollbar-width' / 'scrollbar-color', and vice-versa.",
		Grounding: "Historically, custom scrollbars were styled using WebKit pseudo-elements (::-webkit-scrollbar, ::-webkit-scrollbar-thumb, ::-webkit-scrollbar-track) in Chromium and Safari.\n\n" +
			"However, Gecko (Firefox) strictly enforces the W3C standard (scrollbar-width and scrollbar-color) and deliberately ignores ::-webkit-scrollbar.\n\n" +
			"When developers only write ::-webkit-scrollbar, the scrollbar appears customized in Chrome and Safari, but renders as an unstyled thick grey default scrollbar in Firefox, causing severe visual discordance on dark themes.\n\n" +
			"Charites enforces bidirectional cross-engine pairing, guaranteeing scrollbars render gracefully across Chrome, Firefox, and Safari.",
		BadExamples: []ir.CodeExample{
			{
				Language: "css",
				Comment:  "Declaring only WebKit pseudo-elements (leaves Firefox with unstyled default scrollbar)",
				Code: `.custom-scroll::-webkit-scrollbar {
  width: 6px;
}
.custom-scroll::-webkit-scrollbar-thumb {
  background: var(--muted-foreground);
  border-radius: 9999px;
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "css",
				Comment:  "Declaring both W3C standard properties and WebKit pseudo-elements",
				Code: `.custom-scroll {
  scrollbar-width: thin;
  scrollbar-color: var(--muted-foreground) transparent;
}
.custom-scroll::-webkit-scrollbar {
  width: 6px;
}
.custom-scroll::-webkit-scrollbar-thumb {
  background: var(--muted-foreground);
  border-radius: 9999px;
}`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Firefox Visual Degradation",
				Severity: "MEDIUM",
				Impact:   "Scrollbars appear as bright grey system widgets in Firefox on dark-mode web applications.",
			},
			{
				Vector:   "Layout Shift / Text Clipping",
				Severity: "LOW",
				Impact:   "Layout shift on Firefox when expecting a 6px thin scrollbar but getting a 17px default desktop scrollbar.",
			},
		},
	}
}

// Evaluate memeriksa apakah deklarasi scrollbar tidak lengkap lintas vendor.
func (r *ScrollbarVendorIncompleteRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil {
		return nil
	}

	// 1. Periksa blok style CSS (<style>)
	if node.Tag == "style" {
		content := strings.ToLower(getStyleNodeText(node))
		if content == "" {
			return nil
		}

		hasWebKit := strings.Contains(content, "::-webkit-scrollbar")
		hasW3C := hasW3CScrollbarProperty(content)

		if hasWebKit && !hasW3C {
			return []ir.Diagnostic{
				{
					Line:     node.Span.Line,
					Column:   node.Span.Column,
					Rule:     r.ID(),
					Severity: r.DefaultSeverity(),
					Message:  "Scrollbar customization uses '::-webkit-scrollbar' without W3C standard 'scrollbar-width' or 'scrollbar-color'. Custom scrollbars will be ignored and render as default unstyled widgets in Firefox.",
					Hint:     "Add 'scrollbar-width: thin;' and 'scrollbar-color: <thumb-color> <track-color>;' to the scrollable container rule.",
				},
			}
		}

		if hasW3C && !hasWebKit {
			return []ir.Diagnostic{
				{
					Line:     node.Span.Line,
					Column:   node.Span.Column,
					Rule:     r.ID(),
					Severity: r.DefaultSeverity(),
					Message:  "Scrollbar customization uses W3C standard 'scrollbar-width' or 'scrollbar-color' without '::-webkit-scrollbar' fallback. Older Chromium or Safari WebKit versions will not apply custom sizing.",
					Hint:     "Add '::-webkit-scrollbar' and '::-webkit-scrollbar-thumb' pseudo-elements for legacy WebKit compatibility.",
				},
			}
		}
	}

	// 2. Periksa atribut style inline pada elemen JSX / HTML
	if styleAttr, ok := node.Attributes["style"]; ok {
		lowerStyle := strings.ToLower(styleAttr)
		hasW3C := hasW3CScrollbarProperty(lowerStyle)
		hasWebKit := strings.Contains(lowerStyle, "webkitscrollbar") || strings.Contains(lowerStyle, "webkit-scrollbar")

		if hasW3C && !hasWebKit {
			return []ir.Diagnostic{
				{
					Line:     node.Span.Line,
					Column:   node.Span.Column,
					Rule:     r.ID(),
					Severity: r.DefaultSeverity(),
					Message:  "Inline style defines W3C 'scrollbar-width' or 'scrollbar-color' without WebKit scrollbar pairing.",
					Hint:     "Provide matching WebKit pseudo-element rules in CSS or use Tailwind scrollbar utilities.",
				},
			}
		}
	}

	return nil
}

func hasW3CScrollbarProperty(s string) bool {
	targets := [...]string{"scrollbar-width", "scrollbarwidth", "scrollbar-color", "scrollbarcolor"}
	for _, target := range targets {
		idx := 0
		for {
			pos := strings.Index(s[idx:], target)
			if pos == -1 {
				break
			}
			matchIdx := idx + pos
			// Karakter sebelum target tidak boleh '-' atau alphanumeric (mencegah match pada --custom-var atau my-property)
			validStart := true
			if matchIdx > 0 {
				prev := s[matchIdx-1]
				if prev == '-' || (prev >= 'a' && prev <= 'z') || (prev >= '0' && prev <= '9') {
					validStart = false
				}
			}
			if validStart {
				return true
			}
			idx = matchIdx + len(target)
		}
	}
	return false
}
