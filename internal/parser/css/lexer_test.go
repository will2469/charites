package css_test

import (
	"testing"

	"github.com/will2469/charites/internal/parser/css"
)

func TestLexer_CSSEscapes(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedType  css.TokenType
		expectedValue string
		unescaped     string
	}{
		{
			name:          "Escaped colon in class",
			input:         `\:focus`,
			expectedType:  css.TokenIdent,
			expectedValue: `\:focus`,
			unescaped:     `:focus`,
		},
		{
			name:          "Escaped colon in custom property name",
			input:         `--tw-bg-opacity\:1`,
			expectedType:  css.TokenIdent,
			expectedValue: `--tw-bg-opacity\:1`,
			unescaped:     `--tw-bg-opacity:1`,
		},
		{
			name:          "Hex escape without whitespace",
			input:         `--color\31`,
			expectedType:  css.TokenIdent,
			expectedValue: `--color\31`,
			unescaped:     `--color1`,
		},
		{
			name:          "Hex escape with terminating whitespace",
			input:         `--color\31 00`,
			expectedType:  css.TokenIdent,
			expectedValue: `--color\31 00`,
			unescaped:     `--color100`,
		},
		{
			name:          "6-digit hex escape",
			input:         `--brand\000031`,
			expectedType:  css.TokenIdent,
			expectedValue: `--brand\000031`,
			unescaped:     `--brand1`,
		},
		{
			name:          "Negative dimension",
			input:         `-10px`,
			expectedType:  css.TokenDimension,
			expectedValue: `-10px`,
			unescaped:     `-10px`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := css.NewLexer([]byte(tt.input))
			tok := lexer.NextToken()
			if tok.Type != tt.expectedType {
				t.Fatalf("expected token type %v, got %v", tt.expectedType, tok.Type)
			}
			if tok.Value != tt.expectedValue {
				t.Errorf("expected token value %q, got %q", tt.expectedValue, tok.Value)
			}
			unescaped := css.UnescapeCSS(tok.Value)
			if unescaped != tt.unescaped {
				t.Errorf("expected unescaped %q, got %q", tt.unescaped, unescaped)
			}
		})
	}
}
