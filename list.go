package doable

import "fmt"

type (
	Item interface {
		UID() string
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
	if it, ok := l.l[i.UID()]; ok {
		n += it.n
	}
	l.l[i.UID()] = tuple{item: i, n: n}
}

func (l *List) Clone() *List {
	il := make(items, len(l.l))
	for k, v := range l.l {
		il[k] = v
	}
	return &List{l: il}
}

func (l *List) Count(i Item) int {
	if it, ok := l.l[i.UID()]; ok {
		return it.n
	}
	return 0
}

func (l *List) Del(i Item) error {
	return l.DelN(i, 1)
}

func (l *List) DelN(i Item, n int) error {
	n *= -1

	if it, ok := l.l[i.UID()]; ok {
		n = it.n + n
	}
	if n > 0 {
		l.l[i.UID()] = tuple{item: i, n: n}
		return nil
	} else if n == 0 {
		delete(l.l, i.UID())
		return nil
	} else {
		return fmt.Errorf("Not enough elements: %s < %s", n, l.l[i.UID()])
	}
}

func (l *List) Size() int {
	return len(l.l)
}

type StringItem struct {
	Value string
}

func (i StringItem) UID() string {
	return i.Value
}

func (i StringItem) Match(other Item) bool {
	it, ok := other.(StringItem)
	return ok && i.Value == it.Value
}
