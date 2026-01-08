package gdshader

// BuiltinFunction represents a built-in shader function.
type BuiltinFunction struct {
	Name        string
	Description string
	Signatures  []FunctionSig
}

// FunctionSig represents a function signature.
type FunctionSig struct {
	Params []string // Parameter types
	Return string   // Return type
}

// BuiltinVariable represents a built-in shader variable.
type BuiltinVariable struct {
	Name        string
	Type        string
	Description string
	Stage       string // "vertex", "fragment", "light", or "" for global
	ReadWrite   string // "in", "out", "inout"
}

// BuiltinConstant represents a built-in constant.
type BuiltinConstant struct {
	Name        string
	Type        string
	Description string
	Value       string
}

// ShaderType represents a shader type (spatial, canvas_item, etc.)
type ShaderType string

const (
	ShaderTypeSpatial    ShaderType = "spatial"
	ShaderTypeCanvasItem ShaderType = "canvas_item"
	ShaderTypeParticles  ShaderType = "particles"
	ShaderTypeSky        ShaderType = "sky"
	ShaderTypeFog        ShaderType = "fog"
)

// BuiltinFunctions contains all built-in GLSL-like functions.
var BuiltinFunctions = map[string]*BuiltinFunction{
	// Trigonometric functions
	"radians": {
		Name:        "radians",
		Description: "Converts degrees to radians",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},
	"degrees": {
		Name:        "degrees",
		Description: "Converts radians to degrees",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},
	"sin": {
		Name:        "sin",
		Description: "Returns the sine of the angle",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},
	"cos": {
		Name:        "cos",
		Description: "Returns the cosine of the angle",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},
	"tan": {
		Name:        "tan",
		Description: "Returns the tangent of the angle",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},
	"asin": {
		Name:        "asin",
		Description: "Returns the arc-sine of the parameter",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},
	"acos": {
		Name:        "acos",
		Description: "Returns the arc-cosine of the parameter",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},
	"atan": {
		Name:        "atan",
		Description: "Returns the arc-tangent of the parameter(s)",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"float", "float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},
	"sinh": {
		Name:        "sinh",
		Description: "Returns the hyperbolic sine",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},
	"cosh": {
		Name:        "cosh",
		Description: "Returns the hyperbolic cosine",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},
	"tanh": {
		Name:        "tanh",
		Description: "Returns the hyperbolic tangent",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},
	"asinh": {
		Name:        "asinh",
		Description: "Returns the inverse hyperbolic sine",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},
	"acosh": {
		Name:        "acosh",
		Description: "Returns the inverse hyperbolic cosine",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},
	"atanh": {
		Name:        "atanh",
		Description: "Returns the inverse hyperbolic tangent",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},

	// Exponential functions
	"pow": {
		Name:        "pow",
		Description: "Returns x raised to the power of y",
		Signatures: []FunctionSig{
			{Params: []string{"float", "float"}, Return: "float"},
			{Params: []string{"vec2", "vec2"}, Return: "vec2"},
			{Params: []string{"vec3", "vec3"}, Return: "vec3"},
			{Params: []string{"vec4", "vec4"}, Return: "vec4"},
		},
	},
	"exp": {
		Name:        "exp",
		Description: "Returns e raised to the power of x",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},
	"exp2": {
		Name:        "exp2",
		Description: "Returns 2 raised to the power of x",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},
	"log": {
		Name:        "log",
		Description: "Returns the natural logarithm",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},
	"log2": {
		Name:        "log2",
		Description: "Returns the base-2 logarithm",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},
	"sqrt": {
		Name:        "sqrt",
		Description: "Returns the square root",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},
	"inversesqrt": {
		Name:        "inversesqrt",
		Description: "Returns the inverse square root",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},

	// Common functions
	"abs": {
		Name:        "abs",
		Description: "Returns the absolute value",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
			{Params: []string{"int"}, Return: "int"},
			{Params: []string{"ivec2"}, Return: "ivec2"},
			{Params: []string{"ivec3"}, Return: "ivec3"},
			{Params: []string{"ivec4"}, Return: "ivec4"},
		},
	},
	"sign": {
		Name:        "sign",
		Description: "Returns the sign of the value (-1, 0, or 1)",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
			{Params: []string{"int"}, Return: "int"},
			{Params: []string{"ivec2"}, Return: "ivec2"},
			{Params: []string{"ivec3"}, Return: "ivec3"},
			{Params: []string{"ivec4"}, Return: "ivec4"},
		},
	},
	"floor": {
		Name:        "floor",
		Description: "Returns the largest integer less than or equal to x",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},
	"ceil": {
		Name:        "ceil",
		Description: "Returns the smallest integer greater than or equal to x",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},
	"round": {
		Name:        "round",
		Description: "Returns the nearest integer",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},
	"roundEven": {
		Name:        "roundEven",
		Description: "Returns the nearest even integer",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},
	"trunc": {
		Name:        "trunc",
		Description: "Returns the integer part of x",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},
	"fract": {
		Name:        "fract",
		Description: "Returns the fractional part of x",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},
	"mod": {
		Name:        "mod",
		Description: "Returns x modulo y",
		Signatures: []FunctionSig{
			{Params: []string{"float", "float"}, Return: "float"},
			{Params: []string{"vec2", "float"}, Return: "vec2"},
			{Params: []string{"vec2", "vec2"}, Return: "vec2"},
			{Params: []string{"vec3", "float"}, Return: "vec3"},
			{Params: []string{"vec3", "vec3"}, Return: "vec3"},
			{Params: []string{"vec4", "float"}, Return: "vec4"},
			{Params: []string{"vec4", "vec4"}, Return: "vec4"},
		},
	},
	"min": {
		Name:        "min",
		Description: "Returns the minimum of two values",
		Signatures: []FunctionSig{
			{Params: []string{"float", "float"}, Return: "float"},
			{Params: []string{"vec2", "vec2"}, Return: "vec2"},
			{Params: []string{"vec2", "float"}, Return: "vec2"},
			{Params: []string{"vec3", "vec3"}, Return: "vec3"},
			{Params: []string{"vec3", "float"}, Return: "vec3"},
			{Params: []string{"vec4", "vec4"}, Return: "vec4"},
			{Params: []string{"vec4", "float"}, Return: "vec4"},
			{Params: []string{"int", "int"}, Return: "int"},
			{Params: []string{"uint", "uint"}, Return: "uint"},
		},
	},
	"max": {
		Name:        "max",
		Description: "Returns the maximum of two values",
		Signatures: []FunctionSig{
			{Params: []string{"float", "float"}, Return: "float"},
			{Params: []string{"vec2", "vec2"}, Return: "vec2"},
			{Params: []string{"vec2", "float"}, Return: "vec2"},
			{Params: []string{"vec3", "vec3"}, Return: "vec3"},
			{Params: []string{"vec3", "float"}, Return: "vec3"},
			{Params: []string{"vec4", "vec4"}, Return: "vec4"},
			{Params: []string{"vec4", "float"}, Return: "vec4"},
			{Params: []string{"int", "int"}, Return: "int"},
			{Params: []string{"uint", "uint"}, Return: "uint"},
		},
	},
	"clamp": {
		Name:        "clamp",
		Description: "Clamps x to the range [minVal, maxVal]",
		Signatures: []FunctionSig{
			{Params: []string{"float", "float", "float"}, Return: "float"},
			{Params: []string{"vec2", "vec2", "vec2"}, Return: "vec2"},
			{Params: []string{"vec2", "float", "float"}, Return: "vec2"},
			{Params: []string{"vec3", "vec3", "vec3"}, Return: "vec3"},
			{Params: []string{"vec3", "float", "float"}, Return: "vec3"},
			{Params: []string{"vec4", "vec4", "vec4"}, Return: "vec4"},
			{Params: []string{"vec4", "float", "float"}, Return: "vec4"},
			{Params: []string{"int", "int", "int"}, Return: "int"},
			{Params: []string{"uint", "uint", "uint"}, Return: "uint"},
		},
	},
	"mix": {
		Name:        "mix",
		Description: "Linearly interpolates between x and y",
		Signatures: []FunctionSig{
			{Params: []string{"float", "float", "float"}, Return: "float"},
			{Params: []string{"vec2", "vec2", "float"}, Return: "vec2"},
			{Params: []string{"vec2", "vec2", "vec2"}, Return: "vec2"},
			{Params: []string{"vec3", "vec3", "float"}, Return: "vec3"},
			{Params: []string{"vec3", "vec3", "vec3"}, Return: "vec3"},
			{Params: []string{"vec4", "vec4", "float"}, Return: "vec4"},
			{Params: []string{"vec4", "vec4", "vec4"}, Return: "vec4"},
		},
	},
	"step": {
		Name:        "step",
		Description: "Returns 0.0 if x < edge, otherwise 1.0",
		Signatures: []FunctionSig{
			{Params: []string{"float", "float"}, Return: "float"},
			{Params: []string{"float", "vec2"}, Return: "vec2"},
			{Params: []string{"float", "vec3"}, Return: "vec3"},
			{Params: []string{"float", "vec4"}, Return: "vec4"},
			{Params: []string{"vec2", "vec2"}, Return: "vec2"},
			{Params: []string{"vec3", "vec3"}, Return: "vec3"},
			{Params: []string{"vec4", "vec4"}, Return: "vec4"},
		},
	},
	"smoothstep": {
		Name:        "smoothstep",
		Description: "Performs smooth Hermite interpolation",
		Signatures: []FunctionSig{
			{Params: []string{"float", "float", "float"}, Return: "float"},
			{Params: []string{"float", "float", "vec2"}, Return: "vec2"},
			{Params: []string{"float", "float", "vec3"}, Return: "vec3"},
			{Params: []string{"float", "float", "vec4"}, Return: "vec4"},
			{Params: []string{"vec2", "vec2", "vec2"}, Return: "vec2"},
			{Params: []string{"vec3", "vec3", "vec3"}, Return: "vec3"},
			{Params: []string{"vec4", "vec4", "vec4"}, Return: "vec4"},
		},
	},

	// Geometric functions
	"length": {
		Name:        "length",
		Description: "Returns the length of a vector",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "float"},
			{Params: []string{"vec3"}, Return: "float"},
			{Params: []string{"vec4"}, Return: "float"},
		},
	},
	"distance": {
		Name:        "distance",
		Description: "Returns the distance between two points",
		Signatures: []FunctionSig{
			{Params: []string{"float", "float"}, Return: "float"},
			{Params: []string{"vec2", "vec2"}, Return: "float"},
			{Params: []string{"vec3", "vec3"}, Return: "float"},
			{Params: []string{"vec4", "vec4"}, Return: "float"},
		},
	},
	"dot": {
		Name:        "dot",
		Description: "Returns the dot product of two vectors",
		Signatures: []FunctionSig{
			{Params: []string{"float", "float"}, Return: "float"},
			{Params: []string{"vec2", "vec2"}, Return: "float"},
			{Params: []string{"vec3", "vec3"}, Return: "float"},
			{Params: []string{"vec4", "vec4"}, Return: "float"},
		},
	},
	"cross": {
		Name:        "cross",
		Description: "Returns the cross product of two vectors",
		Signatures: []FunctionSig{
			{Params: []string{"vec3", "vec3"}, Return: "vec3"},
		},
	},
	"normalize": {
		Name:        "normalize",
		Description: "Returns a normalized vector",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},
	"reflect": {
		Name:        "reflect",
		Description: "Reflects a vector about a normal",
		Signatures: []FunctionSig{
			{Params: []string{"float", "float"}, Return: "float"},
			{Params: []string{"vec2", "vec2"}, Return: "vec2"},
			{Params: []string{"vec3", "vec3"}, Return: "vec3"},
			{Params: []string{"vec4", "vec4"}, Return: "vec4"},
		},
	},
	"refract": {
		Name:        "refract",
		Description: "Refracts a vector through a surface",
		Signatures: []FunctionSig{
			{Params: []string{"float", "float", "float"}, Return: "float"},
			{Params: []string{"vec2", "vec2", "float"}, Return: "vec2"},
			{Params: []string{"vec3", "vec3", "float"}, Return: "vec3"},
			{Params: []string{"vec4", "vec4", "float"}, Return: "vec4"},
		},
	},
	"faceforward": {
		Name:        "faceforward",
		Description: "Returns N if dot(Nref, I) < 0, otherwise -N",
		Signatures: []FunctionSig{
			{Params: []string{"float", "float", "float"}, Return: "float"},
			{Params: []string{"vec2", "vec2", "vec2"}, Return: "vec2"},
			{Params: []string{"vec3", "vec3", "vec3"}, Return: "vec3"},
			{Params: []string{"vec4", "vec4", "vec4"}, Return: "vec4"},
		},
	},

	// Matrix functions
	"matrixCompMult": {
		Name:        "matrixCompMult",
		Description: "Component-wise matrix multiplication",
		Signatures: []FunctionSig{
			{Params: []string{"mat2", "mat2"}, Return: "mat2"},
			{Params: []string{"mat3", "mat3"}, Return: "mat3"},
			{Params: []string{"mat4", "mat4"}, Return: "mat4"},
		},
	},
	"transpose": {
		Name:        "transpose",
		Description: "Returns the transpose of a matrix",
		Signatures: []FunctionSig{
			{Params: []string{"mat2"}, Return: "mat2"},
			{Params: []string{"mat3"}, Return: "mat3"},
			{Params: []string{"mat4"}, Return: "mat4"},
		},
	},
	"inverse": {
		Name:        "inverse",
		Description: "Returns the inverse of a matrix",
		Signatures: []FunctionSig{
			{Params: []string{"mat2"}, Return: "mat2"},
			{Params: []string{"mat3"}, Return: "mat3"},
			{Params: []string{"mat4"}, Return: "mat4"},
		},
	},
	"determinant": {
		Name:        "determinant",
		Description: "Returns the determinant of a matrix",
		Signatures: []FunctionSig{
			{Params: []string{"mat2"}, Return: "float"},
			{Params: []string{"mat3"}, Return: "float"},
			{Params: []string{"mat4"}, Return: "float"},
		},
	},
	"outerProduct": {
		Name:        "outerProduct",
		Description: "Returns the outer product of two vectors",
		Signatures: []FunctionSig{
			{Params: []string{"vec2", "vec2"}, Return: "mat2"},
			{Params: []string{"vec3", "vec3"}, Return: "mat3"},
			{Params: []string{"vec4", "vec4"}, Return: "mat4"},
		},
	},

	// Texture functions
	"texture": {
		Name:        "texture",
		Description: "Samples a texture",
		Signatures: []FunctionSig{
			{Params: []string{"sampler2D", "vec2"}, Return: "vec4"},
			{Params: []string{"sampler2D", "vec2", "float"}, Return: "vec4"},
			{Params: []string{"sampler2DArray", "vec3"}, Return: "vec4"},
			{Params: []string{"sampler3D", "vec3"}, Return: "vec4"},
			{Params: []string{"samplerCube", "vec3"}, Return: "vec4"},
		},
	},
	"textureSize": {
		Name:        "textureSize",
		Description: "Returns the size of a texture",
		Signatures: []FunctionSig{
			{Params: []string{"sampler2D", "int"}, Return: "ivec2"},
			{Params: []string{"sampler2DArray", "int"}, Return: "ivec3"},
			{Params: []string{"sampler3D", "int"}, Return: "ivec3"},
			{Params: []string{"samplerCube", "int"}, Return: "ivec2"},
		},
	},
	"textureLod": {
		Name:        "textureLod",
		Description: "Samples a texture with explicit LOD",
		Signatures: []FunctionSig{
			{Params: []string{"sampler2D", "vec2", "float"}, Return: "vec4"},
			{Params: []string{"sampler2DArray", "vec3", "float"}, Return: "vec4"},
			{Params: []string{"sampler3D", "vec3", "float"}, Return: "vec4"},
			{Params: []string{"samplerCube", "vec3", "float"}, Return: "vec4"},
		},
	},
	"textureProj": {
		Name:        "textureProj",
		Description: "Samples a texture with projection",
		Signatures: []FunctionSig{
			{Params: []string{"sampler2D", "vec3"}, Return: "vec4"},
			{Params: []string{"sampler2D", "vec4"}, Return: "vec4"},
		},
	},
	"texelFetch": {
		Name:        "texelFetch",
		Description: "Fetches a single texel",
		Signatures: []FunctionSig{
			{Params: []string{"sampler2D", "ivec2", "int"}, Return: "vec4"},
			{Params: []string{"sampler2DArray", "ivec3", "int"}, Return: "vec4"},
			{Params: []string{"sampler3D", "ivec3", "int"}, Return: "vec4"},
		},
	},
	"textureGrad": {
		Name:        "textureGrad",
		Description: "Samples a texture with explicit gradients",
		Signatures: []FunctionSig{
			{Params: []string{"sampler2D", "vec2", "vec2", "vec2"}, Return: "vec4"},
			{Params: []string{"sampler3D", "vec3", "vec3", "vec3"}, Return: "vec4"},
			{Params: []string{"samplerCube", "vec3", "vec3", "vec3"}, Return: "vec4"},
		},
	},

	// Derivative functions
	"dFdx": {
		Name:        "dFdx",
		Description: "Returns the partial derivative with respect to x",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},
	"dFdy": {
		Name:        "dFdy",
		Description: "Returns the partial derivative with respect to y",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},
	"fwidth": {
		Name:        "fwidth",
		Description: "Returns abs(dFdx) + abs(dFdy)",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "float"},
			{Params: []string{"vec2"}, Return: "vec2"},
			{Params: []string{"vec3"}, Return: "vec3"},
			{Params: []string{"vec4"}, Return: "vec4"},
		},
	},

	// Boolean functions
	"lessThan": {
		Name:        "lessThan",
		Description: "Component-wise less than comparison",
		Signatures: []FunctionSig{
			{Params: []string{"vec2", "vec2"}, Return: "bvec2"},
			{Params: []string{"vec3", "vec3"}, Return: "bvec3"},
			{Params: []string{"vec4", "vec4"}, Return: "bvec4"},
			{Params: []string{"ivec2", "ivec2"}, Return: "bvec2"},
			{Params: []string{"ivec3", "ivec3"}, Return: "bvec3"},
			{Params: []string{"ivec4", "ivec4"}, Return: "bvec4"},
		},
	},
	"greaterThan": {
		Name:        "greaterThan",
		Description: "Component-wise greater than comparison",
		Signatures: []FunctionSig{
			{Params: []string{"vec2", "vec2"}, Return: "bvec2"},
			{Params: []string{"vec3", "vec3"}, Return: "bvec3"},
			{Params: []string{"vec4", "vec4"}, Return: "bvec4"},
			{Params: []string{"ivec2", "ivec2"}, Return: "bvec2"},
			{Params: []string{"ivec3", "ivec3"}, Return: "bvec3"},
			{Params: []string{"ivec4", "ivec4"}, Return: "bvec4"},
		},
	},
	"equal": {
		Name:        "equal",
		Description: "Component-wise equality comparison",
		Signatures: []FunctionSig{
			{Params: []string{"vec2", "vec2"}, Return: "bvec2"},
			{Params: []string{"vec3", "vec3"}, Return: "bvec3"},
			{Params: []string{"vec4", "vec4"}, Return: "bvec4"},
			{Params: []string{"ivec2", "ivec2"}, Return: "bvec2"},
			{Params: []string{"ivec3", "ivec3"}, Return: "bvec3"},
			{Params: []string{"ivec4", "ivec4"}, Return: "bvec4"},
			{Params: []string{"bvec2", "bvec2"}, Return: "bvec2"},
			{Params: []string{"bvec3", "bvec3"}, Return: "bvec3"},
			{Params: []string{"bvec4", "bvec4"}, Return: "bvec4"},
		},
	},
	"notEqual": {
		Name:        "notEqual",
		Description: "Component-wise inequality comparison",
		Signatures: []FunctionSig{
			{Params: []string{"vec2", "vec2"}, Return: "bvec2"},
			{Params: []string{"vec3", "vec3"}, Return: "bvec3"},
			{Params: []string{"vec4", "vec4"}, Return: "bvec4"},
			{Params: []string{"ivec2", "ivec2"}, Return: "bvec2"},
			{Params: []string{"ivec3", "ivec3"}, Return: "bvec3"},
			{Params: []string{"ivec4", "ivec4"}, Return: "bvec4"},
			{Params: []string{"bvec2", "bvec2"}, Return: "bvec2"},
			{Params: []string{"bvec3", "bvec3"}, Return: "bvec3"},
			{Params: []string{"bvec4", "bvec4"}, Return: "bvec4"},
		},
	},
	"any": {
		Name:        "any",
		Description: "Returns true if any component is true",
		Signatures: []FunctionSig{
			{Params: []string{"bvec2"}, Return: "bool"},
			{Params: []string{"bvec3"}, Return: "bool"},
			{Params: []string{"bvec4"}, Return: "bool"},
		},
	},
	"all": {
		Name:        "all",
		Description: "Returns true if all components are true",
		Signatures: []FunctionSig{
			{Params: []string{"bvec2"}, Return: "bool"},
			{Params: []string{"bvec3"}, Return: "bool"},
			{Params: []string{"bvec4"}, Return: "bool"},
		},
	},
	"not": {
		Name:        "not",
		Description: "Component-wise logical NOT",
		Signatures: []FunctionSig{
			{Params: []string{"bvec2"}, Return: "bvec2"},
			{Params: []string{"bvec3"}, Return: "bvec3"},
			{Params: []string{"bvec4"}, Return: "bvec4"},
		},
	},

	// Integer functions
	"floatBitsToInt": {
		Name:        "floatBitsToInt",
		Description: "Reinterprets float bits as int",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "int"},
			{Params: []string{"vec2"}, Return: "ivec2"},
			{Params: []string{"vec3"}, Return: "ivec3"},
			{Params: []string{"vec4"}, Return: "ivec4"},
		},
	},
	"floatBitsToUint": {
		Name:        "floatBitsToUint",
		Description: "Reinterprets float bits as uint",
		Signatures: []FunctionSig{
			{Params: []string{"float"}, Return: "uint"},
			{Params: []string{"vec2"}, Return: "uvec2"},
			{Params: []string{"vec3"}, Return: "uvec3"},
			{Params: []string{"vec4"}, Return: "uvec4"},
		},
	},
	"intBitsToFloat": {
		Name:        "intBitsToFloat",
		Description: "Reinterprets int bits as float",
		Signatures: []FunctionSig{
			{Params: []string{"int"}, Return: "float"},
			{Params: []string{"ivec2"}, Return: "vec2"},
			{Params: []string{"ivec3"}, Return: "vec3"},
			{Params: []string{"ivec4"}, Return: "vec4"},
		},
	},
	"uintBitsToFloat": {
		Name:        "uintBitsToFloat",
		Description: "Reinterprets uint bits as float",
		Signatures: []FunctionSig{
			{Params: []string{"uint"}, Return: "float"},
			{Params: []string{"uvec2"}, Return: "vec2"},
			{Params: []string{"uvec3"}, Return: "vec3"},
			{Params: []string{"uvec4"}, Return: "vec4"},
		},
	},

	// Pack/unpack functions
	"packHalf2x16": {
		Name:        "packHalf2x16",
		Description: "Packs two floats into a uint",
		Signatures:  []FunctionSig{{Params: []string{"vec2"}, Return: "uint"}},
	},
	"unpackHalf2x16": {
		Name:        "unpackHalf2x16",
		Description: "Unpacks a uint into two floats",
		Signatures:  []FunctionSig{{Params: []string{"uint"}, Return: "vec2"}},
	},
	"packUnorm2x16": {
		Name:        "packUnorm2x16",
		Description: "Packs two normalized floats into a uint",
		Signatures:  []FunctionSig{{Params: []string{"vec2"}, Return: "uint"}},
	},
	"unpackUnorm2x16": {
		Name:        "unpackUnorm2x16",
		Description: "Unpacks a uint into two normalized floats",
		Signatures:  []FunctionSig{{Params: []string{"uint"}, Return: "vec2"}},
	},
	"packSnorm2x16": {
		Name:        "packSnorm2x16",
		Description: "Packs two signed normalized floats into a uint",
		Signatures:  []FunctionSig{{Params: []string{"vec2"}, Return: "uint"}},
	},
	"unpackSnorm2x16": {
		Name:        "unpackSnorm2x16",
		Description: "Unpacks a uint into two signed normalized floats",
		Signatures:  []FunctionSig{{Params: []string{"uint"}, Return: "vec2"}},
	},
}

