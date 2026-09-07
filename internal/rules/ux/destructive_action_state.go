package ux

import (
	"regexp"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// StateGateResult mendefinisikan status hasil analisis state-gating konfirmasi.
type StateGateResult int

const (
	// StateGateUnknown berarti bukti state-gating tidak dapat dibuktikan atau berada di luar scope komponen lokal.
	// Kebijakan rule: Unknown diperlakukan sebagai potensi bahaya (tetap emit diagnostik).
	StateGateUnknown StateGateResult = iota

	// StateGateConfirmed berarti aksi destruktif terbukti dilindungi oleh state konfirmasi lokal,
	// dikonsumsi oleh dialog konfirmasi hilir, dan dialog tersebut memiliki aksi konfirmasi.
	StateGateConfirmed

	// StateGateNotConfirmed berarti aksi destruktif terbukti tidak dilindungi konfirmasi
	// (misal mutasi langsung dipanggil, atau dialog tidak memiliki aksi konfirmasi).
	StateGateNotConfirmed
)

// StateBinding merepresentasikan deklarasi useState lokal.
type StateBinding struct {
	StateName  string
	SetterName string
}

var (
	reUseStateDecl = regexp.MustCompile(`(?:const|let|var)\s*\[\s*([a-zA-Z0-9_$]+)\s*,\s*([a-zA-Z0-9_$]+)\s*\]\s*=\s*useState(?:<[^>]*>)?\s*\(`)
	reFuncCall     = regexp.MustCompile(`([a-zA-Z0-9_$.]+)\s*\(`)
)

// AnalyzeStateGating menganalisis apakah elemen interaktif destruktif dilindungi oleh
// alur konfirmasi bertahap (staged confirmation) berbasis useState dalam komponen lokal yang sama.
func AnalyzeStateGating(node *ir.Node) StateGateResult {
	if node == nil || node.Type != ir.NodeElement {
		return StateGateUnknown
	}

	handler := extractHandlerString(node)
	if handler == "" {
		return StateGateNotConfirmed
	}

	callees := extractCallees(handler)
	if hasDirectDestructiveMutationCallee(callees) {
		return StateGateNotConfirmed
	}

	// 2. Ekstraksi kode sumber berkas untuk analisis intra-component scope.
	sourceCode := getFileSourceContent(node)
	if sourceCode == "" {
		return StateGateUnknown
	}

	// 3. Batasi scope ke fungsi/komponen terdekat yang membungkus node.
	scopeStart, scopeEnd, foundScope := findEnclosingComponentScope(sourceCode, node.Span.Line)
	if !foundScope {
		return StateGateUnknown
	}

	// 4. Ekstraksi seluruh deklarasi useState di dalam batas scope komponen lokal.
	scopeText := extractLineRange(sourceCode, scopeStart, scopeEnd)
	bindings := extractStateBindings(scopeText)

	// 5. Identifikasi apakah handler memanggil setter state lokal.
	matchedBinding, setterCallFound := resolveInvokedStateBinding(callees, bindings)
	if !setterCallFound {
		// Handler tidak memanggil setter apapun.
		return StateGateNotConfirmed
	}

	if matchedBinding == nil {
		// Handler memanggil fungsi setX/updateX yang TIDAK dideklarasikan di komponen ini (misal dari props).
		// Sesuai invariant scope: setter di luar scope lokal -> StateGateUnknown.
		return StateGateUnknown
	}

	// 6. Cari elemen dialog konfirmasi di dalam pohon komponen lokal yang sama.
	root := getRootNode(node)
	if root == nil {
		return StateGateUnknown
	}

	confirmed := verifyDownstreamConfirmationConsumer(root, scopeStart, scopeEnd, scopeText, matchedBinding.StateName)
	if confirmed {
		return StateGateConfirmed
	}

	return StateGateNotConfirmed
}

// extractHandlerString mengambil isi atribut handler event (onClick, onPress).
func extractHandlerString(node *ir.Node) string {
	if node == nil || node.Attributes == nil {
		return ""
	}
	for k, v := range node.Attributes {
		kLower := strings.ToLower(k)
		if kLower == "onclick" || kLower == "onpress" {
			return v
		}
	}
	return ""
}

// extractCallees mengekstrak nama-nama fungsi atau metode yang dipanggil dalam string handler.
func extractCallees(handler string) []string {
	matches := reFuncCall.FindAllStringSubmatch(handler, -1)
	if len(matches) == 0 {
		return nil
	}
	callees := make([]string, 0, len(matches))
	for _, m := range matches {
		if len(m) > 1 {
			callees = append(callees, m[1])
		}
	}
	return callees
}

// hasDirectDestructiveMutationCallee mengecek apakah ada pemanggilan fungsi yang merupakan mutasi destruktif langsung.
func hasDirectDestructiveMutationCallee(callees []string) bool {
	for _, callee := range callees {
		if isDirectDestructiveMutationCallee(callee) {
			return true
		}
	}
	return false
}

// isDirectDestructiveMutationCallee membedakan aksi mutasi destruktif langsung dari setter state lokal.
func isDirectDestructiveMutationCallee(callee string) bool {
	// Setter state (setX, updateX) bukan mutasi langsung
	if isLikelyStateSetterName(callee) {
		return false
	}

	// Hilangkan prefix objek jika ada (misal: api.delete -> delete)
	lastDot := strings.LastIndexByte(callee, '.')
	target := callee
	if lastDot != -1 && lastDot+1 < len(callee) {
		target = callee[lastDot+1:]
	}

	targetLower := strings.ToLower(target)
	destructiveKeywords := [...]string{"delete", "remove", "destroy", "purge", "revoke", "hapus"}
	for _, kw := range destructiveKeywords {
		if strings.Contains(targetLower, kw) {
			return true
		}
	}

	if targetLower == "mutate" || strings.HasPrefix(targetLower, "mutate") {
		return true
	}

	return false
}

// isLikelyStateSetterName memeriksa apakah identifier adalah konvensi setter state.
func isLikelyStateSetterName(name string) bool {
	lastDot := strings.LastIndexByte(name, '.')
	if lastDot != -1 && lastDot+1 < len(name) {
		name = name[lastDot+1:]
	}

	if strings.HasPrefix(name, "set") && len(name) > 3 {
		nextChar := name[3]
		if nextChar >= 'A' && nextChar <= 'Z' {
			return true
		}
	}
	if strings.HasPrefix(name, "update") && len(name) > 6 {
		nextChar := name[6]
		if nextChar >= 'A' && nextChar <= 'Z' {
			return true
		}
	}
	return false
}

// extractStateBindings mengekstrak pasangan [state, setter] dari deklarasi useState.
func extractStateBindings(code string) []StateBinding {
	matches := reUseStateDecl.FindAllStringSubmatch(code, -1)
	if len(matches) == 0 {
		return nil
	}
	bindings := make([]StateBinding, 0, len(matches))
	for _, m := range matches {
		if len(m) > 2 {
			bindings = append(bindings, StateBinding{
				StateName:  m[1],
				SetterName: m[2],
			})
		}
	}
	return bindings
}

// resolveInvokedStateBinding mencocokkan callee dalam handler dengan deklarasi useState lokal.
func resolveInvokedStateBinding(callees []string, bindings []StateBinding) (*StateBinding, bool) {
	for _, callee := range callees {
		// Cek apakah callee persis mencocokkan setter yang dideklarasikan
		for i := range bindings {
			if bindings[i].SetterName == callee {
				return &bindings[i], true
			}
		}
	}

	// Cek apakah ada setter yang dipanggil tetapi tidak dideklarasikan di sini
	for _, callee := range callees {
		if isLikelyStateSetterName(callee) {
			return nil, true // setter dipanggil, tetapi tidak bound di scope lokal
		}
	}

	return nil, false
}

// getFileSourceContent mengambil seluruh konten file dari NodeComment di root AST.
func getFileSourceContent(node *ir.Node) string {
	if node == nil {
		return ""
	}
	root := getRootNode(node)
	if root == nil {
		return ""
	}
	for _, ch := range root.Children {
		if ch != nil && ch.Type == ir.NodeComment && len(ch.RawClasses) > 0 {
			return ch.RawClasses
		}
	}
	return ""
}

// getRootNode mengembalikan root node tertinggi dari hierarki AST.
func getRootNode(node *ir.Node) *ir.Node {
	cur := node
	for cur != nil && cur.Parent != nil {
		cur = cur.Parent
	}
	return cur
}

// findEnclosingComponentScope mendeteksi batas baris [startLine, endLine] fungsi/komponen terdekat.
func findEnclosingComponentScope(source string, targetLine int) (int, int, bool) {
	lines := strings.Split(source, "\n")
	totalLines := len(lines)
	if targetLine < 1 || targetLine > totalLines {
		return 1, totalLines, true
	}

	scopeStart := findComponentDeclarationStart(lines, targetLine-1)
	if scopeStart == -1 {
		return 1, totalLines, true
	}

	scopeEnd := findComponentScopeEnd(lines, scopeStart)
	return scopeStart, scopeEnd, true
}

func findComponentDeclarationStart(lines []string, targetIdx int) int {
	for i := targetIdx; i >= 0; i-- {
		lineStr := strings.TrimSpace(lines[i])
		if isComponentFunctionDeclaration(lineStr) {
			return i + 1
		}
	}
	return -1
}

func findComponentScopeEnd(lines []string, scopeStart int) int {
	totalLines := len(lines)
	openParens := 0
	openBraces := 0
	bodyStarted := false

	for i := scopeStart - 1; i < totalLines; i++ {
		lineStr := lines[i]
		for j := 0; j < len(lineStr); j++ {
			c := lineStr[j]
			if !bodyStarted {
				openParens, bodyStarted, openBraces = updateParameterAndBodyState(c, openParens)
				continue
			}

			switch c {
			case '{':
				openBraces++
			case '}':
				openBraces--
				if openBraces == 0 {
					return i + 1
				}
			}
		}
	}

	return totalLines
}

func updateParameterAndBodyState(c byte, openParens int) (int, bool, int) {
	switch c {
	case '(':
		return openParens + 1, false, 0
	case ')':
		if openParens > 0 {
			openParens--
		}
		return openParens, false, 0
	case '{':
		if openParens == 0 {
			return 0, true, 1
		}
	}
	return openParens, false, 0
}

func isComponentFunctionDeclaration(line string) bool {
	if strings.HasPrefix(line, "function ") || strings.HasPrefix(line, "export function ") ||
		strings.HasPrefix(line, "export default function") {
		return true
	}
	if strings.Contains(line, "const ") || strings.Contains(line, "let ") {
		if strings.Contains(line, "=>") || strings.Contains(line, "function") {
			return true
		}
	}
	return false
}

// extractLineRange mengambil teks dari baris startLine hingga endLine (1-indexed).
func extractLineRange(source string, startLine, endLine int) string {
	lines := strings.Split(source, "\n")
	total := len(lines)
	if startLine < 1 {
		startLine = 1
	}
	if endLine > total {
		endLine = total
	}
	if startLine > endLine {
		return ""
	}
	return strings.Join(lines[startLine-1:endLine], "\n")
}

// verifyDownstreamConfirmationConsumer mencari apakah state dikonsumsi oleh primitif konfirmasi hilir.
func verifyDownstreamConfirmationConsumer(root *ir.Node, scopeStart, scopeEnd int, scopeText, stateName string) bool {
	if root == nil {
		return false
	}

	// 1. Cek conditional rendering guard di dalam teks komponen:
	// misal: {confirmIndex !== null && <ActionApprovalDialog ... />}
	// atau:  {pendingDeleteId && <ConfirmDialog ... />}
	hasGuardedJSX := checkConditionalRenderGuard(scopeText, stateName)

	// 2. Telusuri seluruh elemen di dalam komponen lokal
	for node := range root.Walk() {
		if node.Type != ir.NodeElement {
			continue
		}

		// Batasi hanya node dalam scope baris komponen yang sama
		if node.Span.Line < scopeStart || node.Span.Line > scopeEnd {
			continue
		}

		if !isConfirmationPrimitive(node) {
			continue
		}

		// Cek apakah dialog ini mengonsumsi stateName via props visibility atau via guarded JSX
		consumesState := consumesStateVariable(node, stateName) || (hasGuardedJSX && isNearGuardedScope(node, scopeText, stateName))
		if !consumesState {
			continue
		}

		// Verifikasi bahwa dialog memiliki aksi konfirmasi hilir (bukan sekadar dialog informasional)
		if hasDownstreamConfirmationAction(node) {
			return true
		}
	}

	return false
}

// isConfirmationPrimitive memeriksa apakah tag node adalah primitif dialog atau konfirmasi.
func isConfirmationPrimitive(node *ir.Node) bool {
	if node == nil {
		return false
	}
	tag := node.Tag
	tagLower := strings.ToLower(tag)

	if tagLower == "dialog" || tagLower == "alertdialog" {
		return true
	}
	if strings.HasSuffix(tag, "Dialog") || strings.HasSuffix(tag, "Modal") || strings.HasSuffix(tag, "Drawer") {
		return true
	}
	if role, ok := node.GetAttr("role"); ok {
		r := strings.ToLower(cleanAttrValue(role))
		if r == "dialog" || r == "alertdialog" {
			return true
		}
	}
	return false
}

// consumesStateVariable memeriksa apakah elemen mengonsumsi variabel state dalam predikat visibilitas.
func consumesStateVariable(node *ir.Node, stateName string) bool {
	if node == nil || node.Attributes == nil {
		return false
	}

	visibilityKeys := [...]string{"open", "isopen", "show", "visible", "opened"}
	for _, key := range visibilityKeys {
		if val, ok := getAttrCaseInsensitive(node, key); ok {
			if matchesStatePredicate(val, stateName) {
				return true
			}
		}
	}

	return false
}

// matchesStatePredicate mengevaluasi predikat ekspresi konsumsi state sesuai kontrak:
// - state !== null / state != null / null !== state / null != state
// - state === true / state == true / state
// - Boolean(state) / !!state
func matchesStatePredicate(expr, stateName string) bool {
	clean := cleanAttrValue(expr)
	if clean == "" {
		return false
	}
	stateLower := strings.ToLower(stateName)

	// 1. Exact variable name: open={isOpen}
	if clean == stateLower {
		return true
	}

	// 2. Boolean(state) atau !!state
	if clean == "boolean("+stateLower+")" || clean == "!!"+stateLower {
		return true
	}

	noSpace := strings.ReplaceAll(clean, " ", "")

	// 3. state !== null atau state != null
	if noSpace == stateLower+"!==null" || noSpace == stateLower+"!=null" ||
		noSpace == "null!=="+stateLower || noSpace == "null!="+stateLower {
		return true
	}

	// 4. state === true atau state == true
	if noSpace == stateLower+"===true" || noSpace == stateLower+"==true" ||
		noSpace == "true==="+stateLower || noSpace == "true=="+stateLower {
		return true
	}

	if noSpace == "boolean("+stateLower+")" || noSpace == "!!"+stateLower {
		return true
	}

	return false
}

// checkConditionalRenderGuard mengecek pola guard JSX: {state && <Dialog} atau {state !== null && <Dialog}.
func checkConditionalRenderGuard(scopeText, stateName string) bool {
	noSpace := strings.ToLower(strings.ReplaceAll(scopeText, " ", ""))
	stateLower := strings.ToLower(stateName)
	patterns := []string{
		"{" + stateLower + "&&",
		"{" + stateLower + "!==null&&",
		"{" + stateLower + "!=null&&",
		"{" + stateLower + "===true&&",
		"{boolean(" + stateLower + ")&&",
		"{!!" + stateLower + "&&",
	}
	for _, p := range patterns {
		if strings.Contains(noSpace, p) {
			return true
		}
	}
	return false
}

// isNearGuardedScope memverifikasi dialog berada di sekitar baris guarded JSX.
func isNearGuardedScope(node *ir.Node, scopeText, stateName string) bool {
	if !isConfirmationPrimitive(node) {
		return false
	}
	lines := strings.Split(scopeText, "\n")
	stateLower := strings.ToLower(stateName)

	start := node.Span.Line - 5
	if start < 0 {
		start = 0
	}
	end := node.Span.Line
	if end > len(lines) {
		end = len(lines)
	}

	for i := start; i < end; i++ {
		lineClean := strings.ToLower(strings.ReplaceAll(lines[i], " ", ""))
		if strings.Contains(lineClean, stateLower+"&&") || strings.Contains(lineClean, stateLower+"!=null&&") ||
			strings.Contains(lineClean, stateLower+"!==null&&") {
			return true
		}
	}
	return false
}

// hasDownstreamConfirmationAction memverifikasi bahwa dialog memiliki aksi konfirmasi,
// bukan semata-mata dialog informasional.
func hasDownstreamConfirmationAction(node *ir.Node) bool {
	if node == nil {
		return false
	}
	if hasDirectConfirmationAttribute(node) {
		return true
	}
	return hasConfirmationDescendant(node)
}

func hasDirectConfirmationAttribute(node *ir.Node) bool {
	actionAttrs := [...]string{
		"onconfirm", "onapprove", "ondelete", "onaccept",
		"confirmtext", "confirmlabel", "onproceed", "onsubmit",
	}
	for _, attr := range actionAttrs {
		if _, ok := getAttrCaseInsensitive(node, attr); ok {
			return true
		}
	}
	return false
}

func hasConfirmationDescendant(node *ir.Node) bool {
	for child := range node.Walk() {
		if child == node || child.Type != ir.NodeElement {
			continue
		}
		if isConfirmationActionElement(child) {
			return true
		}
	}
	return false
}

func isConfirmationActionElement(child *ir.Node) bool {
	if child.Tag == "AlertDialogAction" || strings.HasSuffix(child.Tag, "ConfirmAction") ||
		strings.HasSuffix(child.Tag, "ApprovalAction") {
		return true
	}
	if isButtonNode(child) && hasConfirmationTextOrHandler(child) {
		return true
	}
	return false
}

// hasConfirmationTextOrHandler mengecek teks atau handler tombol di dalam dialog.
func hasConfirmationTextOrHandler(buttonNode *ir.Node) bool {
	if buttonNode == nil {
		return false
	}

	// Cek teks tombol
	for _, ch := range buttonNode.Children {
		if ch.Type == ir.NodeText {
			txt := strings.ToLower(strings.TrimSpace(ch.RawClasses))
			confirmWords := [...]string{"hapus", "delete", "ya", "konfirmasi", "confirm", "ok", "setuju", "lanjut", "proceed"}
			for _, w := range confirmWords {
				if strings.Contains(txt, w) {
					return true
				}
			}
		}
	}

	// Cek handler tombol
	for k, v := range buttonNode.Attributes {
		kLower := strings.ToLower(k)
		if kLower == "onclick" || kLower == "onpress" {
			vLower := strings.ToLower(cleanAttrValue(v))
			if strings.Contains(vLower, "delete") || strings.Contains(vLower, "hapus") ||
				strings.Contains(vLower, "confirm") || strings.Contains(vLower, "approve") {
				return true
			}
		}
	}

	return false
}
