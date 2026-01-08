package gdshader

// Range represents a source code range.
type Range struct {
	Start Position
	End   Position
}

// Position represents a position in source code.
type Position struct {
	Line   int // 0-indexed
	Column int // 0-indexed
}

// Node is the interface for all AST nodes.
type Node interface {
	GetRange() Range
}

// ShaderDocument represents a parsed .gdshader file.
type ShaderDocument struct {
	ShaderType  *ShaderTypeDecl
	RenderModes *RenderModeDecl
	Structs     []*StructDecl
	Uniforms    []*UniformDecl
	Varyings    []*VaryingDecl
	Constants   []*ConstDecl
	Functions   []*FunctionDecl
	Comments    []*Comment
	Errors      []ParseError
}

// ParseError represents a parsing error.
type ParseError struct {
	Range   Range
	Message string
}

// Comment represents a comment in the source code.
type Comment struct {
	Range Range
	Text  string
	IsDoc bool // true for /** */ doc comments
}

// ShaderTypeDecl represents a shader_type declaration.
type ShaderTypeDecl struct {
	Range Range
	Type  string // "spatial", "canvas_item", "particles", "sky", "fog"
}

func (s *ShaderTypeDecl) GetRange() Range { return s.Range }

// RenderModeDecl represents a render_mode declaration.
type RenderModeDecl struct {
	Range Range
	Modes []string
}

func (r *RenderModeDecl) GetRange() Range { return r.Range }

// StructDecl represents a struct declaration.
type StructDecl struct {
	Range   Range
	Name    string
	Members []*StructMember
}

func (s *StructDecl) GetRange() Range { return s.Range }

// StructMember represents a member of a struct.
type StructMember struct {
	Range Range
	Type  *TypeSpec
	Name  string
}

func (s *StructMember) GetRange() Range { return s.Range }

// UniformDecl represents a uniform variable declaration.
type UniformDecl struct {
	Range        Range
	IsGlobal     bool // global uniform
	Type         *TypeSpec
	Name         string
	Hints        []*Hint
	DefaultValue Expr
	DocComment   string // From /** */ comments
	GroupName    string // From group_uniforms
}

func (u *UniformDecl) GetRange() Range { return u.Range }

// Hint represents a uniform hint.
type Hint struct {
	Range Range
	Name  string
	Args  []Expr
}

func (h *Hint) GetRange() Range { return h.Range }

// VaryingDecl represents a varying variable declaration.
type VaryingDecl struct {
	Range         Range
	Interpolation string // "flat", "smooth", ""
	Type          *TypeSpec
	Name          string
}

func (v *VaryingDecl) GetRange() Range { return v.Range }

// ConstDecl represents a constant declaration.
type ConstDecl struct {
	Range Range
	Type  *TypeSpec
	Name  string
	Value Expr
}

func (c *ConstDecl) GetRange() Range { return c.Range }

// FunctionDecl represents a function declaration.
type FunctionDecl struct {
	Range      Range
	ReturnType *TypeSpec
	Name       string
	Params     []*ParamDecl
	Body       *BlockStmt
}

func (f *FunctionDecl) GetRange() Range { return f.Range }

// ParamDecl represents a function parameter.
type ParamDecl struct {
	Range     Range
	Qualifier string // "in", "out", "inout", "const", ""
	Type      *TypeSpec
	Name      string
}

func (p *ParamDecl) GetRange() Range { return p.Range }

// TypeSpec represents a type specification.
type TypeSpec struct {
	Range     Range
	Precision string // "lowp", "mediump", "highp", ""
	Name      string // "vec3", "mat4", "MyStruct", etc.
	ArraySize Expr   // nil if not array
}

func (t *TypeSpec) GetRange() Range { return t.Range }

// Expr is the interface for all expression nodes.
type Expr interface {
	Node
	exprNode()
}

// LiteralExpr represents a literal value.
type LiteralExpr struct {
	Range Range
	Kind  string // "int", "float", "bool"
	Value string
}

func (l *LiteralExpr) GetRange() Range { return l.Range }
func (l *LiteralExpr) exprNode()       {}

// IdentExpr represents an identifier.
type IdentExpr struct {
	Range Range
	Name  string
}

func (i *IdentExpr) GetRange() Range { return i.Range }
func (i *IdentExpr) exprNode()       {}

// BinaryExpr represents a binary expression.
type BinaryExpr struct {
	Range    Range
	Left     Expr
	Operator string
	Right    Expr
}

func (b *BinaryExpr) GetRange() Range { return b.Range }
func (b *BinaryExpr) exprNode()       {}

// UnaryExpr represents a unary expression.
type UnaryExpr struct {
	Range    Range
	Operator string
	Operand  Expr
	Prefix   bool // true for prefix, false for postfix
}

func (u *UnaryExpr) GetRange() Range { return u.Range }
func (u *UnaryExpr) exprNode()       {}

// CallExpr represents a function call or type constructor.
type CallExpr struct {
	Range Range
	Func  Expr // function name or type constructor
	Args  []Expr
}

