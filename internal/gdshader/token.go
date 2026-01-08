// Package gdshader provides a parser and semantic analyzer for Godot shader files (.gdshader).
package gdshader

// TokenType represents the type of a token.
type TokenType int

const (
	// Special tokens
	TokenEOF TokenType = iota
	TokenError
	TokenNewline

	// Comments
	TokenLineComment  // //
	TokenBlockComment // /* */
	TokenDocComment   // /** */

	// Literals
	TokenIdent
	TokenIntLit   // 123, 0x1F
	TokenFloatLit // 1.5, 1e-3, 1.0f

	// Delimiters
	TokenLParen    // (
	TokenRParen    // )
	TokenLBrace    // {
	TokenRBrace    // }
	TokenLBracket  // [
	TokenRBracket  // ]
	TokenSemicolon // ;
	TokenComma     // ,
	TokenDot       // .
	TokenColon     // :

	// Operators
	TokenAssign           // =
	TokenPlus             // +
	TokenMinus            // -
	TokenStar             // *
	TokenSlash            // /
	TokenPercent          // %
	TokenAmpersand        // &
	TokenPipe             // |
	TokenCaret            // ^
	TokenTilde            // ~
	TokenBang             // !
	TokenQuestion         // ?
	TokenLT               // <
	TokenGT               // >
	TokenLTE              // <=
	TokenGTE              // >=
	TokenEQ               // ==
	TokenNE               // !=
	TokenAnd              // &&
	TokenOr               // ||
	TokenLeftShift        // <<
	TokenRightShift       // >>
	TokenPlusAssign       // +=
	TokenMinusAssign      // -=
	TokenStarAssign       // *=
	TokenSlashAssign      // /=
	TokenPercentAssign    // %=
	TokenAmpAssign        // &=
	TokenPipeAssign       // |=
	TokenCaretAssign      // ^=
	TokenLeftShiftAssign  // <<=
	TokenRightShiftAssign // >>=
	TokenIncrement        // ++
	TokenDecrement        // --

	// Keywords - Shader declaration
	TokenShaderType // shader_type
	TokenRenderMode // render_mode

	// Keywords - Storage qualifiers
	TokenUniform
	TokenVarying
	TokenConst
	TokenGlobal        // global uniform
	TokenGroupUniforms // group_uniforms

	// Keywords - Type qualifiers
	TokenIn
	TokenOut
	TokenInout

	// Keywords - Precision qualifiers
	TokenLowp
	TokenMediump
	TokenHighp

	// Keywords - Interpolation qualifiers
	TokenFlat
	TokenSmooth

	// Keywords - Control flow
	TokenIf
	TokenElse
	TokenFor
	TokenWhile
	TokenDo
	TokenSwitch
	TokenCase
	TokenDefault
	TokenBreak
	TokenContinue
	TokenReturn
	TokenDiscard

	// Keywords - Type definitions
	TokenStruct

	// Keywords - Boolean literals
	TokenTrue
	TokenFalse

	// Keywords - Built-in types
	TokenVoid
	TokenBool
	TokenInt
	TokenUint
	TokenFloat
	TokenVec2
	TokenVec3
	TokenVec4
	TokenBvec2
	TokenBvec3
	TokenBvec4
	TokenIvec2
	TokenIvec3
	TokenIvec4
	TokenUvec2
	TokenUvec3
	TokenUvec4
	TokenMat2
	TokenMat3
	TokenMat4
	TokenSampler2D
	TokenISampler2D
	TokenUSampler2D
	TokenSampler2DArray
	TokenISampler2DArray
	TokenUSampler2DArray
	TokenSampler3D
	TokenISampler3D
	TokenUSampler3D
	TokenSamplerCube
	TokenSamplerCubeArray
	TokenSamplerExternalOES
)

