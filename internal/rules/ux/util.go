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

// cleanAttrValue membersihkan nilai atribut (lowercase, trim whitespace dan tanda kutip).
func cleanAttrValue(v string) string {
	return strings.Trim(strings.TrimSpace(strings.ToLower(v)), "\"'`{}")
}

// getAttrCaseInsensitive mengambil nilai atribut secara case-insensitive untuk beberapa kemungkinan nama key.
func getAttrCaseInsensitive(node *ir.Node, names ...string) (string, bool) {
	if node == nil || node.Attributes == nil {
		return "", false
	}
	for _, target := range names {
		for k, v := range node.Attributes {
			if strings.EqualFold(k, target) {
				return v, true
			}
		}
	}
	return "", false
}

// isInsideHeaderBanner mengecek apakah node berada di dalam <header> atau elemen ber-role "banner".
func isInsideHeaderBanner(node *ir.Node) bool {
	if node == nil {
		return false
	}
	curr := node.Parent
	for curr != nil {
		if curr.Type == ir.NodeElement {
			if curr.Tag == "header" {
				return true
			}
			if role, ok := getAttrCaseInsensitive(curr, "role"); ok && strings.ToLower(role) == "banner" {
				return true
			}
		}
		curr = curr.Parent
	}
	return false
}

// findEnclosingHeaderLink mencari tautan (<a> atau <Link>) pembungkus terdekat di dalam header.
func findEnclosingHeaderLink(node *ir.Node) *ir.Node {
	if node == nil {
		return nil
	}
	curr := node.Parent
	for curr != nil {
		if curr.Type == ir.NodeElement {
			// Berhenti jika sudah mencapai batas header
			if curr.Tag == "header" {
				break
			}
			if role, ok := getAttrCaseInsensitive(curr, "role"); ok && strings.ToLower(role) == "banner" {
				break
			}
			if curr.Tag == "a" || strings.HasSuffix(curr.Tag, "Link") {
				return curr
			}
			if role, ok := getAttrCaseInsensitive(curr, "role"); ok && strings.ToLower(role) == "link" {
				return curr
			}
		}
		curr = curr.Parent
	}
	return nil
}

// isBrandIdentityElement mengecek apakah elemen merepresentasikan identitas merek atau logo situs.
func isBrandIdentityElement(node *ir.Node) bool {
	if node == nil || node.Type != ir.NodeElement {
		return false
	}
	if isNonBrandControlTag(node.Tag) {
		return false
	}
	return isExplicitBrandTag(node.Tag) || hasBrandAttribute(node) || hasBrandClassOrID(node)
}

func isNonBrandControlTag(tag string) bool {
	switch tag {
	case "button", "input", "select", "textarea":
		return true
	default:
		return false
	}
}

func isExplicitBrandTag(tag string) bool {
	tagLower := strings.ToLower(tag)
	return tagLower == "logo" || tagLower == "sitelogo" || tagLower == "brandlogo" || tagLower == "brand"
}

func hasBrandAttribute(node *ir.Node) bool {
	if alt, ok := getAttrCaseInsensitive(node, "alt"); ok {
		altLower := strings.ToLower(alt)
		if strings.Contains(altLower, "logo") || strings.Contains(altLower, "brand") {
			return true
		}
	}
	if label, ok := getAttrCaseInsensitive(node, "aria-label"); ok {
		lblLower := strings.ToLower(label)
		if strings.Contains(lblLower, "logo") || strings.Contains(lblLower, "brand") {
			return true
		}
	}
	if node.Tag == "img" || strings.HasSuffix(node.Tag, "Image") || node.Tag == "svg" {
		if src, ok := getAttrCaseInsensitive(node, "src"); ok {
			srcLower := strings.ToLower(src)
			if strings.Contains(srcLower, "logo") || strings.Contains(srcLower, "brand") {
				return true
			}
		}
	}
	return false
}

