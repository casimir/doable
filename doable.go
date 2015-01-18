package doable

type Node struct {
	deps []*Node
	Item Item
	Nb   int
}

func (n *Node) AddDep(dep ...*Node) {
	n.deps = append(n.deps, dep...)
}

func (n *Node) ListDeps() *List {
	ret := NewList()
	for _, it := range n.deps {
		ret.AddN(it.Item, it.Nb)
	}
	return ret
}

// Tree represents a dependency tree and its context, including available
// dependencies and resolution history if doable or missing dpendencies if not.
type Tree struct {
	Avail *List
	Hist  []*Node
	Miss  *List
	root  *Node
}

func New(root *Node, list *List) *Tree {
	return &Tree{
		Avail: list,
		Hist:  []*Node{},
		root:  root,
	}
}

func (t *Tree) Doable() bool {
	if t.Avail == nil {
		return false
	}
	return t.process(t.root) == nil
}

func (t *Tree) process(n *Node) *Node {
	// Exists already.
	if t.Avail.Count(n.Item) > 0 {
		t.Hist = append(t.Hist, n)
		return nil
	}
	// Doesn't exist and can't be done.
	if n.deps == nil {
		return n
	}
	// 2 loops to avoid memory leak
	for i, d := range n.deps {
		tmp := n.deps[i]
		n.deps[i] = t.process(d)
		if n.deps[i] == nil {
			t.Avail.DelN(tmp.Item, tmp.Nb)
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
		t.Avail.Add(n.Item)
		t.Hist = append(t.Hist, n)
		return nil
	}
	t.Miss = n.ListDeps()
	return n
}
