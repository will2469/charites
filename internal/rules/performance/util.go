package performance

import (
	"bytes"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// cleanAttrVal membersihkan whitespace dan quote luar dari string nilai atribut.
func cleanAttrVal(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') ||
			(s[0] == '`' && s[len(s)-1] == '`') {
			return strings.TrimSpace(s[1 : len(s)-1])
		}
	}
	return s
}

// getFileSourceContent mengekstrak seluruh teks kode sumber berkas dari NodeComment.
func getFileSourceContent(node *ir.Node) string {
	if node == nil {
		return ""
	}
	if node.Type == ir.NodeComment && len(node.RawClasses) > 0 {
		return node.RawClasses
	}
	cur := node
	for cur.Parent != nil {
		cur = cur.Parent
	}
	for _, ch := range cur.Children {
		if ch != nil && ch.Type == ir.NodeComment && len(ch.RawClasses) > 0 {
			return ch.RawClasses
		}
	}
	return ""
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
		if child != nil && child.Type == ir.NodeText {
			sb.WriteString(child.RawClasses)
		}
	}
	res := sb.String()
	if res == "" {
		return node.RawClasses
	}
	return res
}

// isSourceRootOrScript memastikan evaluasi kode sumber script hanya dieksekusi sekali per berkas.
func isSourceRootOrScript(node *ir.Node) bool {
	if node == nil {
		return false
	}
	if node.Type == ir.NodeComment && node.Parent != nil && node.Parent.Parent == nil && len(node.RawClasses) > 0 {
		return true
	}
	if node.Type == ir.NodeElement && strings.EqualFold(node.Tag, "script") {
		return true
	}
	if node.Parent == nil && !hasNodeCommentChild(node) {
		return true
	}
	return false
}

func hasNodeCommentChild(root *ir.Node) bool {
	for _, ch := range root.Children {
		if ch != nil && ch.Type == ir.NodeComment {
			return true
		}
	}
	return false
}

// isMemoizedComponent memeriksa apakah pemanggilan komponen merujuk pada komponen yang di-memoize.
func isMemoizedComponent(tagName string, fileSrc string) bool {
	if len(tagName) == 0 {
		return false
	}
	// Native HTML elements start with lowercase
	if tagName[0] >= 'a' && tagName[0] <= 'z' {
		return false
	}
	if strings.HasPrefix(tagName, "Memo") {
		return true
	}
	if len(fileSrc) == 0 {
		return false
	}

	patterns := [...]string{
		"const " + tagName + " = memo(",
		"const " + tagName + " = React.memo(",
		"const " + tagName + " = React.memo<",
		"const " + tagName + " = memo<",
		"let " + tagName + " = memo(",
		"let " + tagName + " = React.memo(",
		"export const " + tagName + " = memo(",
		"export const " + tagName + " = React.memo(",
	}
	for _, pat := range patterns {
		if strings.Contains(fileSrc, pat) {
			return true
		}
	}
	return false
}

// isInlineLiteralProp memeriksa apakah nilai prop merupakan objek, array, atau fungsi inline instan.
func isInlineLiteralProp(val string) (string, bool) {
	trimmed := strings.TrimSpace(val)
	if len(trimmed) < 2 {
		return "", false
	}
	trimmed = cleanAttrVal(trimmed)

	if strings.HasPrefix(trimmed, "{{") || strings.HasPrefix(trimmed, "{ {") {
		return "object literal", true
	}

	if strings.HasPrefix(trimmed, "{[") || strings.HasPrefix(trimmed, "{ [") {
		return "array literal", true
	}

	if strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}") {
		inner := strings.TrimSpace(trimmed[1 : len(trimmed)-1])
		if strings.HasPrefix(inner, "() =>") || strings.HasPrefix(inner, "()=>") ||
			(strings.Contains(inner, "=>") && (strings.HasPrefix(inner, "(") || strings.HasPrefix(inner, "e =>") || strings.HasPrefix(inner, "event =>"))) ||
			strings.HasPrefix(inner, "function") {
			return "inline function", true
		}
	}

	return "", false
}

