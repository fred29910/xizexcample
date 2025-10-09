package enhanced_router

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TracingMiddleware creates a middleware for OpenTelemetry tracing.
func TracingMiddleware(tracerName string) MiddlewareFunc {
	return func(next HandleFunc) HandleFunc {
		tracer := otel.Tracer(tracerName)
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			path, _ := ctx.Value("path").(string)
			ctx, span := tracer.Start(ctx, path)
			defer span.End()

			resp, err = next(ctx, req)

			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(codes.Ok, "OK")
			}

			return resp, err
		}
	}
}

// AddSpanAttribute adds an attribute to the current span.
func AddSpanAttribute(ctx context.Context, key string, value string) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String(key, value))
}
