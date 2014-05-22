package doable

type Item interface {
	Match(Item) bool
}

type StringItem struct {
	Value string
}

func (i StringItem) Match(other Item) bool {
	it, ok := other.(StringItem)
	return ok && i.Value == it.Value
}