// isIndexKeyViolation memeriksa apakah atribut key merujuk pada parameter indeks .map().
func isIndexKeyViolation(node *ir.Node, fileSrc string) (string, bool) {
	if node == nil || node.Type != ir.NodeElement || node.Attributes == nil {
		return "", false
	}
	keyVal, ok := node.Attributes["key"]
	if !ok {
		return "", false
	}
	cleanKey := strings.TrimSpace(cleanAttrVal(keyVal))
	if cleanKey != "{index}" && cleanKey != "{i}" && cleanKey != "{idx}" &&
		cleanKey != "{ index }" && cleanKey != "{ i }" && cleanKey != "{ idx }" {
		return "", false
	}

	line := node.Span.Line
	lines := strings.Split(fileSrc, "\n")
	if line > 0 && line <= len(lines) {
		start := line - 10
		if start < 0 {
			start = 0
		}
		contextSnippet := strings.Join(lines[start:line], " ")

		if strings.Contains(contextSnippet, "Array.from") || strings.Contains(contextSnippet, "Array(") {
			return "", false
		}

		if isConstCollectionMap(contextSnippet) {
			return "", false
		}
	}

	identifier := strings.Trim(cleanKey, "{} ")
	return identifier, true
}

func isConstCollectionMap(snippet string) bool {
	mapIdx := strings.LastIndex(snippet, ".map(")
	if mapIdx == -1 {
		mapIdx = strings.LastIndex(snippet, ".map (")
	}
	if mapIdx == -1 {
		return false
	}
	prefix := strings.TrimSpace(snippet[:mapIdx])
	identEnd := len(prefix)
	identStart := identEnd
	for identStart > 0 {
		c := prefix[identStart-1]
		if (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
			identStart--
		} else {
			break
		}
	}
	if identStart < identEnd {
		ident := prefix[identStart:identEnd]
		if len(ident) >= 2 && ident[0] >= 'A' && ident[0] <= 'Z' {
			return true
		}
	}
	return false
}

// EffectViolation merepresentasikan temuan efek samping tanpa pembersih.
type EffectViolation struct {
	Line       int
	EffectName string
	Resource   string
}

// findMissingEffectCleanups mendeteksi efek samping yang mengakuisisi resource tanpa fungsi cleanup return.
func findMissingEffectCleanups(fileSrc string) []EffectViolation {
	if len(fileSrc) == 0 {
		return nil
	}
	if !strings.Contains(fileSrc, "useEffect") && !strings.Contains(fileSrc, "useLayoutEffect") {
		return nil
	}

	var violations []EffectViolation
	lines := strings.Split(fileSrc, "\n")
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		effectName := ""
		if strings.Contains(line, "useEffect(") || strings.Contains(line, "useEffect (") {
			effectName = "useEffect"
		} else if strings.Contains(line, "useLayoutEffect(") || strings.Contains(line, "useLayoutEffect (") {
			effectName = "useLayoutEffect"
		}
		if effectName == "" {
			continue
		}

		bodySnippet := extractEffectBody(lines, i)
		if len(bodySnippet) == 0 {
			continue
		}

		res, hasResource := acquiresPersistentResource(bodySnippet)
		if !hasResource {
			continue
		}

		if hasCleanupReturn(bodySnippet) {
			continue
		}

		violations = append(violations, EffectViolation{
			Line:       i + 1,
			EffectName: effectName,
			Resource:   res,
		})
	}

	return violations
}

func extractEffectBody(lines []string, startLine int) string {
	var sb strings.Builder
	depth := 0
	foundOpen := false
	maxLines := startLine + 60
	if maxLines > len(lines) {
		maxLines = len(lines)
	}

	for i := startLine; i < maxLines; i++ {
		l := lines[i]
		sb.WriteString(l)
		sb.WriteByte('\n')
		for j := 0; j < len(l); j++ {
			switch l[j] {
			case '{':
				depth++
				foundOpen = true
			case '}':
				depth--
				if foundOpen && depth <= 0 {
					return sb.String()
				}
			}
		}
	}
	return sb.String()
}

