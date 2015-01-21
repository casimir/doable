package doable

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

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

// String implemenst Stringer for Node.
func (n Node) String() string {
	return fmt.Sprintf("%s (x%d)", n.Item.UID(), n.Nb)
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

// Dump save a representatio of the tree in the given file. The format used is
// DOT (http://en.wikipedia.org/wiki/DOT_%28graph_description_language%29).
func (t *Tree) Dump(path string) error {
	var headBuf bytes.Buffer
	var bodyBuf bytes.Buffer
	dump_rec(t.root, &headBuf, &bodyBuf, 1)

	var buf bytes.Buffer
	buf.WriteString("digraph " + t.root.Item.UID() + " {\n")
	buf.Write(headBuf.Bytes())
	buf.WriteRune('\n')
	buf.Write(bodyBuf.Bytes())
	buf.WriteRune('}')
	return ioutil.WriteFile(path, buf.Bytes(), 0644)
}

func dump_rec(n *Node, head, body *bytes.Buffer, lvl int) {
	label := fmt.Sprintf("  %s%d%d [label=\"%s\"];\n",
		n.Item.UID(), n.Nb, lvl, n)
	head.WriteString(label)

	for _, it := range n.deps {
		line := fmt.Sprintf("  %s%d%d -> %s%d%d;\n",
			n.Item.UID(), n.Nb, lvl,
			it.Item.UID(), it.Nb, lvl+1)
		body.WriteString(line)
		dump_rec(it, head, body, lvl+1)
	}
}
