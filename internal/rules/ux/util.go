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

// hasFeedforwardExplanation memeriksa apakah kontrol interaktif memiliki petunjuk feedforward saat terkunci.
func hasFeedforwardExplanation(node *ir.Node) bool {
	if node == nil {
		return false
	}
	if hasDirectFeedforwardAttributes(node) || isInsideTooltipWrapper(node) {
		return true
	}
	return hasFeedforwardSibling(node)
}

func hasDirectFeedforwardAttributes(node *ir.Node) bool {
	if desc, ok := getAttrCaseInsensitive(node, "aria-describedby", "aria-description"); ok && cleanAttrValue(desc) != "" {
		return true
	}
	if title, ok := getAttrCaseInsensitive(node, "title", "tooltip", "data-tooltip"); ok && cleanAttrValue(title) != "" {
		return true
	}
	return false
}

func isInsideTooltipWrapper(node *ir.Node) bool {
	curr := node.Parent
	for curr != nil {
		tagLower := strings.ToLower(curr.Tag)
		if strings.Contains(tagLower, "tooltip") || strings.Contains(tagLower, "popover") {
			return true
		}
		curr = curr.Parent
	}
	return false
}

func hasFeedforwardSibling(node *ir.Node) bool {
	if node.Parent == nil {
		return false
	}
	describedBy, _ := getAttrCaseInsensitive(node, "aria-describedby")
	describedByClean := cleanAttrValue(describedBy)

	for _, sib := range node.Parent.Children {
		if sib == node || sib.Type != ir.NodeElement {
			continue
		}
		if isExplanationSibling(sib, describedByClean) {
			return true
		}
	}
	return false
}

func isExplanationSibling(sib *ir.Node, describedBy string) bool {
	if describedBy != "" {
		if id, ok := getAttrCaseInsensitive(sib, "id"); ok && cleanAttrValue(id) == describedBy {
			return true
		}
	}
	if role, ok := getAttrCaseInsensitive(sib, "role"); ok {
		r := cleanAttrValue(role)
		if r == "status" || r == "alert" || r == "tooltip" {
			return true
		}
	}
	for _, cls := range sib.Classes {
		base := strings.ToLower(StripVariantsOnlyBase(cls))
		if strings.Contains(base, "hint") || strings.Contains(base, "helper") || strings.Contains(base, "desc") ||
			strings.Contains(base, "feedback") || strings.Contains(base, "muted") || strings.Contains(base, "alert") {
			return true
		}
	}
	return false
}

// hasErrorPresentationInSubtree memeriksa apakah pohon komponen menyediakan elemen presentasi error.
func hasErrorPresentationInSubtree(node *ir.Node) bool {
	if node == nil {
		return false
	}
	for n := range node.Walk() {
		if isErrorDisplayNode(n) || hasErrorTextOrBindingNode(n) {
			return true
		}
	}
	return false
}

func hasErrorTextOrBindingNode(n *ir.Node) bool {
	if n.Type == ir.NodeText {
		txt := strings.ToLower(n.RawClasses)
		if strings.Contains(txt, "error") || strings.Contains(txt, "kesalahan") || strings.Contains(txt, "gagal") {
			return true
		}
	}
	for attrName, attrVal := range n.Attributes {
		if isEventHandlerOrActionAttr(attrName) {
			continue
		}
		nameLower := strings.ToLower(attrName)
		valLower := strings.ToLower(attrVal)
		if strings.Contains(nameLower, "error") || strings.Contains(nameLower, "errormessage") {
			return true
		}
		if strings.Contains(valLower, "error") || strings.Contains(valLower, "iserror") || strings.Contains(valLower, "haserror") {
			return true
		}
	}
	return false
}

func isEventHandlerOrActionAttr(name string) bool {
	lower := strings.ToLower(name)
	return strings.HasPrefix(lower, "on") || lower == "action" || lower == "validate" || lower == "onsubmit"
}

