package doable

import (
	"errors"
	"fmt"
)

type items map[Item]int

type List struct {
	l items
}

func NewList() *List {
	return &List{
		l: make(items),
	}
}

func (c *List) add(i Item, n int) {
	n += c.l[i]
	c.l[i] = n
}

func (c *List) add1(i Item) {
	c.add(i, 1)
}

func (c *List) count(i Item) int {
	return c.l[i]
}

func (c *List) del(i Item, n int) (err error) {
	n = c.l[i] - n
	if n > 0 {
		c.l[i] = n
		return nil
	} else if n == 0 {
		delete(c.l, i)
		return nil
	} else {
		msg := fmt.Sprintf("Not enough elements: %s < %s", n, c.l[i])
		return errors.New(msg)
	}
}

func (c *List) del1(i Item) (err error) {
	return c.del(i, 1)
}

func (c *List) size() int {
	return len(c.l)
}
