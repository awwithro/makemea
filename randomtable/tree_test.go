package randomtable

import (
	"testing"
)

func TestTree(t *testing.T) {
	expected := NewRandomTable()
	expectedRollingTable, err := NewRollingTable("1d6")
	if err != nil {
		t.Error(err)
	}
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
