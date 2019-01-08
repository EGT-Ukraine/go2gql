package tests

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func MakeRequest(schema graphql.Schema, opts *handler.RequestOptions, ctx context.Context) *graphql.Result {
	return graphql.Do(graphql.Params{
		Schema:         schema,
		RequestString:  opts.Query,
		VariableValues: opts.Variables,
		OperationName:  opts.OperationName,
		Context:        ctx,
	})
}

// From https://github.com/graphql-go/graphql/blob/v0.7.7/executor_test.go#L2016
func AssertJSON(t *testing.T, expected string, actual interface{}) {
	var e interface{}
	if err := json.Unmarshal([]byte(expected), &e); err != nil {
		t.Fatalf(err.Error())
	}
	aJSON, err := json.MarshalIndent(actual, "", "  ")
	if err != nil {
		t.Fatalf(err.Error())
	}
	var a interface{}
	if err := json.Unmarshal(aJSON, &a); err != nil {
		t.Fatalf(err.Error())
	}
	if !reflect.DeepEqual(e, a) {
		eNormalizedJSON, err := json.MarshalIndent(e, "", "  ")
		if err != nil {
			t.Fatalf(err.Error())
		}
		t.Fatalf("Expected JSON:\n\n%v\n\nActual JSON:\n\n%v", string(eNormalizedJSON), string(aJSON))
	}
}
