package cls

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// hasSpreadProps memeriksa apakah elemen menggunakan JSX/TSX spread props ({...props}).
func hasSpreadProps(attrs map[string]string) bool {
	for k := range attrs {
		if strings.HasPrefix(k, "{...") {
			return true
		}
	}
	return false
}

// isImageElement mengidentifikasi apakah node merupakan elemen gambar.
func isImageElement(node *ir.Node) bool {
	if node == nil || node.Type != ir.NodeElement {
		return false
	}
	switch strings.ToLower(node.Tag) {
	case "img", "astro-image":
		return true
	}
	switch node.Tag {
	case "Image", "Picture":
		return true
	}
	return false
}

// isEmbedElement mengidentifikasi apakah node merupakan media embed atau iframe.
func isEmbedElement(node *ir.Node) bool {
	if node == nil || node.Type != ir.NodeElement {
		return false
	}
	switch strings.ToLower(node.Tag) {
	case "iframe", "video", "embed", "object":
		return true
	}
	switch node.Tag {
	case "YouTube", "ReactPlayer", "Vimeo":
		return true
	}
	return false
}

// isAdTag memeriksa apakah tag merupakan komponen iklan yang dikenal.
func isAdTag(tag string) bool {
	switch tag {
	case "AdBanner", "GoogleAd", "AdSense", "CarbonAd", "AdUnit", "AdSlot":
		return true
	default:
		return false
	}
}

// hasAdAttribute memeriksa apakah atribut node memiliki penanda slot iklan pihak ketiga.
func hasAdAttribute(attrs map[string]string) bool {
	if attrs == nil {
		return false
	}
	if _, ok := attrs["data-ad-slot"]; ok {
		return true
	}
	if _, ok := attrs["data-ad-client"]; ok {
		return true
	}
	if _, ok := attrs["data-ad-unit"]; ok {
		return true
	}
	if id, ok := attrs["id"]; ok {
		idLower := strings.ToLower(id)
		if strings.HasPrefix(idLower, "ad-") || strings.HasPrefix(idLower, "dfp-ad-") || idLower == "ad" {
			return true
		}
	}
	if slot, ok := attrs["slot"]; ok {
		if strings.ToLower(slot) == "ad" {
			return true
		}
	}
	return false
}

// hasAdClass memeriksa apakah kelas CSS mengandung penanda kontainer iklan.
func hasAdClass(classes []string) bool {
	for _, cls := range classes {
		switch strings.ToLower(cls) {
		case "ad-slot", "ad-container", "advertisement", "sponsor-banner", "ad-banner":
			return true
		}
	}
	return false
}

// isAdContainer mendeteksi apakah node merupakan kontainer penempatan iklan dinamis.
func isAdContainer(node *ir.Node) bool {
	if node == nil || node.Type != ir.NodeElement {
		return false
	}
	return isAdTag(node.Tag) || hasAdAttribute(node.Attributes) || hasAdClass(node.Classes)
}

// isCarouselTrack mendeteksi apakah node merupakan kontainer atau root track slider/carousel.
func isCarouselTrack(node *ir.Node) bool {
	if node == nil || node.Type != ir.NodeElement {
		return false
	}

	// Sinyal 1: Tag penamaan carousel
	switch node.Tag {
	case "Carousel", "Slider", "Swiper", "EmblaCarousel":
		return true
	}
	if strings.HasSuffix(node.Tag, "Carousel") || strings.HasSuffix(node.Tag, "Slider") {
		return true
	}

	// Sinyal 2: Pola scroll snap horizontal
	hasOverflowX := false
	hasSnapX := false
	for _, cls := range node.Classes {
		if cls == "overflow-x-auto" || cls == "overflow-x-scroll" {
			hasOverflowX = true
		}
		if cls == "snap-x" || cls == "snap-mandatory" || strings.HasPrefix(cls, "snap-") {
			hasSnapX = true
		}
	}

	return hasOverflowX && hasSnapX
}

