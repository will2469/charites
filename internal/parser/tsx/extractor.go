package tsx

import (
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

// Extract memindai kode TypeScript/React JSX (.tsx, .jsx), mengekstrak struktur hierarki JSX
// menggunakan Option B Structural Scanning, mengisolasi template literal dinamis,
// dan merakit pohon IR terpadu *ir.Node.
// Mematuhi Zero-Panic Invariant dan recovery semantics pada input malformed.
func Extract(src []byte) (*ir.Node, error) {
	e := newExtractor(src)
	return e.extract(), nil
}

type extractor struct {
	src      []byte
	len      int
	pos      int
	line     int
	col      int
	lastByte byte
	bld      *ir.Builder
}

func newExtractor(src []byte) *extractor {
	return &extractor{
		src:      src,
		len:      len(src),
		pos:      0,
		line:     1,
		col:      1,
		lastByte: ' ',
		bld:      ir.NewBuilder(),
	}
}

func (e *extractor) extract() *ir.Node {
	for e.pos < e.len {
		ch := e.src[e.pos]

		// 1. Komentar baris tunggal: // ...
		if ch == '/' && e.pos+1 < e.len && e.src[e.pos+1] == '/' {
			e.skipLineComment()
			continue
		}

		// 2. Komentar blok: /* ... */ atau {/* ... */}
		if ch == '/' && e.pos+1 < e.len && e.src[e.pos+1] == '*' {
			e.skipBlockComment()
			continue
		}

		// 3. String literal di luar tag JSX
		if ch == '"' || ch == '\'' {
			e.skipStringLiteral(ch)
			continue
		}

		// 4. Template literal JS di luar tag JSX: `...`
		if ch == '`' {
			e.skipJSTemplateLiteral()
			continue
		}

		// 5. Potensi tag JSX: <
		if ch == '<' {
			if e.tryHandleJSX() {
				continue
			}
			// Bukan tag JSX valid (misal operator perbandingan count < 10 atau generics)
			e.advance()
			continue
		}

		e.advance()
	}

	return e.bld.Root()
}

// tryHandleJSX mencoba memproses token JSX. Mengembalikan true jika berhasil mengenali JSX tag/fragment.
func e_isIdentifierChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}

func isIllegalJSXHeaderChar(src []byte, pos, length int) bool {
	c := src[pos]
	if c == ';' || c == ')' || c == '(' || c == ',' {
		return true
	}
	if (c == '&' && pos+1 < length && src[pos+1] == '&') ||
		(c == '|' && pos+1 < length && src[pos+1] == '|') {
		return true
	}
	return false
}

func (e *extractor) isArrowGeneric(tagName string) bool {
	if len(tagName) == 1 && tagName[0] >= 'A' && tagName[0] <= 'Z' {
		peek := e.pos + 1
		for peek < e.len && (e.src[peek] == ' ' || e.src[peek] == '\t' || e.src[peek] == '\n' || e.src[peek] == '\r') {
			peek++
		}
		if peek < e.len && e.src[peek] == '(' {
			return true
		}
	}
	return false
}

func (e *extractor) tryHandleClosing() (handled bool, ok bool) {
	if e.pos >= e.len || e.src[e.pos] != '/' {
		return false, false
	}

	// Fragment closing: </>
	if e.pos+1 < e.len && e.src[e.pos+1] == '>' {
		e.advanceBytes(2)
		e.bld.CloseFragment()
		return true, true
	}

	// Tag penutup: </tag>
	e.advance() // konsumsi '/'
	e.skipWhitespace()
	tagName := e.readIdentifier()
	if tagName == "" {
		return true, false
	}
	for e.pos < e.len && e.src[e.pos] != '>' {
		e.advance()
	}
	if e.pos < e.len && e.src[e.pos] == '>' {
		e.advance()
	}
	e.bld.CloseElement(tagName)
	return true, true
}

type jsxParsedAttrs struct {
	attrs       map[string]string
	rawClasses  string
	classes     []string
	hasDynamic  bool
	selfClosing bool
	isValid     bool
	isBroken    bool
}

func (e *extractor) checkTagDelimiter(tagName string, savedPos, startLine, startCol int, res *jsxParsedAttrs) bool {
	c := e.src[e.pos]
	if c == '<' {
		res.isBroken = true
		res.isValid = true
		return true
	}

	if isIllegalJSXHeaderChar(e.src, e.pos, e.len) {
		e.rollback(savedPos, startLine, startCol)
		return true
	}

	if c == '>' {
		if e.isArrowGeneric(tagName) {
			e.rollback(savedPos, startLine, startCol)
			return true
		}
		e.advance()
		res.isValid = true
		return true
	}

	if c == '/' && e.pos+1 < e.len && e.src[e.pos+1] == '>' {
		e.advanceBytes(2)
		res.selfClosing = true
		res.isValid = true
		return true
	}

	return false
}

