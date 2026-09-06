package astro

import (
	"bytes"
	"strings"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/parser"
)

var voidElements = map[string]bool{
	"area":   true,
	"base":   true,
	"br":     true,
	"col":    true,
	"embed":  true,
	"hr":     true,
	"img":    true,
	"input":  true,
	"link":   true,
	"meta":   true,
	"param":  true,
	"source": true,
	"track":  true,
	"wbr":    true,
}

// Parse memindai berkas komponen Astro (.astro), mengisolasi frontmatter,
// mempertahankan nomor baris asli 1-indexed, dan merakit pohon IR terpadu *ir.Node.
// Mematuhi Zero-Panic Invariant dan recovery semantics pada input malformed.
func Parse(src []byte) (*ir.Node, error) {
	l := newLexer(src)
	return l.parse(), nil
}

type lexer struct {
	src  []byte
	len  int
	pos  int
	line int
	col  int
	bld  *ir.Builder
}

func newLexer(src []byte) *lexer {
	return &lexer{
		src:  src,
		len:  len(src),
		pos:  0,
		line: 1,
		col:  1,
		bld:  ir.NewBuilder(),
	}
}

func (l *lexer) parse() *ir.Node {
	l.skipFrontmatter()

	for l.pos < l.len {
		ch := l.src[l.pos]

		if ch == '<' {
			l.handleTag()
			continue
		}

		if ch == '{' {
			// Ekspresi JSX atau komentar {/* ... */}
			l.handleBraceExpression()
			continue
		}

		// Karakter teks di luar tag
		startLine := l.line
		startCol := l.col
		textStart := l.pos
		for l.pos < l.len && l.src[l.pos] != '<' && l.src[l.pos] != '{' {
			l.advance()
		}
		rawText := string(l.src[textStart:l.pos])
		trimmed := strings.TrimSpace(rawText)
		if trimmed != "" {
			l.bld.AddText(rawText, ir.Span{
				Line:      startLine,
				Column:    startCol,
				EndLine:   l.line,
				EndColumn: l.col,
			})
		}
	}

	return l.bld.Root()
}

// skipFrontmatter melewati blok frontmatter (--- ... ---) jika ada di awal berkas,
// sambil terus memperbarui nomor baris asli (line offset preservation).
func (l *lexer) skipFrontmatter() {
	l.skipWhitespace()
	if l.pos+3 > l.len || l.src[l.pos] != '-' || l.src[l.pos+1] != '-' || l.src[l.pos+2] != '-' {
		return
	}

	l.advanceBytes(3)

	for l.pos < l.len {
		if l.src[l.pos] == '\n' {
			l.advance()
			if l.checkAndConsumeClosingFrontmatter() {
				return
			}
			continue
		}
		l.advance()
	}
}

func (l *lexer) checkAndConsumeClosingFrontmatter() bool {
	if l.pos+3 <= l.len && l.src[l.pos] == '-' && l.src[l.pos+1] == '-' && l.src[l.pos+2] == '-' {
		l.advanceBytes(3)
		for l.pos < l.len && l.src[l.pos] != '\n' {
			l.advance()
		}
		if l.pos < l.len && l.src[l.pos] == '\n' {
			l.advance()
		}
		return true
	}
	return false
}

func (l *lexer) skipUntilChar(target byte) {
	for l.pos < l.len && l.src[l.pos] != target {
		l.advance()
	}
	if l.pos < l.len && l.src[l.pos] == target {
		l.advance()
	}
}

func (l *lexer) handleHTMLComment(startLine, startCol int) {
	l.pos += 3
	l.col += 3
	commentStart := l.pos
	endIdx := bytes.Index(l.src[l.pos:], []byte("-->"))
	if endIdx == -1 {
		raw := string(l.src[commentStart:])
		l.advanceTo(l.len)
		l.bld.AddComment(raw, ir.Span{Line: startLine, Column: startCol, EndLine: l.line, EndColumn: l.col})
		return
	}
	raw := string(l.src[commentStart : l.pos+endIdx])
	l.advanceBytes(endIdx + 3)
	l.bld.AddComment(raw, ir.Span{Line: startLine, Column: startCol, EndLine: l.line, EndColumn: l.col})
}

func (l *lexer) handleCloseElement() {
	l.advance() // konsumsi '/'
	l.skipWhitespace()
	tagName := l.readIdentifier()
	l.skipUntilChar('>')
	l.bld.CloseElement(tagName)
}

func (l *lexer) parseAttributeVal() string {
	l.skipWhitespace()
	if l.pos < l.len && l.src[l.pos] == '=' {
		l.advance() // konsumsi '='
		l.skipWhitespace()
		return l.readAttributeValue()
	}
	return ""
}

type parsedAttrs struct {
	attrs       map[string]string
	rawClasses  string
	classes     []string
	hasDynamic  bool
	selfClosing bool
	ok          bool
}