func acquiresPersistentResource(body string) (string, bool) {
	resources := [...]string{
		"addEventListener",
		"setInterval",
		"setTimeout",
		"ResizeObserver",
		"IntersectionObserver",
		"MutationObserver",
		"WebSocket",
		"AbortController",
	}
	for _, res := range resources {
		if strings.Contains(body, res) {
			return res, true
		}
	}
	return "", false
}

func hasCleanupReturn(body string) bool {
	return strings.Contains(body, "return () =>") ||
		strings.Contains(body, "return ()=>") ||
		strings.Contains(body, "return function") ||
		strings.Contains(body, "removeEventListener") ||
		strings.Contains(body, "clearInterval") ||
		strings.Contains(body, "clearTimeout") ||
		strings.Contains(body, "disconnect()")
}

// isOvercoupledContextProvider memeriksa apakah Context.Provider memuat > 5 properti pada prop value.
func isOvercoupledContextProvider(node *ir.Node) (int, bool) {
	if node == nil || node.Type != ir.NodeElement || node.Attributes == nil {
		return 0, false
	}
	if node.Tag != "Context.Provider" && !strings.HasSuffix(node.Tag, ".Provider") {
		return 0, false
	}

	val, ok := node.Attributes["value"]
	if !ok {
		return 0, false
	}

	trimmed := strings.TrimSpace(cleanAttrVal(val))
	if !strings.HasPrefix(trimmed, "{") || !strings.HasSuffix(trimmed, "}") {
		return 0, false
	}

	inner := strings.TrimSpace(trimmed[1 : len(trimmed)-1])
	if strings.HasPrefix(inner, "{") && strings.HasSuffix(inner, "}") {
		inner = strings.TrimSpace(inner[1 : len(inner)-1])
	}

	count := countObjectProperties(inner)
	if count > 5 {
		return count, true
	}
	return 0, false
}

func countObjectProperties(inner string) int {
	if len(inner) == 0 {
		return 0
	}
	depth := 0
	commas := 0
	hasContent := false

	for i := 0; i < len(inner); i++ {
		c := inner[i]
		switch c {
		case '{', '[', '(':
			depth++
		case '}', ']', ')':
			depth--
		case ',':
			if depth == 0 {
				commas++
			}
		default:
			if c != ' ' && c != '\t' && c != '\n' && c != '\r' {
				hasContent = true
			}
		}
	}

	if !hasContent {
		return 0
	}
	return commas + 1
}

// HeavyImportViolation merepresentasikan temuan impor statis pustaka berbobot besar.
type HeavyImportViolation struct {
	Line   int
	Module string
}

// findStaticHeavyImports mendeteksi pernyataan impor statis modul berukuran besar di tingkat atas.
func findStaticHeavyImports(fileSrc string) []HeavyImportViolation {
	if len(fileSrc) == 0 {
		return nil
	}

	heavyModules := [...]string{
		"monaco-editor",
		"chart.js",
		"echarts",
		"quill",
		"pdfjs-dist",
		"pdfjs",
		"three",
		"datatables.net",
		"xlsx",
		"jspdf",
		"plotly.js",
		"AnalyticalChart",
	}

	var violations []HeavyImportViolation
	lines := strings.Split(fileSrc, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "import ") || strings.HasPrefix(trimmed, "import type ") || strings.HasPrefix(trimmed, "import type{") || strings.Contains(trimmed, "{ type ") {
			continue
		}
		if strings.Contains(trimmed, "lazy(") || strings.Contains(trimmed, "React.lazy(") {
			continue
		}

		for _, mod := range heavyModules {
			if strings.Contains(trimmed, "'"+mod) || strings.Contains(trimmed, "\""+mod) || strings.Contains(trimmed, "/"+mod) {
				violations = append(violations, HeavyImportViolation{
					Line:   i + 1,
					Module: mod,
				})
				break
			}
		}
	}

	return violations
}

// RedundantMemoViolation merepresentasikan temuan pembungkusan useCallback yang hanya dipakai elemen native.
type RedundantMemoViolation struct {
	Line        int
	HandlerName string
	ElementTag  string
}

