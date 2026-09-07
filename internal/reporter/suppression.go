package reporter

import (
	"fmt"
	"path/filepath"
	"strings"
)

// SuppressionDirective mengembalikan sintaks direktif supresi yang valid dan tepat
// berdasarkan ekstensi berkas target (Astro HTML comment vs TSX/JSX line comment vs CSS block comment).
func SuppressionDirective(filePath string, ruleID string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".astro", ".html":
		return fmt.Sprintf("<!-- charites:ignore %s <reason> -->", ruleID)
	case ".css":
		return fmt.Sprintf("/* charites:ignore %s <reason> */", ruleID)
	default: // .tsx, .jsx, .ts, .js, dan bahasa berbasis C-comment lainnya
		return fmt.Sprintf("// charites:ignore %s <reason>", ruleID)
	}
}