// hasDimensionAttributes memeriksa keberadaan atribut 'width' dan 'height' eksplisit.
func hasDimensionAttributes(node *ir.Node) bool {
	if node == nil || node.Attributes == nil {
		return false
	}
	w, hasW := node.Attributes["width"]
	h, hasH := node.Attributes["height"]
	return hasW && hasH && w != "" && h != ""
}

// hasTailwindAspect memeriksa apakah node memiliki utilitas rasio aspek Tailwind.
func hasTailwindAspect(node *ir.Node) bool {
	if node == nil {
		return false
	}
	for _, cls := range node.Classes {
		if cls == "aspect-video" || cls == "aspect-square" {
			return true
		}
		if strings.HasPrefix(cls, "aspect-[") && strings.HasSuffix(cls, "]") {
			return true
		}
		if strings.HasPrefix(cls, "aspect-") && cls != "aspect-auto" {
			return true
		}
	}
	return false
}

// hasTailwindDimensions memeriksa apakah node memiliki dimensi tereservasi via kelas Tailwind.
func hasTailwindDimensions(node *ir.Node) bool {
	if node == nil {
		return false
	}

	if hasTailwindAspect(node) {
		return true
	}

	hasW := false
	hasH := false

	for _, cls := range node.Classes {
		// size-* mencakup width & height sekaligus
		if strings.HasPrefix(cls, "size-") && cls != "size-auto" {
			return true
		}

		if strings.HasPrefix(cls, "w-") && cls != "w-auto" && cls != "w-fit" {
			hasW = true
		}
		if strings.HasPrefix(cls, "h-") && cls != "h-auto" && cls != "h-fit" {
			hasH = true
		}
	}

	return hasW && hasH
}

func hasHeightClass(classes []string) bool {
	for _, cls := range classes {
		if strings.HasPrefix(cls, "size-") && cls != "size-auto" {
			return true
		}
		if strings.HasPrefix(cls, "h-") && cls != "h-auto" && cls != "h-fit" {
			return true
		}
		if strings.HasPrefix(cls, "min-h-") && cls != "min-h-0" {
			return true
		}
		if strings.HasPrefix(cls, "max-h-") {
			return true
		}
	}
	return false
}

func hasHeightStyle(style string) bool {
	if style == "" {
		return false
	}
	styleLower := strings.ToLower(style)
	return strings.Contains(styleLower, "aspect-ratio") ||
		strings.Contains(styleLower, "min-height") ||
		strings.Contains(styleLower, "height")
}

// hasBoundedHeight memeriksa apakah node memiliki pembatas tinggi vertikal.
func hasBoundedHeight(node *ir.Node) bool {
	if node == nil {
		return false
	}
	if hasTailwindAspect(node) || hasHeightClass(node.Classes) {
		return true
	}
	if node.Attributes != nil {
		return hasHeightStyle(node.Attributes["style"])
	}
	return false
}

// hasInlineDimensionStyle memeriksa inline style untuk aspect-ratio atau width+height.
func hasInlineDimensionStyle(node *ir.Node) bool {
	if node == nil || node.Attributes == nil {
		return false
	}
	style, ok := node.Attributes["style"]
	if !ok || style == "" {
		return false
	}
	styleLower := strings.ToLower(style)
	if strings.Contains(styleLower, "aspect-ratio") {
		return true
	}
	return strings.Contains(styleLower, "width") && strings.Contains(styleLower, "height")
}

// hasAncestorDimensionOrAspect memeriksa simpul leluhur hingga maxLevels tingkat
// apakah memiliki utilitas aspect-* atau pembatas ketinggian.
func hasAncestorDimensionOrAspect(node *ir.Node, maxLevels int) bool {
	if node == nil {
		return false
	}
	curr := node.Parent
	levels := 0
	for curr != nil && levels < maxLevels {
		if hasTailwindAspect(curr) || hasBoundedHeight(curr) || hasInlineDimensionStyle(curr) {
			return true
		}
		curr = curr.Parent
		levels++
	}
	return false
}