func hasBrandClassOrID(node *ir.Node) bool {
	for _, cls := range node.Classes {
		base := StripVariantsOnlyBase(cls)
		baseLower := strings.ToLower(base)
		if baseLower == "logo" || baseLower == "brand" || strings.HasPrefix(baseLower, "logo-") ||
			strings.HasSuffix(baseLower, "-logo") || strings.Contains(baseLower, "brand-logo") ||
			baseLower == "site-logo" || baseLower == "site-title" || baseLower == "navbar-brand" {
			return true
		}
	}
	if id, ok := getAttrCaseInsensitive(node, "id"); ok {
		idLower := strings.ToLower(id)
		if strings.Contains(idLower, "logo") || strings.Contains(idLower, "brand") {
			return true
		}
	}
	return false
}

// isNormalizedRootHref memeriksa apakah target tautan mengarah ke root home page ("/").
func isNormalizedRootHref(href string) bool {
	clean := cleanAttrValue(href)
	return isBasicRootHref(clean) || isLocalizedRootHref(clean)
}

func isBasicRootHref(clean string) bool {
	if clean == "/" || clean == "" || clean == "./" {
		return true
	}
	if strings.HasPrefix(clean, "/#") || strings.HasPrefix(clean, "/?") {
		return true
	}
	if strings.HasPrefix(clean, "./#") || strings.HasPrefix(clean, "./?") {
		return true
	}
	return false
}

func isLocalizedRootHref(clean string) bool {
	if !strings.HasPrefix(clean, "/") {
		return false
	}
	remainder := clean[1:]
	slashIdx := strings.IndexByte(remainder, '/')
	queryIdx := strings.IndexAny(remainder, "?#")

	segment := remainder
	if slashIdx != -1 {
		segment = remainder[:slashIdx]
	} else if queryIdx != -1 {
		segment = remainder[:queryIdx]
	}

	if len(segment) != 2 || segment[0] < 'a' || segment[0] > 'z' || segment[1] < 'a' || segment[1] > 'z' {
		return false
	}
	if slashIdx == -1 {
		return true
	}
	sub := remainder[slashIdx+1:]
	return sub == "" || strings.HasPrefix(sub, "?") || strings.HasPrefix(sub, "#")
}

// isFormChunkingContainer mengecek apakah elemen berfungsi sebagai kontainer pembatas form (fieldset, step, tab).
func isFormChunkingContainer(node *ir.Node) bool {
	if node == nil || node.Type != ir.NodeElement {
		return false
	}
	if node.Tag == "fieldset" {
		return true
	}
	tagLower := strings.ToLower(node.Tag)
	switch tagLower {
	case "step", "stepper", "wizardstep", "wizard", "tabpanel", "tabs", "tablist", "accordionitem", "accordion":
		return true
	}
	if role, ok := getAttrCaseInsensitive(node, "role"); ok {
		rLower := strings.ToLower(role)
		if rLower == "group" || rLower == "tabpanel" {
			return true
		}
	}
	return false
}

// isInteractiveFormField mengecek apakah elemen adalah field form interaktif.
// Mengembalikan isField dan logicalGroupKey (jika radio dengan nama tertentu, untuk dikelompokkan).
func isInteractiveFormField(node *ir.Node) (bool, string) {
	if node == nil || node.Type != ir.NodeElement {
		return false, ""
	}

	tagLower := strings.ToLower(node.Tag)

	if tagLower == "button" || strings.HasSuffix(tagLower, "button") {
		return false, ""
	}

	if tagLower == "input" || node.Tag == "Input" {
		typeVal, _ := getAttrCaseInsensitive(node, "type")
		typeClean := cleanAttrValue(typeVal)
		switch typeClean {
		case "hidden", "submit", "button", "reset", "image":
			return false, ""
		case "radio":
			nameVal, _ := getAttrCaseInsensitive(node, "name")
			nameClean := cleanAttrValue(nameVal)
			if nameClean == "" {
				nameClean = "__unnamed_radio__"
			}
			return true, "radio:" + nameClean
		default:
			return true, ""
		}
	}

	if tagLower == "select" || node.Tag == "Select" {
		return true, ""
	}

	if tagLower == "textarea" || node.Tag == "Textarea" {
		return true, ""
	}

	if node.Tag == "Combobox" || node.Tag == "DatePicker" {
		return true, ""
	}

	return false, ""
}

