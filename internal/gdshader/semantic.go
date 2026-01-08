package gdshader

import (
	"fmt"
	"strconv"
	"strings"
)

// SemanticError represents a semantic analysis error.
type SemanticError struct {
	Message string
	Range   Range
}

func (e *SemanticError) Error() string {
	return fmt.Sprintf("%d:%d: %s", e.Range.Start.Line+1, e.Range.Start.Column+1, e.Message)
}

// Symbol represents a symbol in the shader.
type Symbol struct {
	Name       string
	Type       *Type
	Kind       SymbolKind
	Range      Range
	Constant   bool            // For const variables
	ReadOnly   bool            // For built-in input variables
	WriteOnly  bool            // For built-in output variables
	Qualifiers []string        // in, out, inout, uniform, varying, etc.
	Function   *FunctionSymbol // For function symbols
}

// SymbolKind represents the kind of symbol.
type SymbolKind int

const (
	SymbolVariable SymbolKind = iota
	SymbolConstant
	SymbolUniform
	SymbolVarying
	SymbolFunction
	SymbolStruct
	SymbolParameter
	SymbolBuiltinVariable
	SymbolBuiltinFunction
)

// FunctionSymbol holds additional information for function symbols.
type FunctionSymbol struct {
	Params     []*Symbol
	ReturnType *Type
	IsBuiltin  bool
}

// Scope represents a lexical scope for symbol lookup.
type Scope struct {
	parent  *Scope
	symbols map[string]*Symbol
}

func newScope(parent *Scope) *Scope {
	return &Scope{
		parent:  parent,
		symbols: make(map[string]*Symbol),
	}
}

func (s *Scope) define(sym *Symbol) error {
	if existing, ok := s.symbols[sym.Name]; ok {
		return fmt.Errorf("symbol '%s' already defined at line %d", sym.Name, existing.Range.Start.Line+1)
	}
	s.symbols[sym.Name] = sym
	return nil
}

func (s *Scope) lookup(name string) *Symbol {
	if sym, ok := s.symbols[name]; ok {
		return sym
	}
	if s.parent != nil {
		return s.parent.lookup(name)
	}
	return nil
}

// Analyzer performs semantic analysis on a shader AST.
type Analyzer struct {
	doc          *ShaderDocument
	shaderType   ShaderType
	currentScope *Scope
	globalScope  *Scope
	errors       []*SemanticError
	currentFunc  *FunctionDecl
	currentStage string // "vertex", "fragment", "light", etc.
	structs      map[string]*Type
	loopDepth    int
	switchDepth  int
}

// NewAnalyzer creates a new semantic analyzer.
func NewAnalyzer(doc *ShaderDocument) *Analyzer {
	a := &Analyzer{
		doc:     doc,
		structs: make(map[string]*Type),
	}
	a.globalScope = newScope(nil)
	a.currentScope = a.globalScope
	return a
}

// Analyze performs semantic analysis and returns any errors.
func (a *Analyzer) Analyze() []*SemanticError {
	// Determine shader type
	if a.doc.ShaderType != nil {
		a.shaderType = ShaderType(a.doc.ShaderType.Type)
	} else {
		a.addError(Range{Start: Position{Line: 0, Column: 0}}, "missing shader_type declaration")
	}

	// Register built-in constants (ignore redefinition errors for builtins)
	for name, constant := range BuiltinConstants {
		_ = a.globalScope.define(&Symbol{
			Name:     name,
			Type:     TypeFromName(constant.Type),
			Kind:     SymbolConstant,
			Constant: true,
			ReadOnly: true,
		})
	}

	// First pass: register all struct types
	for _, structDecl := range a.doc.Structs {
		a.registerStruct(structDecl)
	}

	// Second pass: register all global symbols (uniforms, varyings, constants, functions)
	for _, uniform := range a.doc.Uniforms {
		a.registerUniform(uniform)
	}
	for _, varying := range a.doc.Varyings {
		a.registerVarying(varying)
	}
	for _, constant := range a.doc.Constants {
		a.registerConstant(constant)
	}
	for _, funcDecl := range a.doc.Functions {
		a.registerFunction(funcDecl)
	}

	// Third pass: analyze function bodies
	for _, funcDecl := range a.doc.Functions {
		a.analyzeFunction(funcDecl)
	}

	return a.errors
}

func (a *Analyzer) addError(rng Range, format string, args ...interface{}) {
	a.errors = append(a.errors, &SemanticError{
		Message: fmt.Sprintf(format, args...),
		Range:   rng,
	})
}

