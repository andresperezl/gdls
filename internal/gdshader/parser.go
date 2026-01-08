package gdshader

import (
	"fmt"
)

// Parser parses GDShader source code into an AST.
type Parser struct {
	lexer    *Lexer
	tokens   []Token
	pos      int
	errors   []ParseError
	comments []*Comment
	lastDoc  string // last doc comment for uniform documentation
}

// Parse parses the input and returns a ShaderDocument.
func Parse(input string) *ShaderDocument {
	p := NewParser(input)
	return p.Parse()
}

// NewParser creates a new Parser for the given input.
func NewParser(input string) *Parser {
	lexer := NewLexer(input)
	tokens := lexer.Tokenize()
	return &Parser{
		lexer:  lexer,
		tokens: tokens,
		pos:    0,
	}
}

// Parse parses the input and returns a ShaderDocument.
func (p *Parser) Parse() *ShaderDocument {
	doc := &ShaderDocument{}

	// Skip initial whitespace/comments
	p.skipNewlinesAndComments()

	// Parse shader_type (required first)
	if p.check(TokenShaderType) {
		doc.ShaderType = p.parseShaderType()
		p.skipNewlinesAndComments()
	}

	// Parse render_mode (optional)
	if p.check(TokenRenderMode) {
		doc.RenderModes = p.parseRenderMode()
		p.skipNewlinesAndComments()
	}

	// Parse declarations
	for !p.isAtEnd() {
		p.skipNewlinesAndComments()
		if p.isAtEnd() {
			break
		}

		decl := p.parseDeclaration()
		if decl == nil {
			// Skip to next semicolon or newline on error
			p.synchronize()
			continue
		}

		switch d := decl.(type) {
		case *StructDecl:
			doc.Structs = append(doc.Structs, d)
		case *UniformDecl:
			doc.Uniforms = append(doc.Uniforms, d)
		case *VaryingDecl:
			doc.Varyings = append(doc.Varyings, d)
		case *ConstDecl:
			doc.Constants = append(doc.Constants, d)
		case *FunctionDecl:
			doc.Functions = append(doc.Functions, d)
		}
	}

	doc.Errors = p.errors
	doc.Comments = p.comments
	return doc
}

// current returns the current token.
func (p *Parser) current() Token {
	if p.pos >= len(p.tokens) {
		return Token{Type: TokenEOF}
	}
	return p.tokens[p.pos]
}

// peek returns the next token without consuming.
func (p *Parser) peek() Token {
	if p.pos+1 >= len(p.tokens) {
		return Token{Type: TokenEOF}
	}
	return p.tokens[p.pos+1]
}

// advance consumes the current token and returns it.
func (p *Parser) advance() Token {
	tok := p.current()
	if p.pos < len(p.tokens) {
		p.pos++
	}
	return tok
}

// check returns true if the current token matches the given type.
func (p *Parser) check(t TokenType) bool {
	return p.current().Type == t
}

// match consumes the current token if it matches any of the given types.
func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

// expect consumes the current token if it matches, otherwise adds an error.
func (p *Parser) expect(t TokenType, msg string) bool {
	if p.check(t) {
		p.advance()
		return true
	}
	p.error(msg)
	return false
}

// isAtEnd returns true if we've reached the end of tokens.
func (p *Parser) isAtEnd() bool {
	return p.current().Type == TokenEOF
}

// error adds a parse error at the current position.
func (p *Parser) error(msg string) {
	tok := p.current()
	p.errors = append(p.errors, ParseError{
		Range: Range{
			Start: Position{Line: tok.Line - 1, Column: tok.Column - 1},
			End:   Position{Line: tok.Line - 1, Column: tok.Column - 1 + len(tok.Literal)},
		},
		Message: msg,
	})
}

// synchronize skips tokens until we find a synchronization point.
func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		// Synchronize on semicolons or declaration keywords
		switch p.current().Type {
		case TokenSemicolon:
			p.advance()
			return
		case TokenStruct, TokenUniform, TokenVarying, TokenConst, TokenVoid,
			TokenBool, TokenInt, TokenUint, TokenFloat, TokenVec2, TokenVec3,
			TokenVec4, TokenMat2, TokenMat3, TokenMat4:
			return
		}
		p.advance()
	}
}

// skipNewlinesAndComments skips newlines and stores comments.
func (p *Parser) skipNewlinesAndComments() {
	for {
		switch p.current().Type {
		case TokenNewline:
			p.advance()
		case TokenLineComment:
			tok := p.advance()
			p.comments = append(p.comments, &Comment{
				Range: p.tokenRange(tok),
				Text:  tok.Literal,
				IsDoc: false,
			})
		case TokenBlockComment:
			tok := p.advance()
			p.comments = append(p.comments, &Comment{
				Range: p.tokenRange(tok),
				Text:  tok.Literal,
				IsDoc: false,
			})
		case TokenDocComment:
			tok := p.advance()
			p.comments = append(p.comments, &Comment{
				Range: p.tokenRange(tok),
				Text:  tok.Literal,
				IsDoc: true,
			})
			p.lastDoc = ExtractDocComment(tok.Literal)
		default:
			return
		}
	}
}

