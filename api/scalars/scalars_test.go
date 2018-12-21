package scalars

import (
	"testing"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGraphQLInt32Scalar(t *testing.T) {
	Convey("Test GraphQLInt32Scalar.Serialize", t, func() {
		So(GraphQLInt32Scalar.Serialize(int32(534)), ShouldEqual, int32(534))
		So(GraphQLInt32Scalar.Serialize(int(534)), ShouldEqual, int32(534))
		So(GraphQLInt32Scalar.Serialize("123"), ShouldEqual, nil)
	})
	Convey("Test GraphQLInt32Scalar.ParseValue", t, func() {
		So(GraphQLInt32Scalar.ParseValue(int32(534)), ShouldEqual, int32(534))
		So(GraphQLInt32Scalar.ParseValue(int64(534)), ShouldEqual, int32(534))
		So(GraphQLInt32Scalar.ParseValue(float32(534)), ShouldEqual, int32(534))
		So(GraphQLInt32Scalar.ParseValue(float64(534)), ShouldEqual, int32(534))
		So(GraphQLInt32Scalar.ParseValue("123"), ShouldEqual, int32(123))
		So(GraphQLInt32Scalar.ParseValue("s123"), ShouldEqual, nil)
		So(GraphQLInt32Scalar.ParseValue(int(534)), ShouldEqual, nil)
	})
	Convey("Test GraphQLInt32Scalar.ParseLiteral", t, func() {
		So(GraphQLInt32Scalar.ParseLiteral(&ast.IntValue{Kind: kinds.IntValue, Value: "534"}), ShouldEqual, int32(534))
		So(GraphQLInt32Scalar.ParseLiteral(&ast.StringValue{Kind: kinds.StringValue, Value: "534"}), ShouldEqual, int32(534))
		So(GraphQLInt32Scalar.ParseLiteral(&ast.StringValue{Kind: kinds.StringValue, Value: "s534"}), ShouldEqual, nil)
		So(GraphQLInt32Scalar.ParseLiteral(&ast.StringValue{Kind: kinds.BooleanValue, Value: "true"}), ShouldEqual, nil)
	})
}

func TestGraphQLInt64Scalar(t *testing.T) {
	Convey("Test GraphQLInt64Scalar.Serialize", t, func() {
		So(GraphQLInt64Scalar.Serialize(int64(534)), ShouldEqual, int64(534))
		So(GraphQLInt64Scalar.Serialize(int32(534)), ShouldEqual, int64(534))
		So(GraphQLInt64Scalar.Serialize(int(534)), ShouldEqual, int64(534))
		So(GraphQLInt64Scalar.Serialize("123"), ShouldEqual, nil)
	})
	Convey("Test GraphQLInt64Scalar.ParseValue", t, func() {
		So(GraphQLInt64Scalar.ParseValue(int(534)), ShouldEqual, nil)
		So(GraphQLInt64Scalar.ParseValue(int32(534)), ShouldEqual, int64(534))
		So(GraphQLInt64Scalar.ParseValue(int64(534)), ShouldEqual, int64(534))
		So(GraphQLInt64Scalar.ParseValue(float32(534)), ShouldEqual, int64(534))
		So(GraphQLInt64Scalar.ParseValue(float64(534)), ShouldEqual, int64(534))
		So(GraphQLInt64Scalar.ParseValue("123"), ShouldEqual, int64(123))
		So(GraphQLInt64Scalar.ParseValue("s123"), ShouldEqual, nil)
	})
	Convey("Test GraphQLInt64Scalar.ParseLiteral", t, func() {
		So(GraphQLInt64Scalar.ParseLiteral(&ast.IntValue{Kind: kinds.IntValue, Value: "534"}), ShouldEqual, int64(534))
		So(GraphQLInt64Scalar.ParseLiteral(&ast.StringValue{Kind: kinds.StringValue, Value: "534"}), ShouldEqual, int64(534))
		So(GraphQLInt64Scalar.ParseLiteral(&ast.StringValue{Kind: kinds.StringValue, Value: "s534"}), ShouldEqual, nil)
		So(GraphQLInt64Scalar.ParseLiteral(&ast.StringValue{Kind: kinds.BooleanValue, Value: "true"}), ShouldEqual, nil)
	})
}

func TestGraphQLUInt64Scalar(t *testing.T) {
	Convey("Test GraphQLUInt64Scalar.Serialize", t, func() {
		So(GraphQLUInt64Scalar.Serialize(uint64(534)), ShouldEqual, uint64(534))
		So(GraphQLUInt64Scalar.Serialize(uint32(534)), ShouldEqual, uint64(534))
		So(GraphQLUInt64Scalar.Serialize(uint(534)), ShouldEqual, uint64(534))
		So(GraphQLUInt64Scalar.Serialize("123"), ShouldEqual, nil)
	})
	Convey("Test GraphQLUInt64Scalar.ParseValue", t, func() {
		So(GraphQLUInt64Scalar.ParseValue(uint32(534)), ShouldEqual, uint64(534))
		So(GraphQLUInt64Scalar.ParseValue(uint64(534)), ShouldEqual, uint64(534))
		So(GraphQLUInt64Scalar.ParseValue(float32(534)), ShouldEqual, uint64(534))
		So(GraphQLUInt64Scalar.ParseValue(float64(534)), ShouldEqual, uint64(534))
		So(GraphQLUInt64Scalar.ParseValue("123"), ShouldEqual, uint64(123))

		So(GraphQLUInt64Scalar.ParseValue(int(534)), ShouldEqual, nil)
		So(GraphQLUInt64Scalar.ParseValue("s123"), ShouldEqual, nil)
	})
	Convey("Test GraphQLUInt64Scalar.ParseLiteral", t, func() {
		So(GraphQLUInt64Scalar.ParseLiteral(&ast.IntValue{Kind: kinds.IntValue, Value: "534"}), ShouldEqual, uint64(534))
		So(GraphQLUInt64Scalar.ParseLiteral(&ast.StringValue{Kind: kinds.StringValue, Value: "534"}), ShouldEqual, uint64(534))
		So(GraphQLUInt64Scalar.ParseLiteral(&ast.StringValue{Kind: kinds.StringValue, Value: "s534"}), ShouldEqual, nil)
		So(GraphQLUInt64Scalar.ParseLiteral(&ast.StringValue{Kind: kinds.BooleanValue, Value: "true"}), ShouldEqual, nil)
	})
}

func TestGraphQLUInt32Scalar(t *testing.T) {
	Convey("Test GraphQLUInt32Scalar.Serialize", t, func() {
		So(GraphQLUInt32Scalar.Serialize(uint32(534)), ShouldEqual, uint32(534))
		So(GraphQLUInt32Scalar.Serialize(uint(534)), ShouldEqual, uint32(534))
		So(GraphQLUInt32Scalar.Serialize("123"), ShouldEqual, nil)
		So(GraphQLUInt32Scalar.Serialize(uint64(534)), ShouldEqual, nil)
	})
	Convey("Test GraphQLUInt32Scalar.ParseValue", t, func() {
		So(GraphQLUInt32Scalar.ParseValue(uint32(534)), ShouldEqual, uint32(534))
		So(GraphQLUInt32Scalar.ParseValue(float32(534)), ShouldEqual, uint32(534))
		So(GraphQLUInt32Scalar.ParseValue(float64(534)), ShouldEqual, uint32(534))
		So(GraphQLUInt32Scalar.ParseValue("123"), ShouldEqual, uint32(123))
		So(GraphQLUInt32Scalar.ParseValue(int(534)), ShouldEqual, nil)
		So(GraphQLUInt32Scalar.ParseValue("s123"), ShouldEqual, nil)
	})
	Convey("Test GraphQLUInt32Scalar.ParseLiteral", t, func() {
		So(GraphQLUInt32Scalar.ParseLiteral(&ast.IntValue{Kind: kinds.IntValue, Value: "534"}), ShouldEqual, uint32(534))
		So(GraphQLUInt32Scalar.ParseLiteral(&ast.StringValue{Kind: kinds.StringValue, Value: "534"}), ShouldEqual, uint32(534))
		So(GraphQLUInt32Scalar.ParseLiteral(&ast.StringValue{Kind: kinds.StringValue, Value: "s534"}), ShouldEqual, nil)
		So(GraphQLUInt32Scalar.ParseLiteral(&ast.StringValue{Kind: kinds.BooleanValue, Value: "true"}), ShouldEqual, nil)
	})
}

func TestGraphQLFloat32Scalar(t *testing.T) {
	Convey("Test GraphQLFloat32Scalar.Serialize", t, func() {
		So(GraphQLFloat32Scalar.Serialize(float32(534)), ShouldEqual, float32(534))
		So(GraphQLFloat32Scalar.Serialize(uint64(534)), ShouldEqual, nil)
	})
	Convey("Test GraphQLFloat32Scalar.ParseValue", t, func() {
		So(GraphQLFloat32Scalar.ParseValue(float32(534)), ShouldEqual, float32(534))
		So(GraphQLFloat32Scalar.ParseValue("123"), ShouldEqual, float32(123))
		So(GraphQLFloat32Scalar.ParseValue(int(534)), ShouldEqual, nil)
		So(GraphQLFloat32Scalar.ParseValue("s123"), ShouldEqual, nil)
	})
	Convey("Test GraphQLFloat32Scalar.ParseLiteral", t, func() {
		So(GraphQLFloat32Scalar.ParseLiteral(&ast.IntValue{Kind: kinds.IntValue, Value: "534"}), ShouldEqual, float32(534))
		So(GraphQLFloat32Scalar.ParseLiteral(&ast.StringValue{Kind: kinds.StringValue, Value: "534"}), ShouldEqual, float32(534))
		So(GraphQLFloat32Scalar.ParseLiteral(&ast.StringValue{Kind: kinds.StringValue, Value: "s534"}), ShouldEqual, nil)
		So(GraphQLFloat32Scalar.ParseLiteral(&ast.StringValue{Kind: kinds.BooleanValue, Value: "true"}), ShouldEqual, nil)
	})
}

func TestGraphQLFloat64Scalar(t *testing.T) {
	Convey("Test GraphQLFloat64Scalar.Serialize", t, func() {
		So(GraphQLFloat64Scalar.Serialize(float32(534)), ShouldEqual, float64(534))
		So(GraphQLFloat64Scalar.Serialize(float64(534)), ShouldEqual, float64(534))
		So(GraphQLFloat64Scalar.Serialize(uint64(534)), ShouldEqual, nil)
	})
	Convey("Test GraphQLFloat64Scalar.ParseValue", t, func() {
		So(GraphQLFloat64Scalar.ParseValue(float32(534)), ShouldEqual, float64(534))
		So(GraphQLFloat64Scalar.ParseValue(float64(534)), ShouldEqual, float64(534))
		So(GraphQLFloat64Scalar.ParseValue("123"), ShouldEqual, float64(123))
		So(GraphQLFloat64Scalar.ParseValue(int(534)), ShouldEqual, nil)
		So(GraphQLFloat64Scalar.ParseValue("s123"), ShouldEqual, nil)
	})
	Convey("Test GraphQLFloat64Scalar.ParseLiteral", t, func() {
		So(GraphQLFloat64Scalar.ParseLiteral(&ast.IntValue{Kind: kinds.IntValue, Value: "534"}), ShouldEqual, float64(534))
		So(GraphQLFloat64Scalar.ParseLiteral(&ast.StringValue{Kind: kinds.StringValue, Value: "534"}), ShouldEqual, float64(534))
		So(GraphQLFloat64Scalar.ParseLiteral(&ast.StringValue{Kind: kinds.StringValue, Value: "s534"}), ShouldEqual, nil)
		So(GraphQLFloat64Scalar.ParseLiteral(&ast.StringValue{Kind: kinds.BooleanValue, Value: "true"}), ShouldEqual, nil)
	})
}

func TestNoDataScalar(t *testing.T) {
	Convey("Test NoDataScalar.Serialize", t, func() {
		So(NoDataScalar.Serialize(float32(534)), ShouldEqual, nil)
	})
	Convey("Test NoDataScalar.ParseValue", t, func() {
		So(NoDataScalar.ParseValue(float32(534)), ShouldEqual, 0)
	})
	Convey("Test NoDataScalar.ParseLiteral", t, func() {
		So(NoDataScalar.ParseLiteral(&ast.StringValue{Kind: kinds.StringValue, Value: "534"}), ShouldEqual, 0)
	})
}
