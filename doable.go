package doable

type Tree struct {
	// List of available items. It should be set before a call to Doable().
	Avail []Item
	Root  *Node
}

func New(root *Node, list []Item) *Tree {
	return &Tree{list, root}
}

func (t *Tree) Doable() bool {
	return t.process(t.Root) == nil
}

func (t *Tree) process(n *Node) *Node {
	for i, it := range t.Avail {
		if it.Match(n.item) {
			t.Avail = append(t.Avail[:i], t.Avail[i+1:]...)
			return nil
		}
	}
	if n.deps == nil {
		return n
	}
	// Mandatory 2 loops to avoid memory leak
	for i, d := range n.deps {
		n.deps[i] = t.process(d)
	}
	tmp := []*Node{}
	for _, d := range n.deps {
		if d != nil {
			tmp = append(tmp, d)
		}
	}
	n.deps = tmp
	if len(n.deps) == 0 {
		return nil
	} else {
		return n
	}
}
