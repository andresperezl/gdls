package gdshader

import (
	"strings"
	"unicode"
)

// Lexer tokenizes GDShader source code.
type Lexer struct {
	input   string
	pos     int  // current position in input
	readPos int  // reading position (after current char)
	ch      byte // current char under examination
	line    int  // current line number (1-indexed)
	column  int  // current column number (1-indexed)
}

// NewLexer creates a new Lexer for the given input.
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// readChar reads the next character and advances the position.
func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0 // EOF
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
	l.column++
	if l.ch == '\n' {
		l.line++
		l.column = 0
	}
}

// peekChar returns the next character without advancing the position.
func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	tok := Token{
		Line:   l.line,
		Column: l.column,
	}

	switch l.ch {
	case 0:
		tok.Type = TokenEOF
		tok.Literal = ""
	case '\n':
		tok.Type = TokenNewline
		tok.Literal = "\n"
		l.readChar()
	case '(':
		tok.Type = TokenLParen
		tok.Literal = "("
		l.readChar()
	case ')':
		tok.Type = TokenRParen
		tok.Literal = ")"
		l.readChar()
	case '{':
		tok.Type = TokenLBrace
		tok.Literal = "{"
		l.readChar()
	case '}':
		tok.Type = TokenRBrace
		tok.Literal = "}"
		l.readChar()
	case '[':
		tok.Type = TokenLBracket
		tok.Literal = "["
		l.readChar()
	case ']':
		tok.Type = TokenRBracket
		tok.Literal = "]"
		l.readChar()
	case ';':
		tok.Type = TokenSemicolon
		tok.Literal = ";"
		l.readChar()
	case ',':
		tok.Type = TokenComma
		tok.Literal = ","
		l.readChar()
	case '.':
		// Check if this is the start of a float literal like .5
		if isDigit(l.peekChar()) {
			tok = l.readNumber()
		} else {
			tok.Type = TokenDot
			tok.Literal = "."
			l.readChar()
		}
	case ':':
		tok.Type = TokenColon
		tok.Literal = ":"
		l.readChar()
	case '?':
		tok.Type = TokenQuestion
		tok.Literal = "?"
		l.readChar()
	case '~':
		tok.Type = TokenTilde
		tok.Literal = "~"
		l.readChar()
	case '+':
		if l.peekChar() == '+' {
			l.readChar()
			tok.Type = TokenIncrement
			tok.Literal = "++"
		} else if l.peekChar() == '=' {
			l.readChar()
			tok.Type = TokenPlusAssign
			tok.Literal = "+="
		} else {
			tok.Type = TokenPlus
			tok.Literal = "+"
		}
		l.readChar()
	case '-':
		if l.peekChar() == '-' {
			l.readChar()
			tok.Type = TokenDecrement
			tok.Literal = "--"
		} else if l.peekChar() == '=' {
			l.readChar()
			tok.Type = TokenMinusAssign
			tok.Literal = "-="
		} else {
			tok.Type = TokenMinus
			tok.Literal = "-"
		}
		l.readChar()
	case '*':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = TokenStarAssign
			tok.Literal = "*="
		} else {
			tok.Type = TokenStar
			tok.Literal = "*"
		}
		l.readChar()
	case '/':
		if l.peekChar() == '/' {
			tok = l.readLineComment()
		} else if l.peekChar() == '*' {
			tok = l.readBlockComment()
		} else if l.peekChar() == '=' {
			l.readChar()
			tok.Type = TokenSlashAssign
			tok.Literal = "/="
			l.readChar()
		} else {
			tok.Type = TokenSlash
			tok.Literal = "/"
			l.readChar()
		}
	case '%':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = TokenPercentAssign
			tok.Literal = "%="
		} else {
			tok.Type = TokenPercent
			tok.Literal = "%"
		}
		l.readChar()
	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			tok.Type = TokenAnd
			tok.Literal = "&&"
		} else if l.peekChar() == '=' {
			l.readChar()
			tok.Type = TokenAmpAssign
			tok.Literal = "&="
		} else {
			tok.Type = TokenAmpersand
			tok.Literal = "&"
		}
		l.readChar()
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			tok.Type = TokenOr
			tok.Literal = "||"
		} else if l.peekChar() == '=' {
			l.readChar()
			tok.Type = TokenPipeAssign
			tok.Literal = "|="
		} else {
			tok.Type = TokenPipe
			tok.Literal = "|"
		}
		l.readChar()
	case '^':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = TokenCaretAssign
			tok.Literal = "^="
		} else {
			tok.Type = TokenCaret
			tok.Literal = "^"
		}
		l.readChar()
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = TokenNE
			tok.Literal = "!="
		} else {
			tok.Type = TokenBang
			tok.Literal = "!"
		}
		l.readChar()
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = TokenEQ
			tok.Literal = "=="
		} else {
			tok.Type = TokenAssign
			tok.Literal = "="
		}
		l.readChar()
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = TokenLTE
			tok.Literal = "<="
			l.readChar()
		} else if l.peekChar() == '<' {
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok.Type = TokenLeftShiftAssign
				tok.Literal = "<<="
			} else {
				tok.Type = TokenLeftShift
				tok.Literal = "<<"
			}
			l.readChar()
		} else {
			tok.Type = TokenLT
			tok.Literal = "<"
			l.readChar()
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = TokenGTE
			tok.Literal = ">="
			l.readChar()
		} else if l.peekChar() == '>' {
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok.Type = TokenRightShiftAssign
				tok.Literal = ">>="
			} else {
				tok.Type = TokenRightShift
				tok.Literal = ">>"
			}
			l.readChar()
		} else {
			tok.Type = TokenGT
			tok.Literal = ">"
			l.readChar()
		}
	default:
		if isLetter(l.ch) || l.ch == '_' {
			tok = l.readIdentifier()
		} else if isDigit(l.ch) {
			tok = l.readNumber()
		} else {
			tok.Type = TokenError
			tok.Literal = string(l.ch)
			l.readChar()
		}
	}

	return tok
}