// BuiltinConstants contains all built-in constants.
var BuiltinConstants = map[string]*BuiltinConstant{
	"PI": {
		Name:        "PI",
		Type:        "float",
		Description: "The mathematical constant pi (3.14159...)",
		Value:       "3.14159265358979323846",
	},
	"TAU": {
		Name:        "TAU",
		Type:        "float",
		Description: "The mathematical constant tau (2 * pi)",
		Value:       "6.28318530717958647692",
	},
	"E": {
		Name:        "E",
		Type:        "float",
		Description: "Euler's number (2.71828...)",
		Value:       "2.71828182845904523536",
	},
}

// GetSpatialVertexBuiltins returns built-in variables for spatial shader vertex stage.
func GetSpatialVertexBuiltins() map[string]*BuiltinVariable {
	return map[string]*BuiltinVariable{
		// Inputs
		"VERTEX":          {Name: "VERTEX", Type: "vec3", Description: "Vertex position in local space", Stage: "vertex", ReadWrite: "inout"},
		"NORMAL":          {Name: "NORMAL", Type: "vec3", Description: "Vertex normal in local space", Stage: "vertex", ReadWrite: "inout"},
		"TANGENT":         {Name: "TANGENT", Type: "vec3", Description: "Vertex tangent in local space", Stage: "vertex", ReadWrite: "inout"},
		"BINORMAL":        {Name: "BINORMAL", Type: "vec3", Description: "Vertex binormal in local space", Stage: "vertex", ReadWrite: "inout"},
		"UV":              {Name: "UV", Type: "vec2", Description: "Primary UV coordinates", Stage: "vertex", ReadWrite: "inout"},
		"UV2":             {Name: "UV2", Type: "vec2", Description: "Secondary UV coordinates", Stage: "vertex", ReadWrite: "inout"},
		"COLOR":           {Name: "COLOR", Type: "vec4", Description: "Vertex color", Stage: "vertex", ReadWrite: "inout"},
		"POINT_SIZE":      {Name: "POINT_SIZE", Type: "float", Description: "Point size for point rendering", Stage: "vertex", ReadWrite: "inout"},
		"INSTANCE_ID":     {Name: "INSTANCE_ID", Type: "int", Description: "Instance ID for instanced rendering", Stage: "vertex", ReadWrite: "in"},
		"VERTEX_ID":       {Name: "VERTEX_ID", Type: "int", Description: "Vertex ID", Stage: "vertex", ReadWrite: "in"},
		"INSTANCE_CUSTOM": {Name: "INSTANCE_CUSTOM", Type: "vec4", Description: "Instance custom data", Stage: "vertex", ReadWrite: "in"},
		// Matrices
		"MODEL_MATRIX":            {Name: "MODEL_MATRIX", Type: "mat4", Description: "Model matrix (world transform)", Stage: "vertex", ReadWrite: "in"},
		"MODEL_NORMAL_MATRIX":     {Name: "MODEL_NORMAL_MATRIX", Type: "mat3", Description: "Normal matrix", Stage: "vertex", ReadWrite: "in"},
		"VIEW_MATRIX":             {Name: "VIEW_MATRIX", Type: "mat4", Description: "View matrix", Stage: "vertex", ReadWrite: "in"},
		"INV_VIEW_MATRIX":         {Name: "INV_VIEW_MATRIX", Type: "mat4", Description: "Inverse view matrix", Stage: "vertex", ReadWrite: "in"},
		"PROJECTION_MATRIX":       {Name: "PROJECTION_MATRIX", Type: "mat4", Description: "Projection matrix", Stage: "vertex", ReadWrite: "inout"},
		"INV_PROJECTION_MATRIX":   {Name: "INV_PROJECTION_MATRIX", Type: "mat4", Description: "Inverse projection matrix", Stage: "vertex", ReadWrite: "in"},
		"MODELVIEW_MATRIX":        {Name: "MODELVIEW_MATRIX", Type: "mat4", Description: "Model-view matrix", Stage: "vertex", ReadWrite: "in"},
		"MODELVIEW_NORMAL_MATRIX": {Name: "MODELVIEW_NORMAL_MATRIX", Type: "mat3", Description: "Model-view normal matrix", Stage: "vertex", ReadWrite: "in"},
		// Camera
		"VIEWPORT_SIZE":          {Name: "VIEWPORT_SIZE", Type: "vec2", Description: "Viewport size in pixels", Stage: "vertex", ReadWrite: "in"},
		"OUTPUT_IS_SRGB":         {Name: "OUTPUT_IS_SRGB", Type: "bool", Description: "True if output is sRGB", Stage: "vertex", ReadWrite: "in"},
		"NODE_POSITION_WORLD":    {Name: "NODE_POSITION_WORLD", Type: "vec3", Description: "Node position in world space", Stage: "vertex", ReadWrite: "in"},
		"CAMERA_POSITION_WORLD":  {Name: "CAMERA_POSITION_WORLD", Type: "vec3", Description: "Camera position in world space", Stage: "vertex", ReadWrite: "in"},
		"CAMERA_DIRECTION_WORLD": {Name: "CAMERA_DIRECTION_WORLD", Type: "vec3", Description: "Camera direction in world space", Stage: "vertex", ReadWrite: "in"},
		"CAMERA_VISIBLE_LAYERS":  {Name: "CAMERA_VISIBLE_LAYERS", Type: "uint", Description: "Camera visible layers bitmask", Stage: "vertex", ReadWrite: "in"},
		// Outputs
		"POSITION": {Name: "POSITION", Type: "vec4", Description: "Output position in clip space", Stage: "vertex", ReadWrite: "out"},
		// Time
		"TIME": {Name: "TIME", Type: "float", Description: "Time since start", Stage: "vertex", ReadWrite: "in"},
	}
}

