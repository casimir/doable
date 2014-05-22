package doable

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNoDep(t *testing.T) {
	Convey("Given an Item and an Item list", t, func() {
		notUsedItem := StringItem{"not used"}
		list := []Item{
			StringItem{"id1"},
			StringItem{"id2"},
			notUsedItem,
		}
		testList := append([]Item{}, list[1:]...)

		Convey("It should return if it is doable", func() {
			tree := New(NewNode(StringItem{"id1"}), nil)
			So(tree.Doable(), ShouldBeFalse)
			tree.Avail = []Item{}
			So(tree.Doable(), ShouldBeFalse)
			tree.Avail = list
			So(tree.Doable(), ShouldBeTrue)
			So(tree.Avail, ShouldResemble, testList)
			tree.Root = NewNode(MockItem{})
			So(tree.Doable(), ShouldBeFalse)
			tree.Root = NewNode(StringItem{"wrongId"})
			So(tree.Doable(), ShouldBeFalse)
		})

		Convey("When it is doable it should consume the item in the list", func() {
			list = []Item{
				StringItem{"id1"},
				notUsedItem,
			}
			tree := New(NewNode(StringItem{"id1"}), list)
			So(tree.Doable(), ShouldBeTrue)
			So(tree.Doable(), ShouldBeFalse)
			So(tree.Avail, ShouldResemble, []Item{notUsedItem})
		})
	})
}

func TestDeps(t *testing.T) {
	Convey("Given an Item and an Item list", t, func() {
		notUsedItem := StringItem{"not used"}
		notHere := StringItem{"not here"}

		Convey("When dependencies are not satisfied  it should return it is not doable", func() {
			list := []Item{
				StringItem{"id1"},
				notUsedItem,
			}

			root := NewNode(StringItem{"root"})
			root.AddDep(NewNode(list[0]), NewNode(notHere))

			tree := New(root, list)

			expected := NewNode(root.item)
			expected.AddDep(NewNode(notHere))

			So(tree.Doable(), ShouldBeFalse)
			So(tree.Avail, ShouldResemble, []Item{notUsedItem})
			So(tree.Root, ShouldResemble, expected)
		})

		Convey("When dependencies are satisfied it should return it is doable", func() {
			list := []Item{
				StringItem{"id0"},
				StringItem{"id1"},
				StringItem{"id2"},
				StringItem{"id3"},
				StringItem{"id4"},
				notUsedItem,
			}

			node1 := NewNode(StringItem{"node 1"})
			node1.AddDep(NewNode(list[2]), NewNode(list[3]))

			node2 := NewNode(StringItem{"node 2"})
			node2.AddDep(NewNode(list[4]), node1)

			root := NewNode(StringItem{"root"})
			root.AddDep(NewNode(list[0]), NewNode(list[1]), node2)

			// root
			//  ↳ id0, id1, node2
			//               ↳ id4, node1
			//                       ↳ id2, id3
			tree := New(root, list)

			So(tree.Doable(), ShouldBeTrue)
			So(tree.Avail, ShouldResemble, []Item{notUsedItem})
		})
	})
}

type MockItem struct{}

func (i MockItem) Match(other Item) bool {
	return false
}

func NewNode(i Item) *Node {
	return &Node{nil, i}
}
