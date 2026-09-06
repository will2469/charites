package ergonomy

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// BottomNavThumbUnreachableRule mendeteksi tombol aksi primer (Call-to-Action) yang ditempatkan secara eksklusif
// pada header atas mobile tanpa alternatif yang dapat dijangkau ibu jari di zona bawah (Fitts's Law / Thumb Zone).
type BottomNavThumbUnreachableRule struct{}

// NewBottomNavThumbUnreachableRule membuat instance baru dari BottomNavThumbUnreachableRule.
func NewBottomNavThumbUnreachableRule() *BottomNavThumbUnreachableRule {
	return &BottomNavThumbUnreachableRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *BottomNavThumbUnreachableRule) ID() string {
	return "ergonomy.bottom-nav-thumb-unreachable"
}

// Description mengembalikan ringkasan aturan.
func (r *BottomNavThumbUnreachableRule) Description() string {
	return "Warns when primary call-to-action (CTA) buttons are exclusively located in the top mobile header without reachable alternatives in the bottom thumb zone"
}

// Category mengembalikan nama kategori rule.
func (r *BottomNavThumbUnreachableRule) Category() string {
	return "ergonomy"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (info).
func (r *BottomNavThumbUnreachableRule) DefaultSeverity() ir.Severity {
	return ir.SeverityInfo
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *BottomNavThumbUnreachableRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Steven Hoober (2017), Designing for Touch & Mobile Thumb Zone Research",
			"Fitts's Law of Motor Movement Ergonomics on Tall Mobile Displays",
			"Apple Human Interface Guidelines (Navigation Bars & Bottom Toolbars)",
			"Google Material Design 3 (Bottom App Bars & Floating Action Buttons)",
		},
		CoreInvariant: "Primary call-to-action controls (e.g. form submissions, checkout confirmations) must be reachable within the lower mobile thumb zone rather than positioned exclusively in top headers.",
		Grounding: "On modern mobile screens (6.1-inch to 6.7-inch+), the top one-third of the screen lies in the 'Hard to Reach' or 'Ow Zone' for one-handed thumb navigation (Steven Hoober's Thumb Zone research).\n\n" +
			"Placing the sole primary submission or action button exclusively in a top navigation header (<header> or 'top-0' container) forces users to awkwardly shift grip or use two hands.\n\n" +
			"Providing a primary CTA in the lower thumb zone (e.g., sticky bottom bar, bottom sheet, or natural form footer) satisfies Fitts's Law and optimizes one-handed mobile usability.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Thumb Strain and Awkward Grip Shifting",
				Severity: "LOW",
				Impact:   "Users on large smartphones experience physical discomfort or drop hazards when repeatedly reaching for top-corner primary actions.",
			},
			{
				Vector:   "Decreased Form Completion Rates",
				Severity: "LOW",
				Impact:   "One-handed mobile users abandon multi-step forms due to friction reaching top submission buttons.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Primary submit button trapped in top sticky header without bottom alternative",
				Code: `<header className="sticky top-0 z-10 flex items-center justify-between p-4 bg-background border-b">
  <button type="button" aria-label="Kembali">
    <ArrowLeft className="w-6 h-6" />
  </button>
  <h1 className="font-semibold text-lg">Edit Profil Warga</h1>
  <button type="submit" className="h-10 px-4 bg-primary text-primary-foreground rounded-xl">
    Simpan
  </button>
</header>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Primary CTA positioned in reachable bottom thumb zone",
				Code: `<header className="sticky top-0 z-10 flex items-center justify-between p-4 bg-background border-b">
  <button type="button" aria-label="Kembali">
    <ArrowLeft className="w-6 h-6" />
  </button>
  <h1 className="font-semibold text-lg">Edit Profil Warga</h1>
</header>
<main className="p-4 pb-24">
  <input name="nama" placeholder="Nama Lengkap" />
</main>
<footer className="fixed bottom-0 inset-x-0 p-4 bg-background border-t pb-[env(safe-area-inset-bottom)]">
  <button type="submit" className="w-full h-12 bg-primary text-primary-foreground rounded-xl font-semibold">
    Simpan Perubahan
  </button>
</footer>`,
			},
		},
	}
}

// Evaluate memeriksa apakah elemen tombol aksi primer ditempatkan di header atas tanpa alternatif bawah.
func (r *BottomNavThumbUnreachableRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	if !isInteractiveButton(node) {
		return nil
	}

	if !isInTopContainer(node) {
		return nil
	}

	if isDesktopOnly(node) {
		return nil
	}

	if isSecondaryOrNavButton(node) {
		return nil
	}

	if !isPrimaryCTA(node) {
		return nil
	}

	root := getDocumentRoot(node)
	if hasBottomThumbAlternative(root, node) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Primary call-to-action (<" + node.Tag + ">) is placed exclusively in the top navigation header, making it hard to reach with one hand on tall mobile screens (Hoober Thumb Zone / Fitts's Law).",
			Hint:     "Provide a primary CTA alternative in the lower thumb zone (e.g. sticky bottom bar, bottom sheet, or footer action) for comfortable one-handed mobile reachability.",
		},
	}
}

func isInteractiveButton(node *ir.Node) bool {
	tagLower := strings.ToLower(node.Tag)
	if tagLower == "button" {
		return true
	}
	if tagLower == "input" && node.Attributes != nil {
		if t, ok := node.Attributes["type"]; ok {
			cleanType := cleanAttrValue(t)
			if cleanType == "submit" || cleanType == "button" {
				return true
			}
		}
	}
	if node.Attributes != nil {
		if role, ok := node.Attributes["role"]; ok {
			cleanRole := cleanAttrValue(role)
			if cleanRole == "button" {
				return true
			}
		}
	}
	return false
}

func isInTopContainer(node *ir.Node) bool {
	for p := node.Parent; p != nil; p = p.Parent {
		if p.Type != ir.NodeElement {
			continue
		}
		if strings.EqualFold(p.Tag, "header") {
			return true
		}
		if hasTopPositionClass(p.Classes, p.RawClasses) {
			return true
		}
	}
	return hasTopPositionClass(node.Classes, node.RawClasses)
}

func hasTopPositionClass(classes []string, rawClasses string) bool {
	if strings.Contains(rawClasses, "top-0") {
		return true
	}
	for _, cls := range classes {
		if cls == "top-0" || strings.HasSuffix(cls, ":top-0") {
			return true
		}
	}
	return false
}

func isDesktopOnly(node *ir.Node) bool {
	if hasDesktopOnlyClass(node.Classes, node.RawClasses) {
		return true
	}
	for p := node.Parent; p != nil; p = p.Parent {
		if p.Type == ir.NodeElement && hasDesktopOnlyClass(p.Classes, p.RawClasses) {
			return true
		}
	}
	return false
}

func hasDesktopOnlyClass(classes []string, rawClasses string) bool {
	if strings.Contains(rawClasses, "hidden md:") ||
		strings.Contains(rawClasses, "hidden sm:") ||
		strings.Contains(rawClasses, "hidden lg:") ||
		strings.Contains(rawClasses, "max-sm:hidden") ||
		strings.Contains(rawClasses, "max-md:hidden") {
		return true
	}
	for _, cls := range classes {
		if strings.HasPrefix(cls, "md:flex") || strings.HasPrefix(cls, "md:block") || strings.HasPrefix(cls, "md:inline") {
			if strings.Contains(rawClasses, "hidden") {
				return true
			}
		}
	}
	return false
}

var secondaryKeywords = [...]string{
	"kembali", "back", "tutup", "close", "batal", "cancel",
	"menu", "search", "cari", "filter", "bantuan", "help",
	"navigasi", "toggle",
}

func containsSecondaryKeyword(text string) bool {
	lower := strings.ToLower(text)
	for _, kw := range secondaryKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

func isSecondaryOrNavButton(node *ir.Node) bool {
	if node.Attributes != nil && isSecondaryByAttrs(node.Attributes) {
		return true
	}
	return isSecondaryByChildren(node.Children)
}

func isSecondaryByAttrs(attrs map[string]string) bool {
	if label, ok := attrs["aria-label"]; ok && containsSecondaryKeyword(label) {
		return true
	}
	if v, ok := attrs["variant"]; ok {
		cleanV := cleanAttrValue(v)
		return cleanV == "ghost" || cleanV == "outline" || cleanV == "secondary" || cleanV == "subtle"
	}
	return false
}

func isSecondaryByChildren(children []*ir.Node) bool {
	for _, child := range children {
		if child.Type != ir.NodeText {
			continue
		}
		if containsSecondaryKeyword(child.RawClasses) {
			return true
		}
		trimmed := strings.TrimSpace(child.RawClasses)
		if trimmed == "<" || trimmed == "x" || trimmed == "X" {
			return true
		}
	}
	return false
}

func isPrimaryCTA(node *ir.Node) bool {
	if node.Attributes != nil {
		if t, ok := node.Attributes["type"]; ok {
			cleanType := cleanAttrValue(t)
			if cleanType == "submit" {
				return true
			}
		}
	}

	if hasPrimaryStyling(node.Classes, node.RawClasses) {
		return true
	}

	return hasPrimaryActionText(node)
}

func hasPrimaryStyling(classes []string, rawClasses string) bool {
	if strings.Contains(rawClasses, "bg-primary") ||
		strings.Contains(rawClasses, "btn-primary") {
		return true
	}
	for _, cls := range classes {
		if cls == "bg-primary" || cls == "btn-primary" {
			return true
		}
	}
	return false
}

var primaryKeywords = [...]string{
	"simpan", "submit", "kirim", "bayar", "checkout", "selesai",
	"daftar", "beli", "lanjut", "save", "next", "confirm", "konfirmasi",
	"pesan", "order", "selesaikan",
}

func containsPrimaryKeyword(text string) bool {
	lower := strings.ToLower(text)
	for _, kw := range primaryKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

func hasPrimaryActionText(node *ir.Node) bool {
	for _, child := range node.Children {
		if child.Type == ir.NodeText && containsPrimaryKeyword(child.RawClasses) {
			return true
		}
	}
	if node.Attributes != nil {
		if ariaLabel, ok := node.Attributes["aria-label"]; ok && containsPrimaryKeyword(ariaLabel) {
			return true
		}
		if val, ok := node.Attributes["value"]; ok && containsPrimaryKeyword(val) {
			return true
		}
	}
	return false
}

func getDocumentRoot(node *ir.Node) *ir.Node {
	curr := node
	for curr.Parent != nil {
		curr = curr.Parent
	}
	return curr
}

func hasBottomThumbAlternative(root *ir.Node, targetButton *ir.Node) bool {
	if root == nil {
		return false
	}

	for curr := range root.Walk() {
		if curr == targetButton {
			continue
		}

		if curr.Type == ir.NodeElement && hasBottomPositionClass(curr.Classes, curr.RawClasses) {
			return true
		}

		if isInteractiveButton(curr) && !isInTopContainer(curr) && isPrimaryCTA(curr) {
			return true
		}
	}

	return false
}

func hasBottomPositionClass(classes []string, rawClasses string) bool {
	if strings.Contains(rawClasses, "bottom-0") ||
		strings.Contains(rawClasses, "bottom-2") ||
		strings.Contains(rawClasses, "bottom-4") ||
		strings.Contains(rawClasses, "bottom-6") ||
		strings.Contains(rawClasses, "bottom-8") ||
		strings.Contains(rawClasses, "bottom-[") {
		return true
	}
	for _, cls := range classes {
		if strings.HasPrefix(cls, "bottom-") || strings.Contains(cls, ":bottom-") {
			return true
		}
	}
	return false
}