// tokenRange creates a Range from a token.
func (p *Parser) tokenRange(tok Token) Range {
	return Range{
		Start: Position{Line: tok.Line - 1, Column: tok.Column - 1},
		End:   Position{Line: tok.Line - 1, Column: tok.Column - 1 + len(tok.Literal)},
	}
}

// parseShaderType parses "shader_type <type>;"
func (p *Parser) parseShaderType() *ShaderTypeDecl {
	start := p.current()
	p.advance() // consume shader_type

	decl := &ShaderTypeDecl{
		Range: p.tokenRange(start),
	}

	if p.check(TokenIdent) {
		decl.Type = p.current().Literal
		p.advance()
	} else {
		p.error("expected shader type (spatial, canvas_item, particles, sky, or fog)")
		return decl
	}

	p.expect(TokenSemicolon, "expected ';' after shader_type")
	return decl
}

// parseRenderMode parses "render_mode mode1, mode2, ...;"
func (p *Parser) parseRenderMode() *RenderModeDecl {
	start := p.current()
	p.advance() // consume render_mode

	decl := &RenderModeDecl{
		Range: p.tokenRange(start),
	}

	for {
		if p.check(TokenIdent) {
			decl.Modes = append(decl.Modes, p.current().Literal)
			p.advance()
		} else {
			p.error("expected render mode identifier")
			break
		}

		if !p.match(TokenComma) {
			break
		}
	}

	p.expect(TokenSemicolon, "expected ';' after render_mode")
	return decl
}

// parseDeclaration parses a top-level declaration.
func (p *Parser) parseDeclaration() interface{} {
	// Check for group_uniforms
	if p.check(TokenGroupUniforms) {
		// Just skip group_uniforms for now, but we could track it
		p.advance()
		if p.check(TokenIdent) {
			p.advance() // group name
			if p.check(TokenDot) {
				p.advance()
				if p.check(TokenIdent) {
					p.advance() // subgroup name
				}
			}
		}
		p.expect(TokenSemicolon, "expected ';' after group_uniforms")
		return nil
	}

	// Check for struct
	if p.check(TokenStruct) {
		return p.parseStructDecl()
	}

	// Check for global uniform
	if p.check(TokenGlobal) {
		p.advance()
		if p.check(TokenUniform) {
			return p.parseUniformDecl(true)
		}
		p.error("expected 'uniform' after 'global'")
		return nil
	}

	// Check for uniform
	if p.check(TokenUniform) {
		return p.parseUniformDecl(false)
	}

	// Check for varying (with optional interpolation qualifier)
	if p.check(TokenFlat) || p.check(TokenSmooth) {
		interp := p.current().Literal
		p.advance()
		if p.check(TokenVarying) {
			return p.parseVaryingDecl(interp)
		}
		p.error("expected 'varying' after interpolation qualifier")
		return nil
	}

	if p.check(TokenVarying) {
		return p.parseVaryingDecl("")
	}

	// Check for const
	if p.check(TokenConst) {
		return p.parseConstDecl()
	}

	// Otherwise, it should be a function or global variable
	return p.parseFunctionOrVar()
}

// parseStructDecl parses a struct declaration.
func (p *Parser) parseStructDecl() *StructDecl {
	start := p.current()
	p.advance() // consume struct

	decl := &StructDecl{
		Range: p.tokenRange(start),
	}

	if p.check(TokenIdent) {
		decl.Name = p.current().Literal
		p.advance()
	} else {
		p.error("expected struct name")
		return decl
	}

	if !p.expect(TokenLBrace, "expected '{' after struct name") {
		return decl
	}

	for !p.check(TokenRBrace) && !p.isAtEnd() {
		p.skipNewlinesAndComments()
		if p.check(TokenRBrace) {
			break
		}

		member := p.parseStructMember()
		if member != nil {
			decl.Members = append(decl.Members, member)
		}
		p.skipNewlinesAndComments()
	}

	p.expect(TokenRBrace, "expected '}' after struct members")
	p.expect(TokenSemicolon, "expected ';' after struct declaration")

	return decl
}

// parseStructMember parses a struct member.
func (p *Parser) parseStructMember() *StructMember {
	typeSpec := p.parseTypeSpec()
	if typeSpec == nil {
		return nil
	}

	member := &StructMember{
		Range: typeSpec.Range,
		Type:  typeSpec,
	}

	if p.check(TokenIdent) {
		member.Name = p.current().Literal
		p.advance()
	} else {
		p.error("expected member name")
		return nil
	}

	p.expect(TokenSemicolon, "expected ';' after struct member")
	return member
}

