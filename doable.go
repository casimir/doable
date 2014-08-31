package doable

type Node struct {
	deps []*Node
	Item Item
	Nb   int
}

func (n *Node) AddDep(dep ...*Node) {
	n.deps = append(n.deps, dep...)
}

type Tree struct {
	Avail *List
	Hist  []*Node
	Root  *Node
}

func New(root *Node, list *List) *Tree {
	return &Tree{
		Avail: list,
		Hist:  []*Node{},
		Root:  root,
	}
}

func (t *Tree) Doable() bool {
	if t.Avail == nil {
		return false
	}
	return t.process(t.Root) == nil
}

func (t *Tree) process(n *Node) *Node {
	if t.Avail.count(n.Item) > 0 {
		t.Hist = append(t.Hist, n)
		return nil
	}
	if n.deps == nil {
		return n
	}
	// 2 loops to avoid memory leak
	for i, d := range n.deps {
		tmp := n.deps[i]
		n.deps[i] = t.process(d)
		if n.deps[i] == nil {
			t.Avail.del(tmp.Item, tmp.Nb)
		}
	}
	tmp := []*Node{}
	for _, d := range n.deps {
		if d != nil {
			tmp = append(tmp, d)
		}
	}
	n.deps = tmp
	if len(n.deps) == 0 {
		t.Avail.add1(n.Item)
		t.Hist = append(t.Hist, n)
		return nil
	}
	return n
}