// detectAutofillCategory mendeteksi apakah input meminta data pribadi/PII, kredensial, atau pembayaran.
func detectAutofillCategory(node *ir.Node) (category string, expectedToken string, isSevere bool, identifier string) {
	if node == nil || node.Type != ir.NodeElement {
		return "", "", false, ""
	}

	tagLower := strings.ToLower(node.Tag)
	if tagLower != "input" && node.Tag != "Input" && tagLower != "textarea" && node.Tag != "Textarea" {
		return "", "", false, ""
	}

	typeVal, _ := getAttrCaseInsensitive(node, "type")
	typeClean := cleanAttrValue(typeVal)
	if typeClean == "" {
		typeClean = "text"
	}
	if isNonTextualInputType(typeClean) {
		return "", "", false, ""
	}

	id := extractSemanticFieldID(node)
	if isSearchOrFilterField(typeClean, id) {
		return "", "", false, ""
	}

	cat, token, severe := matchAutofillPattern(typeClean, id)
	return cat, token, severe, id
}

func isNonTextualInputType(typeClean string) bool {
	switch typeClean {
	case "hidden", "submit", "button", "reset", "checkbox", "radio", "file", "image", "color", "range", "date", "datetime-local", "time", "month", "week":
		return true
	default:
		return false
	}
}

func extractSemanticFieldID(node *ir.Node) string {
	for _, key := range [...]string{"name", "id", "placeholder", "aria-label"} {
		if val, ok := getAttrCaseInsensitive(node, key); ok {
			c := cleanAttrValue(val)
			if c != "" {
				return c
			}
		}
	}
	return ""
}

func isSearchOrFilterField(typeClean string, id string) bool {
	if typeClean == "search" {
		return true
	}
	for _, kw := range [...]string{"search", "cari", "filter", "query", "keyword"} {
		if strings.Contains(id, kw) {
			return true
		}
	}
	return false
}

func matchAutofillPattern(typeClean string, id string) (string, string, bool) {
	// 1. Password
	if typeClean == "password" || strings.Contains(id, "password") || strings.Contains(id, "passwd") || strings.Contains(id, "sandi") {
		return "password", "current-password", true
	}

	// 2. Payment / Credit Card
	for _, kw := range [...]string{"cc-", "cardnumber", "card-number", "card_number", "cardnum", "cvv", "cvc", "ccexp", "exp-date", "exp_date"} {
		if strings.Contains(id, kw) {
			return "payment", "cc-number", true
		}
	}

	// 3. Email
	if typeClean == "email" || strings.Contains(id, "email") || strings.Contains(id, "surel") {
		return "email", "email", false
	}

	// 4. Phone
	if typeClean == "tel" || strings.Contains(id, "phone") || strings.Contains(id, "telepon") || strings.Contains(id, "hp") ||
		strings.Contains(id, "wa") || strings.Contains(id, "whatsapp") || strings.Contains(id, "mobile_number") ||
		strings.Contains(id, "nohp") || strings.Contains(id, "no_hp") {
		return "phone", "tel", false
	}

	// 5. Name / Username
	if id == "username" || id == "user_name" || strings.Contains(id, "username") {
		return "username", "username", false
	}
	for _, kw := range [...]string{"fname", "lname", "firstname", "lastname", "first_name", "last_name", "full_name", "fullname", "nama"} {
		if strings.Contains(id, kw) {
			return "name", "name", false
		}
	}

	// 6. Address / Location
	for _, kw := range [...]string{"address", "street", "alamat", "postal", "zipcode", "zip_code", "kodepos", "city", "country", "kota", "provinsi"} {
		if strings.Contains(id, kw) {
			return "address", "street-address", false
		}
	}

	return "", "", false
}
