package doable

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInOut(t *testing.T) {
	item1 := StringItem{"id1"}
	item2 := StringItem{"id2"}

	Convey("Given a list", t, func() {
		Convey("It should handle I/O operations", func() {
			l := NewList()

			So(l.Size(), ShouldBeZeroValue)
			So(l.Count(item1), ShouldBeZeroValue)

			l.AddN(item1, 3)
			l.Add(item1)
			l.Add(item2)

			So(l.Size(), ShouldEqual, 2)
			So(l.Count(item1), ShouldEqual, 4)
			So(l.Count(item2), ShouldEqual, 1)

			So(l.Clone(), ShouldResemble, l)

			l.DelN(item1, 2)
			l.Del(item1)
			l.Del(item2)

			So(l.Size(), ShouldEqual, 1)
			So(l.Count(item1), ShouldEqual, 1)
		})

		Convey("It should handle error cases", func() {
			l := NewList()
			l.Add(item1)

			So(l.Del(item1), ShouldBeNil)
			So(l.Del(item1), ShouldNotBeNil)
		})
	})
}