// parseUniformDecl parses a uniform declaration.
func (p *Parser) parseUniformDecl(isGlobal bool) *UniformDecl {
	start := p.current()
	p.advance() // consume uniform

	docComment := p.lastDoc
	p.lastDoc = ""

	decl := &UniformDecl{
		Range:      p.tokenRange(start),
		IsGlobal:   isGlobal,
		DocComment: docComment,
	}

	typeSpec := p.parseTypeSpec()
	if typeSpec == nil {
		return decl
	}
	decl.Type = typeSpec

	if p.check(TokenIdent) {
		decl.Name = p.current().Literal
		p.advance()
	} else {
		p.error("expected uniform name")
		return decl
	}

	// Parse array size if present
	if p.check(TokenLBracket) {
		p.advance()
		if !p.check(TokenRBracket) {
			decl.Type.ArraySize = p.parseExpression()
		}
		p.expect(TokenRBracket, "expected ']' after array size")
	}

	// Parse hints
	if p.check(TokenColon) {
		p.advance()
		decl.Hints = p.parseHints()
	}

	// Parse default value
	if p.check(TokenAssign) {
		p.advance()
		decl.DefaultValue = p.parseExpression()
	}

	p.expect(TokenSemicolon, "expected ';' after uniform declaration")
	return decl
}

// parseHints parses uniform hints.
func (p *Parser) parseHints() []*Hint {
	var hints []*Hint

	for p.check(TokenIdent) {

		hint := &Hint{
			Range: p.tokenRange(p.current()),
			Name:  p.current().Literal,
		}
		p.advance()

		// Parse hint arguments if present
		if p.check(TokenLParen) {
			p.advance()
			for !p.check(TokenRParen) && !p.isAtEnd() {
				arg := p.parseExpression()
				if arg != nil {
					hint.Args = append(hint.Args, arg)
				}
				if !p.match(TokenComma) {
					break
				}
			}
			p.expect(TokenRParen, "expected ')' after hint arguments")
		}

		hints = append(hints, hint)

		// Check for more hints or default value
		if !p.match(TokenComma) {
			break
		}
		// If we see '=' after comma, it's a default value, not another hint
		if p.check(TokenAssign) {
			break
		}
	}

	return hints
}

// parseVaryingDecl parses a varying declaration.
func (p *Parser) parseVaryingDecl(interpolation string) *VaryingDecl {
	start := p.current()
	p.advance() // consume varying

	decl := &VaryingDecl{
		Range:         p.tokenRange(start),
		Interpolation: interpolation,
	}

	typeSpec := p.parseTypeSpec()
	if typeSpec == nil {
		return decl
	}
	decl.Type = typeSpec

	if p.check(TokenIdent) {
		decl.Name = p.current().Literal
		p.advance()
	} else {
		p.error("expected varying name")
		return decl
	}

	// Parse array size if present
	if p.check(TokenLBracket) {
		p.advance()
		if !p.check(TokenRBracket) {
			decl.Type.ArraySize = p.parseExpression()
		}
		p.expect(TokenRBracket, "expected ']' after array size")
	}

	p.expect(TokenSemicolon, "expected ';' after varying declaration")
	return decl
}

// parseConstDecl parses a constant declaration.
func (p *Parser) parseConstDecl() *ConstDecl {
	start := p.current()
	p.advance() // consume const

	decl := &ConstDecl{
		Range: p.tokenRange(start),
	}

	typeSpec := p.parseTypeSpec()
	if typeSpec == nil {
		return decl
	}
	decl.Type = typeSpec

	if p.check(TokenIdent) {
		decl.Name = p.current().Literal
		p.advance()
	} else {
		p.error("expected constant name")
		return decl
	}

	if !p.expect(TokenAssign, "expected '=' in constant declaration") {
		return decl
	}

	decl.Value = p.parseExpression()
	p.expect(TokenSemicolon, "expected ';' after constant declaration")
	return decl
}

// parseFunctionOrVar parses a function declaration or global variable.
func (p *Parser) parseFunctionOrVar() interface{} {
	typeSpec := p.parseTypeSpec()
	if typeSpec == nil {
		return nil
	}

	if !p.check(TokenIdent) {
		p.error("expected identifier after type")
		return nil
	}

	name := p.current().Literal
	nameRange := p.tokenRange(p.current())
	p.advance()

	// Check if this is a function
	if p.check(TokenLParen) {
		return p.parseFunctionDecl(typeSpec, name, nameRange)
	}

	// Otherwise it's a global variable (not typically used in shaders, but valid)
	p.error("global variables are not supported; use uniform or const")
	return nil
}

