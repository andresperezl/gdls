package gdshader

import "fmt"

// TypeKind represents the category of a type.
type TypeKind int

const (
	TypeKindVoid TypeKind = iota
	TypeKindBool
	TypeKindInt
	TypeKindUint
	TypeKindFloat
	TypeKindVec2
	TypeKindVec3
	TypeKindVec4
	TypeKindBvec2
	TypeKindBvec3
	TypeKindBvec4
	TypeKindIvec2
	TypeKindIvec3
	TypeKindIvec4
	TypeKindUvec2
	TypeKindUvec3
	TypeKindUvec4
	TypeKindMat2
	TypeKindMat3
	TypeKindMat4
	TypeKindSampler2D
	TypeKindSampler2DArray
	TypeKindSampler3D
	TypeKindSamplerCube
	TypeKindSamplerCubeArray
	TypeKindSamplerExternalOES
	TypeKindIsampler2D
	TypeKindIsampler2DArray
	TypeKindIsampler3D
	TypeKindUsampler2D
	TypeKindUsampler2DArray
	TypeKindUsampler3D
	TypeKindStruct
	TypeKindArray
	TypeKindError // For type errors
)

// Type represents a GDShader type.
type Type struct {
	Kind        TypeKind
	Name        string   // For structs: the struct name
	ArraySize   int      // For arrays: -1 means unsized
	ElementType *Type    // For arrays: the element type
	Fields      []*Field // For structs: the fields
}

// Field represents a struct field.
type Field struct {
	Name string
	Type *Type
}

// String returns the string representation of a type.
func (t *Type) String() string {
	if t == nil {
		return "error"
	}
	switch t.Kind {
	case TypeKindVoid:
		return "void"
	case TypeKindBool:
		return "bool"
	case TypeKindInt:
		return "int"
	case TypeKindUint:
		return "uint"
	case TypeKindFloat:
		return "float"
	case TypeKindVec2:
		return "vec2"
	case TypeKindVec3:
		return "vec3"
	case TypeKindVec4:
		return "vec4"
	case TypeKindBvec2:
		return "bvec2"
	case TypeKindBvec3:
		return "bvec3"
	case TypeKindBvec4:
		return "bvec4"
	case TypeKindIvec2:
		return "ivec2"
	case TypeKindIvec3:
		return "ivec3"
	case TypeKindIvec4:
		return "ivec4"
	case TypeKindUvec2:
		return "uvec2"
	case TypeKindUvec3:
		return "uvec3"
	case TypeKindUvec4:
		return "uvec4"
	case TypeKindMat2:
		return "mat2"
	case TypeKindMat3:
		return "mat3"
	case TypeKindMat4:
		return "mat4"
	case TypeKindSampler2D:
		return "sampler2D"
	case TypeKindSampler2DArray:
		return "sampler2DArray"
	case TypeKindSampler3D:
		return "sampler3D"
	case TypeKindSamplerCube:
		return "samplerCube"
	case TypeKindSamplerCubeArray:
		return "samplerCubeArray"
	case TypeKindSamplerExternalOES:
		return "samplerExternalOES"
	case TypeKindIsampler2D:
		return "isampler2D"
	case TypeKindIsampler2DArray:
		return "isampler2DArray"
	case TypeKindIsampler3D:
		return "isampler3D"
	case TypeKindUsampler2D:
		return "usampler2D"
	case TypeKindUsampler2DArray:
		return "usampler2DArray"
	case TypeKindUsampler3D:
		return "usampler3D"
	case TypeKindStruct:
		return t.Name
	case TypeKindArray:
		if t.ArraySize < 0 {
			return fmt.Sprintf("%s[]", t.ElementType.String())
		}
		return fmt.Sprintf("%s[%d]", t.ElementType.String(), t.ArraySize)
	case TypeKindError:
		return "error"
	default:
		return "unknown"
	}
}

// Equals checks if two types are equal.
func (t *Type) Equals(other *Type) bool {
	if t == nil || other == nil {
		return t == other
	}
	if t.Kind != other.Kind {
		return false
	}
	if t.Kind == TypeKindStruct {
		return t.Name == other.Name
	}
	if t.Kind == TypeKindArray {
		return t.ArraySize == other.ArraySize && t.ElementType.Equals(other.ElementType)
	}
	return true
}

