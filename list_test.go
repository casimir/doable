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
			c := NewList()

			So(c.size(), ShouldBeZeroValue)
			So(c.count(item1), ShouldBeZeroValue)

			c.add(item1, 3)
			c.add1(item2)

			So(c.size(), ShouldEqual, 2)
			So(c.count(item1), ShouldEqual, 3)
			So(c.count(item2), ShouldEqual, 1)

			c.del(item1, 2)
			c.del1(item2)

			So(c.size(), ShouldEqual, 1)
			So(c.count(item1), ShouldEqual, 1)
		})

		Convey("It should handle error cases", func() {
			c := NewList()
			c.add1(item1)

			So(c.del1(item1), ShouldBeNil)
			So(c.del1(item1), ShouldNotBeNil)
		})
	})
}
