package randomtable

import (
	"testing"
)

func TestTree(t *testing.T) {
	expected := NewRandomTable()
	expectedRollingTable := NewRollingTable("1d6")
	tree := NewTree()
	tree.AddTable("Test", &expected)
	tree.AddTable("TestRolling", &expectedRollingTable)
	// actual, _ := tree.GetTable("Test")
	// if actual == nil {
	// 	t.Error("Test table not found")
	// }
	// actualRolling, _ := tree.GetTable("TestRolling")
	// if actualRolling == nil {
	// 	t.Error("TestRolling table not found")
	// }
}
