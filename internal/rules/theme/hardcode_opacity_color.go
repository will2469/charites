package theme

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// OpacityTokenMap memetakan utility warna dengan modifier slash opacity ke token semantik resmi pengganti.
// Sesuai SPEC-03-RULES Section 3.2.A.
var OpacityTokenMap = map[string]string{
	"primary/10":     "primary-light",
	"primary/20":     "primary-light",
	"primary/5":      "primary-subtle",
	"secondary/10":   "muted-light",
	"secondary/5":    "muted-subtle",
	"destructive/10": "destructive-light",
	"destructive/20": "destructive-light",
	"destructive/5":  "destructive-subtle",
	"accent/10":      "accent-light",
	"accent/20":      "accent-light",
	"accent/5":       "accent-subtle",
	"warning/10":     "warning-light",
	"warning/5":      "warning-subtle",
	"muted/10":       "muted-light",
	"muted/5":        "muted-subtle",
	"amber/10":       "amber-light",
	"amber/5":        "amber-subtle",
	"emerald/10":     "emerald-light",
	"emerald/5":      "emerald-subtle",
}

// HardcodeOpacityColorRule mengimplementasikan aturan static analysis "theme.hardcode-opacity-color".
type HardcodeOpacityColorRule struct{}

// NewHardcodeOpacityColorRule membuat instance baru HardcodeOpacityColorRule.
func NewHardcodeOpacityColorRule() *HardcodeOpacityColorRule {
	return &HardcodeOpacityColorRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat <category>.<slug>.
func (r *HardcodeOpacityColorRule) ID() string {
	return "theme.hardcode-opacity-color"
}

// Description mengembalikan penjelasan ringkas maksud dan tujuan rule.
func (r *HardcodeOpacityColorRule) Description() string {
	return "Detects utility classes with hardcoded slash opacity modifiers that have official semantic token replacements"
}

// Category mengembalikan nama kategori rule.
func (r *HardcodeOpacityColorRule) Category() string {
	return "theme"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *HardcodeOpacityColorRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// stripVariants membuang seluruh prefix varian Tailwind (seperti hover:, dark:, md:hover:)
// dan mengembalikan base utility class.
func stripVariants(token string) string {
	lastColon := strings.LastIndexByte(token, ':')
	if lastColon >= 0 && lastColon < len(token)-1 {
		return token[lastColon+1:]
	}
	return token
}

// Evaluate mengevaluasi sebuah node IR dan mendeteksi utility color ber-slash opacity yang memiliki
// pemetaan semantic token pengganti resmi. Mematuhi kontrak pure function dan zero alloc pada node bersih.
func (r *HardcodeOpacityColorRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || len(node.Classes) == 0 {
		return nil
	}

	//nolint:prealloc // zero-alloc on clean nodes required by QUAL-03 allocation invariant
	var diags []ir.Diagnostic
	for _, class := range node.Classes {
		if strings.IndexByte(class, '/') == -1 {
			continue
		}

		base := stripVariants(class)
		var colorSlash string
		switch {
		case strings.HasPrefix(base, "bg-"):
			colorSlash = base[3:]
		case strings.HasPrefix(base, "text-"):
			colorSlash = base[5:]
		case strings.HasPrefix(base, "border-"):
			colorSlash = base[7:]
		case strings.HasPrefix(base, "ring-"):
			colorSlash = base[5:]
		default:
			continue
		}

		replacement, ok := OpacityTokenMap[colorSlash]
		if !ok {
			continue
		}

		diags = append(diags, ir.Diagnostic{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Hardcode opacity color: \"" + class + "\"",
			Hint:     "Use semantic token \"" + replacement + "\".",
		})
	}

	return diags
}
