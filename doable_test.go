package doable

import (
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
	Convey("Given an Item and an Item list", t, func() {
		item0 := StringItem{"id0"}
		item1 := StringItem{"id1"}
		item2 := StringItem{"id2"}
		item3 := StringItem{"id3"}
		item4 := StringItem{"id4"}
		itemNotUsed := StringItem{"not used"}
		itemNotHere := StringItem{"not here"}
		itemRoot := StringItem{"root"}

		Convey("When dependencies are not satisfied it should return it is not doable", func() {
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

type MockItem struct{}

func (i MockItem) ID() string {
	return ""
}

func (i MockItem) Match(other Item) bool {
	return false
}
