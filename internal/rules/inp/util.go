package inp

import (
	"strings"
	"unicode"

	"github.com/will2469/charites/internal/ir"
)

var layoutReadProps = [...]string{
	"offsetWidth",
	"offsetHeight",
	"clientWidth",
	"clientHeight",
	"scrollWidth",
	"scrollHeight",
	"getBoundingClientRect",
	"getComputedStyle",
}

var styleWritePatterns = [...]string{
	".style.",
	".style[",
	".className",
	".classList.add",
	".classList.remove",
	".classList.toggle",
}

// hasLayoutThrashingSequence mendeteksi sekuens penulisan gaya/kelas DOM yang diikuti
// langsung oleh pembacaan properti geometri layout dalam satu alur eksekusi sinkron.
func hasLayoutThrashingSequence(code string) (string, string, bool) {
	if len(code) == 0 {
		return "", "", false
	}

	writePos, writePattern := findFirstStyleWrite(code)
	if writePos == -1 {
		return "", "", false
	}

	// Cari pembacaan layout setelah posisi penulisan
	afterWrite := code[writePos+len(writePattern):]
	readPos, readPattern := findFirstLayoutRead(afterWrite)
	if readPos == -1 {
		return "", "", false
	}

	// Pengecualian: jika terdapat pembatas penjadwalan (requestAnimationFrame atau yield) di antara write dan read
	between := afterWrite[:readPos]
	if strings.Contains(between, "requestAnimationFrame") || strings.Contains(between, "scheduler.yield") {
		return "", "", false
	}

	return writePattern, readPattern, true
}

func findFirstStyleWrite(code string) (int, string) {
	minPos := -1
	matched := ""
	for _, pattern := range styleWritePatterns {
		pos := strings.Index(code, pattern)
		if pos != -1 && (minPos == -1 || pos < minPos) {
			minPos = pos
			matched = pattern
		}
	}
	return minPos, matched
}

func findFirstLayoutRead(code string) (int, string) {
	minPos := -1
	matched := ""
	for _, prop := range layoutReadProps {
		pos := strings.Index(code, prop)
		if pos != -1 && (minPos == -1 || pos < minPos) {
			minPos = pos
			matched = prop
		}
	}
	return minPos, matched
}

// isInteractiveHandlerAttr memeriksa apakah nama atribut merupakan handler event interaktif.
func isInteractiveHandlerAttr(name string) bool {
	switch name {
	case "onClick", "onChange", "onInput", "onSelect", "onKeyDown", "onKeyUp", "onKeyPress", "onPointerDown", "onPointerUp", "onSubmit":
		return true
	default:
		return false
	}
}

// hasHeavySynchronousOps memeriksa apakah kode handler memuat operasi komputasi berat tanpa yield kooperatif.
func hasHeavySynchronousOps(code string) (string, bool) {
	if len(code) == 0 {
		return "", false
	}

	// Pengecualian: jika handler mendelegasikan ke Web Worker atau memuat scheduler.yield / scheduler?.yield / startTransition
	if (strings.Contains(code, "scheduler") && strings.Contains(code, "yield")) ||
		strings.Contains(code, "Worker") ||
		strings.Contains(code, "startTransition") {
		return "", false
	}

	if strings.Contains(code, "JSON.parse(") {
		return "JSON.parse", true
	}
	if strings.Contains(code, ".sort(") {
		return "Array.sort", true
	}
	if strings.Contains(code, "JSON.stringify(") && (strings.Contains(code, "for") || strings.Contains(code, "map")) {
		return "JSON.stringify in loop", true
	}
	if isNestedLoop(code) {
		return "nested loops", true
	}

	return "", false
}

func isNestedLoop(code string) bool {
	firstLoop := strings.Index(code, "for ")
	if firstLoop == -1 {
		firstLoop = strings.Index(code, "for(")
	}
	if firstLoop == -1 {
		return false
	}
	rest := code[firstLoop+4:]
	secondLoop := strings.Index(rest, "for ")
	if secondLoop == -1 {
		secondLoop = strings.Index(rest, "for(")
	}
	return secondLoop != -1
}

