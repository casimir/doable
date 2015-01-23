package doable

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNoDep(t *testing.T) {
	item1 := StringItem{"id1"}
	item2 := StringItem{"id2"}
	itemNotUsed := StringItem{"not used"}

	Convey("Given an Item and an Item list", t, func() {
		list := NewList()
		list.Add(item1)
		list.Add(item2)
		list.Add(itemNotUsed)

		Convey("It should return if it is doable", func() {
			tree := New(&Node{Item: item1, Nb: 1}, nil)
			So(tree.Doable(), ShouldBeFalse)
			So(tree.Hist, ShouldBeEmpty)
			tree.Avail = NewList()
			So(tree.Doable(), ShouldBeFalse)
			So(tree.Hist, ShouldBeEmpty)

			tree.Avail = list
			So(tree.Doable(), ShouldBeTrue)
			So(tree.Avail, ShouldResemble, list)
			tree.root = &Node{Item: MockItem{}, Nb: 1}
			So(tree.Doable(), ShouldBeFalse)
			tree.root = &Node{Item: StringItem{"wrongId"}, Nb: 1}
			So(tree.Doable(), ShouldBeFalse)
		})

		Convey("When it is doable it should update the history", func() {
			list = NewList()
			list.Add(item1)
			list.Add(itemNotUsed)

			testList := NewList()
			testList.Add(item1)
			testList.Add(itemNotUsed)

			tree := New(&Node{Item: item1, Nb: 1}, list)
			hist := []*Node{tree.root}

			So(tree.Doable(), ShouldBeTrue)
			So(tree.Avail, ShouldResemble, testList)
			So(tree.Hist, ShouldResemble, hist)
		})
	})
}

func TestDeps(t *testing.T) {
	item0 := StringItem{"id0"}
	item1 := StringItem{"id1"}
	item2 := StringItem{"id2"}
	item3 := StringItem{"id3"}
	item4 := StringItem{"id4"}
	itemNotUsed := StringItem{"not used"}
	itemNotHere := StringItem{"not here"}
	itemRoot := StringItem{"root"}

	Convey("Given an Item and an Item list", t, func() {
		Convey("When dependencies are missing it should return it is not doable", func() {
			list := NewList()
			list.Add(item1)
			list.Add(itemNotUsed)

			root := &Node{Item: StringItem{"root"}, Nb: 1}
			root.AddDep(
				&Node{Item: item1, Nb: 1},
				&Node{Item: itemNotHere, Nb: 1},
			)

			tree := New(root, list)

			missing := NewList()
			missing.Add(itemNotHere)

			expected := &Node{Item: root.Item, Nb: 1}
			expected.AddDep(&Node{Item: itemNotHere, Nb: 1})

			So(tree.Doable(), ShouldBeFalse)
			So(tree.Avail, ShouldResemble, list)
			So(tree.Miss, ShouldResemble, missing)
			So(tree.root, ShouldResemble, expected)
		})

		Convey("When there is too few dependencies it should return it is not doable", func() {
			list := NewList()
			list.Add(item1)
			list.AddN(item2, 2)

			root := &Node{Item: StringItem{"root"}, Nb: 1}
			root.AddDep(
				&Node{Item: item1, Nb: 1},
				&Node{Item: item2, Nb: 5},
			)

			tree := New(root, list)

			missing := NewList()
			missing.AddN(item2, 3)

			expected := &Node{Item: root.Item, Nb: 1}
			expected.AddDep(&Node{Item: item2, Nb: 3})

			So(tree.Doable(), ShouldBeFalse)
			So(tree.Avail, ShouldResemble, list)
			So(tree.Miss, ShouldResemble, missing)
			So(tree.root, ShouldResemble, expected)
		})

		Convey("When dependencies are satisfied it should return it is doable", func() {
			list := NewList()
			list.Add(item0)
			list.Add(item1)
			list.Add(item2)
			list.Add(item3)
			list.Add(item4)
			list.Add(itemNotUsed)

			testList := NewList()
			testList.Add(itemRoot)
			testList.Add(itemNotUsed)

			node1 := &Node{Item: StringItem{"node 1"}, Nb: 1}
			node1.AddDep(
				&Node{Item: item2, Nb: 1},
				&Node{Item: item3, Nb: 1},
			)

			node2 := &Node{Item: StringItem{"node 2"}, Nb: 1}
			node2.AddDep(&Node{Item: item4, Nb: 1}, node1)

			root := &Node{Item: itemRoot, Nb: 1}
			root.AddDep(
				&Node{Item: item0, Nb: 1},
				&Node{Item: item1, Nb: 1},
				node2,
			)

			// root
			//  ↳ id0, id1, node2
			//               ↳ id4, node1
			//                       ↳ id2, id3
			tree := New(root, list)
			hist := []*Node{
				&Node{Item: item0, Nb: 1},
				&Node{Item: item1, Nb: 1},
				&Node{Item: item4, Nb: 1},
				&Node{Item: item2, Nb: 1},
				&Node{Item: item3, Nb: 1},
				node1,
				node2,
				root,
			}

			So(tree.Doable(), ShouldBeTrue)
			So(tree.Avail, ShouldResemble, testList)
			So(tree.Hist, ShouldResemble, hist)
		})
	})

	Convey("Given an Item and an empty Item list", t, func() {
		root := &Node{Item: itemRoot, Nb: 1}
		root.AddDep(
			&Node{Item: item0, Nb: 1},
			&Node{Item: item1, Nb: 2},
		)
		tree := New(root, NewList())

		So(tree.Doable(), ShouldBeFalse)
		So(tree.Avail, ShouldResemble, NewList())
		So(tree.Hist, ShouldBeEmpty)
		So(tree.Miss, ShouldResemble, root.listDeps(nil))
	})
}