// GetSpatialFragmentBuiltins returns built-in variables for spatial shader fragment stage.
func GetSpatialFragmentBuiltins() map[string]*BuiltinVariable {
	return map[string]*BuiltinVariable{
		// Inputs from vertex
		"VERTEX":       {Name: "VERTEX", Type: "vec3", Description: "Vertex position in view space", Stage: "fragment", ReadWrite: "in"},
		"FRAGCOORD":    {Name: "FRAGCOORD", Type: "vec4", Description: "Fragment coordinates", Stage: "fragment", ReadWrite: "in"},
		"FRONT_FACING": {Name: "FRONT_FACING", Type: "bool", Description: "True if front face", Stage: "fragment", ReadWrite: "in"},
		"NORMAL":       {Name: "NORMAL", Type: "vec3", Description: "Normal in view space", Stage: "fragment", ReadWrite: "inout"},
		"TANGENT":      {Name: "TANGENT", Type: "vec3", Description: "Tangent in view space", Stage: "fragment", ReadWrite: "in"},
		"BINORMAL":     {Name: "BINORMAL", Type: "vec3", Description: "Binormal in view space", Stage: "fragment", ReadWrite: "in"},
		"UV":           {Name: "UV", Type: "vec2", Description: "Primary UV coordinates", Stage: "fragment", ReadWrite: "in"},
		"UV2":          {Name: "UV2", Type: "vec2", Description: "Secondary UV coordinates", Stage: "fragment", ReadWrite: "in"},
		"COLOR":        {Name: "COLOR", Type: "vec4", Description: "Vertex color", Stage: "fragment", ReadWrite: "in"},
		// PBR outputs
		"ALBEDO":                   {Name: "ALBEDO", Type: "vec3", Description: "Albedo color", Stage: "fragment", ReadWrite: "out"},
		"ALPHA":                    {Name: "ALPHA", Type: "float", Description: "Alpha value", Stage: "fragment", ReadWrite: "out"},
		"METALLIC":                 {Name: "METALLIC", Type: "float", Description: "Metallic value", Stage: "fragment", ReadWrite: "out"},
		"ROUGHNESS":                {Name: "ROUGHNESS", Type: "float", Description: "Roughness value", Stage: "fragment", ReadWrite: "out"},
		"SPECULAR":                 {Name: "SPECULAR", Type: "float", Description: "Specular value", Stage: "fragment", ReadWrite: "out"},
		"RIM":                      {Name: "RIM", Type: "float", Description: "Rim lighting intensity", Stage: "fragment", ReadWrite: "out"},
		"RIM_TINT":                 {Name: "RIM_TINT", Type: "float", Description: "Rim tint", Stage: "fragment", ReadWrite: "out"},
		"CLEARCOAT":                {Name: "CLEARCOAT", Type: "float", Description: "Clearcoat intensity", Stage: "fragment", ReadWrite: "out"},
		"CLEARCOAT_ROUGHNESS":      {Name: "CLEARCOAT_ROUGHNESS", Type: "float", Description: "Clearcoat roughness", Stage: "fragment", ReadWrite: "out"},
		"ANISOTROPY":               {Name: "ANISOTROPY", Type: "float", Description: "Anisotropy intensity", Stage: "fragment", ReadWrite: "out"},
		"ANISOTROPY_FLOW":          {Name: "ANISOTROPY_FLOW", Type: "vec2", Description: "Anisotropy flow direction", Stage: "fragment", ReadWrite: "out"},
		"SSS_STRENGTH":             {Name: "SSS_STRENGTH", Type: "float", Description: "Subsurface scattering strength", Stage: "fragment", ReadWrite: "out"},
		"SSS_TRANSMITTANCE_COLOR":  {Name: "SSS_TRANSMITTANCE_COLOR", Type: "vec4", Description: "SSS transmittance color", Stage: "fragment", ReadWrite: "out"},
		"SSS_TRANSMITTANCE_DEPTH":  {Name: "SSS_TRANSMITTANCE_DEPTH", Type: "float", Description: "SSS transmittance depth", Stage: "fragment", ReadWrite: "out"},
		"SSS_TRANSMITTANCE_BOOST":  {Name: "SSS_TRANSMITTANCE_BOOST", Type: "float", Description: "SSS transmittance boost", Stage: "fragment", ReadWrite: "out"},
		"BACKLIGHT":                {Name: "BACKLIGHT", Type: "vec3", Description: "Backlight color", Stage: "fragment", ReadWrite: "out"},
		"AO":                       {Name: "AO", Type: "float", Description: "Ambient occlusion", Stage: "fragment", ReadWrite: "out"},
		"AO_LIGHT_AFFECT":          {Name: "AO_LIGHT_AFFECT", Type: "float", Description: "AO light affect", Stage: "fragment", ReadWrite: "out"},
		"EMISSION":                 {Name: "EMISSION", Type: "vec3", Description: "Emission color", Stage: "fragment", ReadWrite: "out"},
		"NORMAL_MAP":               {Name: "NORMAL_MAP", Type: "vec3", Description: "Normal map", Stage: "fragment", ReadWrite: "out"},
		"NORMAL_MAP_DEPTH":         {Name: "NORMAL_MAP_DEPTH", Type: "float", Description: "Normal map depth", Stage: "fragment", ReadWrite: "out"},
		"ALPHA_SCISSOR_THRESHOLD":  {Name: "ALPHA_SCISSOR_THRESHOLD", Type: "float", Description: "Alpha scissor threshold", Stage: "fragment", ReadWrite: "out"},
		"ALPHA_HASH_SCALE":         {Name: "ALPHA_HASH_SCALE", Type: "float", Description: "Alpha hash scale", Stage: "fragment", ReadWrite: "out"},
		"ALPHA_ANTIALIASING_EDGE":  {Name: "ALPHA_ANTIALIASING_EDGE", Type: "float", Description: "Alpha antialiasing edge", Stage: "fragment", ReadWrite: "out"},
		"ALPHA_TEXTURE_COORDINATE": {Name: "ALPHA_TEXTURE_COORDINATE", Type: "vec2", Description: "Alpha texture coordinate", Stage: "fragment", ReadWrite: "out"},
		"FOG":                      {Name: "FOG", Type: "vec4", Description: "Fog color and density", Stage: "fragment", ReadWrite: "out"},
		// Matrices
		"MODEL_MATRIX":          {Name: "MODEL_MATRIX", Type: "mat4", Description: "Model matrix", Stage: "fragment", ReadWrite: "in"},
		"MODEL_NORMAL_MATRIX":   {Name: "MODEL_NORMAL_MATRIX", Type: "mat3", Description: "Model normal matrix", Stage: "fragment", ReadWrite: "in"},
		"VIEW_MATRIX":           {Name: "VIEW_MATRIX", Type: "mat4", Description: "View matrix", Stage: "fragment", ReadWrite: "in"},
		"INV_VIEW_MATRIX":       {Name: "INV_VIEW_MATRIX", Type: "mat4", Description: "Inverse view matrix", Stage: "fragment", ReadWrite: "in"},
		"PROJECTION_MATRIX":     {Name: "PROJECTION_MATRIX", Type: "mat4", Description: "Projection matrix", Stage: "fragment", ReadWrite: "in"},
		"INV_PROJECTION_MATRIX": {Name: "INV_PROJECTION_MATRIX", Type: "mat4", Description: "Inverse projection matrix", Stage: "fragment", ReadWrite: "in"},
		// Camera
		"VIEWPORT_SIZE":          {Name: "VIEWPORT_SIZE", Type: "vec2", Description: "Viewport size", Stage: "fragment", ReadWrite: "in"},
		"NODE_POSITION_WORLD":    {Name: "NODE_POSITION_WORLD", Type: "vec3", Description: "Node world position", Stage: "fragment", ReadWrite: "in"},
		"CAMERA_POSITION_WORLD":  {Name: "CAMERA_POSITION_WORLD", Type: "vec3", Description: "Camera world position", Stage: "fragment", ReadWrite: "in"},
		"CAMERA_DIRECTION_WORLD": {Name: "CAMERA_DIRECTION_WORLD", Type: "vec3", Description: "Camera world direction", Stage: "fragment", ReadWrite: "in"},
		"CAMERA_VISIBLE_LAYERS":  {Name: "CAMERA_VISIBLE_LAYERS", Type: "uint", Description: "Camera visible layers", Stage: "fragment", ReadWrite: "in"},
		"VIEW":                   {Name: "VIEW", Type: "vec3", Description: "View direction", Stage: "fragment", ReadWrite: "in"},
		// Time
		"TIME": {Name: "TIME", Type: "float", Description: "Time since start", Stage: "fragment", ReadWrite: "in"},
		// Screen
		"SCREEN_UV":         {Name: "SCREEN_UV", Type: "vec2", Description: "Screen UV coordinates", Stage: "fragment", ReadWrite: "in"},
		"SCREEN_PIXEL_SIZE": {Name: "SCREEN_PIXEL_SIZE", Type: "vec2", Description: "Screen pixel size", Stage: "fragment", ReadWrite: "in"},
		// Depth
		"DEPTH": {Name: "DEPTH", Type: "float", Description: "Output depth", Stage: "fragment", ReadWrite: "out"},
	}
}

