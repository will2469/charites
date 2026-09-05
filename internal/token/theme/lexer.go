package theme

import (
	"strings"
)

// Lexer memindai stream bytes CSS menjadi token leksikal secara efisien.
type Lexer struct {
	src  []byte
	idx  int
	line int
	col  int
}

// NewLexer menginisialisasi Lexer baru dari slice byte sumber CSS.
func NewLexer(src []byte) *Lexer {
	return &Lexer{
		src:  src,
		idx:  0,
		line: 1,
		col:  1,
	}
}

func (l *Lexer) makeToken(tokType TokenType, start int, startLoc SourceLocation) Token {
	return Token{
		Type:        tokType,
		Value:       string(l.src[start:l.idx]),
		StartOffset: start,
		EndOffset:   l.idx,
		Span:        SourceSpan{Start: startLoc, End: SourceLocation{Line: l.line, Column: l.col}},
	}
}

// NextToken membaca dan mengembalikan token berikutnya dari buffer.
func (l *Lexer) NextToken() Token {
	if l.idx >= len(l.src) {
		loc := SourceLocation{Line: l.line, Column: l.col}
		return Token{Type: TokenEOF, StartOffset: l.idx, EndOffset: l.idx, Span: SourceSpan{Start: loc, End: loc}}
	}

	startLoc := SourceLocation{Line: l.line, Column: l.col}
	start := l.idx
	b := l.src[l.idx]

	// 1. Whitespace
	if isWhitespace(b) {
		for l.idx < len(l.src) && isWhitespace(l.src[l.idx]) {
			l.advanceByte()
		}
		return l.makeToken(TokenWhitespace, start, startLoc)
	}

	// 2. Comments (/* ... */)
	if b == '/' && l.idx+1 < len(l.src) && l.src[l.idx+1] == '*' {
		l.advanceByte() // '/'
		l.advanceByte() // '*'
		for l.idx < len(l.src) {
			if l.src[l.idx] == '*' && l.idx+1 < len(l.src) && l.src[l.idx+1] == '/' {
				l.advanceByte() // '*'
				l.advanceByte() // '/'
				break
			}
			l.advanceByte()
		}
		return l.makeToken(TokenComment, start, startLoc)
	}

	// 3. Strings ("..." or '...')
	if b == '"' || b == '\'' {
		quote := b
		l.advanceByte() // quote awal
		for l.idx < len(l.src) {
			curr := l.src[l.idx]
			if curr == '\\' {
				l.advanceByte() // skip backslash
				if l.idx < len(l.src) {
					l.advanceByte() // skip escaped character
				}
				continue
			}
			if curr == quote {
				l.advanceByte() // quote akhir
				break
			}
			l.advanceByte()
		}
		return l.makeToken(TokenString, start, startLoc)
	}

	// 4. Single-Character Delimiters
	switch b {
	case ':':
		l.advanceByte()
		return l.makeToken(TokenColon, start, startLoc)
	case ';':
		l.advanceByte()
		return l.makeToken(TokenSemicolon, start, startLoc)
	case ',':
		l.advanceByte()
		return l.makeToken(TokenComma, start, startLoc)
	case '{':
		l.advanceByte()
		return l.makeToken(TokenOpenBrace, start, startLoc)
	case '}':
		l.advanceByte()
		return l.makeToken(TokenCloseBrace, start, startLoc)
	case '(':
		l.advanceByte()
		return l.makeToken(TokenOpenParen, start, startLoc)
	case ')':
		l.advanceByte()
		return l.makeToken(TokenCloseParen, start, startLoc)
	case '[':
		l.advanceByte()
		return l.makeToken(TokenOpenBracket, start, startLoc)
	case ']':
		l.advanceByte()
		return l.makeToken(TokenCloseBracket, start, startLoc)
	}

	// 5. At-Keywords (@media, @layer, @theme, @supports, etc.)
	if b == '@' {
		l.advanceByte() // '@'
		for l.idx < len(l.src) && isIdentChar(l.src[l.idx]) {
			l.advanceByte()
		}
		return l.makeToken(TokenAtKeyword, start, startLoc)
	}

	// 6. Hash / Hex (#123456)
	if b == '#' {
		l.advanceByte() // '#'
		for l.idx < len(l.src) && isHexOrIdentChar(l.src[l.idx]) {
			l.advanceByte()
		}
		return l.makeToken(TokenHash, start, startLoc)
	}

	// 7. Special url(...) function handling
	if (b == 'u' || b == 'U') && l.hasPrefixCI("url(") {
		l.advanceN(4) // "url("
		parenDepth := 1
		for l.idx < len(l.src) && parenDepth > 0 {
			curr := l.src[l.idx]
			if curr == '\\' {
				l.advanceByte()
				if l.idx < len(l.src) {
					l.advanceByte()
				}
				continue
			}
			if curr == '"' || curr == '\'' {
				quote := curr
				l.advanceByte()
				for l.idx < len(l.src) {
					if l.src[l.idx] == '\\' {
						l.advanceByte()
						if l.idx < len(l.src) {
							l.advanceByte()
						}
						continue
					}
					if l.src[l.idx] == quote {
						l.advanceByte()
						break
					}
					l.advanceByte()
				}
				continue
			}
			if curr == '(' {
				parenDepth++
			} else if curr == ')' {
				parenDepth--
				if parenDepth == 0 {
					l.advanceByte()
					break
				}
			}
			l.advanceByte()
		}
		return l.makeToken(TokenIdent, start, startLoc)
	}

	// 8. Identifiers (e.g. --banana, --color-primary, display, oklch, \:focus, --tw-bg-opacity\:1)
	if l.startsIdent() {
		l.consumeIdent()
		return l.makeToken(TokenIdent, start, startLoc)
	}

	// 9. Numbers, Dimensions & Percentages (including signed e.g. -10px, +5%)
	isSignPrefix := (b == '-' || b == '+') && l.idx+1 < len(l.src) && (isDigit(l.src[l.idx+1]) || (l.src[l.idx+1] == '.' && l.idx+2 < len(l.src) && isDigit(l.src[l.idx+2])))
	if isDigit(b) || (b == '.' && l.idx+1 < len(l.src) && isDigit(l.src[l.idx+1])) || isSignPrefix {
		if isSignPrefix {
			l.advanceByte() // consume '+' or '-'
		}
		for l.idx < len(l.src) && (isDigit(l.src[l.idx]) || l.src[l.idx] == '.') {
			l.advanceByte()
		}
		if l.idx < len(l.src) && l.src[l.idx] == '%' {
			l.advanceByte()
			return l.makeToken(TokenPercentage, start, startLoc)
		}
		if l.startsIdent() {
			l.consumeIdent()
			return l.makeToken(TokenDimension, start, startLoc)
		}
		return l.makeToken(TokenNumber, start, startLoc)
	}

	// 10. Fallback Delim
	l.advanceByte()
	return l.makeToken(TokenDelim, start, startLoc)
}