func (a *Analyzer) enterScope() {
	a.currentScope = newScope(a.currentScope)
}

func (a *Analyzer) exitScope() {
	a.currentScope = a.currentScope.parent
}

// registerStruct registers a struct type.
func (a *Analyzer) registerStruct(decl *StructDecl) {
	if _, exists := a.structs[decl.Name]; exists {
		a.addError(decl.Range, "struct '%s' already defined", decl.Name)
		return
	}

	fields := make([]*Field, 0, len(decl.Members))
	for _, member := range decl.Members {
		fieldType := a.resolveType(member.Type)
		if fieldType == nil {
			a.addError(member.Type.Range, "unknown type '%s'", member.Type.Name)
			fieldType = TypeError
		}
		fields = append(fields, &Field{
			Name: member.Name,
			Type: fieldType,
		})
	}

	structType := MakeStructType(decl.Name, fields)
	a.structs[decl.Name] = structType

	_ = a.globalScope.define(&Symbol{
		Name:  decl.Name,
		Type:  structType,
		Kind:  SymbolStruct,
		Range: decl.Range,
	})
}

// registerUniform registers a uniform variable.
func (a *Analyzer) registerUniform(decl *UniformDecl) {
	varType := a.resolveType(decl.Type)
	if varType == nil {
		a.addError(decl.Type.Range, "unknown type '%s'", decl.Type.Name)
		varType = TypeError
	}

	if err := a.globalScope.define(&Symbol{
		Name:       decl.Name,
		Type:       varType,
		Kind:       SymbolUniform,
		Range:      decl.Range,
		ReadOnly:   true,
		Qualifiers: []string{"uniform"},
	}); err != nil {
		a.addError(decl.Range, "%s", err.Error())
	}
}

// registerVarying registers a varying variable.
func (a *Analyzer) registerVarying(decl *VaryingDecl) {
	varType := a.resolveType(decl.Type)
	if varType == nil {
		a.addError(decl.Type.Range, "unknown type '%s'", decl.Type.Name)
		varType = TypeError
	}

	qualifiers := []string{"varying"}
	if decl.Interpolation != "" {
		qualifiers = append(qualifiers, decl.Interpolation)
	}

	if err := a.globalScope.define(&Symbol{
		Name:       decl.Name,
		Type:       varType,
		Kind:       SymbolVarying,
		Range:      decl.Range,
		Qualifiers: qualifiers,
	}); err != nil {
		a.addError(decl.Range, "%s", err.Error())
	}
}

// registerConstant registers a constant variable.
func (a *Analyzer) registerConstant(decl *ConstDecl) {
	varType := a.resolveType(decl.Type)
	if varType == nil {
		a.addError(decl.Type.Range, "unknown type '%s'", decl.Type.Name)
		varType = TypeError
	}

	// Analyze the initializer
	if decl.Value != nil {
		initType := a.analyzeExpr(decl.Value)
		if !varType.Equals(initType) && !CanImplicitlyConvert(initType, varType) {
			a.addError(decl.Range, "cannot initialize '%s' of type '%s' with '%s'",
				decl.Name, varType.String(), initType.String())
		}
	}

	if err := a.globalScope.define(&Symbol{
		Name:       decl.Name,
		Type:       varType,
		Kind:       SymbolConstant,
		Range:      decl.Range,
		Constant:   true,
		ReadOnly:   true,
		Qualifiers: []string{"const"},
	}); err != nil {
		a.addError(decl.Range, "%s", err.Error())
	}
}

// registerFunction registers a function declaration.
func (a *Analyzer) registerFunction(decl *FunctionDecl) {
	returnType := a.resolveType(decl.ReturnType)
	if returnType == nil {
		a.addError(decl.ReturnType.Range, "unknown type '%s'", decl.ReturnType.Name)
		returnType = TypeError
	}

	params := make([]*Symbol, 0, len(decl.Params))
	for _, param := range decl.Params {
		paramType := a.resolveType(param.Type)
		if paramType == nil {
			a.addError(param.Type.Range, "unknown type '%s'", param.Type.Name)
			paramType = TypeError
		}
		params = append(params, &Symbol{
			Name:       param.Name,
			Type:       paramType,
			Kind:       SymbolParameter,
			Range:      param.Range,
			Qualifiers: []string{param.Qualifier},
		})
	}

	funcSym := &Symbol{
		Name:  decl.Name,
		Type:  returnType,
		Kind:  SymbolFunction,
		Range: decl.Range,
		Function: &FunctionSymbol{
			Params:     params,
			ReturnType: returnType,
		},
	}

	if err := a.globalScope.define(funcSym); err != nil {
		a.addError(decl.Range, "%s", err.Error())
	}
}