func (e *extractor) parseJSXAttributes(tagName string, savedPos, startLine, startCol int) jsxParsedAttrs {
	res := jsxParsedAttrs{
		attrs: make(map[string]string),
	}

	for e.pos < e.len {
		e.skipWhitespace()
		if e.pos >= e.len {
			break
		}

		if e.checkTagDelimiter(tagName, savedPos, startLine, startCol, &res) {
			return res
		}

		c := e.src[e.pos]
		if c == '{' {
			e.skipBraceExpression()
			continue
		}

		attrName := e.readAttributeName()
		if attrName == "" {
			e.advance()
			continue
		}

		attrVal := e.parseJSXAttrVal()
		res.attrs[attrName] = attrVal
		if attrName == "className" || attrName == "class" {
			res.rawClasses = attrVal
			res.classes, res.hasDynamic = parser.ExtractClasses(attrVal)
		}
	}

	return res
}

func (e *extractor) parseJSXAttrVal() string {
	e.skipWhitespace()
	if e.pos < e.len && e.src[e.pos] == '=' {
		e.advance() // konsumsi '='
		e.skipWhitespace()
		return e.readAttributeValue()
	}
	return ""
}

func (e *extractor) tryHandleJSX() bool {
	startLine := e.line
	startCol := e.col
	savedPos := e.pos

	e.advance() // konsumsi '<'
	if e.pos >= e.len {
		return false
	}

	if handled, ok := e.tryHandleClosing(); handled {
		return ok
	}

	next := e.src[e.pos]
	if next == '>' {
		e.advance()
		e.bld.OpenFragment(ir.Span{Line: startLine, Column: startCol, EndLine: e.line, EndColumn: e.col})
		return true
	}

	if savedPos > 0 && e_isIdentifierChar(e.src[savedPos-1]) {
		e.rollback(savedPos, startLine, startCol)
		return false
	}

	if !((next >= 'a' && next <= 'z') || (next >= 'A' && next <= 'Z') || next == '_') {
		e.rollback(savedPos, startLine, startCol)
		return false
	}

	tagName := e.readIdentifier()
	if tagName == "" {
		e.rollback(savedPos, startLine, startCol)
		return false
	}

	p := e.parseJSXAttributes(tagName, savedPos, startLine, startCol)
	if !p.isValid {
		return false
	}
	if p.isBroken {
		return true
	}

	span := ir.Span{
		Line:      startLine,
		Column:    startCol,
		EndLine:   e.line,
		EndColumn: e.col,
	}

	isVoid := voidElements[strings.ToLower(tagName)]
	if p.selfClosing || isVoid {
		e.bld.AddSelfClosingElement(tagName, span, p.rawClasses, p.classes, p.attrs, p.hasDynamic)
	} else {
		e.bld.OpenElement(tagName, span, p.rawClasses, p.classes, p.attrs, p.hasDynamic)
	}

	return true
}

func (e *extractor) skipLineComment() {
	e.advance() // '/'
	e.advance() // '/'
	for e.pos < e.len && e.src[e.pos] != '\n' {
		e.advance()
	}
}

func (e *extractor) skipBlockComment() {
	startLine := e.line
	startCol := e.col
	commentStart := e.pos
	e.advance() // '/'
	e.advance() // '*'

	for e.pos+1 < e.len {
		if e.src[e.pos] == '*' && e.src[e.pos+1] == '/' {
			e.advance()
			e.advance()
			raw := string(e.src[commentStart:e.pos])
			if e.bld.StackDepth() > 0 {
				e.bld.AddComment(raw, ir.Span{Line: startLine, Column: startCol, EndLine: e.line, EndColumn: e.col})
			}
			return
		}
		e.advance()
	}
	// Unclosed block comment
	for e.pos < e.len {
		e.advance()
	}
}

func (e *extractor) skipStringLiteral(quote byte) {
	e.advance()
	for e.pos < e.len && e.src[e.pos] != quote {
		if e.src[e.pos] == '\\' && e.pos+1 < e.len {
			e.advance()
		}
		e.advance()
	}
	if e.pos < e.len && e.src[e.pos] == quote {
		e.advance()
	}
}