// findRedundantFunctionMemoizations mendeteksi penggunaan useCallback pada fungsi yang hanya dikonsumsi DOM native.
func findRedundantFunctionMemoizations(fileSrc string) []RedundantMemoViolation {
	if len(fileSrc) == 0 || !strings.Contains(fileSrc, "useCallback") {
		return nil
	}

	var violations []RedundantMemoViolation
	lines := strings.Split(fileSrc, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "export ") {
			continue
		}
		if !strings.Contains(trimmed, "useCallback(") && !strings.Contains(trimmed, "useCallback (") {
			continue
		}

		handler := extractHandlerName(trimmed)
		if handler == "" {
			continue
		}

		if isHandlerConsumedByMemoOrDeps(handler, fileSrc) {
			continue
		}

		tag := findNativeTagConsumer(handler, fileSrc)
		if tag != "" {
			violations = append(violations, RedundantMemoViolation{
				Line:        i + 1,
				HandlerName: handler,
				ElementTag:  tag,
			})
		}
	}

	return violations
}

func extractHandlerName(line string) string {
	parts := strings.Fields(line)
	for idx, p := range parts {
		if (p == "const" || p == "let" || p == "var") && idx+1 < len(parts) {
			name := parts[idx+1]
			name = strings.Trim(name, ":= ")
			return name
		}
	}
	return ""
}

func isHandlerConsumedByMemoOrDeps(handler, fileSrc string) bool {
	// Jika handler dijadikan dependensi hook: [..., handler, ...]
	if strings.Contains(fileSrc, "["+handler+"]") ||
		strings.Contains(fileSrc, ", "+handler+"]") ||
		strings.Contains(fileSrc, "["+handler+",") ||
		strings.Contains(fileSrc, ", "+handler+",") {
		return true
	}

	// Jika handler diteruskan ke komponen React kustom (huruf kapital)
	lines := strings.Split(fileSrc, "\n")
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if !strings.Contains(trimmed, "{"+handler+"}") && !strings.Contains(trimmed, "{ "+handler+" }") {
			continue
		}
		if strings.HasPrefix(trimmed, "<") && len(trimmed) > 1 {
			tagChar := trimmed[1]
			if tagChar >= 'A' && tagChar <= 'Z' {
				return true
			}
		}
	}

	return false
}

func findNativeTagConsumer(handler, fileSrc string) string {
	nativeTags := [...]string{"button", "input", "select", "textarea", "form", "div", "a", "span"}
	lines := strings.Split(fileSrc, "\n")
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if !strings.Contains(trimmed, "{"+handler+"}") && !strings.Contains(trimmed, "{ "+handler+" }") {
			continue
		}
		for _, tag := range nativeTags {
			if strings.Contains(trimmed, "<"+tag+" ") || strings.HasPrefix(trimmed, "<"+tag+">") {
				return tag
			}
		}
	}
	return ""
}

// DerivedStateViolation merepresentasikan sinkronisasi derived state via useEffect.
type DerivedStateViolation struct {
	Line       int
	EffectName string
	StateName  string
}

// findDerivedStateInEffects mendeteksi useEffect yang murni melakukan sinkronisasi derived state.
func findDerivedStateInEffects(fileSrc string) []DerivedStateViolation {
	if len(fileSrc) == 0 || (!strings.Contains(fileSrc, "useEffect") && !strings.Contains(fileSrc, "useLayoutEffect")) {
		return nil
	}

	var violations []DerivedStateViolation
	lines := strings.Split(fileSrc, "\n")
	for i, line := range lines {
		effectName := ""
		if strings.Contains(line, "useEffect(") || strings.Contains(line, "useEffect (") {
			effectName = "useEffect"
		} else if strings.Contains(line, "useLayoutEffect(") || strings.Contains(line, "useLayoutEffect (") {
			effectName = "useLayoutEffect"
		}
		if effectName == "" {
			continue
		}

		body := extractEffectBody(lines, i)
		if body == "" || hasAsyncOrSideEffect(body) {
			continue
		}

		stateName := extractDerivedSetterCall(body)
		if stateName != "" {
			violations = append(violations, DerivedStateViolation{
				Line:       i + 1,
				EffectName: effectName,
				StateName:  stateName,
			})
		}
	}

	return violations
}

