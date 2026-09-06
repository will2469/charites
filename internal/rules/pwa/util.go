package pwa

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// ManifestIcon merepresentasikan entri ikon dalam Web App Manifest.
type ManifestIcon struct {
	Src     string `json:"src"`
	Sizes   string `json:"sizes"`
	Type    string `json:"type"`
	Purpose string `json:"purpose"`
}

// WebAppManifest merepresentasikan struktur data esensial Web App Manifest.
type WebAppManifest struct {
	Name      string
	ShortName string
	StartURL  string
	Display   string
	Icons     []ManifestIcon
	HasIcons  bool
}

type rawManifest struct {
	Name       *string         `json:"name"`
	ShortName  *string         `json:"short_name"`
	ShortNameC *string         `json:"shortName"`
	StartURL   *string         `json:"start_url"`
	StartURLC  *string         `json:"startUrl"`
	Display    *string         `json:"display"`
	Icons      *[]ManifestIcon `json:"icons"`
}

var (
	reBareKey     = regexp.MustCompile(`(?m)([{,]\s*)([a-zA-Z_][a-zA-Z0-9_]*)\s*:`)
	reSingleQuote = regexp.MustCompile(`'([^'\\]*(?:\\.[^'\\]*)*)'`)
	reTrailingSep = regexp.MustCompile(`,\s*([}\]])`)
)

func cleanAttrValue(v string) string {
	s := strings.TrimSpace(v)
	s = strings.Trim(s, "\"'`")
	return strings.TrimSpace(s)
}

func isManifestScript(node *ir.Node) bool {
	if node == nil || node.Type != ir.NodeElement {
		return false
	}
	if !strings.EqualFold(node.Tag, "script") {
		return false
	}
	t, ok := node.GetAttr("type")
	if !ok {
		return false
	}
	return strings.EqualFold(cleanAttrValue(t), "application/manifest+json")
}

func isHeadElement(tag string) bool {
	return strings.EqualFold(tag, "head") || strings.EqualFold(tag, "helmet")
}

func containsWord(s, word string) bool {
	for _, part := range strings.Fields(s) {
		if strings.EqualFold(part, word) {
			return true
		}
	}
	return false
}

func hasManifestLink(headNode *ir.Node) bool {
	if headNode == nil {
		return false
	}
	for child := range headNode.Walk() {
		if child.Type != ir.NodeElement {
			continue
		}
		if !strings.EqualFold(child.Tag, "link") {
			continue
		}
		rel, ok := child.GetAttr("rel")
		if !ok || !containsWord(cleanAttrValue(rel), "manifest") {
			continue
		}
		href, ok := child.GetAttr("href")
		if !ok || cleanAttrValue(href) == "" {
			continue
		}
		return true
	}
	return false
}

func extractScriptText(node *ir.Node) string {
	if node == nil {
		return ""
	}
	var sb strings.Builder
	for _, ch := range node.Children {
		if ch.Type == ir.NodeText {
			sb.WriteString(ch.RawClasses)
		}
	}
	if sb.Len() > 0 {
		return sb.String()
	}
	if val, ok := node.GetAttr("dangerouslysetinnerhtml"); ok {
		return val
	}
	if val, ok := node.GetAttr("set:html"); ok {
		return val
	}
	return node.RawClasses
}

func isQuote(r rune) bool {
	return r == '"' || r == '\'' || r == '`'
}

func updateDepth(r rune, depth int) (int, bool) {
	if r == '{' {
		return depth + 1, false
	}
	if r == '}' {
		newDepth := depth - 1
		return newDepth, newDepth == 0
	}
	return depth, false
}

func findJSONStartIndex(s string) int {
	start := strings.Index(s, "{")
	if start == -1 {
		return -1
	}
	jsIdx := strings.Index(s, "JSON.stringify")
	if jsIdx != -1 {
		innerStart := strings.Index(s[jsIdx:], "{")
		if innerStart != -1 {
			return jsIdx + innerStart
		}
	}
	return start
}

func findMatchingBrace(runes []rune) int {
	depth := 0
	inQuote := false
	var quoteChar rune

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if inQuote {
			if r == '\\' && i+1 < len(runes) {
				i++
				continue
			}
			if r == quoteChar {
				inQuote = false
			}
			continue
		}

		if isQuote(r) {
			inQuote = true
			quoteChar = r
			continue
		}

		var done bool
		depth, done = updateDepth(r, depth)
		if done {
			return i
		}
	}
	return -1
}

func extractJSONObject(s string) string {
	start := findJSONStartIndex(s)
	if start == -1 {
		return ""
	}
	runes := []rune(s[start:])
	end := findMatchingBrace(runes)
	if end == -1 {
		return ""
	}
	return string(runes[:end+1])
}

func relaxedJSToStrictJSON(s string) string {
	normalized := reSingleQuote.ReplaceAllString(s, `"$1"`)
	normalized = reBareKey.ReplaceAllString(normalized, `$1"$2":`)
	normalized = reTrailingSep.ReplaceAllString(normalized, `$1`)
	return normalized
}

func parseManifest(node *ir.Node) (*WebAppManifest, bool) {
	raw := extractScriptText(node)
	if strings.TrimSpace(raw) == "" {
		return nil, false
	}

	objStr := extractJSONObject(raw)
	if objStr == "" {
		return nil, false
	}

	var rawData rawManifest
	if err := json.Unmarshal([]byte(objStr), &rawData); err != nil {
		normalized := relaxedJSToStrictJSON(objStr)
		if nErr := json.Unmarshal([]byte(normalized), &rawData); nErr != nil {
			return fallbackParseManifest(objStr), true
		}
	}

	manifest := &WebAppManifest{}
	if rawData.Name != nil {
		manifest.Name = strings.TrimSpace(*rawData.Name)
	}
	if rawData.ShortName != nil {
		manifest.ShortName = strings.TrimSpace(*rawData.ShortName)
	} else if rawData.ShortNameC != nil {
		manifest.ShortName = strings.TrimSpace(*rawData.ShortNameC)
	}
	if rawData.StartURL != nil {
		manifest.StartURL = strings.TrimSpace(*rawData.StartURL)
	} else if rawData.StartURLC != nil {
		manifest.StartURL = strings.TrimSpace(*rawData.StartURLC)
	}
	if rawData.Display != nil {
		manifest.Display = strings.TrimSpace(*rawData.Display)
	}
	if rawData.Icons != nil {
		manifest.HasIcons = true
		manifest.Icons = *rawData.Icons
	}

	return manifest, true
}

func fallbackParseManifest(s string) *WebAppManifest {
	manifest := &WebAppManifest{}
	if strings.Contains(s, "name") {
		manifest.Name = "fallback"
	}
	if strings.Contains(s, "start_url") || strings.Contains(s, "startUrl") {
		manifest.StartURL = "/"
	}
	if strings.Contains(s, "display") {
		manifest.Display = "standalone"
	}
	if strings.Contains(s, "icons") {
		manifest.HasIcons = true
	}
	return manifest
}
