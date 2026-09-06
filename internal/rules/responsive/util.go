package responsive

import (
	"strconv"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

var breakpointPrefixes = [...]string{
	"sm:", "md:", "lg:", "xl:", "2xl:",
	"max-sm:", "max-md:", "max-lg:", "max-xl:", "max-2xl:",
}

func hasBreakpointPrefix(cls string) bool {
	for _, p := range breakpointPrefixes {
		if strings.HasPrefix(cls, p) {
			return true
		}
	}
	if strings.HasPrefix(cls, "@") && strings.Contains(cls, ":") {
		return true
	}
	return false
}

func isMultiColGrid(cls string) bool {
	if !strings.HasPrefix(cls, "grid-cols-") {
		return false
	}
	suffix := cls[len("grid-cols-"):]
	switch suffix {
	case "3", "4", "5", "6", "7", "8", "9", "10", "11", "12":
		return true
	default:
		return false
	}
}

func isGiantFont(cls string) bool {
	switch cls {
	case "text-5xl", "text-6xl", "text-7xl", "text-8xl", "text-9xl":
		return true
	default:
		return false
	}
}

func hasScrollContainerAncestor(node *ir.Node) bool {
	if node == nil {
		return false
	}
	curr := node.Parent
	for curr != nil {
		for _, cls := range curr.Classes {
			switch cls {
			case "overflow-x-auto", "overflow-x-scroll", "overflow-auto", "overflow-scroll":
				return true
			}
		}
		curr = curr.Parent
	}
	return false
}

func hasResponsiveTableDisplay(node *ir.Node) bool {
	if node == nil {
		return false
	}
	hasHiddenMobile := false
	hasDesktopTable := false
	for _, cls := range node.Classes {
		if cls == "hidden" {
			hasHiddenMobile = true
		}
		if cls == "block" || cls == "flex" || cls == "grid" {
			return true
		}
		if cls == "md:table" || cls == "sm:table" || cls == "lg:table" {
			hasDesktopTable = true
		}
	}
	return hasHiddenMobile && hasDesktopTable
}

func extractStaticPixelWidth(cls string) (int, bool) {
	if hasBreakpointPrefix(cls) {
		return 0, false
	}
	raw := ""
	if strings.HasPrefix(cls, "w-[") && strings.HasSuffix(cls, "]") {
		raw = cls[len("w-[") : len(cls)-1]
	} else if strings.HasPrefix(cls, "min-w-[") && strings.HasSuffix(cls, "]") {
		raw = cls[len("min-w-[") : len(cls)-1]
	}
	if raw == "" || !strings.HasSuffix(raw, "px") {
		return 0, false
	}
	valStr := strings.TrimSuffix(raw, "px")
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, false
	}
	return val, true
}

func hasFluidWidthBoundary(classes []string) bool {
	for _, cls := range classes {
		if cls == "max-w-full" || cls == "max-w-[100%]" || cls == "w-full" {
			return true
		}
	}
	return false
}

func isViewportHeightClass(cls string) bool {
	if hasBreakpointPrefix(cls) {
		return false
	}
	switch cls {
	case "h-screen", "min-h-screen", "max-h-screen", "h-[100vh]", "min-h-[100vh]", "max-h-[100vh]":
		return true
	default:
		return false
	}
}

func isBottomDocked(classes []string) bool {
	hasDockPosition := false
	hasBottomZero := false
	for _, cls := range classes {
		if cls == "fixed" || cls == "sticky" {
			hasDockPosition = true
		}
		if cls == "bottom-0" {
			hasBottomZero = true
		}
	}
	return hasDockPosition && hasBottomZero
}

func hasSafeAreaBottomPadding(classes []string) bool {
	for _, cls := range classes {
		switch cls {
		case "pb-safe", "p-safe", "safe-area-bottom", "pb-[env(safe-area-inset-bottom)]":
			return true
		}
		if strings.Contains(cls, "safe-area-inset-bottom") {
			return true
		}
	}
	return false
}