func hasAsyncOrSideEffect(body string) bool {
	sideEffects := [...]string{
		"fetch(", "axios.", "async ", "await ", ".then(", ".catch(",
		"setTimeout(", "setInterval(", "addEventListener(",
		"document.", "window.", "localStorage.", "sessionStorage.",
		"WebSocket", "ResizeObserver", "IntersectionObserver", ".current",
	}
	for _, se := range sideEffects {
		if strings.Contains(body, se) {
			return true
		}
	}
	return false
}

func extractDerivedSetterCall(body string) string {
	lines := strings.Split(body, "\n")
	setterCount := 0
	foundSetter := ""

	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed == "" || strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") {
			continue
		}
		idx := strings.Index(trimmed, "set")
		if idx != -1 && len(trimmed) > idx+3 {
			charAfterSet := trimmed[idx+3]
			if charAfterSet >= 'A' && charAfterSet <= 'Z' {
				openParen := strings.Index(trimmed[idx:], "(")
				if openParen != -1 {
					foundSetter = trimmed[idx : idx+openParen]
					setterCount++
				}
			}
		}
	}

	if setterCount == 1 {
		return foundSetter
	}
	return ""
}

// UnstableHookViolation merepresentasikan custom hook yang mengembalikan fungsi tidak stabil.
type UnstableHookViolation struct {
	Line         int
	HookName     string
	FunctionName string
}

// findUnstableHookReferences mendeteksi custom hook yang mereturn referensi fungsi tidak stabil.
func findUnstableHookReferences(fileSrc string) []UnstableHookViolation {
	if len(fileSrc) == 0 {
		return nil
	}

	var violations []UnstableHookViolation
	lines := strings.Split(fileSrc, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		hookName := extractCustomHookDeclaration(trimmed)
		if hookName == "" {
			continue
		}

		hookBody := extractFunctionBlock(lines, i)
		if hookBody == "" {
			continue
		}

		unstableFuncs := findUnstableReturnedFunctions(hookBody)
		for _, fn := range unstableFuncs {
			violations = append(violations, UnstableHookViolation{
				Line:         i + 1,
				HookName:     hookName,
				FunctionName: fn,
			})
		}
	}

	return violations
}

func extractCustomHookDeclaration(line string) string {
	if name := extractFunctionHook(line); name != "" {
		return name
	}
	return extractConstHook(line)
}

func extractFunctionHook(line string) string {
	if !strings.Contains(line, "function use") {
		return ""
	}
	idx := strings.Index(line, "use")
	end := strings.Index(line[idx:], "(")
	if end == -1 {
		return ""
	}
	name := strings.TrimSpace(line[idx : idx+end])
	if len(name) > 3 && name[3] >= 'A' && name[3] <= 'Z' {
		return name
	}
	return ""
}

func extractConstHook(line string) string {
	if !strings.Contains(line, "use") || !strings.Contains(line, "=") {
		return ""
	}
	parts := strings.Fields(line)
	for _, p := range parts {
		cand := strings.Trim(p, "=:")
		if len(cand) > 3 && strings.HasPrefix(cand, "use") && cand[3] >= 'A' && cand[3] <= 'Z' {
			return cand
		}
	}
	return ""
}

func extractFunctionBlock(lines []string, startLine int) string {
	var sb strings.Builder
	depth := 0
	foundOpen := false
	maxLines := startLine + 100
	if maxLines > len(lines) {
		maxLines = len(lines)
	}

	for i := startLine; i < maxLines; i++ {
		l := lines[i]
		sb.WriteString(l)
		sb.WriteByte('\n')
		for j := 0; j < len(l); j++ {
			switch l[j] {
			case '{':
				depth++
				foundOpen = true
			case '}':
				depth--
				if foundOpen && depth <= 0 {
					return sb.String()
				}
			}
		}
	}
	return sb.String()
}