// Token represents a lexical token.
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// String returns a string representation of the token type.
func (t TokenType) String() string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenError:
		return "ERROR"
	case TokenNewline:
		return "NEWLINE"
	case TokenLineComment:
		return "LINE_COMMENT"
	case TokenBlockComment:
		return "BLOCK_COMMENT"
	case TokenDocComment:
		return "DOC_COMMENT"
	case TokenIdent:
		return "IDENT"
	case TokenIntLit:
		return "INT"
	case TokenFloatLit:
		return "FLOAT"
	case TokenLParen:
		return "("
	case TokenRParen:
		return ")"
	case TokenLBrace:
		return "{"
	case TokenRBrace:
		return "}"
	case TokenLBracket:
		return "["
	case TokenRBracket:
		return "]"
	case TokenSemicolon:
		return ";"
	case TokenComma:
		return ","
	case TokenDot:
		return "."
	case TokenColon:
		return ":"
	case TokenAssign:
		return "="
	case TokenPlus:
		return "+"
	case TokenMinus:
		return "-"
	case TokenStar:
		return "*"
	case TokenSlash:
		return "/"
	case TokenPercent:
		return "%"
	case TokenAmpersand:
		return "&"
	case TokenPipe:
		return "|"
	case TokenCaret:
		return "^"
	case TokenTilde:
		return "~"
	case TokenBang:
		return "!"
	case TokenQuestion:
		return "?"
	case TokenLT:
		return "<"
	case TokenGT:
		return ">"
	case TokenLTE:
		return "<="
	case TokenGTE:
		return ">="
	case TokenEQ:
		return "=="
	case TokenNE:
		return "!="
	case TokenAnd:
		return "&&"
	case TokenOr:
		return "||"
	case TokenLeftShift:
		return "<<"
	case TokenRightShift:
		return ">>"
	case TokenPlusAssign:
		return "+="
	case TokenMinusAssign:
		return "-="
	case TokenStarAssign:
		return "*="
	case TokenSlashAssign:
		return "/="
	case TokenPercentAssign:
		return "%="
	case TokenAmpAssign:
		return "&="
	case TokenPipeAssign:
		return "|="
	case TokenCaretAssign:
		return "^="
	case TokenLeftShiftAssign:
		return "<<="
	case TokenRightShiftAssign:
		return ">>="
	case TokenIncrement:
		return "++"
	case TokenDecrement:
		return "--"
	case TokenShaderType:
		return "shader_type"
	case TokenRenderMode:
		return "render_mode"
	case TokenUniform:
		return "uniform"
	case TokenVarying:
		return "varying"
	case TokenConst:
		return "const"
	case TokenGlobal:
		return "global"
	case TokenGroupUniforms:
		return "group_uniforms"
	case TokenIn:
		return "in"
	case TokenOut:
		return "out"
	case TokenInout:
		return "inout"
	case TokenLowp:
		return "lowp"
	case TokenMediump:
		return "mediump"
	case TokenHighp:
		return "highp"
	case TokenFlat:
		return "flat"
	case TokenSmooth:
		return "smooth"
	case TokenIf:
		return "if"
	case TokenElse:
		return "else"
	case TokenFor:
		return "for"
	case TokenWhile:
		return "while"
	case TokenDo:
		return "do"
	case TokenSwitch:
		return "switch"
	case TokenCase:
		return "case"
	case TokenDefault:
		return "default"
	case TokenBreak:
		return "break"
	case TokenContinue:
		return "continue"
	case TokenReturn:
		return "return"
	case TokenDiscard:
		return "discard"
	case TokenStruct:
		return "struct"
	case TokenTrue:
		return "true"
	case TokenFalse:
		return "false"
	case TokenVoid:
		return "void"
	case TokenBool:
		return "bool"
	case TokenInt:
		return "int"
	case TokenUint:
		return "uint"
	case TokenFloat:
		return "float"
	case TokenVec2:
		return "vec2"
	case TokenVec3:
		return "vec3"
	case TokenVec4:
		return "vec4"
	case TokenBvec2:
		return "bvec2"
	case TokenBvec3:
		return "bvec3"
	case TokenBvec4:
		return "bvec4"
	case TokenIvec2:
		return "ivec2"
	case TokenIvec3:
		return "ivec3"
	case TokenIvec4:
		return "ivec4"
	case TokenUvec2:
		return "uvec2"
	case TokenUvec3:
		return "uvec3"
	case TokenUvec4:
		return "uvec4"
	case TokenMat2:
		return "mat2"
	case TokenMat3:
		return "mat3"
	case TokenMat4:
		return "mat4"
	case TokenSampler2D:
		return "sampler2D"
	case TokenISampler2D:
		return "isampler2D"
	case TokenUSampler2D:
		return "usampler2D"
	case TokenSampler2DArray:
		return "sampler2DArray"
	case TokenISampler2DArray:
		return "isampler2DArray"
	case TokenUSampler2DArray:
		return "usampler2DArray"
	case TokenSampler3D:
		return "sampler3D"
	case TokenISampler3D:
		return "isampler3D"
	case TokenUSampler3D:
		return "usampler3D"
	case TokenSamplerCube:
		return "samplerCube"
	case TokenSamplerCubeArray:
		return "samplerCubeArray"
	case TokenSamplerExternalOES:
		return "samplerExternalOES"
	default:
		return "UNKNOWN"
	}
}