// parseFunctionDecl parses a function declaration.
func (p *Parser) parseFunctionDecl(returnType *TypeSpec, name string, nameRange Range) *FunctionDecl {
	decl := &FunctionDecl{
		Range: Range{
			Start: returnType.Range.Start,
			End:   nameRange.End,
		},
		ReturnType: returnType,
		Name:       name,
	}

	p.advance() // consume '('

	// Parse parameters
	for !p.check(TokenRParen) && !p.isAtEnd() {
		param := p.parseParamDecl()
		if param != nil {
			decl.Params = append(decl.Params, param)
		}
		if !p.match(TokenComma) {
			break
		}
	}

	p.expect(TokenRParen, "expected ')' after function parameters")

	// Parse body
	p.skipNewlinesAndComments()
	if p.check(TokenLBrace) {
		decl.Body = p.parseBlockStmt()
	} else {
		p.error("expected '{' for function body")
	}

	return decl
}

// parseParamDecl parses a function parameter.
func (p *Parser) parseParamDecl() *ParamDecl {
	param := &ParamDecl{
		Range: p.tokenRange(p.current()),
	}

	// Parse qualifier (in, out, inout, const)
	switch p.current().Type {
	case TokenIn, TokenOut, TokenInout, TokenConst:
		param.Qualifier = p.current().Literal
		p.advance()
	}

	typeSpec := p.parseTypeSpec()
	if typeSpec == nil {
		return nil
	}
	param.Type = typeSpec

	if p.check(TokenIdent) {
		param.Name = p.current().Literal
		p.advance()
	} else {
		p.error("expected parameter name")
		return nil
	}

	// Parse array size if present
	if p.check(TokenLBracket) {
		p.advance()
		if !p.check(TokenRBracket) {
			param.Type.ArraySize = p.parseExpression()
		}
		p.expect(TokenRBracket, "expected ']' after array size")
	}

	return param
}

// parseTypeSpec parses a type specification.
func (p *Parser) parseTypeSpec() *TypeSpec {
	spec := &TypeSpec{
		Range: p.tokenRange(p.current()),
	}

	// Parse precision qualifier
	if p.current().Type.IsPrecision() {
		spec.Precision = p.current().Literal
		p.advance()
	}

	// Parse type name
	if p.current().Type.IsType() {
		spec.Name = p.current().Literal
		p.advance()
	} else if p.check(TokenIdent) {
		// Custom type (struct)
		spec.Name = p.current().Literal
		p.advance()
	} else {
		p.error("expected type name")
		return nil
	}

	return spec
}

// parseBlockStmt parses a block statement.
func (p *Parser) parseBlockStmt() *BlockStmt {
	start := p.current()
	p.advance() // consume '{'

	block := &BlockStmt{
		Range: p.tokenRange(start),
	}

	for !p.check(TokenRBrace) && !p.isAtEnd() {
		p.skipNewlinesAndComments()
		if p.check(TokenRBrace) {
			break
		}

		stmt := p.parseStatement()
		block.Stmts = append(block.Stmts, stmt)
		p.skipNewlinesAndComments()
	}

	p.expect(TokenRBrace, "expected '}' after block")
	return block
}

// parseStatement parses a statement.
func (p *Parser) parseStatement() Stmt {
	switch p.current().Type {
	case TokenLBrace:
		return p.parseBlockStmt()
	case TokenIf:
		return p.parseIfStmt()
	case TokenFor:
		return p.parseForStmt()
	case TokenWhile:
		return p.parseWhileStmt()
	case TokenDo:
		return p.parseDoWhileStmt()
	case TokenSwitch:
		return p.parseSwitchStmt()
	case TokenReturn:
		return p.parseReturnStmt()
	case TokenBreak:
		return p.parseBreakStmt()
	case TokenContinue:
		return p.parseContinueStmt()
	case TokenDiscard:
		return p.parseDiscardStmt()
	case TokenSemicolon:
		tok := p.advance()
		return &EmptyStmt{Range: p.tokenRange(tok)}
	case TokenConst:
		return p.parseVarDeclStmt(true)
	default:
		// Could be a variable declaration or expression statement
		if p.isTypeStart() {
			return p.parseVarDeclStmt(false)
		}
		return p.parseExprStmt()
	}
}

// isTypeStart returns true if the current token could start a type.
func (p *Parser) isTypeStart() bool {
	if p.current().Type.IsType() || p.current().Type.IsPrecision() {
		return true
	}
	// Check for custom type followed by identifier
	if p.check(TokenIdent) && p.peek().Type == TokenIdent {
		return true
	}
	return false
}