func cleanAttrValue(v string) string {
	s := strings.TrimSpace(v)
	s = strings.Trim(s, "\"'`")
	return strings.TrimSpace(s)
}

func hasDeviceWidth(content string) bool {
	return strings.Contains(content, "width=device-width")
}

func hasViewportFitCover(content string) bool {
	return strings.Contains(content, "viewport-fit=cover")
}

func isFlexContainer(node *ir.Node) bool {
	if node == nil || node.Type != ir.NodeElement {
		return false
	}
	for _, cls := range node.Classes {
		if hasBreakpointPrefix(cls) {
			continue
		}
		if cls == "flex" || cls == "inline-flex" {
			return true
		}
	}
	return false
}

func hasFlexChildMinBoundary(node *ir.Node) bool {
	if node == nil {
		return false
	}
	for _, cls := range node.Classes {
		switch cls {
		case "min-w-0", "min-w-full", "overflow-hidden", "overflow-auto", "overflow-x-auto", "w-0":
			return true
		}
		if strings.HasPrefix(cls, "w-") && !hasBreakpointPrefix(cls) {
			if cls != "w-full" && cls != "w-auto" && !strings.Contains(cls, "screen") {
				return true
			}
		}
	}
	return false
}

func hasPotentiallyOverflowingContent(node *ir.Node) bool {
	if node == nil {
		return false
	}
	for _, cls := range node.Classes {
		if cls == "w-full" || cls == "flex-1" || cls == "truncate" {
			return true
		}
	}
	switch node.Tag {
	case "p", "span", "code", "pre", "h1", "h2", "h3", "h4", "h5", "h6":
		return true
	}
	for _, child := range node.Children {
		if child == nil {
			continue
		}
		switch child.Tag {
		case "p", "span", "code", "pre", "h1", "h2", "h3", "h4", "h5", "h6":
			return true
		}
		for _, cls := range child.Classes {
			if cls == "truncate" || cls == "whitespace-nowrap" {
				return true
			}
		}
	}
	return false
}

func isMediaTag(tag string) bool {
	switch tag {
	case "img", "video", "svg", "picture", "canvas", "Image":
		return true
	default:
		return false
	}
}

func extractMediaWidth(node *ir.Node) (int, bool) {
	if node == nil {
		return 0, false
	}
	for _, cls := range node.Classes {
		if w, ok := extractStaticPixelWidth(cls); ok && w > 320 {
			return w, true
		}
	}
	if node.Attributes != nil {
		if wStr, ok := node.Attributes["width"]; ok {
			wStr = cleanAttrValue(wStr)
			wStr = strings.TrimSuffix(wStr, "px")
			if w, err := strconv.Atoi(wStr); err == nil && w > 320 {
				return w, true
			}
		}
	}
	return 0, false
}

func hasResponsiveMediaScaling(node *ir.Node) bool {
	if node == nil {
		return false
	}
	for _, cls := range node.Classes {
		switch cls {
		case "max-w-full", "max-w-[100%]", "w-full":
			return true
		}
	}
	return false
}

func hasNowrap(classes []string) bool {
	for _, cls := range classes {
		if hasBreakpointPrefix(cls) {
			continue
		}
		if cls == "whitespace-nowrap" {
			return true
		}
	}
	return false
}

func hasTextOverflowProtection(classes []string) bool {
	for _, cls := range classes {
		switch cls {
		case "truncate", "overflow-hidden", "overflow-x-auto", "overflow-auto", "break-words", "break-all":
			return true
		}
	}
	return false
}

func hasCodeWrapOrScroll(node *ir.Node) bool {
	if node == nil {
		return false
	}
	for _, cls := range node.Classes {
		switch cls {
		case "break-all", "break-words", "whitespace-normal", "whitespace-pre-wrap":
			return true
		}
	}
	return hasScrollContainerAncestor(node)
}