func (l *Lexer) advanceByte() {
	if l.idx < len(l.src) {
		if l.src[l.idx] == '\n' {
			l.line++
			l.col = 1
		} else {
			l.col++
		}
		l.idx++
	}
}

func (l *Lexer) advanceN(n int) {
	for i := 0; i < n && l.idx < len(l.src); i++ {
		l.advanceByte()
	}
}

func (l *Lexer) hasPrefixCI(prefix string) bool {
	if l.idx+len(prefix) > len(l.src) {
		return false
	}
	sub := string(l.src[l.idx : l.idx+len(prefix)])
	return strings.EqualFold(sub, prefix)
}

func (l *Lexer) startsIdent() bool {
	if l.idx >= len(l.src) {
		return false
	}
	b0 := l.src[l.idx]
	var b1, b2 byte
	if l.idx+1 < len(l.src) {
		b1 = l.src[l.idx+1]
	}
	if l.idx+2 < len(l.src) {
		b2 = l.src[l.idx+2]
	}

	if b0 == '-' {
		return isIdentStart(b1) || b1 == '-' || isValidEscape(b1, b2)
	}
	if isIdentStart(b0) {
		return true
	}
	if b0 == '\\' {
		return isValidEscape(b0, b1)
	}
	return false
}

func (l *Lexer) consumeIdent() {
	for l.idx < len(l.src) {
		b := l.src[l.idx]
		if isIdentChar(b) {
			l.advanceByte()
			continue
		}
		if b == '\\' && l.idx+1 < len(l.src) && isValidEscape(b, l.src[l.idx+1]) {
			l.advanceByte() // consume '\\'
			// Jika diikuti hex digit, konsumsi hingga 6 digit
			if l.idx < len(l.src) && isHexDigit(l.src[l.idx]) {
				hexCount := 0
				for l.idx < len(l.src) && isHexDigit(l.src[l.idx]) && hexCount < 6 {
					l.advanceByte()
					hexCount++
				}
				// Jika diikuti 1 karakter spasi/whitespace, konsumsi sebagai delimiter escape (CSS spec 4.3.9)
				if l.idx < len(l.src) && isWhitespace(l.src[l.idx]) {
					l.advanceByte()
				}
			} else if l.idx < len(l.src) {
				l.advanceByte() // konsumsi byte yang di-escape
			}
			continue
		}
		break
	}
}