func (c *CallExpr) GetRange() Range { return c.Range }
func (c *CallExpr) exprNode()       {}

// IndexExpr represents an array index expression.
type IndexExpr struct {
	Range Range
	Expr  Expr
	Index Expr
}

func (i *IndexExpr) GetRange() Range { return i.Range }
func (i *IndexExpr) exprNode()       {}

// MemberExpr represents a member access expression.
type MemberExpr struct {
	Range  Range
	Expr   Expr
	Member string // Could be swizzle (xyz) or struct member
}

func (m *MemberExpr) GetRange() Range { return m.Range }
func (m *MemberExpr) exprNode()       {}

// TernaryExpr represents a ternary conditional expression.
type TernaryExpr struct {
	Range Range
	Cond  Expr
	Then  Expr
	Else  Expr
}

func (t *TernaryExpr) GetRange() Range { return t.Range }
func (t *TernaryExpr) exprNode()       {}

// ArrayExpr represents an array literal.
type ArrayExpr struct {
	Range    Range
	Elements []Expr
}

func (a *ArrayExpr) GetRange() Range { return a.Range }
func (a *ArrayExpr) exprNode()       {}

// Stmt is the interface for all statement nodes.
type Stmt interface {
	Node
	stmtNode()
}

// BlockStmt represents a block of statements.
type BlockStmt struct {
	Range Range
	Stmts []Stmt
}

func (b *BlockStmt) GetRange() Range { return b.Range }
func (b *BlockStmt) stmtNode()       {}

// ExprStmt represents an expression statement.
type ExprStmt struct {
	Range Range
	Expr  Expr
}

func (e *ExprStmt) GetRange() Range { return e.Range }
func (e *ExprStmt) stmtNode()       {}

// VarDeclStmt represents a variable declaration statement.
type VarDeclStmt struct {
	Range Range
	Const bool
	Type  *TypeSpec
	Decls []*VarDecl
}

func (v *VarDeclStmt) GetRange() Range { return v.Range }
func (v *VarDeclStmt) stmtNode()       {}

// VarDecl represents a single variable in a declaration.
type VarDecl struct {
	Range     Range
	Name      string
	ArraySize Expr // nil if not array
	Init      Expr // nil if no initializer
}

func (v *VarDecl) GetRange() Range { return v.Range }

// IfStmt represents an if statement.
type IfStmt struct {
	Range Range
	Cond  Expr
	Then  Stmt
	Else  Stmt // may be nil or another IfStmt for else-if
}

func (i *IfStmt) GetRange() Range { return i.Range }
func (i *IfStmt) stmtNode()       {}

// ForStmt represents a for loop.
type ForStmt struct {
	Range Range
	Init  Stmt // may be nil
	Cond  Expr // may be nil
	Post  Expr // may be nil (note: expression, not statement)
	Body  Stmt
}

func (f *ForStmt) GetRange() Range { return f.Range }
func (f *ForStmt) stmtNode()       {}

// WhileStmt represents a while loop.
type WhileStmt struct {
	Range Range
	Cond  Expr
	Body  Stmt
}

func (w *WhileStmt) GetRange() Range { return w.Range }
func (w *WhileStmt) stmtNode()       {}

// DoWhileStmt represents a do-while loop.
type DoWhileStmt struct {
	Range Range
	Body  Stmt
	Cond  Expr
}

func (d *DoWhileStmt) GetRange() Range { return d.Range }
func (d *DoWhileStmt) stmtNode()       {}

// SwitchStmt represents a switch statement.
type SwitchStmt struct {
	Range Range
	Expr  Expr
	Cases []*CaseClause
}

func (s *SwitchStmt) GetRange() Range { return s.Range }
func (s *SwitchStmt) stmtNode()       {}

// CaseClause represents a case or default clause.
type CaseClause struct {
	Range  Range
	Values []Expr // nil for default
	Body   []Stmt
}

func (c *CaseClause) GetRange() Range { return c.Range }

// ReturnStmt represents a return statement.
type ReturnStmt struct {
	Range Range
	Value Expr // may be nil
}

func (r *ReturnStmt) GetRange() Range { return r.Range }
func (r *ReturnStmt) stmtNode()       {}

// BreakStmt represents a break statement.
type BreakStmt struct {
	Range Range
}

func (b *BreakStmt) GetRange() Range { return b.Range }
func (b *BreakStmt) stmtNode()       {}

// ContinueStmt represents a continue statement.
type ContinueStmt struct {
	Range Range
}

func (c *ContinueStmt) GetRange() Range { return c.Range }
func (c *ContinueStmt) stmtNode()       {}

// DiscardStmt represents a discard statement.
type DiscardStmt struct {
	Range Range
}

func (d *DiscardStmt) GetRange() Range { return d.Range }
func (d *DiscardStmt) stmtNode()       {}

// EmptyStmt represents an empty statement (just a semicolon).
type EmptyStmt struct {
	Range Range
}

func (e *EmptyStmt) GetRange() Range { return e.Range }
func (e *EmptyStmt) stmtNode()       {}
