// Package parser provides lexing and parsing for TSCN files.
package parser

// TokenType represents the type of a lexical token.
type TokenType int

const (
	TokenEOF TokenType = iota
	TokenError
	TokenNewline
	TokenComment // ; comment

	// Delimiters
	TokenLBracket // [
	TokenRBracket // ]
	TokenLParen   // (
	TokenRParen   // )
	TokenLBrace   // {
	TokenRBrace   // }

	// Operators
	TokenEquals // =
	TokenColon  // :
	TokenComma  // ,
	TokenSlash  // /

	// Literals
	TokenIdent  // identifiers like node, type, Vector3
	TokenString // "string literal"
	TokenNumber // 123, 1.5, -2.5e10, inf, nan
	TokenBool   // true, false
	TokenNull   // null
)

// Token represents a lexical token.
type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
	Offset int
	Length int
}

// String returns the string representation of a token type.
func (t TokenType) String() string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenError:
		return "Error"
	case TokenNewline:
		return "Newline"
	case TokenComment:
		return "Comment"
	case TokenLBracket:
		return "["
	case TokenRBracket:
		return "]"
	case TokenLParen:
		return "("
	case TokenRParen:
		return ")"
	case TokenLBrace:
		return "{"
	case TokenRBrace:
		return "}"
	case TokenEquals:
		return "="
	case TokenColon:
		return ":"
	case TokenComma:
		return ","
	case TokenSlash:
		return "/"
	case TokenIdent:
		return "Ident"
	case TokenString:
		return "String"
	case TokenNumber:
		return "Number"
	case TokenBool:
		return "Bool"
	case TokenNull:
		return "Null"
	default:
		return "Unknown"
	}
}

// Range represents a source location range.
type Range struct {
	Start Position
	End   Position
}

// Position represents a position in source code.
type Position struct {
	Line   int // 0-based line number
	Column int // 0-based column (character offset in line)
	Offset int // byte offset from start of file
}
