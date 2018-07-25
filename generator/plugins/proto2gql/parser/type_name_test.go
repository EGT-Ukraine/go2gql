package parser

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTypeName(t *testing.T) {
	Convey("Test TypeName", t, func() {
		Convey("Should check TypeName.Equal", func() {
			So(TypeName{}.Equal(TypeName{}), ShouldBeTrue)
			So(TypeName{}.Equal(TypeName{"1"}), ShouldBeFalse)
			So(TypeName{"1"}.Equal(TypeName{"1"}), ShouldBeTrue)
			So(TypeName{"1"}.Equal(TypeName{"2"}), ShouldBeFalse)
			So(TypeName{"2", "33"}.Equal(TypeName{"2", "33"}), ShouldBeTrue)
			So(TypeName{"2", "33"}.Equal(TypeName{"2", "44"}), ShouldBeFalse)
		})
		Convey("Should check TypeName.NewSubTypeName", func() {
			So(TypeName{}.NewSubTypeName("1"), ShouldResemble, TypeName{"1"})
			So(TypeName{"1"}.NewSubTypeName("2"), ShouldResemble, TypeName{"1", "2"})
		})
	})
}
