package randomtable

type TextTable struct {
	text string
}

func (t *TextTable) GetItem() string {
	return t.text
}

func (t *TextTable) AddItem(item string, n ...int) {
	t.text += item
}

func NewTextTable() RandomTable {

	t := RandomTable{}
	return t
}
