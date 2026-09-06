package ux

import (
	"strconv"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// StripVariantsOnlyBase mengembalikan utility dasar tanpa mengalokasikan slice varian.
// Mengembalikan substring langsung tanpa alokasi memori (0 B/op, 0 allocs/op).
func StripVariantsOnlyBase(token string) string {
	if !strings.Contains(token, ":") {
		return token
	}

	lastColon := -1
	inBracket := false
	for i := 0; i < len(token); i++ {
		b := token[i]
		switch b {
		case '[':
			inBracket = true
		case ']':
			inBracket = false
		case ':':
			if !inBracket {
				lastColon = i
			}
		}
	}

	if lastColon != -1 && lastColon+1 < len(token) {
		return token[lastColon+1:]
	}
	return token
}

// isIdleAffordanceToken mengecek apakah token berlaku di state idle (bukan semata-mata hover/focus/active).
func isIdleAffordanceToken(token string) bool {
	if strings.HasPrefix(token, "hover:") ||
		strings.HasPrefix(token, "focus:") ||
		strings.HasPrefix(token, "active:") ||
		strings.HasPrefix(token, "focus-visible:") ||
		strings.HasPrefix(token, "focus-within:") ||
		strings.HasPrefix(token, "group-hover:") ||
		strings.HasPrefix(token, "peer-hover:") {
		return false
	}
	return true
}

// parseTailwindSpacingNumber mem-parsing angka skala Tailwind (misal: "4", "0.5", "8", "16") ke nilai float.
func parseTailwindSpacingNumber(s string) (float64, bool) {
	if s == "px" {
		return 0.0625, true
	}
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		inner := s[1 : len(s)-1]
		if strings.HasSuffix(inner, "px") {
			v, err := strconv.ParseFloat(strings.TrimSuffix(inner, "px"), 64)
			if err == nil {
				return v / 16.0, true // normalisasi ke rem
			}
		}
		if strings.HasSuffix(inner, "rem") {
			v, err := strconv.ParseFloat(strings.TrimSuffix(inner, "rem"), 64)
			if err == nil {
				return v, true
			}
		}
	}
	v, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return v * 0.25, true // 4 unit = 1 rem
	}
	return 0, false
}

// extractVerticalGapParent mengembalikan nilai vertical gap (dalam rem) dari kelas parent jika ada.
func extractVerticalGapParent(classes []string) (float64, bool) {
	for _, cls := range classes {
		base := StripVariantsOnlyBase(cls)
		switch {
		case strings.HasPrefix(base, "gap-y-"):
			if val, ok := parseTailwindSpacingNumber(base[len("gap-y-"):]); ok {
				return val, true
			}
		case strings.HasPrefix(base, "gap-") && !strings.HasPrefix(base, "gap-x-"):
			if val, ok := parseTailwindSpacingNumber(base[len("gap-"):]); ok {
				return val, true
			}
		case strings.HasPrefix(base, "space-y-") && !strings.HasPrefix(base, "space-y-reverse"):
			if val, ok := parseTailwindSpacingNumber(base[len("space-y-"):]); ok {
				return val, true
			}
		}
	}
	return 0, false
}

// hasTailwindV3SpaceY mengecek apakah parent memakai utility space-y Tailwind v3.
func hasTailwindV3SpaceY(classes []string) (string, bool) {
	for _, cls := range classes {
		base := StripVariantsOnlyBase(cls)
		if strings.HasPrefix(base, "space-y-") && !strings.HasPrefix(base, "space-y-reverse") {
			return base, true
		}
	}
	return "", false
}

// extractChildMarginTop mengecek margin top pada child class.
func extractChildMarginTop(classes []string) (string, float64, bool) {
	for _, cls := range classes {
		base := StripVariantsOnlyBase(cls)
		if strings.HasPrefix(base, "mt-") {
			if val, ok := parseTailwindSpacingNumber(base[len("mt-"):]); ok {
				return base, val, true
			}
		} else if strings.HasPrefix(base, "my-") {
			if val, ok := parseTailwindSpacingNumber(base[len("my-"):]); ok {
				return base, val, true
			}
		}
	}
	return "", 0, false
}

