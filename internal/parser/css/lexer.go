package css

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

	switch {
	case isWhitespace(b):
		return l.consumeWhitespace(start, startLoc)
	case b == '/' && l.idx+1 < len(l.src) && l.src[l.idx+1] == '*':
		return l.consumeComment(start, startLoc)
	case b == '"' || b == '\'':
		return l.consumeString(b, start, startLoc)
	case b == '@':
		return l.consumeAtKeyword(start, startLoc)
	case b == '#':
		return l.consumeHash(start, startLoc)
	case (b == 'u' || b == 'U') && l.hasPrefixCI("url("):
		return l.consumeURL(start, startLoc)
	case l.startsIdent():
		l.consumeIdent()
		return l.makeToken(TokenIdent, start, startLoc)
	case isNumberStart(b, l.src, l.idx):
		return l.consumeNumber(start, startLoc)
	}

	return l.consumeDelimiter(b, start, startLoc)
}

func (l *Lexer) consumeWhitespace(start int, startLoc SourceLocation) Token {
	for l.idx < len(l.src) && isWhitespace(l.src[l.idx]) {
		l.advanceByte()
	}
	return l.makeToken(TokenWhitespace, start, startLoc)
}

func (l *Lexer) consumeComment(start int, startLoc SourceLocation) Token {
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

func (l *Lexer) consumeString(quote byte, start int, startLoc SourceLocation) Token {
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

func (l *Lexer) consumeAtKeyword(start int, startLoc SourceLocation) Token {
	l.advanceByte() // '@'
	for l.idx < len(l.src) && isIdentChar(l.src[l.idx]) {
		l.advanceByte()
	}
	return l.makeToken(TokenAtKeyword, start, startLoc)
}

func (l *Lexer) consumeHash(start int, startLoc SourceLocation) Token {
	l.advanceByte() // '#'
	for l.idx < len(l.src) && isHexOrIdentChar(l.src[l.idx]) {
		l.advanceByte()
	}
	return l.makeToken(TokenHash, start, startLoc)
}

func (l *Lexer) consumeURL(start int, startLoc SourceLocation) Token {
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
			l.skipQuotedLiteral(curr)
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

func (l *Lexer) skipQuotedLiteral(quote byte) {
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
}

func isNumberStart(b byte, src []byte, idx int) bool {
	if isDigit(b) || (b == '.' && idx+1 < len(src) && isDigit(src[idx+1])) {
		return true
	}
	isSign := b == '-' || b == '+'
	if isSign && idx+1 < len(src) {
		next := src[idx+1]
		return isDigit(next) || (next == '.' && idx+2 < len(src) && isDigit(src[idx+2]))
	}
	return false
}

func (l *Lexer) consumeNumber(start int, startLoc SourceLocation) Token {
	b := l.src[l.idx]
	if b == '-' || b == '+' {
		l.advanceByte()
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

func (l *Lexer) consumeDelimiter(b byte, start int, startLoc SourceLocation) Token {
	l.advanceByte()
	switch b {
	case ':':
		return l.makeToken(TokenColon, start, startLoc)
	case ';':
		return l.makeToken(TokenSemicolon, start, startLoc)
	case ',':
		return l.makeToken(TokenComma, start, startLoc)
	case '{':
		return l.makeToken(TokenOpenBrace, start, startLoc)
	case '}':
		return l.makeToken(TokenCloseBrace, start, startLoc)
	case '(':
		return l.makeToken(TokenOpenParen, start, startLoc)
	case ')':
		return l.makeToken(TokenCloseParen, start, startLoc)
	case '[':
		return l.makeToken(TokenOpenBracket, start, startLoc)
	case ']':
		return l.makeToken(TokenCloseBracket, start, startLoc)
	default:
		return l.makeToken(TokenDelim, start, startLoc)
	}
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
			l.consumeEscapeSequence()
			continue
		}
		break
	}
}

func (l *Lexer) consumeEscapeSequence() {
	l.advanceByte() // consume '\\'
	if l.idx < len(l.src) && isHexDigit(l.src[l.idx]) {
		hexCount := 0
		for l.idx < len(l.src) && isHexDigit(l.src[l.idx]) && hexCount < 6 {
			l.advanceByte()
			hexCount++
		}
		if l.idx < len(l.src) && isWhitespace(l.src[l.idx]) {
			l.advanceByte()
		}
		return
	}
	if l.idx < len(l.src) {
		l.advanceByte()
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
		if i+1 >= n {
			sb.WriteRune('\uFFFD')
			break
		}
		i++ // lewati '\\'
		b := raw[i]
		if isHexDigit(b) {
			r, consumed := parseHexCodePoint(raw, i)
			sb.WriteRune(r)
			i = consumed
			continue
		}
		sb.WriteByte(b)
	}
	return sb.String()
}

func parseHexCodePoint(raw string, start int) (rune, int) {
	n := len(raw)
	i := start
	for i < n && isHexDigit(raw[i]) && (i-start) < 6 {
		i++
	}
	var cp rune
	for j := start; j < i; j++ {
		cp = cp*16 + hexVal(raw[j])
	}
	if cp == 0 || (cp >= 0xD800 && cp <= 0xDFFF) || cp > 0x10FFFF {
		cp = '\uFFFD'
	}
	// Jika diikuti 1 whitespace delimiter escape, lewati
	if i < n && isWhitespace(raw[i]) {
		return cp, i
	}
	return cp, i - 1
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
