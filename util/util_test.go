package util

import "testing"

func TestDeDupe(t *testing.T) {
	a := []string{"one", "one", "two", "three"}
	expceted := []string{"one", "two", "three"}
	acutal := DeDupe(a)
	fail := false
	if len(acutal) != len(expceted) {
		fail = true
	} else {
		for i, v := range acutal {
			if expceted[i] != v {
				fail = true
			}
		}
	}
	if fail {
		t.Errorf("Expected: %v, Got: %v", expceted, acutal)
	}

}

func TestDifference(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{1, 2}
	diff := Difference(a, b)
	if len(diff) != 1 || diff[0] != 3 {
		t.Errorf("Expected: %v, Got: %v", []int{3}, diff)
	}
}