// parseVarDeclStmt parses a variable declaration statement.
func (p *Parser) parseVarDeclStmt(isConst bool) *VarDeclStmt {
	start := p.current()

	stmt := &VarDeclStmt{
		Range: p.tokenRange(start),
		Const: isConst,
	}

	if isConst {
		p.advance() // consume const
	}

	stmt.Type = p.parseTypeSpec()
	if stmt.Type == nil {
		return nil
	}

	// Parse variable declarators
	for {
		decl := &VarDecl{
			Range: p.tokenRange(p.current()),
		}

		if p.check(TokenIdent) {
			decl.Name = p.current().Literal
			p.advance()
		} else {
			p.error("expected variable name")
			break
		}

		// Parse array size if present
		if p.check(TokenLBracket) {
			p.advance()
			if !p.check(TokenRBracket) {
				decl.ArraySize = p.parseExpression()
			}
			p.expect(TokenRBracket, "expected ']' after array size")
		}

		// Parse initializer
		if p.check(TokenAssign) {
			p.advance()
			decl.Init = p.parseExpression()
		}

		stmt.Decls = append(stmt.Decls, decl)

		if !p.match(TokenComma) {
			break
		}
	}

	p.expect(TokenSemicolon, "expected ';' after variable declaration")
	return stmt
}

// parseIfStmt parses an if statement.
func (p *Parser) parseIfStmt() *IfStmt {
	start := p.current()
	p.advance() // consume if

	stmt := &IfStmt{
		Range: p.tokenRange(start),
	}

	p.expect(TokenLParen, "expected '(' after 'if'")
	stmt.Cond = p.parseExpression()
	p.expect(TokenRParen, "expected ')' after if condition")

	p.skipNewlinesAndComments()
	stmt.Then = p.parseStatement()

	p.skipNewlinesAndComments()
	if p.match(TokenElse) {
		p.skipNewlinesAndComments()
		stmt.Else = p.parseStatement()
	}

	return stmt
}

// parseForStmt parses a for statement.
func (p *Parser) parseForStmt() *ForStmt {
	start := p.current()
	p.advance() // consume for

	stmt := &ForStmt{
		Range: p.tokenRange(start),
	}

	p.expect(TokenLParen, "expected '(' after 'for'")

	// Parse init
	if !p.check(TokenSemicolon) {
		if p.isTypeStart() {
			stmt.Init = p.parseVarDeclStmt(false)
		} else {
			stmt.Init = p.parseExprStmt()
		}
	} else {
		p.advance() // consume ';'
	}

	// Parse condition
	if !p.check(TokenSemicolon) {
		stmt.Cond = p.parseExpression()
	}
	p.expect(TokenSemicolon, "expected ';' after for condition")

	// Parse post expression
	if !p.check(TokenRParen) {
		stmt.Post = p.parseExpression()
	}
	p.expect(TokenRParen, "expected ')' after for clauses")

	p.skipNewlinesAndComments()
	stmt.Body = p.parseStatement()

	return stmt
}

// parseWhileStmt parses a while statement.
func (p *Parser) parseWhileStmt() *WhileStmt {
	start := p.current()
	p.advance() // consume while

	stmt := &WhileStmt{
		Range: p.tokenRange(start),
	}

	p.expect(TokenLParen, "expected '(' after 'while'")
	stmt.Cond = p.parseExpression()
	p.expect(TokenRParen, "expected ')' after while condition")

	p.skipNewlinesAndComments()
	stmt.Body = p.parseStatement()

	return stmt
}

// parseDoWhileStmt parses a do-while statement.
func (p *Parser) parseDoWhileStmt() *DoWhileStmt {
	start := p.current()
	p.advance() // consume do

	stmt := &DoWhileStmt{
		Range: p.tokenRange(start),
	}

	p.skipNewlinesAndComments()
	stmt.Body = p.parseStatement()

	p.skipNewlinesAndComments()
	p.expect(TokenWhile, "expected 'while' after do body")
	p.expect(TokenLParen, "expected '(' after 'while'")
	stmt.Cond = p.parseExpression()
	p.expect(TokenRParen, "expected ')' after while condition")
	p.expect(TokenSemicolon, "expected ';' after do-while")

	return stmt
}

// parseSwitchStmt parses a switch statement.
func (p *Parser) parseSwitchStmt() *SwitchStmt {
	start := p.current()
	p.advance() // consume switch

	stmt := &SwitchStmt{
		Range: p.tokenRange(start),
	}

	p.expect(TokenLParen, "expected '(' after 'switch'")
	stmt.Expr = p.parseExpression()
	p.expect(TokenRParen, "expected ')' after switch expression")

	p.skipNewlinesAndComments()
	p.expect(TokenLBrace, "expected '{' after switch")

	for !p.check(TokenRBrace) && !p.isAtEnd() {
		p.skipNewlinesAndComments()
		if p.check(TokenRBrace) {
			break
		}

		clause := p.parseCaseClause()
		if clause != nil {
			stmt.Cases = append(stmt.Cases, clause)
		}
	}

	p.expect(TokenRBrace, "expected '}' after switch cases")
	return stmt
}