func (e *extractor) skipJSTemplateLiteral() {
	e.advance() // konsumsi '`'
	for e.pos < e.len && e.src[e.pos] != '`' {
		if e.src[e.pos] == '\\' && e.pos+1 < e.len {
			e.advance()
			e.advance()
			continue
		}
		if e.src[e.pos] == '$' && e.pos+1 < e.len && e.src[e.pos+1] == '{' {
			e.advance()
			e.advance()
			e.skipBraceExpression()
			continue
		}
		e.advance()
	}
	if e.pos < e.len && e.src[e.pos] == '`' {
		e.advance()
	}
}

func (e *extractor) skipBraceChar(depth *int) {
	c := e.src[e.pos]
	switch c {
	case '{':
		*depth++
		e.advance()
	case '}':
		*depth--
		e.advance()
	case '"', '\'':
		e.skipStringLiteral(c)
	case '`':
		e.skipJSTemplateLiteral()
	default:
		if c == '<' && e.trySkipJSXInBrace() {
			return
		}
		e.advance()
	}
}

func (e *extractor) trySkipJSXInBrace() bool {
	if e.pos+1 < e.len {
		next := e.src[e.pos+1]
		if (next >= 'a' && next <= 'z') || (next >= 'A' && next <= 'Z') || next == '_' || next == '>' || next == '/' {
			return e.tryHandleJSX()
		}
	}
	return false
}

func (e *extractor) skipBraceExpression() {
	if e.pos < e.len && e.src[e.pos] == '{' {
		e.advance()
	}
	depth := 1
	for e.pos < e.len && depth > 0 {
		e.skipBraceChar(&depth)
	}
}

func (e *extractor) readIdentifier() string {
	start := e.pos
	for e.pos < e.len {
		ch := e.src[e.pos]
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') ||
			ch == '_' || ch == '-' || ch == '.' || ch == ':' {
			e.advance()
		} else {
			break
		}
	}
	return string(e.src[start:e.pos])
}

func (e *extractor) readAttributeName() string {
	start := e.pos
	for e.pos < e.len {
		ch := e.src[e.pos]
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') ||
			ch == '_' || ch == '-' || ch == '.' || ch == ':' || ch == '$' || ch == '@' {
			e.advance()
		} else {
			break
		}
	}
	return string(e.src[start:e.pos])
}

func (e *extractor) readQuotedAttribute(quote byte) string {
	start := e.pos
	e.advance()
	for e.pos < e.len && e.src[e.pos] != quote {
		if e.src[e.pos] == '\\' && e.pos+1 < e.len {
			e.advance()
		}
		e.advance()
	}
	if e.pos < e.len && e.src[e.pos] == quote {
		e.advance()
	}
	return string(e.src[start:e.pos])
}

func (e *extractor) readBraceAttribute() string {
	start := e.pos
	e.advance()
	depth := 1
	for e.pos < e.len && depth > 0 {
		c := e.src[e.pos]
		switch c {
		case '{':
			depth++
			e.advance()
		case '}':
			depth--
			e.advance()
		case '"', '\'':
			e.skipStringLiteral(c)
		case '`':
			e.skipJSTemplateLiteral()
		default:
			e.advance()
		}
	}
	return string(e.src[start:e.pos])
}

func (e *extractor) readUnquotedAttribute() string {
	start := e.pos
	for e.pos < e.len {
		c := e.src[e.pos]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '>' || c == '/' {
			break
		}
		e.advance()
	}
	return string(e.src[start:e.pos])
}

func (e *extractor) readAttributeValue() string {
	if e.pos >= e.len {
		return ""
	}

	ch := e.src[e.pos]
	switch ch {
	case '"', '\'':
		return e.readQuotedAttribute(ch)
	case '{':
		return e.readBraceAttribute()
	default:
		return e.readUnquotedAttribute()
	}
}

func (e *extractor) advanceBytes(n int) {
	for i := 0; i < n && e.pos < e.len; i++ {
		e.advance()
	}
}

func (e *extractor) skipWhitespace() {
	for e.pos < e.len {
		ch := e.src[e.pos]
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			e.advance()
		} else {
			break
		}
	}
}

func (e *extractor) advance() {
	if e.pos < e.len {
		ch := e.src[e.pos]
		if ch != ' ' && ch != '\t' && ch != '\n' && ch != '\r' {
			e.lastByte = ch
		}
		if ch == '\n' {
			e.line++
			e.col = 1
		} else {
			e.col++
		}
		e.pos++
	}
}

func (e *extractor) rollback(pos int, line int, col int) {
	e.pos = pos
	e.line = line
	e.col = col
}
