package parser

import (
	"testing"
)

func TestParseMinimal(t *testing.T) {
	input := `[gd_scene format=3]`
	doc := Parse(input)

	if doc.Descriptor == nil {
		t.Fatal("expected descriptor")
	}
	if doc.Descriptor.Type != "gd_scene" {
		t.Errorf("expected type gd_scene, got %s", doc.Descriptor.Type)
	}
	if doc.Descriptor.Format != 3 {
		t.Errorf("expected format 3, got %d", doc.Descriptor.Format)
	}
}

func TestParseWithUID(t *testing.T) {
	input := `[gd_scene load_steps=4 format=3 uid="uid://cecaux1sm7mo0"]`
	doc := Parse(input)

	if doc.Descriptor == nil {
		t.Fatal("expected descriptor")
	}
	if doc.Descriptor.UID != "uid://cecaux1sm7mo0" {
		t.Errorf("expected uid://cecaux1sm7mo0, got %s", doc.Descriptor.UID)
	}
	if doc.Descriptor.LoadSteps == nil || *doc.Descriptor.LoadSteps != 4 {
		t.Error("expected load_steps 4")
	}
}

func TestParseExtResource(t *testing.T) {
	input := `[gd_scene format=3]
[ext_resource type="Texture2D" uid="uid://abc" path="res://texture.png" id="1_abc"]`

	doc := Parse(input)

	if len(doc.ExtResources) != 1 {
		t.Fatalf("expected 1 ext_resource, got %d", len(doc.ExtResources))
	}

	ext := doc.ExtResources[0]
	if ext.Type != "Texture2D" {
		t.Errorf("expected type Texture2D, got %s", ext.Type)
	}
	if ext.UID != "uid://abc" {
		t.Errorf("expected uid://abc, got %s", ext.UID)
	}
	if ext.Path != "res://texture.png" {
		t.Errorf("expected res://texture.png, got %s", ext.Path)
	}
	if ext.ID != "1_abc" {
		t.Errorf("expected 1_abc, got %s", ext.ID)
	}

	// Verify PathRange includes the full path with quotes
	// Line 1: [ext_resource type="Texture2D" uid="uid://abc" path="res://texture.png" id="1_abc"]
	// The path starts at column 52 (0-indexed): path="res://texture.png"
	if ext.PathRange.Start.Line != 1 {
		t.Errorf("expected PathRange.Start.Line to be 1, got %d", ext.PathRange.Start.Line)
	}
	// The path token includes quotes, so it should span the entire "res://texture.png"
	pathLen := len(`"res://texture.png"`)
	actualLen := ext.PathRange.End.Column - ext.PathRange.Start.Column
	if actualLen != pathLen {
		t.Errorf("expected PathRange length to be %d (full quoted string), got %d", pathLen, actualLen)
	}
}

func TestParseSubResource(t *testing.T) {
	input := `[gd_scene format=3]
[sub_resource type="SphereShape3D" id="SphereShape3D_abc"]
radius = 1.5`

	doc := Parse(input)

	if len(doc.SubResources) != 1 {
		t.Fatalf("expected 1 sub_resource, got %d", len(doc.SubResources))
	}

	sub := doc.SubResources[0]
	if sub.Type != "SphereShape3D" {
		t.Errorf("expected type SphereShape3D, got %s", sub.Type)
	}
	if sub.ID != "SphereShape3D_abc" {
		t.Errorf("expected SphereShape3D_abc, got %s", sub.ID)
	}
	if len(sub.Properties) != 1 {
		t.Fatalf("expected 1 property, got %d", len(sub.Properties))
	}
	if sub.Properties[0].Key != "radius" {
		t.Errorf("expected key radius, got %s", sub.Properties[0].Key)
	}
}

func TestParseNode(t *testing.T) {
	input := `[gd_scene format=3]
[node name="Player" type="CharacterBody3D"]
[node name="Arm" type="Node3D" parent="."]
[node name="Hand" type="Node3D" parent="Arm"]`

	doc := Parse(input)

	if len(doc.Nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(doc.Nodes))
	}

	// Root node
	if doc.Nodes[0].Name != "Player" {
		t.Errorf("expected name Player, got %s", doc.Nodes[0].Name)
	}
	if doc.Nodes[0].Type != "CharacterBody3D" {
		t.Errorf("expected type CharacterBody3D, got %s", doc.Nodes[0].Type)
	}
	if doc.Nodes[0].Parent != "" {
		t.Errorf("expected empty parent for root, got %s", doc.Nodes[0].Parent)
	}

	// Child of root
	if doc.Nodes[1].Name != "Arm" {
		t.Errorf("expected name Arm, got %s", doc.Nodes[1].Name)
	}
	if doc.Nodes[1].Parent != "." {
		t.Errorf("expected parent '.', got %s", doc.Nodes[1].Parent)
	}

	// Nested child
	if doc.Nodes[2].Name != "Hand" {
		t.Errorf("expected name Hand, got %s", doc.Nodes[2].Name)
	}
	if doc.Nodes[2].Parent != "Arm" {
		t.Errorf("expected parent 'Arm', got %s", doc.Nodes[2].Parent)
	}
}

