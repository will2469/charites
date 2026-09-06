package lcp

import (
	"strconv"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

func cleanAttrVal(s string) string {
	return strings.Trim(strings.TrimSpace(s), `"'{}`)
}

// isLCPCandidate mengevaluasi apakah elemen media atau kontainer merupakan kandidat LCP
// menggunakan Jalur 1 (Anotasi Eksplisit) atau Jalur 2 (Multi-Signal Scoring Engine Slcp >= 60).
func isLCPCandidate(node *ir.Node) (bool, string) {
	if node == nil || node.Type != ir.NodeElement {
		return false, ""
	}

	if !isMediaTag(node.Tag) {
		return false, ""
	}

	// Abaikan elemen dengan anotasi abaikan eksplisit
	if isIgnoredLCP(node.Attributes) {
		return false, ""
	}

	// 1. Jalur 1: Standar Anotasi Eksplisit (100% Deterministic -> PROVEN)
	if role, ok := node.Attributes["data-perf-role"]; ok {
		cleanRole := cleanAttrVal(role)
		if strings.EqualFold(cleanRole, "hero") || strings.EqualFold(cleanRole, "critical") {
			return true, "PROVEN"
		}
	}
	if cand, ok := node.Attributes["data-lcp-candidate"]; ok {
		if strings.EqualFold(cleanAttrVal(cand), "true") {
			return true, "PROVEN"
		}
	}
	if strings.EqualFold(node.Tag, "HeroMedia") || strings.EqualFold(node.Tag, "HeroBanner") {
		return true, "PROVEN"
	}

	// Periksa apakah leluhur langsung memiliki data-perf-role="hero"
	if isAncestorHero(node) {
		return true, "PROVEN"
	}

	// 2. Jalur 2: Multi-Signal Candidate Scoring Engine (Slcp >= 60)
	if !isMediaTag(node.Tag) {
		return false, ""
	}

	score := calculateCandidateScore(node)
	if score >= 60 {
		return true, "LIKELY"
	}

	return false, ""
}

func isIgnoredLCP(attrs map[string]string) bool {
	if attrs == nil {
		return false
	}
	val, ok := attrs["data-lcp-ignore"]
	return ok && strings.EqualFold(val, "true")
}

func isMediaTag(tag string) bool {
	switch {
	case strings.EqualFold(tag, "img"),
		strings.EqualFold(tag, "picture"),
		strings.EqualFold(tag, "Image"),
		strings.EqualFold(tag, "video"):
		return true
	default:
		return false
	}
}

func isAncestorHero(node *ir.Node) bool {
	cur := node.Parent
	for depth := 0; cur != nil && depth < 3; depth++ {
		if cur.Attributes != nil {
			if role, ok := cur.Attributes["data-perf-role"]; ok && strings.EqualFold(cleanAttrVal(role), "hero") {
				return true
			}
		}
		cur = cur.Parent
	}
	return false
}

func calculateCandidateScore(node *ir.Node) int {
	score := 0

	// Sinyal 1: Topologi Alur Dokumen (+35)
	if node.Span.Line > 0 && node.Span.Line <= 50 {
		score += 35
	}

	// Sinyal 2: Kontainer Semantik (+25)
	if hasSemanticContainer(node) {
		score += 25
	}

	// Sinyal 3: Dimensi Geometris Statis (+20)
	if hasLargeDimensions(node) {
		score += 20
	}

	// Penalti: Elemen Tersembunyi atau di Footer (-50)
	if isHiddenOrFooter(node) {
		score -= 50
	}

	return score
}

func hasSemanticContainer(node *ir.Node) bool {
	if containsHeroKeywords(node.RawClasses) {
		return true
	}

	cur := node.Parent
	for depth := 0; cur != nil && depth < 4; depth++ {
		if strings.EqualFold(cur.Tag, "header") {
			return true
		}
		if cur.Attributes != nil {
			if role, ok := cur.Attributes["role"]; ok && strings.EqualFold(cleanAttrVal(role), "banner") {
				return true
			}
			if id, ok := cur.Attributes["id"]; ok && containsHeroKeywords(cleanAttrVal(id)) {
				return true
			}
		}
		if containsHeroKeywords(cur.RawClasses) {
			return true
		}
		cur = cur.Parent
	}
	return false
}

func containsHeroKeywords(s string) bool {
	if len(s) == 0 {
		return false
	}
	lower := strings.ToLower(s)
	return strings.Contains(lower, "hero") ||
		strings.Contains(lower, "banner") ||
		strings.Contains(lower, "masthead")
}

func hasLargeDimensions(node *ir.Node) bool {
	raw := node.RawClasses
	if strings.Contains(raw, "w-full") ||
		strings.Contains(raw, "aspect-video") ||
		strings.Contains(raw, "aspect-[16/9]") ||
		strings.Contains(raw, "max-w-") ||
		strings.Contains(raw, "object-cover") {
		return true
	}

	if node.Attributes != nil {
		if wStr, ok := node.Attributes["width"]; ok {
			cleanW := cleanAttrVal(wStr)
			if w, err := strconv.Atoi(cleanW); err == nil && w >= 600 {
				return true
			}
		}
	}
	return false
}

func isHiddenOrFooter(node *ir.Node) bool {
	if strings.Contains(node.RawClasses, "hidden") || strings.Contains(node.RawClasses, "invisible") {
		return true
	}

	cur := node.Parent
	for cur != nil {
		if strings.EqualFold(cur.Tag, "footer") {
			return true
		}
		if cur.Attributes != nil {
			if role, ok := cur.Attributes["role"]; ok && strings.EqualFold(cleanAttrVal(role), "dialog") {
				return true
			}
		}
		cur = cur.Parent
	}
	return false
}

// isLazyLoaded memeriksa apakah atribut loading bernilai "lazy".
func isLazyLoaded(attrs map[string]string) bool {
	if attrs == nil {
		return false
	}
	val, ok := attrs["loading"]
	return ok && strings.EqualFold(cleanAttrVal(val), "lazy")
}

// hasHighFetchPriority memeriksa apakah atribut fetchpriority bernilai "high".
func hasHighFetchPriority(attrs map[string]string) bool {
	if attrs == nil {
		return false
	}
	if val, ok := attrs["fetchpriority"]; ok && strings.EqualFold(cleanAttrVal(val), "high") {
		return true
	}
	if val, ok := attrs["fetchPriority"]; ok && strings.EqualFold(cleanAttrVal(val), "high") {
		return true
	}
	return false
}

// hasCSSBackgroundImage mendeteksi penggunaan gambar latar via Tailwind bg-[url(...)] atau style inline.
func hasCSSBackgroundImage(rawClasses string, attrs map[string]string) (string, bool) {
	if url, ok := extractTailwindBgUrl(rawClasses); ok {
		return url, true
	}
	return extractInlineStyleBgUrl(attrs)
}

func extractTailwindBgUrl(rawClasses string) (string, bool) {
	if len(rawClasses) == 0 || !strings.Contains(rawClasses, "bg-[url(") {
		return "", false
	}
	idx := strings.Index(rawClasses, "bg-[url(")
	start := idx + len("bg-[url(")
	end := strings.Index(rawClasses[start:], ")]")
	if end == -1 {
		return "", false
	}
	url := strings.Trim(rawClasses[start:start+end], `'"`)
	if isDecorativePattern(url) {
		return "", false
	}
	return url, true
}

func extractInlineStyleBgUrl(attrs map[string]string) (string, bool) {
	if attrs == nil {
		return "", false
	}
	styleVal, ok := attrs["style"]
	if !ok {
		return "", false
	}
	lower := strings.ToLower(styleVal)
	if !strings.Contains(lower, "background-image") || !strings.Contains(lower, "url(") {
		return "", false
	}
	idx := strings.Index(lower, "url(")
	start := idx + len("url(")
	end := strings.Index(lower[start:], ")")
	if end == -1 {
		return "", false
	}
	url := strings.Trim(styleVal[start:start+end], `'"`)
	if isDecorativePattern(url) {
		return "", false
	}
	return url, true
}

func isDecorativePattern(url string) bool {
	lower := strings.ToLower(url)
	return strings.Contains(lower, "pattern") ||
		strings.Contains(lower, "texture") ||
		strings.Contains(lower, "radial-gradient") ||
		strings.Contains(lower, "linear-gradient") ||
		strings.Contains(lower, "dots.svg") ||
		strings.Contains(lower, "grid.svg")
}

// hasPreloadInHead memeriksa apakah pohon dokumen telah memiliki <link rel="preload" as="image">
// yang sesuai di dalam elemen <head> atau komentar sumber.
func hasPreloadInHead(node *ir.Node, assetURL string) bool {
	root := findDocumentRoot(node)
	if hasCommentPreload(root, assetURL) {
		return true
	}
	return hasLinkElementPreload(root, assetURL)
}

func findDocumentRoot(node *ir.Node) *ir.Node {
	cur := node
	for cur.Parent != nil {
		cur = cur.Parent
	}
	return cur
}

func hasCommentPreload(root *ir.Node, assetURL string) bool {
	cleanURL := cleanAttrVal(assetURL)
	for _, ch := range root.Children {
		if ch != nil && ch.Type == ir.NodeComment && len(ch.RawClasses) > 0 {
			src := ch.RawClasses
			if strings.Contains(src, `rel="preload"`) && strings.Contains(src, `as="image"`) {
				if cleanURL == "" || strings.Contains(src, cleanURL) {
					return true
				}
			}
		}
	}
	return false
}

func hasLinkElementPreload(root *ir.Node, assetURL string) bool {
	cleanURL := cleanAttrVal(assetURL)
	for n := range root.Walk() {
		if isPreloadLinkNode(n, cleanURL) {
			return true
		}
	}
	return false
}

func isPreloadLinkNode(n *ir.Node, cleanURL string) bool {
	if n.Type != ir.NodeElement || !strings.EqualFold(n.Tag, "link") || n.Attributes == nil {
		return false
	}
	rel := cleanAttrVal(n.Attributes["rel"])
	asVal := cleanAttrVal(n.Attributes["as"])
	if !strings.EqualFold(rel, "preload") || !strings.EqualFold(asVal, "image") {
		return false
	}
	href := cleanAttrVal(n.Attributes["href"])
	return cleanURL == "" || strings.Contains(href, cleanURL)
}

// isDelayedDiscoveryLCP mendeteksi elemen LCP yang terindikasi mengalami keterlambatan penemuan (delayed discovery).
func isDelayedDiscoveryLCP(node *ir.Node) (string, bool) {
	if node == nil || node.Type != ir.NodeElement {
		return "", false
	}

	if isStaticHTMLImg(node) {
		return "", false
	}

	if !isCandidateHeroContainer(node) {
		return "", false
	}

	if src, ok := getDynamicDataSrc(node.Attributes); ok {
		return src, true
	}

	return hasCSSBackgroundImage(node.RawClasses, node.Attributes)
}

func isStaticHTMLImg(node *ir.Node) bool {
	if !strings.EqualFold(node.Tag, "img") || node.Attributes == nil {
		return false
	}
	_, ok := node.Attributes["src"]
	return ok
}

func getDynamicDataSrc(attrs map[string]string) (string, bool) {
	if attrs == nil {
		return "", false
	}
	for _, key := range [...]string{"data-bg-src", "data-src"} {
		if raw, ok := attrs[key]; ok {
			src := cleanAttrVal(raw)
			if len(src) > 0 {
				return src, true
			}
		}
	}
	return "", false
}

func isCandidateHeroContainer(node *ir.Node) bool {
	if node.Attributes != nil {
		if role, ok := node.Attributes["data-perf-role"]; ok && strings.EqualFold(cleanAttrVal(role), "hero") {
			return true
		}
	}
	return hasSemanticContainer(node)
}

// isFluidImage memeriksa apakah gambar memiliki ukuran fluida responsif
// (lebar mengikuti kontainer/viewport, misal w-full, max-w-*, aspect-*, atau tanpa atribut width tetap).
func isFluidImage(node *ir.Node) bool {
	if node == nil {
		return false
	}
	// Jika secara eksplisit berdimensi tetap dan bukan w-full, maka bukan fluida
	if _, _, fixed := isFixedDimensionImage(node); fixed {
		return false
	}
	return true
}

// isFixedDimensionImage mendeteksi apakah elemen memiliki dimensi piksel tetap (width/height atribut atau fixed Tailwind classes).
func isFixedDimensionImage(node *ir.Node) (int, int, bool) {
	if node == nil {
		return 0, 0, false
	}

	// Jika memiliki utility fluid w-full atau w-screen, CSS mengambil preseden menjadikan gambar fluida
	raw := strings.ToLower(node.RawClasses)
	if strings.Contains(raw, "w-full") || strings.Contains(raw, "w-screen") {
		return 0, 0, false
	}

	if node.Attributes != nil {
		rawW, hasW := node.Attributes["width"]
		rawH, hasH := node.Attributes["height"]
		if hasW {
			wStr := cleanAttrVal(rawW)
			wStr = strings.TrimSuffix(wStr, "px")
			if w, err := strconv.Atoi(wStr); err == nil && w > 0 {
				h := 0
				if hasH {
					hStr := cleanAttrVal(rawH)
					hStr = strings.TrimSuffix(hStr, "px")
					if val, err2 := strconv.Atoi(hStr); err2 == nil {
						h = val
					}
				}
				return w, h, true
			}
		}
	}

	return parseFixedTailwindDims(node.Classes)
}

func parseFixedTailwindDims(classes []string) (int, int, bool) {
	w, h := 0, 0
	for _, cls := range classes {
		if strings.HasPrefix(cls, "w-[") && strings.HasSuffix(cls, "px]") {
			inner := cls[3 : len(cls)-3]
			if val, err := strconv.Atoi(inner); err == nil {
				w = val
			}
		} else if strings.HasPrefix(cls, "h-[") && strings.HasSuffix(cls, "px]") {
			inner := cls[3 : len(cls)-3]
			if val, err := strconv.Atoi(inner); err == nil {
				h = val
			}
		}
	}
	if w > 0 {
		return w, h, true
	}
	return 0, 0, false
}

// getSrcsetAttribute mengambil nilai atribut srcset atau srcSet (varian JSX).
func getSrcsetAttribute(attrs map[string]string) (string, bool) {
	if attrs == nil {
		return "", false
	}
	if val, ok := attrs["srcset"]; ok {
		return val, true
	}
	if val, ok := attrs["srcSet"]; ok {
		return val, true
	}
	return "", false
}

// hasResponsiveSrcset memeriksa apakah media mendefinisikan varian responsif (srcset dengan width descriptors atau sizes).
func hasResponsiveSrcset(attrs map[string]string) bool {
	if attrs == nil {
		return false
	}
	rawSrcset, okSrcset := getSrcsetAttribute(attrs)
	if !okSrcset || len(cleanAttrVal(rawSrcset)) == 0 {
		return false
	}
	cleaned := cleanAttrVal(rawSrcset)
	hasWidthDesc := strings.Contains(cleaned, "w,") || strings.HasSuffix(cleaned, "w")
	_, hasSizes := attrs["sizes"]
	return hasWidthDesc || hasSizes
}

// isInsidePictureWithResponsiveSource memeriksa apakah <img> berada dalam <picture> dengan <source srcset="...">.
func isInsidePictureWithResponsiveSource(node *ir.Node) bool {
	if node == nil || node.Parent == nil {
		return false
	}
	if !strings.EqualFold(node.Parent.Tag, "picture") {
		return false
	}
	for _, child := range node.Parent.Children {
		if child != nil && strings.EqualFold(child.Tag, "source") && child.Attributes != nil {
			if ss, ok := getSrcsetAttribute(child.Attributes); ok && len(cleanAttrVal(ss)) > 0 {
				return true
			}
		}
	}
	return false
}

// isSVGFile memeriksa apakah sumber gambar merupakan berkas SVG.
func isSVGFile(src string) bool {
	lower := strings.ToLower(cleanAttrVal(src))
	if idx := strings.IndexAny(lower, "?#"); idx != -1 {
		lower = lower[:idx]
	}
	return strings.HasSuffix(lower, ".svg")
}

// isKnownImageCDN memeriksa apakah URL berasal dari CDN gambar modern yang mengotomatisasi negosiasi konten WebP/AVIF.
func isKnownImageCDN(src string) bool {
	lower := strings.ToLower(src)
	cdnHosts := [...]string{
		"cloudinary.com",
		"unsplash.com",
		"imgix.net",
		"cloudflare.com",
		"cdn.sanity.io",
		"ctfassets.net",
		"imagekit.io",
		"shopify.com",
	}
	for _, host := range cdnHosts {
		if strings.Contains(lower, host) {
			return true
		}
	}
	return false
}

// isPictureWithModernSource memeriksa apakah <img> dibungkus dalam <picture> yang memiliki source type webp/avif.
func isPictureWithModernSource(node *ir.Node) bool {
	if node == nil || node.Parent == nil {
		return false
	}
	if !strings.EqualFold(node.Parent.Tag, "picture") {
		return false
	}
	for _, child := range node.Parent.Children {
		if child != nil && strings.EqualFold(child.Tag, "source") && child.Attributes != nil {
			if t, ok := child.Attributes["type"]; ok {
				cleanType := strings.ToLower(cleanAttrVal(t))
				if strings.Contains(cleanType, "webp") || strings.Contains(cleanType, "avif") {
					return true
				}
			}
		}
	}
	return false
}

// isHeavyLegacyRasterFormat memeriksa apakah berkas gambar menggunakan format raster kuno (.png, .bmp, .tiff, .gif) tanpa WebP/AVIF.
func isHeavyLegacyRasterFormat(rawSrc string) (string, bool) {
	cleaned := cleanAttrVal(rawSrc)
	if len(cleaned) == 0 {
		return "", false
	}
	lower := strings.ToLower(cleaned)
	if idx := strings.IndexAny(lower, "?#"); idx != -1 {
		lower = lower[:idx]
	}
	if isKnownImageCDN(lower) {
		return "", false
	}
	for _, ext := range [...]string{".png", ".bmp", ".tiff", ".tif", ".gif"} {
		if strings.HasSuffix(lower, ext) {
			return ext, true
		}
	}
	return "", false
}

// hasDensityDescriptors memeriksa apakah srcset menyertakan deskriptor 1x, 2x, atau 3x.
func hasDensityDescriptors(srcset string) bool {
	cleaned := cleanAttrVal(srcset)
	if len(cleaned) == 0 {
		return false
	}
	parts := strings.Split(cleaned, ",")
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if strings.HasSuffix(trimmed, "1x") || strings.HasSuffix(trimmed, "2x") || strings.HasSuffix(trimmed, "3x") {
			return true
		}
	}
	return false
}

