package doable

type Node struct {
	deps []*Node
	item Item
}

func (n *Node) AddDep(dep ...*Node) {
	n.deps = append(n.deps, dep...)
}
