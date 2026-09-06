package performance

import (
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
