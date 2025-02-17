package tracer

import (
	"context"
	"pkg/helper"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct {
	Ctx  context.Context
	Span trace.Span
}

func NewTracer(ctx context.Context, name string) (Tracer, context.Context) {
	tracer := otel.Tracer(name)
	_, span := tracer.Start(ctx, name)
	t := Tracer{
		Span: span,
	}
	return t, ctx
}

func (t *Tracer) End() {
	t.Span.End()
}

func (t *Tracer) AddAttributes(attributes ...attribute.KeyValue) {
	if attributes == nil {
		return
	}
	t.Span.SetAttributes(attributes...)
}

func (t *Tracer) RecordError(err error, msg string, attributes ...attribute.KeyValue) {

	caller := helper.Caller(2)
	stack := helper.StackTrace(3, 32)

	t.Span.RecordError(err,
		trace.WithTimestamp(time.Now()),
		trace.WithAttributes(
			attribute.String("stack.trace", stack),
		),
	)
	t.Span.SetStatus(1, msg)
	t.AddAttributes(attributes...)
	t.AddAttributes(attribute.String("caller", caller))
}
