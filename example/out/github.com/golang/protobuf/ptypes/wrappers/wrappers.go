// This file was generated by github.com/EGT-Ukraine/go2gql. DO NOT EDIT IT
package wrappers

import (
	context "context"

	scalars "github.com/EGT-Ukraine/go2gql/api/scalars"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	graphql "github.com/graphql-go/graphql"
)

// Enums
// Input object
var DoubleValueInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:   "DoubleValueInput",
	Fields: graphql.InputObjectConfigFieldMap{},
})

func init() {
	DoubleValueInput.AddFieldConfig("value", &graphql.InputObjectFieldConfig{Type: scalars.GraphQLFloat64Scalar, Description: "The double value."})
}

var FloatValueInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:   "FloatValueInput",
	Fields: graphql.InputObjectConfigFieldMap{},
})

func init() {
	FloatValueInput.AddFieldConfig("value", &graphql.InputObjectFieldConfig{Type: scalars.GraphQLFloat32Scalar, Description: "The float value."})
}

var Int64ValueInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:   "Int64ValueInput",
	Fields: graphql.InputObjectConfigFieldMap{},
})

func init() {
	Int64ValueInput.AddFieldConfig("value", &graphql.InputObjectFieldConfig{Type: scalars.GraphQLInt64Scalar, Description: "The int64 value."})
}

var UInt64ValueInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:   "UInt64ValueInput",
	Fields: graphql.InputObjectConfigFieldMap{},
})

func init() {
	UInt64ValueInput.AddFieldConfig("value", &graphql.InputObjectFieldConfig{Type: scalars.GraphQLUInt64Scalar, Description: "The uint64 value."})
}

var Int32ValueInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:   "Int32ValueInput",
	Fields: graphql.InputObjectConfigFieldMap{},
})

func init() {
	Int32ValueInput.AddFieldConfig("value", &graphql.InputObjectFieldConfig{Type: scalars.GraphQLInt32Scalar, Description: "The int32 value."})
}

var UInt32ValueInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:   "UInt32ValueInput",
	Fields: graphql.InputObjectConfigFieldMap{},
})

func init() {
	UInt32ValueInput.AddFieldConfig("value", &graphql.InputObjectFieldConfig{Type: scalars.GraphQLUInt32Scalar, Description: "The uint32 value."})
}

var BoolValueInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:   "BoolValueInput",
	Fields: graphql.InputObjectConfigFieldMap{},
})

func init() {
	BoolValueInput.AddFieldConfig("value", &graphql.InputObjectFieldConfig{Type: graphql.Boolean, Description: "The bool value."})
}

var StringValueInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:   "StringValueInput",
	Fields: graphql.InputObjectConfigFieldMap{},
})

func init() {
	StringValueInput.AddFieldConfig("value", &graphql.InputObjectFieldConfig{Type: graphql.String, Description: "The string value."})
}

var BytesValueInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:   "BytesValueInput",
	Fields: graphql.InputObjectConfigFieldMap{},
})

func init() {
	BytesValueInput.AddFieldConfig("value", &graphql.InputObjectFieldConfig{Type: scalars.GraphQLBytesScalar, Description: "The bytes value."})
}

// Input objects resolvers
func ResolveDoubleValueInput(ctx context.Context, i interface{}) (_ *wrappers.DoubleValue, rerr error) {
	if i == nil {
		return nil, nil
	}
	args, _ := i.(map[string]interface{})
	_ = args
	var result = new(wrappers.DoubleValue)
	result.Value = i.(float64)

	return result, nil
}
func ResolveFloatValueInput(ctx context.Context, i interface{}) (_ *wrappers.FloatValue, rerr error) {
	if i == nil {
		return nil, nil
	}
	args, _ := i.(map[string]interface{})
	_ = args
	var result = new(wrappers.FloatValue)
	result.Value = i.(float32)

	return result, nil
}
func ResolveInt64ValueInput(ctx context.Context, i interface{}) (_ *wrappers.Int64Value, rerr error) {
	if i == nil {
		return nil, nil
	}
	args, _ := i.(map[string]interface{})
	_ = args
	var result = new(wrappers.Int64Value)
	result.Value = i.(int64)

	return result, nil
}
func ResolveUInt64ValueInput(ctx context.Context, i interface{}) (_ *wrappers.UInt64Value, rerr error) {
	if i == nil {
		return nil, nil
	}
	args, _ := i.(map[string]interface{})
	_ = args
	var result = new(wrappers.UInt64Value)
	result.Value = i.(uint64)

	return result, nil
}
func ResolveInt32ValueInput(ctx context.Context, i interface{}) (_ *wrappers.Int32Value, rerr error) {
	if i == nil {
		return nil, nil
	}
	args, _ := i.(map[string]interface{})
	_ = args
	var result = new(wrappers.Int32Value)
	result.Value = i.(int32)

	return result, nil
}
func ResolveUInt32ValueInput(ctx context.Context, i interface{}) (_ *wrappers.UInt32Value, rerr error) {
	if i == nil {
		return nil, nil
	}
	args, _ := i.(map[string]interface{})
	_ = args
	var result = new(wrappers.UInt32Value)
	result.Value = i.(uint32)

	return result, nil
}
func ResolveBoolValueInput(ctx context.Context, i interface{}) (_ *wrappers.BoolValue, rerr error) {
	if i == nil {
		return nil, nil
	}
	args, _ := i.(map[string]interface{})
	_ = args
	var result = new(wrappers.BoolValue)
	result.Value = i.(bool)

	return result, nil
}
func ResolveStringValueInput(ctx context.Context, i interface{}) (_ *wrappers.StringValue, rerr error) {
	if i == nil {
		return nil, nil
	}
	args, _ := i.(map[string]interface{})
	_ = args
	var result = new(wrappers.StringValue)
	result.Value = i.(string)

	return result, nil
}
func ResolveBytesValueInput(ctx context.Context, i interface{}) (_ *wrappers.BytesValue, rerr error) {
	if i == nil {
		return nil, nil
	}
	args, _ := i.(map[string]interface{})
	_ = args
	var result = new(wrappers.BytesValue)
	result.Value = i.([]byte)

	return result, nil
}

// Output objects
// Maps input objects
// Maps input objects resolvers
// Maps output objects
