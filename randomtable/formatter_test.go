package randomtable

import "testing"

func TestHtmlFormatter(t *testing.T) {
	f := HtmlFormatter{}
	actual := f.Format("one", "two")
	expected := "<randomElement table='two'>one</randomElement>"
	if actual != expected {
		t.Errorf("Actual: %v did not equal Expected: %v", actual, expected)
	}
}
