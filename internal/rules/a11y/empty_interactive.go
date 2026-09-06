package a11y

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// EmptyInteractiveRule memastikan tombol atau tautan yang hanya berisi ikon SVG
// memiliki nama aksesibel (aria-label, aria-labelledby, atau title).
type EmptyInteractiveRule struct{}

// NewEmptyInteractiveRule membuat instance baru EmptyInteractiveRule.
func NewEmptyInteractiveRule() *EmptyInteractiveRule {
	return &EmptyInteractiveRule{}
}

// ID mengembalikan Charites Rule ID kanonikal berformat a11y.empty-interactive.
func (r *EmptyInteractiveRule) ID() string {
	return "a11y.empty-interactive"
}

// Description mengembalikan deskripsi ringkas aturan.
func (r *EmptyInteractiveRule) Description() string {
	return "Enforces accessible names on interactive elements (buttons, links) containing only icons or visual elements"
}

// Category mengembalikan nama kategori rule.
func (r *EmptyInteractiveRule) Category() string {
	return "a11y"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *EmptyInteractiveRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *EmptyInteractiveRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"W3C Web Content Accessibility Guidelines (WCAG) 2.2 SC 4.1.2 (Name, Role, Value)",
			"W3C Web Content Accessibility Guidelines (WCAG) 2.2 SC 2.4.4 (Link Purpose in Context)",
			"W3C Accessible Name and Description Computation (AccName) 1.2",
		},
		CoreInvariant: "Interactive controls (<button>, <a>, role=\"button\") must provide an accessible name via text content, aria-label, aria-labelledby, or title.",
		Grounding: "Modern UI designs frequently employ icon-only buttons (such as trash cans for delete, pencils for edit, or magnifying glasses for search) without visible text labels.\n\n" +
			"When developers render icon components (e.g. Lucide, Heroicons, or native <svg>) inside <button> or <a> without aria-label:\n" +
			"1. Unnamed Interactive Elements: Assistive technologies announce \"button\" or \"link\" with zero information about the control's purpose or action.\n" +
			"2. Linter Bypass: Conventional linters often pass these elements because the JSX tree contains child elements (<TrashIcon />), failing to realize that rendered SVGs produce no accessible text.\n" +
			"3. Severe Navigation Barrier: Blind or low-vision users cannot determine which action each button executes without activating it.\n\n" +
			"Charites inspects interactive controls to ensure icon-only elements supply an accessible name.",
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Icon-only button without accessible name",
				Code: `<button onClick={handleDelete} className="size-11 flex items-center justify-center">
  <TrashIcon className="size-5 text-destructive" />
</button>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Icon-only button with explicit aria-label",
				Code: `<button type="button" onClick={handleDelete} aria-label="Hapus dokumen ini" className="size-11 flex items-center justify-center">
  <TrashIcon className="size-5 text-destructive" />
</button>`,
			},
		},
		Risks: []ir.RiskItem{
			{
				Vector:   "Blank Interactive Control",
				Severity: "HIGH",
				Impact:   "Screen readers announce an unlabelled button or link without stating its action.",
			},
			{
				Vector:   "Inaccessible Voice Navigation",
				Severity: "MEDIUM",
				Impact:   "Voice control software cannot target or activate the unlabelled control.",
			},
		},
	}
}

// Evaluate memeriksa apakah kontrol interaktif hanya berisi ikon visual tanpa nama aksesibel.
func (r *EmptyInteractiveRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !isInteractiveTargetTag(node) {
		return nil
	}

	if hasSpreadProps(node.Attributes) || isDecorativeOrHidden(node.Attributes) {
		return nil
	}

	if hasExplicitAccessibleAttribute(node) {
		return nil
	}

	if !isIconOnlyInteractive(node) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Interactive element (<" + node.Tag + ">) contains only icon graphics but lacks an accessible name (aria-label, aria-labelledby, or title).",
			Hint:     "Add aria-label=\"...\" to provide an accessible description for assistive technologies.",
		},
	}
}

func isInteractiveTargetTag(node *ir.Node) bool {
	tagLower := strings.ToLower(node.Tag)
	if tagLower == "button" || tagLower == "a" {
		return true
	}

	if role, ok := node.GetAttr("role"); ok && strings.EqualFold(CleanAttr(role), "button") {
		return true
	}

	return false
}

func hasExplicitAccessibleAttribute(node *ir.Node) bool {
	if label, ok := node.GetAttr("aria-label"); ok && CleanAttr(label) != "" {
		return true
	}
	if labelledby, ok := node.GetAttr("aria-labelledby"); ok && CleanAttr(labelledby) != "" {
		return true
	}
	if title, ok := node.GetAttr("title"); ok && CleanAttr(title) != "" {
		return true
	}
	return false
}

func isIconOnlyInteractive(node *ir.Node) bool {
	hasIcon := false
	hasNonIcon := false

	for _, child := range node.Children {
		if child.Type == ir.NodeText {
			text := strings.TrimSpace(child.RawClasses)
			if text != "" && text != "{" && text != "}" {
				hasNonIcon = true
				break
			}
		}

		if child.Type == ir.NodeElement {
			if isSlotOrDynamicExpression(child) {
				hasNonIcon = true
				break
			}
			if isIconElement(child) {
				hasIcon = true
			} else {
				hasNonIcon = true
				break
			}
		}
	}

	return hasIcon && !hasNonIcon
}

func isSlotOrDynamicExpression(child *ir.Node) bool {
	tagLower := strings.ToLower(child.Tag)
	if tagLower == "slot" {
		return true
	}
	return strings.Contains(child.RawClasses, "children") || strings.Contains(child.Tag, "children")
}

func isIconElement(node *ir.Node) bool {
	tagLower := strings.ToLower(node.Tag)
	if tagLower == "svg" || tagLower == "path" || tagLower == "i" {
		return true
	}
	if strings.HasSuffix(node.Tag, "Icon") || strings.HasPrefix(node.Tag, "Icon") {
		return true
	}
	if strings.HasPrefix(node.Tag, "Lucide") {
		return true
	}
	return false
}
