package parser

import (
	"strconv"
)

// Parser parses TSCN tokens into an AST.
type Parser struct {
	tokens  []Token
	pos     int
	doc     *Document
	current Token
}

// Parse parses TSCN source code and returns a Document.
func Parse(input string) *Document {
	lexer := NewLexer(input)
	tokens := lexer.Tokenize()

	p := &Parser{
		tokens: tokens,
		pos:    0,
		doc: &Document{
			ExtResources: []*ExtResource{},
			SubResources: []*SubResource{},
			Nodes:        []*Node{},
			Connections:  []*Connection{},
			Comments:     []*Comment{},
			Errors:       []ParseError{},
		},
	}

	if len(tokens) > 0 {
		p.current = tokens[0]
	}

	p.parse()
	return p.doc
}

func (p *Parser) parse() {
	for !p.isAtEnd() {
		p.skipNewlines()
		if p.isAtEnd() {
			break
		}

		switch p.current.Type {
		case TokenComment:
			p.parseComment()
		case TokenLBracket:
			p.parseSection()
		case TokenIdent:
			// Could be a property at the top level (for .tres files)
			p.parseProperty()
		default:
			p.addError("unexpected token: " + p.current.Value)
			p.advance()
		}
	}
}

func (p *Parser) parseComment() {
	p.doc.Comments = append(p.doc.Comments, &Comment{
		Range: p.makeRange(p.current),
		Text:  p.current.Value,
	})
	p.advance()
}

func (p *Parser) parseSection() {
	startToken := p.current
	p.advance() // consume '['

	if p.current.Type != TokenIdent {
		p.addError("expected section type after '['")
		p.skipToNextSection()
		return
	}

	sectionType := p.current.Value
	p.advance()

	switch sectionType {
	case "gd_scene", "gd_resource":
		p.parseGdScene(startToken, sectionType)
	case "ext_resource":
		p.parseExtResource(startToken)
	case "sub_resource":
		p.parseSubResource(startToken)
	case "node":
		p.parseNode(startToken)
	case "connection":
		p.parseConnection(startToken)
	case "resource":
		// [resource] section for .tres files - parse as properties
		p.parseResourceSection(startToken)
	default:
		p.addError("unknown section type: " + sectionType)
		p.skipToNextSection()
	}
}

func (p *Parser) parseGdScene(startToken Token, sectionType string) {
	gd := &GdScene{
		Type:   sectionType,
		Format: 3, // Default to format 3
	}

	// Parse key=value pairs in the header
	for p.current.Type != TokenRBracket && !p.isAtEnd() {
		if p.current.Type == TokenIdent {
			key := p.current.Value
			p.advance()

			if p.current.Type != TokenEquals {
				p.addError("expected '=' after key")
				continue
			}
			p.advance()

			switch key {
			case "load_steps":
				if p.current.Type == TokenNumber {
					val, _ := strconv.Atoi(p.current.Value)
					gd.LoadSteps = &val
					p.advance()
				}
			case "format":
				if p.current.Type == TokenNumber {
					gd.Format, _ = strconv.Atoi(p.current.Value)
					p.advance()
				}
			case "uid":
				if p.current.Type == TokenString {
					gd.UID = p.current.Value
					p.advance()
				}
			case "type":
				if p.current.Type == TokenString {
					gd.ResourceType = p.current.Value
					p.advance()
				}
			default:
				// Skip unknown attributes
				p.parseValue()
			}
		} else {
			break
		}
	}

	if p.current.Type == TokenRBracket {
		p.advance()
	}

	gd.Range = Range{
		Start: Position{Line: startToken.Line, Column: startToken.Column, Offset: startToken.Offset},
		End:   Position{Line: p.prevToken().Line, Column: p.prevToken().Column + p.prevToken().Length, Offset: p.prevToken().Offset + p.prevToken().Length},
	}

	p.doc.Descriptor = gd
}