func findUnstableReturnedFunctions(hookBody string) []string {
	returnIdx := strings.LastIndex(hookBody, "return {")
	if returnIdx == -1 {
		return nil
	}

	returnBlock := hookBody[returnIdx:]
	endBrace := strings.Index(returnBlock, "}")
	if endBrace != -1 {
		returnBlock = returnBlock[:endBrace]
	}
	if idx := strings.Index(returnBlock, "{"); idx != -1 {
		returnBlock = returnBlock[idx+1:]
	}

	var unstable []string
	items := strings.Split(returnBlock, ",")
	for _, it := range items {
		trimmed := strings.TrimSpace(it)
		trimmed = strings.Trim(trimmed, ";")
		if trimmed == "" {
			continue
		}

		if strings.Contains(trimmed, "=>") || strings.Contains(trimmed, "function") {
			colonIdx := strings.Index(trimmed, ":")
			if colonIdx != -1 {
				propName := strings.TrimSpace(trimmed[:colonIdx])
				unstable = append(unstable, propName)
			}
			continue
		}

		ident := strings.TrimSpace(trimmed)
		if strings.Contains(ident, ":") {
			colonIdx := strings.Index(ident, ":")
			ident = strings.TrimSpace(ident[colonIdx+1:])
		}
		if isUnmemoizedLocalFunction(ident, hookBody) {
			unstable = append(unstable, ident)
		}
	}

	return unstable
}

func isUnmemoizedLocalFunction(ident, hookBody string) bool {
	lines := strings.Split(hookBody, "\n")
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if strings.Contains(trimmed, "const "+ident+" =") || strings.Contains(trimmed, "let "+ident+" =") {
			if strings.Contains(trimmed, "useCallback") {
				return false
			}
			if strings.Contains(trimmed, "=>") || strings.Contains(trimmed, "function") {
				return true
			}
		}
		if strings.HasPrefix(trimmed, "function "+ident) {
			return true
		}
	}
	return false
}

// hasClientDirective memeriksa apakah node memiliki direktif hidrasi client Astro.
func hasClientDirective(node *ir.Node) (string, bool) {
	if node == nil || len(node.Attributes) == 0 {
		return "", false
	}
	directives := [...]string{
		"client:load",
		"client:idle",
		"client:visible",
		"client:only",
		"client:media",
	}
	for _, dir := range directives {
		if _, ok := node.Attributes[dir]; ok {
			return dir, true
		}
	}
	return "", false
}

// isStaticComponentTag memeriksa apakah nama komponen mengindikasikan komponen statis murni.
func isStaticComponentTag(tag string) bool {
	return strings.Contains(tag, "Static") || strings.Contains(tag, "static")
}

// findNestedIslandOverlap mendeteksi penyarangan pulau interaktif di dalam pulau lain tanpa isolasi slot.
func findNestedIslandOverlap(node *ir.Node) (*ir.Node, string, bool) {
	if _, isIsland := hasClientDirective(node); !isIsland {
		return nil, "", false
	}

	for _, ch := range node.Children {
		if ch == nil {
			continue
		}
		if n, dir, found := checkChildIslandOverlap(ch); found {
			return n, dir, true
		}
	}
	return nil, "", false
}

func checkChildIslandOverlap(curr *ir.Node) (*ir.Node, string, bool) {
	if curr == nil || isSlotIsolated(curr) {
		return nil, "", false
	}
	if dir, hasDir := hasClientDirective(curr); hasDir {
		return curr, dir, true
	}
	for _, ch := range curr.Children {
		if n, dir, found := checkChildIslandOverlap(ch); found {
			return n, dir, true
		}
	}
	return nil, "", false
}

func isSlotIsolated(n *ir.Node) bool {
	if n == nil {
		return false
	}
	if n.Tag == "slot" {
		return true
	}
	if n.Attributes != nil {
		if _, ok := n.Attributes["slot"]; ok {
			return true
		}
	}
	return false
}