// analyzeFunction analyzes a function body.
func (a *Analyzer) analyzeFunction(decl *FunctionDecl) {
	a.currentFunc = decl
	a.enterScope()

	// Determine the stage from function name
	switch decl.Name {
	case "vertex":
		a.currentStage = "vertex"
	case "fragment":
		a.currentStage = "fragment"
	case "light":
		a.currentStage = "light"
	case "start", "process":
		a.currentStage = "compute"
	case "sky":
		a.currentStage = "sky"
	case "fog":
		a.currentStage = "fog"
	default:
		a.currentStage = ""
	}

	// Register built-in variables for this stage
	a.registerBuiltinVariables()

	// Register parameters in function scope
	for _, param := range decl.Params {
		paramType := a.resolveType(param.Type)
		if paramType == nil {
			paramType = TypeError
		}
		_ = a.currentScope.define(&Symbol{
			Name:       param.Name,
			Type:       paramType,
			Kind:       SymbolParameter,
			Range:      param.Range,
			Qualifiers: []string{param.Qualifier},
			ReadOnly:   param.Qualifier == "in",
			WriteOnly:  param.Qualifier == "out",
		})
	}

	// Analyze function body
	if decl.Body != nil {
		a.analyzeStmt(decl.Body)
	}

	a.currentStage = ""
	a.currentFunc = nil
	a.exitScope()
}

// registerBuiltinVariables registers built-in variables for the current shader type and stage.
func (a *Analyzer) registerBuiltinVariables() {
	var builtins map[string]*BuiltinVariable

	switch a.shaderType {
	case ShaderTypeSpatial:
		switch a.currentStage {
		case "vertex":
			builtins = GetSpatialVertexBuiltins()
		case "fragment":
			builtins = GetSpatialFragmentBuiltins()
		case "light":
			builtins = GetSpatialLightBuiltins()
		}
	case ShaderTypeCanvasItem:
		switch a.currentStage {
		case "vertex":
			builtins = GetCanvasItemVertexBuiltins()
		case "fragment":
			builtins = GetCanvasItemFragmentBuiltins()
		case "light":
			builtins = GetCanvasItemLightBuiltins()
		}
	case ShaderTypeParticles:
		switch a.currentStage {
		case "start", "process":
			builtins = GetParticlesBuiltins()
		}
	case ShaderTypeSky:
		switch a.currentStage {
		case "sky":
			builtins = GetSkyBuiltins()
		}
	case ShaderTypeFog:
		switch a.currentStage {
		case "fog":
			builtins = GetFogBuiltins()
		}
	}

	for name, builtin := range builtins {
		varType := TypeFromName(builtin.Type)
		if varType == nil {
			continue // Unknown type, skip
		}
		_ = a.currentScope.define(&Symbol{
			Name:     name,
			Type:     varType,
			Kind:     SymbolBuiltinVariable,
			ReadOnly: builtin.ReadWrite == "in",
		})
	}
}

// resolveType resolves a type specification to a Type.
func (a *Analyzer) resolveType(typeSpec *TypeSpec) *Type {
	if typeSpec == nil {
		return nil
	}

	// Check for built-in types
	if t := TypeFromName(typeSpec.Name); t != nil {
		if typeSpec.ArraySize != nil {
			size := a.evaluateConstExpr(typeSpec.ArraySize)
			return MakeArrayType(t, size)
		}
		return t
	}

	// Check for struct types
	if structType, ok := a.structs[typeSpec.Name]; ok {
		if typeSpec.ArraySize != nil {
			size := a.evaluateConstExpr(typeSpec.ArraySize)
			return MakeArrayType(structType, size)
		}
		return structType
	}

	return nil
}

// evaluateConstExpr evaluates a constant expression and returns its integer value.
func (a *Analyzer) evaluateConstExpr(expr Expr) int {
	switch e := expr.(type) {
	case *LiteralExpr:
		if e.Kind == "int" {
			val, _ := strconv.Atoi(e.Value)
			return val
		}
	case *IdentExpr:
		sym := a.currentScope.lookup(e.Name)
		if sym != nil && sym.Constant {
			// For now, just return -1 for non-literal constants
			return -1
		}
	}
	return -1 // Unsized or error
}