func isWhitespace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\r' || b == '\n' || b == '\f'
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func isHexDigit(b byte) bool {
	return (b >= '0' && b <= '9') || (b >= 'a' && b <= 'f') || (b >= 'A' && b <= 'F')
}

// isValidEscape memeriksa apakah dua byte membentuk CSS escape yang valid (CSS Syntax 3 section 4.3.7).
// Backslash diikuti karakter selain newline adalah valid escape sequence.
func isValidEscape(b1, b2 byte) bool {
	return b1 == '\\' && b2 != '\n' && b2 != '\r' && b2 != '\f'
}

func isIdentStart(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_' || b >= 128
}

func isIdentChar(b byte) bool {
	return isIdentStart(b) || isDigit(b) || b == '-'
}

func isHexOrIdentChar(b byte) bool {
	return isIdentChar(b)
}

// UnescapeCSS mengurai escape sequences pada identifier atau string CSS sesuai CSS Syntax Level 3.
// Contoh: `\:` -> `:`, `\31 ` -> `1`, `\20` -> ` `, `\--` -> `--`.
func UnescapeCSS(raw string) string {
	if !strings.ContainsRune(raw, '\\') {
		return raw
	}
	var sb strings.Builder
	sb.Grow(len(raw))
	n := len(raw)
	for i := 0; i < n; i++ {
		if raw[i] != '\\' {
			sb.WriteByte(raw[i])
			continue
		}
		// Escape dimulai
		if i+1 >= n {
			sb.WriteRune('\uFFFD')
			break
		}
		i++ // lewati '\\'
		b := raw[i]
		if isHexDigit(b) {
			hexStart := i
			for i < n && isHexDigit(raw[i]) && (i-hexStart) < 6 {
				i++
			}
			var cp rune
			for j := hexStart; j < i; j++ {
				cp = cp*16 + hexVal(raw[j])
			}
			if cp == 0 || (cp >= 0xD800 && cp <= 0xDFFF) || cp > 0x10FFFF {
				cp = '\uFFFD'
			}
			sb.WriteRune(cp)
			// Jika diikuti satu karakter whitespace, konsumsi (CSS Syntax 3)
			if i < n && isWhitespace(raw[i]) {
				// Karakter spasi dikonsumsi sebagai pemisah hex escape
			} else {
				i-- // kompensasi loop
			}
			continue
		}
		// Bukan hex escape, tulis karakter literal apa adanya
		sb.WriteByte(b)
	}
	return sb.String()
}

func hexVal(b byte) rune {
	if b >= '0' && b <= '9' {
		return rune(b - '0')
	}
	if b >= 'a' && b <= 'f' {
		return rune(b - 'a' + 10)
	}
	if b >= 'A' && b <= 'F' {
		return rune(b - 'A' + 10)
	}
	return 0
}
