package lsp

import (
	"fmt"
	"strings"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/andresperezl/gdls/internal/analysis"
	"github.com/andresperezl/gdls/internal/gdshader"
	"github.com/andresperezl/gdls/internal/parser"
)

// textDocumentHover handles the textDocument/hover request.
func (s *Server) textDocumentHover(ctx *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	doc := s.workspace.GetDocument(params.TextDocument.URI)
	if doc == nil {
		return nil, nil
	}

	line := int(params.Position.Line)
	col := int(params.Position.Character)

	var hoverInfo string

	switch doc.Type {
	case analysis.DocumentTypeTSCN:
		if doc.TSCNAST == nil {
			return nil, nil
		}
		hoverInfo = s.findTSCNHoverInfo(doc, line, col)
	case analysis.DocumentTypeGDShader:
		if doc.ShaderAST == nil {
			return nil, nil
		}
		hoverInfo = s.findGDShaderHoverInfo(doc, line, col)
	default:
		return nil, nil
	}

	if hoverInfo == "" {
		return nil, nil
	}

	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: hoverInfo,
		},
	}, nil
}

// findTSCNHoverInfo finds hover information for a TSCN document position.
func (s *Server) findTSCNHoverInfo(doc *analysis.Document, line, col int) string {
	ast := doc.TSCNAST

	// Check external resources
	for _, ext := range ast.ExtResources {
		if isInRange(ext.Range, line, col) {
			return formatExtResourceHover(ext)
		}
	}

	// Check sub resources
	for _, sub := range ast.SubResources {
		if isInRange(sub.Range, line, col) {
			// Check if we're on a specific property
			for _, prop := range sub.Properties {
				if isInRange(prop.Range, line, col) {
					return formatPropertyHover(prop, sub.Type)
				}
			}
			return formatSubResourceHover(sub)
		}
	}

	// Check nodes
	for _, node := range ast.Nodes {
		if isInRange(node.Range, line, col) {
			// Check if we're on a specific property
			for _, prop := range node.Properties {
				if isInRange(prop.Range, line, col) {
					return formatPropertyHover(prop, node.Type)
				}
			}
			return formatNodeHover(node, doc)
		}
	}

	// Check connections
	for _, conn := range ast.Connections {
		if isInRange(conn.Range, line, col) {
			return formatConnectionHover(conn)
		}
	}

	// Check descriptor
	if ast.Descriptor != nil && isInRange(ast.Descriptor.Range, line, col) {
		return formatDescriptorHover(ast.Descriptor)
	}

	return ""
}

func formatExtResourceHover(ext *parser.ExtResource) string {
	var sb strings.Builder
	sb.WriteString("### External Resource\n\n")
	sb.WriteString(fmt.Sprintf("**Type:** `%s`\n\n", ext.Type))
	sb.WriteString(fmt.Sprintf("**Path:** `%s`\n\n", ext.Path))
	sb.WriteString(fmt.Sprintf("**ID:** `%s`\n\n", ext.ID))
	if ext.UID != "" {
		sb.WriteString(fmt.Sprintf("**UID:** `%s`\n", ext.UID))
	}
	return sb.String()
}

func formatSubResourceHover(sub *parser.SubResource) string {
	var sb strings.Builder
	sb.WriteString("### Internal Resource\n\n")
	sb.WriteString(fmt.Sprintf("**Type:** `%s`\n\n", sub.Type))
	sb.WriteString(fmt.Sprintf("**ID:** `%s`\n\n", sub.ID))

	// Add type description if available
	if desc := getGodotTypeDescription(sub.Type); desc != "" {
		sb.WriteString(fmt.Sprintf("_%s_\n", desc))
	}

	return sb.String()
}

