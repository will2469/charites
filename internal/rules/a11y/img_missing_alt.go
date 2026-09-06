package a11y

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// ImgMissingAltRule memastikan setiap komponen gambar Astro (<Image />, <Picture />)
// dan tag <img> native pada TSX/MDX menyertakan atribut 'alt' sesuai WCAG 1.1.1.
type ImgMissingAltRule struct{}

// NewImgMissingAltRule membuat instance baru ImgMissingAltRule.
func NewImgMissingAltRule() *ImgMissingAltRule {
	return &ImgMissingAltRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat a11y.img-missing-alt.
func (r *ImgMissingAltRule) ID() string {
	return "a11y.img-missing-alt"
}

// Description mengembalikan deskripsi ringkas aturan.
func (r *ImgMissingAltRule) Description() string {
	return "Enforces required 'alt' attribute on Astro <Image>, <Picture>, and native <img> elements (WCAG 1.1.1)"
}

// Category mengembalikan nama kategori rule.
func (r *ImgMissingAltRule) Category() string {
	return "a11y"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *ImgMissingAltRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ImgMissingAltRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Web Content Accessibility Guidelines (WCAG) 2.2 Success Criterion 1.1.1 (Non-text Content)",
			"HTML Living Standard Section 4.8.4.4 (Requirements for providing alternative text)",
			"Astro Assets Image & Picture Component Specification",
		},
		CoreInvariant: "All visual image elements (<Image />, <Picture />, <img>) must provide an 'alt' attribute or explicit accessibility role.",
		Grounding: "Alternative text ('alt') conveys the purpose, meaning, or narrative of visual images to blind and low-vision users utilizing screen readers, and displays as fallback text when network requests fail or images are blocked.\n\n" +
			"When developers omit the 'alt' attribute entirely:\n" +
			"1. Screen Reader Verbosity Failure: Screen readers resort to reciting the entire raw file URL (e.g. '/assets/images/hero_banner_v2_final_compressed.webp'), creating auditory confusion.\n" +
			"2. SEO & Fallback Degradation: Search crawlers and low-bandwidth users receive no descriptive context.\n" +
			"3. Astro Component Blindspot: Standard ESLint jsx-a11y rules only inspect literal <img> tags, failing to flag Astro <Image /> and <Picture /> components from 'astro:assets'.\n\n" +
			"Charites validates all image tags, including Astro framework asset components, permitting alt=\"\" exclusively for decorative graphics.",
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Astro Image component missing required alt attribute",
				Code:     `<Image src={bannerImage} width={800} height={400} />`,
			},
			{
				Language: "tsx",
				Comment:  "Native img tag without alt",
				Code:     `<img src="/profile.png" className="rounded-full size-10" />`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Informative image with meaningful alternative description",
				Code:     `<Image src={bannerImage} alt="Dashboard analitik desa tahun 2026" width={800} height={400} />`,
			},
			{
				Language: "tsx",
				Comment:  "Decorative image explicitly marked with empty alt",
				Code:     `<img src="/decor-pattern.svg" alt="" className="pointer-events-none opacity-50" />`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Total Information Loss for Screen Readers",
				Severity: "HIGH",
				Impact:   "Blind users cannot understand the visual information or interactive purpose of images.",
			},
			{
				Vector:   "Broken Network Fallback",
				Severity: "MEDIUM",
				Impact:   "When image assets fail to load on slow mobile connections, users see broken boxes without context.",
			},
		},
	}
}

// Evaluate memeriksa keberadaan atribut alt pada elemen gambar.
// Mematuhi 0 B/op, 0 allocs/op pada node bersih (QUAL-03).
func (r *ImgMissingAltRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || !isImageTag(node.Tag) {
		return nil
	}

	if hasAccessibleImageAlt(node.Attributes) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Image element <%s> is missing required 'alt' attribute (WCAG 1.1.1 Non-text Content)", node.Tag),
			Hint:     "Provide a descriptive 'alt=\"...\"' for informative images, or an empty 'alt=\"\"' for decorative images.",
		},
	}
}

func hasAccessibleImageAlt(attrs map[string]string) bool {
	if attrs == nil {
		return false
	}
	if _, hasAlt := attrs["alt"]; hasAlt {
		return true
	}
	if hasAccessibleLabel(attrs) {
		return true
	}
	if isDecorativeOrHidden(attrs) {
		return true
	}
	return hasSpreadProps(attrs)
}

func hasAccessibleLabel(attrs map[string]string) bool {
	if label, ok := attrs["aria-label"]; ok && CleanAttr(label) != "" {
		return true
	}
	if labelledby, ok := attrs["aria-labelledby"]; ok && CleanAttr(labelledby) != "" {
		return true
	}
	return false
}

func isDecorativeOrHidden(attrs map[string]string) bool {
	if hidden, ok := attrs["aria-hidden"]; ok && CleanAttr(hidden) == "true" {
		return true
	}
	if role, ok := attrs["role"]; ok {
		cleanRole := strings.ToLower(CleanAttr(role))
		if cleanRole == "presentation" || cleanRole == "none" {
			return true
		}
	}
	return false
}

func hasSpreadProps(attrs map[string]string) bool {
	for k := range attrs {
		if strings.HasPrefix(k, "{...") {
			return true
		}
	}
	return false
}

func isImageTag(tag string) bool {
	tagLower := strings.ToLower(tag)
	if tagLower == "img" || tagLower == "astro-image" {
		return true
	}
	if tag == "Image" || tag == "Picture" {
		return true
	}
	if strings.HasSuffix(tag, ".Image") || strings.HasSuffix(tag, ".Picture") {
		return true
	}
	return false
}