func isDesktopHidden(classes []string) bool {
	hasHiddenMobile := false
	hasDesktopDisplay := false
	for _, cls := range classes {
		if cls == "hidden" {
			hasHiddenMobile = true
		}
		if strings.HasPrefix(cls, "sm:") || strings.HasPrefix(cls, "md:") || strings.HasPrefix(cls, "lg:") || strings.HasPrefix(cls, "xl:") || strings.HasPrefix(cls, "2xl:") {
			suffix := cls[strings.Index(cls, ":")+1:]
			switch suffix {
			case "block", "inline-block", "inline", "flex", "inline-flex", "grid", "inline-grid", "table":
				hasDesktopDisplay = true
			}
		}
	}
	return hasHiddenMobile && hasDesktopDisplay
}

func isPrimaryAction(node *ir.Node) bool {
	if node == nil || !isActionControlTag(node.Tag) {
		return false
	}
	if hasPrimaryActionAttribute(node) || hasPrimaryActionClass(node.Classes) {
		return true
	}
	return hasPrimaryActionText(node.Children)
}

func isActionControlTag(tag string) bool {
	return tag == "button" || tag == "a" || tag == "input"
}

func hasPrimaryActionAttribute(node *ir.Node) bool {
	if node.Attributes == nil {
		return false
	}
	if node.Attributes["type"] == "submit" {
		return true
	}
	if aria, ok := node.Attributes["aria-label"]; ok {
		return isActionKeyword(strings.ToLower(aria))
	}
	return false
}

func hasPrimaryActionClass(classes []string) bool {
	for _, cls := range classes {
		if cls == "bg-primary" || cls == "bg-destructive" {
			return true
		}
	}
	return false
}

func hasPrimaryActionText(children []*ir.Node) bool {
	for _, child := range children {
		if child != nil && child.Type == ir.NodeText {
			lower := strings.ToLower(strings.TrimSpace(child.RawClasses))
			if isActionKeyword(lower) {
				return true
			}
		}
	}
	return false
}

func isActionKeyword(text string) bool {
	keywords := [...]string{
		"simpan", "save", "bayar", "pay", "checkout", "kirim", "send", "submit",
		"konfirmasi", "confirm", "beli", "buy", "daftar", "register", "selesai", "finish",
	}
	for _, kw := range keywords {
		if strings.Contains(text, kw) {
			return true
		}
	}
	return false
}

func countInteractiveChildren(node *ir.Node) int {
	if node == nil {
		return 0
	}
	count := 0
	for _, child := range node.Children {
		if isInteractiveElement(child) {
			count++
		}
	}
	return count
}

func isInteractiveElement(child *ir.Node) bool {
	if child == nil || child.Type != ir.NodeElement {
		return false
	}
	switch child.Tag {
	case "button", "Button", "a", "Link":
		return true
	case "input":
		if child.Attributes != nil {
			t := child.Attributes["type"]
			return t == "button" || t == "submit" || t == "reset"
		}
	}
	return false
}

func isHorizontalFlexRow(classes []string) bool {
	isFlex := false
	isCol := false
	isWrap := false
	hasScroll := false
	for _, cls := range classes {
		if hasBreakpointPrefix(cls) {
			continue
		}
		switch cls {
		case "flex", "inline-flex":
			isFlex = true
		case "flex-col", "flex-col-reverse":
			isCol = true
		case "flex-wrap", "flex-wrap-reverse":
			isWrap = true
		case "overflow-x-auto", "overflow-x-scroll", "overflow-auto", "overflow-scroll":
			hasScroll = true
		}
	}
	return isFlex && !isCol && !isWrap && !hasScroll
}

func isDynamicViewportHeightClass(cls string) bool {
	if hasBreakpointPrefix(cls) {
		return false
	}
	switch cls {
	case "h-dvh", "min-h-dvh", "max-h-dvh", "h-svh", "min-h-svh", "max-h-svh",
		"h-[100dvh]", "min-h-[100dvh]", "max-h-[100dvh]", "h-[100svh]", "min-h-[100svh]", "max-h-[100svh]":
		return true
	default:
		return false
	}
}

func hasDynamicViewportHeight(classes []string) bool {
	for _, cls := range classes {
		if isDynamicViewportHeightClass(cls) {
			return true
		}
	}
	return false
}