// isLocalRawImage mendeteksi elemen <img> yang merujuk pada aset gambar lokal tanpa modul astro:assets.
func isLocalRawImage(node *ir.Node) (string, bool) {
	if node == nil || node.Type != ir.NodeElement || node.Tag != "img" {
		return "", false
	}
	src, ok := node.Attributes["src"]
	if !ok {
		return "", false
	}
	src = cleanAttrVal(src)
	if src == "" || strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") ||
		strings.HasPrefix(src, "data:") || strings.HasPrefix(src, "//") || strings.HasSuffix(src, ".svg") {
		return "", false
	}
	if strings.HasPrefix(src, "/images/") || strings.HasPrefix(src, "/favicon") || strings.HasPrefix(src, "/assets/") {
		return "", false
	}
	if strings.Contains(src, "assets/") || strings.Contains(src, "../") || strings.Contains(src, "./") || strings.HasPrefix(src, "@/") {
		return src, true
	}
	return "", false
}

// isAggressiveSecondaryPrefetch mendeteksi prefetch agresif pada tautan sekunder atau tautan footer.
func isAggressiveSecondaryPrefetch(node *ir.Node) (string, string, bool) {
	if node == nil || node.Type != ir.NodeElement || node.Tag != "a" {
		return "", "", false
	}
	prefetch, ok := node.Attributes["data-astro-prefetch"]
	if !ok {
		return "", "", false
	}
	prefetch = cleanAttrVal(prefetch)
	if prefetch != "viewport" && prefetch != "load" && prefetch != "" && prefetch != "true" {
		return "", "", false
	}

	href, hasHref := node.Attributes["href"]
	cleanHref := ""
	if hasHref {
		cleanHref = cleanAttrVal(href)
	}

	if isSecondaryHref(cleanHref) {
		return cleanHref, prefetch, true
	}

	cur := node.Parent
	for cur != nil {
		if cur.Type == ir.NodeElement && (cur.Tag == "footer" || cur.Tag == "aside") {
			return cleanHref, prefetch, true
		}
		cur = cur.Parent
	}

	return "", "", false
}

func isSecondaryHref(href string) bool {
	low := strings.ToLower(href)
	secondaryPatterns := [...]string{
		"terms", "privacy", "kebijakan", "syarat", "cookie", "legal", "disclaimer",
	}
	for _, p := range secondaryPatterns {
		if strings.Contains(low, p) {
			return true
		}
	}
	return false
}

// findDynamicClassConcatenation mendeteksi penggabungan string template literal dengan prefiks parsial Tailwind.
func findDynamicClassConcatenation(val string) (string, bool) {
	if !strings.Contains(val, "${") && !strings.Contains(val, " + ") {
		return "", false
	}

	prefixes := [...]string{
		"bg-${", "text-${", "border-${", "p-${", "px-${", "py-${",
		"m-${", "mx-${", "my-${", "mt-${", "mb-${",
		"w-${", "h-${", "gap-${", "grid-cols-${", "flex-${",
		"rounded-${", "shadow-${", "opacity-${", "z-${",
		"\"bg-\" +", "'bg-' +", "\"text-\" +", "'text-' +", "\"p-\" +", "'p-' +",
	}

	for _, p := range prefixes {
		if strings.Contains(val, p) {
			return p, true
		}
	}

	return "", false
}

var arbitraryScaleDuplicates = map[string]string{
	"p-[16px]":         "p-4",
	"p-[1rem]":         "p-4",
	"px-[16px]":        "px-4",
	"px-[1rem]":        "px-4",
	"py-[16px]":        "py-4",
	"py-[1rem]":        "py-4",
	"m-[16px]":         "m-4",
	"m-[1rem]":         "m-4",
	"mt-[16px]":        "mt-4",
	"mt-[1rem]":        "mt-4",
	"mb-[16px]":        "mb-4",
	"mb-[1rem]":        "mb-4",
	"w-[100%]":         "w-full",
	"h-[100%]":         "h-full",
	"w-[100vw]":        "w-screen",
	"h-[100vh]":        "h-screen",
	"text-[16px]":      "text-base",
	"text-[1rem]":      "text-base",
	"rounded-[0px]":    "rounded-none",
	"rounded-[9999px]": "rounded-full",
	"gap-[16px]":       "gap-4",
	"gap-[1rem]":       "gap-4",
	"p-[0px]":          "p-0",
	"m-[0px]":          "m-0",
}

