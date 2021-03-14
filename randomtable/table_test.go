package randomtable

import "testing"

func TestRandomTable(t *testing.T) {
	r := NewRandomTable()
	r.AddItem("Hello")
	if r.GetItem() != "Hello" {
		t.Error("Didn't get expected item")
	}
	all := r.AllItems()
	if all[0] != "Hello" || len(all) != 1 {
		t.Error("Didn't get expected all result")
	}
}