func (p *Parser) parseExtResource(startToken Token) {
	ext := &ExtResource{}

	// Parse key=value pairs in the header
	for p.current.Type != TokenRBracket && !p.isAtEnd() {
		if p.current.Type == TokenIdent {
			key := p.current.Value
			p.advance()

			if p.current.Type != TokenEquals {
				p.addError("expected '=' after key")
				continue
			}
			p.advance()

			switch key {
			case "type":
				if p.current.Type == TokenString {
					ext.Type = p.current.Value
					p.advance()
				}
			case "uid":
				if p.current.Type == TokenString {
					ext.UID = p.current.Value
					p.advance()
				}
			case "path":
				if p.current.Type == TokenString {
					ext.Path = p.current.Value
					ext.PathRange = p.makeRange(p.current)
					p.advance()
				}
			case "id":
				if p.current.Type == TokenString {
					ext.ID = p.current.Value
					p.advance()
				}
			default:
				p.parseValue()
			}
		} else {
			break
		}
	}

	if p.current.Type == TokenRBracket {
		p.advance()
	}

	ext.Range = Range{
		Start: Position{Line: startToken.Line, Column: startToken.Column, Offset: startToken.Offset},
		End:   Position{Line: p.prevToken().Line, Column: p.prevToken().Column + p.prevToken().Length, Offset: p.prevToken().Offset + p.prevToken().Length},
	}

	p.doc.ExtResources = append(p.doc.ExtResources, ext)
}

func (p *Parser) parseSubResource(startToken Token) {
	sub := &SubResource{
		Properties: []*Property{},
	}

	// Parse key=value pairs in the header
	for p.current.Type != TokenRBracket && !p.isAtEnd() {
		if p.current.Type == TokenIdent {
			key := p.current.Value
			p.advance()

			if p.current.Type != TokenEquals {
				p.addError("expected '=' after key")
				continue
			}
			p.advance()

			switch key {
			case "type":
				if p.current.Type == TokenString {
					sub.Type = p.current.Value
					p.advance()
				}
			case "id":
				if p.current.Type == TokenString {
					sub.ID = p.current.Value
					p.advance()
				}
			default:
				p.parseValue()
			}
		} else {
			break
		}
	}

	endToken := p.current
	if p.current.Type == TokenRBracket {
		p.advance()
	}

	// Parse properties until next section
	p.skipNewlines()
subPropLoop:
	for !p.isAtEnd() && p.current.Type != TokenLBracket {
		switch p.current.Type {
		case TokenComment:
			p.parseComment()
		case TokenIdent:
			prop := p.parseProperty()
			if prop != nil {
				sub.Properties = append(sub.Properties, prop)
			}
		case TokenNewline:
			p.advance()
		default:
			break subPropLoop
		}
	}

	// Update end range to include properties
	if len(sub.Properties) > 0 {
		lastProp := sub.Properties[len(sub.Properties)-1]
		endToken = Token{
			Line:   lastProp.Range.End.Line,
			Column: lastProp.Range.End.Column,
			Offset: lastProp.Range.End.Offset,
		}
	}

	sub.Range = Range{
		Start: Position{Line: startToken.Line, Column: startToken.Column, Offset: startToken.Offset},
		End:   Position{Line: endToken.Line, Column: endToken.Column + endToken.Length, Offset: endToken.Offset + endToken.Length},
	}

	p.doc.SubResources = append(p.doc.SubResources, sub)
}