func formatNodeHover(node *parser.Node, doc *analysis.Document) string {
	var sb strings.Builder
	sb.WriteString("### Scene Node\n\n")
	sb.WriteString(fmt.Sprintf("**Name:** `%s`\n\n", node.Name))

	if node.Type != "" {
		sb.WriteString(fmt.Sprintf("**Type:** `%s`\n\n", node.Type))
		if desc := getGodotTypeDescription(node.Type); desc != "" {
			sb.WriteString(fmt.Sprintf("_%s_\n\n", desc))
		}
	} else if node.Instance != nil {
		if ref, ok := node.Instance.(*parser.ResourceRef); ok {
			// Find the external resource to get the path
			for _, ext := range doc.TSCNAST.ExtResources {
				if ext.ID == ref.ID {
					sb.WriteString(fmt.Sprintf("**Instance of:** `%s`\n\n", ext.Path))
					break
				}
			}
		}
	}

	if node.Parent != "" {
		sb.WriteString(fmt.Sprintf("**Parent:** `%s`\n\n", node.Parent))
	} else {
		sb.WriteString("**Parent:** _(scene root)_\n\n")
	}

	if len(node.Groups) > 0 {
		sb.WriteString(fmt.Sprintf("**Groups:** `%s`\n\n", strings.Join(node.Groups, "`, `")))
	}

	return sb.String()
}

func formatPropertyHover(prop *parser.Property, ownerType string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("### Property: `%s`\n\n", prop.Key))

	// Describe the value type
	valueType := describeValueType(prop.Value)
	sb.WriteString(fmt.Sprintf("**Type:** `%s`\n\n", valueType))

	// Show a preview of the value
	valuePreview := formatValuePreview(prop.Value)
	if valuePreview != "" {
		sb.WriteString(fmt.Sprintf("**Value:** `%s`\n", valuePreview))
	}

	return sb.String()
}

func formatConnectionHover(conn *parser.Connection) string {
	var sb strings.Builder
	sb.WriteString("### Signal Connection\n\n")
	sb.WriteString(fmt.Sprintf("**Signal:** `%s`\n\n", conn.Signal))
	sb.WriteString(fmt.Sprintf("**From:** `%s`\n\n", conn.From))
	sb.WriteString(fmt.Sprintf("**To:** `%s`\n\n", conn.To))
	sb.WriteString(fmt.Sprintf("**Method:** `%s`\n", conn.Method))
	return sb.String()
}

func formatDescriptorHover(desc *parser.GdScene) string {
	var sb strings.Builder
	if desc.Type == "gd_scene" {
		sb.WriteString("### Scene File\n\n")
	} else {
		sb.WriteString("### Resource File\n\n")
		if desc.ResourceType != "" {
			sb.WriteString(fmt.Sprintf("**Type:** `%s`\n\n", desc.ResourceType))
		}
	}
	sb.WriteString(fmt.Sprintf("**Format:** `%d` (Godot 4.x)\n\n", desc.Format))
	if desc.LoadSteps != nil {
		sb.WriteString(fmt.Sprintf("**Load Steps:** `%d`\n\n", *desc.LoadSteps))
	}
	if desc.UID != "" {
		sb.WriteString(fmt.Sprintf("**UID:** `%s`\n", desc.UID))
	}
	return sb.String()
}

func describeValueType(v parser.Value) string {
	switch val := v.(type) {
	case *parser.StringValue:
		return "String"
	case *parser.NumberValue:
		if val.IsInt {
			return "int"
		}
		return "float"
	case *parser.BoolValue:
		return "bool"
	case *parser.NullValue:
		return "null"
	case *parser.ArrayValue:
		return "Array"
	case *parser.DictValue:
		return "Dictionary"
	case *parser.TypedValue:
		return val.TypeName
	case *parser.ResourceRef:
		return val.RefType
	case *parser.IdentValue:
		return "Identifier"
	default:
		return "unknown"
	}
}

func formatValuePreview(v parser.Value) string {
	switch val := v.(type) {
	case *parser.StringValue:
		if len(val.Value) > 50 {
			return fmt.Sprintf("\"%s...\"", val.Value[:50])
		}
		return fmt.Sprintf("\"%s\"", val.Value)
	case *parser.NumberValue:
		return val.RawValue
	case *parser.BoolValue:
		return fmt.Sprintf("%v", val.Value)
	case *parser.NullValue:
		return "null"
	case *parser.TypedValue:
		return fmt.Sprintf("%s(...)", val.TypeName)
	case *parser.ResourceRef:
		return fmt.Sprintf("%s(\"%s\")", val.RefType, val.ID)
	case *parser.ArrayValue:
		return fmt.Sprintf("[...] (%d items)", len(val.Values))
	case *parser.DictValue:
		return fmt.Sprintf("{...} (%d entries)", len(val.Entries))
	default:
		return ""
	}
}