// GetSpatialLightBuiltins returns built-in variables for spatial shader light stage.
func GetSpatialLightBuiltins() map[string]*BuiltinVariable {
	return map[string]*BuiltinVariable{
		// From fragment
		"ALBEDO":    {Name: "ALBEDO", Type: "vec3", Description: "Albedo from fragment", Stage: "light", ReadWrite: "in"},
		"ROUGHNESS": {Name: "ROUGHNESS", Type: "float", Description: "Roughness from fragment", Stage: "light", ReadWrite: "in"},
		"METALLIC":  {Name: "METALLIC", Type: "float", Description: "Metallic from fragment", Stage: "light", ReadWrite: "in"},
		"SPECULAR":  {Name: "SPECULAR", Type: "float", Description: "Specular from fragment", Stage: "light", ReadWrite: "in"},
		"BACKLIGHT": {Name: "BACKLIGHT", Type: "vec3", Description: "Backlight from fragment", Stage: "light", ReadWrite: "in"},
		"AO":        {Name: "AO", Type: "float", Description: "AO from fragment", Stage: "light", ReadWrite: "in"},
		// Light info
		"LIGHT":                {Name: "LIGHT", Type: "vec3", Description: "Light direction", Stage: "light", ReadWrite: "in"},
		"LIGHT_COLOR":          {Name: "LIGHT_COLOR", Type: "vec3", Description: "Light color", Stage: "light", ReadWrite: "in"},
		"ATTENUATION":          {Name: "ATTENUATION", Type: "float", Description: "Light attenuation", Stage: "light", ReadWrite: "in"},
		"SHADOW_ATTENUATION":   {Name: "SHADOW_ATTENUATION", Type: "vec3", Description: "Shadow attenuation", Stage: "light", ReadWrite: "in"},
		"LIGHT_IS_DIRECTIONAL": {Name: "LIGHT_IS_DIRECTIONAL", Type: "bool", Description: "Is directional light", Stage: "light", ReadWrite: "in"},
		// View info
		"VIEW":   {Name: "VIEW", Type: "vec3", Description: "View direction", Stage: "light", ReadWrite: "in"},
		"NORMAL": {Name: "NORMAL", Type: "vec3", Description: "Normal in view space", Stage: "light", ReadWrite: "in"},
		// Outputs
		"DIFFUSE_LIGHT":  {Name: "DIFFUSE_LIGHT", Type: "vec3", Description: "Diffuse light output", Stage: "light", ReadWrite: "out"},
		"SPECULAR_LIGHT": {Name: "SPECULAR_LIGHT", Type: "vec3", Description: "Specular light output", Stage: "light", ReadWrite: "out"},
		"ALPHA":          {Name: "ALPHA", Type: "float", Description: "Alpha output", Stage: "light", ReadWrite: "out"},
	}
}

