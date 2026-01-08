package parser

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Lexer tokenizes TSCN source code.
type Lexer struct {
	input       string
	pos         int // current position in input (byte offset)
	line        int // current line number (0-based)
	column      int // current column (0-based)
	start       int // start position of current token
	startLine   int
	startColumn int
}

// NewLexer creates a new lexer for the given input.
func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		pos:    0,
		line:   0,
		column: 0,
	}
}

// Next returns the next token.
func (l *Lexer) Next() Token {
	l.skipWhitespace()

	l.start = l.pos
	l.startLine = l.line
	l.startColumn = l.column

	if l.pos >= len(l.input) {
		return l.makeToken(TokenEOF, "")
	}

	ch := l.peek()

	// Single-character tokens
	switch ch {
	case '[':
		l.advance()
		return l.makeToken(TokenLBracket, "[")
	case ']':
		l.advance()
		return l.makeToken(TokenRBracket, "]")
	case '(':
		l.advance()
		return l.makeToken(TokenLParen, "(")
	case ')':
		l.advance()
		return l.makeToken(TokenRParen, ")")
	case '{':
		l.advance()
		return l.makeToken(TokenLBrace, "{")
	case '}':
		l.advance()
		return l.makeToken(TokenRBrace, "}")
	case '=':
		l.advance()
		return l.makeToken(TokenEquals, "=")
	case ':':
		l.advance()
		return l.makeToken(TokenColon, ":")
	case ',':
		l.advance()
		return l.makeToken(TokenComma, ",")
	case '/':
		l.advance()
		return l.makeToken(TokenSlash, "/")
	case '\n':
		l.advance()
		return l.makeToken(TokenNewline, "\n")
	case ';':
		return l.scanComment()
	case '"':
		return l.scanString()
	case '&':
		// StringName literal: &"name" - treat as a string
		if l.pos+1 < len(l.input) && l.input[l.pos+1] == '"' {
			l.advance() // consume '&'
			return l.scanString()
		}
		l.advance()
		return l.makeToken(TokenError, "&")
	case '-', '+':
		// Could be a number
		if l.pos+1 < len(l.input) {
			next := l.input[l.pos+1]
			if next >= '0' && next <= '9' || next == '.' {
				return l.scanNumber()
			}
		}
		// Otherwise, treat as identifier start (unlikely in TSCN)
		return l.scanIdentifier()
	}

	// Numbers
	if ch >= '0' && ch <= '9' || ch == '.' {
		return l.scanNumber()
	}

	// Identifiers and keywords
	if isIdentStart(ch) {
		return l.scanIdentifier()
	}

	// Unknown character
	l.advance()
	return l.makeToken(TokenError, string(ch))
}

// Tokenize returns all tokens from the input.
func (l *Lexer) Tokenize() []Token {
	var tokens []Token
	for {
		tok := l.Next()
		tokens = append(tokens, tok)
		if tok.Type == TokenEOF {
			break
		}
		// Don't stop on error tokens - continue tokenizing to allow recovery
	}
	return tokens
}

func (l *Lexer) peek() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) peekN(n int) byte {
	if l.pos+n >= len(l.input) {
		return 0
	}
	return l.input[l.pos+n]
}

func (l *Lexer) advance() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	ch := l.input[l.pos]
	l.pos++
	if ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
	return ch
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) {
		ch := l.peek()
		if ch == ' ' || ch == '\t' || ch == '\r' {
			l.advance()
		} else {
			break
		}
	}
}

func (l *Lexer) makeToken(typ TokenType, value string) Token {
	return Token{
		Type:   typ,
		Value:  value,
		Line:   l.startLine,
		Column: l.startColumn,
		Offset: l.start,
		Length: l.pos - l.start,
	}
}

func (l *Lexer) scanComment() Token {
	l.advance() // consume ';'
	start := l.pos
	for l.pos < len(l.input) && l.peek() != '\n' {
		l.advance()
	}
	return l.makeToken(TokenComment, l.input[start:l.pos])
}

