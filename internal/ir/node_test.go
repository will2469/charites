package ir_test

import (
	"strings"
	"testing"
	"unsafe"

	"github.com/will2469/charites/internal/ir"
)

func TestNode_Walk(t *testing.T) {
	// Membangun pohon uji pre-order:
	// root (div)
	// ├── child1 (span)
	// │   └── grandchild1 (text)
	// └── child2 (button)
	root := &ir.Node{
		Type: ir.NodeElement,
		Tag:  "div",
	}
	child1 := &ir.Node{
		Type:   ir.NodeElement,
		Tag:    "span",
		Parent: root,
	}
	grandchild1 := &ir.Node{
		Type:   ir.NodeText,
		Parent: child1,
	}
	child1.Children = []*ir.Node{grandchild1}

	child2 := &ir.Node{
		Type:   ir.NodeElement,
		Tag:    "button",
		Parent: root,
	}
	root.Children = []*ir.Node{child1, child2}

	var visited []string
	for n := range root.Walk() {
		if n.Type == ir.NodeText {
			visited = append(visited, "text")
		} else {
			visited = append(visited, n.Tag)
		}
	}

	expected := []string{"div", "span", "text", "button"}
	if len(visited) != len(expected) {
		t.Fatalf("expected %d nodes visited, got %d", len(expected), len(visited))
	}
	for i, v := range visited {
		if v != expected[i] {
			t.Errorf("at index %d: expected %s, got %s", i, expected[i], v)
		}
	}

	// Test nil node walk
	var nilNode *ir.Node
	nilCount := 0
	for range nilNode.Walk() {
		nilCount++
	}
	if nilCount != 0 {
		t.Errorf("expected 0 visits on nil node walk, got %d", nilCount)
	}
}

func TestNode_EarlyExit(t *testing.T) {
	root := &ir.Node{
		Type: ir.NodeElement,
		Tag:  "root",
	}
	for i := 0; i < 10; i++ {
		root.Children = append(root.Children, &ir.Node{
			Type:   ir.NodeElement,
			Tag:    "child",
			Parent: root,
		})
	}

	visitedCount := 0
	for range root.Walk() {
		visitedCount++
		if visitedCount == 3 {
			break
		}
	}

	if visitedCount != 3 {
		t.Errorf("expected early termination at 3 visits, got %d", visitedCount)
	}
}

func TestNode_ParentChildInvariant(t *testing.T) {
	root := &ir.Node{
		Type: ir.NodeElement,
		Tag:  "main",
	}
	child := &ir.Node{
		Type:   ir.NodeElement,
		Tag:    "section",
		Parent: root,
	}
	root.Children = []*ir.Node{child}

	if root.Parent != nil {
		t.Errorf("expected root.Parent to be nil, got %v", root.Parent)
	}
	if child.Parent != root {
		t.Errorf("expected child.Parent == root")
	}
	if len(root.Children) != 1 || root.Children[0] != child {
		t.Errorf("expected root.Children to contain child")
	}
}

func TestNode_ClassTokenization(t *testing.T) {
	raw := "p-4   bg-primary   text-sm\t\n"
	tokens := strings.Fields(raw)

	node := &ir.Node{
		Type:       ir.NodeElement,
		Tag:        "div",
		RawClasses: raw,
		Classes:    tokens,
	}

	expectedClasses := []string{"p-4", "bg-primary", "text-sm"}
	if len(node.Classes) != len(expectedClasses) {
		t.Fatalf("expected %d classes, got %d", len(expectedClasses), len(node.Classes))
	}
	for i, c := range node.Classes {
		if c != expectedClasses[i] {
			t.Errorf("at index %d: expected class %s, got %s", i, expectedClasses[i], c)
		}
	}
	if node.RawClasses != raw {
		t.Errorf("expected RawClasses to be preserved exactly, got %q", node.RawClasses)
	}
}

func TestNode_SpanIndexing(t *testing.T) {
	span := ir.Span{
		Line:      1,
		Column:    5,
		EndLine:   1,
		EndColumn: 25,
	}

	if span.Line < 1 || span.Column < 1 {
		t.Errorf("span must be 1-indexed, got Line=%d, Column=%d", span.Line, span.Column)
	}
	if span.EndLine < span.Line {
		t.Errorf("EndLine (%d) cannot be less than Line (%d)", span.EndLine, span.Line)
	}
	if span.EndLine == span.Line && span.EndColumn < span.Column {
		t.Errorf("EndColumn (%d) cannot be less than Column (%d) on the same line", span.EndColumn, span.Column)
	}
}

func TestNode_StructSize(t *testing.T) {
	nodeSize := unsafe.Sizeof(ir.Node{})
	if unsafe.Sizeof(uintptr(0)) == 8 { // 64-bit architecture
		if nodeSize > 136 {
			t.Errorf("ir.Node struct size exceeds 136-byte budget on 64-bit target: got %d bytes", nodeSize)
		}
	}
}

func TestNode_Helpers(t *testing.T) {
	node := &ir.Node{
		Type:       ir.NodeElement,
		Tag:        "button",
		RawClasses: "btn btn-primary",
		Classes:    []string{"btn", "btn-primary"},
		Attributes: map[string]string{
			"role":     "button",
			"disabled": "",
		},
	}

	// HasClass
	if !node.HasClass("btn-primary") {
		t.Errorf("expected HasClass('btn-primary') to be true")
	}
	if node.HasClass("btn-secondary") {
		t.Errorf("expected HasClass('btn-secondary') to be false")
	}

	// GetAttr
	val, ok := node.GetAttr("role")
	if !ok || val != "button" {
		t.Errorf("expected GetAttr('role') to return ('button', true), got (%q, %v)", val, ok)
	}
	val, ok = node.GetAttr("nonexistent")
	if ok || val != "" {
		t.Errorf("expected GetAttr('nonexistent') to return ('', false), got (%q, %v)", val, ok)
	}

	// GetAttr on nil attributes
	emptyNode := &ir.Node{}
	val, ok = emptyNode.GetAttr("role")
	if ok || val != "" {
		t.Errorf("expected GetAttr on empty node to return ('', false), got (%q, %v)", val, ok)
	}

	// IsElement
	if !node.IsElement() {
		t.Errorf("expected node.IsElement() with no args to be true for NodeElement")
	}
	if !node.IsElement("div", "button") {
		t.Errorf("expected node.IsElement('div', 'button') to be true")
	}
	if node.IsElement("div", "span") {
		t.Errorf("expected node.IsElement('div', 'span') to be false for button")
	}

	textNode := &ir.Node{Type: ir.NodeText}
	if textNode.IsElement() {
		t.Errorf("expected textNode.IsElement() to be false")
	}

	var nilNode *ir.Node
	if nilNode.IsElement() {
		t.Errorf("expected nilNode.IsElement() to be false")
	}
}

func buildMockSubtree(total int) *ir.Node {
	if total <= 0 {
		return nil
	}
	root := &ir.Node{
		Type: ir.NodeElement,
		Tag:  "root",
	}
	current := root
	for i := 1; i < total; i++ {
		child := &ir.Node{
			Type:   ir.NodeElement,
			Tag:    "item",
			Parent: current,
		}
		current.Children = append(current.Children, child)
		current = child
	}
	return root
}

func BenchmarkNode_Walk(b *testing.B) {
	root := buildMockSubtree(100)
	b.ReportAllocs()

	for b.Loop() {
		count := 0
		for range root.Walk() {
			count++
		}
		if count != 100 {
			b.Fatalf("expected 100 nodes, got %d", count)
		}
	}
}