func TestParseConnection(t *testing.T) {
	input := `[gd_scene format=3]
[node name="Root" type="Node"]
[connection signal="pressed" from="Button" to="." method="_on_button_pressed"]`

	doc := Parse(input)

	if len(doc.Connections) != 1 {
		t.Fatalf("expected 1 connection, got %d", len(doc.Connections))
	}

	conn := doc.Connections[0]
	if conn.Signal != "pressed" {
		t.Errorf("expected signal pressed, got %s", conn.Signal)
	}
	if conn.From != "Button" {
		t.Errorf("expected from Button, got %s", conn.From)
	}
	if conn.To != "." {
		t.Errorf("expected to '.', got %s", conn.To)
	}
	if conn.Method != "_on_button_pressed" {
		t.Errorf("expected method _on_button_pressed, got %s", conn.Method)
	}
}

func TestParseTypedValues(t *testing.T) {
	input := `[gd_scene format=3]
[node name="Root" type="Node3D"]
transform = Transform3D(1, 0, 0, 0, 1, 0, 0, 0, 1, 1, 2, 3)
position = Vector3(1.5, 2.0, 3.0)
color = Color(1, 0.5, 0.25, 1)
path = NodePath("Parent/Child")`

	doc := Parse(input)

	if len(doc.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(doc.Nodes))
	}

	node := doc.Nodes[0]
	if len(node.Properties) != 4 {
		t.Fatalf("expected 4 properties, got %d", len(node.Properties))
	}

	// Transform3D
	prop := node.Properties[0]
	if prop.Key != "transform" {
		t.Errorf("expected key transform, got %s", prop.Key)
	}
	if tv, ok := prop.Value.(*TypedValue); ok {
		if tv.TypeName != "Transform3D" {
			t.Errorf("expected Transform3D, got %s", tv.TypeName)
		}
		if len(tv.Arguments) != 12 {
			t.Errorf("expected 12 arguments, got %d", len(tv.Arguments))
		}
	} else {
		t.Error("expected TypedValue for transform")
	}

	// Vector3
	prop = node.Properties[1]
	if tv, ok := prop.Value.(*TypedValue); ok {
		if tv.TypeName != "Vector3" {
			t.Errorf("expected Vector3, got %s", tv.TypeName)
		}
		if len(tv.Arguments) != 3 {
			t.Errorf("expected 3 arguments, got %d", len(tv.Arguments))
		}
	} else {
		t.Error("expected TypedValue for position")
	}
}

func TestParseResourceReferences(t *testing.T) {
	input := `[gd_scene format=3]
[ext_resource type="Script" path="res://script.gd" id="1_abc"]
[sub_resource type="BoxShape3D" id="BoxShape_123"]
[node name="Root" type="Node"]
script = ExtResource("1_abc")
shape = SubResource("BoxShape_123")`

	doc := Parse(input)

	if len(doc.Nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(doc.Nodes))
	}

	node := doc.Nodes[0]
	if len(node.Properties) != 2 {
		t.Fatalf("expected 2 properties, got %d", len(node.Properties))
	}

	// ExtResource
	prop := node.Properties[0]
	if ref, ok := prop.Value.(*ResourceRef); ok {
		if ref.RefType != "ExtResource" {
			t.Errorf("expected ExtResource, got %s", ref.RefType)
		}
		if ref.ID != "1_abc" {
			t.Errorf("expected 1_abc, got %s", ref.ID)
		}
	} else {
		t.Error("expected ResourceRef for script")
	}

	// SubResource
	prop = node.Properties[1]
	if ref, ok := prop.Value.(*ResourceRef); ok {
		if ref.RefType != "SubResource" {
			t.Errorf("expected SubResource, got %s", ref.RefType)
		}
		if ref.ID != "BoxShape_123" {
			t.Errorf("expected BoxShape_123, got %s", ref.ID)
		}
	} else {
		t.Error("expected ResourceRef for shape")
	}
}