// IsScalar returns true if the type is a scalar type.
func (t *Type) IsScalar() bool {
	switch t.Kind {
	case TypeKindBool, TypeKindInt, TypeKindUint, TypeKindFloat:
		return true
	default:
		return false
	}
}

// IsVector returns true if the type is a vector type.
func (t *Type) IsVector() bool {
	switch t.Kind {
	case TypeKindVec2, TypeKindVec3, TypeKindVec4,
		TypeKindBvec2, TypeKindBvec3, TypeKindBvec4,
		TypeKindIvec2, TypeKindIvec3, TypeKindIvec4,
		TypeKindUvec2, TypeKindUvec3, TypeKindUvec4:
		return true
	default:
		return false
	}
}

// IsMatrix returns true if the type is a matrix type.
func (t *Type) IsMatrix() bool {
	switch t.Kind {
	case TypeKindMat2, TypeKindMat3, TypeKindMat4:
		return true
	default:
		return false
	}
}

// IsSampler returns true if the type is a sampler type.
func (t *Type) IsSampler() bool {
	switch t.Kind {
	case TypeKindSampler2D, TypeKindSampler2DArray, TypeKindSampler3D,
		TypeKindSamplerCube, TypeKindSamplerCubeArray, TypeKindSamplerExternalOES,
		TypeKindIsampler2D, TypeKindIsampler2DArray, TypeKindIsampler3D,
		TypeKindUsampler2D, TypeKindUsampler2DArray, TypeKindUsampler3D:
		return true
	default:
		return false
	}
}

// IsNumeric returns true if the type is numeric (int, uint, or float).
func (t *Type) IsNumeric() bool {
	switch t.Kind {
	case TypeKindInt, TypeKindUint, TypeKindFloat:
		return true
	default:
		return false
	}
}

// IsInteger returns true if the type is an integer type.
func (t *Type) IsInteger() bool {
	switch t.Kind {
	case TypeKindInt, TypeKindUint:
		return true
	default:
		return false
	}
}

// VectorSize returns the number of components for vector types, 0 for non-vectors.
func (t *Type) VectorSize() int {
	switch t.Kind {
	case TypeKindVec2, TypeKindBvec2, TypeKindIvec2, TypeKindUvec2:
		return 2
	case TypeKindVec3, TypeKindBvec3, TypeKindIvec3, TypeKindUvec3:
		return 3
	case TypeKindVec4, TypeKindBvec4, TypeKindIvec4, TypeKindUvec4:
		return 4
	default:
		return 0
	}
}

// MatrixSize returns the dimension for matrix types, 0 for non-matrices.
func (t *Type) MatrixSize() int {
	switch t.Kind {
	case TypeKindMat2:
		return 2
	case TypeKindMat3:
		return 3
	case TypeKindMat4:
		return 4
	default:
		return 0
	}
}

// ComponentType returns the scalar component type for vectors and matrices.
func (t *Type) ComponentType() *Type {
	switch t.Kind {
	case TypeKindVec2, TypeKindVec3, TypeKindVec4:
		return TypeFloat
	case TypeKindBvec2, TypeKindBvec3, TypeKindBvec4:
		return TypeBool
	case TypeKindIvec2, TypeKindIvec3, TypeKindIvec4:
		return TypeInt
	case TypeKindUvec2, TypeKindUvec3, TypeKindUvec4:
		return TypeUint
	case TypeKindMat2, TypeKindMat3, TypeKindMat4:
		return TypeFloat
	default:
		return nil
	}
}

