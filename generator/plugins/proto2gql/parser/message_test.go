package parser

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMessageMethods(t *testing.T) {
	Convey("Test Message.HaveFields", t, func() {
		Convey("Should return true, if there's a normal field", func() {
			So(Message{Fields: []*Field{{}}}.HaveFields(), ShouldBeTrue)
		})
		Convey("Should return true, if there's a map field", func() {
			So(Message{MapFields: []*MapField{{}}}.HaveFields(), ShouldBeTrue)
		})
		Convey("Should return true, if there's a oneof", func() {
			So(Message{OneOffs: []*OneOf{{Fields: []*Field{{}}}}}.HaveFields(), ShouldBeTrue)
		})
		Convey("Should return false, if there no fields", func() {
			So(Message{}.HaveFields(), ShouldBeFalse)
		})
	})
	Convey("Test Message.HaveFieldsExcept", t, func() {
		Convey("Should return true, if there's a normal field", func() {
			msg := Message{Fields: []*Field{
				{Name: "a"},
				{Name: "b"},
			}}
			So(msg.HaveFieldsExcept("a"), ShouldBeTrue)
		})
		Convey("Should return true, if there's a map field", func() {
			msg := Message{MapFields: []*MapField{
				{Name: "a"},
				{Name: "b"},
			}}
			So(msg.HaveFieldsExcept("b"), ShouldBeTrue)
		})
		Convey("Should return true, if there's a oneof", func() {
			msg := Message{OneOffs: []*OneOf{
				{
					Fields: []*Field{
						{Name: "a"},
						{Name: "b"},
						{Name: "c"},
					},
				},
			}}
			So(msg.HaveFieldsExcept("b"), ShouldBeTrue)
		})
		Convey("Should return false, if there's only excepted normal field", func() {
			msg := Message{Fields: []*Field{
				{Name: "a"},
			}}
			So(msg.HaveFieldsExcept("a"), ShouldBeFalse)
		})
		Convey("Should return false, if there's only excepted map field", func() {
			msg := Message{MapFields: []*MapField{
				{Name: "b"},
			}}
			So(msg.HaveFieldsExcept("b"), ShouldBeFalse)
		})

		Convey("Should return false, if there's only excepted oneof field", func() {
			msg := Message{OneOffs: []*OneOf{
				{
					Fields: []*Field{
						{Name: "b"},
					},
				},
			}}
			So(msg.HaveFieldsExcept("b"), ShouldBeFalse)
		})
		Convey("Should return false, if there's no filelds", func() {
			msg := Message{}
			So(msg.HaveFieldsExcept("b"), ShouldBeFalse)
		})
	})
}