func hasClassicalViewportHeight(classes []string) bool {
	for _, cls := range classes {
		if isViewportHeightClass(cls) {
			return true
		}
	}
	return false
}

func isFixedBottom(classes []string) bool {
	hasFixed := false
	hasBottom := false
	for _, cls := range classes {
		if hasBreakpointPrefix(cls) {
			continue
		}
		switch cls {
		case "fixed", "sticky":
			hasFixed = true
		case "bottom-0", "inset-x-0":
			hasBottom = true
		}
	}
	return hasFixed && hasBottom
}

func hasScrollableRegionAncestor(node *ir.Node) bool {
	curr := node
	for curr != nil {
		for _, cls := range curr.Classes {
			switch cls {
			case "overflow-y-auto", "overflow-y-scroll", "overflow-auto", "overflow-scroll":
				return true
			}
		}
		curr = curr.Parent
	}
	return false
}

func hasFormInputDescendant(node *ir.Node) bool {
	if node == nil {
		return false
	}
	for child := range node.Walk() {
		if child != nil && child.Type == ir.NodeElement {
			if child.Tag == "input" || child.Tag == "textarea" || child.Tag == "select" {
				return true
			}
		}
	}
	return false
}

func hasExcessiveMobilePadding(classes []string) bool {
	hasNarrowMaxW := false
	hasHugePadding := false
	hasModeratePadding := false

	for _, cls := range classes {
		if hasBreakpointPrefix(cls) {
			continue
		}
		if cls == "max-w-xs" || cls == "max-w-64" || cls == "max-w-72" {
			hasNarrowMaxW = true
		}
		switch {
		case isHugePaddingClass(cls):
			hasHugePadding = true
		case isModeratePaddingClass(cls):
			hasModeratePadding = true
		}
	}
	return hasHugePadding || (hasNarrowMaxW && hasModeratePadding)
}

func isHugePaddingClass(cls string) bool {
	switch cls {
	case "px-16", "px-20", "px-24", "px-28", "px-32",
		"p-16", "p-20", "p-24", "p-28", "p-32":
		return true
	default:
		return false
	}
}

func isModeratePaddingClass(cls string) bool {
	switch cls {
	case "px-10", "px-12", "px-14", "p-10", "p-12", "p-14":
		return true
	default:
		return false
	}
}

func hasExcessiveGridMinColumn(classes []string) (string, bool) {
	for _, cls := range classes {
		if hasBreakpointPrefix(cls) {
			continue
		}
		if !strings.HasPrefix(cls, "grid-cols-[") {
			continue
		}
		if isOverlargeGridClass(cls) {
			return cls, true
		}
	}
	return "", false
}

func isOverlargeGridClass(cls string) bool {
	lower := strings.ToLower(cls)
	if !strings.Contains(lower, "minmax(") {
		return false
	}
	overlargeSizes := [...]string{
		"350px", "360px", "375px", "380px", "400px", "450px", "500px", "600px",
		"22rem", "24rem", "25rem", "28rem", "30rem", "32rem",
	}
	for _, sz := range overlargeSizes {
		if strings.Contains(lower, sz) {
			return true
		}
	}
	return false
}

func hasConflictingAspectRatio(classes []string) bool {
	hasAspect := false
	hasRigidHeight := false
	hasFluidWidth := false

	for _, cls := range classes {
		if hasBreakpointPrefix(cls) {
			continue
		}
		switch {
		case strings.HasPrefix(cls, "aspect-"):
			hasAspect = true
		case isRigidHeightClass(cls):
			hasRigidHeight = true
		case cls == "w-full" || cls == "max-w-full":
			hasFluidWidth = true
		}
	}
	return hasAspect && hasRigidHeight && !hasFluidWidth
}

func isRigidHeightClass(cls string) bool {
	if strings.HasPrefix(cls, "h-[") && strings.HasSuffix(cls, "px]") {
		return true
	}
	switch cls {
	case "h-64", "h-72", "h-80", "h-96":
		return true
	default:
		return false
	}
}
