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

	// 8. Identifiers (e.g. --banana, --color-primary, display, oklch)
	if isIdentStart(b) || (b == '-' && l.idx+1 < len(l.src) && (isIdentStart(l.src[l.idx+1]) || l.src[l.idx+1] == '-')) {
		for l.idx < len(l.src) && isIdentChar(l.src[l.idx]) {
			l.advanceByte()
		}
		return l.makeToken(TokenIdent, start, startLoc)
	}

	// 9. Numbers, Dimensions & Percentages
	if isDigit(b) || (b == '.' && l.idx+1 < len(l.src) && isDigit(l.src[l.idx+1])) {
		for l.idx < len(l.src) && (isDigit(l.src[l.idx]) || l.src[l.idx] == '.') {
			l.advanceByte()
		}
		if l.idx < len(l.src) && l.src[l.idx] == '%' {
			l.advanceByte()
			return l.makeToken(TokenPercentage, start, startLoc)
		}
		if l.idx < len(l.src) && isIdentStart(l.src[l.idx]) {
			for l.idx < len(l.src) && isIdentChar(l.src[l.idx]) {
				l.advanceByte()
			}
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

func isWhitespace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\r' || b == '\n' || b == '\f'
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
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
