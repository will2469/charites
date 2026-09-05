package theme

import (
	"strings"
)

// Specificity merepresentasikan bobot spesifisitas selektor CSS (A, B, C)
// sesuai W3C CSS Selectors Level 4 (Section 16).
// A = Jumlah ID selectors (#id)
// B = Jumlah Class (.class), Attribute ([attr]), dan Pseudo-class (:root, :hover, :is(), dll)
// C = Jumlah Type/Element (div, html) dan Pseudo-element (::before, ::after, :before, :after)
type Specificity struct {
	A int
	B int
	C int
}

// GreaterThan mengembalikan true jika s memiliki spesifisitas lebih tinggi dari other.
func (s Specificity) GreaterThan(other Specificity) bool {
	if s.A != other.A {
		return s.A > other.A
	}
	if s.B != other.B {
		return s.B > other.B
	}
	return s.C > other.C
}

// Equals mengembalikan true jika s dan other memiliki nilai spesifisitas identik.
func (s Specificity) Equals(other Specificity) bool {
	return s.A == other.A && s.B == other.B && s.C == other.C
}

// Add menjumlahkan dua nilai spesifisitas.
func (s Specificity) Add(other Specificity) Specificity {
	return Specificity{
		A: s.A + other.A,
		B: s.B + other.B,
		C: s.C + other.C,
	}
}

// Max mengembalikan nilai spesifisitas tertinggi antara s dan other.
func (s Specificity) Max(other Specificity) Specificity {
	if s.GreaterThan(other) {
		return s
	}
	return other
}

// ComputeSpecificity mem-parse string selektor CSS dan menghitung bobot spesifisitasnya
// sesuai spesifikasi W3C CSS Selectors Level 4.
// Jika selektor berupa selector list (dipisahkan koma), mengembalikan spesifisitas tertinggi
// dari daftar selektor kompleks (sesuai semantik :is() / ruleset matching).
func ComputeSpecificity(selector string) Specificity {
	p := newSelectorParser([]byte(selector))
	return p.parseSelectorList()
}

type selectorScanner struct {
	src []byte
	idx int
}

func (s *selectorScanner) skipWhitespace() {
	for s.idx < len(s.src) && isWhitespace(s.src[s.idx]) {
		s.idx++
	}
}

