package browser

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// ObsoleteVendorPrefixRule mendeteksi penggunaan prefix vendor usang (-moz-border-radius, dll.)
// serta mendeteksi ketidaklengkapan triad pemotongan teks multi-baris (-webkit-line-clamp).
type ObsoleteVendorPrefixRule struct{}

// NewObsoleteVendorPrefixRule membuat instance baru ObsoleteVendorPrefixRule.
func NewObsoleteVendorPrefixRule() *ObsoleteVendorPrefixRule {
	return &ObsoleteVendorPrefixRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat browser.obsolete-vendor-prefix.
func (r *ObsoleteVendorPrefixRule) ID() string {
	return "browser.obsolete-vendor-prefix"
}

// Description mengembalikan deskripsi ringkas aturan.
func (r *ObsoleteVendorPrefixRule) Description() string {
	return "Detects obsolete CSS vendor prefixes and incomplete -webkit-line-clamp multi-line truncation triads"
}

// Category mengembalikan nama kategori rule.
func (r *ObsoleteVendorPrefixRule) Category() string {
	return "browser"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *ObsoleteVendorPrefixRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ObsoleteVendorPrefixRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C CSS Overflow Module Level 3 (line-clamp & WebKit Triad)",
			"W3C CSS Cascading and Inheritance Level 5 (Standard Property Baselines)",
			"MDN Obsolete and Deprecated Vendor Prefix Specifications",
		},
		CoreInvariant: "Obsolete vendor prefixes should be replaced with W3C standards, and '-webkit-line-clamp' multi-line truncation must include the complete mandatory triad ('display: -webkit-box', '-webkit-box-orient: vertical', and 'overflow: hidden').",
		Grounding: "Modern browser engines have supported standard properties like border-radius, box-shadow, and box-sizing without vendor prefixes for over a decade. Continuing to write dead prefixes (-moz-border-radius, -webkit-box-shadow) pollutes styles and degrades maintainability.\n\n" +
			"Furthermore, multi-line paragraph truncation using '-webkit-line-clamp' strictly requires a 3-part companion triad:\n" +
			"1. display: -webkit-box;\n" +
			"2. -webkit-box-orient: vertical;\n" +
			"3. overflow: hidden;\n\n" +
			"When developers only specify '-webkit-line-clamp: 2' in inline styles (e.g. style={{ WebkitLineClamp: 2 }}) without the triad, text truncation silently fails across all browser engines, causing text to overflow un-truncated.\n\n" +
			"Charites detects dead vendor prefixes and incomplete line-clamp triads, recommending clean Tailwind 'line-clamp-*' utilities.",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Incomplete line-clamp in inline style (fails to truncate silently in all browsers)",
				Code: `<p style={{ WebkitLineClamp: 2, overflow: "hidden" }} className="text-sm text-muted-foreground">
  Pengumuman pelayanan administrasi desa...
</p>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Using Tailwind line-clamp-2 which automatically compiles the complete cross-browser triad",
				Code: `<p className="line-clamp-2 text-sm text-muted-foreground">
  Pengumuman pelayanan administrasi desa...
</p>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Silent Text Truncation Failure",
				Severity: "MEDIUM",
				Impact:   "Multi-line paragraph truncation fails silently, overflowing cards and destroying dashboard layout consistency.",
			},
			{
				Vector:   "Dead Vendor Prefix Clutter",
				Severity: "LOW",
				Impact:   "Dead vendor prefixes clutter CSS output and trigger linter compatibility warnings.",
			},
		},
	}
}

// Evaluate memeriksa apakah terdapat prefix usang atau triad line-clamp yang tidak lengkap.
func (r *ObsoleteVendorPrefixRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil {
		return nil
	}

	var diags []ir.Diagnostic

	// 1. Cek pada atribut style inline pada elemen
	if styleAttr, ok := node.Attributes["style"]; ok {
		diags = append(diags, r.evalInlineStyle(node, styleAttr)...)
	}

	// 2. Cek pada blok style CSS (<style>)
	if node.Tag == "style" {
		diags = append(diags, r.evalStyleBlock(node)...)
	}

	return diags
}

