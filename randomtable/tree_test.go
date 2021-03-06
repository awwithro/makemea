package randomtable

import (
	"reflect"
	"testing"
)

func TestRandomTableTree(t *testing.T) {
	expected := NewRandomTable()
	expectedRollingTable, err := NewRollingTable("1d6")
	if err != nil {
		t.Error(err)
	}
	tree := NewRandomTableTree()
	tree.AddTable("Test", &expected)
	tree.AddTable("TestRolling", &expectedRollingTable)
	actual, _ := tree.GetTable("Test")
	actualRolling, _ := tree.GetTable("TestRolling")
	if !reflect.DeepEqual(actual, &expected) {
		t.Errorf("Table %v didn't match table %v", actual, expected)
	}
	if !reflect.DeepEqual(actualRolling, &expectedRollingTable) {
		t.Errorf("Table %v didn't match table %v", actualRolling, expectedRollingTable)
	}
}