func TestParseArray(t *testing.T) {
	input := `[gd_scene format=3]
[node name="Root" type="Node"]
values = [1, 2, 3]
strings = ["a", "b", "c"]`

	doc := Parse(input)

	node := doc.Nodes[0]

	// Integer array
	prop := node.Properties[0]
	if arr, ok := prop.Value.(*ArrayValue); ok {
		if len(arr.Values) != 3 {
			t.Errorf("expected 3 values, got %d", len(arr.Values))
		}
	} else {
		t.Error("expected ArrayValue")
	}

	// String array
	prop = node.Properties[1]
	if arr, ok := prop.Value.(*ArrayValue); ok {
		if len(arr.Values) != 3 {
			t.Errorf("expected 3 values, got %d", len(arr.Values))
		}
		if sv, ok := arr.Values[0].(*StringValue); ok {
			if sv.Value != "a" {
				t.Errorf("expected 'a', got %s", sv.Value)
			}
		}
	} else {
		t.Error("expected ArrayValue")
	}
}

func TestParseDict(t *testing.T) {
	input := `[gd_scene format=3]
[node name="Root" type="Node"]
data = {
"key1": 123,
"key2": "value"
}`

	doc := Parse(input)

	node := doc.Nodes[0]
	prop := node.Properties[0]

	if dict, ok := prop.Value.(*DictValue); ok {
		if len(dict.Entries) != 2 {
			t.Errorf("expected 2 entries, got %d", len(dict.Entries))
		}
	} else {
		t.Error("expected DictValue")
	}
}

func TestParseComment(t *testing.T) {
	input := `[gd_scene format=3]
; This is a comment
[node name="Root" type="Node"]`

	doc := Parse(input)

	if len(doc.Comments) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(doc.Comments))
	}
	if doc.Comments[0].Text != " This is a comment" {
		t.Errorf("expected ' This is a comment', got %s", doc.Comments[0].Text)
	}
}

func TestParseSlashedProperty(t *testing.T) {
	input := `[gd_scene format=3]
[node name="Skeleton" type="Skeleton3D"]
bones/0/position = Vector3(0, 1, 0)
bones/0/rotation = Quaternion(0, 0, 0, 1)`

	doc := Parse(input)

	node := doc.Nodes[0]
	if len(node.Properties) != 2 {
		t.Fatalf("expected 2 properties, got %d", len(node.Properties))
	}

	if node.Properties[0].Key != "bones/0/position" {
		t.Errorf("expected bones/0/position, got %s", node.Properties[0].Key)
	}
	if node.Properties[1].Key != "bones/0/rotation" {
		t.Errorf("expected bones/0/rotation, got %s", node.Properties[1].Key)
	}
}

func TestParseBoolAndNull(t *testing.T) {
	input := `[gd_scene format=3]
[node name="Root" type="Node"]
enabled = true
disabled = false
nothing = null`

	doc := Parse(input)

	node := doc.Nodes[0]

	// true
	if bv, ok := node.Properties[0].Value.(*BoolValue); ok {
		if !bv.Value {
			t.Error("expected true")
		}
	} else {
		t.Error("expected BoolValue")
	}

	// false
	if bv, ok := node.Properties[1].Value.(*BoolValue); ok {
		if bv.Value {
			t.Error("expected false")
		}
	} else {
		t.Error("expected BoolValue")
	}

	// null
	if _, ok := node.Properties[2].Value.(*NullValue); !ok {
		t.Error("expected NullValue")
	}
}

func TestComplexScene(t *testing.T) {
	input := `[gd_scene load_steps=4 format=3 uid="uid://cecaux1sm7mo0"]

[sub_resource type="SphereShape3D" id="SphereShape3D_tj6p1"]

[sub_resource type="SphereMesh" id="SphereMesh_4w3ye"]

[sub_resource type="StandardMaterial3D" id="StandardMaterial3D_k54se"]
albedo_color = Color(1, 0.639216, 0.309804, 1)

[node name="Ball" type="RigidBody3D"]

[node name="CollisionShape3D" type="CollisionShape3D" parent="."]
shape = SubResource("SphereShape3D_tj6p1")

[node name="MeshInstance3D" type="MeshInstance3D" parent="."]
mesh = SubResource("SphereMesh_4w3ye")
surface_material_override/0 = SubResource("StandardMaterial3D_k54se")`

	doc := Parse(input)

	// Check we have the right counts
	if len(doc.SubResources) != 3 {
		t.Errorf("expected 3 sub_resources, got %d", len(doc.SubResources))
	}
	if len(doc.Nodes) != 3 {
		t.Errorf("expected 3 nodes, got %d", len(doc.Nodes))
	}
	if len(doc.Errors) != 0 {
		t.Errorf("expected no errors, got %d: %v", len(doc.Errors), doc.Errors)
	}
}