func isInRange(r parser.Range, line, col int) bool {
	if line < r.Start.Line || line > r.End.Line {
		return false
	}
	if line == r.Start.Line && col < r.Start.Column {
		return false
	}
	if line == r.End.Line && col > r.End.Column {
		return false
	}
	return true
}

// getGodotTypeDescription returns a brief description for common Godot types.
func getGodotTypeDescription(typeName string) string {
	descriptions := map[string]string{
		// Nodes
		"Node":                "Base class for all scene objects",
		"Node2D":              "A 2D game object",
		"Node3D":              "A 3D game object (formerly Spatial)",
		"Control":             "Base class for all UI-related nodes",
		"Camera2D":            "Camera node for 2D scenes",
		"Camera3D":            "Camera node for 3D scenes",
		"CharacterBody2D":     "2D physics body for character movement",
		"CharacterBody3D":     "3D physics body for character movement",
		"RigidBody2D":         "2D physics body with rigid body dynamics",
		"RigidBody3D":         "3D physics body with rigid body dynamics",
		"StaticBody2D":        "2D physics body that doesn't move",
		"StaticBody3D":        "3D physics body that doesn't move",
		"Area2D":              "2D area for detecting overlaps",
		"Area3D":              "3D area for detecting overlaps",
		"CollisionShape2D":    "2D collision shape for physics",
		"CollisionShape3D":    "3D collision shape for physics",
		"Sprite2D":            "2D sprite node",
		"Sprite3D":            "3D sprite node",
		"MeshInstance3D":      "Instance of a 3D mesh",
		"AnimationPlayer":     "Node for playing animations",
		"AnimationTree":       "Node for blending animations",
		"AudioStreamPlayer":   "Plays audio non-positionally",
		"AudioStreamPlayer2D": "Plays audio with 2D positioning",
		"AudioStreamPlayer3D": "Plays audio with 3D positioning",
		"Label":               "Displays text",
		"Button":              "Clickable button",
		"TextEdit":            "Multi-line text editor",
		"LineEdit":            "Single-line text input",
		"Timer":               "Counts down and emits timeout signal",
		"Path2D":              "Contains a Curve2D path",
		"Path3D":              "Contains a Curve3D path",
		"PathFollow2D":        "Follows a Path2D",
		"PathFollow3D":        "Follows a Path3D",
		"Skeleton3D":          "3D skeleton for mesh deformation",
		"BoneAttachment3D":    "Attaches nodes to skeleton bones",
		"GPUParticles2D":      "2D GPU-accelerated particles",
		"GPUParticles3D":      "3D GPU-accelerated particles",
		"DirectionalLight3D":  "Directional light source",
		"OmniLight3D":         "Omnidirectional point light",
		"SpotLight3D":         "Spotlight",

		// Resources
		"BoxShape3D":         "3D box collision shape",
		"SphereShape3D":      "3D sphere collision shape",
		"CapsuleShape3D":     "3D capsule collision shape",
		"CylinderShape3D":    "3D cylinder collision shape",
		"BoxMesh":            "Box primitive mesh",
		"SphereMesh":         "Sphere primitive mesh",
		"CapsuleMesh":        "Capsule primitive mesh",
		"CylinderMesh":       "Cylinder primitive mesh",
		"PlaneMesh":          "Plane primitive mesh",
		"ArrayMesh":          "Mesh from vertex arrays",
		"StandardMaterial3D": "PBR material for 3D",
		"ShaderMaterial":     "Custom shader material",
		"Texture2D":          "2D texture resource",
		"Animation":          "Animation resource",
		"AnimationLibrary":   "Collection of animations",
		"PackedScene":        "Serialized scene",
		"Script":             "GDScript or other script",
		"AudioStream":        "Audio data resource",
		"Font":               "Font resource",
		"Theme":              "UI theme resource",
	}

	if desc, ok := descriptions[typeName]; ok {
		return desc
	}
	return ""
}

