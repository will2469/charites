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

func hasAppleCapableMeta(headNode *ir.Node) bool {
	if headNode == nil {
		return false
	}
	for child := range headNode.Walk() {
		if isAppleCapableMetaNode(child) {
			return true
		}
	}
	return false
}

func isAppleCapableMetaNode(node *ir.Node) bool {
	if node == nil || node.Type != ir.NodeElement || !strings.EqualFold(node.Tag, "meta") {
		return false
	}
	name, ok := node.GetAttr("name")
	if !ok || !strings.EqualFold(cleanAttrValue(name), "apple-mobile-web-app-capable") {
		return false
	}
	content, ok := node.GetAttr("content")
	return ok && strings.EqualFold(cleanAttrValue(content), "yes")
}

func hasAppleTouchIcon(headNode *ir.Node) bool {
	if headNode == nil {
		return false
	}
	for child := range headNode.Walk() {
		if isAppleTouchIconNode(child) {
			return true
		}
	}
	return false
}

func isAppleTouchIconNode(node *ir.Node) bool {
	if node == nil || node.Type != ir.NodeElement || !strings.EqualFold(node.Tag, "link") {
		return false
	}
	rel, ok := node.GetAttr("rel")
	if !ok {
		return false
	}
	cleanRel := strings.ToLower(cleanAttrValue(rel))
	if !strings.Contains(cleanRel, "apple-touch-icon") {
		return false
	}
	href, ok := node.GetAttr("href")
	return ok && cleanAttrValue(href) != ""
}

func isResourceElement(tag string) bool {
	switch strings.ToLower(tag) {
	case "script", "link", "img", "iframe", "video", "audio", "source", "track", "embed", "object":
		return true
	default:
		return false
	}
}