// GetCanvasItemVertexBuiltins returns built-in variables for canvas_item shader vertex stage.
func GetCanvasItemVertexBuiltins() map[string]*BuiltinVariable {
	return map[string]*BuiltinVariable{
		"VERTEX":             {Name: "VERTEX", Type: "vec2", Description: "Vertex position", Stage: "vertex", ReadWrite: "inout"},
		"UV":                 {Name: "UV", Type: "vec2", Description: "UV coordinates", Stage: "vertex", ReadWrite: "inout"},
		"COLOR":              {Name: "COLOR", Type: "vec4", Description: "Vertex color", Stage: "vertex", ReadWrite: "inout"},
		"POINT_SIZE":         {Name: "POINT_SIZE", Type: "float", Description: "Point size", Stage: "vertex", ReadWrite: "inout"},
		"MODEL_MATRIX":       {Name: "MODEL_MATRIX", Type: "mat4", Description: "Model matrix", Stage: "vertex", ReadWrite: "in"},
		"CANVAS_MATRIX":      {Name: "CANVAS_MATRIX", Type: "mat4", Description: "Canvas matrix", Stage: "vertex", ReadWrite: "in"},
		"SCREEN_MATRIX":      {Name: "SCREEN_MATRIX", Type: "mat4", Description: "Screen matrix", Stage: "vertex", ReadWrite: "in"},
		"INSTANCE_CUSTOM":    {Name: "INSTANCE_CUSTOM", Type: "vec4", Description: "Instance custom data", Stage: "vertex", ReadWrite: "in"},
		"INSTANCE_ID":        {Name: "INSTANCE_ID", Type: "int", Description: "Instance ID", Stage: "vertex", ReadWrite: "in"},
		"VERTEX_ID":          {Name: "VERTEX_ID", Type: "int", Description: "Vertex ID", Stage: "vertex", ReadWrite: "in"},
		"AT_LIGHT_PASS":      {Name: "AT_LIGHT_PASS", Type: "bool", Description: "Is light pass", Stage: "vertex", ReadWrite: "in"},
		"TEXTURE_PIXEL_SIZE": {Name: "TEXTURE_PIXEL_SIZE", Type: "vec2", Description: "Texture pixel size", Stage: "vertex", ReadWrite: "in"},
		"TIME":               {Name: "TIME", Type: "float", Description: "Time", Stage: "vertex", ReadWrite: "in"},
	}
}

