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

// parseRuleOrDeclaration membedakan antara AtRule, StyleRule, atau Declaration pada scope saat ini
// secara grammar-aware dengan melacak kedalaman tanda kurung, kurung siku, dan kurung kurawal.
func (p *Parser) parseRuleOrDeclaration() (Rule, *Declaration) {
	// Lewati titik koma berlebih (e.g. ;;)
	for p.curr.Type == TokenSemicolon {
		p.advance()
	}

	if p.curr.Type == TokenEOF || p.curr.Type == TokenCloseBrace {
		return nil, nil
	}

	// 1. AtRule (@layer, @media, @supports, @theme, @import, dll.)
	if p.curr.Type == TokenAtKeyword {
		return p.parseAtRule(), nil
	}

	startLoc := p.curr.Span.Start
	var tokens []Token

	parenDepth := 0
	bracketDepth := 0
	braceDepth := 0

	for p.curr.Type != TokenEOF {
		// Periksa delimiter pemutus pada kedalaman seimbang (depth 0)
		if parenDepth == 0 && bracketDepth == 0 && braceDepth == 0 {
			if p.curr.Type == TokenOpenBrace {
				// Ini adalah StyleRule (blok selektor terkualifikasi)
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
				// Akhir blok (deklarasi terakhir tanpa titik koma penutup, misal: { color: red })
				if len(tokens) > 0 {
					decl := p.tokensToDeclaration(tokens, startLoc)
					return nil, decl
				}
				return nil, nil
			}
		}

		// Lacak kedalaman nesting
		switch p.curr.Type {
		case TokenOpenParen:
			parenDepth++
		case TokenCloseParen:
			if parenDepth > 0 {
				parenDepth--
			}
		case TokenOpenBracket:
			bracketDepth++
		case TokenCloseBracket:
			if bracketDepth > 0 {
				bracketDepth--
			}
		case TokenOpenBrace:
			braceDepth++
		case TokenCloseBrace:
			if braceDepth > 0 {
				braceDepth--
			}
		}

		tokens = append(tokens, p.curr)
		p.advance()
	}

	if len(tokens) > 0 {
		decl := p.tokensToDeclaration(tokens, startLoc)
		return nil, decl
	}

	return nil, nil
}

func (p *Parser) parseAtRule() Rule {
	startLoc := p.curr.Span.Start
	name := p.curr.Value
	p.advance() // konsumsi AtKeyword

	var preludeTokens []Token
	parenDepth := 0
	bracketDepth := 0
	braceDepth := 0

	for p.curr.Type != TokenEOF {
		if parenDepth == 0 && bracketDepth == 0 && braceDepth == 0 {
			// Jika diakhiri ';' (seperti @import "tailwindcss"; atau @layer utilities;)
			if p.curr.Type == TokenSemicolon {
				prelude := p.sliceTokens(preludeTokens)
				endLoc := p.curr.Span.End
				p.advance() // konsumsi ';'
				return &AtRule{
					Name:    name,
					Prelude: prelude,
					Span:    SourceSpan{Start: startLoc, End: endLoc},
				}
			}

			// Jika diakhiri '{', parse isi blok
			if p.curr.Type == TokenOpenBrace {
				prelude := p.sliceTokens(preludeTokens)
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
		}

		switch p.curr.Type {
		case TokenOpenParen:
			parenDepth++
		case TokenCloseParen:
			if parenDepth > 0 {
				parenDepth--
			}
		case TokenOpenBracket:
			bracketDepth++
		case TokenCloseBracket:
			if bracketDepth > 0 {
				bracketDepth--
			}
		case TokenOpenBrace:
			braceDepth++
		case TokenCloseBrace:
			if braceDepth > 0 {
				braceDepth--
			}
		}

		preludeTokens = append(preludeTokens, p.curr)
		p.advance()
	}

	return &AtRule{Name: name, Prelude: p.sliceTokens(preludeTokens), Span: SourceSpan{Start: startLoc, End: p.curr.Span.End}}
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

	// Cari posisi colon ':' pada kedalaman seimbang (depth 0)
	parenDepth := 0
	bracketDepth := 0
	braceDepth := 0
	colonIdx := -1

	for i, tok := range tokens {
		switch tok.Type {
		case TokenOpenParen:
			parenDepth++
		case TokenCloseParen:
			if parenDepth > 0 {
				parenDepth--
			}
		case TokenOpenBracket:
			bracketDepth++
		case TokenCloseBracket:
			if bracketDepth > 0 {
				bracketDepth--
			}
		case TokenOpenBrace:
			braceDepth++
		case TokenCloseBrace:
			if braceDepth > 0 {
				braceDepth--
			}
		case TokenColon:
			if parenDepth == 0 && bracketDepth == 0 && braceDepth == 0 && colonIdx == -1 {
				colonIdx = i
			}
		}
	}

	if colonIdx <= 0 {
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