// analyzeStmt analyzes a statement.
func (a *Analyzer) analyzeStmt(stmt Stmt) {
	if stmt == nil {
		return
	}

	switch s := stmt.(type) {
	case *BlockStmt:
		a.enterScope()
		for _, st := range s.Stmts {
			a.analyzeStmt(st)
		}
		a.exitScope()

	case *VarDeclStmt:
		a.analyzeVarDeclStmt(s)

	case *ExprStmt:
		a.analyzeExpr(s.Expr)

	case *IfStmt:
		condType := a.analyzeExpr(s.Cond)
		if condType.Kind != TypeKindBool {
			a.addError(s.Cond.GetRange(), "condition must be a boolean expression, got '%s'", condType.String())
		}
		a.analyzeStmt(s.Then)
		if s.Else != nil {
			a.analyzeStmt(s.Else)
		}

	case *ForStmt:
		a.enterScope()
		a.loopDepth++
		if s.Init != nil {
			a.analyzeStmt(s.Init)
		}
		if s.Cond != nil {
			condType := a.analyzeExpr(s.Cond)
			if condType.Kind != TypeKindBool {
				a.addError(s.Cond.GetRange(), "for condition must be a boolean expression")
			}
		}
		if s.Post != nil {
			a.analyzeExpr(s.Post)
		}
		a.analyzeStmt(s.Body)
		a.loopDepth--
		a.exitScope()

	case *WhileStmt:
		condType := a.analyzeExpr(s.Cond)
		if condType.Kind != TypeKindBool {
			a.addError(s.Cond.GetRange(), "while condition must be a boolean expression")
		}
		a.loopDepth++
		a.analyzeStmt(s.Body)
		a.loopDepth--

	case *DoWhileStmt:
		a.loopDepth++
		a.analyzeStmt(s.Body)
		a.loopDepth--
		condType := a.analyzeExpr(s.Cond)
		if condType.Kind != TypeKindBool {
			a.addError(s.Cond.GetRange(), "do-while condition must be a boolean expression")
		}

	case *SwitchStmt:
		exprType := a.analyzeExpr(s.Expr)
		if !exprType.IsInteger() {
			a.addError(s.Expr.GetRange(), "switch expression must be an integer type")
		}
		a.switchDepth++
		for _, c := range s.Cases {
			for _, val := range c.Values {
				caseType := a.analyzeExpr(val)
				if !caseType.IsInteger() {
					a.addError(val.GetRange(), "case value must be an integer")
				}
			}
			for _, st := range c.Body {
				a.analyzeStmt(st)
			}
		}
		a.switchDepth--

	case *ReturnStmt:
		if a.currentFunc == nil {
			a.addError(s.Range, "return statement outside function")
			return
		}
		returnType := a.resolveType(a.currentFunc.ReturnType)
		if s.Value != nil {
			exprType := a.analyzeExpr(s.Value)
			if returnType.Kind == TypeKindVoid {
				a.addError(s.Range, "void function should not return a value")
			} else if !returnType.Equals(exprType) && !CanImplicitlyConvert(exprType, returnType) {
				a.addError(s.Range, "cannot return '%s' from function returning '%s'",
					exprType.String(), returnType.String())
			}
		} else if returnType.Kind != TypeKindVoid {
			a.addError(s.Range, "non-void function must return a value")
		}

	case *BreakStmt:
		if a.loopDepth == 0 && a.switchDepth == 0 {
			a.addError(s.Range, "break statement outside loop or switch")
		}

	case *ContinueStmt:
		if a.loopDepth == 0 {
			a.addError(s.Range, "continue statement outside loop")
		}

	case *DiscardStmt:
		if a.currentStage != "fragment" {
			a.addError(s.Range, "discard can only be used in fragment stage")
		}

	case *EmptyStmt:
		// Nothing to analyze
	}
}

// analyzeVarDeclStmt analyzes a variable declaration statement.
func (a *Analyzer) analyzeVarDeclStmt(s *VarDeclStmt) {
	varType := a.resolveType(s.Type)
	if varType == nil {
		a.addError(s.Type.Range, "unknown type '%s'", s.Type.Name)
		varType = TypeError
	}

	for _, decl := range s.Decls {
		declType := varType
		// Handle array declaration
		if decl.ArraySize != nil {
			size := a.evaluateConstExpr(decl.ArraySize)
			declType = MakeArrayType(varType, size)
		}

		// Check for initializer
		if decl.Init != nil {
			initType := a.analyzeExpr(decl.Init)
			if !declType.Equals(initType) && !CanImplicitlyConvert(initType, declType) {
				a.addError(decl.Range, "cannot initialize '%s' of type '%s' with '%s'",
					decl.Name, declType.String(), initType.String())
			}
		}

		if err := a.currentScope.define(&Symbol{
			Name:     decl.Name,
			Type:     declType,
			Kind:     SymbolVariable,
			Range:    decl.Range,
			Constant: s.Const,
			ReadOnly: s.Const,
		}); err != nil {
			a.addError(decl.Range, "%s", err.Error())
		}
	}
}