// skipWhitespace skips spaces and tabs (but not newlines).
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

// readIdentifier reads an identifier or keyword.
func (l *Lexer) readIdentifier() Token {
	tok := Token{
		Line:   l.line,
		Column: l.column,
	}
	startPos := l.pos
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	tok.Literal = l.input[startPos:l.pos]
	tok.Type = LookupIdent(tok.Literal)
	return tok
}

// readNumber reads an integer or float literal.
func (l *Lexer) readNumber() Token {
	tok := Token{
		Line:   l.line,
		Column: l.column,
	}
	startPos := l.pos
	isFloat := false

	// Check for hex literal
	if l.ch == '0' && (l.peekChar() == 'x' || l.peekChar() == 'X') {
		l.readChar() // consume '0'
		l.readChar() // consume 'x'
		for isHexDigit(l.ch) {
			l.readChar()
		}
		tok.Literal = l.input[startPos:l.pos]
		tok.Type = TokenIntLit
		return tok
	}

	// Read integer part (or start from '.' for .5 style floats)
	if l.ch == '.' {
		isFloat = true
	} else {
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	// Check for decimal point
	if l.ch == '.' && isDigit(l.peekChar()) {
		isFloat = true
		l.readChar() // consume '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	// Check for exponent
	if l.ch == 'e' || l.ch == 'E' {
		isFloat = true
		l.readChar() // consume 'e'
		if l.ch == '+' || l.ch == '-' {
			l.readChar()
		}
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	// Check for float suffix 'f' or 'F'
	if l.ch == 'f' || l.ch == 'F' {
		isFloat = true
		l.readChar()
	}

	// Check for unsigned suffix 'u' or 'U'
	if l.ch == 'u' || l.ch == 'U' {
		l.readChar()
	}

	tok.Literal = l.input[startPos:l.pos]
	if isFloat {
		tok.Type = TokenFloatLit
	} else {
		tok.Type = TokenIntLit
	}
	return tok
}

// readLineComment reads a // comment.
func (l *Lexer) readLineComment() Token {
	tok := Token{
		Type:   TokenLineComment,
		Line:   l.line,
		Column: l.column,
	}
	startPos := l.pos
	// Skip //
	l.readChar()
	l.readChar()
	// Read until end of line or EOF
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	tok.Literal = l.input[startPos:l.pos]
	return tok
}

// readBlockComment reads a /* */ or /** */ comment.
func (l *Lexer) readBlockComment() Token {
	tok := Token{
		Line:   l.line,
		Column: l.column,
	}
	startPos := l.pos
	// Skip /*
	l.readChar()
	l.readChar()

	// Check if it's a doc comment /** */
	isDoc := l.ch == '*' && l.peekChar() != '/'
	if isDoc {
		tok.Type = TokenDocComment
	} else {
		tok.Type = TokenBlockComment
	}

	// Read until */
	for {
		if l.ch == 0 {
			// Unterminated comment
			tok.Type = TokenError
			break
		}
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar() // consume '*'
			l.readChar() // consume '/'
			break
		}
		l.readChar()
	}
	tok.Literal = l.input[startPos:l.pos]
	return tok
}

// Tokenize tokenizes the entire input and returns all tokens.
func (l *Lexer) Tokenize() []Token {
	var tokens []Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == TokenEOF {
			break
		}
	}
	return tokens
}

// TokenizeSkipComments tokenizes the input, skipping comment tokens.
func (l *Lexer) TokenizeSkipComments() []Token {
	var tokens []Token
	for {
		tok := l.NextToken()
		switch tok.Type {
		case TokenLineComment, TokenBlockComment, TokenDocComment:
			continue
		case TokenNewline:
			continue
		default:
			tokens = append(tokens, tok)
		}
		if tok.Type == TokenEOF {
			break
		}
	}
	return tokens
}

// Helper functions

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch))
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isHexDigit(ch byte) bool {
	return isDigit(ch) || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

// ExtractDocComment extracts the text content from a doc comment.
func ExtractDocComment(comment string) string {
	// Remove /** and */
	if len(comment) < 5 {
		return ""
	}
	content := comment[3 : len(comment)-2]

	// Split into lines and clean up
	lines := strings.Split(content, "\n")
	var cleaned []string
	for _, line := range lines {
		// Trim leading whitespace and * characters
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "*")
		line = strings.TrimSpace(line)
		if line != "" {
			cleaned = append(cleaned, line)
		}
	}
	return strings.Join(cleaned, "\n")
}