// hasSkeletonOrFallback memeriksa apakah node memiliki fallback skeleton awal.
func hasSkeletonOrFallback(node *ir.Node) bool {
	if node == nil {
		return false
	}
	for _, child := range node.Children {
		if child.Type == ir.NodeElement {
			if child.Tag == "Skeleton" || strings.Contains(strings.ToLower(child.Tag), "skeleton") {
				return true
			}
			for _, cls := range child.Classes {
				if cls == "skeleton" || cls == "animate-pulse" {
					return true
				}
			}
		}
	}
	return false
}

// getStyleNodeText mengekstrak seluruh teks CSS dari child node NodeText dalam <style>.
func getStyleNodeText(node *ir.Node) string {
	if node == nil {
		return ""
	}
	if len(node.Children) == 0 {
		return node.RawClasses
	}
	var sb strings.Builder
	for _, child := range node.Children {
		if child.Type == ir.NodeText {
			sb.WriteString(child.RawClasses)
		}
	}
	res := sb.String()
	if res == "" {
		return node.RawClasses
	}
	return res
}

// extractFontFaceBlocks mengekstrak blok deklarasi di dalam @font-face { ... }.
func extractFontFaceBlocks(css string) []string {
	if !strings.Contains(css, "@font-face") {
		return nil
	}
	var blocks []string
	idx := 0
	for {
		pos := strings.Index(css[idx:], "@font-face")
		if pos == -1 {
			break
		}
		start := idx + pos + len("@font-face")
		braceOpen := strings.IndexByte(css[start:], '{')
		if braceOpen == -1 {
			break
		}
		blockStart := start + braceOpen + 1
		braceClose := strings.IndexByte(css[blockStart:], '}')
		if braceClose == -1 {
			break
		}
		blockEnd := blockStart + braceClose
		blocks = append(blocks, css[blockStart:blockEnd])
		idx = blockEnd + 1
	}
	return blocks
}

// hasValidFontDisplay memeriksa apakah deklarasi @font-face memiliki descriptor font-display yang sah.
func hasValidFontDisplay(block string) bool {
	lower := strings.ToLower(block)
	pos := strings.Index(lower, "font-display")
	if pos == -1 {
		return false
	}
	after := lower[pos+len("font-display"):]
	colon := strings.IndexByte(after, ':')
	if colon == -1 {
		return false
	}
	val := after[colon+1:]
	semi := strings.IndexByte(val, ';')
	if semi != -1 {
		val = val[:semi]
	}
	val = strings.TrimSpace(val)
	switch val {
	case "swap", "optional", "fallback":
		return true
	default:
		return false
	}
}

// hasLocalFontSource memeriksa apakah blok @font-face merujuk pada font sistem fallback via src: local(...).
func hasLocalFontSource(block string) bool {
	lower := strings.ToLower(block)
	srcPos := strings.Index(lower, "src")
	if srcPos == -1 {
		return false
	}
	return strings.Contains(lower[srcPos:], "local(")
}

// hasFontMetricOverrides memeriksa apakah blok @font-face fallback menyertakan deskriptor penyesuaian metrik.
func hasFontMetricOverrides(block string) bool {
	lower := strings.ToLower(block)
	return strings.Contains(lower, "size-adjust") ||
		strings.Contains(lower, "ascent-override") ||
		strings.Contains(lower, "descent-override")
}

// findExternalFontImports mengekstrak aturan @import CSS yang mengimpor font eksternal.
func findExternalFontImports(css string) []string {
	if !strings.Contains(css, "@import") {
		return nil
	}
	var violations []string
	lines := strings.Split(css, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "@import") {
			continue
		}
		if isExternalFontImport(trimmed) {
			violations = append(violations, trimmed)
		}
	}
	return violations
}

