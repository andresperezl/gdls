package lsp

import (
	"strings"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/andresperezl/gdls/internal/analysis"
)

// textDocumentCompletion handles the textDocument/completion request.
func (s *Server) textDocumentCompletion(ctx *glsp.Context, params *protocol.CompletionParams) (any, error) {
	doc := s.workspace.GetDocument(params.TextDocument.URI)
	if doc == nil {
		return nil, nil
	}

	line := int(params.Position.Line)
	col := int(params.Position.Character)

	// Get the line text to understand context
	lines := strings.Split(doc.Content, "\n")
	if line >= len(lines) {
		return nil, nil
	}
	lineText := lines[line]
	if col > len(lineText) {
		col = len(lineText)
	}
	prefix := lineText[:col]

	// Determine completion context
	items := s.getCompletions(doc, prefix, lineText)

	return &protocol.CompletionList{
		IsIncomplete: false,
		Items:        items,
	}, nil
}

// getCompletions returns completion items based on context.
func (s *Server) getCompletions(doc *analysis.Document, prefix, lineText string) []protocol.CompletionItem {
	items := []protocol.CompletionItem{}

	// Inside a type="" attribute
	if strings.Contains(prefix, "type=\"") && !strings.HasSuffix(prefix, "\"") {
		return s.getNodeTypeCompletions()
	}

	// Inside ExtResource("")
	if strings.Contains(prefix, "ExtResource(\"") && !strings.HasSuffix(prefix, "\")") {
		return s.getExtResourceIDCompletions(doc)
	}

	// Inside SubResource("")
	if strings.Contains(prefix, "SubResource(\"") && !strings.HasSuffix(prefix, "\")") {
		return s.getSubResourceIDCompletions(doc)
	}

	// Inside parent=""
	if strings.Contains(prefix, "parent=\"") && !strings.HasSuffix(prefix, "\"") {
		return s.getNodePathCompletions(doc)
	}

	// After = sign, suggest value types
	if strings.HasSuffix(strings.TrimSpace(prefix), "=") {
		return s.getValueCompletions(doc)
	}

	// At start of line or after newline, suggest property names
	trimmed := strings.TrimSpace(prefix)
	if trimmed == "" || !strings.Contains(lineText, "=") {
		return s.getPropertyCompletions()
	}

	return items
}

// getNodeTypeCompletions returns completions for node types.
func (s *Server) getNodeTypeCompletions() []protocol.CompletionItem {
	nodeTypes := []string{
		// Base
		"Node", "Node2D", "Node3D",
		// 2D Physics
		"CharacterBody2D", "RigidBody2D", "StaticBody2D", "Area2D",
		"CollisionShape2D", "CollisionPolygon2D",
		// 3D Physics
		"CharacterBody3D", "RigidBody3D", "StaticBody3D", "Area3D",
		"CollisionShape3D", "CollisionPolygon3D",
		// Visual 2D
		"Sprite2D", "AnimatedSprite2D", "Polygon2D", "Line2D",
		"TileMap", "TileMapLayer",
		// Visual 3D
		"MeshInstance3D", "MultiMeshInstance3D", "Sprite3D",
		"CSGBox3D", "CSGSphere3D", "CSGCylinder3D", "CSGMesh3D",
		// Cameras
		"Camera2D", "Camera3D",
		// Lights
		"DirectionalLight3D", "OmniLight3D", "SpotLight3D",
		"PointLight2D", "DirectionalLight2D",
		// Audio
		"AudioStreamPlayer", "AudioStreamPlayer2D", "AudioStreamPlayer3D",
		// Animation
		"AnimationPlayer", "AnimationTree",
		// UI
		"Control", "Container", "Panel", "Label", "RichTextLabel",
		"Button", "TextureButton", "LinkButton", "OptionButton", "MenuButton",
		"CheckBox", "CheckButton", "SpinBox", "HSlider", "VSlider",
		"ProgressBar", "TextureProgressBar",
		"TextEdit", "LineEdit", "CodeEdit",
		"Tree", "ItemList", "TabContainer", "TabBar",
		"ScrollContainer", "HBoxContainer", "VBoxContainer", "GridContainer",
		"MarginContainer", "CenterContainer", "AspectRatioContainer",
		"ColorRect", "TextureRect", "NinePatchRect",
		"SubViewport", "SubViewportContainer",
		// Navigation
		"NavigationAgent2D", "NavigationAgent3D",
		"NavigationRegion2D", "NavigationRegion3D",
		// Paths
		"Path2D", "Path3D", "PathFollow2D", "PathFollow3D",
		// Particles
		"GPUParticles2D", "GPUParticles3D", "CPUParticles2D", "CPUParticles3D",
		// Other
		"Timer", "HTTPRequest", "RayCast2D", "RayCast3D",
		"Skeleton2D", "Skeleton3D", "BoneAttachment3D",
		"CanvasLayer", "ParallaxBackground", "ParallaxLayer",
		"WorldEnvironment", "RemoteTransform2D", "RemoteTransform3D",
	}

	items := make([]protocol.CompletionItem, 0, len(nodeTypes))
	for _, t := range nodeTypes {
		kind := protocol.CompletionItemKindClass
		items = append(items, protocol.CompletionItem{
			Label:  t,
			Kind:   &kind,
			Detail: strPtr("Godot Node Type"),
		})
	}
	return items
}