// analyzeExpr analyzes an expression and returns its type.
func (a *Analyzer) analyzeExpr(expr Expr) *Type {
	if expr == nil {
		return TypeError
	}

	switch e := expr.(type) {
	case *LiteralExpr:
		return a.analyzeLiteral(e)

	case *IdentExpr:
		return a.analyzeIdent(e)

	case *BinaryExpr:
		return a.analyzeBinary(e)

	case *UnaryExpr:
		return a.analyzeUnary(e)

	case *TernaryExpr:
		return a.analyzeTernary(e)

	case *CallExpr:
		return a.analyzeCall(e)

	case *IndexExpr:
		return a.analyzeIndex(e)

	case *MemberExpr:
		return a.analyzeMember(e)

	case *ArrayExpr:
		return a.analyzeArrayExpr(e)

	default:
		return TypeError
	}
}

// analyzeLiteral analyzes a literal expression.
func (a *Analyzer) analyzeLiteral(e *LiteralExpr) *Type {
	switch e.Kind {
	case "int":
		return TypeInt
	case "float":
		return TypeFloat
	case "bool":
		return TypeBool
	default:
		return TypeError
	}
}

// analyzeIdent analyzes an identifier expression.
func (a *Analyzer) analyzeIdent(e *IdentExpr) *Type {
	// Check for built-in constants
	if constant, ok := BuiltinConstants[e.Name]; ok {
		return TypeFromName(constant.Type)
	}

	sym := a.currentScope.lookup(e.Name)
	if sym == nil {
		a.addError(e.Range, "undefined symbol '%s'", e.Name)
		return TypeError
	}
	return sym.Type
}

// analyzeBinary analyzes a binary expression.
func (a *Analyzer) analyzeBinary(e *BinaryExpr) *Type {
	leftType := a.analyzeExpr(e.Left)
	rightType := a.analyzeExpr(e.Right)

	op := operatorToTokenType(e.Operator)

	// Handle assignment operators
	if isAssignOp(op) {
		a.checkAssignable(e.Left)
		if op == TokenAssign {
			if !leftType.Equals(rightType) && !CanImplicitlyConvert(rightType, leftType) {
				a.addError(e.Range, "cannot assign '%s' to '%s'", rightType.String(), leftType.String())
			}
		} else {
			// Compound assignment - check the underlying operation is valid
			underlyingOp := compoundToOp(op)
			resultType := BinaryOpResultType(underlyingOp, leftType, rightType)
			if resultType.Kind == TypeKindError {
				a.addError(e.Range, "invalid operands for '%s': '%s' and '%s'",
					e.Operator, leftType.String(), rightType.String())
			}
		}
		return leftType
	}

	resultType := BinaryOpResultType(op, leftType, rightType)
	if resultType.Kind == TypeKindError {
		a.addError(e.Range, "invalid operands for '%s': '%s' and '%s'",
			e.Operator, leftType.String(), rightType.String())
	}
	return resultType
}

// analyzeUnary analyzes a unary expression.
func (a *Analyzer) analyzeUnary(e *UnaryExpr) *Type {
	operandType := a.analyzeExpr(e.Operand)
	op := operatorToTokenType(e.Operator)

	// Pre/post increment/decrement require assignable operand
	if op == TokenIncrement || op == TokenDecrement {
		a.checkAssignable(e.Operand)
	}

	resultType := UnaryOpResultType(op, operandType)
	if resultType.Kind == TypeKindError {
		a.addError(e.Range, "invalid operand for '%s': '%s'",
			e.Operator, operandType.String())
	}
	return resultType
}

// analyzeTernary analyzes a ternary expression.
func (a *Analyzer) analyzeTernary(e *TernaryExpr) *Type {
	condType := a.analyzeExpr(e.Cond)
	if condType.Kind != TypeKindBool {
		a.addError(e.Cond.GetRange(), "ternary condition must be boolean, got '%s'", condType.String())
	}

	thenType := a.analyzeExpr(e.Then)
	elseType := a.analyzeExpr(e.Else)

	if thenType.Equals(elseType) {
		return thenType
	}
	if CanImplicitlyConvert(elseType, thenType) {
		return thenType
	}
	if CanImplicitlyConvert(thenType, elseType) {
		return elseType
	}

	a.addError(e.Range, "incompatible types in ternary expression: '%s' and '%s'",
		thenType.String(), elseType.String())
	return TypeError
}