// isExternalFontImport memeriksa apakah baris @import mengimpor font eksternal remote.
func isExternalFontImport(line string) bool {
	lower := strings.ToLower(line)

	// Whitelist: Tailwind CSS dan berkas stylesheet lokal
	if strings.Contains(lower, "tailwindcss") || strings.Contains(lower, "./") || strings.Contains(lower, "../") {
		return false
	}

	// Wajib merupakan URL eksternal remote
	if !strings.Contains(lower, "http://") && !strings.Contains(lower, "https://") && !strings.Contains(lower, "//") {
		return false
	}

	return strings.Contains(lower, "fonts.googleapis.com") ||
		strings.Contains(lower, "fonts.gstatic.com") ||
		strings.Contains(lower, "fonts.bunny.net") ||
		strings.Contains(lower, "use.typekit.net") ||
		strings.Contains(lower, "font") ||
		strings.Contains(lower, ".woff") ||
		strings.Contains(lower, ".woff2") ||
		strings.Contains(lower, ".ttf")
}

// isIconLigatureElement memeriksa apakah node merupakan elemen ligatur ikon teks.
func isIconLigatureElement(node *ir.Node) bool {
	if node == nil || node.Type != ir.NodeElement {
		return false
	}

	for _, cls := range node.Classes {
		clsLower := strings.ToLower(cls)
		if clsLower == "material-icons" ||
			clsLower == "material-icons-outlined" ||
			clsLower == "material-symbols" ||
			clsLower == "material-symbols-outlined" ||
			clsLower == "font-icon" {
			return true
		}
	}

	return false
}

// hasLockedLigatureBox memeriksa apakah kotak elemen ligatur ikon telah dikunci dimensinya.
func hasLockedLigatureBox(node *ir.Node) bool {
	if node == nil {
		return false
	}

	hasDisplay := false
	hasSize := false
	hasOverflow := false

	for _, cls := range node.Classes {
		switch cls {
		case "inline-block", "block", "inline-flex", "flex":
			hasDisplay = true
		case "overflow-hidden":
			hasOverflow = true
		}
		if strings.HasPrefix(cls, "size-") && cls != "size-auto" {
			hasSize = true
		}
	}

	if !hasSize && hasTailwindDimensions(node) {
		hasSize = true
	}

	return hasDisplay && hasSize && hasOverflow
}

// KeyframeViolation merepresentasikan pelanggaran animasi @keyframes yang memutasi properti geometri.
type KeyframeViolation struct {
	AnimationName string
	Property      string
}

// isGeometryProperty memeriksa apakah nama properti CSS termasuk properti geometri pemicu reflow tata letak CPU.
func isGeometryProperty(prop string) bool {
	prop = strings.TrimSpace(strings.ToLower(prop))
	switch prop {
	case "top", "right", "bottom", "left", "width", "height", "border-width":
		return true
	}
	if strings.HasPrefix(prop, "margin") || strings.HasPrefix(prop, "padding") || strings.HasPrefix(prop, "inset") {
		return true
	}
	if strings.HasPrefix(prop, "border-") && strings.HasSuffix(prop, "-width") {
		return true
	}
	return false
}

// findLayoutTriggerKeyframes mencari animasi @keyframes dalam blok CSS yang memutasi properti geometri.
func findLayoutTriggerKeyframes(css string) []KeyframeViolation {
	if !strings.Contains(css, "@keyframes") && !strings.Contains(css, "@-webkit-keyframes") {
		return nil
	}

	blocks := extractKeyframeBlocks(css)
	if len(blocks) == 0 {
		return nil
	}

	var violations []KeyframeViolation
	for _, block := range blocks {
		props := checkKeyframeGeometry(block.body)
		for _, prop := range props {
			violations = append(violations, KeyframeViolation{
				AnimationName: block.name,
				Property:      prop,
			})
		}
	}
	return violations
}

type keyframeBlock struct {
	name string
	body string
}