// Predefined types.
var (
	TypeVoid  = &Type{Kind: TypeKindVoid}
	TypeBool  = &Type{Kind: TypeKindBool}
	TypeInt   = &Type{Kind: TypeKindInt}
	TypeUint  = &Type{Kind: TypeKindUint}
	TypeFloat = &Type{Kind: TypeKindFloat}
	TypeVec2  = &Type{Kind: TypeKindVec2}
	TypeVec3  = &Type{Kind: TypeKindVec3}
	TypeVec4  = &Type{Kind: TypeKindVec4}
	TypeBvec2 = &Type{Kind: TypeKindBvec2}
	TypeBvec3 = &Type{Kind: TypeKindBvec3}
	TypeBvec4 = &Type{Kind: TypeKindBvec4}
	TypeIvec2 = &Type{Kind: TypeKindIvec2}
	TypeIvec3 = &Type{Kind: TypeKindIvec3}
	TypeIvec4 = &Type{Kind: TypeKindIvec4}
	TypeUvec2 = &Type{Kind: TypeKindUvec2}
	TypeUvec3 = &Type{Kind: TypeKindUvec3}
	TypeUvec4 = &Type{Kind: TypeKindUvec4}
	TypeMat2  = &Type{Kind: TypeKindMat2}
	TypeMat3  = &Type{Kind: TypeKindMat3}
	TypeMat4  = &Type{Kind: TypeKindMat4}
	TypeError = &Type{Kind: TypeKindError}

	TypeSampler2D          = &Type{Kind: TypeKindSampler2D}
	TypeSampler2DArray     = &Type{Kind: TypeKindSampler2DArray}
	TypeSampler3D          = &Type{Kind: TypeKindSampler3D}
	TypeSamplerCube        = &Type{Kind: TypeKindSamplerCube}
	TypeSamplerCubeArray   = &Type{Kind: TypeKindSamplerCubeArray}
	TypeSamplerExternalOES = &Type{Kind: TypeKindSamplerExternalOES}
	TypeIsampler2D         = &Type{Kind: TypeKindIsampler2D}
	TypeIsampler2DArray    = &Type{Kind: TypeKindIsampler2DArray}
	TypeIsampler3D         = &Type{Kind: TypeKindIsampler3D}
	TypeUsampler2D         = &Type{Kind: TypeKindUsampler2D}
	TypeUsampler2DArray    = &Type{Kind: TypeKindUsampler2DArray}
	TypeUsampler3D         = &Type{Kind: TypeKindUsampler3D}
)

// TypeFromName returns the type for a given type name string.
func TypeFromName(name string) *Type {
	switch name {
	case "void":
		return TypeVoid
	case "bool":
		return TypeBool
	case "int":
		return TypeInt
	case "uint":
		return TypeUint
	case "float":
		return TypeFloat
	case "vec2":
		return TypeVec2
	case "vec3":
		return TypeVec3
	case "vec4":
		return TypeVec4
	case "bvec2":
		return TypeBvec2
	case "bvec3":
		return TypeBvec3
	case "bvec4":
		return TypeBvec4
	case "ivec2":
		return TypeIvec2
	case "ivec3":
		return TypeIvec3
	case "ivec4":
		return TypeIvec4
	case "uvec2":
		return TypeUvec2
	case "uvec3":
		return TypeUvec3
	case "uvec4":
		return TypeUvec4
	case "mat2":
		return TypeMat2
	case "mat3":
		return TypeMat3
	case "mat4":
		return TypeMat4
	case "sampler2D":
		return TypeSampler2D
	case "sampler2DArray":
		return TypeSampler2DArray
	case "sampler3D":
		return TypeSampler3D
	case "samplerCube":
		return TypeSamplerCube
	case "samplerCubeArray":
		return TypeSamplerCubeArray
	case "samplerExternalOES":
		return TypeSamplerExternalOES
	case "isampler2D":
		return TypeIsampler2D
	case "isampler2DArray":
		return TypeIsampler2DArray
	case "isampler3D":
		return TypeIsampler3D
	case "usampler2D":
		return TypeUsampler2D
	case "usampler2DArray":
		return TypeUsampler2DArray
	case "usampler3D":
		return TypeUsampler3D
	default:
		return nil // Unknown type (could be struct)
	}
}

// MakeArrayType creates an array type with the given element type and size.
func MakeArrayType(elemType *Type, size int) *Type {
	return &Type{
		Kind:        TypeKindArray,
		ArraySize:   size,
		ElementType: elemType,
	}
}

// MakeStructType creates a struct type with the given name and fields.
func MakeStructType(name string, fields []*Field) *Type {
	return &Type{
		Kind:   TypeKindStruct,
		Name:   name,
		Fields: fields,
	}
}

