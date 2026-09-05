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
	if l.pos+3 > l.len {
		return
	}

	if l.src[l.pos] == '-' && l.src[l.pos+1] == '-' && l.src[l.pos+2] == '-' {
		// Konsumsi '---' pembuka
		l.advance()
		l.advance()
		l.advance()

		// Cari '---' penutup yang berada di baris baru
		for l.pos < l.len {
			if l.src[l.pos] == '\n' {
				l.advance()
				// Periksa apakah setelah newline ada '---'
				if l.pos+3 <= l.len && l.src[l.pos] == '-' && l.src[l.pos+1] == '-' && l.src[l.pos+2] == '-' {
					l.advance()
					l.advance()
					l.advance()
					// Konsumsi sisa baris penutup frontmatter
					for l.pos < l.len && l.src[l.pos] != '\n' {
						l.advance()
					}
					if l.pos < l.len && l.src[l.pos] == '\n' {
						l.advance()
					}
					return
				}
				continue
			}
			l.advance()
		}
	}
}

// handleTag memproses tag pembuka, penutup, fragment, comment, atau doctype.
func (l *lexer) handleTag() {
	startLine := l.line
	startCol := l.col
	l.advance() // konsumsi '<'

	if l.pos >= l.len {
		return
	}

	// 1. Komentar HTML: <!-- ... -->
	if l.matchPrefix("!--") {
		l.pos += 3
		l.col += 3
		commentStart := l.pos
		endIdx := bytes.Index(l.src[l.pos:], []byte("-->"))
		if endIdx == -1 {
			// Komentar tidak ditutup, konsumsi hingga akhir
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

	// 2. Fragment closing: </>
	if l.matchPrefix("/>") {
		l.advance()
		l.advance()
		l.bld.CloseFragment()
		return
	}

	// 3. Tag penutup: </tag>
	if l.src[l.pos] == '/' {
		l.advance() // konsumsi '/'
		l.skipWhitespace()
		tagName := l.readIdentifier()
		// Cari penutup '>'
		for l.pos < l.len && l.src[l.pos] != '>' {
			l.advance()
		}
		if l.pos < l.len && l.src[l.pos] == '>' {
			l.advance()
		}
		l.bld.CloseElement(tagName)
		return
	}

	// 4. Fragment opening: <>
	if l.src[l.pos] == '>' {
		l.advance()
		l.bld.OpenFragment(ir.Span{Line: startLine, Column: startCol, EndLine: l.line, EndColumn: l.col})
		return
	}

	// 5. Doctype / Directive: <! ... >
	if l.src[l.pos] == '!' {
		for l.pos < l.len && l.src[l.pos] != '>' {
			l.advance()
		}
		if l.pos < l.len && l.src[l.pos] == '>' {
			l.advance()
		}
		return
	}

	// 6. Tag Pembuka Elemen: <tag ...>
	tagName := l.readIdentifier()
	if tagName == "" {
		// Malformed: '<' tidak diikuti nama tag yang valid, abaikan karakter '<'
		return
	}

	// Parsing atribut
	attrs := make(map[string]string)
	var rawClasses string
	var classes []string
	var hasDynamic bool
	selfClosing := false

	for l.pos < l.len {
		l.skipWhitespace()
		if l.pos >= l.len {
			break
		}

		ch := l.src[l.pos]

		// Jika menemukan '<' baru sebelum '>' (contoh: <broken <button>):
		// Malformed recovery: buang broken tag, resinkronisasi ke '<' berikutnya
		if ch == '<' {
			return
		}

		if ch == '>' {
			l.advance()
			break
		}

		if ch == '/' && l.pos+1 < l.len && l.src[l.pos+1] == '>' {
			l.advance()
			l.advance()
			selfClosing = true
			break
		}

		// Baca nama atribut
		attrName := l.readAttributeName()
		if attrName == "" {
			l.advance()
			continue
		}

		l.skipWhitespace()
		attrVal := ""

		if l.pos < l.len && l.src[l.pos] == '=' {
			l.advance() // konsumsi '='
			l.skipWhitespace()
			attrVal = l.readAttributeValue()
		}

		attrs[attrName] = attrVal
		if attrName == "class" || attrName == "className" {
			rawClasses = attrVal
			classes, hasDynamic = parser.ExtractClasses(attrVal)
		}
	}

	span := ir.Span{
		Line:      startLine,
		Column:    startCol,
		EndLine:   l.line,
		EndColumn: l.col,
	}

	// Deteksi void element HTML (selalu self-closing)
	isVoid := voidElements[strings.ToLower(tagName)]

	if selfClosing || isVoid {
		l.bld.AddSelfClosingElement(tagName, span, rawClasses, classes, attrs, hasDynamic)
	} else {
		l.bld.OpenElement(tagName, span, rawClasses, classes, attrs, hasDynamic)
	}
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

func (l *lexer) readAttributeValue() string {
	if l.pos >= l.len {
		return ""
	}

	ch := l.src[l.pos]

	// String berkutip ganda: "..."
	if ch == '"' {
		start := l.pos
		l.advance()
		for l.pos < l.len && l.src[l.pos] != '"' {
			if l.src[l.pos] == '\\' && l.pos+1 < l.len {
				l.advance()
			}
			l.advance()
		}
		if l.pos < l.len && l.src[l.pos] == '"' {
			l.advance()
		}
		return string(l.src[start:l.pos])
	}

	// String berkutip tunggal: '...'
	if ch == '\'' {
		start := l.pos
		l.advance()
		for l.pos < l.len && l.src[l.pos] != '\'' {
			if l.src[l.pos] == '\\' && l.pos+1 < l.len {
				l.advance()
			}
			l.advance()
		}
		if l.pos < l.len && l.src[l.pos] == '\'' {
			l.advance()
		}
		return string(l.src[start:l.pos])
	}

	// Ekspresi kurung kurawal JSX: {...}
	if ch == '{' {
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

	// Atribut tanpa kutip
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
