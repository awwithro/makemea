package randomtable

import (
	"testing"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
)

func TestRenderer(t *testing.T) {
	tree := NewTree()
	r := &randomTableRenderer{
		nodeRendererFuncsTmp: map[ast.NodeKind]renderer.NodeRendererFunc{},
		tree:                 tree,
		namespace:            []string{},
		depth:                0,
		currentTableNames:    []string{},
	}
	if r.Namespace() != "" {
		t.Errorf("New Renderer should have an empty namespace. Got: %v", r.Namespace())
	}
	r.Push("test")
	actual := r.Name("test")
	expected := "test/test"
	if actual != expected {
		t.Errorf("Namespace rendered incorrectly. Expected: %v Got: %v", expected, actual)
	}
	r.Pop()
	actual = r.Name("test")
	expected = "test"
	if actual != expected {
		t.Errorf("Namespace rendered incorrectly. Expected: %v Got: %v", expected, actual)
	}
}
