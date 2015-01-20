package doable

import "fmt"

type (
	Item interface {
		ID() string
		Match(Item) bool
	}

	tuple struct {
		item Item
		n    int
	}

	items map[string]tuple

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
	if it, ok := l.l[i.ID()]; ok {
		n += it.n
	}
	l.l[i.ID()] = tuple{item: i, n: n}
}

func (l *List) Clone() *List {
	il := make(items, len(l.l))
	for k, v := range l.l {
		il[k] = v
	}
	return &List{l: il}
}

func (l *List) Count(i Item) int {
	if it, ok := l.l[i.ID()]; ok {
		return it.n
	}
	return 0
}

func (l *List) Del(i Item) (err error) {
	return l.DelN(i, 1)
}

func (l *List) DelN(i Item, n int) (err error) {
	n *= -1

	if it, ok := l.l[i.ID()]; ok {
		n = it.n + n
	}
	if n > 0 {
		l.l[i.ID()] = tuple{item: i, n: n}
		return nil
	} else if n == 0 {
		delete(l.l, i.ID())
		return nil
	} else {
		return fmt.Errorf("Not enough elements: %s < %s", n, l.l[i.ID()])
	}
}

func (l *List) Size() int {
	return len(l.l)
}

type StringItem struct {
	Value string
}

func (i StringItem) ID() string {
	return i.Value
}

func (i StringItem) Match(other Item) bool {
	it, ok := other.(StringItem)
	return ok && i.Value == it.Value
}