// keywords maps keyword strings to their token types.
var keywords = map[string]TokenType{
	// Shader declaration
	"shader_type": TokenShaderType,
	"render_mode": TokenRenderMode,

	// Storage qualifiers
	"uniform":        TokenUniform,
	"varying":        TokenVarying,
	"const":          TokenConst,
	"global":         TokenGlobal,
	"group_uniforms": TokenGroupUniforms,

	// Type qualifiers
	"in":    TokenIn,
	"out":   TokenOut,
	"inout": TokenInout,

	// Precision qualifiers
	"lowp":    TokenLowp,
	"mediump": TokenMediump,
	"highp":   TokenHighp,

	// Interpolation qualifiers
	"flat":   TokenFlat,
	"smooth": TokenSmooth,

	// Control flow
	"if":       TokenIf,
	"else":     TokenElse,
	"for":      TokenFor,
	"while":    TokenWhile,
	"do":       TokenDo,
	"switch":   TokenSwitch,
	"case":     TokenCase,
	"default":  TokenDefault,
	"break":    TokenBreak,
	"continue": TokenContinue,
	"return":   TokenReturn,
	"discard":  TokenDiscard,

	// Type definition
	"struct": TokenStruct,

	// Boolean literals
	"true":  TokenTrue,
	"false": TokenFalse,

	// Built-in types
	"void":               TokenVoid,
	"bool":               TokenBool,
	"int":                TokenInt,
	"uint":               TokenUint,
	"float":              TokenFloat,
	"vec2":               TokenVec2,
	"vec3":               TokenVec3,
	"vec4":               TokenVec4,
	"bvec2":              TokenBvec2,
	"bvec3":              TokenBvec3,
	"bvec4":              TokenBvec4,
	"ivec2":              TokenIvec2,
	"ivec3":              TokenIvec3,
	"ivec4":              TokenIvec4,
	"uvec2":              TokenUvec2,
	"uvec3":              TokenUvec3,
	"uvec4":              TokenUvec4,
	"mat2":               TokenMat2,
	"mat3":               TokenMat3,
	"mat4":               TokenMat4,
	"sampler2D":          TokenSampler2D,
	"isampler2D":         TokenISampler2D,
	"usampler2D":         TokenUSampler2D,
	"sampler2DArray":     TokenSampler2DArray,
	"isampler2DArray":    TokenISampler2DArray,
	"usampler2DArray":    TokenUSampler2DArray,
	"sampler3D":          TokenSampler3D,
	"isampler3D":         TokenISampler3D,
	"usampler3D":         TokenUSampler3D,
	"samplerCube":        TokenSamplerCube,
	"samplerCubeArray":   TokenSamplerCubeArray,
	"samplerExternalOES": TokenSamplerExternalOES,
}

// LookupIdent checks if an identifier is a keyword.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return TokenIdent
}

// IsKeyword returns true if the token type is a keyword.
func (t TokenType) IsKeyword() bool {
	return t >= TokenShaderType && t <= TokenSamplerExternalOES
}

// IsType returns true if the token type represents a built-in type.
func (t TokenType) IsType() bool {
	return t >= TokenVoid && t <= TokenSamplerExternalOES
}

// IsPrecision returns true if the token type is a precision qualifier.
func (t TokenType) IsPrecision() bool {
	return t == TokenLowp || t == TokenMediump || t == TokenHighp
}

// IsInterpolation returns true if the token type is an interpolation qualifier.
func (t TokenType) IsInterpolation() bool {
	return t == TokenFlat || t == TokenSmooth
}

// IsAssignmentOp returns true if the token type is an assignment operator.
func (t TokenType) IsAssignmentOp() bool {
	switch t {
	case TokenAssign, TokenPlusAssign, TokenMinusAssign, TokenStarAssign,
		TokenSlashAssign, TokenPercentAssign, TokenAmpAssign, TokenPipeAssign,
		TokenCaretAssign, TokenLeftShiftAssign, TokenRightShiftAssign:
		return true
	}
	return false
}