// GetCanvasItemFragmentBuiltins returns built-in variables for canvas_item shader fragment stage.
func GetCanvasItemFragmentBuiltins() map[string]*BuiltinVariable {
	return map[string]*BuiltinVariable{
		"FRAGCOORD":          {Name: "FRAGCOORD", Type: "vec4", Description: "Fragment coordinates", Stage: "fragment", ReadWrite: "in"},
		"UV":                 {Name: "UV", Type: "vec2", Description: "UV coordinates", Stage: "fragment", ReadWrite: "in"},
		"COLOR":              {Name: "COLOR", Type: "vec4", Description: "Color output", Stage: "fragment", ReadWrite: "inout"},
		"NORMAL":             {Name: "NORMAL", Type: "vec3", Description: "Normal for 2D lighting", Stage: "fragment", ReadWrite: "out"},
		"NORMAL_MAP":         {Name: "NORMAL_MAP", Type: "vec3", Description: "Normal map", Stage: "fragment", ReadWrite: "out"},
		"NORMAL_MAP_DEPTH":   {Name: "NORMAL_MAP_DEPTH", Type: "float", Description: "Normal map depth", Stage: "fragment", ReadWrite: "out"},
		"TEXTURE":            {Name: "TEXTURE", Type: "sampler2D", Description: "Main texture", Stage: "fragment", ReadWrite: "in"},
		"TEXTURE_PIXEL_SIZE": {Name: "TEXTURE_PIXEL_SIZE", Type: "vec2", Description: "Texture pixel size", Stage: "fragment", ReadWrite: "in"},
		"SCREEN_UV":          {Name: "SCREEN_UV", Type: "vec2", Description: "Screen UV", Stage: "fragment", ReadWrite: "in"},
		"SCREEN_PIXEL_SIZE":  {Name: "SCREEN_PIXEL_SIZE", Type: "vec2", Description: "Screen pixel size", Stage: "fragment", ReadWrite: "in"},
		"POINT_COORD":        {Name: "POINT_COORD", Type: "vec2", Description: "Point coordinate", Stage: "fragment", ReadWrite: "in"},
		"AT_LIGHT_PASS":      {Name: "AT_LIGHT_PASS", Type: "bool", Description: "Is light pass", Stage: "fragment", ReadWrite: "in"},
		"TIME":               {Name: "TIME", Type: "float", Description: "Time", Stage: "fragment", ReadWrite: "in"},
		"SPECULAR_SHININESS": {Name: "SPECULAR_SHININESS", Type: "vec4", Description: "Specular shininess", Stage: "fragment", ReadWrite: "in"},
		"VERTEX":             {Name: "VERTEX", Type: "vec2", Description: "Vertex position", Stage: "fragment", ReadWrite: "in"},
	}
}