func isInsecureResourceURL(rawURL string) bool {
	clean := strings.TrimSpace(cleanAttrValue(rawURL))
	lower := strings.ToLower(clean)
	if !strings.HasPrefix(lower, "http://") {
		return false
	}
	rest := lower[len("http://"):]
	if strings.HasPrefix(rest, "localhost") || strings.HasPrefix(rest, "127.0.0.1") {
		return false
	}
	return true
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

func skipLineComment(runes []rune, i, n int) int {
	i += 2
	for i < n && runes[i] != '\n' {
		i++
	}
	return i
}

func skipBlockComment(runes []rune, i, n int) int {
	i += 2
	for i+1 < n && !(runes[i] == '*' && runes[i+1] == '/') {
		i++
	}
	return i + 2
}

func stripJSComments(code string) string {
	var sb strings.Builder
	runes := []rune(code)
	n := len(runes)
	i := 0
	for i < n {
		if runes[i] == '/' && i+1 < n {
			if runes[i+1] == '/' {
				i = skipLineComment(runes, i, n)
				continue
			}
			if runes[i+1] == '*' {
				i = skipBlockComment(runes, i, n)
				continue
			}
		}
		sb.WriteRune(runes[i])
		i++
	}
	return sb.String()
}

func hasServiceWorkerRegistration(headNode *ir.Node) bool {
	if headNode == nil {
		return false
	}
	root := headNode
	for root.Parent != nil {
		root = root.Parent
	}

	for child := range root.Walk() {
		if child.Type != ir.NodeElement || !strings.EqualFold(child.Tag, "script") {
			continue
		}
		if isSWScriptTag(child) {
			return true
		}
	}
	return false
}

func isSWScriptTag(child *ir.Node) bool {
	if src, ok := child.GetAttr("src"); ok {
		cleanSrc := strings.ToLower(cleanAttrValue(src))
		if strings.Contains(cleanSrc, "sw.js") ||
			strings.Contains(cleanSrc, "service-worker.js") ||
			strings.Contains(cleanSrc, "serviceworker.js") ||
			strings.Contains(cleanSrc, "register-sw.js") ||
			strings.Contains(cleanSrc, "sw-register.js") ||
			strings.Contains(cleanSrc, "pwa.js") {
			return true
		}
	}
	txt := extractScriptText(child)
	return containsSWRegistrationSnippet(txt)
}

func containsSWRegistrationSnippet(txt string) bool {
	return strings.Contains(txt, "serviceWorker.register") ||
		strings.Contains(txt, "navigator.serviceWorker.register") ||
		strings.Contains(txt, "registerRoute") ||
		strings.Contains(txt, "workbox.register") ||
		strings.Contains(txt, "registerSW")
}

func hasSWFetchListener(txt string) bool {
	return strings.Contains(txt, "addEventListener(\"fetch\"") ||
		strings.Contains(txt, "addEventListener('fetch'") ||
		strings.Contains(txt, "addEventListener(`fetch`") ||
		strings.Contains(txt, "self.onfetch") ||
		strings.Contains(txt, "onfetch =")
}

func interceptsFetch(txt string) bool {
	return strings.Contains(txt, "respondWith(") ||
		strings.Contains(txt, "respondWith (")
}

func hasSWOfflineFallback(txt string) bool {
	lower := strings.ToLower(txt)
	if strings.Contains(lower, "caches.match") ||
		strings.Contains(lower, "cache.match") ||
		strings.Contains(lower, "caches.open") ||
		strings.Contains(lower, "caches.keys") ||
		strings.Contains(lower, "caches.has") {
		return true
	}
	if strings.Contains(txt, ".catch(") || strings.Contains(txt, ".catch (") {
		return true
	}
	return strings.Contains(txt, "catch") && (strings.Contains(lower, "offline") || strings.Contains(lower, "fallback") || strings.Contains(lower, "response"))
}

func hasSWFeatureDetection(txt string) bool {
	return strings.Contains(txt, "'serviceWorker' in navigator") ||
		strings.Contains(txt, "\"serviceWorker\" in navigator") ||
		strings.Contains(txt, "'serviceWorker' in window.navigator") ||
		strings.Contains(txt, "\"serviceWorker\" in window.navigator") ||
		strings.Contains(txt, "if (navigator.serviceWorker") ||
		strings.Contains(txt, "if (window.navigator && 'serviceWorker' in window.navigator)") ||
		strings.Contains(txt, "navigator.serviceWorker &&")
}

func hasSWErrorHandling(txt string) bool {
	if strings.Contains(txt, ".catch(") || strings.Contains(txt, ".catch (") {
		return true
	}
	return strings.Contains(txt, "try") && strings.Contains(txt, "catch")
}

func isWorkerScope(txt string) bool {
	return strings.Contains(txt, "addEventListener(\"install\"") ||
		strings.Contains(txt, "addEventListener('install'") ||
		strings.Contains(txt, "addEventListener(\"activate\"") ||
		strings.Contains(txt, "addEventListener('activate'") ||
		strings.Contains(txt, "addEventListener(\"fetch\"") ||
		strings.Contains(txt, "addEventListener('fetch'") ||
		strings.Contains(txt, "self.addEventListener") ||
		strings.Contains(txt, "clients.claim") ||
		strings.Contains(txt, "skipWaiting") ||
		strings.Contains(txt, "importScripts(")
}

func collectForbiddenWorkerAPIs(raw string) []string {
	clean := stripJSComments(raw)
	var forbidden []string

	if strings.Contains(clean, "window.") || strings.Contains(clean, "window[") {
		forbidden = append(forbidden, "window")
	}
	if strings.Contains(clean, "document.") || strings.Contains(clean, "document[") {
		forbidden = append(forbidden, "document")
	}
	if strings.Contains(clean, "localStorage.") || strings.Contains(clean, "localStorage[") {
		forbidden = append(forbidden, "localStorage")
	}
	if strings.Contains(clean, "sessionStorage.") || strings.Contains(clean, "sessionStorage[") {
		forbidden = append(forbidden, "sessionStorage")
	}
	if strings.Contains(clean, "alert(") {
		forbidden = append(forbidden, "alert()")
	}
	if strings.Contains(clean, "confirm(") {
		forbidden = append(forbidden, "confirm()")
	}
	if strings.Contains(clean, "prompt(") {
		forbidden = append(forbidden, "prompt()")
	}

	return forbidden
}