func (p *Parser) parseNode(startToken Token) {
	node := &Node{
		Properties: []*Property{},
	}

	// Parse key=value pairs in the header
	for p.current.Type != TokenRBracket && !p.isAtEnd() {
		if p.current.Type == TokenIdent {
			key := p.current.Value
			p.advance()

			if p.current.Type != TokenEquals {
				p.addError("expected '=' after key")
				continue
			}
			p.advance()

			switch key {
			case "name":
				if p.current.Type == TokenString {
					node.Name = p.current.Value
					p.advance()
				}
			case "type":
				if p.current.Type == TokenString {
					node.Type = p.current.Value
					p.advance()
				}
			case "parent":
				if p.current.Type == TokenString {
					node.Parent = p.current.Value
					p.advance()
				}
			case "instance":
				node.Instance = p.parseValue()
			case "instance_placeholder":
				if p.current.Type == TokenString {
					node.InstancePlaceholder = p.current.Value
					p.advance()
				}
			case "owner":
				if p.current.Type == TokenString {
					node.Owner = p.current.Value
					p.advance()
				}
			case "index":
				if p.current.Type == TokenNumber {
					val, _ := strconv.Atoi(p.current.Value)
					node.Index = &val
					p.advance()
				}
			case "groups":
				val := p.parseValue()
				if arr, ok := val.(*ArrayValue); ok {
					for _, v := range arr.Values {
						if sv, ok := v.(*StringValue); ok {
							node.Groups = append(node.Groups, sv.Value)
						}
					}
				}
			default:
				p.parseValue()
			}
		} else {
			break
		}
	}

	endToken := p.current
	if p.current.Type == TokenRBracket {
		p.advance()
	}

	// Parse properties until next section
	p.skipNewlines()
nodePropLoop:
	for !p.isAtEnd() && p.current.Type != TokenLBracket {
		switch p.current.Type {
		case TokenComment:
			p.parseComment()
		case TokenIdent:
			prop := p.parseProperty()
			if prop != nil {
				node.Properties = append(node.Properties, prop)
			}
		case TokenNewline:
			p.advance()
		default:
			break nodePropLoop
		}
	}

	// Update end range
	if len(node.Properties) > 0 {
		lastProp := node.Properties[len(node.Properties)-1]
		endToken = Token{
			Line:   lastProp.Range.End.Line,
			Column: lastProp.Range.End.Column,
			Offset: lastProp.Range.End.Offset,
		}
	}

	node.Range = Range{
		Start: Position{Line: startToken.Line, Column: startToken.Column, Offset: startToken.Offset},
		End:   Position{Line: endToken.Line, Column: endToken.Column + endToken.Length, Offset: endToken.Offset + endToken.Length},
	}

	p.doc.Nodes = append(p.doc.Nodes, node)
}

func (p *Parser) parseConnection(startToken Token) {
	conn := &Connection{}

	// Parse key=value pairs in the header
	for p.current.Type != TokenRBracket && !p.isAtEnd() {
		if p.current.Type == TokenIdent {
			key := p.current.Value
			p.advance()

			if p.current.Type != TokenEquals {
				p.addError("expected '=' after key")
				continue
			}
			p.advance()

			switch key {
			case "signal":
				if p.current.Type == TokenString {
					conn.Signal = p.current.Value
					p.advance()
				}
			case "from":
				if p.current.Type == TokenString {
					conn.From = p.current.Value
					p.advance()
				}
			case "to":
				if p.current.Type == TokenString {
					conn.To = p.current.Value
					p.advance()
				}
			case "method":
				if p.current.Type == TokenString {
					conn.Method = p.current.Value
					p.advance()
				}
			case "flags":
				if p.current.Type == TokenNumber {
					val, _ := strconv.Atoi(p.current.Value)
					conn.Flags = &val
					p.advance()
				}
			case "binds":
				val := p.parseValue()
				if arr, ok := val.(*ArrayValue); ok {
					conn.Binds = arr.Values
				}
			default:
				p.parseValue()
			}
		} else {
			break
		}
	}

	if p.current.Type == TokenRBracket {
		p.advance()
	}

	conn.Range = Range{
		Start: Position{Line: startToken.Line, Column: startToken.Column, Offset: startToken.Offset},
		End:   Position{Line: p.prevToken().Line, Column: p.prevToken().Column + p.prevToken().Length, Offset: p.prevToken().Offset + p.prevToken().Length},
	}

	p.doc.Connections = append(p.doc.Connections, conn)
}

func (p *Parser) parseResourceSection(startToken Token) {
	// Skip to ]
	for p.current.Type != TokenRBracket && !p.isAtEnd() {
		p.advance()
	}
	if p.current.Type == TokenRBracket {
		p.advance()
	}

	// Parse properties - they become part of the document's resources
	p.skipNewlines()
resPropLoop:
	for !p.isAtEnd() && p.current.Type != TokenLBracket {
		switch p.current.Type {
		case TokenComment:
			p.parseComment()
		case TokenIdent:
			p.parseProperty()
		case TokenNewline:
			p.advance()
		default:
			break resPropLoop
		}
	}
}