// parseCaseClause parses a case or default clause.
func (p *Parser) parseCaseClause() *CaseClause {
	start := p.current()
	clause := &CaseClause{
		Range: p.tokenRange(start),
	}

	if p.match(TokenCase) {
		// Parse case value(s)
		clause.Values = append(clause.Values, p.parseExpression())
	} else if p.match(TokenDefault) {
		// Default has no values
	} else {
		p.error("expected 'case' or 'default'")
		return nil
	}

	p.expect(TokenColon, "expected ':' after case/default")

	// Parse statements until next case/default or end of switch
	for !p.check(TokenCase) && !p.check(TokenDefault) && !p.check(TokenRBrace) && !p.isAtEnd() {
		p.skipNewlinesAndComments()
		if p.check(TokenCase) || p.check(TokenDefault) || p.check(TokenRBrace) {
			break
		}

		stmt := p.parseStatement()
		clause.Body = append(clause.Body, stmt)
		p.skipNewlinesAndComments()
	}

	return clause
}

// parseReturnStmt parses a return statement.
func (p *Parser) parseReturnStmt() *ReturnStmt {
	start := p.current()
	p.advance() // consume return

	stmt := &ReturnStmt{
		Range: p.tokenRange(start),
	}

	if !p.check(TokenSemicolon) {
		stmt.Value = p.parseExpression()
	}

	p.expect(TokenSemicolon, "expected ';' after return")
	return stmt
}

// parseBreakStmt parses a break statement.
func (p *Parser) parseBreakStmt() *BreakStmt {
	start := p.current()
	p.advance() // consume break
	p.expect(TokenSemicolon, "expected ';' after break")
	return &BreakStmt{Range: p.tokenRange(start)}
}

// parseContinueStmt parses a continue statement.
func (p *Parser) parseContinueStmt() *ContinueStmt {
	start := p.current()
	p.advance() // consume continue
	p.expect(TokenSemicolon, "expected ';' after continue")
	return &ContinueStmt{Range: p.tokenRange(start)}
}

// parseDiscardStmt parses a discard statement.
func (p *Parser) parseDiscardStmt() *DiscardStmt {
	start := p.current()
	p.advance() // consume discard
	p.expect(TokenSemicolon, "expected ';' after discard")
	return &DiscardStmt{Range: p.tokenRange(start)}
}

// parseExprStmt parses an expression statement.
func (p *Parser) parseExprStmt() *ExprStmt {
	start := p.current()
	stmt := &ExprStmt{
		Range: p.tokenRange(start),
		Expr:  p.parseExpression(),
	}
	p.expect(TokenSemicolon, "expected ';' after expression")
	return stmt
}

// parseExpression parses an expression.
func (p *Parser) parseExpression() Expr {
	return p.parseTernary()
}

// parseTernary parses a ternary conditional expression.
func (p *Parser) parseTernary() Expr {
	expr := p.parseOr()

	if p.match(TokenQuestion) {
		thenExpr := p.parseExpression()
		p.expect(TokenColon, "expected ':' in ternary expression")
		elseExpr := p.parseTernary()
		expr = &TernaryExpr{
			Range: Range{
				Start: expr.GetRange().Start,
				End:   elseExpr.GetRange().End,
			},
			Cond: expr,
			Then: thenExpr,
			Else: elseExpr,
		}
	}

	return expr
}