func (r *ObsoleteVendorPrefixRule) evalInlineStyle(node *ir.Node, styleAttr string) []ir.Diagnostic {
	lowerStyle := strings.ToLower(styleAttr)
	//nolint:prealloc // zero-alloc on clean nodes required by QUAL-03
	var diags []ir.Diagnostic

	// Cek line-clamp triad
	if hasLineClamp(lowerStyle) && !hasCompleteLineClampTriad(lowerStyle) {
		diags = append(diags, ir.Diagnostic{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Incomplete line-clamp triad. '-webkit-line-clamp' requires 'display: -webkit-box' and '-webkit-box-orient: vertical' to function. Without them, text truncation fails silently across all browsers.",
			Hint:     "Use the standard Tailwind utility 'line-clamp-2' (or line-clamp-3) instead of raw inline styles.",
		})
	}

	// Cek dead vendor prefixes
	for _, deadPrefix := range findDeadPrefixes(lowerStyle) {
		diags = append(diags, ir.Diagnostic{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: ir.SeverityInfo,
			Message:  "Obsolete vendor prefix '" + deadPrefix + "' is deprecated. Modern browsers support standard CSS properties directly.",
			Hint:     "Replace with standard CSS property without vendor prefix or use Tailwind utility classes.",
		})
	}

	return diags
}

func (r *ObsoleteVendorPrefixRule) evalStyleBlock(node *ir.Node) []ir.Diagnostic {
	styleText := getStyleNodeText(node)
	lowerContent := strings.ToLower(styleText)
	if lowerContent == "" {
		return nil
	}

	var diags []ir.Diagnostic

	if strings.Contains(lowerContent, "-webkit-line-clamp") && !hasCompleteLineClampTriad(lowerContent) {
		diags = append(diags, ir.Diagnostic{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Incomplete '-webkit-line-clamp' triad in CSS block. Missing 'display: -webkit-box' or '-webkit-box-orient: vertical'.",
			Hint:     "Add both 'display: -webkit-box;' and '-webkit-box-orient: vertical;' to the rule declaration.",
		})
	}

	deadList := findDeadPrefixes(lowerContent)
	if len(deadList) > 0 {
		lines := strings.Split(styleText, "\n")
		for lineIdx, line := range lines {
			lineLower := strings.ToLower(line)
			for _, deadPrefix := range deadList {
				if strings.Contains(lineLower, deadPrefix) {
					diags = append(diags, ir.Diagnostic{
						Line:     node.Span.Line + lineIdx,
						Column:   node.Span.Column,
						Rule:     r.ID(),
						Severity: ir.SeverityInfo,
						Message:  "Obsolete vendor prefix '" + deadPrefix + "' found in stylesheet.",
						Hint:     "Replace with standard CSS property without vendor prefix.",
					})
				}
			}
		}
	}

	return diags
}

func hasCompleteLineClampTriad(s string) bool {
	hasBox := strings.Contains(s, "display:-webkit-box") ||
		strings.Contains(s, "display: -webkit-box") ||
		strings.Contains(s, "display:'-webkit-box'") ||
		strings.Contains(s, "display:\"-webkit-box\"") ||
		strings.Contains(s, "-webkit-box")

	hasOrient := strings.Contains(s, "boxorient") ||
		strings.Contains(s, "box-orient") ||
		strings.Contains(s, "vertical")

	return hasBox && hasOrient
}

func hasLineClamp(s string) bool {
	return strings.Contains(s, "lineclamp") ||
		strings.Contains(s, "line-clamp") ||
		strings.Contains(s, "-webkit-line-clamp")
}

var deadPrefixProperties = [...]string{
	// Engine mati total: Opera Presto
	"-o-transform",
	"-o-transition",
	"-o-animation",
	"-o-linear-gradient",
	"-o-border-radius",
	"-o-box-shadow",
	"-o-background-size",
	"-o-text-overflow",
	"-o-object-fit",
	"-o-user-select",
	"-o-filter",

	// Engine mati total: KHTML
	"-khtml-opacity",
	"-khtml-user-select",

	// IE / Edge lawas
	"-ms-filter",
	"-ms-transform",
	"-ms-linear-gradient",
	"-ms-border-radius",
	"-ms-box-shadow",

	// Properti yang sekarang 100% standar
	"-moz-transition",
	"-webkit-transition",
	"-moz-transform",
	"-webkit-transform",
	"-moz-animation",
	"-webkit-animation",
	"-moz-linear-gradient",
	"-webkit-gradient",
	"-moz-appearance",
	"-webkit-appearance",
	"-webkit-overflow-scrolling",
	"-moz-border-radius",
	"-webkit-border-radius",
	"-moz-box-shadow",
	"-webkit-box-shadow",
	"-moz-box-sizing",
	"-webkit-box-sizing",
	"-moz-opacity",
}

func findDeadPrefixes(s string) []string {
	var found []string
	for _, p := range deadPrefixProperties {
		if strings.Contains(s, p) {
			found = append(found, p)
		}
	}
	return found
}
