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
