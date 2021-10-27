package randomtable

import (
	"github.com/olekukonko/tablewriter"
)

type TextTable struct {
	text string
}

func (t *TextTable) GetItem() string {
	return t.text
}

func (t *TextTable) AddItem(item string, n ...int) {
	t.text += item
}

func (t TextTable) AllItems() []string {
	return []string{t.text}
}

func (t TextTable) Validate() {

}

func (t TextTable) GetTable(tb *tablewriter.Table, name string) *tablewriter.Table {
	tb.Append([]string{t.text})
	tb.SetHeader([]string{t.text})
	return tb
}

func NewTextTable() TextTable {
	t := TextTable{}
	return t
}