func isErrorDisplayNode(n *ir.Node) bool {
	if n.Type != ir.NodeElement {
		return false
	}
	if role, ok := getAttrCaseInsensitive(n, "role"); ok {
		r := cleanAttrValue(role)
		if r == "alert" || r == "status" {
			return true
		}
	}
	if _, ok := getAttrCaseInsensitive(n, "aria-invalid", "aria-errormessage", "error", "errormessage", "iserror", "isinvalid"); ok {
		return true
	}
	tagLower := strings.ToLower(n.Tag)
	if strings.Contains(tagLower, "error") || tagLower == "alert" || tagLower == "formerror" {
		return true
	}
	for _, cls := range n.Classes {
		base := strings.ToLower(StripVariantsOnlyBase(cls))
		if strings.Contains(base, "error") || strings.Contains(base, "destructive") || strings.Contains(base, "invalid") {
			return true
		}
	}
	return false
}

// hasEmptyStateInSubtree memeriksa apakah ada penanganan cabang state kosong di dalam subtree atau sibling.
func hasEmptyStateInSubtree(node *ir.Node) bool {
	if node == nil {
		return false
	}
	if hasEmptyStateAttributes(node) {
		return true
	}
	for n := range node.Walk() {
		if isEmptyStateIndicatorNode(n) {
			return true
		}
	}
	if node.Parent != nil {
		for _, sib := range node.Parent.Children {
			if sib != node && isEmptyStateIndicatorNode(sib) {
				return true
			}
		}
	}
	return false
}

func hasEmptyStateAttributes(node *ir.Node) bool {
	for k, v := range node.Attributes {
		kLower := strings.ToLower(k)
		if strings.Contains(kLower, "empty") || strings.Contains(kLower, "fallback") {
			return true
		}
		vLower := strings.ToLower(v)
		if strings.Contains(vLower, "length === 0") || strings.Contains(vLower, "length == 0") ||
			strings.Contains(vLower, "length > 0") || strings.Contains(vLower, "isempty") {
			return true
		}
	}
	return false
}

func isEmptyStateIndicatorNode(n *ir.Node) bool {
	if n.Type == ir.NodeElement {
		tagLower := strings.ToLower(n.Tag)
		if strings.Contains(tagLower, "empty") || strings.Contains(tagLower, "nodata") ||
			strings.Contains(tagLower, "noresult") || tagLower == "fallback" {
			return true
		}
		for _, cls := range n.Classes {
			base := strings.ToLower(StripVariantsOnlyBase(cls))
			if strings.Contains(base, "empty") || strings.Contains(base, "zero-state") || strings.Contains(base, "no-data") {
				return true
			}
		}
	}
	if n.Type == ir.NodeText {
		txt := strings.ToLower(n.RawClasses)
		if strings.Contains(txt, "belum ada") || strings.Contains(txt, "tidak ada") ||
			strings.Contains(txt, "kosong") || strings.Contains(txt, "no data") || strings.Contains(txt, "no items") {
			return true
		}
	}
	return false
}

// detectUnthrottledNetworkCall memeriksa apakah handler event memicu network API langsung tanpa debounce/throttle.
func detectUnthrottledNetworkCall(handler string) (string, bool) {
	lower := strings.ToLower(cleanAttrValue(handler))
	if lower == "" {
		return "", false
	}

	if isThrottledOrDebounced(lower) {
		return "", false
	}

	return findDirectNetworkCall(lower)
}

func isThrottledOrDebounced(handler string) bool {
	keywords := [...]string{"debounce", "throttle", "debounced", "throttled", "timeout"}
	for _, kw := range keywords {
		if strings.Contains(handler, kw) {
			return true
		}
	}
	return false
}

func findDirectNetworkCall(handler string) (string, bool) {
	networkCalls := [...]string{
		"fetch(", "axios.", "http.", "api.get(", "api.post(", "api.put(", "api.delete(",
		"query(", "fetchsuggestions(", "searchapi(", "request(",
	}
	for _, call := range networkCalls {
		if strings.Contains(handler, call) {
			return call, true
		}
	}
	return "", false
}