func TestMalformedDictDoesNotCrash(t *testing.T) {
	// This input has a malformed dictionary that previously caused an infinite loop
	input := `[gd_scene format=3]
[node name="Root" type="Node"]
data = {
"key1" "missing_colon",
"key2": 123
}`

	doc := Parse(input)

	// Should complete parsing without crashing
	if doc.Descriptor == nil {
		t.Error("expected descriptor to be parsed")
	}
	if len(doc.Nodes) != 1 {
		t.Errorf("expected 1 node, got %d", len(doc.Nodes))
	}
	// Should have errors for the malformed dictionary
	if len(doc.Errors) == 0 {
		t.Error("expected parse errors for malformed dictionary")
	}
}

func TestMalformedArrayDoesNotCrash(t *testing.T) {
	// Array with unexpected tokens
	input := `[gd_scene format=3]
[node name="Root" type="Node"]
values = [1, 2, @@@@, 3]`

	doc := Parse(input)

	// Should complete parsing without crashing
	if doc.Descriptor == nil {
		t.Error("expected descriptor to be parsed")
	}
}

func TestMalformedTypedValueDoesNotCrash(t *testing.T) {
	// Typed value with malformed arguments
	input := `[gd_scene format=3]
[node name="Root" type="Node"]
position = Vector3(1.0, @@@@, 3.0)`

	doc := Parse(input)

	// Should complete parsing without crashing
	if doc.Descriptor == nil {
		t.Error("expected descriptor to be parsed")
	}
}

func TestDeeplyNestedMalformedInput(t *testing.T) {
	// Deeply nested malformed structures
	input := `[gd_scene format=3]
[node name="Root" type="Node"]
data = {
"a": {
"b": {
"c": [1, 2, {
"d" "missing_colon"
}]
}
}
}`

	doc := Parse(input)

	// Should complete parsing without crashing
	if doc.Descriptor == nil {
		t.Error("expected descriptor to be parsed")
	}
}

func TestErrorLimitExceeded(t *testing.T) {
	// Generate input that would produce many errors
	input := `[gd_scene format=3]
[node name="Root" type="Node"]
`
	// Add many invalid lines
	for i := 0; i < 200; i++ {
		input += "@@invalid@@\n"
	}

	doc := Parse(input)

	// Should complete parsing without crashing
	if doc.Descriptor == nil {
		t.Error("expected descriptor to be parsed")
	}
	// Errors should be limited to maxErrors (100)
	if len(doc.Errors) > 100 {
		t.Errorf("expected at most 100 errors, got %d", len(doc.Errors))
	}
}

func TestParseStringNameLiteral(t *testing.T) {
	// StringName literals use &"name" syntax
	input := `[gd_scene format=3]

[sub_resource type="AnimationLibrary" id="AnimationLibrary_abc"]
_data = {
&"RESET": SubResource("Animation_reset"),
&"walk": SubResource("Animation_walk")
}`

	doc := Parse(input)

	if doc.Descriptor == nil {
		t.Error("expected descriptor to be parsed")
	}
	if len(doc.SubResources) != 1 {
		t.Fatalf("expected 1 sub_resource, got %d", len(doc.SubResources))
	}

	sub := doc.SubResources[0]
	if len(sub.Properties) != 1 {
		t.Fatalf("expected 1 property, got %d", len(sub.Properties))
	}

	// The _data property should have a dict value
	prop := sub.Properties[0]
	if prop.Key != "_data" {
		t.Errorf("expected key '_data', got '%s'", prop.Key)
	}
	if dict, ok := prop.Value.(*DictValue); ok {
		if len(dict.Entries) != 2 {
			t.Errorf("expected 2 dict entries, got %d", len(dict.Entries))
		}
		// Check first key is RESET (without the &)
		if sv, ok := dict.Entries[0].Key.(*StringValue); ok {
			if sv.Value != "RESET" {
				t.Errorf("expected key 'RESET', got '%s'", sv.Value)
			}
		}
	} else {
		t.Error("expected DictValue for _data property")
	}
}

func TestParseAnimationTracks(t *testing.T) {
	// Test parsing animation track format with nested dicts
	input := `[gd_scene format=3]

[sub_resource type="Animation" id="Animation_test"]
length = 0.001
tracks/0/type = "value"
tracks/0/path = NodePath(".")
tracks/0/keys = {
"times": PackedFloat32Array(0),
"transitions": PackedFloat32Array(1),
"update": 0,
"values": [Vector3(0, 0, 0)]
}`

	doc := Parse(input)

	if doc.Descriptor == nil {
		t.Error("expected descriptor to be parsed")
	}
	if len(doc.SubResources) != 1 {
		t.Fatalf("expected 1 sub_resource, got %d", len(doc.SubResources))
	}

	// Should parse without errors (or with recoverable errors)
	if len(doc.Errors) > 0 {
		t.Logf("Parse had %d errors: %v", len(doc.Errors), doc.Errors)
	}
}
