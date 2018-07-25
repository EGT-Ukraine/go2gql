package tracer

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

type Tracer interface {
	CreateChildSpanFromContext(c context.Context, name string) opentracing.Span
	ContextWithSpan(c context.Context, span opentracing.Span) context.Context
}
