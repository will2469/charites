package theme

import (
	"strings"
)

// Parser mem-parse stream token CSS menjadi struktur AST StyleSheet.
type Parser struct {
	src   []byte
	lexer *Lexer
	curr  Token
}

// Parse mem-parse data byte CSS menjadi StyleSheet.
func Parse(src []byte) (*StyleSheet, error) {
	p := &Parser{
		src:   src,
		lexer: NewLexer(src),
	}
	p.advance()

	sheet := &StyleSheet{
		Rules: make([]Rule, 0),
	}

	startLoc := p.curr.Span.Start

	for p.curr.Type != TokenEOF {
		rule, decl := p.parseRuleOrDeclaration()
		if rule != nil {
			sheet.Rules = append(sheet.Rules, rule)
		} else if decl != nil {
			sheet.Declarations = append(sheet.Declarations, *decl)
		}
	}

	sheet.Span = SourceSpan{Start: startLoc, End: p.curr.Span.End}
	return sheet, nil
}

func (p *Parser) advance() {
	for {
		tok := p.lexer.NextToken()
		if tok.Type != TokenWhitespace && tok.Type != TokenComment {
			p.curr = tok
			return
		}
	}
}

func (p *Parser) sliceTokens(tokens []Token) string {
	if len(tokens) == 0 {
		return ""
	}
	start := tokens[0].StartOffset
	end := tokens[len(tokens)-1].EndOffset
	if start >= 0 && end <= len(p.src) && start <= end {
		return strings.TrimSpace(string(p.src[start:end]))
	}
	return tokensToString(tokens)
}

// parseRuleOrDeclaration membedakan antara AtRule, StyleRule, atau Declaration pada scope saat ini.
func (p *Parser) parseRuleOrDeclaration() (Rule, *Declaration) {
	if p.curr.Type == TokenEOF {
		return nil, nil
	}

	// 1. AtRule (@layer, @media, @supports, @theme, @import, dll.)
	if p.curr.Type == TokenAtKeyword {
		return p.parseAtRule(), nil
	}

	// 2. Baca token hingga menemukan ':' atau '{' atau ';' untuk menentukan apakah ini deklarasi atau rule
	startLoc := p.curr.Span.Start
	var tokens []Token

	for p.curr.Type != TokenEOF && p.curr.Type != TokenOpenBrace && p.curr.Type != TokenSemicolon && p.curr.Type != TokenCloseBrace {
		tokens = append(tokens, p.curr)
		p.advance()
	}

	if p.curr.Type == TokenOpenBrace {
		// Ini adalah StyleRule (blok selektor)
		selector := p.sliceTokens(tokens)
		p.advance() // konsumsi '{'
		styleRule := &StyleRule{
			Selector:     selector,
			Declarations: make([]Declaration, 0),
			Rules:        make([]Rule, 0),
		}

		p.parseBlockContents(&styleRule.Declarations, &styleRule.Rules)
		styleRule.Span = SourceSpan{Start: startLoc, End: p.curr.Span.End}
		if p.curr.Type == TokenCloseBrace {
			p.advance() // konsumsi '}'
		}
		return styleRule, nil
	}

	if p.curr.Type == TokenSemicolon {
		// Ini adalah deklarasi properti (prop: val;)
		decl := p.tokensToDeclaration(tokens, startLoc)
		p.advance() // konsumsi ';'
		return nil, decl
	}

	if p.curr.Type == TokenCloseBrace {
		// Akhir blok, jika ada tokens tersisa parse sebagai deklarasi
		if len(tokens) > 0 {
			decl := p.tokensToDeclaration(tokens, startLoc)
			return nil, decl
		}
		return nil, nil
	}

	return nil, nil
}

func (p *Parser) parseAtRule() Rule {
	startLoc := p.curr.Span.Start
	name := p.curr.Value
	p.advance() // konsumsi AtKeyword

	var preludeTokens []Token
	for p.curr.Type != TokenEOF && p.curr.Type != TokenOpenBrace && p.curr.Type != TokenSemicolon {
		preludeTokens = append(preludeTokens, p.curr)
		p.advance()
	}

	prelude := p.sliceTokens(preludeTokens)

	// Jika diakhiri ';' (seperti @import "tailwindcss";)
	if p.curr.Type == TokenSemicolon {
		endLoc := p.curr.Span.End
		p.advance()
		return &AtRule{
			Name:    name,
			Prelude: prelude,
			Span:    SourceSpan{Start: startLoc, End: endLoc},
		}
	}

	// Jika diakhiri '{', parse isi blok
	if p.curr.Type == TokenOpenBrace {
		p.advance() // konsumsi '{'
		atRule := &AtRule{
			Name:         name,
			Prelude:      prelude,
			Declarations: make([]Declaration, 0),
			Rules:        make([]Rule, 0),
		}

		p.parseBlockContents(&atRule.Declarations, &atRule.Rules)
		atRule.Span = SourceSpan{Start: startLoc, End: p.curr.Span.End}
		if p.curr.Type == TokenCloseBrace {
			p.advance() // konsumsi '}'
		}
		return atRule
	}

	return &AtRule{Name: name, Prelude: prelude, Span: SourceSpan{Start: startLoc, End: p.curr.Span.End}}
}

// parseBlockContents mem-parse seluruh deklarasi dan aturan bersarang di dalam kurung kurawal.
func (p *Parser) parseBlockContents(decls *[]Declaration, rules *[]Rule) {
	for p.curr.Type != TokenEOF && p.curr.Type != TokenCloseBrace {
		subRule, subDecl := p.parseRuleOrDeclaration()
		if subRule != nil {
			*rules = append(*rules, subRule)
		} else if subDecl != nil {
			*decls = append(*decls, *subDecl)
		}
	}
}

func tokensToString(tokens []Token) string {
	if len(tokens) == 0 {
		return ""
	}
	var sb strings.Builder
	for i, tok := range tokens {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(tok.Value)
	}
	return strings.TrimSpace(sb.String())
}

func (p *Parser) tokensToDeclaration(tokens []Token, startLoc SourceLocation) *Declaration {
	if len(tokens) == 0 {
		return nil
	}

	// Cari posisi colon ':'
	colonIdx := -1
	for i, tok := range tokens {
		if tok.Type == TokenColon {
			colonIdx = i
			break
		}
	}

	if colonIdx == -1 || colonIdx == 0 {
		return nil
	}

	prop := p.sliceTokens(tokens[:colonIdx])
	val := p.sliceTokens(tokens[colonIdx+1:])

	endLoc := tokens[len(tokens)-1].Span.End

	return &Declaration{
		Property: prop,
		Value:    val,
		RawValue: val,
		Span:     SourceSpan{Start: startLoc, End: endLoc},
	}
}
