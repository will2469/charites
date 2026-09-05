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

func (e *extractor) tryHandleJSX() bool {
	startLine := e.line
	startCol := e.col
	savedPos := e.pos

	e.advance() // konsumsi '<'
	if e.pos >= e.len {
		return false
	}

	next := e.src[e.pos]

	// 1. Fragment closing: </>
	if next == '/' && e.pos+1 < e.len && e.src[e.pos+1] == '>' {
		e.advance()
		e.advance()
		e.bld.CloseFragment()
		return true
	}

	// 2. Tag penutup: </tag>
	if next == '/' {
		e.advance() // konsumsi '/'
		e.skipWhitespace()
		tagName := e.readIdentifier()
		if tagName == "" {
			// Cacat: </ tanpa nama tag
			return false
		}
		// Cari '>' penutup
		for e.pos < e.len && e.src[e.pos] != '>' {
			e.advance()
		}
		if e.pos < e.len && e.src[e.pos] == '>' {
			e.advance()
		}
		e.bld.CloseElement(tagName)
		return true
	}

	// 3. Fragment opening: <>
	if next == '>' {
		e.advance()
		e.bld.OpenFragment(ir.Span{Line: startLine, Column: startCol, EndLine: e.line, EndColumn: e.col})
		return true
	}

	// 4. Disambiguasi pembuka: Jika '<' langsung menempel pada karakter identifier sebelumnya tanpa spasi
	// (contoh: Map<string> atau x<y), ini adalah generics atau perbandingan, bukan tag JSX.
	if savedPos > 0 && e_isIdentifierChar(e.src[savedPos-1]) {
		e.rollback(savedPos, startLine, startCol)
		return false
	}

	// Tag pembuka harus diawali huruf atau underscore
	if !((next >= 'a' && next <= 'z') || (next >= 'A' && next <= 'Z') || next == '_') {
		// Bukan tag JSX (misal: < 10 atau <= atau <<)
		e.rollback(savedPos, startLine, startCol)
		return false
	}

	tagName := e.readIdentifier()
	if tagName == "" {
		e.rollback(savedPos, startLine, startCol)
		return false
	}

	// Parsing atribut hingga '>' atau '/>'
	attrs := make(map[string]string)
	var rawClasses string
	var classes []string
	var hasDynamic bool
	selfClosing := false
	isValidJSX := false

	for e.pos < e.len {
		e.skipWhitespace()
		if e.pos >= e.len {
			break
		}

		c := e.src[e.pos]

		// Jika ada '<' sebelum tag ditutup: recovery (broken tag)
		if c == '<' {
			return true
		}

		// Karakter ilegal di header tag JSX yang menandakan ini bukan JSX
		// (misal: a < b && c > d atau if (a < b) atau const x = <T,>(...) atau const fn = <T>(x: T))
		if c == ';' || c == ')' || c == '(' || c == ',' ||
			(c == '&' && e.pos+1 < e.len && e.src[e.pos+1] == '&') ||
			(c == '|' && e.pos+1 < e.len && e.src[e.pos+1] == '|') {
			e.rollback(savedPos, startLine, startCol)
			return false
		}

		if c == '>' {
			// Disambiguasi: jika nama tag hanya 1 huruf kapital (contoh: <T>) dan langsung diikuti '(',
			// ini adalah generic type parameter arrow function TypeScript: <T>(x: T) => ...
			if len(tagName) == 1 && tagName[0] >= 'A' && tagName[0] <= 'Z' {
				peek := e.pos + 1
				for peek < e.len && (e.src[peek] == ' ' || e.src[peek] == '\t' || e.src[peek] == '\n' || e.src[peek] == '\r') {
					peek++
				}
				if peek < e.len && e.src[peek] == '(' {
					e.rollback(savedPos, startLine, startCol)
					return false
				}
			}
			e.advance()
			isValidJSX = true
			break
		}

		if c == '/' && e.pos+1 < e.len && e.src[e.pos+1] == '>' {
			e.advance()
			e.advance()
			selfClosing = true
			isValidJSX = true
			break
		}

		// Spread attribute: {...props}
		if c == '{' {
			e.skipBraceExpression()
			continue
		}

		attrName := e.readAttributeName()
		if attrName == "" {
			e.advance()
			continue
		}

		e.skipWhitespace()
		attrVal := ""

		if e.pos < e.len && e.src[e.pos] == '=' {
			e.advance() // konsumsi '='
			e.skipWhitespace()
			attrVal = e.readAttributeValue()
		}

		attrs[attrName] = attrVal
		if attrName == "className" || attrName == "class" {
			rawClasses = attrVal
			classes, hasDynamic = parser.ExtractClasses(attrVal)
		}
	}

	if !isValidJSX {
		return true
	}

	span := ir.Span{
		Line:      startLine,
		Column:    startCol,
		EndLine:   e.line,
		EndColumn: e.col,
	}

	isVoid := voidElements[strings.ToLower(tagName)]

	if selfClosing || isVoid {
		e.bld.AddSelfClosingElement(tagName, span, rawClasses, classes, attrs, hasDynamic)
	} else {
		e.bld.OpenElement(tagName, span, rawClasses, classes, attrs, hasDynamic)
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

func (e *extractor) skipBraceExpression() {
	if e.pos < e.len && e.src[e.pos] == '{' {
		e.advance()
	}
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
			// Jika di dalam kurung kurawal ada tag JSX pembuka:
			if c == '<' && e.pos+1 < e.len {
				next := e.src[e.pos+1]
				if (next >= 'a' && next <= 'z') || (next >= 'A' && next <= 'Z') || next == '_' || next == '>' || next == '/' {
					if e.tryHandleJSX() {
						continue
					}
				}
			}
			e.advance()
		}
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

func (e *extractor) readAttributeValue() string {
	if e.pos >= e.len {
		return ""
	}

	ch := e.src[e.pos]

	if ch == '"' || ch == '\'' {
		start := e.pos
		e.advance()
		for e.pos < e.len && e.src[e.pos] != ch {
			if e.src[e.pos] == '\\' && e.pos+1 < e.len {
				e.advance()
			}
			e.advance()
		}
		if e.pos < e.len && e.src[e.pos] == ch {
			e.advance()
		}
		return string(e.src[start:e.pos])
	}

	if ch == '{' {
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

	// Atribut tanpa kutip
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