// extractChildIntraSpacing mengembalikan margin/gap terbesar dari elemen child.
func extractChildIntraSpacing(classes []string) (string, float64, bool) {
	maxVal := 0.0
	maxCls := ""
	found := false

	for _, cls := range classes {
		base := StripVariantsOnlyBase(cls)
		var val float64
		var ok bool
		switch {
		case strings.HasPrefix(base, "mb-"):
			val, ok = parseTailwindSpacingNumber(base[len("mb-"):])
		case strings.HasPrefix(base, "mt-"):
			val, ok = parseTailwindSpacingNumber(base[len("mt-"):])
		case strings.HasPrefix(base, "my-"):
			val, ok = parseTailwindSpacingNumber(base[len("my-"):])
		case strings.HasPrefix(base, "gap-y-"):
			val, ok = parseTailwindSpacingNumber(base[len("gap-y-"):])
		case strings.HasPrefix(base, "gap-") && !strings.HasPrefix(base, "gap-x-"):
			val, ok = parseTailwindSpacingNumber(base[len("gap-"):])
		case strings.HasPrefix(base, "space-y-") && !strings.HasPrefix(base, "space-y-reverse"):
			val, ok = parseTailwindSpacingNumber(base[len("space-y-"):])
		}

		if ok && val > maxVal {
			maxVal = val
			maxCls = base
			found = true
		}
	}

	return maxCls, maxVal, found
}

// isNavigationLandmark mengecek apakah elemen bertindak sebagai landmark navigasi.
func isNavigationLandmark(node *ir.Node) bool {
	if node == nil || node.Type != ir.NodeElement {
		return false
	}
	if node.Tag == "nav" {
		return true
	}
	if role, ok := node.GetAttr("role"); ok && role == "navigation" {
		return true
	}
	return false
}

// isNavLinkNode mengecek apakah node merepresentasikan link navigasi.
func isNavLinkNode(node *ir.Node) bool {
	if node == nil || node.Type != ir.NodeElement {
		return false
	}
	// Jangan hitung tombol aksi!
	if node.Tag == "button" || strings.HasSuffix(node.Tag, "Button") {
		return false
	}
	if role, ok := node.GetAttr("role"); ok {
		if role == "button" || role == "search" || role == "switch" {
			return false
		}
		if role == "link" {
			return true
		}
	}

	if node.Tag == "a" || strings.HasSuffix(node.Tag, "Link") {
		return true
	}
	if _, ok := node.GetAttr("href"); ok {
		return true
	}
	if _, ok := node.GetAttr("to"); ok {
		return true
	}
	return false
}

func isChunkingTagOrAttr(n *ir.Node) bool {
	switch n.Tag {
	case "DropdownMenu", "Dropdown", "NavigationMenu", "Sheet", "Drawer", "Disclosure", "details", "Select":
		return true
	}
	if strings.Contains(n.Tag, "Dropdown") || strings.Contains(n.Tag, "Drawer") || strings.Contains(n.Tag, "Sheet") {
		return true
	}
	if _, ok := n.GetAttr("aria-expanded"); ok {
		return true
	}
	if _, ok := n.GetAttr("aria-haspopup"); ok {
		return true
	}
	return false
}

func hasChunkingTriggerText(n *ir.Node) bool {
	for _, child := range n.Children {
		if child.Type == ir.NodeText {
			txt := strings.ToLower(strings.TrimSpace(child.RawClasses))
			if txt == "lainnya" || txt == "more" || txt == "..." || txt == "menu" {
				return true
			}
		}
	}
	return false
}

// hasChunkingMechanism mengecek apakah di dalam subtree terdapat mekanisme chunking navigasi.
func hasChunkingMechanism(node *ir.Node) bool {
	if node == nil {
		return false
	}

	for n := range node.Walk() {
		if n == node {
			continue
		}
		if isChunkingTagOrAttr(n) || hasChunkingTriggerText(n) {
			return true
		}
	}
	return false
}

// isActionContainer mengecek apakah elemen berfungsi sebagai kontainer kumpulan aksi (dialog footer, action bar, card actions, flex row).
func isActionContainer(node *ir.Node) bool {
	if node == nil || node.Type != ir.NodeElement {
		return false
	}

	tagLower := strings.ToLower(node.Tag)
	if strings.HasSuffix(tagLower, "footer") || strings.HasSuffix(tagLower, "actionbar") ||
		strings.HasSuffix(tagLower, "actions") || strings.HasSuffix(tagLower, "toolbar") {
		return true
	}

	if role, ok := node.GetAttr("role"); ok && (role == "toolbar" || role == "group") {
		return true
	}

	hasFlexOrGrid := false
	for _, cls := range node.Classes {
		base := StripVariantsOnlyBase(cls)
		if base == "flex" || base == "inline-flex" || base == "grid" {
			hasFlexOrGrid = true
			break
		}
	}

	return hasFlexOrGrid
}