// GetCanvasItemLightBuiltins returns built-in variables for canvas_item shader light stage.
func GetCanvasItemLightBuiltins() map[string]*BuiltinVariable {
	return map[string]*BuiltinVariable{
		"FRAGCOORD":            {Name: "FRAGCOORD", Type: "vec4", Description: "Fragment coordinates", Stage: "light", ReadWrite: "in"},
		"NORMAL":               {Name: "NORMAL", Type: "vec3", Description: "Normal", Stage: "light", ReadWrite: "in"},
		"COLOR":                {Name: "COLOR", Type: "vec4", Description: "Color from fragment", Stage: "light", ReadWrite: "in"},
		"UV":                   {Name: "UV", Type: "vec2", Description: "UV coordinates", Stage: "light", ReadWrite: "in"},
		"SPECULAR_SHININESS":   {Name: "SPECULAR_SHININESS", Type: "vec4", Description: "Specular shininess", Stage: "light", ReadWrite: "in"},
		"LIGHT_COLOR":          {Name: "LIGHT_COLOR", Type: "vec4", Description: "Light color", Stage: "light", ReadWrite: "in"},
		"LIGHT_POSITION":       {Name: "LIGHT_POSITION", Type: "vec3", Description: "Light position", Stage: "light", ReadWrite: "in"},
		"LIGHT_DIRECTION":      {Name: "LIGHT_DIRECTION", Type: "vec3", Description: "Light direction", Stage: "light", ReadWrite: "in"},
		"LIGHT_IS_DIRECTIONAL": {Name: "LIGHT_IS_DIRECTIONAL", Type: "bool", Description: "Is directional light", Stage: "light", ReadWrite: "in"},
		"LIGHT_ENERGY":         {Name: "LIGHT_ENERGY", Type: "float", Description: "Light energy", Stage: "light", ReadWrite: "in"},
		"LIGHT_VERTEX":         {Name: "LIGHT_VERTEX", Type: "vec3", Description: "Light vertex", Stage: "light", ReadWrite: "in"},
		"LIGHT":                {Name: "LIGHT", Type: "vec4", Description: "Light output", Stage: "light", ReadWrite: "out"},
		"SHADOW_MODULATE":      {Name: "SHADOW_MODULATE", Type: "vec4", Description: "Shadow modulate", Stage: "light", ReadWrite: "in"},
		"SCREEN_UV":            {Name: "SCREEN_UV", Type: "vec2", Description: "Screen UV", Stage: "light", ReadWrite: "in"},
		"TEXTURE":              {Name: "TEXTURE", Type: "sampler2D", Description: "Main texture", Stage: "light", ReadWrite: "in"},
		"TEXTURE_PIXEL_SIZE":   {Name: "TEXTURE_PIXEL_SIZE", Type: "vec2", Description: "Texture pixel size", Stage: "light", ReadWrite: "in"},
		"POINT_COORD":          {Name: "POINT_COORD", Type: "vec2", Description: "Point coordinate", Stage: "light", ReadWrite: "in"},
		"TIME":                 {Name: "TIME", Type: "float", Description: "Time", Stage: "light", ReadWrite: "in"},
	}
}

// GetParticlesBuiltins returns built-in variables for particles shader.
func GetParticlesBuiltins() map[string]*BuiltinVariable {
	return map[string]*BuiltinVariable{
		"COLOR":              {Name: "COLOR", Type: "vec4", Description: "Particle color", Stage: "start", ReadWrite: "inout"},
		"VELOCITY":           {Name: "VELOCITY", Type: "vec3", Description: "Particle velocity", Stage: "start", ReadWrite: "inout"},
		"MASS":               {Name: "MASS", Type: "float", Description: "Particle mass", Stage: "start", ReadWrite: "inout"},
		"ACTIVE":             {Name: "ACTIVE", Type: "bool", Description: "Is particle active", Stage: "start", ReadWrite: "inout"},
		"RESTART":            {Name: "RESTART", Type: "bool", Description: "Restart flag", Stage: "start", ReadWrite: "in"},
		"CUSTOM":             {Name: "CUSTOM", Type: "vec4", Description: "Custom data", Stage: "start", ReadWrite: "inout"},
		"TRANSFORM":          {Name: "TRANSFORM", Type: "mat4", Description: "Particle transform", Stage: "start", ReadWrite: "inout"},
		"LIFETIME":           {Name: "LIFETIME", Type: "float", Description: "Particle lifetime", Stage: "start", ReadWrite: "in"},
		"DELTA":              {Name: "DELTA", Type: "float", Description: "Delta time", Stage: "start", ReadWrite: "in"},
		"NUMBER":             {Name: "NUMBER", Type: "uint", Description: "Particle number", Stage: "start", ReadWrite: "in"},
		"INDEX":              {Name: "INDEX", Type: "int", Description: "Particle index", Stage: "start", ReadWrite: "in"},
		"EMISSION_TRANSFORM": {Name: "EMISSION_TRANSFORM", Type: "mat4", Description: "Emission transform", Stage: "start", ReadWrite: "in"},
		"RANDOM_SEED":        {Name: "RANDOM_SEED", Type: "uint", Description: "Random seed", Stage: "start", ReadWrite: "in"},
		"TIME":               {Name: "TIME", Type: "float", Description: "Time", Stage: "start", ReadWrite: "in"},
		"INTERPOLATE_TO_END": {Name: "INTERPOLATE_TO_END", Type: "float", Description: "Interpolation to end", Stage: "start", ReadWrite: "in"},
		"AMOUNT_RATIO":       {Name: "AMOUNT_RATIO", Type: "float", Description: "Amount ratio", Stage: "start", ReadWrite: "in"},
	}
}

