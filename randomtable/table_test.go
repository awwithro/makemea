package randomtable

import "testing"

func TestRandomTable(t *testing.T) {
	r := NewRandomTable()
	r.AddItem("Hello")
	if r.GetItem() != "Hello" {
		t.Error("Didn't get expected item")
	}
}