// isButtonNode mengecek apakah node adalah tombol interaktif.
func isButtonNode(node *ir.Node) bool {
	if node == nil || node.Type != ir.NodeElement {
		return false
	}
	if node.Tag == "button" || strings.HasSuffix(node.Tag, "Button") {
		return true
	}
	if role, ok := node.GetAttr("role"); ok && role == "button" {
		return true
	}
	return false
}

// isPrimaryCTAButton mengecek apakah tombol berstatus primary/utama.
func isPrimaryCTAButton(node *ir.Node) bool {
	if !isButtonNode(node) {
		return false
	}

	// Cek eksplisit variant non-primary
	if variant, ok := node.GetAttr("variant"); ok {
		v := strings.ToLower(variant)
		if v == "outline" || v == "ghost" || v == "secondary" || v == "link" || v == "destructive" {
			return false
		}
		if v == "primary" || v == "default" {
			return true
		}
	}

	// Cek kelas styling
	hasPrimaryClass := false
	hasSecondaryOrGhostClass := false

	for _, cls := range node.Classes {
		base := StripVariantsOnlyBase(cls)
		switch base {
		case "bg-primary", "btn-primary", "bg-blue-600", "bg-indigo-600", "bg-black", "bg-neutral-900":
			hasPrimaryClass = true
		case "bg-transparent", "border", "border-input", "bg-secondary", "btn-secondary", "bg-muted", "text-destructive", "bg-destructive":
			hasSecondaryOrGhostClass = true
		}
	}

	if hasSecondaryOrGhostClass {
		return false
	}

	return hasPrimaryClass
}

// isProseContext mengecek apakah node berada dalam konteks bacaan paragraf/artikel (prose).
func isProseContext(node *ir.Node) bool {
	if node == nil {
		return false
	}

	curr := node.Parent
	for curr != nil {
		// Landmark bukan konteks prose
		switch curr.Tag {
		case "nav", "header", "footer", "aside", "table", "thead", "tbody":
			return false
		}

		// Cek tag paragraf atau list
		switch curr.Tag {
		case "p", "li", "blockquote", "article":
			return true
		}

		// Cek kelas penanda prose
		for _, cls := range curr.Classes {
			base := StripVariantsOnlyBase(cls)
			if base == "prose" || base == "rich-text" || strings.HasPrefix(base, "prose-") {
				return true
			}
		}

		curr = curr.Parent
	}

	return false
}

// isStyledAsButton mengecek apakah tag <a> didesain menyerupai tombol fisik.
func isStyledAsButton(node *ir.Node) bool {
	if node == nil {
		return false
	}

	hasButtonClass := false
	for _, cls := range node.Classes {
		base := StripVariantsOnlyBase(cls)
		if base == "btn" || base == "button" || strings.HasPrefix(base, "btn-") {
			return true
		}
		if base == "rounded-md" || base == "rounded-lg" || base == "rounded-full" || base == "px-4" || base == "py-2" {
			hasButtonClass = true
		}
	}

	if role, ok := node.GetAttr("role"); ok && role == "button" {
		return true
	}

	return hasButtonClass
}

// hasPersistentUnderlineAffordance mengecek apakah anchor memiliki underline saat idle.
func hasPersistentUnderlineAffordance(classes []string) bool {
	hasUnderline := false
	hasNoUnderline := false

	for _, cls := range classes {
		base := StripVariantsOnlyBase(cls)
		if !isIdleAffordanceToken(cls) {
			continue // abaikan hover:underline untuk state idle!
		}

		if base == "underline" || strings.HasPrefix(base, "underline-") || base == "border-b" || strings.HasPrefix(base, "decoration-") {
			hasUnderline = true
		}
		if base == "no-underline" {
			hasNoUnderline = true
		}
	}

	if hasNoUnderline {
		return false
	}

	return hasUnderline
}