// analyzeCall analyzes a function call expression.
func (a *Analyzer) analyzeCall(e *CallExpr) *Type {
	funcName := ""
	if ident, ok := e.Func.(*IdentExpr); ok {
		funcName = ident.Name
	} else {
		a.addError(e.Range, "expected function name")
		return TypeError
	}

	// Check for type constructor
	if IsBuiltinTypeName(funcName) {
		return a.analyzeTypeConstructor(funcName, e)
	}

	// Check for built-in function
	if builtin, ok := BuiltinFunctions[funcName]; ok {
		return a.analyzeBuiltinCall(builtin, e)
	}

	// Check for user-defined function
	sym := a.currentScope.lookup(funcName)
	if sym == nil {
		a.addError(e.Range, "undefined function '%s'", funcName)
		return TypeError
	}
	if sym.Kind != SymbolFunction || sym.Function == nil {
		a.addError(e.Range, "'%s' is not a function", funcName)
		return TypeError
	}

	// Check argument count
	if len(e.Args) != len(sym.Function.Params) {
		a.addError(e.Range, "function '%s' expects %d arguments, got %d",
			funcName, len(sym.Function.Params), len(e.Args))
		return sym.Function.ReturnType
	}

	// Check argument types
	for i, arg := range e.Args {
		argType := a.analyzeExpr(arg)
		paramType := sym.Function.Params[i].Type
		qualifier := ""
		if len(sym.Function.Params[i].Qualifiers) > 0 {
			qualifier = sym.Function.Params[i].Qualifiers[0]
		}

		if qualifier == "out" || qualifier == "inout" {
			a.checkAssignable(arg)
		}

		if !paramType.Equals(argType) && !CanImplicitlyConvert(argType, paramType) {
			a.addError(arg.GetRange(), "argument %d: cannot convert '%s' to '%s'",
				i+1, argType.String(), paramType.String())
		}
	}

	return sym.Function.ReturnType
}

// analyzeTypeConstructor analyzes a type constructor call.
func (a *Analyzer) analyzeTypeConstructor(typeName string, e *CallExpr) *Type {
	targetType := TypeFromName(typeName)
	if targetType == nil {
		return TypeError
	}

	if len(e.Args) == 0 {
		a.addError(e.Range, "type constructor requires at least one argument")
		return targetType
	}

	// Analyze all arguments
	argTypes := make([]*Type, len(e.Args))
	for i, arg := range e.Args {
		argTypes[i] = a.analyzeExpr(arg)
	}

	// For scalar types, accept exactly one argument
	if targetType.IsScalar() {
		if len(e.Args) != 1 {
			a.addError(e.Range, "scalar constructor requires exactly one argument")
		}
		return targetType
	}

	// For vector types, validate component count
	if targetType.IsVector() {
		size := targetType.VectorSize()
		totalComponents := 0
		for _, argType := range argTypes {
			if argType.IsScalar() {
				totalComponents++
			} else if argType.IsVector() {
				totalComponents += argType.VectorSize()
			} else {
				a.addError(e.Range, "invalid argument type '%s' for vector constructor", argType.String())
			}
		}
		// Single scalar fills all components
		if len(e.Args) == 1 && argTypes[0].IsScalar() {
			return targetType
		}
		if totalComponents != size {
			a.addError(e.Range, "vector constructor requires %d components, got %d", size, totalComponents)
		}
		return targetType
	}

	// For matrix types
	if targetType.IsMatrix() {
		size := targetType.MatrixSize()
		// Single scalar creates diagonal matrix
		if len(e.Args) == 1 && argTypes[0].IsScalar() {
			return targetType
		}
		// Full construction
		totalComponents := 0
		for _, argType := range argTypes {
			if argType.IsScalar() {
				totalComponents++
			} else if argType.IsVector() {
				totalComponents += argType.VectorSize()
			} else if argType.IsMatrix() {
				totalComponents += argType.MatrixSize() * argType.MatrixSize()
			}
		}
		expected := size * size
		if totalComponents != expected && (len(e.Args) != 1 || !argTypes[0].IsMatrix()) {
			a.addError(e.Range, "matrix constructor requires %d components, got %d", expected, totalComponents)
		}
		return targetType
	}

	return targetType
}