// detectAsyncMutationHandler memeriksa apakah handler memicu mutasi asinkron.
func detectAsyncMutationHandler(handler string) bool {
	lower := strings.ToLower(cleanAttrValue(handler))
	if lower == "" {
		return false
	}
	if strings.Contains(lower, "api.post") || strings.Contains(lower, "api.put") ||
		strings.Contains(lower, "api.delete") || strings.Contains(lower, "api.patch") ||
		strings.Contains(lower, "fetch(") || strings.Contains(lower, "mutate(") ||
		strings.Contains(lower, "mutation(") || strings.Contains(lower, "await ") {
		return true
	}
	return false
}

// hasReentryGuard memeriksa keberadaan guard kunci interaksi (R1).
func hasReentryGuard(node *ir.Node) bool {
	if node == nil {
		return false
	}
	for _, attr := range [...]string{"disabled", "aria-disabled"} {
		if val, ok := getAttrCaseInsensitive(node, attr); ok {
			lower := strings.ToLower(cleanAttrValue(val))
			if lower != "false" && lower != "{false}" {
				return true
			}
		}
	}
	return false
}

// hasPerceivableFeedback memeriksa keberadaan umpan balik visual bahwa proses sedang berjalan (R2).
func hasPerceivableFeedback(node *ir.Node) bool {
	if node == nil {
		return false
	}
	if val, ok := getAttrCaseInsensitive(node, "aria-busy"); ok {
		lower := strings.ToLower(cleanAttrValue(val))
		if lower != "false" && lower != "{false}" {
			return true
		}
	}
	for n := range node.Walk() {
		if isPendingIndicatorNode(n) {
			return true
		}
	}
	return false
}

func isPendingIndicatorNode(n *ir.Node) bool {
	if n.Type == ir.NodeElement {
		tagLower := strings.ToLower(n.Tag)
		if strings.Contains(tagLower, "spinner") || strings.Contains(tagLower, "loader") ||
			strings.Contains(tagLower, "loading") {
			return true
		}
		for _, cls := range n.Classes {
			base := strings.ToLower(StripVariantsOnlyBase(cls))
			if strings.Contains(base, "spin") || strings.Contains(base, "loading") {
				return true
			}
		}
	}
	if n.Type == ir.NodeText {
		txt := strings.ToLower(n.RawClasses)
		if strings.Contains(txt, "memproses") || strings.Contains(txt, "loading") ||
			strings.Contains(txt, "menyimpan") || strings.Contains(txt, "submitting") {
			return true
		}
	}
	return false
}

// detectDestructiveAction memeriksa apakah elemen memicu mutasi atau aksi destruktif.
func detectDestructiveAction(node *ir.Node) (string, bool) {
	if node == nil {
		return "", false
	}
	for attrName, attrVal := range node.Attributes {
		if !isEventHandlerOrActionAttr(attrName) {
			continue
		}
		if act, ok := matchDestructiveKeyword(attrVal); ok {
			return act, true
		}
	}
	if isDestructiveStyledElement(node) {
		for _, child := range node.Children {
			if child.Type == ir.NodeText {
				if act, ok := matchDestructiveKeyword(child.RawClasses); ok {
					return act, true
				}
			}
		}
	}
	return "", false
}

func matchDestructiveKeyword(text string) (string, bool) {
	lower := strings.ToLower(text)
	keywords := [...]string{"delete", "remove", "destroy", "purge", "revoke", "hapus"}
	for _, kw := range keywords {
		if strings.Contains(lower, kw) {
			return kw, true
		}
	}
	return "", false
}

func isDestructiveStyledElement(node *ir.Node) bool {
	for _, cls := range node.Classes {
		base := strings.ToLower(StripVariantsOnlyBase(cls))
		if strings.Contains(base, "destructive") || strings.Contains(base, "danger") ||
			strings.HasPrefix(base, "bg-red-") || strings.HasPrefix(base, "text-red-") {
			return true
		}
	}
	return false
}