// getExtResourceIDCompletions returns completions for external resource IDs.
func (s *Server) getExtResourceIDCompletions(doc *analysis.Document) []protocol.CompletionItem {
	if doc.TSCNAST == nil {
		return nil
	}

	items := make([]protocol.CompletionItem, 0, len(doc.TSCNAST.ExtResources))
	for _, ext := range doc.TSCNAST.ExtResources {
		kind := protocol.CompletionItemKindReference
		items = append(items, protocol.CompletionItem{
			Label:  ext.ID,
			Kind:   &kind,
			Detail: strPtr(ext.Type + " - " + ext.Path),
		})
	}
	return items
}

// getSubResourceIDCompletions returns completions for sub-resource IDs.
func (s *Server) getSubResourceIDCompletions(doc *analysis.Document) []protocol.CompletionItem {
	if doc.TSCNAST == nil {
		return nil
	}

	items := make([]protocol.CompletionItem, 0, len(doc.TSCNAST.SubResources))
	for _, sub := range doc.TSCNAST.SubResources {
		kind := protocol.CompletionItemKindReference
		items = append(items, protocol.CompletionItem{
			Label:  sub.ID,
			Kind:   &kind,
			Detail: strPtr(sub.Type),
		})
	}
	return items
}

// getNodePathCompletions returns completions for node paths.
func (s *Server) getNodePathCompletions(doc *analysis.Document) []protocol.CompletionItem {
	if doc.TSCNAST == nil {
		return nil
	}

	items := []protocol.CompletionItem{}
	kind := protocol.CompletionItemKindValue

	// Add "." for scene root
	items = append(items, protocol.CompletionItem{
		Label:  ".",
		Kind:   &kind,
		Detail: strPtr("Scene root"),
	})

	// Add all node paths
	for _, node := range doc.TSCNAST.Nodes {
		var path string
		if node.Parent == "" {
			continue // Skip root itself
		} else if node.Parent == "." {
			path = node.Name
		} else {
			path = node.Parent + "/" + node.Name
		}

		items = append(items, protocol.CompletionItem{
			Label:  path,
			Kind:   &kind,
			Detail: strPtr(node.Type),
		})
	}

	return items
}

// getValueCompletions returns completions for values after =.
func (s *Server) getValueCompletions(doc *analysis.Document) []protocol.CompletionItem {
	items := []protocol.CompletionItem{}
	kind := protocol.CompletionItemKindFunction

	// Common type constructors
	constructors := []struct {
		label  string
		insert string
		detail string
	}{
		{"Vector2", "Vector2($1, $2)", "2D vector"},
		{"Vector3", "Vector3($1, $2, $3)", "3D vector"},
		{"Vector4", "Vector4($1, $2, $3, $4)", "4D vector"},
		{"Color", "Color($1, $2, $3, $4)", "RGBA color"},
		{"Transform2D", "Transform2D($1, $2, $3, $4, $5, $6)", "2D transform"},
		{"Transform3D", "Transform3D($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)", "3D transform"},
		{"Quaternion", "Quaternion($1, $2, $3, $4)", "Rotation quaternion"},
		{"NodePath", "NodePath(\"$1\")", "Path to a node"},
		{"ExtResource", "ExtResource(\"$1\")", "External resource reference"},
		{"SubResource", "SubResource(\"$1\")", "Internal resource reference"},
	}

	snippetFormat := protocol.InsertTextFormatSnippet
	for _, c := range constructors {
		items = append(items, protocol.CompletionItem{
			Label:            c.label,
			Kind:             &kind,
			Detail:           strPtr(c.detail),
			InsertText:       strPtr(c.insert),
			InsertTextFormat: &snippetFormat,
		})
	}

	// Boolean values
	boolKind := protocol.CompletionItemKindKeyword
	items = append(items, protocol.CompletionItem{
		Label:  "true",
		Kind:   &boolKind,
		Detail: strPtr("Boolean true"),
	})
	items = append(items, protocol.CompletionItem{
		Label:  "false",
		Kind:   &boolKind,
		Detail: strPtr("Boolean false"),
	})
	items = append(items, protocol.CompletionItem{
		Label:  "null",
		Kind:   &boolKind,
		Detail: strPtr("Null value"),
	})

	return items
}

// getPropertyCompletions returns completions for property names.
func (s *Server) getPropertyCompletions() []protocol.CompletionItem {
	items := []protocol.CompletionItem{}
	kind := protocol.CompletionItemKindProperty

	// Common properties
	properties := []struct {
		label  string
		detail string
	}{
		{"transform", "Node transform (Transform2D or Transform3D)"},
		{"position", "Node position (Vector2 or Vector3)"},
		{"rotation", "Node rotation"},
		{"scale", "Node scale (Vector2 or Vector3)"},
		{"visible", "Node visibility"},
		{"modulate", "Color modulation"},
		{"z_index", "2D draw order"},
		{"process_mode", "Processing mode"},
		{"script", "Attached script"},
		{"mesh", "MeshInstance3D mesh"},
		{"shape", "CollisionShape shape"},
		{"texture", "Sprite texture"},
		{"material", "Material override"},
	}

	for _, p := range properties {
		items = append(items, protocol.CompletionItem{
			Label:  p.label,
			Kind:   &kind,
			Detail: strPtr(p.detail),
		})
	}

	return items
}
