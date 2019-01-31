package parser

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func testFileInfo(file *File) *File {
	var Int32Type = &Scalar{file: file, ScalarName: "int32"}

	var StringType = &Scalar{file: file, ScalarName: "string"}

	var RootMessage = file.Messages[0]

	var RootMessage2 = file.Messages[4]

	var RootEnum = file.Enums[0]

	var EmptyMessage = file.Messages[2]

	var NestedMessage = file.Messages[1]

	var NestedEnum = file.Enums[1]

	var NestedNestedEnum = file.Enums[2]

	var MessageWithEmpty = file.Messages[3]

	var CommonCommonEnum = file.Imports[0].Enums[0]

	var CommonCommonMessage = file.Imports[0].Messages[0]

	var ParentScopeEnum = file.Imports[2].Enums[0]

	var Proto2Message = file.Imports[1].Messages[0]

	return &File{
		Services: []*Service{
			{
				Name:          "ServiceExample",
				QuotedComment: `"Service, which do smth"`,
				Methods: []*Method{
					{Name: "getQueryMethod", InputMessage: RootMessage, OutputMessage: RootMessage, QuotedComment: `""`},
					{Name: "mutationMethod", InputMessage: RootMessage2, OutputMessage: NestedMessage, QuotedComment: `"rpc comment"`},
					{Name: "EmptyMsgs", InputMessage: EmptyMessage, OutputMessage: EmptyMessage, QuotedComment: `""`},
					{Name: "MsgsWithEpmty", InputMessage: MessageWithEmpty, OutputMessage: MessageWithEmpty, QuotedComment: `""`},
				},
			},
		},
		Messages: []*Message{
			{
				file:          file,
				Name:          "RootMessage",
				QuotedComment: `""`,
				TypeName:      TypeName{"RootMessage"},
				NormalFields: []*NormalField{
					{Name: "r_msg", Type: NestedMessage, Repeated: true, QuotedComment: `"repeated Message"`},
					{Name: "r_scalar", Type: Int32Type, Repeated: true, QuotedComment: `"repeated Scalar"`},
					{Name: "r_enum", Type: RootEnum, Repeated: true, QuotedComment: `"repeated Enum"`},
					{Name: "r_empty_msg", Type: EmptyMessage, Repeated: true, QuotedComment: `"repeated empty message"`},
					{Name: "n_r_enum", Type: CommonCommonEnum, QuotedComment: `"non-repeated Enum"`},
					{Name: "n_r_scalar", Type: Int32Type, QuotedComment: `"non-repeated Scalar"`},
					{Name: "n_r_msg", Type: CommonCommonMessage, QuotedComment: `"non-repeated Message"`},
					{Name: "scalar_from_context", Type: Int32Type, QuotedComment: `"field from context"`},
					{Name: "n_r_empty_msg", Type: EmptyMessage, QuotedComment: `"non-repeated empty message field"`},
					{Name: "leading_dot", Type: CommonCommonMessage, QuotedComment: `"leading dot in type name"`},
					{Name: "parent_scope", Type: ParentScopeEnum, QuotedComment: `"parent scope"`},
					{Name: "proto2message", Type: Proto2Message, QuotedComment: `""`},
				},
				OneOffs: []*OneOf{
					{Name: "enum_first_oneoff", Fields: []*NormalField{
						{Name: "e_f_o_e", Type: CommonCommonEnum, QuotedComment: `""`},
						{Name: "e_f_o_s", Type: Int32Type, QuotedComment: `""`},
						{Name: "e_f_o_m", Type: CommonCommonMessage, QuotedComment: `""`},
						{Name: "e_f_o_em", Type: EmptyMessage, QuotedComment: `"non-repeated Message"`},
					}},
					{Name: "scalar_first_oneoff", Fields: []*NormalField{
						{Name: "s_f_o_s", Type: Int32Type, QuotedComment: `"non-repeated Scalar"`},
						{Name: "s_f_o_e", Type: RootEnum, QuotedComment: `"non-repeated Enum"`},
						{Name: "s_f_o_mes", Type: RootMessage2, QuotedComment: `"non-repeated Message"`},
						{Name: "s_f_o_m", Type: EmptyMessage, QuotedComment: `"non-repeated Message"`},
					}},
					{Name: "message_first_oneoff", Fields: []*NormalField{
						{Name: "m_f_o_m", Type: RootMessage2, QuotedComment: `"non-repeated Message"`},
						{Name: "m_f_o_s", Type: Int32Type, QuotedComment: `"non-repeated Scalar"`},
						{Name: "m_f_o_e", Type: RootEnum, QuotedComment: `"non-repeated Enum"`},
						{Name: "m_f_o_em", Type: EmptyMessage, QuotedComment: `"non-repeated Message"`},
					}},
					{Name: "empty_first_oneoff", Fields: []*NormalField{
						{Name: "em_f_o_em", Type: EmptyMessage, QuotedComment: `"non-repeated Message"`},
						{Name: "em_f_o_s", Type: Int32Type, QuotedComment: `"non-repeated Scalar"`},
						{Name: "em_f_o_en", Type: RootEnum, QuotedComment: `"non-repeated Enum"`},
						{Name: "em_f_o_m", Type: RootMessage2, QuotedComment: `"non-repeated Message"`},
					}},
				},
				MapFields: []*MapField{
					{
						Name:          "map_enum",
						QuotedComment: `"enum_map\n Map with enum value"`,
						Map: &Map{
							Message:   RootMessage,
							KeyType:   Int32Type,
							ValueType: NestedEnum,
							file:      file,
						},
					},
					{
						Name:          "map_scalar",
						QuotedComment: `"scalar map\n Map with scalar value"`,
						Map: &Map{
							Message:   RootMessage,
							KeyType:   Int32Type,
							ValueType: Int32Type,
							file:      file,
						},
					},
					{
						Name:          "map_msg",
						QuotedComment: `"Map with Message value"`,
						Map: &Map{
							Message:   RootMessage,
							KeyType:   StringType,
							ValueType: NestedMessage,
							file:      file,
						},
					},
					{
						Name:          "ctx_map",
						QuotedComment: `""`,
						Map: &Map{
							Message:   RootMessage,
							KeyType:   StringType,
							ValueType: NestedMessage,
							file:      file,
						},
					},
					{
						Name:          "ctx_map_enum",
						QuotedComment: `""`,
						Map: &Map{
							Message:   RootMessage,
							KeyType:   StringType,
							ValueType: NestedEnum,
							file:      file,
						},
					},
				},
			},
			{
				file:          file,
				Name:          "NestedMessage",
				QuotedComment: `""`,
				TypeName:      TypeName{"RootMessage", "NestedMessage"},
				NormalFields: []*NormalField{
					{Name: "sub_r_enum", Type: NestedEnum, Repeated: true, QuotedComment: `"repeated Enum"`},
					{Name: "sub_sub_r_enum", Type: NestedNestedEnum, Repeated: true, QuotedComment: `"repeated Enum"`},
				},
			},
			{
				file:          file,
				Name:          "Empty",
				QuotedComment: `""`,
				TypeName:      TypeName{"Empty"},
			},
			{
				file:          file,
				Name:          "MessageWithEmpty",
				QuotedComment: `""`,
				TypeName:      TypeName{"MessageWithEmpty"},
				NormalFields: []*NormalField{
					{Name: "empt", Type: EmptyMessage, QuotedComment: `""`},
				},
			},
			{
				file:          file,
				Name:          "RootMessage2",
				QuotedComment: `""`,
				TypeName:      TypeName{"RootMessage2"},
				NormalFields: []*NormalField{
					{Name: "some_field", Type: Int32Type, QuotedComment: `""`},
				},
			},
		},
		Enums: []*Enum{
			{
				file:          file,
				Name:          "RootEnum",
				QuotedComment: `""`,
				TypeName:      TypeName{"RootEnum"},
				Values: []*EnumValue{
					{Name: "RootEnumVal0", Value: 0, QuotedComment: `""`},
					{Name: "RootEnumVal1", Value: 1, QuotedComment: `""`},
					{Name: "RootEnumVal2", Value: 2, QuotedComment: `"It's a RootEnumVal2"`},
				},
			},
			{
				file:          file,
				Name:          "NestedEnum",
				QuotedComment: `""`,
				TypeName:      TypeName{"RootMessage", "NestedEnum"},
				Values: []*EnumValue{
					{Name: "NestedEnumVal0", Value: 0, QuotedComment: `""`},
					{Name: "NestedEnumVal1", Value: 1, QuotedComment: `""`},
				},
			},
			{
				file:          file,
				Name:          "NestedNestedEnum",
				QuotedComment: `""`,
				TypeName:      TypeName{"RootMessage", "NestedMessage", "NestedNestedEnum"},
				Values: []*EnumValue{
					{Name: "NestedNestedEnumVal0", Value: 0, QuotedComment: `""`},
					{Name: "NestedNestedEnumVal1", Value: 1, QuotedComment: `""`},
					{Name: "NestedNestedEnumVal2", Value: 2, QuotedComment: `""`},
					{Name: "NestedNestedEnumVal3", Value: 3, QuotedComment: `""`},
				},
			},
		},
	}
}