// analyzeBuiltinCall analyzes a built-in function call.
func (a *Analyzer) analyzeBuiltinCall(builtin *BuiltinFunction, e *CallExpr) *Type {
	// Analyze all arguments first
	argTypes := make([]*Type, len(e.Args))
	for i, arg := range e.Args {
		argTypes[i] = a.analyzeExpr(arg)
	}

	// Find a matching signature
	for _, sig := range builtin.Signatures {
		if len(sig.Params) != len(argTypes) {
			continue
		}
		matches := true
		for i, paramTypeName := range sig.Params {
			paramType := TypeFromName(paramTypeName)
			if paramType == nil {
				matches = false
				break
			}
			if !paramType.Equals(argTypes[i]) && !CanImplicitlyConvert(argTypes[i], paramType) {
				matches = false
				break
			}
		}
		if matches {
			return TypeFromName(sig.Return)
		}
	}

	// No matching signature found
	argTypeStrs := make([]string, len(argTypes))
	for i, t := range argTypes {
		argTypeStrs[i] = t.String()
	}
	a.addError(e.Range, "no matching overload for '%s(%s)'",
		builtin.Name, strings.Join(argTypeStrs, ", "))

	// Return the first signature's return type as a fallback
	if len(builtin.Signatures) > 0 {
		return TypeFromName(builtin.Signatures[0].Return)
	}
	return TypeError
}

// analyzeIndex analyzes an index expression.
func (a *Analyzer) analyzeIndex(e *IndexExpr) *Type {
	baseType := a.analyzeExpr(e.Expr)
	indexType := a.analyzeExpr(e.Index)

	if !indexType.IsInteger() {
		a.addError(e.Index.GetRange(), "index must be an integer, got '%s'", indexType.String())
	}

	resultType := IndexResultType(baseType)
	if resultType.Kind == TypeKindError {
		a.addError(e.Range, "cannot index type '%s'", baseType.String())
	}
	return resultType
}

// analyzeMember analyzes a member access expression.
func (a *Analyzer) analyzeMember(e *MemberExpr) *Type {
	baseType := a.analyzeExpr(e.Expr)

	// Check for swizzle on vector types
	if baseType.IsVector() {
		resultType, err := ValidateSwizzle(baseType, e.Member)
		if err != nil {
			a.addError(e.Range, "%s", err.Error())
			return TypeError
		}
		return resultType
	}

	// Check for struct field access
	if baseType.Kind == TypeKindStruct {
		for _, field := range baseType.Fields {
			if field.Name == e.Member {
				return field.Type
			}
		}
		a.addError(e.Range, "struct '%s' has no field '%s'", baseType.Name, e.Member)
		return TypeError
	}

	// Matrix column access (like swizzle but single component returns vector)
	if baseType.IsMatrix() {
		// mat[i] returns a column vector, but mat.xyz doesn't make sense
		a.addError(e.Range, "cannot use member access on matrix type, use index instead")
		return TypeError
	}

	a.addError(e.Range, "cannot access member '%s' on type '%s'", e.Member, baseType.String())
	return TypeError
}

// analyzeArrayExpr analyzes an array initialization expression.
func (a *Analyzer) analyzeArrayExpr(e *ArrayExpr) *Type {
	if len(e.Elements) == 0 {
		a.addError(e.Range, "empty array initializer")
		return TypeError
	}

	// All elements must have the same type
	elemType := a.analyzeExpr(e.Elements[0])
	for i := 1; i < len(e.Elements); i++ {
		t := a.analyzeExpr(e.Elements[i])
		if !t.Equals(elemType) && !CanImplicitlyConvert(t, elemType) {
			a.addError(e.Elements[i].GetRange(), "array element type mismatch: expected '%s', got '%s'",
				elemType.String(), t.String())
		}
	}

	return MakeArrayType(elemType, len(e.Elements))
}

