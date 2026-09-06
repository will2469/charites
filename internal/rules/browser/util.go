package browser

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

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