func TestParser_Parse(t *testing.T) {
	Convey("Test Parser.Parse", t, func() {
		parser := Parser{}
		test, err := parser.Parse("../../../../testdata/test.proto", nil, []string{"../../../../testdata"})
		So(err, ShouldBeNil)
		So(test, ShouldNotBeNil)
		test2, err := parser.Parse("../../../../testdata/test2.proto", nil, []string{"../../../../testdata"})
		So(err, ShouldBeNil)
		So(test2, ShouldNotBeNil)
		So(test, ShouldNotEqual, test2)

		Convey("Imports should be the same", func() {
			So(len(test.Imports), ShouldEqual, 3)
			So(len(test2.Imports), ShouldEqual, 1)
			So(test.Imports[0], ShouldEqual, test2.Imports[0])
		})
		Convey("If we trying to parse same File, it should return pointer to parsed one", func() {
			test22, err := parser.Parse("../../../../testdata/test2.proto", nil, []string{"../../../../testdata"})
			So(err, ShouldBeNil)
			So(test22, ShouldEqual, test2)
		})
		f := testFileInfo(test)

		Convey("test.proto Should contains valid enums", func() {
			So(test.Enums, ShouldHaveLength, len(f.Enums))
			for i, enum := range test.Enums {
				validEnum := f.Enums[i]
				Convey("Should contain "+validEnum.Name, func() {
					So(enum.File, ShouldEqual, validEnum.File)
					So(enum.Name, ShouldEqual, validEnum.Name)
					So(enum, ShouldEqual, enum)
					So(enum.File(), ShouldEqual, test)
					So(enum.TypeName, ShouldResemble, validEnum.TypeName)
					So(enum.QuotedComment, ShouldEqual, validEnum.QuotedComment)
					Convey(validEnum.Name+" enum should contains valid values", func() {
						So(enum.Values, ShouldHaveLength, len(validEnum.Values))
						for i, value := range enum.Values {
							validValue := validEnum.Values[i]
							Convey(validEnum.Name+" enum should contains valid "+validValue.Name+" value", func() {
								So(value.Name, ShouldEqual, validValue.Name)
								So(value.Value, ShouldEqual, validValue.Value)
								So(value.QuotedComment, ShouldEqual, validValue.QuotedComment)
							})
						}
					})
				})
			}
		})

		Convey("test.proto Should contains valid messages", func() {
			So(test.Messages, ShouldHaveLength, len(f.Messages))
			for i, msg := range test.Messages {
				validMsg := f.Messages[i]
				Convey("Should have valid parsed "+strings.Join(validMsg.TypeName, "_")+" message ", func() {
					So(msg.File, ShouldEqual, validMsg.File)
					So(msg.Name, ShouldEqual, validMsg.Name)
					So(msg, ShouldEqual, msg)
					So(msg.File(), ShouldEqual, test)
					So(msg.TypeName, ShouldResemble, validMsg.TypeName)
					So(msg.QuotedComment, ShouldEqual, validMsg.QuotedComment)
					So(msg.NormalFields, ShouldHaveLength, len(validMsg.NormalFields))
					for i, fld := range msg.NormalFields {
						validFld := validMsg.NormalFields[i]
						Convey("Should have valid parsed "+strings.Join(validMsg.TypeName, "_")+"."+validFld.Name+" field", func() {
							So(fld.Name, ShouldEqual, validFld.Name)
							So(fld.Repeated, ShouldEqual, validFld.Repeated)
							So(fld.QuotedComment, ShouldEqual, validFld.QuotedComment)
							CompareTypes(fld.Type, validFld.Type)
						})
					}
					So(msg.MapFields, ShouldHaveLength, len(validMsg.MapFields))
					for i, fld := range msg.MapFields {
						validFld := validMsg.MapFields[i]
						Convey("Should have valid parsed "+strings.Join(validMsg.TypeName, "_")+"."+validFld.Name+" field", func() {
							So(fld.Name, ShouldEqual, validFld.Name)
							So(fld.QuotedComment, ShouldEqual, validFld.QuotedComment)
							CompareTypes(fld.Map, validFld.Map)
						})
					}
					So(msg.OneOffs, ShouldHaveLength, len(validMsg.OneOffs))
					for i, oneOf := range msg.OneOffs {
						validOneOf := validMsg.OneOffs[i]
						Convey("Should have valid parsed "+strings.Join(validMsg.TypeName, "_")+"."+validOneOf.Name+" one of", func() {
							So(oneOf.Name, ShouldEqual, validOneOf.Name)
							So(oneOf.Fields, ShouldHaveLength, len(validOneOf.Fields))
							for i, fld := range oneOf.Fields {
								validFld := validOneOf.Fields[i]
								Convey("Should have valid parsed "+strings.Join(validMsg.TypeName, "_")+"."+validOneOf.Name+"."+validFld.Name+" one of field", func() {
									So(fld.Name, ShouldEqual, validFld.Name)
									So(fld.QuotedComment, ShouldEqual, validFld.QuotedComment)
									CompareTypes(fld.Type, validFld.Type)
								})
							}

						})
					}
				})

			}
		})
		Convey("test.proto Should contain valid services", func() {
			So(test.Services, ShouldHaveLength, len(f.Services))
			for i, srv := range test.Services {
				validSrv := f.Services[i]
				Convey("Should have valid parsed "+validSrv.Name+" service ", func() {
					So(srv.Name, ShouldEqual, validSrv.Name)
					So(srv.QuotedComment, ShouldEqual, validSrv.QuotedComment)
					Convey(validSrv.Name+" should contains valid methods", func() {
						So(srv.Methods, ShouldHaveLength, len(validSrv.Methods))
						for i, method := range srv.Methods {
							validMethod := validSrv.Methods[i]
							Convey(validSrv.Name+" should contains valid "+validMethod.Name+" method", func() {
								So(method.Name, ShouldEqual, validMethod.Name)
								So(method.QuotedComment, ShouldEqual, validMethod.QuotedComment)
								Convey(validSrv.Name+"."+validMethod.Name+" should have valid input message type", func() {
									CompareTypes(method.InputMessage, validMethod.InputMessage)
								})
								Convey(validSrv.Name+"."+validMethod.Name+" should have valid output message type", func() {
									CompareTypes(method.OutputMessage, validMethod.OutputMessage)
								})
							})
						}
					})
				})
			}
		})
	})
}

func CompareTypes(t1, t2 Type) {
	So(t1, ShouldNotBeNil)
	So(t2, ShouldNotBeNil)

	switch protoType := t1.(type) {
	case *Scalar:
		So(protoType.ScalarName, ShouldEqual, t2.(*Scalar).ScalarName)
	case *Message:
		So(t1, ShouldEqual, t2)
		So(t1.File(), ShouldEqual, t2.File())
	case *Enum:
		So(t1, ShouldEqual, t2)
		So(t1.File(), ShouldEqual, t2.File())
	case *Map:
		So(t1.(*Map).Message, ShouldEqual, t2.(*Map).Message)
		CompareTypes(t1.(*Map).KeyType, t2.(*Map).KeyType)
		CompareTypes(t1.(*Map).ValueType, t2.(*Map).ValueType)
		So(t1.File(), ShouldEqual, t2.File())
	default:
		panic("Undefined type")
	}
}
