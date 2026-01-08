package parser

// Document represents a parsed TSCN document.
type Document struct {
	Descriptor   *GdScene       // [gd_scene ...] or [gd_resource ...]
	ExtResources []*ExtResource // [ext_resource ...]
	SubResources []*SubResource // [sub_resource ...]
	Nodes        []*Node        // [node ...]
	Connections  []*Connection  // [connection ...]
	Comments     []*Comment     // ; comments
	Errors       []ParseError   // Syntax errors
}

// GdScene represents the file descriptor [gd_scene ...] or [gd_resource ...].
type GdScene struct {
	Range     Range
	Type      string // "gd_scene" or "gd_resource"
	LoadSteps *int
	Format    int
	UID       string
	// For gd_resource
	ResourceType string
}

// ExtResource represents an external resource [ext_resource ...].
type ExtResource struct {
	Range     Range
	Type      string // e.g., "Texture2D", "Material"
	UID       string // uid://...
	Path      string // res://... or relative path
	PathRange Range  // Range of the path string (for go-to-definition)
	ID        string // e.g., "1_7bt6s"
}

// SubResource represents an internal resource [sub_resource ...].
type SubResource struct {
	Range      Range
	Type       string // e.g., "SphereShape3D"
	ID         string // e.g., "SphereShape3D_tj6p1"
	Properties []*Property
}

// Node represents a scene node [node ...].
type Node struct {
	Range               Range
	Name                string
	Type                string // optional (missing for instance nodes)
	Parent              string // "." or "Path/To/Parent", empty for root
	Instance            Value  // ExtResource("id") for instanced scenes
	InstancePlaceholder string
	Owner               string
	Index               *int
	Groups              []string
	Properties          []*Property
}

// Connection represents a signal connection [connection ...].
type Connection struct {
	Range  Range
	Signal string
	From   string // NodePath
	To     string // NodePath
	Method string
	Flags  *int
	Binds  []Value
}

// Property represents a key = value pair.
type Property struct {
	Range    Range
	Key      string // e.g., "transform", "albedo_color", "bones/1/position"
	KeyRange Range  // Range of just the key
	Value    Value
}

// Comment represents a comment line.
type Comment struct {
	Range Range
	Text  string
}

// ParseError represents a parsing error.
type ParseError struct {
	Range   Range
	Message string
}

// Value is an interface for all TSCN value types.
type Value interface {
	valueNode()
	GetRange() Range
}

// StringValue represents a string literal.
type StringValue struct {
	Range Range
	Value string
}

func (v *StringValue) valueNode()      {}
func (v *StringValue) GetRange() Range { return v.Range }

// NumberValue represents a numeric literal (int or float).
type NumberValue struct {
	Range    Range
	Value    float64
	IsInt    bool
	RawValue string // Original text representation
}

func (v *NumberValue) valueNode()      {}
func (v *NumberValue) GetRange() Range { return v.Range }

// BoolValue represents a boolean literal.
type BoolValue struct {
	Range Range
	Value bool
}

func (v *BoolValue) valueNode()      {}
func (v *BoolValue) GetRange() Range { return v.Range }

// NullValue represents a null literal.
type NullValue struct {
	Range Range
}

func (v *NullValue) valueNode()      {}
func (v *NullValue) GetRange() Range { return v.Range }

// ArrayValue represents an array [...].
type ArrayValue struct {
	Range  Range
	Values []Value
}

func (v *ArrayValue) valueNode()      {}
func (v *ArrayValue) GetRange() Range { return v.Range }

// DictValue represents a dictionary {...}.
type DictValue struct {
	Range   Range
	Entries []*DictEntry
}

func (v *DictValue) valueNode()      {}
func (v *DictValue) GetRange() Range { return v.Range }

// DictEntry represents a key-value pair in a dictionary.
type DictEntry struct {
	Range Range
	Key   Value // Usually StringValue or IdentValue
	Value Value
}

// IdentValue represents an unquoted identifier used as a value.
type IdentValue struct {
	Range Range
	Name  string
}

func (v *IdentValue) valueNode()      {}
func (v *IdentValue) GetRange() Range { return v.Range }

// TypedValue represents a typed value like Vector3(1, 2, 3).
type TypedValue struct {
	Range     Range
	TypeName  string // e.g., "Vector3", "Color", "Transform3D"
	TypeRange Range  // Range of just the type name
	Arguments []Value
}

func (v *TypedValue) valueNode()      {}
func (v *TypedValue) GetRange() Range { return v.Range }

// ResourceRef represents ExtResource("id") or SubResource("id").
type ResourceRef struct {
	Range   Range
	RefType string // "ExtResource" or "SubResource"
	ID      string
	IDRange Range // Range of the ID string
}

func (v *ResourceRef) valueNode()      {}
func (v *ResourceRef) GetRange() Range { return v.Range }