func (l *lexer) parseElementAttributes() parsedAttrs {
	res := parsedAttrs{
		attrs: make(map[string]string),
	}

	for l.pos < l.len {
		l.skipWhitespace()
		if l.pos >= l.len || l.src[l.pos] == '<' {
			return res
		}

		ch := l.src[l.pos]
		if ch == '>' {
			l.advance()
			res.ok = true
			return res
		}
		if ch == '/' && l.pos+1 < l.len && l.src[l.pos+1] == '>' {
			l.advanceBytes(2)
			res.selfClosing = true
			res.ok = true
			return res
		}

		if ch == '{' {
			start := l.pos
			l.skipBraceExpression()
			expr := strings.TrimSpace(string(l.src[start:l.pos]))
			if strings.HasPrefix(expr, "{...") {
				res.attrs[expr] = expr
			}
			continue
		}

		attrName := l.readAttributeName()
		if attrName == "" {
			l.advance()
			continue
		}

		attrVal := l.parseAttributeVal()
		res.attrs[attrName] = attrVal
		if attrName == "class" || attrName == "className" {
			res.rawClasses = attrVal
			res.classes, res.hasDynamic = parser.ExtractClasses(attrVal)
		}
	}

	return res
}

func (l *lexer) handleOpenElement(startLine, startCol int) {
	tagName := l.readIdentifier()
	if tagName == "" {
		return
	}

	p := l.parseElementAttributes()
	if !p.ok {
		return
	}

	span := ir.Span{
		Line:      startLine,
		Column:    startCol,
		EndLine:   l.line,
		EndColumn: l.col,
	}

	isVoid := voidElements[strings.ToLower(tagName)]
	switch {
	case p.selfClosing || isVoid:
		l.bld.AddSelfClosingElement(tagName, span, p.rawClasses, p.classes, p.attrs, p.hasDynamic)
	case strings.EqualFold(tagName, "style") || strings.EqualFold(tagName, "script"):
		l.bld.OpenElement(tagName, span, p.rawClasses, p.classes, p.attrs, p.hasDynamic)
		l.consumeRawTextElement(tagName)
	default:
		l.bld.OpenElement(tagName, span, p.rawClasses, p.classes, p.attrs, p.hasDynamic)
	}
}

func (l *lexer) consumeRawTextElement(tagName string) {
	closePrefix := "</" + strings.ToLower(tagName)
	startLine := l.line
	startCol := l.col
	textStart := l.pos

	for l.pos < l.len {
		if l.src[l.pos] == '<' && l.pos+len(closePrefix) <= l.len {
			cand := strings.ToLower(string(l.src[l.pos : l.pos+len(closePrefix)]))
			if cand == closePrefix {
				rawText := string(l.src[textStart:l.pos])
				if strings.TrimSpace(rawText) != "" {
					l.bld.AddText(rawText, ir.Span{
						Line:      startLine,
						Column:    startCol,
						EndLine:   l.line,
						EndColumn: l.col,
					})
				}
				l.advanceBytes(len(closePrefix))
				l.skipUntilChar('>')
				l.bld.CloseElement(tagName)
				return
			}
		}
		l.advance()
	}

	rawText := string(l.src[textStart:l.pos])
	if strings.TrimSpace(rawText) != "" {
		l.bld.AddText(rawText, ir.Span{
			Line:      startLine,
			Column:    startCol,
			EndLine:   l.line,
			EndColumn: l.col,
		})
	}
	l.bld.CloseElement(tagName)
}

// handleTag memproses tag pembuka, penutup, fragment, comment, atau doctype.
func (l *lexer) handleTag() {
	startLine := l.line
	startCol := l.col
	l.advance() // konsumsi '<'

	if l.pos >= l.len {
		return
	}

	if l.matchPrefix("!--") {
		l.handleHTMLComment(startLine, startCol)
		return
	}

	if l.matchPrefix("/>") {
		l.advanceBytes(2)
		l.bld.CloseFragment()
		return
	}

	if l.src[l.pos] == '/' {
		l.handleCloseElement()
		return
	}

	if l.src[l.pos] == '>' {
		l.advance()
		l.bld.OpenFragment(ir.Span{Line: startLine, Column: startCol, EndLine: l.line, EndColumn: l.col})
		return
	}

	if l.src[l.pos] == '!' {
		l.skipUntilChar('>')
		return
	}

	l.handleOpenElement(startLine, startCol)
}

