package css

// TokenType mendefinisikan jenis token CSS hasil pemindaian lexer.
type TokenType uint8

// Daftar jenis token leksikal CSS yang dihasilkan oleh lexer.
const (
	TokenEOF TokenType = iota
	TokenWhitespace
	TokenComment
	TokenIdent
	TokenAtKeyword    // e.g. @layer, @media, @theme
	TokenString       // e.g. "hello;world", 'text'
	TokenHash         // e.g. #123456, #ef4444
	TokenNumber       // e.g. 10, 0.5
	TokenDimension    // e.g. 10px, 0.5rem
	TokenPercentage   // e.g. 10%, 5%
	TokenColon        // :
	TokenSemicolon    // ;
	TokenComma        // ,
	TokenOpenBrace    // {
	TokenCloseBrace   // }
	TokenOpenParen    // (
	TokenCloseParen   // )
	TokenOpenBracket  // [
	TokenCloseBracket // ]
	TokenDelim        // any other single char e.g. /, &, +, >
)

// SourceLocation mencatat koordinat baris dan kolom 1-indexed.
type SourceLocation struct {
	Line   int
	Column int
}

// SourceSpan mencatat rentang posisi awal dan akhir token atau AST node.
type SourceSpan struct {
	Start SourceLocation
	End   SourceLocation
}

// Token merepresentasikan token leksikal CSS tunggal.
type Token struct {
	Type        TokenType
	Value       string
	StartOffset int
	EndOffset   int
	Span        SourceSpan
}