// hasConfirmationGating memeriksa apakah aksi destruktif dilindungi dialog atau konfirmasi.
func hasConfirmationGating(node *ir.Node, handler string) bool {
	if node == nil {
		return false
	}
	lowerHandler := strings.ToLower(cleanAttrValue(handler))
	if strings.Contains(lowerHandler, "confirm(") || strings.Contains(lowerHandler, "window.confirm(") {
		return true
	}
	if hasConfirmationAttribute(node) {
		return true
	}
	return isInsideConfirmationWrapper(node)
}

func hasConfirmationAttribute(node *ir.Node) bool {
	for k, v := range node.Attributes {
		kLower := strings.ToLower(k)
		if strings.Contains(kLower, "confirm") {
			return true
		}
		vLower := strings.ToLower(v)
		if strings.Contains(vLower, "confirm") || strings.Contains(vLower, "step === 2") {
			return true
		}
	}
	return false
}

func isInsideConfirmationWrapper(node *ir.Node) bool {
	curr := node.Parent
	for curr != nil {
		tagLower := strings.ToLower(curr.Tag)
		if strings.Contains(tagLower, "confirm") || strings.Contains(tagLower, "dialog") ||
			strings.Contains(tagLower, "modal") || strings.Contains(tagLower, "alertdialog") {
			return true
		}
		curr = curr.Parent
	}
	return false
}

// detectUnboundedAsyncLoading memeriksa apakah handler mengaktifkan loading sebelum await tanpa me-reset di exit path.
func detectUnboundedAsyncLoading(handler string) bool {
	lower := strings.ToLower(cleanAttrValue(handler))
	if !hasAsyncLoadingStart(lower) {
		return false
	}
	if !strings.Contains(lower, "await ") {
		return false
	}
	if strings.Contains(lower, "finally") {
		return false
	}
	return !hasAsyncLoadingEnd(lower)
}

func hasAsyncLoadingStart(lower string) bool {
	starts := [...]string{
		"setloading(true)", "setisloading(true)", "setpending(true)",
		"setispending(true)", "setsubmitting(true)",
	}
	for _, s := range starts {
		if strings.Contains(lower, s) {
			return true
		}
	}
	return false
}

func hasAsyncLoadingEnd(lower string) bool {
	ends := [...]string{
		"setloading(false)", "setisloading(false)", "setpending(false)",
		"setispending(false)", "setsubmitting(false)",
	}
	for _, e := range ends {
		if strings.Contains(lower, e) {
			return true
		}
	}
	return false
}

// detectSilentCatchSwallow memeriksa apakah blok catch menelan error tanpa memberi umpan balik atau re-throw.
func detectSilentCatchSwallow(handler string) bool {
	lower := strings.ToLower(cleanAttrValue(handler))
	catchIdx := strings.Index(lower, "catch")
	if catchIdx == -1 {
		return false
	}
	catchBody := lower[catchIdx:]

	if hasErrorFeedbackInCatch(catchBody) || strings.Contains(catchBody, "throw ") || strings.Contains(catchBody, "throw;") {
		return false
	}

	return hasSwallowPatternInCatch(catchBody)
}

func hasErrorFeedbackInCatch(catchBody string) bool {
	feedbackKeywords := [...]string{
		"toast.", "toast(", "seterror(", "seterr(", "alert(",
		"banner(", "notification(", "notify(", "reporterror(",
		"message.error(", "dispatch(",
	}
	for _, kw := range feedbackKeywords {
		if strings.Contains(catchBody, kw) {
			return true
		}
	}
	return false
}

func hasSwallowPatternInCatch(catchBody string) bool {
	if strings.Contains(catchBody, "console.log(") || strings.Contains(catchBody, "console.error(") ||
		strings.Contains(catchBody, "console.warn(") {
		return true
	}
	stripped := strings.ReplaceAll(catchBody, " ", "")
	stripped = strings.ReplaceAll(stripped, "\t", "")
	stripped = strings.ReplaceAll(stripped, "\n", "")
	stripped = strings.ReplaceAll(stripped, "\r", "")
	return strings.Contains(stripped, "{}") || strings.Contains(stripped, "{;}")
}