// checkAssignable verifies that an expression can be assigned to.
func (a *Analyzer) checkAssignable(expr Expr) {
	switch e := expr.(type) {
	case *IdentExpr:
		sym := a.currentScope.lookup(e.Name)
		if sym == nil {
			return // Error already reported
		}
		if sym.Constant || sym.ReadOnly {
			a.addError(e.Range, "cannot assign to '%s' (read-only)", e.Name)
		}

	case *IndexExpr:
		a.checkAssignable(e.Expr)

	case *MemberExpr:
		// Swizzle assignment is valid for single components or all different components
		baseType := a.analyzeExpr(e.Expr)
		if baseType.IsVector() {
			// Check for duplicate swizzle components
			seen := make(map[rune]bool)
			for _, ch := range e.Member {
				if seen[ch] {
					a.addError(e.Range, "cannot assign to swizzle with duplicate components")
					return
				}
				seen[ch] = true
			}
		}
		a.checkAssignable(e.Expr)

	default:
		a.addError(expr.GetRange(), "expression is not assignable")
	}
}

// operatorToTokenType converts an operator string to TokenType.
func operatorToTokenType(op string) TokenType {
	switch op {
	case "+":
		return TokenPlus
	case "-":
		return TokenMinus
	case "*":
		return TokenStar
	case "/":
		return TokenSlash
	case "%":
		return TokenPercent
	case "&":
		return TokenAmpersand
	case "|":
		return TokenPipe
	case "^":
		return TokenCaret
	case "<":
		return TokenLT
	case ">":
		return TokenGT
	case "<=":
		return TokenLTE
	case ">=":
		return TokenGTE
	case "==":
		return TokenEQ
	case "!=":
		return TokenNE
	case "&&":
		return TokenAnd
	case "||":
		return TokenOr
	case "<<":
		return TokenLeftShift
	case ">>":
		return TokenRightShift
	case "=":
		return TokenAssign
	case "+=":
		return TokenPlusAssign
	case "-=":
		return TokenMinusAssign
	case "*=":
		return TokenStarAssign
	case "/=":
		return TokenSlashAssign
	case "%=":
		return TokenPercentAssign
	case "&=":
		return TokenAmpAssign
	case "|=":
		return TokenPipeAssign
	case "^=":
		return TokenCaretAssign
	case "<<=":
		return TokenLeftShiftAssign
	case ">>=":
		return TokenRightShiftAssign
	case "!":
		return TokenBang
	case "~":
		return TokenTilde
	case "++":
		return TokenIncrement
	case "--":
		return TokenDecrement
	default:
		return TokenError
	}
}

// isAssignOp returns true if the token is an assignment operator.
func isAssignOp(op TokenType) bool {
	switch op {
	case TokenAssign, TokenPlusAssign, TokenMinusAssign, TokenStarAssign,
		TokenSlashAssign, TokenPercentAssign, TokenAmpAssign, TokenPipeAssign,
		TokenCaretAssign, TokenLeftShiftAssign, TokenRightShiftAssign:
		return true
	default:
		return false
	}
}

// compoundToOp converts a compound assignment operator to its underlying operator.
func compoundToOp(op TokenType) TokenType {
	switch op {
	case TokenPlusAssign:
		return TokenPlus
	case TokenMinusAssign:
		return TokenMinus
	case TokenStarAssign:
		return TokenStar
	case TokenSlashAssign:
		return TokenSlash
	case TokenPercentAssign:
		return TokenPercent
	case TokenAmpAssign:
		return TokenAmpersand
	case TokenPipeAssign:
		return TokenPipe
	case TokenCaretAssign:
		return TokenCaret
	case TokenLeftShiftAssign:
		return TokenLeftShift
	case TokenRightShiftAssign:
		return TokenRightShift
	default:
		return op
	}
}

// TokenName returns a human-readable name for a token type.
func TokenName(t TokenType) string {
	switch t {
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
	case TokenBang:
		return "!"
	case TokenTilde:
		return "~"
	case TokenIncrement:
		return "++"
	case TokenDecrement:
		return "--"
	case TokenAssign:
		return "="
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
	case TokenLeftShift:
		return "<<"
	case TokenRightShift:
		return ">>"
	default:
		return fmt.Sprintf("Token(%d)", t)
	}
}

// GetErrors returns the list of semantic errors.
func (a *Analyzer) GetErrors() []*SemanticError {
	return a.errors
}

// GetSymbolAt returns the symbol at the given position, if any.
func (a *Analyzer) GetSymbolAt(pos Position) *Symbol {
	// This would require position tracking in the AST
	// For now, return nil
	return nil
}

// GetSymbols returns all symbols in the global scope.
func (a *Analyzer) GetSymbols() []*Symbol {
	symbols := make([]*Symbol, 0)
	for _, sym := range a.globalScope.symbols {
		symbols = append(symbols, sym)
	}
	return symbols
}

// GetStructs returns all defined struct types.
func (a *Analyzer) GetStructs() map[string]*Type {
	return a.structs
}