// findGDShaderHoverInfo finds hover information for a GDShader document position.
func (s *Server) findGDShaderHoverInfo(doc *analysis.Document, line, col int) string {
	ast := doc.ShaderAST
	if ast == nil {
		return ""
	}

	// Check shader_type declaration
	if ast.ShaderType != nil && isInGDShaderRange(ast.ShaderType.Range, line, col) {
		return formatShaderTypeHover(ast.ShaderType.Type)
	}

	// Check render_mode declaration
	if ast.RenderModes != nil && isInGDShaderRange(ast.RenderModes.Range, line, col) {
		return formatRenderModeHover(ast.RenderModes)
	}

	// Check uniforms
	for _, uniform := range ast.Uniforms {
		if isInGDShaderRange(uniform.Range, line, col) {
			return formatUniformHover(uniform)
		}
	}

	// Check varyings
	for _, varying := range ast.Varyings {
		if isInGDShaderRange(varying.Range, line, col) {
			return formatVaryingHover(varying)
		}
	}

	// Check constants
	for _, constant := range ast.Constants {
		if isInGDShaderRange(constant.Range, line, col) {
			return formatConstantHover(constant)
		}
	}

	// Check functions
	for _, fn := range ast.Functions {
		if isInGDShaderRange(fn.Range, line, col) {
			return formatFunctionHover(fn)
		}
	}

	// Check structs
	for _, st := range ast.Structs {
		if isInGDShaderRange(st.Range, line, col) {
			return formatStructHover(st)
		}
	}

	// Check for built-in function or constant at position
	// This requires finding the identifier at the position
	hoverInfo := s.findGDShaderBuiltinHover(doc, line, col)
	if hoverInfo != "" {
		return hoverInfo
	}

	return ""
}

// isInGDShaderRange checks if a position is within a GDShader range.
func isInGDShaderRange(r gdshader.Range, line, col int) bool {
	if line < r.Start.Line || line > r.End.Line {
		return false
	}
	if line == r.Start.Line && col < r.Start.Column {
		return false
	}
	if line == r.End.Line && col > r.End.Column {
		return false
	}
	return true
}

func formatShaderTypeHover(shaderType string) string {
	var sb strings.Builder
	sb.WriteString("### Shader Type\n\n")
	sb.WriteString(fmt.Sprintf("**Type:** `%s`\n\n", shaderType))

	descriptions := map[string]string{
		"spatial":     "3D shader for MeshInstance3D and other 3D nodes. Supports vertex, fragment, and light functions.",
		"canvas_item": "2D shader for CanvasItem nodes like Sprite2D, Control, etc. Supports vertex, fragment, and light functions.",
		"particles":   "Shader for GPUParticles2D/3D. Supports start and process functions for particle behavior.",
		"sky":         "Shader for Sky resource. Used for rendering sky backgrounds.",
		"fog":         "Shader for FogVolume. Used for volumetric fog effects.",
	}

	if desc, ok := descriptions[shaderType]; ok {
		sb.WriteString(fmt.Sprintf("_%s_\n", desc))
	}

	return sb.String()
}

func formatRenderModeHover(rm *gdshader.RenderModeDecl) string {
	var sb strings.Builder
	sb.WriteString("### Render Modes\n\n")
	for _, mode := range rm.Modes {
		sb.WriteString(fmt.Sprintf("- `%s`\n", mode))
	}
	return sb.String()
}

func formatUniformHover(uniform *gdshader.UniformDecl) string {
	var sb strings.Builder
	sb.WriteString("### Uniform Variable\n\n")
	sb.WriteString(fmt.Sprintf("**Name:** `%s`\n\n", uniform.Name))
	sb.WriteString(fmt.Sprintf("**Type:** `%s`\n\n", uniform.Type.Name))

	if uniform.IsGlobal {
		sb.WriteString("**Scope:** `global`\n\n")
	}

	if len(uniform.Hints) > 0 {
		sb.WriteString("**Hints:**\n")
		for _, hint := range uniform.Hints {
			sb.WriteString(fmt.Sprintf("- `%s`\n", hint.Name))
		}
	}

	if uniform.DocComment != "" {
		sb.WriteString(fmt.Sprintf("\n_%s_\n", uniform.DocComment))
	}

	return sb.String()
}