func (s *selectorScanner) startsIdent() bool {
	if s.idx >= len(s.src) {
		return false
	}
	b0 := s.src[s.idx]
	var b1, b2 byte
	if s.idx+1 < len(s.src) {
		b1 = s.src[s.idx+1]
	}
	if s.idx+2 < len(s.src) {
		b2 = s.src[s.idx+2]
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

func (s *selectorScanner) consumeIdent() string {
	start := s.idx
	for s.idx < len(s.src) {
		b := s.src[s.idx]
		if isIdentChar(b) {
			s.idx++
			continue
		}
		if b == '\\' && s.idx+1 < len(s.src) && isValidEscape(b, s.src[s.idx+1]) {
			s.idx += 2 // lewati \ dan karakter yang di-escape
			if s.idx-1 < len(s.src) && isHexDigit(s.src[s.idx-1]) {
				hexCount := 1
				for s.idx < len(s.src) && isHexDigit(s.src[s.idx]) && hexCount < 6 {
					s.idx++
					hexCount++
				}
				if s.idx < len(s.src) && isWhitespace(s.src[s.idx]) {
					s.idx++
				}
			}
			continue
		}
		break
	}
	return string(s.src[start:s.idx])
}

// consumeBracket mengonsumsi atribut [...] dengan penanganan string dan escape aman.
func (s *selectorScanner) consumeBracket() string {
	if s.idx >= len(s.src) || s.src[s.idx] != '[' {
		return ""
	}
	start := s.idx
	s.idx++ // consume '['

	var inQuote byte
	for s.idx < len(s.src) {
		b := s.src[s.idx]
		if inQuote != 0 {
			if b == '\\' {
				s.idx += 2
				continue
			}
			if b == inQuote {
				inQuote = 0
			}
			s.idx++
			continue
		}

		if b == '"' || b == '\'' {
			inQuote = b
			s.idx++
			continue
		}

		if b == ']' {
			s.idx++
			return string(s.src[start:s.idx])
		}
		s.idx++
	}

	return string(s.src[start:s.idx])
}

// consumeParentheses mengonsumsi konten di dalam (...) dengan nested depth tracking.
func (s *selectorScanner) consumeParentheses() string {
	if s.idx >= len(s.src) || s.src[s.idx] != '(' {
		return ""
	}
	s.idx++ // consume '('
	start := s.idx

	depth := 1
	var inQuote byte
	for s.idx < len(s.src) {
		b := s.src[s.idx]
		if inQuote != 0 {
			if b == '\\' {
				s.idx += 2
				continue
			}
			if b == inQuote {
				inQuote = 0
			}
			s.idx++
			continue
		}

		switch b {
		case '"', '\'':
			inQuote = b
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				content := string(s.src[start:s.idx])
				s.idx++ // consume ')'
				return content
			}
		}
		s.idx++
	}

	return string(s.src[start:s.idx])
}

type selectorParser struct {
	sc *selectorScanner
}

func newSelectorParser(src []byte) *selectorParser {
	return &selectorParser{
		sc: &selectorScanner{src: src, idx: 0},
	}
}

// parseSelectorList mem-parse daftar selektor (dipisahkan koma) dan mengembalikan spesifisitas maksimum.
func (p *selectorParser) parseSelectorList() Specificity {
	var maxSpec Specificity
	first := true

	for p.sc.idx < len(p.sc.src) {
		p.sc.skipWhitespace()
		if p.sc.idx >= len(p.sc.src) {
			break
		}
		if p.sc.src[p.sc.idx] == ',' {
			p.sc.idx++
			continue
		}

		complexSpec := p.parseComplexSelector()
		if first {
			maxSpec = complexSpec
			first = false
		} else {
			maxSpec = maxSpec.Max(complexSpec)
		}
	}

	return maxSpec
}

// parseComplexSelector mem-parse urutan compound selectors yang dihubungkan kombinator.
func (p *selectorParser) parseComplexSelector() Specificity {
	var total Specificity

	for p.sc.idx < len(p.sc.src) {
		p.sc.skipWhitespace()
		if p.sc.idx >= len(p.sc.src) {
			break
		}
		b := p.sc.src[p.sc.idx]
		if b == ',' {
			break // Akhir dari complex selector dalam daftar
		}

		// Kombinator CSS: '>', '+', '~', '||' (spesifisitas 0)
		if b == '>' || b == '+' || b == '~' {
			p.sc.idx++
			continue
		}
		if b == '|' && p.sc.idx+1 < len(p.sc.src) && p.sc.src[p.sc.idx+1] == '|' {
			p.sc.idx += 2
			continue
		}

		comp := p.parseCompoundComponent()
		total = total.Add(comp)
	}

	return total
}

// parseCompoundComponent mem-parse satu komponen dalam compound selector.
func (p *selectorParser) parseCompoundComponent() Specificity {
	if p.sc.idx >= len(p.sc.src) {
		return Specificity{}
	}

	b := p.sc.src[p.sc.idx]

	switch b {
	case '#':
		p.sc.idx++
		_ = p.sc.consumeIdent()
		return Specificity{A: 1}
	case '.':
		p.sc.idx++
		_ = p.sc.consumeIdent()
		return Specificity{B: 1}
	case '[':
		_ = p.sc.consumeBracket()
		return Specificity{B: 1}
	case '&':
		p.sc.idx++
		return Specificity{B: 1}
	case '*':
		return p.parseUniversalSelector()
	case ':':
		return p.parseColonSelector()
	}

	if p.sc.startsIdent() {
		return p.parseTypeSelector()
	}

	p.sc.idx++
	return Specificity{}
}

func (p *selectorParser) parseUniversalSelector() Specificity {
	p.sc.idx++ // lewati '*'
	if p.sc.idx < len(p.sc.src) && p.sc.src[p.sc.idx] == '|' {
		p.sc.idx++ // lewati '|'
		if p.sc.startsIdent() {
			_ = p.sc.consumeIdent()
			return Specificity{C: 1} // Type selector dengan wildcard namespace *|tag
		}
	}
	return Specificity{}
}

func (p *selectorParser) parseColonSelector() Specificity {
	p.sc.idx++ // lewati ':'
	if p.sc.idx < len(p.sc.src) && p.sc.src[p.sc.idx] == ':' {
		return p.parseDoubleColonPseudoElement()
	}

	ident := strings.ToLower(p.sc.consumeIdent())

	// Legacy CSS2 pseudo-elements yang diizinkan dengan single-colon
	switch ident {
	case "before", "after", "first-line", "first-letter":
		return Specificity{C: 1}
	}

	// Functional pseudo-classes: ident(...)
	if p.sc.idx < len(p.sc.src) && p.sc.src[p.sc.idx] == '(' {
		args := p.sc.consumeParentheses()
		return parseFunctionalPseudoClass(ident, args)
	}

	return Specificity{B: 1}
}

func (p *selectorParser) parseDoubleColonPseudoElement() Specificity {
	p.sc.idx++ // lewati ':' kedua
	ident := strings.ToLower(p.sc.consumeIdent())
	if ident == "slotted" && p.sc.idx < len(p.sc.src) && p.sc.src[p.sc.idx] == '(' {
		inner := p.sc.consumeParentheses()
		innerSpec := ComputeSpecificity(inner)
		return Specificity{C: 1}.Add(innerSpec)
	}
	if p.sc.idx < len(p.sc.src) && p.sc.src[p.sc.idx] == '(' {
		_ = p.sc.consumeParentheses()
	}
	return Specificity{C: 1}
}

func parseFunctionalPseudoClass(ident, args string) Specificity {
	switch ident {
	case "where":
		return Specificity{}
	case "is", "not", "has":
		return ComputeSpecificity(args)
	case "nth-child", "nth-last-child":
		spec := Specificity{B: 1}
		if ofIdx := strings.Index(args, " of "); ofIdx != -1 {
			selectorList := args[ofIdx+4:]
			spec = spec.Add(ComputeSpecificity(selectorList))
		}
		return spec
	default:
		return Specificity{B: 1}
	}
}

func (p *selectorParser) parseTypeSelector() Specificity {
	_ = p.sc.consumeIdent()
	if p.sc.idx < len(p.sc.src) && p.sc.src[p.sc.idx] == '|' {
		p.sc.idx++
		if p.sc.startsIdent() {
			_ = p.sc.consumeIdent()
		}
	}
	return Specificity{C: 1}
}