// parseOr parses logical OR expressions.
func (p *Parser) parseOr() Expr {
	expr := p.parseAnd()

	for p.match(TokenOr) {
		op := "||"
		right := p.parseAnd()
		expr = &BinaryExpr{
			Range: Range{
				Start: expr.GetRange().Start,
				End:   right.GetRange().End,
			},
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}

	return expr
}

// parseAnd parses logical AND expressions.
func (p *Parser) parseAnd() Expr {
	expr := p.parseBitwiseOr()

	for p.match(TokenAnd) {
		op := "&&"
		right := p.parseBitwiseOr()
		expr = &BinaryExpr{
			Range: Range{
				Start: expr.GetRange().Start,
				End:   right.GetRange().End,
			},
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}

	return expr
}

// parseBitwiseOr parses bitwise OR expressions.
func (p *Parser) parseBitwiseOr() Expr {
	expr := p.parseBitwiseXor()

	for p.match(TokenPipe) {
		op := "|"
		right := p.parseBitwiseXor()
		expr = &BinaryExpr{
			Range: Range{
				Start: expr.GetRange().Start,
				End:   right.GetRange().End,
			},
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}

	return expr
}

// parseBitwiseXor parses bitwise XOR expressions.
func (p *Parser) parseBitwiseXor() Expr {
	expr := p.parseBitwiseAnd()

	for p.match(TokenCaret) {
		op := "^"
		right := p.parseBitwiseAnd()
		expr = &BinaryExpr{
			Range: Range{
				Start: expr.GetRange().Start,
				End:   right.GetRange().End,
			},
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}

	return expr
}

// parseBitwiseAnd parses bitwise AND expressions.
func (p *Parser) parseBitwiseAnd() Expr {
	expr := p.parseEquality()

	for p.match(TokenAmpersand) {
		op := "&"
		right := p.parseEquality()
		expr = &BinaryExpr{
			Range: Range{
				Start: expr.GetRange().Start,
				End:   right.GetRange().End,
			},
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}

	return expr
}

// parseEquality parses equality expressions.
func (p *Parser) parseEquality() Expr {
	expr := p.parseRelational()

	for {
		var op string
		if p.match(TokenEQ) {
			op = "=="
		} else if p.match(TokenNE) {
			op = "!="
		} else {
			break
		}

		right := p.parseRelational()
		expr = &BinaryExpr{
			Range: Range{
				Start: expr.GetRange().Start,
				End:   right.GetRange().End,
			},
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}

	return expr
}

// parseRelational parses relational expressions.
func (p *Parser) parseRelational() Expr {
	expr := p.parseShift()

	for {
		var op string
		if p.match(TokenLT) {
			op = "<"
		} else if p.match(TokenGT) {
			op = ">"
		} else if p.match(TokenLTE) {
			op = "<="
		} else if p.match(TokenGTE) {
			op = ">="
		} else {
			break
		}

		right := p.parseShift()
		expr = &BinaryExpr{
			Range: Range{
				Start: expr.GetRange().Start,
				End:   right.GetRange().End,
			},
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}

	return expr
}

// parseShift parses shift expressions.
func (p *Parser) parseShift() Expr {
	expr := p.parseAdditive()

	for {
		var op string
		if p.match(TokenLeftShift) {
			op = "<<"
		} else if p.match(TokenRightShift) {
			op = ">>"
		} else {
			break
		}

		right := p.parseAdditive()
		expr = &BinaryExpr{
			Range: Range{
				Start: expr.GetRange().Start,
				End:   right.GetRange().End,
			},
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}

	return expr
}

// parseAdditive parses additive expressions.
func (p *Parser) parseAdditive() Expr {
	expr := p.parseMultiplicative()

	for {
		var op string
		if p.match(TokenPlus) {
			op = "+"
		} else if p.match(TokenMinus) {
			op = "-"
		} else {
			break
		}

		right := p.parseMultiplicative()
		expr = &BinaryExpr{
			Range: Range{
				Start: expr.GetRange().Start,
				End:   right.GetRange().End,
			},
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}

	return expr
}

// parseMultiplicative parses multiplicative expressions.
func (p *Parser) parseMultiplicative() Expr {
	expr := p.parseUnary()

	for {
		var op string
		if p.match(TokenStar) {
			op = "*"
		} else if p.match(TokenSlash) {
			op = "/"
		} else if p.match(TokenPercent) {
			op = "%"
		} else {
			break
		}

		right := p.parseUnary()
		expr = &BinaryExpr{
			Range: Range{
				Start: expr.GetRange().Start,
				End:   right.GetRange().End,
			},
			Left:     expr,
			Operator: op,
			Right:    right,
		}
	}

	return expr
}

// parseUnary parses unary expressions.
func (p *Parser) parseUnary() Expr {
	switch p.current().Type {
	case TokenBang, TokenTilde, TokenMinus, TokenPlus:
		op := p.current().Literal
		start := p.current()
		p.advance()
		operand := p.parseUnary()
		return &UnaryExpr{
			Range: Range{
				Start: Position{Line: start.Line - 1, Column: start.Column - 1},
				End:   operand.GetRange().End,
			},
			Operator: op,
			Operand:  operand,
			Prefix:   true,
		}
	case TokenIncrement, TokenDecrement:
		op := p.current().Literal
		start := p.current()
		p.advance()
		operand := p.parseUnary()
		return &UnaryExpr{
			Range: Range{
				Start: Position{Line: start.Line - 1, Column: start.Column - 1},
				End:   operand.GetRange().End,
			},
			Operator: op,
			Operand:  operand,
			Prefix:   true,
		}
	}

	return p.parsePostfix()
}

// parsePostfix parses postfix expressions.
func (p *Parser) parsePostfix() Expr {
	expr := p.parsePrimary()

	for {
		if p.match(TokenLParen) {
			// Function call
			var args []Expr
			for !p.check(TokenRParen) && !p.isAtEnd() {
				args = append(args, p.parseExpression())
				if !p.match(TokenComma) {
					break
				}
			}
			endTok := p.current()
			p.expect(TokenRParen, "expected ')' after arguments")
			expr = &CallExpr{
				Range: Range{
					Start: expr.GetRange().Start,
					End:   Position{Line: endTok.Line - 1, Column: endTok.Column},
				},
				Func: expr,
				Args: args,
			}
		} else if p.match(TokenLBracket) {
			// Array index
			index := p.parseExpression()
			endTok := p.current()
			p.expect(TokenRBracket, "expected ']' after index")
			expr = &IndexExpr{
				Range: Range{
					Start: expr.GetRange().Start,
					End:   Position{Line: endTok.Line - 1, Column: endTok.Column},
				},
				Expr:  expr,
				Index: index,
			}
		} else if p.match(TokenDot) {
			// Member access
			if !p.check(TokenIdent) {
				p.error("expected member name after '.'")
				break
			}
			member := p.current().Literal
			endTok := p.current()
			p.advance()
			expr = &MemberExpr{
				Range: Range{
					Start: expr.GetRange().Start,
					End:   Position{Line: endTok.Line - 1, Column: endTok.Column + len(member)},
				},
				Expr:   expr,
				Member: member,
			}
		} else if p.match(TokenIncrement) {
			// Postfix increment
			expr = &UnaryExpr{
				Range:    expr.GetRange(),
				Operator: "++",
				Operand:  expr,
				Prefix:   false,
			}
		} else if p.match(TokenDecrement) {
			// Postfix decrement
			expr = &UnaryExpr{
				Range:    expr.GetRange(),
				Operator: "--",
				Operand:  expr,
				Prefix:   false,
			}
		} else if p.current().Type.IsAssignmentOp() {
			// Assignment
			op := p.current().Literal
			p.advance()
			right := p.parseExpression()
			expr = &BinaryExpr{
				Range: Range{
					Start: expr.GetRange().Start,
					End:   right.GetRange().End,
				},
				Left:     expr,
				Operator: op,
				Right:    right,
			}
		} else {
			break
		}
	}

	return expr
}

// parsePrimary parses primary expressions.
func (p *Parser) parsePrimary() Expr {
	tok := p.current()

	switch tok.Type {
	case TokenTrue:
		p.advance()
		return &LiteralExpr{
			Range: p.tokenRange(tok),
			Kind:  "bool",
			Value: "true",
		}
	case TokenFalse:
		p.advance()
		return &LiteralExpr{
			Range: p.tokenRange(tok),
			Kind:  "bool",
			Value: "false",
		}
	case TokenIntLit:
		p.advance()
		return &LiteralExpr{
			Range: p.tokenRange(tok),
			Kind:  "int",
			Value: tok.Literal,
		}
	case TokenFloatLit:
		p.advance()
		return &LiteralExpr{
			Range: p.tokenRange(tok),
			Kind:  "float",
			Value: tok.Literal,
		}
	case TokenIdent:
		p.advance()
		return &IdentExpr{
			Range: p.tokenRange(tok),
			Name:  tok.Literal,
		}
	case TokenLParen:
		p.advance() // consume '('
		expr := p.parseExpression()
		p.expect(TokenRParen, "expected ')' after expression")
		return expr
	case TokenLBrace:
		// Array initializer { ... }
		return p.parseArrayInitializer()
	default:
		// Check if it's a type constructor
		if tok.Type.IsType() {
			p.advance()
			// Must be followed by '(' for constructor
			if p.check(TokenLParen) {
				p.advance()
				var args []Expr
				for !p.check(TokenRParen) && !p.isAtEnd() {
					args = append(args, p.parseExpression())
					if !p.match(TokenComma) {
						break
					}
				}
				endTok := p.current()
				p.expect(TokenRParen, "expected ')' after constructor arguments")
				return &CallExpr{
					Range: Range{
						Start: p.tokenRange(tok).Start,
						End:   Position{Line: endTok.Line - 1, Column: endTok.Column},
					},
					Func: &IdentExpr{
						Range: p.tokenRange(tok),
						Name:  tok.Literal,
					},
					Args: args,
				}
			}
			// Otherwise just return as identifier
			return &IdentExpr{
				Range: p.tokenRange(tok),
				Name:  tok.Literal,
			}
		}

		p.error(fmt.Sprintf("unexpected token: %s", tok.Literal))
		p.advance()
		return &IdentExpr{
			Range: p.tokenRange(tok),
			Name:  tok.Literal,
		}
	}
}

// parseArrayInitializer parses an array initializer { ... }.
func (p *Parser) parseArrayInitializer() *ArrayExpr {
	start := p.current()
	p.advance() // consume '{'

	arr := &ArrayExpr{
		Range: p.tokenRange(start),
	}

	for !p.check(TokenRBrace) && !p.isAtEnd() {
		arr.Elements = append(arr.Elements, p.parseExpression())
		if !p.match(TokenComma) {
			break
		}
	}

	p.expect(TokenRBrace, "expected '}' after array initializer")
	return arr
}
