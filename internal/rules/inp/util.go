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
	case "onClick", "onKeyDown", "onKeyUp", "onKeyPress", "onPointerDown", "onPointerUp", "onSubmit":
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
