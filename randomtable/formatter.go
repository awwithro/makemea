package randomtable

import "fmt"

// Format will format the first string.
// Second string is the table the item is looked up on
type Formatter interface {
	Format(string, string) string
}

// NoOp Formatter for the base case
type StringFormatter struct {
}

func (s StringFormatter) Format(input string, caller string) string {
	return input
}

// Wraps items in html tags with metadata
type HtmlFormatter struct {
}

func (s HtmlFormatter) Format(input string, caller string) string {
	return fmt.Sprintf("<RandomElement table='%s'>%s</RandomElement>", caller, input)
}

func (t Tree) WithHtmlFormatter() Tree {
	t.formatter = HtmlFormatter{}
	return t
}

func (t Tree) WithStringFormatter() Tree {
	t.formatter = StringFormatter{}
	return t
}