func (l *Lexer) scanString() Token {
	l.advance() // consume opening '"'
	var builder strings.Builder

	for l.pos < len(l.input) {
		ch := l.peek()
		if ch == '"' {
			l.advance() // consume closing '"'
			return l.makeToken(TokenString, builder.String())
		}
		switch ch {
		case '\\':
			l.advance() // consume '\'
			if l.pos < len(l.input) {
				escaped := l.advance()
				switch escaped {
				case 'n':
					builder.WriteByte('\n')
				case 't':
					builder.WriteByte('\t')
				case 'r':
					builder.WriteByte('\r')
				case '\\':
					builder.WriteByte('\\')
				case '"':
					builder.WriteByte('"')
				default:
					builder.WriteByte('\\')
					builder.WriteByte(escaped)
				}
			}
		case '\n':
			// Unterminated string
			return l.makeToken(TokenError, "unterminated string")
		default:
			builder.WriteByte(l.advance())
		}
	}

	return l.makeToken(TokenError, "unterminated string")
}

func (l *Lexer) scanNumber() Token {
	// Handle sign
	if l.peek() == '-' || l.peek() == '+' {
		l.advance()
	}

	// Handle special values: inf, nan
	if l.matchKeyword("inf") || l.matchKeyword("nan") {
		return l.makeToken(TokenNumber, l.input[l.start:l.pos])
	}

	// Integer part
	hasDigits := false
	for l.pos < len(l.input) && l.peek() >= '0' && l.peek() <= '9' {
		l.advance()
		hasDigits = true
	}

	// Decimal part
	if l.peek() == '.' && l.peekN(1) >= '0' && l.peekN(1) <= '9' {
		l.advance() // consume '.'
		for l.pos < len(l.input) && l.peek() >= '0' && l.peek() <= '9' {
			l.advance()
		}
	} else if l.peek() == '.' && !hasDigits {
		// Just a dot, not a number
		return l.makeToken(TokenError, "invalid number")
	}

	// Exponent part
	if l.peek() == 'e' || l.peek() == 'E' {
		l.advance()
		if l.peek() == '-' || l.peek() == '+' {
			l.advance()
		}
		for l.pos < len(l.input) && l.peek() >= '0' && l.peek() <= '9' {
			l.advance()
		}
	}

	return l.makeToken(TokenNumber, l.input[l.start:l.pos])
}

func (l *Lexer) matchKeyword(keyword string) bool {
	if l.pos+len(keyword) > len(l.input) {
		return false
	}
	if l.input[l.pos:l.pos+len(keyword)] == keyword {
		// Make sure it's not part of a longer identifier
		if l.pos+len(keyword) < len(l.input) {
			next, _ := utf8.DecodeRuneInString(l.input[l.pos+len(keyword):])
			if isIdentPart(next) {
				return false
			}
		}
		l.pos += len(keyword)
		l.column += len(keyword)
		return true
	}
	return false
}

func (l *Lexer) scanIdentifier() Token {
	for l.pos < len(l.input) {
		r, size := utf8.DecodeRuneInString(l.input[l.pos:])
		if !isIdentPart(r) {
			break
		}
		l.pos += size
		l.column++
	}

	value := l.input[l.start:l.pos]

	// Check for keywords
	switch value {
	case "true", "false":
		return l.makeToken(TokenBool, value)
	case "null":
		return l.makeToken(TokenNull, value)
	case "inf", "nan":
		return l.makeToken(TokenNumber, value)
	}

	return l.makeToken(TokenIdent, value)
}

func isIdentStart(ch byte) bool {
	return ch >= 'a' && ch <= 'z' ||
		ch >= 'A' && ch <= 'Z' ||
		ch == '_' ||
		ch == '@' // TSCN uses @ for some special identifiers
}

func isIdentPart(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '@'
}