// findDuplicateArbitraryUtility memeriksa apakah ada kelas arbitrary yang menduplikasi skala core bawaan.
func findDuplicateArbitraryUtility(classes []string, rawClasses string) (string, string, bool) {
	if len(classes) == 0 || !strings.Contains(rawClasses, "-[") {
		return "", "", false
	}

	for _, c := range classes {
		if equiv, ok := arbitraryScaleDuplicates[c]; ok {
			return c, equiv, true
		}
	}

	return "", "", false
}

// stripCSSCommentsString menghapus komentar /* ... */ dari CSS string.
func stripCSSCommentsString(src string) string {
	if !strings.Contains(src, "/*") {
		return src
	}

	var buf bytes.Buffer
	buf.Grow(len(src))

	i := 0
	n := len(src)
	for i < n {
		if i+1 < n && src[i] == '/' && src[i+1] == '*' {
			end := strings.Index(src[i+2:], "*/")
			if end == -1 {
				break
			}
			i += 2 + end + 2
			continue
		}
		buf.WriteByte(src[i])
		i++
	}

	return buf.String()
}

// findUntrackedPackageSource mendeteksi berkas CSS root Tailwind v4 yang tidak memuat direktif @source saat mengimpor monorepo package.
func findUntrackedPackageSource(fileSrc string) (string, bool) {
	if len(fileSrc) == 0 {
		return "", false
	}

	clean := stripCSSCommentsString(fileSrc)
	hasImport := strings.Contains(clean, "@import \"tailwindcss\"") ||
		strings.Contains(clean, "@import 'tailwindcss'") ||
		strings.Contains(clean, "@import \"tailwindcss\";") ||
		strings.Contains(clean, "@import 'tailwindcss';")

	if !hasImport {
		return "", false
	}

	if strings.Contains(clean, "@source") {
		return "", false
	}

	return "workspace packages", true
}

// DuplicateUtilityViolation merepresentasikan deklarasi @utility yang menduplikasi utilitas native core.
type DuplicateUtilityViolation struct {
	Line        int
	UtilityName string
	CoreEquiv   string
}

// findDuplicateUtilityDefinitions memeriksa berkas CSS terhadap deklarasi @utility yang menduplikasi core.
func findDuplicateUtilityDefinitions(fileSrc string) []DuplicateUtilityViolation {
	if len(fileSrc) == 0 || !strings.Contains(fileSrc, "@utility") {
		return nil
	}

	var violations []DuplicateUtilityViolation
	lines := strings.Split(fileSrc, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "@utility") {
			continue
		}

		utilName := extractUtilityName(trimmed)
		if utilName == "" {
			continue
		}

		block := extractFunctionBlock(lines, i)
		if equiv := matchCoreUtilityEquivalent(block); equiv != "" {
			violations = append(violations, DuplicateUtilityViolation{
				Line:        i + 1,
				UtilityName: utilName,
				CoreEquiv:   equiv,
			})
		}
	}

	return violations
}

func extractUtilityName(line string) string {
	parts := strings.Fields(line)
	if len(parts) >= 2 && parts[0] == "@utility" {
		name := strings.Trim(parts[1], "{ ")
		return name
	}
	return ""
}

func matchCoreUtilityEquivalent(block string) string {
	clean := strings.ReplaceAll(block, " ", "")
	clean = strings.ReplaceAll(clean, "\t", "")
	clean = strings.ReplaceAll(clean, "\n", "")

	switch {
	case strings.Contains(clean, "display:flex") && strings.Contains(clean, "justify-content:center") && strings.Contains(clean, "align-items:center"):
		return "flex justify-center items-center"
	case strings.Contains(clean, "display:flex") && strings.Contains(clean, "align-items:center"):
		return "flex items-center"
	case strings.Contains(clean, "display:none"):
		return "hidden"
	case strings.Contains(clean, "width:100%"):
		return "w-full"
	case strings.Contains(clean, "height:100%"):
		return "h-full"
	case strings.Contains(clean, "cursor:pointer"):
		return "cursor-pointer"
	case strings.Contains(clean, "font-weight:bold"):
		return "font-bold"
	}

	return ""
}