// GetSkyBuiltins returns built-in variables for sky shader.
func GetSkyBuiltins() map[string]*BuiltinVariable {
	return map[string]*BuiltinVariable{
		"RADIANCE":            {Name: "RADIANCE", Type: "vec3", Description: "Radiance output", Stage: "sky", ReadWrite: "out"},
		"IRRADIANCE":          {Name: "IRRADIANCE", Type: "vec3", Description: "Irradiance output", Stage: "sky", ReadWrite: "out"},
		"FOG":                 {Name: "FOG", Type: "vec4", Description: "Fog output", Stage: "sky", ReadWrite: "out"},
		"AT_CUBEMAP_PASS":     {Name: "AT_CUBEMAP_PASS", Type: "bool", Description: "Is cubemap pass", Stage: "sky", ReadWrite: "in"},
		"AT_HALF_RES_PASS":    {Name: "AT_HALF_RES_PASS", Type: "bool", Description: "Is half res pass", Stage: "sky", ReadWrite: "in"},
		"AT_QUARTER_RES_PASS": {Name: "AT_QUARTER_RES_PASS", Type: "bool", Description: "Is quarter res pass", Stage: "sky", ReadWrite: "in"},
		"EYEDIR":              {Name: "EYEDIR", Type: "vec3", Description: "Eye direction", Stage: "sky", ReadWrite: "in"},
		"HALF_RES_COLOR":      {Name: "HALF_RES_COLOR", Type: "vec4", Description: "Half res color", Stage: "sky", ReadWrite: "in"},
		"QUARTER_RES_COLOR":   {Name: "QUARTER_RES_COLOR", Type: "vec4", Description: "Quarter res color", Stage: "sky", ReadWrite: "in"},
		"SCREEN_UV":           {Name: "SCREEN_UV", Type: "vec2", Description: "Screen UV", Stage: "sky", ReadWrite: "in"},
		"SKY_COORDS":          {Name: "SKY_COORDS", Type: "vec2", Description: "Sky coordinates", Stage: "sky", ReadWrite: "in"},
		"TIME":                {Name: "TIME", Type: "float", Description: "Time", Stage: "sky", ReadWrite: "in"},
		"POSITION":            {Name: "POSITION", Type: "vec3", Description: "World position", Stage: "sky", ReadWrite: "in"},
		"LIGHT0_ENABLED":      {Name: "LIGHT0_ENABLED", Type: "bool", Description: "Light 0 enabled", Stage: "sky", ReadWrite: "in"},
		"LIGHT0_DIRECTION":    {Name: "LIGHT0_DIRECTION", Type: "vec3", Description: "Light 0 direction", Stage: "sky", ReadWrite: "in"},
		"LIGHT0_ENERGY":       {Name: "LIGHT0_ENERGY", Type: "float", Description: "Light 0 energy", Stage: "sky", ReadWrite: "in"},
		"LIGHT0_COLOR":        {Name: "LIGHT0_COLOR", Type: "vec3", Description: "Light 0 color", Stage: "sky", ReadWrite: "in"},
		"LIGHT0_SIZE":         {Name: "LIGHT0_SIZE", Type: "float", Description: "Light 0 size", Stage: "sky", ReadWrite: "in"},
		"LIGHT1_ENABLED":      {Name: "LIGHT1_ENABLED", Type: "bool", Description: "Light 1 enabled", Stage: "sky", ReadWrite: "in"},
		"LIGHT1_DIRECTION":    {Name: "LIGHT1_DIRECTION", Type: "vec3", Description: "Light 1 direction", Stage: "sky", ReadWrite: "in"},
		"LIGHT1_ENERGY":       {Name: "LIGHT1_ENERGY", Type: "float", Description: "Light 1 energy", Stage: "sky", ReadWrite: "in"},
		"LIGHT1_COLOR":        {Name: "LIGHT1_COLOR", Type: "vec3", Description: "Light 1 color", Stage: "sky", ReadWrite: "in"},
		"LIGHT1_SIZE":         {Name: "LIGHT1_SIZE", Type: "float", Description: "Light 1 size", Stage: "sky", ReadWrite: "in"},
		"LIGHT2_ENABLED":      {Name: "LIGHT2_ENABLED", Type: "bool", Description: "Light 2 enabled", Stage: "sky", ReadWrite: "in"},
		"LIGHT2_DIRECTION":    {Name: "LIGHT2_DIRECTION", Type: "vec3", Description: "Light 2 direction", Stage: "sky", ReadWrite: "in"},
		"LIGHT2_ENERGY":       {Name: "LIGHT2_ENERGY", Type: "float", Description: "Light 2 energy", Stage: "sky", ReadWrite: "in"},
		"LIGHT2_COLOR":        {Name: "LIGHT2_COLOR", Type: "vec3", Description: "Light 2 color", Stage: "sky", ReadWrite: "in"},
		"LIGHT2_SIZE":         {Name: "LIGHT2_SIZE", Type: "float", Description: "Light 2 size", Stage: "sky", ReadWrite: "in"},
		"LIGHT3_ENABLED":      {Name: "LIGHT3_ENABLED", Type: "bool", Description: "Light 3 enabled", Stage: "sky", ReadWrite: "in"},
		"LIGHT3_DIRECTION":    {Name: "LIGHT3_DIRECTION", Type: "vec3", Description: "Light 3 direction", Stage: "sky", ReadWrite: "in"},
		"LIGHT3_ENERGY":       {Name: "LIGHT3_ENERGY", Type: "float", Description: "Light 3 energy", Stage: "sky", ReadWrite: "in"},
		"LIGHT3_COLOR":        {Name: "LIGHT3_COLOR", Type: "vec3", Description: "Light 3 color", Stage: "sky", ReadWrite: "in"},
		"LIGHT3_SIZE":         {Name: "LIGHT3_SIZE", Type: "float", Description: "Light 3 size", Stage: "sky", ReadWrite: "in"},
	}
}

// GetFogBuiltins returns built-in variables for fog shader.
func GetFogBuiltins() map[string]*BuiltinVariable {
	return map[string]*BuiltinVariable{
		"WORLD_POSITION":  {Name: "WORLD_POSITION", Type: "vec3", Description: "World position", Stage: "fog", ReadWrite: "in"},
		"OBJECT_POSITION": {Name: "OBJECT_POSITION", Type: "vec3", Description: "Object position", Stage: "fog", ReadWrite: "in"},
		"UVW":             {Name: "UVW", Type: "vec3", Description: "UVW coordinates", Stage: "fog", ReadWrite: "in"},
		"SIZE":            {Name: "SIZE", Type: "vec3", Description: "Size", Stage: "fog", ReadWrite: "in"},
		"SDF":             {Name: "SDF", Type: "float", Description: "Signed distance field", Stage: "fog", ReadWrite: "in"},
		"ALBEDO":          {Name: "ALBEDO", Type: "vec3", Description: "Albedo output", Stage: "fog", ReadWrite: "out"},
		"DENSITY":         {Name: "DENSITY", Type: "float", Description: "Density output", Stage: "fog", ReadWrite: "out"},
		"EMISSION":        {Name: "EMISSION", Type: "vec3", Description: "Emission output", Stage: "fog", ReadWrite: "out"},
		"TIME":            {Name: "TIME", Type: "float", Description: "Time", Stage: "fog", ReadWrite: "in"},
	}
}

// GetBuiltinsForShaderType returns all built-in variables for a given shader type.
func GetBuiltinsForShaderType(shaderType string) map[string]*BuiltinVariable {
	result := make(map[string]*BuiltinVariable)

	switch shaderType {
	case "spatial":
		for k, v := range GetSpatialVertexBuiltins() {
			result[k] = v
		}
		for k, v := range GetSpatialFragmentBuiltins() {
			result[k] = v
		}
		for k, v := range GetSpatialLightBuiltins() {
			result[k] = v
		}
	case "canvas_item":
		for k, v := range GetCanvasItemVertexBuiltins() {
			result[k] = v
		}
		for k, v := range GetCanvasItemFragmentBuiltins() {
			result[k] = v
		}
		for k, v := range GetCanvasItemLightBuiltins() {
			result[k] = v
		}
	case "particles":
		for k, v := range GetParticlesBuiltins() {
			result[k] = v
		}
	case "sky":
		for k, v := range GetSkyBuiltins() {
			result[k] = v
		}
	case "fog":
		for k, v := range GetFogBuiltins() {
			result[k] = v
		}
	}

	return result
}

// UniformHints contains information about uniform hints.
var UniformHints = map[string]string{
	"source_color":                      "Used as albedo or color (sRGB conversion applied)",
	"hint_range":                        "Restricts value to range: hint_range(min, max[, step])",
	"hint_normal":                       "Used as normal map",
	"hint_default_white":                "Default to opaque white",
	"hint_default_black":                "Default to opaque black",
	"hint_default_transparent":          "Default to transparent black",
	"hint_anisotropy":                   "Used as flowmap for anisotropy",
	"hint_roughness_r":                  "Roughness stored in red channel",
	"hint_roughness_g":                  "Roughness stored in green channel",
	"hint_roughness_b":                  "Roughness stored in blue channel",
	"hint_roughness_a":                  "Roughness stored in alpha channel",
	"hint_roughness_normal":             "Roughness guided by normal map",
	"hint_roughness_gray":               "Roughness from grayscale",
	"hint_screen_texture":               "Screen texture sampler",
	"hint_depth_texture":                "Depth texture sampler",
	"hint_normal_roughness_texture":     "Normal roughness texture (Forward+ only)",
	"filter_nearest":                    "Use nearest filtering",
	"filter_linear":                     "Use linear filtering",
	"filter_nearest_mipmap":             "Use nearest filtering with mipmaps",
	"filter_linear_mipmap":              "Use linear filtering with mipmaps",
	"filter_nearest_mipmap_anisotropic": "Use nearest filtering with anisotropic mipmaps",
	"filter_linear_mipmap_anisotropic":  "Use linear filtering with anisotropic mipmaps",
	"repeat_enable":                     "Enable texture repeat",
	"repeat_disable":                    "Disable texture repeat",
	"hint_enum":                         "Display as dropdown: hint_enum(\"Option1\", \"Option2\", ...)",
}