func extractKeyframeBlocks(css string) []keyframeBlock {
	var blocks []keyframeBlock
	idx := 0
	for {
		pos := strings.Index(css[idx:], "@keyframes")
		prefixLen := len("@keyframes")
		if pos == -1 {
			pos = strings.Index(css[idx:], "@-webkit-keyframes")
			prefixLen = len("@-webkit-keyframes")
		}
		if pos == -1 {
			break
		}
		start := idx + pos + prefixLen
		braceOpen := strings.IndexByte(css[start:], '{')
		if braceOpen == -1 {
			break
		}
		name := strings.TrimSpace(css[start : start+braceOpen])
		blockStart := start + braceOpen + 1

		blockEnd := findMatchingBrace(css[blockStart:])
		if blockEnd == -1 {
			break
		}
		blocks = append(blocks, keyframeBlock{
			name: name,
			body: css[blockStart : blockStart+blockEnd],
		})
		idx = blockStart + blockEnd + 1
	}
	return blocks
}

func findMatchingBrace(s string) int {
	depth := 1
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func checkKeyframeGeometry(body string) []string {
	var offending []string
	seen := make(map[string]bool)
	idx := 0
	for {
		open := strings.IndexByte(body[idx:], '{')
		if open == -1 {
			break
		}
		start := idx + open + 1
		close := strings.IndexByte(body[start:], '}')
		if close == -1 {
			break
		}
		decls := body[start : start+close]
		for _, decl := range strings.Split(decls, ";") {
			colon := strings.IndexByte(decl, ':')
			if colon == -1 {
				continue
			}
			prop := strings.TrimSpace(decl[:colon])
			if isGeometryProperty(prop) && !seen[prop] {
				seen[prop] = true
				offending = append(offending, prop)
			}
		}
		idx = start + close + 1
	}
	return offending
}

// isGeometryTailwindClass memeriksa apakah kelas Tailwind memutasi dimensi/geometri.
func isGeometryTailwindClass(cls string) bool {
	for {
		colon := strings.IndexByte(cls, ':')
		if colon == -1 {
			break
		}
		cls = cls[colon+1:]
	}
	if strings.HasPrefix(cls, "w-") || strings.HasPrefix(cls, "h-") ||
		strings.HasPrefix(cls, "p-") || strings.HasPrefix(cls, "m-") ||
		strings.HasPrefix(cls, "px-") || strings.HasPrefix(cls, "py-") ||
		strings.HasPrefix(cls, "pt-") || strings.HasPrefix(cls, "pb-") ||
		strings.HasPrefix(cls, "pl-") || strings.HasPrefix(cls, "pr-") ||
		strings.HasPrefix(cls, "mx-") || strings.HasPrefix(cls, "my-") ||
		strings.HasPrefix(cls, "mt-") || strings.HasPrefix(cls, "mb-") ||
		strings.HasPrefix(cls, "ml-") || strings.HasPrefix(cls, "mr-") ||
		strings.HasPrefix(cls, "top-") || strings.HasPrefix(cls, "bottom-") ||
		strings.HasPrefix(cls, "left-") || strings.HasPrefix(cls, "right-") ||
		strings.HasPrefix(cls, "inset-") {
		return true
	}
	return false
}

// hasPseudoGeometryMutation memeriksa apakah ada kelas yang memutasi geometri pada state interaktif.
func hasPseudoGeometryMutation(node *ir.Node) (string, bool) {
	if node == nil {
		return "", false
	}
	for _, cls := range node.Classes {
		if (strings.HasPrefix(cls, "hover:") || strings.HasPrefix(cls, "focus:") ||
			strings.HasPrefix(cls, "active:") || strings.HasPrefix(cls, "group-hover:") ||
			strings.HasPrefix(cls, "peer-hover:")) && isGeometryTailwindClass(cls) {
			return cls, true
		}
	}
	return "", false
}

// hasContainLayout memeriksa apakah elemen diisolasi dengan properti CSS contain: layout atau contain: strict.
func hasContainLayout(node *ir.Node) bool {
	if node == nil {
		return false
	}
	for _, cls := range node.Classes {
		if cls == "contain-layout" || cls == "contain-strict" {
			return true
		}
	}
	if style, ok := node.Attributes["style"]; ok {
		lower := strings.ToLower(style)
		if strings.Contains(lower, "contain: layout") || strings.Contains(lower, "contain:layout") ||
			strings.Contains(lower, "contain: strict") || strings.Contains(lower, "contain:strict") {
			return true
		}
	}
	return false
}

// hasGeometryTransition memeriksa apakah kelas atau style elemen menargetkan properti geometri dalam transisi.
func hasGeometryTransition(node *ir.Node) (string, bool) {
	if node == nil || hasContainLayout(node) {
		return "", false
	}

	hasTransitionAll := false
	for _, cls := range node.Classes {
		if cls == "transition-all" {
			hasTransitionAll = true
		}
		if strings.HasPrefix(cls, "transition-[") && strings.HasSuffix(cls, "]") {
			target := cls[len("transition-[") : len(cls)-1]
			if isGeometryProperty(target) {
				return target, true
			}
		}
	}

	if hasTransitionAll {
		if offendingClass, ok := hasPseudoGeometryMutation(node); ok {
			return offendingClass, true
		}
	}

	return "", false
}

// findCSSTransitions mencari deklarasi CSS transition pada properti geometri di dalam teks CSS.
func findCSSTransitions(css string) []string {
	if !strings.Contains(css, "transition") {
		return nil
	}
	var violations []string
	seen := make(map[string]bool)
	blocks := extractAllCSSRuleBodies(css)
	for _, body := range blocks {
		bodyLower := strings.ToLower(body)
		if strings.Contains(bodyLower, "contain: layout") || strings.Contains(bodyLower, "contain:layout") ||
			strings.Contains(bodyLower, "contain: strict") || strings.Contains(bodyLower, "contain:strict") {
			continue
		}
		props := extractGeometryTransitionsFromBody(body)
		for _, prop := range props {
			if !seen[prop] {
				seen[prop] = true
				violations = append(violations, prop)
			}
		}
	}
	return violations
}

func extractAllCSSRuleBodies(css string) []string {
	var bodies []string
	idx := 0
	for {
		open := strings.IndexByte(css[idx:], '{')
		if open == -1 {
			break
		}
		start := idx + open + 1
		close := strings.IndexByte(css[start:], '}')
		if close == -1 {
			break
		}
		bodies = append(bodies, css[start:start+close])
		idx = start + close + 1
	}
	return bodies
}

func extractGeometryTransitionsFromBody(body string) []string {
	var props []string
	for _, decl := range strings.Split(body, ";") {
		trimmed := strings.TrimSpace(decl)
		if !strings.HasPrefix(trimmed, "transition:") && !strings.HasPrefix(trimmed, "transition-property:") {
			continue
		}
		colon := strings.IndexByte(trimmed, ':')
		if colon == -1 {
			continue
		}
		val := trimmed[colon+1:]
		for _, tok := range strings.Fields(val) {
			tok = strings.Trim(tok, ",;")
			if isGeometryProperty(tok) {
				props = append(props, tok)
			}
		}
	}
	return props
}

// checkUnstableScrollbarGutter memeriksa apakah selektor root (html, body, :root) mendefinisikan overflow-y: auto tanpa scrollbar-gutter: stable.
func checkUnstableScrollbarGutter(css string) []string {
	lower := strings.ToLower(css)
	if !strings.Contains(lower, "overflow") {
		return nil
	}

	var violations []string
	rules := extractRootCSSRules(css)
	for _, r := range rules {
		bodyLower := strings.ToLower(r.body)
		hasOverflowAuto := strings.Contains(bodyLower, "overflow-y: auto") ||
			strings.Contains(bodyLower, "overflow-y:auto") ||
			strings.Contains(bodyLower, "overflow: auto") ||
			strings.Contains(bodyLower, "overflow:auto")
		hasStableGutter := strings.Contains(bodyLower, "scrollbar-gutter") ||
			strings.Contains(bodyLower, "overflow-y: scroll") ||
			strings.Contains(bodyLower, "overflow-y:scroll") ||
			strings.Contains(bodyLower, "overflow: scroll") ||
			strings.Contains(bodyLower, "overflow:scroll")

		if hasOverflowAuto && !hasStableGutter {
			violations = append(violations, r.selector)
		}
	}
	return violations
}

type cssRule struct {
	selector string
	body     string
}

func extractRootCSSRules(css string) []cssRule {
	var rules []cssRule
	idx := 0
	for {
		open := strings.IndexByte(css[idx:], '{')
		if open == -1 {
			break
		}
		selStart := idx
		selEnd := idx + open
		sel := strings.TrimSpace(css[selStart:selEnd])

		close := strings.IndexByte(css[selEnd+1:], '}')
		if close == -1 {
			break
		}
		body := css[selEnd+1 : selEnd+1+close]

		if isRootSelector(sel) {
			rules = append(rules, cssRule{selector: sel, body: body})
		}
		idx = selEnd + 1 + close + 1
	}
	return rules
}

func isRootSelector(sel string) bool {
	lower := strings.ToLower(sel)
	parts := strings.Split(lower, ",")
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "html" || trimmed == "body" || trimmed == ":root" || trimmed == "html, body" {
			return true
		}
	}
	return false
}