// VectorTypeForSize returns the vector type for a given component type and size.
func VectorTypeForSize(componentType *Type, size int) *Type {
	switch componentType.Kind {
	case TypeKindFloat:
		switch size {
		case 2:
			return TypeVec2
		case 3:
			return TypeVec3
		case 4:
			return TypeVec4
		}
	case TypeKindInt:
		switch size {
		case 2:
			return TypeIvec2
		case 3:
			return TypeIvec3
		case 4:
			return TypeIvec4
		}
	case TypeKindUint:
		switch size {
		case 2:
			return TypeUvec2
		case 3:
			return TypeUvec3
		case 4:
			return TypeUvec4
		}
	case TypeKindBool:
		switch size {
		case 2:
			return TypeBvec2
		case 3:
			return TypeBvec3
		case 4:
			return TypeBvec4
		}
	}
	return TypeError
}

// CanImplicitlyConvert checks if srcType can be implicitly converted to dstType.
func CanImplicitlyConvert(srcType, dstType *Type) bool {
	if srcType == nil || dstType == nil {
		return false
	}
	if srcType.Equals(dstType) {
		return true
	}
	// In GLSL/GDShader, there are limited implicit conversions:
	// - int to float
	// - uint to float
	// - int to uint (when non-negative, but we allow it)
	switch {
	case srcType.Kind == TypeKindInt && dstType.Kind == TypeKindFloat:
		return true
	case srcType.Kind == TypeKindUint && dstType.Kind == TypeKindFloat:
		return true
	case srcType.Kind == TypeKindInt && dstType.Kind == TypeKindUint:
		return true
	// Vector versions
	case srcType.Kind == TypeKindIvec2 && dstType.Kind == TypeKindVec2:
		return true
	case srcType.Kind == TypeKindIvec3 && dstType.Kind == TypeKindVec3:
		return true
	case srcType.Kind == TypeKindIvec4 && dstType.Kind == TypeKindVec4:
		return true
	case srcType.Kind == TypeKindUvec2 && dstType.Kind == TypeKindVec2:
		return true
	case srcType.Kind == TypeKindUvec3 && dstType.Kind == TypeKindVec3:
		return true
	case srcType.Kind == TypeKindUvec4 && dstType.Kind == TypeKindVec4:
		return true
	}
	return false
}

// CanExplicitlyConvert checks if srcType can be explicitly converted to dstType.
func CanExplicitlyConvert(srcType, dstType *Type) bool {
	if CanImplicitlyConvert(srcType, dstType) {
		return true
	}
	// Explicit conversions are allowed between most scalar/vector types
	// as long as they have the same number of components
	if srcType.IsScalar() && dstType.IsScalar() {
		return true
	}
	if srcType.IsVector() && dstType.IsVector() {
		return srcType.VectorSize() == dstType.VectorSize()
	}
	if srcType.IsMatrix() && dstType.IsMatrix() {
		return srcType.MatrixSize() == dstType.MatrixSize()
	}
	// float to int/uint
	if srcType.Kind == TypeKindFloat && (dstType.Kind == TypeKindInt || dstType.Kind == TypeKindUint) {
		return true
	}
	return false
}

// CommonType returns the common type for binary operations, or nil if incompatible.
func CommonType(left, right *Type) *Type {
	if left == nil || right == nil {
		return nil
	}
	if left.Equals(right) {
		return left
	}

	// Numeric promotion: prefer float over int/uint
	if left.IsNumeric() && right.IsNumeric() {
		if left.Kind == TypeKindFloat || right.Kind == TypeKindFloat {
			return TypeFloat
		}
		if left.Kind == TypeKindUint || right.Kind == TypeKindUint {
			return TypeUint
		}
		return TypeInt
	}

	// Vector with scalar: promote to vector
	if left.IsVector() && right.IsScalar() {
		if CanImplicitlyConvert(right, left.ComponentType()) {
			return left
		}
	}
	if right.IsVector() && left.IsScalar() {
		if CanImplicitlyConvert(left, right.ComponentType()) {
			return right
		}
	}

	// Two vectors of same size
	if left.IsVector() && right.IsVector() && left.VectorSize() == right.VectorSize() {
		leftComp := left.ComponentType()
		rightComp := right.ComponentType()
		commonComp := CommonType(leftComp, rightComp)
		if commonComp != nil {
			return VectorTypeForSize(commonComp, left.VectorSize())
		}
	}

	return nil
}

// ValidSwizzleChars contains all valid swizzle characters grouped by set.
var ValidSwizzleChars = map[rune]int{
	'x': 0, 'y': 1, 'z': 2, 'w': 3, // xyzw set
	'r': 0, 'g': 1, 'b': 2, 'a': 3, // rgba set
	's': 0, 't': 1, 'p': 2, 'q': 3, // stpq set
}