func formatVaryingHover(varying *gdshader.VaryingDecl) string {
	var sb strings.Builder
	sb.WriteString("### Varying Variable\n\n")
	sb.WriteString(fmt.Sprintf("**Name:** `%s`\n\n", varying.Name))
	sb.WriteString(fmt.Sprintf("**Type:** `%s`\n\n", varying.Type.Name))

	if varying.Interpolation != "" {
		sb.WriteString(fmt.Sprintf("**Interpolation:** `%s`\n\n", varying.Interpolation))
	}

	sb.WriteString("_Passed between vertex and fragment shaders._\n")
	return sb.String()
}

func formatConstantHover(constant *gdshader.ConstDecl) string {
	var sb strings.Builder
	sb.WriteString("### Constant\n\n")
	sb.WriteString(fmt.Sprintf("**Name:** `%s`\n\n", constant.Name))
	sb.WriteString(fmt.Sprintf("**Type:** `%s`\n\n", constant.Type.Name))
	return sb.String()
}

func formatFunctionHover(fn *gdshader.FunctionDecl) string {
	var sb strings.Builder
	sb.WriteString("### Function\n\n")

	// Build signature
	var params []string
	for _, param := range fn.Params {
		p := param.Type.Name + " " + param.Name
		if param.Qualifier != "" {
			p = param.Qualifier + " " + p
		}
		params = append(params, p)
	}

	sb.WriteString(fmt.Sprintf("```gdshader\n%s %s(%s)\n```\n\n",
		fn.ReturnType.Name, fn.Name, strings.Join(params, ", ")))

	// Special function descriptions
	switch fn.Name {
	case "vertex":
		sb.WriteString("_Runs for each vertex. Used to transform vertex positions._\n")
	case "fragment":
		sb.WriteString("_Runs for each pixel. Used to determine final color._\n")
	case "light":
		sb.WriteString("_Runs for each light affecting a pixel._\n")
	case "start":
		sb.WriteString("_Runs once when a particle spawns._\n")
	case "process":
		sb.WriteString("_Runs each frame for each particle._\n")
	case "sky":
		sb.WriteString("_Runs for each pixel of the sky._\n")
	case "fog":
		sb.WriteString("_Runs for each sample in the fog volume._\n")
	}

	return sb.String()
}

func formatStructHover(st *gdshader.StructDecl) string {
	var sb strings.Builder
	sb.WriteString("### Struct\n\n")
	sb.WriteString(fmt.Sprintf("**Name:** `%s`\n\n", st.Name))

	if len(st.Members) > 0 {
		sb.WriteString("**Members:**\n")
		for _, member := range st.Members {
			sb.WriteString(fmt.Sprintf("- `%s %s`\n", member.Type.Name, member.Name))
		}
	}

	return sb.String()
}

// findGDShaderBuiltinHover finds hover info for built-in functions and constants.
func (s *Server) findGDShaderBuiltinHover(doc *analysis.Document, line, col int) string {
	// Get the word at position
	content := doc.Content
	lines := strings.Split(content, "\n")
	if line >= len(lines) {
		return ""
	}

	lineContent := lines[line]
	if col >= len(lineContent) {
		return ""
	}

	// Find word boundaries
	start := col
	for start > 0 && isIdentChar(lineContent[start-1]) {
		start--
	}
	end := col
	for end < len(lineContent) && isIdentChar(lineContent[end]) {
		end++
	}

	if start == end {
		return ""
	}

	word := lineContent[start:end]

	// Check built-in constants
	if constant, ok := gdshader.BuiltinConstants[word]; ok {
		var sb strings.Builder
		sb.WriteString("### Built-in Constant\n\n")
		sb.WriteString(fmt.Sprintf("```gdshader\nconst %s %s = %s\n```\n\n", constant.Type, constant.Name, constant.Value))
		sb.WriteString(fmt.Sprintf("_%s_\n", constant.Description))
		return sb.String()
	}

	// Check built-in functions
	if fn, ok := gdshader.BuiltinFunctions[word]; ok {
		var sb strings.Builder
		sb.WriteString("### Built-in Function\n\n")

		// Show all overloads
		sb.WriteString("```gdshader\n")
		for _, sig := range fn.Signatures {
			sb.WriteString(fmt.Sprintf("%s %s(%s)\n", sig.Return, fn.Name, strings.Join(sig.Params, ", ")))
		}
		sb.WriteString("```\n\n")

		sb.WriteString(fmt.Sprintf("_%s_\n", fn.Description))
		return sb.String()
	}

	return ""
}

func isIdentChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}