// hasRepeatedStateUpdateInLoop memeriksa apakah terdapat pembaruan state berulang di dalam badan perulangan
// yang memecah batas batching React 18+ (memuat await atau flushSync di dalam loop).
func hasRepeatedStateUpdateInLoop(code string) (string, bool) {
	if len(code) == 0 {
		return "", false
	}

	loops := extractLoopBodies(code)
	for _, body := range loops {
		hasAwait := strings.Contains(body, "await ")
		hasFlushSync := strings.Contains(body, "flushSync(")
		if !hasAwait && !hasFlushSync {
			continue
		}

		setter := findStateSetterCall(body)
		if setter != "" {
			return setter, true
		}
	}

	return "", false
}

func extractLoopBodies(code string) []string {
	var bodies []string
	idx := 0
	for {
		forPos := strings.Index(code[idx:], "for")
		whilePos := strings.Index(code[idx:], "while")
		pos := -1
		prefixLen := 3
		if forPos != -1 && (whilePos == -1 || forPos < whilePos) {
			pos = forPos
			prefixLen = 3
		} else if whilePos != -1 {
			pos = whilePos
			prefixLen = 5
		}
		if pos == -1 {
			break
		}
		start := idx + pos + prefixLen
		braceOpen := strings.IndexByte(code[start:], '{')
		if braceOpen == -1 {
			idx = start
			continue
		}
		bodyStart := start + braceOpen + 1
		bodyEnd := findMatchingBrace(code[bodyStart:])
		if bodyEnd == -1 {
			break
		}
		bodies = append(bodies, code[bodyStart:bodyStart+bodyEnd])
		idx = bodyStart + bodyEnd + 1
	}
	return bodies
}

func findMatchingBrace(s string) int {
	depth := 1
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func findStateSetterCall(code string) string {
	if strings.Contains(code, "setState(") {
		return "setState"
	}
	idx := 0
	for {
		pos := strings.Index(code[idx:], "set")
		if pos == -1 {
			break
		}
		start := idx + pos
		if start+3 < len(code) {
			nextChar := rune(code[start+3])
			if unicode.IsUpper(nextChar) {
				// Cari penutup nama fungsi setFoo(
				paren := strings.IndexByte(code[start:], '(')
				if paren != -1 && paren < 40 {
					name := code[start : start+paren]
					if isValidIdentifier(name) {
						return name
					}
				}
			}
		}
		idx = start + 3
	}
	return ""
}

func isValidIdentifier(name string) bool {
	for _, ch := range name {
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' {
			return false
		}
	}
	return true
}

// hasUnyieldedLongLoop memeriksa fungsi komputasi panjang yang memproses banyak item tanpa batas penjadwalan kooperatif.
func hasUnyieldedLongLoop(code string) (string, bool) {
	if len(code) == 0 {
		return "", false
	}

	hasLoop := strings.Contains(code, "for ") || strings.Contains(code, "for(") || strings.Contains(code, "while")
	if !hasLoop {
		return "", false
	}

	// Pengecualian: jika loop menyertakan yield kooperatif
	if (strings.Contains(code, "scheduler") && strings.Contains(code, "yield")) ||
		strings.Contains(code, "setTimeout") ||
		strings.Contains(code, "requestAnimationFrame") {
		return "", false
	}

	// Deteksi loop komputasi berat
	if strings.Contains(code, "heavyCalculation") || strings.Contains(code, "processLarge") ||
		strings.Contains(code, "items.length") || isNestedLoop(code) {
		return "unyielded loop", true
	}

	return "", false
}