// hasClientOnlyDirective memeriksa apakah atribut memiliki direktif client:only pada Astro.
func hasClientOnlyDirective(attrs map[string]string) bool {
	if attrs == nil {
		return false
	}
	for k := range attrs {
		if k == "client:only" || strings.HasPrefix(k, "client:only") {
			return true
		}
	}
	return false
}

// isHeroIsland memeriksa apakah pulau komponen berada pada area hero pelipatan atas.
func isHeroIsland(node *ir.Node) bool {
	if node == nil {
		return false
	}
	lowerTag := strings.ToLower(node.Tag)
	if strings.Contains(lowerTag, "hero") {
		return true
	}
	if node.Attributes != nil {
		if role, ok := node.Attributes["data-perf-role"]; ok && strings.EqualFold(cleanAttrVal(role), "hero") {
			return true
		}
	}
	if strings.Contains(strings.ToLower(node.RawClasses), "hero") {
		return true
	}
	if isAncestorHero(node) {
		return true
	}
	if node.Parent != nil {
		pTag := strings.ToLower(node.Parent.Tag)
		if pTag == "header" || pTag == "main" {
			return true
		}
	}
	return false
}

// hasFallbackSlot memeriksa apakah ada child node yang memiliki atribut slot="fallback".
func hasFallbackSlot(node *ir.Node) bool {
	if node == nil {
		return false
	}
	for _, child := range node.Children {
		if child != nil && child.Attributes != nil {
			if slot, ok := child.Attributes["slot"]; ok {
				if strings.EqualFold(cleanAttrVal(slot), "fallback") {
					return true
				}
			}
		}
	}
	return false
}