func TestMulti(t *testing.T) {
	item0 := StringItem{"id0"}
	item1 := StringItem{"id1"}
	itemRoot := StringItem{"root"}

	Convey("Given a node with Nb > 1", t, func() {
		list := NewList()
		list.AddN(item0, 2)
		list.AddN(item1, 4)

		testList := NewList()
		testList.Add(itemRoot)

		Convey("It should resolve the dependencies", func() {
			root := &Node{Item: itemRoot, Nb: 1}
			root.AddDep(
				&Node{Item: item0, Nb: 2},
				&Node{Item: item1, Nb: 4},
			)
			tree := New(root, list)

			So(tree.Doable(), ShouldBeTrue)
			So(tree.Avail, ShouldResemble, testList)
		})
	})
}

func TestStringItem(t *testing.T) {
	item1 := StringItem{"id1"}
	item1b := StringItem{"id1"}
	item2 := StringItem{"id2"}

	Convey("It should match similar items", t, func() {
		So(item1.Match(item1), ShouldBeTrue)
		So(item1.Match(item1b), ShouldBeTrue)
		So(item1.Match(item2), ShouldBeFalse)
		So(item1.Match(MockItem{}), ShouldBeFalse)
	})
}

func TestDump(t *testing.T) {
	expected := []byte(exportedTree)
	outFile := os.TempDir() + "/doable_test.dot"

	root := &Node{Item: StringItem{"root"}, Nb: 2}
	dep1 := &Node{Item: StringItem{"item1"}, Nb: 4}
	dep2 := &Node{Item: StringItem{"item2"}, Nb: 1}
	dep3 := &Node{Item: StringItem{"item4"}, Nb: 2}
	dep4 := &Node{Item: StringItem{"item1"}, Nb: 2}

	dep3.AddDep(dep4)
	dep2.AddDep(dep3)
	root.AddDep(dep1, dep2)
	tree := New(root, nil)

	Convey("It should generate a DOT file for the whole tree", t, func() {
		os.Remove(outFile)
		So(tree.Dump(outFile), ShouldBeNil)
		got, _ := ioutil.ReadFile(outFile)
		So(got, ShouldResemble, expected)
	})
}

type MockItem struct{}

func (i MockItem) UID() string {
	return ""
}

func (i MockItem) Match(other Item) bool {
	return false
}

var exportedTree string = `digraph root {
  root21 [label="root (x2)"];
  item142 [label="item1 (x4)"];
  item212 [label="item2 (x1)"];
  item423 [label="item4 (x2)"];
  item124 [label="item1 (x2)"];

  root21 -> item142;
  root21 -> item212;
  item212 -> item423;
  item423 -> item124;
}`
