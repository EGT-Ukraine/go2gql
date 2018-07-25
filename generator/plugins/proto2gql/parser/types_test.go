package parser

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTypeIs(t *testing.T) {
	Convey("Test Type", t, func() {

		Convey("Type is Message", func() {
			msgTyp := &MessageType{Message: &Message{}}
			So(msgTyp.Kind(), ShouldEqual, TypeMessage)
		})

		Convey("Type is Enum", func() {
			enumTyp := &EnumType{Enum: &Enum{}}
			So(enumTyp.Kind(), ShouldEqual, TypeEnum)
		})

		Convey("Type is Map", func() {
			enumTyp := &MapType{Map: &Map{}}
			So(enumTyp.Kind(), ShouldEqual, TypeMap)
		})

		Convey("Type is Scalar", func() {
			enumTyp := &ScalarType{ScalarName: "int"}
			So(enumTyp.Kind(), ShouldEqual, TypeScalar)
		})
	})
}
