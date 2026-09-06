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