// SwizzleSet identifies which set a character belongs to (0=xyzw, 1=rgba, 2=stpq).
var SwizzleSet = map[rune]int{
	'x': 0, 'y': 0, 'z': 0, 'w': 0,
	'r': 1, 'g': 1, 'b': 1, 'a': 1,
	's': 2, 't': 2, 'p': 2, 'q': 2,
}

// ValidateSwizzle checks if a swizzle is valid for a given vector type.
// Returns the result type, or nil if invalid.
func ValidateSwizzle(vectorType *Type, swizzle string) (*Type, error) {
	if !vectorType.IsVector() {
		return nil, fmt.Errorf("swizzle can only be applied to vector types")
	}

	vecSize := vectorType.VectorSize()
	if len(swizzle) == 0 || len(swizzle) > 4 {
		return nil, fmt.Errorf("swizzle must have 1-4 components")
	}

	// Check that all characters are from the same set and valid indices
	seenSet := -1
	for _, ch := range swizzle {
		idx, ok := ValidSwizzleChars[ch]
		if !ok {
			return nil, fmt.Errorf("invalid swizzle character '%c'", ch)
		}
		set := SwizzleSet[ch]
		if seenSet == -1 {
			seenSet = set
		} else if set != seenSet {
			return nil, fmt.Errorf("cannot mix swizzle sets (xyzw, rgba, stpq)")
		}
		if idx >= vecSize {
			return nil, fmt.Errorf("swizzle component '%c' invalid for %s", ch, vectorType.String())
		}
	}

	// Result type has same component type but different size
	compType := vectorType.ComponentType()
	if len(swizzle) == 1 {
		return compType, nil
	}
	return VectorTypeForSize(compType, len(swizzle)), nil
}

// BinaryOpResultType returns the result type of a binary operation.
func BinaryOpResultType(op TokenType, left, right *Type) *Type {
	if left == nil || right == nil || left.Kind == TypeKindError || right.Kind == TypeKindError {
		return TypeError
	}

	switch op {
	// Comparison operators return bool
	case TokenLT, TokenLTE, TokenGT, TokenGTE:
		if left.IsNumeric() && right.IsNumeric() {
			return TypeBool
		}
		if left.IsVector() && right.IsVector() && left.VectorSize() == right.VectorSize() {
			return VectorTypeForSize(TypeBool, left.VectorSize())
		}
		return TypeError

	// Equality operators
	case TokenEQ, TokenNE:
		if left.Equals(right) || CanImplicitlyConvert(left, right) || CanImplicitlyConvert(right, left) {
			if left.IsVector() {
				return VectorTypeForSize(TypeBool, left.VectorSize())
			}
			return TypeBool
		}
		return TypeError

	// Logical operators
	case TokenAnd, TokenOr:
		if left.Kind == TypeKindBool && right.Kind == TypeKindBool {
			return TypeBool
		}
		return TypeError

	// Bitwise operators
	case TokenAmpersand, TokenPipe, TokenCaret, TokenLeftShift, TokenRightShift:
		if left.IsInteger() && right.IsInteger() {
			return CommonType(left, right)
		}
		if left.Kind == TypeKindIvec2 && (right.Kind == TypeKindIvec2 || right.Kind == TypeKindInt) {
			return TypeIvec2
		}
		if left.Kind == TypeKindIvec3 && (right.Kind == TypeKindIvec3 || right.Kind == TypeKindInt) {
			return TypeIvec3
		}
		if left.Kind == TypeKindIvec4 && (right.Kind == TypeKindIvec4 || right.Kind == TypeKindInt) {
			return TypeIvec4
		}
		if left.Kind == TypeKindUvec2 && (right.Kind == TypeKindUvec2 || right.Kind == TypeKindUint) {
			return TypeUvec2
		}
		if left.Kind == TypeKindUvec3 && (right.Kind == TypeKindUvec3 || right.Kind == TypeKindUint) {
			return TypeUvec3
		}
		if left.Kind == TypeKindUvec4 && (right.Kind == TypeKindUvec4 || right.Kind == TypeKindUint) {
			return TypeUvec4
		}
		return TypeError

	// Arithmetic operators
	case TokenPlus, TokenMinus, TokenStar, TokenSlash:
		// Matrix * vector
		if left.IsMatrix() && right.IsVector() {
			matSize := left.MatrixSize()
			vecSize := right.VectorSize()
			if matSize == vecSize {
				return right
			}
		}
		// Vector * matrix
		if left.IsVector() && right.IsMatrix() {
			matSize := right.MatrixSize()
			vecSize := left.VectorSize()
			if matSize == vecSize {
				return left
			}
		}
		// Matrix * matrix
		if left.IsMatrix() && right.IsMatrix() {
			if left.MatrixSize() == right.MatrixSize() {
				return left
			}
		}
		// Matrix * scalar or scalar * matrix
		if left.IsMatrix() && right.IsScalar() && right.Kind == TypeKindFloat {
			return left
		}
		if right.IsMatrix() && left.IsScalar() && left.Kind == TypeKindFloat {
			return right
		}

		// Regular numeric/vector operations
		return CommonType(left, right)

	case TokenPercent:
		// Modulo only for integers
		if left.IsInteger() && right.IsInteger() {
			return CommonType(left, right)
		}
		return TypeError
	}

	return TypeError
}

