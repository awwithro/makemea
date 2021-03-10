package randomtable

import "testing"

func TestRollingTable(t *testing.T) {
	r := NewRollingTable("2d4")

	r.AddItem("Hello", 2, 3, 4, 5, 6, 7, 8)
	if r.GetItem() != "Hello" {
		t.Error("Didn't get expected item")
	}
}