// handleBraceExpression memproses kurung kurawal JSX { ... } atau {/* ... */} di luar tag.
func (l *lexer) handleBraceExpression() {
	startLine := l.line
	startCol := l.col
	l.advance() // konsumsi '{'

	// Komentar JSX: {/* ... */}
	if l.pos+2 <= l.len && l.src[l.pos] == '/' && l.src[l.pos+1] == '*' {
		l.advance()
		l.advance()
		commentStart := l.pos
		endIdx := bytes.Index(l.src[l.pos:], []byte("*/}"))
		if endIdx == -1 {
			raw := string(l.src[commentStart:])
			l.advanceTo(l.len)
			l.bld.AddComment(raw, ir.Span{Line: startLine, Column: startCol, EndLine: l.line, EndColumn: l.col})
			return
		}
		raw := string(l.src[commentStart : l.pos+endIdx])
		l.advanceBytes(endIdx + 3)
		l.bld.AddComment(raw, ir.Span{Line: startLine, Column: startCol, EndLine: l.line, EndColumn: l.col})
		return
	}

	// Ekspresi JS dalam kurung kurawal
	depth := 1
	for l.pos < l.len && depth > 0 {
		ch := l.src[l.pos]
		switch ch {
		case '{':
			depth++
			l.advance()
		case '}':
			depth--
			l.advance()
		case '"', '\'', '`':
			l.skipStringLiteral(ch)
		default:
			l.advance()
		}
	}
}

func (l *lexer) skipBraceExpression() {
	depth := 0
	for l.pos < l.len {
		ch := l.src[l.pos]
		if ch == '{' {
			depth++
			l.advance()
			continue
		}
		if ch == '}' {
			depth--
			l.advance()
			if depth <= 0 {
				return
			}
			continue
		}
		if ch == '"' || ch == '\'' || ch == '`' {
			l.skipStringLiteral(ch)
			continue
		}
		l.advance()
	}
}

func (l *lexer) readIdentifier() string {
	start := l.pos
	for l.pos < l.len {
		ch := l.src[l.pos]
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') ||
			ch == '_' || ch == '-' || ch == '.' || ch == ':' {
			l.advance()
		} else {
			break
		}
	}
	return string(l.src[start:l.pos])
}

func (l *lexer) readAttributeName() string {
	start := l.pos
	for l.pos < l.len {
		ch := l.src[l.pos]
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') ||
			ch == '_' || ch == '-' || ch == '.' || ch == ':' || ch == '@' || ch == '#' || ch == '$' {
			l.advance()
		} else {
			break
		}
	}
	return string(l.src[start:l.pos])
}

func (l *lexer) readQuotedAttribute(quote byte) string {
	start := l.pos
	l.advance()
	for l.pos < l.len && l.src[l.pos] != quote {
		if l.src[l.pos] == '\\' && l.pos+1 < l.len {
			l.advance()
		}
		l.advance()
	}
	if l.pos < l.len && l.src[l.pos] == quote {
		l.advance()
	}
	return string(l.src[start:l.pos])
}

func (l *lexer) readJSXBraceAttribute() string {
	start := l.pos
	l.advance()
	depth := 1
	for l.pos < l.len && depth > 0 {
		c := l.src[l.pos]
		switch c {
		case '{':
			depth++
			l.advance()
		case '}':
			depth--
			l.advance()
		case '"', '\'', '`':
			l.skipStringLiteral(c)
		default:
			l.advance()
		}
	}
	return string(l.src[start:l.pos])
}

func (l *lexer) readUnquotedAttribute() string {
	start := l.pos
	for l.pos < l.len {
		c := l.src[l.pos]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '>' || c == '/' {
			break
		}
		l.advance()
	}
	return string(l.src[start:l.pos])
}

func (l *lexer) readAttributeValue() string {
	if l.pos >= l.len {
		return ""
	}

	ch := l.src[l.pos]
	switch ch {
	case '"', '\'':
		return l.readQuotedAttribute(ch)
	case '{':
		return l.readJSXBraceAttribute()
	default:
		return l.readUnquotedAttribute()
	}
}

func (l *lexer) skipStringLiteral(quote byte) {
	l.advance() // konsumsi quote pembuka
	for l.pos < l.len && l.src[l.pos] != quote {
		if l.src[l.pos] == '\\' && l.pos+1 < l.len {
			l.advance()
		}
		l.advance()
	}
	if l.pos < l.len && l.src[l.pos] == quote {
		l.advance()
	}
}

func (l *lexer) skipWhitespace() {
	for l.pos < l.len {
		ch := l.src[l.pos]
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			l.advance()
		} else {
			break
		}
	}
}

func (l *lexer) matchPrefix(prefix string) bool {
	if l.pos+len(prefix) > l.len {
		return false
	}
	return string(l.src[l.pos:l.pos+len(prefix)]) == prefix
}

func (l *lexer) advance() {
	if l.pos < l.len {
		if l.src[l.pos] == '\n' {
			l.line++
			l.col = 1
		} else {
			l.col++
		}
		l.pos++
	}
}

func (l *lexer) advanceBytes(n int) {
	for i := 0; i < n && l.pos < l.len; i++ {
		l.advance()
	}
}

func (l *lexer) advanceTo(targetPos int) {
	for l.pos < targetPos && l.pos < l.len {
		l.advance()
	}
}