func (p *Parser) parseProperty() *Property {
	if p.current.Type != TokenIdent {
		return nil
	}

	keyStart := Position{Line: p.current.Line, Column: p.current.Column, Offset: p.current.Offset}

	// Parse key (can include slashes for paths like "bones/1/position")
	var key string
keyLoop:
	for {
		switch p.current.Type {
		case TokenIdent, TokenNumber:
			key += p.current.Value
			p.advance()
		case TokenSlash:
			key += "/"
			p.advance()
		default:
			break keyLoop
		}
	}

	keyEnd := Position{Line: p.prevToken().Line, Column: p.prevToken().Column + p.prevToken().Length, Offset: p.prevToken().Offset + p.prevToken().Length}

	if p.current.Type != TokenEquals {
		p.addError("expected '=' after property key")
		return nil
	}
	p.advance()

	value := p.parseValue()
	if value == nil {
		return nil
	}

	return &Property{
		Range: Range{
			Start: keyStart,
			End:   value.GetRange().End,
		},
		Key:      key,
		KeyRange: Range{Start: keyStart, End: keyEnd},
		Value:    value,
	}
}

func (p *Parser) parseValue() Value {
	switch p.current.Type {
	case TokenString:
		val := &StringValue{
			Range: p.makeRange(p.current),
			Value: p.current.Value,
		}
		p.advance()
		return val

	case TokenNumber:
		rawVal := p.current.Value
		floatVal, _ := strconv.ParseFloat(rawVal, 64)
		isInt := !containsAny(rawVal, ".eE") && rawVal != "inf" && rawVal != "nan" && rawVal != "-inf"
		val := &NumberValue{
			Range:    p.makeRange(p.current),
			Value:    floatVal,
			IsInt:    isInt,
			RawValue: rawVal,
		}
		p.advance()
		return val

	case TokenBool:
		val := &BoolValue{
			Range: p.makeRange(p.current),
			Value: p.current.Value == "true",
		}
		p.advance()
		return val

	case TokenNull:
		val := &NullValue{
			Range: p.makeRange(p.current),
		}
		p.advance()
		return val

	case TokenLBracket:
		return p.parseArray()

	case TokenLBrace:
		return p.parseDict()

	case TokenIdent:
		return p.parseTypedValueOrIdent()

	default:
		return nil
	}
}

func (p *Parser) parseArray() Value {
	startToken := p.current
	p.advance() // consume '['

	values := []Value{}
	for p.current.Type != TokenRBracket && !p.isAtEnd() {
		p.skipNewlines()
		if p.current.Type == TokenRBracket {
			break
		}

		val := p.parseValue()
		if val != nil {
			values = append(values, val)
		} else {
			// parseValue returned nil - skip this token to avoid infinite loop
			p.advance()
		}

		p.skipNewlines()
		if p.current.Type == TokenComma {
			p.advance()
		}
	}

	endToken := p.current
	if p.current.Type == TokenRBracket {
		p.advance()
	}

	return &ArrayValue{
		Range: Range{
			Start: Position{Line: startToken.Line, Column: startToken.Column, Offset: startToken.Offset},
			End:   Position{Line: endToken.Line, Column: endToken.Column + endToken.Length, Offset: endToken.Offset + endToken.Length},
		},
		Values: values,
	}
}

func (p *Parser) parseDict() Value {
	startToken := p.current
	p.advance() // consume '{'

	entries := []*DictEntry{}
	for p.current.Type != TokenRBrace && !p.isAtEnd() {
		p.skipNewlines()
		if p.current.Type == TokenRBrace {
			break
		}

		var key Value
		keyStart := p.current
		switch p.current.Type {
		case TokenString:
			key = &StringValue{
				Range: p.makeRange(p.current),
				Value: p.current.Value,
			}
			p.advance()
		case TokenIdent:
			key = &IdentValue{
				Range: p.makeRange(p.current),
				Name:  p.current.Value,
			}
			p.advance()
		default:
			p.addError("expected string or identifier as dictionary key")
			p.advance()
			continue
		}

		if p.current.Type != TokenColon {
			p.addError("expected ':' after dictionary key")
			// Skip to next comma, closing brace, or newline to recover
			for !p.isAtEnd() && p.current.Type != TokenComma && p.current.Type != TokenRBrace && p.current.Type != TokenNewline {
				p.advance()
			}
			continue
		}
		p.advance()

		val := p.parseValue()
		if val != nil {
			entries = append(entries, &DictEntry{
				Range: Range{
					Start: Position{Line: keyStart.Line, Column: keyStart.Column, Offset: keyStart.Offset},
					End:   val.GetRange().End,
				},
				Key:   key,
				Value: val,
			})
		} else {
			// parseValue returned nil - skip to recover
			if p.current.Type != TokenComma && p.current.Type != TokenRBrace {
				p.advance()
			}
		}

		p.skipNewlines()
		if p.current.Type == TokenComma {
			p.advance()
		}
	}

	endToken := p.current
	if p.current.Type == TokenRBrace {
		p.advance()
	}

	return &DictValue{
		Range: Range{
			Start: Position{Line: startToken.Line, Column: startToken.Column, Offset: startToken.Offset},
			End:   Position{Line: endToken.Line, Column: endToken.Column + endToken.Length, Offset: endToken.Offset + endToken.Length},
		},
		Entries: entries,
	}
}