// extractScriptNodeText mengekstrak seluruh teks script dari child node NodeText dalam <script>.
func extractScriptNodeText(node *ir.Node) string {
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

// countClientLoadIslands menghitung jumlah pulau client:load dalam seluruh pohon AST root
// dan mengembalikan simpul client:load pertama yang ditemukan.
func countClientLoadIslands(root *ir.Node) (*ir.Node, int) {
	if root == nil {
		return nil, 0
	}

	var firstNode *ir.Node
	count := 0

	for n := range root.Walk() {
		if n.Type != ir.NodeElement {
			continue
		}
		if _, ok := n.Attributes["client:load"]; ok {
			count++
			if firstNode == nil {
				firstNode = n
			}
		}
	}

	return firstNode, count
}

var exemptIslandKeywords = [...]string{
	"editor",
	"chart",
	"graph",
	"canvas",
	"richtext",
	"codemirror",
	"monaco",
	"terminal",
	"map",
}

func containsSubstringIgnoreCase(s, substrLower string) bool {
	if len(substrLower) == 0 {
		return true
	}
	if len(s) < len(substrLower) {
		return false
	}
	subLen := len(substrLower)
	maxIdx := len(s) - subLen
	for i := 0; i <= maxIdx; i++ {
		match := true
		for j := 0; j < subLen; j++ {
			c := s[i+j]
			if c >= 'A' && c <= 'Z' {
				c += 32
			}
			if c != substrLower[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func isExemptIslandTag(tag string) bool {
	for _, kw := range exemptIslandKeywords {
		if containsSubstringIgnoreCase(tag, kw) {
			return true
		}
	}
	return false
}

func isClientDirective(attrName string) bool {
	return strings.HasPrefix(attrName, "client:")
}

func hasClientDirective(node *ir.Node) bool {
	for k := range node.Attributes {
		if isClientDirective(k) {
			return true
		}
	}
	return false
}

var staticElementTags = [...]string{
	"p", "span", "h1", "h2", "h3", "h4", "h5", "h6",
	"article", "section", "blockquote", "ul", "ol", "li",
	"table", "tbody", "tr", "td", "header", "footer", "aside", "nav", "main",
}

func isStaticElementTag(tag string) bool {
	for _, t := range staticElementTags {
		if strings.EqualFold(tag, t) {
			return true
		}
	}
	return strings.HasPrefix(tag, "Static") || strings.HasPrefix(tag, "static") ||
		strings.HasPrefix(tag, "Text") || strings.HasPrefix(tag, "text")
}

// isHydrationHeavyIsland mengevaluasi apakah pulau client membungkus pohon sub-elemen statis yang masif.
func isHydrationHeavyIsland(node *ir.Node) (int, bool) {
	if node == nil || node.Type != ir.NodeElement {
		return 0, false
	}

	if !hasClientDirective(node) {
		return 0, false
	}

	if isExemptIslandTag(node.Tag) {
		return 0, false
	}

	totalNodes := 0
	staticCount := 0

	for child := range node.Walk() {
		if child == node {
			continue
		}
		totalNodes++
		if child.Type == ir.NodeElement && isStaticElementTag(child.Tag) {
			staticCount++
		}
	}

	if totalNodes >= 12 || (totalNodes >= 8 && staticCount >= 4) {
		return totalNodes, true
	}

	return 0, false
}

// isRenderBlockingScript memeriksa apakah elemen script eksternal memblokir rendering sinkron.
func isRenderBlockingScript(node *ir.Node) bool {
	if node == nil || node.Type != ir.NodeElement || !strings.EqualFold(node.Tag, "script") {
		return false
	}

	srcVal, hasSrc := node.Attributes["src"]
	if !hasSrc || strings.TrimSpace(srcVal) == "" {
		return false
	}

	// Pengecualian non-blocking: defer, async, type="module", type="application/ld+json", type="text/partytown"
	if _, hasDefer := node.Attributes["defer"]; hasDefer {
		return false
	}
	if _, hasAsync := node.Attributes["async"]; hasAsync {
		return false
	}
	if typeVal, hasType := node.Attributes["type"]; hasType {
		if strings.EqualFold(typeVal, "module") ||
			strings.EqualFold(typeVal, "application/ld+json") ||
			strings.EqualFold(typeVal, "text/partytown") {
			return false
		}
	}

	// Dalam Astro, script tanpa is:inline diperlakukan sebagai ESM oleh bundler
	// Jika ada is:inline, script dieksekusi mentah dan memblokir thread jika eksternal
	if _, hasIsInline := node.Attributes["is:inline"]; hasIsInline {
		return true
	}

	// Jika bukan inline tapi merupakan URL eksternal (http/https/cdn)
	if strings.HasPrefix(srcVal, "http://") || strings.HasPrefix(srcVal, "https://") || strings.HasPrefix(srcVal, "//") {
		return true
	}

	return false
}

var urgentSetterKeywords = [...]string{
	"setsearchquery",
	"setquery",
	"setinput",
	"settext",
	"setvalue",
	"setsearch",
	"setkeyword",
	"setterm",
}

var secondaryHeavySetterKeywords = [...]string{
	"setfiltered",
	"setresults",
	"setsearchresults",
	"setlist",
	"setitems",
	"setchartdata",
	"expensivefilter",
	"filteritems",
}

// hasMissingStartTransition memeriksa apakah handler interaksi memicu pembaruan state sekunder berat
// bersamaan dengan input mendesak tanpa dibungkus startTransition.
func hasMissingStartTransition(code string) (string, bool) {
	if len(code) == 0 {
		return "", false
	}

	// Fast pre-filter: jika tidak ada 'target' dan 'set', tidak mungkin ada benturan state update
	if !strings.Contains(code, "target") && !strings.Contains(code, "Target") &&
		!strings.Contains(code, "set") && !strings.Contains(code, "Set") {
		return "", false
	}

	// Pengecualian jika sudah menggunakan startTransition, useTransition, atau useDeferredValue
	if strings.Contains(code, "startTransition") || strings.Contains(code, "useTransition") ||
		strings.Contains(code, "useDeferredValue") || strings.Contains(code, "debounce") {
		return "", false
	}

	hasUrgent := containsSubstringIgnoreCase(code, "target.value")
	if !hasUrgent {
		for _, kw := range urgentSetterKeywords {
			if containsSubstringIgnoreCase(code, kw) {
				hasUrgent = true
				break
			}
		}
	}

	if !hasUrgent {
		return "", false
	}

	for _, heavy := range secondaryHeavySetterKeywords {
		if containsSubstringIgnoreCase(code, heavy) {
			return heavy, true
		}
	}

	return "", false
}

// getScriptOrSourceContent mengambil konten script baik dari elemen <script> maupun dari NodeComment sumber file.
func getScriptOrSourceContent(node *ir.Node) string {
	if node == nil {
		return ""
	}
	if node.Type == ir.NodeElement && strings.EqualFold(node.Tag, "script") {
		return extractScriptNodeText(node)
	}
	if node.Type == ir.NodeComment && len(node.RawClasses) > 0 {
		return node.RawClasses
	}
	if node.Type == ir.NodeElement {
		cur := node
		for cur.Parent != nil {
			cur = cur.Parent
		}
		for _, ch := range cur.Children {
			if ch != nil && ch.Type == ir.NodeComment && len(ch.RawClasses) > 0 {
				return ch.RawClasses
			}
		}
	}
	return ""
}

// findUnboundedEffectDeps mencari pemanggilan useEffect atau useLayoutEffect yang kekurangan array dependensi.
func findUnboundedEffectDeps(code string) (string, int, bool) {
	if len(code) == 0 {
		return "", 0, false
	}
	if !strings.Contains(code, "useEffect") && !strings.Contains(code, "useLayoutEffect") {
		return "", 0, false
	}

	for _, hook := range [...]string{"useLayoutEffect", "useEffect"} {
		idx := 0
		for {
			pos := strings.Index(code[idx:], hook+"(")
			if pos == -1 {
				break
			}
			absPos := idx + pos
			if isHookIdentifierBoundary(code, absPos) {
				argStart := absPos + len(hook) + 1
				if hasSecondArg, lineOffset := checkHookArgs(code[argStart:]); !hasSecondArg {
					line := 1 + strings.Count(code[:absPos], "\n") + lineOffset
					return hook, line, true
				}
			}
			idx = absPos + len(hook) + 1
		}
	}
	return "", 0, false
}

func isHookIdentifierBoundary(code string, pos int) bool {
	if pos == 0 {
		return true
	}
	prev := code[pos-1]
	return !((prev >= 'a' && prev <= 'z') || (prev >= 'A' && prev <= 'Z') || (prev >= '0' && prev <= '9') || prev == '_')
}

func checkHookArgs(argsCode string) (bool, int) {
	parenDepth := 1
	braceDepth := 0
	bracketDepth := 0
	hasComma := false
	lineOffset := 0

	for i := 0; i < len(argsCode); i++ {
		c := argsCode[i]
		if c == '\n' {
			lineOffset++
		}
		switch c {
		case '(':
			parenDepth++
		case ')':
			parenDepth--
			if parenDepth == 0 {
				return hasComma, lineOffset
			}
		case '{':
			braceDepth++
		case '}':
			braceDepth--
		case '[':
			bracketDepth++
		case ']':
			bracketDepth--
		case ',':
			if parenDepth == 1 && braceDepth == 0 && bracketDepth == 0 {
				hasComma = true
			}
		}
	}
	return hasComma, lineOffset
}

// findSyncLayoutEffectViolation memeriksa apakah terdapat komputasi non-geometris berat di dalam useLayoutEffect.
func findSyncLayoutEffectViolation(code string) (string, int, bool) {
	if len(code) == 0 || !strings.Contains(code, "useLayoutEffect") {
		return "", 0, false
	}

	hook := "useLayoutEffect"
	idx := 0
	for {
		pos := strings.Index(code[idx:], hook+"(")
		if pos == -1 {
			break
		}
		absPos := idx + pos
		if isHookIdentifierBoundary(code, absPos) {
			argStart := absPos + len(hook) + 1
			body := extractHookBody(code[argStart:])
			if op, isHeavy := isHeavyLayoutEffectBody(body); isHeavy {
				line := 1 + strings.Count(code[:absPos], "\n")
				return op, line, true
			}
		}
		idx = absPos + len(hook) + 1
	}
	return "", 0, false
}

func extractHookBody(code string) string {
	braceStart := strings.IndexByte(code, '{')
	if braceStart == -1 {
		return ""
	}
	end := findMatchingBrace(code[braceStart+1:])
	if end == -1 {
		return ""
	}
	return code[braceStart+1 : braceStart+1+end]
}

var heavyLayoutOps = [...]string{
	"fetch(",
	"axios",
	"localStorage",
	"sessionStorage",
	"JSON.parse(",
	"async ",
}

var geomMeasurementProps = [...]string{
	"getBoundingClientRect",
	"offsetHeight",
	"offsetWidth",
	"clientHeight",
	"clientWidth",
	"scrollHeight",
	"scrollWidth",
}

func isHeavyLayoutEffectBody(body string) (string, bool) {
	if len(body) == 0 {
		return "", false
	}
	// Pengecualian: jika melakukan pengukuran geometri DOM murni
	for _, geom := range geomMeasurementProps {
		if strings.Contains(body, geom) {
			return "", false
		}
	}

	for _, op := range heavyLayoutOps {
		if strings.Contains(body, op) {
			return op, true
		}
	}

	if strings.Contains(body, "setData") || strings.Contains(body, "setUsers") || strings.Contains(body, "setItems") {
		return "state update without DOM measurement", true
	}

	return "", false
}

// isInlineContextObjectValue memeriksa apakah elemen Context.Provider menerima objek literal inline pada prop value.
func isInlineContextObjectValue(tag string, attrs map[string]string) (string, bool) {
	if attrs == nil {
		return "", false
	}
	if tag != "Context.Provider" && !strings.HasSuffix(tag, ".Provider") {
		return "", false
	}

	val, hasVal := attrs["value"]
	if !hasVal {
		return "", false
	}

	trimmed := strings.TrimSpace(val)
	if strings.HasPrefix(trimmed, "{{") {
		if strings.Contains(trimmed, "useMemo") {
			return "", false
		}
		return tag, true
	}
	if strings.HasPrefix(trimmed, "{") {
		rest := strings.TrimSpace(trimmed[1:])
		if strings.HasPrefix(rest, "{") && !strings.Contains(trimmed, "useMemo") {
			return tag, true
		}
	}

	return "", false
}

// hasExpensiveRenderComputation mendeteksi operasi komputasi berat (seperti chained .filter().sort())
// yang dieksekusi langsung di jalur render tanpa useMemo.
func hasExpensiveRenderComputation(code string) (string, bool) {
	if len(code) == 0 {
		return "", false
	}
	if !strings.Contains(code, ".filter") && !strings.Contains(code, ".sort") {
		return "", false
	}
	if strings.Contains(code, "useMemo") {
		return "", false
	}

	filterPos := strings.Index(code, ".filter(")
	if filterPos != -1 && strings.Contains(code[filterPos:], ".sort(") {
		return ".filter().sort()", true
	}
	sortPos := strings.Index(code, ".sort(")
	if sortPos != -1 && strings.Contains(code[sortPos:], ".filter(") {
		return ".sort().filter()", true
	}

	return "", false
}

// isScrollableContainer memeriksa apakah elemen memiliki kelas CSS kontainer scroll (overflow-y-auto, dll).
func isScrollableContainer(_ string, classes []string, rawClasses string) bool {
	if len(rawClasses) == 0 && len(classes) == 0 {
		return false
	}
	if strings.Contains(rawClasses, "overflow-y-auto") ||
		strings.Contains(rawClasses, "overflow-auto") ||
		strings.Contains(rawClasses, "overflow-scroll") ||
		strings.Contains(rawClasses, "overflow-y-scroll") {
		return true
	}
	for _, cls := range classes {
		if cls == "overflow-y-auto" || cls == "overflow-auto" || cls == "overflow-scroll" || cls == "overflow-y-scroll" {
			return true
		}
	}
	return false
}

// hasUnboundedCollectionMapping mendeteksi apakah kode sumber di sekitar kontainer scroll
// memetakan koleksi dinamis via .map(...) tanpa paginasi (.slice) atau pustaka virtualisasi.
func hasUnboundedCollectionMapping(src string, startLine int) (string, bool) {
	if len(src) == 0 {
		return "", false
	}

	lines := strings.Split(src, "\n")
	startIdx := startLine - 1
	if startIdx < 0 {
		startIdx = 0
	}
	endIdx := startIdx + 40
	if endIdx > len(lines) {
		endIdx = len(lines)
	}

	for i := startIdx; i < endIdx; i++ {
		line := lines[i]
		if !strings.Contains(line, ".map(") {
			continue
		}

		// Pengecualian 1: Pustaka virtualisasi list
		if strings.Contains(line, "virtual") || strings.Contains(line, "Virtual") ||
			strings.Contains(src, "useVirtualizer") || strings.Contains(src, "getVirtualItems") {
			return "", false
		}

		// Pengecualian 2: Paginasi atau pembatasan slice tetap
		if strings.Contains(line, ".slice(") || strings.Contains(line, ".take(") {
			return "", false
		}

		// Pengecualian 3: Static constant array mapping e.g. [1, 2, 3].map
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "[") && strings.Contains(trimmed, "].map(") {
			return "", false
		}

		mapIdx := strings.Index(line, ".map(")
		expr := line[:mapIdx+len(".map()")]
		if len(expr) > 40 {
			expr = expr[len(expr)-40:]
		}
		return strings.TrimSpace(expr), true
	}

	return "", false
}

// isUnconstrainedLayoutScope mendeteksi elemen overlay/drawer/modal besar tanpa isolasi containment
// yang memicu invalidasi layout dokumen secara menyeluruh saat dibuka/ditutup.
func isUnconstrainedLayoutScope(tag string, classes []string, rawClasses string) bool {
	if strings.EqualFold(tag, "dialog") {
		return false
	}
	if len(rawClasses) == 0 && len(classes) == 0 {
		return false
	}

	// Pengecualian: Telah memiliki layout containment eksplisit
	if strings.Contains(rawClasses, "contain-layout") ||
		strings.Contains(rawClasses, "contain-strict") ||
		strings.Contains(rawClasses, "contain-content") ||
		strings.Contains(rawClasses, "contain:layout") ||
		strings.Contains(rawClasses, "contain:strict") {
		return false
	}

	// Wajib berposisi fixed (overlay)
	hasFixed := strings.Contains(rawClasses, "fixed")
	if !hasFixed {
		for _, cls := range classes {
			if cls == "fixed" {
				hasFixed = true
				break
			}
		}
	}
	if !hasFixed {
		return false
	}

	// Memiliki dimensi overlay besar
	hasLargeBounds := strings.Contains(rawClasses, "inset-0") ||
		strings.Contains(rawClasses, "inset-y-0") ||
		strings.Contains(rawClasses, "inset-x-0") ||
		strings.Contains(rawClasses, "w-screen") ||
		strings.Contains(rawClasses, "h-screen") ||
		strings.Contains(rawClasses, "w-96") ||
		strings.Contains(rawClasses, "z-50")
	if !hasLargeBounds {
		return false
	}

	// Memiliki indikator dynamic drawer/modal/overlay atau conditional toggle
	hasOverlayIntent := strings.Contains(rawClasses, "hidden") ||
		strings.Contains(rawClasses, "block") ||
		strings.Contains(rawClasses, "isOpen") ||
		strings.Contains(rawClasses, "open") ||
		strings.Contains(rawClasses, "drawer") ||
		strings.Contains(rawClasses, "modal") ||
		strings.Contains(rawClasses, "sidebar") ||
		strings.Contains(rawClasses, "sheet")
	return hasOverlayIntent
}

// getMissingTouchAction memeriksa apakah elemen menangani gestur sentuhan/pointer
// kustom tanpa mendeklarasikan kebijakan CSS touch-action.
func getMissingTouchAction(_ string, attrs map[string]string, classes []string, rawClasses string) (string, bool) {
	if attrs == nil {
		return "", false
	}

	var gestureHandler string
	for k := range attrs {
		if strings.EqualFold(k, "onPointerDown") ||
			strings.EqualFold(k, "onTouchStart") ||
			strings.EqualFold(k, "onPointerMove") ||
			strings.EqualFold(k, "onTouchMove") {
			gestureHandler = k
			break
		}
	}
	if gestureHandler == "" {
		return "", false
	}

	// Periksa apakah telah memiliki utilitas touch-*
	if strings.Contains(rawClasses, "touch-none") ||
		strings.Contains(rawClasses, "touch-pan-x") ||
		strings.Contains(rawClasses, "touch-pan-y") ||
		strings.Contains(rawClasses, "touch-pan-left") ||
		strings.Contains(rawClasses, "touch-pan-right") ||
		strings.Contains(rawClasses, "touch-pan-up") ||
		strings.Contains(rawClasses, "touch-pan-down") ||
		strings.Contains(rawClasses, "touch-pinch-zoom") ||
		strings.Contains(rawClasses, "touch-manipulation") ||
		strings.Contains(rawClasses, "touch-auto") {
		return "", false
	}

	for _, cls := range classes {
		if strings.HasPrefix(cls, "touch-") {
			return "", false
		}
	}

	if styleVal, ok := attrs["style"]; ok && (strings.Contains(styleVal, "touch-action") || strings.Contains(styleVal, "touchAction")) {
		return "", false
	}

	return gestureHandler, true
}

// findExpensiveStyleMutation mendeteksi mutasi gaya peka-cat imperatif (seperti boxShadow atau filter)
// di dalam penangan interaksi kontinu (onPointerMove, onTouchMove, onScroll, onMouseMove, onWheel).
func findExpensiveStyleMutation(_ string, attrs map[string]string) (string, string, bool) {
	if attrs == nil {
		return "", "", false
	}

	var continuousHandler string
	var handlerCode string
	for k, v := range attrs {
		if strings.EqualFold(k, "onPointerMove") ||
			strings.EqualFold(k, "onTouchMove") ||
			strings.EqualFold(k, "onScroll") ||
			strings.EqualFold(k, "onMouseMove") ||
			strings.EqualFold(k, "onWheel") {
			continuousHandler = k
			handlerCode = v
			break
		}
	}
	if continuousHandler == "" || len(handlerCode) == 0 {
		return "", "", false
	}

	paintProps := [...]string{
		"boxShadow",
		"box-shadow",
		"filter",
		"backdropFilter",
		"backdrop-filter",
		"backgroundImage",
		"background-image",
		"borderImage",
		"mixBlendMode",
	}

	for _, prop := range paintProps {
		if strings.Contains(handlerCode, "style."+prop) ||
			strings.Contains(handlerCode, "style['"+prop+"']") ||
			strings.Contains(handlerCode, "style[\""+prop+"\"]") {
			return continuousHandler, prop, true
		}
	}

	return "", "", false
}