// isDynamicTable memeriksa apakah elemen <table> merender baris data secara dinamis.
func isDynamicTable(node *ir.Node) bool {
	if node == nil {
		return false
	}
	foundDynamic := false
	for curr := range node.Walk() {
		if curr == node {
			continue
		}
		if strings.EqualFold(curr.Tag, "tr") {
			if _, ok := curr.Attributes["key"]; ok {
				foundDynamic = true
				break
			}
		}
	}
	return foundDynamic
}

// hasTableFixed memeriksa apakah tabel memiliki class table-fixed atau style table-layout: fixed.
func hasTableFixed(node *ir.Node) bool {
	if node == nil {
		return false
	}
	for _, cls := range node.Classes {
		if cls == "table-fixed" {
			return true
		}
	}
	if style, ok := node.Attributes["style"]; ok {
		styleLower := strings.ToLower(style)
		if strings.Contains(styleLower, "table-layout: fixed") || strings.Contains(styleLower, "table-layout:fixed") {
			return true
		}
	}
	return false
}

// hasColGroupWithCols memeriksa apakah tabel memiliki child <colgroup> dengan satu atau lebih <col>.
func hasColGroupWithCols(node *ir.Node) bool {
	if node == nil {
		return false
	}
	for _, child := range node.Children {
		if strings.EqualFold(child.Tag, "colgroup") {
			return true
		}
	}
	return false
}

// hasSizedHeaders memeriksa apakah semua elemen <th> dalam <thead> tabel memiliki penentuan lebar statis.
func hasSizedHeaders(node *ir.Node) bool {
	if node == nil {
		return false
	}
	var thNodes []*ir.Node
	for curr := range node.Walk() {
		if curr != node && strings.EqualFold(curr.Tag, "th") {
			thNodes = append(thNodes, curr)
		}
	}
	if len(thNodes) == 0 {
		return false
	}
	for _, th := range thNodes {
		if !hasStaticHeaderWidth(th) {
			return false
		}
	}
	return true
}

func hasStaticHeaderWidth(th *ir.Node) bool {
	if th == nil {
		return false
	}
	if _, ok := th.Attributes["width"]; ok {
		return true
	}
	for _, cls := range th.Classes {
		if strings.HasPrefix(cls, "w-") && cls != "w-auto" {
			return true
		}
	}
	if style, ok := th.Attributes["style"]; ok {
		styleLower := strings.ToLower(style)
		if strings.Contains(styleLower, "width:") {
			return true
		}
	}
	return false
}
