package mobile

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// cleanAttrValue membersihkan nilai atribut dari tanda kutip dan spasi berlebih.
func cleanAttrValue(v string) string {
	return strings.Trim(strings.TrimSpace(strings.ToLower(v)), "\"'`{}")
}

// isDesktopOnly memeriksa apakah elemen atau ancestor-nya secara eksplisit disembunyikan di layar ponsel.
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

// hasDesktopOnlyClass mendeteksi utility responsive yang menyembunyikan elemen pada layar ponsel.
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

// isInteractiveElement memeriksa apakah node merupakan elemen interaktif pengguna.
func isInteractiveElement(node *ir.Node) bool {
	if node == nil || node.Type != ir.NodeElement {
		return false
	}
	tagLower := strings.ToLower(node.Tag)
	switch tagLower {
	case "button", "a", "input", "select", "textarea", "summary":
		return true
	}
	if node.Attributes != nil {
		if role, ok := node.Attributes["role"]; ok {
			cleanRole := cleanAttrValue(role)
			if cleanRole == "button" || cleanRole == "link" {
				return true
			}
		}
	}
	return false
}
