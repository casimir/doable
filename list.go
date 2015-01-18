package doable

import "fmt"

type (
	Item interface {
		Match(Item) bool
	}

	items map[Item]int

	List struct {
		l items
	}
)

func NewList() *List {
	return &List{l: make(items)}
}

func (l *List) Add(i Item) {
	l.AddN(i, 1)
}

func (l *List) AddN(i Item, n int) {
	n += l.l[i]
	l.l[i] = n
}

func (l *List) Clone() *List {
	il := make(items, len(l.l))
	for k, v := range l.l {
		il[k] = v
	}
	return &List{l: il}
}

func (l *List) Count(i Item) int {
	return l.l[i]
}

func (l *List) Del(i Item) (err error) {
	return l.DelN(i, 1)
}

func (l *List) DelN(i Item, n int) (err error) {
	n = l.l[i] - n
	if n > 0 {
		l.l[i] = n
		return nil
	} else if n == 0 {
		delete(l.l, i)
		return nil
	} else {
		return fmt.Errorf("Not enough elements: %s < %s", n, l.l[i])
	}
}

func (l *List) Size() int {
	return len(l.l)
}

type StringItem struct {
	Value string
}

func (i StringItem) Match(other Item) bool {
	it, ok := other.(StringItem)
	return ok && i.Value == it.Value
}