func (p *Parser) parseTypedValueOrIdent() Value {
	startToken := p.current
	name := p.current.Value
	typeRange := p.makeRange(p.current)
	p.advance()

	// Check if it's a function call (typed value)
	if p.current.Type == TokenLParen {
		p.advance() // consume '('

		// Special case for ExtResource and SubResource
		if name == "ExtResource" || name == "SubResource" {
			if p.current.Type == TokenString {
				idRange := p.makeRange(p.current)
				id := p.current.Value
				p.advance()

				endToken := p.current
				if p.current.Type == TokenRParen {
					p.advance()
					endToken = p.prevToken()
				}

				return &ResourceRef{
					Range: Range{
						Start: Position{Line: startToken.Line, Column: startToken.Column, Offset: startToken.Offset},
						End:   Position{Line: endToken.Line, Column: endToken.Column + endToken.Length, Offset: endToken.Offset + endToken.Length},
					},
					RefType: name,
					ID:      id,
					IDRange: idRange,
				}
			}
		}

		// Regular typed value
		args := []Value{}
		for p.current.Type != TokenRParen && !p.isAtEnd() {
			val := p.parseValue()
			if val != nil {
				args = append(args, val)
			} else {
				// parseValue returned nil - skip this token to avoid infinite loop
				if p.current.Type != TokenRParen && p.current.Type != TokenComma {
					p.advance()
				}
			}

			if p.current.Type == TokenComma {
				p.advance()
			}
		}

		endToken := p.current
		if p.current.Type == TokenRParen {
			p.advance()
			endToken = p.prevToken()
		}

		return &TypedValue{
			Range: Range{
				Start: Position{Line: startToken.Line, Column: startToken.Column, Offset: startToken.Offset},
				End:   Position{Line: endToken.Line, Column: endToken.Column + endToken.Length, Offset: endToken.Offset + endToken.Length},
			},
			TypeName:  name,
			TypeRange: typeRange,
			Arguments: args,
		}
	}

	// Plain identifier
	return &IdentValue{
		Range: typeRange,
		Name:  name,
	}
}

func (p *Parser) skipNewlines() {
	for p.current.Type == TokenNewline || p.current.Type == TokenComment {
		if p.current.Type == TokenComment {
			p.parseComment()
		} else {
			p.advance()
		}
	}
}

func (p *Parser) skipToNextSection() {
	for !p.isAtEnd() {
		if p.current.Type == TokenLBracket {
			return
		}
		p.advance()
	}
}

func (p *Parser) advance() {
	if p.pos < len(p.tokens)-1 {
		p.pos++
		p.current = p.tokens[p.pos]
	}
}

func (p *Parser) prevToken() Token {
	if p.pos > 0 {
		return p.tokens[p.pos-1]
	}
	return p.current
}

func (p *Parser) isAtEnd() bool {
	return p.current.Type == TokenEOF
}

const maxErrors = 100 // Limit errors to prevent memory exhaustion on malformed input

func (p *Parser) addError(msg string) {
	if len(p.doc.Errors) >= maxErrors {
		return // Stop accumulating errors after limit
	}
	p.doc.Errors = append(p.doc.Errors, ParseError{
		Range:   p.makeRange(p.current),
		Message: msg,
	})
}

func (p *Parser) makeRange(tok Token) Range {
	return Range{
		Start: Position{Line: tok.Line, Column: tok.Column, Offset: tok.Offset},
		End:   Position{Line: tok.Line, Column: tok.Column + tok.Length, Offset: tok.Offset + tok.Length},
	}
}

func containsAny(s string, chars string) bool {
	for _, c := range chars {
		for _, sc := range s {
			if c == sc {
				return true
			}
		}
	}
	return false
}