// UnaryOpResultType returns the result type of a unary operation.
func UnaryOpResultType(op TokenType, operand *Type) *Type {
	if operand == nil || operand.Kind == TypeKindError {
		return TypeError
	}

	switch op {
	case TokenMinus:
		if operand.IsNumeric() || operand.IsVector() || operand.IsMatrix() {
			return operand
		}
	case TokenBang:
		if operand.Kind == TypeKindBool {
			return TypeBool
		}
		if operand.Kind == TypeKindBvec2 || operand.Kind == TypeKindBvec3 || operand.Kind == TypeKindBvec4 {
			return operand
		}
	case TokenTilde:
		if operand.IsInteger() {
			return operand
		}
		switch operand.Kind {
		case TypeKindIvec2, TypeKindIvec3, TypeKindIvec4,
			TypeKindUvec2, TypeKindUvec3, TypeKindUvec4:
			return operand
		}
	case TokenIncrement, TokenDecrement:
		if operand.IsNumeric() {
			return operand
		}
		if operand.IsVector() {
			compType := operand.ComponentType()
			if compType.IsNumeric() {
				return operand
			}
		}
	}

	return TypeError
}

// IndexResultType returns the result type of an index operation.
func IndexResultType(baseType *Type) *Type {
	if baseType == nil {
		return TypeError
	}

	switch baseType.Kind {
	case TypeKindArray:
		return baseType.ElementType
	case TypeKindVec2, TypeKindVec3, TypeKindVec4:
		return TypeFloat
	case TypeKindIvec2, TypeKindIvec3, TypeKindIvec4:
		return TypeInt
	case TypeKindUvec2, TypeKindUvec3, TypeKindUvec4:
		return TypeUint
	case TypeKindBvec2, TypeKindBvec3, TypeKindBvec4:
		return TypeBool
	case TypeKindMat2:
		return TypeVec2
	case TypeKindMat3:
		return TypeVec3
	case TypeKindMat4:
		return TypeVec4
	default:
		return TypeError
	}
}

// ConstructorArgCount returns the expected number of arguments for a type constructor.
// Returns -1 if variable number of arguments is allowed.
func ConstructorArgCount(t *Type) int {
	switch t.Kind {
	case TypeKindBool, TypeKindInt, TypeKindUint, TypeKindFloat:
		return 1
	case TypeKindVec2, TypeKindIvec2, TypeKindUvec2, TypeKindBvec2:
		return -1 // 1, 2, or 2 scalars
	case TypeKindVec3, TypeKindIvec3, TypeKindUvec3, TypeKindBvec3:
		return -1 // 1, 3, or 3 scalars, or vec2+scalar
	case TypeKindVec4, TypeKindIvec4, TypeKindUvec4, TypeKindBvec4:
		return -1 // 1, 4, or 4 scalars, or combinations
	case TypeKindMat2:
		return -1 // 1, 4, or mat2
	case TypeKindMat3:
		return -1 // 1, 9, or mat3
	case TypeKindMat4:
		return -1 // 1, 16, or mat4
	default:
		return 0
	}
}

// IsBuiltinTypeName checks if a name is a built-in type.
func IsBuiltinTypeName(name string) bool {
	return TypeFromName(name) != nil
}
